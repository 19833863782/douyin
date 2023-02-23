// Package daoimpl
// @Author shaofan
// @Date 2022/5/22
package daoimpl

import (
	"douyin/entity/po"
	"douyin/repositories"
	"gorm.io/gorm"
	"sync"
)

type Feed struct {
}

func (f Feed) InsertBatch(feeds *[]po.Feed, tx *gorm.DB, isTx bool) error {
	var db1 *gorm.DB
	if isTx {
		db1 = tx
	} else {
		db1 = db
	}
	return db1.Omit("id", "create_time", "update_time").Create(feeds).Error
}

func (f Feed) QueryByCondition(feed *po.Feed) ([]po.Feed, error) {
	db1 := db
	if feed.VideoId != 0 {
		db1 = db1.Where("video_id = ?", feed.VideoId)
	}
	if feed.UserId != 0 {
		db1 = db1.Where("user_id = ?", feed.UserId)
	}
	var feeds = make([]po.Feed, 10)
	return feeds, db1.Find(&feeds).Error
}

func (f Feed) DeleteByCondition(feed *[]po.Feed, tx *gorm.DB, isTx bool) error {
	var db1 *gorm.DB
	if isTx {
		db1 = tx
	} else {
		db1 = db
	}
	for i, v := range *feed {
		if i == 0 {
			db1 = db1.Where("video_id = ? and user_id = ?", v.VideoId, v.UserId)
		} else {
			db1 = db1.Or("video_id = ? and user_id = ?", v.VideoId, v.UserId)
		}
	}
	return db1.Delete(po.Feed{}).Error
}

var (
	feed     repositories.Feed
	feedOnce sync.Once
)

func NewFeedDaoInstance() repositories.Feed {
	feedOnce.Do(func() {
		feed = Feed{}
	})
	return feed
}
