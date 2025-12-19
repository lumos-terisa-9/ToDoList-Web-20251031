package services

import (
	"fmt"
	"team_task_hub/backend/internal/models"
	"team_task_hub/backend/internal/repositories"
	"time"
)

// ActivityService 活动服务层
type ActivityService struct {
	activityRepo      *repositories.ActivityRepository
	participationRepo *repositories.ActivityParticipationRepository
}

// NewActivityService 创建活动服务实例
func NewActivityService(activityRepo *repositories.ActivityRepository,
	participationRepo *repositories.ActivityParticipationRepository) *ActivityService {
	return &ActivityService{
		activityRepo:      activityRepo,
		participationRepo: participationRepo,
	}
}

// CreateActivityRequest 创建活动请求参数
type CreateActivityRequest struct {
	Title              string    `json:"title" binding:"required,max=200"`
	Description        string    `json:"description"`
	StartTime          time.Time `json:"start_time" binding:"required"`
	EndTime            time.Time `json:"end_time" binding:"required"`
	ParticipationLimit string    `json:"participation_limit" binding:"omitempty,oneof=public org_only admin_assign"`
}

// CreateActivity 创建活动
func (s *ActivityService) CreateActivity(organizationID uint, req *CreateActivityRequest) error {
	// 验证活动时间逻辑
	if req.EndTime.Before(req.StartTime) {
		return fmt.Errorf("活动结束时间不能早于开始时间")
	}

	// 验证活动时间是否在未来
	if req.StartTime.Before(time.Now()) {
		return fmt.Errorf("活动开始时间不能是过去时间")
	}

	// 创建活动模型
	activity := &models.Activity{
		Title:              req.Title,
		Description:        req.Description,
		OrganizationID:     organizationID,
		StartTime:          req.StartTime,
		EndTime:            req.EndTime,
		ParticipationLimit: req.ParticipationLimit,
	}

	// 保存到数据库
	err := s.activityRepo.Create(activity)
	if err != nil {
		return fmt.Errorf("保存活动数据失败: %v", err)
	}

	return nil
}

// CancelActivity 取消活动
func (s *ActivityService) CancelActivity(activityID uint, reason string) error {
	// 将对应参与记录表中状态改成已取消
	err := s.participationRepo.CancelByActivityID(activityID)
	if err != nil {
		return fmt.Errorf("取消活动参与记录失败: %v", err)
	}

	updateData := map[string]any{
		"status":              "cancelled",
		"cancellation_reason": reason,
	}

	// 将活动表中对应状态改成已取消
	err = s.activityRepo.UpdateStatus(activityID, updateData)
	if err != nil {
		return fmt.Errorf("取消活动失败: %v", err)
	}

	return nil
}

// UpdateActivityRequest 更新活动请求
type UpdateActivityRequest struct {
	Title              string    `json:"title,omitempty" binding:"max=200" example:"更新后的活动标题"`
	Description        string    `json:"description,omitempty" example:"更新后的活动描述"`
	StartTime          time.Time `json:"start_time,omitempty" example:"2024-01-20T09:00:00Z"`
	EndTime            time.Time `json:"end_time,omitempty" example:"2024-01-20T17:00:00Z"`
	ParticipationLimit string    `json:"participation_limit,omitempty" binding:"oneof=public org_only admin_assign" example:"public"`
}

// UpdateActivity 更新活动信息
func (s *ActivityService) UpdateActivity(activityID uint, req *UpdateActivityRequest) error {

	// 验证时间逻辑
	if !req.StartTime.IsZero() && !req.EndTime.IsZero() {
		if req.EndTime.Before(req.StartTime) {
			return fmt.Errorf("活动结束时间不能早于开始时间")
		}
	}

	// 构建更新数据映射
	updateData := make(map[string]any)

	// 只更新提供了的字段（部分更新）
	if req.Title != "" {
		updateData["title"] = req.Title
	}
	if req.Description != "" {
		updateData["description"] = req.Description
	}
	if !req.StartTime.IsZero() {
		updateData["start_time"] = req.StartTime
	}
	if !req.EndTime.IsZero() {
		updateData["end_time"] = req.EndTime
	}
	if req.ParticipationLimit != "" {
		updateData["participation_limit"] = req.ParticipationLimit
	}

	// 如果没有提供任何要更新的字段，直接返回成功
	if len(updateData) == 0 {
		return nil
	}

	// 更新活动信息
	err := s.activityRepo.UpdateStatus(activityID, updateData)
	if err != nil {
		return fmt.Errorf("更新活动失败: %v", err)
	}

	return nil
}

