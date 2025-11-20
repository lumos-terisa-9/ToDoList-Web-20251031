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
                <h4 class="section-title">个人信息</h4>
                <div class="avatar-section">
                  <div class="avatar-container">
                    <img :src="userForm.avatar || '/default-avatar.png'" alt="头像" class="avatar">
                    <input
                      type="file"
                      ref="avatarInput"
                      accept="image/*"
                      @change="handleAvatarUpload"
                      style="display: none"
                    >
                    <button class="avatar-edit-btn" @click="$refs.avatarInput.click()">
                      更换头像
                    </button>
                  </div>
                </div>

                <form @submit.prevent="saveProfile" class="profile-form">
                  <div class="form-group">
                    <label>用户名</label>
                    <input
                      type="text"
                      v-model="userForm.username"
                      placeholder="请输入用户名"
                      required
                    >
                  </div>

                  <div class="form-actions">
                    <button type="submit" class="save-btn" :disabled="loading">
                      {{ loading ? '保存中...' : '保存修改' }}
                    </button>
                  </div>
                </form>
              </div>

              <!-- 密码修改 -->
              <div v-if="activeMenu === 'password'" class="content-section">
                <h4 class="section-title">修改密码</h4>
                <form @submit.prevent="changePassword" class="password-form">
                  <div class="form-group">
                    <label>当前密码</label>
                    <input
                      type="password"
                      v-model="passwordForm.currentPassword"
                      placeholder="请输入当前密码"
                      required
                    >
                  </div>

                  <div class="form-group">
                    <label>新密码</label>
                    <input
                      type="password"
                      v-model="passwordForm.newPassword"
                      placeholder="请输入新密码"
                      required
                    >
                  </div>

                  <div class="form-group">
                    <label>确认新密码</label>
                    <input
                      type="password"
                      v-model="passwordForm.confirmPassword"
                      placeholder="请再次输入新密码"
                      required
                    >
                  </div>

                  <div class="form-actions">
                    <button type="submit" class="save-btn" :disabled="loading">
                      {{ loading ? '修改中...' : '修改密码' }}
                    </button>
                  </div>
                </form>
              </div>

              <!-- 邮箱修改 -->
              <div v-if="activeMenu === 'email'" class="content-section">
                <h4 class="section-title">修改邮箱</h4>
                <form @submit.prevent="changeEmail" class="email-form">
                  <div class="form-group">
                    <label>当前邮箱</label>
                    <input
                      type="email"
                      :value="currentUser?.email || ''"
                      disabled
                      class="disabled-input"
                    >
                  </div>

                  <div class="form-group">
                    <label>新邮箱</label>
                    <div class="email-group">
                      <input
                        type="email"
                        v-model="emailForm.newEmail"
                        placeholder="请输入新邮箱地址"
                        required
                      >
                      <button
                        type="button"
                        class="verification-btn"
                        :disabled="emailCooldown > 0"
                        @click="sendEmailVerification"
                      >
                        {{ emailCooldown > 0 ? `${emailCooldown}s` : '发送验证码' }}
                      </button>
                    </div>
                  </div>

                  <div class="form-group" v-if="showEmailVerification">
                    <label>邮箱验证码</label>
                    <input
                      type="text"
                      v-model="emailForm.verificationCode"
                      placeholder="请输入验证码"
                      required
                    >
                  </div>

                  <div class="form-actions">
                    <button type="submit" class="save-btn" :disabled="loading">
                      {{ loading ? '修改中...' : '修改邮箱' }}
                    </button>
                  </div>
                </form>
              </div>

              <!-- 退出登录 -->
              <div v-if="activeMenu === 'logout'" class="content-section">
                <h4 class="section-title">退出登录</h4>
                <div class="logout-content">
                  <div class="logout-icon">🚪</div>
                  <p class="logout-text">确定要退出登录吗？</p>
                  <p class="logout-desc">退出后需要重新登录才能访问个人页面</p>
                  <button class="logout-btn" @click="handleLogout">
                    确认退出
                  </button>
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
const emailCooldown = ref(0)
const showEmailVerification = ref(false)
const avatarInput = ref(null)
const activeMenu = ref('profile')

// 菜单项配置
const menuItems = ref([
  { key: 'profile', text: '个人信息', icon: '👤' },
  { key: 'password', text: '密码', icon: '🔒' },
  { key: 'email', text: '邮箱', icon: '📧' },
  { key: 'logout', text: '退出登录', icon: '🚪' }
])

// 表单数据
const userForm = ref({
  username: '',
  avatar: ''
})

const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const emailForm = ref({
  newEmail: '',
  verificationCode: ''
})

function close() {
  emit('close')
}

function handleAvatarUpload(event) {
  const file = event.target.files[0]
  if (file) {
    const reader = new FileReader()
    reader.onload = (e) => {
      userForm.value.avatar = e.target.result
    }
    reader.readAsDataURL(file)
  }
}

