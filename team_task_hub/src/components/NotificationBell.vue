<template>
  <div class="notification-container" ref="containerRef">
    <!-- 通知铃铛 -->
    <div class="notification-bell" @click="toggleNotifications" :class="{ 'has-unread': totalNotificationCount > 0 }">
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M18 8C18 6.4087 17.3679 4.88258 16.2426 3.75736C15.1174 2.63214 13.5913 2 12 2C10.4087 2 8.88258 2.63214 7.75736 3.75736C6.63214 4.88258 6 6.4087 6 8C6 15 3 17 3 17H21C21 17 18 15 18 8Z"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M13.73 21C13.5542 21.3031 13.3018 21.5547 12.9982 21.7295C12.6946 21.9044 12.3504 21.9965 12 21.9965C11.6496 21.9965 11.3054 21.9044 11.0018 21.7295C10.6982 21.5547 10.4458 21.3031 10.27 21"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      <span v-if="totalNotificationCount > 0" class="notification-badge">
        {{ totalNotificationCount > 99 ? '99+' : totalNotificationCount }}
      </span>
    </div>

    <!-- 通知面板 -->
    <div v-if="showPanel" class="notification-panel" @mousedown.stop>
      <div class="panel-header">
        <h3>活动通知</h3>
        <button class="close-btn" @click.stop="toggleNotifications">×</button>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
        <p>加载通知中...</p>
      </div>

      <!-- 空状态 -->
      <div v-else-if="!hasNotifications" class="empty-state">
        <div class="empty-icon">🎉</div>
        <p>暂无新通知</p>
        <p class="empty-subtitle">所有活动通知都已处理</p>
      </div>

      <!-- 通知列表 -->
      <div v-else class="notification-list" @scroll.stop @wheel.stop="handleWheel">
        <!-- 已取消活动通知 -->
        <div v-if="cancelledActivities.length > 0" class="notification-section">
          <h4 class="section-title">
            已取消的活动
            <span class="section-count">{{ cancelledActivities.length }}</span>
          </h4>
          <div v-for="activity in cancelledActivities"
               :key="'cancelled-' + activity.id"
               class="notification-item cancelled"
               @click="openActivityModal(activity)">
            <div class="notification-icon">
              <span>❌</span>
            </div>
            <div class="notification-content">
              <div class="notification-title">活动已取消</div>
              <div class="notification-info">
                <span class="org-name">{{ activity.organization?.name || '未知组织' }}</span>
                <span class="activity-title">《{{ activity.title }}》</span>
              </div>
              <div class="notification-meta">
                <span class="notification-time">取消原因：{{ activity.cancellation_reason || '未说明原因' }}</span>
              </div>
              <div class="notification-actions">
                <button class="action-btn danger" @click.stop="deleteCancelled(activity.id)">
                  删除通知
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- 未读活动通知 -->
        <div v-if="unreadActivities.length > 0" class="notification-section">
          <h4 class="section-title">
            新活动通知
            <span class="section-count">{{ unreadActivities.length }}</span>
          </h4>
          <div v-for="activity in unreadActivities"
               :key="'unread-' + activity.id"
               class="notification-item unread"
               @click="openActivityModal(activity)">
            <div class="notification-icon">
              <span>🔔</span>
            </div>
            <div class="notification-content">
              <div class="notification-title">新活动通知</div>
              <div class="notification-info">
                <span class="org-name">{{ activity.organization?.name || '未知组织' }}</span>
                <span class="activity-title">《{{ activity.title }}》</span>
              </div>
              <div class="notification-meta">
                <span class="notification-time">
                  {{ formatActivityTime(activity.start_time) }} - {{ formatActivityTime(activity.end_time) }}
                </span>
                <span class="notification-status status-active">进行中</span>
              </div>
              <div class="notification-actions">
                <button class="action-btn mark-read" @click.stop="markAsRead(activity.id)">
                  标记已读
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 点击外部关闭 - 仅在通知面板外点击时触发 -->
    <div v-if="showPanel" class="notification-overlay" @click="handleOverlayClick"></div>

    <!-- 活动详情弹窗 -->
    <ActivityDetailModal
        v-model:visible="showActivityModal"
        :activity-data="selectedActivityData"
        @close="closeActivityModal"
        @review-submitted="handleReviewSubmitted"
        @review-failed="handleReviewFailed"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import ActivityDetailModal from './ActivityDetailModal.vue'

