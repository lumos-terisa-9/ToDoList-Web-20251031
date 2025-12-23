<template>
  <div class="activity-page">
    <!-- 活动列表区域 -->
    <main class="activity-main">
      <button class="back-btn" @click="goBack">
        <span>‹</span>
        <span class="back-text">返回</span>
      </button>
      <!-- 加载状态 -->
      <div v-if="loading" class="loading-overlay">
        <div class="loading-crystal">
          <div class="crystal-inner"></div>
          <div class="crystal-glow"></div>
        </div>
        <div class="loading-text">加载活动中...</div>
      </div>

      <!-- 三栏活动列表 -->
      <div v-else class="activity-columns">
        <!-- 公开活动 -->
        <section class="activity-column">
          <div class="column-header">
            <div class="column-icon" style="background: rgba(120, 200, 255, 0.15);">
              <span>⚡</span>
            </div>
            <div class="column-title">
              <h2>组织公开活动</h2>
              <div class="column-subtitle">全员可参与 · {{ publicActivities.length }}个</div>
            </div>
          </div>

          <div class="activity-list" v-if="publicActivities.length > 0">
            <div
                v-for="(activity, index) in publicActivities"
                :key="activity.id"
                class="activity-card"
                :class="{
                  'full': activity.status === 'cancelled',
                  'joined': index < publicParticipatedCount
                }"
                @click="viewActivityDetail(activity, 'public')"
            >
              <div class="activity-badge" v-if="activity.status === 'active' && activity.currentParticipants && activity.currentParticipants >= activity.maxParticipants * 0.8">
                热门
              </div>

              <div class="activity-header">
                <div class="activity-type" :class="getActivityTypeClass(activity)">{{ getActivityTypeText(activity) }}</div>
                <div class="activity-time">{{ formatTime(activity.start_time) }}</div>
              </div>

              <div class="activity-title">
                {{ activity.title }}
                <span v-if="index < publicParticipatedCount" class="participated-tag">已参与</span>
              </div>

              <div class="activity-meta">
                <div class="meta-item">
                  <span class="meta-icon">📍</span>
                  {{ activity.location || '未指定地点' }}
                </div>
                <div class="meta-item">
                  <span class="meta-icon">👥</span>
                  <span v-if="activity.currentParticipants && activity.maxParticipants">
                    {{ activity.currentParticipants }}/{{ activity.maxParticipants }}人
                  </span>
                  <span v-else>人数不限</span>
                </div>
              </div>

              <div class="activity-footer">
                <div class="activity-status">
                  <span :class="getStatusClass(activity)">
                    {{ getStatusText(activity) }}
                  </span>
                </div>

                <button
                    class="join-btn"
                    :class="getJoinBtnClass(activity, index)"
                    @click.stop="handleJoin(activity)"
                    :disabled="!canJoinActivity(activity, index)"
                >
                  <span>{{ getJoinBtnText(activity, index) }}</span>

                  <div class="btn-sparkle" v-if="canJoinActivity(activity, index) && index >= publicParticipatedCount">
                    <div class="spark"></div>
                    <div class="spark"></div>
                    <div class="spark"></div>
                  </div>
                </button>
              </div>
            </div>
          </div>

          <div v-else class="empty-state">
            <div class="empty-icon">📅</div>
            <div class="empty-text">暂无公开活动</div>
            <div class="empty-hint">等待组织发布新的公开活动</div>
          </div>
        </section>

        <!-- 内部活动 -->
        <section class="activity-column">
          <div class="column-header">
            <div class="column-icon" style="background: rgba(255, 200, 100, 0.15);">
              <span>🛡️</span>
            </div>
            <div class="column-title">
              <h2>组织内部活动</h2>
              <div class="column-subtitle">仅限成员 · {{ internalActivities.length }}个</div>
            </div>
          </div>

          <div class="activity-list" v-if="internalActivities.length > 0">
            <div
                v-for="(activity, index) in internalActivities"
                :key="activity.id"
                class="activity-card internal"
                :class="{
                  'full': activity.status === 'cancelled',
                  'joined': index < internalParticipatedCount
                }"
                @click="viewActivityDetail(activity, 'internal')"
            >

              <div class="activity-header">
                <div class="activity-type internal">{{ getActivityTypeText(activity) }}</div>
                <div class="activity-time">{{ formatTime(activity.start_time) }}</div>
              </div>

              <div class="activity-title">
                {{ activity.title }}
                <span v-if="index < internalParticipatedCount" class="participated-tag">已参与</span>
              </div>

              <div class="activity-meta">
                <div class="meta-item">
                  <span class="meta-icon">🔒</span>
                  仅内部成员
                </div>
                <div class="meta-item">
                  <span class="meta-icon">👥</span>
                  <span v-if="activity.currentParticipants && activity.maxParticipants">
                    {{ activity.currentParticipants }}/{{ activity.maxParticipants }}人
                  </span>
                  <span v-else>人数不限</span>
                </div>
              </div>

              <div class="activity-footer">
                <div class="activity-status">
                  <span :class="getStatusClass(activity)">
                    {{ getStatusText(activity) }}
                  </span>
                </div>

                <button
                    class="join-btn internal"
                    :class="getJoinBtnClass(activity, index)"
                    @click.stop="handleJoin(activity)"
                    :disabled="!canJoinActivity(activity, index)"
                >
                  <span>{{ getJoinBtnText(activity, index) }}</span>

                  <div class="btn-sparkle" v-if="canJoinActivity(activity, index) && index >= internalParticipatedCount">
                    <div class="spark"></div>
                    <div class="spark"></div>
                    <div class="spark"></div>
                  </div>
                </button>
              </div>
            </div>
          </div>

          <div v-else class="empty-state">
            <div class="empty-icon">🏠</div>
            <div class="empty-text">暂无内部活动</div>
            <div class="empty-hint">等待组织发布内部活动</div>
          </div>
        </section>

        <!-- 专项活动 -->
        <section class="activity-column">
          <div class="column-header">
            <div class="column-icon" style="background: rgba(200, 160, 255, 0.15);">
              <span>🎯</span>
            </div>
            <div class="column-title">
              <h2>专项活动 & 已满活动</h2>
              <div class="column-subtitle">特殊安排 · {{ assignedActivities.length }}个</div>
            </div>
          </div>

          <div class="activity-list" v-if="assignedActivities.length > 0">
            <div
                v-for="(activity, index) in assignedActivities"
                :key="activity.id"
                class="activity-card assigned"
                :class="{
                  'full': activity.status === 'cancelled' || (activity.currentParticipants && activity.maxParticipants && activity.currentParticipants >= activity.maxParticipants),
                  'joined': index < assignedParticipatedCount
                }"
                @click="viewActivityDetail(activity, 'assigned')"
            >

              <div class="activity-header">
                <div class="activity-type assigned">{{ getActivityTypeText(activity) }}</div>
                <div class="activity-time">{{ formatTime(activity.start_time) }}</div>
              </div>

              <div class="activity-title">
                {{ activity.title }}
                <span v-if="index < assignedParticipatedCount" class="participated-tag">已参与</span>
              </div>

              <div class="activity-meta">
                <div class="meta-item">
                  <span class="meta-icon">🎖️</span>
                  {{ getDifficultyText(activity) }}
                </div>
                <div class="meta-item">
                  <span class="meta-icon">👥</span>
                  <span v-if="activity.currentParticipants && activity.maxParticipants">
                    {{ activity.currentParticipants }}/{{ activity.maxParticipants }}人
                    <span v-if="activity.currentParticipants >= activity.maxParticipants" class="full-badge">已满</span>
                  </span>
                  <span v-else>人数不限</span>
                </div>
              </div>

              <div class="activity-footer">
                <div class="activity-status">
                  <span :class="getStatusClass(activity)">
                    {{ getStatusText(activity) }}
                  </span>
                </div>

                <button
                    class="join-btn assigned"
                    :class="getJoinBtnClass(activity, index)"
                    @click.stop="handleJoin(activity)"
                    :disabled="!canJoinActivity(activity, index)"
                >
                  <span>{{ getJoinBtnText(activity, index) }}</span>
                </button>
              </div>
            </div>
          </div>

          <div v-else class="empty-state">
            <div class="empty-icon">🎖️</div>
            <div class="empty-text">暂无专项活动</div>
            <div class="empty-hint">等待组织发布专项活动</div>
          </div>
        </section>
      </div>
    </main>

    <!-- 浮动刷新按钮 -->
    <button class="fab-refresh" @click="refreshActivities" :disabled="loading">
      <span v-if="loading">⟳</span>
      <span v-else>↻</span>
    </button>

    <!-- 新建活动按钮（仅管理员和创建者可见） -->
    <button
        v-if="isAdminOrCreator"
        class="fab-create"
        @click="showCreateDialog = true"
        :disabled="loading"
    >
      <span>+</span>
      <div class="create-tooltip">新建活动</div>
    </button>

    <!-- 创建活动弹窗 -->
    <CreateActivityDialog
        v-if="showCreateDialog"
        :org-id="orgId"
        :org-name="orgName"
        @close="showCreateDialog = false"
        @activity-created="handleActivityCreated"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import CreateActivityDialog from '@/components/CreateActivityDialog.vue'

