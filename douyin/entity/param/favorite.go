// Package param
// @Author shaofan
// @Date 2022/5/13
package param

const (
	DO_LIKE     = 1
	CANCEL_LIKE = 2
)

// Favorite 点赞参数
type Favorite struct {
	VideoID    int  `form:"video_id" binding:"required" msg:"无效的视频标识"`
	ActionType byte `form:"action_type"  binding:"required" msg:"无效的操作类型"`
}

// FavoriteList 点赞列表参数
type FavoriteList struct {
	UserId int `form:"user_id" binding:"required" msg:"无效的用户标识"`
}
