<template>
  <div class="frosted-glass-panel">
    <h3 class="panel-title">{{ title }}</h3>

    <div class="task-input-group">
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

      <ul v-if="tasks.length > 0" class="task-list">

        <li v-for="(task, index) in tasks" :key="task.id || index" class="task-item">

        <span
          class="task-color-dot"
          :style="{ backgroundColor: task.color || '#007aff' }">
        </span>

          <div class="task-info">
            <span class="task-text">{{ task.title }}</span>

            <span v-if="task.start" class="task-time">
            {{ formatTaskTime(task.start) }}
          </span>
            <span v-else-if="task.details" class="task-time">
            {{ task.details.substring(0, 50) }} </span>
          </div>

          <button @click="deleteTask(index)" class="delete-btn">×</button>
        </li>

      </ul>
      <p v-else class="no-tasks">🎉 今天没有待办事项！</p>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, defineProps, defineEmits } from 'vue';

const emit = defineEmits(['request-modal', 'add-task-object']);

const props = defineProps({
  // 接收外部传入的日期或ID，用于区分不同的任务列表
  id: {
    type: [String, Number, Date],
    required: true
  },
  title: {
    type: String,
    default: '我的待办清单'
  }
});

const newTaskText = ref('');
const tasks = ref([]); // 任务列表

// 1. 从 localStorage 加载任务
function loadTasks(id) {
  // 使用不同的键存储任务，以区分不同日期的列表
  const storedTasks = localStorage.getItem(`tasks-${id}`);
  if (storedTasks) {
    let loadedData = JSON.parse(storedTasks);

    // 兼容性检查：如果是旧的字符串数组，进行迁移
    if (loadedData.length > 0 && typeof loadedData[0] === 'string') {
      tasks.value = loadedData.map(title => ({
        id: Date.now() + Math.random(), // 赋予临时 ID
        title: title,
        details: '',
        start: '',
        end: '',
        repeat: [],
        category: 'General',
        color: '#007aff' // 默认颜色
      }));
      // 立即保存新格式，防止下次再次触发迁移
      saveTasks(id, tasks.value);
    } else {
      tasks.value = loadedData; // 加载新的对象数组
    }
  } else {
    tasks.value = [];
  }
}

// 2. 保存任务到 localStorage
function saveTasks(id, currentTasks) {
  localStorage.setItem(`tasks-${id}`, JSON.stringify(currentTasks));
}

// 3. 监听 id 变化，自动加载对应日期的任务
watch(() => props.id, (newId) => {
  loadTasks(newId);
}, { immediate: true });

// 4. 监听 tasks 变化，自动保存
watch(tasks, (newTasks) => {
  saveTasks(props.id, newTasks);
}, { deep: true });

// 添加待办事项 -> 请求打开模态框
function requestNewTaskModal() {
  const text = newTaskText.value.trim();
  // 触发事件，并将输入框内容作为参数传递给父组件
  emit('request-modal', text);
  // 清空输入框，但不在本地添加任务，任务的添加由模态框的保存操作驱动
  newTaskText.value = '';
}

// 删除待办事项
function deleteTask(index) {
  tasks.value.splice(index, 1);
}

function addNewTaskObject(taskObject) {
  // 确保任务对象有一个 ID (如果模态框没有生成)
  if (!taskObject.id) {
    taskObject.id = Date.now();
  }
  tasks.value.unshift(taskObject);
}

// 格式化时间函数，将一个 Date/Time 字符串或对象转换成易读的时间格式
const formatTaskTime = (dateTimeString) => {
  // 假设 dateTimeString 是一个有效的日期时间字符串 (例如 ISO 格式)
  const date = new Date(dateTimeString);

  // 检查日期对象是否有效
  if (isNaN(date.getTime())) {
    return 'Invalid Date';
  }

  // 格式化为 'HH:mm'，例如 '09:30'
  const timeOptions = {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false // 使用 24 小时制
  };

  return date.toLocaleTimeString('zh-CN', timeOptions);
}

defineExpose({
  addNewTaskObject,
  tasks // 也可以暴露 tasks 数组，但暴露方法更安全
});
</script>

