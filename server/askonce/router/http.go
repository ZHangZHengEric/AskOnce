package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	_ "net/http/pprof"
)

func Http(engine *gin.Engine) {
	engine.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))
	//router := engine.Group("/askOnce")

}
