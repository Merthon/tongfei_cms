package main

import (
	"fmt"
	"tonfy_CMS/internal/repository"
)

func main(){
	fmt.Println("正在启动同飞CMS后端服务")

	// 初始化数据库
	repository.InitDB()

	fmt.Println("系统初始化成功")
}