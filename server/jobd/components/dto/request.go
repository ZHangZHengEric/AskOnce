package dto

type ClearTaskReq struct {
	SessionId string `json:"session_id,optional"`
	TaskType  string `json:"task_type,optional"`
}

type DoTaskReq struct {
	SessionId string `json:"session_id,optional"`
	TaskType  string `json:"task_type"`
	Input     string `json:"input,optional"`
	TimeoutMs int    `json:"timeout_ms,optional,default=10000"`
}

type GetTaskReq struct {
	SessionId string `json:"session_id,optional"`
	TaskType  string `json:"task_type"`
	Instance  string `json:"instance"`
}

type UpdateInfoReq struct {
	SessionId string `json:"session_id,optional"`
	TaskId    string `json:"task_id"`
	TaskType  string `json:"task_type,optional"`
	Output    string `json:"output,optional"`
	Status    string `json:"status,optional"`
	Instance  string `json:"instance"`
}

type GetInfoReq struct {
	SessionId string `json:"session_id,optional"`
	TaskId    string `json:"task_id,optional"`
}

type CommitReq struct {
	SessionId string `json:"session_id,optional"`
	TaskType  string `json:"task_type"`
	Input     string `json:"input,optional"`
	TimeoutMs int    `json:"timeout_ms,optional,default=0"`
}

type GetTaskNumReq struct {
	SessionId string `json:"session_id,optional"`
	TaskType  string `json:"task_type,optional"`
}

type AddTaskTypeInfoReq struct {
	SessionId    string `json:"session_id,optional"`
	TaskType     string `json:"task_type"`
	Instance     string `json:"instance"`
	TaskNumLimit int64  `json:"task_num_limit,optional,default=0"`
}

type BatchGetTaskReq struct {
	TaskTypes []string `json:"task_types"`
}
