package handler

import (
	"net/http"
	"tonfy_CMS/internal/model"
	"tonfy_CMS/internal/repository"

	"github.com/labstack/echo/v4"
)

// GetAboutPage 获取页面数据
func GetAboutPage(c echo.Context) error {
	var page model.AboutPage
	if err := repository.DB.First(&page, 1).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "获取数据失败"})
	}
	return c.JSON(http.StatusOK, page)
}

// UpdateAboutPage 更新页面数据
func UpdateAboutPage(c echo.Context) error {
	var req model.AboutPage
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "参数解析失败"})
	}

	var page model.AboutPage
	if err := repository.DB.First(&page, 1).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "找不到记录"})
	}

	// 更新字段
	page.BannerImage = req.BannerImage
	page.BannerTitle = req.BannerTitle
	page.IntroVideo = req.IntroVideo
	page.IntroTextsJSON = req.IntroTextsJSON
	page.ESGPillarsJSON = req.ESGPillarsJSON
	page.ESGBottomCardsJSON = req.ESGBottomCardsJSON
	page.CultureItemsJSON = req.CultureItemsJSON
	page.FootprintMainJSON = req.FootprintMainJSON
	page.FootprintSubsJSON = req.FootprintSubsJSON

	if err := repository.DB.Save(&page).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "更新失败"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "更新成功"})
}