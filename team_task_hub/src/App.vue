<template>
  <div class="app-container">
    <HeaderBar
      :current-user="globalUser"
      @open-login="openLoginModal"
      @show-profile="showProfileModalFunc"
    />

    <div class="scroll-content" ref="scrollContent">
      <router-view v-slot="{ Component }">
        <component
          :is="Component"
          ref="pageComponent"
          @open-login="openLoginModal"
        />
      </router-view>
    </div>

    <!-- 个人信息模态框 - 统一在App.vue中管理 -->
    <UserProfileModal
      :isVisible="showProfileModal"
      @close="showProfileModal = false"
      @update-user="handleUserUpdate"
      @logout="handleLogout"
    />

    <transition name="fade">
      <div v-if="showTopTip" class="scroll-tip top-tip">
        我已经到顶啦
      </div>
    </transition>

    <transition name="fade">
      <div v-if="showBottomTip" class="scroll-tip bottom-tip">
        我已经到底啦
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, provide } from "vue";
import { useRouter } from 'vue-router'
import HeaderBar from './components/HeaderBar.vue'
import UserProfileModal from './components/UserProfileModal.vue'

const router = useRouter()
const scrollContent = ref(null)
const pageComponent = ref(null)

const showTopTip = ref(false)
const showBottomTip = ref(false)
const showProfileModal = ref(false)

// 全局用户状态 - 确保响应式
const globalUser = ref(null)

// 从 localStorage 初始化用户状态
async function initGlobalUser() {
  const userData = localStorage.getItem('currentUser')
  console.log('初始化全局用户状态，localStorage数据:', userData)

  if (userData) {
    try {
      const parsedUser = JSON.parse(userData)
      console.log('解析后的用户数据:', parsedUser)

      // 确保用户数据有username字段
      if (parsedUser && typeof parsedUser === 'object') {
        // 如果是嵌套格式，提取实际用户数据
        const actualUser = parsedUser.data || parsedUser.user || parsedUser

        // 确保有username字段，如果没有则设置默认值
        if (!actualUser.username) {
          if (actualUser.id) {
            actualUser.username = `用户${actualUser.id}`
          } else if (actualUser.email) {
            actualUser.username = actualUser.email.split('@')[0]
          } else {
            actualUser.username = '用户'
          }
        }

        globalUser.value = actualUser
        console.log('设置全局用户状态为:', globalUser.value)
      }
    } catch (e) {
      console.error('解析全局用户数据失败:', e)
      globalUser.value = null
    }
  } else {
    console.log('localStorage中没有用户数据')
    globalUser.value = null
  }
}

