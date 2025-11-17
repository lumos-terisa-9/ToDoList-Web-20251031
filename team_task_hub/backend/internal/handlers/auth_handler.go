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
