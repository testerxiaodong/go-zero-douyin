// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameVideo = "video"

// Video mapped from table <video>
type Video struct {
	ID         int64          `gorm:"column:id;primaryKey;autoIncrement:true;comment:视频id" json:"id"`              // 视频id
	Title      string         `gorm:"column:title;not null;comment:视频标题" json:"title"`                             // 视频标题
	OwnerID    int64          `gorm:"column:owner_id;not null;comment:视频发布者的用户id" json:"owner_id"`                 // 视频发布者的用户id
	PlayURL    string         `gorm:"column:play_url;not null;comment:视频下载地址" json:"play_url"`                     // 视频下载地址
	CoverURL   string         `gorm:"column:cover_url;not null;comment:视频封面地址" json:"cover_url"`                   // 视频封面地址
	CreateTime *time.Time     `gorm:"column:create_time;type:int;type:unsigned;autoCreateTime" json:"create_time"` // 创建时间
	UpdateTime *time.Time     `gorm:"column:update_time;type:int;type:unsigned;autoUpdateTime" json:"update_time"` // 更新时间
	DeleteTime gorm.DeletedAt `gorm:"column:delete_time;comment:删除时间" json:"delete_time"`                          // 删除时间
}

// TableName Video's table name
func (*Video) TableName() string {
	return TableNameVideo
}
