package service

import (
	"askonce/api"
	"askonce/api/jobd"
	"askonce/components"
	"askonce/components/dto"
	"askonce/components/dto/dto_kdb_doc"
	"askonce/components/dto/dto_search"
	"askonce/data"
	"askonce/helpers"
	"askonce/models"
	"askonce/utils"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/russross/blackfriday/v2"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/errors"
	"github.com/xiangtao94/golib/pkg/orm"
	"github.com/xiangtao94/golib/pkg/sse"
	"github.com/xiangtao94/golib/pkg/zlog"
	"golang.org/x/sync/errgroup"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SearchService struct {
	flow.Service
	askInfoDao       *models.AskInfoDao
	kdbDocContentDao *models.KdbDocContentDao

	kdbData      *data.KdbData
	searchData   *data.SearchData
	jobdApi      *jobd.JobdApi
	askAttachDao *models.AskAttachDao
	processDao   *models.AskProcessDao
	userDao      *models.UserDao
}

func (s *SearchService) OnCreate() {
	s.jobdApi = s.Create(new(jobd.JobdApi)).(*jobd.JobdApi)
	s.askInfoDao = s.Create(new(models.AskInfoDao)).(*models.AskInfoDao)
	s.kdbDocContentDao = s.Create(new(models.KdbDocContentDao)).(*models.KdbDocContentDao)
	s.askAttachDao = s.Create(new(models.AskAttachDao)).(*models.AskAttachDao)
	s.userDao = s.Create(new(models.UserDao)).(*models.UserDao)
	s.processDao = s.Create(new(models.AskProcessDao)).(*models.AskProcessDao)
	s.searchData = s.Create(new(data.SearchData)).(*data.SearchData)
	s.kdbData = s.Create(new(data.KdbData)).(*data.KdbData)
}

func (s *SearchService) EchoRes(stage, text string) {
	echoResStr, _ := sonic.MarshalString(dto_search.AskRes{
		Stage: stage,
		Text:  text,
	})
	time.Sleep(80 * time.Millisecond)
	sse.RenderStream(s.GetCtx(), "0", "message", echoResStr)
}

func (s *SearchService) KdbList(req *dto_search.KdbListReq) (res *dto_search.KdbListRes, err error) {
	res = &dto_search.KdbListRes{
		List: make([]dto_search.KdbListItem, 0),
	}
	userInfo, _ := utils.LoginInfo(s.GetCtx())
	all, _, err := s.kdbData.GetKdbList(userInfo.UserId, req.Query, dto.PageParam{
		PageNo:   1,
		PageSize: 10000,
	})
	if err != nil {
		return nil, err
	}
	// 排序
	if req.OrderType == 2 { // 创建时间正序
		sort.Slice(all, func(i, j int) bool {
			return all[i].CreatedAt.Before(all[j].CreatedAt)
		})
	} else if req.OrderType == 3 && len(userInfo.UserId) > 0 { // 最近常用
		latestKdbIds, err := s.askInfoDao.GetUserLatestKdb(userInfo.UserId)
		if err != nil {
			return nil, err
		}
		latestKdbIds = slice.Unique(latestKdbIds)
		weigthMap := make(map[int64]int)
		for i, id := range latestKdbIds {
			weigthMap[id] = len(latestKdbIds) - i
		}
		sort.Slice(all, func(i, j int) bool {
			return weigthMap[all[i].Id] > weigthMap[all[j].Id]
		})
	}
	for _, kdb := range all {
		res.List = append(res.List, dto_search.KdbListItem{
			KdbId:      kdb.Id,
			KdbName:    kdb.Name,
			CreateTime: kdb.CreatedAt.Format(time.DateTime),
		})
	}
	res.Total = int64(len(res.List))

	start, end := utils.SlicePage(req.PageNo, req.PageSize, len(res.List)) //第一页1页显示3条数据
	res.List = res.List[start:end]                                         //  分页后的数据
	return
}

func (s *SearchService) Case(req *dto_search.CaseReq) (res *dto_search.CaseRes, err error) {
	res = &dto_search.CaseRes{
		Cases: make([]string, 0),
	}
	if req.KdbId == 0 {
		hots := helpers.BaiduHotTest(s.GetCtx())
		filterHots := []string{}
		filterHots = hots
		randShuffle(filterHots)
		if len(filterHots) > 5 {
			filterHots = filterHots[0:5]
		}
		res.Cases = filterHots
	} else {
		kdb, _ := s.Create(new(models.KdbDao)).(*models.KdbDao).GetById(req.KdbId)
		if kdb == nil {
			return
		}
		if len(kdb.Setting.Data().KdbAttach.Cases) > 0 {
			filter := kdb.Setting.Data().KdbAttach.Cases
			randShuffle(filter)
			res.Cases = filter
		}
	}
	return
}

func (s *SearchService) Session(req *dto.EmptyReq) (res *dto_search.GenSessionRes, err error) {
	userInfo, _ := utils.LoginInfo(s.GetCtx())
	askInfo, err := s.searchData.CreateSession(userInfo.UserId)
	if err != nil {
		return nil, err
	}
	res = &dto_search.GenSessionRes{
		SessionId: askInfo.SessionId,
	}
	return
}

type AskContext struct {
	ModelType, PromptI18n string
	SessionId             string
	Question              string // 问题
	UserId                string
	AnswerStyle           string
	Kdb                   *models.Kdb
	DbData                *models.AskInfo
	Outline               []jobd.Outline
	searchResult          []dto_search.CommonSearchOutput
	// 执行进度
	Process []*models.AskProcess
}

func (a *AskContext) GetKdbIndex() string {
	if a.Kdb == nil {
		return ""
	}
	return a.Kdb.GetIndexName()
}

//      analyze: "问题分析",
//      webSearch: "全网搜索",
//      kdbSearch: "知识库搜索",
//      summary: "整理答案",
//      finish: "回答完成"

func (a *AskContext) AppendProcess(stageType string, appendText ...string) {
	prefixContent := ""
	if stageType == "search" {
		if a.Kdb != nil {
			prefixContent = "知识库"
			stageType = "kdbSearch"
		} else {
			prefixContent = "互联网"
			stageType = "webSearch"
		}
	}
	a.Process = append(a.Process, &models.AskProcess{
		SessionId: a.SessionId,
		Type:      stageType,
		Time:      time.Now().UnixMilli(),
		Content:   fmt.Sprint(prefixContent, appendText),
		CrudModel: orm.CrudModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	})
}

func (s *SearchService) Ask(req *dto_search.AskReq) (err error) {
	userInfo, _ := utils.LoginInfo(s.GetCtx())
	var kdb *models.Kdb
	askDirect := false
	if req.KdbId > 0 {
		// 校验知识库权限
		kdb, err = s.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeRead)
		if err != nil {
			return err
		}
	} else {
		// 互联网判断是否要搜索
		useRag, _ := s.jobdApi.NetRagAssessment(req.Question)
		if !useRag {
			askDirect = true
		}
	}
	// 判断session是否存在
	askInfoDao := s.Create(new(models.AskInfoDao)).(*models.AskInfoDao)
	askInfo, err := askInfoDao.GetBySessionId(req.SessionId)
	if err != nil {
		return err
	}
	if askInfo == nil {
		return components.ErrorAskSessionNoExist
	}
	askInfo.Question = req.Question
	askInfo.AskType = req.Type
	askInfo.KdbId = req.KdbId
	err = askInfoDao.Update(askInfo)
	if err != nil {
		return err
	}
	user, _ := s.userDao.GetByUserId(userInfo.UserId)
	config := user.Setting.Data()
	modelType := config.ModelType
	promptI18n := "\n 输出使用中文！！"
	if config.Language == "en-us" {
		promptI18n = "\n 输出使用英文！！"
	}
	answerStyle := resolveAnswerStyle(req.Type, req.KdbId)
	askContext := &AskContext{
		ModelType:   modelType,
		PromptI18n:  promptI18n,
		SessionId:   req.SessionId,
		Question:    req.Question,
		Kdb:         kdb,
		DbData:      askInfo,
		UserId:      userInfo.UserId,
		AnswerStyle: answerStyle,
	}
	s.EchoRes("start", "")
	var answer string
	subQuestion := make([]string, 0)
	echoRefers := make([]dto_search.DoReferItem, 0)
	if askDirect {
		answer, echoRefers, err = s.askDirect(askContext)
	} else {
		switch req.Type {
		case "simple":
			answer, echoRefers, err = s.askSimple(askContext)
		case "complex":
			answer, echoRefers, subQuestion, err = s.askComplex(askContext)
		case "research":
			if req.KdbId > 0 {
				answer, echoRefers, subQuestion, err = s.askComplex(askContext)
			} else {
				answer, echoRefers, subQuestion, err = s.askProfessional(askContext)
			}
		default:
			return errors.ErrorParamInvalid
		}
	}
	if err != nil {
		s.LogErrorf("问答报错, %s", err.Error())
		_ = askInfoDao.UpdateById(askInfo.Id, map[string]interface{}{"status": models.AskInfoStatusFail})
		return
	}
	// 保存记录
	go func(entity *SearchService) {
		_ = entity.askRecordUpdate(askInfo, subQuestion, answer, echoRefers)
		_ = entity.processDao.BatchInsert(askContext.Process)
	}(s.CopyWithCtx(s.GetCtx()).(*SearchService))
	askContext.AppendProcess("finish", "回答完成")
	s.EchoRes("done", "问答结束")
	return
}

