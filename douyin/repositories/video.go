// Package repositories
// @Author shaofan
// @Date 2022/5/13
package repositories

import (
	"douyin/entity/po"
	"gorm.io/gorm"
	"time"
)

// Video 视频持久层接口
type Video interface {
	// Insert 						插入
	// video 						视频信息
	// tx							事务操作对象
	// isTx							是否处于事务中
	Insert(tx *gorm.DB, video *po.Video, isTx bool) error

	// QueryBatchIds 				id批量查询
	// videoIds 					视频id集
	// size							查询数量大小
	// @return 						视频列表
	QueryBatchIds(videoIds *[]int, size int) ([]po.Video, error)

	// QueryByConditionTimeDESC 	条件查询并按时间倒序排列
	// condition 					字段匹配查询条件
	// @return 						视频列表
	QueryByConditionTimeDESC(condition *po.Video) (*[]po.Video, error)

	// QueryByLatestTimeDESC 		查询某个时间点之前的视频，时间倒序
	// latestTime					上一次最有一条视频时间
	// @return 						视频列表
	QueryByLatestTimeDESC(latestTime time.Time, size int) (*[]po.Video, error)

	// QueryById 					根据id查询
	// id							视频id
	// @return						视频实体
	QueryById(id int) (*po.Video, error)

	// UpdateByCondition 			条件更新
	// video						更新条件
	UpdateByCondition(video *po.Video, tx *gorm.DB, isTx bool) error

	// QueryVideosByUserId		通过用户id联表查询
	// userId 					用户id
	// @return 					(倒序)视频集合
	QueryVideosByUserId(userId int) (*[]po.Video, error)

	// ChangeFavoriteCount 		修改点赞的数量
	// videoId					视频id
	ChangeFavoriteCount(difference, videoId int, tx *gorm.DB, isTx bool) error

	// ChangeCommentCount		修改评论数量
	// videoId 					视频id
	ChangeCommentCount(difference, videoId int, tx *gorm.DB, isTx bool) error
}
