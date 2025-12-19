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
	codeRepo := repositories.NewVerificationCodeRepository(db)
	activityRepo := repositories.NewActivityRepository(db)
	participationRepo := repositories.NewActivityParticipationRepository(db)

	orgService := services.NewOrganizationService(orgRepo, orgMemberRepo, orgAppRepo, codeRepo)
	activityService := services.NewActivityService(activityRepo, participationRepo)
	orgHandler := handlers.NewOrganizationHandler(orgService, activityService)

	//需要超级管理员权限的路由
	superAdminRoutes := router.Group("/admin/api/organization")
	superAdminRoutes.Use(middleware.AuthMiddleware(authService), middleware.RequireAdmin())
	{
		//获取待处理的申请
		superAdminRoutes.GET("/application/pending", orgHandler.GetPendingApplicationsHandler)

		//更改组织信息
		superAdminRoutes.PATCH("/:orgID/name", orgHandler.UpdateOrganizationNameHandler)
		superAdminRoutes.PATCH("/:orgID/description", orgHandler.UpdateOrganizationDescriptionHandler)
		superAdminRoutes.PATCH("/:orgID/logo", orgHandler.UpdateOrganizationLogoHandler)
	}

	// 组织管理路由组 - 按权限级别分组
	organizationGroup := router.Group("/api/organization")
	{
		// 需要组织创建者权限的路由
		creatorRoutes := organizationGroup.Group("")
		creatorRoutes.Use(middleware.CreateAuthChain(authService, orgService, "creator")...)
		{
			//提拔降级管理员
			creatorRoutes.PATCH("/:orgID/admin/:userID", orgHandler.PromoteToAdminHandler)
			creatorRoutes.DELETE("/:orgID/admin/:userID", orgHandler.CancelAdminHandler)
			//转移组织所有权
			creatorRoutes.PATCH("/ownership/:orgID", orgHandler.TransferOrganizationOwnershipHandler)
		}

		// 需要组织管理员权限的路由
		adminRoutes := organizationGroup.Group("")
		adminRoutes.Use(middleware.CreateAuthChain(authService, orgService, "admin")...)
		{
			adminRoutes.GET("/application/pending-join", orgHandler.GetOrgPendingJoinApplicationsHandler)
			adminRoutes.DELETE("/remove-member", orgHandler.RemoveOrganizationMemberHandler)
			adminRoutes.POST("/change-organization", orgHandler.SubmitChangeOrganizationApplicationHandler)

			//组织定位以及邀请码
			adminRoutes.PUT("/:orgID/location", orgHandler.UpdateOrganizationLocationHandler)
			adminRoutes.POST("/:orgID/invite-codes", orgHandler.CreateCustomInviteCodeHandler)

			//组织活动相关
			adminRoutes.POST("/:orgID/activities", orgHandler.CreateOrganizationActivityHandler)
			adminRoutes.PATCH("/:orgID/activities/:activityID/cancel", orgHandler.CancelActivityHandler)
			adminRoutes.PATCH("/:orgID/activities/:activityID", orgHandler.UpdateActivityHandler)
			adminRoutes.GET("/:orgID/activities/assigned", orgHandler.GetOrgAssignedActivitiesHandler)
			adminRoutes.POST("/activities/batch-assign", orgHandler.BatchAssignActivityHandler)
		}

		// 需要组织成员权限的路由
		memberRoutes := organizationGroup.Group("")
		memberRoutes.Use(middleware.CreateAuthChain(authService, orgService, "member")...)
		{
			//查看内部活动
			memberRoutes.GET("/:orgID/activities/internal", orgHandler.GetOrgInternalActivitiesHandler)

			//搜索组织成员
			memberRoutes.GET("/:orgID/users/search", orgHandler.SearchOrganizationUsersHandler)

			//获取活动参与状态
			memberRoutes.GET("/activities/:activityID/participation-status", orgHandler.GetParticipationStatusHandler)
		}

		// 仅需JWT认证的路由（无特定组织权限要求）
		baseAuthRoutes := organizationGroup.Group("")
		baseAuthRoutes.Use(middleware.AuthMiddleware(authService))
		{
			//创建申请
			baseAuthRoutes.POST("/application/create-organization", orgHandler.SubmitCreateOrganizationApplicationHandler)
			baseAuthRoutes.POST("/application/join-organization", orgHandler.SubmitJoinApplicationHandler)

			//处理申请
			baseAuthRoutes.PATCH("/application/:id/process", orgHandler.ProcessApplicationHandler)

			//个人组织总览
			baseAuthRoutes.GET("/my-organizations", orgHandler.GetUserOrganizationOverviewsHandler)

			//用户获取所有申请
			baseAuthRoutes.GET("/application/my-applications", orgHandler.GetUserApplicationsHandler)

			//用户删除申请
			baseAuthRoutes.DELETE("/application/:id", orgHandler.DeleteUserApplicationHandler)

			//用户查找组织
			baseAuthRoutes.GET("/search", orgHandler.SearchOrganizationsHandler)
			baseAuthRoutes.GET("/:orgID", orgHandler.GetOrganizationByIDHandler)
			baseAuthRoutes.GET("/nearby", orgHandler.SearchNearbyOrganizationsHandler)

			//用户快速加入组织
			baseAuthRoutes.POST("/join-with-code", orgHandler.JoinOrganizationWithCodeHandler)

			//用户查询公开组织活动
			baseAuthRoutes.GET("/:orgID/activities/public", orgHandler.GetOrgPublicActivitiesHandler)

			//查看活动的参与者信息
			baseAuthRoutes.GET("/activities/:activityID/participants", orgHandler.GetActivityParticipantsHandler)

			//参与活动
			baseAuthRoutes.POST("/activities/:activityID/participate", orgHandler.ParticipateActivityHandler)

			//未读活动
			baseAuthRoutes.GET("/me/activities/unread", orgHandler.GetUnreadActivitiesHandler)
			baseAuthRoutes.PATCH("/me/activities/:activityID/mark-as-read", orgHandler.MarkActivityAsReadHandler)

			//已经取消的活动
			baseAuthRoutes.GET("/me/activities/cancelled", orgHandler.GetCancelledActivitiesHandler)
			baseAuthRoutes.DELETE("/me/activities/cancelled/:activityID", orgHandler.DeleteCancelledActivityHandler)
		}
	}
}
