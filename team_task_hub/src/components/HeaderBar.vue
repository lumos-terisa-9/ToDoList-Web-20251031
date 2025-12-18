<template>
  <header class="header-bar">
    <div class="logo">Logo</div>

    <div class="right-section">
      <!-- 动态按钮：在orgmap页面显示"个人界面"，其他页面显示"我的组织" -->
      <button class="header-btn" @click="handleNavigationClick">
        {{ isMapPage ? '个人界面' : '我的组织' }}
      </button>

      <button class="header-btn login-btn" @click="handleUserClick">
        <transition name="fade">
          <span key="login" v-if="!isLoggedIn">登录</span>
          <span key="hello" v-else>你好，{{ displayUsername }}</span>
        </transition>
      </button>
    </div>
  </header>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from "vue";
import { useRoute, useRouter } from 'vue-router'

const props = defineProps({
  currentUser: {
    type: Object,
    default: () => null
  }
})

const route = useRoute()
const router = useRouter()

// 本地用户状态，用于确保响应性
const localUser = ref(null)

// 初始化本地用户状态
function initLocalUser() {
  console.log('初始化HeaderBar用户状态，props:', props.currentUser)
  localUser.value = props.currentUser
}

// 监听props.currentUser的变化
watch(() => props.currentUser, (newUser) => {
  console.log('HeaderBar props用户状态变化:', newUser)
  localUser.value = newUser
}, { deep: true })

// 监听全局用户更新事件
function handleGlobalUserUpdate(event) {
  console.log('HeaderBar收到全局用户更新事件:', event.detail)
  if (event.detail.user) {
    localUser.value = event.detail.user
    console.log('HeaderBar本地用户状态已更新:', localUser.value)
  }
}

// 监听全局登出事件
function handleGlobalUserLogout() {
  console.log('HeaderBar收到全局登出事件')
  localUser.value = null
}

// 直接从本地状态获取用户信息
const currentUser = computed(() => {
  return localUser.value
})

const isLoggedIn = computed(() => {
  const loggedIn = !!currentUser.value
  console.log('HeaderBar登录状态:', loggedIn, '用户:', currentUser.value)
  return loggedIn
})

// 显示的用户名
const displayUsername = computed(() => {
  if (!currentUser.value) return '用户'
  return currentUser.value.username || '用户'
})

// 判断当前是否在 MapPage（orgmap 路由）
const isMapPage = computed(() => {
  return route.name === 'orgmap'
})

function handleUserClick() {
  console.log('HeaderBar用户点击，当前登录状态:', isLoggedIn.value, '用户:', currentUser.value)
  if (isLoggedIn.value) {
    // 显示个人信息模态框
    emit('show-profile')
  } else {
    // 未登录，触发打开登录模态框
    emit('open-login')
  }
}

function handleNavigationClick() {
  if (isMapPage.value) {
    // 在地图页面，点击跳转到个人页面
    router.push({ name: 'personpage' })
  } else {
    // 在其他页面，点击跳转到地图页面
    router.push({ name: 'orgmap' })
  }
}

// 定义事件
const emit = defineEmits(['open-login', 'show-profile'])

onMounted(() => {
  console.log('HeaderBar挂载')
  initLocalUser()

  // 添加事件监听器
  window.addEventListener('global-user-updated', handleGlobalUserUpdate)
  window.addEventListener('global-user-logout', handleGlobalUserLogout)

  // 监听storage变化（跨标签页同步）
  window.addEventListener('storage', (e) => {
    if (e.key === 'currentUser') {
      console.log('HeaderBar检测到currentUser变化，重新获取')
      // 重新从App.vue获取数据（通过props）
      // 这里依赖父组件会更新props
    }
  })
})

onUnmounted(() => {
  // 移除事件监听器
  window.removeEventListener('global-user-updated', handleGlobalUserUpdate)
  window.removeEventListener('global-user-logout', handleGlobalUserLogout)
})
</script>

<style scoped>
.header-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  z-index: 1000;
  background-color: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.3);
  padding: 15px 30px;
  color: white;
  font-size: 18px;
  transition: background 0.3s ease;
}

.logo {
  font-style: italic;
  font-weight: bold;
  letter-spacing: 1px;
}

.right-section {
  display: flex;
  gap: 20px;
  margin-right: 45px;
}

.header-btn {
  background: transparent;
  border: none;
  border-radius: 6px;
  color: white;
  padding: 6px 14px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.header-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: scale(1.05);
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
