<template>
  <div class="activity-detail-common">
    <!-- 顶部导航 -->
    <header class="detail-header">
      <button class="back-btn" @click="$emit('close')">
        <span>‹</span>
        <span class="back-text">返回</span>
      </button>
      <h1 class="detail-title">活动详情</h1>
      <div class="header-badge" :class="getStatusClass(activity)">
        {{ getStatusText(activity) }}
      </div>
    </header>

    <!-- 活动基本信息 -->
    <section class="activity-info-section">
      <div class="section-header">
        <div class="section-icon">📋</div>
        <h2 class="section-title">活动基本信息</h2>
      </div>

      <div class="info-grid">
        <div class="info-item">
          <div class="info-label">活动标题</div>
          <div class="info-value">{{ activity.title }}</div>
        </div>

        <div class="info-item">
          <div class="info-label">活动ID</div>
          <div class="info-value">{{ activity.id }}</div>
        </div>

        <div class="info-item">
          <div class="info-label">活动类型</div>
          <div class="info-value">{{ getActivityTypeText(activity) }}</div>
        </div>

        <div class="info-item">
          <div class="info-label">开始时间</div>
          <div class="info-value">{{ formatTime(activity.start_time) }}</div>
        </div>

        <div class="info-item">
          <div class="info-label">结束时间</div>
          <div class="info-value">{{ formatTime(activity.end_time) }}</div>
        </div>

        <div class="info-item">
          <div class="info-label">参与限制</div>
          <div class="info-value">{{ getParticipationLimitText(activity.participation_limit) }}</div>
        </div>
      </div>

      <div class="info-item full-width">
        <div class="info-label">活动介绍</div>
        <div class="info-value description">{{ activity.description || '暂无介绍' }}</div>
      </div>
    </section>

    <!-- 参与者矩阵（展示用，无复选框） -->
    <section class="participants-section">
      <div class="section-header">
        <div class="section-icon">👥</div>
        <h2 class="section-title">活动参与者</h2>
        <div class="participants-count">
          {{ participants.length }}人参与
          <span v-if="completedUsers.length > 0" class="completed-count">
            · {{ completedUsers.length }}人已完成
          </span>
        </div>
      </div>

      <!-- 参与者网格 -->
      <div class="participants-grid" ref="participantsGrid">
        <div v-if="loadingParticipants" class="loading-participants">
          <div class="loading-spinner"></div>
          加载参与者中...
        </div>

        <div v-else-if="participants.length === 0" class="empty-participants">
          <div class="empty-icon">👤</div>
          <div class="empty-text">暂无参与者</div>
        </div>

        <div v-else class="participants-container">
          <div
              v-for="participant in participants"
              :key="participant.id"
              class="participant-card"
              :class="{ 'completed': isCompleted(participant.id) }"
          >
            <!-- 已完成标记 -->
            <div v-if="isCompleted(participant.id)" class="completed-badge">
              <span>✅</span>
            </div>

            <!-- 参与者头像 -->
            <div class="participant-avatar">
              <img
                  :src="ensureGitHubAvatarUrl(participant.avatar_url)"
                  :alt="participant.username"
                  @error="handleAvatarError"
              >
            </div>

            <!-- 参与者信息 -->
            <div class="participant-info">
              <div class="participant-name">{{ participant.username }}</div>
              <div class="participant-status">
                <span v-if="isCompleted(participant.id)" class="completed-text">已完成</span>
                <span v-else class="participating-text">参与中</span>
              </div>
              <div class="participant-id">ID: {{ participant.id }}</div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 评价区域 -->
    <section v-if="!hasSubmittedReview && activity.status === 'active'" class="review-section">
      <div class="section-header">
        <div class="section-icon">⭐</div>
        <h2 class="section-title">活动评价</h2>
      </div>

      <div class="review-form">
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
              rows="4"
              required
          ></textarea>
          <div v-if="!review.review_text.trim()" class="validation-error">
            评价内容不能为空
          </div>
        </div>

        <button class="submit-review-btn" @click="submitReview" :disabled="submittingReview || !review.review_text.trim()">
          <span v-if="submittingReview">提交中...</span>
          <span v-else>提交评价</span>
        </button>
      </div>
    </section>

    <!-- 已提交评价提示 -->
    <div v-else-if="hasSubmittedReview" class="review-submitted">
      <div class="submitted-icon">✅</div>
      <div class="submitted-text">您已提交评价，感谢您的反馈！</div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'

