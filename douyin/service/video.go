// Package service
// @Author shaofan
// @Date 2022/5/13
package service

import (
	"douyin/entity/bo"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

// Video				视频业务接口
type Video interface {
	// Feed 			获取feed流
	// userId 			用户id，可以为空
	// isLogin 			用户是否登录，用户未登录时用户id无效
	// latestTime		最新投稿时间戳
	// @return 			视频列表
	Feed(userId int, isLogin bool, latestTime int64) ([]bo.Video, int64, error)

	// Publish 			发布视频
	// file 			视频文件
	// userId 			用户id
	// ctx				上下文，用于保存视频到本地
	Publish(ctx *gin.Context, video *multipart.FileHeader, userId int, title string) error

	// VideoList 		查看视频发布列表
	// userId			用户id
	// @return			视频列表
	VideoList(userId int) ([]bo.Video, error)
}
