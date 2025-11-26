<template>
  <div class="map-container">
    <!-- 紧凑的玻璃质感控制面板 -->
    <div class="map-controls glass-panel">
      <div class="control-group">
        <button @click="resetView" class="control-btn" title="重置视图">
          <span class="icon">🗺️</span>
        </button>
        <button @click="zoomIn" class="control-btn" title="放大">
          <span class="icon">🔍</span>
        </button>
        <button @click="zoomOut" class="control-btn" title="缩小">
          <span class="icon">🔎</span>
        </button>
      </div>
      <div class="location-list">
        <div class="list-header">
          <span>组织列表</span>
          <span class="count-badge">{{ locations.length }}</span>
        </div>
        <div
          v-for="loc in locations"
          :key="loc.id"
          class="location-item"
          @click="flyToLocation(loc)"
          :class="{ active: activeLocation?.id === loc.id }"
        >
          <span class="location-dot" :style="{ backgroundColor: getLocationColor(loc.type) }"></span>
          <span class="location-name">{{ loc.name }}</span>
        </div>
      </div>
    </div>

    <!-- 全屏地图容器 -->
    <div id="map" class="map"></div>

    <!-- 紧凑的详情面板 -->
    <div v-if="activeLocation" class="location-detail glass-panel">
      <div class="detail-header">
        <h4>{{ activeLocation.name }}</h4>
        <button class="close-btn" @click="activeLocation = null">×</button>
      </div>
      <div class="detail-content">
        <div class="detail-meta">
          <span class="meta-item">坐标: {{ activeLocation.x }}, {{ activeLocation.y }}</span>
          <span class="meta-item">类型: {{ getTypeName(activeLocation.type) }}</span>
        </div>
        <button class="action-btn compact" @click="openLocation(activeLocation)">
          查看详情
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import L from "leaflet";
import { onMounted, ref, onUnmounted, onBeforeUnmount } from "vue";
import "leaflet/dist/leaflet.css";

// 地图资源
import mapImg from "@/assets/gameMap.jpeg";

// 图标
import demaciaIcon from "@/assets/mapIcon.png";
import noxusIcon from "@/assets/mapIcon.png";
import ioniaIcon from "@/assets/mapIcon.png";
import piltoverIcon from "@/assets/mapIcon.png";

defineExpose({
  isFullScreenPage: true
})

// ----------------------------
// 响应式数据
// ----------------------------
const map = ref(null);
const activeLocation = ref(null);

// 地图图片尺寸（请根据你的图片实际尺寸修改）
const imgWidth = 6000;
const imgHeight = 3374;

//初始缩放尺寸
const initialZoom = ref(null);

// ----------------------------
// 点位数据
// ----------------------------
const locations = ref([
  {id: 1, name: "德玛西亚", x: 22, y: 48, icon: demaciaIcon, type: "kingdom"},
  {id: 2, name: "诺克萨斯", x: 45, y: 34, icon: noxusIcon, type: "empire"},
  {id: 3, name: "艾欧尼亚", x: 78, y: 26, icon: ioniaIcon, type: "region"},
  {id: 4, name: "皮尔特沃夫", x: 52, y: 58, icon: piltoverIcon, type: "city"},
]);

// ----------------------------
// 辅助函数
// ----------------------------
function percentToPx(loc) {
  return [(loc.y / 100) * imgHeight, (loc.x / 100) * imgWidth];
}

function getLocationColor(type) {
  const colors = {
    kingdom: '#ff6b6b',
    empire: '#4ecdc4',
    region: '#455eb7',
    city: '#ffa500'
  };
  return colors[type] || '#666';
}

function getTypeName(type) {
  const names = {
    kingdom: '王国',
    empire: '帝国',
    region: '地区',
    city: '城市'
  };
  return names[type] || '未知';
}

// ----------------------------
// 地图操作函数
// ----------------------------
function openLocation(loc) {
  alert(`你点击了：${loc.name}\nID: ${loc.id}\n坐标: (${loc.x}, ${loc.y})`);
}

