// services/organization_service.go
package services

import (
	"errors"
	"fmt"
	"team_task_hub/backend/internal/cache"
	"team_task_hub/backend/internal/models"
	"team_task_hub/backend/internal/repositories"
	"time"

	"gorm.io/gorm"
)

type OrganizationService struct {
	orgRepo       *repositories.OrganizationRepository
	orgMemberRepo *repositories.OrganizationMemberRepository
	orgAppRepo    *repositories.OrganizationApplicationRepository
}

func NewOrganizationService(
	orgRepo *repositories.OrganizationRepository,
	orgMemberRepo *repositories.OrganizationMemberRepository,
	orgAppRepo *repositories.OrganizationApplicationRepository,
) *OrganizationService {
	return &OrganizationService{
		orgRepo:       orgRepo,
		orgMemberRepo: orgMemberRepo,
		orgAppRepo:    orgAppRepo,
	}
}

// CreateOrganizationWithCache 创建组织并更新缓存
func (s *OrganizationService) CreateOrganizationWithCache(org *models.Organization, creatorUserID uint) error {
	if org.Name == "" {
		return errors.New("组织名称不能为空")
	}
	if creatorUserID == 0 {
		return errors.New("创建者ID无效")
	}

	// 检查组织名称是否已存在
	exists, err := s.orgRepo.ExistsByName(org.Name)
	if err != nil {
		return fmt.Errorf("检查组织名称是否存在时出错: %v", err)
	}
	if exists {
		return errors.New("组织名称已存在")
	}

	// 设置组织信息
	org.CreatorID = creatorUserID

	// 创建组织
	if err := s.orgRepo.Create(org); err != nil {
		return fmt.Errorf("创建组织失败: %v", err)
	}

	// 为创建者添加组织成员记录
	creatorMember := &models.OrganizationMember{
		OrganizationID: org.ID,
		UserID:         creatorUserID,
		Role:           "creator",
		Status:         "active",
	}

	if err := s.orgMemberRepo.Create(creatorMember); err != nil {
		return fmt.Errorf("创建组织成员记录失败: %v", err)
	}

	// 更新缓存
	go func() {
		cache.SetOrganization(org, time.Hour)
		cache.SetOrganizationIDByName(org.Name, org.ID, time.Hour)
		cache.DeleteUserOrganizationNames(creatorUserID)
	}()

	return nil
}

// SubmitJoinApplication 提交加入组织申请
func (s *OrganizationService) SubmitJoinApplication(application *models.OrganizationApplication) error {
	if application.ApplicantUserID == 0 {
		return errors.New("申请人ID无效")
	}
	if application.OrganizationName == "" {
		return errors.New("组织名称不能为空")
	}
	if application.ApplicationReason == "" {
		return errors.New("申请理由不能为空")
	}

	// 检查目标组织是否存在
	org, err := s.orgRepo.GetByName(application.OrganizationName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("目标组织不存在")
		}
		return fmt.Errorf("查询组织失败: %v", err)
	}

	// 检查是否已是组织成员
	exists, err := s.orgMemberRepo.Exists(org.ID, application.ApplicantUserID)
	if err != nil {
		return fmt.Errorf("检查成员关系时出错: %v", err)
	}
	if exists {
		return errors.New("您已是该组织成员")
	}

	// 设置申请类型和状态
	application.ApplicationType = "join_org"
	application.Status = "pending"

	// 创建申请记录
	if err := s.orgAppRepo.Create(application); err != nil {
		return fmt.Errorf("提交申请失败: %v", err)
	}

	return nil
}

// SubmitCreateOrganizationApplication 提交创建组织申请
func (s *OrganizationService) SubmitCreateOrganizationApplication(application *models.OrganizationApplication) error {
	if application.ApplicantUserID == 0 {
		return errors.New("申请人ID无效")
	}
	if application.OrganizationName == "" {
		return errors.New("组织名称不能为空")
	}
	if application.ApplicationReason == "" {
		return errors.New("申请理由不能为空")
	}

	// 检查组织名称是否已存在
	exists, err := s.orgRepo.ExistsByName(application.OrganizationName)
	if err != nil {
		return fmt.Errorf("检查组织名称是否存在时出错: %v", err)
	}
	if exists {
		return errors.New("组织名称已存在")
	}

	// 设置申请类型和状态
	application.ApplicationType = "create_org"
	application.Status = "pending"

	// 创建申请记录
	if err := s.orgAppRepo.Create(application); err != nil {
		return fmt.Errorf("提交创建组织申请失败: %v", err)
	}

	return nil
}

// ProcessApplication 处理组织申请（批准/拒绝）
func (s *OrganizationService) ProcessApplication(applicationID uint, action string, remark string) error {
	// 获取申请
	application, err := s.orgAppRepo.FindByID(applicationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("申请不存在")
		}
		return fmt.Errorf("查询申请失败: %v", err)
	}

	if application.Status != "pending" {
		return errors.New("该申请已被处理")
	}

	newStatus := "approved"
	if action == "reject" {
		newStatus = "rejected"
	}

	// 如果批准，执行相应操作
	if action == "approve" {
		if application.ApplicationType == "create_org" {
			// 批准创建组织
			newOrg := &models.Organization{
				Name:        application.OrganizationName,
				Description: application.ApplicantIntroduction,
				CreatorID:   application.ApplicantUserID,
			}

			if err := s.CreateOrganizationWithCache(newOrg, application.ApplicantUserID); err != nil {
				return fmt.Errorf("创建组织失败: %v", err)
			}

		} else if application.ApplicationType == "join_org" {
			// 批准加入组织
			org, err := s.orgRepo.GetByName(application.OrganizationName)
			if err != nil {
				return fmt.Errorf("查询组织失败: %v", err)
			}

			// 添加成员
			member := &models.OrganizationMember{
				OrganizationID: org.ID,
				UserID:         application.ApplicantUserID,
				Role:           "member",
				Status:         "active",
			}

			if err := s.orgMemberRepo.Create(member); err != nil {
				return fmt.Errorf("添加成员失败: %v", err)
			}

			// 更新用户组织缓存
			go cache.DeleteUserOrganizationNames(application.ApplicantUserID)
		}
	}

	// 更新申请状态
	if err := s.orgAppRepo.UpdateStatus(applicationID, newStatus, remark); err != nil {
		return fmt.Errorf("更新申请状态失败: %v", err)
	}

	return nil
}
