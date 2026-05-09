package model

import "time"

// ServicePage
type ServicePage struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	
	// 第一屏
	Stat1Num    string    `json:"stat1_num" gorm:"type:varchar(50)"`
	Stat1Label  string    `json:"stat1_label" gorm:"type:varchar(255)"`
	Stat2Num    string    `json:"stat2_num" gorm:"type:varchar(50)"`
	Stat2Label  string    `json:"stat2_label" gorm:"type:varchar(255)"`

	// 第二屏
	TeamText1   string    `json:"team_text1" gorm:"type:text"`
	TeamText2   string    `json:"team_text2" gorm:"type:text"`
	TeamImage   string    `json:"team_image" gorm:"type:varchar(255)"` // 存 /uploads/images/...

	// 第四屏
	CardsJSON   string    `json:"cards_json" gorm:"type:text"`

	UpdatedAt   time.Time `json:"updated_at"`
}