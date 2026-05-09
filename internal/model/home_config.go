package model

import "time"

// HomeConfig 主页核心数据统计配置
type HomeConfig struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	
	// Top Row (第一排)
	EmpCount        string    `json:"emp_count"`        // 2500 (总员工)
	RdCount         string    `json:"rd_count"`         // 300  (研发人员)
	LabCount        string    `json:"lab_count"`        // 10   (实验室)

	// Middle Row (第二排)
	ValidationHours string    `json:"validation_hours"` // 15   (百亿小时验证)
	CountriesCount  string    `json:"countries_count"`  // 100  (覆盖国家)

	// Bottom Row (第三排文本)
	BottomText1     string    `json:"bottom_text_1"`    // 50+ Industries Empowered
	BottomText2     string    `json:"bottom_text_2"`    // 2,000+ Global Customers Trust
	BottomText3     string    `json:"bottom_text_3"`    // ±0.1℃ Precision Engineered

	UpdatedAt       time.Time `json:"updated_at"`
}