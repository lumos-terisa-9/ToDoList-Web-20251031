<template>
  <div class="recruit-page">
    <!-- 顶部装饰栏 -->
    <header class="topbar">
      <button class="icon-btn" @click="goBack" aria-label="back">‹</button>

      <div class="brand">
        <div class="logo">{{ recruit.logoText || "𓂀" }}</div>
        <div class="titles">
          <div class="org-name">{{ recruit.orgName }}</div>
          <div class="org-meta">{{ recruit.orgMeta }}</div>
        </div>
      </div>

      <div class="topbar-right">
        <button class="btn ghost" @click="scrollTo('positions')">岗位</button>
        <button class="btn ghost" @click="scrollTo('timeline')">流程</button>
        <button class="btn" @click="scrollTo('apply')">立即报名</button>
      </div>
    </header>

    <!-- HERO -->
    <section class="hero">
      <div class="hero-card">
        <div class="ornate-line"></div>

        <div class="hero-kicker">{{ recruit.heroKicker }}</div>

        <!-- 这里用 v-html 是为了支持 gold 高亮 -->
        <h1 class="hero-title" v-html="recruit.heroTitleHtml"></h1>

        <p class="hero-subtitle">
          {{ recruit.heroSubtitle }}
        </p>

        <div class="hero-actions">
          <button class="btn" @click="scrollTo('apply')">🔖 领取报名卷轴</button>
          <button class="btn ghost" @click="openPreview()">👁️ 先看看我们在做什么</button>
        </div>

        <div class="seal-row">
          <div class="seal">
            <div class="seal-inner">{{ recruit.sealText }}</div>
          </div>
          <div class="seal-text">
            <div class="seal-title">{{ recruit.sealTitle }}</div>
            <div class="seal-desc">{{ recruit.sealDesc }}</div>
          </div>
        </div>

        <div class="ornate-line"></div>
      </div>

      <div class="hero-side">
        <div class="stat-card">
          <div class="stat-title">{{ recruit.sideTitle }}</div>
          <ul class="stat-list">
            <li v-for="x in recruit.sideList" :key="x">{{ x }}</li>
          </ul>
        </div>

        <div class="quote-card">
          <div class="quote-mark">“</div>
          <div class="quote" v-html="recruit.quoteHtml"></div>
          <div class="quote-from">—— {{ recruit.quoteFrom }}</div>
        </div>
      </div>
    </section>

    <!-- 关于我们 -->
    <section class="section parchment" ref="aboutEl">
      <div class="section-header">
        <h2>关于我们</h2>
        <p>{{ recruit.aboutIntro }}</p>
      </div>

      <div class="grid-3">
        <div class="panel" v-for="c in recruit.aboutCards" :key="c.title">
          <div class="panel-title">{{ c.title }}</div>
          <div class="panel-body">{{ c.body }}</div>
        </div>
      </div>
    </section>

    <!-- 四张命运之卡 -->
    <section class="section" ref="positionsEl">
      <div class="section-header">
        <h2 id="positions">{{ recruit.deckSectionTitle }}</h2>
        <p>{{ recruit.deckSectionSubtitle }}</p>
      </div>

      <div class="deck">
        <button
          v-for="d in decks"
          :key="d.key"
          class="card"
          :class="{ active: activeDeck === d.key }"
          @click="activeDeck = d.key"
          type="button"
        >
          <div class="card-top">
            <div class="card-emblem">{{ d.emblem }}</div>
            <div class="card-title">{{ d.title }}</div>
          </div>

          <div class="card-body">
            <div class="card-tagline">{{ d.tagline }}</div>
            <div class="card-desc">{{ d.desc }}</div>

            <div class="card-tags">
              <span v-for="t in d.tags" :key="t" class="tag">{{ t }}</span>
            </div>
          </div>

          <div class="card-foot">
            <span class="card-foot-hint">点击选择</span>
            <span class="card-foot-arrow">›</span>
          </div>
        </button>
      </div>

      <div class="detail parchment">
        <div class="detail-head">
          <div class="detail-title">{{ deckDetail.title }}</div>
          <div class="detail-sub">{{ deckDetail.long }}</div>
        </div>

        <div class="grid-2">
          <div class="panel">
            <div class="panel-title">你会做什么</div>
            <ul class="bullets">
              <li v-for="x in deckDetail.doing" :key="x">{{ x }}</li>
            </ul>
          </div>

          <div class="panel">
            <div class="panel-title">你会学到什么</div>
            <ul class="bullets">
              <li v-for="x in deckDetail.learn" :key="x">{{ x }}</li>
            </ul>
          </div>
        </div>
      </div>
    </section>

    <!-- 岗位卡牌 -->
    <section class="section parchment">
      <div class="section-header">
        <h2>岗位卡牌</h2>
        <p>{{ recruit.jobSectionSubtitle }}</p>
      </div>

      <div class="grid-cards">
        <article
          v-for="p in filteredPositions"
          :key="p.id"
          class="job-card"
          @click="openJob(p)"
        >
          <div class="job-head">
            <div class="job-title">{{ p.title }}</div>
            <div class="job-dept">{{ p.dept }}</div>
          </div>

          <div class="job-tags">
            <span v-for="t in p.tags" :key="t" class="tag small">{{ t }}</span>
          </div>

          <div class="job-brief">{{ p.brief }}</div>

          <div class="job-cta">
            <span class="link">查看详情</span>
            <span class="arrow">›</span>
          </div>
        </article>
      </div>
    </section>

    <!-- 招新流程 -->
    <section class="section" ref="timelineEl">
      <div class="section-header">
        <h2 id="timeline">招新流程</h2>
        <p>{{ recruit.timelineSubtitle }}</p>
      </div>

      <div class="timeline parchment">
        <div v-for="(s, i) in steps" :key="s.title" class="step">
          <div class="step-index">{{ i + 1 }}</div>
          <div class="step-body">
            <div class="step-title">{{ s.title }}</div>
            <div class="step-desc">{{ s.desc }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- FAQ -->
    <section class="section parchment">
      <div class="section-header">
        <h2>常见问题</h2>
        <p>{{ recruit.faqSubtitle }}</p>
      </div>

      <div class="faq">
        <button
          v-for="(f, idx) in faqs"
          :key="f.q"
          class="faq-item"
          @click="toggleFaq(idx)"
          type="button"
        >
          <div class="faq-q">
            <span class="q">Q</span>
            <span class="text">{{ f.q }}</span>
            <span class="chev">{{ openFaq === idx ? "−" : "+" }}</span>
          </div>
          <div class="faq-a" v-if="openFaq === idx">
            <span class="a">A</span>
            <span class="text">{{ f.a }}</span>
          </div>
        </button>
      </div>
    </section>

    <!-- 报名卷轴 -->
    <section class="section" ref="applyEl">
      <div class="section-header">
        <h2 id="apply">领取报名卷轴</h2>
        <p>{{ recruit.applySubtitle }}</p>
      </div>

      <div class="apply parchment">
        <form class="form" @submit.prevent="submit">
          <div class="row">
            <label>
              <div class="label">姓名 / 昵称</div>
              <input v-model.trim="form.name" placeholder="例如：张三" />
            </label>

            <label>
              <div class="label">年级</div>
              <select v-model="form.grade">
                <option value="">请选择</option>
                <option>大一</option>
                <option>大二</option>
                <option>大三</option>
                <option>大四</option>
                <option>研究生</option>
              </select>
            </label>
          </div>

          <div class="row">
            <label>
              <div class="label">专业</div>
              <input v-model.trim="form.major" placeholder="例如：软件工程" />
            </label>

            <label>
              <div class="label">邮箱</div>
              <input v-model.trim="form.email" placeholder="name@example.com" />
            </label>
          </div>

          <div class="row">
            <label>
              <div class="label">选择主线（命运之卡）</div>
              <select v-model="form.deck">
                <option value="">请选择</option>
                <option v-for="d in decks" :key="d.key" :value="d.key">
                  {{ d.title }}
                </option>
              </select>
            </label>

            <label>
              <div class="label">意向岗位</div>
              <select v-model="form.positionId">
                <option value="">请选择</option>
                <option v-for="p in positions" :key="p.id" :value="p.id">
                  {{ p.title }}（{{ p.dept }}）
                </option>
              </select>
            </label>
          </div>

          <label>
            <div class="label">自我介绍 / 你为什么想来</div>
            <textarea
              v-model.trim="form.intro"
              rows="5"
              :placeholder="recruit.introPlaceholder"
            />
          </label>

          <div class="submit-row">
            <button class="btn" type="submit">🕯️ 盖章提交</button>
            <button class="btn ghost" type="button" @click="fillDemo">试填示例</button>
          </div>

          <div class="form-hint">{{ recruit.formHint }}</div>
        </form>
      </div>
    </section>

    <!-- 岗位详情弹窗 -->
    <div v-if="jobModal.open" class="modal-mask" @click.self="closeJob">
      <div class="modal parchment">
        <div class="modal-head">
          <div class="modal-title">{{ jobModal.job?.title }}</div>
          <button class="x" @click="closeJob" type="button">×</button>
        </div>

        <div class="modal-sub">
          <span class="pill">{{ jobModal.job?.dept }}</span>
          <span v-for="t in jobModal.job?.tags || []" :key="t" class="pill ghost">{{ t }}</span>
        </div>

        <div class="grid-2">
          <div class="panel">
            <div class="panel-title">你将负责</div>
            <ul class="bullets">
              <li v-for="x in jobModal.job?.resp || []" :key="x">{{ x }}</li>
            </ul>
          </div>
          <div class="panel">
            <div class="panel-title">我们希望你</div>
            <ul class="bullets">
              <li v-for="x in jobModal.job?.req || []" :key="x">{{ x }}</li>
            </ul>
          </div>
        </div>

        <div class="modal-actions">
          <button class="btn" @click="pickJobAndApply" type="button">就选这个，去报名</button>
          <button class="btn ghost" @click="closeJob" type="button">我再看看</button>
        </div>
      </div>
    </div>

    <!-- 项目预览弹窗 -->
    <div v-if="previewOpen" class="modal-mask" @click.self="previewOpen = false">
      <div class="modal parchment">
        <div class="modal-head">
          <div class="modal-title">我们正在做的事</div>
          <button class="x" @click="previewOpen = false" type="button">×</button>
        </div>

        <div class="panel">
          <div class="panel-title">例：项目清单</div>
          <ul class="bullets">
            <li v-for="x in recruit.previewList" :key="x">{{ x }}</li>
          </ul>
        </div>

        <div class="modal-actions">
          <button class="btn" @click="scrollTo('apply'); previewOpen = false" type="button">
            我想加入
          </button>
          <button class="btn ghost" @click="previewOpen = false" type="button">知道了</button>
        </div>
      </div>
    </div>

    <footer class="footer">
      <div class="footer-inner">
        <div class="footer-left">
          <div class="footer-title">{{ recruit.orgName }}</div>
          <div class="footer-desc">{{ recruit.footerText }}</div>
        </div>
        <div class="footer-right">
          <button class="btn ghost" @click="scrollTo('about')" type="button">返回顶部</button>
          <button class="btn" @click="scrollTo('apply')" type="button">立即报名</button>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watchEffect } from "vue";
