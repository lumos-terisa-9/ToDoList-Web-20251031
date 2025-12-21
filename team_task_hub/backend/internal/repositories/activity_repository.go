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

// // FindByOrgAndParticipation 根据组织ID和参与类型查询活动
// func (r *ActivityRepository) FindByOrgAndParticipation(orgID uint, participationType string) ([]models.Activity, error) {
// 	var activities []models.Activity

// 	// 在数据访问层内部构建查询条件
// 	err := r.db.Model(&models.Activity{}).
// 		Where("organization_id = ?", orgID).
// 		Where("status = ?", "active").
// 		Where("end_time > ?", time.Now()).
// 		Where("participation_limit = ?", participationType).
// 		Find(&activities).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	return activities, nil
// }

// FindByOrgAndParticipation 根据组织ID和参与类型查询活动，并按用户参与状态排序（已参与的活动在前）
func (r *ActivityRepository) FindByOrgAndParticipation(orgID uint, userID uint, participationType string) ([]models.Activity, uint, error) {
	var activities []models.Activity
	var participatedCount int64

	// 1. 构建基础查询：此查询用于获取活动列表
	baseQuery := r.db.Model(&models.Activity{}).
		Select("activities.*, "+
			"COUNT(CASE WHEN activity_participations.user_id = ? THEN 1 ELSE NULL END) > 0 as user_has_participated", userID).
		Joins("LEFT JOIN activity_participations ON activities.id = activity_participations.activity_id").
		Where("activities.organization_id = ?", orgID).
		Where("activities.status = ?", "active").
		Where("activities.end_time > ?", time.Now()).
		Where("activities.participation_limit = ?", participationType).
		Group("activities.id").
		Order("user_has_participated DESC") //按照是否参与来排序

	// 获取用户参与的活动数量（基于基础查询创建新查询，添加额外过滤）
	countQuery := baseQuery.Session(&gorm.Session{}) // 创建基础查询的副本
	err := countQuery.
		Where("activity_participations.user_id = ?", userID).
		Count(&participatedCount).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询活动列表（使用未修改的基础查询，获取所有活动）
	err = baseQuery.Find(&activities).Error
	if err != nil {
		return nil, 0, err
	}

	return activities, uint(participatedCount), nil
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

// FindOrganizationsWithActivitiesInTimeRange 查找用户在时间范围内有活动的组织
func (r *ActivityRepository) FindOrganizationsWithActivitiesInTimeRange(userID uint, startTime, endTime time.Time) ([]models.Organization, error) {
	var organizations []models.Organization

	// 使用JOIN连接三张表：organizations -> activities -> activity_participations
	err := r.db.
		Select("DISTINCT organizations.*"). // 去重，避免同一组织重复出现
		Joins("JOIN activities ON activities.organization_id = organizations.id").
		Joins("JOIN activity_participations ON activity_participations.activity_id = activities.id").
		Where("activity_participations.user_id = ?", userID).
		Where("(activities.start_time BETWEEN ? AND ? OR activities.end_time BETWEEN ? AND ?)",
			startTime, endTime, startTime, endTime).
		Find(&organizations).Error

	if err != nil {
		return nil, fmt.Errorf("查询有活动的组织失败: %v", err)
	}

	return organizations, nil
}

// FindUserActivitiesStartingInRange 查找用户在时间范围内开始的活动
func (r *ActivityRepository) FindUserActivitiesStartingInRange(userID uint, startTime, endTime time.Time) ([]models.Activity, error) {
	var activities []models.Activity

	// 使用JOIN连接活动表和活动参与表，并预加载组织信息
	err := r.db.
		Joins("JOIN activity_participations ON activities.id = activity_participations.activity_id").
		Where("activity_participations.user_id = ?", userID).
		Where("activities.start_time BETWEEN ? AND ?", startTime, endTime).
		Where("activities.status != ?", "cancelled").
		Preload("Organization", func(db *gorm.DB) *gorm.DB { // 对预加载的Organization字段进行定制
			return db.Select("Name", "LogoURL") // 只选择组织的名称和头像URL字段
		}).
		Find(&activities).Error

	if err != nil {
		return nil, fmt.Errorf("查询开始时间范围内的活动失败: %v", err)
	}

	return activities, nil
}

// FindUserActivitiesEndingInRange 查找用户在时间范围内结束的活动
func (r *ActivityRepository) FindUserActivitiesEndingInRange(userID uint, startTime, endTime time.Time) ([]models.Activity, error) {
	var activities []models.Activity

	err := r.db.
		Joins("JOIN activity_participations ON activities.id = activity_participations.activity_id").
		Where("activity_participations.user_id = ?", userID).
		Where("activities.end_time BETWEEN ? AND ?", startTime, endTime).
		Where("activities.status != ?", "cancelled").
		Preload("Organization", func(db *gorm.DB) *gorm.DB { // 对预加载的Organization字段进行定制
			return db.Select("Name", "LogoURL")
		}).
		Find(&activities).Error

	if err != nil {
		return nil, fmt.Errorf("查询结束时间范围内的活动失败: %v", err)
	}

	return activities, nil
}

// FindUserActivitiesCompletedInRange 查找用户在时间范围内完成的活动
func (r *ActivityRepository) FindUserActivitiesCompletedInRange(userID uint, startTime, endTime time.Time) ([]models.Activity, error) {
	var activities []models.Activity

	err := r.db.
		Joins("JOIN activity_participations ON activities.id = activity_participations.activity_id").
		Where("activity_participations.user_id = ?", userID).
		Where("activity_participations.status = 'completed'").
		Where("activity_participations.completed_at BETWEEN ? AND ?", startTime, endTime).
		Preload("Organization", func(db *gorm.DB) *gorm.DB { // 对预加载的Organization字段进行定制
			return db.Select("Name", "LogoURL")
		}).
		Find(&activities).Error

	if err != nil {
		return nil, fmt.Errorf("查询已完成活动失败: %v", err)
	}

	return activities, nil
}

// FindUserActivitiesOverlappingRange 查找用户参与且与时间范围有交集的活动
func (r *ActivityRepository) FindUserActivitiesOverlappingRange(userID uint, startTime, endTime time.Time) ([]models.Activity, error) {
	var activities []models.Activity

	// 活动与时间范围有交集的查询条件：
	// 活动开始时间 <= 范围结束时间 AND 活动结束时间 >= 范围开始时间
	err := r.db.
		Joins("JOIN activity_participations ON activities.id = activity_participations.activity_id").
		Where("activity_participations.user_id = ?", userID).
		Where("activities.start_time <= ? AND activities.end_time >= ?", endTime, startTime).
		Where("activities.status != ?", "cancelled").
		Preload("Organization", func(db *gorm.DB) *gorm.DB { // 对预加载的Organization字段进行定制
			return db.Select("Name", "LogoURL")
		}).
		Find(&activities).Error

	if err != nil {
		return nil, fmt.Errorf("查询时间范围有交集的活动失败: %v", err)
	}

	return activities, nil
}

// 在 activityRepository 结构体中实现
func (r *ActivityRepository) FindUserPendingActivitiesEndingOnDate(userID uint, startTime, endTime time.Time) ([]models.Activity, error) {
	var activities []models.Activity

	err := r.db.
		Joins("JOIN activity_participations ON activities.id = activity_participations.activity_id").
		Where("activity_participations.user_id = ?", userID).
		Where("activity_participations.status = 'pending'").
		Where("activities.end_time BETWEEN ? AND ?", startTime, endTime).
		Preload("Organization", func(db *gorm.DB) *gorm.DB {
			return db.Select("Name", "LogoURL")
		}).
		Find(&activities).Error

	if err != nil {
		return nil, fmt.Errorf("查询指定日期过期待办失败: %v", err)
	}

	return activities, nil
}
