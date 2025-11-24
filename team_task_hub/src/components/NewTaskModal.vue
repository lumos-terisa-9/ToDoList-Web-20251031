<template>
  <transition name="slide-up">
    <div v-if="isVisible" class="modal-container">
      <div class="modal-header">
        <button @click="close" class="header-btn close-btn">×</button>
        <span class="header-title">{{ isEditMode ? 'Edit Event' : 'New Event' }}</span>
        <button v-if="!isEditMode" @click="save" class="header-btn save-btn">✓</button>
      </div>

      <div class="modal-body-wrapper">
        <div class="modal-body">

          <div class="input-group">
            <input type="text" v-model="newEvent.title" placeholder="Title" class="title-input" required>
          </div>

          <div class="input-group details-group">
            <textarea v-model="newEvent.description" placeholder="Add description..." rows="3" class="details-input"></textarea>
          </div>

          <div class="option-section">
            <div class="option-row">
              <span>Start Time</span>
              <div class="date-time-display">
                <input type="datetime-local" v-model="newEvent.start_time" class="datetime-input" required @change="validateRepeatTime">
                <button @click="confirmStartTime" class="confirm-btn" title="确认开始时间">✓</button>
              </div>
            </div>

            <div class="option-row">
              <span>End Time</span>
              <div class="date-time-display">
                <input type="datetime-local" v-model="newEvent.end_time" class="datetime-input" required @change="validateRepeatTime">
                <button @click="confirmEndTime" class="confirm-btn" title="确认结束时间">✓</button>
              </div>
            </div>

            <div class="option-row">
              <span>Urgency</span>
              <select v-model="newEvent.urgency" class="urgency-select">
                <option value="low">Low</option>
                <option value="medium">Medium</option>
                <option value="high">High</option>
              </select>
            </div>

            <div class="option-row repeat-row">
              <span>Repeat</span>
              <div class="repeat-options">
                <select v-model="newEvent.repeat_type" class="repeat-type-select" @change="validateRepeatTime">
                  <option value="none">No Repeat</option>
                  <option value="daily">Daily</option>
                  <option value="weekly">Weekly</option>
                  <option value="monthly">Monthly</option>
                </select>
                <div v-if="newEvent.repeat_type === 'weekly'" class="repeat-extra-options">
                  <div class="repeat-settings-row">
                    <div class="repeat-interval-group">
                      <span>Interval:</span>
                      <input type="number" v-model="newEvent.repeat_interval" min="0" max="52" class="interval-input">
                      <span>weeks</span>
                    </div>
                    <div class="repeat-end-date">
                      <span>Ends:</span>
                      <input type="date" v-model="newEvent.repeat_end_date" class="date-input">
                    </div>
                  </div>
                  <div class="day-selector">
                    <button v-for="day in ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']"
                            :key="day"
                            :class="{ 'selected-day': newEvent.repeatDays.includes(day) }"
                            @click="toggleRepeatDay(day)">
                      {{ day }}
                    </button>
                  </div>
                </div>
                <div v-else-if="newEvent.repeat_type !== 'none'" class="repeat-extra-options">
                  <div class="repeat-settings-row">
                    <div class="repeat-interval-group">
                      <span>Interval:</span>
                      <input type="number" v-model="newEvent.repeat_interval" min="0" max="12" class="interval-input">
                      <span>{{ getIntervalUnit(newEvent.repeat_type) }}</span>
                    </div>
                    <div class="repeat-end-date">
                      <span>Ends:</span>
                      <input type="date" v-model="newEvent.repeat_end_date" class="date-input">
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div class="option-row category-row">
              <span>Category</span>
              <div class="color-palette">
                <span v-for="cat in categories"
                      :key="cat.value"
                      :style="{ backgroundColor: cat.color }"
                      :class="{ 'selected-category': newEvent.category === cat.value }"
                      @click="newEvent.category = cat.value"
                      class="category-dot"
                      :title="cat.name">
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 编辑模式下的按钮 -->
      <div v-if="isEditMode" class="modal-footer">
        <button @click="updateTodo" class="footer-btn update-btn">修改待办</button>
        <button @click="completeTodo" class="footer-btn complete-btn">结束待办</button>
      </div>
    </div>
  </transition>
  <div v-if="isVisible" class="modal-overlay" @click.self="close"></div>