const route = useRoute()
const router = useRouter()

// 基础数据
const API_BASE = 'http://localhost:8080/api'
const orgId = route.params.id
const currentTime = ref('')
const showCreateDialog = ref(false)
const isAdminOrCreator = ref(false)

// 加载状态
const loading = ref(true)

// 活动数据
const publicActivities = ref([])
const internalActivities = ref([])
const assignedActivities = ref([])

// 用户已参与的活动数量
const publicParticipatedCount = ref(0)
const internalParticipatedCount = ref(0)
const assignedParticipatedCount = ref(0)

// 缓存所有活动信息
const allActivitiesCache = ref({
  public: [],
  internal: [],
  assigned: []
})

// 从localStorage获取用户信息
const getCurrentUser = () => {
  try {
    const userData = localStorage.getItem('currentUser')
    if (userData) {
      const parsedData = JSON.parse(userData)

      // 根据提供的响应格式，从data字段中提取用户信息
      if (parsedData.success && parsedData.data) {
        return parsedData.data // 直接返回data字段
      }

      // 如果直接就是用户对象格式，则直接返回
      if (parsedData.id || parsedData.user_id) {
        return parsedData
      }
    }
  } catch (error) {
    console.error('获取用户信息失败:', error)
  }
  return null
}

// 获取用户ID
const getUserIdFromStorage = () => {
  const user = getCurrentUser()
  return user ? user.id : null
}

