package services

import (
	"fmt"
	"math"
	"sort"
	"team_task_hub/backend/internal/models"
	"team_task_hub/backend/internal/repositories"
)

// RecommendationService 推荐服务
type RecommendationService struct {
	orgRepo      *repositories.OrganizationRepository
	embedService *OllamaEmbeddingService
	cache        map[uint][]float64 
}

// RecommendationResult 推荐结果
type RecommendationResult struct {
	models.Organization
	Similarity  float64 `json:"similarity"`
	MatchReason string  `json:"match_reason"`
}

// NewRecommendationService 创建推荐服务
func NewRecommendationService(orgRepo *repositories.OrganizationRepository) *RecommendationService {
	embedService := NewOllamaEmbeddingService()
	return &RecommendationService{
		orgRepo:      orgRepo,
		embedService: embedService,
		cache:        make(map[uint][]float64),
	}
}

// Initialize 初始化服务，预计算所有组织的向量
func (s *RecommendationService) Initialize() error {
	fmt.Println("🚀 正在初始化推荐服务...")

	// 测试Ollama连接
	fmt.Println("🔌 测试Ollama连接...")
	if err := s.embedService.TestConnection(); err != nil {
		return fmt.Errorf("Ollama连接失败: %v", err)
	}
	fmt.Println("✅ Ollama连接成功")

	// 获取所有组织
	orgs, err := s.orgRepo.FindAllForRecommendation()
	if err != nil {
		return fmt.Errorf("获取组织数据失败: %v", err)
	}

	fmt.Printf("📊 找到 %d 个组织，开始向量化...\n", len(orgs))

	// 为每个组织生成向量
	successCount := 0
	for i, org := range orgs {
		if org.Description == "" {
			continue
		}

		embedding, err := s.embedService.GetEmbedding(org.Description)
		if err != nil {
			fmt.Printf("⚠️ 组织 %d (%s) 向量化失败: %v\n", org.ID, org.Name, err)
			continue
		}

		s.cache[org.ID] = embedding
		successCount++

		// 显示进度
		if (i+1)%10 == 0 || i == len(orgs)-1 {
			fmt.Printf("  进度: %d/%d\n", i+1, len(orgs))
		}
	}

	fmt.Printf("✅ 初始化完成，成功向量化 %d/%d 个组织\n", successCount, len(orgs))
	return nil
}

// cosineSimilarity 计算余弦相似度
func (s *RecommendationService) cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0.0
	}

	var dot, normA, normB float64
	for i := 0; i < len(a); i++ {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0.0
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

// RecommendOrganizations 推荐相似组织
func (s *RecommendationService) RecommendOrganizations(query string, limit int) ([]RecommendationResult, error) {
	// 将用户查询向量化
	queryVector, err := s.embedService.GetEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("向量化查询失败: %v", err)
	}

	// 获取所有组织
	allOrgs, err := s.orgRepo.FindAllForRecommendation()
	if err != nil {
		return nil, fmt.Errorf("获取组织列表失败: %v", err)
	}

	// 计算相似度
	var results []RecommendationResult

	for _, org := range allOrgs {
		// 从缓存获取组织向量
		orgVector, exists := s.cache[org.ID]
		if !exists || org.Description == "" {
			continue
		}

		// 计算相似度
		similarity := s.cosineSimilarity(queryVector, orgVector)

		// 只保留相似度较高的结果
		if similarity > 0.6 { // 阈值可根据实际效果调整
			results = append(results, RecommendationResult{
				Organization: org,
				Similarity:   similarity,
				MatchReason:  s.getMatchReason(similarity, org.Description, query),
			})
		}
	}

	// 按相似度降序排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	// 限制返回数量
	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// getMatchReason 生成匹配原因
func (s *RecommendationService) getMatchReason(similarity float64, description, query string) string {
	if similarity > 0.8 {
		return "描述高度相关"
	} else if similarity > 0.6 {
		return "描述非常相似"
	} else if similarity > 0.4 {
		return "描述有一定关联"
	} else {
		return "描述部分相关"
	}
}
