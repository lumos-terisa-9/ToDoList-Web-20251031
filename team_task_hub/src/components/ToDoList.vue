<template>
  <div class="frosted-glass-panel">
    <h3 class="panel-title">{{ title }}</h3>

    <!-- 显示组织列表（仅组织待办） -->
    <div v-if="title === '组织待办' && activeOrganizations.length > 0" class="organizations-list">
      <div class="organizations-header">
        <span class="organizations-title">今日有活动的组织:</span>
        <span class="organizations-count">{{ activeOrganizations.length }}个</span>
      </div>
      <div class="organizations-container">
        <div v-for="org in activeOrganizations" :key="org.id" class="organization-item">
          <img
            :src="ensureGitHubAvatarUrl(org.logo_url)"
            :alt="org.name"
            class="organization-avatar"
            @error="handleAvatarError"
          />
          <span class="organization-name">{{ org.name }}</span>
        </div>
      </div>
    </div>

    <!-- 条件渲染：只有 showInput 为 true 时才显示输入框 -->
    <div v-if="showInput" class="task-input-group">
      <input
        type="text"
        v-model="newTaskText"
        @keyup.enter="requestNewTaskModal" placeholder="手动输入新的待办事项..."
        class="task-input"
      />
      <button @click="requestNewTaskModal" class="add-btn"> <span class="plus-icon">+</span>
      </button>
    </div>

    <div class="task-list-container">
      <ul v-if="filteredTasks.length > 0" class="task-list">
        <li v-for="task in filteredTasks" :key="task.id" class="task-item"
            :class="{ 'cancelled': task.status === 'cancelled' }">
          <!-- 修改任务复选框部分 -->
          <div class="task-checkbox"
               :class="{
                 'completed': task.status === 'completed',
                 'cancelled': task.status === 'cancelled',
                 'clickable': displayMode === 'today' && !task.isComingStart && !task.isComingEnd && task.status !== 'cancelled' && title !== '组织待办',
                 'non-clickable': title === '组织待办' || task.isComingStart || task.isComingEnd
               }"
               @click.stop="displayMode === 'today' && !task.isComingStart && !task.isComingEnd && task.status !== 'cancelled' && title !== '组织待办' ? toggleTaskComplete(task) : null">
            <span v-if="task.status === 'completed'" class="checkmark">✓</span>
            <span v-else-if="task.status === 'cancelled'" class="cancel-mark">×</span>
          </div>
          <span class="task-text"
                :class="{
              'completed': task.status === 'completed',
              'cancelled': task.status === 'cancelled'
            }"
                @click="handleTaskClick(task)">
        {{ task.title }}
            <!-- 如果是组织任务，显示组织信息 -->
        <span v-if="task.organization" class="organization-info">
          <img
            :src="ensureGitHubAvatarUrl(task.organization.logo_url)"
            :alt="task.organization.name"
            class="task-organization-avatar"
            @error="handleAvatarError"
          />
          <span class="task-organization-name">{{ task.organization.name }}</span>
        </span>
      </span>
        </li>
      </ul>
      <p v-else class="no-tasks">
        {{ getEmptyMessage() }}
      </p>
    </div>

    <!-- 修改即将开始/结束的代办部分 -->
    <div v-if="displayMode === 'today'" class="upcoming-tasks-section">
      <h4 class="upcoming-title">即将开始/结束的代办</h4>
      <div class="upcoming-list-container">
        <ul v-if="filteredUpcomingTasks.length > 0" class="upcoming-list">
          <li v-for="task in filteredUpcomingTasks" :key="task.id" class="upcoming-item">
            <div class="upcoming-checkbox"
                 :class="{
                   'completed': task.status === 'completed',
                   'non-clickable': true
                 }">
              <span v-if="task.status === 'completed'" class="upcoming-checkmark">✓</span>
            </div>
            <span class="upcoming-text" :class="{ 'completed': task.status === 'completed' }"
                  @click="handleTaskClick(task)">
          {{ task.title }}
              <!-- 如果是组织任务，显示组织信息 -->
          <span v-if="task.organization" class="organization-info">
            <img
              :src="ensureGitHubAvatarUrl(task.organization.logo_url)"
              :alt="task.organization.name"
              class="task-organization-avatar"
              @error="handleAvatarError"
            />
            <span class="task-organization-name">{{ task.organization.name }}</span>
          </span>
        </span>
            <!-- 修改标识符号显示 -->
            <span v-if="task.isComingStart" class="tag new-tag">new</span>
            <span v-else-if="task.isComingEnd" class="tag alert-tag">!!</span>
          </li>
        </ul>
        <p v-else class="no-upcoming-tasks">暂无即将开始/结束的代办</p>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      加载中...
    </div>

    <!-- 错误状态 -->
    <div v-if="error" class="error-state">
      <p>暂时无法加载待办事项</p>
      <button @click="refreshFromAPI" class="retry-btn">重试</button>
    </div>
  </div>