// 检查token是否有效
async function checkTokenValidity() {
  const token = localStorage.getItem('token')
  if (!token) {
    console.log('未找到token')
    return false
  }

  try {
    console.log('检查token有效性...')
    const response = await fetch('http://localhost:8080/api/auth/me', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    if (response.ok) {
      const userData = await response.json()
      console.log('token有效，用户数据:', userData)

      // 提取实际用户数据
      const actualUser = userData.data || userData.user || userData

      // 确保有username字段
      if (!actualUser.username) {
        if (actualUser.id) {
          actualUser.username = `用户${actualUser.id}`
        } else if (actualUser.email) {
          actualUser.username = actualUser.email.split('@')[0]
        } else {
          actualUser.username = '用户'
        }
      }

      globalUser.value = actualUser
      localStorage.setItem('currentUser', JSON.stringify(actualUser))
      return true
    } else {
      console.log('token无效')
      // token无效，清除本地存储
      localStorage.removeItem('token')
      localStorage.removeItem('currentUser')
      globalUser.value = null
      return false
    }
  } catch (error) {
    console.error('检查token时发生错误:', error)
    return false
  }
}

// 提供全局用户状态给所有子组件
provide('globalUser', {
  user: globalUser,
  setUser: (newUser) => {
    console.log('设置新用户:', newUser)
    // 确保新用户有username字段
    if (newUser && !newUser.username) {
      if (newUser.id) {
        newUser.username = `用户${newUser.id}`
      } else if (newUser.email) {
        newUser.username = newUser.email.split('@')[0]
      } else {
        newUser.username = '用户'
      }
    }
    globalUser.value = newUser
    localStorage.setItem('currentUser', JSON.stringify(newUser))
    console.log('用户已保存到localStorage')
  },
  clearUser: () => {
    console.log('清除用户状态')
    globalUser.value = null
    localStorage.removeItem('currentUser')
    localStorage.removeItem('token')
  }
})

// 添加用户登录事件监听函数
function handleUserLogin(event) {
  console.log('收到用户登录事件:', event.detail)
  const { user, token } = event.detail

  // 更新全局用户状态
  if (user) {
    // 确保用户有username字段
    if (!user.username) {
      if (user.id) {
        user.username = `用户${user.id}`
      } else if (user.email) {
        user.username = user.email.split('@')[0]
      } else {
        user.username = '用户'
      }
    }

    globalUser.value = user
    localStorage.setItem('currentUser', JSON.stringify(user))
    localStorage.setItem('token', token)
    console.log('全局用户状态已更新:', globalUser.value)

    // 触发自定义事件通知其他组件
    window.dispatchEvent(new CustomEvent('global-user-updated', {
      detail: { user }
    }))
  }
}

// 添加storage事件监听（用于跨标签页同步）
function handleStorageChange(e) {
  console.log('storage变化:', e.key)
  if (e.key === 'currentUser') {
    // 当localStorage中的用户数据变化时，重新初始化
    console.log('currentUser发生变化，重新初始化...')
    initGlobalUser()
  } else if (e.key === 'token') {
    console.log('token发生变化')
  }
}

// 监听全局用户更新事件（来自其他组件）
function handleGlobalUserUpdate(event) {
  console.log('收到全局用户更新事件:', event.detail)
  if (event.detail.user) {
    globalUser.value = event.detail.user
    localStorage.setItem('currentUser', JSON.stringify(event.detail.user))
    console.log('通过事件更新全局用户状态:', globalUser.value)
  }
}

function openLoginModal() {
  if (pageComponent.value && pageComponent.value.openLoginModal) {
    pageComponent.value.openLoginModal()
  }
}

// 显示个人信息模态框
function showProfileModalFunc() {
  console.log('showProfileModalFunc被调用，当前用户状态:', globalUser.value)

  // 检查是否有用户登录
  if (globalUser.value) {
    console.log('显示个人信息模态框，当前用户:', globalUser.value)
    showProfileModal.value = true
  } else {
    console.log('用户未登录，显示登录框')
    openLoginModal()
  }
}

// 处理用户信息更新
function handleUserUpdate(updatedUser) {
  console.log('用户信息已更新:', updatedUser)

  // 确保更新后的用户有username字段
  if (updatedUser && !updatedUser.username) {
    if (updatedUser.id) {
      updatedUser.username = `用户${updatedUser.id}`
    } else if (updatedUser.email) {
      updatedUser.username = updatedUser.email.split('@')[0]
    } else {
      updatedUser.username = '用户'
    }
  }

  globalUser.value = updatedUser
  localStorage.setItem('currentUser', JSON.stringify(updatedUser))
  console.log('用户已更新到全局状态和localStorage')

  // 触发全局更新事件
  window.dispatchEvent(new CustomEvent('global-user-updated', {
    detail: { user: updatedUser }
  }))

  // 通知所有子组件用户已更新
  if (pageComponent.value && pageComponent.value.handleUserUpdate) {
    pageComponent.value.handleUserUpdate(updatedUser)
  }
}

// 处理登出
function handleLogout() {
  console.log('用户登出')
  globalUser.value = null
  localStorage.removeItem('token')
  localStorage.removeItem('currentUser')
  showProfileModal.value = false

  // 触发全局登出事件
  window.dispatchEvent(new CustomEvent('global-user-logout'))

  // 跳转到首页
  router.push('/')

  // 通知子组件
  if (pageComponent.value && pageComponent.value.handleLogout) {
    pageComponent.value.handleLogout()
  }
}

function setupScrollContainer() {
  const el = scrollContent.value
  if (!el) return

  const viewportHeight = window.innerHeight
  el.style.height = viewportHeight + 'px'
  el.style.overflowY = 'auto'
}

// 监听窗口大小变化
function handleResize() {
  setupScrollContainer()
}

// 软键盘弹起/收起重新计算高度
function handleVisualResize() {
  setupScrollContainer()
}

// 显示提示文字并自动淡出
function showTip(tipRef) {
  if (tipRef.value) return
  tipRef.value = true
  setTimeout(() => tipRef.value = false, 1500)
}

// Wheel - 检测边界显示提示
function handleWheel(e) {
  const el = scrollContent.value
  if (!el) return

  const scrollTop = el.scrollTop
  const scrollHeight = el.scrollHeight
  const clientHeight = el.clientHeight
  const maxScroll = Math.max(0, scrollHeight - clientHeight)

  // 检测边界 - 修复逻辑
  const atTop = scrollTop <= 0 && e.deltaY < 0
  const atBottom = scrollTop >= maxScroll && e.deltaY > 0

  if (atTop) {
    showTip(showTopTip)
  }
  if (atBottom && maxScroll > 0) { // 确保有可滚动内容
    showTip(showBottomTip)
  }
}

// 添加滚动事件监听来检测边界
function handleScroll() {
  const el = scrollContent.value
  if (!el) return

  const scrollTop = el.scrollTop
  const scrollHeight = el.scrollHeight
  const clientHeight = el.clientHeight
  const maxScroll = Math.max(0, scrollHeight - clientHeight)

  // 检测是否到达底部
  const atBottom = scrollTop >= maxScroll-0.66 && maxScroll > 0

  if (atBottom && !showBottomTip.value) {
    showTip(showBottomTip)
  }
}

onMounted(async () => {
  console.log('App.vue挂载，初始化用户状态...')

  // 先尝试从localStorage初始化
  await initGlobalUser()

  // 如果localStorage有token但没用户数据，或者用户数据不完整，验证token
  const token = localStorage.getItem('token')
  if (token && (!globalUser.value || !globalUser.value.username)) {
    console.log('有token但用户数据不完整，验证token...')
    await checkTokenValidity()
  }

  console.log('初始化后的用户状态:', globalUser.value)

  // 添加事件监听器
  window.addEventListener('user-login', handleUserLogin)
  window.addEventListener('storage', handleStorageChange)
  window.addEventListener('global-user-updated', handleGlobalUserUpdate)

  setupScrollContainer()
  window.addEventListener('resize', handleResize)
  const el = scrollContent.value
  el.addEventListener("wheel", handleWheel, { passive: false })
  el.addEventListener("scroll", handleScroll)
  window.visualViewport?.addEventListener("resize", handleVisualResize)
})

onUnmounted(() => {
  // 移除事件监听器
  window.removeEventListener('user-login', handleUserLogin)
  window.removeEventListener('storage', handleStorageChange)
  window.removeEventListener('global-user-updated', handleGlobalUserUpdate)

  const el = scrollContent.value
  window.removeEventListener("resize", handleResize)
  el.removeEventListener("wheel", handleWheel, { passive: false })
  el.removeEventListener("scroll", handleScroll)
  window.visualViewport?.removeEventListener("resize", handleVisualResize)
})
</script>

<style>
.app-container {
  min-height: 100vh;
  font-family: "Helvetica Neue", Arial, "PingFang SC", "Microsoft YaHei", sans-serif;
  overflow: hidden;
  position: relative;
}

.scroll-content {
  width: 100%;
  /* 高度通过JS动态设置 */
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
  scroll-behavior: smooth;
  box-sizing: border-box;
}

/* ========================================= */
/* === 通用滚动条样式 - 由子页面背景决定 === */
/* ========================================= */

/* 首页样式 - 紫色渐变背景 */
.home-page .scroll-content::-webkit-scrollbar-track {
  background: rgba(39, 145, 251, 0.2);
}

.home-page .scroll-content::-webkit-scrollbar-thumb {
  background: linear-gradient(to bottom, rgba(120, 180, 255, 0.6), rgba(153, 64, 250, 0.6));
}

.home-page .scroll-content::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(to bottom, rgba(140, 200, 255, 0.8), rgba(173, 84, 255, 0.8));
}

