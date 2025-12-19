package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"team_task_hub/backend/internal/models"
	"team_task_hub/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	orgService      *services.OrganizationService
	activityService *services.ActivityService
}

func NewOrganizationHandler(orgService *services.OrganizationService, activityService *services.ActivityService) *OrganizationHandler {
	return &OrganizationHandler{
		orgService:      orgService,
		activityService: activityService,
	}
}

// SuccessResponse 成功响应结构体
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CreateOrganizationApplicationRequest 创建组织申请请求
type CreateOrganizationApplicationRequest struct {
	OrganizationName      string `json:"organization_name" binding:"required,min=1,max=100"`
	ApplicationReason     string `json:"application_reason" binding:"required,min=1"`
	ApplicantIntroduction string `json:"applicant_introduction,omitempty"`
	AttachmentURL         string `json:"attachment_url,omitempty"`
}

// SubmitCreateOrganizationApplicationHandler 提交创建组织申请
// @Summary 提交创建组织申请
// @Description 用户提交创建新组织的申请，需要管理员审批
// @Tags 组织申请
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body CreateOrganizationApplicationRequest true "创建组织申请请求"
// @Success 201 {object} SuccessResponse "申请提交成功" example({"success": true, "message": "创建组织申请提交成功"})
// @Failure 400 {object} ErrorResponse "请求参数错误" example({"success": false, "message": "组织名称已存在"})
// @Failure 401 {object} ErrorResponse "用户未认证"
// @Failure 500 {object} ErrorResponse "系统内部错误" example({"success": false, "message": "提交申请失败"})
// @Router /api/organization/application/create-organization [post]
func (h *OrganizationHandler) SubmitCreateOrganizationApplicationHandler(c *gin.Context) {
	var request struct {
		OrganizationName      string `json:"organization_name" binding:"required,min=1,max=100"`
		ApplicationReason     string `json:"application_reason" binding:"required,min=1"`
		ApplicantIntroduction string `json:"applicant_introduction,omitempty"`
		AttachmentURL         string `json:"attachment_url,omitempty"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "请求参数无效: " + err.Error(),
		})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "用户未认证",
		})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "用户信息不完整",
		})
		return
	}

	application := &models.OrganizationApplication{
		ApplicantUserID:       userID.(uint),
		ApplicantUsername:     username.(string),
		OrganizationName:      request.OrganizationName,
		ApplicationReason:     request.ApplicationReason,
		ApplicantIntroduction: request.ApplicantIntroduction,
		AttachmentURL:         request.AttachmentURL,
	}

	err := h.orgService.SubmitCreateOrganizationApplication(application)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "组织名称已存在" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, ErrorResponse{
			Success: false,
			Message: "提交创建组织申请失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Success: true,
		Message: "创建组织申请提交成功",
	})
}

type JoinOrganizationApplicationRequest struct {
	OrganizationName      string `json:"organization_name" binding:"required,min=1,max=100"`
	ApplicationReason     string `json:"application_reason" binding:"required,min=1"`
	ApplicantIntroduction string `json:"applicant_introduction,omitempty"`
	AttachmentURL         string `json:"attachment_url,omitempty"`
}

// SubmitJoinApplicationHandler 提交加入组织申请
// @Summary 提交加入组织申请
// @Description 用户提交加入现有组织的申请，需要组织管理员审批
// @Tags 组织申请
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body JoinOrganizationApplicationRequest true "加入组织申请请求"
// @Success 201 {object} SuccessResponse "申请提交成功" example({"success": true, "message": "加入组织申请提交成功"})
// @Failure 400 {object} ErrorResponse "请求参数错误" example({"success": false, "message": "目标组织不存在"})
// @Failure 401 {object} ErrorResponse "用户未认证"
// @Failure 500 {object} ErrorResponse "系统内部错误" example({"success": false, "message": "提交申请失败"})
// @Router /api/organization/application/join-organization [post]
func (h *OrganizationHandler) SubmitJoinApplicationHandler(c *gin.Context) {
	var request struct {
		OrganizationName      string `json:"organization_name" binding:"required,min=1,max=100"`
		ApplicationReason     string `json:"application_reason" binding:"required,min=1"`
		ApplicantIntroduction string `json:"applicant_introduction,omitempty"`
		AttachmentURL         string `json:"attachment_url,omitempty"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "请求参数无效: " + err.Error(),
		})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "用户未认证",
		})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "用户信息不完整",
		})
		return
	}

	application := &models.OrganizationApplication{
		ApplicantUserID:       userID.(uint),
		ApplicantUsername:     username.(string),
		OrganizationName:      request.OrganizationName,
		ApplicationReason:     request.ApplicationReason,
		ApplicantIntroduction: request.ApplicantIntroduction,
		AttachmentURL:         request.AttachmentURL,
	}

	err := h.orgService.SubmitJoinApplication(application)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMessage := err.Error()
		if errorMessage == "目标组织不存在" || errorMessage == "您已是该组织成员" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, ErrorResponse{
			Success: false,
			Message: "提交加入组织申请失败: " + errorMessage,
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Success: true,
		Message: "加入组织申请提交成功",
	})
}

type ProcessApplicationRequest struct {
	Action string `json:"action" binding:"required,oneof=approve reject" example:"approve"`
	Remark string `json:"remark,omitempty" example:"申请符合要求"`
}

