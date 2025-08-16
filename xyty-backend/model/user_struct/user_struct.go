package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	// 移除 gorm.Model 嵌入，显式定义所有字段
	ID        uint       `json:"id" gorm:"primarykey"`    // 主键ID
	CreatedAt time.Time  `json:"created_at"`              // 创建时间
	UpdatedAt time.Time  `json:"updated_at"`              // 更新时间
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"` // 软删除标记（*time.Time 兼容 GORM 软删除）

	Username string `json:"username" gorm:"primarykey" gorm:"size:30"`  // 用户名（主键）
	Password string `json:"password" gorm:"size:30"`                    // 密码
	MailAddr string `json:"mail_addr" gorm:"primarykey" gorm:"size:30"` // 邮箱地址（主键）
	Avatar   string `json:"avatar" gorm:"size:255"`                     // 头像字段
}

type UserCode struct {
	gorm.Model
	Username string `json:"username" gorm:"primarykey" gorm:"size:30"`
	MailAddr string `json:"mail_addr" gorm:"primarykey"`
	Code     string `json:"code" gorm:"primarykey"`
}

type VideoRecord struct {
	ID        uint      `json:"id" gorm:"primarykey"`    // 主键ID
	CreatedAt time.Time `json:"created_at"`              // 创建时间
	UpdatedAt time.Time `json:"updated_at"`              // 更新时间
	DeletedAt time.Time `json:"deleted_at" gorm:"index"` // 软删除标记

	Username  string `json:"username" gorm:"index"` // 用户ID
	VideoID   string `json:"video_id"`              // 视频ID
	VideoName string `json:"video_name"`            // 视频名称
	WatchTime int64  `json:"watch_time"`            // 观看时长(秒)
	Progress  int    `json:"progress"`              // 观看进度(百分比)
	WatchDate string `json:"watch_date"`            // 观看日期
}

// ScenarioExperience 情景体验记录结构体
// @Description 存储用户的情景体验记录信息
// @Field ID uint 主键ID
// @Field CreatedAt time.Time 创建时间
// @Field UpdatedAt time.Time 更新时间
// @Field DeletedAt *time.Time 删除时间(软删除标记，可为空)
// @Field Username string 用户ID
// @Field ScenarioID string 情景ID
// @Field ScenarioName string 情景名称
// @Field ExperienceDate string 体验日期
// @Field Score int 评分(可选)
type ScenarioExperience struct {
	// 移除 gorm.Model 嵌入，显式定义所有字段
	ID        uint       `json:"id" gorm:"primarykey"`    // 主键ID
	CreatedAt time.Time  `json:"created_at"`              // 创建时间
	UpdatedAt time.Time  `json:"updated_at"`              // 更新时间
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"` // 软删除标记（*time.Time 兼容 GORM 软删除）

	Username       string `json:"username" gorm:"index"` // 用户ID
	ScenarioID     string `json:"scenario_id"`           // 情景ID
	ScenarioName   string `json:"scenario_name"`         // 情景名称
	ExperienceDate string `json:"experience_date"`       // 体验日期
	Score          int    `json:"score"`                 // 评分(可选)
}
