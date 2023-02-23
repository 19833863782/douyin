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

// Comment 			评论
func Comment(ctx *gin.Context) {
	var commentParam param.Comment
	err := ctx.ShouldBindQuery(&commentParam)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, response.ArgumentError(myerr.ArgumentInvalid(webutil.GetValidMsg(err, commentParam))))
		return
	}
	userId, err := strconv.Atoi(middleware.ThreadLocal.Get().(map[string]string)[config.Config.ThreadLocal.Keys.UserId])
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, response.SystemError)
	}

	err = serviceimpl.NewCommentServiceInstance().Comment(&commentParam, userId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, response.SystemError)
		return
	}
	ctx.JSON(http.StatusOK, response.Ok)
}

// CommentList 		查看评论列表
func CommentList(ctx *gin.Context) {
	var commentListParam param.CommentList
	err := ctx.ShouldBindQuery(&commentListParam)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, response.ArgumentError(myerr.ArgumentInvalid(webutil.GetValidMsg(err, commentListParam))))
		return
	}
	commentList, err := serviceimpl.NewCommentServiceInstance().CommentList(commentListParam.VideoId)
	ctx.JSON(http.StatusOK, response.CommentList{
		Response: response.Ok,
		Data:     *commentList,
	})
}
