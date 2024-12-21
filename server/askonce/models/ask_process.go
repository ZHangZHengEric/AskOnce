package models

import (
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
)

var ProcessTypeNameMap = map[string]string{
	"analyze":   "问题分析",
	"webSearch": "全网搜索",
	"kdbSearch": "知识库搜索",
	"summary":   "整理答案",
	"finish":    "回答完成",
}

type AskProcess struct {
	Id        int64  `gorm:"id; primaryKey;autoIncrement;comment:自增主键"`
	SessionId string `gorm:"type:varchar(128);default:'';comment:会话id"`
	Type      string `gorm:"type:varchar(52);default:'';comment:进度类型"`
	Time      int64  `gorm:"type:int(11);default:0;comment:时间戳"`
	Content   string `gorm:"content"` // 进度内容
	orm.CrudModel
}

func (AskProcess) TableName() string {
	return "ask_process"
}

type AskProcessDao struct {
	flow.Dao
}

func (entity *AskProcessDao) OnCreate() {
	entity.SetTable(AskProcess{}.TableName())
}

func (entity *AskProcessDao) Insert(add *AskProcess) error {
	if add == nil {
		return nil
	}
	return entity.GetDB().Create(add).Error
}

func (entity *AskProcessDao) GetBySessionId(sessionId string) (res []*AskProcess, err error) {
	err = entity.GetDB().Where("session_id = ?", sessionId).Find(&res).Error
	return
}

func (entity *AskProcessDao) BatchInsert(process []*AskProcess) (err error) {
	if process == nil || len(process) == 0 {
		return
	}
	return entity.GetDB().Create(process).Error
}