const props = defineProps({
  activity: {
    type: Object,
    required: true
  },
  orgId: {
    type: [String, Number],
    required: true
  },
  userRole: {
    type: String,
    default: 'Member'
  }
})

const emit = defineEmits(['review-submitted', 'review-failed', 'close'])

// GitHub配置
const GITHUB_CONFIG = {
}

// 状态
const participants = ref([])
const completedUsers = ref([])
const loadingParticipants = ref(false)
const hasSubmittedReview = ref(false)
const submittingReview = ref(false)

// 评价数据
const review = reactive({
  rating: 5,
  review_text: ''
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

  if (!token) {
    console.error('未找到认证令牌')
    return null
  }
  return token
}

// 获取默认头像URL
function getDefaultAvatarUrl() {
  return `https://${GITHUB_CONFIG.username}.github.io/${GITHUB_CONFIG.folder}/default-avatar.png`
}

// 确保头像URL使用GitHub URL
function ensureGitHubAvatarUrl(avatarUrl) {
  if (!avatarUrl) return getDefaultAvatarUrl()

  if (avatarUrl.includes('github.io') || avatarUrl.includes('githubusercontent.com')) {
    return avatarUrl
  }

  if (avatarUrl.startsWith('blob:') || avatarUrl.startsWith('data:') || !avatarUrl.startsWith('http')) {
    return getDefaultAvatarUrl()
  }

  return avatarUrl
}

// 头像加载失败处理
function handleAvatarError(event) {
  event.target.src = getDefaultAvatarUrl()
}

// 格式化时间
function formatTime(timeStr) {
  try {
    const date = new Date(timeStr)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
    // eslint-disable-next-line no-unused-vars
  } catch (error) {
    return '时间待定'
  }
}

// 获取活动类型文本
function getActivityTypeText(activity) {
  const title = activity.title || ''
  if (title.includes('比赛') || title.includes('友谊赛')) return '比赛'
  if (title.includes('训练') || title.includes('练习')) return '训练'
  if (title.includes('选拔')) return '选拔'
  if (title.includes('培训') || title.includes('指导')) return '培训'
  if (title.includes('聚会') || title.includes('交流')) return '交流'
  return '活动'
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

// 检查是否已完成
function isCompleted(userId) {
  return completedUsers.value.includes(parseInt(userId))
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

  // 验证评价内容为必填项
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
        `http://localhost:8080/api/organization/activities/${props.activity.id}/reviews`,
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
        emit('review-submitted', review)
      } else {
        // 根据错误消息判断是否为未参与活动
        let errorMsg = data.message || '提交评价失败'
        if (errorMsg.includes('参与')) {
          errorMsg = '请先参与此活动后再评价'
        }
        emit('review-failed', errorMsg)
      }
    } else {
      const errorText = await response.text()
      let errorMessage = '提交评价失败'

      // 尝试解析错误信息
      try {
        const errorData = JSON.parse(errorText)
        if (errorData.message && errorData.message.includes('参与')) {
          errorMessage = '请先参与此活动后再评价'
        } else {
          errorMessage = errorData.message || '提交评价失败'
        }
      } catch {
        // 如果无法解析JSON，使用原始错误信息
        if (errorText.includes('参与')) {
          errorMessage = '请先参与此活动后再评价'
        }
      }

      emit('review-failed', errorMessage)
    }
  } catch (error) {
    console.error('提交评价失败:', error)
    emit('review-failed', `提交评价失败: ${error.message}`)
  } finally {
    submittingReview.value = false
  }
}

