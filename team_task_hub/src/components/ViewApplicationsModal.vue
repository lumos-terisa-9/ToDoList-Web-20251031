<template>
  <transition name="modal">
    <div v-if="isVisible" class="modal-overlay" @click.self="close">
      <div class="modal-container">
        <div class="modal-header">
          <h3>我的组织申请</h3>
          <button class="close-btn" @click="close">×</button>
        </div>

        <div class="modal-body">
          <!-- 加载状态 -->
          <div v-if="loading" class="loading-state">
            <div class="loading-spinner large"></div>
            <p>正在加载申请记录...</p>
          </div>

          <!-- 空状态 -->
          <div v-else-if="applications.length === 0" class="empty-state">
            <div class="empty-icon">📋</div>
            <p class="empty-text">暂无申请记录</p>
            <p class="empty-subtext">您还没有提交过组织申请</p>
          </div>

          <!-- 申请列表 -->
          <div v-else class="applications-list">
            <div class="stats-summary">
              <div class="stat-item">
                <div class="stat-value">{{ totalApplications }}</div>
                <div class="stat-label">总申请数</div>
              </div>
              <div class="stat-item">
                <div class="stat-value">{{ pendingApplications }}</div>
                <div class="stat-label">审核中</div>
              </div>
              <div class="stat-item">
                <div class="stat-value">{{ processedApplications }}</div>
                <div class="stat-label">已处理</div>
              </div>
            </div>

            <div class="applications-grid">
              <div
                v-for="application in applications"
                :key="application.id"
                class="application-card"
                :class="{
                  'expanded': expandedApplications.includes(application.id),
                  'status-pending': application.status === 'pending',
                  'status-approved': application.status === 'approved',
                  'status-rejected': application.status === 'rejected',
                  'status-cancelled': application.status === 'cancelled'
                }"
              >
                <!-- 折叠状态：只显示申请类型+组织名和状态 -->
                <div
                  class="application-header-collapsed"
                  @click="toggleApplication(application.id)"
                >
                  <div class="collapsed-content">
                    <div class="application-type">
                      <span class="type-icon">
                        {{ application.application_type === 'create_org' ? '🏢' : '➕' }}
                      </span>
                      <span class="type-text">
                        {{ application.application_type === 'create_org' ? '创建组织' : '申请加入组织' }}
                      </span>
                    </div>
                    <div class="org-name-collapsed">{{ application.organization_name }}</div>
                  </div>
                  <div class="application-status" :class="getStatusColor(application.status)">
                    {{ getStatusText(application.status) }}
                  </div>
                  <div class="expand-icon">
                    {{ expandedApplications.includes(application.id) ? '▼' : '▶' }}
                  </div>
                </div>

                <!-- 展开状态：显示详细信息 -->
                <div
                  v-if="expandedApplications.includes(application.id)"
                  class="application-details"
                >
                  <div class="details-section">
                    <div v-if="application.applicant_introduction" class="application-field">
                      <label>申请者介绍：</label>
                      <div class="field-content">{{ application.applicant_introduction }}</div>
                    </div>

                    <div v-if="application.application_reason" class="application-field">
                      <label>申请理由：</label>
                      <div class="field-content">{{ application.application_reason }}</div>
                    </div>

                    <div v-if="application.attachment_url" class="application-field attachment-field">
                      <label>附件：</label>
                      <div class="attachment-container">
                        <a
                          :href="application.attachment_url"
                          target="_blank"
                          class="attachment-link"
                        >
                          <span class="link-icon">📎</span>
                          查看附件
                        </a>
                      </div>
                    </div>

                    <div v-if="application.review_comment" class="application-field">
                      <label>审核意见：</label>
                      <div class="review-comment">{{ application.review_comment }}</div>
                    </div>
                  </div>

                  <div class="application-footer">
                    <div class="application-id">申请ID: {{ application.id }}</div>
                    <div class="application-actions">
                      <!-- 只有待审核状态才显示取消按钮 -->
                      <button
                        v-if="application.status === 'pending'"
                        class="cancel-btn"
                        @click.stop="cancelApplication(application.id)"
                        :disabled="cancellingApplicationId === application.id"
                      >
                        <span v-if="cancellingApplicationId === application.id" class="cancelling-spinner"></span>
                        {{ cancellingApplicationId === application.id ? '取消中...' : '取消申请' }}
                      </button>
                    </div>
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
import { ref, computed, watch } from 'vue'

