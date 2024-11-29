package service

import (
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/errors"
	"github.com/xiangtao94/golib/pkg/orm"
	"jobd/components"
	"jobd/components/defines"
	"jobd/components/dto"
	"jobd/data"
	"jobd/models"
	"time"
)

type WorkerService struct {
	flow.Service
	taskData *data.TaskCache
}

func (entity *WorkerService) OnCreate() {
	entity.taskData = entity.Create(new(data.TaskCache)).(*data.TaskCache)
}

func (entity *WorkerService) GetTaskForWorker(req *dto.GetTaskReq) (resp *dto.GetTaskResp, taskInfo *dto.TaskInfo, err error) {
	resp = &dto.GetTaskResp{
		SessionId: req.SessionId,
	}
	// 1. 查询有没有任务
	length, _ := entity.taskData.GetTodoTaskNum(req.TaskType)
	if length == 0 {
		return nil, nil, errors.NewError(components.MQ_ERR, components.ERR_TASK_QUEUE_EMPTY)
	}
	// 2. 获取任务
	taskInfo, err = entity.taskData.PopInputTask(req.TaskType, int64(defines.GET_TASK_TIMEOUT))
	if err != nil {
		// 查不到直接抛出错误
		return nil, nil, errors.NewError(components.MQ_ERR, components.ERR_LOAD_TASK_ERROR)
	}
	if taskInfo.TimeoutMs != 0 && taskInfo.CommitTime*1000+taskInfo.TimeoutMs >= time.Now().Unix() {
		// 任务超时了，状态设为timeout
		taskInfo.Status = defines.STATUS_TIMEOUT
		return nil, nil, errors.NewError(components.ERROR, components.ERR_TASK_TIMEOUT)
	}
	// 4. 返回任务
	resp.TaskId = taskInfo.TaskId
	resp.TaskType = taskInfo.TaskType
	resp.Input = taskInfo.Input
	// 刷新instance存活时间
	_ = entity.taskData.RefreshInstanceActive(req.Instance, taskInfo.TaskType)
	return
}

func (entity *WorkerService) UpdateTaskForWorker(req *dto.UpdateInfoReq) (resp *dto.UpdateInfoResp, err error) {
	// 更新 输出队列
	taskInfo, err := entity.taskData.UpdateTask(req.TaskId, req.Output, req.Status, req.Instance)
	if err != nil {
		return nil, err
	}
	resp = &dto.UpdateInfoResp{
		SessionId: req.SessionId,
	}
	// 刷新instance存活时间
	_ = entity.taskData.RefreshInstanceActive(req.Instance, req.TaskType)
	if req.Status == defines.STATUS_FINISH || req.Status == defines.STATUS_INPUT_MISMATCH || req.Status == defines.STATUS_EXEC_FAILED {
		err = entity.Create(new(models.TaskRecordDao)).(*models.TaskRecordDao).Insert(&models.TaskRecord{
			TaskType:  taskInfo.TaskType,
			TaskId:    taskInfo.TaskId,
			SessionId: taskInfo.SessionId,
			Status:    taskInfo.Status,
			Instance:  taskInfo.Instance,
			CrudModel: orm.CrudModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		})
		if err != nil {
			entity.LogErrorf("插入历史失败：%s", err.Error())
		}

	}
	return
}

func (entity *WorkerService) BlockBatchGetTaskForWorker(req *dto.BatchGetTaskReq) (resp *dto.GetTaskResp, err error) {
	if len(req.TaskTypes) == 0 {
		return nil, errors.NewError(components.ERROR, components.ERR_TASK_TYPE_NOT_EXISTED)
	}
	// 循环取任务，取到就返回
	// 计算超时
	var Timeout = 360 * time.Second
	now := time.Now()
	closeNotify := entity.GetCtx().Request.Context().Done()
	for {
		if time.Now().After(now.Add(Timeout)) {
			return nil, errors.NewError(components.ERROR, components.ERR_TASK_TIMEOUT)
		}
		for _, taskType := range req.TaskTypes {
			respA, taskInfo, err := entity.GetTaskForWorker(&dto.GetTaskReq{TaskType: taskType})
			select {
			case <-closeNotify:
				entity.LogWarnf("客户端连接断开")
				if taskInfo != nil {
					_ = entity.taskData.RePushInputTask(taskType, taskInfo)
				}
				return nil, nil
			default:
			}
			if err != nil {
				continue
			}
			resp = respA
			return resp, nil
		}
		time.Sleep(50 * time.Millisecond)
	}
}
