// Package daoimpl
// @Author shaofan
// @Date 2022/5/13
package daoimpl

import (
	"douyin/entity/po"
	"douyin/repositories"
	"gorm.io/gorm"
	"sync"
)

type UserImpl struct {
}

func (i UserImpl) ChangeFollowCount(userId, difference int, tx *gorm.DB, isTx bool) error {
	var db1 *gorm.DB
	if isTx {
		db1 = tx
	} else {
		db1 = db
	}
	return db1.Model(&po.User{EntityModel: po.EntityModel{ID: userId}}).UpdateColumn("follow_count", gorm.Expr("follow_count + ?", difference)).Error
}

func (i UserImpl) ChangeFansCount(userId, difference int, tx *gorm.DB, isTx bool) error {
	var db1 *gorm.DB
	if isTx {
		db1 = tx
	} else {
		db1 = db
	}
	return db1.Model(&po.User{EntityModel: po.EntityModel{ID: userId}}).UpdateColumn("follower_count", gorm.Expr("follower_count + ?", difference)).Error
}

var (
	user     repositories.User
	userOnce sync.Once
)

// NewUserDaoInstance 单例模式，下次使用如果有这个实例就用这个实例，没有时再创建
func NewUserDaoInstance() repositories.User {
	userOnce.Do(func() {
		user = UserImpl{}
	})
	return user
}

func (UserImpl) QueryById(userId int) (*po.User, error) {
	db1 := db
	user := po.User{}
	err := db1.Where("id = ?", userId).Find(&user).Error //根据id查数据
	return &user, err
}
func (UserImpl) Insert(tx *gorm.DB, isTx bool, user *po.User) (int, error) {
	var client *gorm.DB
	if isTx {
		client = tx
	} else {
		client = db
	}
	var err error
	err = client.Omit("id", "create_time", "update_time").Create(user).Error
	return (*user).ID, err
}
func (UserImpl) QueryBatchIds(userIds *[]int) (*[]po.User, error) {
	db1 := db
	userList := make([]po.User, len(*userIds)) //根据一堆id查多个用户数据
	var err error
	err = db1.Where("id IN ?", *userIds).Find(&userList).Error
	return &userList, err
}
func (UserImpl) QueryByCondition(user *po.User) (*[]po.User, error) {
	db1 := db
	var err error
	var users []po.User
	if user.ID != 0 {
		db1 = db1.Where(" id=?", user.ID)
	}
	if user.Name != "" {
		db1 = db1.Where(" name=?", user.Name)
	}
	err = db1.Find(&users).Error
	return &users, err
}

//QueryFollows 查询关注列表并且时间倒序
func (i UserImpl) QueryFollows(userId int) (*[]po.User, error) {
	var poUsers = make([]po.User, 0)
	err := db.Raw("SELECT  u.*  FROM dy_user AS u,dy_follow AS f WHERE f.follower_id=? AND u.id=f.follow_id ORDER BY f.update_time DESC", userId).Scan(&poUsers).Error
	if err != nil {
		return nil, err
	}
	return &poUsers, nil
}

//QueryFans 查询粉丝列表并且时间倒序
func (i UserImpl) QueryFans(userId int) (*[]po.User, error) {
	var poUsers = make([]po.User, 0)
	err := db.Raw("SELECT  u.*  FROM dy_user AS u,dy_follow AS f WHERE f.follow_id=? AND u.id=f.follower_id ORDER BY f.update_time DESC", userId).Scan(&poUsers).Error
	if err != nil {
		return nil, err
	}
	return &poUsers, nil
}
