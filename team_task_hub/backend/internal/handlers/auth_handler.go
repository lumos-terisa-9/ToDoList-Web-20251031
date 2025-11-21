package handlers

import (
	"net/http"
	"team_task_hub/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body string true "注册参数" example({"id":1001,"username":"张三","email":"user@example.com","password":"password123","code":"123456"})
// @Success 201 {object} string "注册成功"
// @Failure 400 {object} string "参数错误"
// @Failure 409 {object} string "用户已存在"
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		UserID   uint   `json:"id" binding:"required"`
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		Code     string `json:"code" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数格式错误",
			"error":   err.Error(),
		})
		return
	}

	registerReq := &services.RegisterRequest{
		Userid:   req.UserID,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Code:     req.Code,
	}

	result, err := h.authService.Register(registerReq)
	if err != nil {
		statusCode := http.StatusBadRequest
		errorMessage := err.Error()

		// 冲突错误返回409状态码
		if errorMessage == "邮箱已被注册" {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": errorMessage,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "注册成功",
		"data": gin.H{
			"user_id":  result.UserID,
			"username": result.Username,
			"email":    result.Email,
		},
	})
}

// Login 用户登录
// @Summary 用户登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body string true "登录参数" example({"identifier":"1001","password":"password123"})
// @Success 200 {object} string "登录成功"
// @Failure 400 {object} string "参数错误"
// @Failure 401 {object} string "认证失败"
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Identifier string `json:"identifier" binding:"required"`
		Password   string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数格式错误",
			"error":   err.Error(),
		})
		return
	}

	loginReq := &services.LoginRequest{
		Identifier: req.Identifier,
		Password:   req.Password,
	}

	result, err := h.authService.Login(loginReq)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登录成功",
		"data": gin.H{
			"user": gin.H{
				"user_id":  result.UserID,
				"username": result.Username,
				"email":    result.Email,
			},
			"access_token": result.AccessToken,
			"token_type":   "Bearer",
			"expires_in":   86400,
		},
	})
}

// Logout 用户登出接口
// @Summary 用户登出
// @Description 用户登出，使当前令牌失效。令牌需置于Header中：Authorization: Bearer <token>（bearer和token中间有空格）
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Success 200 {object} services.LogoutResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	tokenString, exists := c.Get("tokenString")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "上下文中未找到认证令牌",
		})
		return
	}

	// 进行类型断言，确保 tokenString 是 string 类型
	tokenStr, ok := tokenString.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "令牌格式错误",
		})
		return
	}

	// 调用服务层使令牌失效（加入黑名单）
	if err := h.authService.Logout(tokenStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "登出失败",
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登出成功",
	})
}

// GetUserProfile 获取当前用户信息
// @Summary 获取当前用户个人信息
// @Description 获取当前已认证登录用户的详细信息，需要有效的JWT令牌。令牌需置于Header中：Authorization: Bearer <token>
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Success 200 {object} map[string]interface{} "成功获取用户信息"
// @Failure 400 {object} string "请求参数错误或用户标识缺失"
// @Failure 401 {object} string "未提供有效的认证令牌"
// @Failure 404 {object} string "用户不存在"
// @Failure 500 {object} string "服务器内部错误或用户标识格式错误"
// @Router /api/auth/me [get]
func (h *AuthHandler) GetUserProfile(c *gin.Context) {
	// 从中间件设置的上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求上下文中未找到用户标识",
		})
		return
	}

	// 类型断言：确保userID是uint类型
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "用户标识格式错误",
		})
		return
	}

	// 调用服务层，通过ID查询用户信息（复用您现有的缓存逻辑）
	user, err := h.authService.FindUserByIDWithCache(userIDUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "用户不存在",
		})
		return
	}

	// 返回用户基本信息
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取用户信息成功",
		"data": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
		},
	})
}
