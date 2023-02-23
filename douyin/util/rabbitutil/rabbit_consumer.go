// Package rabbitutil
// @Author shaofan
// @Date 2022/5/15
package rabbitutil

import (
	"douyin/config"
	"douyin/entity/bo"
	"douyin/entity/po"
	"douyin/entity/rabbitentity"
	"douyin/repositories/daoimpl"
	"douyin/util/obsutil"
	"douyin/util/redisutil"
	"encoding/json"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
	"sync"
)

// 初始化consumer
func initConsumer() error {
	if err := initServer(); err != nil {
		return err
	}
	if err := followConsumer(); err != nil {
		return err
	}
	if err := uploadVideoConsumer(); err != nil {
		return err
	}
	if err := feedVideoConsumer(); err != nil {
		return err
	}
	if err := favoriteConsumer(); err != nil {
		return err
	}
	return nil
}

// 初始化rabbitmq服务器
func initServer() error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	if err := initFeedVideo(channel); err != nil {
		return err
	}
	if err := initUploadVideo(channel); err != nil {
		return err
	}
	if err := initFollow(channel); err != nil {
		return err
	}
	if err := initFavorite(channel); err != nil {
		return err
	}
	if err := channel.Close(); err != nil {
		return err
	}
	return nil
}

// 点赞消费者
func favoriteConsumer() error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	consume, err := channel.Consume(
		config.Config.Rabbit.Queue.Favorite,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	// 协程处理消费
	go func() {
		for msg := range consume {
			var rabbitMSG rabbitentity.RabbitMSG[rabbitentity.Favorite]
			//反序列化
			err := json.Unmarshal(msg.Body, &rabbitMSG)
			failOnErrorFavorite(err, &rabbitMSG)
			data := rabbitMSG.Data
			err = like(&data)
			failOnErrorFavorite(err, &rabbitMSG)
			err = msg.Ack(true)
			failOnErrorFavorite(err, &rabbitMSG)
		}
	}()
	return nil
}

// 修改关注数量消费
func followConsumer() error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	consume, err := channel.Consume(
		config.Config.Rabbit.Queue.Follow,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	// 协程处理消费
	go func() {
		for msg := range consume {
			var rabbitMSG rabbitentity.RabbitMSG[rabbitentity.Follow]
			//反序列化
			err := json.Unmarshal(msg.Body, &rabbitMSG)
			failOnErrorFollow(err, &rabbitMSG)
			changeFollowNumBody := rabbitMSG.Data
			err = follow(&changeFollowNumBody)
			failOnErrorFollow(err, &rabbitMSG)
			err = msg.Ack(true)
			failOnErrorFollow(err, &rabbitMSG)
		}
	}()
	return nil
}

// 上传视频消费
func uploadVideoConsumer() error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	consume, err := channel.Consume(
		config.Config.Rabbit.Queue.DeadUploadVideo,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	// 协程处理消费
	go func() {
		for msg := range consume {
			var rabbitMSG rabbitentity.RabbitMSG[int]
			//反序列化
			err := json.Unmarshal(msg.Body, &rabbitMSG)
			failOnErrorInt(err, &rabbitMSG)
			//查询video数据
			videoId := rabbitMSG.Data
			err = uploadVideo(videoId)
			failOnErrorInt(err, &rabbitMSG)
			// 确认收到消息
			err = msg.Ack(true)
			failOnErrorInt(err, &rabbitMSG)
		}
	}()
	return nil
}

// 投放视频流消费
func feedVideoConsumer() error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	consume, err := channel.Consume(
		config.Config.Rabbit.Queue.DeadFeedVideo,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	// 协程处理消费
	go func() {
		for msg := range consume {
			var rabbitMSG rabbitentity.RabbitMSG[int]
			//反序列化
			err := json.Unmarshal(msg.Body, &rabbitMSG)
			failOnErrorInt(err, &rabbitMSG)
			//查询video数据
			videoId := rabbitMSG.Data
			err = doFeedVideo(videoId)
			failOnErrorInt(err, &rabbitMSG)
			err = msg.Ack(true)
			failOnErrorInt(err, &rabbitMSG)
		}
	}()
	return nil
}

