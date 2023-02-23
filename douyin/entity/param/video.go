// Package param
// @Author shaofan
// @Date 2022/5/22
package param

type Feed struct {
}

// Publish 投稿
type Publish struct {
	Title string `form:"title"  binding:"required" msg:"标题为空标题"`
}

// VideoList 查询用户稿件列表
type VideoList struct {
	UserID int `form:"user_id" binding:"required" msg:"无效的用户标识"`
}