async function saveProfile() {
  if (!userForm.value.username.trim()) {
    alert('请输入用户名')
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
    const userIndex = users.findIndex(u => u.username === props.user.username)
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
  if (!passwordForm.value.currentPassword) {
    alert('请输入当前密码')
    return
  }

  if (!passwordForm.value.newPassword) {
    alert('请输入新密码')
    return
  }

  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    alert('两次输入的新密码不一致')
    return
  }

  // 验证当前密码（模拟）
  if (passwordForm.value.currentPassword !== props.user.password) {
    alert('当前密码错误')
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
    const userIndex = users.findIndex(u => u.username === props.user.username)
    if (userIndex !== -1) {
      users[userIndex].password = passwordForm.value.newPassword
      localStorage.setItem('users', JSON.stringify(users))
    }

    // 重置表单
    passwordForm.value = {
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    }

    alert('密码修改成功！')

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

  if (!emailForm.value.verificationCode) {
    alert('请输入验证码')
    return
  }

  // 验证验证码（模拟）
  const storedCode = '123456' // 模拟验证码
  if (emailForm.value.verificationCode !== storedCode) {
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
    const userIndex = users.findIndex(u => u.username === props.user.username)
    if (userIndex !== -1) {
      users[userIndex].email = emailForm.value.newEmail
      localStorage.setItem('users', JSON.stringify(users))
    }

    // 重置表单
    emailForm.value = {
      newEmail: '',
      verificationCode: ''
    }
    showEmailVerification.value = false

    emit('update-user', updatedUser)
    alert('邮箱修改成功！')

  } catch (error) {
    alert('邮箱修改失败：' + error.message)
  } finally {
    loading.value = false
  }
}

async function sendEmailVerification() {
  if (!emailForm.value.newEmail) {
    alert('请输入新邮箱')
    return
  }

  // 模拟发送验证码
  console.log(`验证码已发送到 ${emailForm.value.newEmail}: 123456`)

  showEmailVerification.value = true

  // 开始倒计时
  emailCooldown.value = 60
  const timer = setInterval(() => {
    emailCooldown.value--
    if (emailCooldown.value <= 0) {
      clearInterval(timer)
    }
  }, 1000)

  alert('验证码已发送到您的邮箱（请在控制台查看）')
}

function handleLogout() {
  // 清除登录状态
  localStorage.removeItem('currentUser')
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
      avatar: newUser.avatar || ''
    }
  }
}, { immediate: true })

// 监听模态框显示状态
watch(() => props.isVisible, (newVal) => {
  if (newVal && props.user) {
    userForm.value = {
      username: props.user.username || '',
      avatar: props.user.avatar || ''
    }
    passwordForm.value = {
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    }
    emailForm.value = {
      newEmail: '',
      verificationCode: ''
    }
    showEmailVerification.value = false
    emailCooldown.value = 0
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
  width: 800px;
  max-width: 90vw;
  max-height: 85vh;
  height: auto;
  min-height: 500px;
  overflow: hidden;
  color: #333;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #eee;
  background: white;
  z-index: 1;
}

.modal-header h3 {
  margin: 0;
  color: #333;
  font-size: 18px;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
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
  width: 200px;
  background: #f8f9fa;
  border-right: 1px solid #e9ecef;
  padding: 20px 0;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 12px 20px;
  cursor: pointer;
  transition: all 0.3s;
  border-left: 3px solid transparent;
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
  font-size: 18px;
  margin-right: 12px;
  width: 24px;
  text-align: center;
}

.menu-text {
  font-size: 14px;
  font-weight: 500;
}

/* 右侧内容 */
.content-area {
  flex: 1;
  padding: 30px;
  overflow-y: auto;
}

.content-section {
  max-width: 400px;
}

.section-title {
  margin: 0 0 25px 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
  padding-bottom: 10px;
  border-bottom: 2px solid #f0f0f0;
}

/* 头像区域 */
.avatar-section {
  display: flex;
  justify-content: center;
  margin-bottom: 30px;
}

.avatar-container {
  text-align: center;
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid #e9ecef;
  margin-bottom: 12px;
}

.avatar-edit-btn {
  padding: 8px 16px;
  background: #409eff;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: background 0.3s;
}

.avatar-edit-btn:hover {
  background: #66b1ff;
}

/* 表单样式 */
.profile-form,
.password-form,
.email-form {
  margin-top: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #333;
  font-size: 14px;
}

.form-group input {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  box-sizing: border-box;
  transition: border 0.3s;
}

.form-group input:focus {
  outline: none;
  border-color: #409eff;
}

.disabled-input {
  background: #f5f5f5;
  color: #999;
  cursor: not-allowed;
}

.email-group {
  display: flex;
  gap: 10px;
}

.email-group input {
  flex: 1;
}

.verification-btn {
  padding: 12px 16px;
  background: #67c23a;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  white-space: nowrap;
  font-size: 14px;
  transition: background 0.3s;
  min-width: 100px;
}

.verification-btn:hover:not(:disabled) {
  background: #85ce61;
}

.verification-btn:disabled {
  background: #c0c4cc;
  cursor: not-allowed;
}

.form-actions {
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.save-btn {
  width: 100%;
  padding: 12px;
  background: #409eff;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: background 0.3s;
}

.save-btn:hover:not(:disabled) {
  background: #66b1ff;
}

.save-btn:disabled {
  background: #c0c4cc;
  cursor: not-allowed;
}

/* 退出登录样式 */
.logout-content {
  text-align: center;
  padding: 40px 20px;
}

.logout-icon {
  font-size: 48px;
  margin-bottom: 20px;
}

.logout-text {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.logout-desc {
  font-size: 14px;
  color: #666;
  margin-bottom: 30px;
}

.logout-btn {
  padding: 12px 40px;
  background: #f56c6c;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: background 0.3s;
}

.logout-btn:hover {
  background: #f78989;
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
    min-height: 400px;
  }

  .profile-layout {
    flex-direction: column;
  }

  .sidebar {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid #e9ecef;
    padding: 10px 0;
  }

  .menu-item {
    padding: 10px 20px;
  }

  .content-area {
    padding: 20px;
  }

  .content-section {
    max-width: none;
  }

  .email-group {
    flex-direction: column;
  }
}
</style>
