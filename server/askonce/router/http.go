package router

import (
	"askonce/components/dto"
	"askonce/components/dto/dto_config"
	"askonce/components/dto/dto_history"
	"askonce/components/dto/dto_kdb"
	"askonce/components/dto/dto_knowledge"
	"askonce/components/dto/dto_search"
	"askonce/components/dto/dto_user"
	"askonce/controllers/web/config"
	"askonce/controllers/web/history"
	"askonce/controllers/web/kdb"
	"askonce/controllers/web/search"
	"askonce/controllers/web/user"
	"askonce/middleware"
	"github.com/gin-gonic/gin"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/sse"
	"net/http"
	_ "net/http/pprof"
)

func Http(engine *gin.Engine) {
	engine.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))
	router := engine.Group("/askonce")
	userGroup := router.Group("user")
	{
		userGroup.POST("/registerByAccount", flow.Use[dto_user.RegisterAccountReq](new(user.RegisterAccountController))) // 账户注册
		userGroup.POST("/loginByAccount", flow.Use[dto_user.LoginAccountReq](new(user.LoginAccountController)))          // 账户登录
		userGroup.GET("/loginInfo", middleware.LoginCheck, flow.Use[dto.EmptyReq](new(user.LoginInfoController)))        // 账户信息
	}
	knowledgeGroup := router.Group("kdb") // 知识库
	{
		// 封面
		knowledgeGroup.GET("covers", middleware.LoginCheck, flow.Use[dto.EmptyReq](new(kdb.CoversController)))

		// 列表
		knowledgeGroup.POST("list", middleware.LoginCheck, flow.Use[dto_kdb.ListReq](new(kdb.ListController)))
		// 新增
		knowledgeGroup.POST("add", middleware.LoginCheck, flow.Use[dto_kdb.AddReq](new(kdb.AddController)))
		// 修改
		knowledgeGroup.POST("update", middleware.LoginCheck, flow.Use[dto_kdb.UpdateReq](new(kdb.UpdateController)))
		// 删除
		knowledgeGroup.POST("delete", middleware.LoginCheck, flow.Use[dto_kdb.DeleteReq](new(kdb.DeleteController)))
		// 详情
		knowledgeGroup.GET("detail", middleware.LoginCheck, flow.Use[dto_kdb.InfoReq](new(kdb.InfoController)))
		// 删除自己与知识库关系
		knowledgeGroup.POST("deleteSelf", middleware.LoginCheck, flow.Use[dto_kdb.DeleteSelfReq](new(kdb.DeleteSelfController)))
		// 判断是否有权限
		knowledgeGroup.POST("auth", middleware.LoginCheck, flow.Use[dto_kdb.AuthReq](new(kdb.AuthController)))

		// 知识库用户列表
		knowledgeGroup.POST("userList", middleware.LoginCheck, flow.Use[dto_kdb.UserListReq](new(kdb.UserListController)))
		// 用户查询
		knowledgeGroup.POST("userQuery", middleware.LoginCheck, flow.Use[dto_kdb.UserQueryReq](new(kdb.UserQueryController)))
		// 知识库用户新增
		knowledgeGroup.POST("userAdd", middleware.LoginCheck, flow.Use[dto_kdb.UserAddReq](new(kdb.UserAddController)))
		// 知识库用户删除
		knowledgeGroup.POST("userDelete", middleware.LoginCheck, flow.Use[dto_kdb.UserDeleteReq](new(kdb.UserDeleteController)))
		// 知识库分享码生成
		knowledgeGroup.POST("shareCodeGen", middleware.LoginCheck, flow.Use[dto_kdb.GenShareCodeReq](new(kdb.GenShareCodeController))) // 知识库用户删除
		// 知识库分享码验证
		knowledgeGroup.POST("shareCodeVerify", middleware.LoginCheck, flow.Use[dto_kdb.VerifyShareCodeReq](new(kdb.VerifyShareCodeController))) // 知识库用户删除
		// 知识库分享码信息
		knowledgeGroup.GET("shareCodeInfo", middleware.LoginCheck, flow.Use[dto_kdb.InfoShareCodeReq](new(kdb.ShareCodeInfoController))) // 知识库用户删除

		// 召回测试
		knowledgeGroup.POST("recall", middleware.LoginCheck, flow.Use[dto_kdb.RecallReq](new(kdb.RecallController)))

		docGroup := knowledgeGroup.Group("doc")
		{
			docGroup.POST("list", middleware.LoginCheck, flow.Use[dto_knowledge.DataListReq](new(kdb.DataListController)))             // 列表
			docGroup.POST("add", middleware.LoginCheck, flow.Use[dto_knowledge.DataAddReq](new(kdb.DataAddController)))                // 新增
			docGroup.POST("batchAdd", middleware.LoginCheck, flow.Use[dto_knowledge.DataBatchAddReq](new(kdb.DataBatchAddController))) // 新增
			docGroup.POST("delete", middleware.LoginCheck, flow.Use[dto_knowledge.DataDeleteReq](new(kdb.DataDeleteController)))       // 删除
			docGroup.POST("redo", middleware.LoginCheck, flow.Use[dto_knowledge.DataRedoReq](new(kdb.DataRedoController)))             // 删除
		}
	}
	configGroup := router.Group("config")
	{
		configGroup.GET("detail", middleware.LoginCheck, flow.Use[dto_config.DetailReq](new(config.DetailController)))
		configGroup.GET("dict", middleware.LoginCheck, flow.Use[dto.EmptyReq](new(config.DictController)))
		configGroup.POST("save", middleware.LoginCheck, flow.Use[dto_config.SaveReq](new(config.SaveController)))
	}
	searchGroup := router.Group("search")
	{
		// 智能搜索 用例
		searchGroup.GET("case", flow.Use[dto.EmptyReq](new(search.CaseController)))
		// 智能搜索 可选知识库列表
		searchGroup.POST("kdbList", middleware.NLIGetLoginInfo, flow.Use[dto_search.KdbListReq](new(search.KdbListController)))
		// 智能搜索 session
		searchGroup.GET("genSession", middleware.NLIGetLoginInfo, flow.Use[dto.EmptyReq](new(search.SessionController)))
		// 智能搜索 踩一下
		searchGroup.POST("unlike", middleware.NLIGetLoginInfo, flow.Use[dto_search.UnlikeReq](new(search.UnlikeController)))
		// 智能搜索
		searchGroup.POST("ask", middleware.NLIGetLoginInfo, sse.UploadEventStream, flow.Use[dto_search.AskReq](new(search.AskController)))
		// 智能搜索 参考
		searchGroup.POST("refer", middleware.NLIGetLoginInfo, flow.Use[dto_search.ReferReq](new(search.ReferController)))
		// 智能搜索 历史
		searchGroup.POST("his", middleware.NLIGetLoginInfo, flow.Use[dto_search.HisReq](new(search.HisController)))
		// 智能搜索 大纲
		searchGroup.POST("outline", middleware.NLIGetLoginInfo, flow.Use[dto_search.OutlineReq](new(search.OutlineController)))
		// 智能搜索 相关
		searchGroup.POST("relation", middleware.NLIGetLoginInfo, flow.Use[dto_search.RelationReq](new(search.RelationController)))
		// 智能搜索 进度
		searchGroup.POST("process", middleware.NLIGetLoginInfo, flow.Use[dto_search.ProcessReq](new(search.ProcessController)))
	}
	historyGroup := router.Group("history")
	{
		historyGroup.POST("ask", middleware.NLIGetLoginInfo, flow.Use[dto_history.AskReq](new(history.AskController)))
	}
}
