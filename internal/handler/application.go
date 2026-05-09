package handler

import (
	"net/http"
	"tonfy_CMS/internal/model"
	"tonfy_CMS/internal/repository"

	"github.com/labstack/echo/v4"
)

// GetAppScenarios 获取所有应用场景列表 (前后台通用)
func GetAppScenarios(c echo.Context) error {
	var apps []model.Application
	// 严格按照 order 字段升序排列
	if err := repository.DB.Order("`order` asc").Find(&apps).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "获取应用场景失败"})
	}
	return c.JSON(http.StatusOK, apps)
}

// UpdateAppScenario 更新单个应用场景 (仅后台用)
func UpdateAppScenario(c echo.Context) error {
	id := c.Param("id")
	var req model.Application
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "参数解析失败"})
	}

	var app model.Application
	if err := repository.DB.First(&app, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "找不到该应用场景记录"})
	}

	// 更新可变字段
	app.Top = req.Top
	app.Title = req.Title
	app.Subtitle = req.Subtitle
	app.Link = req.Link
	
	if req.Background != "" {
		app.Background = req.Background
	}
	if req.Image != "" {
		app.Image = req.Image
	}

	if err := repository.DB.Save(&app).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "保存失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "更新成功"})
}