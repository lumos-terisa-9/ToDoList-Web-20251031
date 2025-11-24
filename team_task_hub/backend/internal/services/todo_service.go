package services

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"team_task_hub/backend/internal/models"
	"team_task_hub/backend/internal/repositories"
	"time"

	"gorm.io/gorm"
)

type TodoService struct {
	todoRepo *repositories.TodoRepository
}

func NewTodoService(todoRepo *repositories.TodoRepository) *TodoService {
	return &TodoService{todoRepo: todoRepo}
}

// CancelTodoAndChildren 在服务层中使用事务取消待办及其子待办
func (s *TodoService) CancelTodoAndChildren(userID uint, title, description string, createAt time.Time) error {
	// 事务封装
	return s.todoRepo.Transaction(func(tx *gorm.DB) error {
		// 查找根待办
		rootTodo, err := s.todoRepo.FindRootTodoByDetailsWithTx(tx, userID, title, description, createAt)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("未找到匹配的待办事项: 用户ID=%d, 标题='%s'", userID, title)
			}
			return fmt.Errorf("查询失败: %v", err)
		}

		// 更新根待办状态
		if err := s.todoRepo.UpdateStatusWithTx(tx, rootTodo.ID, "cancelled"); err != nil {
			return fmt.Errorf("更新根待办状态失败: %v", err)
		}

		// 如果有子节点，批量更新
		if rootTodo.HasChildren {
			rowsAffected, err := s.todoRepo.BatchUpdateChildrenStatusWithTx(tx, rootTodo.ID, "pending", "cancelled")
			if err != nil {
				return fmt.Errorf("批量更新子待办失败: %v", err)
			}
			fmt.Printf("在事务中成功取消 %d 个子待办\n", rowsAffected)
		}
		return nil
	})
}

// CancelTodoByDetails 根据待办详情取消待办
func (s *TodoService) CancelTodoByDetails(userID uint, title, description string, startTime, endTime time.Time) (bool, string) {
	err := s.todoRepo.UpdateTodoStateByDetails(
		userID, title, description, startTime, endTime, "cancelled",
	)

	if err != nil {
		return false, fmt.Sprintf("取消待办失败: %v", err)
	}

	return true, "待办已成功取消"
}

// CreateTodoRequest 创建待办请求
type CreateTodoRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
	Urgency     string    `json:"urgency" binding:"oneof=low medium high" default:"medium"`
	Category    string    `json:"category" binding:"oneof=work study fun personal" default:"personal"`

	// 重复规则
	RepeatType     string    `json:"repeat_type" binding:"oneof=none daily weekly monthly" default:"none"`
	RepeatInterval int       `json:"repeat_interval" default:"1"`
	RepeatEndDate  time.Time `json:"repeat_end_date"`

	// 子待办日期数组
	ChildDates []string `json:"child_dates"`
}

