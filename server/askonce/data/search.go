package data

import (
	"askonce/api/jobd"
	"askonce/api/web_search"
	"askonce/components/dto/dto_search"
	"askonce/es"
	"askonce/helpers"
	"askonce/models"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"gorm.io/datatypes"
	"sort"
	"time"
)

type SearchData struct {
	flow.Data
	jobdApi          *jobd.JobdApi
	webSearchApi     *web_search.WebSearchApi
	askAttachDao     *models.AskAttachDao
	askSubSearchDao  *models.AskSubSearchDao
	kdbDocContentDao *models.KdbDocContentDao
	kdbDocDao        *models.KdbDocDao
	fileDao          *models.FileDao
	kdbData          *KdbData
}

func (entity *SearchData) OnCreate() {
	entity.jobdApi = entity.Create(new(jobd.JobdApi)).(*jobd.JobdApi)
	entity.webSearchApi = entity.Create(new(web_search.WebSearchApi)).(*web_search.WebSearchApi)
	entity.askAttachDao = entity.Create(new(models.AskAttachDao)).(*models.AskAttachDao)
	entity.askSubSearchDao = entity.Create(new(models.AskSubSearchDao)).(*models.AskSubSearchDao)
	entity.kdbDocContentDao = entity.Create(new(models.KdbDocContentDao)).(*models.KdbDocContentDao)
	entity.kdbDocDao = entity.Create(new(models.KdbDocDao)).(*models.KdbDocDao)
	entity.fileDao = entity.Create(new(models.FileDao)).(*models.FileDao)
	entity.kdbData = entity.Create(new(KdbData)).(*KdbData)
}

type DocSearchOptions struct {
	QuerySize  int    // 返回数量
	IndexName  string // 知识库索引名称
	ReturnFull bool   // 返回全文
}

func (s *DocSearchOptions) WithIndex(indexName string) *DocSearchOptions {
	s.IndexName = indexName
	return s
}

func (s *DocSearchOptions) WithQuerySize(querySize int) *DocSearchOptions {
	s.QuerySize = querySize
	return s
}

func (s *DocSearchOptions) WithReturnFull(returnFull bool) *DocSearchOptions {
	s.ReturnFull = returnFull
	return s
}

func (entity *SearchData) DocSearch(sessionId, question string, opts *DocSearchOptions) (results []dto_search.CommonSearchOutput, err error) {
	if opts == nil {
		opts = new(DocSearchOptions)
	}
	if opts.QuerySize == 0 {
		opts.QuerySize = 20
	}
	results = make([]dto_search.CommonSearchOutput, 0)
	var errS error
	if len(opts.IndexName) == 0 { // web搜索
		results, errS = entity.webSearch(question, opts.QuerySize)
		if errS != nil {
			entity.LogErrorf("web搜索报错")
		}
	} else { // 知识库搜索
		// es搜索的片段
		results, errS = entity.esDocSearch(question, opts.IndexName, opts.QuerySize, opts.ReturnFull)
		if errS != nil {
			entity.LogErrorf("es搜索报错")
		}
	}
	if len(results) == 0 {
		return []dto_search.CommonSearchOutput{}, nil
	}
	if sessionId != "" {
		now := time.Now()
		searchResultStr, _ := json.Marshal(results)
		err = entity.askSubSearchDao.Insert(&models.AskSubSearch{
			SessionId:    sessionId,
			SubQuestion:  question,
			SearchResult: searchResultStr,
			CrudModel: orm.CrudModel{
				CreatedAt: now,
				UpdatedAt: now,
			},
		})
	}
	return
}