import { useRoute, useRouter } from "vue-router";

// ✅ 你说你已经写了 RecruitConfig.js，就从这里导入
import { DEFAULT_RECRUIT, ORG_RECRUIT_MAP } from "@/config/RecruitConfig";

const route = useRoute();
const router = useRouter();

function goBack() {
  router.back();
}

// 关键：当前组织 id
const orgId = computed(() => Number(route.params.id));

// 合并默认配置 + 组织覆盖配置
function mergeRecruit(base, patch) {
  return {
    ...base,
    ...patch,
    decks: patch?.decks?.length ? patch.decks : base.decks,
    positions: patch?.positions?.length ? patch.positions : base.positions,
    steps: patch?.steps?.length ? patch.steps : base.steps,
    faqs: patch?.faqs?.length ? patch.faqs : base.faqs,
  };
}

// recruit 就是页面所有数据的唯一来源
const recruit = computed(() => {
  const patch = ORG_RECRUIT_MAP[orgId.value] || {};
  return mergeRecruit(DEFAULT_RECRUIT, patch);
});

// 页面里用到的集合
const decks = computed(() => recruit.value.decks || []);
const positions = computed(() => recruit.value.positions || []);
const steps = computed(() => recruit.value.steps || []);
const faqs = computed(() => recruit.value.faqs || []);

