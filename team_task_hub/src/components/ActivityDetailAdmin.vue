<template>
  <div class="activity-detail-admin">
    <!-- 公共部分 -->
    <ActivityDetailCommon
        ref="commonComponent"
        :activity="activity"
        :org-id="orgId"
        :user-role="userRole"
        @close="$emit('close')"
        @review-submitted="handleReviewSubmitted"
        @review-failed="handleReviewFailed"
    />

    <!-- 管理员特有功能 -->
    <div v-if="activity.status !== 'cancelled'" class="admin-actions">
      <!-- 编辑模式切换 -->
      <div class="mode-switcher">
        <button
            class="mode-btn"
            :class="{ 'active': activeMode === 'view' }"
            @click="activeMode = 'view'"
        >
          <span>👁️</span>
          查看模式
        </button>

        <button
            class="mode-btn"
            :class="{ 'active': activeMode === 'edit' }"
            @click="enterEditMode"
        >
          <span>✏️</span>
          编辑活动
        </button>

        <button
            class="mode-btn"
            :class="{ 'active': activeMode === 'assign' }"
            @click="enterAssignMode"
        >
          <span>👥</span>
          指派人员
        </button>

        <button
            class="mode-btn"
            :class="{ 'active': isSelectMode }"
            @click="toggleSelectMode"
        >
          <span>✅</span>
          标记完成
        </button>
      </div>

      <!-- 编辑活动表单 -->
      <div v-if="activeMode === 'edit'" class="edit-form-section">
        <h3 class="edit-title">编辑活动信息</h3>

        <div class="edit-form">
          <div class="form-row">
            <div class="form-group">
              <label class="form-label">活动标题</label>
              <input
                  v-model="editForm.title"
                  type="text"
                  class="form-input"
                  placeholder="请输入活动标题"
              >
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">活动描述</label>
            <textarea
                v-model="editForm.description"
                class="form-textarea"
                placeholder="请输入活动描述"
                rows="4"
            ></textarea>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">开始时间</label>
              <input
                  v-model="localStartTime"
                  type="datetime-local"
                  class="form-input"
              >
            </div>

            <div class="form-group">
              <label class="form-label">结束时间</label>
              <input
                  v-model="localEndTime"
                  type="datetime-local"
                  class="form-input"
              >
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">参与限制</label>
            <div class="radio-group">
              <label class="radio-option" :class="{ 'selected': editForm.participation_limit === 'public' }">
                <input
                    v-model="editForm.participation_limit"
                    type="radio"
                    value="public"
                    class="radio-input"
                >
                <span class="radio-label">
                  <span class="radio-title">公开活动</span>
                  <span class="radio-description">所有用户可见并可参与</span>
                </span>
              </label>

              <label class="radio-option" :class="{ 'selected': editForm.participation_limit === 'org_only' }">
                <input
                    v-model="editForm.participation_limit"
                    type="radio"
                    value="org_only"
                    class="radio-input"
                >
                <span class="radio-label">
                  <span class="radio-title">组织内部活动</span>
                  <span class="radio-description">仅组织成员可见并可参与</span>
                </span>
              </label>

              <label class="radio-option" :class="{ 'selected': editForm.participation_limit === 'admin_assign' }">
                <input
                    v-model="editForm.participation_limit"
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

          <!-- 修改提示 -->
          <div v-if="!hasChanges && editForm.title" class="no-changes-hint">
            未检测到修改
          </div>

          <div class="form-actions">
            <button class="btn-cancel" @click="cancelEdit">取消</button>
            <button
                class="btn-save"
                @click="saveEdit"
                :disabled="updating || !hasChanges"
                :title="!hasChanges ? '未检测到修改' : ''"
            >
              <span v-if="updating">保存中...</span>
              <span v-else>保存更改</span>
            </button>
          </div>
        </div>
      </div>

      <!-- 指派人员界面 -->
      <div v-if="activeMode === 'assign'" class="assign-section">
        <h3 class="assign-title">指派活动人员</h3>

        <!-- 搜索框 -->
        <div class="search-assign">
          <div class="search-box">
            <input
                v-model="searchQuery"
                type="text"
                class="search-input"
                placeholder="搜索成员（输入用户名、ID或邮箱）..."
                @input="handleSearchInput"
            >
            <span class="search-icon">🔍</span>
            <div class="search-hint" v-if="searching">搜索中...</div>
          </div>
        </div>

        <!-- 成员列表 -->
        <div class="assign-member-list" ref="memberList">
          <div v-if="loadingMembers" class="loading-members">
            <div class="loading-spinner"></div>
            加载成员中...
          </div>

          <div v-else-if="assignableMembers.length === 0" class="empty-members">
            <div class="empty-icon">👥</div>
            <div class="empty-text">未找到成员</div>
          </div>

          <div v-else class="assign-members-grid">
            <div
                v-for="member in assignableMembers"
                :key="member.id"
                class="assign-member-card"
                :class="{ 'already-assigned': isAlreadyAssigned(member.id) }"
                @click="!isAlreadyAssigned(member.id) && toggleAssignSelection(member.id)"
            >
              <div class="member-status">
                <div v-if="isAlreadyAssigned(member.id)" class="already-assigned-badge">
                  <span>✓</span>
                </div>
                <div v-else class="assign-checkbox" :class="{ 'selected': isSelectedForAssign(member.id) }">
                  <div class="checkbox-icon">✓</div>
                </div>
              </div>

              <div class="member-avatar">
                <img
                    :src="ensureGitHubAvatarUrl(member.avatar_url)"
                    :alt="member.username"
                    @error="handleAvatarError"
                >
              </div>

              <div class="member-info">
                <div class="member-name">{{ member.username }}</div>
                <div class="member-email">{{ member.email }}</div>
                <div class="member-id">ID: {{ member.id }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 批量操作 -->
        <div v-if="selectedAssignMembers.length > 0" class="assign-actions">
          <div class="selected-count">
            已选择 {{ selectedAssignMembers.length }} 名成员
          </div>
          <div class="action-buttons">
            <button class="btn-assign" @click="batchAssign" :disabled="assigning">
              <span v-if="assigning">指派中...</span>
              <span v-else>批量指派</span>
            </button>
          </div>
        </div>
      </div>

      <!-- 选择模式下的参与者列表 -->
      <div v-if="isSelectMode" class="select-participants-section">
        <div class="section-header">
          <div class="section-icon">✅</div>
          <h2 class="section-title">选择参与者标记完成</h2>
          <div class="selected-count">
            已选择 {{ selectedParticipants.length }} 名参与者
          </div>
        </div>

        <div class="select-participants-list">
          <div v-if="loadingParticipants" class="loading-participants">
            <div class="loading-spinner"></div>
            加载参与者中...
          </div>

          <div v-else-if="filteredParticipants.length === 0" class="empty-participants">
            <div class="empty-icon">👤</div>
            <div class="empty-text">暂无参与者</div>
          </div>

          <div v-else class="select-participants-container">
            <div
                v-for="participant in filteredParticipants"
                :key="participant.id"
                class="select-participant-card"
                :class="{ 'completed': isCompleted(participant.id) }"
                @click="toggleParticipantSelection(participant)"
            >
              <!-- 选择框（仅显示未完成的参与者） -->
              <div v-if="!isCompleted(participant.id)"
                   class="participant-checkbox"
                   :class="{ 'selected': isSelected(participant.id) }">
                <div class="checkbox-icon">✓</div>
              </div>

              <!-- 完成标记 -->
              <div v-if="isCompleted(participant.id)" class="completed-badge">
                <span>✅</span>
              </div>

              <!-- 参与者头像 -->
              <div class="participant-avatar">
                <img
                    :src="ensureGitHubAvatarUrl(participant.avatar_url)"
                    :alt="participant.username"
                    @error="handleAvatarError"
                >
              </div>

              <!-- 参与者信息 -->
              <div class="participant-info">
                <div class="participant-name">{{ participant.username }}</div>
                <div class="participant-status">
                  <span v-if="isCompleted(participant.id)" class="completed-text">已完成</span>
                  <span v-else class="participating-text">参与中</span>
                </div>
                <div class="participant-id">ID: {{ participant.id }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 批量操作 -->
        <div v-if="selectedParticipants.length > 0" class="complete-actions">
          <div class="action-buttons">
            <button class="btn-mark-complete" @click="markAsComplete" :disabled="markingComplete">
              <span v-if="markingComplete">标记中...</span>
              <span v-else>标记为已完成</span>
            </button>
            <button class="btn-cancel-select" @click="toggleSelectMode">取消选择</button>
          </div>
        </div>
      </div>

      <!-- 危险操作 -->
      <div class="danger-actions">
        <button
            class="btn-complete-activity"
            @click="showCompleteDialog = true"
            v-if="activity.status === 'active'"
        >
          <span>🏁</span>
          完成活动
        </button>
        <button
            class="btn-cancel-activity"
            @click="showCancelDialog = true"
            v-if="activity.status === 'active'"
        >
          <span>⚠️</span>
          取消活动
        </button>
      </div>
    </div>

    <!-- 完成活动弹窗 -->
    <div v-if="showCompleteDialog" class="complete-dialog-overlay" @click.self="showCompleteDialog = false">
      <div class="complete-dialog">
        <h3 class="dialog-title">完成活动</h3>
        <p class="dialog-warning">
          确认要将此活动标记为已完成吗？<br>
          <span style="color: #ffc864;">此操作不可逆！</span>
        </p>

        <div class="dialog-actions">
          <button class="btn-dialog-cancel" @click="showCompleteDialog = false">
            取消
          </button>
          <button
              class="btn-dialog-confirm complete"
              @click="completeActivity"
              :disabled="completing"
          >
            <span v-if="completing">完成中...</span>
            <span v-else>确认完成活动</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 取消活动弹窗 -->
    <div v-if="showCancelDialog" class="cancel-dialog-overlay" @click.self="showCancelDialog = false">
      <div class="cancel-dialog">
        <h3 class="dialog-title">取消活动</h3>
        <p class="dialog-warning">
          此操作不可逆！请输入取消原因：
        </p>

        <div class="cancel-form">
          <textarea
              v-model="cancelReason"
              class="cancel-textarea"
              placeholder="请输入取消原因..."
              rows="4"
              required
          ></textarea>

          <div v-if="showCancelDialog && !cancelReason.trim()" class="validation-error">
            取消原因不能为空
          </div>

          <div class="dialog-actions">
            <button class="btn-dialog-cancel" @click="showCancelDialog = false">
              取消
            </button>
            <button
                class="btn-dialog-confirm"
                @click="cancelActivity"
                :disabled="!cancelReason.trim() || cancelling"
            >
              <span v-if="cancelling">取消中...</span>
              <span v-else>确认取消活动</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 成功提示 -->
    <div v-if="showSuccessMessage" class="success-message">
      <div class="success-icon">✅</div>
      <div class="success-text">{{ successMessage }}</div>
    </div>

    <div v-if="showErrorMessage" class="error-message">
      <div class="error-icon">❌</div>
      <div class="error-text">{{ errorMessage }}</div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import ActivityDetailCommon from './ActivityDetailCommon.vue'

const props = defineProps({
  activity: {
    type: Object,
    required: true
  },
  orgId: {
    type: [String, Number],
    required: true
  },
  orgName: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['close', 'activity-updated'])

// 引用
const commonComponent = ref(null)

// 模式状态
const activeMode = ref('view')
const isSelectMode = ref(false)

// 参与者数据
const participants = ref([])
const completedUsers = ref([])
const loadingParticipants = ref(false)

// 选择模式相关
const selectedParticipants = ref([])
const markingComplete = ref(false)

// 编辑表单
const editForm = reactive({
  title: '',
  description: '',
  start_time: '',
  end_time: '',
  participation_limit: 'public'
})

const localStartTime = ref('')
const localEndTime = ref('')
const updating = ref(false)

// 原始活动数据（用于比较是否有修改）
const originalActivity = ref(null)

// 计算属性：检查表单是否有修改
const hasChanges = computed(() => {
  if (!originalActivity.value) return false

  const activity = originalActivity.value

  // 比较各个字段
  if (editForm.title !== activity.title) return true

  const currentDescription = editForm.description || ''
  const originalDescription = activity.description || ''
  if (currentDescription !== originalDescription) return true

  if (editForm.start_time !== activity.start_time) return true
  if (editForm.end_time !== activity.end_time) return true

  const currentLimit = editForm.participation_limit || 'public'
  const originalLimit = activity.participation_limit || 'public'
  if (currentLimit !== originalLimit) return true

  return false
})

// 计算属性：过滤出未完成的参与者
const filteredParticipants = computed(() => {
  return participants.value.filter(participant => !isCompleted(participant.id))
})

// 指派相关
const searchQuery = ref('')
const searchTimer = ref(null)
const searching = ref(false)
const loadingMembers = ref(false)
const assignableMembers = ref([])
const selectedAssignMembers = ref([])
const alreadyAssignedMembers = ref([])
const assigning = ref(false)

// 完成活动
const showCompleteDialog = ref(false)
const completing = ref(false)

// 取消活动
const showCancelDialog = ref(false)
const cancelReason = ref('')
const cancelling = ref(false)

// 成功消息
const showSuccessMessage = ref(false)
const successMessage = ref('')
const showErrorMessage = ref(false)
const errorMessage = ref('')

// GitHub配置
const GITHUB_CONFIG = {
}

// 用户角色
const userRole = 'Admin'

// 检查是否已完成
function isCompleted(userId) {
  return completedUsers.value.includes(parseInt(userId))
}

// 检查是否被选中
function isSelected(userId) {
  return selectedParticipants.value.includes(parseInt(userId))
}

// 切换参与者选择
function toggleParticipantSelection(participant) {
  if (isCompleted(participant.id)) return

  const userId = parseInt(participant.id)
  const index = selectedParticipants.value.indexOf(userId)

  if (index === -1) {
    selectedParticipants.value.push(userId)
  } else {
    selectedParticipants.value.splice(index, 1)
  }
}

// 切换选择模式
function toggleSelectMode() {
  isSelectMode.value = !isSelectMode.value
  if (isSelectMode.value) {
    // 进入选择模式时加载参与者数据
    fetchParticipantsForSelect()
  } else {
    // 退出选择模式时清空选择
    selectedParticipants.value = []
  }
  activeMode.value = 'view'
}

// 获取参与者数据
// ActivityDetailAdmin.vue - 修改 fetchParticipantsForSelect() 函数
async function fetchParticipantsForSelect() {
  const token = getToken()
  if (!token) {
    console.error('获取参与者失败：未找到token')
    return
  }

  loadingParticipants.value = true

  try {
    // 获取活动参与者
    const participantsResponse = await fetch(
        `http://localhost:8080/api/organization/activities/${props.activity.id}/participants`,
        {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        }
    )

    console.log('参与者响应状态:', participantsResponse.status)

    if (participantsResponse.ok) {
      // ✅ 直接使用 json()，不要用 text()
      const data = await participantsResponse.json()
      console.log('参与者解析数据:', data)

      if (data.success && data.data) {
        // 根据返回格式获取参与者列表
        if (data.data.participants) {
          participants.value = data.data.participants
        } else if (Array.isArray(data.data)) {
          participants.value = data.data
        } else {
          participants.value = []
        }
        console.log('最终参与者列表:', participants.value)
      } else {
        console.warn('参与者数据格式不正确:', data)
        participants.value = []
      }
    } else {
      console.error('获取参与者请求失败:', participantsResponse.status)
      participants.value = []
    }

    // 如果参与者不为空，获取已完成用户
    if (participants.value.length > 0) {
      const userIds = participants.value.map(p => parseInt(p.id))
      console.log('获取已完成用户，用户ID列表:', userIds)

      const completedResponse = await fetch(
          `http://localhost:8080/api/organization/activities/${props.activity.id}/completed-users`,
          {
            method: 'POST',
            headers: {
              'Authorization': `Bearer ${token}`,
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              user_ids: userIds
            })
          }
      )

      console.log('已完成用户响应状态:', completedResponse.status)

      if (completedResponse.ok) {
        // ✅ 直接使用 json()，不要用 text()
        const completedData = await completedResponse.json()
        console.log('已完成用户数据:', completedData)

        if (completedData.success && completedData.data) {
          // 只提取 completed_user_ids 数组
          completedUsers.value = completedData.data.completed_user_ids || []
          console.log('已完成的用户ID数组:', completedUsers.value)
        } else {
          console.warn('获取已完成用户数据失败:', completedData.message)
          completedUsers.value = []
        }
      } else {
        console.warn('获取已完成用户失败:', completedResponse.status)
        completedUsers.value = []
      }
    }
  } catch (error) {
    console.error('获取参与者数据失败:', error)
    participants.value = []
  } finally {
    loadingParticipants.value = false
  }
}

// 标记为完成
async function markAsComplete() {
  if (selectedParticipants.value.length === 0) {
    showError('请先选择要标记的参与者')
    return
  }

  const token = getToken()
  if (!token) {
    showError('请先登录')
    return
  }

  markingComplete.value = true

  try {
    const response = await fetch(
        `http://localhost:8080/api/organization/${props.orgId}/activities/${props.activity.id}/complete-user`,
        {
          method: 'PATCH',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            user_ids: selectedParticipants.value
          })
        }
    )

    if (response.ok) {
      const data = await response.json()
      if (data.success) {
        // 更新本地状态
        completedUsers.value = [...completedUsers.value, ...selectedParticipants.value]
        selectedParticipants.value = []
        showSuccess('参与者标记为已完成')

        // 刷新公共部分的参与者列表
        if (commonComponent.value) {
          commonComponent.value.fetchParticipants()
        }

        // 🔥 重要：刷新选择模式下的参与者列表
        await fetchParticipantsForSelect()  // 这里需要重新调用
      } else {
        throw new Error(data.message || '标记失败')
      }
    } else {
      throw new Error('标记失败')
    }
  } catch (error) {
    console.error('标记完成失败:', error)
    showError(`标记完成失败: ${error.message}`)
  } finally {
    markingComplete.value = false
  }
}

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

