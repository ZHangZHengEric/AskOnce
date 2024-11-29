package service

import (
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"jobd/components/dto"
	"jobd/data"
	"jobd/models"
	"time"
)

type TaskManageService struct {
	flow.Service
	taskData *data.TaskCache
}

func (entity *TaskManageService) OnCreate() {
	entity.taskData = entity.Create(new(data.TaskCache)).(*data.TaskCache)
}

func (entity *TaskManageService) AddTaskTypeInfo(req *dto.AddTaskTypeInfoReq) (resp *dto.AddTaskTypeInfoResp, err error) {
	resp = &dto.AddTaskTypeInfoResp{
		SessionId: req.SessionId,
	}
	taskTypeInfo := &models.TaskType{
		TaskType:     req.TaskType,
		Instance:     req.Instance,
		TaskNumLimit: req.TaskNumLimit,
		CrudModel: orm.CrudModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	err = entity.Create(new(models.TaskTypeDao)).(*models.TaskTypeDao).Insert(taskTypeInfo)
	if err != nil {
		return
	}
	return

}
