package router

import (
	"team_task_hub/backend/internal/handlers"
	"team_task_hub/backend/internal/middleware"
	"team_task_hub/backend/internal/repositories"
	"team_task_hub/backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupOrganizationRoutes 设置组织相关路由
func SetupAIRoutes(router *gin.Engine, db *gorm.DB, authService *services.AuthService) {
	orgRepo := repositories.NewOrganizationRepository(db)
	recommendService := services.NewRecommendationService(orgRepo)
	recommendService.Initialize()
	AIHandler := handlers.NewAIHandler(recommendService)

	AIRoutes := router.Group("/api/ai")
	AIRoutes.Use(middleware.AuthMiddleware(authService))
	{
		AIRoutes.GET("/recommendations/organizations", AIHandler.RecommendOrganizations)
	}
}
