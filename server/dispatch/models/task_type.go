package models

import (
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
)

type TaskType struct {
	Id           int64  `gorm:"id; primaryKey;autoIncrement;comment:自增主键"`
	TaskType     string `gorm:"type:varchar(52);default:'';comment:任务类型"`
	Instance     string `gorm:"type:varchar(52);default:'';comment:任务实例名称"`
	TaskNumLimit int64  `gorm:"type:int(11);default:0;comment:队列大小，0则无上限"`
	orm.CrudModel
}

func (TaskType) TableName() string {
	return "task_type"
}

type TaskTypeDao struct {
	flow.Dao
}

func (entity *TaskTypeDao) OnCreate() {
	entity.Dao.OnCreate()
	entity.SetTable(TaskType{}.TableName())
}

func (entity *TaskTypeDao) GetAll() (list []*TaskType, err error) {
	err = entity.GetDB().Table(entity.GetTable()).Find(&list).Error
	return
}

func (entity *TaskTypeDao) Insert(add *TaskType) error {
	return entity.GetDB().Table(entity.GetTable()).Create(&add).Error
}

func (entity *TaskTypeDao) GetByTaskType(taskType string) (res *TaskType, err error) {
	err = entity.GetDB().Table(entity.GetTable()).Where("task_type = ?", taskType).First(&res).Error
	return
}
