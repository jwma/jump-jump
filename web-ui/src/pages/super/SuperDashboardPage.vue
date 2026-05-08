<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { listSuperUsers, listSuperTenants } from '@/api/super'
import { Building2, Users, Loader2 } from 'lucide-vue-next'

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
      <h1 class="text-2xl font-bold text-gray-900">Super Admin Dashboard</h1>
      <p class="mt-1 text-sm text-gray-500">Platform overview and statistics.</p>
    </div>

    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
      <div class="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-500">Total Tenants</p>
            <p class="mt-1 text-2xl font-bold text-gray-900">
              <span v-if="loading" class="inline-block h-7 w-16 animate-pulse rounded bg-gray-200" />
              <template v-else>{{ totalTenants.toLocaleString() }}</template>
            </p>
          </div>
          <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-50">
            <Building2 class="h-5 w-5 text-blue-600" />
          </div>
        </div>
      </div>

      <div class="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-500">Total Users</p>
            <p class="mt-1 text-2xl font-bold text-gray-900">
              <span v-if="loading" class="inline-block h-7 w-16 animate-pulse rounded bg-gray-200" />
              <template v-else>{{ totalUsers.toLocaleString() }}</template>
            </p>
          </div>
          <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-green-50">
            <Users class="h-5 w-5 text-green-600" />
          </div>
        </div>
      </div>
    </div>

    <div class="rounded-xl border border-gray-200 bg-white shadow-sm">
      <div class="flex items-center justify-between border-b border-gray-100 px-5 py-4">
        <div class="flex items-center gap-2">
          <Building2 class="h-4 w-4 text-gray-500" />
          <h2 class="text-base font-semibold text-gray-900">Tenants</h2>
        </div>
        <router-link
          :to="{ name: 'super-tenants' }"
          class="text-xs font-medium text-blue-600 hover:text-blue-700"
        >
          View all
        </router-link>
      </div>
      <div v-if="loading" class="flex items-center justify-center p-8">
        <Loader2 class="h-5 w-5 animate-spin text-gray-400" />
      </div>
      <div v-else-if="recentTenants.length === 0" class="px-5 py-8 text-center text-sm text-gray-400">
        No tenants found.
      </div>
      <div v-else class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-gray-100 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
              <th class="px-5 py-3">Name</th>
              <th class="px-5 py-3">Slug</th>
              <th class="px-5 py-3">Status</th>
              <th class="hidden px-5 py-3 md:table-cell">Created</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50">
            <tr v-for="t in recentTenants" :key="t.id" class="transition-colors hover:bg-gray-50">
              <td class="px-5 py-3">
                <router-link
                  :to="{ name: 'super-tenant-detail', params: { id: t.id } }"
                  class="font-medium text-blue-600 hover:text-blue-700"
                >
                  {{ t.name }}
                </router-link>
              </td>
              <td class="px-5 py-3 font-mono text-xs text-gray-600">{{ t.slug }}</td>
              <td class="px-5 py-3">
                <span
                  class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium"
                  :class="t.isActive ? 'bg-green-50 text-green-700' : 'bg-gray-100 text-gray-500'"
                >
                  <span class="h-1.5 w-1.5 rounded-full" :class="t.isActive ? 'bg-green-500' : 'bg-gray-400'" />
                  {{ t.isActive ? 'Active' : 'Inactive' }}
                </span>
              </td>
              <td class="hidden whitespace-nowrap px-5 py-3 text-gray-500 md:table-cell">
                {{ formatDate(t.createdAt) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
