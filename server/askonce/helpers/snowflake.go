package helpers

import (
	"github.com/bwmarrin/snowflake"
	"time"
)

var userNode *snowflake.Node

func InitSnowflake(startTime string) {
	workerId, _ := RedisClient.Incr("WorkerId")
	var st time.Time
	// 格式化 1月2号下午3时4分5秒  2006年
	st, err := time.Parse("2006-01-02", startTime)
	snowflake.Epoch = st.UnixNano() / 1e6
	snowflake.NodeBits = 5
	snowflake.StepBits = 2
	machineID := workerId % 32
	userNode, err = snowflake.NewNode(machineID)
	if err != nil || userNode == nil {
		panic("init snowflake.Node failed!")
	}
	return
}

// 生成 64 位的 雪花 ID
func GenUserID() string {
	return userNode.Generate().String()
}

// 生成 64 位的 雪花 ID
func GenID() int64 {
	return userNode.Generate().Int64()
}

// 生成 64 位的 雪花 ID
func GenIDStr() string {
	return userNode.Generate().String()
}
