package service

import (
	"askonce/components/dto/dto_history"
	"askonce/models"
	"askonce/utils"
	"github.com/xiangtao94/golib/flow"
	"time"
)

type HistoryService struct {
	flow.Service
	askInfoDao *models.AskInfoDao
	kdbDao     *models.KdbDao
}

func (entity *HistoryService) OnCreate() {
	entity.askInfoDao = entity.Create(new(models.AskInfoDao)).(*models.AskInfoDao)
	entity.kdbDao = entity.Create(new(models.KdbDao)).(*models.KdbDao)

}

func (entity *HistoryService) Ask(req *dto_history.AskReq) (res *dto_history.AskRes, err error) {
	res = &dto_history.AskRes{
		List: make([]dto_history.AskItem, 0),
	}
	userInfo, _ := utils.LoginInfo(entity.GetCtx())
	if userInfo.UserId == "" {
		return
	}
	asks, cnt, err := entity.askInfoDao.GetListId(models.AskInfoQueryParam{
		UserId:      userInfo.UserId,
		QueryStatus: models.AskInfoStatusSuccess,
		QueryType:   req.QueryType,
		Query:       req.Query,
		PageParam:   req.PageParam,
	})
	if err != nil {
		return nil, err
	}
	res.Total = cnt
	for _, askId := range asks {
		ask, err := entity.askInfoDao.GetBySessionId(askId)
		if err != nil {
			return nil, err
		}
		tmp := dto_history.AskItem{
			SessionId:  ask.SessionId,
			CreateTime: ask.CreatedAt.Format(time.DateTime),
			Question:   ask.Question,
			AskType:    ask.AskType,
		}
		if ask.AskType != "web" {
			kdb, err := entity.kdbDao.GetById(ask.KdbId)
			if err != nil {
				return nil, err
			}
			if kdb != nil {
				tmp.KdbId = kdb.Id
				tmp.KdbName = kdb.Name
			}
		}
		res.List = append(res.List, tmp)
	}
	return
}
