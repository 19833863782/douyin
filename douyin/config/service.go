// Package config
// @Author shaofan
// @Date 2022/5/17
package config

type service struct {
	BigVNum      int    `yaml:"big-v-num"`
	VideoTempDir string `yaml:"video-temp-dir"`
	CoverTempDir string `yaml:"cover-temp-dir"`
	PageSize     int    `yaml:"page-size"`
	FFMPEGPath   string `yaml:"ffmpeg-path"`
	FeedLoop     bool   `yaml:"feed-loop"`
}
