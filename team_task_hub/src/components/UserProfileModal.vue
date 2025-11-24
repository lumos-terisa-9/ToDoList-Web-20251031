<template>
  <transition name="modal">
    <div v-if="isVisible" class="modal-overlay" @click.self="close">
      <div class="modal-container">
        <div class="modal-header">
          <h3>个人信息</h3>
          <button class="close-btn" @click="close">×</button>
        </div>

        <div class="modal-body">
          <div class="profile-layout">
            <!-- 左侧导航 -->
            <div class="sidebar">
              <div
                v-for="item in menuItems"
                :key="item.key"
                :class="['menu-item', { active: activeMenu === item.key }]"
                @click="activeMenu = item.key"
              >
                <span class="menu-icon">{{ item.icon }}</span>
                <span class="menu-text">{{ item.text }}</span>
              </div>
            </div>

            <!-- 右侧内容 -->
            <div class="content-area">
              <!-- 个人信息 -->
              <div v-if="activeMenu === 'profile'" class="content-section">
                <h4 class="section-title">基本信息</h4>

                <!-- 头像区域 -->
                <div class="avatar-section">
                  <div class="avatar-container">
                    <div class="avatar-wrapper">
                      <img :src="currentUser?.avatar_url || '/空白头像.png'" alt="头像" class="avatar">
                      <div class="avatar-overlay">
                        <button class="avatar-edit-btn" @click="$refs.avatarInput.click()">
                          📷 更换
                        </button>
                      </div>
                    </div>
                    <input
                      type="file"
                      ref="avatarInput"
                      accept="image/*"
                      @change="handleAvatarUpload"
                      style="display: none"
                    >
                  </div>
                </div>

                <!-- 用户信息 -->
                <div class="info-grid">
                  <div class="info-item">
                    <label class="info-label">用户ID</label>
                    <div class="info-value masked">
                      {{ maskUserId(currentUser?.id) }}
                    </div>
                  </div>

                  <div class="info-item">
                    <label class="info-label">用户名</label>
                    <input
                      type="text"
                      v-model="userForm.username"
                      class="info-input"
                      placeholder="请输入用户名"
                    >
                  </div>

                  <div class="info-item">
                    <label class="info-label">邮箱</label>
                    <div class="info-value masked">
                      {{ maskEmail(currentUser?.email) }}
                    </div>
                  </div>
                </div>

                <div class="profile-actions">
                  <button class="action-btn primary" @click="saveProfile" :disabled="loading">
                    {{ loading ? '保存中...' : '保存修改' }}
                  </button>
                </div>
              </div>

              <!-- 密码修改 -->
              <div v-if="activeMenu === 'password'" class="content-section">
                <h4 class="section-title">修改密码</h4>

                <div class="security-notice">
                  <div class="notice-icon">🔒</div>
                  <div class="notice-content">
                    <p>为了保障账户安全，修改密码需要邮箱验证</p>
                  </div>
                </div>

                <form @submit.prevent="changePassword" class="security-form">
                  <div class="form-group">
                    <label class="form-label">邮箱验证码</label>
                    <div class="verification-group">
                      <input
                        type="text"
                        v-model="passwordForm.verificationCode"
                        placeholder="请输入验证码"
                        required
                        class="form-input"
                      >
                      <button
                        type="button"
                        class="verification-btn"
                        :disabled="passwordCooldown > 0"
                        @click="sendPasswordVerificationCode"
                      >
                        {{ passwordCooldown > 0 ? `${passwordCooldown}s` : '获取验证码' }}
                      </button>
                    </div>
                  </div>

                  <div class="form-group">
                    <label class="form-label">当前密码</label>
                    <div class="masked-password">
                      ••••••••••
                    </div>
                  </div>

                  <div class="form-group">
                    <label class="form-label">新密码</label>
                    <input
                      type="password"
                      v-model="passwordForm.newPassword"
                      placeholder="请输入新密码"
                      required
                      class="form-input"
                    >
                  </div>

                  <div class="form-group">
                    <label class="form-label">确认新密码</label>
                    <input
                      type="password"
                      v-model="passwordForm.confirmPassword"
                      placeholder="请再次输入新密码"
                      required
                      class="form-input"
                    >
                  </div>

                  <div class="form-actions">
                    <button type="submit" class="action-btn primary" :disabled="loading">
                      {{ loading ? '修改中...' : '确认修改' }}
                    </button>
                  </div>
                </form>
              </div>

              <!-- 邮箱修改 -->
              <div v-if="activeMenu === 'email'" class="content-section">
                <h4 class="section-title">修改邮箱</h4>

                <div class="security-notice">
                  <div class="notice-icon">📧</div>
                  <div class="notice-content">
                    <p>修改邮箱需要原邮箱验证和新邮箱验证</p>
                  </div>
                </div>

                <form @submit.prevent="changeEmail" class="security-form">
                  <div class="form-group">
                    <label class="form-label">当前邮箱</label>
                    <div class="masked-email">
                      {{ maskEmail(currentUser?.email) }}
                    </div>
                  </div>

                  <div class="form-group">
                    <label class="form-label">原邮箱验证码</label>
                    <div class="verification-group">
                      <input
                        type="text"
                        v-model="emailForm.oldVerificationCode"
                        placeholder="请输入原邮箱验证码"
                        required
                        class="form-input"
                      >
                      <button
                        type="button"
                        class="verification-btn"
                        :disabled="oldEmailCooldown > 0"
                        @click="sendOldEmailVerificationCode"
                      >
                        {{ oldEmailCooldown > 0 ? `${oldEmailCooldown}s` : '获取验证码' }}
                      </button>
                    </div>
                  </div>

                  <div class="form-group">
                    <label class="form-label">新邮箱地址</label>
                    <input
                      type="email"
                      v-model="emailForm.newEmail"
                      placeholder="请输入新邮箱地址"
                      required
                      class="form-input"
                    >
                  </div>

                  <div class="form-group">
                    <label class="form-label">新邮箱验证码</label>
                    <div class="verification-group">
                      <input
                        type="text"
                        v-model="emailForm.newVerificationCode"
                        placeholder="请输入新邮箱验证码"
                        required
                        class="form-input"
                      >
                      <button
                        type="button"
                        class="verification-btn"
                        :disabled="newEmailCooldown > 0"
                        @click="sendNewEmailVerificationCode"
                      >
                        {{ newEmailCooldown > 0 ? `${newEmailCooldown}s` : '获取验证码' }}
                      </button>
                    </div>
                  </div>

                  <div class="form-actions">
                    <button type="submit" class="action-btn primary" :disabled="loading">
                      {{ loading ? '修改中...' : '确认修改' }}
                    </button>
                  </div>
                </form>
              </div>

              <!-- 退出登录 -->
              <div v-if="activeMenu === 'logout'" class="content-section">
                <h4 class="section-title">账户安全</h4>
                <div class="logout-content">
                  <div class="logout-icon">🚪</div>
                  <h5 class="logout-title">退出登录</h5>
                  <p class="logout-text">确定要退出当前账户吗？</p>
                  <p class="logout-desc">退出后需要重新登录才能访问个人页面</p>
                  <div class="logout-actions">
                    <button class="action-btn secondary" @click="activeMenu = 'profile'">
                      取消
                    </button>
                    <button class="action-btn danger" @click="handleLogout">
                      确认退出
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const props = defineProps({
  isVisible: Boolean
})

