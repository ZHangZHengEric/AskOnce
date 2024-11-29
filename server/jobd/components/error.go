package components

import "github.com/xiangtao94/golib/pkg/errors"

var OutErrMsg = map[int]string{}

var OutErrMap = map[int]int{}

/*
*
业务校验异常
*/

var ErrorTaskNotExist = errors.Error{
	Code:    100003,
	Message: "任务不存在",
}

var ErrorTaskTimeOut = errors.Error{
	Code:    100004,
	Message: "接口超时",
}

const (
	ERROR = 4000

	MYSQL_ERR = 4001
	MQ_ERR    = 4002
)

const (
	ERR_TASK_TYPE_NOT_EXISTED   = "task type not existed"
	ERR_WAITTING_TASKS_TOO_MANY = "waitting task too many"
	ERR_LOAD_TASK_ERROR         = "load task error"
	ERR_TASK_TIMEOUT            = "task timeout"
	ERR_TASK_QUEUE_EMPTY        = "task queue empty"
	ERR_TASK_NOT_EXISTED        = "task not existed"
)
