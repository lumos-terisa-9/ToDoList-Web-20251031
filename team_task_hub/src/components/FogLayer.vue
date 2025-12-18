<template>
  <canvas ref="canvas" class="fog-layer"></canvas>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from "vue";
import L from "leaflet";

const props = defineProps({
  map: { type: Object, required: true },
  imgWidth: { type: Number, required: true },
  imgHeight: { type: Number, required: true },
  // 永久已探索区域：[{ x, y, radius }]
  discoveredAreas: { type: Array, default: () => [] },
  // 玩家实时位置：{ x, y }
  playerPosition: { type: Object, default: null },
  // 玩家视野半径（百分比）
  playerRadius: { type: Number, default: 12 },
});

const canvas = ref(null);
let ctx = null;

// 离屏云层：只生成一次，拖动时反复复用
let cloudCanvas = null;
let cloudCtx = null;
let lastWidth = 0;
let lastHeight = 0;

function resizeCanvasIfNeeded() {
  if (!canvas.value || !props.map) return;

  const size = props.map.getSize();
  const w = size.x;
  const h = size.y;

  if (canvas.value.width !== w || canvas.value.height !== h) {
    canvas.value.width = w;
    canvas.value.height = h;
  }
  return { w, h };
}

// 百分比 → 世界坐标 → 屏幕像素
function worldToScreen({ x, y }) {
  const lat = (y / 100) * props.imgHeight;
  const lng = (x / 100) * props.imgWidth;
  const latLng = L.latLng(lat, lng);
  const point = props.map.latLngToContainerPoint(latLng);
  return point; // { x, y }
}

/** 在离屏 canvas 上生成一层白色云团（只在尺寸变化时调用） */
function buildCloudLayer(w, h) {
  cloudCanvas = document.createElement("canvas");
  cloudCanvas.width = w;
  cloudCanvas.height = h;
  cloudCtx = cloudCanvas.getContext("2d");

  // 先来一层很淡的白纱
  cloudCtx.fillStyle = "rgba(255,255,255,0.12)";
  cloudCtx.fillRect(0, 0, w, h);

  const cloudCount = 60; // 云团数量，可以按需要调整

  for (let i = 0; i < cloudCount; i++) {
    const cx = Math.random() * w;
    const cy = Math.random() * h;

    const minR = Math.min(w, h) * 0.12;
    const maxR = Math.min(w, h) * 0.28;
    const r = minR + Math.random() * (maxR - minR);

    const gradient = cloudCtx.createRadialGradient(cx, cy, 0, cx, cy, r);
    gradient.addColorStop(0, "rgba(255,255,255,0.85)");
    gradient.addColorStop(0.4, "rgba(255,255,255,0.55)");
    gradient.addColorStop(1, "rgba(255,255,255,0)");

    cloudCtx.fillStyle = gradient;
    cloudCtx.beginPath();
    cloudCtx.arc(cx, cy, r, 0, Math.PI * 2);
    cloudCtx.fill();
  }

  lastWidth = w;
  lastHeight = h;
}

/** destination-out：柔和挖一个洞 */
function eraseCircle(cx, cy, r) {
  const gradient = ctx.createRadialGradient(cx, cy, 0, cx, cy, r);
  gradient.addColorStop(0, "rgba(255,255,255,1)");
  gradient.addColorStop(1, "rgba(255,255,255,0)");

  ctx.save();
  ctx.globalCompositeOperation = "destination-out";
  ctx.fillStyle = gradient;
  ctx.beginPath();
  ctx.arc(cx, cy, r, 0, Math.PI * 2);
  ctx.fill();
  ctx.restore();
}

