package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 直接从上下文中获取角色信息
		roleValue, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "权限信息缺失",
			})
			c.Abort()
			return
		}

		role, ok := roleValue.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "角色信息格式错误",
			})
			c.Abort()
			return
		}

		// 判断：角色是否为管理员
		if role != 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "权限不足：需要管理员权限",
			})
			c.Abort()
			return
		}

		// 是管理员，放行
		c.Next()
	}
}
