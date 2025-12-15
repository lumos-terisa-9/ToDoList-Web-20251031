<template>
  <div class="page">
    <!-- 顶部栏（不改 HeaderBar.vue，只做让位） -->
    <header class="topbar">
      <button class="icon-btn" @click="goBack">‹</button>

      <div class="title">
        <div class="title-main">{{ orgName || "组织信息" }}</div>
        <div class="title-sub">队内信息 · 社团</div>
      </div>

      <div class="topbar-right">
        <input class="search" v-model="keyword" placeholder="搜索成员 / 职位 / 关键词…" />
        <button class="btn ghost" @click="refresh">刷新</button>
      </div>
    </header>

    <main class="main">
      <!-- 左：组织概览 -->
      <section class="col">
        <div class="card hero">
          <div class="hero-row">
            <div class="hero-logo">
              <img v-if="org.logo_url" :src="org.logo_url" alt="logo" />
              <div v-else class="hero-logo-fallback">🏷️</div>
            </div>

            <div class="hero-meta">
              <div class="hero-name">{{ org.org_name || orgName || "未命名组织" }}</div>
              <div class="hero-tags">
                <span class="tag">社团</span>
                <span class="tag muted">创建者ID：{{ org.creator_id ?? "—" }}</span>
              </div>
              <div class="hero-joined">加入时间：{{ joinedText }}</div>
            </div>
          </div>

          <div class="hero-stats">
            <div class="stat">
              <div class="stat-num">{{ stats.member_count }}</div>
              <div class="stat-label">成员</div>
            </div>
            <div class="stat">
              <div class="stat-num">{{ stats.admin_count }}</div>
              <div class="stat-label">管理员</div>
            </div>
            <div class="stat">
              <div class="stat-num">{{ stats.week_activity_count }}</div>
              <div class="stat-label">本周活动</div>
            </div>
            <div class="stat">
              <div class="stat-num">{{ stats.active_rate }}</div>
              <div class="stat-label">活跃度</div>
            </div>
          </div>
        </div>

        <div class="card">
          <div class="card-title">组织简介</div>
          <div class="card-body muted">
            {{ org.intro || "暂无简介（后续可从数据库补充 intro 字段）" }}
          </div>
        </div>

        <div class="card">
          <div class="card-title">公告</div>
          <div class="card-body muted">
            {{ org.notice || "暂无公告" }}
          </div>
        </div>
      </section>

      <!-- 中：职位 / 分工 -->
      <section class="col">
        <div class="card">
          <div class="card-title">职位与分工</div>
          <div class="card-body">
            <div v-if="rolesFiltered.length === 0" class="muted">暂无分工信息</div>

            <div v-for="r in rolesFiltered" :key="r.title" class="role-row">
              <div class="role-title">{{ r.title }}</div>
              <div class="role-person">{{ r.username || "—" }}</div>
              <div class="role-desc muted">{{ r.desc || "" }}</div>
            </div>
          </div>
        </div>

        <div class="card">
          <div class="card-title">队内安排</div>
          <div class="card-body muted">
            {{ org.schedule || "暂无安排（可后续接活动表/日历）" }}
          </div>
        </div>
      </section>

      <!-- 右：成员列表 -->
      <section class="col">
        <div class="card">
          <div class="card-title">成员列表</div>
          <div class="card-body">
            <div v-if="membersFiltered.length === 0" class="muted">暂无成员数据（后端可补 members 列表）</div>

            <div v-for="m in membersFiltered" :key="m.user_id" class="member-row">
              <div class="member-name">{{ m.username }}</div>
              <div class="member-role tag small">{{ m.role }}</div>
              <div class="member-joined muted">{{ m.joined_at || "—" }}</div>
            </div>
          </div>
        </div>

        <div class="card danger-zone">
          <div class="card-title">一次确认</div>
          <div class="card-body">
            <button class="btn danger" @click="leaveOrg">退出组织</button>
          </div>
        </div>
      </section>
    </main>

    <div v-if="loading" class="toast">加载中…</div>
    <div v-if="error" class="toast error">加载失败：{{ error }}</div>
  </div>
</template>

<script setup>
import axios from "axios";
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";

const route = useRoute();
const router = useRouter();

const orgId = computed(() => String(route.params.id || ""));
const orgName = computed(() => String(route.query.name || ""));

const keyword = ref("");
const loading = ref(false);
const error = ref("");

/**
 * 你现在后端 my-organizations 只有：org_name / creator_id / logo_url / joined_at
 * 这里的 org / stats / roles / members 先按“可扩展结构”设计：
 * - 后端以后补字段，前端无需大改
 */
