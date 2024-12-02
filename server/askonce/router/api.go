package router

import (
	"askonce/components/dto/dto_kdb"
	"askonce/controllers/kdb"
	"askonce/middleware"
	"github.com/gin-gonic/gin"
	"github.com/xiangtao94/golib/flow"
)

func API(engine *gin.Engine) {
	router := engine.Group("/askonce/api/v1", middleware.ApiAuth)
	// 知识库
	kdbGroup := router.Group("kdb")
	{
		kdbGroup.POST("add", flow.Use[dto_kdb.AddReq](new(kdb.AddController)))          // 知识库新增
		kdbGroup.POST("update", flow.Use[dto_kdb.UpdateReq](new(kdb.UpdateController))) // 知识库修改
		kdbGroup.GET("info", flow.Use[dto_kdb.InfoReq](new(kdb.InfoController)))        // 知识库详情
		kdbGroup.POST("list", flow.Use[dto_kdb.ListReq](new(kdb.ListController)))       // 知识库列表
	}

}