// GetOrgActivities 根据活动参与限制获取组织活动信息
func (s *ActivityService) GetOrgActivities(orgID uint, participationType string) ([]models.Activity, error) {
	activities, err := s.activityRepo.FindByOrgAndParticipation(orgID, participationType)
	if err != nil {
		return nil, fmt.Errorf("获取组织活动失败，原因：%v", err)
	}

	return activities, nil
}

// GetActivityUsers 根据活动id获取活动参与者基本信息（姓名，账号，头像）
func (s *ActivityService) GetActivityUsers(activityID uint) ([]models.UserInfo, error) {
	userInfos, err := s.participationRepo.FindUsersByActivityID(activityID)
	if err != nil {
		return nil, fmt.Errorf("获取活动的用户信息失败，原因:%v", err)
	}

	return userInfos, nil
}

// ParticipateActivity 参加/领取活动
func (s *ActivityService) ParticipateActivity(activityID, userID uint) error {
	//检查是否参加过该活动
	if exist, err := s.participationRepo.Exists(activityID, userID); err != nil {
		return fmt.Errorf("查询是否存在重复记录失败，原因：%v", err)
	} else if exist {
		return fmt.Errorf("已经参加过本次活动，不能重复参加")
	}

	//在活动参与表中插入新活动参与记录
	activityParticipation := &models.ActivityParticipation{
		ActivityID: activityID,
		UserID:     userID,
	}

	err := s.participationRepo.Create(activityParticipation)
	if err != nil {
		return fmt.Errorf("参加活动失败,插入活动记录表失败，原因：%v", err)
	}

	return nil
}

// ForceAssignActivityToUsers 强制分配活动给一批用户
func (s *ActivityService) ForceAssignActivityToUsers(activityID uint, userIDs []uint) error {
	// 基础验证
	if len(userIDs) == 0 {
		return fmt.Errorf("用户ID列表不能为空")
	}

	// 调用数据访问层执行批量分配
	return s.participationRepo.BatchCreateForcedAssignments(activityID, userIDs)
}

// GetUnreadActivitiesByUserID 获取用户所有未读活动
func (s *ActivityService) GetUnreadActivitiesByUserID(userID uint) ([]models.Activity, error) {
	activities, err := s.activityRepo.FindUnreadActivitiesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户未读活动失败，原因：%v", err)
	}

	return activities, err
}

// MarkActivityAsRead 标记活动为已读
func (s *ActivityService) MarkActivityAsRead(userID, activityID uint) error {
	if err := s.participationRepo.MarkActivityAsRead(userID, activityID); err != nil {
		return fmt.Errorf("将活动标记为已读失败，原因:%v", err)
	}

	return nil
}

// GetCancelledActivitiesByUserID 获取用户已经取消的活动，组织取消的
func (s *ActivityService) GetCancelledActivitiesByUserID(userID uint) ([]models.Activity, error) {
	activities, err := s.activityRepo.FindCancelledActivitiesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户已取消活动失败，原因：%v", err)
	}

	return activities, err
}

// DeleteCancelledActivity 删除用户已取消的活动记录
func (s *ActivityService) DeleteCancelledActivity(userID, activityID uint) error {
	err := s.participationRepo.DeleteCancelledActivity(userID, activityID)
	if err != nil {
		return fmt.Errorf("删除已取消活动失败: %v", err)
	}

	return nil
}

// GetParticipationStatus 批量获取用户的参与状态
func (s *ActivityService) GetParticipationStatus(activityID uint, userIDs []uint) ([]uint, error) {
	var participatedIDs []uint
	// 参数校验
	if len(userIDs) == 0 {
		return participatedIDs, nil // 空请求返回空切片
	}

	// 调用数据访问层，获取已参与此活动的用户ID列表
	participatedIDs, err := s.participationRepo.FindParticipatedUserIDs(activityID, userIDs)
	if err != nil {
		return nil, fmt.Errorf("查询参与状态失败: %v", err)
	}

	return participatedIDs, nil
}

// GetActivityByID 根据ID获取活动详情
func (s *ActivityService) GetActivityByID(id uint) (*models.Activity, error) {
	activity, err := s.activityRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取活动详情失败: %v", err)
	}
	return activity, nil
}

