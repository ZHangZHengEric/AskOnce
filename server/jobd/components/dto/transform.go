package dto

type PageParam struct {
	PageNo   int `json:"pageNo" binding:"required"`
	PageSize int `json:"pageSize" binding:"required"`
}

type TaskInfo struct {
	TaskId     string `json:"taskId"`
	SessionId  string `json:"sessionId"`
	TaskType   string `json:"taskType"`
	CommitTime int64  `json:"commitTime"`
	TimeoutMs  int64  `json:"timeoutMs"`
	Input      string `json:"input"`
	Output     string `json:"output"`
	Instance   string `json:"instance"`
	Status     string `json:"status"`
	UpdateTime int64  `json:"updateTime"`
}

type TaskTypeInfo struct {
	TaskType     string `db:"task_type"`
	ExtendInfo   string `db:"extend_info"`
	TaskNumLimit int    `db:"task_num_limit"` // 队列大小，0则无上限
}
