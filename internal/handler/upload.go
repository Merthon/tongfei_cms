package handler


import(
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
)

// UploadImage处理图片上传
func UploadImage(c echo.Context) error{
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "获取上传文件失败"})
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "打开文件失败"})
	}
	defer src.Close()

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	// 拼接保存路径，对应我们之前创建的 uploads/images 目录
	dstPath := filepath.Join("uploads", "images", filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "在服务器创建文件失败"})
	}
	defer dst.Close()

	// 将上传的文件内容拷贝到我们创建的本地文件中
	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "保存文件失败"})
	}

	// 返回该图片可以直接访问的 URL 路径
	imageUrl := fmt.Sprintf("/uploads/images/%s", filename)
	
	return c.JSON(http.StatusOK, map[string]string{
		"message": "上传成功",
		"url":     imageUrl,
	})
}