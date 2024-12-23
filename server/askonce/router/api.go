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
		// 知识库新增
		kdbGroup.POST("add", flow.Use[dto_kdb.AddReq](new(kdb.AddController)))
		// 知识库修改
		kdbGroup.POST("update", flow.Use[dto_kdb.UpdateReq](new(kdb.UpdateController)))
		// 知识库详情
		kdbGroup.GET("info", flow.Use[dto_kdb.SingleKdbReq](new(kdb.InfoController)))
		// 知识库列表
		kdbGroup.POST("list", flow.Use[dto_kdb.ListReq](new(kdb.ListController)))
		// 知识库删除
		kdbGroup.POST("delete", flow.Use[dto_kdb.KdbDeleteReq](new(kdb.DeleteController)))
		// 知识库文档
		docGroup := kdbGroup.Group("doc")
		{
			// 列表
			docGroup.POST("list", flow.Use[dto_kdb_doc.ListReq](new(kdb.DocListController)))

			docGroup.GET("info", flow.Use[dto_kdb_doc.InfoReq](new(kdb.DocInfoController)))
			// 新增
			docGroup.POST("add", flow.Use[dto_kdb_doc.AddReq](new(kdb.DocAddController)))
			// zip包新增
			docGroup.POST("addByZip", flow.Use[dto_kdb_doc.AddZipReq](new(kdb.DocAddByZipController)))
			// 查询知识库进度
			docGroup.GET("taskProcess", flow.Use[dto_kdb_doc.TaskProcessReq](new(kdb.TaskProcessController)))
			// 任务重做
			docGroup.POST("taskRedo", flow.Use[dto_kdb_doc.TaskRedoReq](new(kdb.TaskRedoController)))
			// 批量文本新增
			docGroup.POST("addByBatchText", flow.Use[dto_kdb_doc.AddByBatchTextReq](new(kdb.DocAddByBatchTextController)))
			// 删除
			docGroup.POST("delete", flow.Use[dto_kdb_doc.DeleteReq](new(kdb.DocDeleteController)))
			// 重做
			docGroup.POST("redo", flow.Use[dto_kdb_doc.RedoReq](new(kdb.DocRedoController)))
		}
		searchGroup := router.Group("search")
		{
			// 网页直搜
			searchGroup.POST("web", flow.Use[dto_search.WebSearchReq](new(search.WebSearchController)))
			// 根据会话继续搜索
			searchGroup.POST("session", flow.Use[dto_search.SessionSearchReq](new(search.SessionSearchController)))
			// 知识库直搜
			searchGroup.POST("kdb", flow.Use[dto_search.KdbSearchReq](new(search.KdbSearchController)))
			// 对话搜索（同步接口）
			searchGroup.POST("chatAskSync", flow.Use[dto_search.ChatAskReq](new(search.ChatAskSyncController)))

			// 问题关注点生成(搜索知识库）
			searchGroup.POST("questionFocus", flow.Use[dto_search.QuestionFocusReq](new(search.QuestionFocusController)))
			// 报告搜索（同步接口，输出docx）
			searchGroup.POST("reportAskSync", flow.Use[dto_search.ReportAskReq](new(search.ReportAskController)))
			// 报告搜索（输出docx）
			searchGroup.POST("reportToDocx", flow.Use[dto_search.ReportDocxReq](new(search.ReportDocxController)))
		}
	}

}
