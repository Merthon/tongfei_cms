package handler

import (
	"net/http"
	"strings"

	"tonfy_CMS/internal/model"
	"tonfy_CMS/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// CheckPermission 权限拦截中间件
func CheckPermission(requiredModule string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userToken, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "未授权的访问"})
			}

			claims, ok := userToken.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "解析Token失败"})
			}

			role, _ := claims["role"].(string)
			modules, _ := claims["modules"].(string)

			// 超级管理员拥有最高权限，直接放行
			if role == "super_admin" {
				return next(c)
			}

			// 如果需要的是超级管理员权限，且当前用户不是，直接拦截
			if requiredModule == "super_admin" {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "仅超级管理员可操作"})
			}

			// 普通管理员，校验其模块权限中是否包含需要的模块
			if strings.Contains(modules, requiredModule) {
				return next(c)
			}

			// 无权限，拦截请求
			return c.JSON(http.StatusForbidden, map[string]string{"error": "您没有此模块的操作权限"})
		}
	}
}
// 定义接收前端创建账号数据的结构体
type CreateEditorRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Modules  string `json:"modules"` 
}

// CreateEditor 超级管理员创建子账号
func CreateEditor(c echo.Context) error {
	req := new(CreateEditorRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "请求数据格式有误"})
	}

	// 1. 基础校验
	if req.Username == "" || req.Password == "" || req.Modules == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "账号、密码和权限模块不能为空"})
	}

	// 2. 检查账号是否已经存在
	var count int64
	repository.DB.Model(&model.AdminUser{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return c.JSON(http.StatusConflict, map[string]string{"error": "该账号名称已存在，请换一个"})
	}

	// 3. 组装新账号（强制将 Role 设为 editor，防止越权创建超管）
	newEditor := model.AdminUser{
		Username: req.Username,
		Password: req.Password, 
		Role:     "editor",
		Modules:  req.Modules,
	}

	// 4. 保存进数据库
	if err := repository.DB.Create(&newEditor).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "数据库保存失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "子账号创建成功！",
		"username": newEditor.Username,
		"modules": newEditor.Modules,
	})
}
// ==========================================
// 获取子账号列表 (查)
// ==========================================
func GetAdminList(c echo.Context) error {
	var users []model.AdminUser
	// 只查普通编辑
	if err := repository.DB.Where("role = ?", "editor").Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "获取账号列表失败"})
	}

	// 安全起见：把返回给前端的密码字段置空
	for i := range users {
		users[i].Password = ""
	}
	return c.JSON(http.StatusOK, users)
}

// ==========================================
// 删除子账号 (删)
// ==========================================
func DeleteEditor(c echo.Context) error {
	id := c.Param("id")
	var user model.AdminUser
	
	if err := repository.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "该账号不存在"})
	}
	
	// 终极安全锁：绝对不允许删除超管
	if user.Role == "super_admin" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "无法删除超级管理员账号"})
	}

	if err := repository.DB.Delete(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "删除失败"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "账号删除成功"})
}

// ==========================================
// 修改子账号权限或密码 (改)
// ==========================================
type UpdateEditorRequest struct {
	Password string `json:"password"` // 如果为空，表示不修改密码
	Modules  string `json:"modules"`
}

func UpdateEditor(c echo.Context) error {
	id := c.Param("id")
	req := new(UpdateEditorRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "请求格式错误"})
	}

	var user model.AdminUser
	if err := repository.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "该账号不存在"})
	}

	if user.Role == "super_admin" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "无法修改超级管理员"})
	}

	if req.Password != "" {
		user.Password = req.Password
	}
	// 更新权限模块
	user.Modules = req.Modules

	if err := repository.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "更新失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "账号更新成功"})
}