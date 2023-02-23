// Package entityutil
// @Author shaofan
// @Date 2022/5/13
// @DESC 实例转换工具
package entityutil

import (
	"douyin/config"
	"douyin/entity/bo"
	"douyin/entity/po"
	"douyin/middleware"
	"douyin/repositories/daoimpl"
	"strconv"
)

// GetCommentBOS 	获取评论BO实例集
// src				评论PO集
// dest 			评论bo集
func GetCommentBOS(src *[]po.Comment, dest *[]bo.Comment) error {
	if *dest == nil || len(*dest) < len(*src) {
		*dest = make([]bo.Comment, len(*src))
	}
	var ids = make([]int, len(*src))
	// 评论集合中所有用户的id集合
	for i, val := range *src {
		ids[i] = val.SenderId
		i++
	}
	// 查询所有用户信息
	userList, err := daoimpl.NewUserDaoInstance().QueryBatchIds(&ids)
	if err != nil {
		return err
	}
	// 将用户信息转为bo
	var userBOS []bo.User
	err = GetUserBOS(userList, &userBOS)
	if err != nil {
		return err
	}
	// 将用户id对应结构化数据存储到映射中
	var userMap = make(map[int]bo.User, len(userBOS))
	for i := range userBOS {
		userMap[userBOS[i].ID] = userBOS[i]
	}

	// 遍历评论po集合，按顺序给bo初始化
	for i, v := range *src {
		copyCommentProperties(&v, &(*dest)[i])
		(*dest)[i].User = userMap[v.SenderId]
	}
	return nil
}

// GetVideoBOS 		获取视频BO实例集
// src				视频PO集
// dest				视频BO集
func GetVideoBOS(src *[]po.Video, dest *[]bo.Video) error {
	if *dest == nil || len(*dest) < len(*src) {
		*dest = make([]bo.Video, len(*src))
	}
	var ids = make([]int, len(*src))
	// 评论集合中所有用户的id集合
	for i, val := range *src {
		ids[i] = val.AuthorId
		i++
	}
	// 查询所有用户信息
	userList, err := daoimpl.NewUserDaoInstance().QueryBatchIds(&ids)
	if err != nil {
		return err
	}
	// 将用户信息转为bo
	var userBOS []bo.User
	err = GetUserBOS(userList, &userBOS)
	if err != nil {
		return err
	}
	// 将用户id对应结构化数据存储到映射中
	var userMap = make(map[int]bo.User, len(userBOS))
	for i := range userBOS {
		userMap[userBOS[i].ID] = userBOS[i]
	}
	// 遍历视频po集合，按顺序给bo初始化
	for i, v := range *src {
		copyVideoProperties(&v, &(*dest)[i])
		(*dest)[i].Author = userMap[v.AuthorId]
	}

	// 如果未登录，不做是否点赞的处理，使用默认值false
	if middleware.ThreadLocal.Get() == nil || middleware.ThreadLocal.Get().(map[string]string)[config.Config.ThreadLocal.Keys.UserId] == "" {
		return nil
	}
	// 已登录的处理，对所有视频判断是否点赞
	currentUserId, err := strconv.Atoi(middleware.ThreadLocal.Get().(map[string]string)[config.Config.ThreadLocal.Keys.UserId])
	if err != nil {
		return err
	}
	// 查询所有点赞的视频
	favoriteVideoIds, err := daoimpl.NewFavoriteDaoInstance().QueryVideoIdsByUserId(currentUserId)
	if err != nil {
		return err
	}
	// 使用map存储点赞视频映射，以空间换时间
	var videoMap = make(map[int]int, len(favoriteVideoIds))
	for _, videoId := range favoriteVideoIds {
		videoMap[videoId] = 1
	}

	// 根据映射获取每个视频是否点赞
	for i, v := range *src {
		_, (*dest)[i].IsFavorite = videoMap[v.ID]
	}
	return nil
}

