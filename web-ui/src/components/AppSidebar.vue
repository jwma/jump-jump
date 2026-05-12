<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useLayoutStore } from '@/stores/layout'
import { UserRole } from '@/types/api'
import {
  LayoutDashboard,
  Link,
  Settings,
  Building2,
  KeyRound,
  LogOut,
  ChevronLeft,
  X,
  Bell,
  ChevronDown,
  ChevronRight,
  Shield,
  Users,
  UserPlus,
} from 'lucide-vue-next'

const { t } = useI18n()
const auth = useAuthStore()
const layout = useLayoutStore()
const route = useRoute()
const router = useRouter()
const tenantDropdownOpen = ref(false)

const menuItems = computed(() => {
  const items: { icon: unknown; label: string; to: { name: string } }[] = []

  if (isSuperRoute.value) {
    items.push(
      { icon: LayoutDashboard, label: t('nav.dashboard'), to: { name: 'super-dashboard' } },
      { icon: Users, label: t('super.users.title'), to: { name: 'super-users' } },
      { icon: Building2, label: t('breadcrumb.tenants'), to: { name: 'super-tenants' } },
    )
  } else {
    items.push(
      { icon: LayoutDashboard, label: t('nav.dashboard'), to: { name: 'dashboard' } },
      { icon: Link, label: t('nav.shortLinks'), to: { name: 'short-links' } },
      { icon: UserPlus, label: t('nav.members'), to: { name: 'members' } },
    )

    if (auth.user?.role === UserRole.Admin) {
      items.push({ icon: Settings, label: t('nav.settings'), to: { name: 'settings' } })
    }
  }

  items.push({ icon: KeyRound, label: t('nav.changePassword'), to: { name: 'change-password' } })

  return items
})

const isSuperRoute = computed(() => String(route.path).startsWith('/super'))

const showTenantSwitcher = computed(() => !isSuperRoute.value)

function isActive(to: { name: string }) {
  if (route.name === to.name) return true
  if (to.name === 'short-links' && String(route.name).startsWith('short-link')) return true
  if (to.name === 'tenants' && String(route.name).startsWith('tenant')) return true
  if (to.name === 'super-tenants' && String(route.name).startsWith('super-tenant')) return true
  return false
}

async function handleTenantSwitch(tenantId: string) {
  tenantDropdownOpen.value = false
  if (tenantId !== auth.currentTenantId) {
    auth.selectTenant(tenantId)
    try {
      await auth.fetchUser()
    } catch {
      auth.clearTenant()
      router.push({ name: 'select-tenant' })
      return
    }
    router.push({ name: 'dashboard' })
  }
}

async function handleLogout() {
  await auth.logout()
  router.push({ name: 'login' })
}

function closeTenantDropdown() {
  tenantDropdownOpen.value = false
  layout.closeMobileMenu()
}

function toggleSuperRoute() {
  if (isSuperRoute.value) {
    router.push({ name: auth.currentTenantId ? 'dashboard' : 'select-tenant' })
  } else {
    router.push({ name: 'super-dashboard' })
  }
}
</script>

