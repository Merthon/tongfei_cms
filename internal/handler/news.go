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
	
	result := repository.DB.Select("id", "title", "date", "description", "image").Order("date DESC, id DESC").Find(&newsList)
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

// CreateNews 添加新闻
func CreateNews(c echo.Context) error {
	// 初始化一个空的news结构体指针
	news := new(model.News)
	if err := c.Bind(news); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "数据格式不正确"})
	}
	result := repository.DB.Create(news)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "保存到数据库失败"})
	}
	return c.JSON(http.StatusCreated, news)
}
// UpdateNews 修改新闻
func UpdateNews(c echo.Context) error {
	id := c.Param("id")

	// 1. 先去数据库里查一下，看看这条新闻存不存在
	var news model.News
	if err := repository.DB.First(&news, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "找不到要修改的新闻"})
	}

	// 2. 将前端传过来的最新 JSON 数据绑定到这条新闻上
	if err := c.Bind(&news); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "数据格式不正确"})
	}

	// 3. 保存回数据库 (GORM 会自动识别主键 ID 并执行 UPDATE 语句)
	repository.DB.Save(&news)

	return c.JSON(http.StatusOK, news)
}
// DeleteNews 删除新闻
func DeleteNews(c echo.Context) error {
	id := c.Param("id")
	result := repository.DB.Delete(&model.News{}, id)
	
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "删除失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "删除成功"})
}