</template>

<script setup>
import { defineProps, defineEmits, ref, watch } from 'vue';

const props = defineProps({
  isVisible: Boolean,
  date: {
    type: Date,
    required: true
  },
  // 添加编辑模式相关属性
  isEditMode: {
    type: Boolean,
    default: false
  },
  editTodoData: {
    type: Object,
    default: null
  }
});
const emit = defineEmits(['close', 'save', 'update', 'complete']);

const API_BASE = 'http://localhost:8080/api'

// 根据后端定义的枚举值
const categories = ref([
  { value: 'work', name: 'Work', color: '#ff3b30' },
  { value: 'study', name: 'Study', color: '#34c759' },
  { value: 'fun', name: 'Fun', color: '#ff9500' },
  { value: 'personal', name: 'Personal', color: '#5856d6' },
]);

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

// 格式化日期为ISO字符串
function formatDateForInput(date) {
  if (!date) return ''
  const d = new Date(date)
  d.setMinutes(d.getMinutes() - d.getTimezoneOffset())
  return d.toISOString().slice(0, 16)
}

// 格式化日期为YYYY-MM-DD
function formatDateForDateInput(date) {
  if (!date) return ''
  const d = new Date(date)
  return d.toISOString().split('T')[0]
}

// 获取重复间隔单位
function getIntervalUnit(repeatType) {
  const units = {
    daily: 'days',
    weekly: 'weeks',
    monthly: 'months'
  }
  return units[repeatType] || ''
}

// 验证重复任务的时间：如果选择了重复，开始和结束时间必须在同一天
function validateRepeatTime() {
  if (newEvent.value.repeat_type !== 'none') {
    const start = new Date(newEvent.value.start_time);
    const end = new Date(newEvent.value.end_time);

    if (start.toDateString() !== end.toDateString()) {
      // 不在同一天，将结束时间设置为开始时间
      newEvent.value.end_time = newEvent.value.start_time;
      alert('重复任务的开始时间和结束时间必须在同一天内，已自动调整结束时间。');
    }
  }
}

// 初始化一个空的新日程对象 - 严格按照后端结构
const initialEvent = {
  title: '',
  description: '',
  start_time: formatDateForInput(props.date),
  end_time: formatDateForInput(new Date(props.date.getTime() + 60 * 60 * 1000)),
  urgency: 'medium',
  category: 'personal',
  repeat_type: 'none',
  repeat_interval: 1,
  repeat_end_date: formatDateForDateInput(new Date(props.date.getTime() + 30 * 24 * 60 * 60 * 1000)), // 默认30天后结束
  repeatDays: [], // 前端使用的重复天数
  child_dates: [] // 子待办日期数组
};

const newEvent = ref({...initialEvent});

// 监听 isVisible 变化，用于重置表单
watch(() => props.isVisible, (val) => {
  if (val) {
    if (props.isEditMode && props.editTodoData) {
      // 编辑模式：填充传入的数据
      newEvent.value = {
        ...props.editTodoData,
        // 确保日期格式正确
        start_time: formatDateForInput(new Date(props.editTodoData.start_time)),
        end_time: formatDateForInput(new Date(props.editTodoData.end_time)),
        repeat_end_date: props.editTodoData.repeat_end_date ? formatDateForDateInput(new Date(props.editTodoData.repeat_end_date)) : ''
      };
    } else {
      // 创建模式：重置表单
      newEvent.value = {
        ...initialEvent,
        start_time: formatDateForInput(props.date),
        end_time: formatDateForInput(new Date(props.date.getTime() + 60 * 60 * 1000)),
        repeat_end_date: formatDateForDateInput(new Date(props.date.getTime() + 30 * 24 * 60 * 60 * 1000))
      };
    }
  }
}, {immediate: true});

// 监听日期变化（创建模式）
watch(() => props.date, (newDate) => {
  if (props.isVisible && !props.isEditMode) {
    newEvent.value.start_time = formatDateForInput(newDate)
    newEvent.value.end_time = formatDateForInput(new Date(newDate.getTime() + 60 * 60 * 1000))
  }
});