// 获取默认头像URL
function getDefaultAvatarUrl() {
  return `https://${GITHUB_CONFIG.username}.github.io/${GITHUB_CONFIG.folder}/default-avatar.png`
}

// 确保头像URL使用GitHub URL
function ensureGitHubAvatarUrl(avatarUrl) {
  if (!avatarUrl) return getDefaultAvatarUrl()

  if (avatarUrl.includes('github.io') || avatarUrl.includes('githubusercontent.com')) {
    return avatarUrl
  }

  if (avatarUrl.startsWith('blob:') || avatarUrl.startsWith('data:') || !avatarUrl.startsWith('http')) {
    return getDefaultAvatarUrl()
  }

  return avatarUrl
}

// 头像加载失败处理
function handleAvatarError(event) {
  event.target.src = getDefaultAvatarUrl()
}

// 处理评价提交
function handleReviewSubmitted() {
  showSuccess('评价提交成功')
}

function handleReviewFailed(errorMessage) {
  if (errorMessage.includes('参与')) {
    showError('请先参与活动再进行评价')
  } else {
    showError(errorMessage)
  }
}

// 显示成功/错误消息
function showSuccess(message) {
  successMessage.value = message
  showSuccessMessage.value = true

  if (showErrorMessage.value) {
    showErrorMessage.value = false
  }

  setTimeout(() => {
    showSuccessMessage.value = false
  }, 3000)
}

