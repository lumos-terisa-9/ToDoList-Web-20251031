// repositories/todo_repository.go
package repositories

import (
	"team_task_hub/backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(todo *models.Todo) error {
	return r.db.Create(todo).Error
}

func (r *TodoRepository) FindByID(id uint) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.First(&todo, id).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *TodoRepository) Update(id uint, updates map[string]any) error {
	result := r.db.Model(&models.Todo{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *TodoRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Todo{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// 查找用户所有未完成的并且仍然有效的任务
func (r *TodoRepository) FindActiveInstances(userID uint) ([]models.Todo, error) {
	var todos []models.Todo
	now := time.Now()

	err := r.db.
		Where("creator_user_id = ?", userID).
		Where("status = 'pending'").
		Where("has_children = false").
		Where("start_time <= ?", now).
		Where("end_time >= ?", now).
		Order("start_time ASC").
		Find(&todos).Error

	return todos, err
}

// FindInstancesStartingInRange 查询在指定时间段内开始的所有任务
func (r *TodoRepository) FindTodosStartingInRange(userID uint, startTime, endTime time.Time) ([]models.Todo, error) {
	var todos []models.Todo
	err := r.db.
		Where("creator_user_id = ?", userID).
		Where("status = 'pending'").
		Where("has_children = false").
		Where("start_time >= ? AND start_time < ?", startTime, endTime).
		Order("start_time ASC").
		Find(&todos).Error
	return todos, err
}

// FindInstancesEndingInRange 查询在指定时间段内结束的所有任务
func (r *TodoRepository) FindTodosEndingInRange(userID uint, startTime, endTime time.Time) ([]models.Todo, error) {
	var todos []models.Todo

	err := r.db.
		Where("creator_user_id = ?", userID).
		Where("status = 'pending'").
		Where("has_children = false").
		Where("end_time >= ? AND end_time <= ?", startTime, endTime).
		Order("end_time ASC").
		Find(&todos).Error
	return todos, err
}

// FindInstancesInTimeRange 在指定时间段内开始或结束的任务
func (r *TodoRepository) FindTodosInTimeRange(userID uint, startTime, endTime time.Time) ([]models.Todo, error) {
	var todos []models.Todo

	err := r.db.
		Where("creator_user_id = ?", userID).
		Where("status = 'pending'").
		Where("has_children = false").
		Where("(start_time BETWEEN ? AND ?) OR (end_time BETWEEN ? AND ?) OR (start_time<? and end_time>?)",
			startTime, endTime, startTime, endTime, startTime, endTime).
		Order("start_time ASC").
		Find(&todos).Error

	return todos, err
}

// FindOneDayExpiredTodos 查找某一天已经过期的待办
func (r *TodoRepository) FindOneDayExpiredTodos(userID uint, startTime, endTime time.Time) ([]models.Todo, error) {
	var todos []models.Todo
	err := r.db.
		Where("creator_user_id = ?", userID).
		Where("status = 'pending'").
		Where("has_children = false").
		Where("endTime between ? and ?", startTime, endTime).
		Order("endTime ASC").
		Find(&todos).Error
	return todos, err
}

// FindCompletedInRange 查询在指定时间段内完成的任务
func (r *TodoRepository) FindCompletedInRange(userID uint, startTime, endTime time.Time) ([]models.Todo, error) {
	var todos []models.Todo

	err := r.db.
		Where("creator_user_id = ?", userID).
		Where("status = 'completed'").
		Where("completed_at >= ? AND completed_at <= ?", startTime, endTime).
		Order("completed_at ASC").
		Find(&todos).Error

	return todos, err
}

// FindRootTodoByDetails 根据用户ID、标题、内容查找根待办事项
func (r *TodoRepository) FindRootTodoByDetails(userID uint, title, description string) (*models.Todo, error) {
	var todo models.Todo

	err := r.db.
		Where("creator_user_id = ?", userID).
		Where("title = ?", title).
		Where("description = ?", description).
		Where("parent_id = 0").
		First(&todo).Error

	if err != nil {
		return nil, err
	}
	return &todo, nil
}

// FindTodoByDetails 根据用户ID、时间范围和内容精确查找待办事项
func (r *TodoRepository) FindTodoByDetails(userID uint, startTime, endTime time.Time, title, description string) (*models.Todo, error) {
	var todo models.Todo

	err := r.db.
		Where("creator_user_id = ?", userID).
		Where("title = ?", title).
		Where("description = ?", description).
		Where("start_time = ?", startTime).
		Where("end_time = ?", endTime).
		First(&todo).Error

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// UpdateStatus 更新单个待办事项的状态
func (r *TodoRepository) UpdateStatus(id uint, status string) error {
	result := r.db.Model(&models.Todo{}).
		Where("id = ?", id).
		Update("status", status)
	return result.Error
}

// UpdateTodoStateByDetails 根据多条件更新待办状态
func (r *TodoRepository) UpdateTodoStateByDetails(userID uint, title, description string, startTime, endTime time.Time, newState string) error {
	result := r.db.Model(&models.Todo{}).
		Where("creator_user_id = ?", userID).
		Where("title = ?", title).
		Where("description = ?", description).
		Where("start_time = ?", startTime).
		Where("end_time = ?", endTime).
		Update("status", newState)

	return result.Error
}

// UpdateTodoStateAndCompletedTime 更新待办状态和完成时间
func (r *TodoRepository) UpdateTodoStateAndCompletedTime(userID uint, title, description string, startTime, endTime time.Time, newState string, completedAt time.Time) error {
	result := r.db.Model(&models.Todo{}).
		Where("creator_user_id = ?", userID).
		Where("title = ?", title).
		Where("description = ?", description).
		Where("start_time = ?", startTime).
		Where("end_time = ?", endTime).
		Updates(map[string]any{
			"status":       newState,
			"completed_at": completedAt,
		})

	return result.Error
}

// BatchUpdateChildrenStatus 批量更新子待办事项状态
func (r *TodoRepository) BatchUpdateChildrenStatus(parentID uint, oldStatus, newStatus string) (int64, error) {
	result := r.db.Model(&models.Todo{}).
		Where("parent_id = ? AND status = ?", parentID, oldStatus).
		Update("status", newStatus)
	return result.RowsAffected, result.Error
}

// FindRootTodoByDetailsWithTx 支持事务的版本：根据用户ID、标题、内容，创建时间查找根待办事项
func (r *TodoRepository) FindRootTodoByDetailsWithTx(tx *gorm.DB, userID uint, title, description string, creatAt time.Time) (*models.Todo, error) {
	var todo models.Todo
	// 使用传入的 tx 进行查询，而非 r.db
	err := tx.
		Where("creator_user_id = ?", userID).
		Where("title = ?", title).
		Where("description = ?", description).
		Where("created_at=?", creatAt).
		Where("parent_id = 0").
		First(&todo).Error

	if err != nil {
		return nil, err
	}
	return &todo, nil
}

// UpdateStatusWithTx 支持事务的版本：更新单个待办事项的状态
func (r *TodoRepository) UpdateStatusWithTx(tx *gorm.DB, id uint, status string) error {
	// 使用传入的 tx 进行更新
	result := tx.Model(&models.Todo{}).
		Where("id = ?", id).
		Update("status", status)
	return result.Error
}

// BatchUpdateChildrenStatusWithTx 支持事务的版本：批量更新子待办事项状态
func (r *TodoRepository) BatchUpdateChildrenStatusWithTx(tx *gorm.DB, parentID uint, oldStatus, newStatus string) (int64, error) {
	// 使用传入的 tx 进行批量更新
	result := tx.Model(&models.Todo{}).
		Where("parent_id = ? AND status = ?", parentID, oldStatus).
		Update("status", newStatus)
	return result.RowsAffected, result.Error
}

// FindExpiredRenewableTodos 查找有"生育能力"的过期子待办
func (r *TodoRepository) FindExpiredRenewableTodos(userID uint) ([]models.Todo, error) {
	var todos []models.Todo
	now := time.Now()
	future := now.Add(7 * 24 * time.Hour)
	err := r.db.
		Where("creator_user_id = ?", userID).
		Where("parent_id != 0").
		Where("status != ?", "cancelled").
		Where("repeat_type != ?", "none").
		Where("end_time < ?", future).
		Find(&todos).Error

	return todos, err
}

// BatchUpdateRepeatType 批量更新重复类型
func (r *TodoRepository) BatchUpdateRepeatType(ids []uint, repeatType string) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Model(&models.Todo{}).
		Where("id IN (?)", ids).
		Update("repeat_type", repeatType).Error
}

// 查询根节点（抽象事件或没有重复事件）
func (r *TodoRepository) FindRootTodos(userID uint) ([]models.Todo, error) {
	var roots []models.Todo
	err := r.db.
		Where("creator_user_id = ? AND parent_id = 0", userID).
		Find(&roots).Error
	return roots, err
}

// 根据抽象事件查找实例事件
func (r *TodoRepository) FindChildrenByParent(parentID uint) ([]models.Todo, error) {
	var children []models.Todo
	err := r.db.
		Where("parent_id = ?", parentID).
		Order("start_time ASC").
		Find(&children).Error
	return children, err
}

func (r *TodoRepository) FindByStatus(userID uint, status string) ([]models.Todo, error) {
	var todos []models.Todo
	err := r.db.
		Where("creator_user_id = ? AND status = ?", userID, status).
		Where("has_children = false").
		Order("start_time ASC").
		Find(&todos).Error
	return todos, err
}

// 批量操作
func (r *TodoRepository) BatchCreate(todos []models.Todo) error {
	if len(todos) == 0 {
		return nil
	}
	return r.db.CreateInBatches(todos, 100).Error
}

// 事务支持
func (r *TodoRepository) Transaction(fn func(*gorm.DB) error) error {
	return r.db.Transaction(fn)
}
