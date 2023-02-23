// Package config
// @Author shaofan
// @Date 2022/5/13
package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

var (
	OutboxExpireTime time.Duration
	InboxExpireTime  time.Duration
)

// Config 配置信息
type config struct {
	DB           database    `yaml:"database"`
	Redis        redis       `yaml:"redis"`
	Rabbit       rabbit      `yaml:"rabbit"`
	Server       server      `yaml:"server"`
	Obs          obs         `yaml:"obs"`
	Service      service     `yaml:"service"`
	ThreadLocal  threadLocal `yaml:"thread-local"`
	Jwt          jwt         `yaml:"jwt"`
	StandardDate string      `yaml:"standard-date"`
}

// Config 全局配置实例
var Config *config

// 读取yml文件
func readConfig() {
	file, err := ioutil.ReadFile("./config/config.yml")
	if err != nil {
		log.Fatalln("读取文件config.yml发生错误", err)
		return
	}
	if err := yaml.Unmarshal(file, Config); err != nil {
		log.Fatalln("解析文件config.yml发生错误", err)
		return
	}
}

// Init 配置初始化
func Init() {
	Config = &config{}
	readConfig()
	var err error
	OutboxExpireTime, err = time.ParseDuration(Config.Redis.ExpireTime.Outbox)
	if err != nil {
		log.Fatalln(err)
	}
	InboxExpireTime, err = time.ParseDuration(Config.Redis.ExpireTime.Inbox)
	if err != nil {
		log.Fatalln(err)
	}
}
