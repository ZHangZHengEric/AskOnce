package data

import (
	"askonce/components"
	"askonce/components/dto"
	"askonce/helpers"
	"askonce/models"
	"encoding/json"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"math/rand"

	"gorm.io/datatypes"

	"time"
)

/*
*
知识库管理
*/
type KdbData struct {
	flow.Service
	kdbDao      *models.KdbDao
	kdbUserDao  *models.KdbUserDao
	kdbShareDao *models.KdbShareDao
	askInfoDao  *models.AskInfoDao
	userDao     *models.UserDao
}

func (k *KdbData) OnCreate() {
	k.kdbDao = flow.Create(k.GetCtx(), new(models.KdbDao))
	k.kdbUserDao = flow.Create(k.GetCtx(), new(models.KdbUserDao))
	k.kdbShareDao = flow.Create(k.GetCtx(), new(models.KdbShareDao))
	k.askInfoDao = flow.Create(k.GetCtx(), new(models.AskInfoDao))
	k.userDao = flow.Create(k.GetCtx(), new(models.UserDao))
}

// 默认的知识库设置
var DefaultKdbSetting = dto.KdbSetting{
	EmbeddingModel: dto.DocEmbeddingModelCommon,
	RetrievalModel: dto.RetrievalSetting{
		SearchMethod:          dto.DocSearchMethodAll,
		TopK:                  10,
		ScoreThresholdEnabled: false,
		ScoreThreshold:        0.3,
		Weights: dto.RetrievalSettingWeights{
			KeywordWeight: 0.5,
			VectorWeight:  0.5,
		},
	},
	KdbAttach: dto.KdbAttach{
		Language:   "zh-cn",
		Cover:      "",
		CoverColor: false,
		Cases:      []string{},
	},
}

// 校验同一个用户下是否有同名知识库
func (k *KdbData) CheckKdbSameName(kdbName string, userId string) (bool, error) {
	kdb, err := k.kdbDao.GetByNameAndUserId(kdbName, userId)
	if err != nil {
		return false, err
	}
	return kdb != nil, nil
}

// 校验知识库、用户是否存在某个权限
func (k *KdbData) CheckKdbAuth(kdbId int64, userId string, authCode int) (*models.Kdb, error) {
	kdb, err := k.kdbDao.GetById(kdbId)
	if err != nil {
		return nil, err
	}
	if kdb == nil {
		return nil, components.ErrorKdbNoOperate
	}
	if kdb.Type == models.KdbTypePublic {
		if authCode == models.AuthTypeRead { // 公开数据集直接访问
			return kdb, nil
		}
		return nil, components.ErrorKdbNoOperate
	}
	r, err := k.kdbUserDao.GetByKdbIdAndUserId(kdb.Id, userId)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, components.ErrorKdbNoOperate
	}
	if r.AuthType < authCode { // 用户操作权限判断
		return nil, components.ErrorKdbNoOperate
	}
	return kdb, nil
}

