<template>
  <transition name="slide-up">
    <div v-if="isVisible" class="modal-container">
      <div class="modal-header">
        <button @click="close" class="header-btn close-btn">×</button>
        <span class="header-title">New Event</span>
        <button @click="save" class="header-btn save-btn">✓</button>
      </div>

      <div class="modal-body-wrapper">
        <div class="modal-body">

          <div class="input-group">
            <input type="text" v-model="newEvent.title" placeholder="Title" class="title-input">
          </div>

          <div class="input-group details-group">
            <textarea v-model="newEvent.details" placeholder="Add details..." rows="3" class="details-input"></textarea>
          </div>

          <div class="option-section">
            <div class="option-row">
              <span>Starts</span>
              <div class="date-time-display">
                <input type="datetime-local" v-model="newEvent.start" class="datetime-input">
                <button @click="confirmStartTime" class="confirm-btn" title="确认开始时间">✓</button>
              </div>
            </div>

            <div class="option-row">
              <span>Ends</span>
              <div class="date-time-display">
                <input type="datetime-local" v-model="newEvent.end" class="datetime-input">
                <button @click="confirmEndTime" class="confirm-btn" title="确认结束时间">✓</button>
              </div>
            </div>

            <div class="option-row repeat-row">
              <span>Repeat</span>
              <div class="day-selector">
                <button v-for="day in ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']"
                        :key="day"
                        :class="{ 'selected-day': newEvent.repeat.includes(day) }"
                        @click="toggleRepeatDay(day)">
                  {{ day }}
                </button>
              </div>
            </div>

            <div class="option-row category-row">
              <span>Category</span>
              <div class="color-palette">
                <span v-for="cat in categories"
                      :key="cat.key"
                      :style="{ backgroundColor: cat.color }"
                      :class="{ 'selected-category': newEvent.category === cat.key }"
                      @click="newEvent.category = cat.key"
                      class="category-dot">
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </transition>
  <div v-if="isVisible" class="modal-overlay" @click.self="close"></div>
</template>

<script setup>
import { defineProps, defineEmits, ref, watch } from 'vue';

const props = defineProps({
  isVisible: Boolean,
  // 接收一个初始数据或 ID，这里简化为接收一个触发事件的标志
});
const emit = defineEmits(['close', 'save']);

// 预设事件种类和颜色
const categories = ref([
  {key: 'Work', color: '#ff3b30'}, // 红
  {key: 'Study', color: '#34c759'}, // 绿
  {key: 'Fun', color: '#ff9500'}, // 橙
  {key: 'Personal', color: '#5856d6'}, // 紫
]);

// 初始化一个空的新日程对象
const initialEvent = {
  title: '',
  details: '',
  start: '', // 初始设置为空字符串，让 datetime-local 处理
  end: '',
  repeat: [], // 存储 ['一', '三'] 这样的值
  category: categories.value[0].key, // 默认第一个类别
};

const newEvent = ref({...initialEvent});

// 监听 isVisible 变化，用于重置表单
watch(() => props.isVisible, (val) => {
  if (val) {
    // 确保每次打开时都使用默认值或继承选中的日期（如果需要更复杂的集成）
    newEvent.value = {
      ...initialEvent,
      // 可以在这里继承 picked 日期（如果 HomeView 将 picked.value 传进来）
      // 暂时简化为清空
    };
  }
}, {immediate: true});


function toggleRepeatDay(day) {
  const index = newEvent.value.repeat.indexOf(day);
  if (index > -1) {
    newEvent.value.repeat.splice(index, 1); // 移除
  } else {
    newEvent.value.repeat.push(day); // 添加
  }
}

// 确认开始时间
function confirmStartTime() {
  if (!newEvent.value.start) {
    alert('请先选择开始时间');
    return;
  }
  // 这里可以添加确认后的逻辑，比如验证时间格式等
  console.log('开始时间已确认:', newEvent.value.start);
}

