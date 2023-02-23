// Package po
// @Author shaofan
// @Date 2022/5/13
package po

// Feed 视频流Po
type Feed struct {
	RelationModel
	UserId  int `json:"user_id" gorm:"user_id;primaryKey"`
	VideoId int `json:"video_id" gorm:"video_id;primaryKey"`
}

func (Feed) TableName() string {
	return "dy_feed"
}
