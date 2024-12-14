package conf

import (
	"askonce/gpt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/xiangtao94/golib/pkg/env"
	"github.com/xiangtao94/golib/pkg/http"
	"github.com/xiangtao94/golib/pkg/middleware"
	"github.com/xiangtao94/golib/pkg/orm"
	"github.com/xiangtao94/golib/pkg/redis"
	"github.com/xiangtao94/golib/pkg/zlog"
	"path/filepath"
)

type SWebConf struct {
	Port       int                           `yaml:"port"`
	AppName    string                        `yaml:"appName"`
	Log        zlog.LogConfig                `yaml:"log"`
	accessConf middleware.AccessLoggerConfig `yaml:"accessConf"`

	Mysql         map[string]orm.MysqlConf   `yaml:"mysql"`
	Redis         map[string]redis.RedisConf `yaml:"redis"`
	ElasticSearch SElasticConf               `yaml:"elastic"`
	MongoDb       map[string]SMongoDBConf    `yaml:"mongodb"`

	Api        map[string]*http.HttpClientConf `yaml:"api"` // 调用三方后台
	MinioConf  SMinioConf                      `yaml:"minioConf"`
	Channel    map[string]gpt.ChannelConf      `yaml:"channel"`
	EsDbConfig string                          `yaml:"esDbConfig"`
}

type SMongoDBConf struct {
	DataBase string `yaml:"database"`
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type SElasticConf struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type SMinioConf struct {
	AK   string `yaml:"ak"`
	SK   string `yaml:"sk"`
	Addr string `yaml:"addr"`
}

var WebConf *SWebConf

func InitConf() {
	envPath := filepath.Join(env.GetRootPath(), "/../../deploy/.env")
	_ = godotenv.Load(envPath)
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
