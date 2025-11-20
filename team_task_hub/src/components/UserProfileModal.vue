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
                      <img :src="userForm.avatar || '/空白头像.png'" alt="头像" class="avatar">
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
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const props = defineProps({
  isVisible: Boolean,
  user: Object
})

const emit = defineEmits(['close', 'update-user', 'logout'])

const loading = ref(false)
const passwordCooldown = ref(0)
const oldEmailCooldown = ref(0)
const newEmailCooldown = ref(0)
const avatarInput = ref(null)
const activeMenu = ref('profile')

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
  username: '',
  avatar: ''
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

// 掩码显示函数
function maskUserId(userId) {
  if (!userId) return '****'
  const str = userId.toString()
  if (str.length <= 4) return str + '****'
  return str.slice(0, 4) + '****'
}

function maskEmail(email) {
  if (!email) return '***@***.***'
  const [name, domain] = email.split('@')
  if (!name || !domain) return '***@***.***'

  const maskedName = name.length > 2
    ? name.charAt(0) + '*'.repeat(Math.min(3, name.length - 2)) + name.charAt(name.length - 1)
    : '*'.repeat(name.length)

  return `${maskedName}@${domain}`
}

function close() {
  emit('close')
}

function handleAvatarUpload(event) {
  const file = event.target.files[0]
  if (file) {
    // 验证文件类型和大小
    if (!file.type.startsWith('image/')) {
      alert('请选择图片文件')
      return
    }
    if (file.size > 2 * 1024 * 1024) { // 2MB
      alert('图片大小不能超过2MB')
      return
    }

    const reader = new FileReader()
    reader.onload = (e) => {
      userForm.value.avatar = e.target.result
    }
    reader.readAsDataURL(file)
  }
}

async function saveProfile() {
  if (!userForm.value.username.trim()) {
    alert('用户名不能为空')
    return
  }

  loading.value = true

  try {
    await new Promise(resolve => setTimeout(resolve, 1000))

    // 模拟保存用户信息
    const updatedUser = {
      ...props.user,
      username: userForm.value.username,
      avatar: userForm.value.avatar
    }

    // 更新本地存储
    localStorage.setItem('currentUser', JSON.stringify(updatedUser))

    // 更新用户列表中的信息
    const users = JSON.parse(localStorage.getItem('users') || '[]')
    const userIndex = users.findIndex(u => u.id === props.user.id)
    if (userIndex !== -1) {
      users[userIndex] = updatedUser
      localStorage.setItem('users', JSON.stringify(users))
    }

    emit('update-user', updatedUser)
    alert('个人信息更新成功！')

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
    await new Promise(resolve => setTimeout(resolve, 1000))

    // 更新密码
    const updatedUser = {
      ...props.user,
      password: passwordForm.value.newPassword
    }

    // 更新本地存储
    localStorage.setItem('currentUser', JSON.stringify(updatedUser))

    // 更新用户列表中的密码
    const users = JSON.parse(localStorage.getItem('users') || '[]')
    const userIndex = users.findIndex(u => u.id === props.user.id)
    if (userIndex !== -1) {
      users[userIndex].password = passwordForm.value.newPassword
      localStorage.setItem('users', JSON.stringify(users))
    }

    // 重置表单
    passwordForm.value = {
      newPassword: '',
      confirmPassword: '',
      verificationCode: ''
    }

    alert('密码修改成功！')
    activeMenu.value = 'profile'

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

  // 验证验证码（模拟）
  const storedCode = '123456' // 模拟验证码
  if (emailForm.value.oldVerificationCode !== storedCode || emailForm.value.newVerificationCode !== storedCode) {
    alert('验证码错误')
    return
  }

  loading.value = true

  try {
    await new Promise(resolve => setTimeout(resolve, 1000))

    // 更新邮箱
    const updatedUser = {
      ...props.user,
      email: emailForm.value.newEmail
    }

    // 更新本地存储
    localStorage.setItem('currentUser', JSON.stringify(updatedUser))

    // 更新用户列表中的邮箱
    const users = JSON.parse(localStorage.getItem('users') || '[]')
    const userIndex = users.findIndex(u => u.id === props.user.id)
    if (userIndex !== -1) {
      users[userIndex].email = emailForm.value.newEmail
      localStorage.setItem('users', JSON.stringify(users))
    }

    // 重置表单
    emailForm.value = {
      newEmail: '',
      oldVerificationCode: '',
      newVerificationCode: ''
    }

    emit('update-user', updatedUser)
    alert('邮箱修改成功！')
    activeMenu.value = 'profile'

  } catch (error) {
    alert('邮箱修改失败：' + error.message)
  } finally {
    loading.value = false
  }
}

// 发送验证码函数
async function sendVerificationCode(email, business, cooldownRef) {
  try {
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

      // 开始倒计时
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
  if (!props.user?.email) {
    alert('未找到邮箱信息')
    return
  }
  await sendVerificationCode(props.user.email, 'reset_password', passwordCooldown)
}

async function sendOldEmailVerificationCode() {
  if (!props.user?.email) {
    alert('未找到邮箱信息')
    return
  }
  await sendVerificationCode(props.user.email, 'change_email', oldEmailCooldown)
}

async function sendNewEmailVerificationCode() {
  if (!emailForm.value.newEmail) {
    alert('请输入新邮箱地址')
    return
  }

  // 邮箱格式验证
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(emailForm.value.newEmail)) {
    alert('请输入有效的邮箱地址')
    return
  }

  await sendVerificationCode(emailForm.value.newEmail, 'change_email', newEmailCooldown)
}

function handleLogout() {
  // 清除登录状态
  localStorage.removeItem('currentUser')
  localStorage.removeItem('token')
  emit('logout')
  close()

  // 跳转到首页
  router.push('/')
  alert('已退出登录')
}

// 监听用户数据变化
watch(() => props.user, (newUser) => {
  if (newUser) {
    userForm.value = {
      username: newUser.username || '',
      avatar: newUser.avatar || '/空白头像.png'
    }
  }
}, { immediate: true })

// 监听模态框显示状态
watch(() => props.isVisible, (newVal) => {
  if (newVal && props.user) {
    userForm.value = {
      username: props.user.username || '',
      avatar: props.user.avatar || '/空白头像.png'
    }
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
