import { createRouter, createWebHistory } from "vue-router";
import PersonPage from "../views/PersonPage.vue";

// 统一用懒加载（推荐），也可以改成静态 import
const HomePage = () => import("../views/HomePage.vue");
const MapPage = () => import("@/views/MapPage.vue");
const OrgPage = () => import("@/views/OrgPage.vue");

const routes = [
  {
    path: "/",
    name: "homepage",
    component: HomePage,
  },
  {
    path: "/personpage",
    name: "personpage",
    component: PersonPage,
  },
  {
    path: "/orgmap",
    name: "orgmap",
    component: MapPage,
  },
  {
    path: "/org/:id",
    name: "Org",
    component: OrgPage,
    props: true,
  },
  {
    path: "/org/:id/info",
    name: "OrgInfo",
    component: () => import("@/views/OrgInfoPage.vue"),
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

router.beforeEach(async (to, from, next) => {
  // 如果访问首页且本地有token，则尝试自动跳转到个人页面
  if (to.name === "homepage") {
    const token = localStorage.getItem("token");

    if (token) {
      try {
        const API_BASE = "http://localhost:8080/api";
        let pureToken = token;

        // 解析token格式
        if (token.startsWith("{")) {
          try {
            const tokenData = JSON.parse(token);
            pureToken =
              tokenData.data?.access_token ||
              tokenData.access_token ||
              tokenData.token ||
              token;
          } catch {
            pureToken = token;
          }
        }

        const response = await fetch(`${API_BASE}/auth/me`, {
          method: "GET",
          headers: {
            Authorization: `Bearer ${pureToken}`,
            "Content-Type": "application/json",
          },
        });

        if (response.ok) {
          next("/personpage");
          return;
        }
      } catch (error) {
        console.log("Token验证失败，停留在首页");
      }
    }
  }

  next();
});

export default router;
