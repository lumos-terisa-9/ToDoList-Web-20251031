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

        <ToDoList
          :id="picked.toDateString()"
          :title="'日期待办事项'"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import MonthCalendar from '@/components/MonthCalendar.vue'
import ToDoList from '@/components/ToDoList.vue'

const picked = ref(new Date())
function onSelect(d) {
  picked.value = d
  console.log('选中日期：', d)
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
  padding: 20px 30px;
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
</style>