func (s *SearchService) askDirect(askContext *AskContext) (answer string, echoRefers []dto_search.DoReferItem, err error) {
	// 开始回答
	answer, echoRefers, err = s.askByDocument(askContext)
	if err != nil {
		err = components.ErrorChatError
		return
	}
	// 生成大纲
	s.askOutline(askContext.SessionId, answer)
	return
}

// 简单搜索
func (s *SearchService) askSimple(askContext *AskContext) (answer string, echoRefers []dto_search.DoReferItem, err error) {
	s.EchoRes("search", "")
	askContext.AppendProcess("search", "搜索开始")
	searchResult, err := s.searchData.DocSearch(askContext.SessionId, askContext.Question, new(data.DocSearchOptions).WithIndex(askContext.GetKdbIndex()))
	if err != nil {
		err = components.ErrorQueryError
		return
	}
	echoContent := fmt.Sprintf("搜索到%v条相关内容", len(searchResult))
	s.EchoRes("search", echoContent)
	err = s.referenceUpdate(askContext.SessionId, searchResult)
	askContext.AppendProcess("search", echoContent)
	// 开始回答
	askContext.searchResult = searchResult
	answer, echoRefers, err = s.askByDocument(askContext)
	if err != nil {
		err = components.ErrorChatError
		return
	}
	// 生成大纲
	s.askOutline(askContext.SessionId, answer)
	return
}

