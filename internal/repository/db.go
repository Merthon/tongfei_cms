package repository

import (
	"log"
	"tonfy_CMS/internal/model" // 注意替换成你自己的 go mod 模块名

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
	err = DB.AutoMigrate(&model.News{}, &model.Product{}, &model.Category{}, &model.Job{}, &model.JobApplication{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("数据库连接并且迁移成功！")
}