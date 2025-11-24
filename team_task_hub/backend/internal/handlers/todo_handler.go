// handlers/todo_handler.go
package handlers

import (
	"net/http"
	"strings"
	"team_task_hub/backend/internal/models"
	"team_task_hub/backend/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)

// TodoHandler 待办处理器
type TodoHandler struct {
	todoService *services.TodoService
}

// NewTodoHandler 构造函数
func NewTodoHandler(todoService *services.TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
	}
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CreateTodo 创建待办事项
// @Summary 创建待办事项
// @Description 创建普通待办或重复待办事项。重复待办会根据规则及ChildDates生成子实例。
// @Tags 待办事项
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body services.CreateTodoRequest true "创建待办请求"
// @Success 201 {object} services.CreateTodoResponse "创建成功" example({"success": true, "message": "待办创建成功", "created_at": "2024-01-15T10:30:00Z"})
// @Failure 400 {object} string "请求参数错误" example({"success": false, "message": "标题不能为空"})
// @Failure 500 {object} string "系统内部错误" example({"success": false, "message": "创建待办失败: 数据库连接错误"})
// @Router /api/todos/createTodo [post]
func (h *TodoHandler) AddTodoHandler(c *gin.Context) {
	var req services.CreateTodoRequest

	// 绑定JSON请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 从认证中间件获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "用户未认证",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "用户ID格式错误",
		})
		return
	}

	// 调用服务层创建待办
	response, err := h.todoService.CreateTodo(userIDUint, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "创建待办失败: " + err.Error(),
		})
		return
	}

	// 根据操作结果设置HTTP状态码
	statusCode := http.StatusCreated
	if !response.Success {
		statusCode = http.StatusBadRequest
	}

	c.JSON(statusCode, response)
}

// RenewTodoResponse 续期响应结构
type RenewTodoResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// RenewExpiredTodos 续期过期待办
// @Summary 续期过期待办
// @Description 为有生育能力的过期待办自动创建下一代实例
// @Tags 待办事项
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Success 200 {object} RenewTodoResponse "续期成功" example({"success": true, "message": "续期完成: 成功 3, 跳过 1, 失败 0"})
// @Failure 400 {object} string "请求参数错误" example({"success": false, "message": "用户未认证"})
// @Failure 500 {object} string "系统内部错误" example({"success": false, "message": "续期操作失败: 数据库连接错误"})
// @Router /api/todos/UpdateTodos [post]
func (h *TodoHandler) UpdateTodos(c *gin.Context) {
	// 从认证中间件获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "用户未认证",
		})
		return
	}

	// 验证用户ID类型
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "用户ID格式错误",
		})
		return
	}

	// 调用服务层续期功能
	success, message := h.todoService.RenewExpiredChildTodos(userIDUint)

	statusCode := http.StatusOK
	if !success {
		statusCode = http.StatusInternalServerError
	}

	c.JSON(statusCode, gin.H{
		"success": success,
		"message": message,
	})
}

// CancelTodoRequest 取消待办请求结构
type CancelTodoRequest struct {
	Title       string    `json:"title" binding:"required" example:"团队会议"`
	Description string    `json:"description,omitempty" example:"项目进度讨论"`
	StartTime   time.Time `json:"start_time" binding:"required" example:"2024-01-15T14:00:00Z"`
	EndTime     time.Time `json:"end_time" binding:"required" example:"2024-01-15T15:00:00Z"`
}

// CancelTodoResponse 取消待办响应结构
type CancelTodoResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CancelTodoByDetails 根据详情取消待办
// @Summary 取消待办事项
// @Description 根据待办的具体信息将其状态改为已取消
// @Tags 待办事项
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body CancelTodoRequest true "取消待办请求"
// @Success 200 {object} CancelTodoResponse "取消成功" example({"success": true, "message": "待办已成功取消"})
// @Failure 400 {object} string "请求参数错误" example({"success": false, "message": "请求参数格式错误"})
// @Failure 401 {object} string "未授权" example({"success": false, "message": "用户未认证"})
// @Failure 500 {object} string "系统内部错误" example({"success": false, "message": "取消待办失败: 数据库错误"})
// @Router /api/todos/cancel [post]
func (h *TodoHandler) CancelTodoByDetails(c *gin.Context) {
	var req CancelTodoRequest

	// 绑定JSON请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 从认证中间件获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "用户未认证",
		})
		return
	}

	// 验证用户ID类型
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "用户ID格式错误",
		})
		return
	}

	// 调用服务层取消待办
	success, message := h.todoService.CancelTodoByDetails(
		userIDUint, req.Title, req.Description, req.StartTime, req.EndTime,
	)

	// 设置HTTP状态码
	statusCode := http.StatusOK
	if !success {
		statusCode = http.StatusInternalServerError
	}

	c.JSON(statusCode, gin.H{
		"success": success,
		"message": message,
	})
}

