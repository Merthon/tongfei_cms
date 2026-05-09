package handler

import (
	"net/http"
	"tonfy_CMS/internal/model"
	"tonfy_CMS/internal/repository"

	"github.com/labstack/echo/v4"
)

// GetServicePage 获取页面数据 (前后台通用)
func GetServicePage(c echo.Context) error {
	var page model.ServicePage
	// 永远只查 ID 为 1 的那条记录
	if err := repository.DB.First(&page, 1).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "获取数据失败"})
	}
	return c.JSON(http.StatusOK, page)
}

// UpdateServicePage 更新页面数据 (仅后台用)
func UpdateServicePage(c echo.Context) error {
	var req model.ServicePage
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "参数解析失败"})
	}

	var page model.ServicePage
	if err := repository.DB.First(&page, 1).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "找不到配置记录"})
	}

	// 更新字段
	page.Stat1Num = req.Stat1Num
	page.Stat1Label = req.Stat1Label
	page.Stat2Num = req.Stat2Num
	page.Stat2Label = req.Stat2Label
	page.TeamText1 = req.TeamText1
	page.TeamText2 = req.TeamText2
	
	// 如果前端传了新图片就更新
	if req.TeamImage != "" {
		page.TeamImage = req.TeamImage
	}
	
	// 更新卡片 JSON
	if req.CardsJSON != "" {
		page.CardsJSON = req.CardsJSON
	}

	// 保存到数据库
	if err := repository.DB.Save(&page).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "更新失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "更新成功"})
}