// 消息补偿机制
func failOnErrorInt(err error, msg *rabbitentity.RabbitMSG[int]) {
	if err != nil {
		msg.ResendCount++
		if int(msg.ResendCount) > config.Config.Rabbit.ResendMax {
			// todo 报警
		}
		handleErrorInt(msg)
		log.Println(err)
	}
}

// 消息补偿机制
func failOnErrorFollow(err error, msg *rabbitentity.RabbitMSG[rabbitentity.Follow]) {
	if err != nil {
		msg.ResendCount++
		if int(msg.ResendCount) > config.Config.Rabbit.ResendMax {
			// todo 报警
		}
		handleErrorFollow(msg)
		log.Println(err)
	}
}

// 消息补偿机制
func failOnErrorFavorite(err error, msg *rabbitentity.RabbitMSG[rabbitentity.Favorite]) {
	if err != nil {
		msg.ResendCount++
		if int(msg.ResendCount) > config.Config.Rabbit.ResendMax {
			// todo 报警
		}
		handleErrorFavorite(msg)
		log.Println(err)
	}
}

// 存储消息补偿信息
func handleErrorInt(msg *rabbitentity.RabbitMSG[int]) {
	// 利用信道加锁
	rabbitentity.ErrorMsgLockChan <- 1
	var rabbitErrorMSG rabbitentity.RabbitErrorMSG
	err := redisutil.Get[rabbitentity.RabbitErrorMSG](config.Config.Redis.Key.ErrorMessage, &rabbitErrorMSG)
	if err != nil {
		log.Println(err)
		<-rabbitentity.ErrorMsgLockChan
		return
	}
	switch msg.Type {
	case rabbitentity.FEED_VIDEO:
		rabbitErrorMSG.FeedVideo = append(rabbitErrorMSG.FeedVideo, *msg)
	case rabbitentity.UPLOAD_VIDEO:
		rabbitErrorMSG.UploadVideo = append(rabbitErrorMSG.UploadVideo, *msg)
	}
	err = redisutil.Set(config.Config.Redis.Key.ErrorMessage, &rabbitErrorMSG)
	if err != nil {
		log.Println(err)
		<-rabbitentity.ErrorMsgLockChan
		return
	}
	<-rabbitentity.ErrorMsgLockChan
}

// 存储消息补偿信息
func handleErrorFollow(msg *rabbitentity.RabbitMSG[rabbitentity.Follow]) {
	// 利用信道加锁
	rabbitentity.ErrorMsgLockChan <- 1
	var rabbitErrorMSG rabbitentity.RabbitErrorMSG
	err := redisutil.Get[rabbitentity.RabbitErrorMSG](config.Config.Redis.Key.ErrorMessage, &rabbitErrorMSG)
	if err != nil {
		log.Println(err)
		<-rabbitentity.ErrorMsgLockChan
		return
	}
	rabbitErrorMSG.Follow = append(rabbitErrorMSG.Follow, *msg)
	err = redisutil.Set(config.Config.Redis.Key.ErrorMessage, &rabbitErrorMSG)
	if err != nil {
		log.Println(err)
		<-rabbitentity.ErrorMsgLockChan
		return
	}
	<-rabbitentity.ErrorMsgLockChan
}

// 存储消息补偿信息
func handleErrorFavorite(msg *rabbitentity.RabbitMSG[rabbitentity.Favorite]) {
	// 利用信道加锁
	rabbitentity.ErrorMsgLockChan <- 1
	var rabbitErrorMSG rabbitentity.RabbitErrorMSG
	err := redisutil.Get[rabbitentity.RabbitErrorMSG](config.Config.Redis.Key.ErrorMessage, &rabbitErrorMSG)
	if err != nil {
		log.Println(err)
		<-rabbitentity.ErrorMsgLockChan
		return
	}
	rabbitErrorMSG.Favorite = append(rabbitErrorMSG.Favorite, *msg)
	err = redisutil.Set(config.Config.Redis.Key.ErrorMessage, &rabbitErrorMSG)
	if err != nil {
		log.Println(err)
		<-rabbitentity.ErrorMsgLockChan
		return
	}
	<-rabbitentity.ErrorMsgLockChan
}

