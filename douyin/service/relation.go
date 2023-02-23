// Package service
// @Author shaofan
// @Date 2022/5/13
package service

import (
	"douyin/entity/bo"
	"douyin/entity/param"
)

// Relation 		关注关系业务接口
type Relation interface {
	// Follow 		关注操作
	// userId 		自己的id
	// toUserId 	对方的id
	// actionType 	操作类型，1-关注，2-取消关注
	Follow(relationParam *param.Relation, userId int) error

	// FollowList 	查看关注列表
	// userId 		用户id
	// @return 		关注的用户列表
	FollowList(userId int) (*[]bo.User, error)

	// FansList 	查看粉丝列表
	// userId 		用户id
	// @return 		粉丝用户列表
	FansList(userId int) (*[]bo.User, error)
}