// 检查用户权限
const checkUserRole = async () => {
  const token = getToken()
  if (!token) {
    console.error('未找到认证令牌')
    return
  }

  const userId = getUserIdFromStorage()
  if (!userId) {
    console.error('未找到用户ID')
    return
  }

  try {
    const response = await fetch(`${API_BASE}/organization/${orgId}/users/${userId}/role`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    if (response.ok) {
      const data = await response.json()
      if (data.success && data.data) {
        const role = data.data.role
        // 如果是Admin或Creator，显示新建活动按钮
        if (role === 'Admin' || role === 'Creator') {
          isAdminOrCreator.value = true
          console.log(`用户角色为${role}，显示新建活动按钮`)
        } else {
          console.log(`用户角色为${role}，无权新建活动`)
        }
      }
    } else {
      console.error('获取用户角色失败:', await response.text())
    }
  } catch (error) {
    console.error('检查用户权限失败:', error)
  }
}

// 活动创建成功后的处理
const handleActivityCreated = () => {
  showCreateDialog.value = false
  // 刷新活动列表
  fetchActivities()
}

// 获取token
const getToken = () => {
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

// 获取活动数据
const fetchActivities = async () => {
  const token = getToken()
  if (!token) {
    alert('请先登录')
    router.push('/login')
    return
  }

  loading.value = true

  try {
    // 并发获取三个接口的数据
    const [publicRes, internalRes, assignedRes] = await Promise.all([
      fetch(`${API_BASE}/organization/${orgId}/activities/public`, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      }),
      fetch(`${API_BASE}/organization/${orgId}/activities/internal`, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      }),
      fetch(`${API_BASE}/organization/${orgId}/activities/assigned`, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      })
    ])

    // 处理公开活动响应
    if (publicRes.ok) {
      const data = await publicRes.json()
      if (data.success && data.data) {
        publicActivities.value = data.data
        publicParticipatedCount.value = data.paticipatedCount || 0
        allActivitiesCache.value.public = data.data // 缓存数据
      } else {
        publicActivities.value = []
        publicParticipatedCount.value = 0
        allActivitiesCache.value.public = []
      }
    } else {
      publicActivities.value = []
      publicParticipatedCount.value = 0
      allActivitiesCache.value.public = []
    }

    // 处理内部活动响应
    if (internalRes.ok) {
      const data = await internalRes.json()
      if (data.success && data.data) {
        internalActivities.value = data.data
        internalParticipatedCount.value = data.paticipatedCount || 0
        allActivitiesCache.value.internal = data.data // 缓存数据
      } else {
        internalActivities.value = []
        internalParticipatedCount.value = 0
        allActivitiesCache.value.internal = []
      }
    } else {
      internalActivities.value = []
      internalParticipatedCount.value = 0
      allActivitiesCache.value.internal = []
    }

    // 处理专项活动响应
    if (assignedRes.ok) {
      const data = await assignedRes.json()
      if (data.success && data.data) {
        assignedActivities.value = data.data
        assignedParticipatedCount.value = data.paticipatedCount || 0
        allActivitiesCache.value.assigned = data.data // 缓存数据
      } else {
        assignedActivities.value = []
        assignedParticipatedCount.value = 0
        allActivitiesCache.value.assigned = []
      }
    } else {
      assignedActivities.value = []
      assignedParticipatedCount.value = 0
      allActivitiesCache.value.assigned = []
    }

    // 将缓存数据存储到sessionStorage，供ActivityDetail使用
    sessionStorage.setItem('activityCache', JSON.stringify({
      orgId: orgId,
      timestamp: Date.now(),
      data: allActivitiesCache.value
    }))

  } catch (error) {
    console.error('获取活动数据失败:', error)
    // 如果出错，清空所有数据
    publicActivities.value = []
    internalActivities.value = []
    assignedActivities.value = []
    publicParticipatedCount.value = 0
    internalParticipatedCount.value = 0
    assignedParticipatedCount.value = 0
  } finally {
    loading.value = false
  }
}