const props = defineProps({
  isVisible: Boolean
})

const emit = defineEmits(['close'])

// API 基础URL
const API_BASE = 'http://localhost:8080/api'

// 响应式数据
const loading = ref(false)
const applications = ref([])
const expandedApplications = ref([]) // 存储展开的申请ID
const cancellingApplicationId = ref(null) // 正在取消的申请ID

// 计算属性
const totalApplications = computed(() => applications.value.length)
const pendingApplications = computed(() =>
  applications.value.filter(app => app.status === 'pending').length
)
const processedApplications = computed(() =>
  applications.value.filter(app => app.status !== 'pending').length
)

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

// 获取申请列表
async function fetchApplications() {
  const token = getToken()
  if (!token) {
    alert('请先登录')
    return
  }

  loading.value = true
  applications.value = []
  expandedApplications.value = [] // 清空展开状态
  cancellingApplicationId.value = null // 重置取消状态

  try {
    console.log('开始获取申请列表...')

    const response = await fetch(`${API_BASE}/organization/application/my-applications`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Accept': 'application/json'
      }
    })

    console.log('申请列表响应状态:', response.status, response.statusText)

    // 获取原始响应文本
    const responseText = await response.text()
    console.log('申请列表原始响应文本:', responseText)

    // 尝试解析JSON
    let data
    try {
      if (responseText.trim()) {
        data = JSON.parse(responseText)
        console.log('申请列表解析成功:', data)
      } else {
        console.log('响应为空，使用默认数据')
        data = { success: false, message: '响应为空' }
      }
    } catch (parseError) {
      console.error('JSON解析失败:', parseError)
      console.error('无法解析的文本:', responseText)
      throw new Error('服务器返回格式错误: ' + responseText.substring(0, 200))
    }

    // 处理响应
    if (response.ok) {
      if (data.success === true) {
        applications.value = data.data || []
        console.log(`成功加载 ${applications.value.length} 条申请记录`)

        // 打印每条申请的详细信息
        applications.value.forEach((app, index) => {
          console.log(`申请 ${index + 1}:`, app)
        })
      } else {
        console.warn('API返回success为false:', data.message)
        alert(data.message || '获取申请列表失败')
      }
    } else {
      console.error('HTTP错误:', response.status, data)
      alert(data.message || `获取申请列表失败: ${response.status}`)
    }

  } catch (error) {
    console.error('获取申请列表失败:', error)

    if (error.message.includes('Failed to fetch')) {
      alert('网络连接失败，请检查网络连接')
    } else if (error.message.includes('未找到认证令牌')) {
      alert('请先登录')
    } else {
      alert(`获取失败: ${error.message}`)
    }
  } finally {
    loading.value = false
    console.log('获取申请列表完成')
  }
}