func (s *SearchService) askComplex(askContext *AskContext) (answer string, echoRefers []dto_search.DoReferItem, questions []string, err error) {
	s.EchoRes("analyze", fmt.Sprintf("开始分析问题：%s", askContext.Question))
	askContext.AppendProcess("analyze", fmt.Sprintf("开始分析问题：%s", askContext.Question))
	splitRes, err := s.jobdApi.SplitQuestion(askContext.Question)
	if err != nil {
		err = components.ErrorJobdError
		return
	}
	questions = splitRes.Questions
	if len(questions) == 0 {
		questions = append(questions, askContext.Question)
	}
	s.EchoRes("analyze", fmt.Sprintf("分析问题为: %s", strings.Join(questions, ";")))
	askContext.AppendProcess("analyze", fmt.Sprintf("分析问题为: %s", strings.Join(questions, ";")))
	if len(questions) == 1 {
		answer, echoRefers, err = s.askSimple(askContext)
		return
	}
	s.EchoRes("search", "")
	askContext.AppendProcess("search", "搜索开始")
	searchResultAll := make([]dto_search.CommonSearchOutput, 0)
	searchResultAllMap := make(map[string][]dto_search.CommonSearchOutput)
	// 处理拆分问题的单个回答
	eg0, _ := errgroup.WithContext(s.GetCtx())
	lock0 := sync.Mutex{}
	for i, subQ := range splitRes.Questions {
		tmpQ := subQ
		tmpSearchContent := splitRes.Questions[i]
		eg0.Go(func() (err error) {
			searchResult, err := s.searchData.DocSearch(askContext.SessionId, tmpSearchContent, new(data.DocSearchOptions).WithIndex(askContext.GetKdbIndex()))
			if err != nil {
				return err
			}
			lock0.Lock()
			searchResultAllMap[tmpQ] = searchResult
			lock0.Unlock()
			return nil
		})
	}
	if err = eg0.Wait(); err != nil {
		err = components.ErrorQueryError
		return
	}
	// 合并答案
	searchResultUnique := make(map[string]bool)
	for _, outputs := range searchResultAllMap {
		for _, output := range outputs {
			unique := base64.StdEncoding.EncodeToString([]byte(output.Url + output.Content))
			if searchResultUnique[unique] {
				continue
			}
			searchResultAll = append(searchResultAll, output)
			searchResultUnique[unique] = true
		}
	}
	s.EchoRes("search", fmt.Sprintf("搜索到%v条相关内容", len(searchResultAll)))
	askContext.AppendProcess("search", fmt.Sprintf("搜索到%v条相关内容", len(searchResultAll)))
	err = s.referenceUpdate(askContext.SessionId, searchResultAll)
	if err != nil {
		err = components.ErrorMysqlError
		return
	}
	// 开始回答
	askContext.searchResult = searchResultAll
	answer, echoRefers, err = s.askByDocument(askContext)
	if err != nil {
		err = components.ErrorChatError
		return
	}
	// 生成大纲
	s.askOutline(askContext.SessionId, answer)
	return
}

func (s *SearchService) askProfessional(askContext *AskContext) (answer string, echoRefers []dto_search.DoReferItem, questions []string, err error) {
	s.EchoRes("analyze", fmt.Sprintf("开始分析问题：%s", askContext.Question))
	askContext.AppendProcess("analyze", fmt.Sprintf("开始分析问题：%s", askContext.Question))
	splitRes, err := s.jobdApi.SplitQuestion(askContext.Question)
	if err != nil {
		err = components.ErrorJobdError
		return
	}
	questions = splitRes.Questions
	if len(questions) == 0 {
		questions = append(questions, askContext.Question)
	}
	s.EchoRes("analyze", fmt.Sprintf("分析问题为: %s", strings.Join(questions, ";")))
	askContext.AppendProcess("analyze", fmt.Sprintf("分析问题为: %s", strings.Join(questions, ";")))
	if len(questions) == 1 {
		answer, echoRefers, err = s.askSimple(askContext)
		return
	}
	s.EchoRes("search", "")
	askContext.AppendProcess("search", "搜索开始")
	searchResultAll := make([]dto_search.CommonSearchOutput, 0)
	searchResultAllMap := make(map[string][]dto_search.CommonSearchOutput)
	// 处理拆分问题的单个回答
	eg0, _ := errgroup.WithContext(s.GetCtx())
	lock0 := sync.Mutex{}
	for i, subQ := range splitRes.Questions {
		tmpQ := subQ
		tmpSearchContent := splitRes.Questions[i]
		eg0.Go(func() (err error) {
			searchResult, err := s.searchData.DocSearch(askContext.SessionId, tmpSearchContent, new(data.DocSearchOptions).WithIndex(askContext.GetKdbIndex()))
			if err != nil {
				return err
			}
			lock0.Lock()
			searchResultAllMap[tmpQ] = searchResult
			lock0.Unlock()
			return nil
		})
	}
	if err = eg0.Wait(); err != nil {
		err = components.ErrorQueryError
		return
	}
	// 合并答案
	searchResultUnique := make(map[string]bool)
	for _, outputs := range searchResultAllMap {
		for _, output := range outputs {
			unique := base64.StdEncoding.EncodeToString([]byte(output.Url + output.Content))
			if searchResultUnique[unique] {
				continue
			}
			searchResultAll = append(searchResultAll, output)
			searchResultUnique[unique] = true
		}
	}
	s.EchoRes("search", fmt.Sprintf("搜索到%v条相关内容", len(searchResultAll)))
	askContext.AppendProcess("search", fmt.Sprintf("搜索到%v条相关内容", len(searchResultAll)))
	err = s.referenceUpdate(askContext.SessionId, searchResultAll)
	if err != nil {
		err = components.ErrorMysqlError
		return
	}
	// 开始回答
	askContext.searchResult = searchResultAll
	answer, echoRefers, err = s.askByDocument(askContext)
	if err != nil {
		err = components.ErrorChatError
		return
	}
	// 生成大纲
	s.askOutline(askContext.SessionId, answer)
	return
}

func IsCompleted(answer string, status string, doneAnswer string) (string, int) {
	begin := len([]rune(doneAnswer))
	if status == "FINISH" { // 对话完成了，返回剩下所有的
		return strings.Replace(answer, doneAnswer, "", 1), begin
	}
	waitJudge := strings.Replace(answer, doneAnswer, "", 1) // 待判断文字
	if len(waitJudge) == 0 {
		return "", 0
	}
	ss1 := strings.Split(waitJudge, "\n\n")
	if len(ss1) >= 2 && len(ss1[0]) > 0 { // 有双换行符，且前面有文字，则直接返回
		return ss1[0] + "\n\n", begin
	}
	ss2 := strings.Split(waitJudge, "\n")
	if len(ss2) >= 2 && len(ss2[0]) > 0 { // 有换行符，且前面有文字，则直接返回
		return ss2[0] + "\n", begin
	}
	r1 := regexp.MustCompile("。|？|！|；|;|!|\\?") // 强句子判断
	matchIndexs := r1.FindAllStringIndex(waitJudge, -1)
	if len(matchIndexs) == 0 { // 没有找到切分的下标， 返回
		return "", begin
	}
	firstMatch := matchIndexs[0]
	matchText := waitJudge[0:firstMatch[1]]
	return matchText, begin
}