// 参加活动
const handleJoin = async (activity) => {
  const token = getToken()
  if (!token) {
    alert('请先登录')
    router.push('/login')
    return
  }

  // 检查是否是已参与的活动
  const isPublic = publicActivities.value.includes(activity)
  const isInternal = internalActivities.value.includes(activity)
  const isAssigned = assignedActivities.value.includes(activity)

  let index = -1
  if (isPublic) index = publicActivities.value.indexOf(activity)
  else if (isInternal) index = internalActivities.value.indexOf(activity)
  else if (isAssigned) index = assignedActivities.value.indexOf(activity)

  if ((isPublic && index < publicParticipatedCount.value) ||
      (isInternal && index < internalParticipatedCount.value) ||
      (isAssigned && index < assignedParticipatedCount.value)) {
    console.log('已参与该活动')
    return
  }

  if (!canJoinActivity(activity, index)) {
    return
  }

  try {
    const response = await fetch(`${API_BASE}/organization/activities/${activity.id}/participate`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        orgId: orgId,
        activityId: activity.id
      })
    })

    const data = await response.json()

    if (response.ok && data.success) {
      // 刷新活动列表
      fetchActivities()
    } else {
      console.log('参加活动失败:', data.message)
    }
  } catch (error) {
    console.error('参加活动失败:', error)
  }
}

// 刷新活动
const refreshActivities = () => {
  fetchActivities()
}

