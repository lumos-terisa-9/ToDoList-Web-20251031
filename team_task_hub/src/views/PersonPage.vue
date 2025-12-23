<template>
  <div class="person-page">
    <div class="notification-wrapper">
      <NotificationBell ref="notificationBell" />
    </div>
    <div class="split-pane">
      <MonthCalendar
        v-model="picked"
        @select="onSelect"
        @date-click="handleDateClick"
        @load-tasks="handleLoadTasks"
        @open-activity-modal="handleOpenActivityModal"
        class="calendar-pane"
        ref="monthCalendar"
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
          @refresh-todos="refreshTodos"
          @edit-task="handleEditTask"
          @open-activity-modal="handleOpenActivityModal"
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
          @refresh-todos="refreshTodos"
          @edit-task="handleEditTask"
          @open-activity-modal="handleOpenActivityModal"
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

    <ActivityDetailModal
      v-model:visible="showActivityModal"
      :activity-data="selectedActivityData"
      @close="closeActivityModal"
      @review-submitted="handleReviewSubmitted"
      @review-failed="handleReviewFailed"
      class="global-activity-modal"
    />
  </div>
</template>

<script setup>
import ActivityDetailModal from "@/components/ActivityDetailModal.vue";

const API_BASE = 'http://localhost:8080/api'
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import MonthCalendar from '@/components/MonthCalendar.vue'
import ToDoList from '@/components/ToDoList.vue'
import NewTaskModal from '@/components/NewTaskModal.vue'
import UserProfileModal from '@/components/UserProfileModal.vue'
import NotificationBell from '@/components/NotificationBell.vue'

const router = useRouter()

const notificationBell = ref(null)
const picked = ref(new Date())
const showModal = ref(false)
const showEditModal = ref(false) // 添加编辑模态框状态
const showProfileModal = ref(false)
const currentUser = ref(null)
const personalTodoList = ref(null)
const orgTodoList = ref(null)
const selectedTask = ref(null) // 添加选中的任务
const monthCalendar = ref(null) // 添加MonthCalendar引用

// 在 PersonPage.vue 的 script 部分添加数据状态
const currentTasks = ref([]);
const currentUpcomingTasks = ref([]);
const currentDisplayMode = ref('today'); // 'today', 'future', 或 'completed'

// 在PersonPage.vue的script中添加
const showActivityModal = ref(false);
const selectedActivityData = ref(null);

// 添加处理函数
function handleOpenActivityModal(task) {
  console.log('收到打开活动详情请求:', task);

  if (task.organization || task.organization_id || task.creator_organ_id > 0) {
    selectedActivityData.value = task;
    showActivityModal.value = true;
  } else {
    console.warn('这不是一个组织任务:', task);
  }
}

function closeActivityModal() {
  showActivityModal.value = false;
  selectedActivityData.value = null;
}

function handleReviewSubmitted(reviewData) {
  console.log('评价提交成功:', reviewData);
  // 可以添加提示
}

function handleReviewFailed(errorMessage) {
  console.error('评价提交失败:', errorMessage);
  // 可以添加错误提示
}

// 处理从日历组件加载的任务数据
function handleLoadTasks(taskData) {
  console.log('接收到任务数据:', taskData);

  currentDisplayMode.value = taskData.type;
  currentTasks.value = taskData.tasks || [];
  currentUpcomingTasks.value = taskData.upcomingTasks || [];

  // 更新ToDoList数据
  updateTodoListData();
}

// 修改 handleDateClick 函数
function handleDateClick(dateInfo) {
  console.log('日期点击:', dateInfo);
  // 这里不再需要加载数据，因为MonthCalendar已经处理了
}

// 刷新待办事项
async function refreshTodos() {
  console.log('刷新待办事项，当前日期:', picked.value);

  // 调用MonthCalendar的重新加载方法
  if (monthCalendar.value) {
    await monthCalendar.value.reloadDate(picked.value);
  }
}

// 新增函数：更新ToDoList数据
function updateTodoListData() {
  console.log('更新ToDoList数据:', currentTasks.value);
  console.log('更新即将开始任务数据:', currentUpcomingTasks.value);

  // 通知子组件更新数据
  if (personalTodoList.value) {
    personalTodoList.value.updateTasks(currentTasks.value);
    personalTodoList.value.updateUpcomingTasks(currentUpcomingTasks.value);
  }
  if (orgTodoList.value) {
    orgTodoList.value.updateTasks(currentTasks.value);
    orgTodoList.value.updateUpcomingTasks(currentUpcomingTasks.value);
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
    return false
  }

  try {
    const response = await fetch('http://localhost:8080/api/auth/me', {
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

      // 初始化加载今日待办
      currentDisplayMode.value = 'today';

      // 先调用更新代办接口
      await checkAndUpdateTodos(token);

      // 自动触发今日日期的点击事件
      if (monthCalendar.value) {
        const today = new Date();
        await monthCalendar.value.reloadDate(today);
      }

      return true;
    } else {
      router.push('/')
      return false;
    }
  } catch (error) {
    console.error('验证用户失败:', error)
    router.push('/')
    return false;
  }
}

// 检查并更新代办
async function checkAndUpdateTodos(token) {
  try {
    // 检查是否需要更新
    const lastUpdate = localStorage.getItem('last_todo_update');
    const today = new Date().toDateString();

    console.log('=== 检查代办更新状态 ===');
    console.log('最后更新时间:', lastUpdate || '无记录');
    console.log('今天日期:', today);

    // 如果没有更新记录或者不是今天更新的，才调用接口
    if (!lastUpdate || lastUpdate !== today) {
      console.log('🔄 开始调用更新代办接口...');

      const response = await fetch(`${API_BASE}/todos/updateTodos`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });

      if (response.ok) {
        const result = await response.json();
        console.log('📊 更新代办接口响应:', result);

        if (result.success) {
          // 更新成功，记录更新时间
          localStorage.setItem('last_todo_update', today);
          console.log('✅ 代办更新成功，已记录更新时间:', today);
        } else {
          console.error('❌ 更新代办失败:', result.message);
        }
      } else {
        console.error('❌ 调用更新代办接口失败:', response.status);
      }
    } else {
      console.log('✅ 今日已更新过代办，无需更新');
    }
  } catch (error) {
    console.error('❌ 检查更新代办失败:', error);
  }
}

function onSelect(d) {
  picked.value = d
  console.log('选中日期：', d)
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
async function handleTaskUpdate() {
  showEditModal.value = false
  await refreshTodos() // 重新加载任务
}

// 处理任务完成
async function handleTaskComplete() {
  showEditModal.value = false
  await refreshTodos() // 重新加载任务
}

async function handleSaveTask(taskData) {
  console.log('准备保存任务到日期:', picked.value, taskData)
  showModal.value = false

  // 保存后刷新待办列表
  await nextTick()
  await refreshTodos()
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

onMounted(async () => {
  await checkAuth();
  // 自动触发今日日期的点击事件
  if (monthCalendar.value) {
    // 模拟点击今日日期
    const today = new Date();
    await monthCalendar.value.reloadDate(today);
  }
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

/* 添加通知组件样式 */
.person-page {
  position: relative; /* 添加相对定位 */
}

.notification-wrapper {
  position: fixed;
  top: 90px;
  right: 100px;
  z-index: 1000;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .notification-wrapper {
    right: 20px;
  }
}
</style>