const API_BASE = 'http://localhost:8080/api'

// 响应式数据
const showPanel = ref(false)
const loading = ref(false)
const unreadActivities = ref([])
const cancelledActivities = ref([])
const containerRef = ref(null)
const showActivityModal = ref(false)
const selectedActivityData = ref(null)

// 计算属性 - 显示所有通知的总数
const totalNotificationCount = computed(() => {
  return unreadActivities.value.length + cancelledActivities.value.length
})

const hasNotifications = computed(() => totalNotificationCount.value > 0)

// 获取token
const getToken = () => {
  let token = localStorage.getItem('token')
  if (token && token.startsWith('{')) {
    try {
      const tokenData = JSON.parse(token)
      token = tokenData.data?.access_token || tokenData.access_token || tokenData.token || token
    } catch {
      // 保持原样
    }
  }
  return token
}

// 打开活动详情弹窗
function openActivityModal(activity) {
  selectedActivityData.value = activity
  showActivityModal.value = true
  showPanel.value = false // 关闭通知面板
}

// 关闭活动详情弹窗
function closeActivityModal() {
  showActivityModal.value = false
  selectedActivityData.value = null
}

// 处理评价提交成功
function handleReviewSubmitted(reviewData) {
  console.log('评价提交成功:', reviewData)
  // 可以在这里添加提示或其他逻辑
}

// 处理评价提交失败
function handleReviewFailed(errorMessage) {
  console.error('评价提交失败:', errorMessage)
  // 可以在这里添加错误提示
}

// 切换通知面板
const toggleNotifications = () => {
  showPanel.value = !showPanel.value
  if (showPanel.value) {
    loadNotifications()
    // 延迟绑定外部点击事件，避免立即触发
    nextTick(() => {
      setTimeout(() => {
        document.addEventListener('click', handleClickOutside)
      }, 10)
    })
  } else {
    document.removeEventListener('click', handleClickOutside)
  }
}

