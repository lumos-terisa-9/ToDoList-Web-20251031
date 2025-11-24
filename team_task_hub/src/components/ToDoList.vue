<template>
  <div class="frosted-glass-panel">
    <h3 class="panel-title">{{ title }}</h3>

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
        <li v-for="task in filteredTasks" :key="task.id" class="task-item">
          <div class="task-checkbox"
               :class="{ 'completed': task.completed }"
               @click.stop="toggleTaskComplete(task)">
            <span v-if="task.completed" class="checkmark">✓</span>
          </div>
          <span class="task-text" :class="{ 'completed': task.completed }"
                @click="openEditModal(task)">
            {{ task.title }}
          </span>
        </li>
      </ul>
      <p v-else class="no-tasks">🎉 今天没有待办事项！</p>
    </div>

    <!-- 即将开始/结束的代办 -->
    <div class="upcoming-tasks-section">
      <h4 class="upcoming-title">即将开始/结束的代办</h4>
      <div class="upcoming-list-container">
        <ul v-if="upcomingTasks.length > 0" class="upcoming-list">
          <li v-for="task in upcomingTasks" :key="task.id" class="upcoming-item">
            <div class="upcoming-checkbox"
                 :class="{ 'completed': task.completed }"
                 @click.stop="toggleUpcomingTaskComplete(task)">
              <span v-if="task.completed" class="upcoming-checkmark">✓</span>
            </div>
            <span class="upcoming-text" :class="{ 'completed': task.completed }"
                  @click="openEditModal(task)">
              {{ task.title }}
            </span>
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
      <button @click="loadTasksFromAPI" class="retry-btn">重试</button>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, defineProps, defineEmits, onMounted, computed } from 'vue';

const emit = defineEmits(['request-modal', 'add-task-object', 'refresh-todos', 'edit-task']);

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
const loading = ref(false);
const error = ref(null);

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

// 根据标题过滤任务：个人待办只显示personal类别，组织待办显示非personal类别
const filteredTasks = computed(() => {
  // 优先使用外部传入的数据，如果没有则使用apiTasks
  const tasksToFilter = externalTasks.value;

  console.log('当前显示模式:', props.displayMode, '任务数量:', tasksToFilter.length);
  console.log('任务数据:', tasksToFilter);

  if (props.title === '个人待办') {
    return tasksToFilter.filter(task => task.category === 'personal');
  } else if (props.title === '组织待办') {
    return tasksToFilter.filter(task => task.category !== 'personal');
  }
  return tasksToFilter;
});

// 从API加载任务
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

        // 暂时将即将开始/结束的代办设置为空数组，后续可以添加逻辑
        upcomingTasks.value = [];
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

// 切换任务完成状态
function toggleTaskComplete(task) {
  console.log('切换任务状态:', task.title, '当前状态:', task.completed);
  task.completed = !task.completed;
  console.log('新状态:', task.completed);
}

// 切换即将开始/结束任务的完成状态
function toggleUpcomingTaskComplete(task) {
  console.log('切换即将开始任务状态:', task.title, '当前状态:', task.completed);
  task.completed = !task.completed;
}

// 打开编辑模态框
function openEditModal(task) {
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

defineExpose({
  refreshFromAPI,
  updateTasks
})

onMounted(() => {
  loadTasksFromAPI()
})

// 监听日期变化重新加载
watch(() => props.date, () => {
  loadTasksFromAPI()
})
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
</style>
