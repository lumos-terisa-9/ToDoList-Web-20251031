package repositories

import (
	"fmt"
	"team_task_hub/backend/internal/models"

	"gorm.io/gorm"
)

// OrganizationMemberRepository 组织成员数据访问层
type OrganizationMemberRepository struct {
	db *gorm.DB
}

// NewOrganizationMemberRepository 构造函数：创建OrganizationMemberRepository实例
func NewOrganizationMemberRepository(db *gorm.DB) *OrganizationMemberRepository {
	return &OrganizationMemberRepository{db: db}
}

// Create 添加组织成员
func (r *OrganizationMemberRepository) Create(member *models.OrganizationMember) error {
	result := r.db.Create(member)
	if result.Error != nil {
		return fmt.Errorf("添加组织成员失败: %v", result.Error)
	}
	return nil
}

// GetByID 根据主键ID查找组织成员
func (r *OrganizationMemberRepository) GetByID(id uint) (*models.OrganizationMember, error) {
	var member models.OrganizationMember
	result := r.db.Preload("User").Preload("Organization").First(&member, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("组织成员不存在")
		}
		return nil, fmt.Errorf("查询组织成员失败: %v", result.Error)
	}
	return &member, nil
}

// GetByOrganizationID 获取组织的所有成员
func (r *OrganizationMemberRepository) GetByOrganizationID(organizationID uint) ([]models.OrganizationMember, error) {
	var members []models.OrganizationMember
	result := r.db.Where("organization_id = ?", organizationID).
		Preload("User").
		Preload("Organization").
		Find(&members)
	if result.Error != nil {
		return nil, fmt.Errorf("查询组织成员失败: %v", result.Error)
	}
	return members, nil
}

// GetByUserID 获取用户加入的所有组织
func (r *OrganizationMemberRepository) GetByUserID(userID uint) ([]models.OrganizationMember, error) {
	var members []models.OrganizationMember
	result := r.db.Where("user_id = ?", userID).
		Preload("Organization").
		Find(&members)
	if result.Error != nil {
		return nil, fmt.Errorf("查询用户组织关系失败: %v", result.Error)
	}
	return members, nil
}

// GetMember 获取特定成员关系
func (r *OrganizationMemberRepository) GetMember(organizationID, userID uint) (*models.OrganizationMember, error) {
	var member models.OrganizationMember
	result := r.db.Where("organization_id = ? AND user_id = ?", organizationID, userID).
		First(&member)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询成员关系失败: %v", result.Error)
	}
	return &member, nil
}

// Update 更新成员信息
func (r *OrganizationMemberRepository) Update(member *models.OrganizationMember) error {
	result := r.db.Save(member)
	if result.Error != nil {
		return fmt.Errorf("更新成员信息失败: %v", result.Error)
	}
	return nil
}

// Delete 删除成员
func (r *OrganizationMemberRepository) Delete(id uint) error {
	result := r.db.Delete(&models.OrganizationMember{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除成员失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("成员记录不存在")
	}
	return nil
}

// RemoveMember 从组织中移除特定成员
func (r *OrganizationMemberRepository) RemoveMember(organizationID, userID uint) error {
	result := r.db.Where("organization_id = ? AND user_id = ?", organizationID, userID).Delete(&models.OrganizationMember{})
	if result.Error != nil {
		return fmt.Errorf("移除成员失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("成员记录不存在")
	}
	return nil
}

// GetMembersByRole 根据角色获取组织成员
func (r *OrganizationMemberRepository) GetMembersByRole(organizationID uint, role string) ([]models.OrganizationMember, error) {
	var members []models.OrganizationMember
	result := r.db.Where("organization_id = ? AND role = ?", organizationID, role).
		Preload("User").
		Find(&members)
	if result.Error != nil {
		return nil, fmt.Errorf("根据角色查询成员失败: %v", result.Error)
	}
	return members, nil
}

// GetActiveMembers 获取活跃成员
func (r *OrganizationMemberRepository) GetActiveMembers(organizationID uint) ([]models.OrganizationMember, error) {
	var members []models.OrganizationMember
	result := r.db.Where("organization_id = ? AND status = 'active'", organizationID).
		Preload("User").
		Find(&members)
	if result.Error != nil {
		return nil, fmt.Errorf("查询活跃成员失败: %v", result.Error)
	}
	return members, nil
}

// Exists 检查成员关系是否存在
func (r *OrganizationMemberRepository) Exists(organizationID, userID uint) (bool, error) {
	var count int64
	result := r.db.Model(&models.OrganizationMember{}).
		Where("organization_id = ? AND user_id = ?", organizationID, userID).
		Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("检查成员关系存在性失败: %v", result.Error)
	}
	return count > 0, nil
}

// IsUserAdmin 检查用户是否是组织管理员
func (r *OrganizationMemberRepository) IsUserAdmin(organizationID, userID uint) (bool, error) {
	var count int64
	result := r.db.Model(&models.OrganizationMember{}).
		Where("organization_id = ? AND user_id = ? AND role IN ('creator', 'admin')", organizationID, userID).
		Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("检查用户管理员权限失败: %v", result.Error)
	}
	return count > 0, nil
}

