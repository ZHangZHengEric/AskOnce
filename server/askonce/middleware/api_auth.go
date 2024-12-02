package middleware

import (
	"askonce/components"
	"askonce/components/defines"
	"askonce/components/dto"
	"askonce/models"
	"github.com/gin-gonic/gin"
	"github.com/xiangtao94/golib/flow"
)

func ApiAuth(ctx *gin.Context) {
	err := ApiAuthCheck(ctx)
	if err != nil {
		flow.RenderJsonFail(ctx, err)
		ctx.Abort()
		return
	}
	ctx.Next()
}

func ApiAuthCheck(ctx *gin.Context) (err error) {
	userName := ctx.GetHeader("User-Source")
	if len(userName) == 0 {
		err = components.ErrorApiAuthError
		return
	}
	entity := flow.Create(ctx, new(models.UserDao))
	user, err := entity.GetByUserName(userName)
	if err != nil {
		return
	}
	if user == nil {
		err = components.ErrorApiAuthError
		return
	}
	userInfo := dto.LoginInfoSession{
		UserId:     user.UserId,
		Account:    user.UserName,
		LoginTime:  0,
		ExpireTime: 0,
	}
	ctx.Set(defines.LoginInfo, userInfo)
	return
}
