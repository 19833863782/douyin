package redisutil

import (
	"douyin/config"
	"github.com/go-redis/redis"
	"log"
)

var RedisDB *redis.Client

// Init redis初始化
func Init() {
	// 建立连接
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Url,
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.DB,
	})
	err := RedisDB.Ping().Err()
	if err != nil {
		log.Fatalln(err)
	}
}
