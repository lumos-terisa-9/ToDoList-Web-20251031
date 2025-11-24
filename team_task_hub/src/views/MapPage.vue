<template>
  <div class="map-container">
    <!-- 地标点位 -->
    <div
      v-for="loc in locations"
      :key="loc.id"
      class="map-point"
      :style="{ top: loc.y + '%', left: loc.x + '%' }"
      @click="openLocation(loc)"
    >
      <img class="icon" :src="loc.icon" />
      <span class="label">{{ loc.name }}</span>
    </div>
  </div>
</template>

<script setup>
import mapImg from '@/assets/gameMap.jpeg'

// 图标必须 import 才能被 Vite 识别
import demaciaIcon from '@/assets/mapIcon.png'
import noxusIcon from '@/assets/mapIcon.png'
import ioniaIcon from '@/assets/mapIcon.png'
import piltoverIcon from '@/assets/mapIcon.png'

// 地图点位（百分比定位）
//name部分需要通过连接数据库把组织名字换上去
const locations = [
  { id: 1, name: '德玛西亚', x: 22, y: 48, icon: demaciaIcon },
  { id: 2, name: '诺克萨斯', x: 45, y: 34, icon: noxusIcon },
  { id: 3, name: '艾欧尼亚', x: 78, y: 26, icon: ioniaIcon },
  { id: 4, name: '皮尔特沃夫', x: 52, y: 58, icon: piltoverIcon }
]

//之后更新逻辑：点击后进入每个组织的详情界面
function openLocation(loc) {
  alert('你点击了：' + loc.name)
}
</script>

<style scoped>
.map-container {
  position: relative;
  width: 100%;
  height: 100vh; /* 让地图真正铺满整个屏幕 */
  background-image: url('@/assets/gameMap.jpeg');
  background-size: cover;
  background-position: center;
  overflow: hidden;
}

/* 地标点位 */
.map-point {
  position: absolute;
  transform: translate(-50%, -50%);
  cursor: pointer;
  text-align: center;
  transition: transform .2s;
  z-index: 10;
}

.map-point:hover {
  transform: translate(-50%, -50%) scale(1.15);
}

.icon {
  width: 52px;
  height: 52px;
}

.label {
  color: white;
  text-shadow: 0 0 5px black;
  display: block;
  margin-top: 4px;
  font-size: 16px;
}
</style>
