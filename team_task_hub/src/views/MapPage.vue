<template>
  <div class="map-container">
    <!-- 控制面板 -->
    <MapControls
      :locations="locations"
      :active-location-id="activeLocation?.id || null"
      :get-location-color="getLocationColor"
      @reset-view="resetView"
      @zoom-in="zoomIn"
      @zoom-out="zoomOut"
      @select-location="flyToLocation"
    />

    <div id="map" class="map"></div>

    <!-- 🌫️ 迷雾组件 -->
    <FogLayer
      v-if="map"
      :map="map"
      :img-width="imgWidth"
      :img-height="imgHeight"
      :locations="locations"
    />

    <!-- 详情面板 -->
    <LocationDetail
      v-if="activeLocation"
      :location="activeLocation"
      :get-type-name="getTypeName"
      @close="activeLocation = null"
      @open="openLocation"
    />

    <!-- 组织管理浮动按钮 -->
    <div class="organization-management">
      <!-- 主按钮 -->
      <button
        class="organization-main-btn"
        :class="{ expanded: showOrganizationOptions }"
        @click="toggleOrganizationOptions"
      >
        <span class="btn-icon">{{ showOrganizationOptions ? '×' : '+' }}</span>
      </button>

      <!-- 选项菜单 -->
      <transition name="fade-slide">
        <div v-if="showOrganizationOptions" class="organization-options">
          <div class="option-item create-option" @click="openCreateOrganizationModal">
            <div class="option-icon">🏢</div>
            <div class="option-text">创建组织</div>
          </div>
          <div class="option-item join-option" @click="openJoinOrganizationModal">
            <div class="option-icon">👥</div>
            <div class="option-text">申请加入</div>
          </div>
          <div class="option-item applications-option" @click="openViewApplicationsModal">
            <div class="option-icon">📋</div>
            <div class="option-text">查看申请</div>
          </div>
        </div>
      </transition>
    </div>

    <!-- 创建组织模态框 -->
    <CreateOrganizationModal
      :is-visible="showCreateOrganizationModal"
      @close="showCreateOrganizationModal = false"
      @created="handleOrganizationCreated"
    />

    <!-- 申请加入组织模态框 -->
    <JoinOrganizationModal
      :is-visible="showJoinOrganizationModal"
      @close="showJoinOrganizationModal = false"
      @joined="handleOrganizationJoined"
    />

    <!-- 查看申请模态框 -->
    <ViewApplicationsModal
      :is-visible="showViewApplicationsModal"
      @close="showViewApplicationsModal = false"
    />
  </div>
</template>

<script setup>
import axios from "axios";
import { useRouter } from "vue-router";
import L from "leaflet";
import { onMounted, ref, onUnmounted, onBeforeUnmount } from "vue";
import "leaflet/dist/leaflet.css";
import MapControls from "@/components/MapControls.vue";
import LocationDetail from "@/components/LocationDetail.vue";
import FogLayer from "@/components/FogLayer.vue";
import CreateOrganizationModal from "@/components/CreateOrganizationModal.vue";
import JoinOrganizationModal from "@/components/JoinOrganizationModal.vue";
import ViewApplicationsModal from "@/components/ViewApplicationsModal.vue";

// 地图资源
import mapImg from "@/assets/gameMap.jpeg";

// 图标
import demaciaIcon from "@/assets/mapIcon.png";
import noxusIcon from "@/assets/mapIcon.png";
import ioniaIcon from "@/assets/mapIcon.png";
import piltoverIcon from "@/assets/mapIcon.png";

// 暴露给父组件的方法
defineExpose({
  isFullScreenPage: true
})

// ----------------------------
// 响应式数据
// ----------------------------
const router = useRouter();
const map = ref(null);
const activeLocation = ref(null);
const showOrganizationOptions = ref(false);
const showCreateOrganizationModal = ref(false);
const showJoinOrganizationModal = ref(false);
const showViewApplicationsModal = ref(false);

