package service

import (
	"askonce/api/jobd"
	"askonce/components"
	"askonce/components/dto"
	"askonce/components/dto/dto_gpt"
	"askonce/components/dto/dto_kdb_doc"
	"askonce/components/dto/dto_search"
	"askonce/data"
	"askonce/helpers"
	"askonce/models"
	"askonce/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/duke-git/lancet/v2/slice"
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
	askInfoDao *models.AskInfoDao

	kdbData      *data.KdbData
	searchData   *data.SearchData
	jobdApi      *jobd.JobdApi
	askAttachDao *models.AskAttachDao
	processDao   *models.AskProcessDao
	userDao      *models.UserDao
	chatData     *data.GptData
}

func (s *SearchService) OnCreate() {
	s.jobdApi = s.Create(new(jobd.JobdApi)).(*jobd.JobdApi)
	s.askInfoDao = s.Create(new(models.AskInfoDao)).(*models.AskInfoDao)
	s.askAttachDao = s.Create(new(models.AskAttachDao)).(*models.AskAttachDao)
	s.userDao = s.Create(new(models.UserDao)).(*models.UserDao)
	s.processDao = s.Create(new(models.AskProcessDao)).(*models.AskProcessDao)
	s.searchData = s.Create(new(data.SearchData)).(*data.SearchData)
	s.kdbData = s.Create(new(data.KdbData)).(*data.KdbData)
	s.chatData = s.Create(new(data.GptData)).(*data.GptData)
}

func (s *SearchService) EchoRes(stage, text string) {
	echoResStr, _ := sonic.MarshalString(dto_search.AskRes{
		Stage: stage,
		Text:  text,
	})
	time.Sleep(50 * time.Millisecond)
	sse.RenderStream(s.GetCtx(), "0", "message", echoResStr)
}

