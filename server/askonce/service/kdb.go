package service

import (
	"askonce/components"
	"askonce/components/dto"
	"askonce/components/dto/dto_kdb"
	"askonce/data"
	"askonce/models"
	"askonce/utils"
	"github.com/xiangtao94/golib/flow"
	"time"
)

type KdbService struct {
	flow.Service
	kdbCoverDao *models.KdbCoverDao
	kdbDocDao   *models.KdbDocDao

	kdbData  *data.KdbData
	userData *data.UserData
}

func (k *KdbService) OnCreate() {
	k.kdbCoverDao = flow.Create(k.GetCtx(), new(models.KdbCoverDao))
	k.kdbDocDao = flow.Create(k.GetCtx(), new(models.KdbDocDao))

	k.kdbData = flow.Create(k.GetCtx(), new(data.KdbData))
	k.userData = flow.Create(k.GetCtx(), new(data.UserData))
}

func (k *KdbService) Add(req *dto_kdb.AddReq) (res *dto_kdb.AddRes, err error) {
	userInfo, _ := utils.LoginInfo(k.GetCtx())
	exist, err := k.kdbData.CheckKdbSameName(req.Name, userInfo.Account)
	if err != nil {
		return
	}
	if exist {
		return nil, components.ErrorKdbExist
	}
	add, err := k.kdbData.AddKdb(req.Name, req.Intro, userInfo)
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
		typeNum := 0
		if kdb.Type == models.KdbTypePublic {
			typeNum = 1
		}
		tmpI := dto_kdb.ListItem{
			Id:           kdb.Id,
			Name:         kdb.Name,
			CreateTime:   kdb.CreatedAt.Format(time.DateTime),
			DataSource:   kdb.DataSource,
			DocNum:       0,
			Cover:        "",
			DefaultColor: false,
			Creator:      kdb.Creator,
			Type:         typeNum,
			Intro:        kdb.Intro,
		}
		if kdb.Setting.Data().KdbAttach.CoverId != 0 {
			cover, _ := k.kdbCoverDao.GetById(kdb.Setting.Data().KdbAttach.CoverId)
			if cover != nil {
				tmpI.Cover = cover.Url
			}
		}
		res.List = append(res.List, tmpI)
	}
	return
}

func (k *KdbService) Info(kdbId int64) (res *dto_kdb.InfoRes, err error) {
	userInfo, _ := utils.LoginInfo(k.GetCtx())
	kdb, err := k.kdbData.CheckKdbAuth(kdbId, userInfo.UserId, models.AuthTypeRead)
	if err != nil {
		return
	}
	kdbSetting := kdb.Setting.Data()

	res = &dto_kdb.InfoRes{
		KdbId:          kdb.Id,
		Name:           kdb.Name,
		Intro:          kdb.Intro,
		CreatedAt:      kdb.CreatedAt.Unix(),
		CreatedBy:      kdb.Creator,
		UpdatedAt:      kdb.UpdatedAt.Unix(),
		DataSourceType: dto_kdb.DataSourceType(kdb.DataSource),
		WordCount:      0,
		DocumentCount:  0,
		KdbSetting: dto.KdbSetting{
			RetrievalModel: kdbSetting.RetrievalModel,
		},
	}
	if kdbSetting.KdbAttach.CoverId != 0 {
		cover, _ := k.kdbCoverDao.GetById(kdbSetting.KdbAttach.CoverId)
		if cover != nil {
			res.Cover = cover.Url
		}
	}
	return
}

func (k *KdbService) Delete(kdbId int64) (res interface{}, err error) {
	userInfo, _ := utils.LoginInfo(k.GetCtx())
	kdb, err := k.kdbData.CheckKdbAuth(kdbId, userInfo.UserId, models.AuthTypeSuperAdmin)
	if err != nil {
		return
	}
	err = k.kdbData.DeleteKdb(userInfo.UserId, kdb)
	if err != nil {
		return
	}
	return
}

func (k *KdbService) DeleteRelation(kdbId int64) (res interface{}, err error) {
	userInfo, err := utils.LoginInfo(k.GetCtx())
	if err != nil {
		return nil, err
	}
	kdb, err := k.kdbData.CheckKdbAuth(kdbId, userInfo.UserId, models.AuthTypeSuperAdmin)
	if err != nil {
		return
	}
	err = k.kdbData.DeleteKdbUser(kdb.Id, []string{userInfo.UserId})
	if err != nil {
		return
	}
	return
}

func (k *KdbService) Covers(req *dto.EmptyReq) (res *dto_kdb.CoversRes, err error) {
	userInfo, err := utils.LoginInfo(k.GetCtx())
	if err != nil {
		return nil, err
	}
	res = &dto_kdb.CoversRes{List: make([]dto_kdb.CoverItem, 0)}
	covers, err := k.Create(new(models.KdbCoverDao)).(*models.KdbCoverDao).GetAll()
	if err != nil {
		return nil, err
	}
	for _, cover := range covers {
		if len(cover.UserId) == 0 && userInfo.UserId != cover.UserId {
			continue
		}
		res.List = append(res.List, dto_kdb.CoverItem{
			Id:   cover.Id,
			Type: cover.Type,
			Url:  cover.Url,
		})
	}
	return
}

