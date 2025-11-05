<template>
  <div class="app-background">
    <HeaderBar />

    <div class="scroll-content"
         ref="scrollContent"
         @touchstart="startTouch"
         @touchmove="moveTouch"
         @touchend="endTouch"
         @mousedown.prevent="startMouse"
         @mousemove="moveMouse"
         @mouseup="endMouse"
         @mouseleave="handleMouseLeave">
      <!-- 用 router-view 显示 HomeView.vue 等页面 -->
      <router-view />
    </div>
<!--      <main class="page-content">-->
<!--        <h1>欢迎来到team_task_hub</h1>-->
<!--        <p>using vue3</p>-->
<!--        <div style="height:1000px;"></div>-->
<!--      </main>-->
<!--    </div>-->

    <transition name="fade">
      <div v-if="showTopTip" class="scroll-tip top-tip" :style="{ transform: `translate(-50%, ${topTipOffset}px)` }">
        我已经到顶啦
      </div>
    </transition>

    <transition name="fade">
      <div v-if="showBottomTip" class="scroll-tip bottom-tip" :style="{ transform: `translate(-50%, ${bottomTipOffset}px)` }">
        我已经到底啦
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from "vue";
import HeaderBar from './components/HeaderBar.vue'

// 节流函数（无需装 lodash，手写 16 ms 够用）
const throttle = (fn, delay = 16) => {
  let last = 0
  return (...args) => {
    const now = Date.now()
    if (now - last > delay) {
      last = now
      fn.apply(this, args)
    }
  }
}

const scrollContent = ref(null)

const startY = ref(0)
const offsetY = ref(0)
const isDragging = ref(false)
const isBouncing = ref(false)

const showTopTip = ref(false)
const showBottomTip = ref(false)
const topTipOffset = ref(0)
const bottomTipOffset = ref(0)

function setupScrollContainer() {
  const el = scrollContent.value
  if (!el) return

  const viewportHeight = window.innerHeight
  el.style.height = viewportHeight + 'px'
  el.style.overflowY = 'auto'
}

// 监听窗口大小变化
function handleResize() {
  setupScrollContainer()
}

// 软键盘弹起/收起重新计算高度
function handleVisualResize() {
  setupScrollContainer()
}

onMounted(() => {
  setupScrollContainer()
  window.addEventListener('resize', handleResize)
  // 添加全局鼠标抬起监听，解决鼠标移出窗口的问题
  // 3. 滚轮性能优化
  const el = scrollContent.value
  el.addEventListener("wheel", handleWheel, { passive: false })
  window.visualViewport?.addEventListener("resize", handleVisualResize)

  onUnmounted(() => {
    window.removeEventListener("resize", handleResize)
    el.removeEventListener("wheel", handleWheel, { passive: false })
    window.visualViewport?.removeEventListener("resize", handleVisualResize)
  })
})

// 全局鼠标抬起处理

// 鼠标移出元素处理
function handleMouseLeave() {
  if (isDragging.value && !isBouncing.value) {
    releaseDrag()
  }
}

// 显示提示文字并自动淡出
function showTip(tipRef) {
  if (tipRef.value) return
  tipRef.value = true
  setTimeout(() => tipRef.value = false, 1500)
}

// Touch
function startTouch(e) {
  if (isBouncing.value) return
  isDragging.value = true
  startY.value = e.touches[0].clientY
}
function moveTouchRaw(e) {
  if (!isDragging.value || isBouncing.value) return
  const dy = e.touches[0].clientY - startY.value
  handleDrag(dy)
}
const moveTouch = throttle(moveTouchRaw, 16)

function endTouch() {
  if (!isDragging.value || isBouncing.value) return
  releaseDrag()
}

// Mouse
function startMouse(e) {
  if (isBouncing.value) return
  isDragging.value = true
  startY.value = e.clientY
}
function moveMouseRaw(e) {
  if (!isDragging.value || isBouncing.value) return
  const dy = e.clientY - startY.value
  handleDrag(dy)
}
const moveMouse = throttle(moveMouseRaw, 16)

function endMouse() {
  if (!isDragging.value || isBouncing.value) return
  releaseDrag()
}

// Wheel - 只在边界时阻止默认行为
function handleWheel(e) {
  if (isBouncing.value) {
    e.preventDefault()
    return
  }

  const el = scrollContent.value
  if (!el) return

  const scrollTop = el.scrollTop
  const scrollHeight = el.scrollHeight
  const clientHeight = el.clientHeight
  const maxScroll = Math.max(0, scrollHeight - clientHeight)

  // 检测边界 - 只在边界情况下阻止默认行为
  const atTop = scrollTop <= 0 && e.deltaY < 0
  const atBottom = scrollTop >= maxScroll && e.deltaY > 0 && maxScroll > 0

  if (atTop || atBottom) {
    e.preventDefault()
    simulateDrag(-e.deltaY)
  }
  // 其他情况让浏览器正常处理滚动
}