const emit = defineEmits(['close', 'update-user', 'logout'])

const loading = ref(false)
const passwordCooldown = ref(0)
const oldEmailCooldown = ref(0)
const newEmailCooldown = ref(0)
const avatarInput = ref(null)
const activeMenu = ref('profile')
const currentUser = ref(null)

// API 基础URL
const API_BASE = 'http://localhost:8080/api'

// 菜单项配置
const menuItems = ref([
  { key: 'profile', text: '基本信息', icon: '👤' },
  { key: 'password', text: '修改密码', icon: '🔒' },
  { key: 'email', text: '修改邮箱', icon: '📧' },
  { key: 'logout', text: '退出登录', icon: '🚪' }
])

// 表单数据
const userForm = ref({
  username: ''
})

const passwordForm = ref({
  newPassword: '',
  confirmPassword: '',
  verificationCode: ''
})

const emailForm = ref({
  newEmail: '',
  oldVerificationCode: '',
  newVerificationCode: ''
})

// 获取当前用户信息
async function fetchCurrentUser() {
  // 从本地浏览器获取token
  let token = localStorage.getItem('token')
  console.log('从localStorage获取原始token:', token)

  // 如果token是JSON字符串，解析它
  if (token && token.startsWith('{')) {
    try {
      const tokenData = JSON.parse(token)
      console.log('解析token数据:', tokenData)

      // 提取纯token字符串
      if (tokenData.data && tokenData.data.access_token) {
        token = tokenData.data.access_token
      } else if (tokenData.access_token) {
        token = tokenData.access_token
      } else if (tokenData.token) {
        token = tokenData.token
      }
      console.log('提取后的纯token:', token)
    } catch (error) {
      console.error('解析token失败:', error)
      return null
    }
  }

  if (!token) {
    console.error('未找到认证令牌')
    return null
  }

  try {
    console.log('开始调用 /auth/me 接口...')
    const response = await fetch(`${API_BASE}/auth/me`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    console.log('接口响应状态:', response.status)

    if (response.ok) {
      const result = await response.json()
      console.log('接口返回数据:', result)

      // 处理头像URL，确保使用GitHub URL
      const userData = result.data || result
      if (userData.avatar_url) {
        userData.avatar_url = ensureGitHubAvatarUrl(userData.avatar_url)
      }

      return userData
    } else {
      console.error('获取用户信息失败:', response.status)
      const errorText = await response.text()
      console.error('错误详情:', errorText)
      return null
    }
  } catch (error) {
    console.error('获取用户信息请求失败:', error)
    return null
  }
}

// 用户登出
async function logoutUser() {
  // 从本地浏览器获取token
  let token = localStorage.getItem('token')

  // 如果token是JSON字符串，解析它
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
      return false
    }
  }

  if (!token) {
    console.error('未找到认证令牌')
    return false
  }

  try {
    console.log('开始登出，使用token:', token)
    const response = await fetch(`${API_BASE}/auth/logout`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    console.log('登出响应状态:', response.status)

    if (response.ok) {
      const result = await response.json()
      console.log('登出成功:', result)
      return true
    } else {
      console.error('登出失败:', response.status)
      const errorText = await response.text()
      console.error('登出错误详情:', errorText)
      return false
    }
  } catch (error) {
    console.error('登出请求失败:', error)
    return false
  }
}