// 格式化时间为本地显示格式
const formatTime = (timeStr) => {
  try {
    const date = new Date(timeStr)
    return date.toLocaleString('zh-CN', {
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
const getActivityTypeText = (activity) => {
  const title = activity.title || ''
  if (title.includes('比赛') || title.includes('友谊赛')) return '比赛'
  if (title.includes('训练') || title.includes('练习')) return '训练'
  if (title.includes('选拔')) return '选拔'
  if (title.includes('培训') || title.includes('指导')) return '培训'
  if (title.includes('聚会') || title.includes('交流')) return '交流'
  return '活动'
}

// 获取活动类型class
const getActivityTypeClass = (activity) => {
  const type = getActivityTypeText(activity)
  switch (type) {
    case '比赛': return 'competition'
    case '训练': return 'training'
    case '选拔': return 'selection'
    case '培训': return 'training'
    case '交流': return 'communication'
    default: return 'default'
  }
}

// 获取状态文本
const getStatusText = (activity) => {
  if (!activity) return '未知状态'
  switch (activity.status) {
    case 'active': return '进行中'
    case 'completed': return '已完成'
    case 'cancelled': return '已取消'
    default: return '未知状态'
  }
}

// 获取状态class
const getStatusClass = (activity) => {
  if (!activity) return 'status-unknown'
  switch (activity.status) {
    case 'active': return 'status-active'
    case 'completed': return 'status-completed'
    case 'cancelled': return 'status-cancelled'
    default: return 'status-unknown'
  }
}

// 获取难度文本
const getDifficultyText = (activity) => {
  if (!activity) return '中级'
  if (activity.title && activity.title.includes('高级')) return '高级'
  if (activity.title && activity.title.includes('中级')) return '中级'
  if (activity.title && activity.title.includes('基础') || activity.title.includes('入门')) return '初级'
  return '中级'
}

// 获取参加按钮文本
const getJoinBtnText = (activity, index) => {
  if (!activity) return '未知'

  // 检查是否已参与
  const isPublic = publicActivities.value.includes(activity)
  const isInternal = internalActivities.value.includes(activity)
  const isAssigned = assignedActivities.value.includes(activity)

  if ((isPublic && index < publicParticipatedCount.value) ||
      (isInternal && index < internalParticipatedCount.value) ||
      (isAssigned && index < assignedParticipatedCount.value)) {
    return '已参与'
  }

  if (activity.status === 'cancelled') return '已取消'
  if (activity.status === 'completed') return '已结束'
  if (activity.currentParticipants && activity.maxParticipants && activity.currentParticipants >= activity.maxParticipants) {
    return '已满'
  }
  if (activity.participation_limit === 'admin_assign') return '申请参加'
  return '参加'
}

// 获取参加按钮class
const getJoinBtnClass = (activity, index) => {
  if (!activity) return { 'disabled': true }

  // 检查是否已参与
  const isPublic = publicActivities.value.includes(activity)
  const isInternal = internalActivities.value.includes(activity)
  const isAssigned = assignedActivities.value.includes(activity)

  const isParticipated = (isPublic && index < publicParticipatedCount.value) ||
      (isInternal && index < internalParticipatedCount.value) ||
      (isAssigned && index < assignedParticipatedCount.value)

  const baseClass = {
    'joined': isParticipated || activity.status === 'completed',
    'disabled': !canJoinActivity(activity, index) || isParticipated
  }

  if (activity.participation_limit === 'org_only') {
    baseClass.internal = true
  } else if (activity.participation_limit === 'assigned') {
    baseClass.assigned = true
  }

  return baseClass
}

// 是否可以参加活动
const canJoinActivity = (activity, index) => {
  if (!activity) return false

  // 检查是否已参与
  const isPublic = publicActivities.value.includes(activity)
  const isInternal = internalActivities.value.includes(activity)
  const isAssigned = assignedActivities.value.includes(activity)

  if ((isPublic && index < publicParticipatedCount.value) ||
      (isInternal && index < internalParticipatedCount.value) ||
      (isAssigned && index < assignedParticipatedCount.value)) {
    return false
  }

  if (activity.status !== 'active') return false
  if (activity.currentParticipants && activity.maxParticipants && activity.currentParticipants >= activity.maxParticipants) {
    return false
  }
  return true
}

// 查看活动详情
const viewActivityDetail = (activity, type) => {
  console.log('跳转到活动详情:', {
    activityId: activity.id,
    orgId: orgId,
    activity: activity,
    type: type
  })

  router.push({
    name: 'ActivityDetail',
    params: {
      orgId: orgId,
      activityId: activity.id,
      fromType: type // 添加类型信息，帮助在ActivityDetail中定位活动
    },
    query: {
      timestamp: Date.now() // 添加时间戳避免缓存问题
    }
  })
}

// 返回上一页
const goBack = () => {
  router.push({
    name: 'OrgPage',
    params: { id: orgId }
  })
}

// 更新时间显示
const updateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    weekday: 'short'
  })
}

// 生命周期
onMounted(() => {
  updateTime()
  fetchActivities()
  checkUserRole()  // 检查用户权限

  // 更新时间
  const timer = setInterval(updateTime, 60000)

  onUnmounted(() => {
    clearInterval(timer)
  })
})
</script>

<style scoped>
.activity-page {
  min-height: 100vh;
  background: var(--bg, #0f1419);
  color: var(--text, #fff);
  position: relative;
  overflow-x: hidden;
  font-family: ui-sans-serif, system-ui, -apple-system, "PingFang SC", "Microsoft YaHei",
  "Helvetica Neue", Arial, "Noto Sans", "Apple Color Emoji", "Segoe UI Emoji";
}

/* 奇幻背景 */
.activity-page::before {
  content: "";
  position: fixed;
  inset: -200px;
  background:
      radial-gradient(900px 520px at 10% 10%, rgba(120, 200, 255, 0.08), transparent 60%),
      radial-gradient(760px 520px at 90% 20%, rgba(255, 200, 100, 0.06), transparent 62%),
      radial-gradient(980px 820px at 50% 90%, rgba(200, 160, 255, 0.05), transparent 68%),
      repeating-linear-gradient(135deg, rgba(255, 255, 255, 0.02) 0 1px, transparent 1px 12px);
  pointer-events: none;
  z-index: 0;
}

/* 修改后的顶部栏 - 只保留返回按钮 */
.back-btn {
  position: absolute;
  top: 95px;      /* 距离顶部缝隙 */
  left: 40px;     /* 距离左侧缝隙 */
  z-index: 1000;
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

/* 活动主区域 */
.activity-main {
  position: relative;
  z-index: 1;
  padding: 24px;
  max-width: 1600px;
  margin: 0 auto;
  padding-top: 140px;
}

/* 加载状态 */
.loading-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 20, 25, 0.9);
  backdrop-filter: blur(8px);
  display: grid;
  place-items: center;
  z-index: 2000;
}

.loading-crystal {
  width: 80px;
  height: 80px;
  position: relative;
  animation: crystalFloat 3s ease-in-out infinite;
}

.crystal-inner {
  width: 100%;
  height: 100%;
  background: conic-gradient(
      from 0deg,
      rgba(120, 200, 255, 0.8),
      rgba(200, 160, 255, 0.8),
      rgba(255, 200, 100, 0.8),
      rgba(120, 200, 255, 0.8)
  );
  clip-path: polygon(50% 0%, 100% 50%, 50% 100%, 0% 50%);
  animation: crystalRotate 4s linear infinite;
}

.crystal-glow {
  position: absolute;
  inset: -20px;
  background: radial-gradient(circle at center,
  rgba(120, 200, 255, 0.3) 0%,
  transparent 70%);
  animation: crystalPulse 2s ease-in-out infinite;
}

.loading-text {
  margin-top: 24px;
  font-size: 16px;
  color: rgba(255, 255, 255, 0.8);
  animation: textPulse 1.5s ease-in-out infinite;
}

/* 活动列布局 */
.activity-columns {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
  margin-top: 16px;
}

@media (max-width: 1200px) {
  .activity-columns {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .activity-columns {
    grid-template-columns: 1fr;
  }
}

/* 列样式 */
.activity-column {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 20px;
  padding: 20px;
  backdrop-filter: blur(10px);
  transition: transform 0.3s ease;
}

.activity-column:hover {
  transform: translateY(-4px);
  border-color: rgba(255, 255, 255, 0.15);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
}

.column-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.column-icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: grid;
  place-items: center;
  font-size: 28px;
  border: 1px solid rgba(255, 255, 255, 0.15);
}

.column-title h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 800;
  letter-spacing: 0.02em;
}

.column-subtitle {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
  margin-top: 4px;
}

/* 活动卡片 */
.activity-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  max-height: 600px;
  overflow-y: auto;
  padding-right: 8px;
}

