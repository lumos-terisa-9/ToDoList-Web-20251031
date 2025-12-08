<template>
  <canvas ref="canvas" class="fog-layer"></canvas>
</template>

<script setup>
import {ref, onMounted, onBeforeUnmount, watch} from "vue";
import L from "leaflet";

const props = defineProps({
  map: {type: Object, required: true},
  imgWidth: {type: Number, required: true},
  imgHeight: {type: Number, required: true},
  discoveredAreas: {type: Array, default: () => []},
  playerPosition: {type: Object, default: null},
  playerRadius: {type: Number, default: 12},
});

// ---------------------------
// Canvas 准备
// ---------------------------
const canvas = ref(null);
let ctx = null;

function resizeCanvas() {
  if (!canvas.value) return;
  const size = props.map.getSize();
  canvas.value.width = size.x;
  canvas.value.height = size.y;
}

// 百分比 → 世界坐标 → 屏幕像素
function worldToScreen({x, y}) {
  const lat = (y / 100) * props.imgHeight;
  const lng = (x / 100) * props.imgWidth;
  const point = props.map.latLngToContainerPoint([lat, lng]);
  return point;
}

// ---------------------------
// 🌥️ 生成白色云纹理（核心）
// ---------------------------
// 🌥️ 白色云团层：一堆柔光白色圆叠加，形成云纹理
function drawCloudLayer() {
  const w = canvas.value.width;
  const h = canvas.value.height;

  // 先来一层很淡的白纱
  ctx.fillStyle = "rgba(255,255,255,0.12)";
  ctx.fillRect(0, 0, w, h);

  const cloudCount = 60; // 云团数量，可以按喜好调多调少

  for (let i = 0; i < cloudCount; i++) {
    const cx = Math.random() * w;
    const cy = Math.random() * h;

    // 云团半径：范围内随机
    const minR = Math.min(w, h) * 0.12;
    const maxR = Math.min(w, h) * 0.28;
    const r = minR + Math.random() * (maxR - minR);

    const gradient = ctx.createRadialGradient(cx, cy, 0, cx, cy, r);
    gradient.addColorStop(0, "rgba(255,255,255,0.85)");
    gradient.addColorStop(0.4, "rgba(255,255,255,0.55)");
    gradient.addColorStop(1, "rgba(255,255,255,0)");

    ctx.fillStyle = gradient;
    ctx.beginPath();
    ctx.arc(cx, cy, r, 0, Math.PI * 2);
    ctx.fill();
  }
}

// ---------------------------
// 🌟 柔光擦出视野圈
// ---------------------------
function eraseCircle(cx, cy, r) {
  const gradient = ctx.createRadialGradient(cx, cy, 0, cx, cy, r);
  gradient.addColorStop(0, "rgba(255,255,255,1)");
  gradient.addColorStop(1, "rgba(255,255,255,0)");

  ctx.save();
  ctx.globalCompositeOperation = "destination-out"; // 用擦除方式
  ctx.fillStyle = gradient;
  ctx.beginPath();
  ctx.arc(cx, cy, r, 0, Math.PI * 2);
  ctx.fill();
  ctx.restore();
}

// ---------------------------
// 🔄 主渲染
// ---------------------------
function updateFog() {
  if (!canvas.value) return;
  resizeCanvas();

  const { width, height } = canvas.value;

  // 1. 清空
  ctx.clearRect(0, 0, width, height);

  // 2. 画一整层白色云雾（有明显云团）
  drawCloudLayer();

  // 3. 中心挖一个大圈，让画面中心更清晰一点（可以不要，看你喜好）
  const centerX = width / 2;
  const centerY = height / 2;
  eraseCircle(centerX, centerY, width * 0.3);

  // 4. 永久已探索区域：每个区域擦一圈
  props.discoveredAreas.forEach((area) => {
    const pt = worldToScreen(area);
    const r = (area.radius / 100) * width;
    eraseCircle(pt.x, pt.y, r);
  });

  // 5. 玩家实时视野
  if (props.playerPosition) {
    const pt = worldToScreen(props.playerPosition);
    const r = (props.playerRadius / 100) * width;
    eraseCircle(pt.x, pt.y, r);
  }
}

// ---------------------------
// 事件监听
// ---------------------------
let moveHandler = null;
let zoomHandler = null;

onMounted(() => {
  ctx = canvas.value.getContext("2d");
  resizeCanvas();
  updateFog();

  moveHandler = () => updateFog();
  zoomHandler = () => updateFog();

  props.map.on("move", moveHandler);
  props.map.on("zoom", zoomHandler);

  window.addEventListener("resize", updateFog);
});

watch(() => props.discoveredAreas, updateFog, {deep: true});
watch(() => props.playerPosition, updateFog);

onBeforeUnmount(() => {
  props.map.off("move", moveHandler);
  props.map.off("zoom", zoomHandler);
  window.removeEventListener("resize", updateFog);
});
</script>

<style scoped>
.fog-layer {
  position: absolute;
  inset: 0;
  pointer-events: none;
  /* 全局 z-index 交给 App.vue 控制 */
}
</style>