// 模拟拖拽动画
function simulateDrag(dy) {
  if (isBouncing.value) return // 已经在回弹中，不再触发新的回弹

  const el = scrollContent.value
  if (!el) return

  const scrollTop = el.scrollTop
  const scrollHeight = el.scrollHeight
  const clientHeight = el.clientHeight
  const maxScroll = Math.max(0, scrollHeight - clientHeight)

  const atTop = scrollTop <= 0 && dy > 0
  const atBottom = scrollTop >= maxScroll && dy < 0 && maxScroll > 0

  if (atTop || atBottom) {
    isBouncing.value = true // 开始回弹

    // 保存当前滚动位置，用于固定滚动条
    const currentScrollTop = el.scrollTop

    offsetY.value = dy / 2
    el.style.transform = `translateY(${offsetY.value}px)`

    // 固定滚动条位置但不隐藏
    el.scrollTop = currentScrollTop

    if (atTop) {
      showTip(showTopTip)
      topTipOffset.value = offsetY.value / 2
    }
    if (atBottom) {
      showTip(showBottomTip)
      bottomTipOffset.value = offsetY.value / 2
    }

    // 自动回弹
    setTimeout(() => {
      el.style.transition = 'transform 1s ease'
      el.style.transform = 'translateY(0)'
      setTimeout(() => {
        el.style.transition = ''
        // 恢复滚动位置
        el.scrollTop = currentScrollTop
        offsetY.value = 0
        topTipOffset.value = 0
        bottomTipOffset.value = 0
        isBouncing.value = false // 回弹结束
      }, 1000) // 改为1000ms，与过渡时间保持一致
    }, 50)
  }
}

// 拖拽逻辑
function handleDrag(dy) {
  if (isBouncing.value) return

  const el = scrollContent.value
  if (!el) return

  const scrollTop = el.scrollTop
  const scrollHeight = el.scrollHeight
  const clientHeight = el.clientHeight
  const maxScroll = Math.max(0, scrollHeight - clientHeight)

  const atTop = scrollTop <= 0 && dy > 0   // 顶部向下
  const atBottom = scrollTop >= maxScroll && dy < 0 && maxScroll > 0 // 底部向上

  if (atTop || atBottom) {
    // 保存当前滚动位置，用于固定滚动条
    const currentScrollTop = el.scrollTop

    offsetY.value = dy / 2
    el.style.transform = `translateY(${offsetY.value}px)`

    // 固定滚动条位置但不隐藏
    el.scrollTop = currentScrollTop

    if (atTop) {
      showTip(showTopTip)
      topTipOffset.value = offsetY.value / 2
    }
    if (atBottom) {
      showTip(showBottomTip)
      bottomTipOffset.value = offsetY.value / 2
    }
  }
}

// 回弹
function releaseDrag() {
  if (isBouncing.value) return // 回弹中不处理释放

  const el = scrollContent.value
  if (!el) return
  isBouncing.value = true // 开始回弹

  // 保存当前滚动位置
  const currentScrollTop = el.scrollTop

  el.style.transition = 'transform 1s ease'
  el.style.transform = 'translateY(0)'
  setTimeout(() => {
    el.style.transition = ''
    // 恢复滚动位置
    el.scrollTop = currentScrollTop
    isBouncing.value = false // 回弹结束
  }, 1000) // 改为1000ms，与过渡时间保持一致

  offsetY.value = 0
  topTipOffset.value = 0
  bottomTipOffset.value = 0
  isDragging.value = false
}
</script>

<style>
.app-background {
  min-height: 100vh;
  background: linear-gradient(to bottom, #012855, #0e59b8);
  font-family: "Helvetica Neue", Arial, "PingFang SC", "Microsoft YaHei", sans-serif;
  overflow: hidden;
  position: relative;
}

.scroll-content {
  width: 100%;
  /* 高度通过JS动态设置 */
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
  scroll-behavior: smooth;
  /*添加 60px 的顶部内边距，为 HeaderBar 留出空间
   *设置 box-sizing，确保 padding 不会增加 JS 设置的 height
   */
  padding-top: 60px;
  box-sizing: border-box;
}

/* 自定义细蓝色滚动条 */
.scroll-content::-webkit-scrollbar {
  width: 6px; /* 滚动条宽度 */
}

.scroll-content::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.1); /* 轨道背景 */
  border-radius: 3px;
}

.scroll-content::-webkit-scrollbar-thumb {
  background: linear-gradient(to bottom, #1e90ff, #0066cc); /* 蓝色渐变 */
  border-radius: 3px;
  transition: background 0.3s ease;
}

.scroll-content::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(to bottom, #4da6ff, #0080ff); /* 悬停时更亮的蓝色 */
}

/* Firefox 滚动条样式 */
.scroll-content {
  scrollbar-width: thin;
  scrollbar-color: #1e90ff rgba(255, 255, 255, 0.1);
}

.page-content {
  padding: 24px;
  color: white;
  margin-top: 120px;
  text-align: center;
}

.scroll-tip {
  position: fixed;
  left: 50%;
  color: white;
  font-size: 13px;
  background: rgba(0,0,0,0.3);
  padding: 6px 12px;
  border-radius: 8px;
  z-index: 2000;
  pointer-events: none;
}

.top-tip { top: 100px; }
.bottom-tip { bottom: 50px; }

.fade-enter-active, .fade-leave-active { transition: opacity 0.5s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
