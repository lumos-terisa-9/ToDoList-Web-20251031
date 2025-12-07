// handlers/organization_handler.go
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
	orgService *services.OrganizationService
}

func NewOrganizationHandler(orgService *services.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		orgService: orgService,
	}
}

// SuccessResponse 成功响应结构体
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// SubmitCreateOrganizationApplicationHandler 提交创建组织申请
// @Summary 提交创建组织申请
// @Description 用户提交创建新组织的申请，需要管理员审批
// @Tags 组织申请
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body models.OrganizationApplication true "创建组织申请请求"
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

// SubmitJoinApplicationHandler 提交加入组织申请
// @Summary 提交加入组织申请
// @Description 用户提交加入现有组织的申请，需要组织管理员审批
// @Tags 组织申请
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body models.OrganizationApplication true "加入组织申请请求"
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

// ProcessApplicationHandler 处理组织申请
// @Summary 处理组织申请
// @Description 系统/组织管理员审批或拒绝组织创建/加入申请
// @Tags 组织申请
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param id path int true "申请ID"
// @Param request body object true "处理请求" example({"action":"approve","remark":"申请符合要求"})
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

// RemoveOrganizationMemberHandler 移除组织成员
// @Summary 移除组织成员
// @Description 组织管理员或创建者从组织中移除成员
// @Tags 组织管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body object true "移除请求" example({"org_id": 1, "user_id": 2})
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
// @Tags 组织管理
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
// @Tags 组织管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param id path int true "组织ID" minimum(1)
// @Success 200 {object} SuccessResponse "获取成功" example({"success": true, "message": "获取组织成功", "data": {"id": 1, "name": "技术部"}})
// @Failure 400 {object} ErrorResponse "请求参数错误" example({"success": false, "message": "无效的组织ID"})
// @Failure 404 {object} ErrorResponse "组织不存在"
// @Failure 500 {object} ErrorResponse "系统内部错误" example({"success": false, "message": "查询组织失败: 组织不存在"})
// @Router /api/organization/{id} [get]
func (h *OrganizationHandler) GetOrganizationByIDHandler(c *gin.Context) {
	// 获取并验证路径参数
	orgIDStr := c.Param("id")
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

// UpdateOrganizationNameHandler 更新组织名称
// @Summary 更新组织名称
// @Description 管理员手动更新指定组织的名称信息，example({"new_name": "新技术研发部"})
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param id path int true "组织ID"
// @Param request body object true "更新请求" example({"new_name": "新技术研发部"})
// @Success 200 {object} SuccessResponse "更新成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /admin/api/organization/{id}/name [patch]
func (h *OrganizationHandler) UpdateOrganizationNameHandler(c *gin.Context) {
	// 从路径参数获取组织ID
	orgIDStr := c.Param("id")
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

// UpdateOrganizationDescriptionHandler 更新组织描述
// @Summary 更新组织描述
// @Description 管理员手动更新指定组织的描述信息，example({"new_description": "这是更新后的组织描述"})
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param id path int true "组织ID"
// @Param request body object true "更新请求" example({"new_description": "这是更新后的组织描述"})
// @Success 200 {object} SuccessResponse "更新成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /admin/api/organization/{id}/description [patch]
func (h *OrganizationHandler) UpdateOrganizationDescriptionHandler(c *gin.Context) {
	orgIDStr := c.Param("id")
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

// UpdateOrganizationLogoHandler 更新组织Logo
// @Summary 更新组织Logo
// @Description 管理员手动更新指定组织的Logo图片URL，example({"new_logo_url": "https://example.com/logo.png"})
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param id path int true "组织ID"
// @Param request body object true "更新请求" example({"new_logo_url": "https://example.com/logo.png"})
// @Success 200 {object} SuccessResponse "更新成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /admin/api/organization/{id}/logo [patch]
func (h *OrganizationHandler) UpdateOrganizationLogoHandler(c *gin.Context) {
	orgIDStr := c.Param("id")
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

// TransferOrganizationOwnershipHandler 转移组织所有权
// @Summary 转移组织所有权
// @Description 组织者将指定组织的所有权转移给其他成员，无需管理员审批
// @Tags 组织管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param id path int true "组织ID"
// @Param request body object true "转移请求" example({"new_creator_id": 123})
// @Success 200 {object} SuccessResponse "转移成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "系统内部错误"
// @Router /api/organization/{id}/ownership [patch]
func (h *OrganizationHandler) TransferOrganizationOwnershipHandler(c *gin.Context) {
	orgIDStr := c.Param("id")
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