function showError(message) {
  errorMessage.value = message
  showErrorMessage.value = true

  if (showSuccessMessage.value) {
    showSuccessMessage.value = false
  }

  setTimeout(() => {
    showErrorMessage.value = false
  }, 3000)
}

// 进入编辑模式
function enterEditMode() {
  activeMode.value = 'edit'
  isSelectMode.value = false

  // 保存原始活动数据
  originalActivity.value = {
    title: props.activity.title,
    description: props.activity.description || '',
    start_time: props.activity.start_time,
    end_time: props.activity.end_time,
    participation_limit: props.activity.participation_limit || 'public'
  }

  // 填充编辑表单
  editForm.title = props.activity.title
  editForm.description = props.activity.description || ''
  editForm.start_time = props.activity.start_time
  editForm.end_time = props.activity.end_time
  editForm.participation_limit = props.activity.participation_limit || 'public'

  // 设置本地时间输入
  localStartTime.value = convertAPIToDateTimeLocal(props.activity.start_time)
  localEndTime.value = convertAPIToDateTimeLocal(props.activity.end_time)
}

// 取消编辑
function cancelEdit() {
  activeMode.value = 'view'
  originalActivity.value = null
}

// 保存编辑
async function saveEdit() {
  // 如果没有修改，直接返回
  if (!hasChanges.value) {
    console.log('没有检测到修改，不保存')
    return
  }

  const token = getToken()
  if (!token) {
    showError('请先登录')
    return
  }

  // 验证必填项
  if (!editForm.title.trim()) {
    showError('活动标题不能为空')
    return
  }

  // 转换时间格式
  if (localStartTime.value) {
    editForm.start_time = convertDateTimeLocalToAPI(localStartTime.value)
  }

  if (localEndTime.value) {
    editForm.end_time = convertDateTimeLocalToAPI(localEndTime.value)
  }

  // 验证时间格式
  if (!validateAPITimeFormat(editForm.start_time)) {
    showError('开始时间格式不正确')
    return
  }

  if (!validateAPITimeFormat(editForm.end_time)) {
    showError('结束时间格式不正确')
    return
  }

  // 验证时间逻辑
  const start = new Date(editForm.start_time)
  const end = new Date(editForm.end_time)
  if (start >= end) {
    showError('结束时间必须晚于开始时间')
    return
  }

  updating.value = true

  try {
    // 构建更新数据
    const updateData = {
      title: editForm.title,
      description: editForm.description || '',
      start_time: editForm.start_time,
      end_time: editForm.end_time,
      participation_limit: editForm.participation_limit || 'public'
    }

    console.log('发送的更新数据:', updateData)

    const url = `http://localhost:8080/api/organization/${props.orgId}/activities/${props.activity.id}`
    const response = await fetch(url, {
      method: 'PATCH',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(updateData)
    })

    if (response.ok) {
      const data = await response.json()
      if (data.success) {
        activeMode.value = 'view'
        originalActivity.value = null
        showSuccess('活动信息更新成功')
        emit('activity-updated')
      } else {
        throw new Error(data.message || '更新失败')
      }
    } else {
      const errorText = await response.text()
      throw new Error(`HTTP ${response.status}: ${errorText}`)
    }
  } catch (error) {
    console.error('更新活动失败:', error)
    showError(`更新活动失败: ${error.message}`)
  } finally {
    updating.value = false
  }
}

