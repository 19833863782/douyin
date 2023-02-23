// Package config
// @Author shaofan
// @Date 2022/5/16
package config

// obs连接配置
type obs struct {
	AccessKey string `yaml:"access-key"`
	SecretKey string `yaml:"secret-key"`
	EndPoint  string `yaml:"end-point"`
	Location  string `yaml:"location"`
	Buckets   struct {
		Video string `yaml:"video"`
		Cover string `yaml:"cover"`
	} `yaml:"buckets"`
}
