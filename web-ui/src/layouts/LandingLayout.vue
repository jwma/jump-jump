<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useTheme } from '@/composables/useTheme'
import { Link2, Languages, Sun, Moon, Monitor, Menu, X } from 'lucide-vue-next'

const { t, locale } = useI18n()
const auth = useAuthStore()
const { theme } = useTheme()
const scrolled = ref(false)
const langDropdownOpen = ref(false)
const themeDropdownOpen = ref(false)
const mobileMenuOpen = ref(false)

function setLocale(lang: string) {
  locale.value = lang
  localStorage.setItem('locale', lang)
  langDropdownOpen.value = false
}

function setTheme(value: 'light' | 'dark' | 'system') {
  theme.value = value
  themeDropdownOpen.value = false
}

function onScroll() {
  scrolled.value = window.scrollY > 16
}

function closeDropdowns() {
  langDropdownOpen.value = false
  themeDropdownOpen.value = false
}

function closeMobileMenu() {
  mobileMenuOpen.value = false
  closeDropdowns()
}

onMounted(() => {
  onScroll()
  window.addEventListener('scroll', onScroll, { passive: true })
  document.addEventListener('click', closeDropdowns)
})
onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
  document.removeEventListener('click', closeDropdowns)
})
</script>

<template>
  <div class="min-h-screen bg-white dark:bg-gray-950">
    <!-- Navbar -->
    <nav
      class="fixed inset-x-0 top-0 z-50 transition-all duration-200"
      :class="
        scrolled
          ? 'border-b border-gray-200 bg-white/80 shadow-sm backdrop-blur-md dark:border-gray-800 dark:bg-gray-950/80'
          : 'bg-transparent'
      "
    >
      <div class="mx-auto flex h-16 max-w-6xl items-center justify-between px-4 sm:px-6 lg:px-8">
        <router-link :to="{ name: 'landing' }" class="flex cursor-pointer items-center gap-2">
          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-600">
            <Link2 class="h-5 w-5 text-white" />
          </div>
          <span class="text-lg font-bold text-gray-900 dark:text-white">Jump Jump</span>
        </router-link>

        <!-- Desktop nav (sm+) -->
        <div class="hidden items-center gap-3 sm:flex">
          <!-- Theme toggle -->
          <div class="relative">
            <button
              class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-300"
              :aria-label="t('topbar.toggleTheme')"
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
                class="flex w-full items-center gap-2 px-3 py-2 text-sm text-gray-600 transition-colors hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-gray-800"
                :class="{ 'bg-gray-50 font-medium dark:bg-gray-800': theme === 'light' }"
                @click="setTheme('light')"
              >
                <Sun class="h-4 w-4" /> {{ t('topbar.light') }}
              </button>
              <button
                class="flex w-full items-center gap-2 px-3 py-2 text-sm text-gray-600 transition-colors hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-gray-800"
                :class="{ 'bg-gray-50 font-medium dark:bg-gray-800': theme === 'dark' }"
                @click="setTheme('dark')"
              >
                <Moon class="h-4 w-4" /> {{ t('topbar.dark') }}
              </button>
              <button
                class="flex w-full items-center gap-2 px-3 py-2 text-sm text-gray-600 transition-colors hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-gray-800"
                :class="{ 'bg-gray-50 font-medium dark:bg-gray-800': theme === 'system' }"
                @click="setTheme('system')"
              >
                <Monitor class="h-4 w-4" /> {{ t('topbar.system') }}
              </button>
            </div>
          </div>

          <!-- Language switcher -->
          <div class="relative">
            <button
              class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-300"
              :aria-label="t('topbar.switchLanguage')"
              @click.stop="langDropdownOpen = !langDropdownOpen"
            >
              <Languages class="h-5 w-5" />
            </button>
            <div
              v-if="langDropdownOpen"
              class="absolute right-0 top-full z-50 mt-2 w-28 overflow-hidden rounded-xl border border-gray-200 bg-white shadow-lg dark:border-gray-700 dark:bg-gray-900"
              @click.stop
            >
              <button
                class="flex w-full items-center gap-2 px-3 py-2 text-sm text-gray-600 transition-colors hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-gray-800"
                :class="{ 'bg-gray-50 font-medium dark:bg-gray-800': locale === 'en' }"
                @click="setLocale('en')"
              >
                English
              </button>
              <button
                class="flex w-full items-center gap-2 px-3 py-2 text-sm text-gray-600 transition-colors hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-gray-800"
                :class="{ 'bg-gray-50 font-medium dark:bg-gray-800': locale === 'zh' }"
                @click="setLocale('zh')"
              >
                中文
              </button>
            </div>
          </div>

          <router-link
            v-if="auth.isLoggedIn"
            :to="{ name: 'dashboard' }"
            class="cursor-pointer rounded-lg px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-800"
          >
            {{ t('landing.nav.dashboard') }}
          </router-link>
          <template v-else>
            <router-link
              :to="{ name: 'login' }"
              class="cursor-pointer rounded-lg px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-800"
            >
              {{ t('landing.nav.login') }}
            </router-link>
            <router-link
              :to="{ name: 'login' }"
              class="cursor-pointer rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700"
            >
              {{ t('landing.nav.getStarted') }}
            </router-link>
          </template>
        </div>

        <!-- Mobile hamburger button -->
        <button
          class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-800 sm:hidden"
          :aria-label="mobileMenuOpen ? 'Close menu' : 'Open menu'"
          @click="mobileMenuOpen = !mobileMenuOpen"
        >
          <X v-if="mobileMenuOpen" class="h-6 w-6" />
          <Menu v-else class="h-6 w-6" />
        </button>
      </div>

      <!-- Mobile menu dropdown -->
      <div
        v-if="mobileMenuOpen"
        class="border-t border-gray-200 bg-white px-4 pb-4 pt-2 dark:border-gray-800 dark:bg-gray-950 sm:hidden"
      >
        <div class="flex items-center justify-between">
          <!-- Theme toggle -->
          <span class="text-sm text-gray-500 dark:text-gray-400">
            {{ t('topbar.toggleTheme') }}
          </span>
          <div class="flex items-center gap-1">
            <button
              v-for="opt in ['light', 'dark', 'system'] as const"
              :key="opt"
              class="rounded-lg px-3 py-1.5 text-sm transition-colors"
              :class="
                theme === opt
                  ? 'bg-blue-50 font-medium text-blue-600 dark:bg-blue-950 dark:text-blue-400'
                  : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-800'
              "
              @click="setTheme(opt)"
            >
              {{ t(`topbar.${opt}`) }}
            </button>
          </div>
        </div>

        <div class="mt-3 flex items-center justify-between">
          <!-- Language switcher -->
          <span class="text-sm text-gray-500 dark:text-gray-400">
            {{ t('topbar.switchLanguage') }}
          </span>
          <div class="flex items-center gap-1">
            <button
              v-for="loc in ['en', 'zh']"
              :key="loc"
              class="rounded-lg px-3 py-1.5 text-sm transition-colors"
              :class="
                locale === loc
                  ? 'bg-blue-50 font-medium text-blue-600 dark:bg-blue-950 dark:text-blue-400'
                  : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-800'
              "
              @click="setLocale(loc)"
            >
              {{ loc === 'en' ? 'English' : '中文' }}
            </button>
          </div>
        </div>

        <div class="mt-4 flex flex-col gap-2 border-t border-gray-200 pt-4 dark:border-gray-800">
          <router-link
            v-if="auth.isLoggedIn"
            :to="{ name: 'dashboard' }"
            class="cursor-pointer rounded-lg px-4 py-2.5 text-center text-sm font-medium text-gray-700 transition-colors hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-800"
            @click="closeMobileMenu"
          >
            {{ t('landing.nav.dashboard') }}
          </router-link>
          <template v-else>
            <router-link
              :to="{ name: 'login' }"
              class="cursor-pointer rounded-lg px-4 py-2.5 text-center text-sm font-medium text-blue-600 transition-colors hover:bg-blue-50 dark:text-blue-400 dark:hover:bg-blue-950"
              @click="closeMobileMenu"
            >
              {{ t('landing.nav.login') }}
            </router-link>
            <router-link
              :to="{ name: 'login' }"
              class="cursor-pointer rounded-lg bg-blue-600 px-4 py-2.5 text-center text-sm font-medium text-white transition-colors hover:bg-blue-700"
              @click="closeMobileMenu"
            >
              {{ t('landing.nav.getStarted') }}
            </router-link>
          </template>
        </div>
      </div>
    </nav>

    <!-- Main -->
    <main>
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <!-- Footer -->
    <footer class="border-t border-gray-200 bg-gray-50 dark:border-gray-800 dark:bg-gray-900">
      <div class="mx-auto max-w-6xl px-4 py-12 sm:px-6 lg:px-8">
        <div class="flex flex-col items-center gap-4 sm:flex-row sm:justify-between">
          <div class="flex items-center gap-2">
            <div class="flex h-7 w-7 items-center justify-center rounded-md bg-blue-600">
              <Link2 class="h-4 w-4 text-white" />
            </div>
            <span class="text-sm text-gray-600 dark:text-gray-400">
              {{ t('landing.footer.description') }}
            </span>
          </div>
          <div class="flex items-center gap-6">
            <a
              :href="t('landing.footer.githubUrl')"
              target="_blank"
              rel="noopener noreferrer"
              class="cursor-pointer text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300"
            >
              {{ t('landing.footer.github') }}
            </a>
            <span class="text-sm text-gray-400 dark:text-gray-500">
              {{ t('landing.footer.copyright') }}
            </span>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
