package data

import (
	"askonce/components/dto"
	"askonce/models"
	"github.com/xiangtao94/golib/flow"
)

type KdbDocData struct {
	flow.Service
	kdbDocDao *models.KdbDocDao
}

func (k *KdbDocData) OnCreate() {
	k.kdbDocDao = flow.Create(k.GetCtx(), new(models.KdbDocDao))
}

func (k *KdbDocData) GetDocList(kdbId int64, queryName string, pageParam dto.PageParam) (list []*models.KdbDoc, cnt int64, err error) {
	list = make([]*models.KdbDoc, 0)
	list, cnt, err = k.kdbDocDao.GetList(kdbId, queryName, pageParam)
	return
}
