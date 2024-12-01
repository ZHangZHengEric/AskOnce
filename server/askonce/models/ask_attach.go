package models

import (
	"askonce/components"
	"github.com/pkg/errors"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type AskAttach struct {
	Id        int64          `gorm:"id; primaryKey;autoIncrement;comment:自增主键"`
	SessionId string         `gorm:"type:varchar(128);default:'';comment:会话id"`
	Reference datatypes.JSON `gorm:"type:json;column:reference"` //  参考引用
	Outline   datatypes.JSON `gorm:"type:json;column:outline"`   //  大纲
	Relation  datatypes.JSON `gorm:"type:json;column:relation"`  //  相关
	orm.CrudModel
}

func (AskAttach) TableName() string {
	return "ask_attach"
}

type AskAttachDao struct {
	flow.Dao
}

func (entity *AskAttachDao) OnCreate() {
	entity.SetTable(AskAttach{}.TableName())
}

func (entity *AskAttachDao) Insert(add *AskAttach) error {
	if add == nil {
		return nil
	}
	return entity.GetDB().Create(add).Error
}

func (entity *AskAttachDao) GetBySessionId(sessionId string) (res *AskAttach, err error) {
	err = entity.GetDB().Where("session_id = ?", sessionId).First(&res).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		return nil, nil
	}
	return
}

func (entity *AskAttachDao) UpdateBySessionId(sessionId string, update map[string]interface{}) error {
	update["updated_at"] = time.Now()
	db := entity.GetDB()
	err := db.Where("session_id = ?", sessionId).Updates(update).Error
	if err != nil {
		return components.ErrorMysqlError
	}
	return nil
}
