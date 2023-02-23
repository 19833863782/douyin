// Package repositories
// @Author shaofan
// @Date 2022/5/13
package repositories

import (
	"douyin/entity/po"
	"gorm.io/gorm"
)

// Favorite 点赞持久层接口
type Favorite interface {
	// Insert 					增加点赞
	// favorite 				一条点赞数据
	Insert(favorite *po.Favorite, tx *gorm.DB, isTx bool) error

	// QueryVideoIdsByUserId 	通过用户id查询视频id列表
	// userId 					用户id
	// @return 					视频id集
	QueryVideoIdsByUserId(userId int) ([]int, error)

	// DeleteByCondition		条件删除点赞数据
	// favorite					删除条件
	DeleteByCondition(favorite *po.Favorite, tx *gorm.DB, isTx bool) error
}