// 取消申请
async function cancelApplication(applicationId) {
  const token = getToken()
  if (!token) {
    alert('请先登录')
    return
  }

  // 确认操作
  if (!confirm('确定要取消这条申请吗？此操作不可撤销。')) {
    return
  }

  cancellingApplicationId.value = applicationId

  try {
    console.log(`开始取消申请 ${applicationId}...`)

    const response = await fetch(`${API_BASE}/organization/application/${applicationId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      }
    })

    console.log('取消申请响应状态:', response.status, response.statusText)

    // 获取原始响应文本
    const responseText = await response.text()
    console.log('取消申请原始响应文本:', responseText)

    // 尝试解析JSON
    let data
    try {
      if (responseText.trim()) {
        data = JSON.parse(responseText)
        console.log('取消申请解析成功:', data)
      } else {
        console.log('响应为空')
        data = {}
      }
    } catch (parseError) {
      console.error('JSON解析失败:', parseError)
      throw new Error('服务器返回格式错误')
    }

    // 处理响应
    if (response.ok) {
      if (data.success === true || response.status === 200 || response.status === 204) {
        alert('申请已成功取消')

        // 更新本地数据：将申请状态改为已取消
        const index = applications.value.findIndex(app => app.id === applicationId)
        if (index !== -1) {
          applications.value[index].status = 'cancelled'

          // 如果该申请是展开的，关闭它
          const expandedIndex = expandedApplications.value.indexOf(applicationId)
          if (expandedIndex !== -1) {
            expandedApplications.value.splice(expandedIndex, 1)
          }
        }

        // 刷新列表
        await fetchApplications()
      } else {
        console.warn('API返回success为false:', data.message)
        alert(data.message || '取消申请失败')
      }
    } else {
      console.error('HTTP错误:', response.status, data)
      alert(data.message || `取消申请失败: ${response.status}`)
    }

  } catch (error) {
    console.error('取消申请失败:', error)

    if (error.message.includes('Failed to fetch')) {
      alert('网络连接失败，请检查网络连接')
    } else if (error.message.includes('未找到认证令牌')) {
      alert('请先登录')
    } else {
      alert(`取消失败: ${error.message}`)
    }
  } finally {
    cancellingApplicationId.value = null
    console.log('取消申请操作完成')
  }
}

// 切换申请展开状态
function toggleApplication(applicationId) {
  const index = expandedApplications.value.indexOf(applicationId)
  if (index === -1) {
    // 展开
    expandedApplications.value.push(applicationId)
  } else {
    // 折叠
    expandedApplications.value.splice(index, 1)
  }
}

// 状态相关函数
function getStatusColor(status) {
  const colors = {
    pending: 'status-pending',
    approved: 'status-approved',
    rejected: 'status-rejected',
    cancelled: 'status-cancelled'
  }
  return colors[status] || 'status-unknown'
}

function getStatusText(status) {
  const texts = {
    pending: '审核中',
    approved: '已通过',
    rejected: '已拒绝',
    cancelled: '已取消'
  }
  return texts[status] || '未知状态'
}

// 关闭模态框
function close() {
  emit('close')
}

// 监听模态框显示状态
watch(() => props.isVisible, (newVal) => {
  if (newVal) {
    console.log('开始获取申请列表...')
    fetchApplications()
  }
})
</script>

<style scoped>
/* 模态框基础样式 */
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
  z-index: 3000;
  backdrop-filter: blur(4px);
}

.modal-container {
  background: white;
  border-radius: 20px;
  width: 800px;
  max-width: 90vw;
  max-height: 85vh;
  overflow-y: auto; /* 移除overflow-x来去掉水平滚动条 */
  color: #333;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

/* 移除水平滚动条 */
.modal-container::-webkit-scrollbar {
  width: 0;
  height: 0;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 28px;
  border-bottom: 1px solid #f0f0f0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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
  color: white;
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

.modal-body {
  padding: 28px;
}

/* 加载状态 */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px;
  color: #718096;
}

.loading-spinner.large {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(102, 126, 234, 0.3);
  border-radius: 50%;
  border-top-color: #667eea;
  animation: spin 0.8s linear infinite;
}

.loading-state p {
  margin-top: 16px;
  font-size: 15px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px;
  text-align: center;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 20px;
  opacity: 0.5;
}

.empty-text {
  font-size: 18px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.empty-subtext {
  font-size: 14px;
  color: #718096;
  max-width: 300px;
  line-height: 1.5;
}

/* 统计摘要 */
.stats-summary {
  display: flex;
  justify-content: space-around;
  background: linear-gradient(135deg, #f6f8ff 0%, #f0f4ff 100%);
  border-radius: 16px;
  padding: 20px;
  margin-bottom: 24px;
  border: 1px solid #e2e8f0;
}

.stat-item {
  text-align: center;
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #667eea;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 13px;
  color: #718096;
  font-weight: 500;
}

/* 申请网格 */
.applications-grid {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.application-card {
  background: white;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  transition: all 0.3s;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.application-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  border-color: #cbd5e0;
}

/* 状态边框样式 */
.application-card.status-pending {
  border-left: 4px solid #f6ad55;
}

.application-card.status-approved {
  border-left: 4px solid #48bb78;
}

.application-card.status-rejected {
  border-left: 4px solid #f56565;
}

.application-card.status-cancelled {
  border-left: 4px solid #a0aec0;
}

/* 折叠状态头部 */
.application-header-collapsed {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  cursor: pointer;
  transition: background-color 0.3s;
  user-select: none;
}

.application-header-collapsed:hover {
  background-color: #f8fafc;
}

.collapsed-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.application-type {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #2d3748;
}

.type-icon {
  font-size: 16px;
}

.type-text {
  font-size: 15px;
}

.org-name-collapsed {
  font-size: 14px;
  color: #4a5568;
  font-weight: 500;
}

.application-status {
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin: 0 12px;
}

.status-pending {
  background: #fffaf0;
  color: #dd6b20;
}

.status-approved {
  background: #f0fff4;
  color: #38a169;
}

.status-rejected {
  background: #fff5f5;
  color: #e53e3e;
}

.status-cancelled {
  background: #f7fafc;
  color: #718096;
}

.status-unknown {
  background: #edf2f7;
  color: #4a5568;
}

.expand-icon {
  color: #a0aec0;
  font-size: 14px;
  transition: transform 0.3s;
}

.expanded .expand-icon {
  transform: rotate(90deg);
}

/* 展开的详细信息 */
.application-details {
  border-top: 1px solid #f0f0f0;
  padding: 20px;
  background: #fafafa;
}

.details-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* 申请内容 */
.application-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.application-field label {
  font-weight: 600;
  color: #4a5568;
  font-size: 13px;
  min-width: 80px;
}

.field-content {
  color: #2d3748;
  font-size: 14px;
  line-height: 1.5;
  background: #f7fafc;
  padding: 10px 14px;
  border-radius: 8px;
  border-left: 3px solid #cbd5e0;
}

/* 附件字段特殊样式 */
.attachment-field {
  margin-top: 4px;
}

.attachment-container {
  margin-top: 2px; /* 减小与标签的距离 */
}

.attachment-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #4299e1;
  text-decoration: none;
  font-size: 14px;
  padding: 6px 10px; /* 减小内边距 */
  border-radius: 6px;
  background: #f0f9ff;
  border: 1px solid #bee3f8;
  transition: all 0.3s;
  max-width: fit-content;
}

.attachment-link:hover {
  background: #ebf8ff;
  color: #3182ce;
  border-color: #90cdf4;
  transform: translateY(-1px);
  box-shadow: 0 2px 6px rgba(66, 153, 225, 0.2);
}

.link-icon {
  font-size: 14px;
}

.review-comment {
  color: #2d3748;
  font-size: 14px;
  line-height: 1.5;
  background: #f0f9ff;
  padding: 10px 14px;
  border-radius: 8px;
  border-left: 3px solid #4299e1;
  font-style: italic;
}

/* 申请底部 */
.application-footer {
  margin-top: 20px;
  padding-top: 12px;
  border-top: 1px solid #e2e8f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.application-id {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  letter-spacing: 0.5px;
  font-size: 12px;
  color: #a0aec0;
}

/* 取消按钮样式 */
.application-actions {
  display: flex;
  gap: 8px;
}

.cancel-btn {
  padding: 8px 16px;
  background: linear-gradient(135deg, #f56565 0%, #ed8936 100%);
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  gap: 6px;
  box-shadow: 0 2px 4px rgba(237, 137, 54, 0.3);
}

.cancel-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #e53e3e 0%, #dd6b20 100%);
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(237, 137, 54, 0.4);
}

.cancel-btn:active:not(:disabled) {
  transform: translateY(0);
  box-shadow: 0 1px 2px rgba(237, 137, 54, 0.3);
}

.cancel-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  background: linear-gradient(135deg, #cbd5e0 0%, #a0aec0 100%);
  box-shadow: none;
}

/* 取消加载动画 */
.cancelling-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 0.8s linear infinite;
  display: inline-block;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .modal-container {
    width: 95vw;
    margin: 10px;
  }

  .modal-header {
    padding: 20px 24px;
  }

  .modal-body {
    padding: 24px;
  }

  .stats-summary {
    flex-direction: column;
    gap: 16px;
  }

  .stat-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .application-header-collapsed {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .application-status {
    align-self: flex-start;
    margin: 0;
  }

  .expand-icon {
    position: absolute;
    right: 20px;
    top: 16px;
  }

  .application-footer {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
}

@media (max-width: 480px) {
  .modal-header h3 {
    font-size: 18px;
  }

  .empty-icon {
    font-size: 48px;
  }

  .empty-text {
    font-size: 16px;
  }

  .application-details {
    padding: 16px;
  }

  .field-content,
  .review-comment {
    padding: 8px 12px;
    font-size: 13px;
  }

  .attachment-link {
    padding: 5px 8px;
    font-size: 13px;
  }

  .cancel-btn {
    padding: 6px 12px;
    font-size: 12px;
  }
}
</style>
