
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"tonfy_CMS/internal/model" // 【注意：把 tongfei-cms 换成你实际的 go mod 名字！】

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("🚀 老兵的数据搬家战车启动中...")

	// 1. 连接到你的 SQLite 数据库 (确保 cms.db 在根目录)
	db, err := gorm.Open(sqlite.Open("cms.db"), &gorm.Config{})
	if err != nil {
		panic("❌ 连接数据库失败: " + err.Error())
	}

	// 2. 读取总目录 products.json
	productsJSONPath := "./products/products.json"
	jsonFile, err := os.Open(productsJSONPath)
	if err != nil {
		panic("❌ 找不到 ./products/products.json 文件，请检查路径！")
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	
	// 定义一个临时结构体来解析外层的行业和产品数组
	var categories map[string]struct {
		Products []string `json:"产品"`
	}
	if err := json.Unmarshal(byteValue, &categories); err != nil {
		panic("❌ products.json 解析失败: " + err.Error())
	}

	// 3. 开始疯狂遍历和吸入数据
	successCount := 0
	for categoryName, categoryData := range categories {
		for _, productName := range categoryData.Products {
			
			// 去找单品的 data.json
			dataPath := fmt.Sprintf("./products/%s/data.json", productName)
			dataBytes, err := ioutil.ReadFile(dataPath)
			if err != nil {
				fmt.Printf("⚠️  警告: 找不到文件 %s，已跳过\n", dataPath)
				continue
			}

			// 解析单品详情数据
			var detailData map[string]interface{}
			json.Unmarshal(dataBytes, &detailData)

			// 安全地提取我们需要的关键字段 (防空指针)
			name := ""
			if val, ok := detailData["产品-名称"].(string); ok {
				name = val
			}
			mainImage := ""
			if val, ok := detailData["产品-主图"].(string); ok {
				mainImage = val
			}
			fileUrl := ""
			if val, ok := detailData["产品-文件"].(string); ok {
				fileUrl = val
			}

			// 把整个详情对象，重新压缩成 JSON 字符串，准备存入数据库
			detailBytes, _ := json.Marshal(detailData)

			// 检查是否已经迁移过，防止重复运行脚本导致数据重复
			var existing model.Product
			if db.Where("model_name = ?", productName).First(&existing).Error == nil {
				fmt.Printf("⏭️  产品 [%s] 已在库中，跳过\n", productName)
				continue
			}

			// 4. 组装终极产品对象，入库！
			newProduct := model.Product{
				Category:   categoryName,
				Name:       name,
				ModelName:  productName, // 这里存文件夹名，方便以后路由匹配
				MainImage:  mainImage,
				FileUrl:    fileUrl,
				DetailData: string(detailBytes),
			}

			if err := db.Create(&newProduct).Error; err != nil {
				fmt.Printf("❌ 插入 [%s] 失败: %v\n", productName, err)
			} else {
				fmt.Printf("✅ 成功迁移: %s -> %s\n", categoryName, productName)
				successCount++
			}
		}
	}

	fmt.Printf("\n🎉 报告：数据搬家大获全胜！共成功迁移了 %d 个产品到数据库中！\n", successCount)
}