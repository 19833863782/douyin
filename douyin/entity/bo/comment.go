// Package bo
// @Author shaofan
// @Date 2022/5/13
package bo

// Comment 评论BO
type Comment struct {
	ID         int    `json:"id"`
	User       User   `json:"user"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
}
