<template>
  <div class="month-calendar">
    <div class="calendar-header">
      <button @click="prevMonth" class="nav-btn">‹</button>
      <span class="header-title">{{ year }} 年 {{ month + 1 }} 月</span>
      <button @click="nextMonth" class="nav-btn">›</button>
    </div>

    <!-- 添加分类信息显示 -->
    <div class="categories-display">
      <div v-for="cat in categories" :key="cat.key" class="category-item">
        <span class="category-dot" :style="{ backgroundColor: cat.color }"></span>
        <span class="category-name">{{ cat.key }}</span>
      </div>
    </div>

    <table class="calendar-table">
      <thead>
      <tr>
        <th v-for="d in days" :key="d">{{ d }}</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="(week, wi) in weeks" :key="wi">

        <td v-for="(date, di) in week"
            :key="di"
            :class="{
                empty: isNaN(date.getDate()),
                today: isToday(date),
                selected: isSelected(date)
              }"
            @click="selectDate(date)">

            <span class="day-number" v-if="!isNaN(date.getDate())">
              {{ date.getDate() }}
            </span>

        </td>
      </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup>
import { ref, computed, defineEmits, defineProps } from 'vue'

const API_BASE = 'http://localhost:8080/api'
const props = defineProps({
  modelValue: { type: Date, required: true }
})
const emit = defineEmits(['update:modelValue', 'select', 'date-click', 'load-tasks'])

const current = ref(new Date(props.modelValue))
const showActivityModal = ref(false)
const selectedActivityData = ref(null)

// 活动缓存
const activityCache = new Map()

// 将星期表头改为英文大写以匹配图片
const days = ['SUN', 'MON', 'TUE', 'WED', 'THU', 'FRI', 'SAT']

// 分类信息
const categories = ref([
  {key: 'Work', color: '#ff3b30'}, // 红
  {key: 'Study', color: '#34c759'}, // 绿
  {key: 'Fun', color: '#ff9500'}, // 橙
  {key: 'Personal', color: '#5856d6'}, // 紫
])

const year = computed(() => current.value.getFullYear())
const month = computed(() => current.value.getMonth())

const weeks = computed(() => {
  const first = new Date(year.value, month.value, 1)
  const last = new Date(year.value, month.value + 1, 0)
  const result = []
  let week = []
  // 补前空格
  for (let i = 0; i < first.getDay(); i++) week.push(new Date(NaN))
  for (let d = 1; d <= last.getDate(); d++) {
    week.push(new Date(year.value, month.value, d))
    if (week.length === 7) {
      result.push(week)
      week = []
    }
  }
  // 补后空格
  if (week.length) {
    while (week.length < 7) week.push(new Date(NaN))
    result.push(week)
  }
  return result
})

function openActivityModal(activityData) {
  console.log('MonthCalendar接收到活动数据:', activityData);

  // 触发事件给父组件，让父组件显示弹窗
  emit('open-activity-modal', activityData);
}

// 缓存活动数据
function cacheActivityData(activity) {
  if (activity && activity.id) {
    activityCache.set(activity.id, activity)
  }
}

// 缓存活动列表数据
function cacheActivities(activities) {
  if (Array.isArray(activities)) {
    activities.forEach(activity => {
      cacheActivityData(activity)
    })
  }
}

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

// 修复时区问题：将日期转换为本地日期字符串 (YYYY-MM-DD)
function formatDateToLocalString(date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}

// 清除所有缓存
function clearAllCache() {
  // 只清除当前日期的缓存，保留其他缓存键（如token等）
  localStorage.removeItem('current_date_tasks');
  localStorage.removeItem('current_date_info');
  console.log('已清除所有任务缓存');
}

// 缓存当前日期任务数据
function cacheCurrentDateTasks(date, tasks, upcomingTasks, type, activeOrganizations) {
  const cacheData = {
    date: date.toISOString(),
    tasks: tasks,
    upcomingTasks: upcomingTasks || [],
    type: type,
    activeOrganizations: activeOrganizations || [],
    timestamp: new Date().getTime()
  };
  localStorage.setItem('current_date_tasks', JSON.stringify(cacheData));
  localStorage.setItem('current_date_info', JSON.stringify({
    date: date.toISOString(),
    type: type,
    timestamp: new Date().getTime()
  }));
  console.log(`已缓存${type}任务数据:`, cacheData);
}

// 检查是否是今天之后的日期
function isAfterToday(date) {
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const compareDate = new Date(date);
  compareDate.setHours(0, 0, 0, 0);
  return compareDate > today;
}

// ==================== 个人代办接口 ====================

