<template>
  <div v-if="visible" class="modal-overlay" @click.self="closeModal">
    <div class="modal-container">
      <!-- 关闭按钮 -->
      <button class="close-btn" @click="closeModal">×</button>

      <!-- 加载状态 -->
      <div v-if="loading" class="loading-container">
        <div class="loading-spinner"></div>
        <p>加载活动详情中...</p>
      </div>

      <!-- 错误状态 -->
      <div v-else-if="error" class="error-container">
        <div class="error-icon">⚠️</div>
        <p class="error-message">{{ error }}</p>
        <button class="retry-btn" @click="closeModal">关闭</button>
      </div>

      <!-- 活动详情内容 -->
      <div v-else-if="activity" class="activity-detail-modal">
        <!-- 顶部状态栏 -->
        <div class="modal-header">
          <div class="activity-status" :class="getStatusClass(activity)">
            {{ getStatusText(activity) }}
          </div>
          <h2 class="activity-title">{{ activity.title }}</h2>
          <div class="activity-id">ID: {{ activity.id }}</div>
        </div>

        <!-- 基本信息 -->
        <div class="info-section">
          <h3 class="section-title">
            <span class="section-icon">📋</span>
            基本信息
          </h3>
          <div class="info-grid">
            <div class="info-row">
              <span class="info-label">开始时间：</span>
              <span class="info-value">{{ formatTime(activity.start_time) }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">结束时间：</span>
              <span class="info-value">{{ formatTime(activity.end_time) }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">参与限制：</span>
              <span class="info-value">{{ getParticipationLimitText(activity.participation_limit) }}</span>
            </div>
          </div>
        </div>

        <!-- 活动描述 -->
        <div class="description-section">
          <h3 class="section-title">
            <span class="section-icon">📝</span>
            活动描述
          </h3>
          <div class="description-content">
            {{ activity.description || '暂无描述' }}
          </div>
        </div>

        <!-- 组织信息 -->
        <div v-if="activity.organization" class="organization-section">
          <h3 class="section-title">
            <span class="section-icon">🏢</span>
            组织信息
          </h3>
          <div class="organization-info">
            <div class="org-row">
              <span class="info-label">组织名称：</span>
              <span class="info-value">{{ activity.organization.name }}</span>
            </div>
            <div v-if="activity.organization.description" class="org-description">
              <span class="info-label">组织描述：</span>
              <div class="info-value">{{ activity.organization.description }}</div>
            </div>
          </div>
        </div>

        <!-- 额外信息（取消原因等） -->
        <div v-if="activity.status === 'cancelled' && activity.cancellation_reason" class="extra-section">
          <h3 class="section-title">
            <span class="section-icon">⚠️</span>
            取消原因
          </h3>
          <div class="cancellation-reason">
            {{ activity.cancellation_reason }}
          </div>
        </div>

        <!-- 评价区域 -->
        <div v-if="showReviewSection" class="review-section">
          <h3 class="section-title">
            <span class="section-icon">⭐</span>
            活动评价
          </h3>

          <!-- 已提交评价提示 -->
          <div v-if="hasSubmittedReview" class="review-submitted">
            <div class="submitted-icon">✅</div>
            <div class="submitted-text">您已提交评价，感谢您的反馈！</div>
          </div>

          <!-- 评价表单 -->
          <div v-else class="review-form">
            <div class="rating-input">
              <div class="rating-label">评分（0-10分）</div>
              <div class="rating-slider">
                <input
                  v-model="review.rating"
                  type="range"
                  min="0"
                  max="10"
                  step="1"
                  class="slider"
                  @input="updateRatingVisual"
                >
                <div class="rating-display">
                  <span class="rating-value">{{ review.rating }}</span>
                  <span class="rating-max">/ 10</span>
                </div>
              </div>
              <div class="rating-stars">
                <span
                  v-for="star in 10"
                  :key="star"
                  class="star"
                  :class="{ 'active': star <= review.rating }"
                  @click="review.rating = star"
                >★</span>
              </div>
            </div>

            <div class="review-text-input">
              <div class="review-label">
                评价内容<span class="required-mark">*</span>
              </div>
              <textarea
                v-model="review.review_text"
                class="review-textarea"
                placeholder="请分享您的活动体验..."
                rows="3"
                required
              ></textarea>
              <div v-if="!review.review_text.trim()" class="validation-error">
                评价内容不能为空
              </div>
            </div>

            <button
              class="submit-review-btn"
              @click="submitReview"
              :disabled="submittingReview || !review.review_text.trim()"
            >
              <span v-if="submittingReview">提交中...</span>
              <span v-else>提交评价</span>
            </button>
          </div>
        </div>

        <!-- 在评价区域上方或下方添加 -->
        <div v-if="reviewSuccess" class="review-success-toast">
          <div class="success-icon">✅</div>
          <div class="success-message">{{ reviewSuccessMessage }}</div>
        </div>

        <!-- 操作按钮 -->
        <div class="modal-actions">
          <button class="close-action-btn" @click="closeModal">关闭</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    required: true
  },
  activityData: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close', 'review-submitted', 'review-failed'])

// 状态
const loading = ref(false)
const error = ref(null)
const activity = ref(null)
const hasSubmittedReview = ref(false)
const submittingReview = ref(false)
const reviewSuccess = ref(false)
const reviewSuccessMessage = ref('')

// 评价数据
const review = reactive({
  rating: 5,
  review_text: ''
})

// 计算属性
const showReviewSection = computed(() => {
  return activity.value &&
    activity.value.status === 'active' &&
    !hasSubmittedReview.value
})

// 获取token
function getToken() {
  let token = localStorage.getItem('token')
  if (token && token.startsWith('{')) {
    try {
      const tokenData = JSON.parse(token)
      if (tokenData.data && tokenData.data.access_token) {
        token = tokenData.data.access_token
      } else if (tokenData.access_token) {
        token = tokenData.access_token
      } else if (tokenData.token) {
        token = tokenData.token
      }
    } catch (error) {
      console.error('解析token失败:', error)
      return null
    }
  }
  return token
}

// 检查评价状态
async function checkReviewStatus() {
  if (!activity.value) return

  const token = getToken()
  if (!token) return

  try {
    // 调用接口检查用户是否已提交评价
    const response = await fetch(
      `http://localhost:8080/api/organization/activities/${activity.value.id}/review-status`,
      {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      }
    )

    if (response.ok) {
      const data = await response.json()
      if (data.success) {
        hasSubmittedReview.value = data.has_reviewed || false
      }
    }
  } catch (err) {
    console.error('检查评价状态失败:', err)
  }
}

// 格式化时间
function formatTime(timeStr) {
  try {
    if (!timeStr || timeStr === '1900-01-01T00:00:00+08:00') {
      return '未设置'
    }
    const date = new Date(timeStr)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return '时间待定'
  }
}

// 获取参与限制文本
function getParticipationLimitText(limit) {
  switch (limit) {
    case 'public': return '公开活动'
    case 'org_only': return '组织内部'
    case 'admin_assign': return '专项指派'
    default: return '未知'
  }
}

// 获取状态文本
function getStatusText(activity) {
  switch (activity.status) {
    case 'active': return '进行中'
    case 'completed': return '已完成'
    case 'cancelled': return '已取消'
    default: return '未知'
  }
}

// 获取状态class
function getStatusClass(activity) {
  switch (activity.status) {
    case 'active': return 'status-active'
    case 'completed': return 'status-completed'
    case 'cancelled': return 'status-cancelled'
    default: return 'status-unknown'
  }
}

// 更新评分视觉
function updateRatingVisual(event) {
  review.rating = parseFloat(event.target.value)
}

// 提交评价
async function submitReview() {
  const token = getToken()
  if (!token) {
    emit('review-failed', '请先登录')
    return
  }

  if (!review.review_text || !review.review_text.trim()) {
    emit('review-failed', '评价内容不能为空')
    return
  }

  if (review.rating < 0 || review.rating > 10) {
    emit('review-failed', '评分必须在0-10之间')
    return
  }

  submittingReview.value = true

  try {
    const response = await fetch(
      `http://localhost:8080/api/organization/activities/${activity.value.id}/reviews`,
      {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          rating: review.rating,
          review_text: review.review_text.trim()
        })
      }
    )

    if (response.ok) {
      const data = await response.json()
      if (data.success) {
        hasSubmittedReview.value = true
        reviewSuccess.value = true
        reviewSuccessMessage.value = '评价提交成功！'

        // 3秒后自动关闭提示
        setTimeout(() => {
          reviewSuccess.value = false
        }, 3000)

        emit('review-submitted', review)
      }
    } else {
      const errorText = await response.text()
      let errorMessage = '提交评价失败'

      try {
        const errorData = JSON.parse(errorText)
        if (errorData.message && errorData.message.includes('参与')) {
          errorMessage = '请先参与此活动后再评价'
        } else {
          errorMessage = errorData.message || '提交评价失败'
        }
      } catch {
        if (errorText.includes('参与')) {
          errorMessage = '请先参与此活动后再评价'
        }
      }

      emit('review-failed', errorMessage)
    }
  } catch (err) {
    console.error('提交评价失败:', err)
    emit('review-failed', `提交评价失败: ${err.message}`)
  } finally {
    submittingReview.value = false
  }
}

