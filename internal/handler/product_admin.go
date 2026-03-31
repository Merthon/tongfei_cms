package handler

import (
	"net/http"

	"tonfy_CMS/internal/model"      
	"tonfy_CMS/internal/repository" 

	"github.com/labstack/echo/v4"
)

// 1. 获取后台产品列表 (按时间倒序)
func GetAdminProductList(c echo.Context) error {
	var products []model.Product
	if err := repository.DB.Order("created_at desc").Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "获取产品列表失败"})
	}
	return c.JSON(http.StatusOK, products)
}

// 2. 获取单个产品详情 (用于回显到编辑弹窗)
func GetAdminProductDetail(c echo.Context) error {
	id := c.Param("id")
	var product model.Product
	if err := repository.DB.First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "产品不存在"})
	}
	return c.JSON(http.StatusOK, product)
}

// 3. 新增产品
func CreateProduct(c echo.Context) error {
	var product model.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "数据格式错误"})
	}

	if err := repository.DB.Create(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "保存产品失败"})
	}
	return c.JSON(http.StatusOK, product)
}

// 4. 更新产品
func UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	var product model.Product

	// 先看看数据库里有没有这个产品
	if err := repository.DB.First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "找不到要更新的产品"})
	}

	// 绑定前端传来的新数据
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "数据格式错误"})
	}

	// 保存更新
	if err := repository.DB.Save(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "更新产品失败"})
	}
	return c.JSON(http.StatusOK, product)
}

// 5. 删除产品
func DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	if err := repository.DB.Delete(&model.Product{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "删除产品失败"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "删除成功"})
}