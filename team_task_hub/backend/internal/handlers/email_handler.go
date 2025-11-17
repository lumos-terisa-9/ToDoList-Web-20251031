package handlers

import (
	"net/http"
	"team_task_hub/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// EmailHandler 邮箱验证码相关的HTTP处理器
type EmailHandler struct {
	emailService *services.EmailService
}

// NewEmailHandler 创建邮箱处理器
func NewEmailHandler(emailService *services.EmailService) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
	}
}

// SendVerificationCodeRequest 发送验证码请求结构体
type SendVerificationCodeRequest struct {
	Email    string `json:"email" binding:"required,email"` // 邮箱地址
	Business string `json:"business" binding:"required"`    // 业务类型: register, reset_password, change_email, join_org
}

// VerifyCodeRequest 验证验证码请求结构体
type VerifyCodeRequest struct {
	Email    string `json:"email" binding:"required,email"` // 邮箱地址
	Code     string `json:"code" binding:"required,len=6"`  // 6位验证码
	Business string `json:"business" binding:"required"`    // 业务类型
}

// SendVerificationCode 发送邮箱验证码
// @Summary 发送邮箱验证码
// @Description 向指定邮箱发送6位数字验证码，支持多种业务场景
// @Tags 邮箱验证
// @Accept json
// @Produce json
// @Param request body SendVerificationCodeRequest true "发送验证码请求参数"
// @Success 200 {object} string "验证码发送成功"
// @Failure 400 {object} string "请求参数错误"
// @Failure 429 {object} string "请求过于频繁"
// @Failure 500 {object} string "服务器内部错误"
// @Router /api/email/send-verification [post]
func (h *EmailHandler) SendVerificationCode(c *gin.Context) {
	var req SendVerificationCodeRequest

	// 参数绑定和基础验证
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数格式错误",
			"error":   err.Error(),
		})
		return
	}

	// 2. 验证业务类型
	validBusinessTypes := map[string]bool{
		"register":       true,
		"reset_password": true,
		"change_email":   true,
		"join_org":       true,
	}
	if !validBusinessTypes[req.Business] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "不支持的业务类型，请使用: register, reset_password, change_email, join_org",
		})
		return
	}

	// 3. 调用EmailService发送验证码
	if err := h.emailService.SendVerificationCode(req.Email, req.Business); err != nil {
		errorMessage := err.Error()
		statusCode := http.StatusBadRequest

		// 频率限制错误返回429状态码
		if errorMessage == "请求过于频繁，请60秒后再试" {
			statusCode = http.StatusTooManyRequests
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": errorMessage,
		})
		return
	}

	// 4. 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "验证码已发送到您的邮箱，5分钟内有效",
		"data": gin.H{
			"email":    req.Email,
			"business": req.Business,
		},
	})
}

// VerifyCode 验证邮箱验证码
// @Summary 验证邮箱验证码
// @Description 验证邮箱和验证码是否匹配且有效，支持多种业务场景
// @Tags 邮箱验证
// @Accept json
// @Produce json
// @Param request body VerifyCodeRequest true "验证验证码请求参数"
// @Success 200 {object} string "验证码验证成功"
// @Failure 400 {object} string "验证码无效或已过期"
// @Failure 500 {object} string "服务器内部错误"
// @Router /api/email/verify [post]
func (h *EmailHandler) VerifyCode(c *gin.Context) {
	var req VerifyCodeRequest

	// 参数绑定和基础验证
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数格式错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证业务类型
	validBusinessTypes := map[string]bool{
		"register":       true,
		"reset_password": true,
		"change_email":   true,
		"join_org":       true,
	}
	if !validBusinessTypes[req.Business] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "不支持的业务类型",
		})
		return
	}

	// 调用EmailService验证验证码
	isValid, verificationCode, err := h.emailService.VerifyCode(req.Email, req.Code, req.Business)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "验证码验证失败",
			"error":   err.Error(),
		})
		return
	}

	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "验证码无效或已过期",
		})
		return
	}

	// 返回验证成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "验证码验证成功",
		"data": gin.H{
			"email":       req.Email,
			"business":    req.Business,
			"valid_until": verificationCode.ExpiresAt,
		},
	})
}

// TestSMTPConnection 测试SMTP连接（管理员功能）
// @Summary 测试SMTP连接
// @Description 测试邮箱SMTP服务器连接状态（仅用于调试）
// @Tags 系统管理
// @Produce json
// @Success 200 {object} string "SMTP连接正常"
// @Failure 500 {object} string "SMTP连接失败"
// @Router /api/admin/email/test-connection [get]
func (h *EmailHandler) TestSMTPConnection(c *gin.Context) {
	if err := h.emailService.TestConnection(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "SMTP连接测试失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "SMTP连接正常",
	})
}

// CleanupExpiredCodes 清理过期验证码（管理员功能）
// @Summary 清理过期验证码
// @Description 清理数据库中已过期的验证码记录（定时任务调用）
// @Tags 系统管理
// @Produce json
// @Success 200 {object} string "清理完成"
// @Failure 500 {object} string "清理失败"
// @Router /api/admin/email/cleanup [post]
func (h *EmailHandler) CleanupExpiredCodes(c *gin.Context) {
	if err := h.emailService.CleanupExpiredCodes(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "清理过期验证码失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "过期验证码清理完成",
	})
}

// GetEmailStats 获取邮箱服务统计信息（管理员功能）
// @Summary 获取邮箱服务统计
// @Description 获取邮箱验证码发送统计信息
// @Tags 系统管理
// @Produce json
// @Success 200 {object} string "统计信息"
// @Router /api/admin/email/stats [get]
func (h *EmailHandler) GetEmailStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"today_sent":      0,
			"success_rate":    "100%",
			"active_sessions": 0,
		},
	})
}
