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

    <!-- 连接测试面板 -->
    <div class="connection-test-panel">
      <h4>🔗 后端连接测试</h4>
      <button @click="testBackendConnection" :disabled="testing">
        {{ testing ? '测试中...' : '测试连接' }}
      </button>
      <div v-if="testResult" class="test-result" :class="testResult.status">
        <span v-if="testResult.status === 'success'">✅ {{ testResult.message }}</span>
        <span v-else>❌ {{ testResult.message }}</span>
        <pre v-if="testResult.data">{{ JSON.stringify(testResult.data, null, 2) }}</pre>
      </div>
    </div>

    <!-- 个人信息模态框 -->
    <UserProfileModal
      :isVisible="showProfileModal"
      :user="currentUser"
      @close="showProfileModal = false"
      @update-user="handleUserUpdate"
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
const testing = ref(false)
const testResult = ref(null)

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

// 显示个人信息模态框
function showProfileModalFunc() {
  showProfileModal.value = true
}

// 暴露方法给父组件
defineExpose({
  showProfileModal: showProfileModalFunc
})

// 后端连接测试
const testBackendConnection = async () => {
  testing.value = true
  testResult.value = null

  try {
    // 测试健康检查端点
    const response = await fetch('http://localhost:8080/health')

    if (!response.ok) {
      throw new Error(`HTTP错误! 状态: ${response.status}`)
    }

    const data = await response.json()
    testResult.value = {
      status: 'success',
      message: `连接成功! 后端状态: ${data.status}`,
      data: data
    }

    console.log('✅ 后端连接测试成功:', data)
  } catch (error) {
    testResult.value = {
      status: 'error',
      message: `连接失败: ${error.message}`,
      data: null
    }
    console.error('❌ 后端连接测试失败:', error)
  } finally {
    testing.value = false
  }
}

onMounted(() => {
  checkAuth()
})
</script>

<style scoped>
.person-page {
  min-height: 100vh;
  background: linear-gradient(to bottom, #0e59b8, #16b1f4);
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
  padding: 90px 120px 40px 120px;
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

.connection-test-panel {
  position: fixed;
  bottom: 20px;
  right: 20px;
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid #dcdfe6;
  border-radius: 8px;
  padding: 15px;
  max-width: 300px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  z-index: 1000;
  font-size: 14px;
}

.connection-test-panel h4 {
  margin: 0 0 10px 0;
  color: #303133;
}

.connection-test-panel button {
  background: #409eff;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.connection-test-panel button:disabled {
  background: #c0c4cc;
  cursor: not-allowed;
}

.connection-test-panel button:hover:not(:disabled) {
  background: #66b1ff;
}

.test-result {
  margin-top: 10px;
  padding: 8px;
  border-radius: 4px;
  font-size: 12px;
}

.test-result.success {
  background: #f0f9ff;
  color: #67c23a;
  border: 1px solid #b3e19d;
}

.test-result.error {
  background: #fef0f0;
  color: #f56c6c;
  border: 1px solid #fbc4c4;
}

.test-result pre {
  background: rgba(0, 0, 0, 0.05);
  padding: 5px;
  border-radius: 3px;
  margin-top: 5px;
  font-size: 10px;
  overflow-x: auto;
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

  .connection-test-panel {
    position: relative;
    bottom: auto;
    right: auto;
    max-width: none;
    margin: 20px;
  }
}
</style>
