// Package serviceimpl
// @Author shaofan
// @Date 2022/5/13
package serviceimpl

import (
	"bytes"
	"douyin/config"
	"douyin/entity/bo"
	"douyin/entity/po"
	"douyin/repositories/daoimpl"
	"douyin/service"
	"douyin/util/entityutil"
	"douyin/util/obsutil"
	"douyin/util/rabbitutil"
	"douyin/util/redisutil"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type Video struct {
}

func (v Video) Feed(userId int, isLogin bool, latestTime int64) ([]bo.Video, int64, error) {
	videoDao := daoimpl.NewVideoDaoInstance()
	// 视频集合
	var videos []po.Video
	// 从redis中取出的视频id集合
	var videoIds []int
	// 收件箱的feed
	var inbox []bo.Feed
	// 发件箱的feed
	var outbox []bo.Feed
	// 可能出现的错误
	var err error
	// 如果用户已登录，查询redis中的feed流
	if isLogin {
		wait := sync.WaitGroup{}
		wait.Add(2)
		go func(inbox *[]bo.Feed, userId int) {
			defer wait.Done()
			*inbox, err = fromInbox(userId)
		}(&inbox, userId)
		go func(outbox *[]bo.Feed, userId int) {
			defer wait.Done()
			*outbox, err = fromOutbox(userId)
		}(&outbox, userId)
		wait.Wait()
		if err != nil {
			return nil, 0, err
		}
		videoIds, err = mergeBox(&inbox, &outbox, latestTime)
		if err != nil {
			return nil, 0, err
		}
	}
	// 得到视频id集合，从数据库中查询视频数据
	videos, err = videoDao.QueryBatchIds(&videoIds, config.Config.Service.PageSize)
	if err != nil {
		return nil, 0, err
	}
	// 获取查询视频的时间条件，如果从redis从查到了数据，则为数据最后一条的时间，否则为前端传来的时间
	// 如果redis中查到数据，还需比较两个时间取最早的时间，防止刷视频的过程中得到feed导致重新刷一遍已经刷过的视频
	var timeCondition time.Time
	if len(videos) == 0 || time.UnixMilli(latestTime).Before(videos[len(videos)-1].CreateTime) {
		timeCondition = time.UnixMilli(latestTime)
	} else {
		timeCondition = videos[len(videos)-1].CreateTime
	}
	// 如果用户未登录，或者feed流中的视频不足一页，通过时间条件在数据库中查询补足一页
	if len(videos) < config.Config.Service.PageSize {
		videos1, err := videoDao.QueryByLatestTimeDESC(timeCondition, config.Config.Service.PageSize-len(videos))
		if err != nil {
			return nil, 0, err
		}
		//拼接到原来数据后面
		videos = append(videos, *videos1...)
	}
	// 转换视频bo
	var videoBOS = make([]bo.Video, len(videos))
	err = entityutil.GetVideoBOS(&videos, &videoBOS)
	if err != nil {
		return nil, 0, err
	}
	if isLogin {
		// 取出部分后，重新存入数据
		tx, err := clearInbox(&inbox, userId)
		if err != nil {
			if tx != nil {
				tx.Rollback()
			}
			return nil, 0, err
		}
		// 提交事务
		tx.Commit()
	}
	// 获得下一次请求的时间
	nextTime := latestTime
	// 是否循环播放视频
	if config.Config.Service.FeedLoop {
		if len(videos) == config.Config.Service.PageSize {
			nextTime = videos[len(videos)-1].CreateTime.UnixMilli()
		} else {
			nextTime = time.Now().UnixMilli()
		}
	} else {
		if len(videos) > 0 && time.UnixMilli(latestTime).After(videos[len(videos)-1].CreateTime) {
			nextTime = videos[len(videos)-1].CreateTime.UnixMilli()
		}
	}

	return videoBOS, nextTime, nil
}

// Publish check token then save upload file to public directory
func (v Video) Publish(c *gin.Context, video *multipart.FileHeader, userId int, title string) error {
	// 检查临时路径是否存在，不存在则创建
	err := checkDIR()
	if err != nil {
		return err
	}
	now := time.Now().Format(config.Config.StandardDate)
	// 将视频保存到本地
	videoName := obsutil.ParseFileName(filepath.Base(video.Filename))
	videoSaveFile := filepath.Join(config.Config.Service.VideoTempDir, now, videoName)
	if err := c.SaveUploadedFile(video, videoSaveFile); err != nil {
		return err
	}
	// 生成封面文件
	coverSaveFile, err := generateCover(videoSaveFile)
	if err != nil {
		return err
	}
	// 数据入库
	var videoDB = daoimpl.NewVideoDaoInstance()
	var tx = daoimpl.Begin()
	var playUrl string
	var coverUrl string
	if config.Config.Server.WithProxy {
		playUrl = config.Config.Server.Protocol + "://" + config.Config.Server.Proxy + "/" + videoSaveFile
		coverUrl = config.Config.Server.Protocol + "://" + config.Config.Server.Proxy + "/" + coverSaveFile
	}
	videoPo := po.Video{
		PlayUrl:       playUrl,
		CoverUrl:      coverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		AuthorId:      userId,
		Title:         title,
	}
	err = videoDB.Insert(tx, &videoPo, true)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = rabbitutil.UploadVideo(videoPo.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = rabbitutil.FeedVideo(videoPo.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (v Video) VideoList(userId int) ([]bo.Video, error) {
	// 查询数据库获取投稿列表
	poVideoList, err := daoimpl.NewVideoDaoInstance().QueryByConditionTimeDESC(&po.Video{AuthorId: userId})
	if err != nil {
		return nil, err
	}
	var boVideoList = make([]bo.Video, len(*poVideoList))
	// po列表转bo
	err = entityutil.GetVideoBOS(poVideoList, &boVideoList)
	if err != nil {
		return nil, err
	}
	return boVideoList, nil
}

// 从自己的收件箱获取
func fromInbox(userId int) ([]bo.Feed, error) {
	var feedBOS = make([]bo.Feed, 0)
	err := redisutil.ZRevRange[bo.Feed](config.Config.Redis.Key.Inbox+strconv.Itoa(userId), &feedBOS)
	if err != nil {
		return nil, err
	}
	// 如果收件箱不存在，则从数据库中查询
	if len(feedBOS) == 0 {
		feedDao := daoimpl.NewFeedDaoInstance()
		feedPOS, err := feedDao.QueryByCondition(&po.Feed{UserId: userId})
		if err != nil {
			return nil, err
		}
		entityutil.GetFeedBOS(&feedPOS, &feedBOS)
	}
	return feedBOS, nil
}

// 从大v的发件箱获取Feed流
// userId 			当前用户id
// @return			当前用户所有关注的大v的发件箱排序后的切片
func fromOutbox(userId int) ([]bo.Feed, error) {
	userDao := daoimpl.NewUserDaoInstance()
	follows, err := userDao.QueryFollows(userId)
	if err != nil {
		return nil, err
	}
	var feeds []bo.Feed
	var length = 0
	for _, user := range *follows {
		if user.FollowerCount >= config.Config.Service.BigVNum {
			// 这里不能初始化分配，初始化分配会导致分配到默认值，而不是空
			var feed []bo.Feed
			err := redisutil.ZRevRange[bo.Feed](config.Config.Redis.Key.Outbox+strconv.Itoa(user.ID), &feed)
			if err != nil {
				return nil, err
			}
			feeds, err = mergeFeeds(&feeds, &feed)
			if err != nil {
				return nil, err
			}
			length += len(feed)
			if length > config.Config.Service.PageSize {
				return feeds, nil
			}
		}
	}
	return feeds, nil
}

// 合并发件箱和收件箱的feed流，得到查询的id集合
// inbox 用户收件箱
// outbox 大v发件箱
// @return id集合
func mergeBox(inbox *[]bo.Feed, outbox *[]bo.Feed, latestTime int64) ([]int, error) {
	var videoIds = make([]int, config.Config.Service.PageSize)
	var inboxIndex, outboxIndex, resIndex = 0, 0, 0

	for inboxIndex < len(*inbox) && outboxIndex < len(*outbox) {
		inboxTime := (*inbox)[inboxIndex].CreateTime
		outboxTime := (*outbox)[outboxIndex].CreateTime
		if inboxTime.After(outboxTime) {
			videoIds[resIndex] = (*inbox)[inboxIndex].VideoId
			resIndex++
			inboxIndex++
		} else {
			// 只对发件箱中上一次latestTime之前的视频进行投送
			if outboxTime.UnixMilli() < latestTime {
				videoIds[resIndex] = (*outbox)[outboxIndex].VideoId
				resIndex++
			}
			outboxIndex++
		}
		// 当结果id集已经满了一页
		if resIndex == config.Config.Service.PageSize {
			*inbox = (*inbox)[inboxIndex:]
			return videoIds, nil
		}
	}
	for inboxIndex < len(*inbox) {
		videoIds[resIndex] = (*inbox)[inboxIndex].VideoId
		resIndex++
		inboxIndex++
		if resIndex == config.Config.Service.PageSize {
			// 去掉没有选中的信息，为后面清理redis中的收件箱做准备
			*inbox = (*inbox)[:inboxIndex]
			return videoIds, nil
		}
	}
	// 去掉没有选中的信息，为后面加锁管理redis中的收件箱做准备
	*inbox = (*inbox)[:inboxIndex]
	for outboxIndex < len(*outbox) {
		outboxTime := (*outbox)[outboxIndex].CreateTime
		// 只对发件箱中上一次latestTime之前的视频进行投送
		if outboxTime.UnixMilli() < latestTime {
			videoIds[resIndex] = (*outbox)[outboxIndex].VideoId
			resIndex++
		}
		outboxIndex++
		if resIndex == config.Config.Service.PageSize {
			return videoIds, nil
		}
	}
	return videoIds, nil
}

// 按顺序合并两个feed切片
func mergeFeeds(feed1 *[]bo.Feed, feed2 *[]bo.Feed) ([]bo.Feed, error) {
	var index1, index2, resIndex = 0, 0, 0
	var feeds = make([]bo.Feed, len(*feed1)+len(*feed2))

	for index1 < len(*feed1) && index2 < len(*feed2) {
		time1 := (*feed1)[index1].CreateTime
		time2 := (*feed2)[index2].CreateTime
		if time1.After(time2) {
			feeds[resIndex] = (*feed1)[index1]
			index1++
		} else {
			feeds[resIndex] = (*feed2)[index2]
			index2++
		}
		resIndex++
	}
	for index1 < len(*feed1) {
		feeds[resIndex] = (*feed1)[index1]
		resIndex++
		index1++
	}
	for index2 < len(*feed2) {
		feeds[resIndex] = (*feed2)[index2]
		resIndex++
		index2++
	}
	return feeds, nil
}

// 清理收件箱，用户查看一次收件箱后，将收件箱中已经查看过的视频清除
// trash 已经查看过的feed对象
func clearInbox(trash *[]bo.Feed, userId int) (*gorm.DB, error) {
	var err error
	// 清理数据库中的数据
	tx := daoimpl.Begin()
	var trashPOS = make([]po.Feed, len(*trash))
	for i, v := range *trash {
		trashPOS[i] = po.Feed{VideoId: v.VideoId, UserId: userId}
	}
	// 如果存在需要删除的feed，则进入数据库进行删除
	if len(*trash) > 0 {
		err = daoimpl.NewFeedDaoInstance().DeleteByCondition(&trashPOS, tx, true)
		if err != nil {
			return tx, err
		}
	}
	// 刷新收件箱的过期时间
	err = redisutil.Expire(config.Config.Redis.Key.Inbox+strconv.Itoa(userId), config.InboxExpireTime)
	if err != nil {
		return tx, nil
	}
	// 删除垃圾数据
	err = redisutil.ZRem[bo.Feed](config.Config.Redis.Key.Inbox+strconv.Itoa(userId),
		trash,
		false,
		nil)
	if err != nil {
		return tx, err
	}
	return tx, err
}

// 生成视频封面
func generateCover(path string) (string, error) {
	// 获得两周前的日期
	now := time.Now().Format(config.Config.StandardDate)
	coverPath := filepath.Join(config.Config.Service.CoverTempDir+now, obsutil.ParseFileName("xxx.jpg"))
	cmd := exec.Command(config.Config.Service.FFMPEGPath+"ffmpeg", "-i", path,
		"-vframes", "1", "-f", "singlejpeg", coverPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println(stderr.String())
		return "", err
	}
	return coverPath, nil
}

// 检查视频和封面的临时目录，不存在则创建
func checkDIR() error {
	// 获得两周前的日期
	now := time.Now().Format(config.Config.StandardDate)
	// 拼接两周前的视频和文件路径
	videoPath, err := filepath.Abs(config.Config.Service.VideoTempDir + now)
	if err != nil {
		return err
	}
	coverPath, err := filepath.Abs(config.Config.Service.CoverTempDir + now)
	if err != nil {
		return err
	}

	// 判断视频路径是否存在
	_, err = os.Stat(videoPath)
	// 视频目录不存在
	if os.IsNotExist(err) {
		err := os.MkdirAll(videoPath, os.ModeDir)
		if err != nil {
			return err
		}
	} else {
		if err != nil {
			return err
		}
	}

	// 判断视频路径是否存在并删除
	_, err = os.Stat(coverPath)
	// 视频目录不存在
	if os.IsNotExist(err) {
		err := os.MkdirAll(coverPath, os.ModeDir)
		if err != nil {
			return err
		}
	} else {
		if err != nil {
			return err
		}
	}

	return nil
}

var (
	video     service.Video
	videoOnce sync.Once
)

func NewVideoServiceInstance() service.Video {
	videoOnce.Do(func() {
		video = Video{}
	})
	return video
}
