// Package rabbitentity
// @Author shaofan
// @Date 2022/6/5
package rabbitentity

type Favorite struct {
	UserId     int  `json:"user_id"`
	VideoId    int  `json:"video_id"`
	IsFavorite bool `json:"is_favorite"`
}