func (s *SearchService) referenceUpdate(sessionId string, searchResult []dto_search.CommonSearchOutput) (err error) {
	askAttach, err := s.askAttachDao.GetBySessionId(sessionId)
	if err != nil {
		return err
	}
	if askAttach == nil {
		return
	}
	refers := make([]dto_search.CommonSearchOutput, 0)
	_ = json.Unmarshal(askAttach.Reference, &refers)
	nMap := make(map[string]dto_search.CommonSearchOutput)
	for _, output := range searchResult {
		if output.Form == "web" {
			nMap[output.Title+output.Url] = output
		} else {
			nMap[strconv.FormatInt(output.DocSegmentId, 10)] = output
		}
	}
	for _, refer := range refers {
		if refer.Form == "web" {
			if _, ok := nMap[refer.Title+refer.Url]; ok {
				delete(nMap, refer.Title+refer.Url)
			}
		} else {
			if _, ok := nMap[strconv.FormatInt(refer.DocSegmentId, 10)]; ok {
				delete(nMap, strconv.FormatInt(refer.DocSegmentId, 10))
			}
		}
	}
	for _, v := range nMap {
		refers = append(refers, v)
	}
	refersAfterStr, _ := json.Marshal(refers)
	err = s.askAttachDao.UpdateBySessionId(sessionId, map[string]interface{}{"reference": refersAfterStr})
	if err != nil {
		return
	}
	return nil
}

func (s *SearchService) referDo(begin int, needReference string, searchResult []dto_search.CommonSearchOutput, setting dto.KdbSetting) (output []dto_search.DoReferItem, err error) {
	output = make([]dto_search.DoReferItem, 0)
	if len(searchResult) == 0 {
		return
	}
	referStrList := []string{}
	for _, o := range searchResult {
		referStrList = append(referStrList, o.Content)
	}

	referenceRes, err := s.jobdApi.ResultAddReference(needReference, referStrList, setting.ReferenceThreshold)
	if err != nil {
		zlog.Errorf(s.GetCtx(), "ResultAddReference error %s", err.Error())
		return
	}
	for _, referenceMap := range referenceRes.ReferenceMap {
		if len(referenceMap.IndexRange) == 2 { // 处理总结每段文字的引用点
			// todo 找寻referIndexEnd 后第一个句号或者换行
			r1 := regexp.MustCompile("。|？|！|；|;|!|\\?|\n") // 强句子判断
			matchIndexs := r1.FindStringIndex(needReference)
			length := len([]rune(needReference))
			if len(matchIndexs) > 0 { //找到了
				ttt := needReference[0:matchIndexs[0]]
				length = len([]rune(ttt))
			}
			if length < referenceMap.IndexRange[1] {
				length = referenceMap.IndexRange[1]
			}
			t := dto_search.DoReferItem{
				Start:       begin + referenceMap.IndexRange[0],
				End:         begin + referenceMap.IndexRange[1],
				NumberIndex: begin + length,
				Refers:      nil,
			}
			refers := []dto_search.DoReferReferItem{}
			if len(referenceMap.ReferenceList) == 0 {
				continue
			}
			for _, index := range referenceMap.ReferenceList {
				index2, ok := referenceRes.ReferenceListSelectIndex[strconv.Itoa(index)]
				if !ok { // 文档引用来源 字段中不存在 跳过
					continue
				}
				if len(index2) == 2 {
					refers = append(refers, dto_search.DoReferReferItem{
						Index:      index,
						ReferStart: index2[0],
						ReferEnd:   index2[1],
					})
				}
			}
			if len(refers) == 0 {
				continue
			}
			t.Refers = refers
			output = append(output, t)

		}
	}
	return output, nil
}

func (s *SearchService) History(req *dto_search.HisReq) (res *dto_search.HisRes, err error) {
	askInfo, err := s.askInfoDao.GetBySessionId(req.SessionId)
	if err != nil {
		return
	}
	res = &dto_search.HisRes{}
	if askInfo == nil {
		return
	}
	if askInfo.Status == models.AskInfoStatusFail {
		return nil, components.ErrorAskSessionError
	}
	var mapResult = make(map[string]string)
	_ = json.Unmarshal(askInfo.Answer, &mapResult)
	unlike := false
	if askInfo.LikeStatus == 2 {
		unlike = true
	}
	res = &dto_search.HisRes{
		SessionId:    askInfo.SessionId,
		Unlike:       unlike,
		Question:     askInfo.Question,
		Result:       mapResult["new"],
		ResultRefers: mapResult["refers"],
	}

	return
}

func (s *SearchService) Reference(req *dto_search.ReferReq) (res *dto_search.ReferenceRes, err error) {
	askInfo, err := s.askInfoDao.GetBySessionId(req.SessionId)
	if err != nil {
		return
	}
	res = &dto_search.ReferenceRes{
		List: make([]dto_search.CommonSearchOutput, 0),
	}
	if askInfo == nil {
		return
	}
	var mapResult = make(map[string]string)
	_ = json.Unmarshal(askInfo.Answer, &mapResult)
	askAttach, err := s.askAttachDao.GetBySessionId(req.SessionId)
	if err != nil {
		return
	}
	if askAttach == nil {
		return
	}
	refers := make([]dto_search.CommonSearchOutput, 0)
	_ = json.Unmarshal(askAttach.Reference, &refers)
	res.List = refers
	return
}

