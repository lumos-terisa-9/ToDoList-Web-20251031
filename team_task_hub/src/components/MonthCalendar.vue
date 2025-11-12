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
import { ref, computed } from 'vue'

const props = defineProps({
  modelValue: { type: Date, required: true }
})
const emit = defineEmits(['update:modelValue', 'select'])

// 脚本 (Logic) 部分与你原来的一致，无需修改
const current = ref(new Date(props.modelValue))
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

function prevMonth() {
  current.value = new Date(year.value, month.value - 1, 1)
}
function nextMonth() {
  current.value = new Date(year.value, month.value + 1, 1)
}
function isToday(date) {
  if (isNaN(date)) return false
  const t = new Date()
  return date.getFullYear() === t.getFullYear() &&
    date.getMonth() === t.getMonth() &&
    date.getDate() === t.getDate()
}
function isSelected(date) {
  if (isNaN(date)) return false
  const d = props.modelValue
  return date.getFullYear() === d.getFullYear() &&
    date.getMonth() === d.getMonth() &&
    date.getDate() === d.getDate()
}
function selectDate(date) {
  if (isNaN(date)) return
  emit('update:modelValue', date)
  emit('select', date)
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
  height: 65px;  /* 固定高度，使其更像卡片 */
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
