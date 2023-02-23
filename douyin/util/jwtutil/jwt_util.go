// Package jwtutil
// @Author shaofan
// @Date 2022/5/13
// @DESC jwt生成与校验工具
package jwtutil

import (
	"douyin/config"
	"github.com/dgrijalva/jwt-go"
	"time"
)

//加密秘钥
var mySigningKey []byte

func InitJWT() {
	mySigningKey = []byte(config.Config.Jwt.SecretKey)
}

// CreateJWT 	生成token
// id 			用户id
// @return 		token序列
func CreateJWT(id int) (string, error) {
	//加密方法，使用的是hs256加密
	//加密用户id
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Unix() + int64(config.Config.Jwt.ExpireTime), //过期时间: 一天
	})
	//使用秘钥加密token
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseJWT 	解析token
// token 		需要解析的token
// @return 		如果token可用，则返回token中的用户id
func ParseJWT(token string) (int, error) {
	//获取原始token
	tokenOriginal, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	//可以判断token是否过期
	if err != nil {
		return 0, err
	}
	//jwt中数值类型默认为float64
	id := tokenOriginal.Claims.(jwt.MapClaims)["id"].(float64)
	//返回int类型的id
	return int(id), nil
}