// ProcessApplicationHandler 处理组织申请
// @Summary 处理组织申请
// @Description 系统/组织管理员审批或拒绝组织创建/加入申请
// @Tags 组织申请
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param id path int true "申请ID"
// @Param request body ProcessApplicationRequest true "处理请求" example({"action":"approve","remark":"申请符合要求"})
// @Success 200 {object} SuccessResponse "处理成功" example({"success": true, "message": "申请处理成功"})
// @Failure 400 {object} ErrorResponse "请求参数错误" example({"success": false, "message": "申请不存在"})
// @Failure 401 {object} ErrorResponse "用户未认证"
// @Failure 500 {object} ErrorResponse "系统内部错误" example({"success": false, "message": "处理申请失败"})
// @Router /api/organization/application/{id}/process [patch]
func (h *OrganizationHandler) ProcessApplicationHandler(c *gin.Context) {
	var request struct {
		Action string `json:"action" binding:"required,oneof=approve reject"`
		Remark string `json:"remark,omitempty"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "请求参数无效: " + err.Error(),
		})
		return
	}

	applicationIDStr := c.Param("id")
	applicationID, err := strconv.ParseUint(applicationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "无效的申请ID",
		})
		return
	}

	processedBy, _ := c.Get("userID")
	err = h.orgService.ProcessApplication(uint(applicationID), request.Action, request.Remark, processedBy.(uint))
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMessage := err.Error()
		if errorMessage == "申请不存在" || errorMessage == "该申请已被处理" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, ErrorResponse{
			Success: false,
			Message: "处理申请失败: " + errorMessage,
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "申请处理成功",
	})
}

// GetUserOrganizationOverviewsHandler 获取用户组织总览
// @Summary 获取用户组织总览
// @Description 获取用户加入的所有组织信息，包括组织基本信息和加入时间
// @Tags 组织管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Success 200 {object} SuccessResponse "获取成功" example({"success": true, "message": "获取组织列表成功"})
// @Failure 401 {object} ErrorResponse "用户未认证"
// @Failure 500 {object} ErrorResponse "系统内部错误" example({"success": false, "message": "获取组织列表失败"})
// @Router /api/organization/my-organizations [get]
func (h *OrganizationHandler) GetUserOrganizationOverviewsHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "用户未认证",
		})
		return
	}

	overviews, err := h.orgService.GetUserOrganizationOverviews(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "获取组织列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取组织列表成功",
		"data":    overviews,
	})
}

// GetPendingApplicationsHandler 获取待处理申请（支持过滤）
// @Summary 获取待处理申请列表
// @Description 管理员获取待处理申请，可通过type参数过滤
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param type query string false "申请类型: create_org（创建组织）,change_org（修改组织）,join_org（加入组织）"
// @Success 200 {object} SuccessResponse "获取成功"
// @Router /admin/api/organization/application/pending [get]
func (h *OrganizationHandler) GetPendingApplicationsHandler(c *gin.Context) {
	appType := c.Query("type")

	applications, err := h.orgService.GetPendingApplications(appType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "获取申请列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取申请列表成功",
		"data":    applications,
	})
}

// GetOrgPendingJoinApplicationsHandler 获取组织的待处理加入申请
// @Summary 获取组织的待处理加入申请
// @Description 组织管理员获取该组织的待处理加入申请
// @Tags 组织申请
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param orgName query string true "组织名称"
// @Success 200 {object} SuccessResponse "获取成功" example({"success": true, "message": "获取加入申请成功"})
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 401 {object} ErrorResponse "用户未认证"
// @Failure 500 {object} ErrorResponse "系统内部错误" example({"success": false, "message": "获取申请失败"})
// @Router /api/orgnization/application/pending-join [get]
func (h *OrganizationHandler) GetOrgPendingJoinApplicationsHandler(c *gin.Context) {
	orgName := c.Query("orgName")
	if orgName == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "组织名称不能为空",
		})
		return
	}

	applications, err := h.orgService.OrgGetPendingJoinApplications(orgName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "获取加入组织申请失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取加入组织申请成功",
		"data":    applications,
	})
}

// GetUserApplicationsHandler 获取用户的申请列表
// @Summary 获取用户的申请列表
// @Description 获取用户提交的所有组织相关申请
// @Tags 组织申请
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Success 200 {object} SuccessResponse "获取成功" example({"success": true, "message": "获取申请列表成功"})
// @Failure 401 {object} ErrorResponse "用户未认证"
// @Failure 500 {object} ErrorResponse "系统内部错误" example({"success": false, "message": "获取申请列表失败"})
// @Router /api/orgnization/application/my-applications [get]
func (h *OrganizationHandler) GetUserApplicationsHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "用户未认证",
		})
		return
	}

	applications, err := h.orgService.GetUserApplications(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "获取申请列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取申请列表成功",
		"data":    applications,
	})
}

// DeleteUserApplicationHandler 删除用户申请
// @Summary 删除用户申请
// @Description 用户删除自己提交的申请
// @Tags 组织申请
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param id path int true "申请ID"
// @Success 200 {object} SuccessResponse "删除成功" example({"success": true, "message": "申请删除成功"})
// @Failure 400 {object} ErrorResponse "请求参数错误" example({"success": false, "message": "申请ID无效"})
// @Failure 401 {object} ErrorResponse "用户未认证"
// @Failure 500 {object} ErrorResponse "系统内部错误" example({"success": false, "message": "删除申请失败"})
// @Router /api/organization/application/{id} [delete]
func (h *OrganizationHandler) DeleteUserApplicationHandler(c *gin.Context) {
	applicationIDStr := c.Param("id")
	applicationID, err := strconv.ParseUint(applicationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "无效的申请ID",
		})
		return
	}

	err = h.orgService.DeleteUserApplication(uint(applicationID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "删除申请失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "申请删除成功",
	})
}

type RemoveOrganizationMemberRequest struct {
	OrgID  uint `json:"org_id" binding:"required" example:"1"`
	UserID uint `json:"user_id" binding:"required" example:"2"`
}

// RemoveOrganizationMemberHandler 移除组织成员
// @Summary 移除组织成员
// @Description 组织管理员或创建者从组织中移除成员
// @Tags 组织管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body RemoveOrganizationMemberRequest true "移除请求" example({"org_id": 1, "user_id": 2})
// @Success 200 {object} SuccessResponse "移除成功" example({"success": true, "message": "成员移除成功"})
// @Failure 400 {object} ErrorResponse "请求参数错误" example({"success": false, "message": "成员关系不存在"})
// @Failure 401 {object} ErrorResponse "用户未认证"
// @Failure 500 {object} ErrorResponse "系统内部错误" example({"success": false, "message": "移除成员失败"})
// @Router /api/organization/remove-member [delete]
func (h *OrganizationHandler) RemoveOrganizationMemberHandler(c *gin.Context) {
	var request struct {
		OrgID  uint `json:"org_id" binding:"required"`
		UserID uint `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "请求参数无效: " + err.Error(),
		})
		return
	}

	removedBy, _ := c.Get("userID")
	err := h.orgService.RemoveOrganizationMember(request.OrgID, request.UserID, removedBy.(uint))
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMessage := err.Error()
		if errorMessage == "成员关系不存在" || errorMessage == "不能移除组织创建者" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, ErrorResponse{
			Success: false,
			Message: "移除成员失败: " + errorMessage,
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "成员移除成功",
	})
}

