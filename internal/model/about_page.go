package model

import "time"

type AboutPage struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	
	// 1. 首屏 Banner
	BannerImage        string    `json:"banner_image" gorm:"type:varchar(255)"`
	BannerTitle        string    `json:"banner_title" gorm:"type:text"` // 存带 <br> 的标题

	// 2. 公司简介
	IntroVideo         string    `json:"intro_video" gorm:"type:varchar(255)"`
	IntroTextsJSON     string    `json:"intro_texts_json" gorm:"type:text"` // 存储多段介绍文案的 JSON 数组

	// 3. ESG 战略 (3个支柱 + 4个底卡)
	ESGPillarsJSON     string    `json:"esg_pillars_json" gorm:"type:text"`
	ESGBottomCardsJSON string    `json:"esg_bottom_cards_json" gorm:"type:text"`

	// 4. 企业文化 (3个卡片)
	CultureItemsJSON   string    `json:"culture_items_json" gorm:"type:text"`

	// 5. 全球足迹 (1个大卡 + 3个小卡)
	FootprintMainJSON  string    `json:"footprint_main_json" gorm:"type:text"`
	FootprintSubsJSON  string    `json:"footprint_subs_json" gorm:"type:text"`

	UpdatedAt          time.Time `json:"updated_at"`
}