func (s *SearchService) Outline(req *dto_search.OutlineReq) (res *dto_search.OutlineRes, err error) {
	askInfo, err := s.askInfoDao.GetBySessionId(req.SessionId)
	if err != nil {
		return
	}
	res = &dto_search.OutlineRes{
		List: make([]dto_search.OutlineItem, 0),
	}
	if askInfo == nil {
		return
	}
	var mapResult = make(map[string]string)
	_ = json.Unmarshal(askInfo.Answer, &mapResult)
	askAttach, err := s.askAttachDao.GetBySessionId(req.SessionId)
	if err != nil {
		return
	}
	if askAttach == nil {
		return
	}
	outlines := make([]dto_search.OutlineItem, 0)
	_ = json.Unmarshal(askAttach.Outline, &outlines)
	res.List = outlines
	return
}

func (s *SearchService) Unlike(req *dto_search.UnlikeReq) (res interface{}, err error) {
	askInfoDao := s.Create(new(models.AskInfoDao)).(*models.AskInfoDao)
	askInfo, err := askInfoDao.GetBySessionId(req.SessionId)
	if err != nil {
		return
	}
	if askInfo == nil {
		return
	}
	reasonsStr, _ := json.Marshal(req.Reasons)
	err = askInfoDao.UpdateById(askInfo.Id, map[string]interface{}{"like_status": 2, "like_reasons": reasonsStr})
	if err != nil {
		return
	}
	return
}

func (s *SearchService) askByDocument(req *AskContext) (answer string, echoRefers []dto_search.DoReferItem, err error) {
	req.AppendProcess("summary", fmt.Sprintf("回答问题【%s】答案开始", req.Question))
	searchResult := req.searchResult
	// 生成答案 + 引用
	alreadyReferAnswer := ""
	wg := sync.WaitGroup{}
	lock := sync.Mutex{}
	first := true
	currentAnswer := ""
	err = s.jobdApi.AnswerByDocuments(req.SessionId, req.Question, req.AnswerStyle, searchResult, func(jobdRes jobd.JobdCommonRes) error {
		if first {
			s.EchoRes("generate", "")
			first = false
		}
		chatAnswer := jobd.AnswerByDocumentsRes{}
		_ = json.Unmarshal([]byte(jobdRes.Output), &chatAnswer)
		// 对话展示逻辑
		echoAnswer := chatAnswer.Answer
		currentAnswer = currentAnswer + echoAnswer
		s.EchoRes("appendText", echoAnswer)
		answer = currentAnswer
		if len(searchResult) > 0 {
			// 引用判断逻辑
			needReference, begin := IsCompleted(currentAnswer, jobdRes.Status, alreadyReferAnswer)
			if len(needReference) > 0 {
				// 重新查一次数据库对引用
				attach, err := s.askAttachDao.GetBySessionId(req.SessionId)
				if err != nil {
					return err
				}
				if attach != nil {
					tmpResult := make([]dto_search.CommonSearchOutput, 0)
					_ = json.Unmarshal(attach.Reference, &tmpResult)
					if len(tmpResult) > len(searchResult) {
						s.EchoRes("refreshSearch", fmt.Sprintf("再次搜索到%v条相关内容", len(tmpResult)-len(searchResult)))
						searchResult = tmpResult
					}
				}
				s.LogInfof("完整句子: %s。开始位置: %v", needReference, begin)
				wg.Add(1)
				alreadyReferAnswer = alreadyReferAnswer + needReference
				go func(entity *SearchService, begin int, needRefer string, searchResult []dto_search.CommonSearchOutput) {
					defer wg.Done()
					kdbSetting := dto.KdbSetting{}
					if req.Kdb != nil {
						kdbSetting = req.Kdb.Setting.Data()
					}
					aa, errA := entity.referDo(begin, needRefer, searchResult, kdbSetting)
					if errA != nil {
						return
					}
					if len(aa) == 0 {
						return
					}
					lock.Lock()
					echoRefers = append(echoRefers, aa...)
					lock.Unlock()
					sort.Slice(echoRefers, func(i, j int) bool {
						return echoRefers[i].Start < echoRefers[j].Start
					})
					// 	合并一次
					echoRefers = mergeItems(echoRefers)
					if len(echoRefers) == 0 {
						return
					}
					aaStr, _ := json.Marshal(echoRefers)
					entity.EchoRes("refer", string(aaStr))
				}(s.CopyWithCtx(s.GetCtx()).(*SearchService), begin, needReference, searchResult)
			}
		}
		return nil
	})
	wg.Wait()
	req.AppendProcess("summary", fmt.Sprintf("回答问题【%s】答案结束", req.Question))
	s.EchoRes("complete", answer)
	return
}

func (s *SearchService) askOutline(sessionId string, answer string) {
	if answer == "" {
		return
	}
	outlineRes, errA := s.jobdApi.AnswerOutline(answer)
	if errA != nil {
		s.LogErrorf("生成大纲失败%s", errA.Error())
		return
	}
	outlineStr, _ := json.Marshal(outlineRes.AnswerOutline)
	err := s.askAttachDao.UpdateBySessionId(sessionId, map[string]interface{}{"outline": outlineStr})
	if err != nil {
		return
	}
	s.EchoRes("outline", "done")
}

func (s *SearchService) askRecordUpdate(askInfo *models.AskInfo, questions []string, answer string, echoRefers []dto_search.DoReferItem) (err error) {
	askInfo.SubQuestion = questions
	resultMap := make(map[string]string)
	resultMap["new"] = answer
	echoRefersStr, _ := json.Marshal(echoRefers)
	resultMap["refers"] = string(echoRefersStr)
	resultMapStr, _ := json.Marshal(resultMap)
	askInfo.Answer = resultMapStr
	askInfo.Status = models.AskInfoStatusSuccess
	err = s.askInfoDao.UpdateEntity(askInfo)
	return
}