// SearchOrganizationsHandler 搜索组织
// @Summary 搜索组织
// @Description 根据组织名称进行模糊搜索，返回匹配的组织简略信息列表
// @Tags 组织查询
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param keyword query string true "搜索关键词"
// @Success 200 {object} SuccessResponse "搜索成功" example({"success": true, "message": "搜索成功", "data": [{"id": 1, "name": "软工羽队"}]})
// @Failure 400 {object} ErrorResponse "请求参数错误" example({"success": false, "message": "搜索关键词不能为空"})
// @Failure 401 {object} ErrorResponse "用户未认证"
// @Failure 500 {object} ErrorResponse "系统内部错误" example({"success": false, "message": "搜索失败: 数据库连接错误"})
// @Router /api/organization/search [get]
func (h *OrganizationHandler) SearchOrganizationsHandler(c *gin.Context) {
	// 从查询参数获取搜索关键词
	keyword := c.Query("keyword")

	// 验证关键词不能为空
	if strings.TrimSpace(keyword) == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "搜索关键词不能为空",
		})
		return
	}

	// 调用服务层进行搜索
	results, err := h.orgService.FindOrgInfoByName(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "搜索失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "搜索成功",
		"data":    results,
	})
}

// GetOrganizationByIDHandler 根据组织ID获取组织详情
// @Summary 获取组织详情
// @Description 根据组织ID获取组织的完整详细信息。该接口集成了缓存，高频请求将直接返回缓存结果。
// @Tags 组织查询
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param orgID path int true "组织ID" minimum(1)
// @Success 200 {object} SuccessResponse "获取成功" example({"success": true, "message": "获取组织成功", "data": {"id": 1, "name": "技术部"}})
// @Failure 400 {object} ErrorResponse "请求参数错误" example({"success": false, "message": "无效的组织ID"})
// @Failure 404 {object} ErrorResponse "组织不存在"
// @Failure 500 {object} ErrorResponse "系统内部错误" example({"success": false, "message": "查询组织失败: 组织不存在"})
// @Router /api/organization/{orgID} [get]
func (h *OrganizationHandler) GetOrganizationByIDHandler(c *gin.Context) {
	// 获取并验证路径参数
	orgIDStr := c.Param("orgID")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil || orgID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的组织ID，必须为正整数",
		})
		return
	}

	// 调用服务层方法
	organization, err := h.orgService.FindOrgByID(uint(orgID))
	if err != nil {
		// 根据服务层返回的错误信息细化HTTP状态码
		statusCode := http.StatusInternalServerError
		errorMessage := err.Error()

		// 判断是否为"未找到"的错误
		if strings.Contains(errorMessage, "不存在") || organization == nil {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": "获取组织信息失败: " + errorMessage,
		})
		return
	}

	//  返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取组织成功",
		"data":    organization,
	})
}

type UpdateOrganizationNameRequest struct {
	NewName string `json:"new_name" binding:"required,min=1,max=100" example:"新技术研发部"`
}

