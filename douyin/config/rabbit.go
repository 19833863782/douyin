// Package config
// @Author shaofan
// @Date 2022/5/13
package config

// 消息队列相关配置参数
type rabbit struct {
	Url string `yaml:"url"`

	ResendMax int `yaml:"resend-max"`
	Queue     struct {
		Follow          string `yaml:"follow"`
		Favorite        string `yaml:"favorite"`
		UploadVideo     string `yaml:"upload-video"`
		FeedVideo       string `yaml:"feed-video"`
		DeadUploadVideo string `yaml:"dead-upload-video"`
		DeadFeedVideo   string `yaml:"dead-feed-video"`
	} `yaml:"queue"`

	TTL struct {
		UploadVideo int `yaml:"upload-video"`
		FeedVideo   int `yaml:"feed-video"`
	} `yaml:"ttl"`

	Key struct {
		Follow      string `yaml:"follow"`
		Favorite    string `yaml:"favorite"`
		UploadVideo string `yaml:"upload-video"`
		FeedVideo   string `yaml:"feed-video"`
	} `yaml:"key"`

	Exchange struct {
		ServiceExchange     string `yaml:"service-exchange"`
		DeadServiceExchange string `yaml:"dead-service-exchange"`
	} `yaml:"exchange"`
}