// GetUserBOS 		获取用户BO实例集
// src				用户PO集
// dest 			用户BO集
func GetUserBOS(users *[]po.User, dest *[]bo.User) error {
	if *dest == nil || len(*dest) < len(*users) {
		*dest = make([]bo.User, len(*users))
	}

	// 如果没有线程变量或者线程变量中没有用户id，表示没有登录，IsFollow字段设置为false
	if middleware.ThreadLocal.Get() == nil || middleware.ThreadLocal.Get().(map[string]string)[config.Config.ThreadLocal.Keys.UserId] == "" {
		for i, v := range *users {
			copyUserProperties(&v, &(*dest)[i])
			(*dest)[i].IsFollow = false
		}
		return nil
	}
	// 已登录的处理
	// 获取当前用户id
	var currentUserId int
	userId := middleware.ThreadLocal.Get().(map[string]string)[config.Config.ThreadLocal.Keys.UserId]
	var err error
	currentUserId, err = strconv.Atoi(userId)
	if err != nil {
		return err
	}

	// 查询用户的关注的集合
	allFollowsId, err := daoimpl.NewRelationDaoInstance().QueryFollowIdByFansId(currentUserId)
	if err != nil {
		return err
	}
	// 使用映射存储关注者的id，用空间换时间
	var followsMap = make(map[int]int, len(allFollowsId))
	for _, follow := range allFollowsId {
		followsMap[follow] = currentUserId //key=关注的人的id；value=目前用户id
	}

	// 遍历原切片，通过映射得到bo集合
	for i, v := range *users {
		copyUserProperties(&v, &(*dest)[i])
		_, (*dest)[i].IsFollow = followsMap[v.ID]
	}
	return nil
}

// GetUserBO 		获取单个用户BO实例
// src				用户PO
// dest				用户BO
func GetUserBO(src *po.User, dest *bo.User) error {
	// 判断是否登录
	if middleware.ThreadLocal.Get() == nil || middleware.ThreadLocal.Get().(map[string]string)[config.Config.ThreadLocal.Keys.UserId] == "" {
		(*dest).IsFollow = false
	} else {
		userId := middleware.ThreadLocal.Get().(map[string]string)[config.Config.ThreadLocal.Keys.UserId]
		uid, err := strconv.Atoi(userId)
		if err != nil {
			return err
		}
		var poFollow = po.Follow{
			FollowId:   (*src).ID,
			FollowerId: uid,
		}
		poFollows, err := daoimpl.NewRelationDaoInstance().QueryByCondition(&poFollow)
		if err != nil {
			return err
		}
		dest.IsFollow = len(*poFollows) != 0
	}
	// 复制用户相同属性
	copyUserProperties(src, dest)
	return nil
}

// GetFeedBOS 	将FeedPo集合转化为FeedBo集合
// src			FeedPO集合
// dest			FeedBO集合
func GetFeedBOS(src *[]po.Feed, dest *[]bo.Feed) {
	if *dest == nil || cap(*dest) < len(*src) {
		temp := make([]bo.Feed, len(*src))
		dest = &temp
	}
	for index, feed := range *src {
		(*dest)[index] = bo.Feed{VideoId: feed.VideoId, CreateTime: feed.CreateTime}
	}
}

// 复制用户属性
func copyUserProperties(src *po.User, dest *bo.User) {
	dest.ID = src.ID
	dest.Name = src.Name
	dest.FollowCount = src.FollowCount
	dest.FollowerCount = src.FollowerCount
}

// 复制评论属性
func copyCommentProperties(src *po.Comment, dest *bo.Comment) {
	dest.ID = src.ID
	dest.CreateDate = src.CreateTime.Format("01-02")
	dest.Content = src.Content
}

// 复制视频属性
func copyVideoProperties(src *po.Video, dest *bo.Video) {
	dest.ID = src.ID
	dest.CoverUrl = src.CoverUrl
	dest.PlayUrl = src.PlayUrl
	dest.CommentCount = src.CommentCount
	dest.FavoriteCount = src.FavoriteCount
	dest.Title = src.Title
}
