// Package config
// @Author shaofan
// @Date 2022/5/13
package config

import "time"

// 数据库连接，gorm初始化配置
type database struct {
	DSN  string `yaml:"dsn"`
	Pool struct {
		ConnMaxLifetime time.Duration `yaml:"connMaxLifetime"`
		ConnMaxIdleTime time.Duration `yaml:"connMaxIdleTime"`
		MaxOpenConn     int           `yaml:"maxOpenConn"`
		MaxIdleConn     int           `yaml:"maxIdleConn"`
	} `yaml:"pool"`
}
