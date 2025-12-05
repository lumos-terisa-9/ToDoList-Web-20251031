// handlers/organization_handler.go
package handlers

import (
	"net/http"
	"strconv"
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
// @Router /api/organization/applications/join-organization [post]
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
// @Router /api/organization/applications/{id}/process [patch]
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

// GetPendingCreateOrgApplicationsHandler 获取待处理的创建组织申请
// @Summary 获取待处理的创建组织申请
// @Description 管理员获取所有待处理的创建组织申请
// @Tags 组织申请
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Success 200 {object} SuccessResponse "获取成功" example({"success": true, "message": "获取创建组织申请成功"})
// @Failure 401 {object} ErrorResponse "用户未认证"
// @Failure 500 {object} ErrorResponse "系统内部错误" example({"success": false, "message": "获取申请失败"})
// @Router /api/orgnization/application/pending-create [get]
func (h *OrganizationHandler) GetPendingCreateOrgApplicationsHandler(c *gin.Context) {
	applications, err := h.orgService.GetPendingCreateOrgApplications()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "获取创建组织申请失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取创建组织申请成功",
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
// @Router /api/organization/application/delete:{id} [delete]
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
