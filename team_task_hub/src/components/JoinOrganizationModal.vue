<template>
  <transition name="modal">
    <div v-if="isVisible" class="modal-overlay" @click.self="close">
      <div class="modal-container">
        <div class="modal-header">
          <h3>申请加入组织</h3>
          <button class="close-btn" @click="close">×</button>
        </div>

        <div class="modal-body">
          <!-- 标签页切换 -->
          <div class="tab-container">
            <div class="tab-header">
              <button
                class="tab-btn"
                :class="{ active: activeTab === 'nearby' }"
                @click="activeTab = 'nearby'"
              >
                <span class="tab-icon">📍</span>
                <span class="tab-text">附近组织</span>
              </button>
              <button
                class="tab-btn"
                :class="{ active: activeTab === 'search' }"
                @click="activeTab = 'search'"
              >
                <span class="tab-icon">🔍</span>
                <span class="tab-text">搜索组织</span>
              </button>
            </div>

            <div class="tab-content">
              <!-- 附近组织标签页 -->
              <div v-if="activeTab === 'nearby'" class="tab-pane">
                <div class="search-section">
                  <button
                    class="location-btn wide-action-btn"
                    @click="showLocationFeature"
                    :class="{ 'has-location': hasUserLocation }"
                    :disabled="fetchingLocation"
                  >
                    <span v-if="fetchingLocation" class="loading-spinner small"></span>
                    <span v-else class="btn-icon">📍</span>
                    <span class="btn-text">
                      {{ fetchingLocation ? '获取位置中...' : (hasUserLocation ? '已获取当前位置' : '获取当前位置') }}
                    </span>
                    <span v-if="hasUserLocation && !fetchingLocation" class="location-status-indicator">
                      ✓
                    </span>
                  </button>

                  <div v-if="hasUserLocation" class="location-info">
                    <div class="location-coords">
                      纬度: {{ userLocation.latitude.toFixed(6) }},
                      经度: {{ userLocation.longitude.toFixed(6) }}
                    </div>
                    <div v-if="userLocation.accuracy" class="location-accuracy">
                      精度: ±{{ userLocation.accuracy.toFixed(1) }}米
                    </div>
                  </div>

                  <button
                    class="search-btn primary wide-action-btn"
                    @click="searchNearbyOrganizations"
                    :disabled="loading.nearby || !hasUserLocation"
                  >
                    <span v-if="loading.nearby" class="loading-spinner small"></span>
                    {{ loading.nearby ? '搜索中...' : '搜索附近组织' }}
                    <span v-if="!hasUserLocation" class="search-hint">
                      （请先获取位置）
                    </span>
                  </button>
                </div>

                <!-- 搜索结果 -->
                <div class="results-section">
                  <h4 class="results-title" v-if="nearbyOrganizations.length > 0">
                    找到 {{ nearbyOrganizations.length }} 个附近组织
                  </h4>

                  <div v-if="loading.nearby" class="loading-state">
                    <div class="loading-spinner"></div>
                    <p>正在搜索附近组织...</p>
                  </div>

                  <div v-else-if="nearbyOrganizations.length === 0 && hasSearched" class="empty-state">
                    <div class="empty-icon">🏢</div>
                    <p class="empty-text">该位置附近暂无组织</p>
                    <p class="empty-subtext">
                      当前位置: {{ userLocation.latitude.toFixed(4) }}, {{ userLocation.longitude.toFixed(4) }}
                    </p>
                  </div>

                  <div v-else-if="!hasSearched" class="empty-state">
                    <div class="empty-icon">📍</div>
                    <p class="empty-text">获取位置后即可搜索附近组织</p>
                    <p class="empty-subtext">
                      点击上方"获取当前位置"按钮开始
                    </p>
                  </div>

                  <div v-else class="organizations-list">
                    <div
                      v-for="org in nearbyOrganizations"
                      :key="org.id"
                      class="organization-card-compact"
                    >
                      <div class="org-info-compact">
                        <div class="org-name-compact">{{ org.name }}</div>
                      </div>
                      <div class="org-actions-compact">
                        <button
                          class="action-btn-compact"
                          @click="openJoinDialog(org)"
                        >
                          <span class="btn-icon">➕</span>
                          <span class="btn-text">加入</span>
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 搜索组织标签页 -->
              <div v-else class="tab-pane">
                <div class="search-section">
                  <div class="search-input-group-wide">
                    <input
                      type="text"
                      v-model="searchKeyword"
                      placeholder="输入组织名称或关键词"
                      class="search-input-wide"
                    />
                    <button class="search-btn-wide" @click="searchOrganizations">
                      搜索
                    </button>
                  </div>
                </div>

                <!-- 搜索结果 -->
                <div class="results-section">
                  <h4 class="results-title" v-if="searchResults.length > 0">
                    找到 {{ searchResults.length }} 个组织
                  </h4>

                  <div v-if="loading.search" class="loading-state">
                    <div class="loading-spinner"></div>
                    <p>正在搜索组织...</p>
                  </div>

                  <div v-else-if="searchResults.length === 0 && hasSearchedForOrganizations" class="empty-state">
                    <div class="empty-icon">🔍</div>
                    <p class="empty-text">未找到相关组织</p>
                    <p class="empty-subtext">尝试使用其他关键词搜索</p>
                  </div>

                  <div v-else-if="!hasSearchedForOrganizations" class="empty-state">
                    <div class="empty-icon">🔍</div>
                    <p class="empty-text">输入关键词搜索组织</p>
                    <p class="empty-subtext">支持组织名称、关键词等搜索</p>
                  </div>

                  <div v-else class="organizations-list">
                    <div
                      v-for="org in searchResults"
                      :key="org.id"
                      class="organization-card-compact"
                    >
                      <div class="org-info-compact">
                        <div class="org-name-compact">{{ org.name }}</div>
                      </div>
                      <div class="org-actions-compact">
                        <button
                          class="action-btn-compact"
                          @click="openJoinDialog(org)"
                        >
                          <span class="btn-icon">➕</span>
                          <span class="btn-text">加入</span>
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 加入组织对话框 -->
          <div v-if="showJoinDialog" class="join-dialog-overlay">
            <div class="join-dialog-wide">
              <div class="join-dialog-header">
                <h4>加入 {{ selectedOrg?.name }}</h4>
                <button class="close-btn small" @click="closeJoinDialog">×</button>
              </div>

              <div class="join-dialog-body">
                <div class="join-options">
                  <button
                    class="join-option"
                    :class="{ 'active': joinMethod === 'code' }"
                    @click="joinMethod = 'code'"
                  >
                    <div class="option-icon">🔑</div>
                    <div class="option-content">
                      <div class="option-title">使用邀请码加入</div>
                      <div class="option-desc">已有组织邀请码？直接输入即可加入</div>
                    </div>
                  </button>

                  <button
                    class="join-option"
                    :class="{ 'active': joinMethod === 'application' }"
                    @click="joinMethod = 'application'"
                  >
                    <div class="option-icon">📝</div>
                    <div class="option-content">
                      <div class="option-title">申请加入</div>
                      <div class="option-desc">提交申请，等待组织管理员审核</div>
                    </div>
                  </button>
                </div>

                <!-- 邀请码加入表单 -->
                <div v-if="joinMethod === 'code'" class="join-form">
                  <div class="form-group">
                    <label class="form-label">邀请码</label>
                    <input
                      type="text"
                      v-model="invitationCode"
                      placeholder="请输入邀请码"
                      class="form-input"
                    />
                  </div>
                  <div class="form-actions">
                    <button
                      class="action-btn secondary"
                      @click="closeJoinDialog"
                    >
                      取消
                    </button>
                    <button
                      class="action-btn primary"
                      @click="joinWithCode"
                      :disabled="loading.join"
                    >
                      <span v-if="loading.join" class="loading-spinner"></span>
                      {{ loading.join ? '加入中...' : '确认加入' }}
                    </button>
                  </div>
                </div>

                <!-- 申请加入表单 -->
                <div v-else class="join-form-no-scroll">
                  <div class="form-group">
                    <label class="form-label">
                      <span class="label-text">申请者介绍</span>
                      <span class="label-required">*</span>
                    </label>
                    <textarea
                      v-model="applicationData.applicant_introduction"
                      placeholder="请介绍一下您自己..."
                      class="form-textarea"
                      rows="3"
                    ></textarea>
                    <div class="char-count">
                      已输入 {{ applicationData.applicant_introduction.length }} 个字符（至少1个字符）
                    </div>
                  </div>

                  <div class="form-group">
                    <label class="form-label">
                      <span class="label-text">申请理由</span>
                      <span class="label-required">*</span>
                    </label>
                    <textarea
                      v-model="applicationData.application_reason"
                      placeholder="请说明为什么想加入这个组织..."
                      class="form-textarea"
                      rows="4"
                    ></textarea>
                    <div class="char-count">
                      已输入 {{ applicationData.application_reason.length }} 个字符（至少1个字符）
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
                        <div v-if="!applicationData.attachment_url">
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

                  <div class="form-actions">
                    <button
                      class="action-btn secondary"
                      @click="closeJoinDialog"
                    >
                      取消
                    </button>
                    <button
                      class="action-btn primary"
                      @click="submitApplication"
                      :disabled="loading.application"
                    >
                      <span v-if="loading.application" class="loading-spinner"></span>
                      {{ loading.application ? '提交中...' : '提交申请' }}
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
import { ref, reactive, watch, computed } from 'vue'

