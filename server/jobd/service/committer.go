package service

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/errors"
	"github.com/xiangtao94/golib/pkg/sse"
	"jobd/components"
	"jobd/components/defines"
	"jobd/components/dto"
	"jobd/data"
	"jobd/helpers"
	"time"
)

type CommitterService struct {
	flow.Service
	taskData *data.TaskCache
}

func (entity *CommitterService) OnCreate() {
	entity.taskData = entity.Create(new(data.TaskCache)).(*data.TaskCache)
}

// 异步提交任务
func (entity *CommitterService) Commit(req *dto.CommitReq) (resp *dto.CommitResp, err error) {
	if len(req.TaskType) == 0 {
		return nil, errors.NewError(components.ERROR, "taskType 不正确")
	}
	resp = &dto.CommitResp{
		SessionId: req.SessionId,
	}
	newTaskInfo := new(dto.TaskInfo)
	newTaskInfo.Input = req.Input
	newTaskInfo.CommitTime = time.Now().Unix()
	newTaskInfo.TaskId = fmt.Sprintf("%s_%v", req.TaskType, helpers.GenTaskId())
	newTaskInfo.Status = defines.STATUS_WAITTING
	newTaskInfo.TaskType = req.TaskType
	newTaskInfo.SessionId = req.SessionId
	taskId := newTaskInfo.TaskId
	if err != nil {
		return nil, err
	}
	// 2. 写入消息队列
	err = entity.taskData.CommitTask(req.TaskType, newTaskInfo)
	if err != nil {
		return nil, err
	}
	resp.TaskId = taskId
	return
}

func (entity *CommitterService) DoTaskStream(req *dto.DoTaskReq) (err error) {
	// 计算超时
	var Timeout time.Duration
	if req.TimeoutMs == 0 {
		Timeout = 3600000 * time.Millisecond
	} else {
		Timeout = time.Duration(req.TimeoutMs) * time.Millisecond
	}
	now := time.Now()
	if len(req.SessionId) == 0 {
		req.SessionId = uuid.New().String()
	}

	newTaskInfo := new(dto.TaskInfo)
	newTaskInfo.Input = req.Input
	newTaskInfo.CommitTime = time.Now().Unix()
	newTaskInfo.TaskId = fmt.Sprintf("%s_%v", req.TaskType, helpers.GenTaskId())
	newTaskInfo.Status = defines.STATUS_WAITTING
	newTaskInfo.TaskType = req.TaskType
	newTaskInfo.SessionId = req.SessionId
	taskId := newTaskInfo.TaskId
	// 2. 写入消息队列
	err = entity.taskData.CommitTask(req.TaskType, newTaskInfo)
	if err != nil {
		if err != nil {
			return err
		}
	}
	var ret *dto.TaskInfo
	for {
		ret, err = entity.taskData.PopOutputTask(taskId, int64(1*time.Second))
		if err != nil {
			if err != nil {
				return err
			}
		}
		select {
		case <-entity.GetCtx().Request.Context().Done():
			entity.LogWarnf("客户端连接断开")
			return nil
		default:
		}
		if time.Now().After(now.Add(Timeout)) {
			return errors.NewError(400, components.ERR_TASK_TIMEOUT)
		}
		if ret != nil {
			taskInfo := ret
			resp := dto.DoTaskResp{}
			resp.Output = taskInfo.Output
			resp.Status = taskInfo.Status
			resp.TaskId = taskInfo.TaskId
			resp.SessionId = req.SessionId
			resp.TaskType = taskInfo.TaskType
			str, _ := json.Marshal(resp)
			sse.RenderStream(entity.GetCtx(), "", "", string(str))
			if taskInfo.Status == defines.STATUS_FINISH || taskInfo.Status == defines.STATUS_CANCEL || taskInfo.Status == defines.STATUS_EXEC_FAILED {
				break
			}
		}
		continue
	}
	return
}

