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
	if err := repository.DB.Order("sort_order DESC, created_at DESC").Find(&products).Error; err != nil {
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

// 3.31 批量更新产品排序
func UpdateProductsSort(c echo.Context) error {
	// 切片接数据
	var payload struct {
		IDs []uint `json:"ids"`
	}

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "参数格式错误"})
	}

	// 开启事务：保证这一组产品的权重更新要么全成功，要么全回滚
	tx := repository.DB.Begin()

	// 遍历这个切片，索引越小的，赋予的权重越高。
	total := len(payload.IDs)
	for index, id := range payload.IDs {
		sortOrder := total - index

		// 更新数据库
		if err := tx.Model(&model.Product{}).Where("id = ?", id).Update("sort_order", sortOrder).Error; err != nil {
			tx.Rollback()
			return c.JSON(500, map[string]string{"error": "数据库保存失败"})
		}
	}

	tx.Commit()
	return c.JSON(200, map[string]string{"message": "排序保存成功"})
}
