<template>
  <div class="page">
    <!-- 顶部组织栏（HUD） -->
    <header class="topbar">
      <button class="icon-btn" @click="goBack" aria-label="back">‹</button>

      <div class="brand">
        <div class="logo">{{ org.logoText }}</div>
        <div class="titles">
          <div class="org-name">{{ org.name }}</div>
          <div class="org-meta">成员 {{ org.memberCount }} · 负责人 {{ org.owner }}</div>
        </div>
      </div>

      <div class="topbar-right">
        <input
          class="search"
          v-model="keyword"
          placeholder="全局搜索（成员 / 任务 / 活动）"
        />
        <button class="btn ghost" @click="toggleTheme">
          {{ theme === "dark" ? "浅色" : "深色" }}
        </button>
        <div class="user">你好，{{ userName }}</div>
      </div>
    </header>

    <!-- HUD 导航条（奇幻苏丹卡牌） -->
    <section class="tabs-center">
      <nav class="tarot-tabs" aria-label="sections">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          class="tarot-tab"
          :class="{ 'is-active': pressedKey === tab.key }"
          @click="onTabClick(tab)"
          type="button"
        >
          <div class="tarot-roman">{{ tab.roman }}</div>
          <div class="tarot-glyph" aria-hidden="true"></div>
          <div class="tarot-title">{{ tab.label }}</div>
        </button>
      </nav>
    </section>

    <!-- 主体 -->
    <main class="main" v-if="false">
      <!-- 左列 -->
      <section class="column">
        <BlockCard :title="layout[0].title">
          <ul v-if="layout[0].type === 'list'" class="list">
            <li v-for="x in layout[0].content" :key="x">{{ x }}</li>
          </ul>
          <div v-else class="muted">{{ layout[0].content }}</div>
        </BlockCard>
      </section>

      <!-- 中列 -->
      <section class="column">
        <BlockCard :title="layout[1].title">
          <div class="muted">{{ layout[1].content }}</div>
        </BlockCard>
      </section>

      <!-- 右列 -->
      <section class="column">
        <BlockCard :title="layout[2].title">
          <div class="muted">{{ layout[2].content }}</div>
        </BlockCard>
      </section>

      <!-- 固定右侧 -->
      <section class="column right">
        <BlockCard title="公告">
          <div class="muted">{{ org.blocks.notice }}</div>
          <div class="muted" style="margin-top: 10px">
            加入时间：{{ joinedTimeText }}
          </div>
        </BlockCard>

        <button class="fab" aria-label="create">＋</button>

        <div class="danger-zone">
          <div class="danger-title">危险操作</div>
          <div class="muted">一次确认</div>
          <button class="btn danger">退出组织</button>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";

/* ---------------- 基础 ---------------- */
const route = useRoute();
const router = useRouter();
const userName = "曾子桐";
const keyword = ref("");
const pressedKey = ref(null);

function goBack() {
  router.push({ name: 'orgmap' });
}

/* ---------------- 主题切换 ---------------- */
const theme = ref("dark");

onMounted(() => {
  const saved = localStorage.getItem("theme");
  if (saved) theme.value = saved;
  document.documentElement.setAttribute("data-theme", theme.value);
});

watch(theme, (v) => {
  localStorage.setItem("theme", v);
  document.documentElement.setAttribute("data-theme", v);
});

function toggleTheme() {
  theme.value = theme.value === "dark" ? "light" : "dark";
}

/* ---------------- 组织数据 ---------------- */
const ORG_CONFIG = {
  1: {
    name: "羽毛球队",
    logoText: "🏸",
    memberCount: 34,
    owner: "队长",
    blocks: {
      departments: ["队长 / 副队", "训练负责人", "外联负责人", "后勤负责人"],
      activities: "每周二、周四晚训练；校内友谊赛。",
      management: "训练计划、经费报销、器材管理。",
      recruit: "欢迎新生加入，热爱运动即可。",
      notice: "今晚训练照常。",
    },
  },
};

const org = computed(() => ORG_CONFIG[route.params.id] || ORG_CONFIG[1]);

const joinedTimeText = computed(() => {
  const t = route.query.joinTime;
  if (!t || String(t).startsWith("0001")) return "暂无";
  return String(t).slice(0, 10);
});

