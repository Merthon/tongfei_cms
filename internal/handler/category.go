package handler

import (
	"net/http"
	"tonfy_CMS/internal/model"    
	"tonfy_CMS/internal/repository" 

	"github.com/labstack/echo/v4"
)

// GetCategories 获取所有分类 (前台和后台通用，因为都要按排序展示)
func GetCategories(c echo.Context) error {
	var categories []model.Category
	//取出来的时候，强行按照 sort_order降序排好队！
	if err := repository.DB.Order("sort_order DESC, created_at ASC").Find(&categories).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "获取分类失败"})
	}
	return c.JSON(http.StatusOK, categories)
}

// CreateCategory 创建新行业分类
func CreateCategory(c echo.Context) error {
	var cat model.Category
	if err := c.Bind(&cat); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "参数解析失败"})
	}
	if err := repository.DB.Create(&cat).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "保存失败"})
	}
	return c.JSON(http.StatusOK, cat)
}

// UpdateCategory 更新行业分类 (改名字、改文案、换头图)
func UpdateCategory(c echo.Context) error {
	id := c.Param("id")
	var cat model.Category
	// 先看看数据库里有没有这个行业
	if err := repository.DB.First(&cat, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "分类不存在"})
	}
	// 把前端传来的新数据覆盖上去
	if err := c.Bind(&cat); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "参数解析失败"})
	}
	repository.DB.Save(&cat)
	return c.JSON(http.StatusOK, cat)
}

// DeleteCategory 删除行业分类
func DeleteCategory(c echo.Context) error {
	id := c.Param("id")
	repository.DB.Delete(&model.Category{}, id)
	return c.JSON(http.StatusOK, map[string]string{"message": "删除成功"})
}

// UpdateCategoriesSort 批量更新行业排序
func UpdateCategoriesSort(c echo.Context) error {
	var payload struct {
		IDs []uint `json:"ids"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "参数错误"})
	}

	tx := repository.DB.Begin()
	total := len(payload.IDs)
	for index, id := range payload.IDs {
		sortOrder := total - index
		if err := tx.Model(&model.Category{}).Where("id = ?", id).Update("sort_order", sortOrder).Error; err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "更新失败"})
		}
	}
	tx.Commit()
	return c.JSON(http.StatusOK, map[string]string{"message": "排序已保存"})
}