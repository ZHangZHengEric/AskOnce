package data

import (
	"askonce/api/jobd"
	"askonce/api/web_search"
	"askonce/components/dto/dto_search"
	"askonce/conf"
	"askonce/helpers"
	"askonce/models"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"gorm.io/datatypes"
	"sort"
	"strings"
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
	gptData          *GptData
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
	entity.gptData = entity.Create(new(GptData)).(*GptData)
}

func (entity *SearchData) SearchFromWebOrKnowledge(sessionId, question string, kdbId int64, userId string) (results []dto_search.CommonSearchOutput, err error) {
	results = make([]dto_search.CommonSearchOutput, 0)
	if kdbId == 0 { // web搜索
		searchList, err := entity.webSearchApi.Search(question)
		if err != nil {
			return nil, err
		}
		if len(searchList) >= 10 {
			searchList = searchList[:10]
		}
		for _, resp := range searchList {
			results = append(results, dto_search.CommonSearchOutput{
				Title:   resp.Title,
				Url:     resp.Url,
				Content: resp.Content,
			})
		}
	} else { // 知识库搜索
		kdb, err := entity.kdbData.CheckKdbAuth(kdbId, userId, models.AuthTypeRead)
		if err != nil {
			return nil, err
		}
		// es搜索的片段
		esSearchResult, err := entity.CommonEsSearch(EsCommonSearch{
			IndexName: kdb.GetIndexName(),
			Query:     question,
		})
		if err != nil {
			entity.LogErrorf("es搜索报错")
		}
		for _, result := range esSearchResult {
			results = append(results, dto_search.CommonSearchOutput{
				Title:   result.Title,
				Url:     result.Url,
				Content: result.Content,
			})
		}
	}
	if len(results) > 0 {
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

type EsCommonSearch struct {
	IndexName string
	Query     string
	Size      int
}

type EsCommonSearchResult struct {
	Id          int64
	Title       string
	Url         string
	Content     string
	FullContent string
	Score       float32
}

func (entity *SearchData) CommonEsSearch(input EsCommonSearch) (res []*EsCommonSearchResult, err error) {
	embRes, err := entity.gptData.Embedding([]string{input.Query})
	if err != nil {
		return
	}
	querySize := 10
	esDbConfigStr := strings.Replace(conf.WebConf.EsDbConfig, "@indexName@", input.IndexName, 1)
	recalls, err := entity.jobdApi.EsSearch(embRes[0], input.Query, querySize, esDbConfigStr)
	if err != nil {
		return
	}
	if len(recalls) == 0 {
		return
	}
	sort.Slice(recalls, func(i, j int) bool {
		return recalls[i].Score > recalls[j].Score
	})
	// 去重一次
	uniqueMap := make(map[string]int)
	recalls2 := make([]jobd.ESSearchOutput, 0)
	for _, recall := range recalls {
		uniqueCode := fmt.Sprintf("%v%s", recall.Source.DocId, recall.Source.DocContent)
		if _, ok := uniqueMap[uniqueCode]; !ok {
			recalls2 = append(recalls2, recall)
			uniqueMap[uniqueCode] = 1
		}
	}
	var dataIds []int64
	dataSearchMap := make(map[int]jobd.SearchOutputSource)
	for i, result := range recalls2 {
		dataIds = append(dataIds, result.Source.DocId)
		dataSearchMap[i] = result.Source
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
	for i, result := range recalls2 {
		dataIdInt := result.Source.DocId
		ddd, ok := docMap[dataIdInt]
		if !ok {
			continue
		}
		out := &EsCommonSearchResult{}
		out.Id = result.Source.DocId
		out.Title = ddd.DocName
		out.Url = filePathMap[ddd.SourceId]
		out.Content = appendText(dataSearchMap[i], dataContentMap[dataIdInt])
		out.FullContent = dataContentMap[dataIdInt]
		out.Score = result.Score
		res = append(res, out)
	}
	returnSize := 10
	if input.Size > 0 {
		returnSize = input.Size
	}
	if len(res) <= returnSize {
		return
	} else {
		res = res[0 : returnSize-1]
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

func (entity *SearchData) SearchFromWeb(sessionId string, question string) (results []dto_search.CommonSearchOutput, err error) {
	searchList, err := entity.webSearchApi.Search(question)
	if err != nil {
		return nil, err
	}
	if len(searchList) >= 10 {
		searchList = searchList[:10]
	}
	for _, resp := range searchList {
		results = append(results, dto_search.CommonSearchOutput{
			Title:   resp.Title,
			Url:     resp.Url,
			Content: resp.Content,
		})
	}
	if len(results) > 0 {
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
	return results, nil
}

func appendText(source jobd.SearchOutputSource, fullContent string) string {
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