/* ---------------- 卡牌与布局 ---------------- */
const tabs = [
  { key: "info", label: "队内信息", roman: "I", route: "OrgInfo", match: ["OrgInfo"] },
  { key: "activity", label: "活动", roman: "II", route: "ActivityPage", match: ["ActivityPage"]},
  { key: "manage", label: "组织管理", roman: "III", match: ["OrgManage"] },
  { key: "recruit", label: "招新", roman: "IV", match: ["OrgRecruit"] },
];

const activeKey = computed(() => {
  const name = route.name ? String(route.name) : "";
  const hit = tabs.find(t => (t.match || []).includes(name));
  return hit ? hit.key : "info"; // 默认回到 info
});

const activeTab = ref("info");

function onTabClick(tab) {
  //  点一下闪烁反馈（不常亮）
  pressedKey.value = tab.key;
  setTimeout(() => (pressedKey.value = null), 180);
  if (tab.route) {
    router.push({
      name: tab.route,
      params: { id: route.params.id },
      query: route.query,
    });
  } else {
    activeTab.value = tab.key;
  }
}

const layout = computed(() => {
  const b = org.value.blocks;
  const map = {
    info: [
      { title: "队内信息", type: "list", content: b.departments },
      { title: "成员情况", type: "text", content: "成员分组、出勤情况等" },
      { title: "职责分工", type: "text", content: "负责人及职责说明" },
    ],
    activity: [
      { title: "活动", type: "text", content: b.activities },
      { title: "近期安排", type: "text", content: "训练 / 比赛 / 聚餐" },
      { title: "历史记录", type: "text", content: "过往活动记录" },
    ],
    manage: [
      { title: "组织管理", type: "text", content: b.management },
      { title: "权限审批", type: "text", content: "成员审批、权限控制" },
      { title: "资源管理", type: "text", content: "经费、器材、文件" },
    ],
    recruit: [
      { title: "招新", type: "text", content: b.recruit },
      { title: "报名流程", type: "text", content: "报名 / 面试 / 录取" },
      { title: "常见问题", type: "text", content: "训练强度、时间等" },
    ],
  };
  return map[activeTab.value];
});

/* ---------------- BlockCard ---------------- */
const BlockCard = {
  props: { title: String },
  template: `
    <div class="card">
      <div class="card-title">{{ title }}</div>
      <div class="card-body"><slot /></div>
    </div>
  `,
};
</script>

<style scoped>
/* ========== Theme Tokens ========== */
:global(:root) {
  --bg: #0f1419;
  --text: #fff;
  --muted: rgba(255, 255, 255, 0.72);

  --panel: rgba(255, 255, 255, 0.06);
  --border: rgba(255, 255, 255, 0.12);

  --gold: rgba(255, 215, 120, 0.85);
  --gold-soft: rgba(255, 215, 120, 0.18);
  --glow: rgba(255, 215, 120, 0.45);

  --shadow: rgba(0, 0, 0, 0.55);
}

:global(html[data-theme="light"]) {
  --bg: #f6f7fb;
  --text: #111;
  --muted: rgba(0, 0, 0, 0.62);

  --panel: rgba(0, 0, 0, 0.04);
  --border: rgba(0, 0, 0, 0.10);

  --gold: rgba(190, 120, 30, 0.85);
  --gold-soft: rgba(190, 120, 30, 0.18);
  --glow: rgba(190, 120, 30, 0.35);

  --shadow: rgba(0, 0, 0, 0.18);
}

/* ========== Page + Fantasy Background ========== */
.page {
  min-height: 100vh;
  background: var(--bg);
  color: var(--text);
  display: flex;
  flex-direction: column;

  position: relative;
  overflow: hidden;
  font-family: ui-sans-serif, system-ui, -apple-system, "PingFang SC", "Microsoft YaHei",
  "Helvetica Neue", Arial, "Noto Sans", "Apple Color Emoji", "Segoe UI Emoji";
}

/* 背景：光晕 + 细网纹（奇幻 UI 氛围） */
.page::before {
  content: "";
  position: absolute;
  inset: -140px;
  background:
    radial-gradient(900px 520px at 18% 18%, rgba(255, 200, 120, 0.10), transparent 60%),
    radial-gradient(760px 520px at 78% 28%, rgba(160, 210, 255, 0.07), transparent 62%),
    radial-gradient(980px 820px at 55% 88%, rgba(255, 215, 120, 0.06), transparent 68%),
    repeating-linear-gradient(135deg, rgba(255, 255, 255, 0.028) 0 1px, transparent 1px 10px);
  pointer-events: none;
  z-index: 0;
}

