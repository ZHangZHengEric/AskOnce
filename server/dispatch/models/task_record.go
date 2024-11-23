package models

import (
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"hash/crc32"
)

type TaskRecord struct {
	Id        int64  `gorm:"id; primaryKey;autoIncrement;comment:自增主键"`
	TaskType  string `gorm:"type:varchar(52);default:'';comment:任务类型"`
	Instance  string `gorm:"type:varchar(52);default:'';comment:任务实例名称"`
	TaskId    string `gorm:"type:varchar(52);default:'';comment:任务Id"`
	SessionId string `gorm:"type:varchar(512);default:'';comment:任务sessionId"`
	Status    string `gorm:"type:varchar(52);default:'';comment:状态"`
	orm.CrudModel
}

type TaskRecordDao struct {
	flow.Dao
}

func (entity *TaskRecordDao) OnCreate() {
	entity.SetTable("task_record")
}

func (entity *TaskRecordDao) Insert(add *TaskRecord) error {
	return entity.GetDB().Table(entity.GetPartitionTable(int64(crc32.ChecksumIEEE([]byte(add.TaskType))))).Create(&add).Error
}
