// Package router -----------------------------
// @file      : commands.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/11 17:19
// -------------------------------------------
package router

import (
	"askonce/controllers"
	"github.com/gin-gonic/gin"
	"github.com/xiangtao94/golib/pkg/job/cycle"
	"time"
)

func Tasks(engine *gin.Engine) {
	// 定时任务
	startCycle(engine)
}

func startCycle(engine *gin.Engine) {
	cronJob := cycle.InitCycle(engine)
	cronJob.AddFunc(1*time.Second, controllers.BuildWaitingDoc)
	cronJob.Start()
}