.page::after {
  content: "";
  position: absolute;
  inset: 0;
  background:
    radial-gradient(circle at 28% 42%, rgba(255, 255, 255, 0.05), transparent 55%),
    radial-gradient(circle at 72% 60%, rgba(255, 255, 255, 0.03), transparent 55%);
  opacity: 0.35;
  mix-blend-mode: overlay;
  pointer-events: none;
  z-index: 0;
}

.topbar,
.hud-tabs,
.main {
  position: relative;
  z-index: 1;
}

/* ========== Topbar HUD ========== */
.topbar {
  position: sticky;
  top: 0;
  z-index: 1000;

  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;

  border-bottom: 1px solid var(--border);
  background: linear-gradient(180deg, rgba(0, 0, 0, 0.35), rgba(0, 0, 0, 0.12));
  backdrop-filter: blur(12px);
}

.icon-btn {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.14);
  background: rgba(255, 255, 255, 0.08);
  color: var(--text);
  font-size: 20px;
  cursor: pointer;
  transition: transform 0.12s ease, background 0.12s ease, border-color 0.12s ease;
}

.icon-btn:hover {
  background: rgba(255, 255, 255, 0.12);
  border-color: rgba(255, 215, 120, 0.22);
  transform: translateY(-1px);
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 240px;
}

.logo {
  width: 44px;
  height: 44px;
  border-radius: 14px;
  display: grid;
  place-items: center;

  background: radial-gradient(circle at 30% 20%, rgba(255, 255, 255, 0.18), transparent 60%),
  linear-gradient(180deg, rgba(255, 215, 120, 0.20), rgba(255, 255, 255, 0.05));
  border: 1px solid rgba(255, 215, 120, 0.22);
  box-shadow: 0 10px 22px rgba(0, 0, 0, 0.35);
}

.titles .org-name {
  font-weight: 900;
  letter-spacing: 0.06em;
}

.titles .org-meta {
  font-size: 12px;
  color: var(--muted);
  margin-top: 2px;
}

.topbar-right {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 10px;
}

.search {
  width: 320px;
  max-width: 34vw;
  padding: 10px 14px;

  border-radius: 14px;
  border: 1px solid rgba(255, 255, 255, 0.14);
  background: rgba(255, 255, 255, 0.07);
  color: var(--text);
  outline: none;

  backdrop-filter: blur(10px);
  transition: border-color 0.12s ease, box-shadow 0.12s ease;
}

.search::placeholder {
  color: var(--muted);
}

.search:focus {
  border-color: var(--gold);
  box-shadow: 0 0 0 3px rgba(255, 215, 128, 0.16);
}

.btn {
  padding: 10px 12px;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.14);
  background: rgba(255, 255, 255, 0.10);
  color: var(--text);
  cursor: pointer;
  transition: transform 0.12s ease, background 0.12s ease, border-color 0.12s ease;
}

.btn:hover {
  background: rgba(255, 255, 255, 0.14);
  border-color: rgba(255, 215, 120, 0.22);
  transform: translateY(-1px);
}

.btn:active {
  transform: translateY(0);
}

.user {
  font-size: 12px;
  color: var(--muted);
  white-space: nowrap;
}

/* ========== HUD Tabs (Sultan Cards) ========== */
.hud-tabs {
  display: flex;
  justify-content: center;
  gap: 14px;
  padding: 12px 16px 10px;
  border-bottom: 1px solid var(--border);

  background: linear-gradient(180deg, rgba(0, 0, 0, 0.22), rgba(0, 0, 0, 0));
  backdrop-filter: blur(10px);
}

.tab-label {
  display: block;
  font-size: 14px;
  font-weight: 900;
  letter-spacing: 0.12em;
}

.tab-sub {
  display: block;
  margin-top: 6px;
  font-size: 11px;
  letter-spacing: 0.22em;
  opacity: 0.65;
}

/* 奇幻苏丹卡牌（更像游戏 UI） */
.sultan-card {
  position: relative;
  width: 220px;
  height: 78px;
  border-radius: 22px;
  border: 1px solid rgba(255, 215, 140, 0.25);

  background:
    radial-gradient(circle at 20% 15%, rgba(255, 255, 255, 0.14), transparent 55%),
    radial-gradient(circle at 80% 65%, rgba(255, 215, 120, 0.08), transparent 55%),
    linear-gradient(180deg, rgba(28, 32, 38, 0.92), rgba(15, 18, 22, 0.92));

  color: var(--text);
  cursor: pointer;
  text-align: left;
  padding: 14px 16px;

  box-shadow:
    0 18px 36px rgba(0, 0, 0, 0.48),
    inset 0 1px 0 rgba(255, 255, 255, 0.10),
    inset 0 -6px 12px rgba(0, 0, 0, 0.55);

  transition: transform 0.12s ease, box-shadow 0.12s ease, border-color 0.12s ease;
}