// 加载通知
const loadNotifications = async () => {
  const token = getToken()
  if (!token) return

  loading.value = true
  try {
    // 获取已取消活动
    const cancelledResponse = await fetch(`${API_BASE}/organization/me/activities/cancelled`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    // 获取未读活动
    const unreadResponse = await fetch(`${API_BASE}/organization/me/activities/unread`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    // 处理已取消活动响应
    if (cancelledResponse.ok) {
      const data = await cancelledResponse.json()
      console.log('已取消活动数据:', data)
      if (data.success && data.data && data.data.activities) {
        cancelledActivities.value = data.data.activities
      } else {
        cancelledActivities.value = []
      }
    } else {
      console.warn('获取已取消活动失败:', cancelledResponse.status)
      cancelledActivities.value = []
    }

    // 处理未读活动响应
    if (unreadResponse.ok) {
      const data = await unreadResponse.json()
      console.log('未读活动数据:', data)
      if (data.success && data.data && data.data.activities) {
        // 从未读列表中过滤掉已取消的活动，避免重复显示
        const cancelledIds = new Set(cancelledActivities.value.map(a => a.id))
        unreadActivities.value = data.data.activities.filter(activity =>
            activity.status === 'active' && !cancelledIds.has(activity.id)
        )
      } else {
        unreadActivities.value = []
      }
    } else {
      console.warn('获取未读活动失败:', unreadResponse.status)
      unreadActivities.value = []
    }

  } catch (error) {
    console.error('加载通知失败:', error)
    cancelledActivities.value = []
    unreadActivities.value = []
  } finally {
    loading.value = false
  }
}

// 标记为已读（仅用于未读活动通知）
const markAsRead = async (activityId) => {
  const token = getToken()
  if (!token) return

  try {
    const response = await fetch(`${API_BASE}/organization/me/activities/${activityId}/mark-as-read`, {
      method: 'PATCH',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    if (response.ok) {
      // 从未读列表中移除
      unreadActivities.value = unreadActivities.value.filter(activity => activity.id !== activityId)
    }
  } catch (error) {
    console.error('标记已读失败:', error)
  }
}

// 删除已取消活动通知
const deleteCancelled = async (activityId) => {
  const token = getToken()
  if (!token) return

  try {
    const response = await fetch(`${API_BASE}/organization/me/activities/cancelled/${activityId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    if (response.ok) {
      // 从已取消列表中移除
      cancelledActivities.value = cancelledActivities.value.filter(activity => activity.id !== activityId)
    }
  } catch (error) {
    console.error('删除已取消活动失败:', error)
  }
}

// 格式化时间
const formatActivityTime = (timeStr) => {
  try {
    const date = new Date(timeStr)
    return date.toLocaleString('zh-CN', {
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return '时间待定'
  }
}

// 处理鼠标滚轮事件 - 修复滚动问题
const handleWheel = (event) => {
  // 允许在通知列表内滚动
  event.stopPropagation()
  const listElement = event.currentTarget
  const maxScroll = listElement.scrollHeight - listElement.clientHeight

  // 如果已经到达顶部且还在向上滚动，或者到达底部且还在向下滚动，阻止默认行为
  if ((listElement.scrollTop <= 0 && event.deltaY < 0) ||
      (listElement.scrollTop >= maxScroll && event.deltaY > 0)) {
    event.preventDefault()
  }
}

// 处理遮罩层点击
const handleOverlayClick = (event) => {
  // 只有在点击遮罩层本身时才关闭
  if (event.target.className === 'notification-overlay') {
    toggleNotifications()
  }
}

// 处理外部点击
const handleClickOutside = (event) => {
  if (!containerRef.value || !containerRef.value.contains(event.target)) {
    showPanel.value = false
    document.removeEventListener('click', handleClickOutside)
  }
}

// 生命周期
onMounted(() => {
  // 每30秒检查一次新通知
  loadNotifications()
  const interval = setInterval(loadNotifications, 30000)

  onUnmounted(() => {
    clearInterval(interval)
    document.removeEventListener('click', handleClickOutside)
  })
})

// 暴露方法给父组件
defineExpose({
  refreshNotifications: loadNotifications
})
</script>

<style scoped>
.notification-container {
  position: relative;
  z-index: 1000;
  display: inline-block; /* 改为行内块，避免影响布局 */
}

/* 通知铃铛样式 */
.notification-bell {
  position: relative;
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  color: white;
}

.notification-bell:hover {
  background: rgba(255, 255, 255, 0.15);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.notification-bell.has-unread {
  animation: bell-shake 0.5s ease-in-out;
}

.notification-badge {
  position: absolute;
  top: -5px;
  right: -5px;
  background: linear-gradient(135deg, #ff6b6b, #ff4757);
  color: white;
  font-size: 12px;
  font-weight: 600;
  min-width: 20px;
  height: 20px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 6px;
  box-shadow: 0 2px 8px rgba(255, 107, 107, 0.4);
  animation: badge-pulse 2s infinite;
}

/* 通知面板样式 - 修复层级和交互 */
.notification-panel {
  position: absolute;
  top: 55px;
  right: 0;
  width: 420px;
  max-height: 500px;
  background: rgba(15, 20, 25, 0.98);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  overflow: hidden;
  animation: panel-slide 0.3s ease-out;
  z-index: 1001; /* 确保在遮罩层之上 */
  pointer-events: auto; /* 确保可以接收鼠标事件 */
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.panel-header h3 {
  margin: 0;
  color: white;
  font-size: 18px;
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.6);
  font-size: 24px;
  cursor: pointer;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  pointer-events: auto; /* 确保按钮可以点击 */
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

/* 加载状态 */
.loading-state {
  padding: 40px 20px;
  text-align: center;
  color: rgba(255, 255, 255, 0.6);
  pointer-events: none;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(255, 255, 255, 0.1);
  border-top-color: #007bff;
  border-radius: 50%;
  margin: 0 auto 16px;
  animation: spin 1s linear infinite;
}

/* 空状态 */
.empty-state {
  padding: 60px 20px;
  text-align: center;
  color: rgba(255, 255, 255, 0.6);
  pointer-events: none;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-subtitle {
  font-size: 14px;
  margin-top: 8px;
  color: rgba(255, 255, 255, 0.4);
}

/* 通知列表 - 隐藏滚动条但保留滚动功能 */
.notification-list {
  max-height: 380px;
  overflow-y: auto;
  padding: 0;
  /* 隐藏滚动条但保留滚动功能 */
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE and Edge */
  overscroll-behavior: contain; /* 防止滚动传播 */
  -webkit-overflow-scrolling: touch; /* 平滑滚动 */
  pointer-events: auto; /* 确保可以滚动 */
}

.notification-list::-webkit-scrollbar {
  display: none; /* Chrome, Safari, Opera */
}

/* 鼠标滚轮支持 */
.notification-list:hover {
  cursor: default;
}

.notification-section {
  padding: 0 20px 20px;
}

.section-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 20px 0 16px;
  color: white;
  font-size: 14px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.section-count {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.8);
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
}

/* 通知项 */
.notification-item {
  display: flex;
  gap: 12px;
  padding: 16px;
  margin-bottom: 12px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  border-left: 4px solid #007bff;
  transition: all 0.2s ease;
  pointer-events: auto; /* 确保可以点击 */
  cursor: pointer; /* 添加手形光标 */
}

.notification-item:hover {
  background: rgba(255, 255, 255, 0.08);
  transform: translateY(-2px);
}

.notification-item.unread {
  border-left-color: #007bff;
  background: rgba(0, 123, 255, 0.1);
}

.notification-item.unread:hover {
  background: rgba(0, 123, 255, 0.15);
}

.notification-item.cancelled {
  border-left-color: #ff6b6b;
  background: rgba(255, 107, 107, 0.1);
}

.notification-item.cancelled:hover {
  background: rgba(255, 107, 107, 0.15);
}

.notification-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 16px;
}

.notification-item.unread .notification-icon {
  background: rgba(0, 123, 255, 0.2);
}

.notification-item.cancelled .notification-icon {
  background: rgba(255, 107, 107, 0.2);
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  color: white;
  font-weight: 600;
  margin-bottom: 8px;
  line-height: 1.4;
}

.notification-info {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  font-size: 13px;
  flex-wrap: wrap;
}

.org-name {
  color: #4dabf7;
  font-weight: 500;
  background: rgba(77, 171, 247, 0.1);
  padding: 2px 8px;
  border-radius: 6px;
}

.activity-title {
  color: rgba(255, 255, 255, 0.9);
  font-weight: 500;
}

.notification-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
}

.notification-time {
  display: flex;
  align-items: center;
  gap: 4px;
}

.notification-status {
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
}

.status-active {
  background: rgba(100, 200, 100, 0.2);
  color: #64c864;
}

.notification-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  padding: 6px 12px;
  border-radius: 8px;
  border: none;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  pointer-events: auto; /* 确保按钮可以点击 */
}

.action-btn.mark-read {
  background: rgba(0, 123, 255, 0.2);
  color: #4dabf7;
  border: 1px solid rgba(0, 123, 255, 0.3);
}

.action-btn.mark-read:hover {
  background: rgba(0, 123, 255, 0.3);
}

.action-btn.danger {
  background: rgba(255, 107, 107, 0.2);
  color: #ff8787;
  border: 1px solid rgba(255, 107, 107, 0.3);
}

.action-btn.danger:hover {
  background: rgba(255, 107, 107, 0.3);
}

/* 遮罩层 - 修复层级 */
.notification-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1000; /* 在通知面板之下 */
  background: transparent; /* 透明，不遮挡内容 */
  cursor: default;
}

/* 动画 */
@keyframes bell-shake {
  0%, 100% { transform: rotate(0deg); }
  25% { transform: rotate(-15deg); }
  75% { transform: rotate(15deg); }
}

@keyframes badge-pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.1); }
}

@keyframes panel-slide {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 响应式设计 */
@media (max-width: 480px) {
  .notification-panel {
    width: 320px;
    right: -138px;
  }
}
</style>