function toggleRepeatDay(day) {
  const index = newEvent.value.repeatDays.indexOf(day);
  if (index > -1) {
    newEvent.value.repeatDays.splice(index, 1);
  } else {
    newEvent.value.repeatDays.push(day);
  }
}

function confirmStartTime() {
  if (!newEvent.value.start_time) {
    alert('请先选择开始时间');
    return;
  }
  console.log('开始时间已确认:', newEvent.value.start_time);
}

function confirmEndTime() {
  if (!newEvent.value.end_time) {
    alert('请先选择结束时间');
    return;
  }
  console.log('结束时间已确认:', newEvent.value.end_time);
}

function close() {
  emit('close');
}

async function save() {
  if (!newEvent.value.title) {
    alert('Title is required!');
    return;
  }

  if (!newEvent.value.start_time || !newEvent.value.end_time) {
    alert('Start time and end time are required!');
    return;
  }

  const token = getToken()
  if (!token) {
    alert('未找到认证令牌')
    return
  }

  // 获取当前用户ID
  const currentUserId = localStorage.getItem('currentUserId')
  if (!currentUserId) {
    alert('无法获取用户信息，请重新登录')
    return
  }

  try {
    // 准备发送到后端的数据 - 严格按照后端结构
    const taskData = {
      title: newEvent.value.title,
      description: newEvent.value.description || '',
      start_time: new Date(newEvent.value.start_time).toISOString(),
      end_time: new Date(newEvent.value.end_time).toISOString(),
      urgency: newEvent.value.urgency,
      category: newEvent.value.category,
      repeat_type: newEvent.value.repeat_type,
      repeat_interval: newEvent.value.repeat_interval,
      repeat_end_date: newEvent.value.repeat_end_date ? new Date(newEvent.value.repeat_end_date).toISOString() : null,
      child_dates: newEvent.value.child_dates || [],
      creator_id: parseInt(currentUserId) // 添加创建者ID
    }

    console.log('发送任务数据:', taskData)

    const response = await fetch(`${API_BASE}/todos/createTodo`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(taskData)
    })

    if (response.ok) {
      const result = await response.json()
      if (result.success) {
        console.log('任务创建成功:', result)
        emit('save', { ...newEvent.value, id: result.data?.id })
        emit('close')
      } else {
        alert(result.message || '创建失败')
      }
    } else {
      const errorResult = await response.json()
      console.error('创建任务失败:', response.status, errorResult)
      alert(`创建任务失败: ${errorResult.message || '未知错误'}`)
    }
  } catch (error) {
    console.error('调用创建任务接口失败:', error)
    alert('创建任务失败，请检查网络连接')
  }
}

// 修改待办：先删除再创建
async function updateTodo() {
  // 先删除原待办
  const deleteSuccess = await deleteTodo(props.editTodoData.id);
  if (deleteSuccess) {
    // 再创建新待办
    await save();
  } else {
    alert('删除原待办失败，无法修改');
  }
}

// 删除待办 - 修复：使用正确的接口
async function deleteTodo() {
  const token = getToken()
  if (!token) {
    alert('未找到认证令牌')
    return false
  }

  try {
    // 使用取消待办接口，传递必要的参数
    const response = await fetch(`${API_BASE}/todos/cancel`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        title: props.editTodoData.title,
        content: props.editTodoData.description || ''
      })
    })

    if (response.ok) {
      const result = await response.json()
      if (result.success) {
        console.log('取消待办成功')
        return true
      } else {
        alert(result.message || '取消待办失败')
        return false
      }
    } else {
      const errorResult = await response.json()
      console.error('取消待办失败:', response.status, errorResult)
      alert(`取消待办失败: ${errorResult.message || '未知错误'}`)
      return false
    }
  } catch (error) {
    console.error('调用取消待办接口失败:', error)
    alert('取消待办失败，请检查网络连接')
    return false
  }
}