const props = defineProps({
  isVisible: Boolean
})

const emit = defineEmits(['close', 'joined'])

// API 基础URL
const API_BASE = 'http://localhost:8080/api'

// GitHub配置
const GITHUB_CONFIG = {
}

// 响应式数据
const hasSearchedForOrganizations = ref(false)
const fetchingLocation = ref(false) // 位置获取中状态
const activeTab = ref('nearby')
const loading = reactive({
  nearby: false,
  search: false,
  join: false,
  application: false
})
const uploading = ref(false)
const uploadProgress = ref(0)
const fileInput = ref(null)
const uploadedFileName = ref('')
const hasSearched = ref(false)

// 用户位置相关
const userLocation = ref(null) // 移除默认值，初始为null

// 搜索相关
const searchKeyword = ref('')
const nearbyOrganizations = ref([])
const searchResults = ref([])

// 加入组织相关
const showJoinDialog = ref(false)
const selectedOrg = ref(null)
const joinMethod = ref('code')
const invitationCode = ref('')
const applicationData = reactive({
  applicant_introduction: '',
  application_reason: '',
  attachment_url: '',
  organization_name: ''
})

// 计算属性
const hasUserLocation = computed(() => {
  return userLocation.value !== null &&
    userLocation.value.latitude !== undefined &&
    userLocation.value.longitude !== undefined
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
    applicationData.attachment_url = fileUrl
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
  applicationData.attachment_url = ''
  uploadedFileName.value = ''
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

// 获取当前位置
async function getCurrentLocation() {
  return new Promise((resolve, reject) => {
    // 检查浏览器是否支持Geolocation API
    if (!navigator.geolocation) {
      reject(new Error('您的浏览器不支持地理位置功能'))
      return
    }

    // 设置获取位置选项
    const options = {
      enableHighAccuracy: true, // 尝试获取高精度位置
      timeout: 300000, // 30秒
      maximumAge: 60000 // 允许使用1分钟内缓存的位置
    }

    // 获取当前位置
    navigator.geolocation.getCurrentPosition(
      (position) => {
        resolve({
          latitude: position.coords.latitude,
          longitude: position.coords.longitude,
          accuracy: position.coords.accuracy
        })
      },
      (error) => {
        let errorMessage

        switch(error.code) {
          case error.PERMISSION_DENIED:
            errorMessage = '位置请求被拒绝。请允许浏览器访问位置信息。'
            break
          case error.POSITION_UNAVAILABLE:
            errorMessage = '位置信息不可用。请检查您的设备位置服务。'
            break
          case error.TIMEOUT:
            errorMessage = '获取位置超时。请确保设备有良好的网络连接。'
            break
          default:
            errorMessage = '获取位置失败。请稍后重试。'
        }

        reject(new Error(errorMessage))
      },
      options
    )
  })
}

// 获取当前位置功能
async function showLocationFeature() {
  if (fetchingLocation.value) return // 防止重复点击

  try {
    fetchingLocation.value = true // 开始获取

    const location = await getCurrentLocation()

    // 保存位置信息
    userLocation.value = {
      latitude: location.latitude,
      longitude: location.longitude,
      accuracy: location.accuracy
    }

    console.log('位置获取成功:', userLocation.value)
  } catch (error) {
    console.error('获取位置失败:', error)
    alert(`获取位置失败: ${error.message}`)
  } finally {
    fetchingLocation.value = false // 无论成功失败都结束加载
  }
}

// 搜索附近组织
async function searchNearbyOrganizations() {
  // 检查是否已获取位置
  if (!hasUserLocation.value) {
    alert('请先获取当前位置')
    return
  }

  const token = getToken()
  if (!token) {
    alert('请先登录')
    return
  }

  loading.nearby = true
  nearbyOrganizations.value = []
  hasSearched.value = true

  try {
    const params = {
      lat: userLocation.value.latitude,
      lng: userLocation.value.longitude
    }

    console.log('搜索附近组织，使用位置:', params)

    // 方案1A：改回使用查询参数（最简单）
    const response = await fetch(`${API_BASE}/organization/nearby?lat=${params.lat}&lng=${params.lng}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    const data = await response.json()
    console.log('附近组织响应:', data)

    if (response.ok) {
      nearbyOrganizations.value = data.data?.organizations || []

      if (nearbyOrganizations.value.length === 0) {
        console.log('该位置附近没有找到组织')
      }
    } else {
      alert(data.message || '搜索附近组织失败')
    }
  } catch (error) {
    console.error('搜索附近组织失败:', error)
    alert('搜索失败，请检查网络连接')
  } finally {
    loading.nearby = false
  }
}

// 搜索组织
async function searchOrganizations() {
  if (!searchKeyword.value.trim()) {
    searchResults.value = [] // 清空结果
    alert('请输入搜索关键词')
    return
  }

  // 设置已搜索标志
  hasSearchedForOrganizations.value = true

  const token = getToken()
  if (!token) {
    alert('请先登录')
    return
  }

  loading.search = true
  searchResults.value = []

  try {
    console.log('搜索组织，关键词:', searchKeyword.value)

    // 根据API文档，使用GET请求，keyword作为查询参数
    const url = new URL(`${API_BASE}/organization/search`)
    url.searchParams.append('keyword', searchKeyword.value.trim())

    console.log('请求URL:', url.toString())

    const response = await fetch(url.toString(), {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Accept': 'application/json'
      }
      // 注意：GET请求没有body
    })

    console.log('响应状态:', response.status, response.statusText)

    // 先获取响应文本
    const responseText = await response.text()
    console.log('原始响应:', responseText)

    // 解析JSON
    let data
    try {
      data = JSON.parse(responseText)
      console.log('解析后的数据:', data)
    } catch (parseError) {
      console.error('JSON解析失败:', parseError)
      throw new Error('服务器返回格式错误: ' + responseText.substring(0, 100))
    }

    // 根据响应格式处理
    if (response.ok) {
      if (data.success === true) {
        searchResults.value = data.data || []
        console.log(`搜索到 ${searchResults.value.length} 个组织`)
      } else {
        alert(data.message || '搜索失败')
      }
    } else {
      // HTTP状态码错误
      alert(data.message || `搜索失败: ${response.status} ${response.statusText}`)
    }

  } catch (error) {
    console.error('搜索组织失败:', error)

    if (error.message.includes('Failed to fetch')) {
      alert('网络连接失败，请检查网络连接')
    } else {
      alert(`搜索失败: ${error.message}`)
    }
  } finally {
    loading.search = false
  }
}

// 打开加入对话框
function openJoinDialog(org) {
  selectedOrg.value = org
  applicationData.organization_name = org.name
  showJoinDialog.value = true
  joinMethod.value = 'code' // 默认选择邀请码加入
}

// 关闭加入对话框
function closeJoinDialog() {
  showJoinDialog.value = false
  selectedOrg.value = null
  invitationCode.value = ''
  applicationData.applicant_introduction = ''
  applicationData.application_reason = ''
  applicationData.attachment_url = ''
  applicationData.organization_name = ''
  uploadedFileName.value = ''
  uploading.value = false
  uploadProgress.value = 0
}

// 表单验证
function validateApplication() {
  if (!applicationData.applicant_introduction.trim()) {
    alert('请输入申请者介绍')
    return false
  }

  if (!applicationData.application_reason.trim()) {
    alert('请输入申请理由')
    return false
  }

  return true
}

// 使用邀请码加入组织
async function joinWithCode() {
  if (!invitationCode.value.trim()) {
    alert('请输入邀请码')
    return
  }

  if (invitationCode.value.trim().length !== 6) {
    alert('邀请码必须为6位')
    return
  }

  const token = getToken()
  if (!token) {
    alert('请先登录')
    return
  }

  loading.join = true

  try {
    const requestData = {
      code: invitationCode.value,
      organization_id: selectedOrg.value.id
    }

    console.log('使用邀请码加入组织，数据:', requestData)

    const response = await fetch(`${API_BASE}/organization/join-with-code`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(requestData)
    })

    const data = await response.json()
    console.log('邀请码加入响应:', data)

    if (response.ok) {
      alert('成功加入组织！')
      closeJoinDialog()
      emit('joined')
    } else {
      alert(data.message || '加入组织失败')
    }
  } catch (error) {
    console.error('使用邀请码加入组织失败:', error)
    alert('加入失败，请检查网络连接')
  } finally {
    loading.join = false
  }
}

// 提交加入申请
async function submitApplication() {
  if (!validateApplication()) {
    return
  }

  const token = getToken()
  if (!token) {
    alert('请先登录')
    return
  }

  loading.application = true

  try {
    console.log('提交加入申请，数据:', applicationData)

    const response = await fetch(`${API_BASE}/organization/application/join-organization`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(applicationData)
    })

    const data = await response.json()
    console.log('加入申请响应:', data)

    if (response.ok) {
      alert('申请已提交，等待组织管理员审核！')
      closeJoinDialog()
      emit('joined')
    } else {
      alert(data.message || '提交申请失败')
    }
  } catch (error) {
    console.error('提交加入申请失败:', error)
    alert('提交失败，请检查网络连接')
  } finally {
    loading.application = false
  }
}

// 关闭模态框
function close() {
  closeJoinDialog()
  emit('close')
}

// 监听模态框显示状态
watch(() => props.isVisible, (newVal) => {
  if (!newVal) {
    closeJoinDialog()
    searchKeyword.value = ''
    nearbyOrganizations.value = []
    searchResults.value = []
    hasSearched.value = false
    hasSearchedForOrganizations.value = false
  }
})
</script>

<style scoped>
/* 宽按钮统一样式 */
.wide-action-btn {
  width: 95%;
  max-width: 700px;
  margin: 0 auto 16px auto;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 16px 28px;
  border: none;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 6px 16px rgba(102, 126, 234, 0.3);
  position: relative;
}

/* 位置按钮已获取状态 */
.wide-action-btn.has-location {
  background: linear-gradient(135deg, #48bb78 0%, #38a169 100%);
}

.wide-action-btn:not(.has-location) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.wide-action-btn:hover:not(:disabled) {
  transform: translateY(-3px);
  box-shadow: 0 12px 24px rgba(102, 126, 234, 0.4);
}

.wide-action-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 位置状态指示器 */
.location-status-indicator {
  margin-left: 10px;
  font-size: 16px;
  font-weight: bold;
  color: white;
  background: rgba(255, 255, 255, 0.2);
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 位置信息显示 */
.location-info {
  width: 95%;
  max-width: 700px;
  margin: 0 auto 16px auto;
  padding: 12px 16px;
  background: #f0f9ff;
  border: 1px solid #bee3f8;
  border-radius: 8px;
  font-size: 12px;
  color: #2c5282;
}

.location-coords {
  font-weight: 500;
  margin-bottom: 4px;
}

.location-accuracy {
  color: #718096;
  font-size: 11px;
}

/* 搜索提示 */
.search-hint {
  margin-left: 10px;
  font-size: 12px;
  font-weight: normal;
  opacity: 0.8;
}

/* 搜索按钮样式 */
.search-btn.primary {
  background: linear-gradient(135deg, #4299e1 0%, #667eea 100%);
}

/* 搜索输入组 - 宽版 */
.search-input-group-wide {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
  width: 100%;
  max-width: 500px;
  margin-left: auto;
  margin-right: auto;
}

.search-input-wide {
  flex: 1;
  padding: 12px 16px;
  border: 2px solid #e2e8f0;
  border-radius: 12px;
  font-size: 14px;
  transition: all 0.3s;
}

.search-input-wide:focus {
  outline: none;
  border-color: #4299e1;
  box-shadow: 0 0 0 3px rgba(66, 153, 225, 0.1);
}

.search-btn-wide {
  padding: 12px 24px;
  border: none;
  background: linear-gradient(135deg, #4299e1 0%, #667eea 100%);
  color: white;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
  min-width: 80px;
}

.search-btn-wide:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(66, 153, 225, 0.4);
}

/* 紧凑的组织卡片 */
.organization-card-compact {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  background: #f7fafc;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  transition: all 0.3s;
  width: 95%;  /* 增加宽度 */
  max-width: 600px;  /* 增加最大宽度 */
  margin: 0 auto 12px auto;  /* 增加间距 */
}

.organization-card-compact:hover {
  background: #edf2f7;
  border-color: #cbd5e0;
  transform: translateX(2px);
}

.org-info-compact {
  flex: 1;
  overflow: hidden;
}

.org-name-compact {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.org-actions-compact {
  margin-left: 12px;
}

.action-btn-compact {
  padding: 6px 12px;
  border: none;
  background: linear-gradient(135deg, #48bb78 0%, #38a169 100%);
  color: white;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  gap: 4px;
}

.action-btn-compact:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(72, 187, 120, 0.3);
}

/* 加宽的加入对话框 */
.join-dialog-overlay {
  position: fixed;
  top: 10vh; /* 增加上边缘距离 */
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: flex-start; /* 顶部对齐 */
  z-index: 2100;
  backdrop-filter: blur(4px);
  padding-top: 5vh; /* 增加顶部内边距 */
}

.join-dialog-wide {
  background: white;
  border-radius: 16px;
  width: 500px;
  max-width: 90vw;
  max-height: 80vh;
  overflow-y: auto;
  scrollbar-width: none;
  -ms-overflow-style: none;

  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

/* 无滚动条的加入表单 */
.join-form-no-scroll {
  animation: slideIn 0.3s ease;
  overflow-y: hidden; /* 移除滚动条 */
  max-height: none;
}

/* 激活状态的加入选项 */
.join-option.active {
  border-color: #4299e1;
  background: #f0f7ff;
  border-width: 2px;
}

.join-option.active .option-icon {
  background: rgba(66, 153, 225, 0.1);
  color: #4299e1;
}

/* 字符计数样式 */
.char-count {
  font-size: 11px;
  color: #718096;
  margin-top: 2px;
  text-align: right;
}

/* 表单标签 */
.form-label {
  display: flex;
  align-items: center;
  gap: 4px;
  font-weight: 500;
  color: #333;
  font-size: 13px;
  margin-bottom: 6px;
}

.label-required {
  color: #f56565;
}

.label-optional {
  color: #a0aec0;
  font-size: 11px;
}

/* 上传区域样式 */
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

/* 其余样式保持不变 */

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
}

.modal-container {
  background: white;
  border-radius: 20px;
  width: 600px;
  max-width: 90vw;
  max-height: 85vh;
  overflow-y: auto;
  scrollbar-width: none;
  -ms-overflow-style: none;
  color: #333;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 28px;
  border-bottom: 1px solid #f0f0f0;
  background: linear-gradient(135deg, #4299e1 0%, #667eea 100%);
  border-radius: 20px 20px 0 0;
  color: white;
}

.modal-header h3 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.close-btn {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #c3c2c2;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background 0.3s;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.close-btn.small {
  width: 28px;
  height: 28px;
  font-size: 20px;
}

.modal-body {
  padding: 28px;
  position: relative;
}

/* 标签页样式 */
.tab-container {
  display: flex;
  flex-direction: column;
}

.tab-header {
  display: flex;
  gap: 4px;
  background: #f7fafc;
  border-radius: 12px;
  padding: 4px;
  margin-bottom: 24px;
}

.tab-btn {
  flex: 1;
  padding: 12px 16px;
  border: none;
  background: transparent;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  color: #718096;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.tab-btn.active {
  background: white;
  color: #4299e1;
  box-shadow: 0 2px 8px rgba(66, 153, 225, 0.2);
}

.tab-icon {
  font-size: 16px;
}

.tab-text {
  font-weight: 500;
}

/* 搜索部分 */
.search-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 24px;
  text-align: center;
}

/* 结果部分 */
.results-section {
  margin-top: 16px;
}

.results-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f0f0f0;
}

/* 加载状态 */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  color: #718096;
}

.loading-state p {
  margin-top: 12px;
  font-size: 14px;
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  text-align: center;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-text {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.empty-subtext {
  font-size: 14px;
  color: #718096;
}

/* 组织列表 */
.organizations-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 300px;
  overflow-y: auto;
}

/* 加入对话框 */
.join-dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #f0f0f0;
  background: #f7fafc;
  border-radius: 16px 16px 0 0;
}

.join-dialog-header h4 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.join-dialog-body {
  padding: 24px;
}

/* 加入选项 */
.join-options {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 24px;
}

.join-option {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  border: 2px solid #e2e8f0;
  border-radius: 12px;
  background: white;
  cursor: pointer;
  transition: all 0.3s;
  text-align: left;
}

.join-option:hover {
  border-color: #4299e1;
  background: #f0f7ff;
  transform: translateY(-2px);
}

.option-icon {
  font-size: 24px;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f7fafc;
  border-radius: 10px;
}

.option-content {
  flex: 1;
}

.option-title {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.option-desc {
  font-size: 13px;
  color: #718096;
}

/* 表单样式 */
.join-form {
  animation: slideIn 0.3s ease;
}

.form-group {
  margin-bottom: 20px;
}

.form-input,
.form-textarea {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #e2e8f0;
  border-radius: 12px;
  font-size: 14px;
  transition: all 0.3s;
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: #4299e1;
  box-shadow: 0 0 0 3px rgba(66, 153, 225, 0.1);
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
}

.form-actions {
  display: flex;
  gap: 12px;
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid #f0f0f0;
}

.action-btn {
  flex: 1;
  padding: 14px 24px;
  border: none;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.action-btn.primary {
  background: linear-gradient(135deg, #4299e1 0%, #667eea 100%);
  color: white;
}

.action-btn.primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(66, 153, 225, 0.4);
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

/* 加载动画 */
.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 0.8s linear infinite;
}

.loading-spinner.small {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 0.8s linear infinite;
  margin-right: 8px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 动画效果 */
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

/* 滚动条样式 */
.modal-container::-webkit-scrollbar,
.organizations-list::-webkit-scrollbar,
.join-dialog-wide::-webkit-scrollbar {
  width: 6px;
}

.modal-container::-webkit-scrollbar-track,
.organizations-list::-webkit-scrollbar-track,
.join-dialog-wide::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.modal-container::-webkit-scrollbar-thumb,
.organizations-list::-webkit-scrollbar-thumb,
.join-dialog-wide::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.modal-container::-webkit-scrollbar-thumb:hover,
.organizations-list::-webkit-scrollbar-thumb:hover,
.join-dialog-wide::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .modal-header {
    padding: 20px 24px;
  }

  .modal-body {
    padding: 24px;
  }

  .tab-header {
    flex-direction: column;
  }

  .form-actions {
    flex-direction: column;
  }

  .wide-action-btn,
  .search-input-group-wide {
    max-width: 100%;
  }

  .organization-card-compact {
    width: 95%;
  }

  .join-dialog-wide {
    width: 95vw;
  }

  .location-info {
    width: 95%;
    max-width: 100%;
  }
}

@media (max-width: 480px) {
  .modal-header h3 {
    font-size: 18px;
  }

  .tab-btn {
    padding: 10px 12px;
    font-size: 13px;
  }

  .search-input-group-wide {
    flex-direction: column;
  }

  .search-btn-wide {
    width: 100%;
  }

  .wide-action-btn {
    padding: 14px 20px;
    font-size: 14px;
  }

  .location-status-indicator {
    margin-left: 6px;
    width: 20px;
    height: 20px;
    font-size: 14px;
  }
}
</style>
