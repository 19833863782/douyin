// Package config
// @Author shaofan
// @Date 2022/5/13
package config

type server struct {
	Port      string `yaml:"port"`
	Name      string `yaml:"name"`
	Protocol  string `yaml:"protocol"`
	IP        string `yaml:"ip"`
	Proxy     string `yaml:"proxy"`
	WithProxy bool   `yaml:"with-proxy"`
}
