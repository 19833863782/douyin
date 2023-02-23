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

var favoriteService = serviceimpl.NewFavoriteServiceInstance()

// Like 			点赞与取消赞操作
func Like(ctx *gin.Context) {

	var favorite param.Favorite

	err := ctx.ShouldBindQuery(&favorite)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, response.ArgumentError(myerr.ArgumentInvalid(webutil.GetValidMsg(err, favorite))))
		return
	}
	userId, err := strconv.Atoi(middleware.ThreadLocal.Get().(map[string]string)[config.Config.ThreadLocal.Keys.UserId])
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, response.SystemError)
	}
	err = favoriteService.Like(&favorite, userId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, response.SystemError)
		return
	}
	ctx.JSON(http.StatusOK, response.Ok)
}

// FavoriteList 	查看点赞列表
func FavoriteList(ctx *gin.Context) {
	var favoriteListParam param.FavoriteList

	err := ctx.ShouldBindQuery(&favoriteListParam)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, response.ArgumentError(myerr.ArgumentInvalid(webutil.GetValidMsg(err, favoriteListParam))))
		return
	}

	videoList, err := favoriteService.FavoriteList(favoriteListParam.UserId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, response.SystemError)
		return
	}
	ctx.JSON(http.StatusOK, response.FavoriteList{
		Response: response.Ok,
		Data:     videoList,
	})
	return
}