func (entity *SearchData) esDocSearch(question string, indexName string, querySize int, returnFull bool) (res []dto_search.CommonSearchOutput, err error) {
	embRes, err := helpers.EmbeddingGpt.CreateEmbedding(entity.GetCtx(), []string{question})
	if err != nil {
		return
	}
	recalls, err := es.CommonDocumentSearch[*es.DocDocument](entity.GetCtx(), indexName, question, embRes[0], 20)
	if err != nil {
		return
	}
	if len(recalls) == 0 {
		return
	}
	recallsParsedRes, err := entity.jobdApi.SearchResultPostProcess(question, recalls)
	if err != nil {
		return
	}
	recalls = recallsParsedRes.SearchResult
	sort.Slice(recalls, func(i, j int) bool {
		return recalls[i].Score > recalls[j].Score
	})
	dataSearchMap := make(map[int]*es.DocDocument)
	var dataIds []int64
	for i, result := range recalls {
		dataIds = append(dataIds, result.DocId)
		dataSearchMap[i] = result
	}
	dataContents, err := entity.kdbDocContentDao.GetByDataIds(dataIds)
	if err != nil {
		return nil, err
	}
	dataContentMap := make(map[int64]string)
	for _, content := range dataContents {
		dataContentMap[content.DocId] = content.Content
	}
	docs, err := entity.kdbDocDao.GetByIds(dataIds)
	if err != nil {
		return nil, err
	}
	docMap := make(map[int64]*models.KdbDoc)
	fileIds := make([]string, 0)
	for _, da := range docs {
		docMap[da.Id] = da
		fileIds = append(fileIds, da.SourceId)
	}
	files, err := entity.fileDao.GetByIds(fileIds)
	if err != nil {
		return nil, err
	}
	filePathMap := make(map[string]string)
	for _, file := range files {
		filePathMap[file.Id] = file.Path
	}
	for i, result := range recalls {
		ddd, ok := docMap[result.DocId]
		if !ok {
			continue
		}
		out := dto_search.CommonSearchOutput{}
		out.DocSegmentId = result.DocSegmentId
		out.DocId = result.DocId
		out.Title = ddd.DocName
		out.Url = filePathMap[ddd.SourceId]
		out.Metadata = ddd.Metadata
		out.Content = appendText(dataSearchMap[i], dataContentMap[result.DocId])
		out.Score = result.Score
		out.Form = "kdb"
		if returnFull {
			out.FullContent = dataContentMap[result.DocId]
		}
		res = append(res, out)
	}
	if len(res) >= querySize {
		res = res[:querySize]
	}
	return
}

func (entity *SearchData) CreateSession(userId string) (add *models.AskInfo, err error) {
	askInfoDao := entity.Create(new(models.AskInfoDao)).(*models.AskInfoDao)
	sessionId := fmt.Sprintf("%x", md5.Sum([]byte(helpers.GenIDStr())))
	add = &models.AskInfo{
		SessionId: sessionId,
		Question:  "",
		UserId:    userId,
		CrudModel: orm.CrudModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	err = askInfoDao.Insert(add)
	err = entity.askAttachDao.Insert(&models.AskAttach{
		SessionId: sessionId,
		Reference: datatypes.JSON("[]"),
		Outline:   datatypes.JSON("[]"),
		Relation:  datatypes.JSON("{}"),
		CrudModel: orm.CrudModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	})
	return
}

func (entity *SearchData) webSearch(question string, querySize int) (results []dto_search.CommonSearchOutput, err error) {
	searchList, err := entity.webSearchApi.Search(question)
	if err != nil {
		return nil, err
	}
	if len(searchList) >= querySize {
		searchList = searchList[:querySize]
	}
	for _, resp := range searchList {
		results = append(results, dto_search.CommonSearchOutput{
			Title:   resp.Title,
			Url:     resp.Url,
			Content: resp.Content,
			Form:    "web",
		})
	}
	return results, nil
}

func appendText(source *es.DocDocument, fullContent string) string {
	prefixIndex := source.Start
	suffixIndex := source.End
	full := []rune(fullContent)
	prefixIndex = prefixIndex - 256
	if prefixIndex < 0 {
		prefixIndex = 0
	}
	suffixIndex = suffixIndex + 256
	if suffixIndex > len(full)-1 {
		suffixIndex = len(full) - 1
	}
	if suffixIndex < prefixIndex {
		suffixIndex = prefixIndex
	}
	return string(full[prefixIndex:suffixIndex])
}
