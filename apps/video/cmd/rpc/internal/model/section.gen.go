// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"gorm.io/gorm"
)

const TableNameSection = "section"

// Section mapped from table <section>
type Section struct {
	ID         int64          `gorm:"column:id;primaryKey;autoIncrement:true;comment:分区id" json:"id"`              // 分区id
	Name       string         `gorm:"column:name;not null;comment:分区名" json:"name"`                                // 分区名
	CreateTime int64          `gorm:"column:create_time;type:int;type:unsigned;autoCreateTime" json:"create_time"` // 创建时间
	UpdateTime int64          `gorm:"column:update_time;type:int;type:unsigned;autoUpdateTime" json:"update_time"` // 更新时间
	DeleteTime gorm.DeletedAt `gorm:"column:delete_time;comment:删除时间" json:"delete_time"`                          // 删除时间
}

// TableName Section's table name
func (*Section) TableName() string {
	return TableNameSection
}