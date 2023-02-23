// Package po
// @Author shaofan
// @Date 2022/5/13
package po

// Video 视频PO
type Video struct {
	EntityModel
	PlayUrl       string `json:"play_url" gorm:"play_url;not null"`
	CoverUrl      string `json:"cover_url" gorm:"cover_url;not null"`
	FavoriteCount int    `json:"favorite_count" gorm:"favorite_count;not null"`
	CommentCount  int    `json:"comment_count" gorm:"comment_count;not null"`
	AuthorId      int    `json:"author_id" gorm:"author_id;not null"`
	Title         string `json:"title" gorm:"title not null"`
}

func (Video) TableName() string {
	return "dy_video"
}