function flyToLocation(loc) {
  if (!map.value) return;

  const [py, px] = percentToPx(loc);
  const targetLatLng = L.latLng(py, px);

  // 第一次调用时记录“初始缩放等级”（一般就是 fitBounds 之后的 zoom）
  if (initialZoom.value === null) {
    initialZoom.value = map.value.getZoom();
  }

  // 每次点击都从同一个 baseZoom 开始重新计算，不再沿用当前 zoom
  let zoom = initialZoom.value;

  // 最大允许缩放：在初始基础上只允许略微放大，避免放大过头
  const maxZoomAllowed = Math.min(
    map.value.getMaxZoom(),
    initialZoom.value + 0.8      // ★ 想再小可以改成 0.5
  );

  // 用“整张图片”的边界来判断是否会出界，而不是当前视图的 getBounds()
  const imageBounds = L.latLngBounds(
    [0, 0],
    [imgHeight, imgWidth]
  );

  function canCenterAt(z) {
    const size = map.value.getSize();
    const halfW = size.x / 2;
    const halfH = size.y / 2;

    const proj = map.value.project(targetLatLng, z);
    const topLeft = proj.subtract([halfW, halfH]);
    const bottomRight = proj.add([halfW, halfH]);

    const tlLatLng = map.value.unproject(topLeft, z);
    const brLatLng = map.value.unproject(bottomRight, z);

    // 只要视图四角都还在整张图片范围内，就认为可以在这个 zoom 居中
    return imageBounds.contains(tlLatLng) && imageBounds.contains(brLatLng);
  }

  // 从初始 zoom 开始，能不放大就不放大；不够的话再一点点放大
  while (zoom < maxZoomAllowed && !canCenterAt(zoom)) {
    zoom += 0.25; // 小步放大，避免一下子 zoom 很大
  }

  const finalZoom = zoom;

  // 使用计算好的 finalZoom，直接飞到目标点
  map.value.flyTo(targetLatLng, finalZoom, {
    duration: 0.8
  });

  activeLocation.value = loc;
}


function resetView() {
  if (!map.value) return;
  const bounds = [[0, 0], [imgHeight, imgWidth]];
  map.value.fitBounds(bounds);
  activeLocation.value = null;
}

function zoomIn() {
  if (!map.value) return;
  map.value.zoomIn();
}

function zoomOut() {
  if (!map.value) return;
  map.value.zoomOut();
}

// 窗口大小变化时重新调整地图
function handleResize() {
  if (map.value) {
    setTimeout(() => {
      map.value.invalidateSize();
    }, 100);
  }
}

// ----------------------------
// 初始化地图
// ----------------------------
onMounted(() => {
  try {
    // 原始边界
    const bounds = [
      [0, 0],
      [imgHeight, imgWidth]
    ];

    // 扩展边界，避免缩放时露出黑边
    const paddedBounds = L.latLngBounds(
      [-50, -50],
      [imgHeight + 50, imgWidth + 50]
    );

    map.value = L.map("map", {
      crs: L.CRS.Simple,
      minZoom: -2.5,
      maxZoom: 4,
      zoomControl: false,
      attributionControl: false,
      preferCanvas: true,
      fadeAnimation: false,
      zoomAnimation: false,
      maxBounds: paddedBounds,        // ★ 改动点
      maxBoundsViscosity: 1.0
    });

    // 图层
    const imageLayer = L.imageOverlay(mapImg, bounds, {
      className: 'map-image-layer'
    }).addTo(map.value);

    L.control.zoom({ position: 'bottomright'}).addTo(map.value);

    map.value.fitBounds(bounds);

    // 你的 marker 逻辑保持不变
    locations.value.forEach((loc) => {
      const [py, px] = percentToPx(loc);
      const icon = L.divIcon({
        html: `
          <div class="custom-marker" data-location-id="${loc.id}">
            <div class="marker-pin" style="background-color: ${getLocationColor(loc.type)}">
              <img src="${loc.icon}" alt="${loc.name}" />
            </div>
            <div class="marker-label">${loc.name}</div>
          </div>
        `,
        className: 'custom-div-icon',
        iconSize: [36, 36],
        iconAnchor: [18, 36],
      });

      const marker = L.marker([py, px], { icon })
        .addTo(map.value)
        .on("click", () => {
          activeLocation.value = loc;
        });

      marker.on('mouseover', () => {
        marker.getElement().classList.add('marker-hover');
      });
      marker.on('mouseout', () => {
        marker.getElement().classList.remove('marker-hover');
      });
    });

    window.addEventListener('resize', handleResize);

  } catch (error) {
    console.error('地图初始化失败:', error);
  }
});


// 清理资源
onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize);
});

onUnmounted(() => {
  if (map.value) {
    map.value.remove();
    map.value = null;
  }
});
</script>

<style scoped>
/* 重置样式确保全屏 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.map-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  background: #0f1419;
}

.map {
  width: 100%;
  height: 100%;
  background: #0f1419;
}

/* 玻璃质感面板基础样式 */
.glass-panel {
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(12px) saturate(180%);
  -webkit-backdrop-filter: blur(12px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  box-shadow:
    0 8px 32px rgba(0, 0, 0, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

/* 紧凑的地图控制面板 - 调整位置避免被HeaderBar遮挡 */
.map-controls {
  position: absolute;
  top: 70px; /* HeaderBar高度约48px + 额外间距 */
  left: 16px;
  z-index: 1000;
  padding: 12px;
  min-width: 160px;
  max-width: 200px;
  max-height: calc(100vh - 90px); /* 限制高度，避免超出屏幕 */
  overflow-y: auto;
}

.control-group {
  display: flex;
  gap: 6px;
  margin-bottom: 12px;
  justify-content: center;
}

.control-btn {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.2);
  color: white;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.control-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: scale(1.1);
}

.control-btn .icon {
  font-size: 14px;
  filter: brightness(0) invert(1);
}

/* 组织列表 */
.location-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 8px;
  color: white;
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 4px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.count-badge {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  padding: 2px 6px;
  border-radius: 8px;
  font-size: 10px;
}

.location-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: white;
  font-size: 12px;
}

.location-item:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateX(2px);
}