// ActivityDetailCommon.vue
// 获取参与者列表
async function fetchParticipants() {
  const token = getToken()
  if (!token) {
    console.error('获取参与者失败：未找到token')
    return
  }

  loadingParticipants.value = true

  try {
    // 获取活动参与者
    const participantsResponse = await fetch(
        `http://localhost:8080/api/organization/activities/${props.activity.id}/participants`,
        {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        }
    )

    console.log('参与者响应状态：', participantsResponse.status)
    const responseText = await participantsResponse.text()
    console.log('参与者响应原始数据：', responseText)

    if (participantsResponse.ok) {
      const data = JSON.parse(responseText)
      console.log('参与者解析数据：', data)

      if (data.success && data.data) {
        // 根据返回格式获取参与者列表
        if (data.data.participants) {
          participants.value = data.data.participants
        } else if (Array.isArray(data.data)) {
          participants.value = data.data
        } else {
          participants.value = []
        }
        console.log('最终参与者列表：', participants.value)
      } else {
        console.warn('参与者数据格式不正确：', data)
        participants.value = []
      }
    } else {
      console.error('获取参与者请求失败：', participantsResponse.status)
      participants.value = []
    }

    // 如果参与者不为空，获取已完成用户
    if (participants.value.length > 0) {
      const userIds = participants.value.map(p => parseInt(p.id))
      console.log('获取已完成用户，用户ID列表：', userIds)

      const completedResponse = await fetch(
          `http://localhost:8080/api/organization/activities/${props.activity.id}/completed-users`,
          {
            method: 'POST',
            headers: {
              'Authorization': `Bearer ${token}`,
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              user_ids: userIds
            })
          }
      )

      console.log('已完成用户响应状态：', completedResponse.status)

      if (completedResponse.ok) {
        const completedData = await completedResponse.json()
        console.log('已完成用户数据：', completedData)

        if (completedData.success && completedData.data) {
          // ✅ 关键修改：只提取 completed_user_ids 数组
          completedUsers.value = completedData.data.completed_user_ids || []
          console.log('已完成的用户ID数组：', completedUsers.value)
        } else {
          console.warn('获取已完成用户数据失败：', completedData.message)
          completedUsers.value = []
        }
      } else {
        console.warn('获取已完成用户失败：', completedResponse.status)
        completedUsers.value = []
      }
    }
  } catch (error) {
    console.error('获取参与者数据失败:', error)
    participants.value = []
  } finally {
    loadingParticipants.value = false
  }
}

// 初始化
onMounted(() => {
  fetchParticipants()
})

// 暴露方法给父组件
defineExpose({
  fetchParticipants
})
</script>

<style scoped>
.activity-detail-common {
  min-height: 100vh;
  background: linear-gradient(135deg, #0f1419 0%, #1a2029 100%);
  color: white;
  padding: 20px;
  padding-bottom: 40px;
  border-radius: 24px;
}

/* 顶部导航 */
.detail-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 32px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 12px;
  color: white;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.12);
  transform: translateX(-2px);
}

.back-text {
  font-size: 14px;
}