// 从datetime-local转换为API格式
function convertDateTimeLocalToAPI(localDateTime) {
  if (!localDateTime) return ''

  const date = new Date(localDateTime)
  const year = date.getUTCFullYear()
  const month = String(date.getUTCMonth() + 1).padStart(2, '0')
  const day = String(date.getUTCDate()).padStart(2, '0')
  const hours = String(date.getUTCHours()).padStart(2, '0')
  const minutes = String(date.getUTCMinutes()).padStart(2, '0')
  const seconds = '00'

  return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}Z`
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

// 验证时间格式是否为API要求的格式
function validateAPITimeFormat(timeStr) {
  const pattern = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$/
  return pattern.test(timeStr)
}

// 进入指派模式
async function enterAssignMode() {
  activeMode.value = 'assign'
  isSelectMode.value = false

  // 加载已指派成员
  await fetchAlreadyAssignedMembers()

  // 初始搜索所有成员
  searchQuery.value = ''
  await searchMembers('')
}

// 搜索成员
async function searchMembers(namePrefix = '') {
  const token = getToken()
  if (!token) return

  loadingMembers.value = true
  searching.value = true

  try {
    let url = `http://localhost:8080/api/organization/${props.orgId}/users/search`
    if (namePrefix && namePrefix.trim()) {
      const encodedPrefix = encodeURIComponent(namePrefix.trim())
      url += `?name_prefix=${encodedPrefix}`
    }

    console.log('搜索成员URL:', url)

    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    console.log('搜索成员响应状态:', response.status)

    if (response.ok) {
      const data = await response.json()
      console.log('搜索成员响应数据:', data)
      if (data.success && data.data && data.data.users) {
        assignableMembers.value = data.data.users
      } else {
        assignableMembers.value = []
      }
    } else {
      console.error('搜索成员失败:', response.status)
      assignableMembers.value = []
    }
  } catch (error) {
    console.error('搜索成员失败:', error)
    assignableMembers.value = []
  } finally {
    loadingMembers.value = false
    searching.value = false
  }
}