// 完成待办 - 修复：传递正确的参数
async function completeTodo() {
  const token = getToken()
  if (!token) {
    alert('未找到认证令牌')
    return
  }

  try {
    // 根据错误信息，需要传递这些必填字段
    const completeData = {
      title: props.editTodoData.title,
      description: props.editTodoData.description || 'No description',
      start_time: props.editTodoData.start_time,
      end_time: props.editTodoData.end_time
    }

    console.log('完成待办参数:', completeData)

    const response = await fetch(`${API_BASE}/todos/complete`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(completeData)
    })

    if (response.ok) {
      const result = await response.json()
      if (result.success) {
        console.log('完成任务成功')
        emit('complete')
        emit('close')
      } else {
        alert(result.message || '完成任务失败')
      }
    } else {
      const errorResult = await response.json()
      console.error('完成任务失败:', response.status, errorResult)
      alert(`完成任务失败: ${errorResult.message || '未知错误'}`)
    }
  } catch (error) {
    console.error('调用完成任务接口失败:', error)
    alert('完成任务失败，请检查网络连接')
  }
}
</script>

<style scoped>
/* 原有样式保持不变，新增底部按钮样式 */
.modal-footer {
  display: flex;
  justify-content: space-around;
  padding: 15px;
  border-top: 1px solid #38383a;
}

.footer-btn {
  padding: 12px 30px;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.3s;
  min-width: 140px;
  max-height: 40px;
}

.update-btn {
  background-color: #007aff;
  color: white;
}

.update-btn:hover {
  background-color: #0056b3;
}

.complete-btn {
  background-color: #ff3b30;
  color: white;
}

.complete-btn:hover {
  background-color: #d70015;
}

/* 原有样式保持不变，新增底部按钮样式 */
.modal-footer {
  display: flex;
  justify-content: space-around;
  padding: 15px;
  border-top: 1px solid #38383a;
}

.footer-btn {
  padding: 8px 20px;
  border: none;
  border-radius: 8px;
  font-size: 0.8rem;
  cursor: pointer;
  transition: background-color 0.3s;
}

.update-btn {
  background-color: #007aff;
  color: white;
}

.update-btn:hover {
  background-color: #0056b3;
}

.complete-btn {
  background-color: #34c759;
  color: white;
}

.complete-btn:hover {
  background-color: #2aa44f;
}

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
  padding: 10px 14px;
  border-bottom: 1px solid #38383a;
  position: relative;
}

.header-title {
  position: absolute; /* 改为绝对定位 */
  left: 50%; /* 居中 */
  transform: translateX(-50%); /* 居中 */
  color: white;
  font-size: 1.0rem;
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

.date-input {
  background: #3a3a3c;
  color: #0a84ff;
  border: 1px solid #555;
  border-radius: 4px;
  padding: 4px 8px;
  font-size: 0.9rem;
  outline: none;
}

/* 调整下拉框样式，使背景更深 */
.urgency-select, .repeat-type-select {
  background: #2c2c2e;
  color: white;
  border: 1px solid #444;
  border-radius: 6px;
  padding: 6px 12px;
  font-size: 0.9rem;
  outline: none;
  min-width: 120px;
}

.urgency-select:focus, .repeat-type-select:focus {
  border-color: #0a84ff;
}

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

.repeat-row {
  align-items: center;
}

.repeat-options {
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: flex-end;
  text-align: right;
}

.repeat-extra-options {
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: flex-end;
}

.repeat-settings-row {
  display: flex;
  align-items: center;
  gap: 15px;
  justify-content: flex-end;
}

.repeat-interval-group {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.9rem;
}

.repeat-end-date {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.9rem;
}

.interval-input {
  background: #3a3a3c;
  color: white;
  border: 1px solid #555;
  border-radius: 4px;
  padding: 4px 8px;
  width: 50px;
  font-size: 0.9rem;
  outline: none;
}

.day-selector {
  display: flex;
  gap: 5px;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.day-selector button {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #2c2c2e;
  color: #8e8e93;
  border: 1px solid #444;
  cursor: pointer;
  font-size: 0.7rem;
  font-weight: 600;
  transition: all 0.2s;
}

.day-selector button.selected-day {
  background-color: #0a84ff;
  color: white;
  border-color: #0a84ff;
}

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
  border-color: white;
}

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
