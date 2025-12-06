// services/organization_service.go
package services

import (
	"errors"
	"fmt"
	"log"
	"strings"
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

// CreateOrganization 直接创建组织
func (s *OrganizationService) CreateOrganization(org *models.Organization, creatorUserID uint) error {
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

	// 为创建者添加组织成员记录，角色为 'creator'
	creatorMember := &models.OrganizationMember{
		OrganizationID: org.ID,
		UserID:         creatorUserID,
		Role:           "creator",
		Status:         "active",
		JoinedAt:       time.Now(),
	}

	if err := s.orgMemberRepo.Create(creatorMember); err != nil {
		return fmt.Errorf("创建组织成员记录失败: %v", err)
	}
	//删除缓存
	cache.DeleteUserOrganizationOverviews(creatorUserID)
	return nil
}

// JoinOrganization 直接加入组织
func (s *OrganizationService) JoinOrganization(orgID, userID uint) error {
	// 检查用户是否已是组织成员
	exists, err := s.orgMemberRepo.Exists(orgID, userID)
	if err != nil {
		return fmt.Errorf("检查成员关系时出错: %v", err)
	}
	if exists {
		return errors.New("用户已是该组织成员")
	}

	// 添加成员
	member := &models.OrganizationMember{
		OrganizationID: orgID,
		UserID:         userID,
	}

	if err := s.orgMemberRepo.Create(member); err != nil {
		return fmt.Errorf("添加成员失败: %v", err)
	}

	cache.DeleteUserOrganizationOverviews(userID)

	return nil
}

// SubmitCreateOrganizationApplication 提交创建组织申请
func (s *OrganizationService) SubmitCreateOrganizationApplication(application *models.OrganizationApplication) error {
	// 检查组织名称是否已存在
	exists, err := s.orgRepo.ExistsByName(application.OrganizationName)
	if err != nil {
		return fmt.Errorf("检查组织名称是否存在时出错: %v", err)
	}
	if exists {
		return errors.New("组织名称已存在")
	}

	// 设置申请类型
	application.ApplicationType = "create_org"

	// 创建申请记录
	if err := s.orgAppRepo.Create(application); err != nil {
		return fmt.Errorf("提交创建组织申请失败: %v", err)
	}

	return nil
}

// SubmitJoinApplication 提交加入组织申请
func (s *OrganizationService) SubmitJoinApplication(application *models.OrganizationApplication) error {
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

	// 设置申请信息
	application.ApplicationType = "join_org"

	// 创建申请记录
	if err := s.orgAppRepo.Create(application); err != nil {
		return fmt.Errorf("提交加入组织申请失败: %v", err)
	}

	return nil
}

// ProcessApplication 处理组织申请
func (s *OrganizationService) ProcessApplication(applicationID uint, action string, remark string, processedBy uint) error {
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

			if err := s.CreateOrganization(newOrg, application.ApplicantUserID); err != nil {
				return fmt.Errorf("创建组织失败: %v", err)
			}

		} else if application.ApplicationType == "join_org" {
			// 批准加入组织
			org, err := s.orgRepo.GetByName(application.OrganizationName)
			if err != nil {
				return fmt.Errorf("查询组织失败: %v", err)
			}
			s.JoinOrganization(org.ID, application.ApplicantUserID)
		}
	}

	// 更新申请状态和备注
	if err := s.orgAppRepo.UpdateStatus(applicationID, newStatus, remark); err != nil {
		return fmt.Errorf("更新申请状态失败: %v", err)
	}

	return nil
}

// GetUserOrganizationOverviews 获取用户的组织总览
func (s *OrganizationService) GetUserOrganizationOverviews(userID uint) ([]cache.UserOrganizationOverview, error) {
	// 先尝试从缓存获取
	cachedOverviews, err := cache.GetUserOrganizationOverviews(userID)
	if err == nil && cachedOverviews != nil {
		return cachedOverviews, nil
	}

	// 缓存未命中，从数据库查询
	overviews, err := s.getUserOrganizationOverviewsFromDB(userID)
	if err != nil {
		return nil, err
	}

	// 异步更新缓存
	go cache.SetUserOrganizationOverviews(userID, overviews, 30*time.Minute)

	return overviews, nil
}

// getUserOrganizationOverviewsFromDB 从数据库获取用户组织总览
func (s *OrganizationService) getUserOrganizationOverviewsFromDB(userID uint) ([]cache.UserOrganizationOverview, error) {
	// 获取用户的所有组织成员关系
	members, err := s.orgMemberRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户组织关系失败: %v", err)
	}

	var overviews []cache.UserOrganizationOverview

	for _, member := range members {
		// 使用预加载的组织信息构建总览
		overview := cache.UserOrganizationOverview{
			OrganizationName: member.Organization.Name,
			CreatorID:        member.Organization.CreatorID,
			LogoURL:          member.Organization.LogoURL,
			JoinedAt:         member.JoinedAt,
		}

		overviews = append(overviews, overview)
	}

	return overviews, nil
}

// invalidateUserOrganizationCache 使用户组织缓存失效
func (s *OrganizationService) invalidateUserOrganizationCache(userID uint) {
	if err := cache.DeleteUserOrganizationOverviews(userID); err != nil {
		log.Printf("删除用户组织缓存失败 userID=%d: %v", userID, err)
	}
}

