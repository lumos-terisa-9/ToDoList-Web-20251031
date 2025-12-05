package repositories

import (
	"team_task_hub/backend/internal/models"

	"gorm.io/gorm"
)

type OrganizationApplicationRepository struct {
	db *gorm.DB
}

func NewOrganizationApplicationRepository(db *gorm.DB) *OrganizationApplicationRepository {
	return &OrganizationApplicationRepository{db: db}
}

// Create 创建新的组织申请
func (r *OrganizationApplicationRepository) Create(application *models.OrganizationApplication) error {
	result := r.db.Create(application)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindByID 根据申请ID查找申请
func (r *OrganizationApplicationRepository) FindByID(id uint) (*models.OrganizationApplication, error) {
	var application models.OrganizationApplication
	result := r.db.First(&application, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &application, nil
}

// FindByApplicantAndStatus 根据申请用户ID查找申请列表
func (r *OrganizationApplicationRepository) FindByApplicant(applicantUserID uint) ([]models.OrganizationApplication, error) {
	var applications []models.OrganizationApplication

	result := r.db.Where("applicant_user_id=?", applicantUserID).Find(&applications)

	if result.Error != nil {
		return nil, result.Error
	}
	return applications, nil
}

// FindByOrganizationAndStatus 根据组织名称和状态查找申请列表
func (r *OrganizationApplicationRepository) FindByOrganizationAndStatus(organizationName string, status string) ([]models.OrganizationApplication, error) {
	var applications []models.OrganizationApplication
	query := r.db.Where("organization_name = ?", organizationName)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	result := query.Order("id DESC").Find(&applications)
	if result.Error != nil {
		return nil, result.Error
	}
	return applications, nil
}

// FindPendingApplications 查找所有待处理的申请（管理员用）
func (r *OrganizationApplicationRepository) FindPendingApplications() ([]models.OrganizationApplication, error) {
	var applications []models.OrganizationApplication
	result := r.db.Where("status = ?", "pending").Order("id DESC").Find(&applications)
	if result.Error != nil {
		return nil, result.Error
	}
	return applications, nil
}

// FindByTypeAndStatus 根据申请类型和状态查找
func (r *OrganizationApplicationRepository) FindByTypeAndStatus(appType, status string) ([]models.OrganizationApplication, error) {
	var applications []models.OrganizationApplication

	result := r.db.Where("application_type = ? and status = ?", appType, status).Find(&applications)
	if result.Error != nil {
		return nil, result.Error
	}

	return applications, nil
}

// FindByOrgNameAndTypeAndStatus 根据组织名字获取加入组织申请
func (r *OrganizationApplicationRepository) FindByOrgNameAndTypeAndStatus(orgName, appType, status string) ([]models.OrganizationApplication, error) {
	var applications []models.OrganizationApplication
	result := r.db.Where("organization_name=? and application_type=? and status=?", orgName, appType, status).Find(&applications)
	if result.Error != nil {
		return nil, result.Error
	}
	return applications, nil
}

// Update 更新申请信息
func (r *OrganizationApplicationRepository) Update(application *models.OrganizationApplication) error {
	result := r.db.Save(application)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateStatus 更新申请状态（审批操作）
func (r *OrganizationApplicationRepository) UpdateStatus(id uint, status string, adminRemark string) error {
	result := r.db.Model(&models.OrganizationApplication{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":       status,
			"admin_remark": adminRemark,
		})
	return result.Error
}

// Delete 删除申请
func (r *OrganizationApplicationRepository) Delete(id uint) error {
	result := r.db.Delete(&models.OrganizationApplication{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// ExistsPendingApplication 检查用户是否已有针对某组织的待处理申请
func (r *OrganizationApplicationRepository) ExistsPendingApplication(applicantUserID uint, organizationName string, appType string) (bool, error) {
	var count int64
	result := r.db.Model(&models.OrganizationApplication{}).
		Where("applicant_user_id = ? AND organization_name = ? AND application_type = ? AND status = ?",
			applicantUserID, organizationName, appType, "pending").
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// GetApplicationsWithUsers 获取申请列表并预加载申请人信息
func (r *OrganizationApplicationRepository) GetApplicationsWithUsers(status string) ([]models.OrganizationApplication, error) {
	var applications []models.OrganizationApplication

	query := r.db.Preload("Applicant")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	result := query.Order("id DESC").Find(&applications)
	if result.Error != nil {
		return nil, result.Error
	}
	return applications, nil
}

// GetApplicationStats 获取申请统计信息
func (r *OrganizationApplicationRepository) GetApplicationStats() (map[string]int64, error) {
	stats := make(map[string]int64)

	var total, pending, approved, rejected int64

	// 总数
	r.db.Model(&models.OrganizationApplication{}).Count(&total)

	// 各状态数量
	r.db.Model(&models.OrganizationApplication{}).Where("status = ?", "pending").Count(&pending)
	r.db.Model(&models.OrganizationApplication{}).Where("status = ?", "approved").Count(&approved)
	r.db.Model(&models.OrganizationApplication{}).Where("status = ?", "rejected").Count(&rejected)

	stats["total"] = total
	stats["pending"] = pending
	stats["approved"] = approved
	stats["rejected"] = rejected

	return stats, nil
}
