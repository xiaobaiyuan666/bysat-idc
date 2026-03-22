package auth

import (
	"net/http"

	appErrors "idc-finance/internal/platform/errors"
	"idc-finance/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

const DemoToken = "phase1-admin-token"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token       string   `json:"token"`
	DisplayName string   `json:"displayName"`
	Roles       []string `json:"roles"`
}

func RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/login", login)
}

func login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "登录参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	if req.Username != "admin" || req.Password != "Admin123!" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":      "AUTH_FAILED",
			"message":   "账号或密码错误",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(LoginResponse{
		Token:       DemoToken,
		DisplayName: "系统管理员",
		Roles:       []string{"super-admin"},
	}, middleware.GetRequestID(c)))
}
