// Package controller
// @Author shaofan
// @Date 2022/5/13
package controller

import (
	"douyin/entity/myerr"
	"douyin/entity/param"
	"douyin/entity/response"
	"douyin/service/serviceimpl"
	"douyin/util/webutil"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Register 		用户注册
func Register(context *gin.Context) {
	var user param.User
	err := context.ShouldBindQuery(&user)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusBadRequest, response.ArgumentError(myerr.ArgumentInvalid(webutil.GetValidMsg(err, user))))
		return
	}
	userId, token, err := serviceimpl.NewUserService().Register(user)
	if err != nil { //注册失败
		log.Println(err)
		context.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
	} else { //注册成功
		context.JSON(http.StatusOK, response.Register{
			Response: response.Ok,
			UserId:   userId,
			Token:    token,
		})
	}
}

// Login 			用户登录
func Login(context *gin.Context) {
	var user param.User
	err := context.ShouldBindQuery(&user)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusBadRequest, response.ArgumentError(myerr.ArgumentInvalid(webutil.GetValidMsg(err, user))))
		return
	}
	userId, token, err := serviceimpl.NewUserService().Login(user)
	if err != nil { //登录失败
		log.Println(err)
		context.JSON(http.StatusProxyAuthRequired, response.ErrorResponse(err))
		return
	}
	context.JSON(http.StatusOK, response.Register{
		Response: response.Ok,
		UserId:   userId,
		Token:    token,
	})
}

// UserInfo 		查看用户信息
func UserInfo(context *gin.Context) {
	var userInfoParam param.UserInfo
	err := context.ShouldBindQuery(&userInfoParam)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ArgumentError(myerr.ArgumentInvalid(webutil.GetValidMsg(err, userInfoParam))))
		return
	}
	user, err := serviceimpl.NewUserService().UserInfo(userInfoParam.UserId)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse(err))
	} else {
		context.JSON(http.StatusOK, response.UserInfo{
			Response: response.Ok,
			User:     *user,
		})
	}
}
