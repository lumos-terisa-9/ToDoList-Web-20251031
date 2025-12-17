<template>
  <transition name="modal">
    <div v-if="isVisible" class="modal-overlay" @click.self="close">
      <div class="modal-container">
        <div class="modal-header">
          <h3>创建新组织</h3>
          <button class="close-btn" @click="close">×</button>
        </div>

        <div class="modal-body">
          <form @submit.prevent="handleSubmit" class="organization-form">
            <!-- 组织名称 -->
            <div class="form-group">
              <label class="form-label">
                <span class="label-text">组织名称</span>
                <span class="label-required">*</span>
              </label>
              <input
                type="text"
                v-model="formData.organization_name"
                placeholder="请输入组织名称"
                required
                class="form-input"
                :class="{ 'error': errors.organization_name }"
              />
              <div v-if="errors.organization_name" class="error-message">
                {{ errors.organization_name }}
              </div>
              <div class="char-count">
                已输入 {{ formData.organization_name.length }} 个字符（至少1个字符）
              </div>
            </div>

            <!-- 申请者介绍 -->
            <div class="form-group">
              <label class="form-label">
                <span class="label-text">组织介绍</span>
                <span class="label-required">*</span>
              </label>
              <textarea
                v-model="formData.applicant_introduction"
                placeholder="请介绍一下组织..."
                required
                class="form-textarea"
                :class="{ 'error': errors.applicant_introduction }"
                rows="3"
              ></textarea>
              <div v-if="errors.applicant_introduction" class="error-message">
                {{ errors.applicant_introduction }}
              </div>
              <div class="char-count">
                已输入 {{ formData.applicant_introduction.length }} 个字符（至少1个字符）
              </div>
            </div>

            <!-- 申请理由 -->
            <div class="form-group">
              <label class="form-label">
                <span class="label-text">申请理由</span>
                <span class="label-required">*</span>
              </label>
              <textarea
                v-model="formData.application_reason"
                placeholder="请说明创建这个组织的原因..."
                required
                class="form-textarea"
                :class="{ 'error': errors.application_reason }"
                rows="4"
              ></textarea>
              <div v-if="errors.application_reason" class="error-message">
                {{ errors.application_reason }}
              </div>
              <div class="char-count">
                已输入 {{ formData.application_reason.length }} 个字符（至少1个字符）
              </div>
            </div>

            <!-- 附件上传 -->
            <div class="form-group">
              <label class="form-label">
                <span class="label-text">附件上传</span>
                <span class="label-optional">（选填）</span>
              </label>
              <div class="upload-area" @click="handleUploadClick" @dragover.prevent @drop.prevent="handleDrop">
                <input
                  type="file"
                  ref="fileInput"
                  @change="handleFileSelect"
                  style="display: none"
                />
                <div class="upload-icon">📎</div>
                <div class="upload-text">
                  <div v-if="!formData.attachment_url">
                    点击或拖拽文件到此处上传
                    <div class="upload-hint">支持图片、PDF、文档等格式，最大2MB</div>
                  </div>
                  <div v-else class="upload-success">
                    <div class="file-info">
                      <span class="file-icon">📄</span>
                      <span class="file-name">{{ uploadedFileName }}</span>
                    </div>
                    <button type="button" class="remove-btn" @click.stop="removeFile">移除</button>
                  </div>
                </div>
              </div>
              <div v-if="uploading" class="upload-progress">
                <div class="progress-bar">
                  <div class="progress-fill" :style="{ width: uploadProgress + '%' }"></div>
                </div>
                <div class="progress-text">上传中... {{ uploadProgress }}%</div>
              </div>
            </div>

            <!-- 表单操作 -->
            <div class="form-actions">
              <button type="button" class="action-btn secondary" @click="close">
                取消
              </button>
              <button type="submit" class="action-btn primary" :disabled="loading">
                <span v-if="loading" class="loading-spinner"></span>
                {{ loading ? '创建中...' : '确认创建' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'

const props = defineProps({
  isVisible: Boolean
})

const emit = defineEmits(['close', 'created'])

// API 基础URL
const API_BASE = 'http://localhost:8080/api'

// GitHub配置
const GITHUB_CONFIG = {
}

// 响应式数据
const loading = ref(false)
const uploading = ref(false)
const uploadProgress = ref(0)
const errors = reactive({})
const fileInput = ref(null)
const uploadedFileName = ref('')

// 表单数据
const formData = reactive({
  organization_name: '',
  applicant_introduction: '',
  application_reason: '',
  attachment_url: ''
})

// 获取处理后的token
function getToken() {
  let token = localStorage.getItem('token')
  console.log('从localStorage获取原始token:', token)

  if (token && token.startsWith('{')) {
    try {
      const tokenData = JSON.parse(token)
      console.log('解析token数据:', tokenData)

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
  return token
}

// 文件转Base64
function fileToBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (e) => resolve(e.target.result)
    reader.onerror = (error) => reject(error)
    reader.readAsDataURL(file)
  })
}

// 上传文件到GitHub
async function uploadToGitHub(file) {
  try {
    // 验证文件大小（2MB）
    if (file.size > 2 * 1024 * 1024) {
      throw new Error('文件大小不能超过2MB')
    }

    // 将文件转换为Base64
    const base64Data = await fileToBase64(file)
    const cleanBase64 = base64Data.split(',')[1]

    // 生成唯一的文件名
    const fileExtension = file.name.split('.').pop() || 'bin'
    const fileName = `org_attachment_${Date.now()}_${Math.random().toString(36).substr(2, 9)}.${fileExtension}`

    // 构造API URL
    const apiUrl = `https://api.github.com/repos/${GITHUB_CONFIG.username}/${GITHUB_CONFIG.repo}/contents/${GITHUB_CONFIG.folder}/${fileName}`

    // 请求数据
    const requestData = {
      message: `Upload organization attachment: ${fileName}`,
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

    // 返回文件URL
    return `https://${GITHUB_CONFIG.username}.github.io/${GITHUB_CONFIG.folder}/${fileName}`

  } catch (error) {
    console.error('GitHub上传错误:', error)
    throw error
  }
}

// 表单验证
function validateForm() {
  let isValid = true
  errors.organization_name = ''
  errors.applicant_introduction = ''
  errors.application_reason = ''

  if (!formData.organization_name.trim()) {
    errors.organization_name = '组织名称不能为空'
    isValid = false
  } else if (formData.organization_name.trim().length < 1) {
    errors.organization_name = '组织名称至少需要1个字符'
    isValid = false
  }

  if (!formData.applicant_introduction.trim()) {
    errors.applicant_introduction = '申请者介绍不能为空'
    isValid = false
  } else if (formData.applicant_introduction.trim().length < 1) {
    errors.applicant_introduction = '申请者介绍至少需要1个字符'
    isValid = false
  }

  if (!formData.application_reason.trim()) {
    errors.application_reason = '申请理由不能为空'
    isValid = false
  } else if (formData.application_reason.trim().length < 1) {
    errors.application_reason = '申请理由至少需要1个字符'
    isValid = false
  }

  return isValid
}

// 文件上传处理
function handleUploadClick() {
  if (uploading.value) return
  fileInput.value.click()
}

async function handleFileSelect(event) {
  const file = event.target.files[0]
  if (!file) return

  await uploadFile(file)
}

async function handleDrop(event) {
  event.preventDefault()
  const file = event.dataTransfer.files[0]
  if (!file) return

  await uploadFile(file)
}

async function uploadFile(file) {
  try {
    uploading.value = true
    uploadProgress.value = 10 // 模拟进度

    // 上传文件到GitHub
    const fileUrl = await uploadToGitHub(file)

    uploadProgress.value = 100
    formData.attachment_url = fileUrl
    uploadedFileName.value = file.name

    console.log('文件上传成功，URL:', fileUrl)
  } catch (error) {
    console.error('文件上传失败:', error)
    alert(`文件上传失败: ${error.message}`)
  } finally {
    uploading.value = false
    uploadProgress.value = 0
  }
}

function removeFile() {
  formData.attachment_url = ''
  uploadedFileName.value = ''
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

// 提交表单
async function handleSubmit() {
  if (!validateForm()) {
    return
  }

  const token = getToken()
  if (!token) {
    alert('请先登录')
    return
  }

  loading.value = true

  try {
    console.log('开始创建组织请求...')
    console.log('请求数据:', JSON.stringify(formData, null, 2))

    const response = await fetch(`${API_BASE}/organization/application/create-organization`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(formData)
    })

    const data = await response.json()
    console.log('创建组织响应:', data)

    if (response.ok) {
      alert('创建组织申请已提交，等待审核！')
      resetForm()
      emit('created')
    } else {
      alert(data.message || '创建组织失败，请稍后重试')
    }
  } catch (error) {
    console.error('创建组织请求失败:', error)
    alert('创建组织失败，请检查网络连接')
  } finally {
    loading.value = false
  }
}

// 重置表单
function resetForm() {
  formData.organization_name = ''
  formData.applicant_introduction = ''
  formData.application_reason = ''
  formData.attachment_url = ''
  uploadedFileName.value = ''
  Object.keys(errors).forEach(key => errors[key] = '')
}

// 关闭模态框
function close() {
  resetForm()
  emit('close')
}

// 监听模态框显示状态
watch(() => props.isVisible, (newVal) => {
  if (!newVal) {
    resetForm()
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
  backdrop-filter: blur(4px);
  padding-top: 1vh; /* 增加上边缘距离1.5倍 */
}

.modal-container {
  background: white;
  border-radius: 20px;
  width: 500px;
  max-width: 90vw;
  max-height: 70vh; /* 缩短弹窗长度 */
  overflow-y: hidden; /* 去掉滚动条 */
  overflow-x: hidden; /* 确保水平方向也没有滚动条 */
  color: #333;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px; /* 减少内边距 */
  border-bottom: 1px solid #f0f0f0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 20px 20px 0 0;
  color: white;
  flex-shrink: 0;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px; /* 减小字体大小 */
  font-weight: 600;
}

.close-btn {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  font-size: 22px; /* 减小关闭按钮 */
  cursor: pointer;
  color: white;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background 0.3s;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.modal-body {
  padding: 20px 24px; /* 减少内边距 */
  overflow-y: auto; /* 只在内容区域允许滚动 */
  flex: 1;
}

.organization-form {
  display: flex;
  flex-direction: column;
  gap: 16px; /* 减小间距 */
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px; /* 减小间距 */
}

.form-label {
  display: flex;
  align-items: center;
  gap: 4px;
  font-weight: 500;
  color: #333;
  font-size: 13px; /* 减小字体大小 */
}

.label-required {
  color: #f56565;
}

.label-optional {
  color: #a0aec0;
  font-size: 11px; /* 减小字体大小 */
}

.form-input,
.form-textarea {
  padding: 10px 14px; /* 减少内边距 */
  border: 2px solid #e2e8f0;
  border-radius: 10px; /* 减小圆角 */
  font-size: 13px; /* 减小字体大小 */
  transition: all 0.3s;
  background: white;
  color: #333;
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-input.error,
.form-textarea.error {
  border-color: #f56565;
}

.form-input.error:focus,
.form-textarea.error:focus {
  box-shadow: 0 0 0 3px rgba(245, 101, 101, 0.1);
}

.form-textarea {
  resize: vertical;
  min-height: 70px; /* 减小最小高度 */
  line-height: 1.5;
}

.char-count {
  font-size: 11px;
  color: #718096;
  margin-top: 2px;
  text-align: right;
}

.error-message {
  color: #f56565;
  font-size: 11px;
  margin-top: 2px;
}

.upload-area {
  border: 2px dashed #e2e8f0;
  border-radius: 10px;
  padding: 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
  background: #f8fafc;
}

.upload-area:hover {
  border-color: #667eea;
  background: #f0f7ff;
}

.upload-icon {
  font-size: 32px;
  margin-bottom: 10px;
  color: #a0aec0;
}

.upload-text {
  font-size: 13px;
  color: #718096;
}

.upload-hint {
  font-size: 11px;
  color: #a0aec0;
  margin-top: 4px;
}

.upload-success {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.file-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-icon {
  font-size: 20px;
}

.file-name {
  font-size: 13px;
  color: #333;
  font-weight: 500;
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.remove-btn {
  padding: 4px 10px;
  border: 1px solid #e2e8f0;
  background: white;
  color: #f56565;
  border-radius: 6px;
  font-size: 11px;
  cursor: pointer;
  transition: all 0.3s;
}

.remove-btn:hover {
  background: #fed7d7;
  border-color: #f56565;
}

.upload-progress {
  margin-top: 10px;
}

.progress-bar {
  height: 6px;
  background: #e2e8f0;
  border-radius: 3px;
  overflow: hidden;
  margin-bottom: 4px;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 3px;
  transition: width 0.3s;
}

.progress-text {
  font-size: 11px;
  color: #718096;
  text-align: center;
}

.form-actions {
  display: flex;
  gap: 10px; /* 减小间距 */
  margin-top: 16px; /* 减小间距 */
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
  flex-shrink: 0;
}

.action-btn {
  flex: 1;
  padding: 12px 20px; /* 减少内边距 */
  border: none;
  border-radius: 10px; /* 减小圆角 */
  font-size: 13px; /* 减小字体大小 */
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px; /* 减小间距 */
}

.action-btn.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.action-btn.primary:hover:not(:disabled) {
  transform: translateY(-1px); /* 减小动画幅度 */
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3); /* 减小阴影 */
}

.action-btn.primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.action-btn.secondary {
  background: #f7fafc;
  color: #4a5568;
  border: 2px solid #e2e8f0;
}

.action-btn.secondary:hover {
  background: #edf2f7;
  border-color: #cbd5e0;
}

.loading-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
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
@media (max-width: 640px) {
  .modal-header {
    padding: 16px 20px;
  }

  .modal-body {
    padding: 16px 20px;
  }

  .form-actions {
    flex-direction: column;
  }

  .action-btn {
    padding: 10px 16px;
  }
}
</style>
