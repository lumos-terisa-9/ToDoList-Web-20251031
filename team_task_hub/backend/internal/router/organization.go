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
			creatorRoutes.PATCH("/:id/ownership", orgHandler.TransferOrganizationOwnershipHandler)
		}

		// 需要组织管理员权限的路由
		adminRoutes := organizationGroup.Group("")
		adminRoutes.Use(middleware.CreateAuthChain(authService, orgService, "admin")...)
		{
			adminRoutes.GET("/application/pending-join", orgHandler.GetOrgPendingJoinApplicationsHandler)
			adminRoutes.DELETE("/remove-member", orgHandler.RemoveOrganizationMemberHandler)
			adminRoutes.POST("/change-organization", orgHandler.SubmitChangeOrganizationApplicationHandler)
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
			//创建申请
			baseAuthRoutes.POST("/application/create-organization", orgHandler.SubmitCreateOrganizationApplicationHandler)
			baseAuthRoutes.POST("/application/join-organization", orgHandler.SubmitJoinApplicationHandler)

			//处理申请
			baseAuthRoutes.PATCH("/application/:id/process", orgHandler.ProcessApplicationHandler)

			//个人组织总览
			baseAuthRoutes.GET("/my-organizations", orgHandler.GetUserOrganizationOverviewsHandler)

			//获取待处理的申请
			baseAuthRoutes.GET("/application/pending", orgHandler.GetPendingApplicationsHandler)

			//用户获取所有申请
			baseAuthRoutes.GET("/application/my-applications", orgHandler.GetUserApplicationsHandler)

			//用户删除申请
			baseAuthRoutes.DELETE("/application/:id", orgHandler.DeleteUserApplicationHandler)

			//用户查找组织
			baseAuthRoutes.GET("/search", orgHandler.SearchOrganizationsHandler)
			baseAuthRoutes.GET("/:id", orgHandler.GetOrganizationByIDHandler)

			//更改组织信息
			baseAuthRoutes.PATCH("/:id/name", orgHandler.UpdateOrganizationNameHandler)
			baseAuthRoutes.PATCH("/:id/description", orgHandler.UpdateOrganizationDescriptionHandler)
			baseAuthRoutes.PATCH("/:id/logo", orgHandler.UpdateOrganizationLogoHandler)
		}
	}
}
