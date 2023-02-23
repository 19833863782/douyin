// Package response
// @Author shaofan
// @Date 2022/5/19
package response

import "douyin/entity/bo"

type Feed struct {
	Response
	NextTime  int64      `json:"next_time"`
	VideoList []bo.Video `json:"video_list"`
}

type VideoList struct {
	Response
	VideoList []bo.Video `json:"video_list"`
}

type PubVideo struct {
	Response
}
