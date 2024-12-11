// Package controllers -----------------------------
// @file      : job.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/11 17:21
// -------------------------------------------
package controllers

import (
	"askonce/service"
	"github.com/gin-gonic/gin"
	"github.com/xiangtao94/golib/flow"
)

func BuildWaitingDoc(ctx *gin.Context) error {
	s := flow.Create(ctx, new(service.KdbDocService))
	return s.BuildWaitingDoc()
}
