package main

import (
	"dispatch/conf"
	"dispatch/helpers"
	"dispatch/models"
	"dispatch/router"
	"github.com/gin-gonic/gin"
	"github.com/xiangtao94/golib"
	"github.com/xiangtao94/golib/flow"
)

func main() {
	// 1.全局变量初始化
	helpers.PreInit()
	defer helpers.Clear()
	// 2 启动器创建
	engine := gin.New()
	golib.Bootstraps(engine, conf.WebConf)
	// 3 初始化资源
	helpers.InitResource()
	// 4.初始化http服务路由
	router.Http(engine)
	// 5 初始化建库建表
	helpers.MysqlClient.AutoMigrate(&models.TaskType{}, &models.TaskRecord{})
	// 6.框架启动
	flow.Start(engine, conf.WebConf, func(engine *gin.Engine) (err error) {
		flow.SetDefaultDBClient(helpers.MysqlClient)
		flow.SetDefaultRedisClient(helpers.RedisClient)
		return nil
	})
}