// UpdateOrganizationNameHandler 更新组织名称
// @Summary 更新组织名称
// @Description 管理员手动更新指定组织的名称信息，example({"new_name": "新技术研发部"})
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param orgID path int true "组织ID"
// @Param request body UpdateOrganizationNameRequest true "更新请求" example({"new_name": "新技术研发部"})
// @Success 200 {object} SuccessResponse "更新成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /admin/api/organization/{orgID}/name [patch]
func (h *OrganizationHandler) UpdateOrganizationNameHandler(c *gin.Context) {
	// 从路径参数获取组织ID
	orgIDStr := c.Param("orgID")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil || orgID == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "无效的组织ID",
		})
		return
	}

	// 绑定请求体
	var request struct {
		NewName string `json:"new_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "请求参数无效: " + err.Error(),
		})
		return
	}

	// 调用服务层
	err = h.orgService.UpdateOrganizationName(uint(orgID), request.NewName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "更新组织名称失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "组织名称更新成功",
	})
}

type UpdateOrganizationDescriptionRequest struct {
	NewDescription string `json:"new_description" binding:"required,max=500" example:"这是更新后的组织描述"`
}

// UpdateOrganizationDescriptionHandler 更新组织描述
// @Summary 更新组织描述
// @Description 管理员手动更新指定组织的描述信息，example({"new_description": "这是更新后的组织描述"})
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param orgID path int true "组织ID"
// @Param request body UpdateOrganizationDescriptionRequest true "更新请求" example({"new_description": "这是更新后的组织描述"})
// @Success 200 {object} SuccessResponse "更新成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /admin/api/organization/{orgID}/description [patch]
func (h *OrganizationHandler) UpdateOrganizationDescriptionHandler(c *gin.Context) {
	orgIDStr := c.Param("orgID")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil || orgID == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "无效的组织ID",
		})
		return
	}

	var request struct {
		NewDescription string `json:"new_description" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "请求参数无效: " + err.Error(),
		})
		return
	}

	err = h.orgService.UpdateOrganizationDescription(uint(orgID), request.NewDescription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "更新组织描述失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "组织描述更新成功",
	})
}

// UpdateOrganizationLogoRequest 更新组织Logo请求
type UpdateOrganizationLogoRequest struct {
	NewLogoURL string `json:"new_logo_url" binding:"required,url,max=255" example:"https://example.com/logo.png"`
}

// UpdateOrganizationLogoHandler 更新组织Logo
// @Summary 更新组织Logo
// @Description 管理员手动更新指定组织的Logo图片URL，example({"new_logo_url": "https://example.com/logo.png"})
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param orgID path int true "组织ID"
// @Param request body UpdateOrganizationLogoRequest true "更新请求" example({"new_logo_url": "https://example.com/logo.png"})
// @Success 200 {object} SuccessResponse "更新成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /admin/api/organization/{orgID}/logo [patch]
func (h *OrganizationHandler) UpdateOrganizationLogoHandler(c *gin.Context) {
	orgIDStr := c.Param("orgID")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil || orgID == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "无效的组织ID",
		})
		return
	}

	var request struct {
		NewLogoURL string `json:"new_logo_url" binding:"max=255"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "请求参数无效: " + err.Error(),
		})
		return
	}

	err = h.orgService.UpdateOrganizationLogo(uint(orgID), request.NewLogoURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "更新组织Logo失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "组织Logo更新成功",
	})
}

// UpdateOrganizationLocationRequest 更新组织位置请求结构体
type UpdateOrganizationLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" binding:"required,min=-180,max=180"`
}

// UpdateOrganizationLocationHandler 更新组织地理位置
// @Summary 更新组织地理位置
// @Description 根据经纬度为组织生成8位精度的Geohash编码
// @Tags 组织管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param orgID path int true "组织ID"
// @Param request body UpdateOrganizationLocationRequest true "位置更新请求"
// @Success 200 {object} SuccessResponse "更新成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 401 {object} ErrorResponse "用户未认证"
// @Failure 404 {object} ErrorResponse "组织不存在"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /api/organization/{orgID}/location [put]
func (h *OrganizationHandler) UpdateOrganizationLocationHandler(c *gin.Context) {
	// 获取组织ID
	orgIDStr := c.Param("orgID")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil || orgID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的组织ID格式",
		})
		return
	}
	// 绑定和验证请求体
	var request UpdateOrganizationLocationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数无效: " + err.Error(),
		})
		return
	}
	// 调用服务层方法
	err = h.orgService.UpdateOrganizationLocation(uint(orgID), request.Latitude, request.Longitude)
	if err != nil {
		// 根据错误类型细化HTTP状态码
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()

		switch {
		case strings.Contains(errorMsg, "组织不存在"):
			statusCode = http.StatusNotFound
		case strings.Contains(errorMsg, "纬度值无效") || strings.Contains(errorMsg, "经度值无效"):
			statusCode = http.StatusBadRequest
		case strings.Contains(errorMsg, "权限"):
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": "更新组织位置失败: " + errorMsg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "组织地理位置更新成功",
	})
}

// TransferOrganizationOwnershipRequest 转移组织所有权请求
type TransferOrganizationOwnershipRequest struct {
	NewCreatorID uint `json:"new_creator_id" binding:"required,min=1" example:"123"`
}

