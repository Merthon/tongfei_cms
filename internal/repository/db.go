package repository

import (
	"fmt"
	"log"
	"tonfy_CMS/internal/model" 

	"github.com/glebarez/sqlite"
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
	err = DB.AutoMigrate(&model.News{}, &model.Product{}, &model.Category{}, &model.Job{}, &model.JobApplication{}, &model.Banner{}, &model.ContactMessage{}, &model.AdminUser{}, &model.ServicePage{})
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
	// 初始化 ServicePage 默认数据
    var serviceCount int64
    DB.Model(&model.ServicePage{}).Count(&serviceCount)
    if serviceCount == 0 {
        defaultCards := `[
            {"title":"Fast & Reliable Response","items":["Dedicated Customer Interface & Accountable Personnel","Standardized Internal Collaboration Processes","Defined Response Timeframes"]},
            {"title":"Dependable Delivery","items":["Multi-Regional Manufacturing & Capacity Network","Integrated Supply Chain Coordination","Project-Oriented Production Management"]},
            {"title":"Expert Engineering Support","items":["Dedicated Application & Systems Engineering Teams","Mature Design & Validation Tools","Deep Expertise Across Multiple Industries"]},
            {"title":"Worldwide & Localized Support","items":["Dedicated Local Teams & Authorized Partners","Multi-Timezone & Multi-Language Assistance","Global Spare Parts & Logistics"]},
            {"title":"Tailored Solutions","items":["Modular Product & Technology Platforms","Integrated Engineering & R&D Collaboration","Fast Prototyping & Solution"]},
            {"title":"Sustained Reliability","items":["Robust Quality Management System","End-to-End Lifecycle Testing & Validation","Continuous Improvement with Closed-Loop Feedback"]}
        ]`

        defaultPage := model.ServicePage{
            ID:         1,
            Stat1Num:   "49",
            Stat1Label: "Global Service Centers",
            Stat2Num:   "17",
            Stat2Label: "Countries & Regions Coverage",
            TeamText1:  "TONFY is powered by a highly experienced service team...",
            TeamText2:  "With strong engineering expertise, standardized service processes...",
            TeamImage:  "/assets/images/service/sec01.webp",
            CardsJSON:  defaultCards,
        }
        DB.Create(&defaultPage)
        fmt.Println("🚀 已自动生成服务支持页默认配置！")
    }
}