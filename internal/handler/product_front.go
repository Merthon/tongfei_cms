package handler

import (
	"encoding/json"
	"net/http"

	"tonfy_CMS/internal/model"      
	"tonfy_CMS/internal/repository" 

	"github.com/labstack/echo/v4"
)

// GetProductsJson 伪装成总目录 products.json
// GetProductsJson 伪装成总目录 products.json
func GetProductsJson(c echo.Context) error {
    var products []model.Product
    
    // 🚨 【核心修复】：加上 Order 指令！让权重高的排前面，权重一样的按创建时间排！
    // 注意：既然要用 created_at 排序，最好把它也 select 出来防错。
    err := repository.DB.Select("category", "model_name", "sort_order", "created_at").
        Order("sort_order DESC, created_at DESC").
        Find(&products).Error
        
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "获取产品目录失败"})
    }

    result := make(map[string]map[string][]string)

    for _, p := range products {
        if _, exists := result[p.Category]; !exists {
            result[p.Category] = map[string][]string{
                "产品": {},
            }
        }
        result[p.Category]["产品"] = append(result[p.Category]["产品"], p.ModelName)
    }

    return c.JSON(http.StatusOK, result)
}

// GetProductDataJson 伪装成单品文件夹里的 data.json
func GetProductDataJson(c echo.Context) error {
	modelName := c.Param("modelName")

	var product model.Product
	// 1. 找当前产品
	if err := repository.DB.Where("model_name = ?", modelName).First(&product).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "找不到该产品的数据"})
	}

	var detailData map[string]interface{}
	if product.DetailData == "" {
		detailData = make(map[string]interface{})
	} else if err := json.Unmarshal([]byte(product.DetailData), &detailData); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "产品数据解析失败"})
	}

	var relatedProducts []model.Product
	// 去数据库查：行业相同，但排除掉当前自己这个产品，最多取 8 个
	repository.DB.Select("model_name", "main_image", "sort_order", "created_at").
        Where("category = ? AND id != ?", product.Category, product.ID).
        Order("sort_order DESC, created_at DESC"). // 排队魔法！
        Limit(8).Find(&relatedProducts)
	// 准备两个空数组
	var recommendLinks []string
	var recommendImages []string

	// 遍历查出来的兄弟产品，把它们的名字和主图塞进数组
	for _, p := range relatedProducts {
		recommendLinks = append(recommendLinks, p.ModelName)
		// 注意这里：如果是新上传的图片(带/uploads)，就直接用；如果是老图片，就拼上相对路径
		img := p.MainImage
		if len(img) > 0 && img[0] != '/' && img[:4] != "http" {
			img = "../" + p.ModelName + "/" + img
		}
		recommendImages = append(recommendImages, img)
	}

	detailData["推荐-链接"] = recommendLinks
	detailData["推荐-图片"] = recommendImages
	// =========================================================

	return c.JSON(http.StatusOK, detailData)
}