// GetActivitiesByOrganization 获取组织的所有活动
func (s *ActivityService) GetActivitiesByOrganization(orgID uint) ([]models.Activity, error) {
	activities, err := s.activityRepo.GetByOrganizationID(orgID)
	if err != nil {
		return nil, fmt.Errorf("获取组织活动列表失败: %v", err)
	}
	return activities, nil
}

// JoinActivity 用户参与活动
func (s *ActivityService) JoinActivity(activityID, userID uint) error {
	// 检查活动是否存在
	activity, err := s.activityRepo.GetByID(activityID)
	if err != nil {
		return fmt.Errorf("活动不存在: %v", err)
	}

	// 检查活动是否已开始
	if activity.StartTime.Before(time.Now()) {
		return fmt.Errorf("活动已开始，无法参与")
	}

	// 检查用户是否已参与该活动
	exists, err := s.participationRepo.Exists(activityID, userID)
	if err != nil {
		return fmt.Errorf("检查参与状态失败: %v", err)
	}
	if exists {
		return fmt.Errorf("您已参与此活动")
	}

	// 创建参与记录
	participation := &models.ActivityParticipation{
		ActivityID: activityID,
		UserID:     userID,
	}

	err = s.participationRepo.Create(participation)
	if err != nil {
		return fmt.Errorf("参与活动失败: %v", err)
	}

	return nil
}

// SubmitReview 提交活动评价
func (s *ActivityService) SubmitReview(activityID, userID uint, rating int, reviewText string) error {
	// 验证评分范围
	if rating < 1 || rating > 10 {
		return fmt.Errorf("评分必须在1到10之间")
	}

	// 查找参与记录
	participation, err := s.participationRepo.GetByActivityAndUser(activityID, userID)
	if err != nil {
		return fmt.Errorf("查询参与记录失败: %v", err)
	}
	if participation == nil {
		return fmt.Errorf("未找到您的参与记录")
	}

	// 检查活动是否已结束（只有参与过的活动才能评价）
	activity, err := s.activityRepo.GetByID(activityID)
	if err != nil {
		return fmt.Errorf("获取活动信息失败: %v", err)
	}

	if activity.EndTime.After(time.Now()) {
		return fmt.Errorf("活动尚未结束，无法评价")
	}

	// 更新评价信息
	participation.Rating = rating
	participation.ReviewText = reviewText

	err = s.participationRepo.Update(participation)
	if err != nil {
		return fmt.Errorf("提交评价失败: %v", err)
	}

	return nil
}

// GetActivityParticipants 获取活动参与者列表
func (s *ActivityService) GetActivityParticipants(activityID uint) ([]models.ActivityParticipation, error) {
	participants, err := s.participationRepo.GetByActivityID(activityID)
	if err != nil {
		return nil, fmt.Errorf("获取参与者列表失败: %v", err)
	}
	return participants, nil
}

// GetUserActivities 获取用户参与的活动列表
func (s *ActivityService) GetUserActivities(userID uint) ([]models.ActivityParticipation, error) {
	participations, err := s.participationRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户活动列表失败: %v", err)
	}
	return participations, nil
}

// GetActivityWithParticipants 获取活动详情及参与者信息
func (s *ActivityService) GetActivityWithParticipants(activityID uint) (*models.Activity, []models.ActivityParticipation, error) {
	// 获取活动详情
	activity, err := s.activityRepo.GetByID(activityID)
	if err != nil {
		return nil, nil, fmt.Errorf("获取活动详情失败: %v", err)
	}

	// 获取参与者列表
	participants, err := s.participationRepo.GetByActivityID(activityID)
	if err != nil {
		return nil, nil, fmt.Errorf("获取参与者列表失败: %v", err)
	}

	return activity, participants, nil
}

// CalculateActivityAverageRating 计算活动平均评分
func (s *ActivityService) CalculateActivityAverageRating(activityID uint) (float64, int, error) {
	participants, err := s.participationRepo.GetByActivityID(activityID)
	if err != nil {
		return 0, 0, fmt.Errorf("获取参与者列表失败: %v", err)
	}

	if len(participants) == 0 {
		return 0, 0, nil
	}

	totalRating := 0
	ratingCount := 0

	for _, participant := range participants {
		if participant.Rating > 0 { // 只计算已评分的
			totalRating += participant.Rating
			ratingCount++
		}
	}

	if ratingCount == 0 {
		return 0, 0, nil
	}

	averageRating := float64(totalRating) / float64(ratingCount)
	return averageRating, ratingCount, nil
}
