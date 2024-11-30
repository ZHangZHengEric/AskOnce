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
	SessionId    string         `gorm:"session_id"`                             // sessionId
	SubQuestion  string         `gorm:"sub_question"`                           // 子问题
	SearchResult datatypes.JSON `gorm:"column:search_result"`                   //  搜索结果
	orm.CrudModel
}

type AskSubSearchDao struct {
	flow.Dao
}

func (entity *AskSubSearchDao) OnCreate() {
	entity.SetTable("ask_sub_search")
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
