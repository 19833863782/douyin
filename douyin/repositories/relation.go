// Package repositories
// @Author shaofan
// @Date 2022/5/13
package repositories

import (
	"douyin/entity/po"
	"gorm.io/gorm"
)

// Relation 关注持久层接口
type Relation interface {

	// Insert 					增加关注关系
	// follow 					关注信息
	Insert(follow *po.Follow, tx *gorm.DB, isTx bool) error

	// DeleteByCondition 		条件删除关注关系
	// follow 					删除条件
	DeleteByCondition(follow *po.Follow, tx *gorm.DB, isTx bool) error

	//QueryByCondition          条件查询
	//follow					查询条件
	//@return					follow集合
	QueryByCondition(follow *po.Follow) (*[]po.Follow, error)

	// QueryFollowIdByFansId 	通过粉丝id查询关注id集
	// fansId 					关注id
	// @return 					关注id集
	QueryFollowIdByFansId(fansId int) ([]int, error)

	//QueryFansIdByFollowId 	通过关注id查询粉丝id集
	// followId 				关注id
	// @return 					粉丝id集
	QueryFansIdByFollowId(followId int) ([]int, error)
}