// 锚点
const aboutEl = ref(null);
const positionsEl = ref(null);
const timelineEl = ref(null);
const applyEl = ref(null);

function scrollTo(anchor) {
  const map = { about: aboutEl, positions: positionsEl, timeline: timelineEl, apply: applyEl };
  const el = map[anchor]?.value;
  if (el?.scrollIntoView) el.scrollIntoView({ behavior: "smooth", block: "start" });
}

// 卡牌选择
const activeDeck = ref("");

watchEffect(() => {
  const first = decks.value[0]?.key;
  if (!first) return;
  if (!activeDeck.value) activeDeck.value = first;
  if (!decks.value.some((d) => d.key === activeDeck.value)) activeDeck.value = first;
});

const deckDetail = computed(() => {
  const d = decks.value.find((x) => x.key === activeDeck.value) || decks.value[0] || {};
  return {
    title: d.title || "",
    long: d.long || "",
    doing: Array.isArray(d.doing) ? d.doing : [],
    learn: Array.isArray(d.learn) ? d.learn : [],
  };
});

// 岗位过滤：如果 position 有 deckKey，就按 activeDeck 过滤；没有 deckKey 视为通用
const filteredPositions = computed(() => {
  return (positions.value || []).filter((p) => !p.deckKey || p.deckKey === activeDeck.value);
});

