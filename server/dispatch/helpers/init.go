package helpers

import (
	"dispatch/conf"
	"github.com/xiangtao94/golib/pkg/zlog"
)

func PreInit() {
	// 初始化配置
	conf.InitConf()
}

func InitResource() {
	// 初始化mysql
	InitMysql()
	// 初始化redis
	InitRedis()
	InitSnowflake("2024-01-01")
}

func Clear() {
	// 服务结束时的清理工作，对应 InitResource() 初始化的资源
	zlog.CloseLogger()
}
