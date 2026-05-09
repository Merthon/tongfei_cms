package handler

import (
	"net/http"
	"tonfy_CMS/internal/model"
	"tonfy_CMS/internal/repository"

	"github.com/labstack/echo/v4"
)

// GetHomeConfig 获取主页数据
func GetHomeConfig(c echo.Context) error {
	var config model.HomeConfig
	if err := repository.DB.First(&config, 1).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "获取主页数据失败"})
	}
	return c.JSON(http.StatusOK, config)
}

// UpdateHomeConfig 更新主页数据
func UpdateHomeConfig(c echo.Context) error {
	var req model.HomeConfig
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "参数解析失败"})
	}

	var config model.HomeConfig
	if err := repository.DB.First(&config, 1).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "找不到数据记录"})
	}

	config.EmpCount = req.EmpCount
	config.RdCount = req.RdCount
	config.LabCount = req.LabCount
	config.ValidationHours = req.ValidationHours
	config.CountriesCount = req.CountriesCount
	config.BottomText1 = req.BottomText1
	config.BottomText2 = req.BottomText2
	config.BottomText3 = req.BottomText3

	if err := repository.DB.Save(&config).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "保存失败"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "主页数据更新成功"})
}