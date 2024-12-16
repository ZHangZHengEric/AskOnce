package router

import (
	"askonce/components/dto/dto_kdb"
	"askonce/components/dto/dto_kdb_doc"
	"askonce/components/dto/dto_search"
	"askonce/controllers/kdb"
	"askonce/controllers/search"
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
		kdbGroup.GET("info", flow.Use[dto_kdb.SingleKdbReq](new(kdb.InfoController)))   // 知识库详情
		kdbGroup.POST("list", flow.Use[dto_kdb.ListReq](new(kdb.ListController)))       // 知识库列表
		// 删除
		kdbGroup.POST("delete", flow.Use[dto_kdb.SingleKdbReq](new(kdb.DeleteController))) // 删除

		docGroup := kdbGroup.Group("doc")
		{
			// 列表
			docGroup.POST("list", flow.Use[dto_kdb_doc.ListReq](new(kdb.DocListController)))
			// 新增
			docGroup.POST("add", flow.Use[dto_kdb_doc.AddReq](new(kdb.DocAddController)))
			// zip包新增
			docGroup.POST("addByZip", flow.Use[dto_kdb_doc.AddZipReq](new(kdb.DocAddByZipController)))
			// 查询知识库进度
			docGroup.POST("taskProcess", flow.Use[dto_kdb_doc.LoadProcessReq](new(kdb.TaskProcessController)))
			// 批量文本新增
			docGroup.POST("addByBatchText", flow.Use[dto_kdb_doc.AddByBatchTextReq](new(kdb.DocAddByBatchTextController)))
			// 删除
			docGroup.POST("delete", flow.Use[dto_kdb_doc.DeleteReq](new(kdb.DocDeleteController)))
			// 重做
			docGroup.POST("redo", flow.Use[dto_kdb_doc.RedoReq](new(kdb.DocRedoController)))
		}
		searchGroup := router.Group("search")
		{
			//  网页直搜
			searchGroup.POST("web", flow.Use[dto_search.WebSearchReq](new(search.WebSearchController)))
			searchGroup.POST("kdb", flow.Use[dto_search.KdbSearchReq](new(search.KdbSearchController)))
			// 智能搜索
			searchGroup.POST("chatAskSync", flow.Use[dto_search.ChatAskReq](new(search.ChatAskSyncController)))
		}
	}

}
