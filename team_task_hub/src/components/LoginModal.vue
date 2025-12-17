<template>
  <transition name="modal">
    <div v-if="isVisible" class="modal-overlay" @click.self="close">
      <div class="modal-container">
        <div class="modal-header">
          <h3>{{ isLogin ? '登录' : '注册' }}</h3>
          <button class="close-btn" @click="close">×</button>
        </div>

        <div class="modal-body">
          <!-- 登录表单 -->
          <form v-if="isLogin" @submit.prevent="handleLogin">
            <div class="form-group">
              <input
                type="text"
                v-model="loginForm.identifier"
                placeholder="用户ID或邮箱"
                required
              >
            </div>
            <div class="form-group">
              <input
                type="password"
                v-model="loginForm.password"
                placeholder="密码"
                required
              >
            </div>
            <button type="submit" class="submit-btn" :disabled="loading">
              {{ loading ? '登录中...' : '登录' }}
            </button>
          </form>

          <!-- 注册表单 -->
          <form v-else @submit.prevent="handleRegister">
            <div class="form-group">
              <input
                type="text"
                v-model="registerForm.id"
                placeholder="用户ID"
                required
              >
            </div>
            <div class="form-group">
              <input
                type="text"
                v-model="registerForm.username"
                placeholder="用户名"
                required
              >
            </div>
            <div class="form-group">
              <input
                type="email"
                v-model="registerForm.email"
                placeholder="邮箱"
                required
              >
            </div>
            <div class="form-group verification-group">
              <input
                type="text"
                v-model="registerForm.code"
                placeholder="验证码"
                required
              >
              <button
                type="button"
                class="verification-btn"
                :disabled="verificationCooldown > 0 || !registerForm.email"
                @click="sendVerificationCode"
              >
                {{ verificationCooldown > 0 ? `${verificationCooldown}s` : '获取验证码' }}
              </button>
            </div>
            <div class="form-group">
              <input
                type="password"
                v-model="registerForm.password"
                placeholder="密码"
                required
              >
            </div>
            <div class="form-group">
              <input
                type="password"
                v-model="registerForm.confirmPassword"
                placeholder="确认密码"
                required
              >
            </div>
            <button type="submit" class="submit-btn" :disabled="loading">
              {{ loading ? '注册中...' : '注册' }}
            </button>
          </form>

          <div class="switch-mode">
            <span @click="toggleMode">
              {{ isLogin ? '没有账号？立即注册' : '已有账号？立即登录' }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  isVisible: Boolean
})

const emit = defineEmits(['close', 'login-success'])

const router = useRouter()
const isLogin = ref(true)
const loading = ref(false)
const verificationCooldown = ref(0)

// API 基础URL
const API_BASE = 'http://localhost:8080/api'

const loginForm = ref({
  identifier: '',
  password: ''
})

const registerForm = ref({
  id: '',
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  code: ''
})

function close() {
  emit('close')
  resetForms()
  isLogin.value = true
}

function resetForms() {
  loginForm.value = { identifier: '', password: '' }
  registerForm.value = {
    id: '',
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
    code: ''
  }
}

function toggleMode() {
  isLogin.value = !isLogin.value
  resetForms()
}

