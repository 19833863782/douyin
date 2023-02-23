// Package response
// @Author shaofan
// @Date 2022/5/19
package response

import "douyin/entity/bo"

type FollowList struct {
	Response
	Data []bo.User `json:"user_list"`
}

type FansList struct {
	Response
	Data []bo.User `json:"user_list"`
}
