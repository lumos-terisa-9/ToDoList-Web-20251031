<template>
  <div class="wrap">
    <input
        class="search"
        :value="modelValue"
        @input="$emit('update:modelValue', $event.target.value)"
        placeholder="搜索：姓名 / 邮箱"
        @keyup.enter="$emit('search')"
    />

    <button class="btn" @click="$emit('search')" :disabled="loading">
      {{ loading ? "搜索中…" : "搜索" }}
    </button>

    <button class="btn ghost" @click="$emit('show-all')" :disabled="loading">
      全部
    </button>
  </div>
</template>

<script setup>
defineProps({
  modelValue: { type: String, default: "" },
  loading: { type: Boolean, default: false },
});
defineEmits(["update:modelValue", "search", "show-all"]);
</script>

<style scoped>
.wrap{ display:flex; align-items:center; gap:10px; }
.search{
  width: 320px; max-width: 40vw; padding: 10px 14px;
  border-radius: 14px; border: 1px solid rgba(255,255,255,.14);
  background: rgba(255,255,255,.07); color:#fff; outline:none;
}
.btn{
  padding:10px 12px; border-radius:12px;
  border:1px solid rgba(255,255,255,.14);
  background:rgba(255,255,255,.10); color:#fff; cursor:pointer;
  white-space: nowrap;
}
.btn.ghost{ background: transparent; }
</style>
