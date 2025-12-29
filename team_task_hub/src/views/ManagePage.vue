<template>
  <div class="page">
    <header class="topbar">
      <button class="icon-btn" @click="goBack" aria-label="back">‹</button>

      <div class="brand">
        <div class="logo">🛡️</div>
        <div class="titles">
          <div class="org-name">组织管理</div>
          <div class="org-meta">成员 {{ members.length }} · 我的角色 {{ myRole }}</div>
        </div>
      </div>

      <div class="topbar-right">
        <OrgMemberSearchBar
          v-model="q"
          :loading="loading"
          @search="onSearchClick"
          @show-all="showAllMembers"
        />

        <button class="btn" v-if="canManage" @click="openAddMember">
          ＋ 添加成员
        </button>
      </div>
    </header>

    <main class="main">
      <section class="panel">
        <div class="panel-title">成员列表</div>

        <OrgMemberTable
          :members="members"
          :loading="loading"
          :error="error"
          :can-manage="canManage"
          :my-user-id="myUserId"
          @remove="removeMember"
          @set-role="({ userId, role }) => setRole(userId, role)"
        />
      </section>
    </main>

    <AddMemberModal
      :visible="showAdd"
      :value="addValue"
      @close="showAdd = false"
      @confirm="confirmAddFromModal"
    />
  </div>
</template>

