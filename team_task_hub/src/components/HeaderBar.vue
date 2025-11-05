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

  /* === 【修改点】增强的毛玻璃效果 === */

  /* 1. 使用非常透明的浅色背景，让后面的深色背景透出来 */
  background-color: rgba(255, 255, 255, 0.15); /* 透明度为 15% 的白色 */

  /* 2. 略微增加模糊度 */
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px); /* 兼容 Safari */

  /* 3. 添加一个细小的底部分界线，增加玻璃的质感和悬浮感 */
  border-bottom: 1px solid rgba(255, 255, 255, 0.3);

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
