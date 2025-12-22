package repositories

import (
	"fmt"
	"team_task_hub/backend/internal/models"

	"gorm.io/gorm"
)

// OrganizationRepository 组织数据访问层
// 负责所有与organizations表相关的数据库操作
type OrganizationRepository struct {
	db *gorm.DB
}

// NewOrganizationRepository 构造函数：创建OrganizationRepository实例
func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

// Create 创建新组织
func (r *OrganizationRepository) Create(organization *models.Organization) error {
	result := r.db.Create(organization)
	if result.Error != nil {
		return fmt.Errorf("创建组织失败: %v", result.Error)
	}
	return nil
}

// SearchOrganizationsByName 根据组织名称模糊搜索，仅返回ID和名称
func (r *OrganizationRepository) SearchOrganizationsByName(keyword string, limit int) ([]models.OrgInfo, error) {
	var simpleOrgs []models.OrgInfo

	if keyword == "" {
		return simpleOrgs, nil
	}

	// 使用左前缀匹配以利用数据库索引
	searchPattern := keyword + "%"

	result := r.db.Model(&models.Organization{}).
		Select("id", "name").
		Where("name LIKE ?", searchPattern).
		Limit(limit).
		Find(&simpleOrgs)

	if result.Error != nil {
		return nil, fmt.Errorf("模糊搜索组织失败: %v", result.Error)
	}
	return simpleOrgs, nil
}

// FindByGeohashPrefix 通过Geohash前缀查找组织
func (r *OrganizationRepository) FindByGeohashPrefix(prefix string) ([]models.OrgInfo, error) {
	var simpleOrgs []models.OrgInfo

	// 使用LIKE查询进行前缀匹配，利用数据库索引提高性能
	err := r.db.Model(&models.Organization{}).
		Select("id", "name").
		Where("location_code LIKE ?", prefix+"%").
		Find(&simpleOrgs).Error

	if err != nil {
		return nil, fmt.Errorf("通过位置码前缀查询组织失败，原因：%v", err)
	}

	return simpleOrgs, nil
}

// GetByID 根据主键ID查找组织
func (r *OrganizationRepository) GetByID(id uint) (*models.Organization, error) {
	var organization models.Organization
	result := r.db.Preload("Creator").First(&organization, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("组织不存在")
		}
		return nil, fmt.Errorf("查询组织失败: %v", result.Error)
	}
	organization.Creator.PasswordHash = ""
	return &organization, nil
}

// GetByCreatorID 根据创建者ID查找组织列表
func (r *OrganizationRepository) GetByCreatorID(creatorID uint) ([]models.Organization, error) {
	var organizations []models.Organization
	result := r.db.Where("creator_id = ?", creatorID).Preload("Creator").Find(&organizations)
	if result.Error != nil {
		return nil, fmt.Errorf("根据创建者查询组织失败: %v", result.Error)
	}
	return organizations, nil
}

// GetByLocationCode 根据位置代码查找组织
func (r *OrganizationRepository) GetByLocationCode(locationCode string) (*models.Organization, error) {
	var organization models.Organization
	result := r.db.Where("location_code = ?", locationCode).Preload("Creator").First(&organization)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("根据位置代码查询组织失败: %v", result.Error)
	}
	organization.Creator.PasswordHash = ""
	return &organization, nil
}

// GetByName 根据名称查找组织
func (r *OrganizationRepository) GetByName(name string) (*models.Organization, error) {
	var organization models.Organization
	result := r.db.Where("name = ?", name).Preload("Creator").First(&organization)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("根据名称查询组织失败: %v", result.Error)
	}
	organization.Creator.PasswordHash = ""
	return &organization, nil
}

// Update 更新组织信息
func (r *OrganizationRepository) Update(organization *models.Organization) error {
	result := r.db.Save(organization)
	if result.Error != nil {
		return fmt.Errorf("更新组织失败: %v", result.Error)
	}
	return nil
}

// Delete 删除组织
func (r *OrganizationRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Organization{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除组织失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("组织不存在")
	}
	return nil
}

// GetAll 获取所有组织列表（分页）
func (r *OrganizationRepository) GetAll(page, pageSize int) ([]models.Organization, error) {
	var organizations []models.Organization
	offset := (page - 1) * pageSize
	result := r.db.Preload("Creator").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&organizations)
	if result.Error != nil {
		return nil, fmt.Errorf("查询组织列表失败: %v", result.Error)
	}
	return organizations, nil
}

// GetByStatus 根据状态获取组织列表
func (r *OrganizationRepository) GetByStatus(status string) ([]models.Organization, error) {
	var organizations []models.Organization
	result := r.db.Where("status = ?", status).Preload("Creator").Find(&organizations)
	if result.Error != nil {
		return nil, fmt.Errorf("根据状态查询组织失败: %v", result.Error)
	}
	return organizations, nil
}

// ExistsByName 检查组织名称是否已存在
func (r *OrganizationRepository) ExistsByName(name string) (bool, error) {
	var count int64
	result := r.db.Model(&models.Organization{}).Where("name = ?", name).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("检查组织名称存在性失败: %v", result.Error)
	}
	return count > 0, nil
}

// ExistsByLocationCode 检查位置代码是否已存在
func (r *OrganizationRepository) ExistsByLocationCode(locationCode string) (bool, error) {
	var count int64
	result := r.db.Model(&models.Organization{}).Where("location_code = ?", locationCode).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("检查位置代码存在性失败: %v", result.Error)
	}
	return count > 0, nil
}

// Count 获取组织总数
func (r *OrganizationRepository) Count() (int64, error) {
	var count int64
	result := r.db.Model(&models.Organization{}).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("统计组织数量失败: %v", result.Error)
	}
	return count, nil
}

// CountByStatus 根据状态统计组织数量
func (r *OrganizationRepository) CountByStatus(status string) (int64, error) {
	var count int64
	result := r.db.Model(&models.Organization{}).Where("status = ?", status).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("根据状态统计组织数量失败: %v", result.Error)
	}
	return count, nil
}

// UpdateStatus 更新组织状态
func (r *OrganizationRepository) UpdateStatus(id uint, status string) error {
	result := r.db.Model(&models.Organization{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("更新组织状态失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("组织不存在")
	}
	return nil
}

// FindAllForRecommendation 获取所有组织的基本信息
func (r *OrganizationRepository) FindAllForRecommendation() ([]models.Organization, error) {
	var orgs []models.Organization
	// 只选择推荐计算所需的字段，避免不必要的数据传输
	err := r.db.Select("id", "name", "description").Find(&orgs).Error
	if err != nil {
		return nil, err
	}
	return orgs, nil
}
