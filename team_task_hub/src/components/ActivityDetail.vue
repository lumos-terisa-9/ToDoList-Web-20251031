<template>
  <div class="activity-detail-page">
    <!-- 加载状态 -->
    <div v-if="loading" class="loading-overlay">
      <div class="loading-crystal">
        <div class="crystal-inner"></div>
        <div class="crystal-glow"></div>
      </div>
      <div class="loading-text">加载活动详情中...</div>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="!activity" class="error-overlay">
      <div class="error-icon">⚠️</div>
      <div class="error-text">加载活动详情失败</div>
      <button class="back-btn" @click="handleClose">返回</button>
    </div>

    <!-- 根据角色显示不同界面 -->
    <div v-else>
      <!-- Admin或Creator界面 -->
      <ActivityDetailAdmin
          v-if="userRole === 'Admin' || userRole === 'Creator'"
          :activity="activity"
          :org-id="orgId"
          :org-name="orgName"
          @close="handleClose"
          @activity-updated="handleActivityUpdated"
      />

      <!-- Member或None角色界面 -->
      <ActivityDetailMember
          v-else
          :activity="activity"
          :org-id="orgId"
          :org-name="orgName"
          @close="handleClose"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ActivityDetailAdmin from './ActivityDetailAdmin.vue'
import ActivityDetailMember from './ActivityDetailMember.vue'

const route = useRoute()
const router = useRouter()

// 基础数据
const orgId = route.params.orgId
const orgName = "羽毛球队"
const activityId = route.params.activityId
const fromType = route.params.fromType || 'public'

// 状态
const loading = ref(true)
const userRole = ref('')
const activity = ref(null)

// 从sessionStorage获取缓存的ActivityPage数据
function getCachedActivities() {
  try {
    const cacheData = sessionStorage.getItem('activityCache')
    if (cacheData) {
      const parsedData = JSON.parse(cacheData)

      // 检查缓存是否过期（5分钟内有效）
      const cacheTime = parsedData.timestamp || 0
      const currentTime = Date.now()

      if (currentTime - cacheTime < 5 * 60 * 1000 && parsedData.orgId === orgId) {
        return parsedData.data
      }
    }
  } catch (error) {
    console.error('读取缓存数据失败:', error)
  }
  return null
}

// 根据ID从缓存中查找活动
function findActivityFromCache() {
  const cacheData = getCachedActivities()

  if (!cacheData) {
    console.log('未找到缓存数据')
    return null
  }

  console.log('缓存数据:', cacheData)

  // 优先从fromType指定的类型中查找
  if (fromType && cacheData[fromType]) {
    const found = cacheData[fromType].find(item => item.id.toString() === activityId)
    if (found) {
      console.log('从指定类型中找到活动:', fromType, found)
      return found
    }
  }

  // 如果没找到，从所有类型中查找
  const allActivities = [
    ...(cacheData.public || []),
    ...(cacheData.internal || []),
    ...(cacheData.assigned || [])
  ]

  const found = allActivities.find(item => item.id.toString() === activityId)
  if (found) {
    console.log('从所有活动中找到:', found)
    return found
  }

  console.log('未找到活动ID:', activityId)
  return null
}

// 从localStorage获取当前用户
function getCurrentUser() {
  try {
    const userData = localStorage.getItem('currentUser')
    if (userData) {
      const parsedData = JSON.parse(userData)
      if (parsedData.success && parsedData.data) {
        return parsedData.data
      }
      if (parsedData.id || parsedData.user_id) {
        return parsedData
      }
    }
  } catch (error) {
    console.error('获取用户信息失败:', error)
  }
  return null
}

// 检查用户角色
async function checkUserRole() {
  const token = getToken()
  if (!token) {
    console.error('未找到认证令牌')
    userRole.value = 'Member'
    return
  }

  const user = getCurrentUser()
  if (!user || !user.id) {
    console.error('未找到用户ID')
    userRole.value = 'Member'
    return
  }

  try {
    const response = await fetch(`http://localhost:8080/api/organization/${orgId}/users/${user.id}/role`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    if (response.ok) {
      const data = await response.json()
      if (data.success && data.data) {
        userRole.value = data.data.role
      } else {
        userRole.value = 'Member'
      }
    } else {
      userRole.value = 'Member'
    }
  } catch (error) {
    console.error('检查用户权限失败:', error)
    userRole.value = 'Member'
  }
}

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

// 处理关闭
function handleClose() {
  router.back()
}

// 处理活动更新
function handleActivityUpdated() {
  // 清除缓存，让用户返回ActivityPage时重新加载
  sessionStorage.removeItem('activityCache')
  handleClose()
}

// 初始化
onMounted(async () => {
  loading.value = true
  console.log('ActivityDetail初始化:', {
    orgId,
    activityId,
    fromType,
    routeParams: route.params
  })

  // 1. 先从缓存中查找活动
  const cachedActivity = findActivityFromCache()
  if (cachedActivity) {
    activity.value = cachedActivity
    console.log('从缓存中获取到活动:', activity.value)
  } else {
    console.log('未找到缓存的活动数据')
  }

  // 2. 检查用户角色
  await checkUserRole()

  loading.value = false
})
</script>

<style scoped>
/* 样式保持不变 */
.activity-detail-page {
  min-height: 100vh;
  background: #0f1419;
  position: relative;
  padding: 80px 100px 10px 100px;
}

.loading-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 20, 25, 0.9);
  backdrop-filter: blur(8px);
  display: grid;
  place-items: center;
  z-index: 2000;
}

.error-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 20, 25, 0.9);
  backdrop-filter: blur(8px);
  display: grid;
  place-items: center;
  z-index: 2000;
  color: white;
  text-align: center;
}

.error-icon {
  font-size: 64px;
  margin-bottom: 20px;
}

.error-text {
  font-size: 20px;
  margin-bottom: 30px;
  color: rgba(255, 255, 255, 0.8);
}

.back-btn {
  padding: 12px 32px;
  border-radius: 12px;
  background: rgba(120, 200, 255, 0.15);
  border: 1px solid rgba(120, 200, 255, 0.3);
  color: #78c8ff;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.back-btn:hover {
  background: rgba(120, 200, 255, 0.25);
  border-color: rgba(120, 200, 255, 0.5);
  transform: translateY(-2px);
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
</style>
