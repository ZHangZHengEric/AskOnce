package models

import (
	"askonce/components"
	"askonce/components/dto"
	"github.com/pkg/errors"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

const (
	AskInfoStatusFail    = 2
	AskInfoStatusSuccess = 1
)

// AskInfo  提问定义表
type AskInfo struct {
	Id          int64                       `gorm:"id; primaryKey;autoIncrement;comment:自增主键"`
	SessionId   string                      `gorm:"type:varchar(128);default:'';comment:会话id"`
	Question    string                      `gorm:"type:varchar(1024);default:'';comment:问题"`
	AskType     string                      `gorm:"type:varchar(52);default:'';comment:问题类型"`
	KdbId       int64                       `gorm:"type:int(11);default:0;comment:知识库id"`
	SubQuestion datatypes.JSONSlice[string] `gorm:"type:json;sub_question;comment:切分的问题"`
	Answer      datatypes.JSON              `gorm:"type:json;answer;comment:答案"`
	UserId      string                      `gorm:"type:varchar(128);default:'';comment:用户id"`
	Status      int64                       `gorm:"type:int(11);default:0;comment:状态"`
	LikeStatus  int64                       `gorm:"type:int(11);default:0;comment:喜欢状态"`
	orm.CrudModel
}

func (k AskInfo) TableName() string {
	return "ask_info"
}

type AskInfoDao struct {
	flow.Dao
}

func (entity *AskInfoDao) OnCreate() {
	entity.SetTable(AskInfo{}.TableName())
}

func (entity *AskInfoDao) Insert(add *AskInfo) error {
	if add == nil {
		return nil
	}
	return entity.GetDB().Create(add).Error
}

func (entity *AskInfoDao) GetBySessionId(sessionId string) (res *AskInfo, err error) {
	err = entity.GetDB().Where("session_id = ?", sessionId).First(&res).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		return nil, nil
	}
	return
}

func (entity *AskInfoDao) UpdateById(id int64, update map[string]interface{}) error {
	update["updated_at"] = time.Now()
	db := entity.GetDB()
	err := db.Where("id = ?", id).Updates(update).Error
	if err != nil {
		return components.ErrorMysqlError
	}
	return nil
}

func (entity *AskInfoDao) UpdateEntity(update *AskInfo) error {
	db := entity.GetDB()
	err := db.Model(&AskInfo{}).Where("id = ?", update.Id).Updates(update).Error
	if err != nil {
		return components.ErrorMysqlError
	}
	return nil
}

type AskInfoQueryParam struct {
	UserId      int64
	QueryStatus int
	QueryType   string
	Query       string
	dto.PageParam
}

func (entity *AskInfoDao) GetListId(query AskInfoQueryParam) (list []string, cnt int64, err error) {
	db := entity.GetDB().Model(&AskInfo{})
	db = db.Where("user_id = ?", query.UserId)
	if query.QueryStatus > 0 {
		db = db.Where("status = ?", query.QueryStatus)
	}
	if len(query.QueryType) > 0 {
		db = db.Where("ask_type = ?", query.QueryType)
	}
	if len(query.Query) > 0 {
		db = db.Where("question like ?", "%"+query.Query+"%")
	}
	db = db.Count(&cnt)
	if err = db.Offset((query.PageNo-1)*query.PageSize).Limit(query.PageSize).Order("created_at desc").Pluck("session_id", &list).Error; err != nil {
		entity.LogErrorf("AskInfoDao GetList err:%s", err.Error())
		return nil, 0, components.ErrorMysqlError
	}
	return
}

func (entity *AskInfoDao) DeleteByUserIdAndKdbId(userId string, kdbId int64) error {
	db := entity.GetDB().Model(&AskInfo{})
	db = db.Where("user_id = ? and kdb_id = ?", userId, kdbId)
	err := db.Delete(&AskInfo{}).Error
	return err
}

func (entity *AskInfoDao) GetUserLatestKdb(userId string) (kdbIds []int64, err error) {
	err = entity.GetDB().Model(&AskInfo{}).Where("user_id = ? and kdb_id > 0", userId).
		Select("kdb_id").Order("created_at desc").Scan(&kdbIds).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		return nil, nil
	}
	return
}
