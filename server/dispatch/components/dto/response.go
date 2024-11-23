package dto

type ClearTaskResp struct {
	SessionId string `json:"session_id,optional"`
	ErrorMsg  string `json:"error_msg,optional"`
}

type DoTaskResp struct {
	SessionId string `json:"session_id,optional"`
	TaskId    string `json:"task_id,optional"`
	TaskType  string `json:"task_type,optional"`
	Output    string `json:"output,optional"`
	Status    string `json:"status,optional"`
}

type GetInfoResp struct {
	SessionId string `json:"session_id,optional"`
	TaskId    string `json:"task_id,optional"`
	TaskType  string `json:"task_type,optional"`
	Input     string `json:"input,optional"`
	Output    string `json:"output,optional"`
	Status    string `json:"status,optional"`
}

type CommitResp struct {
	SessionId string `json:"session_id,optional"`
	TaskId    string `json:"task_id"`
}

type GetTaskNumResp struct {
	SessionId string `json:"session_id,optional"`
	ErrorMsg  string `json:"error_msg,optional"`
	WaitTasks int    `json:"wait_tasks"`
}

type GetTaskResp struct {
	SessionId string `json:"session_id,optional"`
	TaskId    string `json:"task_id"`
	TaskType  string `json:"task_type"`
	Input     string `json:"input,optional"`
}

type UpdateInfoResp struct {
	SessionId string `json:"session_id,optional"`
	ErrorMsg  string `json:"error_msg,optional"`
}

type AddTaskTypeInfoResp struct {
	SessionId string `json:"session_id,optional"`
	ErrorMsg  string `json:"error_msg,optional"`
}

type UpdateTaskTypeInfoResp struct {
	SessionId string `json:"session_id,optional"`
	ErrorMsg  string `json:"error_msg,optional"`
}

type GetAllTaskTypeInfoResp struct {
	SessionId      string                 `json:"session_id,optional"`
	ErrorMsg       string                 `json:"error_msg,optional"`
	TaskTypeInfoMp map[string]interface{} `json:"task_type_info_mp,optional"`
}

type GetAllUnFinishedTaskResp struct {
	SessionId string   `json:"session_id,optional"`
	TaskIds   []string `json:"task_ids"`
}

type JobdWsResponse struct {
	SessionId string
	ErrMsg    string
	Success   bool
}
