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

	"github.com/mmcloughlin/geohash"
	"gorm.io/gorm"
)

type OrganizationService struct {
	orgRepo       *repositories.OrganizationRepository
	orgMemberRepo *repositories.OrganizationMemberRepository
	orgAppRepo    *repositories.OrganizationApplicationRepository
	codeRepo      *repositories.VerificationCodeRepository
}

func NewOrganizationService(
	orgRepo *repositories.OrganizationRepository,
	orgMemberRepo *repositories.OrganizationMemberRepository,
	orgAppRepo *repositories.OrganizationApplicationRepository,
	codeRepo *repositories.VerificationCodeRepository,
) *OrganizationService {
	return &OrganizationService{
		orgRepo:       orgRepo,
		orgMemberRepo: orgMemberRepo,
		orgAppRepo:    orgAppRepo,
		codeRepo:      codeRepo,
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
		JoinedAt:       time.Now(),
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

// SubmitChangeOrganizationApplication 提交修改组织申请
func (s *OrganizationService) SubmitChangeOrganizationApplication(application *models.OrganizationApplication) error {
	// 设置申请类型
	application.ApplicationType = "change_org"

	// 创建申请记录
	if err := s.orgAppRepo.Create(application); err != nil {
		return fmt.Errorf("提交修改组织申请失败: %v", err)
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
func (s *OrganizationService) GetPendingApplications(appType string) ([]models.OrganizationApplication, error) {
	applications, err := s.orgAppRepo.FindByTypeAndStatus(appType, "pending")
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

// GetPendingChangeOrgApplications 管理员获取修改组织申请列表
func (s *OrganizationService) GetPendingChangeOrgApplications() ([]models.OrganizationApplication, error) {
	applications, err := s.orgAppRepo.FindByTypeAndStatus("change_org", "pending")
	if err != nil {
		return nil, fmt.Errorf("管理员获取修改组织申请失败，失败原因:%v", err)
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

// UpdateOrganizationName 更新组织名称
func (s *OrganizationService) UpdateOrganizationName(orgID uint, newName string) error {
	// 获取当前组织信息
	org, err := s.orgRepo.GetByID(orgID)
	if err != nil {
		return fmt.Errorf("查询组织失败: %v", err)
	}

	// 更新字段
	org.Name = newName

	// 持久化到数据库
	err = s.orgRepo.Update(org)
	if err != nil {
		return fmt.Errorf("更新组织名称失败: %v", err)
	}

	// 清理缓存，保证后续读取的数据是最新的
	go cache.DeleteOrganizationInfo(orgID)

	return nil
}

// UpdateOrganizationDescription 更新组织描述
func (s *OrganizationService) UpdateOrganizationDescription(orgID uint, newDescription string) error {
	org, err := s.orgRepo.GetByID(orgID)
	if err != nil {
		return fmt.Errorf("查询组织失败: %v", err)
	}

	org.Description = newDescription

	err = s.orgRepo.Update(org)
	if err != nil {
		return fmt.Errorf("更新组织描述失败: %v", err)
	}

	go cache.DeleteOrganizationInfo(orgID)
	return nil
}

// UpdateOrganizationLogo 更新组织Logo
func (s *OrganizationService) UpdateOrganizationLogo(orgID uint, newLogoURL string) error {
	org, err := s.orgRepo.GetByID(orgID)
	if err != nil {
		return fmt.Errorf("查询组织失败: %v", err)
	}

	org.LogoURL = newLogoURL

	err = s.orgRepo.Update(org)
	if err != nil {
		return fmt.Errorf("更新组织Logo失败: %v", err)
	}

	go cache.DeleteOrganizationInfo(orgID)
	return nil
}

// UpdateOrganizationLocation 根据经纬度更新组织的9位Geohash地理位置编码
func (s *OrganizationService) UpdateOrganizationLocation(orgID uint, latitude, longitude float64) error {

	org, err := s.orgRepo.GetByID(orgID)
	if err != nil {
		return fmt.Errorf("查询组织失败: %v", err)
	}
	if org == nil {
		return fmt.Errorf("组织不存在")
	}

	// 验证经纬度有效性（经纬度由前端浏览器获取，一般有效）
	if latitude < -90 || latitude > 90 {
		return fmt.Errorf("纬度值无效，应在-90到90之间")
	}
	if longitude < -180 || longitude > 180 {
		return fmt.Errorf("经度值无效，应在-180到180之间")
	}

	// 计算9位精度的Geohash编码
	locationCode := geohash.EncodeWithPrecision(latitude, longitude, 9)

	// 更新组织的位置编码
	org.LocationCode = locationCode
	err = s.orgRepo.Update(org)
	if err != nil {
		return fmt.Errorf("更新组织位置编码失败: %v", err)
	}

	// 清理缓存（异步）
	go cache.DeleteOrganizationInfo(orgID)

	return nil
}

// TransferOrganizationOwnership 转移组织所有权
func (s *OrganizationService) TransferOrganizationOwnership(orgID uint, newCreatorID uint, processedBy uint) error {
	// 查询组织信息
	org, err := s.orgRepo.GetByID(orgID)
	if err != nil {
		return fmt.Errorf("查询组织失败: %v", err)
	}

	// 记录原创建者ID
	oldCreatorID := org.CreatorID

	// 更新组织表的创建者字段
	org.CreatorID = newCreatorID
	err = s.orgRepo.Update(org)
	if err != nil {
		return fmt.Errorf("更新组织创建者失败: %v", err)
	}

	// 将原创建者（当前操作者）在成员表中的角色降级为 'admin'
	err = s.orgMemberRepo.UpdateRole(orgID, oldCreatorID, "admin")
	if err != nil {
		return fmt.Errorf("降级原先组织者失败，原因：%v", err)
	}

	// 将新创建者在成员表中的角色提升为 'creator'
	err = s.orgMemberRepo.UpdateRole(orgID, newCreatorID, "creator")
	if err != nil {
		return fmt.Errorf("创建新组织者失败，原因：%v", err)
	}

	// 清理缓存
	go cache.DeleteOrganizationInfo(orgID)

	return nil
}

// CreateAdmin 组织者创建新的管理员
func (s *OrganizationService) CreateAdmin(orgID, userID uint) error {
	err := s.orgMemberRepo.UpdateRole(orgID, userID, "admin")
	if err != nil {
		return fmt.Errorf("创建新管理员失败，原因:%v", err)
	}
	return nil
}

// CancelAdmin 组织者取消管理员（降级）
func (s *OrganizationService) CancelAdmin(orgID, userID uint) error {
	err := s.orgMemberRepo.UpdateRole(orgID, userID, "member")
	if err != nil {
		return fmt.Errorf("取消管理员失败，原因：%v", err)
	}
	return nil
}

// FindNearbyOrganizationsByGeohashPrefix 通过Geohash前缀查询附近组织
func (s *OrganizationService) FindNearbyOrganizationsByGeohashPrefix(userLat, userLng float64) ([]models.OrgInfo, error) {
	// 验证经纬度有效性
	if userLat < -90 || userLat > 90 || userLng < -180 || userLng > 180 {
		return nil, fmt.Errorf("无效的经纬度参数")
	}

	// 生成用户位置的Geohash编码（9位精度）
	fullGeohash := geohash.EncodeWithPrecision(userLat, userLng, 9)

	// 取前8位作为查询前缀（约±19米精度）
	searchPrefix := fullGeohash[:8]

	// 使用前缀查询数据库
	orgInfos, err := s.orgRepo.FindByGeohashPrefix(searchPrefix)
	if err != nil {
		return nil, fmt.Errorf("数据库查询失败: %v", err)
	}

	return orgInfos, nil
}

// CreateCustomInviteCode 管理员为组织创建自定义邀请码
func (s *OrganizationService) CreateCustomInviteCode(orgID uint, customCode string) error {
	// 验证自定义验证码格式（6位数字）
	if !isValidCustomCode(customCode) {
		return fmt.Errorf("验证码必须是6位数字")
	}

	// 创建自定义验证码记录
	verificationCode := &models.VerificationCode{
		Code:           customCode,
		Business:       "join_organization",
		OrganizationID: orgID,
		ExpiresAt:      time.Now().Add(3 * time.Minute),
		Used:           false,
	}

	// 保存到数据库
	err := s.codeRepo.Create(verificationCode)
	if err != nil {
		return fmt.Errorf("创建自定义邀请码失败: %v", err)
	}

	return nil
}

// 验证自定义验证码格式（6位数字）
func isValidCustomCode(code string) bool {
	if len(code) != 6 {
		return false
	}
	// 检查是否全部为数字
	for _, char := range code {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

// VerifyAndJoinOrganization 验证邀请码并加入组织
func (s *OrganizationService) VerifyAndJoinOrganization(userID uint, organizationID uint, code string) error {
	// 查询有效的组织加入验证码
	validCode, err := s.codeRepo.FindValidJoinOrganizationCode(organizationID, code)
	if err != nil {
		return fmt.Errorf("验证码查询失败: %v", err)
	}
	if validCode == nil {
		return fmt.Errorf("邀请码无效、已过期或类型不正确")
	}

	// 检查用户是否已是组织成员
	isMember, err := s.orgMemberRepo.Exists(organizationID, userID)
	if err != nil {
		return fmt.Errorf("检查成员关系失败: %v", err)
	}
	if isMember {
		return fmt.Errorf("您已是该组织成员，无需重复加入")
	}

	// 将用户加入组织
	newMember := &models.OrganizationMember{
		OrganizationID: organizationID,
		UserID:         userID,
		Role:           "member",
		Status:         "active",
		JoinedAt:       time.Now(),
	}

	err = s.orgMemberRepo.Create(newMember)
	if err != nil {
		return fmt.Errorf("加入组织失败: %v", err)
	}

	cache.DeleteUserOrganizationOverviews(userID)
	return nil
}
