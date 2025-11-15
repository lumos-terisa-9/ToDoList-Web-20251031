// internal/router/email.go
package router

import (
	"team_task_hub/backend/internal/handlers"
	"team_task_hub/backend/internal/repositories"
	"team_task_hub/backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupEmailRoutes 设置邮箱相关路由
func SetupEmailRoutes(r *gin.RouterGroup, db *gorm.DB, emailConfig *services.EmailConfig) {
	// 初始化依赖
	codeRepo := repositories.NewVerificationCodeRepository(db)
	emailService := services.NewEmailService(emailConfig, codeRepo)
	emailHandler := handlers.NewEmailHandler(emailService)

	// 邮箱验证路由组
	email := r.Group("/email")
	{
		email.POST("/send-verification", emailHandler.SendVerificationCode) // 发送验证码
		email.POST("/verify", emailHandler.VerifyCode)                      // 验证验证码
	}

	// 管理员路由组（需要权限验证）
	admin := r.Group("/admin")
	admin.GET("/email/test-connection", emailHandler.TestSMTPConnection) // 测试SMTP连接
	admin.POST("/email/cleanup", emailHandler.CleanupExpiredCodes)       // 清理过期验证码
	admin.GET("/email/stats", emailHandler.GetEmailStats)                // 获取统计信
}
