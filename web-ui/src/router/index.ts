import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { UserRole } from '@/types/api'

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
          meta: { title: 'Dashboard' },
        },
        {
          path: 'short-links',
          name: 'short-links',
          component: () => import('@/pages/ShortLinksPage.vue'),
          meta: { title: 'Short Links' },
        },
        {
          path: 'short-links/:id',
          name: 'short-link-detail',
          component: () => import('@/pages/ShortLinkDetailPage.vue'),
          meta: { title: 'Detail' },
        },
        {
          path: 'config',
          name: 'config',
          component: () => import('@/pages/ConfigPage.vue'),
          meta: { title: 'System Config', requiredRole: UserRole.Admin },
        },
        {
          path: 'tenants',
          name: 'tenants',
          component: () => import('@/pages/TenantsPage.vue'),
          meta: { title: 'Tenants', requiredRole: UserRole.Admin },
        },
        {
          path: 'preferences',
          name: 'preferences',
          component: () => import('@/pages/PreferencesPage.vue'),
          meta: { title: 'Preferences' },
        },
        {
          path: 'change-password',
          name: 'change-password',
          component: () => import('@/pages/ChangePasswordPage.vue'),
          meta: { title: 'Change Password' },
        },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/pages/NotFoundPage.vue'),
      meta: { requiresAuth: false },
    },
  ],
  scrollBehavior() {
    return { top: 0 }
  },
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

  if (to.meta.requiredRole && auth.user?.role !== to.meta.requiredRole) {
    return { name: 'dashboard' }
  }

  if (to.name === 'login' && auth.isLoggedIn) {
    return { name: 'dashboard' }
  }
})

export default router
