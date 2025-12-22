package handlers

import (
	"net/http"
	"strconv"
	"team_task_hub/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// AIHandler AI推荐处理器
type AIHandler struct {
	recommendationService *services.RecommendationService
}

// NewAIHandler 创建处理器
func NewAIHandler(recommendationService *services.RecommendationService) *AIHandler {
	return &AIHandler{
		recommendationService: recommendationService,
	}
}

// RecommendOrganizations 推荐相似组织
// @Summary 智能推荐组织
// @Description 根据用户描述，智能推荐相似的组织
// @Tags AI
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param query query string true "用户查询描述" example("找一个技术社区")
// @Param limit query int false "返回数量限制" default(5)
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/ai/recommendations/organizations [get]
func (h *AIHandler) RecommendOrganizations(c *gin.Context) {
	// 获取查询参数
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "查询参数不能为空",
		})
		return
	}

	// 获取数量限制
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if err != nil || limit <= 0 {
		limit = 5
	}
	if limit > 20 {
		limit = 20
	}

	// 调用推荐服务
	recommendations, err := h.recommendationService.RecommendOrganizations(query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "推荐失败: " + err.Error(),
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "推荐成功",
		"data": gin.H{
			"query":           query,
			"recommendations": recommendations,
			"count":           len(recommendations),
		},
	})
}
