package models

import (
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
)

var ProcessTypeNameMap = map[string]string{
	"analyze":   "问题分析",
	"webSearch": "全网搜索",
	"vdbSearch": "知识库搜索",
	"summary":   "整理答案",
	"finish":    "回答完成",
}

type AskProcess struct {
	Id        int64  `gorm:"id; primaryKey;autoIncrement" json:"id"` // 自增主键
	SessionId string `gorm:"session_id"`                             // sessionId
	Type      string `gorm:"type"`                                   // 进度类型
	Time      int64  `gorm:"time"`                                   // 时间戳
	Content   string `gorm:"content"`                                // 进度内容
	orm.CrudModel
}

type AskProcessDao struct {
	flow.Dao
}

func (entity *AskProcessDao) OnCreate() {
	entity.SetTable("ask_process")
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
