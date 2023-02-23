// Package response
// @Author shaofan
// @Date 2022/5/19
package response

import "douyin/entity/bo"

type FavoriteList struct {
	Response
	Data []bo.Video `json:"video_list"`
}
