package repositories

import (
	"fmt"
	"team_task_hub/backend/internal/models"

	"gorm.io/gorm"
)

// UserRepository 用户数据访问层
// 负责所有与users表相关的数据库操作
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 构造函数：创建UserRepository实例
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建新用户
func (r *UserRepository) Create(user *models.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		return fmt.Errorf("创建用户失败: %v", result.Error)
	}
	return nil
}

// FindByID 根据主键ID查找用户
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %v", result.Error)
	}
	return &user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 邮箱不存在，不是错误
		}
		return nil, fmt.Errorf("查询用户失败: %v", result.Error)
	}
	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询用户失败: %v", result.Error)
	}
	return &user, nil
}

// Update 更新用户信息
func (r *UserRepository) Update(user *models.User) error {
	result := r.db.Save(user)
	if result.Error != nil {
		return fmt.Errorf("更新用户失败: %v", result.Error)
	}
	return nil
}

// Delete 删除用户
func (r *UserRepository) Delete(id uint) error {
	result := r.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除用户失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在")
	}
	return nil
}

// FindAll 获取所有用户列表
func (r *UserRepository) FindAll(page, pageSize int) ([]models.User, error) {
	var users []models.User
	offset := (page - 1) * pageSize
	result := r.db.Offset(offset).Limit(pageSize).Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("查询用户列表失败: %v", result.Error)
	}
	return users, nil
}

// ExistsById 检查id（学号）是否存在
func (r *UserRepository) ExistsById(id uint) (bool, error) {
	var count int64
	result := r.db.Model(&models.User{}).Where("id = ?", id).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("检查ID存在性失败: %v", result.Error)
	}
	return count > 0, nil
}

// ExistsByEmail 检查邮箱是否已存在
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	result := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("检查邮箱存在性失败: %v", result.Error)
	}
	return count > 0, nil
}

// ExistsByUsername 检查用户名是否已存在
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	result := r.db.Model(&models.User{}).Where("username = ?", username).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("检查用户名存在性失败: %v", result.Error)
	}
	return count > 0, nil
}

// UpdatePassword 更新用户密码
func (r *UserRepository) UpdatePassword(userID uint, newPasswordHash string) error {
	result := r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("password_hash", newPasswordHash)

	if result.Error != nil {
		return fmt.Errorf("更新密码失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在")
	}
	return nil
}

// UpdateAvatar 更新用户头像
// userID: 用户ID
// avatarURL: 新的头像URL
func (r *UserRepository) UpdateAvatar(userID uint, avatarURL string) error {
	result := r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("avatar_url", avatarURL)

	if result.Error != nil {
		return fmt.Errorf("更新头像失败: %v", result.Error)
	}
	return nil
}

// Count 获取用户总数
func (r *UserRepository) Count() (int64, error) {
	var count int64
	result := r.db.Model(&models.User{}).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("统计用户数量失败: %v", result.Error)
	}
	return count, nil
}

// FindByIDs 根据多个ID批量查找用户
func (r *UserRepository) FindByIDs(ids []uint) ([]models.User, error) {
	var users []models.User
	result := r.db.Where("id IN ?", ids).Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("批量查询用户失败: %v", result.Error)
	}
	return users, nil
}

// SoftDelete 软删除用户（如果模型有DeletedAt字段）
func (r *UserRepository) SoftDelete(id uint) error {
	result := r.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("软删除用户失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在")
	}
	return nil
}

// Restore 恢复软删除的用户
func (r *UserRepository) Restore(id uint) error {
	result := r.db.Unscoped().Model(&models.User{}).Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return fmt.Errorf("恢复用户失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在")
	}
	return nil
}