.location-item.active {
  background: rgba(78, 205, 196, 0.3);
  border-left: 2px solid #4ecdc4;
}

.location-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.location-name {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 紧凑的详情面板 - 同样调整位置 */
.location-detail {
  position: absolute;
  top: 70px; /* HeaderBar高度约48px + 额外间距 */
  right: 16px;
  z-index: 1000;
  padding: 12px;
  min-width: 180px;
  max-width: 220px;
  color: white;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.detail-header h4 {
  margin: 0;
  font-size: 14px;
  color: white;
}

.close-btn {
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.7);
  font-size: 16px;
  cursor: pointer;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.detail-content {
  font-size: 12px;
}

.detail-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-bottom: 10px;
}

.meta-item {
  color: rgba(255, 255, 255, 0.8);
}

.action-btn.compact {
  width: 100%;
  padding: 6px 12px;
  background: rgba(255, 107, 107, 0.8);
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 11px;
  transition: background 0.2s ease;
}

.action-btn.compact:hover {
  background: rgba(255, 82, 82, 0.9);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .map-controls {
    top: 60px; /* 移动端也相应调整 */
    min-width: 140px;
    max-width: 160px;
    padding: 8px;
  }

  .location-detail {
    top: 60px; /* 移动端也相应调整 */
    min-width: 160px;
    max-width: 200px;
  }

  .control-btn {
    width: 28px;
    height: 28px;
  }
}

@media (max-width: 480px) {
  .map-controls {
    top: 55px; /* 更小的屏幕进一步调整 */
    left: 8px;
    min-width: 120px;
    max-width: 140px;
  }

  .location-detail {
    top: 55px; /* 更小的屏幕进一步调整 */
    right: 8px;
    min-width: 140px;
    max-width: 180px;
  }
}
</style>

<style>
/* Leaflet 标记样式优化 */
.custom-div-icon {
  background: none !important;
  border: none !important;
}

.custom-marker {
  text-align: center;
  transition: all 0.3s ease;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.3));
}

.marker-pin {
  width: 32px;
  height: 32px;
  border-radius: 50% 50% 50% 0;
  transform: rotate(-45deg);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.marker-pin img {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  transform: rotate(45deg);
  border: 2px solid white;
}

.marker-label {
  position: absolute;
  top: -20px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(0, 0, 0, 0.8);
  color: white;
  padding: 3px 8px;
  border-radius: 4px;
  font-size: 11px;
  white-space: nowrap;
  opacity: 0;
  transition: opacity 0.3s ease;
  pointer-events: none;
  backdrop-filter: blur(4px);
}

.custom-marker:hover .marker-label {
  opacity: 1;
}

.marker-hover .marker-pin {
  transform: rotate(-45deg) scale(1.1);
}

/* 重要：让 Leaflet 完全覆盖且不冲突 */
.leaflet-container {
  background: #000;
  width: 100%;
  height: 100%;
}

/* 调整 Leaflet 默认控件样式 */
.leaflet-control-zoom {
  border: none !important;
  background: rgba(255, 255, 255, 0.1) !important;
  backdrop-filter: blur(8px);
  border-radius: 8px !important;
  overflow: hidden;
  margin-bottom: 20px !important;
  margin-right: 10px !important;
}

.leaflet-control-zoom a {
  background: rgba(255, 255, 255, 0.2) !important;
  color: white !important;
  border: none !important;
  border-radius: 0 !important;
  width: 30px !important;
  height: 30px !important;
  line-height: 30px !important;
}

.leaflet-control-zoom a:hover {
  background: rgba(255, 255, 255, 0.3) !important;
}

.leaflet-control-zoom a:first-child {
  border-bottom: 1px solid rgba(255, 255, 255, 0.1) !important;
}

/* 隐藏 Leaflet 水印 */
.leaflet-control-attribution {
  display: none !important;
}

/* 确保地图层正确显示 */
.leaflet-map-pane,
.leaflet-tile,
.leaflet-marker-icon,
.leaflet-marker-shadow,
.leaflet-tile-container,
.leaflet-image-layer,
.leaflet-pane > svg path,
.leaflet-tile-pane {
  max-width: none !important;
  max-height: none !important;
}
</style>