// CancelTodoAndChildrenRequest 取消待办及其子待办请求结构
type CancelTodoAndChildrenRequest struct {
	Title       string    `json:"title" binding:"required" example:"每周团队会议"`
	Description string    `json:"description,omitempty" example:"项目进度讨论"`
	CreatedAt   time.Time `json:"createdAt"`
}

// CancelTodoAndChildren 取消待办及其子待办
// @Summary 取消待办及其所有子待办
// @Description 根据待办标题和描述取消父待办及其所有子待办（使用事务确保一致性）
// @Tags 待办事项
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body CancelTodoAndChildrenRequest true "取消待办请求"
// @Success 200 {object} CancelTodoResponse "取消成功" example({"success": true, "message": "待办及其子待办已成功取消"})
// @Failure 400 {object} string "请求参数错误" example({"success": false, "message": "请求参数格式错误"})
// @Failure 404 {object} string "未找到待办" example({"success": false, "message": "未找到匹配的待办事项"})
// @Failure 500 {object} string "系统内部错误" example({"success": false, "message": "取消操作失败: 数据库错误"})
// @Router /api/todos/cancel-with-children [post]
func (h *TodoHandler) CancelTodoAndChildren(c *gin.Context) {
	var req CancelTodoAndChildrenRequest

	// 绑定JSON请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 从认证中间件获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "用户未认证",
		})
		return
	}

	// 验证用户ID类型
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "用户ID格式错误",
		})
		return
	}

	// 调用服务层取消待办及其子待办
	err := h.todoService.CancelTodoAndChildren(userIDUint, req.Title, req.Description, req.CreatedAt)

	if err != nil {
		// 根据错误类型设置不同的HTTP状态码
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()

		if strings.Contains(errorMsg, "未找到匹配的待办事项") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(errorMsg, "请求参数") {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": errorMsg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "待办及其所有子待办已成功取消",
	})
}

// CompleteTodoRequest 完成待办请求结构
type CompleteTodoRequest struct {
	Title       string    `json:"title" binding:"required" example:"团队会议"`
	Description string    `json:"description" binding:"required" example:"项目进度讨论"`
	StartTime   time.Time `json:"start_time" binding:"required" example:"2024-01-15T14:00:00Z"`
	EndTime     time.Time `json:"end_time" binding:"required" example:"2024-01-15T15:00:00Z"`
}

// CompleteTodoResponse 完成待办响应结构
type CompleteTodoResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CompleteTodoByDetails 完成待办事项
// @Summary 完成待办事项
// @Description 根据待办的具体信息将其状态改为已完成，并自动记录完成时间
// @Tags 待办事项
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body CompleteTodoRequest true "完成待办请求"
// @Success 200 {object} CompleteTodoResponse "完成成功" example({"success": true, "message": "待办已完成于 2024-01-15 14:30"})
// @Failure 400 {object} string "请求参数错误" example({"success": false, "message": "标题和内容不能为空"})
// @Failure 500 {object} string "系统内部错误" example({"success": false, "message": "完成待办失败: 数据库错误"})
// @Router /api/todos/complete [post]
func (h *TodoHandler) CompleteTodoByDetails(c *gin.Context) {
	var req CompleteTodoRequest

	// 绑定JSON请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 从认证中间件获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "用户未认证",
		})
		return
	}

	// 验证用户ID类型
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "用户ID格式错误",
		})
		return
	}

	// 调用服务层完成待办
	success, message := h.todoService.CompleteTodoByDetails(
		userIDUint, req.Title, req.Description, req.StartTime, req.EndTime,
	)

	if !success {
		// 根据错误消息设置HTTP状态码
		statusCode := http.StatusInternalServerError
		if strings.Contains(message, "不能为空") || strings.Contains(message, "不能晚于") {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
	})
}

type TodoListResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Todos   []models.Todo `json:"todos"`
	Count   int           `json:"count"`
}

