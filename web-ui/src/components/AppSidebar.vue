<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useLayoutStore } from '@/stores/layout'
import { UserRole } from '@/types/api'
import {
  LayoutDashboard,
  Link,
  Settings,
  Building2,
  UserCog,
  KeyRound,
  LogOut,
  ChevronLeft,
  X,
} from 'lucide-vue-next'

const auth = useAuthStore()
const layout = useLayoutStore()
const route = useRoute()
const router = useRouter()

const menuItems = computed(() => {
  const items = [
    { icon: LayoutDashboard, label: 'Dashboard', to: { name: 'dashboard' } },
    { icon: Link, label: 'Short Links', to: { name: 'short-links' } },
  ]

  if (auth.user?.role === UserRole.Admin) {
    items.push(
      { icon: Settings, label: 'System Config', to: { name: 'config' } },
      { icon: Building2, label: 'Tenants', to: { name: 'tenants' } },
    )
  }

  items.push(
    { icon: UserCog, label: 'Preferences', to: { name: 'preferences' } },
    { icon: KeyRound, label: 'Change Password', to: { name: 'change-password' } },
  )

  return items
})

function isActive(to: { name: string }) {
  if (route.name === to.name) return true
  if (to.name === 'short-links' && String(route.name).startsWith('short-link')) return true
  return false
}

async function handleLogout() {
  await auth.logout()
  router.push({ name: 'login' })
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
      'fixed left-0 top-0 z-50 flex h-full flex-col border-r bg-white transition-all duration-300 lg:relative lg:z-auto',
      layout.sidebarCollapsed ? 'w-16' : 'w-60',
      layout.mobileMenuOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0',
    ]"
  >
    <!-- Header -->
    <div class="flex h-14 items-center justify-between border-b px-3">
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
          class="whitespace-nowrap text-base font-semibold text-gray-900"
        >
          Jump Jump
        </span>
      </router-link>

      <!-- Mobile close -->
      <button
        class="rounded p-1 text-gray-400 hover:text-gray-600 lg:hidden"
        @click="layout.closeMobileMenu"
      >
        <X class="h-5 w-5" />
      </button>
      <!-- Desktop collapse -->
      <button
        class="hidden rounded p-1 text-gray-400 hover:text-gray-600 lg:block"
        @click="layout.toggleSidebar"
      >
        <ChevronLeft
          :class="['h-4 w-4 transition-transform', layout.sidebarCollapsed && 'rotate-180']"
        />
      </button>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 overflow-y-auto px-2 py-3">
      <ul class="space-y-1">
        <li v-for="item in menuItems" :key="item.label">
          <router-link
            :to="item.to"
            :class="[
              'flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors',
              isActive(item.to)
                ? 'bg-blue-50 text-blue-700'
                : 'text-gray-600 hover:bg-gray-100 hover:text-gray-900',
            ]"
            :title="layout.sidebarCollapsed ? item.label : undefined"
            @click="layout.closeMobileMenu"
          >
            <component :is="item.icon" class="h-5 w-5 shrink-0" />
            <span v-show="!layout.sidebarCollapsed">{{ item.label }}</span>
          </router-link>
        </li>
      </ul>
    </nav>

    <!-- Logout -->
    <div class="border-t px-2 py-3">
      <button
        class="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-gray-600 transition-colors hover:bg-red-50 hover:text-red-600"
        @click="handleLogout"
      >
        <LogOut class="h-5 w-5 shrink-0" />
        <span v-show="!layout.sidebarCollapsed">Logout</span>
      </button>
    </div>
  </aside>
</template>