// 搜索输入处理
function handleSearchInput() {
  if (searchTimer.value) {
    clearTimeout(searchTimer.value)
  }

  searching.value = true
  searchTimer.value = setTimeout(() => {
    searchMembers(searchQuery.value.trim())
  }, 300)
}

// 获取已指派成员
async function fetchAlreadyAssignedMembers() {
  const token = getToken()
  if (!token) return

  try {
    const url = `http://localhost:8080/api/organization/activities/${props.activity.id}/participants`
    console.log('获取已指派成员接口调用信息：', {
      url,
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    const response = await fetch(url, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    console.log('已指派成员响应状态：', response.status)
    const responseText = await response.text()
    console.log('已指派成员响应原始数据：', responseText)

    if (response.ok) {
      const data = JSON.parse(responseText)
      console.log('已指派成员解析数据：', data)

      if (data.success && data.data) {
        if (data.data.participants) {
          alreadyAssignedMembers.value = data.data.participants.map(p => parseInt(p.id))
        } else if (Array.isArray(data.data)) {
          alreadyAssignedMembers.value = data.data.map(p => parseInt(p.id))
        } else if (Array.isArray(data)) {
          alreadyAssignedMembers.value = data.map(p => parseInt(p.id))
        } else {
          console.warn('无法识别的参与者数据格式：', data.data)
          alreadyAssignedMembers.value = []
        }
      } else {
        console.warn('获取参与者数据失败：', data.message)
        alreadyAssignedMembers.value = []
      }
    } else {
      console.error('获取参与者请求失败：', response.status)
      alreadyAssignedMembers.value = []
    }
  } catch (error) {
    console.error('获取已指派成员失败:', error)
    alreadyAssignedMembers.value = []
  }
}

// 检查是否已指派
function isAlreadyAssigned(memberId) {
  return alreadyAssignedMembers.value.includes(parseInt(memberId))
}

// 检查是否选中指派
function isSelectedForAssign(memberId) {
  return selectedAssignMembers.value.includes(parseInt(memberId))
}

// 切换指派选择
function toggleAssignSelection(memberId) {
  if (isAlreadyAssigned(memberId)) return

  const userId = parseInt(memberId)
  const index = selectedAssignMembers.value.indexOf(userId)

  if (index === -1) {
    selectedAssignMembers.value.push(userId)
  } else {
    selectedAssignMembers.value.splice(index, 1)
  }
}

// 批量指派
async function batchAssign() {
  if (selectedAssignMembers.value.length === 0) {
    showError('请先选择要指派的成员')
    return
  }

  const token = getToken()
  if (!token) {
    showError('请先登录')
    return
  }

  assigning.value = true

  try {
    const response = await fetch(
        `http://localhost:8080/api/organization/${props.orgId}/activities/batch-assign`,
        {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            activity_id: parseInt(props.activity.id),
            user_ids: selectedAssignMembers.value
          })
        }
    )

    if (response.ok) {
      const data = await response.json()
      if (data.success) {
        // 更新已指派成员列表
        alreadyAssignedMembers.value = [...alreadyAssignedMembers.value, ...selectedAssignMembers.value]
        selectedAssignMembers.value = []
        showSuccess('成员指派成功')

        // 刷新参与者列表
        if (commonComponent.value) {
          commonComponent.value.fetchParticipants()
        }
      } else {
        throw new Error(data.message || '指派失败')
      }
    } else {
      throw new Error('指派失败')
    }
  } catch (error) {
    console.error('批量指派失败:', error)
    showError(`批量指派失败: ${error.message}`)
  } finally {
    assigning.value = false
  }
}

// 完成活动
async function completeActivity() {
  const token = getToken()
  if (!token) {
    showError('请先登录')
    return
  }

  completing.value = true

  try {
    const response = await fetch(
        `http://localhost:8080/api/organization/${props.orgId}/activities/${props.activity.id}/complete-activity`,
        {
          method: 'PATCH',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        }
    )

    console.log('完成活动响应状态:', response.status)
    const responseText = await response.text()
    console.log('完成活动响应原始数据:', responseText)

    if (response.ok) {
      let data
      try {
        data = JSON.parse(responseText)
      } catch (e) {
        console.error('解析响应失败:', e)
        throw new Error('服务器响应格式错误')
      }

      if (data.success) {
        showCompleteDialog.value = false
        showSuccess('活动已标记为已完成')
        emit('activity-updated')
      } else {
        throw new Error(data.message || '完成活动失败')
      }
    } else {
      let errorMessage = '完成活动失败'
      try {
        const errorData = JSON.parse(responseText)
        errorMessage = errorData.message || `HTTP ${response.status}: ${responseText}`
      } catch {
        errorMessage = `HTTP ${response.status}: ${responseText}`
      }
      throw new Error(errorMessage)
    }
  } catch (error) {
    console.error('完成活动失败:', error)
    showError(`完成活动失败: ${error.message}`)
  } finally {
    completing.value = false
  }
}

// 取消活动
async function cancelActivity() {
  if (!cancelReason.value.trim()) {
    showError('请输入取消原因')
    return
  }

  const token = getToken()
  if (!token) {
    showError('请先登录')
    return
  }

  cancelling.value = true

  try {
    const response = await fetch(
        `http://localhost:8080/api/organization/${props.orgId}/activities/${props.activity.id}/cancel`,
        {
          method: 'PATCH',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            reason: cancelReason.value
          })
        }
    )

    console.log('取消活动响应状态:', response.status)
    const responseText = await response.text()
    console.log('取消活动响应原始数据:', responseText)

    if (response.ok) {
      let data
      try {
        data = JSON.parse(responseText)
      } catch (e) {
        console.error('解析响应失败:', e)
        throw new Error('服务器响应格式错误')
      }

      if (data.success) {
        showCancelDialog.value = false
        cancelReason.value = ''
        showSuccess('活动已取消')
        emit('activity-updated')
      } else {
        throw new Error(data.message || '取消失败')
      }
    } else {
      let errorMessage = '取消失败'
      try {
        const errorData = JSON.parse(responseText)
        errorMessage = errorData.message || `HTTP ${response.status}: ${responseText}`
      } catch {
        errorMessage = `HTTP ${response.status}: ${responseText}`
      }
      throw new Error(errorMessage)
    }
  } catch (error) {
    console.error('取消活动失败:', error)
    showError(`取消活动失败: ${error.message}`)
  } finally {
    cancelling.value = false
  }
}
</script>

