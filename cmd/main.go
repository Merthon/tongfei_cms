package main

import (
	"fmt"
	
	"tonfy_CMS/internal/handler"
	"tonfy_CMS/internal/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo-jwt/v4"
)

func main() {
	fmt.Println("正在启动同飞后端服务")
	
	// 1. 初始化数据库
	repository.InitDB()

	// 2. 初始化 Echo 实例
	e := echo.New()
	
	// 3. 挂载官方极其好用的基础中间件
	e.Use(middleware.Logger())  
	e.Use(middleware.Recover()) 
	
	// 4. 配置跨域 CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // 允许所有前端跨域请求
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	e.Static("/uploads", "uploads")
	e.Static("/admin", "admin-ui")
	e.Static("/products", "products")

	// 5. 配置 API 路由组
	publicApi := e.Group("/api")
	publicApi.POST("/login", handler.Login)           // 登录接口必须公开
	publicApi.GET("/news", handler.GetNewsList)       // 前台看新闻列表是公开的
	publicApi.GET("/news/:id", handler.GetNewsDetail) // 前台看新闻详情是公开的

	// 产品
	publicApi.GET("/front/products.json", handler.GetProductsJson)
	publicApi.GET("/front/products/:modelName/data.json", handler.GetProductDataJson)
	publicApi.GET("/front/categories", handler.GetCategories)

	// 2. 受保护的后台
	adminApi := e.Group("/api/admin")
	// 给 adminApi 这个组加上 JWT 
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

	// 6. 启动服务，监听 8080 端口
	e.Logger.Fatal(e.Start(":8080"))
}