// FAQ
const openFaq = ref(0);
function toggleFaq(i) {
  openFaq.value = openFaq.value === i ? -1 : i;
}

// 岗位弹窗
const jobModal = ref({ open: false, job: null });
function openJob(job) {
  jobModal.value = { open: true, job };
}
function closeJob() {
  jobModal.value = { open: false, job: null };
}
function pickJobAndApply() {
  if (jobModal.value.job) form.value.positionId = jobModal.value.job.id;
  closeJob();
  scrollTo("apply");
}

// 预览弹窗
const previewOpen = ref(false);
function openPreview() {
  previewOpen.value = true;
}

// 表单
const form = ref({
  name: "",
  grade: "",
  major: "",
  email: "",
  deck: "",
  positionId: "",
  intro: "",
});

function fillDemo() {
  form.value = {
    name: "Terisa_Z",
    grade: "大三",
    major: "软件工程",
    email: "terisa@example.com",
    deck: activeDeck.value,
    positionId: filteredPositions.value[0]?.id || positions.value[0]?.id || "",
    intro: "我希望在组织里稳定投入：能协作、能复盘、遇到问题不逃跑。希望学期末有可展示的成果。",
  };
}

function validateEmail(v) {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v);
}

function submit() {
  const v = form.value;
  if (!v.deck) v.deck = activeDeck.value;

  if (!v.name || !v.grade || !v.major || !v.email || !v.deck || !v.positionId || !v.intro) {
    alert("请把卷轴填完整：姓名/年级/专业/邮箱/方向/岗位/自我介绍都需要。");
    return;
  }
  if (!validateEmail(v.email)) {
    alert("邮箱格式不对哦。");
    return;
  }

  alert(`✅ 已盖章提交（演示版）\n组织：${recruit.value.orgName}\n方向：${v.deck}\n岗位：${v.positionId}`);
}

onMounted(() => {
  openFaq.value = 0;
});
</script>

<style scoped>
/* ========= 基础 ========= */
.recruit-page {
  min-height: 100vh;
  color: rgba(255, 245, 225, 0.92);
  background:
    radial-gradient(1200px 800px at 20% 10%, rgba(255, 210, 120, 0.10), transparent 60%),
    radial-gradient(900px 700px at 80% 20%, rgba(180, 40, 70, 0.16), transparent 60%),
    radial-gradient(1200px 900px at 50% 90%, rgba(20, 10, 20, 0.85), rgba(8, 6, 10, 0.95)),
    linear-gradient(180deg, rgba(20, 8, 12, 1), rgba(6, 6, 10, 1));
  position: relative;
  overflow-x: hidden;
}