func (k *KdbData) AddKdb(kdbName, kdbIntro string, user dto.LoginInfo) (add *models.Kdb, err error) {
	now := time.Now()
	add = &models.Kdb{
		Name:       kdbName,
		Intro:      kdbIntro,
		Setting:    datatypes.NewJSONType(DefaultKdbSetting),
		Type:       models.KdbTypePrivate,
		DataSource: models.DataSourceFile,
		Creator:    user.UserName,
		CrudModel: orm.CrudModel{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	db := helpers.MysqlClient.WithContext(k.GetCtx())
	k.kdbDao.SetDB(db)
	k.kdbUserDao.SetDB(db)
	tx := db.Begin()
	err = k.kdbDao.Insert(add)
	if err != nil {
		tx.Rollback()
		return
	}
	err = k.kdbUserDao.Insert(&models.KdbUser{
		KdbId:    add.Id,
		UserId:   user.UserId,
		AuthType: models.AuthTypeSuperAdmin,
		CrudModel: orm.CrudModel{
			CreatedAt: now,
			UpdatedAt: now,
		},
	})
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func (k *KdbData) UpdateKdb(kdb *models.Kdb, kdbName, kdbIntro string, kdbSetting *dto.KdbSetting) (err error) {
	updateMap := map[string]interface{}{}
	updateMap["name"] = kdbName
	updateMap["intro"] = kdbIntro
	if kdbSetting != nil {
		// 先不更新embedding
		kdbSetting.EmbeddingModel = kdb.Setting.Data().EmbeddingModel
		settingStr, _ := json.Marshal(*kdbSetting)
		updateMap["setting"] = settingStr
	}
	err = k.kdbDao.Update(kdb.Id, updateMap)
	if err != nil {
		return
	}
	return
}

// 获取可用的kdb列表
func (k *KdbData) GetKdbList(userId string, query string, param dto.PageParam) (list []*models.Kdb, cnt int64, err error) {
	userRelation, err := k.kdbUserDao.GetByUserId(userId)
	if err != nil {
		return
	}
	kdbIds := make([]int64, 0)
	for _, ur := range userRelation {
		kdbIds = append(kdbIds, ur.KdbId)
	}
	pubIds, err := k.kdbDao.GetPubIds()
	if err != nil {
		return
	}
	kdbIds = append(kdbIds, pubIds...)
	list, cnt, err = k.kdbDao.GetList(kdbIds, query, param)
	if err != nil {
		return
	}
	return
}

func (k *KdbData) DeleteKdb(userId string, kdb *models.Kdb) (err error) {
	db := helpers.MysqlClient.WithContext(k.GetCtx())
	k.kdbDao.SetDB(db)
	k.kdbUserDao.SetDB(db)
	k.askInfoDao.SetDB(db)
	tx := db.Begin()
	err = k.kdbDao.DeleteById(kdb.Id)
	if err != nil {
		tx.Rollback()
		return
	}
	err = k.kdbUserDao.DeleteByKdbIdAndUserId(kdb.Id, userId)
	if err != nil {
		tx.Rollback()
		return
	}
	err = k.askInfoDao.DeleteByUserIdAndKdbId(userId, kdb.Id)
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit().Error
	return
}

func (k *KdbData) DeleteKdbRelation(userId string, kdbId int64) (err error) {
	err = k.kdbUserDao.DeleteById(kdbId)
	if err != nil {
		return
	}
	err = k.askInfoDao.DeleteByUserIdAndKdbId(userId, kdbId)
	if err != nil {
		return
	}
	return
}

func (k *KdbData) GetKdbAuthType(userId string, kdbId int64) (authType int, err error) {
	kdb, err := k.kdbDao.GetById(kdbId)
	if err != nil {
		return 0, err
	}
	if kdb == nil {
		return 0, components.ErrorKdbNoOperate
	}
	if kdb.Type == models.KdbTypePublic {
		authType = models.AuthTypeRead
		return
	}
	r, err := k.kdbUserDao.GetByKdbIdAndUserId(kdb.Id, userId)
	if err != nil {
		return 0, err
	}
	if r == nil {
		return
	}
	return r.AuthType, nil
}

// 查询kdb下用户id
func (k *KdbData) QueryKdbUserRelation(kdbId int64, authType int) (list []*models.KdbUser, err error) {
	userRelations, err := k.kdbUserDao.GetByKdbId(kdbId)
	if err != nil {
		return
	}
	for _, relation := range userRelations {
		if authType != 0 && authType != relation.AuthType {
			continue
		}
		list = append(list, relation)
	}
	return
}

func (k *KdbData) AddKdbUser(kdb *models.Kdb, userIds []string, authType int) (err error) {
	for _, userId := range userIds {
		vv, err := k.kdbUserDao.GetByKdbIdAndUserId(kdb.Id, userId)
		if err != nil {
			return err
		}
		if vv != nil && vv.AuthType == authType {
			continue
		}
		err = k.kdbUserDao.Insert(&models.KdbUser{
			KdbId:    kdb.Id,
			UserId:   userId,
			AuthType: authType,
			CrudModel: orm.CrudModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		})
		if err != nil {
			return err
		}
	}
	return
}

func (k *KdbData) DeleteKdbUser(kdbId int64, userIds []string) (err error) {
	for _, userId := range userIds {
		err = k.kdbUserDao.DeleteByKdbIdAndUserId(kdbId, userId)
		if err != nil {
			return
		}
		err = k.askInfoDao.DeleteByUserIdAndKdbId(userId, kdbId)
		if err != nil {
			return
		}
	}

	return
}

func (k *KdbData) AddKdbShareCode(kdb *models.Kdb, userId string, authType int) (shareCode string, err error) {
	shareCode = randStr(15)
	err = k.kdbShareDao.Insert(&models.KdbShare{
		Id:        0,
		ShareCode: shareCode,
		KdbId:     kdb.Id,
		UserId:    userId,
		AuthType:  authType,
		CrudModel: orm.CrudModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	})
	if err != nil {
		return "", err
	}
	return
}

func (k *KdbData) VerifyKdbShareCode(userId string, shardCode string) (err error) {
	share, err := k.kdbShareDao.GetByShareCode(shardCode)
	if err != nil {
		return
	}
	if share == nil {
		err = components.ErrorShareEmpty
		return
	}
	relation, err := k.kdbUserDao.GetByKdbIdAndUserId(share.KdbId, userId)
	if err != nil {
		return
	}
	if relation != nil {
		return
	}
	err = k.kdbUserDao.Insert(&models.KdbUser{
		KdbId:    share.KdbId,
		UserId:   userId,
		AuthType: share.AuthType,
		CrudModel: orm.CrudModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	})
	return
}

func (k *KdbData) GetKdbShareCode(shareCode string) (kdb *models.Kdb, kdbShare *models.KdbShare, err error) {

	kdbShare, err = k.kdbShareDao.GetByShareCode(shareCode)
	if err != nil {
		return
	}
	if kdbShare == nil {
		err = components.ErrorShareEmpty
		return
	}
	kdb, err = k.kdbDao.GetById(kdbShare.KdbId)
	if err != nil {
		return
	}
	return
}

func randStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}