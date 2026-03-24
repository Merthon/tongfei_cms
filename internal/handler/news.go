package handler

import (
	"net/http"

	"tonfy_CMS/internal/model"
	"tonfy_CMS/internal/repository"

	"github.com/labstack/echo/v4"
)

// GetNewsList 获取新闻列表
func GetNewsList(c echo.Context) error {
	var newsList []model.News
	
	result := repository.DB.Select("id", "title", "date", "description", "image").Find(&newsList)
	if result.Error != nil {
		// Echo 直接 return c.JSON，非常简洁
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "获取新闻列表失败"})
	}

	return c.JSON(http.StatusOK, newsList)
}

// GetNewsDetail 获取单条新闻详情
func GetNewsDetail(c echo.Context) error {
	id := c.Param("id")

	var news model.News
	result := repository.DB.First(&news, id)
	
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "新闻不存在"})
	}

	return c.JSON(http.StatusOK, news)
}