.activity-list::-webkit-scrollbar {
  width: 4px;
}

.activity-list::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 2px;
}

.activity-list::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 2px;
}

.activity-card {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 18px;
  position: relative;
  transition: all 0.3s ease;
  animation: cardAppear 0.6s ease-out;
  cursor: pointer;
}

.activity-card:hover {
  background: rgba(255, 255, 255, 0.06);
  border-color: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

.activity-card.full {
  opacity: 0.7;
  filter: grayscale(0.3);
}

.activity-card.internal {
  border-left: 4px solid rgba(255, 200, 100, 0.5);
}

.activity-card.assigned {
  border-left: 4px solid rgba(200, 160, 255, 0.5);
}

/* 已参与卡片样式 */
.activity-card.joined {
  background: linear-gradient(135deg,
  rgba(100, 200, 100, 0.08),
  rgba(100, 200, 100, 0.03));
  border-color: rgba(100, 200, 100, 0.2);
}

.activity-card.joined .activity-title {
  color: #64c864;
}

.activity-card.joined::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, #64c864, transparent);
  opacity: 0.6;
}

.activity-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.05em;
}

.activity-badge.internal {
  background: rgba(255, 200, 100, 0.15);
  border: 1px solid rgba(255, 200, 100, 0.3);
  color: #ffc864;
}

.activity-badge.assigned {
  background: rgba(200, 160, 255, 0.15);
  border: 1px solid rgba(200, 160, 255, 0.3);
  color: #c8a0ff;
}

