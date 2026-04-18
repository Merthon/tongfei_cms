package handler

import (
	"net/http"
	"time"

	"tonfy_CMS/internal/model"
	"tonfy_CMS/internal/repository"
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

    // 1. 去数据库里找这个账号
    var user model.AdminUser
    // 注意：这里为了兼容你之前的测试，暂时对比明文密码。
    // 在真实企业上线前，密码应该用 bcrypt 加密，这里咱们先跑通逻辑！
    if err := repository.DB.Where("username = ? AND password = ?", req.Username, req.Password).First(&user).Error; err != nil {
        // 如果找不到人，或者密码不对，就会走到这里
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "账号或密码错误"})
    }

    // 2. 登录成功，生成 JWT Token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user":    user.Username,
        "role":    user.Role,    
        "modules": user.Modules, 
        "exp":     time.Now().Add(time.Hour * 72).Unix(), 
    })

    t, err := token.SignedString(JwtSecret)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "生成Token失败"})
    }

    // 3. 返回给前端（顺便把角色和权限明文返回，方便前端隐藏菜单）
    return c.JSON(http.StatusOK, map[string]string{
        "message": "登录成功",
        "token":   t,
        "role":    user.Role,
        "modules": user.Modules,
    })
}