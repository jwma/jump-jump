<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { listTenants } from '@/api/tenant'
import type { Tenant } from '@/types/api'
import { Plus, Search, Building2, Eye } from 'lucide-vue-next'

defineOptions({ name: 'TenantsPage' })

const router = useRouter()
const { t } = useI18n()

const tenants = ref<Tenant[]>([])
const loading = ref(false)
const search = ref('')

const filteredTenants = computed(() => {
  if (!search.value) return tenants.value
  const q = search.value.toLowerCase()
  return tenants.value.filter(
    (tItem) =>
      tItem.name.toLowerCase().includes(q) ||
      tItem.slug.toLowerCase().includes(q) ||
      tItem.id.toLowerCase().includes(q),
  )
})

async function fetchTenants() {
  loading.value = true
  try {
    const data = await listTenants()
    tenants.value = data || []
  } catch {
    // error handled by interceptor
  } finally {
    loading.value = false
  }
}

function formatDate(d: string) {
  return new Date(d).toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

onMounted(fetchTenants)
</script>

<template>
  <div>
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-xl font-semibold text-gray-900 dark:text-white">{{ t('tenants.title') }}</h1>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          {{ t('tenants.count', { count: tenants.length, suffix: tenants.length !== 1 ? 's' : '' }) }}
        </p>
      </div>
      <router-link
        :to="{ name: 'tenant-create' }"
        class="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700"
      >
        <Plus class="h-4 w-4" />
        {{ t('tenants.createTenant') }}
      </router-link>
    </div>

    <!-- Search -->
    <div class="mt-4 flex flex-col gap-3 rounded-lg border bg-white dark:bg-gray-900 p-4 sm:flex-row sm:items-center">
      <div class="relative flex-1">
        <Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400 dark:text-gray-500" />
        <input
          v-model="search"
          type="text"
          :placeholder="t('tenants.searchPlaceholder')"
          class="w-full rounded-md border border-gray-300 dark:border-gray-600 py-2 pl-9 pr-3 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
        />
      </div>
    </div>

    <!-- Table -->
    <div class="mt-4 overflow-hidden rounded-lg border bg-white dark:bg-gray-900 dark:bg-gray-900">
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr
              class="border-b bg-gray-50 dark:bg-gray-800/50 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400 dark:text-gray-500"
            >
              <th class="px-4 py-3">{{ t('tenants.name') }}</th>
              <th class="px-4 py-3">{{ t('tenants.slug') }}</th>
              <th class="px-4 py-3">{{ t('tenants.status') }}</th>
              <th class="hidden px-4 py-3 lg:table-cell">{{ t('tenants.created') }}</th>
              <th class="px-4 py-3 text-right">{{ t('tenants.actions') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y">
            <tr v-if="loading">
              <td colspan="5" class="px-4 py-8 text-center text-gray-400 dark:text-gray-500">{{ t('common.loading') }}</td>
            </tr>
            <tr v-else-if="filteredTenants.length === 0">
              <td colspan="5" class="px-4 py-8 text-center text-gray-400 dark:text-gray-500">
                <div class="flex flex-col items-center gap-2">
                  <Building2 class="h-8 w-8 text-gray-300 dark:text-gray-600" />
                  <span>{{ t('tenants.noTenants') }}</span>
                </div>
              </td>
            </tr>
            <tr
              v-for="tenant in filteredTenants"
              :key="tenant.id"
              class="transition-colors hover:bg-gray-50 dark:hover:bg-gray-700 dark:bg-gray-800/50"
            >
              <td class="px-4 py-3">
                <div class="flex items-center gap-2">
                  <div
                    class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-blue-50 text-xs font-bold text-blue-600"
                  >
                    {{ tenant.name.charAt(0).toUpperCase() }}
                  </div>
                  <div>
                    <p class="font-medium text-gray-900 dark:text-white">{{ tenant.name }}</p>
                    <p class="text-xs text-gray-400 dark:text-gray-500">{{ tenant.id }}</p>
                  </div>
                </div>
              </td>
              <td class="px-4 py-3">
                <span class="font-mono text-sm text-gray-600 dark:text-gray-400 dark:text-gray-500">{{ tenant.slug }}</span>
              </td>
              <td class="px-4 py-3">
                <span
                  class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium"
                  :class="
                    tenant.isActive
                      ? 'bg-green-50 text-green-700'
                      : 'bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400 dark:text-gray-500'
                  "
                >
                  <span
                    class="h-1.5 w-1.5 rounded-full"
                    :class="tenant.isActive ? 'bg-green-500' : 'bg-gray-400'"
                  />
                  {{ tenant.isActive ? t('common.active') : t('common.inactive') }}
                </span>
              </td>
              <td class="hidden whitespace-nowrap px-4 py-3 text-gray-500 dark:text-gray-400 dark:text-gray-500 lg:table-cell">
                {{ formatDate(tenant.createdAt) }}
              </td>
              <td class="px-4 py-3">
                <div class="flex items-center justify-end gap-1">
                  <button
                    class="rounded p-1.5 text-gray-400 dark:text-gray-500 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800 dark:bg-gray-800 hover:text-gray-600 dark:text-gray-400 dark:text-gray-500 dark:hover:text-gray-300 dark:text-gray-600"
                    :title="t('tenants.viewDetails')"
                    @click="router.push({ name: 'tenant-detail', params: { id: tenant.id } })"
                  >
                    <Eye class="h-4 w-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