func (s *SearchService) Relation(req *dto_search.RelationReq) (res *dto_search.RelationRes, err error) {
	askInfo, err := s.askInfoDao.GetBySessionId(req.SessionId)
	if err != nil {
		return
	}
	res = &dto_search.RelationRes{
		EventsInfo: make([]dto_search.RelateEventsInfo, 0),
		PeopleInfo: make([]dto_search.RelatePeopleInfo, 0),
		OrgsInfo:   make([]dto_search.RelateOrgInfo, 0),
	}
	if askInfo == nil {
		return
	}
	//var mapResult = make(map[string]string)
	//_ = json.Unmarshal(askInfo.Answer, &mapResult)
	//askAttach, err := s.askAttachDao.GetBySessionId(req.SessionId)
	//if err != nil {
	//	return
	//}
	//if askAttach == nil {
	//	return
	//}
	//relationObj := &jobd.GenerateRelateInfoRes{}
	//_ = json.Unmarshal(askAttach.Relation, &relationObj)
	//for _, info := range relationObj.PeopleInfo {
	//	res.PeopleInfo = append(res.PeopleInfo, dto_search.RelatePeopleInfo{
	//		PersonName:     info.PersonName,
	//		PersonDescribe: info.PersonDescribe,
	//	})
	//}
	//for _, info := range relationObj.EventsInfo {
	//	res.EventsInfo = append(res.EventsInfo, dto_search.RelateEventsInfo{
	//		EventName:     info.EventName,
	//		EventDate:     info.EventDate,
	//		EventDescribe: info.EventDescribe,
	//	})
	//}
	//for _, info := range relationObj.OrgsInfo {
	//	res.OrgsInfo = append(res.OrgsInfo, dto_search.RelateOrgInfo{
	//		OrgName:     info.OrgName,
	//		OrgDescribe: info.OrgDescribe,
	//	})
	//}
	return
}

func (s *SearchService) Process(req *dto_search.ProcessReq) (res *dto_search.ProcessRes, err error) {

	askInfo, err := s.askInfoDao.GetBySessionId(req.SessionId)
	if err != nil {
		return
	}
	res = &dto_search.ProcessRes{List: make([]dto_search.ProcessItem, 0)}
	if askInfo == nil {
		return
	}
	process, err := s.processDao.GetBySessionId(askInfo.SessionId)
	if err != nil {
		return
	}
	for _, askProcess := range process {
		res.List = append(res.List, dto_search.ProcessItem{
			StageName: models.ProcessTypeNameMap[askProcess.Type],
			StageType: askProcess.Type,
			Content:   askProcess.Content,
			Time:      askProcess.Time,
		})
	}
	return
}

func (s *SearchService) Recall(req *dto_kdb_doc.RecallReq) (res *dto_kdb_doc.RecallRes, err error) {
	userInfo, _ := utils.LoginInfo(s.GetCtx())

	kdb, err := s.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeRead)
	if err != nil {
		return nil, err
	}
	res = &dto_kdb_doc.RecallRes{
		List: make([]dto_kdb_doc.RecallItem, 0),
	}

	// es搜索的片段
	esSearchResult, err := s.searchData.DocSearch("", req.Query, new(data.DocSearchOptions).WithIndex(kdb.GetIndexName()))
	if err != nil {
		return
	}
	var docIds []int64
	for _, result := range esSearchResult {
		docIds = append(docIds, result.DocId)
	}
	dataContents, err := s.kdbDocContentDao.GetByDataIds(docIds)
	if err != nil {
		return nil, err
	}
	dataContentMap := make(map[int64]string)
	for _, content := range dataContents {
		dataContentMap[content.DocId] = content.Content
	}
	for _, result := range esSearchResult {
		res.List = append(res.List, dto_kdb_doc.RecallItem{
			DataName:      result.Title,
			DataPath:      result.Url,
			SearchContent: result.Content,
			DataContent:   dataContentMap[result.DocId],
		})
	}
	return
}

func mergeItems(items []dto_search.DoReferItem) []dto_search.DoReferItem {
	if len(items) == 0 {
		return nil
	}

	mergedItems := make([]dto_search.DoReferItem, 0)
	currentItem := items[0]

	for i := 1; i < len(items); i++ {
		if items[i].NumberIndex == currentItem.NumberIndex && items[i].Start+1 == currentItem.End {
			currentItem.End = items[i].End
		} else {
			mergedItems = append(mergedItems, currentItem)
			currentItem = items[i]
		}
	}

	mergedItems = append(mergedItems, currentItem)

	return mergedItems
}

func (s *SearchService) AskSync(req *dto_search.ChatAskReq) (res *dto_search.AskSyncRes, err error) {
	userInfo, _ := utils.LoginInfo(s.GetCtx())
	// 校验知识库权限
	kdb, err := s.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeRead)
	if err != nil {
		return nil, err
	}
	askInfo, err := s.searchData.CreateSession(userInfo.UserId)
	if err != nil {
		return nil, err
	}
	_ = s.askInfoDao.UpdateById(askInfo.Id, map[string]interface{}{"question": req.Question, "ask_type": "detail", "kdb_id": req.KdbId})
	user, _ := s.userDao.GetByUserId(userInfo.UserId)
	config := user.Setting.Data()
	modelType := config.ModelType
	promptI18n := "\n 输出使用中文！！"
	if config.Language == "en-us" {
		promptI18n = "\n 输出使用英文！！"
	}
	askContext := &AskContext{
		ModelType:  modelType,
		PromptI18n: promptI18n,
		SessionId:  askInfo.SessionId,
		Question:   req.Question,
		Kdb:        kdb,
		DbData:     askInfo,
		UserId:     userInfo.UserId,
	}
	askContext.AnswerStyle = "detail"
	answer, answerRefer, searchResult, err := s.AskSyncDo(askContext)
	if err != nil {
		s.LogErrorf("问答报错, %s", err.Error())
		s.askInfoDao.UpdateById(askInfo.Id, map[string]interface{}{"status": models.AskInfoStatusFail})
		return
	}
	askAttach, err := s.askAttachDao.GetBySessionId(askInfo.SessionId)
	if err != nil {
		return
	}
	if askAttach == nil {
		return
	}
	refers := make([]dto_search.CommonSearchOutput, 0)
	_ = json.Unmarshal(askAttach.Reference, &refers)
	res = &dto_search.AskSyncRes{
		Answer:       answer,
		AnswerRefer:  answerRefer,
		SearchResult: searchResult,
	}
	return
}

