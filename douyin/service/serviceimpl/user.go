// Package serviceimpl
// @Author shaofan
// @Date 2022/5/13
package serviceimpl

import (
	"douyin/config"
	"douyin/entity/bo"
	"douyin/entity/myerr"
	"douyin/entity/param"
	"douyin/entity/po"
	"douyin/repositories/daoimpl"
	"douyin/service"
	"douyin/util/encryptionutil"
	"douyin/util/entityutil"
	"douyin/util/jwtutil"
	"douyin/util/redisutil"
	"strconv"
	"sync"
	"time"
)

type UserServiceImpl struct {
}

var (
	user     service.User
	userOnce sync.Once
)

func NewUserService() service.User {
	userOnce.Do(func() {
		user = UserServiceImpl{}
	})
	return user
}
func (UserServiceImpl) UserInfo(userId int) (*bo.User, error) {
	userBo := bo.User{}
	userPo, err := daoimpl.NewUserDaoInstance().QueryById(userId) //调用dao层根据id查用户方法
	if err != nil {
		return nil, err
	}
	if userPo == nil {
		return nil, myerr.UserNotFound
	}
	err = entityutil.GetUserBO(userPo, &userBo)
	if err != nil {
		return nil, err
	}
	return &userBo, nil
}
func (UserServiceImpl) Register(userParam param.User) (int, string, error) {
	userPo := po.User{}
	userPo.Name = userParam.UserName
	users, err := daoimpl.NewUserDaoInstance().QueryByCondition(&userPo)
	if err != nil {
		return 0, "", err
	}
	if len(*users) != 0 {
		return 0, "", myerr.UserNameExist
	}
	userPo.Password, err = encryptionutil.Encryption(userParam.Password) //调用md5密码加密工具方法
	if err != nil {
		return -6, "", err
	}
	userPo.FollowCount = 0 //初始关注数和粉丝数都应是0
	userPo.FollowerCount = 0
	userId, err := daoimpl.NewUserDaoInstance().Insert(nil, false, &userPo) //执行插入并返回用户id
	jwt, err := jwtutil.CreateJWT(userId)
	if err != nil {
		return 0, "", err
	}
	err = saveTokenToRedis(userId, jwt)
	if err != nil {
		return 0, "", err
	}
	return userId, jwt, nil
}
func (UserServiceImpl) Login(userParam param.User) (int, string, error) {
	userBo := po.User{}
	userBo.Name = userParam.UserName
	users, err := daoimpl.NewUserDaoInstance().QueryByCondition(&userBo)
	if err != nil {
		return 0, "", err
	}
	if len(*users) == 0 {
		return 0, "", myerr.UserNotFound
	}
	tt, err := encryptionutil.EncryptionCompare(userParam.Password, (*users)[0].Password) //调用md5加密对比工具方法
	if err != nil {
		return 0, "", err
	}
	if tt == false {
		return -1, "", myerr.LoginError
	}
	jwt, err := jwtutil.CreateJWT((*users)[0].ID)
	if err != nil {
		return 0, "", err
	}
	err = saveTokenToRedis((*users)[0].ID, jwt)
	if err != nil {
		return 0, "", err
	}
	return (*users)[0].ID, jwt, nil
}
func saveTokenToRedis(userId int, token string) error {
	tokenExpireTime, err := time.ParseDuration(config.Config.Redis.ExpireTime.Token)
	if err != nil {
		return err
	}
	err = redisutil.SetWithExpireTime(config.Config.Redis.Key.Token+strconv.Itoa(userId), &token, tokenExpireTime)
	if err != nil {
		return err
	}
	return nil
}
