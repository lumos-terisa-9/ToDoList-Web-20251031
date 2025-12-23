<template>
  <div class="create-activity-dialog-overlay" @click="closeDialog">
    <div class="create-activity-dialog" @click.stop>
      <!-- 对话框头部 -->
      <div class="dialog-header">
        <h2 class="dialog-title">创建新活动</h2>
        <div class="dialog-subtitle">{{ orgName }} · 请输入活动信息</div>
        <button class="close-btn" @click="closeDialog">×</button>
      </div>

      <!-- 表单区域 -->
      <div class="dialog-body">
        <div class="form-container">
          <!-- 基本表单 -->
          <div class="form-section">
            <h3 class="section-title">活动基本信息</h3>

            <div class="form-group">
              <label class="form-label">
                <span class="label-text">活动标题</span>
                <span class="label-required">*</span>
              </label>
              <input
                v-model="formData.title"
                type="text"
                class="form-input"
                placeholder="请输入活动标题"
                :class="{ 'error': errors.title }"
              >
              <div v-if="errors.title" class="error-message">{{ errors.title }}</div>
            </div>

            <div class="form-group">
              <label class="form-label">
                <span class="label-text">活动描述</span>
              </label>
              <textarea
                v-model="formData.description"
                class="form-textarea"
                placeholder="请输入活动描述（可选）"
                rows="3"
              ></textarea>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label">
                  <span class="label-text">开始时间</span>
                  <span class="label-required">*</span>
                </label>
                <input
                  v-model="localStartTime"
                  type="datetime-local"
                  class="form-input"
                  :class="{ 'error': errors.start_time }"
                >
                <div v-if="errors.start_time" class="error-message">{{ errors.start_time }}</div>
              </div>

              <div class="form-group">
                <label class="form-label">
                  <span class="label-text">结束时间</span>
                  <span class="label-required">*</span>
                </label>
                <input
                  v-model="localEndTime"
                  type="datetime-local"
                  class="form-input"
                  :class="{ 'error': errors.end_time }"
                >
                <div v-if="errors.end_time" class="error-message">{{ errors.end_time }}</div>
              </div>
            </div>

            <div class="form-group">
              <label class="form-label">
                <span class="label-text">参与限制</span>
                <span class="label-required">*</span>
              </label>
              <div class="radio-group">
                <label class="radio-option" :class="{ 'selected': formData.participation_limit === 'public' }">
                  <input
                    v-model="formData.participation_limit"
                    type="radio"
                    value="public"
                    class="radio-input"
                  >
                  <span class="radio-label">
                    <span class="radio-title">公开活动</span>
                    <span class="radio-description">所有用户可见并可参与</span>
                  </span>
                </label>

                <label class="radio-option" :class="{ 'selected': formData.participation_limit === 'org_only' }">
                  <input
                    v-model="formData.participation_limit"
                    type="radio"
                    value="org_only"
                    class="radio-input"
                  >
                  <span class="radio-label">
                    <span class="radio-title">组织内部活动</span>
                    <span class="radio-description">仅组织成员可见并可参与</span>
                  </span>
                </label>

                <label class="radio-option" :class="{ 'selected': formData.participation_limit === 'admin_assign' }">
                  <input
                    v-model="formData.participation_limit"
                    type="radio"
                    value="admin_assign"
                    class="radio-input"
                  >
                  <span class="radio-label">
                    <span class="radio-title">专项活动</span>
                    <span class="radio-description">管理员指定参与人员</span>
                  </span>
                </label>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 对话框底部 -->
      <div class="dialog-footer">
        <button class="btn-cancel" @click="closeDialog">取消</button>
        <button
          class="btn-submit"
          @click="submitForm"
          :disabled="submitting"
        >
          <span v-if="submitting">创建中...</span>
          <span v-else>创建活动</span>
          <div v-if="!submitting" class="submit-sparkle">
            <div class="spark"></div>
            <div class="spark"></div>
            <div class="spark"></div>
          </div>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch, onMounted } from 'vue'