const apiBaseUrl = "http://localhost:8080"; // 你的后端地址
// 建一个 axios 实例，统一配置 baseURL 和 token
const api = axios.create({
  baseURL: apiBaseUrl,
  timeout: 10000,
});

// 每次请求自动带上 token（根据后端的校验方式调整）
api.interceptors.request.use((config) => {
  const raw = localStorage.getItem("token");
  if (raw) {
    try {
      const obj = JSON.parse(raw);
      const accessToken = obj?.data?.access_token;
      if (accessToken) {
        config.headers.Authorization = `Bearer ${accessToken}`;
      }
    } catch (e) {
      // 如果 raw 本身就是 token 字符串（兼容老版本）
      config.headers.Authorization = `Bearer ${raw}`;
    }
  }
  return config;
});

// 地图图片尺寸
const imgWidth = 6000;
const imgHeight = 3374;

//初始缩放尺寸
const initialZoom = ref(null);

// ----------------------------
// 点位数据（只负责坐标/type，名字和加入时间从后端填）
// ----------------------------
const locations = ref([
  { id: 1, x: 22, y: 48, icon: demaciaIcon, type: "kingdom", name: "", joinTime: null },
  { id: 2, x: 45, y: 34, icon: noxusIcon,   type: "empire",  name: "", joinTime: null },
  { id: 3, x: 78, y: 26, icon: ioniaIcon,   type: "region",  name: "", joinTime: null },
  { id: 4, x: 52, y: 58, icon: piltoverIcon,type: "city",    name: "", joinTime: null },
]);

// ----------------------------
// 辅助函数
// ----------------------------

// 从后端加载所有"可在地图展示的组织"，名字&加入时间都来自数据库
async function loadOrgInfoFromBackend() {
  try {
    const resp = await api.get("/api/organization/my-organizations");
    const orgList = resp.data?.data;

    if (!Array.isArray(orgList)) return;

    locations.value = locations.value.map((loc, idx) => {
      const org = orgList[idx]; // 按顺序塞进点位
      if (!org) return loc;

      return {
        ...loc,
        name: org.org_name || loc.name,
        joinTime: org.joined_at || null,
        logoUrl: org.logo_url || null,
        creatorId: org.creator_id ?? null,
      };
    });

    console.log("加载后的地标：", locations.value);
  } catch (err) {
    console.error("加载组织信息失败：", err);
  }
}

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
// 组织管理相关函数
// ----------------------------
function toggleOrganizationOptions() {
  showOrganizationOptions.value = !showOrganizationOptions.value;
}

function openCreateOrganizationModal() {
  showOrganizationOptions.value = false;
  showCreateOrganizationModal.value = true;
}

function openJoinOrganizationModal() {
  showOrganizationOptions.value = false;
  showJoinOrganizationModal.value = true;
}

function openViewApplicationsModal() {
  showOrganizationOptions.value = false;
  showViewApplicationsModal.value = true;
}

function handleOrganizationCreated() {
  showCreateOrganizationModal.value = false;
  // 可以在这里添加创建成功后的处理逻辑
}

function handleOrganizationJoined() {
  showJoinOrganizationModal.value = false;
  // 可以在这里添加加入成功后的处理逻辑
  // 如果查看申请模态框是打开的，可以刷新申请列表
}

// ----------------------------
// 地图操作函数
// ----------------------------
function openLocation(loc) {
  router.push({
    name: "Org",              // 路由 name（你路由里要配成 Org）
    params: { id: loc.id },   // /org/:id
    query: {                  // 可选：带一些展示信息
      name: loc.name || "",
      joinTime: loc.joinTime || "",
    },
  });
}

