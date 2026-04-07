package model

import "time"

// ContactMessage 客户留言表
type ContactMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Company   string    `json:"company"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`    // 对应前端的 tel
	Email     string    `json:"email"`
	City      string    `json:"city"`     
	Industry  string    `json:"industry"` 
	Message   string    `gorm:"type:text" json:"message"` // 对应前端的 content
	Status    string    `json:"status" gorm:"default:'Unread'"`
	CreatedAt time.Time `json:"created_at"`
}