// 同步回答
func (s *SearchService) AskSyncDo(askContext *AskContext) (answer string, echoRefers []dto_search.DoReferItem, searchResult []dto_search.CommonSearchOutput, err error) {
	askContext.AppendProcess("search", "搜索开始")
	searchResult, err = s.searchData.DocSearch(askContext.SessionId, askContext.Question, new(data.DocSearchOptions).WithIndex(askContext.GetKdbIndex()).WithReturnFull(true))
	if err != nil {
		err = components.ErrorQueryError
		return
	}
	err = s.referenceUpdate(askContext.SessionId, searchResult)
	if err != nil {
		err = components.ErrorMysqlError
		return
	}
	askContext.AppendProcess("search", fmt.Sprintf("知识库搜索到%v条相关内容", len(searchResult)))
	askContext.AppendProcess("summary", fmt.Sprintf("回答问题【%s】答案开始", askContext.Question))
	answer, err = s.jobdApi.AnswerByDocumentsSync(askContext.SessionId, askContext.Question, askContext.AnswerStyle, searchResult, askContext.Outline)
	if err != nil {
		return
	}
	askContext.AppendProcess("summary", fmt.Sprintf("回答问题【%s】答案结束", askContext.Question))
	referResult := uniqueDoc(searchResult)
	echoRefers, err = s.referDo(0, answer, referResult, askContext.Kdb.Setting.Data())
	if err != nil {
		return
	}
	askContext.AppendProcess("finish", "回答完成")
	// 保存记录
	go func(entity *SearchService) {
		_ = entity.askRecordUpdate(askContext.DbData, []string{askContext.Question}, answer, echoRefers)
	}(s.CopyWithCtx(s.GetCtx()).(*SearchService))
	searchResult = referResult
	return
}

// 根据文档去重
func uniqueDoc(result []dto_search.CommonSearchOutput) []dto_search.CommonSearchOutput {
	sort.Slice(result, func(i, j int) bool {
		return result[i].DocSegmentId < result[j].DocSegmentId
	})
	existMap := make(map[int64]dto_search.CommonSearchOutput)
	for _, item := range result {
		if existing, ok := existMap[item.DocId]; ok {
			// 更新内容
			existing.Content = existing.Content + "\n\n" + item.Content
			existMap[item.DocId] = existing
		} else {
			existMap[item.DocId] = item
		}
	}
	// 转换结果为切片
	var output []dto_search.CommonSearchOutput
	for _, item := range existMap {
		output = append(output, item)
	}
	sort.Slice(output, func(i, j int) bool {
		return output[i].DocId < output[j].DocId
	})
	return output
}

func (s *SearchService) WebSearch(req *dto_search.WebSearchReq) (res interface{}, err error) {
	searchResult, err := s.searchData.DocSearch(req.SessionId, req.Question, nil)
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

func (s *SearchService) SessionSearch(req *dto_search.SessionSearchReq) (res *dto_search.SessionSearchRes, err error) {
	askInfo, err := s.askInfoDao.GetBySessionId(req.SessionId)
	if err != nil {
		return nil, err
	}
	if askInfo == nil {
		return nil, components.ErrorAskSessionNoExist
	}
	kdb := &models.Kdb{}
	if askInfo.KdbId != 0 {
		kdb, err = s.kdbData.CheckKdbAuth(askInfo.KdbId, askInfo.UserId, models.AuthTypeRead)
		if err != nil {
			return nil, err
		}
		if kdb.DataType == models.DataSourceDatabase {
			return nil, components.ErrorDbSearchError
		}
	}
	searchResult, err := s.searchData.DocSearch(req.SessionId, req.Question, new(data.DocSearchOptions).WithIndex(kdb.GetIndexName()))
	if err != nil {
		return nil, err
	}
	err = s.referenceUpdate(req.SessionId, searchResult)
	res = &dto_search.SessionSearchRes{
		SearchResult: searchResult,
	}
	return res, nil
}

func (s *SearchService) KdbSearch(req *dto_search.KdbSearchReq) (res *dto_search.KdbSearchRes, err error) {
	userInfo, _ := utils.LoginInfo(s.GetCtx())
	kdb, err := s.kdbData.CheckKdbAuthByName(req.KdbName, userInfo, models.AuthTypeRead)
	if err != nil {
		return nil, err
	}
	res = &dto_search.KdbSearchRes{
		SearchResult: make([]dto_search.CommonSearchOutput, 0),
	}
	// es搜索的片段
	esSearchResult, err := s.searchData.DocSearch("", req.Question, new(data.DocSearchOptions).WithIndex(kdb.GetIndexName()))
	if err != nil {
		return nil, components.ErrorQueryError
	}
	res.SearchResult = append(res.SearchResult, esSearchResult...)
	return

}

func (s *SearchService) QuestionFocus(req *dto_search.QuestionFocusReq) (res *dto_search.QuestionFocusRes, err error) {
	userInfo, _ := utils.LoginInfo(s.GetCtx())
	kdb, err := s.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeRead)
	if err != nil {
		return nil, err
	}
	res = &dto_search.QuestionFocusRes{
		Focus: make([]string, 0),
	}
	// es搜索的片段
	esSearchResult, err := s.searchData.DocSearch("", req.Question, new(data.DocSearchOptions).WithIndex(kdb.GetIndexName()))
	if err != nil {
		return nil, components.ErrorQueryError
	}
	outlines, err := s.jobdApi.QuestionOutline(req.Question, esSearchResult)
	if err != nil {
		return nil, components.ErrorJobdError
	}
	for _, o := range outlines.Outline {
		if o.Level == "h1" {
			res.Focus = append(res.Focus, o.Content)
		}
	}
	return
}