</template>

<script setup>
import { ref, defineProps, defineEmits, onMounted, computed, watch } from 'vue';

const emit = defineEmits(['request-modal', 'add-task-object', 'refresh-todos', 'edit-task', 'open-activity-modal']);

const props = defineProps({
  id: {
    type: [String, Number, Date],
    required: true
  },
  title: {
    type: String,
    default: '我的待办清单'
  },
  date: {
    type: Date,
    required: true
  },
  showInput: {
    type: Boolean,
    default: true
  },
  displayMode: {
    type: String,
    default: 'today'
  }
});

const API_BASE = 'http://localhost:8080/api'

const newTaskText = ref('');
const apiTasks = ref([]);
const externalTasks = ref([]); // 外部传入的任务数据
const upcomingTasks = ref([]); // 即将开始/结束的代办
const activeOrganizations = ref([]); // 今日有活动的组织列表
const loading = ref(false);
const error = ref(null);

// GitHub配置
const GITHUB_CONFIG = {
}

// 获取空状态消息
function getEmptyMessage() {
  switch (props.displayMode) {
    case 'today':
      return '🎉 今天没有待办事项！';
    case 'future':
      return '🎉 该日期没有开始的待办事项';
    case 'completed':
      return '🎉 该日期没有已完成的事项';
    default:
      return '🎉 没有待办事项！';
  }
}

// 获取token的通用函数
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
  return token
}

// 根据创建者ID过滤任务：个人待办显示个人创建的任务，组织待办显示组织创建的任务
const filteredTasks = computed(() => {
  // 优先使用外部传入的数据，如果没有则使用apiTasks
  const tasksToFilter = externalTasks.value.length > 0 ? externalTasks.value : apiTasks.value;

  console.log('当前显示模式:', props.displayMode, '任务数量:', tasksToFilter.length);
  console.log('任务数据:', tasksToFilter);

  if (props.title === '个人待办') {
    // 个人待办：creator_organ_id为0且creator_user_id不为0，或者organization_id不存在
    return tasksToFilter.filter(task =>
      (task.creator_organ_id === 0 && task.creator_user_id !== 0) ||
      (!task.organization_id && !task.organization)
    );
  } else if (props.title === '组织待办') {
    // 组织待办：creator_organ_id不为0 或 organization_id存在 或 organization对象存在
    return tasksToFilter.filter(task =>
      task.creator_organ_id !== 0 ||
      (task.organization_id && task.organization_id !== 0) ||
      task.organization
    );
  }
  return tasksToFilter;
});

// 根据创建者ID过滤并排序即将开始/结束的任务
const filteredUpcomingTasks = computed(() => {
  let filteredTasks = [];

  if (props.title === '个人待办') {
    // 个人待办：creator_organ_id为0且creator_user_id不为0，或者organization_id不存在
    filteredTasks = upcomingTasks.value.filter(task =>
      (task.creator_organ_id === 0 && task.creator_user_id !== 0) ||
      (!task.organization_id && !task.organization)
    );
  } else if (props.title === '组织待办') {
    // 组织待办：creator_organ_id不为0 或 organization_id存在 或 organization对象存在
    filteredTasks = upcomingTasks.value.filter(task =>
      task.creator_organ_id !== 0 ||
      (task.organization_id && task.organization_id !== 0) ||
      task.organization
    );
  } else {
    filteredTasks = upcomingTasks.value;
  }

  // 按时间排序：时间近的排在上面
  return filteredTasks.sort((a, b) => {
    const timeA = new Date(a.sortTime || a.start_time || a.startTime || a.end_time || a.endTime || a.createdAt);
    const timeB = new Date(b.sortTime || b.start_time || b.startTime || b.end_time || b.endTime || b.createdAt);
    return timeA - timeB;
  });
});

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

