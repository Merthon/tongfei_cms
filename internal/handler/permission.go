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
			// 1. 获取 JWT 中间件解析后的 Token 数据
			// 注意："user" 是 Echo JWT 中间件默认的 Context Key
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

			// 2. 超级管理员拥有最高权限，直接放行
			if role == "super_admin" {
				return next(c)
			}

			// 3. 如果需要的是超级管理员权限，且当前用户不是，直接拦截
			if requiredModule == "super_admin" {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "仅超级管理员可操作"})
			}

			// 4. 普通管理员，校验其模块权限中是否包含需要的模块
			if strings.Contains(modules, requiredModule) {
				return next(c)
			}

			// 5. 无权限，拦截请求
			return c.JSON(http.StatusForbidden, map[string]string{"error": "您没有此模块的操作权限"})
		}
	}
}
// 定义接收前端创建账号数据的结构体
type CreateEditorRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Modules  string `json:"modules"` // 前端传过来的权限字符串，比如 "news,product"
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