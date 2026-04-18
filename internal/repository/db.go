package repository

import (
	"fmt"
	"log"
	"tonfy_CMS/internal/model" 

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// 初始化数据库
func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("cms.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 迁移
	err = DB.AutoMigrate(&model.News{}, &model.Product{}, &model.Category{}, &model.Job{}, &model.JobApplication{}, &model.Banner{}, &model.ContactMessage{}, &model.AdminUser{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("数据库连接并且迁移成功！")

	var count int64
    DB.Model(&model.AdminUser{}).Count(&count)
    if count == 0 {
        boss := model.AdminUser{
            Username: "admin",
            Password: "123456", // 保持你之前的测试密码
            Role:     "super_admin",
            Modules:  "all",
        }
        DB.Create(&boss)
        fmt.Println("🚀 已自动生成超级管理员账号: admin / 123456")
    }
}