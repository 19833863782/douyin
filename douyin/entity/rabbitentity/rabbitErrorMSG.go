// Package rabbitentity
// @Author shaofan
// @Date 2022/5/31
package rabbitentity

type RabbitErrorMSG struct {
	FeedVideo   []RabbitMSG[int]      `json:"feed_video"`
	UploadVideo []RabbitMSG[int]      `json:"upload_video"`
	Follow      []RabbitMSG[Follow]   `json:"follow"`
	Favorite    []RabbitMSG[Favorite] `json:"favorite"`
}

var ErrorMsgLockChan = make(chan int, 1)