// GetPendingApplications 获取待处理申请列表
func (s *OrganizationService) GetPendingApplications() ([]models.OrganizationApplication, error) {
	applications, err := s.orgAppRepo.FindPendingApplications()
	if err != nil {
		return nil, fmt.Errorf("查询待处理申请失败: %v", err)
	}
	return applications, nil
}

// GetPendingCreateOrgApplications 管理员获取创建组织申请列表
func (s *OrganizationService) GetPendingCreateOrgApplications() ([]models.OrganizationApplication, error) {
	applications, err := s.orgAppRepo.FindByTypeAndStatus("create_org", "pending")
	if err != nil {
		return nil, fmt.Errorf("管理员获取创建组织申请失败，失败原因:%v", err)
	}
	return applications, nil
}

// OrgGetPendingJoinApplications 组织获取加入组织申请列表
func (s *OrganizationService) OrgGetPendingJoinApplications(orgName string) ([]models.OrganizationApplication, error) {
	applications, err := s.orgAppRepo.FindByOrgNameAndTypeAndStatus(orgName, "join_org", "pending")
	if err != nil {
		return nil, fmt.Errorf("组织获取加入组织申请失败，失败原因:%v", err)
	}
	return applications, nil
}

// GetUserApplications 获取用户的申请列表
func (s *OrganizationService) GetUserApplications(userID uint) ([]models.OrganizationApplication, error) {
	applications, err := s.orgAppRepo.FindByApplicant(userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户申请失败: %v", err)
	}
	return applications, nil
}

// DeleteUserApplication 删除用户申请
func (s *OrganizationService) DeleteUserApplication(appID uint) error {
	err := s.orgAppRepo.Delete(appID)
	if err != nil {
		return fmt.Errorf("删除申请失败，原因:%v", err)
	}
	return nil
}

// GetApplicationDetail 获取申请详情
func (s *OrganizationService) GetApplicationDetail(applicationID uint) (*models.OrganizationApplication, error) {
	application, err := s.orgAppRepo.FindByID(applicationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("申请不存在")
		}
		return nil, fmt.Errorf("查询申请详情失败: %v", err)
	}
	return application, nil
}

// RemoveOrganizationMember 移除组织成员
func (s *OrganizationService) RemoveOrganizationMember(orgID, userID uint, removedBy uint) error {
	// 检查成员关系是否存在
	member, err := s.orgMemberRepo.GetMember(orgID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("成员关系不存在")
		}
		return fmt.Errorf("查询成员关系失败: %v", err)
	}

	// 不能移除创建者
	if member.Role == "creator" {
		return errors.New("不能移除组织创建者")
	}

	if err := s.orgMemberRepo.RemoveMember(orgID, userID); err != nil {
		return fmt.Errorf("删除成员失败，原因：%v", err)
	}

	s.invalidateUserOrganizationCache(userID)

	return nil
}

// IsOrganizationCreator 验证用户是否是组织创建者
func (s *OrganizationService) IsOrganizationCreator(orgID, userID uint) (bool, error) {
	org, err := s.orgRepo.GetByID(orgID)
	if err != nil {
		return false, fmt.Errorf("查询组织失败: %v", err)
	}
	return org.CreatorID == userID, nil
}

// IsOrganizationAdmin 验证用户是否是组织管理员或创建者
func (s *OrganizationService) IsOrganizationAdmin(orgID, userID uint) (bool, error) {
	// 检查是否是创建者
	isCreator, err := s.IsOrganizationCreator(orgID, userID)
	if err != nil {
		return false, err
	}
	if isCreator {
		return true, nil
	}

	// 检查是否是管理员
	member, err := s.orgMemberRepo.GetMember(orgID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("查询成员关系失败: %v", err)
	}

	return member.Role == "admin", nil
}

// IsOrganizationMember 验证用户是否是组织成员
func (s *OrganizationService) IsOrganizationMember(orgID, userID uint) (bool, error) {
	exists, err := s.orgMemberRepo.Exists(orgID, userID)
	if err != nil {
		return false, fmt.Errorf("检查成员关系失败: %v", err)
	}
	return exists, nil
}

// FindOrgInfoByName 根据关键词模糊化查询组织
func (s *OrganizationService) FindOrgInfoByName(name string) ([]models.OrgInfo, error) {
	cleanName := strings.TrimSpace(name)
	if cleanName == "" {
		return nil, errors.New("查询字段不能为空")
	}

	orgInfos, err := s.orgRepo.SearchOrganizationsByName(cleanName, 20)
	if err != nil {
		return nil, fmt.Errorf("查询组织info失败，原因：%v", err)
	}

	return orgInfos, nil
}

// FindOrgByID 根据id查询组织基本信息（带缓存）
func (s *OrganizationService) FindOrgByID(id uint) (*models.Organization, error) {
	if cacheOrg, err := cache.GetOrganizationInfo(id); cacheOrg != nil && err == nil {
		return cacheOrg, nil
	}
	org, err := s.orgRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("通过id查询组织出错，原因：%v", err)
	}

	go cache.SetOrganizationInfo(org, 30*time.Minute)
	return org, nil
}
