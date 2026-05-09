package model

import "time"

type Application struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Anchor     string    `json:"anchor" gorm:"type:varchar(50)"`     // HTML中的ID锚点，如 EnergyStorage
	Top        string    `json:"top" gorm:"type:varchar(255)"`      // 顶部小字
	Title      string    `json:"title" gorm:"type:varchar(255)"`    // 主标题
	Subtitle   string    `json:"subtitle" gorm:"type:text"`         // 描述（支持HTML换行）
	Link       string    `json:"link" gorm:"type:varchar(255)"`     // 跳转链接
	Background string    `json:"background" gorm:"type:varchar(255)"` // 背景大图路径
	Image      string    `json:"image" gorm:"type:varchar(255)"`      // 侧边产品图路径
	Order      int       `json:"order" gorm:"type:int;default:0"`    // 排序
	UpdatedAt  time.Time `json:"updated_at"`
}