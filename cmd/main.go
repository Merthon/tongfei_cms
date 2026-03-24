package main

import (
	"fmt"
	
	"tonfy_CMS/internal/handler"
	"tonfy_CMS/internal/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	fmt.Println("正在启动同飞 CMS 后端服务 (Echo 驱动)...")
	
	// 1. 初始化数据库
	repository.InitDB()

	// 2. 初始化 Echo 实例
	e := echo.New()
	
	// 3. 挂载官方极其好用的基础中间件
	e.Use(middleware.Logger())  // 自动打印每一次请求的超美观日志
	e.Use(middleware.Recover()) // 防止因为某个接口报错导致整个程序崩溃
	
	// 4. 配置跨域 CORS (官方提供，一行代码搞定)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // 允许所有前端跨域请求
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// 5. 配置 API 路由组
	api := e.Group("/api")
	api.GET("/news", handler.GetNewsList)
	api.GET("/news/:id", handler.GetNewsDetail)

	// 6. 启动服务，监听 8080 端口
	e.Logger.Fatal(e.Start(":8080"))
}