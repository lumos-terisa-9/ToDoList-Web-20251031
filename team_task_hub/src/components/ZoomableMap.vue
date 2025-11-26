<template>
  <div id="map" class="map"></div>
</template>

<script setup>
import L from "leaflet"
import { onMounted } from "vue"

// 地图 URL —— 换成你自己的
const mapUrl = "https://mapgenie.io/wild-hearts/maps/akikure-canyon?x=-0.742058656536102&y=0.7262284077758494&zoom=10.788850297817369"

// 地图原图大小（非常重要）
// ⚠️ 你必须知道地图图片的宽度和高度（像素）
const imgWidth = 4096
const imgHeight = 4096

onMounted(() => {
  // 创建地图（不使用经纬度）
  const map = L.map("map", {
    crs: L.CRS.Simple,     // 使用简单坐标系
    minZoom: -2,
    maxZoom: 4,
    zoomSnap: 0.1,
    zoomDelta: 0.5
  })

  // 设置地图 bounds（左上到右下）
  const bounds = [[0, 0], [imgHeight, imgWidth]]

  // 绑定图片覆盖层
  L.imageOverlay(mapUrl, bounds).addTo(map)

  // 设置初始视图
  map.fitBounds(bounds)

  // 示例：添加一个点位（可换为数据库数据）
  L.marker([2000, 1500]).addTo(map).bindPopup("示例位置")
})
</script>

<style>
.map {
  width: 100%;
  height: 100vh;
}
</style>
