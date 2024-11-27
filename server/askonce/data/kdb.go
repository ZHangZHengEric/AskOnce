package data

import (
	"askonce/components"
	"askonce/components/dto"
	"askonce/helpers"
	"askonce/models"
	"encoding/json"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"

	"gorm.io/datatypes"

	"time"
)

type KdbData struct {
	flow.Service
	kdbDao     *models.KdbDao
	kdbUserDao *models.KdbUserDao
}

func (k *KdbData) OnCreate() {
	k.kdbDao = flow.Create(k.GetCtx(), new(models.KdbDao))
	k.kdbUserDao = flow.Create(k.GetCtx(), new(models.KdbUserDao))
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

func (k *KdbData) AddKdb(kdbName, kdbIntro string, userId string) (add *models.Kdb, err error) {
	now := time.Now()
	add = &models.Kdb{
		Name:       kdbName,
		Intro:      kdbIntro,
		Setting:    datatypes.NewJSONType(DefaultKdbSetting),
		Type:       models.KdbTypePrivate,
		DataSource: models.DataSourceFile,
		Creator:    userId,
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
		UserId:   userId,
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
