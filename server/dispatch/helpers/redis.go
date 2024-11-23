package helpers

import (
	"dispatch/conf"
	"github.com/xiangtao94/golib/pkg/redis"
)

// 推荐，直接使用
var RedisClient *redis.Redis

// 初始化redis
func InitRedis() {
	c := conf.WebConf.Redis["default"]
	var err error
	RedisClient, err = redis.InitRedisClient(c)
	if err != nil || RedisClient == nil {
		panic("init redis failed!")
	}
}

func CloseRedis() {
	_ = RedisClient.Close()
}
