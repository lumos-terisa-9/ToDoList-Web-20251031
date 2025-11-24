import { createRouter, createWebHistory } from 'vue-router'
import PersonPage from '../views/PersonPage.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'homepage',
      component: () => import('../views/HomePage.vue'), // 新的首页
    },
    {
      path: '/personpage',
      name: 'personpage',
      component: PersonPage, // 原来的日历页面
    },
    {
      path: '/orgmap',
      name: 'orgmap',
      component: () => import('@/views/MapPage.vue')
    }
  ],
})

export default router
