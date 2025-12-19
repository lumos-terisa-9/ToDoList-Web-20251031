package repositories

import (
	"fmt"
	"team_task_hub/backend/internal/models"

	"gorm.io/gorm"
)

// ActivityParticipationRepository 活动参与数据访问层
// 负责所有与activity_participation表相关的数据库操作
type ActivityParticipationRepository struct {
	db *gorm.DB
}

// NewActivityParticipationRepository 构造函数：创建ActivityParticipationRepository实例
func NewActivityParticipationRepository(db *gorm.DB) *ActivityParticipationRepository {
	return &ActivityParticipationRepository{db: db}
}

// Create 创建新的活动参与记录
func (r *ActivityParticipationRepository) Create(participation *models.ActivityParticipation) error {
	result := r.db.Create(participation)
	if result.Error != nil {
		return fmt.Errorf("创建活动参与记录失败: %v", result.Error)
	}
	return nil
}

// GetByID 根据主键ID查找活动参与记录
func (r *ActivityParticipationRepository) GetByID(id uint) (*models.ActivityParticipation, error) {
	var participation models.ActivityParticipation
	result := r.db.Preload("Activity").Preload("User").First(&participation, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("活动参与记录不存在")
		}
		return nil, fmt.Errorf("查询活动参与记录失败: %v", result.Error)
	}
	return &participation, nil
}

// GetByActivityID 根据活动ID查找所有参与记录
func (r *ActivityParticipationRepository) GetByActivityID(activityID uint) ([]models.ActivityParticipation, error) {
	var participations []models.ActivityParticipation
	result := r.db.Where("activity_id = ?", activityID).Preload("User").Find(&participations)
	if result.Error != nil {
		return nil, fmt.Errorf("根据活动查询参与记录失败: %v", result.Error)
	}
	return participations, nil
}

// GetByUserID 根据用户ID查找其所有活动参与记录
func (r *ActivityParticipationRepository) GetByUserID(userID uint) ([]models.ActivityParticipation, error) {
	var participations []models.ActivityParticipation
	result := r.db.Where("user_id = ?", userID).Preload("Activity").Find(&participations)
	if result.Error != nil {
		return nil, fmt.Errorf("根据用户查询活动参与记录失败: %v", result.Error)
	}
	return participations, nil
}

// GetByActivityAndUser 根据活动ID和用户ID查找特定的参与记录
func (r *ActivityParticipationRepository) GetByActivityAndUser(activityID, userID uint) (*models.ActivityParticipation, error) {
	var participation models.ActivityParticipation
	result := r.db.Where("activity_id = ? AND user_id = ?", activityID, userID).Preload("Activity").Preload("User").First(&participation)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询用户活动参与记录失败: %v", result.Error)
	}
	return &participation, nil
}

// Update 更新活动参与记录（如更新评分和评价）
func (r *ActivityParticipationRepository) Update(participation *models.ActivityParticipation) error {
	result := r.db.Save(participation)
	if result.Error != nil {
		return fmt.Errorf("更新活动参与记录失败: %v", result.Error)
	}
	return nil
}

// Delete 删除活动参与记录
func (r *ActivityParticipationRepository) Delete(id uint) error {
	result := r.db.Delete(&models.ActivityParticipation{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除活动参与记录失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("活动参与记录不存在")
	}
	return nil
}

// Exists 检查用户是否已参与某个活动
func (r *ActivityParticipationRepository) Exists(activityID, userID uint) (bool, error) {
	var count int64
	result := r.db.Model(&models.ActivityParticipation{}).Where("activity_id = ? AND user_id = ?", activityID, userID).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("检查活动参与记录存在性失败: %v", result.Error)
	}
	return count > 0, nil
}

// CancelByActivityID 取消某个活动的所有参与记录（软删除，更新状态）
func (r *ActivityParticipationRepository) CancelByActivityID(activityID uint) error {
	result := r.db.Model(&models.ActivityParticipation{}).
		Where("activity_id = ?", activityID).
		Updates(map[string]any{
			"status": "cancelled",
		})
	if result.Error != nil {
		return fmt.Errorf("取消活动参与记录失败: %v", result.Error)
	}
	return nil
}

// CancelByActivityAndUser 取消特定用户对某个活动的参与（软删除，更新状态）
func (r *ActivityParticipationRepository) CancelByActivityAndUser(activityID, userID uint) error {
	result := r.db.Model(&models.ActivityParticipation{}).
		Where("activity_id = ? AND user_id = ?", activityID, userID).
		Updates(map[string]any{
			"status": "cancelled",
		})
	if result.Error != nil {
		return fmt.Errorf("取消用户活动参与失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到该用户对指定活动的参与记录")
	}
	return nil
}
