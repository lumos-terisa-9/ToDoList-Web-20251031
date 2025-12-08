<template>
  <div class="map-controls glass-panel">
    <div class="control-group">
      <button @click="$emit('reset-view')" class="control-btn" title="重置视图">
        <span class="icon">🗺️</span>
      </button>
      <button @click="$emit('zoom-in')" class="control-btn" title="放大">
        <span class="icon">🔍</span>
      </button>
      <button @click="$emit('zoom-out')" class="control-btn" title="缩小">
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
        @click="$emit('select-location', loc)"
        :class="{ active: activeLocationId === loc.id }"
      >
        <span
          class="location-dot"
          :style="{ backgroundColor: getLocationColor(loc.type) }"
        ></span>
        <span class="location-name">{{ loc.name }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  locations: {
    type: Array,
    required: true
  },
  activeLocationId: {
    type: [Number, String, null],
    default: null
  },
  getLocationColor: {
    type: Function,
    required: true
  }
});

defineEmits(['reset-view', 'zoom-in', 'zoom-out', 'select-location']);
</script>

<!-- 这里可以先不写样式，复用 MapPage 里原来的样式 -->
