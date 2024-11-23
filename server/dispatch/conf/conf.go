package conf

import (
	"github.com/gin-gonic/gin"
	"github.com/xiangtao94/golib/pkg/env"
	"github.com/xiangtao94/golib/pkg/http"
	"github.com/xiangtao94/golib/pkg/middleware"
	"github.com/xiangtao94/golib/pkg/orm"
	"github.com/xiangtao94/golib/pkg/redis"
	"github.com/xiangtao94/golib/pkg/zlog"
)

type SWebConf struct {
	Port    int            `yaml:"port"`
	AppName string         `yaml:"appName"`
	Log     zlog.LogConfig `yaml:"log"`

	Mysql map[string]orm.MysqlConf   `yaml:"mysql"`
	Redis map[string]redis.RedisConf `yaml:"redis"`

	Api map[string]*http.HttpClientConf `yaml:"api"` // 调用三方后台

	accessConf middleware.AccessLoggerConfig `yaml:"accessConf"`
}

var WebConf *SWebConf

func InitConf() {
	// load from yaml
	env.LoadConf("default.yaml", "mount", &WebConf)
}

func (s *SWebConf) GetZlogConf() zlog.LogConfig {
	return s.Log
}

func (s *SWebConf) GetAccessLogConf() middleware.AccessLoggerConfig {
	return s.accessConf
}

func (s *SWebConf) GetHandleRecoveryFunc() gin.RecoveryFunc {
	return nil
}

func (s *SWebConf) GetAppName() string {
	return s.AppName
}

func (s *SWebConf) GetPort() int {
	return s.Port
}
