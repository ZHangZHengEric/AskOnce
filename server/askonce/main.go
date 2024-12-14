// Package backend -----------------------------
// @file      : main.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/11/21 16:55
// -------------------------------------------

package main

import (
	"askonce/conf"
	"askonce/helpers"
	"askonce/router"
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
	router.API(engine)
	// 6.框架启动
	flow.Start(engine, conf.WebConf, func(engine *gin.Engine) (err error) {
		flow.SetDefaultDBClient(helpers.MysqlClient)
		flow.SetDefaultRedisClient(helpers.RedisClient)
		// 初始化任务
		router.Tasks(engine)
		return nil
	})
}
