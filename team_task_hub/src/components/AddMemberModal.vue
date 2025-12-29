<template>
  <div v-if="visible" class="modal-mask" @click.self="$emit('close')">
    <div class="modal">
      <div class="modal-title">添加成员</div>

      <div class="form">
        <label class="label">用户 ID（或邮箱）</label>
        <input class="input" v-model="localValue" placeholder="例如：123 或 xxx@xx.com" />
        <div class="muted" style="margin-top:8px;">
          后端支持哪种就用哪种：userId 或 email
        </div>
      </div>

      <div class="modal-actions">
        <button class="btn ghost" @click="$emit('close')">取消</button>
        <button class="btn" @click="$emit('confirm', localValue)">确认添加</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from "vue";

const props = defineProps({
  visible: { type: Boolean, default: false },
  value: { type: String, default: "" },
});

defineEmits(["close", "confirm"]);

const localValue = ref(props.value);

watch(() => props.visible, (v) => {
  if (v) localValue.value = props.value || "";
});
</script>

<style scoped>
.muted{ opacity:.7; }
.modal-mask{
  position:fixed; inset:0; background:rgba(0,0,0,.55);
  display:grid; place-items:center; z-index:2000;
}
.modal{
  width:min(520px, 92vw);
  border-radius:18px;
  background:rgba(20,24,30,.98);
  border:1px solid rgba(255,255,255,.12);
  padding:14px;
}
.modal-title{ font-weight:900; letter-spacing:.08em; margin-bottom:10px; }
.label{ display:block; font-size:12px; opacity:.75; margin-bottom:6px; }
.input{
  width:100%;
  padding:10px 12px;
  border-radius:12px;
  border:1px solid rgba(255,255,255,.14);
  background:rgba(255,255,255,.06);
  color:#fff; outline:none; box-sizing:border-box;
}
.modal-actions{ display:flex; justify-content:flex-end; gap:10px; margin-top:12px; }
.btn{
  padding:10px 12px; border-radius:12px;
  border:1px solid rgba(255,255,255,.14);
  background:rgba(255,255,255,.10); color:#fff; cursor:pointer;
}
.btn.ghost{ background:transparent; }
</style>
