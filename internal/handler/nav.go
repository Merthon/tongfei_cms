package handler

import (
	"net/http"
	"tonfy_CMS/internal/model"
	"tonfy_CMS/internal/repository"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm" 
)

// GetFullNavData 获取导航栏所有动态数据 (分类+产品+配置)
func GetFullNavData(c echo.Context) error {
	// 1. 获取分类及每个分类下的前5个产品
	var categories []model.Category
	repository.DB.Preload("Products", func(db *gorm.DB) *gorm.DB {
		return db.Order("id desc").Limit(5) // 每个行业只取最新5个
	}).Find(&categories)

	// 2. 获取导航文案配置
	var config model.NavConfig
	repository.DB.First(&config, 1)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"categories": categories,
		"config":     config,
	})
}

func UpdateNavConfig(c echo.Context) error {
	var req model.NavConfig
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "参数解析失败"})
	}

	var config model.NavConfig
	if err := repository.DB.First(&config, 1).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "找不到导航配置记录"})
	}

	// 更新 Service 板块字段
	config.ServiceTitle1 = req.ServiceTitle1
	config.ServiceDesc1 = req.ServiceDesc1
	config.ServiceLink1 = req.ServiceLink1
	config.ServiceTitle2 = req.ServiceTitle2
	config.ServiceDesc2 = req.ServiceDesc2
	config.ServiceLink2 = req.ServiceLink2
	config.ServiceTel = req.ServiceTel
	config.ServiceEmail = req.ServiceEmail

	// 更新 About 板块字段
	config.AboutTitle1 = req.AboutTitle1
	config.AboutLink1 = req.AboutLink1
	config.AboutTitle2 = req.AboutTitle2
	config.AboutLink2 = req.AboutLink2
	config.AboutTitle3 = req.AboutTitle3
	config.AboutLink3 = req.AboutLink3
	config.AboutTitle4 = req.AboutTitle4
	config.AboutLink4 = req.AboutLink4

	// 更新统计数据字段
	config.TickerSymbol = req.TickerSymbol
	config.EmployeeCount = req.EmployeeCount
	config.RdEngineerCount = req.RdEngineerCount

	if err := repository.DB.Save(&config).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "更新导航配置失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "导航配置更新成功"})
}