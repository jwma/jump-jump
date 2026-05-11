import { createRouter, createWebHistory } from 'vue-router'
import NProgress from 'nprogress'
import { useAuthStore } from '@/stores/auth'
import { UserRole } from '@/types/api'

NProgress.configure({ showSpinner: false })

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
      path: '/select-tenant',
      name: 'select-tenant',
      component: () => import('@/pages/SelectTenantPage.vue'),
      meta: { requiresAuth: true, requiresTenant: false },
    },
    {
      path: '/super',
      component: () => import('@/layouts/DefaultLayout.vue'),
      meta: { requiresAuth: true, requiresSuper: true },
      redirect: { name: 'super-dashboard' },
      children: [
        {
          path: 'dashboard',
          name: 'super-dashboard',
          component: () => import('@/pages/super/SuperDashboardPage.vue'),
          meta: { title: 'Super Dashboard' },
        },
        {
          path: 'users',
          name: 'super-users',
          component: () => import('@/pages/super/SuperUsersPage.vue'),
          meta: { title: 'User Management' },
        },
        {
          path: 'tenants',
          name: 'super-tenants',
          component: () => import('@/pages/super/SuperTenantsPage.vue'),
          meta: { title: 'Tenant Management' },
        },
        {
          path: 'tenants/:id',
          name: 'super-tenant-detail',
          component: () => import('@/pages/super/SuperTenantDetailPage.vue'),
          meta: { title: 'Tenant Detail' },
        },
      ],
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
          meta: { title: 'Dashboard', requiresTenant: true },
        },
        {
          path: 'short-links',
          name: 'short-links',
          component: () => import('@/pages/ShortLinksPage.vue'),
          meta: { title: 'Short Links', requiresTenant: true },
        },
        {
          path: 'short-links/create',
          name: 'short-link-create',
          component: () => import('@/pages/ShortLinkCreatePage.vue'),
          meta: { title: 'Create Short Link', requiresTenant: true },
        },
        {
          path: 'short-links/:id',
          name: 'short-link-detail',
          component: () => import('@/pages/ShortLinkDetailPage.vue'),
          meta: { title: 'Detail', requiresTenant: true },
        },
        {
          path: 'short-links/:id/edit',
          name: 'short-link-edit',
          component: () => import('@/pages/ShortLinkEditPage.vue'),
          meta: { title: 'Edit Short Link', requiresTenant: true },
        },
        {
          path: 'members',
          name: 'members',
          component: () => import('@/pages/MembersPage.vue'),
          meta: { title: 'Members', requiresTenant: true },
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('@/pages/SettingsPage.vue'),
          meta: { title: 'Settings', requiredRole: UserRole.Admin, requiresTenant: true },
        },
        {
          path: 'tenants',
          name: 'tenants',
          component: () => import('@/pages/TenantsPage.vue'),
          meta: { title: 'Tenants', requiredRole: UserRole.Admin, requiresTenant: true },
        },
        {
          path: 'tenants/create',
          name: 'tenant-create',
          component: () => import('@/pages/TenantCreatePage.vue'),
          meta: { title: 'Create Tenant', requiredRole: UserRole.Admin, requiresTenant: true },
        },
        {
          path: 'tenants/:id',
          name: 'tenant-detail',
          component: () => import('@/pages/TenantDetailPage.vue'),
          meta: { title: 'Tenant Detail', requiredRole: UserRole.Admin, requiresTenant: true },
        },
        {
          path: 'preferences',
          name: 'preferences',
          component: () => import('@/pages/PreferencesPage.vue'),
        },
        {
          path: 'change-password',
          name: 'change-password',
          component: () => import('@/pages/ChangePasswordPage.vue'),
          meta: { title: 'Change Password' },
        },
        {
          path: 'invitations',
          name: 'invitations',
          component: () => import('@/pages/InvitationsPage.vue'),
          meta: { title: 'Invitations' },
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
  NProgress.start()
  const auth = useAuthStore()

  if (to.meta.requiresAuth !== false && !auth.isLoggedIn) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  if (auth.isLoggedIn && !auth.authUser) {
    try {
      await auth.fetchAuthInfo()
    } catch {
      auth.clearAuth()
      return { name: 'login', query: { redirect: to.fullPath } }
    }
  }

  if (auth.isLoggedIn && !auth.user && auth.currentTenantId) {
    try {
      await auth.fetchUser()
    } catch {
      auth.clearTenant()
      return { name: 'select-tenant' }
    }
  }

  // Fetch invitations once per session
  if (auth.isLoggedIn && !auth.invitationsLoaded) {
    auth.fetchInvitations().catch(() => {})
  }

  if (to.meta.requiresTenant && !auth.currentTenantId && !auth.isSuper) {
    return { name: 'select-tenant' }
  }

  if (to.meta.requiredRole && auth.user?.role !== to.meta.requiredRole) {
    return { name: 'dashboard' }
  }

  if (to.meta.requiresSuper && !auth.isSuper) {
    return { name: auth.currentTenantId ? 'dashboard' : 'select-tenant' }
  }

  if (to.name === 'login' && auth.isLoggedIn) {
    return { name: auth.currentTenantId ? 'dashboard' : 'select-tenant' }
  }
})

router.afterEach(() => {
  NProgress.done()
})

router.onError(() => {
  NProgress.done()
})

export default router
