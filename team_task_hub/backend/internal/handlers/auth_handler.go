package handlers

import (
	"net/http"
	"strings"
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
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 从Authorization头获取令牌
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "未提供认证令牌",
		})
		return
	}

	// 解析Bearer Token的格式
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "令牌格式错误，应为：Bearer <token>",
		})
		return
	}

	tokenString := parts[1] // 提取出真正的Token字符串

	// 调用服务层使令牌失效
	if err := h.authService.Logout(tokenString); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "登出失败",
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, services.LogoutResponse{
		Message: "登出成功",
		Success: true,
	})
}