// GetTodayTodos 获取今日待办事项
// @Summary 获取今日待办事项
// @Description 查询与今天时间段有交集的所有未完成待办事项
// @Tags 待办查询
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Success 200 {object} TodoListResponse "查询成功" example({"success": true, "message": "查询成功", "todos": [...], "count": 5})
// @Failure 401 {object} string "未授权" example({"success": false, "message": "用户未认证"})
// @Failure 500 {object} string "系统内部错误" example({"success": false, "message": "查询今日待办失败: 数据库错误"})
// @Router /api/todos/todayTodos [get]
func (h *TodoHandler) GetTodayTodos(c *gin.Context) {
	// 从认证中间件获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "用户未认证",
		})
		return
	}

	// 验证用户ID类型
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "用户ID格式错误",
		})
		return
	}

	// 调用服务层获取今日待办
	todos, err := h.todoService.GetTodayTodos(userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "查询今日待办失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "查询成功",
		"todos":   todos,
		"count":   len(todos),
	})
}

// CompletedTodoResponse 已完成待办响应结构
type CompletedTodoResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Date    string        `json:"date"`
	Todos   []models.Todo `json:"todos"`
	Count   int           `json:"count"`
}

// GetCompletedTodosByDate 获取指定日期完成的待办事项
// @Summary 获取指定日期完成的待办事项
// @Description 查询在指定日期完成的所有待办事项
// @Tags 待办查询
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param date query string true "查询日期" example("2024-01-15")
// @Success 200 {object} CompletedTodoResponse "查询成功" example({"success": true, "message": "查询成功", "date": "2024-01-15", "todos": [...], "count": 3})
// @Failure 400 {object} string "请求参数错误" example({"success": false, "message": "日期参数不能为空"})
// @Failure 401 {object} string "未授权" example({"success": false, "message": "用户未认证"})
// @Failure 500 {object} string "系统内部错误" example({"success": false, "message": "查询已完成待办失败: 数据库错误"})
// @Router /api/todos/completed_todo [get]
func (h *TodoHandler) GetCompletedTodosByDate(c *gin.Context) {
	// 从查询参数获取日期
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "日期参数不能为空",
		})
		return
	}

	// 从认证中间件获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "用户未认证",
		})
		return
	}

	// 验证用户ID类型
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "用户ID格式错误",
		})
		return
	}

	// 调用服务层获取指定日期完成的待办
	todos, err := h.todoService.GetCompletedTodosByDate(userIDUint, dateStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "查询已完成待办失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "查询成功",
		"date":    dateStr,
		"todos":   todos,
		"count":   len(todos),
	})
}

// GetUpcomingStartingTodos 获取未来七天会开始的待办
// @Summary 获取未来七天会开始的待办事项
// @Description 查询在未来七天内会开始的所有待办事项
// @Tags 待办查询
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Success 200 {object} TodoListResponse "查询成功" example({"success": true, "message": "查询成功", "todos": [...], "count": 5})
// @Failure 401 {object} string "未授权" example({"success": false, "message": "用户未认证"})
// @Failure 500 {object} string "系统内部错误" example({"success": false, "message": "查询未来七天开始的待办失败: 数据库错误"})
// @Router /api/todos/coming-startTodos [get]
func (h *TodoHandler) GetComingStartTodos(c *gin.Context) {
	// 从认证中间件获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "用户未认证",
		})
		return
	}

	// 验证用户ID类型
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "用户ID格式错误",
		})
		return
	}

	// 调用服务层获取未来七天会开始的待办
	todos, err := h.todoService.GetTodosStartingInNext7Days(userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "查询未来七天开始的待办失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "查询成功",
		"todos":   todos,
		"count":   len(todos),
	})
}

// GetUpcomingEndingTodos 获取未来七天会结束的待办
// @Summary 获取未来七天会结束的待办事项
// @Description 查询在未来七天内会结束的所有待办事项
// @Tags 待办查询
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Success 200 {object} TodoListResponse "查询成功" example({"success": true, "message": "查询成功", "todos": [...], "count": 3})
// @Failure 401 {object} string "未授权" example({"success": false, "message": "用户未认证"})
// @Failure 500 {object} string "系统内部错误" example({"success": false, "message": "查询未来七天结束的待办失败: 数据库错误"})
// @Router /api/todos/coming-endTodos [get]
func (h *TodoHandler) GetComingEndTodos(c *gin.Context) {
	// 从认证中间件获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "用户未认证",
		})
		return
	}

	// 验证用户ID类型
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "用户ID格式错误",
		})
		return
	}

	// 调用服务层获取未来七天会结束的待办
	todos, err := h.todoService.GetTodosEndingInNext7Days(userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "查询未来七天结束的待办失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "查询成功",
		"todos":   todos,
		"count":   len(todos),
	})
}
