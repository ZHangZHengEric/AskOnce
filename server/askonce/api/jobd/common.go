package jobd

import (
	"askonce/components"
	"askonce/conf"
	"encoding/json"
	"github.com/bytedance/sonic"
	jsoniter "github.com/json-iterator/go"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/errors"
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
	httpRes, err := entity.Client.HttpPost(entity.GetCtx(), "/jobd/committer/DoTaskInner", reqOpts)
	if err != nil {
		return output, components.ErrorJobdError
	}

	apiRes, err := entity.handel(httpRes)
	if err != nil {
		return output, components.ErrorJobdError
	}
	// 接口报错处理
	jobdRes := JobdCommonRes{}
	if err = entity.decodeApiResponse(&jobdRes, apiRes, err); err != nil {
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
	httpRes, err := entity.Client.HttpPost(entity.GetCtx(), "/jobd/committer/DoTaskInner", reqOpts)
	if err != nil {
		return output, components.ErrorJobdError
	}

	apiRes, err := entity.handel(httpRes)
	if err != nil {
		return output, components.ErrorJobdError
	}
	// 接口报错处理
	jobdRes := JobdCommonRes{}
	if err = entity.decodeApiResponse(&jobdRes, apiRes, err); err != nil {
		return
	}
	if jobdRes.Status != STATUS_FINISH {
		err = components.ErrorJobdError
		return
	}
	return jobdRes.Output, nil
}

type OldRes struct {
	ErrNo  int                 `json:"errNo"`
	ErrMsg string              `json:"errMsg"`
	Data   jsoniter.RawMessage `json:"data"`
}

func (entity *JobdApi) handel(res *http.HttpResult) (*flow.ApiRes, error) {
	httpRes := OldRes{}
	if len(res.Response) > 0 {
		e := jsoniter.Unmarshal(res.Response, &httpRes)
		if e != nil {
			// 限制一下错误日志打印的长度，2k
			data := res.Response
			if len(data) > 2000 {
				data = data[0:2000]
			}
			// 返回数据json unmarshal失败，打印错误日志
			entity.LogErrorf("http response json unmarshal failed, response:%s, err:%s", string(data), e)
			return nil, e
		}
	}
	if httpRes.ErrNo != 0 {
		entity.LogErrorf("http call has error,  errNo:%d, errMsg:%s", httpRes.ErrNo, httpRes.ErrMsg)
	}
	apiRes := &flow.ApiRes{
		Code:    httpRes.ErrNo,
		Message: httpRes.ErrMsg,
		Data:    httpRes.Data,
	}
	return apiRes, nil
}

func (entity *JobdApi) decodeApiResponse(outPut interface{}, data *flow.ApiRes, err error) error {
	if err != nil {
		return err
	}

	if data.Code != 0 {
		return errors.Error{
			Code:    data.Code,
			Message: data.Message,
		}
	}
	if len(data.Data) > 0 {
		// 解析数据
		if err = json.Unmarshal(data.Data, outPut); err != nil {
			entity.LogErrorf("api error, api response unmarshal, data:%s, err:%+v", data.Data, err.Error())
			return errors.ErrorSystemError
		}

	}
	return nil
}
