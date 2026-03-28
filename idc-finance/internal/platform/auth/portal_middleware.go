package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const PortalDemoToken = "phase1-portal-token"

func PortalGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "缺少门户登录凭证",
			})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		if token != PortalDemoToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "门户登录凭证无效",
			})
			return
		}

		customerID, customerName := currentPortalIdentity()
		c.Set("customerID", customerID)
		c.Set("customerName", customerName)
		c.Next()
	}
}
