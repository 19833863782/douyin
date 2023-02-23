// Package config
// @Author shaofan
// @Date 2022/5/18
package config

// 线程变量相关配置
type threadLocal struct {
	Keys struct {
		UserId string `yaml:"user-id"`
	} `yaml:"keys"`
}
