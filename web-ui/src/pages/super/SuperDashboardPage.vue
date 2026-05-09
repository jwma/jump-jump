<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { listSuperUsers, listSuperTenants } from '@/api/super'
import { Building2, Users, Loader2 } from 'lucide-vue-next'

const { t } = useI18n()

defineOptions({ name: 'SuperDashboardPage' })

const loading = ref(true)
const totalTenants = ref(0)
const totalUsers = ref(0)
const recentTenants = ref<{ id: string; name: string; slug: string; isActive: boolean; createdAt: string }[]>([])

async function fetchDashboardData() {
  loading.value = true
  try {
    const [usersData, tenantsData] = await Promise.all([
      listSuperUsers(1, 1),
      listSuperTenants(1, 100),
    ])
    totalUsers.value = usersData.total
    totalTenants.value = tenantsData.total

    recentTenants.value = (tenantsData.tenants || [])
      .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
      .slice(0, 8)
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
  })
}

onMounted(fetchDashboardData)
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ t('super.dashboard.title') }}</h1>
      <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{ t('super.dashboard.subtitle') }}</p>
    </div>

    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
      <div class="rounded-xl border border-gray-200 bg-white p-5 shadow-sm dark:border-gray-700 dark:bg-gray-900">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400">{{ t('super.dashboard.totalTenants') }}</p>
            <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">
              <span v-if="loading" class="inline-block h-7 w-16 animate-pulse rounded bg-gray-200 dark:bg-gray-700" />
              <template v-else>{{ totalTenants.toLocaleString() }}</template>
            </p>
          </div>
          <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-50">
            <Building2 class="h-5 w-5 text-blue-600" />
          </div>
        </div>
      </div>

      <div class="rounded-xl border border-gray-200 bg-white p-5 shadow-sm dark:border-gray-700 dark:bg-gray-900">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400">{{ t('super.dashboard.totalUsers') }}</p>
            <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">
              <span v-if="loading" class="inline-block h-7 w-16 animate-pulse rounded bg-gray-200 dark:bg-gray-700" />
              <template v-else>{{ totalUsers.toLocaleString() }}</template>
            </p>
          </div>
          <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-green-50">
            <Users class="h-5 w-5 text-green-600" />
          </div>
        </div>
      </div>
    </div>

    <div class="rounded-xl border border-gray-200 bg-white shadow-sm dark:border-gray-700 dark:bg-gray-900">
      <div class="flex items-center justify-between border-b border-gray-100 px-5 py-4 dark:border-gray-800">
        <div class="flex items-center gap-2">
          <Building2 class="h-4 w-4 text-gray-500 dark:text-gray-400" />
          <h2 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('super.dashboard.tenants') }}</h2>
        </div>
        <router-link
          :to="{ name: 'super-tenants' }"
          class="text-xs font-medium text-blue-600 hover:text-blue-700"
        >
          {{ t('common.viewAll') }}
        </router-link>
      </div>
      <div v-if="loading" class="flex items-center justify-center p-8">
        <Loader2 class="h-5 w-5 animate-spin text-gray-400" />
      </div>
      <div v-else-if="recentTenants.length === 0" class="px-5 py-12 text-center">
        <Building2 class="mx-auto h-10 w-10 text-gray-300 dark:text-gray-600" />
        <p class="mt-3 text-sm font-medium text-gray-500 dark:text-gray-400">{{ t('super.dashboard.noTenants') }}</p>
        <p class="mt-1 text-sm text-gray-400 dark:text-gray-500">{{ t('super.dashboard.tenantsWillAppear') }}</p>
        <router-link
          :to="{ name: 'super-tenants' }"
          class="mt-3 inline-flex items-center gap-1.5 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700"
        >
          {{ t('super.dashboard.viewAllTenants') }}
        </router-link>
      </div>
      <div v-else class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-gray-100 dark:border-gray-800 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400">
              <th class="px-5 py-3">{{ t('common.name') }}</th>
              <th class="px-5 py-3">{{ t('common.slug') }}</th>
              <th class="px-5 py-3">{{ t('common.status') }}</th>
              <th class="hidden px-5 py-3 md:table-cell">{{ t('common.created') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50 dark:divide-gray-800">
            <tr v-for="tenant in recentTenants" :key="tenant.id" class="transition-colors hover:bg-gray-50 dark:hover:bg-gray-800">
              <td class="px-5 py-3">
                <router-link
                  :to="{ name: 'super-tenant-detail', params: { id: tenant.id } }"
                  class="font-medium text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-300"
                >
                  {{ tenant.name }}
                </router-link>
              </td>
              <td class="px-5 py-3 font-mono text-xs text-gray-600 dark:text-gray-400">{{ tenant.slug }}</td>
              <td class="px-5 py-3">
                <span
                  class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium"
                  :class="tenant.isActive ? 'bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-400' : 'bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400'"
                >
                  <span class="h-1.5 w-1.5 rounded-full" :class="tenant.isActive ? 'bg-green-500' : 'bg-gray-400'" />
                  {{ tenant.isActive ? t('common.active') : t('common.inactive') }}
                </span>
              </td>
              <td class="hidden whitespace-nowrap px-5 py-3 text-gray-500 dark:text-gray-400 md:table-cell">
                {{ formatDate(tenant.createdAt) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
