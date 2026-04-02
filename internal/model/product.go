package model

import "time"

// Product 定义了产品在数据库中的长相
type Product struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Category  string    `json:"category"`    // 所属行业，例如 "Machinery Manufacturing"
	Name      string    `json:"name"`        // 产品名称，例如 "AIR/WATER HEAT EXCHANGER"
	ModelName string    `json:"model_name"`  // 产品型号/目录名，例如 "air_water-heat_exchanger_mwa" (用于匹配你前端以前的路径)
	// 3.31加入拖动排序
	SortOrder int    `json:"sort_order" gorm:"default:0"`
	MainImage string    `json:"main_image"`  // 产品主图的相对路径
	FileUrl   string    `json:"file_url"`    // 产品说明书 ZIP 的相对路径
	
	// 首页展示控制开关
	IsBanner   bool     `json:"is_banner"`   // 是否出现在首页最顶部的轮播图中
	IsFeatured bool     `json:"is_featured"` // 是否作为该行业的主推产品（显示在首页右侧的3个坑位）

	// 在后台全部压缩成一个 JSON 字符串存进这个文本字段里！前端要的时候原样吐出去。
	DetailData string   `gorm:"type:text" json:"detail_data"`

	// 自动管理的时间戳
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}