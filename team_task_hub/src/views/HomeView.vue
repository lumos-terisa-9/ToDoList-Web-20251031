<template>
  <div class="home-container">
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
    <NewTaskModal
      :isVisible="showModal"
      @close="showModal = false"
      @save="handleSaveTask"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import MonthCalendar from '@/components/MonthCalendar.vue'
import ToDoList from '@/components/ToDoList.vue'
import NewTaskModal from '@/components/NewTaskModal.vue'

const picked = ref(new Date())
function onSelect(d) {
  picked.value = d
  console.log('选中日期：', d)
}

// 控制模态框的显示状态
const showModal = ref(false)

// 监听 ToDoList 的事件来打开模态框
function handleNewTaskRequest() {
  showModal.value = true
}

// 处理模态框保存的逻辑
function handleSaveTask(taskData) {
  console.log('准备保存任务到日期:', picked.value, taskData)
  // 实际应用中，您需要在这里编写逻辑，将 taskData 添加到对应 picked 日期的 tasks 数组中 (ToDoList 内部的任务数组)
  showModal.value = false
}
const testing = ref(false)
const testResult = ref(null)

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
</script>

<style scoped>
.home-container {
  height: calc(100vh - 60px);
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  /* 【移除】这里的背景设置，让 App.vue 的蓝色渐变透过来 */
  /* background: linear-gradient(to bottom, #012855, #0e59b8); */
}

/* 分栏容器 */
.split-pane {
  display: flex;
  flex: 1;
  align-items: stretch;
  /* 【新增】卡片外边距，让卡片看起来是浮动的 */
  padding: 20px 80px; /* 核心修改点：两侧间距增大！！ */
  gap: 20px; /* 卡片之间的间距 */
}

/* 左侧：日历卡片 */
.calendar-pane {
  flex: 1;
  max-width: 50%;

  /* === 【核心修改】整体毛玻璃效果 === */
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);

  /* 【修改】统一圆角，使其成为完整的浮动卡片 */
  border-radius: 16px;
  overflow-y: auto;
  padding: 16px;
  box-sizing: border-box;
}

/* 右侧：任务面板卡片 */
.task-pane {
  flex: 1;
  /* 【修改】统一使用与左侧相同的毛玻璃效果 */
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);

  color: #fff;
  /* 【修改】统一圆角 */
  border-radius: 16px;
  padding: 32px;
  overflow-y: auto;
  box-sizing: border-box;
}

/* 【删除】分隔线 */
.split-pane::before {
  display: none;
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
</style>
