<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()

async function handleLogout() {
  await auth.logout()
  router.push({ name: 'login' })
}
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <header class="border-b bg-white">
      <div class="mx-auto flex h-14 max-w-7xl items-center justify-between px-4">
        <div class="flex items-center gap-6">
          <h1 class="text-lg font-semibold">Jump Jump</h1>
          <nav class="flex gap-4">
            <RouterLink to="/" class="text-sm text-gray-600 hover:text-gray-900"
              >Dashboard</RouterLink
            >
            <RouterLink to="/short-links" class="text-sm text-gray-600 hover:text-gray-900"
              >Short Links</RouterLink
            >
            <RouterLink to="/config" class="text-sm text-gray-600 hover:text-gray-900"
              >Config</RouterLink
            >
          </nav>
        </div>
        <div class="flex items-center gap-4">
          <span class="text-sm text-gray-500">{{ auth.username }}</span>
          <button class="text-sm text-gray-500 hover:text-gray-900" @click="handleLogout">
            Logout
          </button>
        </div>
      </div>
    </header>
    <main class="mx-auto max-w-7xl px-4 py-6">
      <RouterView />
    </main>
  </div>
</template>