// 加载活动详情
async function loadActivityDetail() {
  console.log('活动详情数据:', activity.value);
  console.log('活动状态:', activity.value?.status);
  console.log('是否是组织活动:', activity.value?.organization ? '是' : '否');
  loading.value = true
  error.value = null

  try {
    // 直接从props中获取活动数据
    if (props.activityData) {
      activity.value = props.activityData
      checkReviewStatus()
    } else {
      error.value = '未找到活动信息'
    }
  } catch (err) {
    console.error('加载活动详情失败:', err)
    error.value = err.message || '加载失败，请重试'
  } finally {
    loading.value = false
  }
}

// 关闭模态框
function closeModal() {
  emit('close')
}

// 监听visible变化
watch(() => props.visible, (newVal) => {
  if (newVal) {
    // 重置状态
    error.value = null
    review.rating = 5
    review.review_text = ''
    hasSubmittedReview.value = false

    // 加载活动详情
    loadActivityDetail()
  } else {
    // 关闭时重置数据
    activity.value = null
  }
})

// 监听activityData变化
watch(() => props.activityData, () => {
  if (props.visible) {
    loadActivityDetail()
  }
})
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(5px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 20px;
}

.modal-container {
  position: relative;
  width: 100%;
  max-width: 800px;
  max-height: 90vh;
  background: linear-gradient(135deg, #0f1419 0%, #1a2029 100%);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
}

.close-btn {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: white;
  font-size: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 10;
  transition: all 0.3s ease;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: rotate(90deg);
}

.loading-container,
.error-container {
  padding: 60px 40px;
  text-align: center;
  color: white;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(255, 255, 255, 0.1);
  border-top-color: #007bff;
  border-radius: 50%;
  margin: 0 auto 20px;
  animation: spin 1s linear infinite;
}

.error-icon {
  font-size: 48px;
  margin-bottom: 20px;
  opacity: 0.8;
}

.error-message {
  margin-bottom: 20px;
  color: rgba(255, 255, 255, 0.8);
}

.retry-btn {
  padding: 10px 24px;
  background: rgba(0, 123, 255, 0.2);
  border: 1px solid rgba(0, 123, 255, 0.3);
  color: #4dabf7;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.retry-btn:hover {
  background: rgba(0, 123, 255, 0.3);
}

.activity-detail-modal {
  padding: 40px;
  overflow-y: auto;
  max-height: calc(90vh - 80px);
  scrollbar-width: none;         /* Firefox */
  -ms-overflow-style: none;
}

.modal-header {
  margin-bottom: 32px;
  text-align: center;
}

.activity-status {
  display: inline-block;
  padding: 8px 16px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 12px;
}

.status-active {
  background: rgba(100, 200, 100, 0.15);
  border: 1px solid rgba(100, 200, 100, 0.3);
  color: #64c864;
}

.status-completed {
  background: rgba(120, 200, 255, 0.15);
  border: 1px solid rgba(120, 200, 255, 0.3);
  color: #78c8ff;
}

.status-cancelled {
  background: rgba(255, 100, 100, 0.15);
  border: 1px solid rgba(255, 100, 100, 0.3);
  color: #ff6464;
}

.activity-title {
  font-size: 28px;
  font-weight: 700;
  margin: 0 0 8px 0;
  color: white;
  line-height: 1.3;
}

.activity-id {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
  font-family: 'SF Mono', monospace;
}

.info-section,
.description-section,
.organization-section,
.extra-section,
.review-section {
  background: rgba(255, 255, 255, 0.03);
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 20px;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.section-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 18px;
  font-weight: 600;
  color: white;
  margin: 0 0 20px 0;
}

.section-icon {
  font-size: 20px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.info-row,
.org-row {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.info-row:last-child,
.org-row:last-child {
  border-bottom: none;
}

.info-label {
  color: rgba(255, 255, 255, 0.6);
  font-size: 14px;
}

.info-value {
  color: white;
  font-size: 14px;
  font-weight: 500;
}

.description-content {
  color: rgba(255, 255, 255, 0.9);
  line-height: 1.6;
  font-size: 15px;
  white-space: pre-line;
}

.org-description {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.cancellation-reason {
  color: rgba(255, 100, 100, 0.9);
  line-height: 1.6;
  font-size: 15px;
  padding: 12px;
  background: rgba(255, 100, 100, 0.1);
  border-radius: 8px;
  white-space: pre-line;
}

.review-submitted {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 20px;
  background: rgba(100, 200, 100, 0.1);
  border-radius: 12px;
  border: 1px solid rgba(100, 200, 100, 0.3);
}

.submitted-icon {
  font-size: 24px;
}

.submitted-text {
  font-size: 16px;
  font-weight: 600;
  color: #64c864;
}

.review-form {
  max-width: 600px;
  margin: 0 auto;
}

.rating-input {
  margin-bottom: 24px;
}

.rating-label {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
  margin-bottom: 12px;
  font-weight: 500;
}

.rating-slider {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 12px;
}

.slider {
  flex: 1;
  height: 6px;
  border-radius: 3px;
  background: rgba(255, 255, 255, 0.1);
  outline: none;
  -webkit-appearance: none;
}

.slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #c8a0ff;
  cursor: pointer;
  border: 2px solid white;
}

.rating-display {
  min-width: 60px;
  text-align: center;
  font-weight: 600;
}

.rating-value {
  font-size: 24px;
  color: #c8a0ff;
}

.rating-max {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
}

.rating-stars {
  display: flex;
  gap: 4px;
  justify-content: center;
}

.star {
  font-size: 20px;
  color: rgba(255, 255, 255, 0.3);
  cursor: pointer;
  transition: all 0.2s ease;
}

.star:hover {
  transform: scale(1.2);
  color: rgba(255, 215, 0, 0.8);
}

.star.active {
  color: #ffd700;
}

.review-text-input {
  margin-bottom: 24px;
}

.review-label {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
  margin-bottom: 8px;
  font-weight: 500;
}

.required-mark {
  color: #ff6464;
  margin-left: 4px;
}

.review-textarea {
  width: 100%;
  padding: 12px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: white;
  font-size: 14px;
  font-family: inherit;
  resize: vertical;
  transition: all 0.3s ease;
}

.review-textarea:focus {
  outline: none;
  border-color: rgba(200, 160, 255, 0.5);
  background: rgba(200, 160, 255, 0.05);
}

.validation-error {
  color: #ff6464;
  font-size: 12px;
  margin-top: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.validation-error::before {
  content: "⚠️";
  font-size: 12px;
}

.submit-review-btn {
  width: 100%;
  padding: 14px;
  border-radius: 12px;
  border: 1px solid rgba(100, 200, 100, 0.3);
  background: linear-gradient(135deg,
  rgba(100, 200, 100, 0.15),
  rgba(100, 200, 100, 0.08));
  color: #64c864;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.submit-review-btn:hover:not(:disabled) {
  background: linear-gradient(135deg,
  rgba(100, 200, 100, 0.25),
  rgba(100, 200, 100, 0.15));
  border-color: rgba(100, 200, 100, 0.5);
  transform: translateY(-2px);
}

.submit-review-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.modal-actions {
  margin-top: 32px;
  text-align: center;
}

.close-action-btn {
  padding: 12px 32px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  color: white;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.close-action-btn:hover {
  background: rgba(255, 255, 255, 0.15);
  transform: translateY(-2px);
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* 响应式调整 */
@media (max-width: 768px) {
  .modal-container {
    margin: 10px;
    max-height: 85vh;
  }

  .activity-detail-modal {
    padding: 24px;
    max-height: calc(85vh - 48px);
  }

  .activity-title {
    font-size: 24px;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }
}

/* 评价成功提示样式 */
.review-success-toast {
  position: fixed;
  top: 20px;
  right: 20px;
  background: rgba(100, 200, 100, 0.95);
  border: 1px solid #64c864;
  border-radius: 12px;
  padding: 16px 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  color: white;
  font-weight: 600;
  z-index: 100000;
  animation: slideIn 0.3s ease-out;
  box-shadow: 0 4px 20px rgba(100, 200, 100, 0.3);
}

.success-icon {
  font-size: 20px;
}

.success-message {
  font-size: 14px;
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}
</style>