async function handleLogout() {
  loading.value = true

  try {
    // 调用后端登出接口
    const logoutSuccess = await logoutUser()

    // 无论登出是否成功，都清除本地存储的登录状态，
    localStorage.removeItem('token')
    localStorage.removeItem('currentUser')
    currentUser.value = null

    emit('logout')
    close()

    // 跳转到首页
    router.push('/')

    if (logoutSuccess) {
      alert('已成功退出登录')
    } else {
      alert('已退出登录（服务器登出失败，但已清除本地状态）')
    }

  } catch (error) {
    console.error('登出过程中发生错误:', error)
    // 发生错误时仍然清除本地状态
    localStorage.removeItem('token')
    localStorage.removeItem('currentUser')
    currentUser.value = null

    emit('logout')
    close()

    router.push('/')
    alert('已退出登录')
  } finally {
    loading.value = false
  }
}

// 更新用户名
async function updateUsername(username) {
  // 从本地浏览器获取token
  let token = localStorage.getItem('token')

  // 如果token是JSON字符串，解析它
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
      return { success: false, message: '令牌格式错误' }
    }
  }

  if (!token) {
    console.error('未找到认证令牌')
    return { success: false, message: '未找到认证令牌' }
  }

  try {
    const response = await fetch(`${API_BASE}/auth/change_userName`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ username })
    })

    const result = await response.json()
    console.log('更新用户名响应:', result)

    if (response.ok) {
      return { success: true, data: result }
    } else {
      return { success: false, message: result.message || '用户名更新失败' }
    }
  } catch (error) {
    console.error('更新用户名请求失败:', error)
    return { success: false, message: '网络请求失败' }
  }
}

// 配置
const GITHUB_CONFIG = {
}

// 获取默认头像URL - 修正为GitHub Pages格式
function getDefaultAvatarUrl() {
  // GitHub Pages仓库的特殊URL格式：https://用户名.github.io/文件夹/文件名
  return `https://${GITHUB_CONFIG.username}.github.io/${GITHUB_CONFIG.folder}/default-avatar.png`
}

// 确保头像URL使用GitHub URL
function ensureGitHubAvatarUrl(avatarUrl) {
  if (!avatarUrl) return getDefaultAvatarUrl()

  // 如果已经是GitHub URL，直接返回
  if (avatarUrl.includes('github.io') || avatarUrl.includes('githubusercontent.com')) {
    return avatarUrl
  }

  // 如果是本地URL或无效URL，返回默认头像
  if (avatarUrl.startsWith('blob:') || avatarUrl.startsWith('data:') || !avatarUrl.startsWith('http')) {
    return getDefaultAvatarUrl()
  }

  return avatarUrl
}

// 更新头像
async function updateAvatar(avatarUrl) {
  // 从本地浏览器获取token
  let token = localStorage.getItem('token')

  // 如果token是JSON字符串，解析它
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
      return { success: false, message: '令牌格式错误' }
    }
  }

  if (!token) {
    console.error('未找到认证令牌')
    return { success: false, message: '未找到认证令牌' }
  }

  try {
    const response = await fetch(`${API_BASE}/auth/change_avatar`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ avatar_url: avatarUrl })
    })

    const result = await response.json()
    console.log('更新头像响应:', result)

    if (response.ok) {
      return { success: true, data: result }
    } else {
      return { success: false, message: result.message || '头像更新失败' }
    }
  } catch (error) {
    console.error('更新头像请求失败:', error)
    return { success: false, message: '网络请求失败' }
  }
}