// CreateTodoResponse 创建待办响应
type CreateTodoResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateTodo 创建待办（支持普通和重复待办）
func (s *TodoService) CreateTodo(userID uint, req *CreateTodoRequest) (*CreateTodoResponse, error) {
	// 验证请求参数
	if err := s.validateCreateRequest(req); err != nil {
		return &CreateTodoResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	// 判断是普通待办还是重复待办
	if req.RepeatType == "none" || len(req.ChildDates) == 0 {
		return s.createSingleTodo(userID, req)
	} else {
		return s.createRepeatingTodo(userID, req)
	}
}

// createSingleTodo 创建普通待办
func (s *TodoService) createSingleTodo(userID uint, req *CreateTodoRequest) (*CreateTodoResponse, error) {
	todo := &models.Todo{
		Title:         req.Title,
		Description:   req.Description,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		Urgency:       req.Urgency,
		Category:      req.Category,
		Status:        "pending",
		ParentID:      0,
		HasChildren:   false,
		RepeatType:    "none",
		CreatorUserID: userID,
	}

	if err := s.todoRepo.Create(todo); err != nil {
		return &CreateTodoResponse{
			Success: false,
			Message: "创建待办失败: " + err.Error(),
		}, err
	}

	return &CreateTodoResponse{
		Success:   true,
		Message:   "待办创建成功",
		CreatedAt: time.Now(),
	}, nil
}

// createRepeatingTodo 创建重复待办
func (s *TodoService) createRepeatingTodo(userID uint, req *CreateTodoRequest) (*CreateTodoResponse, error) {
	response := &CreateTodoResponse{
		Success:   false,
		CreatedAt: time.Now(),
	}

	// 创建父待办,作为抽象待办存在，产生实例待办
	parentTodo := &models.Todo{
		Title:          req.Title,
		Description:    req.Description,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		Urgency:        req.Urgency,
		Category:       req.Category,
		Status:         "pending",
		ParentID:       0,
		HasChildren:    true,
		RepeatType:     req.RepeatType,
		RepeatInterval: req.RepeatInterval,
		RepeatEndDate:  req.RepeatEndDate,
		CreatorUserID:  userID,
	}

	if err := s.todoRepo.Create(parentTodo); err != nil {
		response.Message = "创建父待办失败: " + err.Error()
		return response, err
	}

	// 生成子待办实例
	childTodos, err := s.generateChildTodos(parentTodo, req)
	if err != nil {
		response.Message = "生成子待办失败: " + err.Error()
		return response, err
	}

	// 批量创建子待办
	if len(childTodos) > 0 {
		if err := s.todoRepo.BatchCreate(childTodos); err != nil {
			response.Message = "创建子待办失败: " + err.Error()
			return response, err
		}

		response.Message = fmt.Sprintf("重复待办创建成功，生成 %d 个子实例", len(childTodos))
	} else {
		response.Message = "重复待办创建成功，但未生成子实例"
	}

	response.Success = true
	return response, nil
}

// generateChildTodos 根据前端提供的日期数组生成子待办
func (s *TodoService) generateChildTodos(parent *models.Todo, req *CreateTodoRequest) ([]models.Todo, error) {
	var childTodos []models.Todo

	if len(req.ChildDates) == 0 {
		return childTodos, nil
	}

	// 解析父待办的时间部分
	parentStartTime := parent.StartTime
	parentEndTime := parent.EndTime

	// 提取父待办的时间部分
	parentStartTimeOfDay := time.Date(0, 1, 1,
		parentStartTime.Hour(), parentStartTime.Minute(), parentStartTime.Second(), 0, time.UTC)

	// 计算任务持续时间，用于计算子待办的结束时间
	taskDuration := parentEndTime.Sub(parentStartTime)

	// 处理前端提供的日期数组
	for _, dateStr := range req.ChildDates {
		// 清理和验证日期字符
		dateStr = strings.TrimSpace(dateStr)
		if dateStr == "" {
			continue
		}

		// 解析日期
		childDate, err := ParseDateString(dateStr)
		if err != nil {
			log.Printf("警告: 跳过无效日期 %s: %v", dateStr, err)
			continue
		}

		// 组合日期和时间：使用子待办的日期 + 父待办的时间
		childStartTime := time.Date(
			childDate.Year(), childDate.Month(), childDate.Day(),
			parentStartTimeOfDay.Hour(), parentStartTimeOfDay.Minute(), parentStartTimeOfDay.Second(), 0,
			parentStartTime.Location(),
		)

		childEndTime := childStartTime.Add(taskDuration)

		// 检查是否超过重复结束日期
		if !req.RepeatEndDate.IsZero() && childStartTime.After(req.RepeatEndDate) {
			log.Printf("跳过超过结束日期的实例: %s", childStartTime.Format("2006-01-02"))
			continue
		}

		// 创建子待办实例
		childTodo := models.Todo{
			Title:          parent.Title,
			Description:    parent.Description,
			StartTime:      childStartTime,
			EndTime:        childEndTime,
			CreatedAt:      parent.CreatedAt,
			Urgency:        parent.Urgency,
			Category:       parent.Category,
			Status:         "pending",
			ParentID:       parent.ID,
			HasChildren:    false,
			RepeatType:     parent.RepeatType,
			RepeatInterval: parent.RepeatInterval,
			RepeatEndDate:  parent.RepeatEndDate,
			CreatorUserID:  parent.CreatorUserID,
		}
		childTodos = append(childTodos, childTodo)
	}

	log.Printf("为父待办 %d 生成 %d 个子待办实例", parent.ID, len(childTodos))
	return childTodos, nil
}

// validateCreateRequest 验证创建请求
func (s *TodoService) validateCreateRequest(req *CreateTodoRequest) error {
	if req.Title == "" || req.Description == "" {
		return errors.New("标题和内容不能为空")
	}

	if req.StartTime.IsZero() || req.EndTime.IsZero() {
		return errors.New("开始时间和结束时间不能为空")
	}

	if req.StartTime.After(req.EndTime) {
		return errors.New("开始时间不能晚于结束时间")
	}

	if req.StartTime.Before(time.Now().Add(-24 * time.Hour)) {
		return errors.New("不能创建过去时间的待办")
	}

	return nil
}

// RenewExpiredChildTodos 为有生育能力的过期待办创建下一代
func (s *TodoService) RenewExpiredChildTodos(userID uint) (bool, string) {
	// 查找有生育能力的过期子待办
	expiredTodos, err := s.todoRepo.FindExpiredRenewableTodos(userID)
	if err != nil {
		return false, fmt.Sprintf("查询失败: %v", err)
	}

	if len(expiredTodos) == 0 {
		return true, "未找到需要续期的过期子待办"
	}

	var newInstances []models.Todo
	var renewedIDs []uint
	renewedCount := 0
	skippedCount := 0
	failedCount := 0

	// 为每个待办创建下一代实例
	for _, todo := range expiredTodos {
		result, newTodo, err := s.createNextGeneration(&todo)
		if err != nil {
			failedCount++
			log.Printf("续期失败: %s (ID: %d) - %v", todo.Title, todo.ID, err)
			continue
		}

		switch result {
		case "renewed":
			renewedCount++
			renewedIDs = append(renewedIDs, todo.ID)
			newInstances = append(newInstances, *newTodo)
		case "skipped":
			skippedCount++
		}
	}

	// 批量创建新一代实例
	if len(newInstances) > 0 {
		if err := s.todoRepo.BatchCreate(newInstances); err != nil {
			return false, fmt.Sprintf("创建新实例失败: %v", err)
		}
	}

	// 将已续期的待办设为"无生育能力"
	if len(renewedIDs) > 0 {
		if err := s.todoRepo.BatchUpdateRepeatType(renewedIDs, "none"); err != nil {
			return false, fmt.Sprintf("更新重复规则失败: %v", err)
		}
	}

	message := fmt.Sprintf("续期完成: 成功 %d, 跳过 %d, 失败 %d", renewedCount, skippedCount, failedCount)
	return true, message
}

// createNextGeneration 为待办创建下一代实例
func (s *TodoService) createNextGeneration(parent *models.Todo) (string, *models.Todo, error) {
	// 使用CalculateNextInstance计算下一个实例时间
	nextStart, exists := parent.CalculateNextInstance(time.Now())
	if !exists {
		return "skipped", nil, fmt.Errorf("无后续实例可用")
	}

	// 计算任务持续时间
	taskDuration := parent.EndTime.Sub(parent.StartTime)
	nextEnd := nextStart.Add(taskDuration)

	// 创建新一代实例
	newTodo := &models.Todo{
		Title:          parent.Title,
		Description:    parent.Description,
		StartTime:      nextStart,
		EndTime:        nextEnd,
		CreatedAt:      parent.CreatedAt,
		Status:         "pending",
		Urgency:        parent.Urgency,
		Category:       parent.Category,
		ParentID:       parent.ParentID,
		HasChildren:    false,
		RepeatType:     parent.RepeatType,
		RepeatInterval: parent.RepeatInterval,
		RepeatEndDate:  parent.RepeatEndDate,
		CreatorUserID:  parent.CreatorUserID,
		CreatorOrganID: parent.CreatorOrganID,
	}

	return "renewed", newTodo, nil
}

// CompleteTodoByDetails 完成待办（自动设置完成时间）
func (s *TodoService) CompleteTodoByDetails(userID uint, title, description string, startTime, endTime time.Time) (bool, string) {
	// 参数验证
	if err := s.validateCompleteRequest(title, description, startTime, endTime); err != nil {
		return false, err.Error()
	}

	completedAt := time.Now()

	// 调用仓库层更新状态和完成时间
	err := s.todoRepo.UpdateTodoStateAndCompletedTime(
		userID, title, description, startTime, endTime, "completed", completedAt,
	)

	if err != nil {
		return false, fmt.Sprintf("完成待办失败: %v", err)
	}

	return true, fmt.Sprintf("待办已完成于 %s", completedAt.Format("2006-01-02 15:04"))
}

// validateCompleteRequest 验证完成待办请求参数
func (s *TodoService) validateCompleteRequest(title, description string, startTime, endTime time.Time) error {

	if title == "" || description == "" {
		return fmt.Errorf("标题和内容不能为空")
	}

	if startTime.IsZero() || endTime.IsZero() {
		return fmt.Errorf("开始时间和结束时间不能为空")
	}

	if startTime.After(endTime) {
		return fmt.Errorf("开始时间不能晚于结束时间")
	}

	return nil
}

// GetTodayTodos 查找今日待办
func (s *TodoService) GetTodayTodos(userID uint) ([]models.Todo, error) {

	startOfToday, endOfToday := TodayRange()

	todos, err := s.todoRepo.FindTodosInTimeRange(userID, startOfToday, endOfToday)
	if err != nil {
		return nil, fmt.Errorf("查询今日待办失败: %v", err)
	}

	return todos, nil
}

// GetCompletedTodosByDate 获取指定日期完成的所有待办
func (s *TodoService) GetCompletedTodosByDate(userID uint, dateStr string) ([]models.Todo, error) {
	// 解析日期字符串
	date, err := ParseDateString(dateStr)
	if err != nil {
		return nil, err
	}

	//获取指定日期的开始和结束
	startOfDay, endOfDay := DayRange(date)

	// 调用仓库层函数查询在指定日期范围内完成的任务
	todos, err := s.todoRepo.FindCompletedInRange(userID, startOfDay, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("查询已完成待办失败: %v", err)
	}

	return todos, nil
}

// GetTodosStartingInNext7Days 获取未来七天内会开始的待办
func (s *TodoService) GetTodosStartingInNext7Days(userID uint) ([]models.Todo, error) {
	// 获取当前时间和未来七天的时间范围
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startTime = startTime.AddDate(0, 0, 1)
	endTime := startTime.AddDate(0, 0, 7)

	todos, err := s.todoRepo.FindTodosStartingInRange(userID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("查询未来七天开始的待办失败: %v", err)
	}

	return todos, nil
}

// GetTodosEndingInNext7Days 获取未来七天内会结束的待办
func (s *TodoService) GetTodosEndingInNext7Days(userID uint) ([]models.Todo, error) {
	// 获取当前时间和未来七天的时间范围
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endTime := now.AddDate(0, 0, 7)

	todos, err := s.todoRepo.FindTodosEndingInRange(userID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("查询未来七天结束的待办失败: %v", err)
	}

	return todos, nil
}
