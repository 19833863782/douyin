// Package param
/**
 * @Author yg
 * @Date 2022-05-17
 * @Description
 **/
package param

const (
	DO_FOLLOW     = 1
	CANCEL_FOLLOW = 2
)

//Relation 关注参数
type Relation struct {
	ToUserID   int  `json:"to_user_id" form:"to_user_id" binding:"required" msg:"无对方id"`
	ActionType byte `json:"action_type" form:"action_type"  binding:"required" msg:"无效的操作类型"`
}

type FollowList struct {
	UserID int `form:"user_id" binding:"required" msg:"无效的用户标识"`
}
