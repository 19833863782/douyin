// Package repositories
// @Author shaofan
// @Date 2022/5/22
package repositories

import (
	"douyin/entity/po"
	"gorm.io/gorm"
)

// Feed feed流持久层接口
type Feed interface {
	// InsertBatch 批量插入
	InsertBatch(feeds *[]po.Feed, tx *gorm.DB, isTx bool) error

	// QueryByCondition 条件查询
	QueryByCondition(feed *po.Feed) ([]po.Feed, error)

	// DeleteByCondition 条件删除
	DeleteByCondition(feed *[]po.Feed, tx *gorm.DB, isTx bool) error
}
