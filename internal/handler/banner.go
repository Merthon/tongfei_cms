package handler

import (
	"net/http"
	"os"         
	"strings"    

	"tonfy_CMS/internal/model"
	"tonfy_CMS/internal/repository"

	"github.com/labstack/echo/v4"
)

// =================  前台接口 (公开) =================

// GetFrontBanners 获取前台展示的 Banner 
func GetFrontBanners(c echo.Context) error {
	var banners []model.Banner
	if err := repository.DB.Where("is_active = ?", true).Order("sort_order DESC, created_at ASC").Find(&banners).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "获取数据失败"})
	}
	return c.JSON(http.StatusOK, banners)
}

// ================= 后台接口 (受保护) =================

// GetAdminBanners 获取后台管理列表
func GetAdminBanners(c echo.Context) error {
	var banners []model.Banner
	if err := repository.DB.Order("sort_order DESC, created_at ASC").Find(&banners).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "获取数据失败"})
	}
	return c.JSON(http.StatusOK, banners)
}

func CreateBanner(c echo.Context) error {
	var banner model.Banner
	if err := c.Bind(&banner); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "参数解析失败"})
	}
	if err := repository.DB.Create(&banner).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "保存失败"})
	}
	return c.JSON(http.StatusOK, banner)
}

func UpdateBanner(c echo.Context) error {
	id := c.Param("id")
	var banner model.Banner
	if err := repository.DB.First(&banner, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "数据不存在"})
	}
	if err := c.Bind(&banner); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "参数解析失败"})
	}
	repository.DB.Save(&banner)
	return c.JSON(http.StatusOK, banner)
}

func DeleteBanner(c echo.Context) error {
	id := c.Param("id")

	// 1. 先把这条记录查出来，获取它的媒体文件路径
	var banner model.Banner
	if err := repository.DB.First(&banner, id).Error; err == nil {
		// 2. 物理删除文件。
		// 数据库里存的 MediaUrl 是带斜杠的绝对路由（例如: "/uploads/123.mp4"）
		// 我们必须用 strings.TrimPrefix 把它头上的 "/" 砍掉，变成相对路径 "uploads/123.mp4" 才能删
		if banner.MediaUrl != "" {
			filePath := strings.TrimPrefix(banner.MediaUrl, "/")
			os.Remove(filePath) // 尝试删除文件，即使没找到也不会让程序崩溃
		}
	}

	// 3. 彻底删除数据库记录
	repository.DB.Delete(&model.Banner{}, id)
	return c.JSON(http.StatusOK, map[string]string{"message": "素材与物理文件已彻底删除"})
}

func UpdateBannersSort(c echo.Context) error {
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
		if err := tx.Model(&model.Banner{}).Where("id = ?", id).Update("sort_order", sortOrder).Error; err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "更新排序失败"})
		}
	}
	tx.Commit()
	return c.JSON(http.StatusOK, map[string]string{"message": "排序已保存"})
}