const org = ref({
  org_name: "",
  creator_id: null,
  logo_url: "",
  joined_at: "",
  intro: "",
  notice: "",
  schedule: "",
});

const stats = ref({
  member_count: "—",
  admin_count: "—",
  week_activity_count: "—",
  active_rate: "—",
});

const roles = ref([
  { title: "队长", username: "", desc: "统筹训练与队内事务" },
  { title: "训练负责人", username: "", desc: "安排训练计划与分组" },
  { title: "外联负责人", username: "", desc: "联系比赛与合作" },
  { title: "后勤负责人", username: "", desc: "器材与后勤保障" },
]);

const members = ref([]);

/** axios 实例：带 token（跟你 MapPage 一样逻辑） */
const api = axios.create({
  baseURL: "http://localhost:8080",
  timeout: 10000,
});

api.interceptors.request.use((config) => {
  const raw = localStorage.getItem("token");
  if (raw) {
    try {
      const obj = JSON.parse(raw);
      const accessToken = obj?.data?.access_token;
      if (accessToken) config.headers.Authorization = `Bearer ${accessToken}`;
    } catch {
      config.headers.Authorization = `Bearer ${raw}`;
    }
  }
  return config;
});

/**
 * 拉取组织信息（你现在没有 org_id 字段，所以这里先用 name 做兜底）
 * ✅ 当前可用：从 /api/organization/my-organizations 找到第 orgId 个（和 MapPage 的逻辑一致）
 * ⭐ 更推荐：你后端以后加一个 /api/organization/:id/info，前端这里改 1 行即可
 */
async function fetchOrgInfo() {
  loading.value = true;
  error.value = "";
  try {
    // 先复用你现有接口：my-organizations
    const resp = await api.get("/api/organization/my-organizations");
    const list = resp.data?.data;
    if (Array.isArray(list)) {
      const idx = Math.max(0, Number(orgId.value) - 1); // id=1..4 -> index=0..3
      const o = list[idx] || list[0];
      if (o) {
        org.value.org_name = o.org_name || orgName.value || "";
        org.value.creator_id = o.creator_id ?? null;
        org.value.logo_url = o.logo_url || "";
        org.value.joined_at = o.joined_at || "";
      }
    }
    // stats / members 先占位：等后端补字段再填
  } catch (e) {
    error.value = e?.message || String(e);
  } finally {
    loading.value = false;
  }
}

const joinedText = computed(() => {
  const t = org.value.joined_at;
  if (!t || String(t).startsWith("0001-01-01")) return "暂无";
  const d = new Date(t);
  if (Number.isNaN(d.getTime())) return String(t);
  const pad = (n) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;
});

const rolesFiltered = computed(() => {
  const k = keyword.value.trim();
  if (!k) return roles.value;
  return roles.value.filter((r) => `${r.title}${r.username}${r.desc}`.includes(k));
});

const membersFiltered = computed(() => {
  const k = keyword.value.trim();
  if (!k) return members.value;
  return members.value.filter((m) => `${m.username}${m.role}${m.joined_at}`.includes(k));
});

function goBack() {
  router.back();
}

function refresh() {
  fetchOrgInfo();
}

function leaveOrg() {
  alert("这里接后端退出组织 API：成功后 router.push('/orgmap')");
}

onMounted(fetchOrgInfo);
</script>

<style scoped>
/* ===== Theme tokens（深/浅色自动适配：跟你前面用的一套） ===== */
:global(:root) {
  --bg: #0f1419;
  --text: #ffffff;
  --muted: rgba(255, 255, 255, 0.78);

  --panel: rgba(255, 255, 255, 0.10);
  --panel-border: rgba(255, 255, 255, 0.10);
  --divider: rgba(255, 255, 255, 0.08);

  --btn: rgba(255, 255, 255, 0.14);
  --btn-ghost: rgba(255, 255, 255, 0.10);
  --btn-ghost-border: rgba(255, 255, 255, 0.14);

  --input-bg: rgba(255, 255, 255, 0.06);
  --input-border: rgba(255, 255, 255, 0.12);

  --danger: rgba(255, 107, 107, 0.8);
  --glow: rgba(255, 215, 128, 0.45);
}

:global(html[data-theme="light"]) {
  --bg: #f6f7fb;
  --text: #0f1419;
  --muted: rgba(15, 20, 25, 0.70);

  --panel: rgba(0, 0, 0, 0.04);
  --panel-border: rgba(0, 0, 0, 0.08);
  --divider: rgba(0, 0, 0, 0.08);

  --btn: rgba(0, 0, 0, 0.06);
  --btn-ghost: rgba(0, 0, 0, 0.04);
  --btn-ghost-border: rgba(0, 0, 0, 0.10);

  --input-bg: rgba(0, 0, 0, 0.03);
  --input-border: rgba(0, 0, 0, 0.10);

  --danger: rgba(220, 60, 60, 0.85);
  --glow: rgba(255, 180, 80, 0.35);
}

