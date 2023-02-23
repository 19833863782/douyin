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
	"douyin/util/obsutil"
	"douyin/util/webutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Feed 		刷视频
func Feed(ctx *gin.Context) {
	var latestTime int64
	var timeParam = ctx.Query("latest_time")
	var err error
	videoService := serviceimpl.NewVideoServiceInstance()
	if timeParam == "" {
		latestTime = time.Now().UnixMilli()
	} else {
		latestTime, err = strconv.ParseInt(timeParam, 10, 64)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, response.SystemError)
			return
		}
	}
	var userId = 0
	var isLogin = false
	if middleware.ThreadLocal.Get() != nil {
		isLogin = true
		userId, err = strconv.Atoi(middleware.ThreadLocal.Get().(map[string]string)[config.Config.ThreadLocal.Keys.UserId])
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, response.SystemError)
			return
		}
	}
	feed, nextTime, err := videoService.Feed(userId, isLogin, latestTime)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, response.Feed{
		VideoList: feed,
		NextTime:  nextTime,
	})
}

// Publish 		投稿视频
func Publish(ctx *gin.Context) {

	// 通过线程获取投稿人id
	authorId, err := strconv.Atoi(middleware.ThreadLocal.Get().(map[string]string)[config.Config.ThreadLocal.Keys.UserId])
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, response.SystemError)
		return
	}

	// 通过请求参数获取视频标题
	var publishParam param.Publish
	err = ctx.Bind(&publishParam)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, response.ArgumentError(myerr.ArgumentInvalid(webutil.GetValidMsg(err, publishParam))))
		return
	}

	// 获取视频文件
	video, err := ctx.FormFile("data")
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(myerr.VideoNotFound))
		return
	}
	if !obsutil.IsVideo(video.Filename) {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(myerr.FileError))
	}
	// 发布
	err = serviceimpl.NewVideoServiceInstance().Publish(ctx, video, authorId, publishParam.Title)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusProxyAuthRequired, response.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, response.PubVideo{
		Response: response.Ok,
	})
}

// VideoList 	查看视频发布列表
func VideoList(ctx *gin.Context) {
	// 通过请求参数获取对象id
	var videoList param.VideoList
	err := ctx.ShouldBindQuery(&videoList)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, response.ArgumentError(myerr.ArgumentInvalid(webutil.GetValidMsg(err, videoList))))
		return
	}

	boVideos, err := serviceimpl.NewVideoServiceInstance().VideoList(videoList.UserID)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, response.VideoList{
		Response:  response.Ok,
		VideoList: boVideos,
	})
}
