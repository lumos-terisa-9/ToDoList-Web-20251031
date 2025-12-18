<template>
  <div class="page">
    <!-- 顶部组织栏 -->
    <header class="topbar">
      <button class="icon-btn" @click="goBack">‹</button>

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
          {{ theme === 'dark' ? '浅色' : '深色' }}
        </button>
        <div class="user">你好，{{ userName }}</div>
      </div>
    </header>

    <!-- 主体 -->
    <main class="main">
      <!-- 苏丹风卡牌菜单 -->
      <div class="sultan-menu">
        <button
            v-for="tab in tabs"
            :key="tab.key"
            class="sultan-card"
            :class="{ active: activeTab === tab.key }"
            @click="activeTab = tab.key"
        >
          {{ tab.label }}
        </button>
      </div>

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

        <button class="fab">＋</button>

        <div class="confirm">
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

function goBack() {
  router.back();
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
  { key: "info", label: "队内信息" },
  { key: "activity", label: "活动" },
  { key: "manage", label: "组织管理" },
  { key: "recruit", label: "招新" },
];

const activeTab = ref("info");

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
/* ========== Theme ========== */
:global(:root) {
  --bg: #0f1419;
  --text: #fff;
  --muted: rgba(255,255,255,.75);
  --panel: rgba(255,255,255,.1);
  --border: rgba(255,255,255,.1);
  --glow: rgba(255,215,128,.45);
}
:global(html[data-theme="light"]) {
  --bg: #f6f7fb;
  --text: #111;
  --muted: rgba(0,0,0,.6);
  --panel: rgba(0,0,0,.05);
  --border: rgba(0,0,0,.1);
  --glow: rgba(255,180,80,.35);
}

/* ========== Layout ========== */
.page {
  min-height: 100vh;
  padding-top: 70px;
  background: var(--bg);
  color: var(--text);
  display: flex;
  flex-direction: column;
}

.topbar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--border);
}

.main {
  flex: 1;
  position: relative;
  display: grid;
  grid-template-columns: 1fr 1fr 1fr 0.9fr;
  gap: 16px;
  padding: 16px;
}

/* ========== Sultan Cards ========== */
.sultan-menu {
  position: absolute;
  left: 50%;
  top: 55%;
  transform: translate(-50%, -50%);

  display: flex;
  flex-direction: row;   /* ✅ 改成横向排列（左->右） */
  gap: 22px;             /* 卡牌间距 */

  z-index: 50;
}

/* =========================
   巴洛克手绘风 · 卡牌按钮
========================= */

.sultan-card {
  position: relative;
  width: 230px;
  height: 90px;
  border-radius: 28px 24px 30px 22px;
  border: 2px solid rgba(210, 180, 120, 0.45);
  background:
      radial-gradient(
          circle at top left,
          rgba(255, 255, 255, 0.12),
          rgba(255, 255, 255, 0.02) 60%
      ),
      linear-gradient(
          180deg,
          rgba(40, 45, 50, 0.95),
          rgba(20, 24, 28, 0.95)
      );
  color: var(--text);
  font-size: 17px;
  font-weight: 900;
  letter-spacing: 0.12em;
  cursor: pointer;
  box-shadow:
      0 20px 36px rgba(0, 0, 0, 0.45),
      inset 0 2px 4px rgba(255, 255, 255, 0.12),
      inset 0 -3px 6px rgba(0, 0, 0, 0.6);
  transition:
      transform 0.15s ease,
      box-shadow 0.15s ease,
      background 0.15s ease,
      border-color 0.15s ease;
}

.sultan-card::before {
  content: "";
  position: absolute;
  inset: 8px;
  border-radius: 22px 18px 24px 16px;
  border: 1px dashed rgba(220, 190, 140, 0.35);
  pointer-events: none;
}

.sultan-card::after {
  content: "";
  position: absolute;
  inset: 0;
  border-radius: inherit;
  background:
      radial-gradient(
          circle at 30% 20%,
          rgba(255, 255, 255, 0.08),
          transparent 55%
      );
  pointer-events: none;
}

.sultan-card:hover {
  transform: translateY(-2px) rotate(-0.2deg);
  border-color: rgba(240, 210, 150, 0.75);
}

.sultan-card:active {
  transform: translateY(1px);
}

.sultan-card.active {
  border-color: rgba(255, 215, 120, 0.95);
  box-shadow:
      0 24px 48px rgba(0, 0, 0, 0.55),
      0 0 0 2px rgba(255, 215, 120, 0.6),
      0 0 28px rgba(255, 215, 120, 0.55),
      inset 0 2px 6px rgba(255, 255, 255, 0.18),
      inset 0 -4px 8px rgba(0, 0, 0, 0.6);
}

/* ========== Cards ========== */
.card {
  background: var(--panel);
  border: 1px solid var(--border);
  border-radius: 18px;
  padding: 14px;
}
.card-title {
  font-weight: 700;
  margin-bottom: 8px;
}
.muted {
  color: var(--muted);
}

/* ========== Others ========== */
.column {
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.column.right {
  position: relative;
}
.fab {
  position: absolute;
  right: 0;
  top: 45%;
  width: 56px;
  height: 56px;
  border-radius: 50%;
}
.confirm {
  margin-top: auto;
}
/* 顶栏右侧：输入框 + 按钮 + 用户名 */
.topbar-right {
  position: fixed;
  top: 70px;          /* HeaderBar 高度，和你 page padding-top 一致 */
  right: 24px;        /* 离右侧一点边距 */

  display: flex;
  align-items: center;
  gap: 10px;

  z-index: 1200;      /* 比 HeaderBar(1000) 还高，保证不被挡 */
}

/* 搜索框：恢复“之前的胶囊玻璃风格” */
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
}

.search::placeholder {
  color: var(--muted);
}

/* 深浅色都好看的 focus */
.search:focus {
  border-color: var(--glow);
  box-shadow: 0 0 0 3px rgba(255, 215, 128, 0.18);
}

/* 顶栏按钮恢复“玻璃按钮” */
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
  transform: translateY(-1px);
}

.btn:active {
  transform: translateY(0);
}

/* 用户名显示更像之前那版 */
.user {
  font-size: 12px;
  color: var(--muted);
  white-space: nowrap;
}
</style>
