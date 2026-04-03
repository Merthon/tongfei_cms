package model

import "time"

type Job struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Title            string    `json:"title"`
	Department       string    `json:"department"`
	Location         string    `json:"location"`
	Responsibilities string    `gorm:"type:text" json:"responsibilities"`
	IsActive         bool      `json:"is_active" gorm:"default:true"`
	SortOrder        int       `json:"sort_order" gorm:"default:0"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type JobApplication struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Position      string    `json:"position"`
	ApplicantName string    `json:"applicant_name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	CoverLetter   string    `gorm:"type:text" json:"cover_letter"`
	ResumeFileUrl string    `json:"resume_file_url"`
	Status        string    `json:"status" gorm:"default:'Unread'"`
	CreatedAt     time.Time `json:"created_at"`
}