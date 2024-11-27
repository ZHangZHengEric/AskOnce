package service

import (
	"askonce/components"
	"askonce/components/dto"
	"askonce/data"
	"askonce/models"
	"askonce/utils"
	"github.com/xiangtao94/golib/flow"
	"time"
)

type KdbService struct {
	flow.Service
	kdbData *data.KdbData
}

func (k *KdbService) OnCreate() {
	k.kdbData = flow.Create(k.GetCtx(), new(data.KdbData))
}

func (k *KdbService) Add(req *dto_kdb.AddReq) (res *dto_kdb.AddRes, err error) {
	userInfo, _ := utils.LoginInfo(k.GetCtx())
	exist, err := k.kdbData.CheckKdbSameName(req.Name, userInfo.UserId)
	if err != nil {
		return
	}
	if exist {
		return nil, components.ErrorKdbExist
	}
	add, err := k.kdbData.AddKdb(req.Name, req.Intro, userInfo.UserId)
	if err != nil {
		return
	}
	res = &dto_kdb.AddRes{
		KdbId: add.Id,
	}
	return
}

func (k *KdbService) Update(req *dto_kdb.UpdateReq) (res any, err error) {
	userInfo, _ := utils.LoginInfo(k.GetCtx())
	kdb, err := k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeSuperAdmin)
	if err != nil {
		return
	}
	err = k.kdbData.UpdateKdb(kdb, req.Name, req.Intro, req.KdbSetting)
	if err != nil {
		return
	}
	return
}

func (k *KdbService) List(req *dto_kdb.ListReq) (res *dto_kdb.ListResp, err error) {
	userInfo, _ := utils.LoginInfo(k.GetCtx())
	kdbs, cnt, err := k.kdbData.GetKdbList(userInfo.UserId, req.QueryName, req.PageParam)
	if err != nil {
		return nil, err
	}
	res = &dto_kdb.ListResp{
		Total: cnt,
		List:  make([]dto_kdb.ListItem, 0),
	}
	for _, kdb := range kdbs {
		res.List = append(res.List, dto_kdb.ListItem{
			Id:         kdb.Id,
			Name:       kdb.Name,
			DataSource: kdb.DataSource,
			CreateTime: kdb.CreatedAt.Format(time.DateTime),
			Type:       kdb.Type,
		})
	}
	return
}

func (k *KdbService) Info(req *dto_kdb.InfoReq) (res *dto_kdb.InfoRes, err error) {
	userInfo, _ := utils.LoginInfo(k.GetCtx())
	kdb, err := k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeRead)
	if err != nil {
		return
	}
	kdbSetting := kdb.Setting.Data()
	res = &dto_kdb.InfoRes{
		Id:             kdb.Id,
		Name:           kdb.Name,
		Description:    kdb.Intro,
		CreatedAt:      kdb.CreatedAt.Unix(),
		CreatedBy:      kdb.Creator,
		UpdatedAt:      kdb.UpdatedAt.Unix(),
		DataSourceType: dto_kdb.DataSourceType(kdb.DataSource),
		WordCount:      0,
		DocumentCount:  0,
		KdbSetting: dto.KdbSetting{
			EmbeddingModel: kdbSetting.EmbeddingModel,
			RetrievalModel: kdbSetting.RetrievalModel,
		},
	}
	return
}