/* ===== Layout ===== */
.page {
  min-height: 100vh;
  padding-top: 70px; /* 让 HeaderBar 不挡住 */
  background: var(--bg);
  color: var(--text);
}

.topbar {
  display: flex;
  align-items: center;
  gap: 12px;

  padding: 14px 16px;
  border-bottom: 1px solid var(--divider);
}

.icon-btn {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  border: none;
  cursor: pointer;
  background: var(--btn);
  color: var(--text);
  font-size: 18px;
}

.title {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
}

.title-main {
  font-size: 18px;
  font-weight: 800;
}

.title-sub {
  font-size: 12px;
  color: var(--muted);
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.search {
  width: 300px;
  max-width: 32vw;
  padding: 10px 14px;
  border-radius: 14px;
  border: 1px solid var(--input-border);
  background: var(--input-bg);
  color: var(--text);
  outline: none;
}

.search::placeholder {
  color: var(--muted);
}

.search:focus {
  border-color: var(--glow);
  box-shadow: 0 0 0 3px rgba(255, 215, 128, 0.18);
}

.btn {
  padding: 10px 12px;
  border-radius: 12px;
  border: none;
  cursor: pointer;
  background: var(--btn);
  color: var(--text);
}

.btn.ghost {
  background: var(--btn-ghost);
  border: 1px solid var(--btn-ghost-border);
}

.btn.danger {
  width: 100%;
  background: var(--danger);
}

/* ===== Content grid ===== */
.main {
  padding: 16px;
  display: grid;
  grid-template-columns: 1.2fr 1fr 1fr;
  gap: 14px;
}

.col {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.card {
  border-radius: 18px;
  padding: 14px;
  background: var(--panel);
  border: 1px solid var(--panel-border);
  backdrop-filter: blur(12px);
}

.card-title {
  font-weight: 800;
  margin-bottom: 10px;
}

.card-body {
  font-size: 13px;
}

.muted {
  color: var(--muted);
  line-height: 1.6;
}

/* ===== Hero ===== */
.hero-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.hero-logo {
  width: 54px;
  height: 54px;
  border-radius: 16px;
  overflow: hidden;
  background: var(--btn-ghost);
  display: flex;
  align-items: center;
  justify-content: center;
}

.hero-logo img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.hero-logo-fallback {
  font-size: 22px;
}

.hero-name {
  font-size: 18px;
  font-weight: 900;
}

.hero-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 4px;
}

.tag {
  padding: 4px 8px;
  border-radius: 999px;
  border: 1px solid var(--panel-border);
  background: var(--btn-ghost);
  font-size: 12px;
}

.tag.small {
  padding: 2px 6px;
  font-size: 11px;
}

.hero-joined {
  margin-top: 6px;
  color: var(--muted);
  font-size: 12px;
}

.hero-stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 10px;
  margin-top: 12px;
}

.stat {
  border-radius: 14px;
  padding: 10px;
  border: 1px solid var(--panel-border);
  background: var(--btn-ghost);
}

.stat-num {
  font-size: 16px;
  font-weight: 900;
}

.stat-label {
  font-size: 12px;
  color: var(--muted);
  margin-top: 2px;
}

/* ===== Role / Member list ===== */
.role-row {
  display: grid;
  grid-template-columns: 90px 1fr;
  gap: 8px 10px;
  padding: 10px 0;
  border-bottom: 1px solid var(--divider);
}

.role-row:last-child {
  border-bottom: none;
}

.role-title {
  font-weight: 800;
}

.role-person {
  font-weight: 700;
}

.role-desc {
  grid-column: 1 / -1;
}

.member-row {
  display: grid;
  grid-template-columns: 1fr auto auto;
  gap: 10px;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid var(--divider);
}

.member-row:last-child {
  border-bottom: none;
}

.danger-zone .btn.danger {
  margin-top: 8px;
}

/* ===== Toast ===== */
.toast {
  position: fixed;
  left: 50%;
  bottom: 18px;
  transform: translateX(-50%);
  padding: 10px 12px;
  border-radius: 12px;
  border: 1px solid var(--panel-border);
  background: var(--panel);
  color: var(--text);
  backdrop-filter: blur(12px);
}

.toast.error {
  border-color: rgba(255, 107, 107, 0.55);
}
</style>
