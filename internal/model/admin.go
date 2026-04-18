package model

import "gorm.io/gorm"

// 后台管理员账号表
type AdminUser struct {
    gorm.Model
    Username string `json:"username" gorm:"unique;not null"`
    Password string `json:"-"` // 密码绝对不能通过 JSON 泄露给前端
    
    // 权限核心系统
    Role     string `json:"role" gorm:"type:varchar(20);default:'editor'"` 
    Modules  string `json:"modules" gorm:"type:varchar(255)"`             
}