.sultan-card::before {
  content: "";
  position: absolute;
  inset: 8px;
  border-radius: 16px;
  border: 1px dashed rgba(255, 215, 140, 0.22);
  opacity: 0.9;
  pointer-events: none;
}

.sultan-card::after {
  content: "";
  position: absolute;
  left: 14px;
  right: 14px;
  bottom: 10px;
  height: 3px;
  border-radius: 999px;
  background: rgba(255, 215, 140, 0.12);
  box-shadow: 0 0 12px rgba(255, 215, 120, 0.08);
  pointer-events: none;
}

.sultan-card:hover {
  transform: translateY(-2px);
  border-color: rgba(255, 215, 140, 0.42);
  box-shadow:
    0 22px 44px rgba(0, 0, 0, 0.55),
    0 0 0 1px rgba(255, 215, 140, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.12),
    inset 0 -8px 16px rgba(0, 0, 0, 0.58);
}

.sultan-card:active {
  transform: translateY(0px);
}

.sultan-card.active {
  border-color: rgba(255, 215, 120, 0.72);
  box-shadow:
    0 26px 50px rgba(0, 0, 0, 0.62),
    0 0 0 2px rgba(255, 215, 120, 0.18),
    0 0 26px rgba(255, 215, 120, 0.25),
    inset 0 1px 0 rgba(255, 255, 255, 0.16),
    inset 0 -10px 18px rgba(0, 0, 0, 0.62);
}

.sultan-card.active::after {
  background: linear-gradient(
    90deg,
    rgba(255, 215, 120, 0),
    rgba(255, 215, 120, 0.95),
    rgba(255, 215, 120, 0)
  );
  box-shadow: 0 0 18px rgba(255, 215, 120, 0.35);
}

/* ========== Main Grid ========== */
.main {
  flex: 1;
  display: grid;
  grid-template-columns: 1fr 1fr 1fr 0.9fr;
  gap: 16px;
  padding: 16px;

  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
}

/* ========== Panels / Cards ========== */
.card {
  position: relative;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.07), rgba(255, 255, 255, 0.04));
  border: 1px solid rgba(255, 255, 255, 0.10);
  border-radius: 18px;
  padding: 14px;
  box-shadow: 0 16px 30px rgba(0, 0, 0, 0.35);
}

.card::before {
  content: "";
  position: absolute;
  inset: 8px;
  border-radius: 14px;
  border: 1px solid rgba(255, 215, 120, 0.10);
  pointer-events: none;
  opacity: 0.9;
}

.card-title {
  font-weight: 900;
  letter-spacing: 0.08em;
  margin-bottom: 10px;
}

.card-body {
  line-height: 1.6;
}

.muted {
  color: var(--muted);
}

