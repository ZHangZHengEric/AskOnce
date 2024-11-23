package data

import (
	"dispatch/components"
	"dispatch/components/defines"
	"dispatch/components/dto"
	"dispatch/helpers"
	"encoding/json"
	"errors"
	"github.com/xiangtao94/golib/flow"
	"strings"
	"time"
)

// key= taskType的 队列
type TaskCache struct {
	flow.Redis
}

// 输入队列
const TaskInputKey = "InputQueue"

const TaskOutputKey = "OutputQueue"

const TaskInfoKey = "TaskInfo"

const InstanceActive = "InstanceActive"

func (entity *TaskCache) CommitTask(taskType string, taskInfo *dto.TaskInfo) (err error) {
	taskJsonT, err := json.Marshal(taskInfo)
	taskJson := string(taskJsonT)
	if err != nil {
		return err
	}
	// 输入的队列
	_, err = helpers.RedisClient.LPush(entity.FormatCacheKey("%s:%s", TaskInputKey, taskType), taskJson)
	if err != nil {
		return err
	}
	// 输出的队列
	_, err = helpers.RedisClient.LPush(entity.FormatCacheKey("%s:%s", TaskOutputKey, taskInfo.TaskId), taskJson)
	if err != nil {
		return err
	}
	// 任务详情实时缓存
	err = helpers.RedisClient.Set(entity.FormatCacheKey("%s:%s", TaskInfoKey, taskInfo.TaskId), taskJson, flow.EXPIRE_TIME_1_HOUR)
	if err != nil {
		return err
	}
	return nil
}

func (entity *TaskCache) UpdateTask(taskId string, output string, status string, instance string) (taskInfo *dto.TaskInfo, err error) {
	// 获取已有信息
	out, err := helpers.RedisClient.Get(entity.FormatCacheKey("%s:%s", TaskInfoKey, taskId))
	if err != nil {
		entity.LogErrorf("获取缓存失败%s", err.Error())
		return
	}
	if out == nil {
		return nil, components.ErrorTaskNotExist
	}
	err = json.Unmarshal(out, &taskInfo)
	taskInfo.Output = output
	taskInfo.Status = status
	taskInfo.UpdateTime = time.Now().Unix()
	taskInfo.Instance = instance
	taskJson, err := json.Marshal(taskInfo)
	if err != nil {
		return
	}
	// 更新信息
	_, err = helpers.RedisClient.LPush(entity.FormatCacheKey("%s:%s", TaskOutputKey, taskId), string(taskJson))
	if err != nil {
		return
	}
	err = helpers.RedisClient.Set(entity.FormatCacheKey("%s:%s", TaskInfoKey, taskId), string(taskJson), flow.EXPIRE_TIME_1_HOUR)
	if err != nil {
		return
	}
	if status == defines.STATUS_FINISH {
		_, err = helpers.RedisClient.Expire(entity.FormatCacheKey("%s:%s", TaskOutputKey, taskId), flow.EXPIRE_TIME_15_MINUTE)
		if err != nil {
			return nil, err
		}
	}
	return
}

func (entity *TaskCache) PopInputTask(taskType string, timeOut int64) (taskInfo *dto.TaskInfo, err error) {
	out, err := helpers.RedisClient.BRPop(entity.FormatCacheKey("%s:%s", TaskInputKey, taskType), timeOut)
	if err != nil {
		return nil, err
	}
	taskInfo = new(dto.TaskInfo)
	err = json.Unmarshal(out[1], &taskInfo)
	return taskInfo, err
}

func (entity *TaskCache) ClearInputTask(taskType string) (err error) {
	_, err = helpers.RedisClient.Del(entity.GetCtx(), entity.FormatCacheKey("%s:%s", TaskInputKey, taskType))
	if err != nil {
		return nil
	}
	return
}

func (entity *TaskCache) PopOutputTask(taskId string, timeOut int64) (ret *dto.TaskInfo, err error) {
	out, err := helpers.RedisClient.BRPop(entity.FormatCacheKey("%s:%s", TaskOutputKey, taskId), timeOut)
	if err != nil {
		return nil, nil
	}
	ret = new(dto.TaskInfo)
	err = json.Unmarshal(out[1], &ret)
	return ret, err
}

func (entity *TaskCache) PopOutputTaskV2(taskId string) (ret *dto.TaskInfo, err error) {
	out, err := helpers.RedisClient.RPop(entity.FormatCacheKey("%s:%s", TaskOutputKey, taskId))
	if err != nil {
		return nil, nil
	}
	if len(out) == 0 {
		return nil, nil
	}
	ret = new(dto.TaskInfo)
	err = json.Unmarshal(out, &ret)
	return ret, err
}

func (entity *TaskCache) GetTodoTaskNum(taskType string) (length int, err error) {
	length, err = helpers.RedisClient.LLen(entity.FormatCacheKey("%s:%s", TaskInputKey, taskType))
	if err != nil {
		return 0, errors.New(err.Error())
	}
	return
}

func (entity *TaskCache) ClearTodoTask(taskType string) error {
	_, err := helpers.RedisClient.Del(entity.GetCtx(), entity.FormatCacheKey("%s:%s", TaskInputKey, taskType))
	if err != nil {
		return err
	}
	return nil
}

func (entity *TaskCache) GetAllUnFinishedTask() (res []string, err error) {
	_, keys, err := helpers.RedisClient.SScan(entity.FormatCacheKey("%s", TaskInfoKey), 0, "*", 0)
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		res = append(res, strings.Replace(key, TaskInfoKey+":", "", -1))
	}
	return
}

func (entity *TaskCache) GetByTaskId(taskId string) (ret *dto.TaskInfo, err error) {
	// 获取已有信息
	out, err := helpers.RedisClient.Get(entity.FormatCacheKey("%s:%s", TaskInfoKey, taskId))
	if err != nil {
		entity.LogErrorf("获取缓存失败%s", err.Error())
		return
	}
	if out == nil {
		return nil, components.ErrorTaskNotExist
	}
	err = json.Unmarshal(out, &ret)
	return
}

type InstanceActiveItem struct {
	TaskType   string `json:"taskType"`
	UpdateTIme int64  `json:"UpdateTIme"`
}

func (entity *TaskCache) RefreshInstanceActive(instance string, taskType string) (err error) {
	if len(instance) == 0 {
		return
	}
	now := time.Now().Unix()
	tmp := InstanceActiveItem{
		TaskType:   taskType,
		UpdateTIme: now,
	}
	tmpStr, _ := json.Marshal(tmp)
	err = helpers.RedisClient.Set(entity.FormatCacheKey("%s:%s", InstanceActive, instance), tmpStr, flow.EXPIRE_TIME_30_MINUTE)
	if err != nil {
		return
	}
	return
}

func (entity *TaskCache) RePushInputTask(taskType string, taskInfo *dto.TaskInfo) (err error) {
	// 输入的队列
	taskJsonT, err := json.Marshal(taskInfo)
	taskJson := string(taskJsonT)

	if err != nil {
		return err
	}
	_, err = helpers.RedisClient.LPush(entity.FormatCacheKey("%s:%s", TaskInputKey, taskType), taskJson)
	if err != nil {
		return err
	}
	return
}
