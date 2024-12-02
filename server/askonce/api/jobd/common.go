package jobd

import (
	"askonce/components"
	"askonce/conf"
	"github.com/bytedance/sonic"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/http"
)

type JobdApi struct {
	flow.Api
}

func (entity *JobdApi) OnCreate() {
	entity.EncodeType = http.EncodeJson
	entity.Client = conf.WebConf.Api["jobd"]
}

const (
	STATUS_FINISH         = "FINISH"
	STATUS_EXEC_FAILED    = "EXEC_FAILED"
	STATUS_WAITTING       = "WAITTING"
	STATUS_RUNNING        = "RUNNING"
	STATUS_TIMEOUT        = "TIMEOUT"
	STATUS_INPUT_MISMATCH = "INPUT_MISMATCH"
)

type JobdCommonReq struct {
	SessionId string `json:"session_id,optional"`
	TaskType  string `json:"task_type"`
	Input     string `json:"input,optional"`
	TimeoutMs int    `json:"timeout_ms,optional,default=10000"`
}

type JobdCommonRes struct {
	SessionId string `json:"session_id,optional"`
	TaskId    string `json:"task_id,optional"`
	TaskType  string `json:"task_type,optional"`
	Output    string `json:"output,optional"`
	Status    string `json:"status,optional"`
}

// param
// k 业务入参
// v 业务出参
func doTaskProcess[K any, V any](entity *JobdApi, taskType string, input K, timeoutMs int) (output V, err error) {
	inputStr, _ := sonic.MarshalString(input)
	jobdReq := JobdCommonReq{
		TaskType:  taskType,
		Input:     inputStr,
		TimeoutMs: timeoutMs,
	}
	reqOpts := http.HttpRequestOptions{
		RequestBody: jobdReq,
		Encode:      entity.GetEncodeType(),
	}
	apiRes, err := entity.ApiPostWithOpts("/jobd/committer/DoTask", reqOpts)
	if err != nil {
		return output, components.ErrorJobdError
	}
	// 接口报错处理
	jobdRes := JobdCommonRes{}
	if err = entity.DecodeApiResponse(&jobdRes, apiRes, err); err != nil {
		return output, components.ErrorJobdError
	}
	if jobdRes.Status != STATUS_FINISH {
		err = components.ErrorJobdError
		return
	}
	_ = sonic.Unmarshal([]byte(jobdRes.Output), &output)
	return output, nil
}

// param
// k 业务入参
// v 业务出参
func doTaskProcessString[K any](entity *JobdApi, taskType string, input K, timeoutMs int) (output string, err error) {
	inputStr, _ := sonic.MarshalString(input)
	jobdReq := JobdCommonReq{
		TaskType:  taskType,
		Input:     inputStr,
		TimeoutMs: timeoutMs,
	}
	reqOpts := http.HttpRequestOptions{
		RequestBody: jobdReq,
		Encode:      entity.GetEncodeType(),
	}
	apiRes, err := entity.ApiPostWithOpts("/jobd/committer/DoTask", reqOpts)
	if err != nil {
		return output, components.ErrorJobdError
	}

	// 接口报错处理
	jobdRes := JobdCommonRes{}
	if err = entity.DecodeApiResponse(&jobdRes, apiRes, err); err != nil {
		return
	}
	if jobdRes.Status != STATUS_FINISH {
		err = components.ErrorJobdError
		return
	}
	return jobdRes.Output, nil
}
