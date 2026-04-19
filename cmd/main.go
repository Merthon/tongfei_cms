package main

import (
	"fmt"
	"log"

	"tonfy_CMS/internal/handler"
	"tonfy_CMS/internal/repository"

	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("未找到 .env 文件")
	}
	fmt.Println("正在启动同飞后端服务")

	// 初始化数据库
	repository.InitDB()

	// 初始化 Echo 实例
	e := echo.New()

	// 挂载官方极其好用的基础中间件
	e.Use(middleware.Recover())
	// 配置跨域 CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // 允许所有前端跨域请求
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	e.Static("/uploads", "uploads")
	e.Static("/admin", "admin-ui")
	e.Static("/products", "products")

	// ==========================================
	// 配置公开 API 路由组 (无权限拦截)
	// ==========================================
	publicApi := e.Group("/api")
	publicApi.POST("/login", handler.Login)           // 登录接口
	publicApi.GET("/news", handler.GetNewsList)       // 新闻列表接口
	publicApi.GET("/news/:id", handler.GetNewsDetail) // 新闻详情接口
	
	// 产品
	publicApi.GET("/front/products.json", handler.GetProductsJson)
	publicApi.GET("/front/products/:modelName/data.json", handler.GetProductDataJson)
	publicApi.GET("/front/categories", handler.GetCategories)

	// 职位
	publicApi.POST("/apply", handler.SubmitApplication)
	publicApi.GET("/front/jobs", handler.GetFrontJobs)
	
	// 联系我们
	publicApi.POST("/contact", handler.SubmitContact)
	
	// 发送产品邮件
	publicApi.POST("/product/send-manual", handler.SendProductManual)
	
	// 首页banner
	publicApi.GET("/front/banners", handler.GetFrontBanners)

	// ==========================================
	// 受保护的后台 API 路由组 (含权限拦截)
	// ==========================================
	adminApi := e.Group("/api/admin")
	// adminApi 这个组加上 JWT 
	adminApi.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: handler.JwtSecret,
	}))
	adminApi.POST("/account/create_editor", handler.CreateEditor, handler.CheckPermission("super_admin"))
	adminApi.GET("/account/list", handler.GetAdminList, handler.CheckPermission("super_admin"))
    adminApi.DELETE("/account/:id", handler.DeleteEditor, handler.CheckPermission("super_admin"))
    adminApi.PUT("/account/:id", handler.UpdateEditor, handler.CheckPermission("super_admin"))

	// 公共接口 (只要登录了就能用，不加模块拦截器)
	adminApi.POST("/upload", handler.UploadImage)

	// 新闻模块 (需要 news 权限)
	adminApi.POST("/news", handler.CreateNews, handler.CheckPermission("news"))
	adminApi.PUT("/news/:id", handler.UpdateNews, handler.CheckPermission("news"))
	adminApi.DELETE("/news/:id", handler.DeleteNews, handler.CheckPermission("news"))
	
	// 排序接口 (产品权限)
	adminApi.PUT("/products/sort", handler.UpdateProductsSort, handler.CheckPermission("product"))
	
	// 产品管理 (需要 product 权限)
	adminApi.GET("/products", handler.GetAdminProductList, handler.CheckPermission("product"))
	adminApi.GET("/products/:id", handler.GetAdminProductDetail, handler.CheckPermission("product"))
	adminApi.POST("/products", handler.CreateProduct, handler.CheckPermission("product"))
	adminApi.PUT("/products/:id", handler.UpdateProduct, handler.CheckPermission("product"))
	adminApi.DELETE("/products/:id", handler.DeleteProduct, handler.CheckPermission("product"))
	
	// 行业管理 (归属产品模块，需要 product 权限)
	adminApi.GET("/categories", handler.GetCategories, handler.CheckPermission("product"))
	adminApi.POST("/categories", handler.CreateCategory, handler.CheckPermission("product"))
	adminApi.PUT("/categories/:id", handler.UpdateCategory, handler.CheckPermission("product"))
	adminApi.DELETE("/categories/:id", handler.DeleteCategory, handler.CheckPermission("product"))
	adminApi.PUT("/categories/sort", handler.UpdateCategoriesSort, handler.CheckPermission("product"))
	
	// 职位管理 (需要 job 权限)
	adminApi.GET("/jobs", handler.GetAdminJobs, handler.CheckPermission("job"))
	adminApi.POST("/jobs", handler.CreateJob, handler.CheckPermission("job"))
	adminApi.PUT("/jobs/:id", handler.UpdateJob, handler.CheckPermission("job"))
	adminApi.DELETE("/jobs/:id", handler.DeleteJob, handler.CheckPermission("job"))
	adminApi.PUT("/jobs/sort", handler.UpdateJobsSort, handler.CheckPermission("job"))
	adminApi.GET("/applications", handler.GetApplications, handler.CheckPermission("job"))
	adminApi.PUT("/applications/:id/status", handler.UpdateApplicationStatus, handler.CheckPermission("job"))
	adminApi.DELETE("/applications/:id", handler.DeleteApplication, handler.CheckPermission("job"))

	// 首页banner管理 (需要 banner 权限)
	adminApi.GET("/banners", handler.GetAdminBanners, handler.CheckPermission("banner"))
	adminApi.POST("/banners", handler.CreateBanner, handler.CheckPermission("banner"))
	adminApi.PUT("/banners/:id", handler.UpdateBanner, handler.CheckPermission("banner"))
	adminApi.DELETE("/banners/:id", handler.DeleteBanner, handler.CheckPermission("banner"))
	adminApi.PUT("/banners/sort", handler.UpdateBannersSort, handler.CheckPermission("banner"))

	// ==========================================
	// 配置日志轮转规则
	// ==========================================
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

	// 启动服务，监听 8080 端口
	e.Logger.Fatal(e.Start(":8080"))
}