const props = defineProps({
  orgId: {
    type: [String, Number],
    required: true
  },
  orgName: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['close', 'activity-created'])

// 表单数据
const formData = reactive({
  title: '',
  description: '',
  start_time: '',
  end_time: '',
  participation_limit: 'public'
})

// 本地时间输入（datetime-local格式）
const localStartTime = ref('')
const localEndTime = ref('')

// 验证错误
const errors = reactive({
  title: '',
  start_time: '',
  end_time: ''
})

// 提交状态
const submitting = ref(false)

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

// 时间格式转换函数
const formatToFullISO = (datetimeLocalStr) => {
  if (!datetimeLocalStr) return null;

  // 如果已经是完整格式（包含时区），直接返回
  if (datetimeLocalStr.includes('+') || datetimeLocalStr.includes('Z')) {
    return datetimeLocalStr;
  }

  // 如果是 datetime-local 格式 (YYYY-MM-DDTHH:MM)，添加秒和时区
  const date = new Date(datetimeLocalStr);
  return date.toISOString(); // 这会转换为 UTC 时间，如 2025-11-28T17:00:00.000Z
};

// 从datetime-local转换为API格式
function convertDateTimeLocalToAPI(localDateTime) {
  return formatToFullISO(localDateTime)
}

// 从API格式转换为datetime-local输入值
function convertAPIToDateTimeLocal(apiDateTime) {
  if (!apiDateTime) return ''

  try {
    const date = new Date(apiDateTime)
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')

    return `${year}-${month}-${day}T${hours}:${minutes}`
  } catch (error) {
    console.error('转换时间格式失败:', error)
    return ''
  }
}

// 验证时间格式
function validateAPITimeFormat(timeStr) {
  if (!timeStr) return false

  // 检查是否为有效的ISO格式
  try {
    const date = new Date(timeStr)
    if (isNaN(date.getTime())) return false

    // 确保年份有效
    const year = date.getUTCFullYear()
    return year > 1000 && year < 9999
    // eslint-disable-next-line no-unused-vars
  } catch (error) {
    return false
  }
}

// 设置默认时间
function setDefaultTimes() {
  const today = new Date()
  const tomorrow = new Date(today)
  tomorrow.setDate(tomorrow.getDate() + 1)

  // 设置本地时间输入
  localStartTime.value = convertAPIToDateTimeLocal(today.toISOString())
  localEndTime.value = convertAPIToDateTimeLocal(tomorrow.toISOString())

  // 同步到API格式
  formData.start_time = convertDateTimeLocalToAPI(localStartTime.value)
  formData.end_time = convertDateTimeLocalToAPI(localEndTime.value)

  console.log('设置的默认时间:', {
    start: formData.start_time,
    end: formData.end_time
  })
}

// 表单验证
function validateForm() {
  let isValid = true

  // 清空错误信息
  Object.keys(errors).forEach(key => errors[key] = '')

  // 验证标题
  if (!formData.title.trim()) {
    errors.title = '活动标题不能为空'
    isValid = false
  }

  // 验证开始时间
  if (!formData.start_time || !validateAPITimeFormat(formData.start_time)) {
    errors.start_time = '请选择有效的开始时间'
    isValid = false
  }

  // 验证结束时间
  if (!formData.end_time || !validateAPITimeFormat(formData.end_time)) {
    errors.end_time = '请选择有效的结束时间'
    isValid = false
  }

  // 验证时间逻辑
  if (isValid) {
    const start = new Date(formData.start_time)
    const end = new Date(formData.end_time)

    if (start >= end) {
      errors.end_time = '结束时间必须晚于开始时间'
      isValid = false
    }
  }

  return isValid
}

// 提交表单
async function submitForm() {
  if (!validateForm()) {
    return
  }

  const token = getToken()
  if (!token) {
    alert('请先登录')
    return
  }

  submitting.value = true

  try {
    // 构建活动数据，严格按照API要求的格式
    const activityData = {
      title: formData.title,
      description: formData.description || '',
      start_time: formData.start_time,  // 已经是正确格式
      end_time: formData.end_time,      // 已经是正确格式
      participation_limit: formData.participation_limit
    }

    console.log('提交的活动数据:', activityData)

    // 1. 创建活动
    const createResponse = await fetch(`http://localhost:8080/api/organization/${props.orgId}/activities`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(activityData)
    })

    const responseData = await createResponse.json()

    if (!createResponse.ok || !responseData.success) {
      throw new Error(responseData.message || '创建活动失败')
    }

    // 提交成功
    emit('activity-created')

  } catch (error) {
    console.error('提交失败:', error)
    alert(`创建活动失败: ${error.message}`)
  } finally {
    submitting.value = false
  }
}

// 关闭对话框
function closeDialog() {
  emit('close')
}

// 监听本地时间变化，更新API格式时间
watch(localStartTime, (newValue) => {
  if (newValue) {
    formData.start_time = convertDateTimeLocalToAPI(newValue)
  }
})

watch(localEndTime, (newValue) => {
  if (newValue) {
    formData.end_time = convertDateTimeLocalToAPI(newValue)
  }
})

