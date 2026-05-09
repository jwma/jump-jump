<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { ChevronRight, Home } from 'lucide-vue-next'

const route = useRoute()

interface BreadcrumbItem {
  label: string
  to?: string
}

function routeTitle(name: string): string {
  const map: Record<string, string> = {
    dashboard: t('breadcrumb.dashboard'),
    'short-links': t('breadcrumb.shortLinks'),
    'short-link-create': t('breadcrumb.create'),
    'short-link-detail': t('breadcrumb.detail'),
    'short-link-edit': t('breadcrumb.edit'),
    config: t('settings.title'),
    tenants: t('breadcrumb.tenants'),
    preferences: t('breadcrumb.preferences'),
    'change-password': t('breadcrumb.changePassword'),
  }
  return map[name] || name
}

function parentRoute(name: string): { label: string; to: string } | undefined {
  const map: Record<string, { label: string; to: string }> = {
    'short-link-create': { label: t('breadcrumb.shortLinks'), to: '/short-links' },
    'short-link-detail': { label: t('breadcrumb.shortLinks'), to: '/short-links' },
    'short-link-edit': { label: t('breadcrumb.shortLinks'), to: '/short-links' },
    'tenant-create': { label: t('breadcrumb.tenants'), to: '/tenants' },
    'tenant-detail': { label: t('breadcrumb.tenants'), to: '/tenants' },
    'super-tenant-detail': { label: t('breadcrumb.tenants'), to: '/super/tenants' },
  }
  return map[name]
}

const items = computed<BreadcrumbItem[]>(() => {
  const matched = route.matched.filter((r) => r.meta?.title || r.name)
  if (matched.length === 0) return []

  const crumbs: BreadcrumbItem[] = []
  for (let i = 0; i < matched.length; i++) {
    const r = matched[i]
    const label = (r.meta.title as string) || routeTitle(r.name as string) || (r.name as string)
    if (i < matched.length - 1) {
      crumbs.push({ label, to: r.path })
    } else {
      const parent = parentRoute(r.name as string)
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