function flyToLocation(loc) {
  if (!map.value) return;

  const [py, px] = percentToPx(loc);
  const targetLatLng = L.latLng(py, px);

  if (initialZoom.value === null) {
    initialZoom.value = map.value.getZoom();
  }

  let zoom = initialZoom.value;
  const maxZoomAllowed = Math.min(
    map.value.getMaxZoom(),
    initialZoom.value + 0.8
  );

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

    return imageBounds.contains(tlLatLng) && imageBounds.contains(brLatLng);
  }

  while (zoom < maxZoomAllowed && !canCenterAt(zoom)) {
    zoom += 0.25;
  }

  const finalZoom = zoom;

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
onMounted(async () => {
  await loadOrgInfoFromBackend();

  try {
    const bounds = [
      [0, 0],
      [imgHeight, imgWidth]
    ];

    const paddedBounds = L.latLngBounds(
      [-50, -50],
      [imgHeight + 50, imgWidth + 50]
    );

    map.value = L.map("map", {
      crs: L.CRS.Simple,
      minZoom: -2.5,
      maxZoom: 20,
      zoomControl: false,
      attributionControl: false,
      preferCanvas: true,
      fadeAnimation: false,
      zoomAnimation: false,
      maxBounds: paddedBounds,
      maxBoundsViscosity: 1.0
    });

    const imageLayer = L.imageOverlay(mapImg, bounds, {
      className: 'map-image-layer'
    }).addTo(map.value);

    L.control.zoom({ position: 'bottomright'}).addTo(map.value);

    map.value.fitBounds(bounds);

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

<!-- 添加组织管理相关样式 -->
<style>
/* 组织管理浮动按钮 */
.organization-management {
  position: absolute;
  bottom: 100px;
  right: 30px;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.organization-main-btn {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  color: white;
  font-size: 28px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 32px rgba(102, 126, 234, 0.4);
  transition: all 0.3s ease;
  position: relative;
  z-index: 1001;
}

.organization-main-btn:hover {
  transform: scale(1.1);
  box-shadow: 0 12px 40px rgba(102, 126, 234, 0.6);
}

.organization-main-btn.expanded {
  transform: rotate(45deg);
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.organization-main-btn .btn-icon {
  font-weight: 300;
  transition: transform 0.3s ease;
}

/* 选项菜单 */
.organization-options {
  margin-top: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(12px) saturate(180%);
  -webkit-backdrop-filter: blur(12px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  padding: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  min-width: 180px;
}

.option-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: white;
}

.option-item:hover {
  transform: translateX(-4px);
  background: rgba(255, 255, 255, 0.1);
}

.option-item.create-option:hover {
  background: rgba(102, 126, 234, 0.3);
}

.option-item.join-option:hover {
  background: rgba(66, 153, 225, 0.3);
}

.option-item.applications-option:hover {
  background: rgba(245, 158, 11, 0.3); /* 橙色系 */
}

.option-icon {
  font-size: 20px;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
}

.option-text {
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
}

/* 动画效果 */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.3s ease;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .organization-management {
    bottom: 80px;
    right: 20px;
  }

  .organization-main-btn {
    width: 56px;
    height: 56px;
    font-size: 24px;
  }

  .organization-options {
    min-width: 160px;
  }
}

@media (max-width: 480px) {
  .organization-management {
    bottom: 70px;
    right: 16px;
  }

  .organization-main-btn {
    width: 52px;
    height: 52px;
    font-size: 22px;
  }
}
</style>

<!-- 原有的样式保持不变 -->
<style>
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
  top: 70px;
  left: 16px;
  z-index: 1000;
  padding: 12px;
  min-width: 160px;
  max-width: 200px;
  max-height: calc(100vh - 90px);
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
  top: 70px;
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
    top: 60px;
    min-width: 140px;
    max-width: 160px;
    padding: 8px;
  }

  .location-detail {
    top: 60px;
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
    top: 55px;
    left: 8px;
    min-width: 120px;
    max-width: 140px;
  }

  .location-detail {
    top: 55px;
    right: 8px;
    min-width: 140px;
    max-width: 180px;
  }
}

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