.activity-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.activity-type {
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.activity-type.competition {
  background: rgba(255, 100, 100, 0.15);
  color: #ff6464;
}

.activity-type.training,
.activity-type.selection {
  background: rgba(120, 200, 255, 0.15);
  color: #78c8ff;
}

.activity-type.communication,
.activity-type.default {
  background: rgba(100, 200, 100, 0.15);
  color: #64c864;
}

.activity-type.internal {
  background: rgba(255, 200, 100, 0.15);
  color: #ffc864;
}

.activity-type.assigned {
  background: rgba(200, 160, 255, 0.15);
  color: #c8a0ff;
}

.activity-time {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
  font-family: 'SF Mono', monospace;
}

.activity-title {
  font-size: 18px;
  font-weight: 700;
  margin-bottom: 12px;
  line-height: 1.4;
  display: flex;
  align-items: center;
  gap: 8px;
}

/* 已参与标签 */
.participated-tag {
  padding: 2px 8px;
  border-radius: 8px;
  background: rgba(100, 200, 100, 0.15);
  color: #64c864;
  font-size: 11px;
  font-weight: 600;
  border: 1px solid rgba(100, 200, 100, 0.3);
}

.activity-meta {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.7);
}

.meta-icon {
  opacity: 0.8;
}

.full-badge {
  margin-left: 6px;
  padding: 2px 6px;
  border-radius: 8px;
  background: rgba(255, 80, 80, 0.15);
  color: #ff5050;
  font-size: 11px;
  font-weight: 600;
}

.activity-description {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
  line-height: 1.5;
  margin-bottom: 16px;
  padding: 10px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 8px;
  border-left: 3px solid rgba(200, 160, 255, 0.3);
}

.activity-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.activity-status {
  font-size: 12px;
}

.status-active {
  color: #64c864;
  padding: 2px 8px;
  background: rgba(100, 200, 100, 0.1);
  border-radius: 8px;
}

.status-completed {
  color: rgba(255, 255, 255, 0.5);
  padding: 2px 8px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
}

.status-cancelled {
  color: #ff6464;
  padding: 2px 8px;
  background: rgba(255, 100, 100, 0.1);
  border-radius: 8px;
}

/* 参加按钮 */
.join-btn {
  position: relative;
  padding: 8px 20px;
  border-radius: 12px;
  border: 1px solid rgba(120, 200, 255, 0.3);
  background: linear-gradient(135deg,
  rgba(120, 200, 255, 0.15),
  rgba(120, 200, 255, 0.08));
  color: #78c8ff;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
  overflow: hidden;
  min-width: 80px;
}

.join-btn:hover:not(.disabled):not(.joined) {
  background: linear-gradient(135deg,
  rgba(120, 200, 255, 0.25),
  rgba(120, 200, 255, 0.15));
  border-color: rgba(120, 200, 255, 0.5);
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(120, 200, 255, 0.2);
}

.join-btn.internal {
  border-color: rgba(255, 200, 100, 0.3);
  background: linear-gradient(135deg,
  rgba(255, 200, 100, 0.15),
  rgba(255, 200, 100, 0.08));
  color: #ffc864;
}

.join-btn.internal:hover:not(.disabled):not(.joined) {
  background: linear-gradient(135deg,
  rgba(255, 200, 100, 0.25),
  rgba(255, 200, 100, 0.15));
  border-color: rgba(255, 200, 100, 0.5);
  box-shadow: 0 8px 20px rgba(255, 200, 100, 0.2);
}

.join-btn.assigned {
  border-color: rgba(200, 160, 255, 0.3);
  background: linear-gradient(135deg,
  rgba(200, 160, 255, 0.15),
  rgba(200, 160, 255, 0.08));
  color: #c8a0ff;
}

.join-btn.assigned:hover:not(.disabled):not(.joined) {
  background: linear-gradient(135deg,
  rgba(200, 160, 255, 0.25),
  rgba(200, 160, 255, 0.15));
  border-color: rgba(200, 160, 255, 0.5);
  box-shadow: 0 8px 20px rgba(200, 160, 255, 0.2);
}

/* 已参与按钮样式 */
.join-btn.joined {
  background: linear-gradient(135deg,
  rgba(100, 200, 100, 0.15),
  rgba(100, 200, 100, 0.08));
  border-color: rgba(100, 200, 100, 0.3);
  color: #64c864;
  cursor: default;
}

.join-btn.disabled {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.4);
  cursor: not-allowed;
}

.join-btn.disabled:hover {
  transform: none;
  box-shadow: none;
}

