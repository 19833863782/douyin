// Package route
// @Author shaofan
// @Date 2022/5/13
package route

import (
	"douyin/config"
	"douyin/controller"
	"douyin/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

// InitRoute 初始化接口路由
func InitRoute() *gin.Engine {
	route := gin.Default()
	route.MaxMultipartMemory = 8 << 20
	route.StaticFS("/public", http.Dir("./public"))

	withAUTH := route.Group(config.Config.Server.Name, middleware.JWTAuth, middleware.SaveUserId)
	withAUTH.GET("/feed", controller.Feed)
	user1 := withAUTH.Group("/user")
	{
		user1.GET("/", controller.UserInfo)
	}
	publish := withAUTH.Group("/publish")
	{
		publish.GET("/list/", controller.VideoList)
		publish.POST("/action/", controller.Publish)
	}
	favorite := withAUTH.Group("/favorite")
	{
		favorite.GET("/list/", controller.FavoriteList)
		favorite.POST("/action/", controller.Like)
	}
	comment := withAUTH.Group("/comment")
	{
		comment.GET("/list/", controller.CommentList)
		comment.POST("/action/", controller.Comment)
	}
	relation := withAUTH.Group("/relation")
	{
		relation.GET("/follow/list/", controller.FollowList)
		relation.GET("/follower/list/", controller.FansList)
		relation.POST("/action/", controller.Follow)
	}

	withoutAUTH := route.Group(config.Config.Server.Name)
	user2 := withoutAUTH.Group("/user/")
	{
		user2.POST("/register/", controller.Register)
		user2.POST("/login/", controller.Login)
	}
	return route
}
