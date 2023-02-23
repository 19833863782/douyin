// Package serviceimpl
// @Author shaofan
// @Date 2022/5/13
package serviceimpl

import (
	"douyin/entity/bo"
	"douyin/entity/param"
	"douyin/entity/po"
	"douyin/repositories/daoimpl"
	"douyin/service"
	"douyin/util/entityutil"
	"douyin/util/rabbitutil"
	"sync"
)

type Favorite struct {
}

func (f Favorite) Like(favoriteParam *param.Favorite, userId int) error {
	err := rabbitutil.Favorite(userId, favoriteParam.VideoID, favoriteParam.ActionType == param.DO_LIKE)
	if err != nil {
		return err
	}
	return nil
}

func (f Favorite) FavoriteList(userId int) ([]bo.Video, error) {
	var poVideos *[]po.Video
	poVideos, err := daoimpl.NewVideoDaoInstance().QueryVideosByUserId(userId)
	if err != nil {
		return nil, err
	}
	var boVideos = make([]bo.Video, len(*poVideos))
	err = entityutil.GetVideoBOS(poVideos, &boVideos)
	if err != nil {
		return nil, err
	}
	return boVideos, nil
}

var (
	favorite     service.Favorite
	favoriteOnce sync.Once
)

func NewFavoriteServiceInstance() service.Favorite {
	favoriteOnce.Do(func() {
		favorite = Favorite{}
	})
	return favorite
}
