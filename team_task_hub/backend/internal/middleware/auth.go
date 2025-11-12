package middleware

import "github.com/gin-gonic/gin"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 简单的认证中间件，后续完善
		c.Next()
	}
}
