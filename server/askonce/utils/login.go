package utils

import (
	"askonce/components/defines"
	"askonce/components/dto"
	"github.com/gin-gonic/gin"
	"github.com/xiangtao94/golib/pkg/errors"
	"github.com/xiangtao94/golib/pkg/zlog"
)

func LoginInfo(ctx *gin.Context) (loginInfo dto.LoginInfo, err error) {
	data, has := ctx.Get(defines.LoginInfo)
	if !has {
		err = errors.ErrorUserNotLogin
		return
	}
	loginInfo, ok := data.(dto.LoginInfo)
	if !ok {
		zlog.Warnf(ctx, "LoginInfo unmarsh data failed data: %+v", data)
		err = errors.ErrorSystemError
	}
	return loginInfo, err
}