// 加载今日待办
async function loadTodayTodos() {
  const token = getToken();
  if (!token) return null;

  try {
    console.log('开始调用今日待办接口...');
    const response = await fetch(`${API_BASE}/todos/todayTodos`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    if (response.ok) {
      const result = await response.json();
      console.log('今日待办接口返回的全部内容:', result);

      if (result.success && result.todos) {
        // 缓存活动数据（如果是组织活动）
        result.todos.forEach(todo => {
          if (todo.creator_organ_id && todo.creator_organ_id > 0) {
            cacheActivityData(todo)
          }
        });

        return result.todos;
      }
    } else {
      console.error('获取今日待办失败:', response.status);
    }
  } catch (error) {
    console.error('调用今日待办接口失败:', error);
  }
  return null;
}

// 加载某一天开始的待办
async function loadOneDayTodos(date) {
  const token = getToken();
  if (!token) return null;

  try {
    // 修复时区问题：使用本地日期格式
    const dateStr = formatDateToLocalString(date);
    console.log('开始调用某一天开始的待办接口...', dateStr);
    const response = await fetch(`${API_BASE}/todos/Get-OneDayTodos?date=${dateStr}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    if (response.ok) {
      const result = await response.json();
      console.log('某一天开始的待办接口返回的全部内容:', result);
      if (result.success && result.todos) {
        // 缓存活动数据（如果是组织活动）
        result.todos.forEach(todo => {
          if (todo.creator_organ_id && todo.creator_organ_id > 0) {
            cacheActivityData(todo)
          }
        });

        return result.todos;
      }
    } else {
      console.error('获取某一天开始的待办失败:', response.status);
    }
  } catch (error) {
    console.error('调用某一天开始的待办接口失败:', error);
  }
  return null;
}

// 加载已完成待办
async function loadCompletedTodos(date) {
  const token = getToken();
  if (!token) return null;

  try {
    // 修复时区问题：使用本地日期格式
    const dateStr = formatDateToLocalString(date);
    console.log('开始调用已完成待办接口，日期:', dateStr);

    const response = await fetch(`${API_BASE}/todos/completed_todo?date=${dateStr}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    if (response.ok) {
      const result = await response.json();
      console.log('已完成待办接口返回的全部内容:', result);
      if (result.success && result.todos) {
        // 缓存活动数据（如果是组织活动）
        result.todos.forEach(todo => {
          if (todo.creator_organ_id && todo.creator_organ_id > 0) {
            cacheActivityData(todo)
          }
        });

        return result.todos;
      }
    } else {
      console.error('获取已完成待办失败:', response.status);
    }
  } catch (error) {
    console.error('调用已完成待办接口失败:', error);
  }
  return null;
}

// 加载即将开始的待办
async function loadComingStartTodos() {
  const token = getToken();
  if (!token) return null;

  try {
    console.log('开始调用即将开始的待办接口...');
    const response = await fetch(`${API_BASE}/todos/coming-startTodos`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    if (response.ok) {
      const result = await response.json();
      console.log('即将开始的待办接口返回的全部内容:', result);
      if (result.success && result.todos) {
        // 缓存活动数据（如果是组织活动）
        result.todos.forEach(todo => {
          if (todo.creator_organ_id && todo.creator_organ_id > 0) {
            cacheActivityData(todo)
          }
        });

        return result.todos;
      }
    } else {
      console.error('获取即将开始的待办失败:', response.status);
    }
  } catch (error) {
    console.error('调用即将开始的待办接口失败:', error);
  }
  return null;
}

// 加载即将结束的待办
async function loadComingEndTodos() {
  const token = getToken();
  if (!token) return null;

  try {
    console.log('开始调用即将结束的待办接口...');
    const response = await fetch(`${API_BASE}/todos/coming-endTodos`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    if (response.ok) {
      const result = await response.json();
      console.log('即将结束的待办接口返回的全部内容:', result);
      if (result.success && result.todos) {
        // 缓存活动数据（如果是组织活动）
        result.todos.forEach(todo => {
          if (todo.creator_organ_id && todo.creator_organ_id > 0) {
            cacheActivityData(todo)
          }
        });

        return result.todos;
      }
    } else {
      console.error('获取即将结束的待办失败:', response.status);
    }
  } catch (error) {
    console.error('调用即将结束的待办接口失败:', error);
  }
  return null;
}

// ==================== 组织代办接口 ====================

// 加载今日组织待办
async function loadTodayOrganizationTodos() {
  const token = getToken();
  if (!token) return null;

  try {
    console.log('开始调用今日组织待办接口...');
    const response = await fetch(`${API_BASE}/todos/activities/today`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    if (response.ok) {
      const result = await response.json();
      console.log('今日组织待办接口返回的全部内容:', result);
      if (result.success && result.activities) {
        // 转换数据结构，添加creator_organ_id字段
        const activities = result.activities.map(activity => ({
          ...activity,
          creator_organ_id: activity.organization_id,
          creator_user_id: 0,
          // 为兼容性保留原字段
          organization: activity.organization
        }));

        // 缓存活动数据
        cacheActivities(activities);

        return activities;
      }
    } else {
      console.error('获取今日组织待办失败:', response.status);
    }
  } catch (error) {
    console.error('调用今日组织待办接口失败:', error);
  }
  return null;
}

// 加载某一天开始的组织待办
async function loadOneDayOrganizationTodos(date) {
  const token = getToken();
  if (!token) return null;

  try {
    const dateStr = formatDateToLocalString(date);
    console.log('开始调用某一天开始的组织待办接口...', dateStr);
    const response = await fetch(`${API_BASE}/todos/activities/starting-on-date?date=${dateStr}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    if (response.ok) {
      const result = await response.json();
      console.log('某一天开始的组织待办接口返回的全部内容:', result);
      if (result.success && result.activities) {
        // 转换数据结构，添加creator_organ_id字段
        const activities = result.activities.map(activity => ({
          ...activity,
          creator_organ_id: activity.organization_id,
          creator_user_id: 0,
          organization: activity.organization
        }));

        // 缓存活动数据
        cacheActivities(activities);

        return activities;
      }
    } else {
      console.error('获取某一天开始的组织待办失败:', response.status);
    }
  } catch (error) {
    console.error('调用某一天开始的组织待办接口失败:', error);
  }
  return null;
}

// 加载某一天过期的组织待办
async function loadOneDayExpiredOrganizationTodos(date) {
  const token = getToken();
  if (!token) return null;

  try {
    const dateStr = formatDateToLocalString(date);
    console.log('开始调用某一天过期的组织待办接口...', dateStr);
    const response = await fetch(`${API_BASE}/todos/activities/expiring-on-date?date=${dateStr}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    if (response.ok) {
      const result = await response.json();
      console.log('某一天过期的组织待办接口返回的全部内容:', result);
      if (result.success && result.activities) {
        // 转换数据结构，添加creator_organ_id字段
        const activities = result.activities.map(activity => ({
          ...activity,
          creator_organ_id: activity.organization_id,
          creator_user_id: 0,
          organization: activity.organization
        }));

        // 缓存活动数据
        cacheActivities(activities);

        return activities;
      }
    } else {
      console.error('获取某一天过期的组织待办失败:', response.status);
    }
  } catch (error) {
    console.error('调用某一天过期的组织待办接口失败:', error);
  }
  return null;
}

// 加载某一天完成的组织待办
async function loadCompletedOrganizationTodos(date) {
  const token = getToken();
  if (!token) return null;

  try {
    const dateStr = formatDateToLocalString(date);
    console.log('开始调用某一天完成的组织待办接口...', dateStr);
    const response = await fetch(`${API_BASE}/todos/activities/completed-on-date?date=${dateStr}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    if (response.ok) {
      const result = await response.json();
      console.log('某一天完成的组织待办接口返回的全部内容:', result);
      if (result.success && result.activities) {
        // 转换数据结构，添加creator_organ_id字段
        const activities = result.activities.map(activity => ({
          ...activity,
          creator_organ_id: activity.organization_id,
          creator_user_id: 0,
          organization: activity.organization
        }));

        // 缓存活动数据
        cacheActivities(activities);

        return activities;
      }
    } else {
      console.error('获取某一天完成的组织待办失败:', response.status);
    }
  } catch (error) {
    console.error('调用某一天完成的组织待办接口失败:', error);
  }
  return null;
}

// 加载未来七天即将开始的组织待办
async function loadComingStartOrganizationTodos() {
  const token = getToken();
  if (!token) return null;

  try {
    console.log('开始调用未来七天即将开始的组织待办接口...');
    const response = await fetch(`${API_BASE}/todos/activities/upcoming-starting`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    if (response.ok) {
      const result = await response.json();
      console.log('未来七天即将开始的组织待办接口返回的全部内容:', result);
      if (result.success && result.activities) {
        // 转换数据结构，添加creator_organ_id字段
        const activities = result.activities.map(activity => ({
          ...activity,
          creator_organ_id: activity.organization_id,
          creator_user_id: 0,
          organization: activity.organization
        }));

        // 缓存活动数据
        cacheActivities(activities);

        return activities;
      }
    } else {
      console.error('获取未来七天即将开始的组织待办失败:', response.status);
    }
  } catch (error) {
    console.error('调用未来七天即将开始的组织待办接口失败:', error);
  }
  return null;
}

// 加载未来七天即将结束的组织待办
async function loadComingEndOrganizationTodos() {
  const token = getToken();
  if (!token) return null;

  try {
    console.log('开始调用未来七天即将结束的组织待办接口...');
    const response = await fetch(`${API_BASE}/todos/activities/upcoming-ending`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    if (response.ok) {
      const result = await response.json();
      console.log('未来七天即将结束的组织待办接口返回的全部内容:', result);
      if (result.success && result.activities) {
        // 转换数据结构，添加creator_organ_id字段
        const activities = result.activities.map(activity => ({
          ...activity,
          creator_organ_id: activity.organization_id,
          creator_user_id: 0,
          organization: activity.organization
        }));

        // 缓存活动数据
        cacheActivities(activities);

        return activities;
      }
    } else {
      console.error('获取未来七天即将结束的组织待办失败:', response.status);
    }
  } catch (error) {
    console.error('调用未来七天即将结束的组织待办接口失败:', error);
  }
  return null;
}

// 加载今日有活动的组织列表
async function loadActiveOrganizations() {
  const token = getToken();
  if (!token) return null;

  try {
    console.log('开始调用今日有活动的组织列表接口...');
    const response = await fetch(`${API_BASE}/todos/organizations/today`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    if (response.ok) {
      const result = await response.json();
      console.log('今日有活动的组织列表接口返回的全部内容:', result);
      if (result.success && result.organizations) {
        return result.organizations;
      }
    } else {
      console.error('获取今日有活动的组织列表失败:', response.status);
    }
  } catch (error) {
    console.error('调用今日有活动的组织列表接口失败:', error);
  }
  return null;
}

function prevMonth() {
  current.value = new Date(year.value, month.value - 1, 1);
}

function nextMonth() {
  current.value = new Date(year.value, month.value + 1, 1);
}

function isToday(date) {
  if (isNaN(date)) return false;
  const t = new Date();
  return date.getFullYear() === t.getFullYear() &&
    date.getMonth() === t.getMonth() &&
    date.getDate() === t.getDate();
}

function isSelected(date) {
  if (isNaN(date)) return false;
  const d = props.modelValue;
  return date.getFullYear() === d.getFullYear() &&
    date.getMonth() === d.getMonth() &&
    date.getDate() === d.getDate();
}

async function selectDate(date) {
  if (isNaN(date)) return;

  // 更新选中的日期
  emit('update:modelValue', date);
  emit('select', date);

  // 触发日期点击事件，传递日期和是否是今天的标志
  const dateInfo = {
    date: date,
    isToday: isToday(date)
  };
  emit('date-click', dateInfo);

  // 清除之前的所有缓存
  clearAllCache();

  // 根据日期类型加载任务（同时加载个人和组织代办）
  let personalTasks = [];
  let organizationTasks = [];
  let upcomingPersonalTasks = [];
  let upcomingOrganizationTasks = [];
  let type = '';

  // 获取有活动的组织列表（仅今天）
  let activeOrgs = [];
  if (dateInfo.isToday) {
    activeOrgs = await loadActiveOrganizations() || [];
  }

  if (dateInfo.isToday) {
    // 今天：加载今日待办 + 已完成待办 + 即将开始/结束的待办
    console.log('加载今日数据');
    type = 'today';

    // 个人代办
    const todayPersonalTodos = await loadTodayTodos(date) || [];
    const completedPersonalTodos = await loadCompletedTodos(date) || [];
    personalTasks = [...todayPersonalTodos, ...completedPersonalTodos];

    const comingStartPersonal = await loadComingStartTodos() || [];
    const comingEndPersonal = await loadComingEndTodos() || [];

    // 组织代办
    const todayOrganizationTodos = await loadTodayOrganizationTodos() || [];
    const completedOrganizationTodos = await loadCompletedOrganizationTodos(date) || [];
    organizationTasks = [...todayOrganizationTodos, ...completedOrganizationTodos];

    const comingStartOrganization = await loadComingStartOrganizationTodos() || [];
    const comingEndOrganization = await loadComingEndOrganizationTodos() || [];

    // 处理重复任务并添加标识
    // 个人待办即将任务处理
    const personalStartIds = new Set(comingStartPersonal.map(task => task.id));
    const uniquePersonalEnd = comingEndPersonal.filter(task => !personalStartIds.has(task.id));

    const personalStartWithTag = comingStartPersonal.map(task => ({
      ...task,
      creator_organ_id: 0, // 确保是个人任务
      creator_user_id: task.creator_user_id || 1, // 确保有用户ID
      isComingStart: true,
      isComingEnd: false,
      sortTime: task.start_time || task.startTime || task.createdAt
    }));

    const personalEndWithTag = uniquePersonalEnd.map(task => ({
      ...task,
      creator_organ_id: 0, // 确保是个人任务
      creator_user_id: task.creator_user_id || 1,
      isComingStart: false,
      isComingEnd: true,
      sortTime: task.end_time || task.endTime || task.createdAt
    }));

    upcomingPersonalTasks = [...personalStartWithTag, ...personalEndWithTag]
      .sort((a, b) => new Date(a.sortTime) - new Date(b.sortTime));

    // 组织待办即将任务处理
    const organizationStartIds = new Set(comingStartOrganization.map(task => task.id));
    const uniqueOrganizationEnd = comingEndOrganization.filter(task => !organizationStartIds.has(task.id));

    const organizationStartWithTag = comingStartOrganization.map(task => ({
      ...task,
      creator_organ_id: task.creator_organ_id || task.organization_id || 1, // 确保有组织ID
      creator_user_id: 0, // 组织任务用户ID为0
      isComingStart: true,
      isComingEnd: false,
      sortTime: task.start_time || task.startTime || task.createdAt
    }));

    const organizationEndWithTag = uniqueOrganizationEnd.map(task => ({
      ...task,
      creator_organ_id: task.creator_organ_id || task.organization_id || 1,
      creator_user_id: 0,
      isComingStart: false,
      isComingEnd: true,
      sortTime: task.end_time || task.endTime || task.createdAt
    }));

    upcomingOrganizationTasks = [...organizationStartWithTag, ...organizationEndWithTag]
      .sort((a, b) => new Date(a.sortTime) - new Date(b.sortTime));
  } else if (isAfterToday(date)) {
    // 今天之后：加载某一天开始的待办
    console.log('加载未来日期数据');
    type = 'future';
    personalTasks = await loadOneDayTodos(date) || [];
    organizationTasks = await loadOneDayOrganizationTodos(date) || [];
  } else {
    // 今天之前：加载已完成待办 + 过期待办
    console.log('加载过去日期数据');
    type = 'completed';

    // 个人代办
    const completedPersonalTodos = await loadCompletedTodos(date) || [];
    const expiredPersonalTodos = await loadOneDayExpiredTodos(date) || [];

    // 合并个人代办结果并去重
    const personalTodosMap = new Map();
    completedPersonalTodos.forEach(todo => {
      personalTodosMap.set(todo.id, {
        ...todo,
        creator_organ_id: 0, // 确保是个人任务
        creator_user_id: todo.creator_user_id || 1,
        isExpired: false,
        status: 'completed'
      });
    });

    expiredPersonalTodos.forEach(todo => {
      if (!personalTodosMap.has(todo.id)) {
        personalTodosMap.set(todo.id, {
          ...todo,
          creator_organ_id: 0, // 确保是个人任务
          creator_user_id: todo.creator_user_id || 1,
          isExpired: true,
          status: todo.status || 'expired'
        });
      } else {
        personalTodosMap.set(todo.id, {
          ...personalTodosMap.get(todo.id),
          isExpired: false,
          status: 'completed'
        });
      }
    });

    personalTasks = Array.from(personalTodosMap.values());

    // 组织代办
    const completedOrganizationTodos = await loadCompletedOrganizationTodos(date) || [];
    const expiredOrganizationTodos = await loadOneDayExpiredOrganizationTodos(date) || [];

    // 合并组织代办结果并去重
    const organizationTodosMap = new Map();
    completedOrganizationTodos.forEach(todo => {
      organizationTodosMap.set(todo.id, {
        ...todo,
        creator_organ_id: todo.creator_organ_id || todo.organization_id || 1,
        creator_user_id: 0, // 组织任务用户ID为0
        isExpired: false,
        status: 'completed'
      });
    });

    expiredOrganizationTodos.forEach(todo => {
      if (!organizationTodosMap.has(todo.id)) {
        organizationTodosMap.set(todo.id, {
          ...todo,
          creator_organ_id: todo.creator_organ_id || todo.organization_id || 1,
          creator_user_id: 0,
          isExpired: true,
          status: todo.status || 'expired'
        });
      } else {
        organizationTodosMap.set(todo.id, {
          ...organizationTodosMap.get(todo.id),
          isExpired: false,
          status: 'completed'
        });
      }
    });

    organizationTasks = Array.from(organizationTodosMap.values());
  }

  // 合并所有任务
  const allTasks = [...personalTasks, ...organizationTasks];
  const allUpcomingTasks = [...upcomingPersonalTasks, ...upcomingOrganizationTasks];

  // 缓存当前日期数据
  cacheCurrentDateTasks(date, allTasks, allUpcomingTasks, type, activeOrgs);

  // 将任务数据传递给父组件
  emit('load-tasks', {
    date: date,
    tasks: allTasks,
    upcomingTasks: allUpcomingTasks,
    type: type,
    activeOrganizations: activeOrgs
  });
}

// 暴露方法给父组件
defineExpose({
  clearCache: clearAllCache,
  reloadDate: async (date) => {
    const token = getToken();
    if (!token) {
      console.error('未找到认证令牌');
      return;
    }

    // 先检查并更新代办
    await checkAndUpdateTodos(token);

    // 清除缓存
    clearAllCache();

    // 重新加载任务
    let personalTasks = [];
    let organizationTasks = [];
    let upcomingPersonalTasks = [];
    let upcomingOrganizationTasks = [];
    let type = '';

    // 获取有活动的组织列表（仅今天）
    let activeOrgs = [];
    if (isToday(date)) {
      activeOrgs = await loadActiveOrganizations() || [];
    }

    if (isToday(date)) {
      type = 'today';

      // 个人代办
      personalTasks = await loadTodayTodos(date) || [];
      const comingStartPersonal = await loadComingStartTodos() || [];
      const comingEndPersonal = await loadComingEndTodos() || [];

      // 组织代办
      organizationTasks = await loadTodayOrganizationTodos() || [];
      const comingStartOrganization = await loadComingStartOrganizationTodos() || [];
      const comingEndOrganization = await loadComingEndOrganizationTodos() || [];

      // 处理即将开始/结束的个人任务
      const personalStartIds = new Set(comingStartPersonal.map(task => task.id));
      const uniquePersonalEnd = comingEndPersonal.filter(task => !personalStartIds.has(task.id));

      const personalStartWithTag = comingStartPersonal.map(task => ({
        ...task,
        creator_organ_id: 0,
        creator_user_id: task.creator_user_id || 1,
        isComingStart: true,
        isComingEnd: false,
        sortTime: task.start_time || task.startTime || task.createdAt
      }));

      const personalEndWithTag = uniquePersonalEnd.map(task => ({
        ...task,
        creator_organ_id: 0,
        creator_user_id: task.creator_user_id || 1,
        isComingStart: false,
        isComingEnd: true,
        sortTime: task.end_time || task.endTime || task.createdAt
      }));

      upcomingPersonalTasks = [...personalStartWithTag, ...personalEndWithTag]
        .sort((a, b) => new Date(a.sortTime) - new Date(b.sortTime));

      // 处理即将开始/结束的组织任务
      const organizationStartIds = new Set(comingStartOrganization.map(task => task.id));
      const uniqueOrganizationEnd = comingEndOrganization.filter(task => !organizationStartIds.has(task.id));

      const organizationStartWithTag = comingStartOrganization.map(task => ({
        ...task,
        creator_organ_id: task.creator_organ_id || task.organization_id || 1,
        creator_user_id: 0,
        isComingStart: true,
        isComingEnd: false,
        sortTime: task.start_time || task.startTime || task.createdAt
      }));

      const organizationEndWithTag = uniqueOrganizationEnd.map(task => ({
        ...task,
        creator_organ_id: task.creator_organ_id || task.organization_id || 1,
        creator_user_id: 0,
        isComingStart: false,
        isComingEnd: true,
        sortTime: task.end_time || task.endTime || task.createdAt
      }));

      upcomingOrganizationTasks = [...organizationStartWithTag, ...organizationEndWithTag]
        .sort((a, b) => new Date(a.sortTime) - new Date(b.sortTime));
    } else if (isAfterToday(date)) {
      type = 'future';
      personalTasks = await loadOneDayTodos(date) || [];
      organizationTasks = await loadOneDayOrganizationTodos(date) || [];
    } else {
      type = 'completed';

      // 个人代办
      const completedPersonalTodos = await loadCompletedTodos(date) || [];
      const expiredPersonalTodos = await loadOneDayExpiredTodos(date) || [];

      // 合并个人代办结果并去重
      const personalTodosMap = new Map();
      completedPersonalTodos.forEach(todo => {
        personalTodosMap.set(todo.id, {
          ...todo,
          creator_organ_id: 0,
          creator_user_id: todo.creator_user_id || 1,
          isExpired: false,
          status: 'completed'
        });
      });

      expiredPersonalTodos.forEach(todo => {
        if (!personalTodosMap.has(todo.id)) {
          personalTodosMap.set(todo.id, {
            ...todo,
            creator_organ_id: 0,
            creator_user_id: todo.creator_user_id || 1,
            isExpired: true,
            status: todo.status || 'expired'
          });
        } else {
          personalTodosMap.set(todo.id, {
            ...personalTodosMap.get(todo.id),
            isExpired: false,
            status: 'completed'
          });
        }
      });

      personalTasks = Array.from(personalTodosMap.values());

      // 组织代办
      const completedOrganizationTodos = await loadCompletedOrganizationTodos(date) || [];
      const expiredOrganizationTodos = await loadOneDayExpiredOrganizationTodos(date) || [];

      // 合并组织代办结果并去重
      const organizationTodosMap = new Map();
      completedOrganizationTodos.forEach(todo => {
        organizationTodosMap.set(todo.id, {
          ...todo,
          creator_organ_id: todo.creator_organ_id || todo.organization_id || 1,
          creator_user_id: 0,
          isExpired: false,
          status: 'completed'
        });
      });

      expiredOrganizationTodos.forEach(todo => {
        if (!organizationTodosMap.has(todo.id)) {
          organizationTodosMap.set(todo.id, {
            ...todo,
            creator_organ_id: todo.creator_organ_id || todo.organization_id || 1,
            creator_user_id: 0,
            isExpired: true,
            status: todo.status || 'expired'
          });
        } else {
          organizationTodosMap.set(todo.id, {
            ...organizationTodosMap.get(todo.id),
            isExpired: false,
            status: 'completed'
          });
        }
      });

      organizationTasks = Array.from(organizationTodosMap.values());
    }

    // 合并所有任务
    const allTasks = [...personalTasks, ...organizationTasks];
    const allUpcomingTasks = [...upcomingPersonalTasks, ...upcomingOrganizationTasks];

    // 缓存数据
    cacheCurrentDateTasks(date, allTasks, allUpcomingTasks, type, activeOrgs);

    // 将任务数据传递给父组件
    emit('load-tasks', {
      date: date,
      tasks: allTasks,
      upcomingTasks: allUpcomingTasks,
      type: type,
      activeOrganizations: activeOrgs
    });

    return { tasks: allTasks, upcomingTasks: allUpcomingTasks };
  },
  openActivityModal // 暴露打开活动详情弹窗的方法
});

// 检查并更新代办
async function checkAndUpdateTodos(token) {
  try {
    // 检查是否需要更新
    const lastUpdate = localStorage.getItem('last_todo_update');
    const today = new Date().toDateString();

    // 如果没有更新记录或者不是今天更新的，才调用接口
    if (!lastUpdate || lastUpdate !== today) {
      console.log('开始调用更新代办接口...');

      const response = await fetch(`${API_BASE}/todos/updateTodos`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });

      if (response.ok) {
        const result = await response.json();
        console.log('更新代办接口响应:', result);

        if (result.success) {
          // 更新成功，记录更新时间
          localStorage.setItem('last_todo_update', today);
          console.log('代办更新成功，已记录更新时间:', today);
        } else {
          console.error('更新代办失败:', result.message);
        }
      } else {
        console.error('调用更新代办接口失败:', response.status);
      }
    } else {
      console.log('今天已经更新过代办，跳过更新');
    }
  } catch (error) {
    console.error('检查更新代办失败:', error);
  }
}

// 加载某一天过期的待办
async function loadOneDayExpiredTodos(date) {
  const token = getToken();
  if (!token) return null;

  try {
    // 使用相同的本地日期格式
    const dateStr = formatDateToLocalString(date);
    console.log('开始调用某一天过期的待办接口...', dateStr);

    const response = await fetch(`${API_BASE}/todos/get-OneDayExpiredTodos?date=${dateStr}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    if (response.ok) {
      const result = await response.json();
      console.log('某一天过期的待办接口返回的全部内容:', result);
      if (result.success && result.todos) {
        // 缓存活动数据（如果是组织活动）
        result.todos.forEach(todo => {
          if (todo.creator_organ_id && todo.creator_organ_id > 0) {
            cacheActivityData(todo)
          }
        });

        return result.todos;
      }
    } else {
      console.error('获取某一天过期的待办失败:', response.status);
    }
  } catch (error) {
    console.error('调用某一天过期的待办接口失败:', error);
  }
  return null;
}
</script>

<style scoped>
/* 卡片式网格布局 */

.month-calendar {
  width: 100%;
  text-align: center;
  color: #fff; /* 适应深色背景，文本改为白色 */
}

/* 头部 (月份和切换按钮) */
.calendar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px; /* 增加与日历主体的间距 */
  padding: 5px 10px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 10px;
}

.header-title {
  font-size: 1.5rem; /* 字体更大 */
  font-weight: 600; /* 字体加粗 */
}

.nav-btn {
  /* 背景为玻璃材质 */
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.4);
  color: #fff;
  border-radius: 50%;
  width: 32px;
  height: 32px;
  font-size: 1.2rem;
  cursor: pointer;
  transition: all 0.3s;
}

.nav-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  border-color: rgba(255, 255, 255, 0.6);
}

/* 分类信息显示 */
.categories-display {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-bottom: 20px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.category-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.category-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  display: inline-block;
}

.category-name {
  font-size: 0.9rem;
  color: rgba(255, 255, 255, 0.8);
}

/* 日历表格主体 */
.calendar-table {
  width: 100%;
  /* * 关键：
   * 1. border-collapse: separate  -> 允许单元格有间距
   * 2. border-spacing: 8px       -> 设置单元格之间的间距
   */
  border-collapse: separate;
  border-spacing: 8px;
  table-layout: fixed;
}

/* 星期表头 (SUN, MON...) */
.calendar-table th {
  font-weight: 600; /* 加粗 */
  font-size: 0.8rem; /* 字体小一点 */
  color: rgba(255, 255, 255, 0.7);/* 颜色为浅白色 */
  text-transform: uppercase; /* 大写 */
  padding-bottom: 8px; /* 与日期卡片拉开距离 */
}

/* 日期单元格 (卡片) */
.calendar-table td {
  height: 50px;  /* 固定高度，使其更像卡片 */
  text-align: left; /* 内容左上角对齐 */
  vertical-align: top; /* 内容左上角对齐 */

  /* === 【核心修改】每个单元格的毛玻璃效果 === */
  background: rgba(255, 255, 255, 0.05); /* 更透明的背景 */
  backdrop-filter: blur(5px); /* 单元格自身的轻微模糊 */
  -webkit-backdrop-filter: blur(5px);
  border: 1px solid rgba(255, 255, 255, 0.2); /* 浅边框 */

  border-radius: 16px;
  padding: 8px;
  cursor: pointer;
  transition: all 0.3s ease;

  /* 添加轻微阴影以营造卡片感 */
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.03);
}

/* 空白单元格 (用于占位) */
.calendar-table td.empty {
  background: transparent; /* 透明背景 */
  border: none;
  box-shadow: none;
  cursor: default;
  /* * visibility: hidden
   * 保留单元格的布局位置，但使其不可见
   * (这比 v-if 好，因为它能保持网格对齐)
   */
  visibility: hidden;
}

/* "今天" 的样式 */
.calendar-table td.today {
  background-color: #e0f0ff; /* 你的旧 "today" 颜色 */
  border-color: #007bff; /* 蓝色边框高亮 */
}

/* "选中" 的样式 */
.calendar-table td.selected {
  background-color: #007bff;
  border-color: #007bff;
  color: #fff; /* 选中时，内部文字变白 */
  box-shadow: 0 4px 10px rgba(0, 123, 255, 0.3);
}

/* 鼠标悬停效果 (非空白格) */
.calendar-table td:not(.empty):hover {
  transform: translateY(-3px); /* 悬停时轻微上浮 */
  /* 悬停时背景更亮 */
  background: rgba(255, 255, 255, 0.15);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15);
}

/* 让选中状态立即生效，即使鼠标还在上面 */
.calendar-table td.selected,
.calendar-table td.selected:hover {
  background-color: #007bff !important;
  border-color: #007bff !important;
  color: #fff !important;
  transform: none !important; /* 取消悬停上浮 */
  box-shadow: 0 4px 10px rgba(0, 123, 255, 0.3) !important;
}


/* 日期数字 */
.day-number {
  font-size: 1rem;
  font-weight: 600;
  color: #fff;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
}

/* "今天" 状态下的日期数字 */
.today .day-number {
  color: #007bff; /* 蓝色 */
}

/* "选中" 状态下的日期数字 */
.selected .day-number {
  color: #fff; /* 白色 */
}

/* "任务" 样式占位 (就像图片中的 "8 tasks") */
/*
.task-info {
  margin-top: 8px;
  font-size: 0.8rem;
  color: #007bff;
}
.selected .task-info {
  color: #fff;
  opacity: 0.8;
}
*/
</style>
