package router

import (
	"team_task_hub/backend/internal/handlers"
	"team_task_hub/backend/internal/middleware"
	"team_task_hub/backend/internal/repositories"
	"team_task_hub/backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupAuthRoutes 设置认证路由
func SetupAuthRoutes(r *gin.RouterGroup, db *gorm.DB, emailService *services.EmailService, jwtSecret string) {
	// 初始化依赖
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, emailService, jwtSecret)
	authHandler := handlers.NewAuthHandler(authService)

	// 认证路由组
	auth := r.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	protected := r.Group("/auth")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		protected.POST("/logout", authHandler.Logout)
		protected.GET("/me", authHandler.GetUserProfile)
		protected.PUT("/change_userName", authHandler.UpdateUserName)
		protected.PUT("/change_avatar", authHandler.UpdateAvatar)
		protected.PUT("/change_email", authHandler.UpdateEmail)
		protected.PUT("/change_password", authHandler.UpdatePassword)
	}
}
