package model

import "time"

type Banner struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	MediaType  string    `json:"media_type"` // "video" 或 "image"
	MediaUrl   string    `json:"media_url"`  
	LinkUrl    string    `json:"link_url"`   
	
	// 让文字也全部从数据库读取
	MainTitle  string    `json:"main_title"` 
	SubTitle   string    `json:"sub_title"`  
	
	IsActive   bool      `json:"is_active" gorm:"default:true"`
	SortOrder  int       `json:"sort_order" gorm:"default:0"`
	CreatedAt  time.Time `json:"created_at"`
}