func (s *SearchService) ReportAsk(req *dto_search.ReportAskReq) (res *dto_search.ReportAskRes, err error) {
	userInfo, _ := utils.LoginInfo(s.GetCtx())
	kdb, err := s.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeRead)
	if err != nil {
		return nil, err
	}
	askInfo, err := s.searchData.CreateSession(userInfo.UserId)
	if err != nil {
		return nil, err
	}
	_ = s.askInfoDao.UpdateById(askInfo.Id, map[string]interface{}{"question": req.Question, "ask_type": "professional_no_more_qa", "kdb_id": req.KdbId})
	user, _ := s.userDao.GetByUserId(userInfo.UserId)
	config := user.Setting.Data()
	modelType := config.ModelType
	promptI18n := "\n 输出使用中文！！"
	if config.Language == "en-us" {
		promptI18n = "\n 输出使用英文！！"
	}
	askContext := &AskContext{
		ModelType:  modelType,
		PromptI18n: promptI18n,
		SessionId:  askInfo.SessionId,
		Question:   req.Question,
		Kdb:        kdb,
		DbData:     askInfo,
		UserId:     userInfo.UserId,
	}
	askContext.AnswerStyle = "professional_no_more_qa"
	if len(req.Focus) > 0 {
		outline := make([]jobd.Outline, 0, len(req.Focus))
		outline = append(outline, jobd.Outline{
			Level:      "h1",
			TitleLevel: "#",
			Content:    req.Subject,
		})
		for _, f := range req.Focus {
			outline = append(outline, jobd.Outline{
				Level:      "h2",
				TitleLevel: "##",
				Content:    f,
			})
		}
		askContext.Outline = outline
	}
	answer, answerRefer, searchResult, err := s.AskSyncDo(askContext)
	if err != nil {
		s.LogErrorf("问答报错, %s", err.Error())
		s.askInfoDao.UpdateById(askInfo.Id, map[string]interface{}{"status": models.AskInfoStatusFail})
		return
	}
	res = &dto_search.ReportAskRes{
		Answer:       answer,
		AnswerRefer:  answerRefer,
		SearchResult: searchResult,
	}
	return
}

func (s *SearchService) ReportDocx(req *dto_search.ReportDocxReq) (res *dto_search.ReportDocxRes, err error) {
	goUnoApi := flow.Create(s.GetCtx(), new(api.GoUnoApi))

	html, _ := annotateHTML(req.OriginAnswer, req.Answer, req.AnswerRefer, req.SearchResult)

	file, err := goUnoApi.HtmlToDocx(fmt.Sprintf("%s.html", req.DocName), html)
	if err != nil {
		return nil, err
	}
	res = &dto_search.ReportDocxRes{
		DocxUrl: file,
	}
	return
}

func (s *SearchService) KdbDatabaseSearch(req *dto_search.KdbDatabaseSearchReq) (res interface{}, err error) {
	askInfo, err := s.askInfoDao.GetBySessionId(req.SessionId)
	if err != nil {
		return nil, err
	}
	if askInfo == nil {
		return nil, components.ErrorAskSessionNoExist
	}
	kdb, err := s.kdbData.CheckKdbAuth(askInfo.KdbId, askInfo.UserId, models.AuthTypeRead)
	if err != nil {
		return nil, err
	}
	if kdb.DataType != models.DataSourceDatabase {
		return nil, components.ErrorDbSearchError
	}
	searchResult, err := s.searchData.DatabaseSearch(req.SessionId, req.Question, req.DatabaseType, req.TermsParam)
	if err != nil {
		return nil, err
	}
	res = &dto_search.DatabaseSearchRes{
		SearchResult: searchResult,
	}
	return
}

// markdownToHTML renders Markdown using blackfriday/v2
func markdownToHTML(answer string) string {
	htmlBytes := blackfriday.Run([]byte(answer))
	return string(htmlBytes)
}

// annotateHTML adds references to the rendered Markdown HTML
func annotateHTML(answer string, pAnswer string, refers []dto_search.DoReferItem, searchResult []dto_search.CommonSearchOutput) (string, map[string]string) {
	answerRunes := []rune(answer)
	referHtmlMap := make(map[string]string)
	orderS := []string{}
	for _, refer := range refers {
		// Add referenced text with annotations
		referencedText := string(answerRunes[refer.Start:refer.End])
		var tmp bytes.Buffer
		if header, found := extractMarkdownHeaderFormat(referencedText); found {
			fmt.Printf("提取的 Markdown 标题: %s\n", header)
			referencedText = strings.Replace(referencedText, header, "", 1)
		}
		tmp.WriteString("<span style='color:blue;'>")
		supStr := ""
		// Add all reference details
		for _, ref := range refer.Refers {
			if ref.Index < len(searchResult) {
				showNum := ref.Index + 1
				supStr = supStr + fmt.Sprintf("<sup>[%d]</sup> ", showNum)
			}
		}
		// Close annotation and mark referenced text
		tmp.WriteString(fmt.Sprintf("%s%s</span>", referencedText, supStr))
		referHtmlMap[referencedText] = tmp.String()
		orderS = append(orderS, referencedText)
	}
	for _, v := range orderS {
		answer = strings.Replace(answer, v, referHtmlMap[v], 1)
	}
	if len(pAnswer) > 0 {
		for _, v := range orderS {
			pAnswer = strings.Replace(pAnswer, fmt.Sprintf("[%s]()", v), referHtmlMap[v], 1)
		}
		answer = pAnswer
	}
	html := markdownToHTML(answer)

	return html, referHtmlMap
}

func randShuffle(slice []string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

func resolveAnswerStyle(typeN string, kdbId int64) string {
	switch typeN {
	case "simple":
		return "simplify"
	case "complex":
		return "detailed"
	case "research":
		if kdbId > 0 {
			return "detailed"
		} else {
			return "professional"
		}
	default:
		return "simplify"
	}
}

// 判断文本开头是否包含 Markdown 标题并返回 Markdown 标题格式（#）
func extractMarkdownHeaderFormat(text string) (string, bool) {
	// 正则表达式：匹配文本开头的 Markdown 标题格式（仅提取 # 和空格部分）
	re := regexp.MustCompile(`^#{1,6}\s+`)
	match := re.FindString(text)

	if match != "" {
		// 返回匹配到的标题格式（# 和空格部分）
		return match, true
	}
	return "", false
}