// 获取今日有活动的组织列表（仅用于组织待办）
async function loadActiveOrganizations() {
  if (props.title !== '组织待办') return;

  const token = getToken()
  if (!token) {
    console.error('未找到认证令牌')
    return
  }

  try {
    console.log('开始调用获取今日有活动的组织列表接口...');
    const response = await fetch(`${API_BASE}/todos/organizations/today`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    if (response.ok) {
      const result = await response.json()
      console.log('获取今日有活动的组织列表接口返回的全部内容:', result)
      if (result.success && result.organizations) {
        activeOrganizations.value = result.organizations
        console.log('今日有活动的组织列表:', activeOrganizations.value)
      } else {
        console.warn('获取组织列表返回数据格式异常:', result)
      }
    } else {
      const errorText = await response.text()
      console.error('获取组织列表失败:', response.status, errorText)
    }
  } catch (error) {
    console.error('调用获取组织列表接口失败:', error)
  }
}

// 更新组织列表的方法（由父组件调用）
function updateOrganizations(organizations) {
  console.log('更新组织列表数据:', organizations);
  activeOrganizations.value = organizations || [];
  console.log('更新后的组织列表:', activeOrganizations.value);
}

// 从API加载任务（现在主要用于刷新）
async function loadTasksFromAPI() {
  const token = getToken()
  if (!token) {
    console.error('未找到认证令牌')
    return
  }

  loading.value = true
  error.value = null
  try {
    const response = await fetch(`${API_BASE}/todos/todayTodos`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    if (response.ok) {
      const result = await response.json()
      console.log('获取待办成功:', result)
      if (result.success && result.todos) {
        // 确保每个任务都有 completed 属性
        apiTasks.value = result.todos.map(task => ({
          ...task,
          completed: task.completed || false
        }))
      } else {
        console.warn('获取待办返回数据格式异常:', result)
      }
    } else {
      const errorText = await response.text()
      console.error('获取待办失败:', response.status, errorText)
      error.value = `加载失败: ${response.status} ${errorText}`
    }
  } catch (error) {
    console.error('调用待办接口失败:', error)
    error.value = '网络错误，请检查连接'
  } finally {
    loading.value = false
  }
}

// 完成任务
async function completeTask(task) {
  const token = getToken()
  if (!token) {
    console.error('未找到认证令牌')
    return false
  }

  try {
    // 直接使用点击的任务数据
    const requestBody = {
      "description": task.description || task.content || "",
      "end_time": task.end_time || task.endTime || "",
      "start_time": task.start_time || task.startTime || "",
      "title": task.title
    };

    console.log('完成任务请求参数:', requestBody);

    const response = await fetch(`${API_BASE}/todos/complete`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(requestBody)
    })

    const result = await response.json();
    console.log('完成任务响应:', result);

    if (response.ok && result.success) {
      console.log('完成任务成功:', result.message);
      return true;
    } else {
      console.error('完成任务失败:', result.message);
      return false;
    }
  } catch (error) {
    console.error('调用完成任务接口失败:', error);
    return false;
  }
}

// 取消完成任务
async function cancelCompletedTask(task) {
  const token = getToken()
  if (!token) {
    console.error('未找到认证令牌')
    return false
  }

  try {
    // 直接使用点击的任务数据
    const requestBody = {
      "description": task.description || task.content || "",
      "end_time": task.end_time || task.endTime || "",
      "start_time": task.start_time || task.startTime || "",
      "title": task.title
    };

    console.log('取消完成任务请求参数:', requestBody);

    const response = await fetch(`${API_BASE}/todos/cancel-completedTodo`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(requestBody)
    })

    const result = await response.json();
    console.log('取消完成任务响应:', result);

    if (response.ok && result.success) {
      console.log('取消完成任务成功:', result.message);
      return true;
    } else {
      console.error('取消完成任务失败:', result.message);
      return false;
    }
  } catch (error) {
    console.error('调用取消完成任务接口失败:', error);
    return false;
  }
}