// 使用GitHub API上传头像 - 修正URL返回
async function uploadToGitHub(file) {
  try {
    // 将文件转换为Base64
    const base64Data = await fileToBase64(file)
    const cleanBase64 = base64Data.split(',')[1] // 移除data:image/jpeg;base64,前缀

    // 生成唯一的文件名
    const fileExtension = file.type.split('/')[1] || 'png'
    const fileName = `avatar_${Date.now()}_${Math.random().toString(36).substr(2, 9)}.${fileExtension}`

    // 构造API URL - GitHub Pages仓库的特殊路径
    const apiUrl = `https://api.github.com/repos/${GITHUB_CONFIG.username}/${GITHUB_CONFIG.repo}/contents/${GITHUB_CONFIG.folder}/${fileName}`

    // 请求数据
    const requestData = {
      message: `Upload avatar: ${fileName}`,
      content: cleanBase64,
      branch: GITHUB_CONFIG.branch
    }

    console.log('开始上传到GitHub:', apiUrl)
    console.log('文件信息:', {
      name: fileName,
      size: file.size,
      type: file.type
    })

    const response = await fetch(apiUrl, {
      method: 'PUT',
      headers: {
        'Authorization': `token ${GITHUB_CONFIG.token}`,
        'Content-Type': 'application/json',
        'Accept': 'application/vnd.github.v3+json'
      },
      body: JSON.stringify(requestData)
    })

    console.log('GitHub API响应状态:', response.status)

    if (!response.ok) {
      const errorData = await response.json()
      console.error('GitHub API错误详情:', errorData)
      throw new Error(`GitHub上传失败: ${errorData.message}`)
    }

    const result = await response.json()
    console.log('GitHub上传成功:', result)

    // 重要：GitHub Pages仓库的特殊URL格式
    // 对于 snow04c.github.io 仓库，访问URL是：https://snow04c.github.io/avatars/文件名
    // 不需要在URL中包含仓库名，因为这是用户页面仓库
    return `https://${GITHUB_CONFIG.username}.github.io/${GITHUB_CONFIG.folder}/${fileName}`

  } catch (error) {
    console.error('GitHub上传错误:', error)
    throw error
  }
}

// 文件转Base64的工具函数
function fileToBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (e) => resolve(e.target.result)
    reader.onerror = (error) => reject(error)
    reader.readAsDataURL(file)
  })
}

// 修改后的头像上传函数 - 添加更多错误处理和调试信息
async function handleAvatarUpload(event) {
  const file = event.target.files[0]
  if (!file) return

  // 验证文件类型和大小
  if (!file.type.startsWith('image/')) {
    alert('请选择图片文件（JPEG、PNG、GIF等）')
    return
  }

  if (file.size > 2 * 1024 * 1024) {
    alert('图片大小不能超过2MB')
    return
  }

  loading.value = true

  try {
    // 显示本地预览
    const previewUrl = URL.createObjectURL(file)
    console.log('本地预览URL:', previewUrl)

    // 1. 上传到GitHub获取公网URL
    console.log('开始上传头像到GitHub...')
    const githubAvatarUrl = await uploadToGitHub(file)
    console.log('GitHub头像URL:', githubAvatarUrl)

    // 2. 使用GitHub URL更新头像到后端
    const result = await updateAvatar(githubAvatarUrl)

    if (result.success) {
      // 更新本地用户数据
      await initUserData()
      alert('头像更新成功！待系统审核')

      // 释放预览URL
      URL.revokeObjectURL(previewUrl)
    } else {
      alert(result.message || '头像更新失败')
    }

  } catch (error) {
    console.error('头像上传失败:', error)

    // 详细的错误信息
    let errorMessage = '头像上传失败'
    if (error.message.includes('GitHub上传失败')) {
      if (error.message.includes('bad credentials')) {
        errorMessage = 'GitHub Token无效，请检查token权限'
      } else if (error.message.includes('not found')) {
        errorMessage = 'GitHub仓库不存在或无权访问'
      } else {
        errorMessage = `GitHub上传失败: ${error.message}`
      }
    } else if (error.message.includes('Network Error')) {
      errorMessage = '网络连接失败，请检查网络设置'
    }

    alert(errorMessage)
  } finally {
    loading.value = false
    // 清空文件输入
    event.target.value = ''
  }
}

