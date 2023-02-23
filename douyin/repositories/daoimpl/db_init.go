// Package daoimpl
// @Author shaofan
// @Date 2022/5/13
package daoimpl

import (
	"douyin/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var db *gorm.DB

// Init 初始化数据库连接
func Init() {
	dbConf := config.Config.DB
	// 全局变量注意这里不要使用:=符号，这个符号会创建一个局部变量，无法给全局变量赋值
	var err error
	db, err = gorm.Open(mysql.Open(dbConf.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("数据库连接错误 %v", err)
		return
	}
	pool, err := db.DB()
	if err != nil {
		log.Fatalf("连接池初始化错误 %v", err)
		return
	}
	// 连接池配置
	pool.SetConnMaxIdleTime(dbConf.Pool.ConnMaxIdleTime)
	pool.SetMaxOpenConns(dbConf.Pool.MaxOpenConn)
	pool.SetMaxIdleConns(dbConf.Pool.MaxIdleConn)
	pool.SetConnMaxLifetime(dbConf.Pool.ConnMaxLifetime)
}

// Begin 开启事务
func Begin() *gorm.DB {
	return db.Begin()
}
