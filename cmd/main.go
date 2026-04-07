package main

import (
	"fmt"
	"log" 
	
	"tonfy_CMS/internal/handler"
	"tonfy_CMS/internal/repository"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo-jwt/v4"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	err := godotenv.Load()
    if err != nil {
        log.Println("未找到 .env 文件")
    }
	fmt.Println("正在启动同飞后端服务")
	
	// 1. 初始化数据库
	repository.InitDB()

	// 2. 初始化 Echo 实例
	e := echo.New()
	
	// 3. 挂载官方极其好用的基础中间件
	e.Use(middleware.Recover()) 
	// 4. 配置跨域 CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // 允许所有前端跨域请求
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	e.Static("/uploads", "uploads")
	e.Static("/admin", "admin-ui")
	e.Static("/products", "products")

	// 配置 API 路由组
	publicApi := e.Group("/api")
	publicApi.POST("/login", handler.Login)           // 登录接口
	publicApi.GET("/news", handler.GetNewsList)       // 新闻列表接口
	publicApi.GET("/news/:id", handler.GetNewsDetail) // 新闻详情接口
	// 产品
	publicApi.GET("/front/products.json", handler.GetProductsJson)
	publicApi.GET("/front/products/:modelName/data.json", handler.GetProductDataJson)
	publicApi.GET("/front/categories", handler.GetCategories)

	//职位
	publicApi.POST("/apply", handler.SubmitApplication)
	publicApi.GET("/front/jobs", handler.GetFrontJobs)
	//联系我们
	publicApi.POST("/contact", handler.SubmitContact)
	//发送产品邮件
	publicApi.POST("/product/send-manual", handler.SendProductManual)
	//首页banner
	publicApi.GET("/front/banners", handler.GetFrontBanners)

	// 受保护的后台
	adminApi := e.Group("/api/admin")
	// adminApi 这个组加上 JWT 
	adminApi.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: handler.JwtSecret,
	}))

	// 把增、删、改放进受保护的组里
	adminApi.POST("/news", handler.CreateNews)
	adminApi.PUT("/news/:id", handler.UpdateNews)
	adminApi.DELETE("/news/:id", handler.DeleteNews)
	adminApi.POST("/upload", handler.UploadImage)
    //排序接口
	adminApi.PUT("/products/sort", handler.UpdateProductsSort)
	// 产品管理
	adminApi.GET("/products", handler.GetAdminProductList)
	adminApi.GET("/products/:id", handler.GetAdminProductDetail)
	adminApi.POST("/products", handler.CreateProduct)
	adminApi.PUT("/products/:id", handler.UpdateProduct)
	adminApi.DELETE("/products/:id", handler.DeleteProduct)
	// 行业管理
	adminApi.GET("/categories", handler.GetCategories)
	adminApi.POST("/categories", handler.CreateCategory)
	adminApi.PUT("/categories/:id", handler.UpdateCategory)
	adminApi.DELETE("/categories/:id", handler.DeleteCategory)
	adminApi.PUT("/categories/sort", handler.UpdateCategoriesSort)
	// 职位管理
	adminApi.GET("/jobs", handler.GetAdminJobs)
    adminApi.POST("/jobs", handler.CreateJob)
    adminApi.PUT("/jobs/:id", handler.UpdateJob)
    adminApi.DELETE("/jobs/:id", handler.DeleteJob)
    adminApi.PUT("/jobs/sort", handler.UpdateJobsSort)
	adminApi.GET("/applications", handler.GetApplications)
    adminApi.PUT("/applications/:id/status", handler.UpdateApplicationStatus)
	adminApi.DELETE("/applications/:id", handler.DeleteApplication)

	// 首页banner管理
	adminApi.GET("/banners", handler.GetAdminBanners)
    adminApi.POST("/banners", handler.CreateBanner)
    adminApi.PUT("/banners/:id", handler.UpdateBanner)
    adminApi.DELETE("/banners/:id", handler.DeleteBanner)
    adminApi.PUT("/banners/sort", handler.UpdateBannersSort)
	
	//运行日志
	// 配置日志轮转规则
    logWriter := &lumberjack.Logger{
        Filename:   "logs/system.log", // 存放在项目根目录的 logs 文件夹下
        MaxSize:    10,                // 每个日志文件最大 10 MB
        MaxBackups: 30,                // 最多保留 30 个旧文件
        MaxAge:     30,                // 最多保留 30 天 (超过天数自动物理删除)
        Compress:   true,              // 自动把旧日志压缩成 .gz 文件，极其省硬盘
    }

    // 将切割引擎接入 Echo 框架
    e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
        Output: logWriter,
        Format: "[${time_rfc3339}] status=${status} latency=${latency_human} ip=${remote_ip} method=${method} uri=${uri} error=${error}\n",
    }))

	// 6. 启动服务，监听 8080 端口
	e.Logger.Fatal(e.Start(":8080"))
}