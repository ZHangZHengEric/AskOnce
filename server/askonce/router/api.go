package router

import (
	"askonce/components/dto/dto_kdb"
	"askonce/controllers/api"
	"github.com/gin-gonic/gin"
	"github.com/xiangtao94/golib/flow"
)

func API(engine *gin.Engine) {
	router := engine.Group("/askOnce/api/v1")
	// 知识库
	kdbGroup := router.Group("kdb")
	{
		kdbGroup.POST("add", flow.Use[dto_kdb.AddReq](new(api.AddController)))          // 知识库新增
		kdbGroup.POST("update", flow.Use[dto_kdb.UpdateReq](new(api.UpdateController))) // 知识库修改
		kdbGroup.GET("info", flow.Use[dto_kdb.InfoReq](new(api.InfoController)))        // 知识库详情
		kdbGroup.POST("list", flow.Use[dto_kdb.ListReq](new(api.ListController)))       // 知识库列表
	}

}
