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
                <input v-if="newEvent.repeat_type === 'none'" type="datetime-local" v-model="newEvent.start_time" class="datetime-input" required @change="validateRepeatTime">
                <input v-else type="time" v-model="newEvent.start_time_time" class="time-input" required>
                <button @click="confirmStartTime" class="confirm-btn" title="确认开始时间">✓</button>
              </div>
            </div>

            <div class="option-row">
              <span>End Time</span>
              <div class="date-time-display">
                <input v-if="newEvent.repeat_type === 'none'" type="datetime-local" v-model="newEvent.end_time" class="datetime-input" required @change="validateRepeatTime">
                <input v-else type="time" v-model="newEvent.end_time_time" class="time-input" required>
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
                <select v-model="newEvent.repeat_type" class="repeat-type-select" @change="onRepeatTypeChange">
                  <option value="none">No Repeat</option>
                  <option value="daily">Daily</option>
                  <option value="weekly">Weekly</option>
                  <option value="monthly">Monthly</option>
                </select>

                <!-- 每周重复 - 星期选择 -->
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

                <!-- 修改每月重复部分 -->
                <div v-else-if="newEvent.repeat_type === 'monthly'" class="repeat-extra-options">
                  <div class="repeat-settings-row">
                    <div class="repeat-interval-group">
                      <span>Interval:</span>
                      <input type="number" v-model="newEvent.repeat_interval" min="0" max="12" class="interval-input">
                      <span>months</span>
                    </div>
                    <div class="repeat-end-date">
                      <span>Ends:</span>
                      <input type="date" v-model="newEvent.repeat_end_date" class="date-input">
                    </div>
                  </div>
                  <div class="month-day-selector">
                    <div class="month-days-grid">
                      <button v-for="day in 31" :key="day"
                              :class="{ 'selected-month-day': newEvent.monthDays.includes(day) }"
                              @click="toggleMonthDay(day)">
                        {{ day }}
                      </button>
                    </div>
                  </div>
                </div>

                <!-- 每日重复 -->
                <div v-else-if="newEvent.repeat_type === 'daily'" class="repeat-extra-options">
                  <div class="repeat-settings-row">
                    <div class="repeat-interval-group">
                      <span>Interval:</span>
                      <input type="number" v-model="newEvent.repeat_interval" min="0" max="365" class="interval-input">
                      <span>days</span>
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
        <button @click="cancelTodo" class="footer-btn cancel-btn">取消当日代办</button>
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
// function getIntervalUnit(repeatType) {
//   const units = {
//     daily: 'days',
//     weekly: 'weeks',
//     monthly: 'months'
//   }
//   return units[repeatType] || ''
// }

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

// 初始化一个空的新日程对象
const initialEvent = {
  title: '',
  description: '',
  start_time: formatDateForInput(props.date),
  end_time: formatDateForInput(new Date(props.date.getTime() + 60 * 60 * 1000)),
  start_time_time: '09:00',
  end_time_time: '10:00',
  urgency: 'medium',
  category: 'personal',
  repeat_type: 'none',
  repeat_interval: 1,
  repeat_end_date: formatDateForDateInput(new Date(props.date.getTime() + 30 * 24 * 60 * 60 * 1000)),
  repeatDays: [], // 每周重复的星期
  monthDays: [], // 每月重复的日期
  child_dates: [] // 子待办日期数组
};

const newEvent = ref({...initialEvent});

