package model

import "time"

// Product 定义了产品在数据库中的长相
type Product struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Category  string    `json:"category"`    
	Name      string    `json:"name"`        
	ModelName string    `json:"model_name"`  
	// 加入拖动排序
	SortOrder int    `json:"sort_order" gorm:"default:0"`
	MainImage string    `json:"main_image"`  
	FileUrl   string    `json:"file_url"`    
	
	// 首页展示控制开关
	IsBanner   bool     `json:"is_banner"`  
	IsFeatured bool     `json:"is_featured"` 

	// 在后台全部压缩成一个 JSON 字符串
	DetailData string   `gorm:"type:text" json:"detail_data"`

	// 自动管理的时间戳
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Category 行业分类数据模型
type Category struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Name             string    `json:"name"`              // 行业名称 
	ShortDescription string    `json:"short_description"` // 主页短文案
	BannerSubtitle   string    `json:"banner_subtitle"`   // 列表页长文案
	BannerImage      string    `json:"banner_image"`      // 行业专属 Banner 图路径
	SortOrder        int       `json:"sort_order" gorm:"default:0"` // 拖拽排序权重
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}