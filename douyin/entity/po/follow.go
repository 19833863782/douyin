// Package po
// @Author shaofan
// @Date 2022/5/13
package po

// Follow 关注关系PO
type Follow struct {
	RelationModel
	FollowId   int `json:"follow_id" gorm:"follow_id;primaryKey"`
	FollowerId int `json:"follower_id" gorm:"follower_id;primaryKey"`
}

func (Follow) TableName() string {
	return "dy_follow"
}