// 监听 isVisible 变化，用于重置表单
// 监听 isVisible 变化，用于重置表单
watch(() => props.isVisible, (val) => {
  if (val) {
    if (props.isEditMode && props.editTodoData) {
      // 编辑模式：填充传入的数据
      const editData = props.editTodoData;

      // 处理时间格式
      const startDate = new Date(editData.start_time);
      const endDate = new Date(editData.end_time);

      // 处理 monthly 数据
      let monthDays = [];
      if (editData.monthDays && Array.isArray(editData.monthDays)) {
        monthDays = editData.monthDays;
      } else if (editData.repeat_type === 'monthly') {
        // 如果没有 monthDays 数据，使用开始日期的日期作为默认
        const startDateForMonth = new Date(editData.start_time);
        monthDays = [startDateForMonth.getDate()];
      }

      // 处理 repeatDays 数据
      let repeatDays = [];
      if (editData.repeatDays && Array.isArray(editData.repeatDays)) {
        repeatDays = editData.repeatDays;
      } else if (editData.repeat_type === 'weekly') {
        // 如果没有 repeatDays 数据，使用开始日期的星期作为默认
        const startDateForWeek = new Date(editData.start_time);
        const dayNames = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
        repeatDays = [dayNames[startDateForWeek.getDay()]];
      }

      newEvent.value = {
        ...editData,
        // 确保日期格式正确
        start_time: formatDateForInput(startDate),
        end_time: formatDateForInput(endDate),
        start_time_time: `${String(startDate.getHours()).padStart(2, '0')}:${String(startDate.getMinutes()).padStart(2, '0')}`,
        end_time_time: `${String(endDate.getHours()).padStart(2, '0')}:${String(endDate.getMinutes()).padStart(2, '0')}`,
        repeat_end_date: editData.repeat_end_date ? formatDateForDateInput(new Date(editData.repeat_end_date)) : '',
        repeatDays: repeatDays,
        monthDays: monthDays,
        child_dates: editData.child_dates || []
      };

      console.log('编辑模式数据加载完成:', {
        repeat_type: newEvent.value.repeat_type,
        repeatDays: newEvent.value.repeatDays,
        monthDays: newEvent.value.monthDays
      });
    } else {
      // 创建模式：重置表单
      newEvent.value = {
        ...initialEvent,
        start_time: formatDateForInput(props.date),
        end_time: formatDateForInput(new Date(props.date.getTime() + 60 * 60 * 1000)),
        repeat_end_date: formatDateForDateInput(new Date(props.date.getTime() + 30 * 24 * 60 * 60 * 1000)),
        // 设置默认的重复日期
        monthDays: [props.date.getDate()],
        repeatDays: ['Mon']
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

// 重复类型改变时的处理
function onRepeatTypeChange() {
  if (newEvent.value.repeat_type !== 'none') {
    // 保存时间部分
    const startDate = new Date(newEvent.value.start_time);
    const endDate = new Date(newEvent.value.end_time);

    newEvent.value.start_time_time = `${String(startDate.getHours()).padStart(2, '0')}:${String(startDate.getMinutes()).padStart(2, '0')}`;
    newEvent.value.end_time_time = `${String(endDate.getHours()).padStart(2, '0')}:${String(endDate.getMinutes()).padStart(2, '0')}`;
  }

  // 初始化重复设置
  if (newEvent.value.repeat_type === 'weekly') {
    newEvent.value.repeatDays = ['Mon'];
  } else if (newEvent.value.repeat_type === 'monthly') {
    newEvent.value.monthDays = [props.date.getDate()];
  }
}

function toggleRepeatDay(day) {
  const index = newEvent.value.repeatDays.indexOf(day);
  if (index > -1) {
    newEvent.value.repeatDays.splice(index, 1);
  } else {
    newEvent.value.repeatDays.push(day);
  }
}

// 切换每月日期
function toggleMonthDay(day) {
  const index = newEvent.value.monthDays.indexOf(day);
  if (index > -1) {
    newEvent.value.monthDays.splice(index, 1);
  } else {
    newEvent.value.monthDays.push(day);
  }
  newEvent.value.monthDays.sort((a, b) => a - b);
}

function confirmStartTime() {
  if (!newEvent.value.start_time && !newEvent.value.start_time_time) {
    alert('请先选择开始时间');
    return;
  }
  console.log('开始时间已确认:', newEvent.value.repeat_type === 'none' ? newEvent.value.start_time : newEvent.value.start_time_time);
}

function confirmEndTime() {
  if (!newEvent.value.end_time && !newEvent.value.end_time_time) {
    alert('请先选择结束时间');
    return;
  }
  console.log('结束时间已确认:', newEvent.value.repeat_type === 'none' ? newEvent.value.end_time : newEvent.value.end_time_time);
}

function close() {
  emit('close');
}

const formatToFullISO = (datetimeLocalStr) => {
  if (!datetimeLocalStr) return null;

  // 如果已经是完整格式（包含时区），直接返回
  if (datetimeLocalStr.includes('+') || datetimeLocalStr.includes('Z')) {
    return datetimeLocalStr;
  }

  // 如果是 datetime-local 格式 (YYYY-MM-DDTHH:MM)，添加秒和时区
  const date = new Date(datetimeLocalStr);
  return date.toISOString(); // 这会转换为 UTC 时间，如 2025-11-28T17:00:00.000Z
};

function generateChildDates() {
  if (newEvent.value.repeat_type === 'none') {
    return [];
  }

  const childDates = [];
  const today = new Date();
  const repeatEndDate = newEvent.value.repeat_end_date ? new Date(newEvent.value.repeat_end_date) : new Date(today.getFullYear() + 1, today.getMonth(), today.getDate());

  // 设置日期范围限制
  const weeklyEndDate = new Date(today);
  weeklyEndDate.setDate(today.getDate() + 7); // 未来7天

  const monthlyEndDate = new Date(today);
  monthlyEndDate.setMonth(today.getMonth() + 1); // 未来1个月

  // let currentDate = new Date(today);

  switch (newEvent.value.repeat_type) {
    case 'daily':
      // 只生成今天的一个日期
      { const todayStr = today.toISOString().split('T')[0];
      childDates.push(todayStr);
      break; }

    case 'weekly':
      { const dayMap = { 'Mon': 1, 'Tue': 2, 'Wed': 3, 'Thu': 4, 'Fri': 5, 'Sat': 6, 'Sun': 0 };
      const selectedDays = newEvent.value.repeatDays.map(day => dayMap[day]);

      // 生成未来7天内选中的日期
      for (let i = 0; i <= 7; i++) {
        const checkDate = new Date(today);
        checkDate.setDate(today.getDate() + i);

        // 取重复结束日期和7天限制的较小值
        const endDate = repeatEndDate < weeklyEndDate ? repeatEndDate : weeklyEndDate;
        if (checkDate > endDate) break;

        if (selectedDays.includes(checkDate.getDay())) {
          const dateStr = checkDate.toISOString().split('T')[0];
          childDates.push(dateStr);
        }
      }
      break; }

    case 'monthly':
      // 确保 monthDays 有值
      { if (!newEvent.value.monthDays || newEvent.value.monthDays.length === 0) {
        console.warn('monthly 重复未选择日期');
        return [];
      }

      // 生成未来1个月内选中的日期
      // 取重复结束日期和1个月限制的较小值
      const monthlyEnd = repeatEndDate < monthlyEndDate ? repeatEndDate : monthlyEndDate;

      for (let dayOffset = 0; dayOffset <= 31; dayOffset++) { // 最多31天
        const checkDate = new Date(today);
        checkDate.setDate(today.getDate() + dayOffset);

        if (checkDate > monthlyEnd) break;

        // 检查是否在选中的日期中
        const dayOfMonth = checkDate.getDate();
        if (newEvent.value.monthDays.includes(dayOfMonth)) {
          const dateStr = checkDate.toISOString().split('T')[0];
          childDates.push(dateStr);
        }
      }
      break; }
  }

  console.log('生成的子日期:', childDates);
  return childDates;
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
    // 生成子日期
    const childDates = generateChildDates();
    const firstChildDate = childDates.length > 0 ? childDates[0] : new Date().toISOString().split('T')[0];
    const taskData = {
      title: newEvent.value.title,
      description: newEvent.value.description || '',
      start_time: formatToFullISO(`${firstChildDate}T${newEvent.value.start_time_time || '00:00'}:00`),
      end_time: formatToFullISO(`${firstChildDate}T${newEvent.value.end_time_time || '00:00'}:00`),
      urgency: newEvent.value.urgency,
      category: newEvent.value.category,
      repeat_type: newEvent.value.repeat_type,
      repeat_interval: newEvent.value.repeat_interval,
      repeat_end_date: newEvent.value.repeat_end_date ? formatToFullISO(newEvent.value.repeat_end_date + 'T00:00:00') : null,
      child_dates: childDates
    }

    // 详细打印所有数据
    console.log('=== 创建代办 - 完整数据打印 ===');
    console.log('基础信息:');
    console.log('  - title:', taskData.title);
    console.log('  - description:', taskData.description);
    console.log('  - start_time:', taskData.start_time);
    console.log('  - end_time:', taskData.end_time);
    console.log('  - urgency:', taskData.urgency);
    console.log('  - category:', taskData.category);
    console.log('重复设置:');
    console.log('  - repeat_type:', taskData.repeat_type);
    console.log('  - repeat_interval:', taskData.repeat_interval);
    console.log('  - repeat_end_date:', taskData.repeat_end_date);
    console.log('子日期数据 (child_dates):');
    // 修改子日期数据的打印逻辑
    console.log('子日期数据 (child_dates):');
    if (childDates.length > 0) {
      childDates.forEach((date, index) => {
        console.log(`  [${index}] date: ${date}`);
      });
      console.log(`  共生成 ${childDates.length} 个子日期`);
    } else {
      console.log('  无子日期（非重复任务或重复设置无效）');
    }
    console.log('表单原始数据 (newEvent):', JSON.parse(JSON.stringify(newEvent.value)));
    console.log('=== 数据打印结束 ===');

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

// 修改待办：先调用 cancel-with-children 再调用创建接口
async function updateTodo() {
  if (!newEvent.value.title) {
    alert('Title is required!');
    return;
  }

  const token = getToken()
  if (!token) {
    alert('未找到认证令牌')
    return
  }

  try {
    console.log('=== 修改待办 - 开始 ===');
    console.log('=== 所有 props 数据 ===');
    console.log('props.isVisible:', props.isVisible);
    console.log('props.date:', props.date);
    console.log('props.isEditMode:', props.isEditMode);
    console.log('props.editTodoData:', props.editTodoData);

    // 详细打印 editTodoData 的所有字段
    if (props.editTodoData) {
      console.log('=== props.editTodoData 详细字段 ===');
      for (const key in props.editTodoData) {
        console.log(`  ${key}:`, props.editTodoData[key]);
      }
    }

    // 先保存原始数据到临时变量
    const originalCreatedAt = props.editTodoData.created_at;
    const originalTitle = props.editTodoData.title;
    const originalDescription = props.editTodoData.description || "";

    if (!originalCreatedAt) {
      alert('无法获取任务的创建时间');
      return;
    }

    // 第一步：调用 cancel-with-children 接口结束原任务
    const cancelRequestBody = {
      "createdAt": originalCreatedAt,
      "description": originalDescription,
      "title": originalTitle
    };

    console.log('=== 调用 cancel-with-children 接口 ===');
    console.log('请求URL:', `${API_BASE}/todos/cancel-with-children`);
    console.log('请求方法: POST');
    console.log('请求头:', {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    });
    console.log('请求体参数:', JSON.stringify(cancelRequestBody, null, 2));
    console.log('cancelRequestBody 详细字段:');
    console.log('  - createdAt:', cancelRequestBody.createdAt);
    console.log('  - description:', cancelRequestBody.description);
    console.log('  - title:', cancelRequestBody.title);

    const cancelResponse = await fetch(`${API_BASE}/todos/cancel-with-children`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(cancelRequestBody)
    });

    const cancelResult = await cancelResponse.json();
    console.log('cancel-with-children 响应状态:', cancelResponse.status);
    console.log('cancel-with-children 响应结果:', cancelResult);

    if (!cancelResponse.ok || !cancelResult.success) {
      alert(cancelResult.message || '结束原任务失败');
      return;
    }

    // 第二步：调用创建代办接口创建新任务（编辑后的任务）
    const childDates = generateChildDates();
    const firstChildDate = childDates.length > 0 ? childDates[0] : new Date().toISOString().split('T')[0];
    const createRequestBody = {
      "category": newEvent.value.category || "personal",
      "child_dates": generateChildDates(),
      "description": newEvent.value.description || "",
      "end_time": formatToFullISO(`${firstChildDate}T${newEvent.value.end_time_time || '00:00'}:00`),
      "repeat_interval": newEvent.value.repeat_interval || 1,
      "repeat_type": newEvent.value.repeat_type || "none",
      "start_time": formatToFullISO(`${firstChildDate}T${newEvent.value.start_time_time || '00:00'}:00`),
      "title": newEvent.value.title,
      "urgency": newEvent.value.urgency || "medium"
    };

    console.log('=== 调用 createTodo 接口 ===');
    console.log('请求URL:', `${API_BASE}/todos/createTodo`);
    console.log('请求方法: POST');
    console.log('请求头:', {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    });
    console.log('请求体参数:', JSON.stringify(createRequestBody, null, 2));
    console.log('createRequestBody 详细字段:');
    console.log('  - category:', createRequestBody.category);
    console.log('  - child_dates:', createRequestBody.child_dates);
    console.log('  - description:', createRequestBody.description);
    console.log('  - end_time:', createRequestBody.end_time);
    console.log('  - repeat_end_date:', createRequestBody.repeat_end_date);
    console.log('  - repeat_interval:', createRequestBody.repeat_interval);
    console.log('  - repeat_type:', createRequestBody.repeat_type);
    console.log('  - start_time:', createRequestBody.start_time);
    console.log('  - title:', createRequestBody.title);
    console.log('  - urgency:', createRequestBody.urgency);

    const createResponse = await fetch(`${API_BASE}/todos/createTodo`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(createRequestBody)
    });

    const createResult = await createResponse.json();
    console.log('createTodo 响应状态:', createResponse.status);
    console.log('createTodo 响应结果:', createResult);

    if (createResponse.ok && createResult.success) {
      console.log('=== 修改待办 - 成功 ===');
      emit('update', { ...newEvent.value, id: createResult.data?.id })
      emit('close')
    } else {
      alert(createResult.message || '创建新任务失败');
    }

  } catch (error) {
    console.error('编辑任务失败:', error);
    alert('编辑任务失败，请检查网络连接');
  }
}

// 结束待办 - 使用 cancel-with-children 接口
async function completeTodo() {
  const token = getToken()
  if (!token) {
    alert('未找到认证令牌')
    return
  }

  try {
    console.log('=== 结束待办 - 开始 ===');
    console.log('=== 所有 props 数据 ===');
    console.log('props.isVisible:', props.isVisible);
    console.log('props.date:', props.date);
    console.log('props.isEditMode:', props.isEditMode);
    console.log('props.editTodoData:', props.editTodoData);

    // 详细打印 editTodoData 的所有字段
    if (props.editTodoData) {
      console.log('=== props.editTodoData 详细字段 ===');
      for (const key in props.editTodoData) {
        console.log(`  ${key}:`, props.editTodoData[key]);
      }
    }

    // 使用正确的 created_at 字段
    const createdAt = props.editTodoData.created_at;

    if (!createdAt) {
      alert('无法获取任务的创建时间');
      return;
    }

    // 调用 cancel-with-children 接口结束代办
    const requestBody = {
      "createdAt": createdAt, // 使用 created_at 字段
      "description": props.editTodoData.description || "",
      "title": props.editTodoData.title
    };

    console.log('=== 调用 cancel-with-children 接口 ===');
    console.log('请求URL:', `${API_BASE}/todos/cancel-with-children`);
    console.log('请求方法: POST');
    console.log('请求头:', {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    });
    console.log('请求体参数:', JSON.stringify(requestBody, null, 2));
    console.log('requestBody 详细字段:');
    console.log('  - createdAt:', requestBody.createdAt);
    console.log('  - description:', requestBody.description);
    console.log('  - title:', requestBody.title);

    const response = await fetch(`${API_BASE}/todos/cancel-with-children`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(requestBody)
    });

    const result = await response.json();
    console.log('cancel-with-children 响应状态:', response.status);
    console.log('cancel-with-children 响应结果:', result);

    if (response.ok && result.success) {
      console.log('=== 结束待办 - 成功 ===');
      emit('complete')
      emit('close')
    } else {
      alert(result.message || '结束代办失败');
    }

  } catch (error) {
    console.error('调用结束代办接口失败:', error);
    alert('结束代办失败，请检查网络连接');
  }
}

// 取消待办 - 调用 /api/todos/cancel 接口
async function cancelTodo() {
  const token = getToken()
  if (!token) {
    alert('未找到认证令牌')
    return
  }

  try {
    console.log('=== 取消待办 - 开始 ===');
    console.log('取消的任务数据 (props.editTodoData):', props.editTodoData);

    // 调用 cancel 接口取消待办
    const requestBody = {
      "description": props.editTodoData.description || "",
      "end_time": props.editTodoData.end_time,
      "start_time": props.editTodoData.start_time,
      "title": props.editTodoData.title
    };

    console.log('=== 调用 cancel 接口 ===');
    console.log('请求URL:', `${API_BASE}/todos/cancel`);
    console.log('请求方法: POST');
    console.log('请求头:', {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    });
    console.log('请求体参数:', JSON.stringify(requestBody, null, 2));
    console.log('requestBody 详细字段:');
    console.log('  - description:', requestBody.description);
    console.log('  - end_time:', requestBody.end_time);
    console.log('  - start_time:', requestBody.start_time);
    console.log('  - title:', requestBody.title);

    const response = await fetch(`${API_BASE}/todos/cancel`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(requestBody)
    });

    const result = await response.json();
    console.log('cancel 响应状态:', response.status);
    console.log('cancel 响应结果:', result);

    if (response.ok && result.success) {
      console.log('=== 取消待办 - 成功 ===');
      emit('complete') // 使用 complete 事件刷新列表
      emit('close')
    } else {
      alert(result.message || '取消待办失败');
    }

  } catch (error) {
    console.error('调用取消待办接口失败:', error);
    alert('取消待办失败，请检查网络连接');
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
  /* 隐藏滚动条但保持滚动功能 */
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE and Edge */
}

.modal-body-wrapper::-webkit-scrollbar {
  display: none; /* Chrome, Safari and Opera */
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

.time-input {
  background: #3a3a3c;
  color: #0a84ff;
  border: 1px solid #555;
  border-radius: 5px;
  padding: 5px;
  font-size: 0.9rem;
  outline: none;
  width: 120px;
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

/* 修改每月日期选择器对齐 */
.month-day-selector {
  width: 100%;
  display: flex;
  justify-content: flex-end; /* 右对齐 */
  margin: 10px 0; /* 上下间隙 */
  padding: 0 15px; /* 左右间隙 */
}

.month-days-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 5px;
  max-width: 280px;
}

.month-days-grid button {
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

.month-days-grid button.selected-month-day {
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