<style scoped>
.activity-detail-admin {
  position: relative;
}

/* 管理员操作区域 */
.admin-actions {
  margin: 24px;
  padding: 24px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

/* 模式切换 */
.mode-switcher {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.mode-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.mode-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

.mode-btn.active {
  background: rgba(120, 200, 255, 0.15);
  border-color: rgba(120, 200, 255, 0.3);
  color: #78c8ff;
}

.mode-btn.active:hover {
  background: rgba(120, 200, 255, 0.2);
}

/* 编辑表单 */
.edit-form-section {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 24px;
  border: 1px solid rgba(120, 200, 255, 0.2);
}

.edit-title {
  margin: 0 0 20px 0;
  font-size: 18px;
  font-weight: 600;
  color: #78c8ff;
}

.edit-form {
  max-width: 800px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 48px;
  margin-bottom: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.8);
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
}

/* 单选按钮组 */
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
}

.radio-option:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.2);
}

.radio-option.selected {
  background: rgba(120, 200, 255, 0.1);
  border-color: rgba(120, 200, 255, 0.3);
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

/* 修改提示 */
.no-changes-hint {
  position: absolute;
  top: -30px;
  right: 0;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  animation: fadeIn 0.3s ease;
  background: rgba(255, 255, 255, 0.05);
  padding: 4px 12px;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 表单操作按钮 */
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  margin-top: 32px;
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

.btn-save {
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
  min-width: 120px;
}

.btn-save:hover:not(:disabled) {
  background: linear-gradient(135deg,
  rgba(100, 200, 100, 0.25),
  rgba(100, 200, 100, 0.15));
  border-color: rgba(100, 200, 100, 0.5);
  transform: translateY(-2px);
}

.btn-save:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: rgba(255, 255, 255, 0.05) !important;
  border-color: rgba(255, 255, 255, 0.1) !important;
  color: rgba(255, 255, 255, 0.4) !important;
}

/* 选择模式下的参与者列表 */
.select-participants-section {
  background: rgba(255, 255, 255, 0.03);
  border-radius: 20px;
  padding: 24px;
  margin-bottom: 24px;
  border: 1px solid rgba(200, 160, 255, 0.2);
}

.section-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
}

.section-icon {
  font-size: 24px;
}

.section-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #c8a0ff;
}

.selected-count {
  margin-left: auto;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
  background: rgba(200, 160, 255, 0.1);
  padding: 4px 12px;
  border-radius: 8px;
}

.select-participants-list {
  min-height: 200px;
}

.loading-participants {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 60px;
  color: rgba(255, 255, 255, 0.6);
}

.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.1);
  border-top-color: rgba(200, 160, 255, 0.8);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.empty-participants {
  text-align: center;
  padding: 60px 20px;
  color: rgba(255, 255, 255, 0.4);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 12px;
  opacity: 0.3;
}

.empty-text {
  font-size: 16px;
}

/* 选择模式参与者容器 */
.select-participants-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 16px;
  max-height: 300px;
  overflow-y: auto;
  padding-right: 8px;
}

.select-participants-container::-webkit-scrollbar {
  width: 4px;
}

.select-participants-container::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 2px;
}

.select-participants-container::-webkit-scrollbar-thumb {
  background: rgba(200, 160, 255, 0.3);
  border-radius: 2px;
}

/* 选择模式参与者卡片 */
.select-participant-card {
  position: relative;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  transition: all 0.3s ease;
  cursor: pointer;
}

.select-participant-card:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(200, 160, 255, 0.3);
  transform: translateY(-2px);
}

.select-participant-card.completed {
  opacity: 0.8;
  background: rgba(100, 200, 100, 0.08);
  border-color: rgba(100, 200, 100, 0.2);
  cursor: default;
}

/* 选择框 */
.participant-checkbox {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 20px;
  height: 20px;
  border-radius: 6px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  background: rgba(255, 255, 255, 0.05);
  display: grid;
  place-items: center;
  transition: all 0.3s ease;
}

.participant-checkbox.selected {
  background: rgba(200, 160, 255, 0.8);
  border-color: rgba(200, 160, 255, 1);
}

