package test

import (
	"github.com/xiangtao94/golib/pkg/env"
	"jobd/helpers"
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
		dir := getSourcePath(0)
		env.SetAppName("testing")
		env.SetRootPath(dir + "/..")
		helpers.PreInit()
		helpers.InitResource()
		Ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	})
}

func getSourcePath(skip int) string {
	_, filename, _, _ := runtime.Caller(skip)
	return path.Dir(filename)
}