/* 轻微“噪点/纸纹”效果（纯 CSS） */
.recruit-page::before {
  content: "";
  position: fixed;
  inset: 0;
  pointer-events: none;
  opacity: 0.10;
  background-image:
    repeating-linear-gradient(0deg, rgba(255,255,255,0.06), rgba(255,255,255,0.06) 1px, transparent 1px, transparent 3px),
    repeating-linear-gradient(90deg, rgba(255,255,255,0.04), rgba(255,255,255,0.04) 1px, transparent 1px, transparent 4px);
  mix-blend-mode: overlay;
}

/* ========= 顶栏 ========= */
.topbar {
  position: sticky;
  top: 0;
  z-index: 50;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 18px;
  backdrop-filter: blur(10px);
  background: linear-gradient(180deg, rgba(20, 10, 14, 0.85), rgba(10, 8, 12, 0.55));
  border-bottom: 1px solid rgba(255, 210, 120, 0.18);
}

.icon-btn {
  width: 42px;
  height: 38px;
  border-radius: 12px;
  border: 1px solid rgba(255, 210, 120, 0.22);
  background: rgba(0, 0, 0, 0.25);
  color: rgba(255, 245, 225, 0.92);
  font-size: 20px;
  cursor: pointer;
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 240px;
}

.logo {
  width: 40px;
  height: 40px;
  border-radius: 14px;
  display: grid;
  place-items: center;
  background: radial-gradient(circle at 30% 30%, rgba(255, 220, 150, 0.22), rgba(120, 60, 20, 0.25));
  border: 1px solid rgba(255, 210, 120, 0.30);
  box-shadow: 0 10px 25px rgba(0,0,0,0.35);
}

.titles .org-name {
  font-weight: 700;
  letter-spacing: 0.4px;
}

.titles .org-meta {
  font-size: 12px;
  opacity: 0.78;
  margin-top: 2px;
}

.topbar-right {
  display: flex;
  gap: 10px;
  align-items: center;
}

.btn {
  cursor: pointer;
  border: 1px solid rgba(255, 210, 120, 0.30);
  background: linear-gradient(180deg, rgba(255, 210, 120, 0.20), rgba(120, 60, 20, 0.18));
  color: rgba(255, 245, 225, 0.96);
  padding: 10px 14px;
  border-radius: 14px;
  box-shadow: 0 12px 26px rgba(0,0,0,0.35);
  transition: transform 160ms ease, box-shadow 160ms ease, filter 160ms ease;
}
.btn:hover {
  transform: translateY(-1px);
  filter: brightness(1.06);
  box-shadow: 0 16px 34px rgba(0,0,0,0.42);
}
.btn.ghost {
  background: rgba(0,0,0,0.22);
}

/* ========= 通用布局 ========= */
.section {
  max-width: 1120px;
  margin: 0 auto;
  padding: 40px 18px;
}
.section-header h2 {
  margin: 0 0 8px;
  font-size: 26px;
  letter-spacing: 0.6px;
}
.section-header p {
  margin: 0;
  opacity: 0.80;
  line-height: 1.7;
}

/* 羊皮纸容器（“苏丹 UI”关键） */
.parchment {
  background:
    radial-gradient(900px 400px at 20% 10%, rgba(255, 220, 150, 0.18), transparent 60%),
    radial-gradient(800px 380px at 80% 30%, rgba(255, 210, 120, 0.12), transparent 60%),
    linear-gradient(180deg, rgba(35, 20, 18, 0.70), rgba(18, 12, 14, 0.78));
  border: 1px solid rgba(255, 210, 120, 0.22);
  border-radius: 22px;
  box-shadow: 0 18px 50px rgba(0,0,0,0.45);
  position: relative;
}
.parchment::after {
  content: "";
  position: absolute;
  inset: 10px;
  border-radius: 16px;
  border: 1px dashed rgba(255, 210, 120, 0.18);
  pointer-events: none;
}

/* ========= HERO ========= */
.hero {
  max-width: 1120px;
  margin: 0 auto;
  padding: 26px 18px 8px;
  display: grid;
  gap: 18px;
  grid-template-columns: 1.4fr 0.8fr;
}
@media (max-width: 920px) {
  .hero { grid-template-columns: 1fr; }
}

