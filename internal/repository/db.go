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
	err = DB.AutoMigrate(&model.News{}, &model.Product{}, &model.Category{}, &model.Job{}, &model.JobApplication{}, &model.Banner{}, &model.ContactMessage{}, &model.AdminUser{}, &model.ServicePage{}, &model.Application{}, &model.AboutPage{})
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
	// 初始化 Application 默认数据
	var applicationCount int64
	DB.Model(&model.Application{}).Count(&applicationCount)
	if applicationCount == 0 {
		defaultApplication := model.Application{
			Anchor:     "energy-storage",
			Top:        "Energy Storage Solutions",
			Title:      "Advanced Energy Storage Systems",
			Subtitle:   "Reliable and efficient energy storage for modern applications.",
			Link:       "/energy-storage",
			Background: "/assets/images/application/bg_energy_storage.webp",
			Image:      "/assets/images/application/img_energy_storage.webp",
			Order:      1,
		}
		DB.Create(&defaultApplication)
		fmt.Println("🚀 已自动生成 Application 默认配置！")
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

	var appCount int64
	DB.Model(&model.Application{}).Count(&appCount)
	// ================= 强制刷新应用场景数据 =================
    // 1. 强制清空旧的残留数据（极其安全，只清空应用场景，不影响你的产品和新闻）
    DB.Where("1 = 1").Delete(&model.Application{})

    // 2. 重新导入 10 条完整的真实数据
    apps := []model.Application{
        {Anchor: "EnergyStorage", Order: 1, Top: "Provide professional and stable industrial temperature control products and solutions for various industries", Title: "Energy Storage", Subtitle: "Energy Capture, Storage, Transfer, and Application<br>TONFY Contributes to Every Stage, Making Temperature Controllable and Energy Safer.", Link: "./products.html?i=Energy Storage", Background: "./assets/images/applications/EnergyStorage/hy.webp", Image: "./assets/images/applications/EnergyStorage/cp.webp"},
        {Anchor: "Laser", Order: 2, Top: "Provide professional and stable industrial temperature control products and solutions for various industries", Title: "Laser", Subtitle: "Monitoring Every Temperature Change,<br>Ensuring Precise Control of Energy Release.", Link: "./products.html?i=Laser", Background: "./assets/images/applications/Laser/hy.webp", Image: "./assets/images/applications/Laser/cp.webp"},
        {Anchor: "MachineryManufacturing", Order: 3, Top: "Provide professional and stable industrial temperature control products and solutions for various industries", Title: "Machinery Manufacturing", Subtitle: "Precision Temperature Control<br>20+ Years of Expertise – Showcasing China’s Manufacturing Excellence to the World.", Link: "./products.html?i=Machinery Manufacturing", Background: "./assets/images/applications/MachineryManufacturing/hy.webp", Image: "./assets/images/applications/MachineryManufacturing/cp.webp"},
        {Anchor: "DataCenter", Order: 4, Top: "Provide professional and stable industrial temperature control products and solutions for various industries", Title: "Data Center", Subtitle: "Liquid Cooling Technology & Smart Solutions:<br>Powering the Transformation and Growth of the Data Center Industry.", Link: "./products.html?i=Data Center", Background: "./assets/images/applications/DataCenter/hy.webp", Image: "./assets/images/applications/DataCenter/cp.webp"},
        {Anchor: "Semiconductor", Order: 5, Top: "Provide professional and stable industrial temperature control products and solutions for various industries", Title: "Semiconductor", Subtitle: "Precision & Reliability:<br>Delivering Products for the Efficient, Uninterrupted Operation of Semiconductor Equipment.", Link: "./products.html?i=Semiconductor", Background: "./assets/images/applications/Semiconductor/hy.webp", Image: "./assets/images/applications/Semiconductor/cp.webp"},
        {Anchor: "PowerElectronics", Order: 6, Top: "Provide professional and stable industrial temperature control products and solutions for various industries", Title: "Power Electronics", Subtitle: "Enhancing Daily Life.<br>Providing Secure and Stable Operating Environments for the Industry.", Link: "./products.html?i=Power Electronics", Background: "./assets/images/applications/PowerElectronics/hy.webp", Image: "./assets/images/applications/PowerElectronics/cp.webp"},
        {Anchor: "HydrogenEnergy", Order: 7, Top: "Provide professional and stable industrial temperature control products and solutions for various industries", Title: "Hydrogen Energy", Subtitle: "Full-Process Temperature Management:<br>Supporting the Clean Energy Transition.", Link: "./products.html?i=Hydrogen Energy", Background: "./assets/images/applications/HydrogenEnergy/hy.webp", Image: "./assets/images/applications/HydrogenEnergy/cp.webp"},
        {Anchor: "MedicalDevice", Order: 8, Top: "Provide professional and stable industrial temperature control products and solutions for various industries", Title: "Medical Device", Subtitle: "Full-Process Temperature Management:<br>Advancing the Clean Energy Transition.", Link: "./products.html?i=Medical Device", Background: "./assets/images/applications/MedicalDevice/hy.webp", Image: "./assets/images/applications/MedicalDevice/cp.webp"},
        {Anchor: "OilGas", Order: 9, Top: "Provide professional and stable industrial temperature control products and solutions for various industries", Title: "Oil & Gas Storage and Transportation", Subtitle: "Ensuring Temperature Stability of Every Drop of Energy,<br>Covering the Entire Storage and Transportation Chain.", Link: "#", Background: "./assets/images/applications/OilGasStorage/hy.webp", Image: "./assets/images/applications/OilGasStorage/cp.webp"},
        {Anchor: "NewEnergyVehicles", Order: 10, Top: "Provide professional and stable industrial temperature control products and solutions for various industries", Title: "New Energy Vehicles (Charging & Swapping)", Subtitle: "Safeguarding Energy Futures, Protecting Your Every Charge.", Link: "#", Background: "./assets/images/applications/NewEnergyVehicles/hy.webp", Image: "./assets/images/applications/NewEnergyVehicles/cp.webp"},
    }
    
    for _, app := range apps {
        DB.Create(&app)
    }
    fmt.Println(" 已强制清空旧数据，并成功导入 10 条真实应用场景！")
    // ======================================================
	var aboutCount int64
DB.Model(&model.AboutPage{}).Count(&aboutCount)
if aboutCount == 0 {
    // 准备各个板块的初始 JSON 数据
    introTexts := `["TONFY, a global leader in industrial temperature control...","The company established its first subsidiary, ATF Cooling GmbH...","TONFY boasts an annual production capacity of 300,000 units...","Committed to delivering professional and stable industrial...","Its products and solutions are mainly applied in fields..."]`
    
    esgPillars := `[
        {"header":"Environment","subtitle":"Product, Factory, Supply Chain<br>Green Practices","items":["Eco-friendly product design energy efficiency low-GWP refrigerants","Green Factories low-carbon manufacturing","Lifecycle impact minimization design — operation — disposal"]},
        {"header":"Social","subtitle":"Employee, Customer, Community<br>Responsibility","items":["Safe & healthy workplace (ISO 45001)","Employee supply chain management","Customer support & engagement","ISO 9001 / 14001 / 50001"]},
        {"header":"Governance","subtitle":"International Management System<br>& Risk Control","items":["Compliance with CE, RoHS, REACH","Transparent reporting & ESG disclosure","Risk management across global operations"]}
    ]`
    
    esgBottom := `[
        {"title":"Product Certifications<br>& Standards","desc":"CE, UL, IEC/EN, IEEN, ISO Series"},
        {"title":"Global Service Network","desc":"Cross-timezone, Multi-language Spare Parts Guarantee"},
        {"title":"Supply Chain Transparency<br>& Responsibility","desc":"CE, UL, IEC/EN, IEEN, ISO Series"},
        {"title":"Data & Targets Disclosure","desc":"Carbon Reduction Energy Efficiency improvement Sustainability Report"}
    ]`

    cultureItems := `[
        {"img":"./assets/images/about/sec2-1.webp","title":"Vision","desc":"Leading thermal solutions for future industry."},
        {"img":"./assets/images/about/sec2-2.webp","title":"Mission","desc":"Reliable temperature control for crucial systems."},
        {"img":"./assets/images/about/sec2-3.webp","title":"Commitment","desc":"Sustainable growth and responsive service network."}
    ]`

    footprintMain := `{"img":"./assets/images/about/sec3_1.webp","title":"Sanhe, Hebei","desc":"TONFY Headquarters"}`
    
    footprintSubs := `[
        {"img":"./assets/images/about/sec3-2.webp","year":"2017","title":"Germany","desc":"European Subsidiary Established"},
        {"img":"./assets/images/about/sec3-3.webp","year":"2025","title":"Singapore","desc":"APAC Headquarters"},
        {"img":"./assets/images/about/sec3-4.webp","year":"2025","title":"Thailand","desc":"Regional Manufacturing Hub"}
    ]`

    DB.Create(&model.AboutPage{
        ID: 1,
        BannerImage: "./assets/images/about/about_bg.webp",
        BannerTitle: "GLOBAL LEADER IN INDUSTRIAL<br>TEMPERATURE CONTROL SOLUTIONS",
        IntroVideo: "./assets/videos/english.mp4",
        IntroTextsJSON: introTexts,
        ESGPillarsJSON: esgPillars,
        ESGBottomCardsJSON: esgBottom,
        CultureItemsJSON: cultureItems,
        FootprintMainJSON: footprintMain,
        FootprintSubsJSON: footprintSubs,
    })
    fmt.Println("🚀 已初始化关于我们页面的默认数据！")
}
}