// 初始化
onMounted(() => {
  setDefaultTimes()
})
</script>

<style scoped>
.create-activity-dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(8px);
  display: grid;
  place-items: center;
  z-index: 3000;
  padding: 20px;
  animation: fadeIn 0.3s ease;
}

.create-activity-dialog {
  width: 100%;
  max-width: 900px;
  max-height: 90vh;
  background: rgba(20, 25, 30, 0.95);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 40px 80px rgba(0, 0, 0, 0.5);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  animation: slideUp 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}

.dialog-header {
  padding: 24px 32px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  background: linear-gradient(180deg,
  rgba(255, 255, 255, 0.05) 0%,
  rgba(255, 255, 255, 0.02) 100%);
  position: relative;
}

.dialog-title {
  margin: 0;
  font-size: 28px;
  font-weight: 900;
  background: linear-gradient(135deg, #fff 0%, rgba(255, 255, 255, 0.8) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: 0.02em;
}

.dialog-subtitle {
  margin-top: 8px;
  color: rgba(255, 255, 255, 0.6);
  font-size: 14px;
}

.close-btn {
  position: absolute;
  top: 24px;
  right: 24px;
  width: 40px;
  height: 40px;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.8);
  font-size: 24px;
  cursor: pointer;
  transition: all 0.2s ease;
  display: grid;
  place-items: center;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 100, 100, 0.3);
  color: #ff6464;
  transform: rotate(90deg);
}

.dialog-body {
  flex: 1;
  overflow-y: auto;
  padding: 32px;
}

.form-container {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.form-section {
  background: rgba(255, 255, 255, 0.02);
  border-radius: 16px;
  padding: 24px;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.section-title {
  margin: 0 0 20px 0;
  font-size: 18px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.9);
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.8);
}

.label-required {
  color: #ff6464;
}

.form-input,
.form-textarea {
  width: 100%;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  color: white;
  font-size: 14px;
  transition: all 0.3s ease;
  box-sizing: border-box;  /* 关键：确保padding包含在宽度内 */
  max-width: 100%;         /* 防止超出父容器 */
}

.form-input:hover,
.form-textarea:hover {
  border-color: rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.07);
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: rgba(120, 200, 255, 0.5);
  background: rgba(120, 200, 255, 0.05);
  box-shadow: 0 0 0 3px rgba(120, 200, 255, 0.1);
}

.form-input.error,
.form-textarea.error {
  border-color: rgba(255, 100, 100, 0.5);
  background: rgba(255, 100, 100, 0.05);
}

.error-message {
  margin-top: 6px;
  color: #ff6464;
  font-size: 12px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.radio-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.radio-option {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.radio-option:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.2);
}

.radio-option.selected {
  background: rgba(120, 200, 255, 0.1);
  border-color: rgba(120, 200, 255, 0.3);
}

.radio-option.selected::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg,
  rgba(120, 200, 255, 0.8),
  rgba(200, 160, 255, 0.8));
}

.radio-input {
  display: none;
}

.radio-label {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
}

.radio-title {
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
}

.radio-description {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
}

/* 对话框底部 */
.dialog-footer {
  padding: 24px 32px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  background: rgba(0, 0, 0, 0.3);
}

.btn-cancel {
  padding: 12px 24px;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.8);
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-cancel:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.2);
}

.btn-submit {
  position: relative;
  padding: 12px 32px;
  border-radius: 12px;
  border: 1px solid rgba(100, 200, 100, 0.3);
  background: linear-gradient(135deg,
  rgba(100, 200, 100, 0.15),
  rgba(100, 200, 100, 0.08));
  color: #64c864;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
  overflow: hidden;
  min-width: 120px;
}

.btn-submit:hover:not(:disabled) {
  background: linear-gradient(135deg,
  rgba(100, 200, 100, 0.25),
  rgba(100, 200, 100, 0.15));
  border-color: rgba(100, 200, 100, 0.5);
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(100, 200, 100, 0.2);
}

.btn-submit:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.submit-sparkle {
  position: absolute;
  inset: 0;
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.btn-submit:hover:not(:disabled) .submit-sparkle {
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

/* 动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(60px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
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
  .create-activity-dialog {
    max-height: 80vh;
  }

  .dialog-body {
    padding: 20px;
  }

  .form-section {
    padding: 16px;
  }

  .form-row {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .dialog-footer {
    padding: 16px 20px;
    flex-direction: column;
  }

  .btn-cancel,
  .btn-submit {
    width: 100%;
    text-align: center;
  }
}
</style>
