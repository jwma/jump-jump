<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useLayoutStore } from '@/stores/layout'
import { useTheme } from '@/composables/useTheme'
import GlobalSearch from '@/components/GlobalSearch.vue'
import {
  Menu,
  Bell,
  LogOut,
  KeyRound,
  ChevronDown,
  Clock,
  X,
  Sun,
  Moon,
  Monitor,
} from 'lucide-vue-next'

const auth = useAuthStore()
const layout = useLayoutStore()
const router = useRouter()
const { theme } = useTheme()

const userDropdownOpen = ref(false)
const notifDropdownOpen = ref(false)
const themeDropdownOpen = ref(false)

const pendingInvitations = computed(() =>
  auth.invitations.filter((i) => i.status === 'pending'),
)

async function handleLogout() {
  userDropdownOpen.value = false
  await auth.logout()
  router.push({ name: 'login' })
}

function closeDropdowns(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (!target.closest('.topbar-dropdown')) {
    userDropdownOpen.value = false
    notifDropdownOpen.value = false
    themeDropdownOpen.value = false
  }
}

onMounted(() => document.addEventListener('click', closeDropdowns))
onBeforeUnmount(() => document.removeEventListener('click', closeDropdowns))
</script>

<template>
  <header class="flex h-14 shrink-0 items-center justify-between border-b border-gray-200 bg-white px-4 dark:border-gray-700 dark:bg-gray-900">
    <div class="flex items-center gap-3">
      <!-- Mobile menu toggle -->
      <button
        class="rounded p-1.5 text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-800 lg:hidden"
        @click="layout.toggleMobileMenu"
      >
        <Menu class="h-5 w-5" />
      </button>
      <!-- Breadcrumb slot -->
      <slot name="breadcrumb" />
    </div>

    <div class="flex items-center gap-3">
      <!-- Global Search -->
      <GlobalSearch />
    </div>

    <div class="flex items-center gap-2">
      <!-- Theme toggle -->
      <div class="topbar-dropdown relative">
        <button
          class="rounded p-1.5 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600 dark:hover:bg-gray-800 dark:hover:text-gray-300"
          aria-label="Toggle theme"
          @click.stop="themeDropdownOpen = !themeDropdownOpen"
        >
          <Sun v-if="theme === 'light'" class="h-5 w-5" />
          <Moon v-else-if="theme === 'dark'" class="h-5 w-5" />
          <Monitor v-else class="h-5 w-5" />
        </button>
        <div
          v-if="themeDropdownOpen"
          class="absolute right-0 top-full z-50 mt-2 w-36 overflow-hidden rounded-xl border border-gray-200 bg-white shadow-lg dark:border-gray-700 dark:bg-gray-900"
          @click.stop
        >
          <button
            class="flex w-full items-center gap-2 px-3 py-2 text-sm text-gray-600 hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-gray-800"
            :class="{ 'bg-gray-50 dark:bg-gray-800': theme === 'light' }"
            @click="theme = 'light'; themeDropdownOpen = false"
          >
            <Sun class="h-4 w-4" /> Light
          </button>
          <button
            class="flex w-full items-center gap-2 px-3 py-2 text-sm text-gray-600 hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-gray-800"
            :class="{ 'bg-gray-50 dark:bg-gray-800': theme === 'dark' }"
            @click="theme = 'dark'; themeDropdownOpen = false"
          >
            <Moon class="h-4 w-4" /> Dark
          </button>
          <button
            class="flex w-full items-center gap-2 px-3 py-2 text-sm text-gray-600 hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-gray-800"
            :class="{ 'bg-gray-50 dark:bg-gray-800': theme === 'system' }"
            @click="theme = 'system'; themeDropdownOpen = false"
          >
            <Monitor class="h-4 w-4" /> System
          </button>
        </div>
      </div>

      <!-- Notification center -->
      <div class="topbar-dropdown relative">
        <button
          class="relative rounded p-1.5 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600 dark:hover:bg-gray-800 dark:hover:text-gray-300"
          aria-label="Open notifications"
          aria-haspopup="menu"
          :aria-expanded="notifDropdownOpen"
          @click.stop="notifDropdownOpen = !notifDropdownOpen"
        >
          <Bell class="h-5 w-5" />
          <span
            v-if="auth.pendingInvitationCount > 0"
            class="absolute -right-0.5 -top-0.5 flex h-4 items-center justify-center rounded-full bg-red-500 px-1 text-[10px] font-bold text-white"
          >
            {{ auth.pendingInvitationCount > 9 ? '9+' : auth.pendingInvitationCount }}
          </span>
        </button>

        <div
          v-if="notifDropdownOpen"
          class="absolute right-0 top-full z-50 mt-2 w-80 overflow-hidden rounded-xl border border-gray-200 bg-white shadow-lg dark:border-gray-700 dark:bg-gray-900"
          @click.stop
        >
          <div class="flex items-center justify-between border-b border-gray-200 dark:border-gray-700 px-4 py-3">
            <h3 class="text-sm font-semibold text-gray-900 dark:text-white">Notifications</h3>
            <button
              class="rounded p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-600 dark:hover:bg-gray-800 dark:hover:text-gray-300"
              aria-label="Close notifications"
              @click="notifDropdownOpen = false"
            >
              <X class="h-4 w-4" />
            </button>
          </div>

          <div class="max-h-80 overflow-y-auto">
            <template v-if="pendingInvitations.length > 0">
              <div
                v-for="inv in pendingInvitations"
                :key="inv.id"
                class="flex items-start gap-3 border-b border-gray-50 dark:border-gray-800 px-4 py-3 transition-colors hover:bg-gray-50 dark:hover:bg-gray-800"
              >
                <div class="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-amber-100">
                  <Bell class="h-4 w-4 text-amber-600" />
                </div>
                <div class="min-w-0 flex-1">
                  <p class="text-sm text-gray-700 dark:text-gray-300">
                    Invited to <span class="font-medium">{{ inv.tenantName }}</span>
                  </p>
                  <p class="mt-0.5 text-xs text-gray-400 dark:text-gray-500">
                    <Clock class="mr-0.5 inline h-3 w-3" />
                    Invited by {{ inv.inviterUsername }}
                  </p>
                </div>
              </div>
            </template>
            <div v-else class="px-4 py-8 text-center">
              <Bell class="mx-auto h-8 w-8 text-gray-200 dark:text-gray-700" />
              <p class="mt-2 text-sm text-gray-400 dark:text-gray-500">No new notifications</p>
            </div>
          </div>

          <router-link
            v-if="pendingInvitations.length > 0"
            :to="{ name: 'invitations' }"
            class="block border-t border-gray-200 dark:border-gray-700 px-4 py-2.5 text-center text-sm font-medium text-blue-600 transition-colors hover:bg-gray-50 dark:hover:bg-gray-800"
            @click="notifDropdownOpen = false"
          >
            View all invitations
          </router-link>
        </div>
      </div>

      <!-- User avatar dropdown -->
      <div class="topbar-dropdown relative">
        <button
          class="flex items-center gap-2 rounded-lg px-2 py-1.5 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800"
          aria-label="User menu"
          aria-haspopup="menu"
          :aria-expanded="userDropdownOpen"
          @click.stop="userDropdownOpen = !userDropdownOpen"
        >
          <div
            class="flex h-8 w-8 items-center justify-center rounded-full bg-blue-100 text-sm font-medium text-blue-700"
          >
            {{ auth.username.charAt(0).toUpperCase() }}
          </div>
          <span class="hidden text-sm font-medium text-gray-700 dark:text-gray-300 sm:inline">{{ auth.username }}</span>
          <ChevronDown class="hidden h-3.5 w-3.5 text-gray-400 sm:block" :class="{ 'rotate-180': userDropdownOpen }" />
        </button>

        <div
          v-if="userDropdownOpen"
          class="absolute right-0 top-full z-50 mt-2 w-56 overflow-hidden rounded-xl border border-gray-200 bg-white shadow-lg dark:border-gray-700 dark:bg-gray-900"
          @click.stop
        >
          <div class="border-b border-gray-200 dark:border-gray-700 px-4 py-3">
            <p class="text-sm font-medium text-gray-900 dark:text-white">{{ auth.username }}</p>
            <p v-if="auth.currentTenant" class="mt-0.5 text-xs text-gray-400 dark:text-gray-500">{{ auth.currentTenant.tenantName }}</p>
          </div>
          <div class="py-1">
            <button
              class="flex w-full items-center gap-2.5 px-4 py-2 text-sm text-gray-600 transition-colors hover:bg-gray-50 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-white"
              @click="userDropdownOpen = false; router.push({ name: 'change-password' })"
            >
              <KeyRound class="h-4 w-4 text-gray-400" />
              Change Password
            </button>
            <button
              class="flex w-full items-center gap-2.5 px-4 py-2 text-sm text-red-600 transition-colors hover:bg-red-50 dark:hover:bg-red-900/20"
              @click="handleLogout"
            >
              <LogOut class="h-4 w-4" />
              Logout
            </button>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>
