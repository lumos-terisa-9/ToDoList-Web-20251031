import { createRouter, createWebHistory } from 'vue-router'
import PersonPage from '../views/PersonPage.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'homepage',
      component: () => import('../views/HomePage.vue'),
    },
    {
      path: '/personpage',
      name: 'personpage',
      component: PersonPage,
    },
    {
      path: '/orgmap',
      name: 'orgmap',
      component: () => import('@/views/MapPage.vue')
    },
    {
      path: '/org/:id',
      name: 'OrgPage', // 你 MapPage 里之前用 name:"Org" 就能匹配
      component: () => import('@/views/OrgPage.vue'),
      props: true,
    },
    {
      path: '/org/:id/info',
      name: 'OrgInfo',
      component: () => import('@/views/OrgInfoPage.vue'),
    },
    {
      path: '/org/:id/activities',
      name: 'ActivityPage',
      component: () => import('@/views/ActivityPage.vue'),
      props: true,
    },
    {
      // 添加 ActivityDetail 路由
      path: '/org/:orgId/activity/:activityId',
      name: 'ActivityDetail',
      component: () => import('@/components/ActivityDetail.vue'),
      props: true,
    }
  ],
})

router.beforeEach(async (to, from, next) => {
  // 如果访问首页且本地有token，则尝试自动跳转到个人页面
  if (to.name === 'homepage') {
    const token = localStorage.getItem('token')

    if (token) {
      try {
        const API_BASE = 'http://localhost:8080/api'
        let pureToken = token

        // 解析token格式
        if (token.startsWith('{')) {
          try {
            const tokenData = JSON.parse(token)
            pureToken = tokenData.data?.access_token || tokenData.access_token || tokenData.token || token
          } catch {
            pureToken = token
          }
        }

        const response = await fetch(`${API_BASE}/auth/me`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${pureToken}`,
            'Content-Type': 'application/json'
          }
        })

        if (response.ok) {
          next('/personpage')
          return
        }
        // eslint-disable-next-line no-unused-vars
      } catch (error) {
        console.log('Token验证失败，停留在首页')
      }
    }
  }

  // 其他情况正常导航，，
  next()
})

export default router