// 切换任务完成状态
async function toggleTaskComplete(task) {
  // 组织待办不可更改，只可查看
  if (props.title === '组织待办') {
    console.log('组织待办不可更改，只可查看');
    return;
  }

  console.log('切换任务状态:', task.title, '当前状态:', task.status);
  console.log('任务详细信息:', task);

  // 根据status字段判断当前状态
  const originalStatus = task.status;
  const newStatus = task.status === 'completed' ? 'pending' : 'completed';
  task.status = newStatus;

  try {
    let success = false;
    if (newStatus === 'completed') {
      // 如果新状态是完成，调用完成任务接口
      success = await completeTask(task);
    } else {
      // 如果新状态是未完成，调用取消完成任务接口
      success = await cancelCompletedTask(task);
    }

    if (!success) {
      // 如果接口调用失败，回退到原始状态
      task.status = originalStatus;
      console.log('接口调用失败，状态已回退');
    } else {
      // 接口调用成功，但不刷新整个列表
      console.log('任务状态更新成功，保持显示');
    }
  } catch (error) {
    console.error('切换任务状态失败:', error);
    // 如果接口调用失败，回退到原始状态
    task.status = originalStatus;
  }
}

// 添加更新任务的方法
function updateTasks(tasks) {
  console.log('更新任务数据:', tasks);
  // 清空外部任务数据
  externalTasks.value = [];
  // 添加新的任务数据
  externalTasks.value = tasks.map(task => ({
    ...task,
    completed: task.completed || false
  }));
  console.log('更新后的externalTasks:', externalTasks.value);
}

// 更新即将开始的任务
function updateUpcomingTasks(tasks) {
  console.log('更新即将开始任务数据:', tasks);
  upcomingTasks.value = tasks.map(task => ({
    ...task,
    completed: task.completed || false
  }));
}

// 打开编辑模态框
function openEditModal(task) {
  // 组织待办只可查看，不可编辑
  if (props.title === '组织待办') {
    console.log('组织待办只可查看，不可编辑');
    return;
  }

  console.log('打开编辑模态框:', task);
  emit('edit-task', task);
}

// 添加待办事项 -> 请求打开模态框
function requestNewTaskModal() {
  const text = newTaskText.value.trim()
  // 触发事件，让父组件打开创建待办模态框
  emit('request-modal', text);
  // 清空输入框
  newTaskText.value = '';
}

// 刷新API数据
function refreshFromAPI() {
  loadTasksFromAPI()
}

// 监听日期变化，如果是今天，重新加载组织列表
watch(() => props.date, (newDate) => {
  const today = new Date();
  const isToday = newDate.getDate() === today.getDate() &&
    newDate.getMonth() === today.getMonth() &&
    newDate.getFullYear() === today.getFullYear();

  if (props.title === '组织待办' && isToday) {
    loadActiveOrganizations();
  }
});

defineExpose({
  refreshFromAPI,
  updateTasks,
  updateUpcomingTasks,
  updateOrganizations,
  loadActiveOrganizations
})

onMounted(() => {
  // 如果是组织待办，加载活动组织列表
  if (props.title === '组织待办') {
    loadActiveOrganizations();
  }
})

// 新增函数：处理任务点击
function handleTaskClick(task) {
  // 组织待办点击应该打开活动详情弹窗
  if (props.title === '组织待办') {
    console.log('打开组织活动详情:', task);
    // 发射事件给父组件
    emit('open-activity-modal', task);
  } else {
    // 个人待办保持原有逻辑
    console.log('打开编辑模态框:', task);
    emit('edit-task', task);
  }
}
</script>

<style scoped>
/* 保持原有的样式不变，只添加新样式 */
.frosted-glass-panel {
  background-color: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 4px 30px rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  padding: 12px; /* 减小内边距 */
  color: #fff;
  margin-top: 12px; /* 减小外边距 */
  z-index: 10;
  pointer-events: auto;
}

.panel-title {
  font-size: 1.3rem; /* 减小字体大小 */
  font-weight: 600;
  margin-top: 0;
  margin-bottom: 16px; /* 减小间距 */
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
}