<script setup>
import axios from "axios";
import { ref, computed, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";

// ✅如果你的组件不在 src/components/org，请改成你真实路径
import OrgMemberSearchBar from "@/components/OrgMemberSearchBar.vue";
import OrgMemberTable from "@/components/OrgMemberTable.vue";
import AddMemberModal from "@/components/AddMemberModal.vue";

const route = useRoute();
const router = useRouter();
const API_BASE = "http://localhost:8080";

const orgId = computed(() => route.params.id);

// ------- 登录用户 id（从 localStorage 取，兼容多种结构）-------
const cuRaw = localStorage.getItem("currentUser");
const cu = cuRaw ? JSON.parse(cuRaw) : {};
const actual = cu.data || cu.user || cu;
const myUserId = Number(actual.id || actual.user_id || actual.userId || actual.uid || 0);

// ------- 页面状态 -------
const myRole = ref("Member");
const canManage = ref(false);

const members = ref([]);
const loading = ref(false);
const error = ref("");

const q = ref("");

// ------- 缓存每个 user 的角色（刷新页面会重新查，但同次页面内不会重复请求）-------
const roleCache = new Map(); // userId -> "Creator"/"Admin"/"Member"

// ------- 基础 -------
function goBack() {
  if (window.history.length > 1) return router.back();
  router.push({ name: "OrgPage", params: { id: route.params.id }, query: route.query });
}

function authHeaders() {
  const token = localStorage.getItem("token");
  return token ? { Authorization: `Bearer ${token}` } : {};
}

// ------- 获取某个用户的 role（做法2：search 后逐个补角色）-------
async function fetchUserRole(userId) {
  if (roleCache.has(userId)) return roleCache.get(userId);

  const res = await axios.get(
    `${API_BASE}/api/organization/${orgId.value}/users/${userId}/role`,
    { headers: authHeaders() }
  );

  const role = String(res.data?.data?.role || "Member").trim();
  roleCache.set(userId, role);
  return role;
}

// ------- 获取我的 role & 权限 -------
async function fetchMyRole() {
  if (!myUserId) {
    myRole.value = "Member";
    canManage.value = false;
    return;
  }

  try {
    const res = await axios.get(
      `${API_BASE}/api/organization/${orgId.value}/users/${myUserId}/role`,
      { headers: authHeaders() }
    );

    const role = String(res.data?.data?.role || "Member").trim();
    myRole.value = role;

    const r = role.toLowerCase();
    canManage.value = ["creator", "admin"].includes(r);
  } catch (e) {
    myRole.value = "Member";
    canManage.value = false;
  }
}

// ------- 搜索成员（空关键字 => 后端返回全量）-------
async function searchMembers(keyword = "") {
  const kw = String(keyword).trim();

  loading.value = true;
  error.value = "";
  try {
    const params = {};
    if (kw) params.keyword = kw; // ✅你后端参数名是 keyword；空则不传，后端返回全量

    const res = await axios.get(
      `${API_BASE}/api/organization/${orgId.value}/users/search`,
      { headers: authHeaders(), params }
    );

    const users = res.data?.data?.users || [];

    // 先占位
    members.value = users.map(u => ({
      userId: u.id,
      name: u.username,
      email: u.email,
      avatarUrl: u.avatar_url,
      role: "Member",
      joinedAt: "",
    }));

    // 并发补全 role（做法2）
    await Promise.all(
      members.value.map(async (m) => {
        try {
          const r = await fetchUserRole(m.userId);
          console.log("role for", m.userId, m.name, "=>", r);
          m.role = r;
        } catch (e) {
          console.log("role failed for", m.userId, m.name, e?.response?.status, e?.response?.data || e.message);
          m.role = "Member";
        }
      })
    );
  } catch (e) {
    error.value = e?.response?.data?.message || e.message || "unknown";
    members.value = [];
  } finally {
    loading.value = false;
  }
}

// ------- UI 事件 -------
async function onSearchClick() {
  await searchMembers(q.value);
}

async function showAllMembers() {
  q.value = "";
  await searchMembers("");
}

// ------- 添加成员（你后端接口未给齐，先占位）-------
const showAdd = ref(false);
const addValue = ref("");

function openAddMember() {
  addValue.value = "";
  showAdd.value = true;
}

async function confirmAddFromModal(val) {
  addValue.value = val;
  await confirmAdd();
}

async function confirmAdd() {
  const v = addValue.value.trim();
  if (!v) return alert("请输入 userId 或邮箱");

  // TODO：等你给“添加成员”真实接口再接
  showAdd.value = false;
  await searchMembers(q.value);
}

// ------- 删除成员（等你给接口）-------
async function removeMember(userId) {
  alert("删除成员 API 还没对齐后端");
}

// ------- 设角色：目前只实现“提拔管理员”（PATCH）-------
async function setRole(userId, role) {
  // 你表格组件会传 role='admin' 或 'member'，目前后端只给了提拔管理员接口
  const r = String(role).toLowerCase();
  if (r !== "admin") {
    alert("后端还没给取消管理员接口，暂时只能提拔为管理员");
    return;
  }

  if (!confirm("确定将该成员提拔为管理员？")) return;

  try {
    await axios.patch(
      `${API_BASE}/api/organization/${orgId.value}/admin/${userId}`,
      null, // ✅PATCH 没有 body 用 null 更稳
      { headers: authHeaders() }
    );

    // 本地更新 + 缓存更新（刷新前立即生效）
    roleCache.set(userId, "Admin");
    const m = members.value.find(x => x.userId === userId);
    if (m) m.role = "Admin";

    alert("成员已成功提拔为管理员");
  } catch (e) {
    alert(e?.response?.data?.message || e.message || "提拔失败");
  }
}

// ------- 页面加载：先拿我的权限，再拉成员列表 -------
onMounted(async () => {
  await fetchMyRole();
  await searchMembers("");
});
</script>

<style scoped>
:global(:root) {
  --hb: 80px;
}

.page {
  min-height: 100vh;
  background: var(--bg, #0f1419);
  color: var(--text, #fff);
  padding-top: var(--hb, 80px);
  box-sizing: border-box;
}

.topbar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.12);
  position: sticky;
  top: var(--hb, 80px);
  z-index: 10;
  backdrop-filter: blur(10px);
  background: rgba(0, 0, 0, 0.2);
}

.topbar-right {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 10px;
}

.icon-btn {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.14);
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
  cursor: pointer;
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo {
  width: 44px;
  height: 44px;
  border-radius: 14px;
  display: grid;
  place-items: center;
  background: rgba(255, 255, 255, 0.1);
}

.org-name {
  font-weight: 900;
  letter-spacing: 0.06em;
}

.org-meta {
  font-size: 12px;
  opacity: 0.7;
  margin-top: 2px;
}

.main {
  padding: 16px;
  max-width: 1100px;
  margin: 0 auto;
  box-sizing: border-box;
}

.panel {
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 18px;
  padding: 14px;
  background: rgba(255, 255, 255, 0.05);
}

.panel-title {
  font-weight: 900;
  letter-spacing: 0.08em;
  margin-bottom: 10px;
}

.btn {
  padding: 10px 12px;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.14);
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  cursor: pointer;
  white-space: nowrap;
}
</style>
