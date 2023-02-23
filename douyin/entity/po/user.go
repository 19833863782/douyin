// Package po
// @Author shaofan
// @Date 2022/5/13
package po

// User 用户PO
type User struct {
	EntityModel
	Name          string `json:"name" gorm:"name;not null"`
	Password      string `json:"password" gorm:"password;not null"`
	FollowCount   int    `json:"follow_count" gorm:"follow_count;not null"`
	FollowerCount int    `json:"follower_count" gorm:"follower_count;not null"`
}

func (User) TableName() string {
	return "dy_user"
}
