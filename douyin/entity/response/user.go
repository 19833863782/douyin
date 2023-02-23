// Package response
// @Author shaofan
// @Date 2022/5/19
package response

import "douyin/entity/bo"

type Register struct {
	Response
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

type Login struct {
	Response
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

type UserInfo struct {
	Response
	User bo.User `json:"user"`
}