.checkbox-icon {
  font-size: 12px;
  color: white;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.participant-checkbox.selected .checkbox-icon {
  opacity: 1;
}

/* 完成标记 */
.completed-badge {
  position: absolute;
  top: 12px;
  left: 12px;
  font-size: 16px;
  background: rgba(100, 200, 100, 0.2);
  border-radius: 50%;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 参与者头像 */
.participant-avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  overflow: hidden;
  margin-bottom: 12px;
  border: 3px solid rgba(255, 255, 255, 0.1);
}

.select-participant-card.completed .participant-avatar {
  border-color: rgba(100, 200, 100, 0.4);
}

.participant-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* 参与者信息 */
.participant-info {
  text-align: center;
  width: 100%;
}

.participant-name {
  font-weight: 600;
  font-size: 14px;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.participant-status {
  margin-bottom: 4px;
}

.completed-text {
  font-size: 12px;
  color: #64c864;
  padding: 2px 6px;
  background: rgba(100, 200, 100, 0.15);
  border-radius: 8px;
}

.participating-text {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
}

.participant-id {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  font-family: 'SF Mono', monospace;
}

/* 批量操作 */
.complete-actions {
  display: flex;
  justify-content: flex-end;
  padding-top: 20px;
  margin-top: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.action-buttons {
  display: flex;
  gap: 12px;
}

.btn-mark-complete {
  padding: 10px 24px;
  border-radius: 12px;
  border: 1px solid rgba(200, 160, 255, 0.3);
  background: linear-gradient(135deg,
  rgba(200, 160, 255, 0.15),
  rgba(200, 160, 255, 0.08));
  color: #c8a0ff;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-mark-complete:hover:not(:disabled) {
  background: linear-gradient(135deg,
  rgba(200, 160, 255, 0.25),
  rgba(200, 160, 255, 0.15));
  border-color: rgba(200, 160, 255, 0.5);
  transform: translateY(-2px);
}

.btn-mark-complete:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-cancel-select {
  padding: 10px 20px;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.8);
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-cancel-select:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.2);
}

/* 指派区域样式（保持不变） */
.assign-section {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 24px;
  border: 1px solid rgba(200, 160, 255, 0.2);
}

.assign-title {
  margin: 0 0 20px 0;
  font-size: 18px;
  font-weight: 600;
  color: #c8a0ff;
}

.search-assign {
  margin-bottom: 24px;
}

.search-box {
  position: relative;
}

.search-input {
  width: 100%;
  padding: 12px 16px 12px 44px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  color: white;
  font-size: 14px;
  transition: all 0.3s ease;
}

.search-input:focus {
  outline: none;
  border-color: rgba(200, 160, 255, 0.5);
  background: rgba(200, 160, 255, 0.05);
}

.search-icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  color: rgba(255, 255, 255, 0.5);
  font-size: 16px;
}

.search-hint {
  position: absolute;
  right: 16px;
  top: 50%;
  transform: translateY(-50%);
  color: rgba(200, 160, 255, 0.7);
  font-size: 12px;
  animation: pulse 1.5s infinite;
}

.assign-member-list {
  max-height: 400px;
  overflow-y: auto;
  margin-bottom: 20px;
}

.loading-members,
.empty-members {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
  color: rgba(255, 255, 255, 0.6);
}

.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.1);
  border-top-color: rgba(200, 160, 255, 0.8);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-right: 12px;
}

.empty-icon {
  font-size: 32px;
  margin-right: 12px;
  opacity: 0.3;
}

.assign-members-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}

.assign-member-card {
  position: relative;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  transition: all 0.3s ease;
}

.assign-member-card:not(.already-assigned) {
  cursor: pointer;
}

.assign-member-card:not(.already-assigned):hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(200, 160, 255, 0.3);
  transform: translateY(-2px);
}

.assign-member-card.already-assigned {
  opacity: 0.7;
  background: rgba(100, 200, 100, 0.05);
  border-color: rgba(100, 200, 100, 0.2);
}

.member-status {
  flex-shrink: 0;
}

.already-assigned-badge {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: rgba(100, 200, 100, 0.2);
  border: 1px solid rgba(100, 200, 100, 0.4);
  display: grid;
  place-items: center;
  color: #64c864;
  font-size: 12px;
}

.assign-checkbox {
  width: 20px;
  height: 20px;
  border-radius: 6px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  background: rgba(255, 255, 255, 0.05);
  display: grid;
  place-items: center;
  transition: all 0.3s ease;
}

.assign-checkbox.selected {
  background: rgba(200, 160, 255, 0.8);
  border-color: rgba(200, 160, 255, 1);
}

.checkbox-icon {
  font-size: 12px;
  color: white;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.assign-checkbox.selected .checkbox-icon {
  opacity: 1;
}

.member-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
  border: 2px solid rgba(255, 255, 255, 0.1);
}

.member-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.member-info {
  flex: 1;
  min-width: 0;
}

.member-name {
  font-weight: 600;
  font-size: 14px;
  margin-bottom: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.member-email {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.6);
  margin-bottom: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.member-id {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.5);
  font-family: 'SF Mono', monospace;
}

.assign-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  margin-top: 20px;
}

.action-buttons {
  display: flex;
  gap: 12px;
}

.btn-assign {
  padding: 10px 24px;
  border-radius: 12px;
  border: 1px solid rgba(200, 160, 255, 0.3);
  background: linear-gradient(135deg,
  rgba(200, 160, 255, 0.15),
  rgba(200, 160, 255, 0.08));
  color: #c8a0ff;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-assign:hover:not(:disabled) {
  background: linear-gradient(135deg,
  rgba(200, 160, 255, 0.25),
  rgba(200, 160, 255, 0.15));
  border-color: rgba(200, 160, 255, 0.5);
  transform: translateY(-2px);
}

.btn-assign:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 危险操作 */
.danger-actions {
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  gap: 12px;
}

