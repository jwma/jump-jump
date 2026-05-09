<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Home } from 'lucide-vue-next'

const router = useRouter()
</script>

<template>
  <div class="flex min-h-[60vh] flex-col items-center justify-center text-center">
    <h1 class="text-7xl font-bold text-gray-200 dark:text-gray-700">404</h1>
    <p class="mt-4 text-lg text-gray-500 dark:text-gray-400">Page not found</p>
    <p class="mt-1 text-sm text-gray-400 dark:text-gray-500">
      The page you are looking for does not exist or has been moved.
    </p>
    <button
      class="mt-6 inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700"
      @click="router.push({ name: 'dashboard' })"
    >
      <Home class="h-4 w-4" />
      Back to Dashboard
    </button>
  </div>
</template>
