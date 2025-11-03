<template>
  <div class="home-container">
    <h1 class="title">欢迎使用本系统</h1>
    <p class="subtitle">这是主页内容区域，你可以在这里加入输入框、按钮等功能区。</p>
  </div>
</template>

<script setup>
// 暂时不需要逻辑
</script>

<style scoped>
.home-container {
  margin-top: 120px; /* 防止内容贴住 Header */
  text-align: center;
  color: white;
  padding: 20px;
}

.title {
  font-size: 32px;
  font-weight: bold;
  margin-bottom: 10px;
}

.subtitle {
  font-size: 18px;
  opacity: 0.9;
}
</style>
