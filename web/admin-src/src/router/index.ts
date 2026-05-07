import { createRouter, createWebHistory } from 'vue-router'
import { getToken, getTenantId } from '@/utils/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/login/index.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/select-tenant',
      name: 'SelectTenant',
      component: () => import('@/views/select-tenant/index.vue'),
      meta: { requiresAuth: true, requiresTenant: false },
    },
    {
      path: '/invitations',
      name: 'Invitations',
      component: () => import('@/views/invitations/index.vue'),
      meta: { requiresAuth: true, requiresTenant: false },
    },
    {
      path: '/',
      component: () => import('@/components/layout/AppLayout.vue'),
      meta: { requiresAuth: true, requiresTenant: true },
      children: [
        {
          path: '',
          name: 'Dashboard',
          component: () => import('@/views/dashboard/index.vue'),
        },
        {
          path: 'change-password',
          name: 'ChangePassword',
          component: () => import('@/views/change-password/index.vue'),
        },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

router.beforeEach((to) => {
  const token = getToken()

  if (to.meta.requiresAuth && !token) {
    return { name: 'Login', query: { redirect: to.fullPath } }
  }

  if (to.name === 'Login' && token) {
    return { name: 'SelectTenant' }
  }

  if (to.meta.requiresTenant && token && !getTenantId()) {
    return { name: 'SelectTenant' }
  }
})

export default router
