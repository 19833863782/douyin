// Package myerr
// @Author shaofan
// @Date 2022/5/19
package myerr

import "errors"

var (
	VideoNotFound = errors.New("视频不存在")
	UserNameExist = errors.New("用户名已存在")
	UserNotFound  = errors.New("用户不存在")
	InvalidToken  = errors.New("用户凭证无效")
	NoPermission  = errors.New("无权限")
	LoginError    = errors.New("用户名或密码错误")
	FileError     = errors.New("文件格式错误")
)

// ArgumentInvalid 参数无效
func ArgumentInvalid(message string) error {
	return errors.New(message)
}
