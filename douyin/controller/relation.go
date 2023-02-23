// Package controller
// @Author shaofan
// @Date 2022/5/13
package controller

import (
	"douyin/config"
	"douyin/entity/myerr"
	"douyin/entity/param"
	"douyin/entity/response"
	"douyin/middleware"
	"douyin/service/serviceimpl"
	"douyin/util/webutil"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

var relationService = serviceimpl.NewRelationServiceInstance()

// Follow 		关注与取关
func Follow(ctx *gin.Context) {

	var relation param.Relation

	err := ctx.ShouldBindQuery(&relation)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(myerr.ArgumentInvalid(webutil.GetValidMsg(err, relation))))
		return
	}

	userId, err := strconv.Atoi(middleware.ThreadLocal.Get().(map[string]string)[config.Config.ThreadLocal.Keys.UserId])
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, response.SystemError)
	}
	err = relationService.Follow(&relation, userId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusForbidden, response.SystemError)
		return
	}
	ctx.JSON(http.StatusOK, response.Ok)
}

// FollowList 	查看关注列表
func FollowList(ctx *gin.Context) {
	var followListParam param.FollowList
	err := ctx.ShouldBindQuery(&followListParam)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, response.ArgumentError(myerr.ArgumentInvalid(webutil.GetValidMsg(err, followListParam))))
		return
	}
	userList, err := relationService.FollowList(followListParam.UserID)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusForbidden, response.SystemError)
		return
	}
	ctx.JSON(http.StatusOK, response.FollowList{
		Response: response.Ok,
		Data:     *userList,
	})
	return
}

// FansList 	查看粉丝列表
func FansList(ctx *gin.Context) {
	var fansListParam param.FollowList
	err := ctx.ShouldBindQuery(&fansListParam)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, response.ArgumentError(myerr.ArgumentInvalid(webutil.GetValidMsg(err, fansListParam))))
		return
	}
	userList, err := relationService.FansList(fansListParam.UserID)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusForbidden, response.SystemError)
		return
	}
	ctx.JSON(http.StatusOK, response.FansList{
		Response: response.Ok,
		Data:     *userList,
	})
	return
}