.hero-card {
  padding: 26px 22px;
  border-radius: 26px;
  border: 1px solid rgba(255, 210, 120, 0.26);
  background:
    radial-gradient(700px 350px at 30% 0%, rgba(255, 210, 120, 0.16), transparent 60%),
    linear-gradient(180deg, rgba(40, 18, 18, 0.78), rgba(12, 10, 14, 0.88));
  box-shadow: 0 22px 60px rgba(0,0,0,0.55);
  position: relative;
}

.ornate-line {
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(255, 210, 120, 0.50), transparent);
  margin: 10px 0 14px;
}

.hero-kicker {
  font-size: 12px;
  opacity: 0.80;
  letter-spacing: 1.6px;
  text-transform: uppercase;
}

.hero-title {
  margin: 10px 0 10px;
  font-size: 34px;
  line-height: 1.18;
}
.gold {
  color: rgba(255, 210, 120, 0.95);
  text-shadow: 0 10px 30px rgba(255, 200, 120, 0.10);
}

.hero-subtitle {
  margin: 0;
  opacity: 0.86;
  line-height: 1.85;
}

.hero-actions {
  margin-top: 16px;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.seal-row {
  margin-top: 18px;
  display: flex;
  gap: 14px;
  align-items: center;
}
.seal {
  width: 54px;
  height: 54px;
  border-radius: 18px;
  background: radial-gradient(circle at 30% 30%, rgba(255, 210, 120, 0.35), rgba(120, 20, 30, 0.55));
  border: 1px solid rgba(255, 210, 120, 0.35);
  box-shadow: 0 18px 40px rgba(0,0,0,0.45);
  display: grid;
  place-items: center;
  transform: rotate(-6deg);
}
.seal-inner {
  font-weight: 800;
  font-size: 11px;
  opacity: 0.92;
  letter-spacing: 1px;
}
.seal-text .seal-title {
  font-weight: 700;
  margin-bottom: 2px;
}
.seal-text .seal-desc {
  opacity: 0.78;
  font-size: 13px;
}

.hero-side {
  display: grid;
  gap: 14px;
}
.stat-card, .quote-card {
  padding: 18px 18px;
  border-radius: 22px;
  border: 1px solid rgba(255, 210, 120, 0.18);
  background: rgba(0,0,0,0.26);
  box-shadow: 0 18px 40px rgba(0,0,0,0.40);
}
.stat-title {
  font-weight: 800;
  margin-bottom: 8px;
}
.stat-list {
  margin: 0;
  padding-left: 18px;
  line-height: 1.85;
  opacity: 0.86;
}
.quote-mark {
  font-size: 34px;
  opacity: 0.55;
  margin-bottom: 6px;
}
.quote {
  line-height: 1.85;
  opacity: 0.90;
}
.quote-from {
  margin-top: 10px;
  font-size: 12px;
  opacity: 0.70;
}

/* ========= 面板/网格 ========= */
.grid-3 {
  margin-top: 18px;
  display: grid;
  gap: 14px;
  grid-template-columns: repeat(3, 1fr);
}
.grid-2 {
  display: grid;
  gap: 14px;
  grid-template-columns: repeat(2, 1fr);
}
@media (max-width: 920px) {
  .grid-3 { grid-template-columns: 1fr; }
  .grid-2 { grid-template-columns: 1fr; }
}

.panel {
  padding: 16px 16px;
  border-radius: 18px;
  border: 1px solid rgba(255, 210, 120, 0.16);
  background: rgba(0,0,0,0.22);
}
.panel-title {
  font-weight: 800;
  margin-bottom: 8px;
}
.panel-body {
  opacity: 0.84;
  line-height: 1.85;
}

.bullets {
  margin: 0;
  padding-left: 18px;
  line-height: 1.85;
  opacity: 0.86;
}

/* ========= 卡组（四张命运卡） ========= */
.deck {
  margin-top: 18px;
  display: grid;
  gap: 14px;
  grid-template-columns: repeat(4, 1fr);
}
@media (max-width: 1050px) {
  .deck { grid-template-columns: repeat(2, 1fr); }
}
@media (max-width: 640px) {
  .deck { grid-template-columns: 1fr; }
}

.card {
  text-align: left;
  cursor: pointer;
  padding: 16px 16px;
  border-radius: 24px;
  border: 1px solid rgba(255, 210, 120, 0.20);
  background:
    radial-gradient(260px 140px at 30% 0%, rgba(255, 210, 120, 0.18), transparent 60%),
    linear-gradient(180deg, rgba(30, 16, 18, 0.74), rgba(10, 10, 14, 0.88));
  box-shadow: 0 20px 55px rgba(0,0,0,0.50);
  transition: transform 160ms ease, filter 160ms ease, border-color 160ms ease;
  position: relative;
}
.card:hover {
  transform: translateY(-2px) rotate(-0.2deg);
  filter: brightness(1.06);
  border-color: rgba(255, 210, 120, 0.34);
}
.card.active {
  border-color: rgba(255, 210, 120, 0.55);
  box-shadow:
    0 24px 65px rgba(0,0,0,0.58),
    0 0 0 1px rgba(255, 210, 120, 0.15) inset;
}

.card-top {
  display: flex;
  align-items: center;
  gap: 10px;
}
.card-emblem {
  width: 38px;
  height: 38px;
  border-radius: 14px;
  display: grid;
  place-items: center;
  border: 1px solid rgba(255, 210, 120, 0.22);
  background: rgba(0,0,0,0.25);
}
.card-title {
  font-weight: 900;
  letter-spacing: 0.4px;
}
.card-body {
  margin-top: 10px;
}
.card-tagline {
  font-size: 13px;
  opacity: 0.84;
  margin-bottom: 8px;
}
.card-desc {
  opacity: 0.78;
  line-height: 1.75;
  font-size: 13px;
}
.card-tags {
  margin-top: 10px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.tag {
  padding: 6px 10px;
  border-radius: 999px;
  border: 1px solid rgba(255, 210, 120, 0.18);
  background: rgba(0,0,0,0.18);
  font-size: 12px;
  opacity: 0.88;
}
.tag.small {
  padding: 5px 9px;
  font-size: 11px;
}
.card-foot {
  margin-top: 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  opacity: 0.75;
}
.card-foot-arrow {
  font-size: 18px;
}

/* 详情 */
.detail {
  margin-top: 16px;
  padding: 18px;
  border-radius: 22px;
}
.detail-head {
  margin-bottom: 12px;
}
.detail-title {
  font-size: 18px;
  font-weight: 900;
}
.detail-sub {
  opacity: 0.82;
  line-height: 1.85;
  margin-top: 6px;
}

/* ========= 岗位卡 ========= */
.grid-cards {
  margin-top: 18px;
  display: grid;
  gap: 14px;
  grid-template-columns: repeat(3, 1fr);
}
@media (max-width: 980px) {
  .grid-cards { grid-template-columns: repeat(2, 1fr); }
}
@media (max-width: 640px) {
  .grid-cards { grid-template-columns: 1fr; }
}

.job-card {
  padding: 16px 16px;
  border-radius: 24px;
  border: 1px solid rgba(255, 210, 120, 0.16);
  background: rgba(0,0,0,0.18);
  box-shadow: 0 18px 45px rgba(0,0,0,0.45);
  cursor: pointer;
  transition: transform 160ms ease, filter 160ms ease, border-color 160ms ease;
}
.job-card:hover {
  transform: translateY(-2px);
  filter: brightness(1.05);
  border-color: rgba(255, 210, 120, 0.30);
}
.job-head {
  display: flex;
  justify-content: space-between;
  gap: 10px;
}
.job-title {
  font-weight: 900;
  line-height: 1.2;
}
.job-dept {
  font-size: 12px;
  opacity: 0.72;
  white-space: nowrap;
}
.job-tags {
  margin: 10px 0 8px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.job-brief {
  opacity: 0.82;
  line-height: 1.75;
  min-height: 44px;
}
.job-cta {
  margin-top: 12px;
  display: flex;
  justify-content: space-between;
  opacity: 0.80;
}
.link {
  color: rgba(255, 210, 120, 0.92);
}
.arrow { font-size: 18px; }

/* ========= 流程 ========= */
.timeline {
  margin-top: 16px;
  padding: 18px;
  border-radius: 22px;
}
.step {
  display: flex;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid rgba(255, 210, 120, 0.14);
}
.step:last-child { border-bottom: none; }
.step-index {
  width: 34px;
  height: 34px;
  border-radius: 14px;
  display: grid;
  place-items: center;
  border: 1px solid rgba(255, 210, 120, 0.22);
  background: rgba(0,0,0,0.18);
  font-weight: 900;
}
.step-title {
  font-weight: 900;
  margin-bottom: 4px;
}
.step-desc {
  opacity: 0.82;
  line-height: 1.75;
}

/* ========= FAQ ========= */
.faq {
  margin-top: 16px;
  display: grid;
  gap: 10px;
}
.faq-item {
  text-align: left;
  padding: 14px 14px;
  border-radius: 20px;
  border: 1px solid rgba(255, 210, 120, 0.16);
  background: rgba(0,0,0,0.18);
  cursor: pointer;
}
.faq-q, .faq-a {
  display: flex;
  gap: 10px;
  align-items: baseline;
}
.faq-q .q, .faq-a .a {
  width: 22px;
  height: 22px;
  border-radius: 10px;
  display: grid;
  place-items: center;
  border: 1px solid rgba(255, 210, 120, 0.20);
  background: rgba(0,0,0,0.20);
  font-size: 12px;
  font-weight: 900;
  opacity: 0.90;
}
.faq-q .text { font-weight: 800; }
.faq-a { margin-top: 10px; opacity: 0.86; line-height: 1.75; }
.chev {
  margin-left: auto;
  opacity: 0.8;
  font-size: 18px;
}

/* ========= 表单 ========= */
.apply {
  margin-top: 16px;
  padding: 18px;
  border-radius: 22px;
}
.form label {
  display: block;
  margin-bottom: 12px;
}
.row {
  display: grid;
  gap: 12px;
  grid-template-columns: repeat(2, 1fr);
}
@media (max-width: 820px) {
  .row { grid-template-columns: 1fr; }
}

.label {
  font-weight: 800;
  margin-bottom: 6px;
  opacity: 0.92;
}

input, select, textarea {
  width: 100%;
  box-sizing: border-box;
  border-radius: 16px;
  border: 1px solid rgba(255, 210, 120, 0.18);
  background: rgba(0,0,0,0.22);
  color: rgba(255, 245, 225, 0.92);
  padding: 10px 12px;
  outline: none;
}
input:focus, select:focus, textarea:focus {
  border-color: rgba(255, 210, 120, 0.38);
  box-shadow: 0 0 0 3px rgba(255, 210, 120, 0.10);
}

.submit-row {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  margin-top: 8px;
}
.form-hint {
  margin-top: 10px;
  font-size: 12px;
  opacity: 0.74;
  line-height: 1.7;
}

/* ========= 弹窗 ========= */
.modal-mask {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.60);
  display: grid;
  place-items: center;
  z-index: 80;
  padding: 18px;
}
.modal {
  width: min(920px, 100%);
  padding: 18px;
  border-radius: 22px;
}
.modal-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 10px;
}
.modal-title {
  font-size: 18px;
  font-weight: 900;
}
.x {
  width: 38px;
  height: 38px;
  border-radius: 14px;
  border: 1px solid rgba(255, 210, 120, 0.22);
  background: rgba(0,0,0,0.22);
  color: rgba(255, 245, 225, 0.92);
  cursor: pointer;
  font-size: 20px;
}
.modal-sub {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
}
.pill {
  padding: 6px 10px;
  border-radius: 999px;
  border: 1px solid rgba(255, 210, 120, 0.18);
  background: rgba(0,0,0,0.18);
  font-size: 12px;
}
.pill.ghost {
  opacity: 0.85;
}
.modal-actions {
  margin-top: 14px;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  justify-content: flex-end;
}

/* ========= 页脚 ========= */
.footer {
  padding: 26px 18px 40px;
  border-top: 1px solid rgba(255, 210, 120, 0.14);
  background: rgba(0,0,0,0.10);
}
.footer-inner {
  max-width: 1120px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}
.footer-title {
  font-weight: 900;
}
.footer-desc {
  margin-top: 4px;
  opacity: 0.75;
  font-size: 13px;
}
</style>
