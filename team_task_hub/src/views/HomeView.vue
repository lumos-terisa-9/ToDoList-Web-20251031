<template>
  <div class="home-container">
    <div class="split-pane">
      <!-- 左侧：日历 -->
      <MonthCalendar
        v-model="picked"
        @select="onSelect"
        class="calendar-pane"
      />

      <!-- 右侧：任务区 -->
      <div class="task-pane">
        <h2>任务面板</h2>
        <p>已选日期：{{ picked.toLocaleDateString() }}</p>
        <div class="task-content">
          <p>这里可以展示任务列表、表单或统计信息。</p>
          <div style="height: 800px;"></div>
        </div>
      </div>
    </div>
  </div>

  <div class="home-layout">
    <h1 style="color:#fff">日历应该出现在下方</h1>
    <!-- 不要写空冒号，想传属性就写完整 -->
    <MonthCalendar v-model="picked" @select="onSelect" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import MonthCalendar from '@/components/MonthCalendar.vue'

const picked = ref(new Date())
function onSelect(d) {
  picked.value = d
  console.log('选中日期：', d)
}
</script>

<style scoped>
.home-container {
  height: calc(100vh - 60px); /* 留出 HeaderBar 高度 */
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  background: linear-gradient(to bottom, #012855, #0e59b8);
}

/* 分栏容器 */
.split-pane {
  display: flex;
  flex: 1;
  align-items: stretch;
}

/* 左侧：日历 */
.calendar-pane {
  flex: 1;
  max-width: 50%;
  background: #fff;
  border-radius: 16px 0 0 16px;
  overflow-y: auto;
  padding: 16px;
  box-sizing: border-box;
}

/* 右侧：任务面板 */
.task-pane {
  flex: 1;
  background: linear-gradient(to bottom right, #0e59b8, #013a7c);
  color: #fff;
  border-radius: 0 16px 16px 0;
  padding: 32px;
  overflow-y: auto;
  box-sizing: border-box;
}

/* 分隔线 */
.split-pane::before {
  content: "";
  width: 1px;
  background: rgba(255,255,255,0.2);
  align-self: stretch;
}

/* 任务内容滚动区 */
.task-content {
  margin-top: 16px;
  background: rgba(255,255,255,0.05);
  border-radius: 8px;
  padding: 16px;
}
</style>