.btn-complete-activity {
  padding: 12px 24px;
  border-radius: 12px;
  border: 1px solid rgba(255, 200, 100, 0.3);
  background: linear-gradient(135deg,
  rgba(255, 200, 100, 0.15),
  rgba(255, 200, 100, 0.08));
  color: #ffc864;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.btn-complete-activity:hover {
  background: linear-gradient(135deg,
  rgba(255, 200, 100, 0.25),
  rgba(255, 200, 100, 0.15));
  border-color: rgba(255, 200, 100, 0.5);
  transform: translateY(-2px);
}

.btn-cancel-activity {
  padding: 12px 24px;
  border-radius: 12px;
  border: 1px solid rgba(255, 100, 100, 0.3);
  background: linear-gradient(135deg,
  rgba(255, 100, 100, 0.15),
  rgba(255, 100, 100, 0.08));
  color: #ff6464;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.btn-cancel-activity:hover {
  background: linear-gradient(135deg,
  rgba(255, 100, 100, 0.25),
  rgba(255, 100, 100, 0.15));
  border-color: rgba(255, 100, 100, 0.5);
  transform: translateY(-2px);
}

/* 弹窗样式 */
.complete-dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(8px);
  display: grid;
  place-items: center;
  z-index: 4000;
  padding: 20px;
}

.complete-dialog {
  width: 100%;
  max-width: 450px;
  background: rgba(30, 35, 40, 0.95);
  border-radius: 20px;
  padding: 32px;
  border: 1px solid rgba(255, 200, 100, 0.2);
  box-shadow: 0 20px 60px rgba(255, 200, 100, 0.1);
}

.cancel-dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(8px);
  display: grid;
  place-items: center;
  z-index: 4000;
  padding: 20px;
}

.cancel-dialog {
  width: 100%;
  max-width: 500px;
  background: rgba(30, 35, 40, 0.95);
  border-radius: 20px;
  padding: 32px;
  border: 1px solid rgba(255, 100, 100, 0.2);
  box-shadow: 0 20px 60px rgba(255, 100, 100, 0.1);
}

.dialog-title {
  margin: 0 0 16px 0;
  font-size: 20px;
  font-weight: 700;
}

.complete-dialog .dialog-title {
  color: #ffc864;
}

.cancel-dialog .dialog-title {
  color: #ff6464;
}

.dialog-warning {
  margin: 0 0 24px 0;
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
  line-height: 1.5;
}

.cancel-textarea {
  width: 100%;
  padding: 16px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 100, 100, 0.3);
  border-radius: 12px;
  color: white;
  font-size: 14px;
  font-family: inherit;
  resize: vertical;
  margin-bottom: 24px;
  transition: all 0.3s ease;
}

.cancel-textarea:focus {
  outline: none;
  border-color: rgba(255, 100, 100, 0.5);
  background: rgba(255, 100, 100, 0.05);
}

.validation-error {
  color: #ff6464;
  font-size: 12px;
  margin-top: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
  animation: fadeIn 0.3s ease;
}

.validation-error::before {
  content: "⚠️";
  font-size: 12px;
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 16px;
}

.btn-dialog-cancel {
  padding: 10px 20px;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.8);
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-dialog-cancel:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.2);
}

.btn-dialog-confirm {
  padding: 10px 24px;
  border-radius: 12px;
  border: 1px solid rgba(255, 100, 100, 0.3);
  background: linear-gradient(135deg,
  rgba(255, 100, 100, 0.2),
  rgba(255, 100, 100, 0.1));
  color: #ff6464;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-dialog-confirm.complete {
  border: 1px solid rgba(255, 200, 100, 0.3);
  background: linear-gradient(135deg,
  rgba(255, 200, 100, 0.2),
  rgba(255, 200, 100, 0.1));
  color: #ffc864;
}

.btn-dialog-confirm.complete:hover:not(:disabled) {
  background: linear-gradient(135deg,
  rgba(255, 200, 100, 0.3),
  rgba(255, 200, 100, 0.2));
  border-color: rgba(255, 200, 100, 0.5);
}

.btn-dialog-confirm:hover:not(:disabled) {
  background: linear-gradient(135deg,
  rgba(255, 100, 100, 0.3),
  rgba(255, 100, 100, 0.2));
  border-color: rgba(255, 100, 100, 0.5);
}

.btn-dialog-confirm:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 成功/错误消息 */
.success-message {
  position: fixed;
  top: 84px;
  right: 24px;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 24px;
  background: rgba(100, 200, 100, 0.15);
  border: 1px solid rgba(100, 200, 100, 0.3);
  border-radius: 12px;
  backdrop-filter: blur(10px);
  z-index: 5000;
  animation: slideIn 0.3s ease, fadeOut 0.3s ease 2.7s forwards;
}

.success-icon {
  font-size: 20px;
}

.success-text {
  font-weight: 600;
  color: #64c864;
}

.error-message {
  position: fixed;
  top: 84px;
  right: 24px;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 24px;
  background: rgba(255, 100, 100, 0.15);
  border: 1px solid rgba(255, 100, 100, 0.3);
  border-radius: 12px;
  backdrop-filter: blur(10px);
  z-index: 5000;
  animation: slideIn 0.3s ease, fadeOut 0.3s ease 2.7s forwards;
}

.error-icon {
  font-size: 20px;
  color: #ff6464;
}

.error-text {
  font-weight: 600;
  color: #ff6464;
}

/* 动画 */
@keyframes pulse {
  0%, 100% {
    opacity: 0.5;
  }
  50% {
    opacity: 1;
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

@keyframes fadeOut {
  from {
    opacity: 1;
  }
  to {
    opacity: 0;
  }
}

/* 响应式调整 */
@media (max-width: 768px) {
  .admin-actions {
    margin: 16px;
    padding: 16px;
  }

  .mode-switcher {
    flex-direction: column;
  }

  .mode-btn {
    width: 100%;
    justify-content: center;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .select-participants-container {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  }

  .assign-members-grid {
    grid-template-columns: 1fr;
  }

  .assign-actions {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .danger-actions {
    flex-direction: column;
  }

  .btn-complete-activity,
  .btn-cancel-activity {
    width: 100%;
    justify-content: center;
  }
}
</style>
