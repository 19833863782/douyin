// Package po
// @Author shaofan
// @Date 2022/5/14
package po

import "time"

// EntityModel 实体模型
type EntityModel struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	CreateTime time.Time `json:"create_time" gorm:"create_time;not null"`
	UpdateTime time.Time `json:"update_time" gorm:"update_time;not null"`
}

// RelationModel 关系模型
type RelationModel struct {
	CreateTime time.Time `json:"create_time" gorm:"create_time;not null"`
	UpdateTime time.Time `json:"update_time" gorm:"update_time;not null"`
}
