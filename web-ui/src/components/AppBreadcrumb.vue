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
  'short-link-create': 'Create',
  'short-link-detail': 'Detail',
  'short-link-edit': 'Edit',
  config: 'System Config',
  tenants: 'Tenants',
  preferences: 'Preferences',
  'change-password': 'Change Password',
}

const parentRoutes: Record<string, { label: string; to: string }> = {
  'short-link-create': { label: 'Short Links', to: '/short-links' },
  'short-link-detail': { label: 'Short Links', to: '/short-links' },
  'short-link-edit': { label: 'Short Links', to: '/short-links' },
  'tenant-create': { label: 'Tenants', to: '/tenants' },
  'tenant-detail': { label: 'Tenants', to: '/tenants' },
  'super-tenant-detail': { label: 'Tenants', to: '/super/tenants' },
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
      const parent = parentRoutes[r.name as string]
      if (parent) {
        crumbs.push({ label: parent.label, to: parent.to })
      }
      crumbs.push({ label })
    }
  }
  return crumbs
})
</script>

<template>
  <nav v-if="items.length > 0" class="flex items-center gap-1 text-sm">
    <router-link to="/" class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300">
      <Home class="h-4 w-4" />
    </router-link>
    <template v-for="(item, i) in items" :key="i">
      <ChevronRight class="h-3.5 w-3.5 text-gray-300 dark:text-gray-600" />
      <router-link v-if="item.to" class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200" :to="item.to">
        {{ item.label }}
      </router-link>
      <span v-else class="font-medium text-gray-900 dark:text-white">{{ item.label }}</span>
    </template>
  </nav>
</template>
