// Package daoimpl
// @Author shaofan
// @Date 2022/5/13
package daoimpl

import (
	"douyin/entity/po"
	"douyin/repositories"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Video struct {
}

func (v Video) ChangeCommentCount(difference, videoId int, tx *gorm.DB, isTx bool) error {
	var db1 *gorm.DB
	if isTx {
		db1 = tx
	} else {
		db1 = db
	}
	if err := db1.Model(&po.Video{EntityModel: po.EntityModel{ID: videoId}}).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", difference)).Error; err != nil {
		return err
	}
	return nil
}

func (v Video) ChangeFavoriteCount(difference, videoId int, tx *gorm.DB, isTx bool) error {
	var db1 *gorm.DB
	if isTx {
		db1 = tx
	} else {
		db1 = db
	}
	if err := db1.Model(&po.Video{EntityModel: po.EntityModel{ID: videoId}}).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", difference)).Error; err != nil {
		return err
	}
	return nil
}

func (v Video) QueryVideosByUserId(userId int) (*[]po.Video, error) {
	var poVideos []po.Video
	err := db.Raw("SELECT v.* FROM dy_video v,dy_favorite f WHERE v.`id`= f.`video_id` AND f.`user_id` = ? ORDER BY f.`create_time` DESC", userId).Scan(&poVideos).Error
	if err != nil {
		return nil, err
	}
	return &poVideos, err
}

func (v Video) UpdateByCondition(video *po.Video, tx *gorm.DB, isTx bool) error {
	var client *gorm.DB
	if isTx {
		client = tx
	} else {
		client = db
	}
	return client.Model(video).Updates(video).Error
}

func (v Video) QueryById(id int) (*po.Video, error) {
	db1 := db
	var video po.Video
	if id != 0 {
		db1 = db1.Where("id = ?", id)
	}
	err := db1.Find(&video).Error
	return &video, err
}

func (v Video) Insert(tx *gorm.DB, video *po.Video, isTx bool) error {
	var db1 *gorm.DB
	if isTx {
		db1 = tx
	} else {
		db1 = db
	}
	return db1.Omit("id", "create_time", "update_time").Create(video).Error
}

func (v Video) QueryBatchIds(videoIds *[]int, size int) ([]po.Video, error) {
	var videos = make([]po.Video, len(*videoIds))
	return videos, db.Where("id in (?)", *videoIds).Order("create_time DESC").Limit(size).Find(&videos).Error
}

func (v Video) QueryByConditionTimeDESC(condition *po.Video) (*[]po.Video, error) {
	db1 := db
	var videos []po.Video
	if condition.ID != 0 {
		db1 = db1.Where("id = ?", condition.ID)
	}
	if condition.AuthorId != 0 {
		db1 = db1.Where("author_id = ?", condition.AuthorId)
	}
	if condition.Title != "" {
		db1 = db1.Where("name = ?", condition.Title)
	}
	err := db1.Order("create_time desc").Find(&videos).Error
	return &videos, err
}

func (v Video) QueryByLatestTimeDESC(latestTime time.Time, size int) (*[]po.Video, error) {
	db1 := db
	var videos []po.Video
	db1 = db1.Where("create_time < ?", latestTime)
	err := db1.Order("create_time desc").Limit(size).Find(&videos).Error
	return &videos, err
}

var (
	video     repositories.Video
	videoOnce sync.Once
)

func NewVideoDaoInstance() repositories.Video {
	videoOnce.Do(func() {
		video = Video{}
	})
	return video
}