// TransferOrganizationOwnershipHandler 转移组织所有权
// @Summary 转移组织所有权
// @Description 组织者将指定组织的所有权转移给其他成员，无需管理员审批
// @Tags 组织管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param orgID path int true "组织ID"
// @Param request body TransferOrganizationOwnershipRequest true "转移请求" example({"new_creator_id": 123})
// @Success 200 {object} SuccessResponse "转移成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /api/organization/{orgID}/ownership [patch]
func (h *OrganizationHandler) TransferOrganizationOwnershipHandler(c *gin.Context) {
	orgIDStr := c.Param("orgID")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil || orgID == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "无效的组织ID",
		})
		return
	}

	var request struct {
		NewCreatorID uint `json:"new_creator_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "请求参数无效: " + err.Error(),
		})
		return
	}

	// 从JWT token中获取当前操作者ID
	processedBy, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "用户未认证",
		})
		return
	}

	err = h.orgService.TransferOrganizationOwnership(
		uint(orgID),
		request.NewCreatorID,
		processedBy.(uint),
	)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()

		// 根据错误类型细化状态码
		if strings.Contains(errorMsg, "无权") ||
			strings.Contains(errorMsg, "不是该组织成员") {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, ErrorResponse{
			Success: false,
			Message: "转移组织所有权失败: " + errorMsg,
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "组织所有权转移成功",
	})
}

// SubmitChangeOrgAppRequest 提交组织变更申请请求结构
type SubmitChangeOrgAppRequest struct {
	OrganizationID        uint   `json:"organization_id" binding:"required,min=1"`
	OrganizationName      string `json:"organization_name" binding:"required,min=1,max=100"`
	ApplicationReason     string `json:"application_reason" binding:"required,min=1"`
	ApplicantIntroduction string `json:"applicant_introduction,omitempty"`
	AttachmentURL         string `json:"attachment_url,omitempty" max:"255"`
}

// SubmitChangeOrganizationApplicationHandler 提交组织变更申请
// @Summary 提交组织变更申请
// @Description 用户提交组织信息变更申请，需要管理员审批
// @Tags 组织申请
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body SubmitChangeOrgAppRequest true "变更申请请求"
// @Success 201 {object} SuccessResponse "申请提交成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 401 {object} ErrorResponse "用户未认证"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /api/application/change-organization [post]
func (h *OrganizationHandler) SubmitChangeOrganizationApplicationHandler(c *gin.Context) {
	// 绑定和验证请求数据
	var request SubmitChangeOrgAppRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "请求参数无效: " + err.Error(),
		})
		return
	}

	// 从JWT token中获取申请人信息
	applicantUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "用户未认证",
		})
		return
	}

	applicantUsername, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "用户信息不完整",
		})
		return
	}

	// 创建申请记录
	application := &models.OrganizationApplication{
		ApplicantUserID:       applicantUserID.(uint),
		ApplicantUsername:     applicantUsername.(string),
		OrganizationID:        request.OrganizationID,
		OrganizationName:      request.OrganizationName,
		ApplicationReason:     request.ApplicationReason,
		ApplicantIntroduction: request.ApplicantIntroduction,
		AttachmentURL:         request.AttachmentURL,
	}

	// 调用服务层提交申请
	err := h.orgService.SubmitChangeOrganizationApplication(application)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "提交变更申请失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusCreated, SuccessResponse{
		Success: true,
		Message: "组织变更申请提交成功，等待管理员审批",
	})
}

// PromoteToAdminHandler 提拔组织成员为管理员
// @Summary 提拔组织成员为管理员
// @Description 组织创建者将指定普通成员提拔为组织管理员
// @Tags 组织管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param orgID path int true "组织ID"
// @Param userID path int true "要提拔的用户ID"
// @Success 200 {object} SuccessResponse "提拔成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 403 {object} ErrorResponse "权限不足（非组织创建者操作）"
// @Failure 404 {object} ErrorResponse "组织或用户不存在"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /api/organization/{orgID}/admin/{userID} [patch]
func (h *OrganizationHandler) PromoteToAdminHandler(c *gin.Context) {
	// 解析并验证路径参数
	orgIDStr := c.Param("orgID")
	userIDStr := c.Param("userID")

	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil || orgID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的组织ID格式",
		})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的用户ID格式",
		})
		return
	}

	err = h.orgService.CreateAdmin(uint(orgID), uint(userID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": "提拔成员失败: " + errorMsg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "成员已成功提拔为管理员",
	})
}

// CancelAdminHandler 取消管理员权限（降级为普通成员）
// @Summary 取消管理员权限
// @Description 组织创建者将指定管理员降级为普通组织成员
// @Tags 组织管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param orgID path int true "组织ID"
// @Param userID path int true "要取消管理员权限的用户ID"
// @Success 200 {object} SuccessResponse "取消管理员权限成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 403 {object} ErrorResponse "权限不足（非组织创建者操作）"
// @Failure 404 {object} ErrorResponse "组织或用户不存在"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /api/organization/{orgID}/admin/{userID} [delete]
func (h *OrganizationHandler) CancelAdminHandler(c *gin.Context) {
	// 解析并验证路径参数
	orgIDStr := c.Param("orgID")
	userIDStr := c.Param("userID")

	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil || orgID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的组织ID格式",
		})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的用户ID格式",
		})
		return
	}

	// 调用服务层的CancelAdmin方法
	err = h.orgService.CancelAdmin(uint(orgID), uint(userID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": "取消管理员权限失败: " + errorMsg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "管理员权限已成功取消，该用户现为普通成员",
	})
}

// SearchNearbyOrganizationsHandler 查询附近组织
// @Summary 查询附近组织
// @Description 根据用户经纬度生成Geohash前缀，查询约50米范围内的附近组织
// @Tags 组织查询
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param lat query number true "纬度" minimum(-90) maximum(90)
// @Param lng query number true "经度" minimum(-180) maximum(180)
// @Success 200 {object} SuccessResponse "查询成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /api/organization/nearby [get]
func (h *OrganizationHandler) SearchNearbyOrganizationsHandler(c *gin.Context) {
	// 1. 获取和验证查询参数
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "纬度参数格式错误",
		})
		return
	}

	lng, err := strconv.ParseFloat(c.Query("lng"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "经度参数格式错误",
		})
		return
	}

	// 调用服务层方法
	organizations, err := h.orgService.FindNearbyOrganizationsByGeohashPrefix(lat, lng)
	if err != nil {
		// 根据错误类型细化HTTP状态码
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()

		if strings.Contains(errorMsg, "无效的经纬度参数") {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": "查询附近组织失败: " + errorMsg,
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "查询成功",
		"data": gin.H{
			"organizations": organizations,
			"count":         len(organizations),
		},
	})
}

// CreateInviteCodeRequest 创建邀请码请求结构
type CreateInviteCodeRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

// CreateCustomInviteCodeHandler 创建自定义组织邀请码
// @Summary 创建组织邀请码
// @Description 管理员为指定组织创建自定义6位数字邀请码，有效期为3分钟
// @Tags 组织管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param orgID path int true "组织ID"
// @Param request body CreateInviteCodeRequest true "邀请码创建请求"
// @Success 201 {object} SuccessResponse "创建成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 401 {object} ErrorResponse "用户未认证"
// @Failure 403 {object} ErrorResponse "权限不足"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /api/organization/{orgID}/invite-codes [post]
func (h *OrganizationHandler) CreateCustomInviteCodeHandler(c *gin.Context) {
	// 获取并验证组织ID
	orgIDStr := c.Param("orgID")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil || orgID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的组织ID格式",
		})
		return
	}

	// 绑定和验证请求体
	var request struct {
		Code string `json:"code" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数无效: " + err.Error(),
		})
		return
	}

	// 调用服务层方法
	err = h.orgService.CreateCustomInviteCode(uint(orgID), request.Code)
	if err != nil {
		// 根据错误类型细化HTTP状态码
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()

		if strings.Contains(errorMsg, "验证码必须是6位数字") {
			statusCode = http.StatusBadRequest
		} else if strings.Contains(errorMsg, "组织不存在") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(errorMsg, "权限") {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": "创建邀请码失败: " + errorMsg,
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "组织邀请码创建成功，有效期3分钟",
	})
}

