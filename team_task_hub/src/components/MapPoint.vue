<template>
  <div class="point" :style="style" @click="$emit('click')">
    <img class="icon" :src="icon" />
    <div class="label">{{ name }}</div>
  </div>
</template>

<script setup>
const props = defineProps({
  x: Number,
  y: Number,
  icon: String,
  name: String
})

const style = {
  position: "absolute",
  left: props.x + "%",
  top: props.y + "%",
  transform: "translate(-50%, -50%)"
}
</script>

<style scoped>
.point {
  text-align: center;
  cursor: pointer;
  transition: transform .25s ease;
}

.point:hover {
  transform: translate(-50%, -50%) scale(1.25);
}

.icon {
  width: 60px;
  filter: drop-shadow(0 0 8px #fff5);
}

.label {
  margin-top: 4px;
  color: white;
  text-shadow: 0 0 6px black;
  font-size: 16px;
}
</style>
