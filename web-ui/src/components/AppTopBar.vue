<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { useLayoutStore } from '@/stores/layout'
import { Menu } from 'lucide-vue-next'

const auth = useAuthStore()
const layout = useLayoutStore()
</script>

<template>
  <header class="flex h-14 shrink-0 items-center justify-between border-b bg-white px-4">
    <div class="flex items-center gap-3">
      <!-- Mobile menu toggle -->
      <button
        class="rounded p-1.5 text-gray-500 hover:bg-gray-100 lg:hidden"
        @click="layout.toggleMobileMenu"
      >
        <Menu class="h-5 w-5" />
      </button>
      <!-- Breadcrumb slot -->
      <slot name="breadcrumb" />
    </div>

    <div class="flex items-center gap-3">
      <div class="flex items-center gap-2">
        <div
          class="flex h-8 w-8 items-center justify-center rounded-full bg-blue-100 text-sm font-medium text-blue-700"
        >
          {{ auth.username.charAt(0).toUpperCase() }}
        </div>
        <span class="hidden text-sm font-medium text-gray-700 sm:inline">{{ auth.username }}</span>
      </div>
    </div>
  </header>
</template>