<template>
  <!-- Mobile overlay -->
  <div
    v-if="layout.mobileMenuOpen"
    class="fixed inset-0 z-40 bg-black/50 lg:hidden"
    @click="layout.closeMobileMenu"
  />

  <aside
    :class="[
      'fixed left-0 top-0 z-50 flex h-full flex-col border-r border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-900 transition-all duration-300 lg:relative lg:z-auto',
      layout.sidebarCollapsed ? 'w-16' : 'w-60',
      layout.mobileMenuOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0',
    ]"
  >
    <!-- Header -->
    <div
      :class="[
        'flex h-14 items-center border-b border-gray-200 dark:border-gray-700',
        layout.sidebarCollapsed
          ? 'justify-between px-3 lg:justify-center lg:px-0'
          : 'justify-between px-3',
      ]"
    >
      <router-link
        to="/"
        class="flex items-center gap-2 overflow-hidden"
        @click="layout.closeMobileMenu"
      >
        <div
          class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-blue-600 text-sm font-bold text-white"
        >
          JJ
        </div>
        <span
          v-show="!layout.sidebarCollapsed"
          class="whitespace-nowrap text-base font-semibold text-gray-900 dark:text-white"
        >
          {{ t('auth.appName') }}
        </span>
      </router-link>

      <!-- Mobile close -->
      <button
        class="rounded p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 lg:hidden"
        @click="layout.closeMobileMenu"
      >
        <X class="h-5 w-5" />
      </button>
      <!-- Desktop collapse (only when expanded) -->
      <button
        v-show="!layout.sidebarCollapsed"
        class="hidden rounded p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 lg:block"
        @click="layout.toggleSidebar"
      >
        <ChevronLeft class="h-4 w-4 transition-transform" />
      </button>
    </div>

    <!-- Expand toggle (only when collapsed, desktop) -->
    <div v-if="layout.sidebarCollapsed" class="sidebar-tooltip-wrapper">
      <button
        class="absolute -right-3 top-4 z-10 hidden h-6 w-6 items-center justify-center rounded-full border border-gray-200 bg-white shadow-sm text-gray-400 hover:text-gray-600 dark:border-gray-600 dark:bg-gray-800 dark:hover:text-gray-300 lg:flex"
        :aria-label="t('nav.expandSidebar')"
        @click="layout.toggleSidebar"
      >
        <ChevronLeft class="h-3 w-3 rotate-180" />
      </button>
      <span class="sidebar-tooltip" style="left: calc(100% + 16px); top: 16px; transform: none">{{
        t('nav.expandSidebar')
      }}</span>
    </div>

    <!-- Tenant switcher (hidden on super admin routes) -->
    <div
      v-if="!layout.sidebarCollapsed && showTenantSwitcher"
      class="border-b border-gray-200 dark:border-gray-700 px-3 py-2"
    >
      <button
        class="flex w-full items-center justify-between rounded-lg px-2 py-1.5 text-left text-sm hover:bg-gray-100 dark:hover:bg-gray-800"
        @click="tenantDropdownOpen = !tenantDropdownOpen"
      >
        <div class="min-w-0 flex-1">
          <div class="truncate text-xs text-gray-400 dark:text-gray-500">
            {{ t('nav.currentTenant') }}
          </div>
          <div class="truncate font-medium text-gray-700 dark:text-gray-300">
            {{
              auth.currentTenant?.tenantName ||
              (auth.isSuper ? t('nav.superAdmin') : t('nav.notSelected'))
            }}
          </div>
        </div>
        <ChevronDown class="h-4 w-4 shrink-0 text-gray-400" />
      </button>

      <div v-if="tenantDropdownOpen" class="mt-1 space-y-0.5">
        <button
          v-for="tenant in auth.tenants"
          :key="tenant.tenantId"
          class="flex w-full items-center justify-between rounded-md px-2 py-1.5 text-sm hover:bg-gray-100 dark:hover:bg-gray-800"
          :class="
            auth.currentTenantId === tenant.tenantId
              ? 'bg-blue-50 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
              : 'text-gray-600 dark:text-gray-400'
          "
          @click="handleTenantSwitch(tenant.tenantId)"
        >
          <span class="truncate">{{ tenant.tenantName }}</span>
          <span class="text-xs text-gray-400 dark:text-gray-500">{{ tenant.role }}</span>
        </button>
        <router-link
          :to="{ name: 'select-tenant' }"
          class="flex items-center gap-1 rounded-md px-2 py-1.5 text-xs text-gray-400 hover:bg-gray-50 hover:text-gray-600 dark:hover:bg-gray-800 dark:hover:text-gray-300"
          @click="closeTenantDropdown"
        >
          <ChevronRight class="h-3 w-3" />
          {{ t('nav.manageTenants') }}
        </router-link>
      </div>
    </div>

    <!-- Super admin toggle -->
    <div
      v-if="auth.isSuper"
      class="sidebar-tooltip-wrapper border-b border-gray-200 dark:border-gray-700 px-2 py-2"
    >
      <button
        class="flex w-full items-center gap-2 rounded-lg px-2 py-1.5 text-sm transition-colors"
        :class="
          isSuperRoute
            ? 'bg-indigo-50 text-indigo-700 dark:bg-indigo-900/30 dark:text-indigo-400'
            : 'text-gray-500 hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-200'
        "
        @click="toggleSuperRoute"
      >
        <Shield class="h-4 w-4 shrink-0" />
        <span v-show="!layout.sidebarCollapsed">{{
          isSuperRoute ? t('nav.backToApp') : t('nav.superAdmin')
        }}</span>
      </button>
      <span v-if="layout.sidebarCollapsed" class="sidebar-tooltip">{{
        isSuperRoute ? t('nav.backToApp') : t('nav.superAdmin')
      }}</span>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 overflow-y-auto px-2 py-3">
      <ul class="space-y-1">
        <li v-for="item in menuItems" :key="item.label" class="sidebar-tooltip-wrapper">
          <router-link
            :to="item.to"
            :aria-label="layout.sidebarCollapsed ? item.label : undefined"
            :class="[
              'flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors',
              isActive(item.to)
                ? 'bg-blue-50 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
                : 'text-gray-600 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-white',
            ]"
            @click="layout.closeMobileMenu"
          >
            <component :is="item.icon" class="h-5 w-5 shrink-0" />
            <span v-show="!layout.sidebarCollapsed">{{ item.label }}</span>
          </router-link>
          <span v-if="layout.sidebarCollapsed" class="sidebar-tooltip">{{ item.label }}</span>
        </li>
      </ul>
    </nav>

    <!-- Invitation badge + Logout -->
    <div class="border-t border-gray-200 dark:border-gray-700 px-2 py-3">
      <div v-if="auth.pendingInvitationCount > 0" class="sidebar-tooltip-wrapper">
        <router-link
          :to="{ name: 'invitations' }"
          class="mb-1 flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-amber-600 transition-colors hover:bg-amber-50 dark:hover:bg-amber-900/20"
          @click="layout.closeMobileMenu"
        >
          <span class="relative">
            <Bell class="h-5 w-5 shrink-0" />
            <span
              class="absolute -right-1.5 -top-1.5 flex h-4 min-w-4 items-center justify-center rounded-full bg-red-500 px-1 text-[10px] font-bold text-white"
            >
              {{ auth.pendingInvitationCount > 9 ? '9+' : auth.pendingInvitationCount }}
            </span>
          </span>
          <span v-show="!layout.sidebarCollapsed">{{ t('nav.invitations') }}</span>
        </router-link>
        <span v-if="layout.sidebarCollapsed" class="sidebar-tooltip">{{
          t('nav.invitations')
        }}</span>
      </div>
      <div class="sidebar-tooltip-wrapper">
        <button
          class="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-gray-600 transition-colors hover:bg-red-50 hover:text-red-600 dark:text-gray-400 dark:hover:bg-red-900/20 dark:hover:text-red-400"
          @click="handleLogout"
        >
          <LogOut class="h-5 w-5 shrink-0" />
          <span v-show="!layout.sidebarCollapsed">{{ t('nav.logout') }}</span>
        </button>
        <span v-if="layout.sidebarCollapsed" class="sidebar-tooltip">{{ t('nav.logout') }}</span>
      </div>
    </div>
  </aside>
</template>

<style scoped>
.sidebar-tooltip-wrapper {
  position: relative;
}

.sidebar-tooltip {
  position: absolute;
  left: calc(100% + 8px);
  top: 50%;
  transform: translateY(-50%);
  white-space: nowrap;
  border-radius: 6px;
  background: #1f2937;
  color: #fff;
  font-size: 12px;
  line-height: 1;
  padding: 6px 10px;
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.15s ease;
  z-index: 100;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.sidebar-tooltip::before {
  content: '';
  position: absolute;
  right: 100%;
  top: 50%;
  transform: translateY(-50%);
  border: 5px solid transparent;
  border-right-color: #1f2937;
}

.sidebar-tooltip-wrapper:hover .sidebar-tooltip,
.sidebar-tooltip-wrapper:focus-within .sidebar-tooltip {
  opacity: 1;
}
</style>
