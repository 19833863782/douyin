// Package rabbitutil
// @Author shaofan
// @Date 2022/5/15
package rabbitutil

import (
	"douyin/config"
	"github.com/streadway/amqp"
	"log"
)

var (
	conn *amqp.Connection
)

// Init rabbitmq初始化
func Init() {
	// 建立连接
	var err error
	conn, err = amqp.Dial(config.Config.Rabbit.Url)
	if err != nil {
		log.Fatalln(err)
	}
	// 开始监听消费
	if err := initConsumer(); err != nil {
		log.Fatalln(err)
	}
}

// 声明交换机
func initExchange(exchange string, channel *amqp.Channel) error {
	err := channel.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		return err
	}
	return nil
}

// 声明队列
func initQueue(queue string, args amqp.Table, channel *amqp.Channel) error {
	_, err := channel.QueueDeclare(queue, true, false, false, false, args)
	if err != nil {
		return err
	}
	return nil
}

// 声明绑定
func initBinding(exchange, queue, key string, channel *amqp.Channel) error {
	err := channel.QueueBind(queue, key, exchange, false, nil)
	if err != nil {
		return err
	}
	return nil
}

// 声明整套消息路径
func producerInit(exchange, queue, key string, args amqp.Table, channel *amqp.Channel) error {
	// 声明交换机
	err := initExchange(exchange, channel)
	if err != nil {
		return err
	}
	// 声明队列
	err = initQueue(queue, args, channel)
	if err != nil {
		return err
	}
	// 声明绑定
	err = initBinding(
		exchange,
		queue,
		key,
		channel,
	)
	if err != nil {
		return err
	}
	return nil
}

// 初始化投放视频流消息队列
func initFeedVideo(channel *amqp.Channel) error {
	if err := producerInit(
		config.Config.Rabbit.Exchange.DeadServiceExchange,
		config.Config.Rabbit.Queue.DeadFeedVideo,
		config.Config.Rabbit.Key.FeedVideo,
		nil,
		channel,
	); err != nil {
		return err
	}
	// 声明服务端
	if err := producerInit(
		config.Config.Rabbit.Exchange.ServiceExchange,
		config.Config.Rabbit.Queue.FeedVideo,
		config.Config.Rabbit.Key.FeedVideo,
		amqp.Table{
			"x-message-ttl":             config.Config.Rabbit.TTL.FeedVideo,
			"x-dead-letter-exchange":    config.Config.Rabbit.Exchange.DeadServiceExchange,
			"x-dead-letter-routing-key": config.Config.Rabbit.Key.FeedVideo,
		},
		channel,
	); err != nil {
		return err
	}
	return nil
}

// 初始化上传视频消息队列
func initUploadVideo(channel *amqp.Channel) error {
	if err := producerInit(
		config.Config.Rabbit.Exchange.DeadServiceExchange,
		config.Config.Rabbit.Queue.DeadUploadVideo,
		config.Config.Rabbit.Key.UploadVideo,
		nil,
		channel,
	); err != nil {
		return err
	}
	// 声明服务端
	if err := producerInit(
		config.Config.Rabbit.Exchange.ServiceExchange,
		config.Config.Rabbit.Queue.UploadVideo,
		config.Config.Rabbit.Key.UploadVideo,
		amqp.Table{
			"x-message-ttl":             config.Config.Rabbit.TTL.UploadVideo,
			"x-dead-letter-exchange":    config.Config.Rabbit.Exchange.DeadServiceExchange,
			"x-dead-letter-routing-key": config.Config.Rabbit.Key.UploadVideo,
		},
		channel,
	); err != nil {
		return err
	}
	return nil
}

// 初始化关注处理的消息队列
func initFollow(channel *amqp.Channel) error {
	// 声明服务端
	if err := producerInit(
		config.Config.Rabbit.Exchange.ServiceExchange,
		config.Config.Rabbit.Queue.Follow,
		config.Config.Rabbit.Key.Follow,
		nil,
		channel,
	); err != nil {
		return err
	}
	return nil
}

// 初始化点赞处理的消息队列
func initFavorite(channel *amqp.Channel) error {
	if err := producerInit(
		config.Config.Rabbit.Exchange.ServiceExchange,
		config.Config.Rabbit.Queue.Favorite,
		config.Config.Rabbit.Key.Favorite,
		nil,
		channel,
	); err != nil {
		return err
	}
	return nil
}
