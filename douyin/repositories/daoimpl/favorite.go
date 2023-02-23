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

type Favorite struct {
}

func (f Favorite) Insert(favorite *po.Favorite, tx *gorm.DB, isTx bool) error {
	var db1 *gorm.DB
	if isTx {
		db1 = tx
	} else {
		db1 = db
	}
	if err := db1.Select([]string{"video_id", "user_id"}).Create(favorite).Error; err != nil {
		return err
	}
	return nil
}

func (f Favorite) QueryVideoIdsByUserId(userId int) ([]int, error) {
	var videoIds = make([]int, 5)
	var favorites = make([]po.Favorite, 5)
	err := db.Select("video_id").Where("user_id = ?", userId).Find(&favorites).Error
	if err != nil {
		return nil, err
	}
	for _, video := range favorites {
		videoIds = append(videoIds, video.VideoId)
	}
	return videoIds, nil
}

func (f Favorite) DeleteByCondition(favorite *po.Favorite, tx *gorm.DB, isTx bool) error {
	var db1 *gorm.DB
	if isTx {
		db1 = tx
	} else {
		db1 = db
	}
	videoId := favorite.VideoId
	userId := favorite.UserId
	if videoId != 0 {
		db1 = db1.Where("video_id = ?", videoId)
	}
	if userId != 0 {
		db1 = db1.Where("user_id = ?", userId)
	}
	err := db1.Delete(&po.Favorite{}).Error
	if err != nil {
		return err
	}
	return nil
}

var (
	favorite     repositories.Favorite
	favoriteOnce sync.Once
)

func NewFavoriteDaoInstance() repositories.Favorite {
	favoriteOnce.Do(func() {
		favorite = Favorite{}
	})
	return favorite
}
