<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { listSuperTenants, setSuperTenantStatus } from '@/api/super'
import type { SuperTenant } from '@/types/api'
import {
  Building2,
  Eye,
  ChevronLeft,
  ChevronRight,
  Loader2,
} from 'lucide-vue-next'

defineOptions({ name: 'SuperTenantsPage' })

const router = useRouter()

const tenants = ref<SuperTenant[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = 20
const total = ref(0)
const toggleLoading = ref<string | null>(null)

async function fetchTenants() {
  loading.value = true
  try {
    const data = await listSuperTenants(page.value, pageSize)
    tenants.value = data.tenants || []
    total.value = data.total
  } catch {
    // error handled by interceptor
  } finally {
    loading.value = false
  }
}

function goToPage(p: number) {
  const maxPage = Math.ceil(total.value / pageSize) || 1
  if (p < 1 || p > maxPage) return
  page.value = p
  fetchTenants()
}

async function handleToggleStatus(tenant: SuperTenant) {
  toggleLoading.value = tenant.id
  try {
    await setSuperTenantStatus(tenant.id, { isActive: !tenant.isActive })
    tenant.isActive = !tenant.isActive
  } catch {
    // error handled by interceptor
  } finally {
    toggleLoading.value = null
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
    <div>
      <h1 class="text-xl font-semibold text-gray-900">Tenant Management</h1>
      <p class="mt-1 text-sm text-gray-500">{{ total }} tenant{{ total !== 1 ? 's' : '' }} total</p>
    </div>

    <!-- Table -->
    <div class="mt-4 overflow-hidden rounded-lg border bg-white">
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b bg-gray-50 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
              <th class="px-4 py-3">Name</th>
              <th class="px-4 py-3">Slug</th>
              <th class="px-4 py-3">Status</th>
              <th class="hidden px-4 py-3 lg:table-cell">Created</th>
              <th class="px-4 py-3 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y">
            <tr v-if="loading">
              <td colspan="5" class="px-4 py-8 text-center text-gray-400">
                <Loader2 class="inline h-5 w-5 animate-spin" />
              </td>
            </tr>
            <tr v-else-if="tenants.length === 0">
              <td colspan="5" class="px-4 py-8 text-center text-gray-400">
                <div class="flex flex-col items-center gap-2">
                  <Building2 class="h-8 w-8 text-gray-300" />
                  <span>No tenants found.</span>
                </div>
              </td>
            </tr>
            <tr
              v-for="tenant in tenants"
              :key="tenant.id"
              class="transition-colors hover:bg-gray-50"
            >
              <td class="px-4 py-3">
                <div class="flex items-center gap-2">
                  <div
                    class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-blue-50 text-xs font-bold text-blue-600"
                  >
                    {{ tenant.name.charAt(0).toUpperCase() }}
                  </div>
                  <div>
                    <p class="font-medium text-gray-900">{{ tenant.name }}</p>
                    <p class="text-xs text-gray-400">{{ tenant.id }}</p>
                  </div>
                </div>
              </td>
              <td class="px-4 py-3">
                <span class="font-mono text-sm text-gray-600">{{ tenant.slug }}</span>
              </td>
              <td class="px-4 py-3">
                <button
                  :disabled="toggleLoading === tenant.id"
                  class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium transition-colors"
                  :class="[
                    tenant.isActive ? 'bg-green-50 text-green-700 hover:bg-green-100' : 'bg-gray-100 text-gray-500 hover:bg-gray-200',
                    toggleLoading === tenant.id ? 'opacity-50' : 'cursor-pointer',
                  ]"
                  @click="handleToggleStatus(tenant)"
                >
                  <Loader2 v-if="toggleLoading === tenant.id" class="h-3 w-3 animate-spin" />
                  <span v-else class="h-1.5 w-1.5 rounded-full" :class="tenant.isActive ? 'bg-green-500' : 'bg-gray-400'" />
                  {{ tenant.isActive ? 'Active' : 'Inactive' }}
                </button>
              </td>
              <td class="hidden whitespace-nowrap px-4 py-3 text-gray-500 lg:table-cell">
                {{ formatDate(tenant.createdAt) }}
              </td>
              <td class="px-4 py-3">
                <div class="flex items-center justify-end gap-1">
                  <button
                    class="rounded p-1.5 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600"
                    title="View details"
                    @click="router.push({ name: 'super-tenant-detail', params: { id: tenant.id } })"
                  >
                    <Eye class="h-4 w-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="total > pageSize" class="flex items-center justify-between border-t px-4 py-3">
        <p class="text-xs text-gray-500">
          Page {{ page }} of {{ Math.ceil(total / pageSize) }} ({{ total }} total)
        </p>
        <div class="flex items-center gap-1">
          <button
            :disabled="page <= 1"
            class="rounded p-1.5 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600 disabled:opacity-50"
            @click="goToPage(page - 1)"
          >
            <ChevronLeft class="h-4 w-4" />
          </button>
          <button
            :disabled="page >= Math.ceil(total / pageSize)"
            class="rounded p-1.5 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600 disabled:opacity-50"
            @click="goToPage(page + 1)"
          >
            <ChevronRight class="h-4 w-4" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
