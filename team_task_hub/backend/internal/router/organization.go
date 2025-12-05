// router/organization_router.go
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
func SetupOrganizationRoutes(router *gin.Engine, db *gorm.DB, authService *services.AuthService) {
	// 初始化依赖
	orgRepo := repositories.NewOrganizationRepository(db)
	orgMemberRepo := repositories.NewOrganizationMemberRepository(db)
	orgAppRepo := repositories.NewOrganizationApplicationRepository(db)

	orgService := services.NewOrganizationService(orgRepo, orgMemberRepo, orgAppRepo)
	orgHandler := handlers.NewOrganizationHandler(orgService)

	// 组织管理路由组 - 按权限级别分组
	organizationGroup := router.Group("/api/organization")
	{
		// 需要组织创建者权限的路由
		creatorRoutes := organizationGroup.Group("")
		creatorRoutes.Use(middleware.CreateAuthChain(authService, orgService, "creator")...)
		{

		}

		// 需要组织管理员权限的路由
		adminRoutes := organizationGroup.Group("")
		adminRoutes.Use(middleware.CreateAuthChain(authService, orgService, "admin")...)
		{
			adminRoutes.GET("/application/pending-join", orgHandler.GetOrgPendingJoinApplicationsHandler)
			adminRoutes.DELETE("/remove-member", orgHandler.RemoveOrganizationMemberHandler)
		}

		// 需要组织成员权限的路由
		memberRoutes := organizationGroup.Group("")
		memberRoutes.Use(middleware.CreateAuthChain(authService, orgService, "member")...)
		{

		}

		// 仅需JWT认证的路由（无特定组织权限要求）
		baseAuthRoutes := organizationGroup.Group("")
		baseAuthRoutes.Use(middleware.CreateAuthChain(authService, orgService, "")...)
		{
			baseAuthRoutes.POST("/application/create-organization", orgHandler.SubmitCreateOrganizationApplicationHandler)
			baseAuthRoutes.POST("/applications/join-organization", orgHandler.SubmitJoinApplicationHandler)
			baseAuthRoutes.PATCH("/applications/{id}/process", orgHandler.ProcessApplicationHandler)
			baseAuthRoutes.GET("/my-organizations", orgHandler.GetUserOrganizationOverviewsHandler)
			baseAuthRoutes.GET("/application/pending-create", orgHandler.GetPendingCreateOrgApplicationsHandler)
			baseAuthRoutes.GET("/application/my-applications", orgHandler.GetUserApplicationsHandler)
			baseAuthRoutes.DELETE("/application/delete:{id}", orgHandler.DeleteUserApplicationHandler)
		}
	}
}
