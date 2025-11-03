<template>
  <header class="header-bar">
    <div class="logo">Logo</div>

    <div class="right-section">
      <button class="header-btn" @click="handleOrgClick">我的组织</button>

      <button class="header-btn login-btn" @click="handleLoginClick">
        <transition name="fade">
          <span key="login" v-if="!isLoggedIn">登录</span>
          <span key="hello" v-else>你好，{{ username }}</span>
        </transition>
      </button>
    </div>
  </header>
</template>

<script setup>
import { ref } from "vue";

const isLoggedIn = ref(false);
const username = ref("小明");

function handleLoginClick() {
  if (!isLoggedIn.value) {
    isLoggedIn.value = true;
  } else {
    alert("个人信息 功能暂未上线");
  }
}

function handleOrgClick() {
  alert("我的组织 功能暂未上线");
}
</script>

<style scoped>
.header-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  z-index: 1000;

  /* 纯色背景（深蓝色）*/
  background-color: rgba(0, 43, 92, 0.95);
  backdrop-filter: blur(10px); /* 可保留玻璃感模糊 */

  padding: 15px 30px; /* Header 左右间距 */
  color: white;
  font-size: 18px;
  transition: background 0.3s ease;
}

/* Logo 字体斜体 */
.logo {
  font-style: italic;
  font-weight: bold;
  letter-spacing: 1px;
}

/* 右侧按钮水平排列 */
.right-section {
  display: flex;
  gap: 20px; /* 按钮间距 */
  margin-right: 45px; /* 整体左移 */
}

/* 顶部按钮样式：去掉边框*/
.header-btn {
  background: transparent;
  border: none;
  border-radius: 6px;
  color: white;
  padding: 6px 14px;
  cursor: pointer;
  transition: all 0.3s ease;
  /*font-family: "KaiTi", "楷体", serif; 按钮字体改为楷体 */
}

/* hover 动画 */
.header-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: scale(1.05);
}

/* 登录按钮淡入淡出动画 */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
