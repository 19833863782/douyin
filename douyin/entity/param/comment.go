// Package param
// @Author shaofan
// @Date 2022/5/13
package param

const (
	DO_COMMENT     = 1
	DELETE_COMMENT = 2
)

// Comment 上传与删除参数
type Comment struct {
	VideoID     int    `form:"video_id" `
	ActionType  byte   `form:"action_type" binding:"required" msg:"无效的操作类型"`
	CommentText string `form:"comment_text"`
	CommentId   int    `form:"comment_id"`
}

// CommentList 查询评论列表
type CommentList struct {
	VideoId int `form:"video_id" binding:"required" msg:"无效的参数"`
}
