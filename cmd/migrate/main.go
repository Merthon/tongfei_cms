package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"tonfy_CMS/internal/model"      // 替换成你的 go mod 名字
	"tonfy_CMS/internal/repository" // 替换成你的 go mod 名字
)

// OldNews 定义一个临时结构体，用来解析你那个 news.json
type OldNews struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Image       string `json:"image"`
	// Link 字段我们不需要存进数据库，所以这里可以忽略
}

func main() {
	fmt.Println("🚀 开始执行 34 条老新闻的自动化搬家脚本...")

	// 1. 初始化数据库连接
	repository.InitDB()

	// 2. 读取 news.json 文件
	jsonFile, err := os.Open("assets/data/news.json")
	if err != nil {
		log.Fatalf("❌ 找不到 news.json 文件，请检查路径: %v", err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var oldNewsList []OldNews
	// 将 JSON 数组解析到我们的临时结构体切片中
	if err := json.Unmarshal(byteValue, &oldNewsList); err != nil {
		log.Fatalf("❌ 解析 JSON 失败: %v", err)
	}

	fmt.Printf("✅ 成功读取到 %d 条老新闻元数据，开始逐条缝合正文...\n", len(oldNewsList))

	// 3. 循环遍历这 34 条数据
	successCount := 0
	for _, old := range oldNewsList {
		// 拼装这篇新闻对应的 HTML 正文文件路径，例如: assets/data/articles/1.html
		htmlPath := fmt.Sprintf("assets/data/articles/%d.html", old.ID)
		
		var contentHtml string
		// 尝试读取对应的 HTML 文件
		htmlContent, err := os.ReadFile(htmlPath)
		if err != nil {
			fmt.Printf("⚠️  警告: 找不到 ID 为 %d 的正文文件 (%s)，已将其正文设为空\n", old.ID, htmlPath)
			contentHtml = "<p>暂无正文内容</p>" // 如果真找不到，给个默认兜底
		} else {
			contentHtml = string(htmlContent) // 成功读取到了正文！
		}

		// 4. 将老数据完美映射到我们全新的 GORM 数据模型上
		newRecord := model.News{
			Title:       old.Title,
			Date:        old.Date,
			Description: old.Description,
			Image:       old.Image, // 老图片的路径原封不动存进去 (比如 ./assets/...)
			Content:     contentHtml, // 缝合进来的 HTML 正文
		}

		// 5. 插入 SQLite 数据库！
		if err := repository.DB.Create(&newRecord).Error; err != nil {
			fmt.Printf("❌ 插入 ID %d 失败: %v\n", old.ID, err)
		} else {
			successCount++
		}
	}

	fmt.Printf("🎉 搬家大功告成！成功将 %d 条新闻存入数据库！\n", successCount)
}