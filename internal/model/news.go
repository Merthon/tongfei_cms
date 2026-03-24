package model

import (
	"gorm.io/gorm"
)

// News 代表了数据库中的新闻表
type News struct {
	gorm.Model           // GORM 默认提供 ID, CreatedAt, UpdatedAt, DeletedAt
	Title       string `gorm:"type:varchar(255);not null" json:"title"`
	Date        string `gorm:"type:varchar(50)" json:"date"`               // 对应你 JSON 里的 "2025-8-15"
	Description string `gorm:"type:text" json:"description"`               // 副标题/简述
	Image       string `gorm:"type:varchar(255)" json:"image"`             // 封面图路径
	Content     string `gorm:"type:text" json:"content"`                   // 核心：存放那段 HTML 正文代码
}