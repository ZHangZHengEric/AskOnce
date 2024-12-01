package test

import (
	"askonce/conf"
	"askonce/helpers"
	"github.com/xiangtao94/golib"
	"github.com/xiangtao94/golib/pkg/env"
	"net/http/httptest"
	"path"
	"runtime"
	"sync"

	"github.com/gin-gonic/gin"
)

var once = sync.Once{}
var Ctx *gin.Context

// Init 基础资源初始化
func Init() {
	once.Do(func() {
		engine := gin.New()
		dir := getSourcePath(0)
		env.SetAppName("testing")
		env.SetRootPath(dir + "/..")
		helpers.PreInit()
		Ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
		golib.Bootstraps(engine, conf.WebConf)
		helpers.InitResource()
	})
}

func getSourcePath(skip int) string {
	_, filename, _, _ := runtime.Caller(skip)
	return path.Dir(filename)
}
