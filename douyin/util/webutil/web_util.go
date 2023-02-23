// Package webutil
// @Author shaofan
// @Date 2022/5/18
package webutil

import (
	"github.com/go-playground/validator/v10"
	"reflect"
)

// GetValidMsg  通过错误获取自定义的提示信息
func GetValidMsg(err error, obj interface{}) string {
	//通过反射获取结构体
	getObj := reflect.TypeOf(obj)
	//取得错误信息
	if errs, ok := err.(validator.ValidationErrors); ok {
		//遍历所有校验错误
		for _, e := range errs {
			//遍历结构体中的字段
			for i := 0; i < getObj.NumField(); i++ {
				//当结构体中某个字段和出错的字段相同时，返回字段标签中的msg，这个msg就是自定义的错误提示
				if getObj.Field(i).Name == e.Field() {
					return getObj.Field(i).Tag.Get("msg")
				}
			}
		}
	}
	//如果没有找到该字段直接返回错误
	return err.Error()
}
