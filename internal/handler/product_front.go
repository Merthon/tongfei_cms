package handler

import (
	"encoding/json"
	"net/http"

	"tonfy_CMS/internal/model"      
	"tonfy_CMS/internal/repository" 

	"github.com/labstack/echo/v4"
)

// GetProductsJson 伪装成总目录 products.json
func GetProductsJson(c echo.Context) error {
	var products []model.Product
	
	// 去数据库里查出所有的产品（为了速度，只查所属行业和型号名字）
	if err := repository.DB.Select("category", "model_name").Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "获取产品目录失败"})
	}

	// 拼装成前端需要的那个嵌套 JSON 格式
	// 目标格式: {"行业名": {"产品": ["型号1", "型号2"]}}
	result := make(map[string]map[string][]string)

	for _, p := range products {
		// 如果这个行业还没在 map 里，就初始化它
		if _, exists := result[p.Category]; !exists {
			result[p.Category] = map[string][]string{
				"产品": {},
			}
		}
		// 把产品型号塞进对应的行业数组里
		result[p.Category]["产品"] = append(result[p.Category]["产品"], p.ModelName)
	}

	// 返回给前端！
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
	repository.DB.Select("model_name", "main_image").
		Where("category = ? AND id != ?", product.Category, product.ID).
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