func (entity *CommitterService) DoTask(req *dto.DoTaskReq) (resp *dto.DoTaskResp, err error) {
	if len(req.TaskType) == 0 {
		return nil, errors.NewError(components.ERROR, "taskType 不正确")
	}
	// 计算超时
	var Timeout time.Duration
	if req.TimeoutMs == 0 {
		Timeout = 3600000 * time.Millisecond
	} else {
		Timeout = time.Duration(req.TimeoutMs) * time.Millisecond
	}
	now := time.Now()
	if len(req.SessionId) == 0 {
		req.SessionId = uuid.New().String()
	}
	resp = &dto.DoTaskResp{
		SessionId: req.SessionId,
	}
	newTaskInfo := new(dto.TaskInfo)
	newTaskInfo.Input = req.Input
	newTaskInfo.CommitTime = time.Now().Unix()
	newTaskInfo.TaskId = fmt.Sprintf("%s_%v", req.TaskType, helpers.GenTaskId())
	newTaskInfo.Status = defines.STATUS_WAITTING
	newTaskInfo.TaskType = req.TaskType
	newTaskInfo.SessionId = req.SessionId
	taskId := newTaskInfo.TaskId
	// 2. 写入消息队列
	err = entity.taskData.CommitTask(req.TaskType, newTaskInfo)
	if err != nil {
		return nil, err
	}
	var ret *dto.TaskInfo
	for {
		ret, err = entity.taskData.PopOutputTaskV2(taskId)
		if err != nil {
			return nil, err
		}
		select {
		case <-entity.GetCtx().Request.Context().Done():
			entity.LogWarnf("客户端连接断开")
			return nil, nil
		default:
		}
		if time.Now().After(now.Add(Timeout)) {
			return nil, errors.NewError(components.ERROR, components.ERR_TASK_TIMEOUT)
		}
		if ret != nil {
			if ret.Status == defines.STATUS_FINISH {
				break
			}
			if ret.Status != defines.STATUS_RUNNING && ret.Status != defines.STATUS_WAITTING {
				return nil, errors.NewError(components.ERROR, ret.Status+"|"+ret.Output)

			}
		}
		time.Sleep(10 * time.Millisecond)
		continue
	}
	if ret == nil {
		return
	}
	resp.Output = ret.Output
	resp.SessionId = ret.SessionId
	resp.TaskId = taskId
	resp.Status = ret.Status
	resp.TaskType = ret.TaskType
	return
}

func (entity *CommitterService) GetInfo(req *dto.GetInfoReq) (resp *dto.GetInfoResp, err error) {
	resp = &dto.GetInfoResp{
		SessionId: req.SessionId,
	}
	taskInfo, err := entity.taskData.GetByTaskId(req.TaskId)
	if err != nil {
		return nil, errors.NewError(components.ERROR, components.ERR_TASK_NOT_EXISTED)
	}
	if taskInfo == nil {
		return &dto.GetInfoResp{}, nil
	}
	resp.Output = taskInfo.Output
	resp.TaskId = req.TaskId
	resp.Status = taskInfo.Status
	resp.TaskType = taskInfo.TaskType
	return
}

func (entity *CommitterService) BlockGetInfo(req *dto.GetInfoReq) (resp *dto.GetInfoResp, err error) {
	resp = &dto.GetInfoResp{
		SessionId: req.SessionId,
	}
	ch := make(chan *dto.TaskInfo, 1)
	defer close(ch)
	for {
		ret, err := entity.taskData.PopOutputTask(req.TaskId, int64(defines.DEFAULT_TASK_TIMEOUT))
		if err != nil {
			return nil, err
		}
		if ret == nil {
			continue
		}
		if ret.Status == defines.STATUS_FINISH {
			ch <- ret
			break
		}
		if ret.Status != defines.STATUS_RUNNING && ret.Status != defines.STATUS_WAITTING {
			return nil, errors.NewError(components.ERROR, ret.Status+"|"+ret.Output)
		}
	}
	select {
	case <-time.After(defines.DEFAULT_TASK_TIMEOUT):
		return nil, errors.NewError(components.ERROR, components.ERR_TASK_TIMEOUT)
	case ret := <-ch:
		resp.Output = ret.Output
		resp.TaskId = req.TaskId
		resp.Status = ret.Status
		resp.TaskType = ret.TaskType
		return resp, nil
	}
}
