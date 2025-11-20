<template>
  <div class="app-container">
    <HeaderBar @open-login="openLoginModal" @show-profile="showProfileModal" />

    <div class="scroll-content" ref="scrollContent">
      <router-view v-slot="{ Component }">
        <component
          :is="Component"
          ref="pageComponent"
          @open-login="openLoginModal"
        />
      </router-view>
    </div>

    <transition name="fade">
      <div v-if="showTopTip" class="scroll-tip top-tip">
        我已经到顶啦
      </div>
    </transition>

    <transition name="fade">
      <div v-if="showBottomTip" class="scroll-tip bottom-tip">
        我已经到底啦
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from "vue";
import HeaderBar from './components/HeaderBar.vue'

const scrollContent = ref(null)
const pageComponent = ref(null)

const showTopTip = ref(false)
const showBottomTip = ref(false)

function openLoginModal() {
  if (pageComponent.value && pageComponent.value.openLoginModal) {
    pageComponent.value.openLoginModal()
  }
}

function showProfileModal() {
  if (pageComponent.value && pageComponent.value.showProfileModal) {
    pageComponent.value.showProfileModal()
  }
}

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

// 显示提示文字并自动淡出
function showTip(tipRef) {
  if (tipRef.value) return
  tipRef.value = true
  setTimeout(() => tipRef.value = false, 1500)
}

// Wheel - 检测边界显示提示
function handleWheel(e) {
  const el = scrollContent.value
  if (!el) return

  const scrollTop = el.scrollTop
  const scrollHeight = el.scrollHeight
  const clientHeight = el.clientHeight
  const maxScroll = Math.max(0, scrollHeight - clientHeight)

  // 检测边界 - 修复逻辑
  const atTop = scrollTop <= 0 && e.deltaY < 0
  const atBottom = scrollTop >= maxScroll && e.deltaY > 0

  if (atTop) {
    showTip(showTopTip)
  }
  if (atBottom && maxScroll > 0) { // 确保有可滚动内容
    showTip(showBottomTip)
  }
}

// 添加滚动事件监听来检测边界
function handleScroll() {
  const el = scrollContent.value
  if (!el) return

  const scrollTop = el.scrollTop
  const scrollHeight = el.scrollHeight
  const clientHeight = el.clientHeight
  const maxScroll = Math.max(0, scrollHeight - clientHeight)

  // 检测是否到达底部
  const atBottom = scrollTop >= maxScroll-0.66 && maxScroll > 0

  if (atBottom && !showBottomTip.value) {
    showTip(showBottomTip)
  }
}

onMounted(() => {
  setupScrollContainer()
  window.addEventListener('resize', handleResize)
  const el = scrollContent.value
  el.addEventListener("wheel", handleWheel, { passive: false })
  el.addEventListener("scroll", handleScroll)
  window.visualViewport?.addEventListener("resize", handleVisualResize)

  onUnmounted(() => {
    window.removeEventListener("resize", handleResize)
    el.removeEventListener("wheel", handleWheel, { passive: false })
    el.removeEventListener("scroll", handleScroll)
    window.visualViewport?.removeEventListener("resize", handleVisualResize)
  })
})
</script>

<style>
.app-container {
  min-height: 100vh;
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
  box-sizing: border-box;
}

/* ========================================= */
/* === 通用滚动条样式 - 由子页面背景决定 === */
/* ========================================= */

/* 首页样式 - 紫色渐变背景 */
.home-page .scroll-content::-webkit-scrollbar-track {
  background: rgba(39, 145, 251, 0.2);
}

.home-page .scroll-content::-webkit-scrollbar-thumb {
  background: linear-gradient(to bottom, rgba(120, 180, 255, 0.6), rgba(153, 64, 250, 0.6));
}

.home-page .scroll-content::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(to bottom, rgba(140, 200, 255, 0.8), rgba(173, 84, 255, 0.8));
}

.home-page .scroll-content {
  scrollbar-color: rgba(120, 180, 255, 0.6) rgba(39, 145, 251, 0.2);
}

/* 个人页面样式 - 蓝色渐变背景 */
.person-page .scroll-content::-webkit-scrollbar-track {
  background: rgba(14, 89, 184, 0.2);
}

.person-page .scroll-content::-webkit-scrollbar-thumb {
  background: linear-gradient(to bottom, rgba(30, 144, 255, 0.6), rgba(22, 177, 244, 0.6));
}

.person-page .scroll-content::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(to bottom, rgba(77, 166, 255, 0.8), rgba(52, 197, 255, 0.8));
}

.person-page .scroll-content {
  scrollbar-color: rgba(30, 144, 255, 0.6) rgba(14, 89, 184, 0.2);
}

/* 基础滚动条样式 */
.scroll-content::-webkit-scrollbar {
  width: 6px;
}

.scroll-content::-webkit-scrollbar-track {
  border-radius: 3px;
  transition: background 0.3s ease;
}

.scroll-content::-webkit-scrollbar-thumb {
  border-radius: 3px;
  transition: background 0.3s ease;
}

.scroll-content {
  scrollbar-width: thin;
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
  font-size: 14px;
  background: rgba(34, 34, 34, 0.6);
  padding: 8px 16px;
  border-radius: 8px;
  z-index: 2000;
  pointer-events: none;
  white-space: nowrap;
  line-height: 1.2;
  height: auto;
  min-height: auto;
}

.top-tip {
  top: 100px;
  transform: translateX(-70%);
}

.bottom-tip {
  bottom: 45px;
  transform: translateX(-70%);
}

.fade-enter-active, .fade-leave-active { transition: opacity 0.5s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