func (s *SearchService) saveRes(sessionId, stage string, text string) {
	err := s.processDao.Insert(&models.AskProcess{
		SessionId: sessionId,
		Type:      stage,
		Content:   text,
		Time:      time.Now().UnixMilli(),
		CrudModel: orm.CrudModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	})
	if err != nil {
		s.LogErrorf("processDao.Insert error: %s", err.Error())
	}
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
	start, end := utils.SlicePage(req.PageNo, req.PageSize, len(res.List)) //第一页1页显示3条数据
	res.List = res.List[start:end]                                         //分页后的数据
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

func randShuffle(slice []string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
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
	KdbId                 int64  // 为kdb时有值
	DbData                *models.AskInfo
	UserId                string
	AnswerStyle           string
}

func (s *SearchService) Ask(req *dto_search.AskReq) (err error) {
	// 文本校验
	green, _ := helpers.TextCheck(s.GetCtx(), req.Question)
	if !green {
		return components.ErrorTextCheckError
	}
	userInfo, _ := utils.LoginInfo(s.GetCtx())
	// 校验知识库权限
	if req.KdbId > 0 {
		_, err := s.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeRead)
		if err != nil {
			return err
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
	err = askInfoDao.UpdateById(askInfo.Id, map[string]interface{}{"question": req.Question, "ask_type": req.Type, "kdb_id": req.KdbId})
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
	askContext := AskContext{
		ModelType:  modelType,
		PromptI18n: promptI18n,
		SessionId:  req.SessionId,
		Question:   req.Question,
		KdbId:      req.KdbId,
		DbData:     askInfo,
		UserId:     userInfo.UserId,
	}
	s.EchoRes("start", "")
	askDirect := false
	if req.KdbId == 0 { // 互联网判断是否要搜索
		judgeRes, err := s.jobdApi.NetRagAssessment(req.Question)
		if err != nil {
			s.LogErrorf("NetRagAssessment err %s", err.Error())
		} else {
			if judgeRes.Result {
				askDirect = true
			}
		}
	}
	if askDirect {
		err = s.AskDirect(askContext)
	} else {
		switch req.Type {
		case "simple":
			askContext.AnswerStyle = "simplify"
			err = s.AskSimple(askContext)
		case "complex":
			askContext.AnswerStyle = "detailed"
			err = s.AskComplex(askContext)
		case "research":
			if req.KdbId > 0 {
				askContext.AnswerStyle = "detailed"
				err = s.AskComplex(askContext)
			} else {
				askContext.AnswerStyle = "detailed_no_chapter"
				err = s.AskResearch(askContext)
			}
		default:
			return errors.ErrorParamInvalid
		}
	}
	if err != nil {
		s.LogErrorf("问答报错, %s", err.Error())
		askInfoDao.UpdateById(askInfo.Id, map[string]interface{}{"status": models.AskInfoStatusFail})
		return
	}
	s.EchoRes("done", "")
	return
}

func (s *SearchService) AskDirect(req AskContext) (err error) {
	s.saveRes(req.SessionId, "summary", "整理答案开始")
	// 开始回答
	answer, echoRefers, err := s.askChat(req, nil, nil)
	if err != nil {
		return components.ErrorChatError
	}
	s.saveRes(req.SessionId, "summary", "整理答案结束")
	s.saveRes(req.SessionId, "finish", "回答完成")
	// 保存记录
	go func(entity *SearchService) {
		_ = entity.askRecordUpdate(req.DbData, []string{req.Question}, answer, echoRefers)
	}(s.CopyWithCtx(s.GetCtx()).(*SearchService))
	// 生成大纲
	s.askOutline(req.SessionId, answer)
	s.EchoRes("outline", "done")
	return
}

// 简单搜索
func (s *SearchService) AskSimple(req AskContext) (err error) {
	s.EchoRes("search", "")
	if req.KdbId == 0 {
		s.saveRes(req.SessionId, "webSearch", "开始搜索互联网")
	} else {
		s.saveRes(req.SessionId, "vdbSearch", "开始搜索知识库")
	}
	searchResult := make([]dto_search.CommonSearchOutput, 0)
	searchResult, err = s.searchData.SearchFromWebOrKnowledge(req.SessionId, req.Question, req.KdbId, req.UserId)
	if err != nil {
		return components.ErrorQueryError
	}
	if len(searchResult) == 0 {
		return components.ErrorQueryEmpty
	}
	searchResultStr, _ := json.Marshal(searchResult)
	err = s.askAttachDao.UpdateBySessionId(req.SessionId, map[string]interface{}{"reference": searchResultStr})
	if err != nil {
		return components.ErrorMysqlError
	}
	s.EchoRes("search", fmt.Sprintf("搜索到%v条相关内容", len(searchResult)))
	if req.KdbId == 0 {
		s.saveRes(req.SessionId, "webSearch", fmt.Sprintf("互联网搜索到%v条相关内容", len(searchResult)))
	} else {
		s.saveRes(req.SessionId, "vdbSearch", fmt.Sprintf("知识库搜索到%v条相关内容", len(searchResult)))
	}
	if req.AnswerStyle == "detailed" {
		go func(entity *SearchService) {
			first, err := entity.jobdApi.GenerateRelateInfo(req.Question, searchResult)
			if err != nil {
				entity.LogErrorf("GenerateRelateInfo error, %s", err.Error())
			}
			second, err := entity.jobdApi.DeduplicationRelateInfo(first)
			if err != nil {
				entity.LogErrorf("DeduplicationRelateInfo error, %s", err.Error())
			}
			relateStr, _ := json.Marshal(second)
			err = entity.askAttachDao.UpdateBySessionId(req.SessionId, map[string]interface{}{"relation": relateStr})
			if err != nil {
				return
			}
			entity.EchoRes("relate", "done")
		}(s.CopyWithCtx(s.GetCtx()).(*SearchService))
	}
	// 处理prompt
	promptRes, err := s.jobdApi.SimpleQAConstruct(req.Question, req.AnswerStyle, searchResult)
	if err != nil {
		return components.ErrorJobdError
	}
	s.saveRes(req.SessionId, "summary", "整理答案开始")
	// 开始回答
	answer, echoRefers, err := s.askChat(req, promptRes, searchResult)
	if err != nil {
		return components.ErrorChatError
	}
	s.saveRes(req.SessionId, "summary", "整理答案结束")
	s.saveRes(req.SessionId, "finish", "回答完成")
	// 保存记录
	go func(entity *SearchService) {
		_ = entity.askRecordUpdate(req.DbData, []string{req.Question}, answer, echoRefers)
	}(s.CopyWithCtx(s.GetCtx()).(*SearchService))
	// 生成大纲
	s.askOutline(req.SessionId, answer)
	s.EchoRes("outline", "done")
	return
}

func (s *SearchService) AskComplex(req AskContext) (err error) {
	s.EchoRes("analyze", fmt.Sprintf("开始分析问题：%s", req.Question))
	s.saveRes(req.SessionId, "analyze", fmt.Sprintf("开始分析问题：%s", req.Question))
	splitRes, err := s.jobdApi.SplitQuestion(req.Question)
	if err != nil {
		return components.ErrorJobdError
	}
	if len(splitRes.SubTitles) == 0 {
		s.EchoRes("analyze", fmt.Sprintf("分析问题为: %s", req.Question))
		s.saveRes(req.SessionId, "analyze", fmt.Sprintf("分析问题为: %s", req.Question))
	} else {
		s.EchoRes("analyze", fmt.Sprintf("分析问题为: %s ", strings.Join(splitRes.SubTitles, ";")))
		s.saveRes(req.SessionId, "analyze", fmt.Sprintf("分析问题为: %s ", strings.Join(splitRes.SubTitles, ";")))
	}
	if len(splitRes.SearchContents) <= 1 {
		return s.AskSimple(req)
	}
	s.EchoRes("search", "")
	if req.KdbId == 0 {
		s.saveRes(req.SessionId, "webSearch", "开始搜索互联网")
	} else {
		s.saveRes(req.SessionId, "vdbSearch", "开始搜索知识库")
	}
	searchResultAll := make([]dto_search.CommonSearchOutput, 0)
	searchResultAllMap := make(map[string][]dto_search.CommonSearchOutput)
	// 处理拆分问题的单个回答
	eg0, _ := errgroup.WithContext(s.GetCtx())
	lock0 := sync.Mutex{}
	for i, subQ := range splitRes.SubTitles {
		tmpQ := subQ
		tmpSearchContent := splitRes.SearchContents[i]
		eg0.Go(func() (err error) {
			searchResult, err := s.searchData.SearchFromWebOrKnowledge(req.SessionId, tmpSearchContent, req.KdbId, req.UserId)
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
		return components.ErrorQueryError
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
	if req.KdbId == 0 {
		s.saveRes(req.SessionId, "webSearch", fmt.Sprintf("互联网搜索到%v条相关内容", len(searchResultAll)))
	} else {
		s.saveRes(req.SessionId, "vdbSearch", fmt.Sprintf("知识库搜索到%v条相关内容", len(searchResultAll)))
	}
	oreferenceStr, _ := json.Marshal(searchResultAll)
	err = s.askAttachDao.UpdateBySessionId(req.SessionId, map[string]interface{}{"reference": oreferenceStr})
	if err != nil {
		return
	}
	go func(entity *SearchService) {
		first, err := entity.jobdApi.GenerateRelateInfo(req.Question, searchResultAll)
		if err != nil {
			entity.LogErrorf("GenerateRelateInfo error, %s", err.Error())
		}
		second, err := entity.jobdApi.DeduplicationRelateInfo(first)
		if err != nil {
			entity.LogErrorf("DeduplicationRelateInfo error, %s", err.Error())
		}
		relateStr, _ := json.Marshal(second)
		err = entity.askAttachDao.UpdateBySessionId(req.SessionId, map[string]interface{}{"relation": relateStr})
		if err != nil {
			return
		}
		entity.EchoRes("relation", "done")
	}(s.CopyWithCtx(s.GetCtx()).(*SearchService))

	subAnswerAllMap := make(map[int]string)
	eg1, _ := errgroup.WithContext(s.GetCtx())
	lock1 := sync.Mutex{}
	for i, subQ := range splitRes.SubTitles {
		tmpQ := subQ
		tmpSearchResult := searchResultAllMap[tmpQ]
		tmpIndex := i
		eg1.Go(func() (err error) {
			// 处理prompt
			promptRes, err := s.jobdApi.SimpleQAConstruct(tmpQ, req.AnswerStyle, tmpSearchResult)
			if err != nil {
				return
			}

			prompt := tmpQ
			var temperature float64
			var maxNewTokens int
			if len(promptRes.Prompt) > 0 {
				prompt = promptRes.Prompt
				temperature = promptRes.GenerateParam.Temperature
				maxNewTokens = promptRes.GenerateParam.MaxNewTokens
			}
			chatReq := &dto_gpt.ChatCompletionReq{
				Messages: []dto_gpt.ChatCompletionMessage{
					{
						Role:    "Human",
						Content: prompt + req.PromptI18n,
					},
				},
				Temperature: temperature,
				MaxTokens:   maxNewTokens,
			}
			subAnswer, _, err := s.chatData.ChatSync(req.ModelType, req.UserId, chatReq)
			if err != nil {
				return
			}
			lock1.Lock()
			subAnswerAllMap[tmpIndex] = subAnswer
			lock1.Unlock()
			return nil
		})
	}
	if err = eg1.Wait(); err != nil {
		return components.ErrorChatError
	}
	s.saveRes(req.SessionId, "summary", "整理答案开始")
	var subAnswerAll []string
	for i := range splitRes.SubTitles {
		subAnswerAll = append(subAnswerAll, subAnswerAllMap[i])
	}
	promptRes, err := s.jobdApi.MergeAnswers(req.Question, splitRes.SubTitles, subAnswerAll)
	if err != nil {
		return components.ErrorChatError
	}
	// 开始回答
	answer, echoRefers, err := s.askChat(req, promptRes, searchResultAll)
	if err != nil {
		return components.ErrorChatError
	}
	s.saveRes(req.SessionId, "summary", "整理答案结束")
	// 保存记录
	go func(entity *SearchService) {
		_ = entity.askRecordUpdate(req.DbData, splitRes.SubTitles, answer, echoRefers)
	}(s.CopyWithCtx(s.GetCtx()).(*SearchService))
	// 生成大纲
	s.askOutline(req.SessionId, answer)
	s.EchoRes("outline", "done")

	return
}

func (s *SearchService) AskResearch(req AskContext) (err error) {
	s.EchoRes("search", "")
	s.saveRes(req.SessionId, "webSearch", "开始搜索互联网")
	searchResult := make([]dto_search.CommonSearchOutput, 0)
	searchResult, err = s.searchData.SearchFromWebOrKnowledge(req.SessionId, req.Question, req.KdbId, req.UserId)
	if err != nil {
		return components.ErrorQueryError
	}
	if len(searchResult) == 0 {
		return components.ErrorQueryEmpty
	}
	s.EchoRes("search", fmt.Sprintf("搜索到%v条相关内容", len(searchResult)))
	s.saveRes(req.SessionId, "webSearch", fmt.Sprintf("互联网搜索到%v条相关内容", len(searchResult)))
	oreferenceStr, _ := json.Marshal(searchResult)
	err = s.askAttachDao.UpdateBySessionId(req.SessionId, map[string]interface{}{"reference": oreferenceStr})
	if err != nil {
		return
	}
	keyPointsRes, err := s.jobdApi.GenerateAnswerKeyPoints(req.Question, searchResult)
	if err != nil {
		return components.ErrorChatError
	}
	keyPoints := keyPointsRes.AnswerKeyPoints
	if len(keyPoints) == 0 {
		return components.ErrorQueryEmpty
	}
	moreKeyPoints := make([]string, 0)
	wg3 := sync.WaitGroup{}
	wg3.Add(1)
	go func(entity *SearchService) {
		defer func() {
			wg3.Done()
		}()
		moreKeyPointRes, err := entity.jobdApi.GenerateMoreKeyPoints(req.Question, keyPointsRes.AnswerKeyPoints)
		if err != nil {
			entity.LogErrorf("GenerateMoreKeyPoints err %v", err.Error())
			return
		}
		if len(moreKeyPointRes.MoreKeyPoints) > 0 {
			moreKeyPoints = append(moreKeyPoints, moreKeyPointRes.MoreKeyPoints...)
		}
	}(s.CopyWithCtx(s.GetCtx()).(*SearchService))
	// 并发处理相关信息
	go func(entity *SearchService) {
		first, err := entity.jobdApi.GenerateRelateInfo(req.Question, searchResult)
		if err != nil {
			entity.LogErrorf("GenerateRelateInfo error, %s", err.Error())
		}
		second, err := entity.jobdApi.DeduplicationRelateInfo(first)
		if err != nil {
			entity.LogErrorf("DeduplicationRelateInfo error, %s", err.Error())
		}
		relateStr, _ := json.Marshal(second)
		err = entity.askAttachDao.UpdateBySessionId(req.SessionId, map[string]interface{}{"relation": relateStr})
		if err != nil {
			return
		}
		entity.EchoRes("relate", "done")
	}(s.CopyWithCtx(s.GetCtx()).(*SearchService))
	// 按顺序输出
	s.saveRes(req.SessionId, "summary", "整理答案开始")
	s.EchoRes("generate", "")
	answerAll := ""
	echoReferAll := []DoReferItem{}
	beginIndex := 0
	kpointTree, allPath := jobd.ReturnTree(keyPoints)

	answerAll, echoReferAll, beginIndex, err = s.treeAskForResearch(kpointTree, answerAll, echoReferAll, beginIndex, req, searchResult, keyPoints)
	if err != nil {
		return components.ErrorChatError
	}
	wg3.Wait()
	morePointDone := false
	searchResultUnique := make(map[string]bool)
	for _, output := range searchResult {
		unique := base64.StdEncoding.EncodeToString([]byte(output.Url + output.Content))
		if searchResultUnique[unique] {
			continue
		}
		searchResultUnique[unique] = true
	}
	for _, point := range moreKeyPoints {
		point = strings.Replace(point, "#", "", -1)
		point = strings.TrimSpace(point)
		srTmp, err := s.searchData.SearchFromWebOrKnowledge(req.SessionId, req.Question, req.KdbId, req.UserId)
		if err != nil {
			continue
		}
		for _, output := range srTmp {
			unique := base64.StdEncoding.EncodeToString([]byte(output.Url + output.Content))
			if searchResultUnique[unique] {
				continue
			}
			searchResultUnique[unique] = true
			searchResult = append(searchResult, output)
		}
		orferenceStr2, _ := json.Marshal(searchResult)
		err = s.askAttachDao.UpdateBySessionId(req.SessionId, map[string]interface{}{"reference": orferenceStr2})
		if err != nil {
			continue
		}
		s.EchoRes("refreshSearch", fmt.Sprintf("再次搜索到%v条相关内容", len(searchResult)))

		// 处理prompt
		promptRes, err := s.jobdApi.SimpleQAConstruct(point, "simplify", searchResult)
		if err != nil {
			continue
		}
		zhanwei := ""
		if morePointDone {
			zhanwei = fmt.Sprintf("\n\n ## %s\n\n", point)
		} else {
			zhanwei = fmt.Sprintf("\n\n# %s\n\n ## %s\n\n", "更多问题", point)
			morePointDone = true
		}
		s.EchoRes("appendText", zhanwei)
		zhanweil := len([]rune(zhanwei))
		// 开始回答
		beginIndex = beginIndex + zhanweil
		answer, echoRefers, err := s.askChatForResearch(req, promptRes, searchResult, beginIndex, echoReferAll)
		if err != nil {
			return components.ErrorChatError
		}
		answerAll = answerAll + zhanwei + answer
		beginIndex = beginIndex + len([]rune(answer))
		echoReferAll = echoRefers

	}
	s.EchoRes("complete", answerAll)
	s.saveRes(req.SessionId, "summary", "整理答案结束")
	// 保存记录
	go func(entity *SearchService) {
		_ = entity.askRecordUpdate(req.DbData, allPath, answerAll, echoReferAll)
	}(s.CopyWithCtx(s.GetCtx()).(*SearchService))
	// 生成大纲
	s.askOutline(req.SessionId, answerAll)
	s.EchoRes("outline", "done")

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

type DoReferItem struct {
	Start       int                `json:"start"`
	End         int                `json:"end"`
	NumberIndex int                `json:"numberIndex"`
	Refers      []DoReferReferItem `json:"refers"`
}

type DoReferReferItem struct {
	Index      int `json:"index"`
	ReferStart int `json:"referStart"`
	ReferEnd   int `json:"referEnd"`
}

func (s *SearchService) referDo(begin int, needReference string, searchResult []dto_search.CommonSearchOutput) (output []DoReferItem, err error) {

	referStrList := []string{}
	for _, o := range searchResult {
		referStrList = append(referStrList, o.Content)
	}

	referenceRes, err := s.jobdApi.AtomResultReference(needReference, referStrList)
	if err != nil {
		zlog.Errorf(s.GetCtx(), "AtomResultReference error %s", err.Error())
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
			t := DoReferItem{
				Start:       begin + referenceMap.IndexRange[0],
				End:         begin + referenceMap.IndexRange[1],
				NumberIndex: begin + length,
				Refers:      nil,
			}
			refers := []DoReferReferItem{}
			if len(referenceMap.ReferenceList) == 0 {
				continue
			}
			for _, index := range referenceMap.ReferenceList {
				index2, ok := referenceRes.ReferenceListSelectIndex[strconv.Itoa(index)]
				if !ok { // 文档引用来源 字段中不存在 跳过
					continue
				}
				if len(index2) == 2 {
					refers = append(refers, DoReferReferItem{
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
		List: make([]dto_search.ReferenceItem, 0),
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
	refers := make([]dto_search.ReferenceItem, 0)
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

func (s *SearchService) askChat(req AskContext, promptRes *jobd.SimpleQAConstructRes, searchResult []dto_search.CommonSearchOutput) (answer string, echoRefers []DoReferItem, err error) {
	prompt := req.Question
	var temperature float64
	var maxNewTokens int
	if promptRes != nil && len(promptRes.Prompt) > 0 {
		prompt = promptRes.Prompt
		temperature = promptRes.GenerateParam.Temperature
		maxNewTokens = promptRes.GenerateParam.MaxNewTokens
	}
	chatReq := &dto_gpt.ChatCompletionReq{
		Messages: []dto_gpt.ChatCompletionMessage{
			{
				Role:    "Human",
				Content: prompt + req.PromptI18n,
			},
		},
		Temperature: temperature,
		MaxTokens:   maxNewTokens,
	}
	// 生成答案 + 引用
	alreadyReferAnswer := ""
	wg := sync.WaitGroup{}
	lock := sync.Mutex{}
	first := true
	err = s.chatData.Chat(req.ModelType, req.UserId, chatReq, func(offset int, chatAnswer data.ChatAnswer) error {
		if first {
			s.EchoRes("generate", "")
			first = false
		}
		currentAnswer := chatAnswer.Answer
		// 对话展示逻辑
		echoAnswer := strings.Replace(currentAnswer, answer, "", 1)
		if len([]rune(echoAnswer)) <= 10 && chatAnswer.Status != "FINISH" {
			return nil
		}
		s.EchoRes("appendText", echoAnswer)
		answer = chatAnswer.Answer
		if len(searchResult) > 0 {
			// 引用判断逻辑
			needReference, begin := IsCompleted(currentAnswer, chatAnswer.Status, alreadyReferAnswer)
			if len(needReference) > 0 {
				s.LogInfof("完整句子: %s。开始位置: %v", needReference, begin)
				wg.Add(1)
				alreadyReferAnswer = alreadyReferAnswer + needReference
				go func(entity *SearchService, begin int, needRefer string, searchResult []dto_search.CommonSearchOutput) {
					defer wg.Done()
					aa, errA := entity.referDo(begin, needRefer, searchResult)
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
	s.EchoRes("complete", answer)

	return
}

func (s *SearchService) askChatForResearch(req AskContext, promptRes *jobd.SimpleQAConstructRes, searchResult []dto_search.CommonSearchOutput, startIndex int, echoReferAll []DoReferItem) (answer string, echoRefers []DoReferItem, err error) {
	echoRefers = append(echoRefers, echoReferAll...)
	prompt := req.Question
	var temperature float64
	var maxNewTokens int
	if promptRes != nil && len(promptRes.Prompt) > 0 {
		prompt = promptRes.Prompt
		temperature = promptRes.GenerateParam.Temperature
		maxNewTokens = promptRes.GenerateParam.MaxNewTokens
	}
	chatReq := &dto_gpt.ChatCompletionReq{
		Messages: []dto_gpt.ChatCompletionMessage{
			{
				Role:    "Human",
				Content: prompt + req.PromptI18n,
			},
		},
		Temperature: temperature,
		MaxTokens:   maxNewTokens,
	}
	// 生成答案 + 引用
	alreadyReferAnswer := ""
	wg := sync.WaitGroup{}
	lock := sync.Mutex{}
	err = s.chatData.Chat(req.ModelType, req.UserId, chatReq, func(offset int, chatAnswer data.ChatAnswer) error {
		currentAnswer := chatAnswer.Answer
		// 对话展示逻辑
		echoAnswer := strings.Replace(currentAnswer, answer, "", 1)
		if len([]rune(echoAnswer)) <= 10 && chatAnswer.Status != "FINISH" {
			return nil
		}
		s.EchoRes("appendText", echoAnswer)
		answer = chatAnswer.Answer
		if len(searchResult) > 0 {
			// 引用判断逻辑
			needReference, begin := IsCompleted(currentAnswer, chatAnswer.Status, alreadyReferAnswer)
			begin = begin + startIndex
			if len(needReference) > 0 {
				s.LogInfof("完整句子: %s。开始位置: %v", needReference, begin)
				wg.Add(1)
				alreadyReferAnswer = alreadyReferAnswer + needReference
				go func(entity *SearchService, begin int, needRefer string, searchResult []dto_search.CommonSearchOutput) {
					defer wg.Done()
					aa, errA := entity.referDo(begin, needRefer, searchResult)
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
}

func (s *SearchService) askRecordUpdate(askInfo *models.AskInfo, questions []string, answer string, echoRefers []DoReferItem) (err error) {
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
	var mapResult = make(map[string]string)
	_ = json.Unmarshal(askInfo.Answer, &mapResult)
	askAttach, err := s.askAttachDao.GetBySessionId(req.SessionId)
	if err != nil {
		return
	}
	if askAttach == nil {
		return
	}
	relationObj := &jobd.GenerateRelateInfoRes{}
	_ = json.Unmarshal(askAttach.Relation, &relationObj)
	for _, info := range relationObj.PeopleInfo {
		res.PeopleInfo = append(res.PeopleInfo, dto_search.RelatePeopleInfo{
			PersonName:     info.PersonName,
			PersonDescribe: info.PersonDescribe,
		})
	}
	for _, info := range relationObj.EventsInfo {
		res.EventsInfo = append(res.EventsInfo, dto_search.RelateEventsInfo{
			EventName:     info.EventName,
			EventDate:     info.EventDate,
			EventDescribe: info.EventDescribe,
		})
	}
	for _, info := range relationObj.OrgsInfo {
		res.OrgsInfo = append(res.OrgsInfo, dto_search.RelateOrgInfo{
			OrgName:     info.OrgName,
			OrgDescribe: info.OrgDescribe,
		})
	}
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

func (s *SearchService) treeAskForResearch(tree []*jobd.KeyPointNode, answerAll string, echoReferAll []DoReferItem, beginIndex int, req AskContext, searchResult []dto_search.CommonSearchOutput, keyPoints []jobd.AnswerKeyPointsItem) (string, []DoReferItem, int, error) {
	if len(tree) == 0 {
		return answerAll, echoReferAll, beginIndex, nil
	}
	var err error
	for _, node := range tree {
		answerAll, echoReferAll, beginIndex, err = s.treeAskForResearch(node.Children, answerAll, echoReferAll, beginIndex, req, searchResult, keyPoints)
		// 处理prompt
		promptRes, err := s.jobdApi.ResearchQAConstruct(node.FullPath, req.AnswerStyle, searchResult, node.Content, keyPoints)
		if err != nil {
			return answerAll, echoReferAll, beginIndex, components.ErrorJobdError
		}
		zhanwei := ""
		if node.Level == "h1" {
			zhanwei = fmt.Sprintf("\n\n # %s\n\n", node.Content)
		} else if node.Level == "h2" {
			zhanwei = fmt.Sprintf("\n\n ## %s\n\n", node.Content)
		} else if node.Level == "h3" {
			zhanwei = fmt.Sprintf("\n\n ### %s\n\n", node.Content)
		}
		s.EchoRes("appendText", zhanwei)
		zhanweil := len([]rune(zhanwei))
		// 开始回答
		beginIndex = beginIndex + zhanweil
		answer, echoRefers, err := s.askChatForResearch(req, promptRes, searchResult, beginIndex, echoReferAll)
		if err != nil {
			return answerAll, echoReferAll, beginIndex, components.ErrorChatError
		}
		answerAll = answerAll + zhanwei + answer
		beginIndex = beginIndex + len([]rune(answer))
		echoReferAll = append(echoReferAll, echoRefers...)
	}
	return answerAll, echoReferAll, beginIndex, err
}

func (s *SearchService) Recall(req *dto_kdb_doc.RecallReq) (res *dto_kdb_doc.RecallRes, err error) {
	userInfo, err := utils.LoginInfo(s.GetCtx())
	if err != nil {
		return nil, err
	}
	kdb, err := s.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeRead)
	if err != nil {
		return nil, err
	}
	res = &dto_kdb_doc.RecallRes{
		List: make([]dto_kdb_doc.RecallItem, 0),
	}

	// es搜索的片段
	esSearchResult, err := s.searchData.CommonEsSearch(data.EsCommonSearch{
		IndexName: kdb.GetIndexName(),
		Query:     req.Query,
	})
	if err != nil {
		return nil, components.ErrorQueryEmpty
	}
	for _, result := range esSearchResult {
		res.List = append(res.List, dto_kdb_doc.RecallItem{
			DataName:      result.Title,
			DataPath:      result.Url,
			SearchContent: result.Content,
			DataContent:   result.FullContent,
		})
	}
	return
}

func mergeItems(items []DoReferItem) []DoReferItem {
	if len(items) == 0 {
		return nil
	}

	mergedItems := make([]DoReferItem, 0)
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