// 更新密码
async function updatePassword(passwordData) {
  // 从本地浏览器获取token
  let token = localStorage.getItem('token')

  // 如果token是JSON字符串，解析它
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
      return { success: false, message: '令牌格式错误' }
    }
  }

  if (!token) {
    console.error('未找到认证令牌')
    return { success: false, message: '未找到认证令牌' }
  }

  try {
    const requestData = {
      email: currentUser.value?.email,
      code: passwordData.verificationCode,
      new_password: passwordData.newPassword
    }

    console.log('密码修改请求数据:', JSON.stringify(requestData, null, 2))

    const response = await fetch(`${API_BASE}/auth/change_password`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(requestData)
    })

    console.log('密码修改响应状态:', response.status)

    let result
    const responseText = await response.text()

    if (responseText) {
      try {
        result = JSON.parse(responseText)
        console.log('更新密码响应:', result)
        // eslint-disable-next-line no-unused-vars
      } catch (parseError) {
        console.log('密码修改响应不是JSON格式:', responseText)
        result = { rawResponse: responseText }
      }
    } else {
      console.log('密码修改响应为空')
      result = {}
    }

    if (response.ok) {
      return { success: true, data: result }
    } else {
      const errorMessage = result.message || result.error || '密码修改失败'
      console.error('密码修改失败详情:', result)
      return {
        success: false,
        message: errorMessage,
        status: response.status
      }
    }
  } catch (error) {
    console.error('修改密码请求失败:', error)
    return { success: false, message: '网络请求失败' }
  }
}

// 更新邮箱
async function updateEmail(emailData) {
  // 从本地浏览器获取token
  let token = localStorage.getItem('token')

  // 如果token是JSON字符串，解析它
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
      return { success: false, message: '令牌格式错误' }
    }
  }

  if (!token) {
    console.error('未找到认证令牌')
    return { success: false, message: '未找到认证令牌' }
  }

  try {
    const requestData = {
      new_email: emailData.newEmail,
      old_email_code: emailData.oldVerificationCode,
      new_email_code: emailData.newVerificationCode
    }

    console.log('邮箱修改请求数据:', JSON.stringify(requestData, null, 2))

    const response = await fetch(`${API_BASE}/auth/change_email`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(requestData)
    })

    console.log('邮箱修改响应状态:', response.status)

    let result
    const responseText = await response.text()

    if (responseText) {
      try {
        result = JSON.parse(responseText)
        console.log('更新邮箱响应:', result)
        // eslint-disable-next-line no-unused-vars
      } catch (jsonError) {
        console.log('邮箱修改响应不是JSON格式:', responseText)
        result = { rawResponse: responseText }
      }
    } else {
      console.log('邮箱修改响应为空')
      result = {}
    }

    if (response.ok) {
      // 邮箱修改成功，后端返回新的令牌
      let newToken = ''

      // 处理不同类型的响应格式
      if (typeof result === 'string') {
        // 如果直接返回token字符串
        newToken = result
      } else if (result.token) {
        // 如果返回的是对象且包含token字段
        newToken = result.token
      } else if (result.data && result.data.token) {
        // 如果返回的是 { data: { token: ... } } 格式
        newToken = result.data.token
      } else {
        // 其他格式，使用原始响应
        newToken = responseText
      }

      console.log('获取到的新令牌:', newToken)

      // 保存新令牌到本地存储
      localStorage.setItem('token', newToken)
      console.log('新令牌已保存到localStorage')

      return {
        success: true,
        data: result,
        newToken: newToken
      }
    } else {
      const errorMessage = result.message || result.error || '邮箱修改失败'
      console.error('邮箱修改失败详情:', result)
      return {
        success: false,
        message: errorMessage,
        status: response.status
      }
    }
  } catch (error) {
    console.error('修改邮箱请求失败:', error)
    return { success: false, message: '网络请求失败' }
  }
}

// 掩码显示函数
function maskUserId(userId) {
  if (!userId) return '****'
  const str = userId.toString()
  if (str.length <= 4) return str
  return str.slice(0, -4) + '****'
}

function maskEmail(email) {
  if (!email) return '***@***.***'
  const [name, domain] = email.split('@')
  if (!name || !domain) return '***@***.***'

  if (name.length <= 6) {
    return '*'.repeat(name.length) + '@' + domain
  }

  return '*'.repeat(6) + name.slice(6) + '@' + domain
}

function close() {
  emit('close')
}

async function saveProfile() {
  if (!userForm.value.username.trim()) {
    alert('用户名不能为空')
    return
  }

  loading.value = true

  try {
    const result = await updateUsername(userForm.value.username)

    if (result.success) {
      await initUserData()
      alert('用户名更新成功！')
    } else {
      alert(result.message)
    }

  } catch (error) {
    alert('保存失败：' + error.message)
  } finally {
    loading.value = false
  }
}

