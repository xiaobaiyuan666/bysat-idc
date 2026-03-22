package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AdminGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "缺少登录凭证",
			})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		if token != DemoToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "登录凭证无效",
			})
			return
		}

		c.Set("adminID", int64(1))
		c.Set("adminRoles", []string{"super-admin"})
		c.Set("adminName", "系统管理员")
		c.Next()
	}
}