.detail-title {
  flex: 1;
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  background: linear-gradient(135deg, #fff 0%, rgba(255, 255, 255, 0.8) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.header-badge {
  padding: 8px 16px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
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

/* 通用区块样式 */
.section-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
}

.section-icon {
  font-size: 24px;
}

.section-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

/* 活动信息 */
.activity-info-section {
  background: rgba(255, 255, 255, 0.03);
  border-radius: 20px;
  padding: 24px;
  margin-bottom: 24px;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.info-item {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  padding: 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.info-item.full-width {
  grid-column: 1 / -1;
}

.info-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
  margin-bottom: 8px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.info-value {
  font-size: 16px;
  font-weight: 600;
}

.info-value.description {
  font-size: 14px;
  line-height: 1.6;
  color: rgba(255, 255, 255, 0.9);
  font-weight: normal;
}

/* 参与者区域 */
.participants-section {
  background: rgba(255, 255, 255, 0.03);
  border-radius: 20px;
  padding: 24px;
  margin-bottom: 24px;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.participants-count {
  margin-left: auto;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
}

.completed-count {
  color: #64c864;
}

.participants-grid {
  min-height: 200px;
}

.loading-participants {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 60px;
  color: rgba(255, 255, 255, 0.6);
}

.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.1);
  border-top-color: rgba(200, 160, 255, 0.8);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.empty-participants {
  text-align: center;
  padding: 60px 20px;
  color: rgba(255, 255, 255, 0.4);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 12px;
  opacity: 0.3;
}

.empty-text {
  font-size: 16px;
}

/* 参与者容器 */
.participants-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 16px;
  max-height: 400px;
  overflow-y: auto;
  padding-right: 8px;
}

.participants-container::-webkit-scrollbar {
  width: 4px;
}

.participants-container::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 2px;
}

.participants-container::-webkit-scrollbar-thumb {
  background: rgba(200, 160, 255, 0.3);
  border-radius: 2px;
}

/* 参与者卡片 */
.participant-card {
  position: relative;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  transition: all 0.3s ease;
}

.participant-card:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

.participant-card.completed {
  opacity: 0.9;
  background: rgba(100, 200, 100, 0.08);
  border-color: rgba(100, 200, 100, 0.2);
}

/* 完成标记 */
.completed-badge {
  position: absolute;
  top: 12px;
  left: 12px;
  font-size: 16px;
  background: rgba(100, 200, 100, 0.2);
  border-radius: 50%;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 参与者头像 */
.participant-avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  overflow: hidden;
  margin-bottom: 12px;
  border: 3px solid rgba(255, 255, 255, 0.1);
}

.participant-card.completed .participant-avatar {
  border-color: rgba(100, 200, 100, 0.4);
}

.participant-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* 参与者信息 */
.participant-info {
  text-align: center;
  width: 100%;
}

.participant-name {
  font-weight: 600;
  font-size: 14px;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.participant-status {
  margin-bottom: 4px;
}

.completed-text {
  font-size: 12px;
  color: #64c864;
  padding: 2px 6px;
  background: rgba(100, 200, 100, 0.15);
  border-radius: 8px;
}

.participating-text {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
}

.participant-id {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  font-family: 'SF Mono', monospace;
}

/* 评价区域 */
.review-section {
  background: rgba(255, 255, 255, 0.03);
  border-radius: 20px;
  padding: 24px;
  margin-bottom: 24px;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.review-form {
  max-width: 600px;
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
}

.star {
  font-size: 24px;
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
  margin-bottom: 12px;
  font-weight: 500;
}

.review-textarea {
  width: 100%;
  padding: 16px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
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

.submit-review-btn {
  padding: 14px 32px;
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
  width: 100%;
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

/* 已提交评价提示 */
.review-submitted {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 32px;
  background: rgba(100, 200, 100, 0.1);
  border-radius: 20px;
  border: 1px solid rgba(100, 200, 100, 0.3);
  margin: 24px 0;
}

.submitted-icon {
  font-size: 24px;
}

.submitted-text {
  font-size: 16px;
  font-weight: 600;
  color: #64c864;
}

/* 动画 */
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* 响应式调整 */
@media (max-width: 768px) {
  .activity-detail-common {
    padding: 16px;
  }

  .detail-header {
    flex-wrap: wrap;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .participants-container {
    grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  }

  .rating-slider {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }
}

/* 必填标记 */
.required-mark {
  color: #ff6464;
  margin-left: 4px;
}

/* 验证错误提示 */
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
</style>