/* 按钮火花效果 */
.btn-sparkle {
  position: absolute;
  inset: 0;
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.join-btn:hover:not(.disabled):not(.joined) .btn-sparkle {
  opacity: 1;
}

.spark {
  position: absolute;
  width: 2px;
  height: 12px;
  background: linear-gradient(to bottom,
  transparent,
  rgba(255, 255, 255, 0.8),
  transparent);
  animation: sparkFall 1s linear infinite;
}

.spark:nth-child(1) {
  left: 20%;
  animation-delay: 0s;
}

.spark:nth-child(2) {
  left: 50%;
  animation-delay: 0.2s;
}

.spark:nth-child(3) {
  left: 80%;
  animation-delay: 0.4s;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: rgba(255, 255, 255, 0.4);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.3;
}

.empty-text {
  font-size: 16px;
  margin-bottom: 8px;
}

.empty-hint {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.25);
}

/* 刷新按钮 */
.fab-refresh {
  position: fixed;
  bottom: 32px;
  right: 32px;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  border: 1px solid rgba(120, 200, 255, 0.3);
  background: linear-gradient(135deg,
  rgba(120, 200, 255, 0.15),
  rgba(120, 200, 255, 0.08));
  color: #78c8ff;
  font-size: 24px;
  cursor: pointer;
  z-index: 100;
  transition: all 0.3s ease;
  display: grid;
  place-items: center;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

.fab-refresh:hover:not(:disabled) {
  transform: rotate(180deg) scale(1.1);
  box-shadow: 0 12px 40px rgba(120, 200, 255, 0.3);
}

.fab-refresh:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 新建活动按钮样式 */
.fab-create {
  position: fixed;
  bottom: 32px;
  right: 100px;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  border: 1px solid rgba(100, 200, 100, 0.3);
  background: linear-gradient(135deg,
  rgba(100, 200, 100, 0.15),
  rgba(100, 200, 100, 0.08));
  color: #64c864;
  font-size: 28px;
  font-weight: 300;
  cursor: pointer;
  z-index: 100;
  transition: all 0.3s ease;
  display: grid;
  place-items: center;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  overflow: hidden;
}

.fab-create::before {
  content: '';
  position: absolute;
  inset: 0;
  background: radial-gradient(circle at center,
  rgba(100, 200, 100, 0.2) 0%,
  transparent 70%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.fab-create:hover:not(:disabled) {
  transform: scale(1.1);
  box-shadow: 0 12px 40px rgba(100, 200, 100, 0.3);
}

.fab-create:hover:not(:disabled)::before {
  opacity: 1;
}

.fab-create:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.fab-create span {
  transition: transform 0.3s ease;
}

.fab-create:hover:not(:disabled) span {
  transform: rotate(90deg);
}

.create-tooltip {
  position: absolute;
  bottom: 70px;
  right: 0;
  padding: 6px 12px;
  background: rgba(0, 0, 0, 0.8);
  color: white;
  font-size: 12px;
  border-radius: 6px;
  white-space: nowrap;
  opacity: 0;
  transform: translateY(10px);
  transition: all 0.3s ease;
  pointer-events: none;
}

.fab-create:hover .create-tooltip {
  opacity: 1;
  transform: translateY(0);
}

/* 动画 */
@keyframes crystalFloat {
  0%, 100% {
    transform: translateY(0) rotate(0deg);
  }
  50% {
    transform: translateY(-20px) rotate(180deg);
  }
}

@keyframes crystalRotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

@keyframes crystalPulse {
  0%, 100% {
    opacity: 0.5;
    transform: scale(1);
  }
  50% {
    opacity: 0.8;
    transform: scale(1.1);
  }
}

@keyframes textPulse {
  0%, 100% {
    opacity: 0.6;
  }
  50% {
    opacity: 1;
  }
}

@keyframes cardAppear {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes sparkFall {
  0% {
    transform: translateY(-20px) rotate(45deg);
    opacity: 0;
  }
  20% {
    opacity: 1;
  }
  80% {
    opacity: 1;
  }
  100% {
    transform: translateY(20px) rotate(45deg);
    opacity: 0;
  }
}

/* 响应式调整 */
@media (max-width: 768px) {
  .detail-header {
    margin: 16px;
    padding: 16px;
  }

  .activity-main {
    padding: 16px;
  }

  .activity-columns {
    gap: 16px;
  }

  .activity-column {
    padding: 16px;
  }

  .fab-refresh {
    bottom: 24px;
    right: 24px;
    width: 48px;
    height: 48px;
    font-size: 20px;
  }

  .fab-create {
    bottom: 24px;
    right: 80px;
    width: 48px;
    height: 48px;
    font-size: 24px;
  }
}
</style>
