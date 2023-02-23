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

type Relation struct {
}

func (r Relation) Insert(follow *po.Follow, tx *gorm.DB, isTx bool) error {
	var db1 *gorm.DB
	if isTx {
		db1 = tx
	} else {
		db1 = db
	}
	err := db1.Select([]string{"follow_id", "follower_id"}).Create(follow).Error
	if err != nil {
		return err
	}
	return nil
}

func (r Relation) DeleteByCondition(follow *po.Follow, tx *gorm.DB, isTx bool) error {
	var db1 *gorm.DB
	if isTx {
		db1 = tx
	} else {
		db1 = db
	}

	if follow.FollowId != 0 {
		db1 = db1.Where("follow_id = ?", follow.FollowId)
	}
	if follow.FollowerId != 0 {
		db1 = db1.Where("follower_id = ?", follow.FollowerId)
	}
	err := db1.Delete(&po.Follow{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r Relation) QueryByCondition(follow *po.Follow) (*[]po.Follow, error) {
	followId := follow.FollowId
	followerId := follow.FollowerId
	db1 := db
	if followId != 0 {
		db1 = db1.Where("follow_id = ?", followId)
	}
	if followerId != 0 {
		db1 = db1.Where("follower_id = ?", followerId)
	}
	var follows []po.Follow
	err := db1.Find(&follows).Error
	if err != nil {
		return nil, err
	}
	return &follows, nil
}

func (r Relation) QueryFollowIdByFansId(fansId int) ([]int, error) {
	var userIds []int
	var follows []po.Follow
	result := db.Select("follow_id").Where("follower_id = ?", fansId).Find(&follows)
	err := result.Error
	if err != nil {
		return nil, err
	}
	for _, follow := range follows {
		userIds = append(userIds, follow.FollowId)
	}
	return userIds, nil
}

func (r Relation) QueryFansIdByFollowId(followId int) ([]int, error) {
	var userIds []int
	var follows []po.Follow
	result := db.Select("follower_id").Where("follow_id = ?", followId).Find(&follows)
	err := result.Error
	if err != nil {
		return nil, err
	}
	for _, follow := range follows {
		userIds = append(userIds, follow.FollowerId)
	}
	return userIds, nil
}

var (
	relation     repositories.Relation
	relationOnce sync.Once
)

func NewRelationDaoInstance() repositories.Relation {
	relationOnce.Do(func() {
		relation = Relation{}
	})
	return relation
}
