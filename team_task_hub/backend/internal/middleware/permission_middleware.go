// middleware/permission_middleware.go
package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"team_task_hub/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// 权限验证中间件

// RequireOrganizationCreator 要求用户是组织创建者
func RequireOrganizationCreator(orgService *services.OrganizationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户ID
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "用户未认证",
			})
			c.Abort()
			return
		}

		// 从URL参数获取组织ID
		orgID, err := getOrganizationIDFromParam(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// 验证用户是否是组织创建者
		isCreator, err := orgService.IsOrganizationCreator(orgID, userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "权限验证失败: " + err.Error(),
			})
			c.Abort()
			return
		}

		if !isCreator {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "权限不足：需要组织创建者权限",
			})
			c.Abort()
			return
		}

		// 权限验证通过，继续处理
		c.Next()
	}
}

// RequireOrganizationAdmin 要求用户是组织管理员或创建者
func RequireOrganizationAdmin(orgService *services.OrganizationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户ID
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "用户未认证",
			})
			c.Abort()
			return
		}

		// 从URL参数获取组织ID
		orgID, err := getOrganizationIDFromParam(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// 验证用户是否是组织管理员或创建者
		isAdmin, err := orgService.IsOrganizationAdmin(orgID, userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "权限验证失败: " + err.Error(),
			})
			c.Abort()
			return
		}

		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "权限不足：需要组织管理员权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireOrganizationMember 要求用户是组织成员（包括创建者、管理员、普通成员）
func RequireOrganizationMember(orgService *services.OrganizationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "用户未认证",
			})
			c.Abort()
			return
		}

		// 从URL参数获取组织ID
		orgID, err := getOrganizationIDFromParam(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// 验证用户是否是组织成员
		isMember, err := orgService.IsOrganizationMember(orgID, userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "权限验证失败: " + err.Error(),
			})
			c.Abort()
			return
		}

		if !isMember {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "权限不足：需要组织成员权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// getOrganizationIDFromParam 从URL参数中获取组织ID
func getOrganizationIDFromParam(c *gin.Context) (uint, error) {
	orgIDStr := c.Param("orgID")
	if orgIDStr == "" {
		return 0, fmt.Errorf("组织ID参数缺失")
	}

	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("无效的组织ID格式")
	}

	return uint(orgID), nil
}

// 链式组合中间件

// WithAuth 组合JWT认证和权限验证的高阶函数
func WithAuth(authService *services.AuthService, permissionMiddleware gin.HandlerFunc) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		AuthMiddleware(authService), // 先执行JWT认证
		permissionMiddleware,        // 再执行权限验证
	}
}

// CreateAuthChain 创建认证链的函数
func CreateAuthChain(authService *services.AuthService, orgService *services.OrganizationService, permissionLevel string) []gin.HandlerFunc {
	switch permissionLevel {
	case "creator":
		return WithAuth(authService, RequireOrganizationCreator(orgService))
	case "admin":
		return WithAuth(authService, RequireOrganizationAdmin(orgService))
	case "member":
		return WithAuth(authService, RequireOrganizationMember(orgService))
	default:
		return []gin.HandlerFunc{AuthMiddleware(authService)}
	}
}
