<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { ChevronRight, Home } from 'lucide-vue-next'

const route = useRoute()

interface BreadcrumbItem {
  label: string
  to?: string
}

const routeTitles: Record<string, string> = {
  dashboard: 'Dashboard',
  'short-links': 'Short Links',
  'short-link-detail': 'Detail',
  config: 'System Config',
  tenants: 'Tenants',
  preferences: 'Preferences',
  'change-password': 'Change Password',
}

const items = computed<BreadcrumbItem[]>(() => {
  const matched = route.matched.filter((r) => r.meta?.title || r.name)
  if (matched.length === 0) return []

  const crumbs: BreadcrumbItem[] = []
  for (let i = 0; i < matched.length; i++) {
    const r = matched[i]
    const label = (r.meta.title as string) || routeTitles[r.name as string] || (r.name as string)
    if (i < matched.length - 1) {
      crumbs.push({ label, to: r.path })
    } else {
      crumbs.push({ label })
    }
  }
  return crumbs
})
</script>

<template>
  <nav v-if="items.length > 0" class="flex items-center gap-1 text-sm">
    <router-link to="/" class="text-gray-400 hover:text-gray-600">
      <Home class="h-4 w-4" />
    </router-link>
    <template v-for="(item, i) in items" :key="i">
      <ChevronRight class="h-3.5 w-3.5 text-gray-300" />
      <router-link v-if="item.to" class="text-gray-500 hover:text-gray-700" :to="item.to">
        {{ item.label }}
      </router-link>
      <span v-else class="font-medium text-gray-900">{{ item.label }}</span>
    </template>
  </nav>
</template>
