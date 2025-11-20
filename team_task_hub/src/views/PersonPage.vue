<template>
  <div class="person-page">
    <div class="split-pane">
      <MonthCalendar
        v-model="picked"
        @select="onSelect"
        class="calendar-pane"
      />

      <div class="task-pane">
        <h2>任务面板</h2>
        <p>已选日期：{{ picked.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' }) }}</p>

        <!-- 个人待办 -->
        <ToDoList
          :id="picked.toDateString()"
          :title="'个人待办'"
          @request-modal="handleNewTaskRequest"
          class="todo-list-item"
        />
        <!-- 组织待办 -->
        <ToDoList
          :id="'org-' + picked.toDateString()"
          :title="'组织待办'"
          :show-input="false"
          class="todo-list-item"
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

    <NewTaskModal
      :isVisible="showModal"
      @close="showModal = false"
      @save="handleSaveTask"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import MonthCalendar from '@/components/MonthCalendar.vue'
import ToDoList from '@/components/ToDoList.vue'
import NewTaskModal from '@/components/NewTaskModal.vue'
import UserProfileModal from '@/components/UserProfileModal.vue'

const router = useRouter()

const picked = ref(new Date())
const showModal = ref(false)
const showProfileModal = ref(false)
const currentUser = ref(null)

// 检查登录状态
function checkAuth() {
  const userData = localStorage.getItem('currentUser')
  if (!userData) {
    // 未登录，跳转到首页
    router.push('/')
    return
  }
  currentUser.value = JSON.parse(userData)
}

function onSelect(d) {
  picked.value = d
  console.log('选中日期：', d)
}

function handleNewTaskRequest() {
  showModal.value = true
}

function handleSaveTask(taskData) {
  console.log('准备保存任务到日期:', picked.value, taskData)
  showModal.value = false
}

function handleUserUpdate(updatedUser) {
  currentUser.value = updatedUser
  // 更新本地存储
  localStorage.setItem('currentUser', JSON.stringify(updatedUser))
}

function handleLogout() {
  // 清除用户数据
  currentUser.value = null
  // 跳转到首页
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
