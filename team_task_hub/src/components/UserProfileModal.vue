<template>
  <transition name="modal">
    <div v-if="isVisible" class="modal-overlay" @click.self="close">
      <div class="modal-container">
        <div class="modal-header">
          <h3>个人信息</h3>
          <button class="close-btn" @click="close">×</button>
        </div>

        <div class="modal-body">
          <!-- 头像编辑 -->
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

          <!-- 用户信息表单 -->
          <form @submit.prevent="saveProfile">
            <div class="form-group">
              <label>用户名</label>
              <input
                type="text"
                v-model="userForm.username"
                required
              >
            </div>

            <div class="form-group">
              <label>邮箱</label>
              <div class="email-group">
                <input
                  type="email"
                  v-model="userForm.email"
                  required
                >
                <button
                  type="button"
                  class="verification-btn"
                  :disabled="emailCooldown > 0"
                  @click="sendEmailVerification"
                >
                  {{ emailCooldown > 0 ? `${emailCooldown}s` : '验证' }}
                </button>
              </div>
            </div>

            <div class="form-group" v-if="showEmailVerification">
              <label>邮箱验证码</label>
              <input
                type="text"
                v-model="userForm.emailVerificationCode"
                placeholder="请输入验证码"
              >
            </div>

            <div class="form-group">
              <label>新密码（留空不修改）</label>
              <input
                type="password"
                v-model="userForm.newPassword"
                placeholder="新密码"
              >
            </div>

            <div class="form-group" v-if="userForm.newPassword">
              <label>确认新密码</label>
              <input
                type="password"
                v-model="userForm.confirmPassword"
                placeholder="确认新密码"
              >
            </div>

            <div class="form-actions">
              <button type="submit" class="save-btn" :disabled="loading">
                {{ loading ? '保存中...' : '保存修改' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  isVisible: Boolean,
  user: Object
})

const emit = defineEmits(['close', 'update-user'])

const loading = ref(false)
const emailCooldown = ref(0)
const showEmailVerification = ref(false)
const avatarInput = ref(null)

const userForm = ref({
  username: '',
  email: '',
  avatar: '',
  newPassword: '',
  confirmPassword: '',
  emailVerificationCode: ''
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
  if (userForm.value.newPassword && userForm.value.newPassword !== userForm.value.confirmPassword) {
    alert('两次输入的密码不一致')
    return
  }

  if (showEmailVerification.value && !userForm.value.emailVerificationCode) {
    alert('请输入邮箱验证码')
    return
  }

  loading.value = true

  try {
    await new Promise(resolve => setTimeout(resolve, 1000))

    // 模拟保存用户信息
    const updatedUser = {
      ...props.user,
      username: userForm.value.username,
      email: userForm.value.email,
      avatar: userForm.value.avatar
    }

    if (userForm.value.newPassword) {
      updatedUser.password = userForm.value.newPassword
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
    close()
    alert('个人信息更新成功！')

  } catch (error) {
    alert('保存失败：' + error.message)
  } finally {
    loading.value = false
  }
}

async function sendEmailVerification() {
  if (!userForm.value.email) {
    alert('请输入邮箱')
    return
  }

  // 模拟发送验证码
  const code = Math.random().toString().slice(2, 8)
  console.log(`邮箱验证码已发送到 ${userForm.value.email}: ${code}`)

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

// 监听用户数据变化
watch(() => props.user, (newUser) => {
  if (newUser) {
    userForm.value = {
      username: newUser.username || '',
      email: newUser.email || '',
      avatar: newUser.avatar || '',
      newPassword: '',
      confirmPassword: '',
      emailVerificationCode: ''
    }
  }
}, { immediate: true })

// 监听模态框显示状态
watch(() => props.isVisible, (newVal) => {
  if (newVal && props.user) {
    userForm.value = {
      username: props.user.username || '',
      email: props.user.email || '',
      avatar: props.user.avatar || '',
      newPassword: '',
      confirmPassword: '',
      emailVerificationCode: ''
    }
    showEmailVerification.value = false
    emailCooldown.value = 0
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
  width: 500px;
  max-width: 90vw;
  max-height: 80vh;
  overflow-y: auto;
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

.avatar-section {
  display: flex;
  justify-content: center;
  margin-bottom: 24px;
}

.avatar-container {
  text-align: center;
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #eee;
}

.avatar-edit-btn {
  margin-top: 8px;
  padding: 6px 12px;
  background: #409eff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-weight: 500;
  color: #333;
}

.form-group input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  box-sizing: border-box;
}

.email-group {
  display: flex;
  gap: 10px;
}

.email-group input {
  flex: 1;
}

.verification-btn {
  padding: 10px 16px;
  background: #67c23a;
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

.form-actions {
  margin-top: 24px;
}

.save-btn {
  width: 100%;
  padding: 12px;
  background: #409eff;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 16px;
  cursor: pointer;
}

.save-btn:disabled {
  background: #c0c4cc;
  cursor: not-allowed;
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