<style scoped>
/* ========================================= */
/* === 1. 玻璃效果 (毛玻璃 / Frosted Glass) === */
/* ========================================= */
.frosted-glass-panel {
  /* 背景是**半透明白色**，这是实现模糊效果的基础 */
  background-color: rgba(255, 255, 255, 0.2);

  /* === 核心 CSS 属性：背景模糊 === */
  /* backdrop-filter 会模糊这个元素**后面**的内容 */
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px); /* 兼容 Safari */

  /* 增加 iOS 般的质感：一个细小的白色边框 */
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 4px 30px rgba(0, 0, 0, 0.1);

  border-radius: 12px;
  padding: 16px;
  color: #fff; /* 文本颜色为白色，以适应深色背景 */
  margin-top: 16px;
  z-index: 10; /* 确保在最上层 */
  pointer-events: auto; /* 重新启用鼠标事件 */
}

.panel-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin-top: 0;
  margin-bottom: 20px;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2); /* 轻微阴影，让文字更突出 */
}

/* ========================================= */
/* === 2. 输入框样式 === */
/* ========================================= */
.task-input-group {
  display: flex;
  margin-bottom: 20px;
}

.task-input {
  flex-grow: 1;
  padding: 10px 15px;
  border: none;
  /* 输入框背景颜色略微透明，但比面板不透明一点 */
  background-color: rgba(255, 255, 255, 0.85);
  border-radius: 8px 0 0 8px;
  font-size: 1rem;
  color: #333;
  outline: none;
  transition: background-color 0.3s;
}

.task-input::placeholder {
  color: #888;
}

.add-btn {
  background-color: #007aff; /* iOS 蓝色 */
  color: white;
  border: none;
  padding: 0 15px;
  border-radius: 0 8px 8px 0;
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
  font-size: 1.5rem;
  font-weight: bold;
  line-height: 1;
}

/* 【新增样式】任务信息布局 */
.task-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: rgba(255, 255, 255, 0.1);
  padding: 12px;
  margin-bottom: 10px;
  border-radius: 8px;
  transition: background-color 0.3s;
}

.task-info {
  flex-grow: 1;
  display: flex;
  flex-direction: column; /* 标题和时间垂直排列 */
  text-align: left;
}

.task-text {
  font-weight: 500;
  line-height: 1.4;
}

.task-time {
  font-size: 0.8rem;
  opacity: 0.7;
  margin-top: 2px;
}

/* 【新增样式】任务类别颜色点 */
.task-color-dot {
  display: inline-block;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  margin-right: 10px;
  flex-shrink: 0; /* 阻止它被压缩 */
}

/* ========================================= */
/* === 3. 列表样式 === */
/* ========================================= */
.task-list-container {
  max-height: 400px; /* 限制高度并允许滚动 */
  overflow-y: auto;
  padding-right: 8px; /* 为滚动条留出空间 */
}

/* 自定义滚动条 */
.task-list-container::-webkit-scrollbar {
  width: 6px;
}
.task-list-container::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 3px;
}
.task-list-container::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.5); /* 浅色透明滑块 */
  border-radius: 3px;
}

.task-list {
  list-style-type: none;
  padding: 0;
  margin: 0;
}

.task-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  /* 待办项背景也使用玻璃材质，但更透明 */
  background-color: rgba(255, 255, 255, 0.1);
  padding: 12px;
  margin-bottom: 10px;
  border-radius: 8px;
  transition: background-color 0.3s;
}

.task-item:hover {
  background-color: rgba(255, 255, 255, 0.2);
}

.task-text {
  flex-grow: 1;
  text-align: left;
  line-height: 1.4;
}

.delete-btn {
  background: none;
  border: none;
  color: #ff3b30; /* iOS 红色 */
  font-weight: bold;
  font-size: 1.5rem;
  line-height: 1;
  margin-left: 10px;
  cursor: pointer;
  padding: 0 5px;
  transition: opacity 0.3s;
}

.delete-btn:hover {
  opacity: 0.8;
}

.no-tasks {
  text-align: center;
  font-style: italic;
  opacity: 0.7;
  padding: 20px 0;
}
</style>