async function changePassword() {
  if (!passwordForm.value.newPassword) {
    alert('请输入新密码')
    return
  }

  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    alert('两次输入的新密码不一致')
    return
  }

  if (!passwordForm.value.verificationCode) {
    alert('请输入邮箱验证码')
    return
  }

  loading.value = true

  try {
    const result = await updatePassword(passwordForm.value)

    if (result.success) {
      // 重置表单
      passwordForm.value = {
        newPassword: '',
        confirmPassword: '',
        verificationCode: ''
      }

      alert('密码修改成功！请重新登录')

      // 清除本地存储的登录状态
      localStorage.removeItem('token')
      localStorage.removeItem('currentUser')
      currentUser.value = null

      emit('logout')
      close()

      // 跳转到首页
      router.push('/')

    } else {
      alert(result.message)
    }

  } catch (error) {
    alert('密码修改失败：' + error.message)
  } finally {
    loading.value = false
  }
}

async function changeEmail() {
  if (!emailForm.value.newEmail) {
    alert('请输入新邮箱')
    return
  }

  if (!emailForm.value.oldVerificationCode || !emailForm.value.newVerificationCode) {
    alert('请输入完整的验证码')
    return
  }

  loading.value = true

  try {
    const result = await updateEmail(emailForm.value)

    if (result.success) {
      // 重置表单
      emailForm.value = {
        newEmail: '',
        oldVerificationCode: '',
        newVerificationCode: ''
      }

      // 如果有新令牌，重新获取用户信息
      if (result.newToken) {
        console.log('使用新令牌重新获取用户信息...')
        await initUserData()
      }

      alert('邮箱修改成功！')
      activeMenu.value = 'profile'
    } else {
      alert(result.message)
    }

  } catch (error) {
    alert('邮箱修改失败：' + error.message)
  } finally {
    loading.value = false
  }
}

