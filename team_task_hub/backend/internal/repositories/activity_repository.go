package repositories

import (
	"fmt"
	"team_task_hub/backend/internal/models"
	"time"

	"gorm.io/gorm"
)

// ActivityRepository 活动数据访问层
// 负责所有与activities表相关的数据库操作
type ActivityRepository struct {
	db *gorm.DB
}

// NewActivityRepository 构造函数：创建ActivityRepository实例
func NewActivityRepository(db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

// Create 创建新活动
func (r *ActivityRepository) Create(activity *models.Activity) error {
	result := r.db.Create(activity)
	if result.Error != nil {
		return fmt.Errorf("创建活动失败: %v", result.Error)
	}
	return nil
}

// GetByID 根据主键ID查找活动
func (r *ActivityRepository) GetByID(id uint) (*models.Activity, error) {
	var activity models.Activity
	result := r.db.Preload("Organization").First(&activity, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("活动不存在")
		}
		return nil, fmt.Errorf("查询活动失败: %v", result.Error)
	}
	return &activity, nil
}

// GetByOrganizationID 根据组织ID查找活动列表
func (r *ActivityRepository) GetByOrganizationID(organizationID uint) ([]models.Activity, error) {
	var activities []models.Activity
	result := r.db.Where("organization_id = ?", organizationID).Preload("Organization").Find(&activities)
	if result.Error != nil {
		return nil, fmt.Errorf("根据组织查询活动失败: %v", result.Error)
	}
	return activities, nil
}

// Update 更新活动信息
func (r *ActivityRepository) Update(activity *models.Activity) error {
	result := r.db.Save(activity)
	if result.Error != nil {
		return fmt.Errorf("更新活动失败: %v", result.Error)
	}
	return nil
}

// Delete 删除活动
func (r *ActivityRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Activity{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除活动失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("活动不存在")
	}
	return nil
}

// UpdateStatus 更新活动状态及相关字段
func (r *ActivityRepository) UpdateStatus(activityID uint, updateData map[string]any) error {
	// 使用 Updates 方法同时更新多个字段
	result := r.db.Model(&models.Activity{}).Where("id = ?", activityID).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("活动不存在")
	}
	return nil
}

// FindByOrgAndParticipation 根据组织ID和参与类型查询活动
func (r *ActivityRepository) FindByOrgAndParticipation(orgID uint, participationType string) ([]models.Activity, error) {
	var activities []models.Activity

	// 在数据访问层内部构建查询条件
	err := r.db.Model(&models.Activity{}).
		Where("organization_id = ?", orgID).
		Where("status = ?", "active").
		Where("end_time > ?", time.Now()).
		Where("participation_limit = ?", participationType).
		Find(&activities).Error

	if err != nil {
		return nil, err
	}

	return activities, nil
}

// 查询某个活动的用户信息
func (r *ActivityParticipationRepository) FindUsersByActivityID(activityID uint) ([]models.UserInfo, error) {
	var userInfos []models.UserInfo

	// 构建查询，通过JOIN直接关联到users表并选择特定字段
	db := r.db.Model(&models.ActivityParticipation{}).
		Joins("LEFT JOIN users ON activity_participations.user_id = users.id").
		Where("activity_participations.activity_id = ?", activityID)

	// 直接查询并扫描到userInfos中
	err := db.Select("users.id, users.username, users.avatar_url").Scan(&userInfos).Error
	if err != nil {
		return nil, err
	}
	return userInfos, nil
}

// FindUnreadActivitiesByUserID 根据用户ID查询所有未读活动
func (r *ActivityRepository) FindUnreadActivitiesByUserID(userID uint) ([]models.Activity, error) {
	var activities []models.Activity

	// 直接查询用户所有未读活动
	err := r.db.Model(&models.Activity{}).
		Select("activities.*").
		Joins("INNER JOIN activity_participations ON activities.id = activity_participations.activity_id").
		Where("activity_participations.user_id = ? AND activity_participations.is_unread = ?", userID, true).
		Find(&activities).Error

	return activities, err
}

// FindCancelledActivitiesByUserID 查询用户已取消的活动
func (r *ActivityRepository) FindCancelledActivitiesByUserID(userID uint) ([]models.Activity, error) {
	var activities []models.Activity

	// 查询用户参与的活动，且状态为cancelled
	err := r.db.Model(&models.Activity{}).
		Select("activities.*").
		Joins("INNER JOIN activity_participations ON activities.id = activity_participations.activity_id").
		Where("activity_participations.user_id = ? AND activity_participations.status = ?", userID, "cancelled").
		Find(&activities).Error

	return activities, err
}
