package defines

import (
	"github.com/xiangtao94/golib/flow"
	"time"
)

const (
	COOKIE_DEFAULT_AGE = flow.EXPIRE_TIME_1_WEEK
)

const (
	STATUS_ENQUE_ERR      = "ENQUE_ERROR"
	STATUS_FINISH         = "FINISH"
	STATUS_CANCEL         = "CANCEL"
	STATUS_EXEC_FAILED    = "EXEC_FAILED"
	STATUS_WAITTING       = "WAITTING"
	STATUS_RUNNING        = "RUNNING"
	STATUS_TIMEOUT        = "TIMEOUT"
	STATUS_INPUT_MISMATCH = "INPUT_MISMATCH"
)

const (
	GET_TASK_TIMEOUT     = 1000 * time.Millisecond // 取任务超时
	CYCLE_QUERY_WAIT     = 1000 * time.Millisecond // 查询任务是否完成的轮询时间
	DEFAULT_TASK_TIMEOUT = 100 * time.Second       // 一个任务的最长执行时间
)
