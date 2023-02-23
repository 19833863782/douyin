// Package repositories
// @Author shaofan
// @Date 2022/5/13
package repositories

import (
	"douyin/entity/po"
	"gorm.io/gorm"
)

// Comment 评论持久层接口
type Comment interface {
	// Insert 				增加评论
	// comment 				一条评论数据
	Insert(comment *po.Comment, tx *gorm.DB, isTx bool) error

	// QueryByConditionOrderByTime 	条件查询评论
	// comment 				查询条件，针对条件中的非空值查找
	// @return 				评论列表
	QueryByConditionOrderByTime(comment *po.Comment) (*[]po.Comment, error)

	// UpdateByCondition	条件更新评论数据
	// comment				新的评论数据
	UpdateByCondition(comment *po.Comment, tx *gorm.DB, isTx bool) error
}