// UpdateRole 更新成员角色
func (r *OrganizationMemberRepository) UpdateRole(organizationID, userID uint, role string) error {
	result := r.db.Model(&models.OrganizationMember{}).
		Where("organization_id = ? AND user_id = ?", organizationID, userID).
		Update("role", role)

	if result.Error != nil {
		return fmt.Errorf("更新成员角色失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("成员记录不存在")
	}
	return nil
}

// FindUserInfosByOrgAndNamePrefix 根据组织ID和用户名前缀查询用户基本信息
func (r *OrganizationMemberRepository) FindUserInfosByOrgAndNamePrefix(orgID uint, namePrefix string) ([]models.UserInfo, error) {
	var userInfos []models.UserInfo

	// 构建基础查询：关联组织成员表和用户表，并筛选指定组织的活跃成员
	db := r.db.Model(&models.OrganizationMember{}).
		Joins("JOIN users ON organization_members.user_id = users.id").
		Where("organization_members.organization_id = ?", orgID)

	// 如果 namePrefix 不为空，则添加模糊查询条件；如果为空，则此条件不生效，查询所有成员
	if namePrefix != "" {
		db = db.Where("users.username LIKE ?", namePrefix+"%")
	}

	// 执行查询
	err := db.Select("users.id, users.username, users.avatar_url, users.email").
		Order("users.username ASC"). // 按用户名排序
		Scan(&userInfos).Error

	if err != nil {
		return nil, err
	}

	return userInfos, nil
}

// UpdateStatus 更新成员状态
func (r *OrganizationMemberRepository) UpdateStatus(organizationID, userID uint, status string) error {
	result := r.db.Model(&models.OrganizationMember{}).
		Where("organization_id = ? AND user_id = ?", organizationID, userID).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("更新成员状态失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("成员记录不存在")
	}
	return nil
}

// UpdateScore 更新成员积分
func (r *OrganizationMemberRepository) UpdateScore(organizationID, userID uint, score int) error {
	result := r.db.Model(&models.OrganizationMember{}).
		Where("organization_id = ? AND user_id = ?", organizationID, userID).
		Update("score", score)

	if result.Error != nil {
		return fmt.Errorf("更新成员积分失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("成员记录不存在")
	}
	return nil
}

// CountMembers 统计组织成员数量
func (r *OrganizationMemberRepository) CountMembers(organizationID uint) (int64, error) {
	var count int64
	result := r.db.Model(&models.OrganizationMember{}).
		Where("organization_id = ?", organizationID).
		Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("统计成员数量失败: %v", result.Error)
	}
	return count, nil
}

// CountMembersByRole 根据角色统计成员数量
func (r *OrganizationMemberRepository) CountMembersByRole(organizationID uint, role string) (int64, error) {
	var count int64
	result := r.db.Model(&models.OrganizationMember{}).
		Where("organization_id = ? AND role = ?", organizationID, role).
		Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("根据角色统计成员数量失败: %v", result.Error)
	}
	return count, nil
}

// CountMembersByStatus 根据状态统计成员数量
func (r *OrganizationMemberRepository) CountMembersByStatus(organizationID uint, status string) (int64, error) {
	var count int64
	result := r.db.Model(&models.OrganizationMember{}).
		Where("organization_id = ? AND status = ?", organizationID, status).
		Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("根据状态统计成员数量失败: %v", result.Error)
	}
	return count, nil
}
