<template>
  <div class="month-calendar">
    <div class="calendar-header">
      <button @click="prevMonth">‹</button>
      <span>{{ year }} 年 {{ month + 1 }} 月</span>
      <button @click="nextMonth">›</button>
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
            :class="{ today: isToday(date), selected: isSelected(date) }"
            @click="selectDate(date)">
          {{ date.getDate() }}
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

const current = ref(new Date(props.modelValue))
const days = ['日', '一', '二', '三', '四', '五', '六']

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
.month-calendar {
  width: 100%;
  text-align: center;
}

.calendar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-weight: bold;
}

.calendar-table {
  width: 100%;
  border-collapse: collapse;
}

.calendar-table th {
  font-weight: normal;
  color: #666;
}

.calendar-table td {
  width: 40px;
  height: 40px;
  text-align: center;
  border-radius: 6px;
  cursor: pointer;
}

.calendar-table td.today {
  background-color: #e0f0ff;
  font-weight: bold;
}

.calendar-table td.selected {
  background-color: #007bff;
  color: #fff;
}
</style>
