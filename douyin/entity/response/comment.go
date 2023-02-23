// Package response
// @Author shaofan
// @Date 2022/5/19
package response

import "douyin/entity/bo"

type CommentList struct {
	Response
	Data []bo.Comment `json:"comment_list"`
}
