// Package douyin
// @Author chenxi
// @Date 2023/2/19
// @DESC
package main

import (
	"douyin/config"
	"douyin/job"
	"douyin/repositories/daoimpl"
	route2 "douyin/route"
	"douyin/util/jwtutil"
	"douyin/util/rabbitutil"
	"douyin/util/redisutil"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	_ "net/http/pprof"
)

var route *gin.Engine

func main() {
	go func() {
		err := http.ListenAndServe("localhost:8888", nil)
		if err != nil {
			return
		}
	}()
	err := route.Run(":" + config.Config.Server.Port)
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	// 配置文件初始化
	config.Init()
	// 数据库初始化
	daoimpl.Init()
	// 路由初始化
	route = route2.InitRoute()
	// jwt初始化
	jwtutil.InitJWT()
	// redis初始化
	redisutil.Init()
	// 消息队列初始化
	rabbitutil.Init()
	// 开始定时任务
	job.StartJob()
}