// 擦除一个矩形区域，并且四周带柔和过渡
function eraseElementArea(selector, padding = 10, soften = 25) {
  if (!canvas.value || !ctx) return;

  const el = document.querySelector(selector);
  if (!el) return;

  const rect = el.getBoundingClientRect();
  const canvasRect = canvas.value.getBoundingClientRect();

  // 屏幕 → canvas 坐标
  const x = rect.left - canvasRect.left - padding;
  const y = rect.top - canvasRect.top - padding;
  const w = rect.width + padding * 2;
  const h = rect.height + padding * 2;

  // 内部矩形：完全擦除
  ctx.save();
  ctx.globalCompositeOperation = "destination-out";
  ctx.fillStyle = "rgba(255,255,255,1)";
  ctx.fillRect(x, y, w, h);
  ctx.restore();

  // 边缘柔光擦除（让矩形边缘有柔和过渡）
  const edges = soften;

  ctx.save();
  ctx.globalCompositeOperation = "destination-out";

  // 上边缘
  const topGrad = ctx.createLinearGradient(x, y - edges, x, y);
  topGrad.addColorStop(0, "rgba(255,255,255,0)");
  topGrad.addColorStop(1, "rgba(255,255,255,1)");
  ctx.fillStyle = topGrad;
  ctx.fillRect(x, y - edges, w, edges);

  // 下边缘
  const bottomGrad = ctx.createLinearGradient(x, y + h, x, y + h + edges);
  bottomGrad.addColorStop(0, "rgba(255,255,255,1)");
  bottomGrad.addColorStop(1, "rgba(255,255,255,0)");
  ctx.fillStyle = bottomGrad;
  ctx.fillRect(x, y + h, w, edges);

  // 左边缘
  const leftGrad = ctx.createLinearGradient(x - edges, y, x, y);
  leftGrad.addColorStop(0, "rgba(255,255,255,0)");
  leftGrad.addColorStop(1, "rgba(255,255,255,1)");
  ctx.fillStyle = leftGrad;
  ctx.fillRect(x - edges, y, edges, h);

  // 右边缘
  const rightGrad = ctx.createLinearGradient(x + w, y, x + w + edges, y);
  rightGrad.addColorStop(0, "rgba(255,255,255,1)");
  rightGrad.addColorStop(1, "rgba(255,255,255,0)");
  ctx.fillStyle = rightGrad;
  ctx.fillRect(x + w, y, edges, h);

  ctx.restore();
}

/** 主更新函数：不再随机生成云，只复用 cloudCanvas，避免闪烁 */
function updateFog() {
  if (!canvas.value || !props.map) return;
  const size = resizeCanvasIfNeeded();
  if (!size) return;
  const { w, h } = size;

  if (!ctx) {
    ctx = canvas.value.getContext("2d");
  }

  // 如果离屏云层不存在或尺寸变化，重新生成一次
  if (!cloudCanvas || w !== lastWidth || h !== lastHeight) {
    buildCloudLayer(w, h);
  }

  // 1. 把离屏云层画到可见 canvas 上
  ctx.globalCompositeOperation = "source-over";
  ctx.clearRect(0, 0, w, h);
  ctx.drawImage(cloudCanvas, 0, 0);

  // 2. 中心挖一个大洞，让画面中心更亮一些（可选）
  const centerX = w / 2;
  const centerY = h / 2;
  eraseCircle(centerX, centerY, w * 0.3);

  // 3. 永久已探索区域
  props.discoveredAreas.forEach((area) => {
    const pt = worldToScreen(area);
    const r = (area.radius / 100) * w;
    eraseCircle(pt.x, pt.y, r);
  });

  // 4. 玩家实时视野圈
  if (props.playerPosition) {
    const pt = worldToScreen(props.playerPosition);
    const r = (props.playerRadius / 100) * w;
    eraseCircle(pt.x, pt.y, r);
  }

  //挖掉控制面板UI区域
  eraseElementArea(".map-controls", 8);
  eraseElementArea(".location-detail", 8);
  eraseElementArea(".leaflet-control-zoom", 8, 20);
}

// 用 rAF 做一个简单的节流，避免拖动时调用过于频繁
let pending = false;

function scheduleUpdate() {
  if (pending) return;
  pending = true;
  requestAnimationFrame(() => {
    pending = false;
    updateFog();
  });
}


let moveHandler = null;
let zoomHandler = null;

onMounted(() => {
  if (!props.map) return;
  ctx = canvas.value.getContext("2d");

  updateFog();

  // 拖动过程中用 move 事件，但内部用 rAF 节流
  moveHandler = () => scheduleUpdate();
  zoomHandler = () => scheduleUpdate();

  props.map.on("move", moveHandler);
  props.map.on("zoom", zoomHandler);

  window.addEventListener("resize", scheduleUpdate);
});

watch(
  () => props.discoveredAreas,
  () => scheduleUpdate(),
  {deep: true}
);

watch(
  () => props.playerPosition,
  () => scheduleUpdate(),
  {deep: true}
);

onBeforeUnmount(() => {
  if (props.map) {
    props.map.off("move", moveHandler);
    props.map.off("zoom", zoomHandler);
  }
  window.removeEventListener("resize", scheduleUpdate);
});
</script>

<style scoped>
.fog-layer {
  position: absolute;
  inset: 0;
  pointer-events: none;
  /* z-index 交给 App.vue 的全局样式控制 */
}
</style>
