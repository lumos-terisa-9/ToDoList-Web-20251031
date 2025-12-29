<template>
  <div>
    <div v-if="loading" class="muted">加载中…</div>
    <div v-else-if="error" class="muted">加载失败：{{ error }}</div>

    <div v-else class="table">
      <div class="row head">
        <div>成员</div>
        <div>邮箱</div>
        <div>角色</div>
        <div>加入时间</div>
        <div style="text-align:right;">操作</div>
      </div>

      <div class="row" v-for="m in members" :key="m.userId">
        <div class="name">
          {{ m.name }}
          <span v-if="m.userId === myUserId" class="me">（我）</span>
        </div>

        <div class="muted">{{ m.email || "-" }}</div>

        <div>
          <span class="badge" :class="roleClass(m.role)">{{ roleText(m.role) }}</span>
        </div>

        <div class="muted">{{ (m.joinedAt || '').slice(0,10) || "-" }}</div>

        <div class="actions">
          <template v-if="canManage">
            <button
              class="btn small"
              v-if="normRole(m.role) !== 'admin'"
              @click="$emit('set-role', { userId: m.userId, role: 'admin' })"
              :disabled="normRole(m.role) === 'creator'"
            >
              设为管理员
            </button>

            <button
              class="btn small ghost"
              v-else
              @click="$emit('set-role', { userId: m.userId, role: 'member' })"
              :disabled="normRole(m.role) === 'creator'"
            >
              取消管理员
            </button>

            <button
              class="btn small danger"
              @click="$emit('remove', m.userId)"
              :disabled="normRole(m.role) === 'creator' || m.userId === myUserId"
            >
              删除成员
            </button>
          </template>

          <span v-else class="muted">无权限</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  members: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
  error: { type: String, default: "" },
  canManage: { type: Boolean, default: false },
  myUserId: { type: Number, default: 0 },
});
defineEmits(["remove", "set-role"]);

function roleText(role) {
  const r = String(role).toLowerCase();
  if (r === "owner") return "负责人";
  if (r === "admin") return "管理员";
  if (r === "creator") return "创建者";
  if (r === "admin") return "管理员";
  return "成员";
}

function normRole(role) {
  return String(role || "member").trim().toLowerCase();
}

function roleClass(role) {
  const r = normRole(role);
  // 复用你原来的 badge.owner 样式给 Creator
  return r === "creator" ? "owner" : r === "admin" ? "admin" : "member";
}

</script>

<style scoped>
.muted{ opacity:.7; }

.table{ display:flex; flex-direction:column; gap:8px; }
.row{
  display:grid;
  grid-template-columns: 1.2fr 1.2fr .8fr 1fr 1.6fr;
  gap:10px;
  padding:10px 8px;
  border-radius:12px;
}
.row.head{
  opacity:.75; font-size:12px; letter-spacing:.06em;
  border:1px solid rgba(255,255,255,.08);
}
.row:not(.head){
  border:1px solid rgba(255,255,255,.08);
}
.name{ font-weight:700; }
.me{ font-size:12px; opacity:.7; margin-left:6px; }
.actions{ display:flex; justify-content:flex-end; gap:8px; flex-wrap:wrap; }

.btn{
  padding:10px 12px; border-radius:12px;
  border:1px solid rgba(255,255,255,.14);
  background:rgba(255,255,255,.10); color:#fff; cursor:pointer;
}
.btn.ghost{ background:transparent; }
.btn.small{ padding:8px 10px; border-radius:10px; font-size:12px; }
.btn.danger{ border-color: rgba(255,90,90,.35); background: rgba(255,90,90,.12); }

.badge{
  display:inline-block; padding:6px 10px;
  border-radius:999px; font-size:12px;
  border:1px solid rgba(255,255,255,.14);
}
.badge.member{ background: rgba(255,255,255,.06); }
.badge.admin{ background: rgba(255,215,120,.12); border-color: rgba(255,215,120,.25); }
.badge.owner{ background: rgba(160,210,255,.12); border-color: rgba(160,210,255,.25); }

@media (max-width: 900px){
  .row{ grid-template-columns: 1.2fr 1.2fr .8fr 1fr; }
  .actions{ justify-content:flex-start; }
}
@media (max-width: 640px){
  .row{ grid-template-columns:1fr; gap:6px; }
  .row.head{ display:none; }
}
</style>
