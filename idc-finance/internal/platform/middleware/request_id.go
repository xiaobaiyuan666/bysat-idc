package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

const RequestIDKey = "requestID"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		buffer := make([]byte, 8)
		_, _ = rand.Read(buffer)
		requestID := hex.EncodeToString(buffer)
		c.Set(RequestIDKey, requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	value, exists := c.Get(RequestIDKey)
	if !exists {
		return ""
	}
	requestID, _ := value.(string)
	return requestID
}
