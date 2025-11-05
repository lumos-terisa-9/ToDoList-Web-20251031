<template>
  <div class="frosted-glass-panel">
    <h3 class="panel-title">{{ title }}</h3>

    <div class="task-input-group">
      <input
        type="text"
        v-model="newTaskText"
        @keyup.enter="addTask"
        placeholder="手动输入新的待办事项..."
        class="task-input"
      />
      <button @click="addTask" class="add-btn">
        <span class="plus-icon">+</span>
      </button>
    </div>

    <div class="task-list-container">
      <ul v-if="tasks.length > 0" class="task-list">
        <li v-for="(task, index) in tasks" :key="index" class="task-item">
          <span class="task-text">{{ task }}</span>
          <button @click="deleteTask(index)" class="delete-btn">×</button>
        </li>
      </ul>
      <p v-else class="no-tasks">🎉 今天没有待办事项！</p>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue';

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
  tasks.value = storedTasks ? JSON.parse(storedTasks) : [];
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


// 添加待办事项
function addTask() {
  const text = newTaskText.value.trim();
  if (text) {
    tasks.value.unshift(text); // 新任务放在最前面
    newTaskText.value = '';
  }
}

// 删除待办事项
function deleteTask(index) {
  tasks.value.splice(index, 1);
}
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
