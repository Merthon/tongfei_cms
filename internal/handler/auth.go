package handler

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// 定义用于接收前端账号密码的结构体
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// JWT 签名密钥
var JwtSecret = []byte("tongfei-cms-super-secret-key")

// Login 处理后台登录
func Login(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "请求格式错误"})
	}

	// 测试方案：固定账号密码
	if req.Username == "admin" && req.Password == "123456" {
		// 登录成功，生成 JWT Token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user": "admin",
			"exp":  time.Now().Add(time.Hour * 72).Unix(), // Token 72小时过期
		})

		// 用密钥签名
		t, err := token.SignedString(JwtSecret)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "生成Token失败"})
		}

		// 返回给前端
		return c.JSON(http.StatusOK, map[string]string{
			"message": "登录成功",
			"token":   t,
		})
	}

	return c.JSON(http.StatusUnauthorized, map[string]string{"error": "账号或密码错误"})
}