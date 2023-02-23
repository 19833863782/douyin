// Package middleware
// @Author shaofan
// @Date 2022/5/13
package middleware

import (
	"douyin/config"
	"douyin/entity/myerr"
	"douyin/entity/response"
	"douyin/util/jwtutil"
	"douyin/util/redisutil"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// JWTAuth jwt鉴权
func JWTAuth(ctx *gin.Context) {
	var err error
	var tokenString string
	//判断投稿接口
	path := ctx.Request.URL.Path
	isFeed := strings.Contains(path, "/feed")
	isPublish := strings.Contains(path, "/publish/action")
	if isPublish {
		//投稿接口 通过form-data获取token
		tokenString = ctx.PostForm("token")
	} else if isFeed {
		// feed流接口，只用获取token，且没有token不拦截
		tokenString = ctx.Query("token")
		if tokenString == "" {
			ctx.Next()
			return
		}
	} else {
		//获取token,没有token需要拦截
		tokenString = ctx.Query("token")
	}

	//解析token
	err, uid := parseToken(ctx, tokenString)
	if err != nil {
		log.Println(err)
		return
	}
	if !isFeed {
		//从redis判断token是否有效
		err = tokenValid(ctx, uid)
		if err != nil {
			log.Println(err)
			return
		}
	}
	//存入uid，供线程变量中间件操作
	ctx.Set(config.Config.ThreadLocal.Keys.UserId, uid)
	//继续执行下面的程序
	ctx.Next()
}

//解析token
func parseToken(ctx *gin.Context, token string) (error, int) {
	//解析出token中的用户id
	uid, err := jwtutil.ParseJWT(token)
	if err != nil {
		//token过期
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(myerr.InvalidToken))
		ctx.Abort()
	}
	return err, uid
}

//判断redis中token是否有效
func tokenValid(ctx *gin.Context, userId int) error {
	var redisToken string
	err := redisutil.Get[string](config.Config.Redis.Key.Token+strconv.Itoa(userId), &redisToken)
	if err != nil || redisToken == "" {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(myerr.InvalidToken))
		ctx.Abort()
		return errors.New("redis中token不存在")
	}
	return err
}
