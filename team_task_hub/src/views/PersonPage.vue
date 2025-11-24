<template>
  <div class="person-page">
    <div class="split-pane">
      <MonthCalendar
        v-model="picked"
        @select="onSelect"
        @date-click="handleDateClick"
        class="calendar-pane"
      />

      <div class="task-pane">
        <h2>任务面板</h2>
        <p>已选日期：{{ picked.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' }) }}</p>

        <!-- 个人待办 -->
        <ToDoList
          :id="picked.toDateString()"
          :title="'个人待办'"
          :date="picked"
          :display-mode="currentDisplayMode"
          @request-modal="handleNewTaskRequest"
          @refresh-todos="loadTodayTodos"
          @edit-task="handleEditTask"
          class="todo-list-item"
          ref="personalTodoList"
        />
        <!-- 组织待办 -->
        <ToDoList
          :id="'org-' + picked.toDateString()"
          :title="'组织待办'"
          :date="picked"
          :display-mode="currentDisplayMode"
          :show-input="false"
          @refresh-todos="loadTodayTodos"
          @edit-task="handleEditTask"
          class="todo-list-item"
          ref="orgTodoList"
        />
      </div>
    </div>

    <!-- 个人信息模态框 -->
    <UserProfileModal
      :isVisible="showProfileModal"
      :user="currentUser"
      @close="showProfileModal = false"
      @update-user="handleUserUpdate"
      @logout="handleLogout"
    />

    <!-- 创建待办模态框 -->
    <NewTaskModal
      :isVisible="showModal"
      :date="picked"
      @close="showModal = false"
      @save="handleSaveTask"
    />

    <!-- 编辑待办模态框 - 在全局显示 -->
    <NewTaskModal
      :isVisible="showEditModal"
      :date="picked"
      :isEditMode="true"
      :editTodoData="selectedTask"
      @close="showEditModal = false"
      @update="handleTaskUpdate"
      @complete="handleTaskComplete"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import MonthCalendar from '@/components/MonthCalendar.vue'
import ToDoList from '@/components/ToDoList.vue'
import NewTaskModal from '@/components/NewTaskModal.vue'
import UserProfileModal from '@/components/UserProfileModal.vue'

const router = useRouter()

const picked = ref(new Date())
const showModal = ref(false)
const showEditModal = ref(false) // 添加编辑模态框状态
const showProfileModal = ref(false)
const currentUser = ref(null)
const personalTodoList = ref(null)
const orgTodoList = ref(null)
const selectedTask = ref(null) // 添加选中的任务

const API_BASE = 'http://localhost:8080/api'
// 在 PersonPage.vue 的 script 部分添加数据状态
const todayTasks = ref([]);
const completedTasks = ref([]);
const currentDisplayMode = ref('today'); // 'today' 或 'completed'

// 修改 handleDateClick 函数
function handleDateClick(dateInfo) {
  console.log('日期点击:', dateInfo);

  if (dateInfo.isToday) {
    // 点击今天，调用第二个接口：获取今日待办
    currentDisplayMode.value = 'today';
    loadTodayTodos();
  } else {
    // 点击其他日期，调用第一个接口：获取指定日期完成的待办
    currentDisplayMode.value = 'completed';
    loadCompletedTodos(dateInfo.date);
  }
}

// 修改 loadTodayTodos 函数
async function loadTodayTodos() {
  const token = getToken()
  if (!token) return

  try {
    console.log('开始调用今日待办接口...')
    const response = await fetch(`${API_BASE}/todos/todayTodos`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    if (response.ok) {
      const result = await response.json()
      console.log('今日待办:', result)

      if (result.success && result.todos) {
        todayTasks.value = result.todos;
        // 通知子组件刷新数据
        updateTodoListData();
      }
    } else {
      console.error('获取今日待办失败:', response.status)
    }
  } catch (error) {
    console.error('调用今日待办接口失败:', error)
  }
}

// 修改 loadCompletedTodos 函数
async function loadCompletedTodos(date) {
  const token = getToken()
  if (!token) return

  try {
    // 格式化日期为 YYYY-MM-DD
    const dateStr = date.toISOString().split('T')[0];
    console.log('开始调用已完成待办接口，日期:', dateStr);

    const response = await fetch(`${API_BASE}/todos/completed_todo?date=${dateStr}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    if (response.ok) {
      const result = await response.json()
      console.log('已完成待办:', result)
      if (result.success && result.todos) {
        completedTasks.value = result.todos;
        // 通知子组件刷新数据
        updateTodoListData();
      }
    } else {
      console.error('获取已完成待办失败:', response.status)
    }
  } catch (error) {
    console.error('调用已完成待办接口失败:', error)
  }
}

// 新增函数：更新ToDoList数据
function updateTodoListData() {
  let tasksToDisplay = [];

  if (currentDisplayMode.value === 'today') {
    tasksToDisplay = todayTasks.value;
  } else {
    tasksToDisplay = completedTasks.value;
  }

  // 通知子组件更新数据
  if (personalTodoList.value) {
    personalTodoList.value.updateTasks(tasksToDisplay);
  }
  if (orgTodoList.value) {
    orgTodoList.value.updateTasks(tasksToDisplay);
  }
}
// 获取token的通用函数
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

// 检查登录状态
async function checkAuth() {
  const token = getToken()
  if (!token) {
    router.push('/')
    return
  }

  try {
    const response = await fetch(`${API_BASE}/auth/me`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    if (response.ok) {
      const userData = await response.json()
      currentUser.value = userData
      localStorage.setItem('currentUser', JSON.stringify(userData))

      // 加载今日待办
      await loadTodayTodos()
    } else {
      router.push('/')
    }
  } catch (error) {
    console.error('验证用户失败:', error)
    router.push('/')
  }
}

function onSelect(d) {
  picked.value = d
  console.log('选中日期：', d)
  // 日期变化时重新加载待办
  loadTodayTodos()
}

function handleNewTaskRequest() {
  showModal.value = true
}

// 处理编辑任务
function handleEditTask(task) {
  console.log('处理编辑任务:', task)
  selectedTask.value = task
  showEditModal.value = true
}

// 处理任务更新
function handleTaskUpdate() {
  showEditModal.value = false
  loadTodayTodos() // 重新加载任务
}

// 处理任务完成
function handleTaskComplete() {
  showEditModal.value = false
  loadTodayTodos() // 重新加载任务
}

async function handleSaveTask(taskData) {
  console.log('准备保存任务到日期:', picked.value, taskData)
  showModal.value = false

  // 保存后刷新待办列表
  await nextTick()
  await loadTodayTodos()
}

function handleUserUpdate(updatedUser) {
  currentUser.value = updatedUser
  localStorage.setItem('currentUser', JSON.stringify(updatedUser))
}

function handleLogout() {
  currentUser.value = null
  localStorage.removeItem('token')
  localStorage.removeItem('currentUser')
  router.push('/')
}

// 显示个人信息模态框
function showProfileModalFunc() {
  showProfileModal.value = true
}

// 暴露方法给父组件
defineExpose({
  showProfileModal: showProfileModalFunc
})

onMounted(() => {
  checkAuth()
})
</script>

<style scoped>
.person-page {
  min-height: 100vh;
  background: linear-gradient(160deg, #f9ccfd -10%, #1076fb 95%);
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
}

/* 分栏容器 */
.split-pane {
  display: flex;
  flex: 1;
  align-items: stretch;
  padding: 90px 160px 60px 160px;
  gap: 20px;
  min-height: calc(100vh - 140px);
}

/* 左侧：日历卡片 */
.calendar-pane {
  flex: 1;
  max-width: 50%;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  overflow-y: auto;
  padding: 16px;
  box-sizing: border-box;
}

/* 右侧：任务面板卡片 */
.task-pane {
  flex: 1;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: #fff;
  border-radius: 16px;
  padding: 32px;
  overflow-y: auto;
  box-sizing: border-box;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .split-pane {
    padding: 90px 20px 20px 20px;
    flex-direction: column;
  }

  .calendar-pane {
    max-width: 100%;
  }
}
</style>
