package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xiangtao94/golib/flow"
	"jobd/components/dto"
	"jobd/controllers"
	"jobd/controllers/committer"
	"jobd/controllers/worker"
	"jobd/middleware"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func Http(engine *gin.Engine) {
	engine.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))
	router := engine.Group("/jobd")

	committerGroup := router.Group("committer")
	{
		committerGroup.POST("DoTask", flow.Use[dto.DoTaskReq](new(committer.DoTaskCtl)))
		committerGroup.POST("Commit", flow.Use[dto.CommitReq](new(committer.CommitCtl)))
		committerGroup.POST("GetInfo", flow.Use[dto.GetInfoReq](new(committer.GetInfoCtl)))
		committerGroup.POST("BlockGetInfo", flow.Use[dto.GetInfoReq](new(committer.BlockGetInfoCtl)))
	}

	workerGroup := router.Group("worker")
	{
		workerGroup.POST("GetTask", flow.Use[dto.GetTaskReq](new(worker.GetTaskCtl)))
		workerGroup.POST("BlockGetTask", flow.Use[dto.GetTaskReq](new(worker.BlockGetTaskCtl)))
		workerGroup.POST("BlockBatchGetTask", flow.Use[dto.BatchGetTaskReq](new(worker.BlockBatchGetTaskCtl)))
		workerGroup.POST("UpdateInfo", middleware.TaskTimeout(10*time.Second), flow.Use[dto.UpdateInfoReq](new(worker.UpdateInfoCtl)))
	}

	apiGroup := router.Group("api")
	{
		apiGroup.POST("AddTaskTypeInfo", flow.Use[dto.AddTaskTypeInfoReq](new(controllers.AddTaskTypeInfoCtl)))
	}
}