// JoinWithCodeRequest 邀请码加入组织请求结构
type JoinWithCodeRequest struct {
	OrganizationID uint   `json:"organization_id" binding:"required,min=1"`
	Code           string `json:"code" binding:"required,len=6"`
}

// JoinOrganizationWithCodeHandler 使用邀请码加入组织
// @Summary 使用邀请码加入组织
// @Description 用户通过提交有效的邀请码直接加入指定组织（免审核）
// @Tags 加入组织
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body JoinWithCodeRequest true "加入请求"
// @Success 200 {object} SuccessResponse "加入成功"
// @Failure 400 {object} ErrorResponse "请求参数错误或邀请码无效"
// @Failure 409 {object} ErrorResponse "用户已是组织成员"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /api/organization/join-with-code [post]
func (h *OrganizationHandler) JoinOrganizationWithCodeHandler(c *gin.Context) {
	// 绑定和验证请求参数
	var request struct {
		OrganizationID uint   `json:"organization_id" binding:"required,min=1"`
		Code           string `json:"code" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数无效: " + err.Error(),
		})
		return
	}

	// 从JWT Token中获取当前用户ID
	currentUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "用户未认证",
		})
		return
	}

	// 调用服务层验证邀请码并加入组织
	err := h.orgService.VerifyAndJoinOrganization(
		currentUserID.(uint),
		request.OrganizationID,
		request.Code,
	)
	if err != nil {
		// 根据错误类型细化HTTP状态码
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()

		switch {
		case strings.Contains(errorMsg, "邀请码无效") ||
			strings.Contains(errorMsg, "验证码查询失败"):
			statusCode = http.StatusBadRequest
		case strings.Contains(errorMsg, "已是该组织成员"):
			statusCode = http.StatusConflict
		case strings.Contains(errorMsg, "权限"):
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": "加入组织失败: " + errorMsg,
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "成功加入组织",
	})
}

// SearchOrganizationUsersHandler 根据组织ID和用户名前缀搜索用户
// @Summary 搜索组织成员
// @Description 根据组织ID和用户名前缀搜索匹配的组织成员基本信息列表，如果前缀为空或者没有前缀，返回组织所有成员基本信息，仅组织成员可调用
// @Tags 组织管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param orgID path int true "组织ID"
// @Param name_prefix query string false "用户名前缀"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/organization/{orgID}/users/search [get]
func (h *OrganizationHandler) SearchOrganizationUsersHandler(c *gin.Context) {
	// 获取并验证路径参数
	orgIDStr := c.Param("orgID")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil || orgID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的组织ID格式",
		})
		return
	}

	// 获取查询参数
	namePrefix := c.Query("name_prefix")

	// 调用服务层
	userInfos, err := h.orgService.GetOrgUserInfosByNamePrefix(uint(orgID), namePrefix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "搜索用户失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "搜索成功",
		"data": gin.H{
			"users":       userInfos,
			"total_count": len(userInfos),
		},
	})
}

//组织活动相关接口

// CreateOrganizationActivityHandler 创建组织活动
// @Summary 创建组织活动
// @Description 在指定组织下创建新的活动
// @Tags 活动管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param orgID path int true "组织ID"
// @Param request body services.CreateActivityRequest true "创建活动请求"
// @Success 201 {object} SuccessResponse "创建成功"
// @Router /api/organization/{orgID}/activities [post]
func (h *OrganizationHandler) CreateOrganizationActivityHandler(c *gin.Context) {
	orgIDStr := c.Param("orgID")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil || orgID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的组织ID格式",
		})
		return
	}

	var req services.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数无效: " + err.Error(),
		})
		return
	}

	// 调用活动服务层
	err = h.activityService.CreateActivity(uint(orgID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "创建活动失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "活动创建成功",
	})
}

// CancelActivityRequest 定义取消活动的请求体结构
type CancelActivityRequest struct {
	Reason string `json:"reason" binding:"required,min=1,max=500" example:"因天气原因取消"`
}

// CancelActivityHandler 取消组织活动
// @Summary 取消组织活动
// @Description 将指定组织的活动状态标记为已取消，并记录取消原因
// @Tags 活动管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param orgID path int true "组织ID"
// @Param activityID path int true "活动ID"
// @Param request body CancelActivityRequest true "取消活动请求"
// @Success 200 {object} SuccessResponse "取消成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /api/organization/{orgID}/activities/{activityID}/cancel [patch]
func (h *OrganizationHandler) CancelActivityHandler(c *gin.Context) {
	// 获取并验证路径参数
	activityIDStr := c.Param("activityID")
	activityID, err := strconv.ParseUint(activityIDStr, 10, 32)
	if err != nil || activityID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的活动ID格式",
		})
		return
	}

	// 绑定并验证请求体
	var req CancelActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数无效: " + err.Error(),
		})
		return
	}

	// 调用活动服务层
	err = h.activityService.CancelActivity(uint(activityID), req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "取消活动失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "活动已成功取消",
	})
}

// UpdateActivityHandler 更新活动信息
// @Summary 更新活动信息
// @Description 更新指定活动的详细信息，支持部分更新（只更新提供的字段）
// @Tags 活动管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param orgID path int true "组织ID"
// @Param activityID path int true "活动ID"
// @Param request body services.UpdateActivityRequest true "更新活动请求"
// @Success 200 {object} SuccessResponse "更新成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /api/organization/{orgID}/activities/{activityID} [patch]
func (h *OrganizationHandler) UpdateActivityHandler(c *gin.Context) {
	// 获取并验证路径参数
	activityIDStr := c.Param("activityID")

	activityID, err := strconv.ParseUint(activityIDStr, 10, 32)
	if err != nil || activityID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的活动ID格式",
		})
		return
	}

	// 绑定并验证请求体
	var req services.UpdateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数无效: " + err.Error(),
		})
		return
	}

	// 调用服务层
	err = h.activityService.UpdateActivity(uint(activityID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "更新活动失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "活动更新成功",
	})
}

// GetOrgPublicActivitiesHandler 获取组织的公开活动列表
// @Summary 获取组织的公开活动列表
// @Description 获取指定组织下未结束、未取消的公开活动列表，任何已登录用户可访问
// @Tags 活动查询
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param orgID path int true "组织ID"
// @Success 200 {object} SuccessResponse "获取成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /api/organization/{orgID}/activities/public [get]
func (h *OrganizationHandler) GetOrgPublicActivitiesHandler(c *gin.Context) {
	//传入并且验证参数
	orgIDStr := c.Param("orgID")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil || orgID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的组织ID格式"})
		return
	}

	// 调用服务层，明确查询公开活动
	activities, err := h.activityService.GetOrgActivities(uint(orgID), "public")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取公开活动列表失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取公开活动列表成功",
		"data":    activities,
	})
}

// GetOrgInternalActivitiesHandler 获取组织的内部活动（仅成员可见）
// @Summary 获取组织的内部活动列表
// @Description 获取指定组织下未结束、未取消的仅组织成员可参与的活动列表，必须是组织成员才能访问
// @Tags 活动查询
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param orgID path int true "组织ID"
// @Success 200 {object} SuccessResponse "获取成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 403 {object} ErrorResponse "非组织成员，无权访问"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /api/organization/{orgID}/activities/internal [get]
func (h *OrganizationHandler) GetOrgInternalActivitiesHandler(c *gin.Context) {
	//传入并且验证参数
	orgIDStr := c.Param("orgID")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil || orgID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的组织ID格式"})
		return
	}

	// 调用服务层，明确查询非公开活动
	activities, err := h.activityService.GetOrgActivities(uint(orgID), "org_only")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取内部活动列表失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取内部活动列表成功",
		"data":    activities,
	})
}

// GetOrgAssignedActivitiesHandler 获取组织的指派活动（高权限可见）
// @Summary 获取组织的指派活动列表
// @Description 获取指定组织下未结束、未取消的仅由管理员指派参与的活动列表，需要组织管理员权限
// @Tags 活动管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param orgID path int true "组织ID"
// @Success 200 {object} SuccessResponse "获取成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 403 {object} ErrorResponse "非管理员，无权访问"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /api/organization/{orgID}/activities/assigned [get]
func (h *OrganizationHandler) GetOrgAssignedActivitiesHandler(c *gin.Context) {
	//传入并且验证参数
	orgIDStr := c.Param("orgID")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil || orgID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的组织ID格式"})
		return
	}

	// 调用服务层，明确查询高权限活动
	activities, err := h.activityService.GetOrgActivities(uint(orgID), "admin_assign")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取指派活动列表失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取指派活动列表成功",
		"data":    activities,
	})
}

// GetActivityParticipantsHandler 获取活动的参与者基本信息列表
// @Summary 获取活动的参与者基本信息列表
// @Description 根据活动ID，获取参与该活动的用户基本信息（包括姓名、账号、头像）
// @Tags 活动管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param orgID path int true "组织ID"
// @Param activityID path int true "活动ID"
// @Success 200 {object} SuccessResponse "获取成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /api/organization/activities/{activityID}/participants [get]
func (h *OrganizationHandler) GetActivityParticipantsHandler(c *gin.Context) {
	// 获取并验证路径参数
	activityIDStr := c.Param("activityID")
	activityID, err := strconv.ParseUint(activityIDStr, 10, 32)
	if err != nil || activityID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的活动ID格式",
		})
		return
	}

	// 调用服务层
	userInfos, err := h.activityService.GetActivityUsers(uint(activityID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取活动参与者失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取活动参与者成功",
		"data": gin.H{
			"participants": userInfos,
			"total_count":  len(userInfos),
		},
	})
}

// ParticipateActivityHandler 用户参加活动
// @Summary 用户参加活动/领取活动
// @Description 用户参加指定活动，系统会检查是否重复参加，也就是领取待办的意思
// @Tags 活动参与
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param orgID path int true "组织ID"
// @Param activityID path int true "活动ID"
// @Success 201 {object} SuccessResponse "参加成功"
// @Failure 400 {object} ErrorResponse "请求参数错误或用户已参加"
// @Failure 404 {object} ErrorResponse "活动不存在"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /api/organization/activities/{activityID}/participate [post]
func (h *OrganizationHandler) ParticipateActivityHandler(c *gin.Context) {
	// 获取并验证路径参数（活动ID）
	activityIDStr := c.Param("activityID")
	activityID, err := strconv.ParseUint(activityIDStr, 10, 32)
	if err != nil || activityID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的活动ID格式",
		})
		return
	}

	// 从JWT中间件设置的上下文中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "无法识别用户身份，请重新登录",
		})
		return
	}

	// 调用服务层
	err = h.activityService.ParticipateActivity(uint(activityID), userID.(uint))
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()

		// 根据服务层返回的错误信息细化HTTP状态码
		switch {
		case strings.Contains(errorMsg, "已经参加过"):
			statusCode = http.StatusBadRequest
		case strings.Contains(errorMsg, "活动不存在"):
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": "参加活动失败: " + errorMsg,
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "成功参加活动",
	})
}

// GetParticipationStatusRequest 定义请求体结构
type GetParticipationStatusRequest struct {
	UserIDs []uint `json:"user_ids" binding:"required,min=1" example:"[101,102,103]"`
}

// GetParticipationStatusHandler 批量获取用户参与状态
// @Summary 批量获取用户参与状态
// @Description 根据活动ID和用户ID列表，返回已参与该活动的用户ID列表，用于管理员分配活动的时候后端告诉它某些人已经参与，无法重复参与
// @Tags 活动管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param activityID path int true "活动ID"
// @Param request body GetParticipationStatusRequest true "用户ID列表"
// @Success 200 {object} SuccessResponse "查询成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /api/organization/activities/{activityID}/participation-status [post]
func (h *OrganizationHandler) GetParticipationStatusHandler(c *gin.Context) {
	// 获取并验证路径参数（活动ID）
	activityIDStr := c.Param("activityID")
	activityID, err := strconv.ParseUint(activityIDStr, 10, 32)
	if err != nil || activityID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的活动ID格式",
		})
		return
	}

	// 绑定并验证请求体
	var req GetParticipationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数无效: " + err.Error(),
		})
		return
	}

	// 调用服务层
	participatedIDs, err := h.activityService.GetParticipationStatus(uint(activityID), req.UserIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "查询参与状态失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "查询成功",
		"data": gin.H{
			"participated_user_ids": participatedIDs,
			"total_count":           len(participatedIDs),
		},
	})
}