.column {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.column.right {
  position: relative;
}

/* list */
.list {
  padding-left: 18px;
  margin: 0;
}
.list li {
  margin: 6px 0;
  color: var(--muted);
}

/* ========== FAB (game-like) ========== */
.fab {
  position: sticky;
  top: 150px; /* 滚动时像 UI 组件黏在侧边 */
  align-self: flex-end;

  width: 56px;
  height: 56px;
  border-radius: 50%;
  border: 1px solid rgba(255, 215, 120, 0.35);
  background:
    radial-gradient(circle at 30% 20%, rgba(255, 255, 255, 0.18), transparent 60%),
    linear-gradient(180deg, rgba(255, 215, 120, 0.18), rgba(255, 255, 255, 0.06));
  color: var(--text);
  font-size: 26px;
  cursor: pointer;
  box-shadow:
    0 18px 40px rgba(0, 0, 0, 0.55),
    0 0 22px rgba(255, 215, 120, 0.18);
  transition: transform 0.12s ease, box-shadow 0.12s ease;
}

.fab:hover {
  transform: translateY(-2px);
  box-shadow:
    0 22px 46px rgba(0, 0, 0, 0.62),
    0 0 28px rgba(255, 215, 120, 0.28);
}

/* ========== Danger Zone ========== */
.danger-zone {
  margin-top: auto;
  padding: 14px;
  border-radius: 18px;
  border: 1px solid rgba(255, 90, 90, 0.25);
  background: linear-gradient(180deg, rgba(255, 90, 90, 0.08), rgba(255, 255, 255, 0.02));
  box-shadow: 0 16px 30px rgba(0, 0, 0, 0.35);
  position: relative;
}

.danger-zone::before {
  content: "";
  position: absolute;
  inset: 8px;
  border-radius: 14px;
  border: 1px dashed rgba(255, 90, 90, 0.20);
  pointer-events: none;
}

.danger-title {
  font-weight: 900;
  letter-spacing: 0.10em;
  margin-bottom: 8px;
}

.btn.danger {
  width: 100%;
  margin-top: 10px;
  border-color: rgba(255, 90, 90, 0.35);
  background: rgba(255, 90, 90, 0.12);
}

.btn.danger:hover {
  background: rgba(255, 90, 90, 0.16);
  border-color: rgba(255, 90, 90, 0.55);
}

/* ========= Tarot Tabs (replace Sultan Tabs) ========= */
.tarot-tabs{
  display: flex;
  justify-content: center;
  gap: 26px;         /* 卡牌间距更大 */
  padding: 0;        /* 不再像导航条那样加上下边距 */
  border: 0;         /* 不要底部线 */
  background: none;  /* 不要顶部渐变背景 */
  backdrop-filter: none;
  flex-wrap: wrap;   /* 小屏自动换行 */
  perspective: 900px;
}

/* 按键本体：羊皮纸 + 金棕描边 + 轻纹理 */
.tarot-tab{
  width: 220px;
  height: 400px;
  border: 0;
  cursor: pointer;
  border-radius: 18px;
  position: relative;
  padding: 14px 14px 12px;
  text-align: center;
  color: rgba(35, 25, 15, 0.92);
  transform-style: preserve-3d;
  will-change: transform;

  background:
    /* 角落云纹（兼容版：不用双 stop） */
    radial-gradient(circle at 10% 10%,
    rgba(168,120,70,.18) 0%,
    rgba(168,120,70,.18) 14%,
    transparent 15%
    ),
    radial-gradient(circle at 18% 18%,
    rgba(168,120,70,.14) 0%,
    rgba(168,120,70,.14) 10%,
    transparent 11%
    ),
    radial-gradient(circle at 90% 10%,
    rgba(168,120,70,.18) 0%,
    rgba(168,120,70,.18) 14%,
    transparent 15%
    ),
    radial-gradient(circle at 82% 18%,
    rgba(168,120,70,.14) 0%,
    rgba(168,120,70,.14) 10%,
    transparent 11%
    ),
    radial-gradient(circle at 10% 90%,
    rgba(168,120,70,.18) 0%,
    rgba(168,120,70,.18) 14%,
    transparent 15%
    ),
    radial-gradient(circle at 18% 82%,
    rgba(168,120,70,.14) 0%,
    rgba(168,120,70,.14) 10%,
    transparent 11%
    ),
    radial-gradient(circle at 90% 90%,
    rgba(168,120,70,.18) 0%,
    rgba(168,120,70,.18) 14%,
    transparent 15%
    ),
    radial-gradient(circle at 82% 82%,
    rgba(168,120,70,.14) 0%,
    rgba(168,120,70,.14) 10%,
    transparent 11%
    ),

      /* 纸张高光 + 暗角（同理改成 3 stop，最稳） */
    radial-gradient(circle at 25% 20%,
    rgba(255,255,255,.55) 0%,
    rgba(255,255,255,.55) 60%,
    transparent 61%
    ),
    radial-gradient(circle at 80% 70%,
    rgba(0,0,0,.06) 0%,
    rgba(0,0,0,.06) 55%,
    transparent 56%
    ),

      /* 纸纹 */
    repeating-linear-gradient(135deg, rgba(0,0,0,.03) 0 1px, transparent 1px 10px),

      /* 底色 */
    linear-gradient(180deg, #efe7d9, #e6dccb);

  box-shadow:
    0 14px 28px rgba(0,0,0,.30),      /* 主投影 */
    0 2px 0 rgba(255,255,255,.35),    /* 上沿高光 */
    inset 0 1px 0 rgba(255,255,255,.75);

  transition: transform .14s ease, box-shadow .14s ease, filter .14s ease;
}

/* 双层描边（外框+内框） */
.tarot-tab::before{
  content:"";
  position:absolute;
  inset: 6px;
  border-radius: 14px;
  border: 2px solid rgba(168, 120, 70, .55);
  box-shadow: inset 0 0 0 1px rgba(255,255,255,.35);
  pointer-events:none;
}
.tarot-tab::after{
  content:"";
  position:absolute;
  inset: 12px;
  border-radius: 10px;
  border: 1px solid rgba(168, 120, 70, .28);
  pointer-events:none;
}

/* 顶部罗马数字 */
.tarot-roman{
  font-family: ui-serif, "Times New Roman", "Songti SC", serif;
  font-weight: 800;
  letter-spacing: .14em;
  font-size: 13px;
  opacity: .9;
}

/* 中央星芒：放射线 + 星标 */
.tarot-glyph{
  width: 68px;
  height: 68px;
  margin: 10px auto 8px;
  border-radius: 50%;
  position: relative;

  background:
    radial-gradient(circle at center, rgba(168,120,70,.14), transparent 62%),
    repeating-conic-gradient(
      from 0deg,
      rgba(35,25,15,.18) 0 2deg,
      transparent 2deg 18deg
    );
}
.tarot-glyph::before{
  content:"✦";
  position:absolute;
  inset:0;
  display:grid;
  place-items:center;
  font-size: 30px;
  color: rgba(168,120,70,.95);
  text-shadow: 0 0 10px rgba(168,120,70,.18);
}

/* 底部标题：塔罗牌感 */
.tarot-title{
  margin-top: 2px;
  font-family: ui-serif, "Times New Roman", "Songti SC", serif;
  font-weight: 900;
  letter-spacing: .40em;
  font-size: 16px;
}

/* hover/active */
.tarot-tab:hover{
  transform:
    translateY(-8px)
    scale(1.02)
    rotateX(2deg)
    rotateY(-2deg);

  box-shadow:
    0 26px 48px rgba(0,0,0,.40),
    0 8px 24px rgba(0,0,0,.25),
    0 0 18px rgba(255,215,120,.25),
    inset 0 1px 0 rgba(255,255,255,.85);

  filter: saturate(1.08);
}

.tarot-tab:active{
  transform:
    translateY(-3px)
    scale(0.99)
    rotateX(0deg)
    rotateY(0deg);
}

.tarot-tab{
  position: relative;
  z-index: 1;
}

.tarot-tab.is-active::before{
  content: "";
  position: absolute;
  inset: -14px;                 /* 向外扩散 */
  border-radius: 28px;
  background:
    radial-gradient(
      ellipse at center,
      rgba(255,215,120,.45),
      rgba(255,215,120,.18) 45%,
      transparent 70%
    );
  filter: blur(14px);
  z-index: -1;                  /* 在卡牌后面 */
}

/* 当前高亮（对应 activeKey） */
.tarot-tab.is-active{
  transform: translateY(-6px);
  box-shadow:
    0 30px 56px rgba(0,0,0,.45),
    0 0 0 2px rgba(255,215,120,.35),
    0 0 32px rgba(255,215,120,.55),   /* 外圈金色光 */
    inset 0 1px 0 rgba(255,255,255,.95);
}

.tarot-tab.is-active .tarot-glyph{
  box-shadow: 0 0 18px rgba(255,215,120,.55);
}
.tarot-tab.is-active .tarot-glyph::before{
  text-shadow:
    0 0 6px rgba(255,215,120,.85),
    0 0 14px rgba(255,215,120,.65);
}

/* ========== Tabs Center (把塔罗牌放到页面中间) ========== */
.tabs-center{
  flex: 1;                 /* 关键：占满 topbar 下面的剩余高度 */
  display: flex;
  align-items: center;     /* 垂直居中 */
  justify-content: center; /* 水平居中 */
  padding: 24px 16px;
}

/* ========== Responsive ========== */
@media (max-width: 1200px) {
  .main {
    grid-template-columns: 1fr 1fr;
  }
  .hud-tabs {
    flex-wrap: wrap;
  }
  .sultan-card {
    width: 240px;
  }
}

@media (max-width: 760px) {
  .topbar {
    flex-wrap: wrap;
    gap: 10px;
  }
  .topbar-right {
    width: 100%;
    justify-content: space-between;
  }
  .search {
    width: 100%;
    max-width: unset;
    flex: 1;
  }
  .main {
    grid-template-columns: 1fr;
  }
  .sultan-card {
    width: 100%;
    max-width: 520px;
  }
}

</style>
