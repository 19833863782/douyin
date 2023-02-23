// Package daoimpl
// @Author shaofan
// @Date 2022/5/13
package daoimpl

import (
	"douyin/entity/po"
	"douyin/repositories"
	"gorm.io/gorm"
	"sync"
)

type Comment struct {
}

func (c Comment) Insert(comment *po.Comment, tx *gorm.DB, isTx bool) error {
	var db1 *gorm.DB
	if isTx {
		db1 = tx
	} else {
		db1 = db
	}
	return db1.Omit("ID", "create_time", "update_time").Create(comment).Error
}

func (c Comment) QueryByConditionOrderByTime(comment *po.Comment) (*[]po.Comment, error) {
	var comments []po.Comment
	if comment.ID != 0 {
		return &comments, db.Where("id=?", comment.ID).Find(&comments).Error
	}
	db1 := db
	if comment.SenderId != 0 {
		db1 = db1.Where("sender_id = ?", comment.SenderId)
	}
	if comment.VideoId != 0 {
		db1 = db1.Where("video_id = ?", comment.VideoId)
	}
	if comment.Status != 0 {
		db1 = db1.Where("status = ?", comment.Status)
	}
	return &comments, db1.Model(comment).Order("create_time desc").Find(&comments).Error
}

func (c Comment) UpdateByCondition(comment *po.Comment, tx *gorm.DB, isTx bool) error {
	var db1 *gorm.DB
	if isTx {
		db1 = tx
	} else {
		db1 = db
	}
	return db1.Model(comment).Updates(comment).Error
}

var (
	comment     repositories.Comment
	commentOnce sync.Once
)

// NewCommentDaoInstance 获取评论dao实例
func NewCommentDaoInstance() repositories.Comment {
	commentOnce.Do(func() {
		comment = Comment{}
	})
	return comment
}
