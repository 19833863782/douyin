// Package serviceimpl
// @Author shaofan
// @Date 2022/5/13
package serviceimpl

import (
	"douyin/entity/bo"
	"douyin/entity/param"
	"douyin/entity/po"
	"douyin/repositories/daoimpl"
	"douyin/service"
	"douyin/util/entityutil"
	"douyin/util/rabbitutil"
	"sync"
)

type Relation struct {
}

func (r Relation) Follow(relationParam *param.Relation, userId int) error {
	err := rabbitutil.Follow(userId, relationParam.ToUserID, relationParam.ActionType == param.DO_FOLLOW)
	if err != nil {
		return err
	}
	return nil
}

func (r Relation) FollowList(userId int) (*[]bo.User, error) {
	var userPOS *[]po.User
	userPOS, err := daoimpl.NewUserDaoInstance().QueryFollows(userId)
	if err != nil {
		return nil, err
	}
	var userBOS = make([]bo.User, len(*userPOS))
	err = entityutil.GetUserBOS(userPOS, &userBOS)
	if err != nil {
		return nil, err
	}
	return &userBOS, nil
}

func (r Relation) FansList(userId int) (*[]bo.User, error) {

	var userPOS *[]po.User
	userPOS, err := daoimpl.NewUserDaoInstance().QueryFans(userId)
	if err != nil {
		return nil, err
	}
	var userBOS = make([]bo.User, len(*userPOS))
	err = entityutil.GetUserBOS(userPOS, &userBOS)
	if err != nil {
		return nil, err
	}
	return &userBOS, nil
}

var (
	relation     service.Relation
	relationOnce sync.Once
)

func NewRelationServiceInstance() service.Relation {
	relationOnce.Do(func() {
		relation = Relation{}
	})
	return relation
}
