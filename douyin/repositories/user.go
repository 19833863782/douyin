// Package repositories
// @Author shaofan
// @Date 2022/5/13
package repositories

import (
	"douyin/entity/po"
	"gorm.io/gorm"
)

// User 用户持久层接口
type User interface {
	// Insert 			插入一条数据
	// user  			用户数据
	Insert(tx *gorm.DB, isTx bool, user *po.User) (int, error)

	// QueryById 		通过id查询
	// userId 			用户id
	// @return 			用户信息
	QueryById(userId int) (*po.User, error)

	// QueryBatchIds 	id批量查询
	// userIds 			用户id集
	// @return			用户集
	QueryBatchIds(userIds *[]int) (*[]po.User, error)

	// QueryByCondition		通过已有的属性查询
	// user				用户
	// @return 			用户切片
	QueryByCondition(user *po.User) (*[]po.User, error)

	//QueryFollows 查询关注列表并且时间倒序
	QueryFollows(userId int) (*[]po.User, error)

	//QueryFans 查询粉丝列表并且时间倒序
	QueryFans(userId int) (*[]po.User, error)

	// ChangeFollowCount 修改关注数量
	ChangeFollowCount(userId, difference int, tx *gorm.DB, isTx bool) error

	// ChangeFansCount 修改粉丝数量
	ChangeFansCount(userId, difference int, tx *gorm.DB, isTx bool) error
}
