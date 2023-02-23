// Package service
// @Author shaofan
// @Date 2022/5/13
package service

import (
	"douyin/entity/bo"
	"douyin/entity/param"
)

// Favorite				点赞业务接口
type Favorite interface {
	// Like 			点赞操作
	// favoriteParam 	点赞参数
	// userId			当前操作的用户id
	Like(favoriteParam *param.Favorite, userId int) error

	// FavoriteList 	点赞列表
	// userId 			用户id
	FavoriteList(userId int) ([]bo.Video, error)
}
