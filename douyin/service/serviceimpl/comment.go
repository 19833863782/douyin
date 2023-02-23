// Package serviceimpl
// @Author shaofan
// @Date 2022/5/13
package serviceimpl

import (
	"douyin/config"
	"douyin/entity/bo"
	"douyin/entity/myerr"
	"douyin/entity/param"
	"douyin/entity/po"
	"douyin/middleware"
	"douyin/repositories/daoimpl"
	"douyin/service"
	"douyin/util/entityutil"
	"strconv"
	"sync"
)

type Comment struct {
}

func (c Comment) Comment(commentParam *param.Comment, userId int) error {
	err := validVideoExistence(commentParam.VideoID)
	if err != nil {
		return err
	}
	if commentParam.ActionType == param.DO_COMMENT {
		return doComment(commentParam, userId)
	} else {
		return deleteComment(commentParam)
	}
}

// 发布评论
func doComment(commentParam *param.Comment, userId int) error {
	tx := daoimpl.Begin()
	wg := sync.WaitGroup{}
	wg.Add(2)
	var err error
	go func() {
		defer wg.Done()
		var comment po.Comment
		comment.SenderId = userId
		comment.VideoId = commentParam.VideoID
		comment.Content = commentParam.CommentText
		comment.Status = po.NORMAL
		err = daoimpl.NewCommentDaoInstance().Insert(&comment, tx, true)
	}()

	go func() {
		defer wg.Done()
		err = daoimpl.NewVideoDaoInstance().ChangeCommentCount(1, commentParam.VideoID, tx, true)
	}()

	wg.Wait()
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 删除评论
func deleteComment(commentParam *param.Comment) error {
	var err error
	err = validUser(commentParam.CommentId)
	if err != nil {
		return err
	}
	tx := daoimpl.Begin()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		commentDao := daoimpl.NewCommentDaoInstance()
		var comment po.Comment
		comment.ID = commentParam.CommentId
		comment.Status = po.DELETE
		err = commentDao.UpdateByCondition(&comment, tx, true)
	}()

	go func() {
		defer wg.Done()
		err = daoimpl.NewVideoDaoInstance().ChangeCommentCount(1, commentParam.VideoID, tx, true)
	}()

	wg.Wait()
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil

}

// 校验视频是否存在
func validVideoExistence(videoId int) error {
	videoDao := daoimpl.NewVideoDaoInstance()
	video, err := videoDao.QueryById(videoId)
	if err != nil {
		return err
	}
	// 视频信息不为nil，初始化id为0，为0表示不存在
	if video.ID == 0 {
		return myerr.VideoNotFound
	}
	return nil
}

func (c Comment) CommentList(videoId int) (*[]bo.Comment, error) {
	err2 := validVideoExistence(videoId)
	if err2 != nil {
		return nil, err2
	}
	//查询
	commentDao := daoimpl.NewCommentDaoInstance()
	var comment = new(po.Comment)
	comment.VideoId = videoId
	comment.Status = po.NORMAL
	comments, err := commentDao.QueryByConditionOrderByTime(comment)
	if err != nil {
		return nil, err
	}
	//转换
	var commentBos []bo.Comment
	err = entityutil.GetCommentBOS(comments, &commentBos)
	if err != nil {
		return nil, err
	}
	return &commentBos, nil
}

func validUser(commentId int) error {
	comment, err := daoimpl.NewCommentDaoInstance().QueryByConditionOrderByTime(&po.Comment{EntityModel: po.EntityModel{ID: commentId}})
	if err != nil {
		return err
	}
	currentUserId, err := strconv.Atoi(middleware.ThreadLocal.Get().(map[string]string)[config.Config.ThreadLocal.Keys.UserId])
	if err != nil {
		return err
	}
	if (*comment)[0].SenderId != currentUserId {
		return myerr.NoPermission
	}
	return nil
}

var (
	comment     service.Comment
	commentOnce sync.Once
)

// NewCommentServiceInstance 获取评论service实例
func NewCommentServiceInstance() service.Comment {
	commentOnce.Do(func() {
		comment = Comment{}
	})
	return comment
}