async function handleLogin() {
  if (!loginForm.value.identifier || !loginForm.value.password) {
    alert('请填写完整的登录信息')
    return
  }

  loading.value = true

  try {
    console.log('开始登录请求...')
    const response = await fetch(`${API_BASE}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        identifier: loginForm.value.identifier,
        password: loginForm.value.password
      })
    })

    console.log('登录响应状态:', response.status)
    const data = await response.json()
    console.log('登录响应数据:', data)

    if (response.ok && data.success) {
      // 登录成功 - 根据提供的响应格式提取数据
      const responseData = data.data
      console.log('提取到的数据:', responseData)

      // 提取token
      const token = responseData.access_token
      console.log('Token:', token)

      // 提取用户信息 - 关键修改点
      const userData = responseData.user || {}
      console.log('用户数据:', userData)

      // 确保userData有正确的字段名
      const formattedUser = {
        id: userData.user_id || userData.id,
        username: userData.username,
        email: userData.email
      }

      console.log('格式化后的用户数据:', formattedUser)

      // 验证是否有用户名
      if (!formattedUser.username) {
        console.error('登录响应中没有用户名字段')
        alert('登录成功，但无法获取用户名信息')
        return
      }

      // 保存令牌到本地存储
      localStorage.setItem('token', token)

      // 保存完整的格式化用户数据
      localStorage.setItem('currentUser', JSON.stringify(formattedUser))

      console.log('登录信息已保存到localStorage:', {
        token: token.substring(0, 20) + '...',
        user: formattedUser
      })

      // 触发自定义事件通知App.vue更新状态
      window.dispatchEvent(new CustomEvent('user-login', {
        detail: {
          user: formattedUser,
          token: token
        }
      }))

      emit('login-success', formattedUser)
      close()

      // 跳转到个人页面
      router.push('/personpage')
    } else {
      // 登录失败
      alert(data.message || '登录失败，请检查用户名和密码')
    }
  } catch (error) {
    console.error('登录请求失败:', error)
    alert('登录失败，请检查网络连接')
  } finally {
    loading.value = false
  }
}

async function handleRegister() {
  // 前端验证
  if (!registerForm.value.id || !registerForm.value.username || !registerForm.value.email ||
    !registerForm.value.password || !registerForm.value.code) {
    alert('请填写所有必填字段')
    return
  }

  if (registerForm.value.password !== registerForm.value.confirmPassword) {
    alert('两次输入的密码不一致')
    return
  }

  if (registerForm.value.password.length < 6) {
    alert('密码长度至少6位')
    return
  }

  loading.value = true

  try {
    const response = await fetch(`${API_BASE}/auth/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        id: parseInt(registerForm.value.id),
        username: registerForm.value.username,
        email: registerForm.value.email,
        password: registerForm.value.password,
        code: registerForm.value.code
      })
    })

    const data = await response.json()

    if (response.ok) {
      alert('注册成功！请登录')
      isLogin.value = true
      resetForms()
    } else {
      alert(data.message || '注册失败')
    }
  } catch (error) {
    console.error('注册请求失败:', error)
    alert('注册失败，请检查网络连接')
  } finally {
    loading.value = false
  }
}

async function sendVerificationCode() {
  if (!registerForm.value.email) {
    alert('请输入邮箱')
    return
  }

  // 邮箱格式验证
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(registerForm.value.email)) {
    alert('请输入有效的邮箱地址')
    return
  }

  try {
    const response = await fetch(`${API_BASE}/email/send-verification`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: registerForm.value.email,
        business: 'register' // 注册业务
      })
    })

    const data = await response.json()

    if (response.ok) {
      alert('验证码已发送到您的邮箱，请查收')

      // 开始倒计时
      verificationCooldown.value = 60
      const timer = setInterval(() => {
        verificationCooldown.value--
        if (verificationCooldown.value <= 0) {
          clearInterval(timer)
        }
      }, 1000)
    } else {
      alert(data.message || '验证码发送失败')
    }
  } catch (error) {
    console.error('发送验证码失败:', error)
    alert('验证码发送失败，请检查网络连接')
  }
}

// 监听模态框显示状态
watch(() => props.isVisible, (newVal) => {
  if (newVal) {
    resetForms()
  }
})
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 2000;
}

.modal-container {
  background: white;
  border-radius: 12px;
  width: 400px;
  max-width: 90vw;
  color: #333;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #eee;
}

.modal-header h3 {
  margin: 0;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group input {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  box-sizing: border-box;
}

.verification-group {
  display: flex;
  gap: 10px;
}

.verification-group input {
  flex: 1;
}

.verification-btn {
  padding: 12px 16px;
  background: #409eff;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  white-space: nowrap;
}

.verification-btn:disabled {
  background: #c0c4cc;
  cursor: not-allowed;
}

.submit-btn {
  width: 100%;
  padding: 12px;
  background: #409eff;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 16px;
  cursor: pointer;
  margin-bottom: 16px;
}

.submit-btn:disabled {
  background: #c0c4cc;
  cursor: not-allowed;
}

.switch-mode {
  text-align: center;
}

.switch-mode span {
  color: #409eff;
  cursor: pointer;
  font-size: 14px;
}

.modal-enter-active, .modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from, .modal-leave-to {
  opacity: 0;
}

.modal-enter-active .modal-container,
.modal-leave-active .modal-container {
  transition: transform 0.3s ease;
}

.modal-enter-from .modal-container,
.modal-leave-to .modal-container {
  transform: scale(0.9);
}
</style>
