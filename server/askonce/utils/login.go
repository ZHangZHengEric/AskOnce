package utils

import (
	"askonce/components/defines"
	"askonce/components/dto"
	"github.com/gin-gonic/gin"
	"github.com/xiangtao94/golib/pkg/errors"
	"github.com/xiangtao94/golib/pkg/zlog"
	"net"
	"strings"
)

func LoginInfo(ctx *gin.Context) (loginInfo dto.LoginInfoSession, err error) {
	data, has := ctx.Get(defines.LoginInfo)
	if !has {
		err = errors.ErrorUserNotLogin
		return
	}
	loginInfo, ok := data.(dto.LoginInfoSession)
	if !ok {
		zlog.Warnf(ctx, "LoginInfo unmarsh data failed data: %+v", data)
		err = errors.ErrorSystemError
	}
	return loginInfo, err
}

func GetCookieDomain(host string) (domain string) {
	hostArr := strings.Split(host, ":")
	if len(hostArr) != 0 {
		if net.ParseIP(hostArr[0]) == nil {
			//域名
			domain = hostArr[0]
		}
	}
	return domain
}