// 确认结束时间
function confirmEndTime() {
  if (!newEvent.value.end) {
    alert('请先选择结束时间');
    return;
  }
  // 这里可以添加确认后的逻辑，比如验证时间格式等
  console.log('结束时间已确认:', newEvent.value.end);
}

function close() {
  emit('close');
}

function save() {
  if (!newEvent.value.title) {
    alert('Title is required!'); // 简单的验证
    return;
  }
  emit('save', {...newEvent.value});
}
</script>

<style scoped>
/* --- Modal Container & Overlay (与上次类似，略去重复代码) --- */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 2000;
}

.modal-container {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  background-color: #1c1c1e;
  color: white;
  border-radius: 15px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
  display: flex;
  flex-direction: column;
  z-index: 2001;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  border-bottom: 1px solid #38383a;
}

.header-title {
  color: white;
  font-size: 1.1rem;
  font-weight: 600;
}

.header-btn {
  background: none;
  border: none;
  color: #0a84ff;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 5px 10px;
  font-weight: 600;
}

.close-btn {
  font-size: 1.8rem;
}

.save-btn {
  font-size: 1.1rem;
}

.modal-body-wrapper {
  flex-grow: 1;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
}

.modal-body {
  padding: 15px 15px 40px;
  color: white;
}

/* --- Input Fields --- */
.input-group {
  background-color: #2c2c2e;
  border-radius: 10px;
  margin-bottom: 10px;
  padding: 0 15px;
}

.title-input {
  width: 100%;
  padding: 12px 0;
  border: none;
  background: none;
  color: white;
  font-size: 1.2rem;
  font-weight: 600;
  outline: none;
  border-bottom: 1px solid #38383a;
}

.details-group {
  margin-bottom: 20px;
}

.details-input {
  width: 100%;
  padding: 10px 0;
  border: none;
  background: none;
  color: white;
  font-size: 1rem;
  outline: none;
  resize: vertical;
}


/* --- Options Section --- */
.option-section {
  background-color: #2c2c2e;
  border-radius: 10px;
  margin-bottom: 20px;
}

.option-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 15px;
  border-bottom: 1px solid #38383a;
  font-size: 1rem;
}

.option-row:last-child {
  border-bottom: none;
}

/* --- Date/Time Input (简化处理) --- */
.date-time-display {
  display: flex;
  align-items: center;
  gap: 8px;
}

.datetime-input {
  background: #3a3a3c;
  color: #0a84ff;
  border: 1px solid #555;
  border-radius: 5px;
  padding: 5px;
  font-size: 0.9rem;
  outline: none;
}

/* 确认按钮样式 */
.confirm-btn {
  background: #0a84ff;
  color: white;
  border: none;
  border-radius: 50%;
  width: 24px;
  height: 24px;
  font-size: 0.8rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s;
}

.confirm-btn:hover {
  background: #0070e0;
}

/* --- Repeat (Days) --- */
.repeat-row {
  align-items: flex-start;
}

.day-selector {
  display: flex;
  gap: 8px;
}

.day-selector button {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background: #3a3a3c;
  color: #8e8e93;
  border: none;
  cursor: pointer;
  font-size: 0.7rem;
  font-weight: 600;
  transition: all 0.2s;
}

.day-selector button.selected-day {
  background-color: #0a84ff;
  color: white;
}


/* --- Category (Color) --- */
.category-row {
  font-weight: 500;
}

.color-palette {
  display: flex;
  gap: 10px;
}

.category-dot {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  cursor: pointer;
  border: 2px solid transparent;
  transition: border-color 0.2s;
}

.category-dot.selected-category {
  border-color: white; /* 选中时加白边 */
}

/* --- Animation --- */
.slide-up-enter-active, .slide-up-leave-active {
  transition: transform 0.4s ease, opacity 0.4s ease;
}

.slide-up-enter-from, .slide-up-leave-to {
  opacity: 0;
  transform: translate(-50%, 100vh);
}

.slide-up-enter-to, .slide-up-leave-from {
  opacity: 1;
  transform: translate(-50%, -50%);
}
</style>