// 发送验证码函数
async function sendVerificationCode(email, business, cooldownRef) {
  try {
    console.log(`发送验证码: email=${email}, business=${business}`)

    const response = await fetch(`${API_BASE}/email/send-verification`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: email,
        business: business
      })
    })

    const data = await response.json()

    if (response.ok) {
      alert('验证码已发送到您的邮箱，请查收')

      cooldownRef.value = 60
      const timer = setInterval(() => {
        cooldownRef.value--
        if (cooldownRef.value <= 0) {
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

async function sendPasswordVerificationCode() {
  if (!currentUser.value?.email) {
    alert('未找到邮箱信息')
    return
  }
  // 修改密码使用 'change_password' 业务
  await sendVerificationCode(currentUser.value.email, 'change_password', passwordCooldown)
}

async function sendOldEmailVerificationCode() {
  if (!currentUser.value?.email) {
    alert('未找到邮箱信息')
    return
  }
  // 原邮箱验证使用 'change_email_old' 业务
  await sendVerificationCode(currentUser.value.email, 'change_email_old', oldEmailCooldown)
}

async function sendNewEmailVerificationCode() {
  if (!emailForm.value.newEmail) {
    alert('请输入新邮箱地址')
    return
  }

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(emailForm.value.newEmail)) {
    alert('请输入有效的邮箱地址')
    return
  }

  // 新邮箱验证使用 'change_email_new' 业务
  await sendVerificationCode(emailForm.value.newEmail, 'change_email_new', newEmailCooldown)
}

// 初始化用户数据
async function initUserData() {
  // 检查本地是否有token
  const token = localStorage.getItem('token')
  if (token) {
    console.log('找到token，开始获取用户信息...')
    const userData = await fetchCurrentUser()
    if (userData) {
      currentUser.value = userData
      userForm.value.username = userData.username || ''
      console.log('初始化用户数据成功:', userData)
      emit('update-user', userData)
    } else {
      console.error('获取用户数据失败')
    }
  } else {
    console.error('未找到token，用户未登录')
  }
}

// 监听模态框显示状态
watch(() => props.isVisible, async (newVal) => {
  console.log('个人信息模态框状态变化:', newVal)
  if (newVal) {
    await initUserData()

    passwordForm.value = {
      newPassword: '',
      confirmPassword: '',
      verificationCode: ''
    }
    emailForm.value = {
      newEmail: '',
      oldVerificationCode: '',
      newVerificationCode: ''
    }
    passwordCooldown.value = 0
    oldEmailCooldown.value = 0
    newEmailCooldown.value = 0
    activeMenu.value = 'profile'
  }
})

// 组件挂载时初始化
onMounted(() => {
  if (props.isVisible) {
    initUserData()
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
  width: 750px; /* 缩小宽度 */
  max-width: 90vw;
  max-height: 85vh;
  height: auto;
  min-height: 500px;
  overflow: hidden;
  color: #333;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px; /* 缩小内边距 */
  border-bottom: 1px solid #f0f0f0;
  background: white;
  z-index: 1;
}

.modal-header h3 {
  margin: 0;
  color: #333;
  font-size: 16px; /* 缩小字体 */
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  font-size: 20px; /* 缩小字体 */
  cursor: pointer;
  color: #666;
  padding: 0;
  width: 28px; /* 缩小尺寸 */
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background-color 0.2s;
}

.close-btn:hover {
  background: #f5f5f5;
}

.modal-body {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.profile-layout {
  display: flex;
  width: 100%;
  height: 100%;
}

/* 左侧导航 */
.sidebar {
  width: 180px; /* 缩小宽度 */
  background: #f8f9fa;
  border-right: 1px solid #e9ecef;
  padding: 16px 0; /* 缩小内边距 */
  overflow-y: auto;
  max-height: 100%;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 12px 16px; /* 缩小内边距 */
  cursor: pointer;
  transition: all 0.3s;
  border-left: 3px solid transparent;
  margin: 2px 0; /* 缩小间距 */
}

.menu-item:hover {
  background: #e9ecef;
}

.menu-item.active {
  background: #e3f2fd;
  border-left-color: #2196f3;
  color: #2196f3;
}

.menu-icon {
  font-size: 16px; /* 缩小图标 */
  margin-right: 10px; /* 缩小间距 */
  width: 20px;
  text-align: center;
}

.menu-text {
  font-size: 13px; /* 缩小字体 */
  font-weight: 500;
}

/* 右侧内容 */
.content-area {
  flex: 1;
  padding: 20px; /* 缩小内边距 */
  overflow-y: auto;
  background: #fff;
  max-height: 100%;
  display: flex;
  flex-direction: column;
}

.content-section {
  max-width: 100%;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.section-title {
  margin: 0 0 20px 0; /* 缩小间距 */
  font-size: 16px; /* 缩小字体 */
  font-weight: 600;
  color: #333;
  padding-bottom: 10px; /* 缩小间距 */
  border-bottom: 2px solid #f0f0f0;
}

/* 头像区域 */
.avatar-section {
  display: flex;
  justify-content: center;
  margin-bottom: 20px; /* 缩小间距 */
}

.avatar-container {
  text-align: center;
}

.avatar-wrapper {
  position: relative;
  display: inline-block;
  border-radius: 50%;
  overflow: hidden;
}

.avatar {
  width: 120px; /* 缩小头像 */
  height: 120px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #e9ecef; /* 缩小边框 */
  transition: all 0.3s;
}

.avatar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s;
}

.avatar-wrapper:hover .avatar-overlay {
  opacity: 1;
}

.avatar-edit-btn {
  background: rgba(255, 255, 255, 0.9);
  color: #333;
  border: none;
  border-radius: 16px; /* 缩小圆角 */
  padding: 4px 8px; /* 缩小内边距 */
  font-size: 10px; /* 缩小字体 */
  cursor: pointer;
  transition: background 0.3s;
}

.avatar-edit-btn:hover {
  background: white;
}

/* 信息网格 */
.info-grid {
  display: flex;
  flex-direction: column;
  gap: 16px; /* 缩小间距 */
  margin-bottom: 20px; /* 缩小间距 */
  flex: 1;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 6px; /* 缩小间距 */
}

.info-label {
  font-size: 13px; /* 缩小字体 */
  font-weight: 500;
  color: #666;
}

.info-value {
  font-size: 14px; /* 缩小字体 */
  color: #333;
  padding: 6px 0; /* 缩小内边距 */
}

.info-value.masked {
  color: #888;
  font-family: 'Courier New', monospace;
}

.info-input {
  padding: 10px; /* 缩小内边距 */
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  font-size: 13px; /* 缩小字体 */
  background: white;
  transition: border 0.3s;
}

.info-input:focus {
  outline: none;
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1);
}

.masked-email {
  padding: 10px; /* 缩小内边距 */
  background: #f5f5f5;
  border-radius: 6px;
  color: #666;
  font-family: 'Courier New', monospace;
  font-size: 13px; /* 缩小字体 */
}

.masked-password {
  padding: 10px; /* 缩小内边距 */
  background: #f5f5f5;
  border-radius: 6px;
  color: #666;
  font-family: 'Courier New', monospace;
  letter-spacing: 2px;
  font-size: 13px; /* 缩小字体 */
}

/* 安全通知 */
.security-notice {
  display: flex;
  align-items: center;
  gap: 10px; /* 缩小间距 */
  background: #f0f7ff;
  border: 1px solid #d0e3ff;
  border-radius: 6px; /* 缩小圆角 */
  padding: 12px; /* 缩小内边距 */
  margin-bottom: 20px; /* 缩小间距 */
}

.notice-icon {
  font-size: 18px; /* 缩小图标 */
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 20px; /* 缩小高度 */
}

.notice-content p {
  margin: 0;
  font-size: 13px; /* 缩小字体 */
  color: #333;
  line-height: 1.4;
}

/* 表单样式 */
.security-form {
  margin-top: 16px; /* 缩小间距 */
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.form-group {
  margin-bottom: 16px; /* 缩小间距 */
}

.form-label {
  display: block;
  margin-bottom: 6px; /* 缩小间距 */
  font-weight: 500;
  color: #333;
  font-size: 13px; /* 缩小字体 */
}

.form-input {
  width: 100%;
  padding: 10px; /* 缩小内边距 */
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  font-size: 13px; /* 缩小字体 */
  box-sizing: border-box;
  transition: border 0.3s;
}

.form-input:focus {
  outline: none;
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1);
}

.verification-group {
  display: flex;
  gap: 8px; /* 缩小间距 */
}

.verification-group .form-input {
  flex: 1;
}

.verification-btn {
  padding: 10px 12px; /* 缩小内边距 */
  background: #409eff;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  white-space: nowrap;
  font-size: 12px; /* 缩小字体 */
  transition: background 0.3s;
  min-width: 90px; /* 缩小宽度 */
}

.verification-btn:hover:not(:disabled) {
  background: #66b1ff;
}

.verification-btn:disabled {
  background: #c0c4cc;
  cursor: not-allowed;
}

/* 按钮样式 */
.profile-actions,
.form-actions {
  margin-top: auto;
  padding-top: 16px; /* 缩小间距 */
  border-top: 1px solid #f0f0f0;
  flex-shrink: 0;
}

.action-btn {
  padding: 10px 20px; /* 缩小内边距 */
  border: none;
  border-radius: 6px;
  font-size: 13px; /* 缩小字体 */
  cursor: pointer;
  transition: all 0.3s;
  font-weight: 500;
}

.action-btn.primary {
  background: #409eff;
  color: white;
  width: 100%;
}

.action-btn.primary:hover:not(:disabled) {
  background: #66b1ff;
  transform: translateY(-1px);
}

.action-btn.primary:disabled {
  background: #c0c4cc;
  cursor: not-allowed;
  transform: none;
}

.action-btn.secondary {
  background: #f5f5f5;
  color: #333;
  border: 1px solid #e0e0e0;
}

.action-btn.secondary:hover {
  background: #e8e8e8;
}

.action-btn.danger {
  background: #f56c6c;
  color: white;
}

.action-btn.danger:hover {
  background: #f78989;
}

/* 退出登录样式 */
.logout-content {
  text-align: center;
  padding: 30px 16px; /* 缩小内边距 */
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.logout-icon {
  font-size: 48px; /* 缩小图标 */
  margin-bottom: 16px; /* 缩小间距 */
  opacity: 0.7;
}

.logout-title {
  font-size: 16px; /* 缩小字体 */
  font-weight: 600;
  color: #333;
  margin-bottom: 10px; /* 缩小间距 */
}

.logout-text {
  font-size: 14px; /* 缩小字体 */
  color: #333;
  margin-bottom: 6px; /* 缩小间距 */
}

.logout-desc {
  font-size: 13px; /* 缩小字体 */
  color: #666;
  margin-bottom: 24px; /* 缩小间距 */
}

.logout-actions {
  display: flex;
  gap: 10px; /* 缩小间距 */
  justify-content: center;
}

.logout-actions .action-btn {
  min-width: 80px; /* 缩小宽度 */
}

/* 模态框动画 */
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

/* 响应式设计 */
@media (max-width: 768px) {
  .modal-container {
    width: 95vw;
    max-height: 90vh;
    min-height: 450px; /* 移动端也缩小 */
  }

  .profile-layout {
    flex-direction: column;
  }

  .sidebar {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid #e9ecef;
    padding: 10px 0;
    max-height: none;
    min-height: auto;
  }

  .menu-item {
    padding: 10px 16px;
    display: inline-flex;
    margin: 0 4px;
    border-left: none;
    border-bottom: 2px solid transparent;
  }

  .menu-item.active {
    border-left: none;
    border-bottom-color: #2196f3;
  }

  .content-area {
    padding: 16px;
    max-height: none;
    min-height: auto;
  }

  .content-section {
    max-width: none;
  }

  .verification-group {
    flex-direction: column;
  }

  .logout-actions {
    flex-direction: column;
  }

  .logout-actions .action-btn {
    width: 100%;
  }
}

/* 滚动条样式 */
.sidebar::-webkit-scrollbar,
.content-area::-webkit-scrollbar {
  width: 4px; /* 缩小滚动条 */
}

.sidebar::-webkit-scrollbar-track,
.content-area::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 2px; /* 缩小圆角 */
}

.sidebar::-webkit-scrollbar-thumb,
.content-area::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 2px; /* 缩小圆角 */
}

.sidebar::-webkit-scrollbar-thumb:hover,
.content-area::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>
