<template>
  <div class="location-detail glass-panel" v-if="location">
    <div class="detail-header">
      <h4>{{ location.name }}</h4>
      <button class="close-btn" @click="$emit('close')">×</button>
    </div>
    <div class="detail-content">
      <div class="detail-meta">
        <span class="meta-item">加入时间: {{ joinedTimeText }}</span>
        <span class="meta-item">类型: 社团</span>
      </div>
      <button class="action-btn compact" @click="$emit('open', location)">
        查看详情
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  location: {
    type: Object,
    default: null
  },
});

defineEmits(["close", "open"]);

// 格式化 joined_at / joinTime
const joinedTimeText = computed(() => {
  const t = props.location?.joinTime;
  if (!t) return "暂无";

  // 你后端现在会给 "0001-01-01T00:00:00Z" 这种默认值
  if (typeof t === "string" && t.startsWith("0001-01-01")) return "暂无";

  const d = new Date(t);
  if (Number.isNaN(d.getTime())) return String(t);

  // 输出 YYYY-MM-DD（你想要带时间我也可以给你改）
  const pad = (n) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;
});
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
