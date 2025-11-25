package router

import (
	"team_task_hub/backend/internal/handlers"
	"team_task_hub/backend/internal/middleware"
	"team_task_hub/backend/internal/repositories"
	"team_task_hub/backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupTodoRoutes 设置待办路由
func SetupTodoRoutes(router *gin.Engine, db *gorm.DB, authService *services.AuthService) {
	// 待办事项路由组
	todoRepo := repositories.NewTodoRepository(db)
	todoService := services.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)

	todoGroup := router.Group("/api/todos")
	todoGroup.Use(middleware.AuthMiddleware(authService))
	{
		todoGroup.POST("/createTodo", todoHandler.AddTodoHandler)
		todoGroup.POST("/updateTodos", todoHandler.UpdateTodos)
		todoGroup.POST("/cancel", todoHandler.CancelTodoByDetails)
		todoGroup.POST("/cancel-with-children", todoHandler.CancelTodoAndChildren)
		todoGroup.POST("/complete", todoHandler.CompleteTodoByDetails)
		todoGroup.GET("/todayTodos", todoHandler.GetTodayTodos)
		todoGroup.GET("/completed_todo", todoHandler.GetCompletedTodosByDate)
		todoGroup.GET("/coming-startTodos", todoHandler.GetComingStartTodos)
		todoGroup.GET("/coming-endTodos", todoHandler.GetComingEndTodos)
		todoGroup.POST("/cancel-completedTodo", todoHandler.CancelCompletedTodo)
		todoGroup.GET("/Get-OneDayTodos", todoHandler.GetOneDayTodos)
		todoGroup.GET("/get-OneDayExpiredTodos", todoHandler.GetOneDayExpiredTodos)
	}
}
