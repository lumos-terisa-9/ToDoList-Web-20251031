<template>
  <div class="location-detail glass-panel" v-if="location">
    <div class="detail-header">
      <h4>{{ location.name }}</h4>
      <button class="close-btn" @click="$emit('close')">×</button>
    </div>
    <div class="detail-content">
      <div class="detail-meta">
        <span class="meta-item">坐标: {{ location.x }}, {{ location.y }}</span>
        <span class="meta-item">类型: {{ typeName }}</span>
      </div>
      <button class="action-btn compact" @click="$emit('open', location)">
        查看详情
      </button>
    </div>
  </div>
</template>

<script setup>
import {computed} from "vue";   // ✅ 必须要有

const props = defineProps({
  location: {
    type: Object,
    default: null
  },
  getTypeName: {
    type: Function,
    required: true
  }
});

const typeName = computed(() =>
  props.location ? props.getTypeName(props.location.type) : ""
);

defineEmits(["close", "open"]);
</script>

<style scoped>
.location-detail {
  position: absolute;
  top: 70px;
  right: 16px;
  z-index: 1000; /* ✅ 一定要比迷雾层高 */
  padding: 12px;
  min-width: 180px;
  max-width: 220px;
  color: white;
}

/* 下面这些可以直接照你原来 MapPage 里的搬过来 */

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
</style>
