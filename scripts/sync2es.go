package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Video mapped from table <video>
type Video struct {
	ID         int64          `gorm:"column:id;primaryKey;autoIncrement:true;comment:视频id" json:"id"` // 视频id
	Title      string         `gorm:"column:title;not null;comment:视频标题" json:"title"`                // 视频标题
	SectionID  int64          `gorm:"column:section_id;not null;comment:视频分区id" json:"section_id"`    // 视频分区id
	TagIds     string         `gorm:"column:tag_ids;not null" json:"tag_ids"`
	OwnerID    int64          `gorm:"column:owner_id;not null;comment:视频发布者id" json:"owner_id"`                    // 视频发布者id
	PlayURL    string         `gorm:"column:play_url;not null;comment:视频下载地址" json:"play_url"`                     // 视频下载地址
	CoverURL   string         `gorm:"column:cover_url;not null;comment:封面下载地址" json:"cover_url"`                   // 封面下载地址
	CreateTime int64          `gorm:"column:create_time;type:int;type:unsigned;autoCreateTime" json:"create_time"` // 创建时间
	UpdateTime int64          `gorm:"column:update_time;type:int;type:unsigned;autoUpdateTime" json:"update_time"` // 更新时间
	DeleteTime gorm.DeletedAt `gorm:"column:delete_time;comment:删除时间" json:"delete_time"`                          // 删除时间
}

// TableName Video's table name
func (*Video) TableName() string {
	return "video"
}

// User mapped from table <user>
type User struct {
	ID         int64          `gorm:"column:id;primaryKey;autoIncrement:true;comment:用户id" json:"id"`              // 用户id
	Username   string         `gorm:"column:username;not null;comment:用户名" json:"username"`                        // 用户名
	Password   string         `gorm:"column:password;not null;comment:密码" json:"password"`                         // 密码
	CreateTime int64          `gorm:"column:create_time;type:int;type:unsigned;autoCreateTime" json:"create_time"` // 创建时间
	UpdateTime int64          `gorm:"column:update_time;type:int;type:unsigned;autoUpdateTime" json:"update_time"` // 更新时间
	DeleteTime gorm.DeletedAt `gorm:"column:delete_time;comment:删除时间" json:"delete_time"`                          // 删除时间
}

// TableName User's table name
func (*User) TableName() string {
	return "user"
}
func main() {
	// 连接视频信息所在数据库
	dsn := "root:my-secret-pw@tcp(127.0.0.1:3306)/go_zero_douyin?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// 查询所有视频信息
	var videos []Video
	videoResult := db.Find(&videos)
	if videoResult.Error != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}

	videoClient := resty.New()
	// 遍历视频信息，创建es文档
	for _, video := range videos {
		body := fmt.Sprintf(`{"video_id":%d}`, video.ID)
		_, err := videoClient.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjU4MDQ2NTksImlhdCI6MTY5NDI2ODY1OSwiand0VXNlcklkIjo0fQ.NHtNwiMk7t6vXOACsyyole_fcMwt5vgmXmIoKlmFKWE").
			SetBody(body).
			Post("http://127.0.0.1:1003/video/v1/sync")
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("同步视频信息到es成功")

	// 查询所有用户信息
	var users []User
	userResult := db.Find(&users)
	if userResult.Error != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}

	userClient := resty.New()
	// 遍历用户信息，创建es文档
	for _, user := range users {
		body := fmt.Sprintf(`{"user_id":%d}`, user.ID)
		_, err := userClient.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjU4MDQ2NTksImlhdCI6MTY5NDI2ODY1OSwiand0VXNlcklkIjo0fQ.NHtNwiMk7t6vXOACsyyole_fcMwt5vgmXmIoKlmFKWE").
			SetBody(body).
			Post("http://127.0.0.1:1002/user/v1/sync")
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("同步用户信息到es成功")
}
