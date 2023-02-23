// Package rabbitentity
// @Author shaofan
// @Date 2022/5/14
package rabbitentity

const (
	FOLLOW       = '1'
	UPLOAD_VIDEO = '2'
	FEED_VIDEO   = '3'
	FAVORITE     = '4'
)

// RabbitMSG 消息队列标准消息体
type RabbitMSG[T RabbitType] struct {
	Type byte `json:"type"`
	// 消息体
	Data T `json:"data"`
	// 重发次数
	ResendCount uint8 `json:"resend_count"`
}

type RabbitType interface {
	int | Follow | Favorite
}
