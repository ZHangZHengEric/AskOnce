package middleware

import (
	"askonce/components"
	"askonce/components/defines"
	"askonce/data"
	"github.com/gin-gonic/gin"
	"github.com/xiangtao94/golib/flow"
)

func LoginCheck(ctx *gin.Context) {
	err := SessionCheckLogin(ctx)
	if err != nil {
		flow.RenderJsonFail(ctx, err)
		ctx.Abort()
		return
	}
	ctx.Next()
}

// 兼容非必须登陆接口获取已登陆用户信息
func NLIGetLoginInfo(ctx *gin.Context) {
	_ = SessionCheckLogin(ctx)
}

func SessionCheckLogin(ctx *gin.Context) (err error) {
	sessionKey, _ := ctx.Cookie(defines.COOKIE_KEY)
	if len(sessionKey) == 0 {
		err = components.ErrorNotLogin
		return
	}
	entity := flow.Create(ctx, new(data.SessionCache))
	userInfo, err := entity.GetSession(sessionKey)
	if err != nil {
		return
	}
	ctx.Set(defines.LoginInfo, userInfo)
	return
}