/* 组织列表样式 */
.organizations-list {
  margin-bottom: 16px;
  padding: 12px;
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.organizations-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.organizations-title {
  font-size: 0.9rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
}

.organizations-count {
  font-size: 0.8rem;
  color: #4ecdc4;
  background-color: rgba(78, 205, 196, 0.2);
  padding: 2px 6px;
  border-radius: 10px;
}

.organizations-container {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.organization-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background-color: rgba(255, 255, 255, 0.15);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;
}

.organization-item:hover {
  background-color: rgba(255, 255, 255, 0.25);
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.organization-avatar {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  object-fit: cover;
  border: 1.5px solid rgba(255, 255, 255, 0.4);
}

.organization-name {
  font-size: 0.85rem;
  color: #fff;
  font-weight: 500;
}

/* 任务中的组织信息样式 */
.organization-info {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-left: 8px;
  padding: 2px 6px;
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  font-size: 0.75rem;
}

.task-organization-avatar {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.task-organization-name {
  color: rgba(255, 255, 255, 0.8);
}

.task-input-group {
  display: flex;
  margin-bottom: 16px; /* 减小间距 */
}

.task-input {
  flex-grow: 1;
  padding: 8px 12px; /* 减小内边距 */
  border: none;
  background-color: rgba(255, 255, 255, 0.85);
  border-radius: 6px 0 0 6px; /* 减小圆角 */
  font-size: 0.9rem; /* 减小字体大小 */
  color: #333;
  outline: none;
  transition: background-color 0.3s;
}

.task-input::placeholder {
  color: #888;
  font-size: 0.9rem; /* 减小占位符字体大小 */
}

.add-btn {
  background-color: #007aff;
  color: white;
  border: none;
  padding: 0 12px; /* 减小内边距 */
  border-radius: 0 6px 6px 0; /* 减小圆角 */
  cursor: pointer;
  transition: background-color 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.add-btn:hover {
  background-color: #005bb5;
}

.plus-icon {
  font-size: 1.3rem; /* 减小图标大小 */
  font-weight: bold;
  line-height: 1;
}

.task-list-container {
  max-height: 200px; /* 减小高度 */
  overflow-y: auto;
  padding-right: 6px; /* 减小内边距 */
  margin-bottom: 12px; /* 添加底部间距 */
}

.task-list-container::-webkit-scrollbar {
  width: 4px; /* 减小滚动条宽度 */
}
.task-list-container::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 2px; /* 减小圆角 */
}
.task-list-container::-webkit-scrollbar-thumb {
  background: rgba(44, 43, 43, 0.5);
  border-radius: 2px; /* 减小圆角 */
}

.task-list {
  list-style-type: none;
  padding: 0;
  margin: 0;
}

.task-item {
  display: flex;
  align-items: center;
  padding: 8px 6px; /* 减小内边距 */
  margin-bottom: 6px; /* 减小间距 */
  transition: background-color 0.3s;
  gap: 10px; /* 减小间距 */
  cursor: pointer;
  border-radius: 4px; /* 减小圆角 */
}

.task-item:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

/* 任务复选框 - 减小大小 */
.task-checkbox {
  width: 14px; /* 进一步减小 */
  height: 14px; /* 进一步减小 */
  border: 1.5px solid #8e8e93; /* 减小边框宽度 */
  border-radius: 2px; /* 减小圆角 */
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.2s;
  cursor: pointer;
  background: transparent;
  position: relative;
}

.task-checkbox:hover {
  border-color: #34c759;
}

.task-checkbox.completed {
  background-color: #34c759;
  border-color: #34c759;
}

.checkmark {
  color: white;
  font-weight: bold;
  font-size: 10px; /* 减小对勾大小 */
  line-height: 1;
}

.task-text {
  font-weight: 500;
  line-height: 1.3; /* 减小行高 */
  flex-grow: 1;
  cursor: pointer;
  padding: 2px 0; /* 减小内边距 */
  font-size: 0.9rem; /* 减小字体大小 */
}

.task-text.completed {
  opacity: 0.6;
}

.no-tasks {
  text-align: center;
  font-style: italic;
  opacity: 0.7;
  padding: 16px 0; /* 减小内边距 */
  font-size: 0.9rem; /* 减小字体大小 */
}

/* 即将开始/结束的代办样式 */
.upcoming-tasks-section {
  margin-top: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.2);
  padding-top: 12px;
}

.upcoming-title {
  font-size: 1.1rem; /* 比主标题稍小 */
  font-weight: 600;
  margin-bottom: 10px;
  color: rgba(255, 255, 255, 0.9);
}

.upcoming-list-container {
  max-height: 150px; /* 比主列表稍矮 */
  overflow-y: auto;
  padding-right: 6px;
}

.upcoming-list-container::-webkit-scrollbar {
  width: 4px;
}
.upcoming-list-container::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 2px;
}
.upcoming-list-container::-webkit-scrollbar-thumb {
  background: rgba(44, 43, 43, 0.5);
  border-radius: 2px;
}

.upcoming-list {
  list-style-type: none;
  padding: 0;
  margin: 0;
}

.upcoming-item {
  display: flex;
  align-items: center;
  padding: 6px 4px; /* 比主列表更紧凑 */
  margin-bottom: 4px;
  transition: background-color 0.3s;
  gap: 8px;
  cursor: pointer;
  border-radius: 3px;
}

.upcoming-item:hover {
  background-color: rgba(255, 255, 255, 0.08);
}

.upcoming-checkbox {
  width: 12px; /* 更小的复选框 */
  height: 12px;
  border: 1.5px solid #8e8e93;
  border-radius: 2px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.2s;
  cursor: pointer;
  background: transparent;
}

.upcoming-checkbox:hover {
  border-color: #34c759;
}

.upcoming-checkbox.completed {
  background-color: #34c759;
  border-color: #34c759;
}

.upcoming-checkmark {
  color: white;
  font-weight: bold;
  font-size: 9px;
  line-height: 1;
}

.upcoming-text {
  font-weight: 500;
  line-height: 1.2;
  flex-grow: 1;
  cursor: pointer;
  padding: 1px 0;
  font-size: 0.85rem; /* 更小的字体 */
}

.upcoming-text.completed {
  opacity: 0.6;
}

.no-upcoming-tasks {
  text-align: center;
  font-style: italic;
  opacity: 0.7;
  padding: 12px 0;
  font-size: 0.85rem; /* 更小的字体 */
}

.loading-state {
  text-align: center;
  padding: 8px; /* 减小内边距 */
  opacity: 0.7;
  font-size: 0.9rem; /* 减小字体大小 */
}

.error-state {
  text-align: center;
  padding: 16px; /* 减小内边距 */
  color: #ff6b6b;
  font-size: 0.9rem; /* 减小字体大小 */
}

.retry-btn {
  background: #4ecdc4;
  color: white;
  border: none;
  padding: 6px 12px; /* 减小内边距 */
  border-radius: 4px; /* 减小圆角 */
  cursor: pointer;
  margin-top: 8px; /* 减小间距 */
  font-size: 0.85rem; /* 减小字体大小 */
}

.retry-btn:hover {
  background: #45b7af;
}

.tag {
  margin-left: 8px;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: bold;
}

.new-tag {
  background-color: #007bff;
  color: white;
}

.alert-tag {
  background-color: #dc3545;
  color: white;
}

/* 确保代办项布局正确 */
.upcoming-item {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.upcoming-text {
  flex: 1;
}

/* 可点击的复选框样式 */
.task-checkbox.clickable {
  cursor: pointer;
}

.task-checkbox.clickable:hover {
  background-color: #f0f0f0;
}

/* 不可点击的复选框样式 */
.task-checkbox:not(.clickable),
.upcoming-checkbox.non-clickable {
  cursor: not-allowed;
  opacity: 0.6;
}

/* 取消待办样式 */
.task-item.cancelled {
  opacity: 0.6;
}

.task-checkbox.cancelled {
  background-color: #6c757d;
  cursor: not-allowed;
  text-decoration: line-through;
}

.task-text.cancelled {
  color: #6c757d;
  text-decoration: line-through;
}

.cancel-mark {
  color: white;
  font-weight: bold;
}

/* 模态框底部按钮布局 */
.modal-footer {
  display: flex;
  justify-content: space-between;
  padding: 16px;
  border-top: 1px solid #e0e0e0;
}

.footer-btn {
  flex: 1;
  margin: 0 4px;
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.update-btn {
  background-color: #007bff;
  color: white;
}

.cancel-btn {
  background-color: #ff9500; /* 橙色 */
  color: white;
}

.complete-btn {
  background-color: #28a745;
  color: white;
}

.footer-btn:hover {
  opacity: 0.9;
}
</style>