.home-page .scroll-content {
  scrollbar-color: rgba(120, 180, 255, 0.6) rgba(39, 145, 251, 0.2);
}

/* 个人页面样式 - 蓝色渐变背景 */
.person-page .scroll-content::-webkit-scrollbar-track {
  background: rgba(14, 89, 184, 0.2);
}

.person-page .scroll-content::-webkit-scrollbar-thumb {
  background: linear-gradient(to bottom, rgba(30, 144, 255, 0.6), rgba(22, 177, 244, 0.6));
}

.person-page .scroll-content::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(to bottom, rgba(77, 166, 255, 0.8), rgba(52, 197, 255, 0.8));
}

.person-page .scroll-content {
  scrollbar-color: rgba(30, 144, 255, 0.6) rgba(14, 89, 184, 0.2);
}

/* 基础滚动条样式 */
.scroll-content::-webkit-scrollbar {
  width: 6px;
}

.scroll-content::-webkit-scrollbar-track {
  border-radius: 3px;
  transition: background 0.3s ease;
}

.scroll-content::-webkit-scrollbar-thumb {
  border-radius: 3px;
  transition: background 0.3s ease;
}

.scroll-content {
  scrollbar-width: thin;
  scrollbar-gutter: stable;
}

.page-content {
  padding: 24px;
  color: white;
  margin-top: 120px;
  text-align: center;
}

.scroll-tip {
  position: fixed;
  left: 50%;
  color: white;
  font-size: 14px;
  background: rgba(34, 34, 34, 0.6);
  padding: 8px 16px;
  border-radius: 8px;
  z-index: 2000;
  pointer-events: none;
  white-space: nowrap;
  line-height: 1.2;
  height: auto;
  min-height: auto;
}

.top-tip {
  top: 100px;
  transform: translateX(-70%);
}

.bottom-tip {
  bottom: 45px;
  transform: translateX(-70%);
}

.fade-enter-active, .fade-leave-active { transition: opacity 0.5s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

/* ============ 地图 & 迷雾的全局层级控制 ============ */

/* 地图容器本身（Leaflet 会塞各种 pane 进去） */
.map {
  position: relative; /* 确保绝对定位的雾可以以它为参考 */
}

/* 迷雾层：盖在地图上，但在 UI 面板之下 */
.fog-layer {
  position: absolute;
  inset: 0;
  pointer-events: none;  /* 不拦截鼠标事件 */
  z-index: 400;          /* 比地图瓦片/overlay 高 */
}

/* 左侧控制面板 + 右侧详情面板：永远在最上层 */
.map-controls,
.location-detail {
  position: absolute;
  z-index: 1000;         /* 明确高于迷雾 */
}

/* Leaflet 的 marker pane 提高层级，始终在雾上面 */
.leaflet-marker-pane {
  z-index: 700 !important;
}
</style>