func (k *KdbService) Auth(kdbId int64) (res *dto_kdb.AuthRes, err error) {
	userInfo, err := utils.LoginInfo(k.GetCtx())
	if err != nil {
		return nil, err
	}
	authType, err := k.kdbData.GetKdbAuthType(userInfo.UserId, kdbId)
	if err != nil {
		return nil, err
	}
	res = &dto_kdb.AuthRes{
		AuthType: authType,
	}
	return
}

func (k *KdbService) UserList(req *dto_kdb.UserListReq) (res *dto_kdb.UserListRes, err error) {
	userInfo, err := utils.LoginInfo(k.GetCtx())
	if err != nil {
		return nil, err
	}
	kdb, err := k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeWrite)
	if err != nil {
		return nil, err
	}
	res = &dto_kdb.UserListRes{
		List:  make([]dto_kdb.UserListItem, 0),
		Total: 0,
	}
	userRelations, err := k.kdbData.QueryKdbUserRelation(kdb.Id, req.AuthType)
	if err != nil {
		return nil, err
	}
	userIds := make([]string, 0)
	joinTimeMap := map[string]string{}
	for _, relation := range userRelations {
		if relation.UserId == userInfo.UserId {
			continue
		}
		userIds = append(userIds, relation.UserId)
		joinTimeMap[relation.UserId] = relation.CreatedAt.Format(time.DateTime)
	}
	if len(userIds) == 0 {
		return
	}
	users, cnt, err := k.userData.QueryUserList(userIds, req.QueryName, req.PageParam)
	if err != nil {
		return nil, err
	}
	res.Total = cnt
	for _, user := range users {
		res.List = append(res.List, dto_kdb.UserListItem{
			UserId:   user.UserId,
			UserName: user.UserName,
			JoinTime: joinTimeMap[user.UserId],
		})
	}
	return
}

func (k *KdbService) UserQuery(req *dto_kdb.UserQueryReq) (res *dto_kdb.UserQueryRes, err error) {
	userInfo, err := utils.LoginInfo(k.GetCtx())
	if err != nil {
		return nil, err
	}
	users, _, err := k.userData.QueryUserList([]string{}, req.QueryName, dto.PageParam{
		PageNo:   1,
		PageSize: 10000,
	})
	if err != nil {
		return nil, err
	}
	res = &dto_kdb.UserQueryRes{
		List: make([]dto_kdb.UserQueryItem, 0),
	}
	for _, user := range users {
		if user.UserId == userInfo.UserId {
			continue
		}

		res.List = append(res.List, dto_kdb.UserQueryItem{
			UserId:   user.UserId,
			UserName: user.UserName,
		})
	}
	res.Total = int64(len(res.List))
	return
}

func (k *KdbService) UserAdd(req *dto_kdb.UserAddReq) (res interface{}, err error) {
	userInfo, err := utils.LoginInfo(k.GetCtx())
	if err != nil {
		return nil, err
	}

	kdb, err := k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeWrite)
	if err != nil {
		return nil, err
	}
	err = k.kdbData.AddKdbUser(kdb, req.UserIds, req.AuthType)
	return
}

func (k *KdbService) UserDelete(req *dto_kdb.UserDeleteReq) (res interface{}, err error) {
	userInfo, err := utils.LoginInfo(k.GetCtx())
	if err != nil {
		return nil, err
	}
	kdb, err := k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeWrite)
	if err != nil {
		return nil, err
	}
	err = k.kdbData.DeleteKdbUser(kdb.Id, req.UserIds)
	return
}

func (k *KdbService) GenShareCode(req *dto_kdb.GenShareCodeReq) (res *dto_kdb.GenShareCodeRes, err error) {
	userInfo, err := utils.LoginInfo(k.GetCtx())
	if err != nil {
		return nil, err
	}
	kdb, err := k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeSuperAdmin)
	if err != nil {
		return nil, err
	}
	shareCode, err := k.kdbData.AddKdbShareCode(kdb, userInfo.UserId, req.AuthType)
	if err != nil {
		return nil, err
	}
	res = &dto_kdb.GenShareCodeRes{
		ShareCode: shareCode,
	}
	return
}

func (k *KdbService) VerifyShareCode(req *dto_kdb.VerifyShareCodeReq) (res interface{}, err error) {
	userInfo, err := utils.LoginInfo(k.GetCtx())
	if err != nil {
		return nil, err
	}
	err = k.kdbData.VerifyKdbShareCode(userInfo.UserId, req.ShareCode)
	if err != nil {
		return nil, components.ErrorShareEmpty
	}
	return
}

func (k *KdbService) ShareCodeInfo(req *dto_kdb.InfoShareCodeReq) (res *dto_kdb.ShareCodeInfoRes, err error) {
	kdb, kdbShare, err := k.kdbData.GetKdbShareCode(req.ShareCode)
	if err != nil {
		return nil, components.ErrorShareEmpty
	}
	res = &dto_kdb.ShareCodeInfoRes{
		Creator:  kdb.Creator,
		KdbName:  kdb.Name,
		AuthType: kdbShare.AuthType,
	}
	return
}
