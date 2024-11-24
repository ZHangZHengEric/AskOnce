package worker

import (
	"dispatch/components/dto"
	"dispatch/service"
	"github.com/xiangtao94/golib/flow"
)

type GetTaskCtl struct {
	flow.Controller
}

func (entity *GetTaskCtl) Action(req *dto.GetTaskReq) (res interface{}, err error) {
	s := flow.Create(entity.GetCtx(), new(service.WorkerService))
	resp, _, err := s.GetTaskForWorker(req)
	if err != nil {
		return nil, err
	}
	entity.LogInfof("worker %v 获取任务成功，task_type:%v , session_id:%v, task_id:%v", entity.GetCtx().ClientIP(), req.TaskType, resp.SessionId, resp.TaskId)
	return resp, nil
}

type BlockGetTaskCtl struct {
	flow.Controller
}

func (entity *BlockGetTaskCtl) Action(req *dto.GetTaskReq) (res interface{}, err error) {
	s := flow.Create(entity.GetCtx(), new(service.WorkerService))
	resp, err := s.BlockBatchGetTaskForWorker(&dto.BatchGetTaskReq{TaskTypes: []string{req.TaskType}})
	if err != nil {
		return nil, err
	}
	entity.LogInfof("worker %v 获取任务成功，task_type:%v , session_id:%v, task_id:%v", entity.GetCtx().ClientIP(), req.TaskType, resp.SessionId, resp.TaskId)
	return resp, nil
}

type BlockBatchGetTaskCtl struct {
	flow.Controller
}

func (entity *BlockBatchGetTaskCtl) Action(req *dto.BatchGetTaskReq) (res interface{}, err error) {
	s := flow.Create(entity.GetCtx(), new(service.WorkerService))
	resp, err := s.BlockBatchGetTaskForWorker(req)
	if err != nil {
		return nil, err
	}
	entity.LogInfof("worker %v 获取任务成功，task_type:%v, session_id:%v, task_id:%v", entity.GetCtx().ClientIP(), resp.TaskType, resp.SessionId, resp.TaskId)
	return resp, nil
}

type UpdateInfoCtl struct {
	flow.Controller
}

func (entity *UpdateInfoCtl) Action(req *dto.UpdateInfoReq) (res interface{}, err error) {
	s := flow.Create(entity.GetCtx(), new(service.WorkerService))
	resp, err := s.UpdateTaskForWorker(req)
	if err != nil {
		return nil, err
	}
	entity.LogInfof("worker %v 更新任务成功，task_type:%v , session_id:%v, task_id:%v", entity.GetCtx().ClientIP(), req.TaskType, resp.SessionId)
	return resp, nil
}
