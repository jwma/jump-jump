<script setup lang="ts">
import AppSidebar from '@/components/AppSidebar.vue'
import AppTopBar from '@/components/AppTopBar.vue'
import AppBreadcrumb from '@/components/AppBreadcrumb.vue'
import ShortcutHelp from '@/components/ShortcutHelp.vue'
import { useKeyboardShortcuts } from '@/composables/useKeyboardShortcuts'

useKeyboardShortcuts()
</script>

<template>
  <div class="flex h-screen overflow-hidden bg-gray-50">
    <AppSidebar />

    <div class="flex flex-1 flex-col overflow-hidden">
      <AppTopBar>
        <template #breadcrumb>
          <AppBreadcrumb />
        </template>
      </AppTopBar>

      <main class="flex-1 overflow-y-auto p-4 md:p-6">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>

    <ShortcutHelp />
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
