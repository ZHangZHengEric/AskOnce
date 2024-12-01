package models

import (
	"github.com/pkg/errors"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AskSubSearch struct {
	Id           int64          `gorm:"id; primaryKey;autoIncrement" json:"id"` //  自增主键
	SessionId    string         `gorm:"type:varchar(128);default:'';comment:会话id"`
	SubQuestion  string         `gorm:"type:varchar(1024);default:'';comment:子问题"` // 子问题
	SearchResult datatypes.JSON `gorm:"type:json;column:search_result"`            //  搜索结果
	orm.CrudModel
}

func (AskSubSearch) TableName() string {
	return "ask_sub_search"
}

type AskSubSearchDao struct {
	flow.Dao
}

func (entity *AskSubSearchDao) OnCreate() {
	entity.SetTable(AskSubSearch{}.TableName())
}

func (entity *AskSubSearchDao) Insert(add *AskSubSearch) error {
	if add == nil {
		return nil
	}
	return entity.GetDB().Create(add).Error
}

func (entity *AskSubSearchDao) GetBySessionId(sessionId string) (res *AskSubSearch, err error) {
	err = entity.GetDB().Where("session_id = ?", sessionId).First(&res).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		return nil, nil
	}
	return
}
