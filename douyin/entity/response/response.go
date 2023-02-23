// Package response
// @Author shaofan
// @Date 2022/5/19
package response

import "douyin/entity/myerr"

type Response struct {
	Code    int    `json:"status_code"`
	Message string `json:"status_msg"`
}

// SystemError 系统错误
var SystemError = Response{Code: -1, Message: "系统维护中"}

// Ok 正常无数据响应
var Ok = Response{Code: 0, Message: "ok"}

// ErrorResponse 获取错误响应
// message 错误信息
func ErrorResponse(err error) Response {
	switch err {
	case myerr.UserNameExist:
		return Response{Code: -1, Message: "该账号已被注册"}
	case myerr.UserNotFound:
		return Response{Code: -1, Message: "该用户不存在"}
	case myerr.NoPermission:
		return Response{Code: -1, Message: "您没有权限操作"}
	case myerr.InvalidToken:
		return Response{Code: -1, Message: "您的凭证无效"}
	case myerr.VideoNotFound:
		return Response{Code: -1, Message: "该视频不存在"}
	case myerr.LoginError:
		return Response{Code: -1, Message: "用户名或密码错误"}
	case myerr.FileError:
		return Response{Code: -1, Message: "文件格式错误"}
	default:
		return Response{Code: -1, Message: "系统维护中"}
	}
}

// ArgumentError 参数错误
func ArgumentError(err error) Response {
	return Response{-1, err.Error()}
}
