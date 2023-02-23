// Package config
// @Author shaofan
// @Date 2022/5/13
package config

type redis struct {
	Url      string `yaml:"url"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`

	Key struct {
		Token         string `yaml:"token"`
		MessageBackup string `yaml:"message-backup"`
		ErrorMessage  string `yaml:"error-message"`
		Outbox        string `yaml:"outbox"`
		Inbox         string `yaml:"inbox"`
	} `yaml:"key"`

	ExpireTime struct {
		Token  string `yaml:"token"`
		Outbox string `yaml:"outbox"`
		Inbox  string `yaml:"inbox"`
	} `yaml:"expire-time"`
}