// 处理点赞操作
func like(body *rabbitentity.Favorite) error {
	var err error
	tx := daoimpl.Begin()
	wait := sync.WaitGroup{}
	wait.Add(2)
	if body.IsFavorite {
		err = doLike(body.VideoId, body.UserId, &wait, tx)
	} else {
		err = cancelLike(body.VideoId, body.UserId, &wait, tx)
	}
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 点赞视频
func doLike(videoId, userId int, wait *sync.WaitGroup, tx *gorm.DB) error {
	var err error
	//点赞视频
	go func() {
		defer wait.Done()
		err = daoimpl.NewFavoriteDaoInstance().Insert(&po.Favorite{VideoId: videoId, UserId: userId}, tx, true)
		if err != nil {
			return
		}
	}()
	//增加视频点赞数
	go func() {
		defer wait.Done()
		err = daoimpl.NewVideoDaoInstance().ChangeFavoriteCount(1, videoId, tx, true)
	}()
	wait.Wait()
	if err != nil {
		return err
	}
	return nil
}

// 取消点赞视频
func cancelLike(videoId, userId int, wait *sync.WaitGroup, tx *gorm.DB) error {
	var err error
	//取消点赞视频
	go func() {
		defer wait.Done()
		err = daoimpl.NewFavoriteDaoInstance().DeleteByCondition(&po.Favorite{VideoId: videoId, UserId: userId}, tx, true)
		if err != nil {
			return
		}
	}()
	//减少视频点赞数
	go func() {
		defer wait.Done()
		err = daoimpl.NewVideoDaoInstance().ChangeFavoriteCount(-1, videoId, tx, true)
	}()
	wait.Wait()
	if err != nil {
		return err
	}
	return nil
}

// 关注操作
func follow(body *rabbitentity.Follow) error {
	var tx = daoimpl.Begin()
	var err error
	wg := sync.WaitGroup{}
	wg.Add(3)
	if body.IsFollow {
		err = doFollow(body.UserId, body.ToUserId, &wg, tx)
	} else {
		err = cancelFollow(body.UserId, body.ToUserId, &wg, tx)
	}
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 进行关注
func doFollow(userId, toUserId int, wg *sync.WaitGroup, tx *gorm.DB) error {
	var err error
	go func() {
		defer wg.Done()
		err = daoimpl.NewRelationDaoInstance().Insert(&po.Follow{FollowId: toUserId, FollowerId: userId}, tx, true)
	}()
	go func() {
		defer wg.Done()
		err = daoimpl.NewUserDaoInstance().ChangeFollowCount(userId, 1, tx, true)
	}()
	go func() {
		defer wg.Done()
		err = daoimpl.NewUserDaoInstance().ChangeFansCount(toUserId, 1, tx, true)
	}()
	wg.Wait()
	if err != nil {
		return err
	}
	return nil
}

// 取消关注
func cancelFollow(userId, toUserId int, wg *sync.WaitGroup, tx *gorm.DB) error {
	var err error
	go func() {
		defer wg.Done()
		err = daoimpl.NewRelationDaoInstance().DeleteByCondition(&po.Follow{FollowId: toUserId, FollowerId: userId}, tx, true)
	}()
	go func() {
		defer wg.Done()
		err = daoimpl.NewUserDaoInstance().ChangeFollowCount(userId, -1, tx, true)
	}()
	go func() {
		defer wg.Done()
		err = daoimpl.NewUserDaoInstance().ChangeFansCount(toUserId, -1, tx, true)
	}()
	wg.Wait()
	if err != nil {
		return err
	}
	return nil
}

// 上传视频
func uploadVideo(videoId int) error {
	video, err := daoimpl.NewVideoDaoInstance().QueryById(videoId)
	if err != nil {
		return err
	}
	// 如果是视频不存在，结束处理
	if video.ID == 0 {
		return nil
	}
	var playUrl string
	var coverUrl string
	if config.Config.Server.WithProxy {
		playUrl = video.PlayUrl[strings.Index(video.PlayUrl, config.Config.Server.Proxy)+len(config.Config.Server.Proxy)+1:]
		coverUrl = video.CoverUrl[strings.Index(video.CoverUrl, config.Config.Server.Proxy)+len(config.Config.Server.Proxy)+1:]
	} else {
		playUrl = video.PlayUrl[strings.Index(video.PlayUrl, config.Config.Server.Port)+len(config.Config.Server.Port)+1:]
		coverUrl = video.CoverUrl[strings.Index(video.PlayUrl, config.Config.Server.Port)+len(config.Config.Server.Port)+1:]
	}
	//上传视频
	playUrl, err = obsutil.Upload(playUrl, config.Config.Obs.Buckets.Video)
	video.PlayUrl = playUrl
	if err != nil {
		return err
	}
	// 上传封面
	coverUrl, err = obsutil.Upload(coverUrl, config.Config.Obs.Buckets.Cover)
	if err != nil {
		return err
	}
	video.CoverUrl = coverUrl

	// 更新数据库
	err = daoimpl.NewVideoDaoInstance().UpdateByCondition(video, nil, false)
	if err != nil {
		return err
	}
	return nil
}

// 投放视频流
func doFeedVideo(videoId int) error {
	video, err := daoimpl.NewVideoDaoInstance().QueryById(videoId)
	if err != nil {
		return err
	}
	// 如果是视频不存在，结束处理
	if video.ID == 0 {
		return nil
	}
	// 获取投稿人的信息
	sender, err := daoimpl.NewUserDaoInstance().QueryById(video.AuthorId)
	if err != nil {
		return err
	}
	// 取得需要存入redis的value
	var value = []redis.Z{{Score: float64(video.CreateTime.UnixMilli()), Member: bo.Feed{VideoId: videoId, CreateTime: video.CreateTime}}}
	//大v用户
	if sender.FollowerCount >= config.Config.Service.BigVNum {
		err = redisutil.ZAddWithExpireTime(config.Config.Redis.Key.Outbox+strconv.Itoa(sender.ID),
			value,
			config.OutboxExpireTime,
			false,
			nil)
		if err != nil {
			return err
		}
	} else if sender.FollowerCount > 0 {
		//普通用户
		userIds, err := daoimpl.NewRelationDaoInstance().QueryFansIdByFollowId(video.AuthorId)
		if err != nil {
			return err
		}
		users, err := daoimpl.NewUserDaoInstance().QueryBatchIds(&userIds)
		if err != nil {
			return err
		}
		// 创建redis事务处理
		var pipeline = redisutil.Begin()
		// feed集合，用户持久化
		var feeds = make([]po.Feed, 0)
		for _, user := range *users {
			// 增加到用户的发件箱中,向用户投放不影响用户的收件箱过期时间
			err = redisutil.ZAdd(config.Config.Redis.Key.Inbox+strconv.Itoa(user.ID),
				value,
				true,
				pipeline)
			if err != nil {
				err1 := pipeline.Discard()
				if err1 != nil {
					return err1
				}
				return err
			}
			// 记录需要入库的数据
			feeds = append(feeds, po.Feed{UserId: user.ID, VideoId: videoId})
		}
		// 开始feed持久化，对数据入库
		feedDao := daoimpl.NewFeedDaoInstance()
		tx := daoimpl.Begin()
		err = feedDao.InsertBatch(&feeds, tx, true)
		if err != nil {
			tx.Rollback()
			err1 := pipeline.Discard()
			if err1 != nil {
				return err1
			}
			return err
		}
		_, err = pipeline.Exec()
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
	}
	return nil
}
