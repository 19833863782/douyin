// Package job
// @Author shaofan
// @Date 2022/5/14
package job

import (
	"douyin/config"
	"douyin/entity/bo"
	"douyin/entity/rabbitentity"
	"douyin/util/rabbitutil"
	"douyin/util/redisutil"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"path/filepath"
	"time"
)

// StartJob 开启任务调度
func StartJob() {
	clearOutBoxJob()
	clearLocalVideoJob()
	handleErrorMSGJob()
}

// 返回一个支持至 秒 级别的 cron
func newWithSeconds() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}

// 清理用户发件箱
func clearOutBoxJob() {
	c := newWithSeconds()
	_, err := c.AddFunc("0 0 0 * * *", clearOutBox)
	if err != nil {
		log.Fatalln(err)
	}
	c.Start()
}

// 处理错误消息
func handleErrorMSGJob() {
	c := newWithSeconds()
	_, err := c.AddFunc("@every 10s", handleErrorMSG)
	if err != nil {
		log.Fatalln(err)
		return
	}
	c.Start()
}

// 清理本地临时存储的视频
func clearLocalVideoJob() {
	c := newWithSeconds()
	_, err := c.AddFunc("0 0 1 * * *", clearLocalVideo)
	if err != nil {
		log.Fatalln(err)
		return
	}
	c.Start()
}

// 清理用户发件箱
func clearOutBox() {
	var outBoxes = make([]string, 0)
	err := redisutil.Keys(config.Config.Redis.Key.Outbox, &outBoxes)
	if err != nil {
		log.Println(err)
		return
	}
	for _, key := range outBoxes {
		var feeds = make([]bo.Feed, 0)
		if err != nil {
			log.Println(err)
			return
		}
		// 获取发件箱中的数据
		err = redisutil.ZRevRange[bo.Feed](key, &feeds)
		if err != nil {
			log.Println(err)
			return
		}
		var index int
		var feed bo.Feed
		// 对发件箱中的数据进行排查
		for index, feed = range feeds {
			// 如果发布时间增加七天大于当前，则退出循环
			if feed.CreateTime.AddDate(0, 0, 7).After(time.Now()) {
				break
			}
		}
		// 根据index来获取需要清理的数据
		feeds = feeds[:index]
		err = redisutil.ZRem[bo.Feed](key, &feeds, false, nil)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// 处理错误消息
func handleErrorMSG() {
	// 利用信道加锁
	rabbitentity.ErrorMsgLockChan <- 1
	p := redisutil.Begin()
	var msgS rabbitentity.RabbitErrorMSG
	err := redisutil.GetAndDelete[rabbitentity.RabbitErrorMSG](config.Config.Redis.Key.ErrorMessage, &msgS, true, p)
	<-rabbitentity.ErrorMsgLockChan
	if err != nil {
		log.Println(err)
		err := p.Discard()
		if err != nil {
			log.Println(err)
			return
		}
		return
	}
	for _, msg := range msgS.Favorite {
		err := rabbitutil.Publish[rabbitentity.Favorite](&msg,
			config.Config.Rabbit.Exchange.ServiceExchange,
			config.Config.Rabbit.Key.Follow,
		)
		if err != nil {
			log.Println(err)
			err := p.Discard()
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
	}
	for _, msg := range msgS.Follow {
		err := rabbitutil.Publish[rabbitentity.Follow](&msg,
			config.Config.Rabbit.Exchange.ServiceExchange,
			config.Config.Rabbit.Key.Follow,
		)
		if err != nil {
			log.Println(err)
			err := p.Discard()
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
	}
	for _, msg := range msgS.FeedVideo {
		err := rabbitutil.Publish[int](&msg,
			config.Config.Rabbit.Exchange.ServiceExchange,
			config.Config.Rabbit.Key.FeedVideo)
		if err != nil {
			log.Println(err)
			err := p.Discard()
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
	}
	for _, msg := range msgS.UploadVideo {
		err := rabbitutil.Publish[int](&msg,
			config.Config.Rabbit.Exchange.ServiceExchange,
			config.Config.Rabbit.Key.UploadVideo)
		if err != nil {
			log.Println(err)
			err := p.Discard()
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
	}
	_, err = p.Exec()
	if err != nil {
		log.Println(err)
		return
	}

}

// 清理本地临时存储的视频
func clearLocalVideo() {
	// 获得两周前的日期
	now := time.Now().AddDate(0, 0, -7).Format(config.Config.StandardDate)
	// 拼接两周前的视频和文件路径
	videoPath, err := filepath.Abs(filepath.Join(config.Config.Service.VideoTempDir, now))
	if err != nil {
		log.Println(err)
		return
	}
	coverPath, err := filepath.Abs(filepath.Join(config.Config.Service.CoverTempDir, now))
	if err != nil {
		log.Println(err)
		return
	}
	// 判断视频路径是否存在并删除
	_, err = os.Stat(videoPath)
	if err != nil {
		log.Println(err)
		return
	}
	err = os.RemoveAll(videoPath)
	if err != nil {
		log.Println(err)
		return
	}
	// 判断封面路径是否存在并删除
	_, err = os.Stat(coverPath)
	if err != nil {
		log.Println(err)
		return
	}
	err = os.RemoveAll(coverPath)
	if err != nil {
		log.Println(err)
		return
	}
}
