import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/pages/LoginPage.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/',
      component: () => import('@/layouts/DefaultLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'dashboard',
          component: () => import('@/pages/DashboardPage.vue'),
        },
        {
          path: 'short-links',
          name: 'short-links',
          component: () => import('@/pages/ShortLinksPage.vue'),
        },
        {
          path: 'short-links/:id',
          name: 'short-link-detail',
          component: () => import('@/pages/ShortLinkDetailPage.vue'),
        },
        {
          path: 'config',
          name: 'config',
          component: () => import('@/pages/ConfigPage.vue'),
        },
      ],
    },
  ],
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()

  if (to.meta.requiresAuth !== false && !auth.isLoggedIn) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  if (auth.isLoggedIn && !auth.user) {
    try {
      await auth.fetchUser()
    } catch {
      auth.clearAuth()
      return { name: 'login', query: { redirect: to.fullPath } }
    }
  }

  if (to.name === 'login' && auth.isLoggedIn) {
    return { name: 'dashboard' }
  }
})

export default router
