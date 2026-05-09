<script setup lang="ts">
import { ref, onMounted, nextTick, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getTenant, listDomains, addDomain, removeDomain } from '@/api/tenant'
import { ApiError } from '@/api/http'
import type { Tenant, TenantDomain } from '@/types/api'
import {
  ArrowLeft,
  Globe,
  Plus,
  Trash2,
  Star,
  Loader2,
  ExternalLink,
} from 'lucide-vue-next'

defineOptions({ name: 'TenantDetailPage' })

const router = useRouter()
const route = useRoute()
const tenantId = route.params.id as string

const tenant = ref<Tenant | null>(null)
const domains = ref<TenantDomain[]>([])
const loading = ref(false)
const notFound = ref(false)
const pageError = ref('')

const domainsLoading = ref(false)

// Add domain form
const newDomain = ref('')
const newDomainIsDefault = ref(false)
const addingDomain = ref(false)

// Delete domain
const confirmDeleteDomain = ref<string | null>(null)
const deleting = ref(false)
const dialogEl = ref<HTMLElement | null>(null)

async function fetchTenant() {
  loading.value = true
  try {
    tenant.value = await getTenant(tenantId)
  } catch (e) {
    if (e instanceof ApiError && e.status === 404) {
      notFound.value = true
    } else {
      const msg = e instanceof Error ? e.message : ''
      pageError.value = msg || 'Failed to load tenant.'
    }
  } finally {
    loading.value = false
  }
}

async function fetchDomains() {
  domainsLoading.value = true
  try {
    const data = await listDomains(tenantId)
    domains.value = data || []
  } catch {
    domains.value = []
  } finally {
    domainsLoading.value = false
  }
}

async function handleAddDomain() {
  if (!newDomain.value.trim()) return
  addingDomain.value = true
  try {
    const data = await addDomain(tenantId, {
      domain: newDomain.value.trim(),
      isDefault: newDomainIsDefault.value,
    })
    domains.value = data || []
    newDomain.value = ''
    newDomainIsDefault.value = false
  } catch {
    // error handled by interceptor
  } finally {
    addingDomain.value = false
  }
}

async function handleDeleteDomain(domain: string) {
  deleting.value = true
  try {
    const data = await removeDomain(tenantId, domain)
    domains.value = data || []
  } catch {
    // error handled by interceptor
  } finally {
    deleting.value = false
    confirmDeleteDomain.value = null
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

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    confirmDeleteDomain.value = null
  }
}

watch(confirmDeleteDomain, (open) => {
  if (open) {
    document.addEventListener('keydown', handleKeydown)
    nextTick(() => {
      dialogEl.value?.focus()
    })
  } else {
    document.removeEventListener('keydown', handleKeydown)
  }
})

onMounted(() => {
  fetchTenant()
  fetchDomains()
})
</script>

<template>
  <div>
    <div class="flex items-center gap-3">
      <button
        class="rounded p-1.5 text-gray-400 dark:text-gray-500 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800 dark:bg-gray-800 hover:text-gray-600 dark:text-gray-400 dark:text-gray-500 dark:hover:text-gray-300 dark:text-gray-600"
        @click="router.push({ name: 'tenants' })"
      >
        <ArrowLeft class="h-5 w-5" />
      </button>
      <div class="flex-1">
        <h1 class="text-xl font-semibold text-gray-900 dark:text-white">Tenant Detail</h1>
        <p v-if="tenant" class="mt-0.5 text-sm text-gray-500 dark:text-gray-400 dark:text-gray-500">
          <span class="font-mono font-medium text-gray-700 dark:text-gray-300 dark:text-gray-600">{{ tenant.id }}</span>
        </p>
      </div>
    </div>

    <div v-if="loading" class="mt-6 text-center text-gray-400 dark:text-gray-500">
      <Loader2 class="inline h-5 w-5 animate-spin" />
    </div>

    <div v-else-if="notFound" class="mt-6 text-center text-gray-400 dark:text-gray-500">
      Tenant not found.
      <button
        class="ml-2 text-blue-600 hover:underline"
        @click="router.push({ name: 'tenants' })"
      >
        Back to list
      </button>
    </div>

    <div v-else-if="pageError" class="mt-6 rounded-lg border border-red-200 bg-red-50 p-4">
      <p class="text-sm text-red-600">{{ pageError }}</p>
    </div>

    <template v-else-if="tenant">
      <!-- Basic info card -->
      <div class="mt-4 rounded-lg border bg-white dark:bg-gray-900 p-5">
        <div class="flex items-center gap-3 border-b border-gray-100 dark:border-gray-800 pb-4">
          <div
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-blue-50 text-sm font-bold text-blue-600"
          >
            {{ tenant.name.charAt(0).toUpperCase() }}
          </div>
          <div>
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ tenant.name }}</h2>
            <p class="text-sm text-gray-500 dark:text-gray-400 dark:text-gray-500">{{ tenant.slug }}</p>
          </div>
          <span
            class="ml-auto inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium"
            :class="tenant.isActive ? 'bg-green-50 text-green-700' : 'bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400 dark:text-gray-500'"
          >
            <span
              class="h-1.5 w-1.5 rounded-full"
              :class="tenant.isActive ? 'bg-green-500' : 'bg-gray-400'"
            />
            {{ tenant.isActive ? 'Active' : 'Inactive' }}
          </span>
        </div>
        <div class="mt-4 grid gap-4 sm:grid-cols-2">
          <div>
            <label class="text-xs font-medium uppercase text-gray-400 dark:text-gray-500">Tenant ID</label>
            <p class="mt-1 font-mono text-sm text-gray-900 dark:text-white">{{ tenant.id }}</p>
          </div>
          <div>
            <label class="text-xs font-medium uppercase text-gray-400 dark:text-gray-500">Slug</label>
            <p class="mt-1 font-mono text-sm text-gray-900 dark:text-white">{{ tenant.slug }}</p>
          </div>
          <div>
            <label class="text-xs font-medium uppercase text-gray-400 dark:text-gray-500">Created</label>
            <p class="mt-1 text-sm text-gray-700 dark:text-gray-300 dark:text-gray-600">{{ formatDate(tenant.createdAt) }}</p>
          </div>
          <div>
            <label class="text-xs font-medium uppercase text-gray-400 dark:text-gray-500">Updated</label>
            <p class="mt-1 text-sm text-gray-700 dark:text-gray-300 dark:text-gray-600">{{ formatDate(tenant.updatedAt) }}</p>
          </div>
        </div>
      </div>

      <!-- Domain management -->
      <div class="mt-6">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <Globe class="h-5 w-5 text-gray-500 dark:text-gray-400 dark:text-gray-500" />
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">Domains</h2>
            <span class="rounded-full bg-gray-100 dark:bg-gray-800 px-2 py-0.5 text-xs font-medium text-gray-600 dark:text-gray-400 dark:text-gray-500">
              {{ domains.length }}
            </span>
          </div>
        </div>

        <!-- Add domain form -->
        <div class="mt-4 rounded-lg border bg-white dark:bg-gray-900 p-4">
          <div class="flex flex-col gap-3 sm:flex-row sm:items-end">
            <div class="flex-1">
              <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300 dark:text-gray-600">Domain</label>
              <input
                v-model="newDomain"
                type="text"
                placeholder="example.com"
                class="w-full rounded-md border border-gray-300 dark:border-gray-600 px-3 py-2 font-mono text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
                @keydown.enter.prevent="handleAddDomain"
              />
            </div>
            <div class="flex items-center gap-3">
              <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 dark:text-gray-600">
                <input
                  v-model="newDomainIsDefault"
                  type="checkbox"
                  class="h-4 w-4 rounded border-gray-300 dark:border-gray-600 text-blue-600 focus:ring-blue-500"
                />
                Default
              </label>
              <button
                :disabled="!newDomain.trim() || addingDomain"
                class="inline-flex items-center gap-1.5 rounded-md bg-blue-600 px-3 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:opacity-50"
                @click="handleAddDomain"
              >
                <Loader2 v-if="addingDomain" class="h-3.5 w-3.5 animate-spin" />
                <Plus v-else class="h-3.5 w-3.5" />
                Add
              </button>
            </div>
          </div>
        </div>

        <!-- Domain list -->
        <div v-if="domainsLoading" class="mt-4 flex items-center justify-center rounded-lg border bg-white dark:bg-gray-900 py-8">
          <Loader2 class="h-5 w-5 animate-spin text-gray-400 dark:text-gray-500" />
        </div>
        <div v-else-if="domains.length > 0" class="mt-4 overflow-hidden rounded-lg border bg-white dark:bg-gray-900 dark:bg-gray-900">
          <table class="w-full text-sm">
            <thead>
              <tr
                class="border-b bg-gray-50 dark:bg-gray-800/50 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400 dark:text-gray-500"
              >
                <th class="px-4 py-3">Domain</th>
                <th class="px-4 py-3">Default</th>
                <th class="hidden px-4 py-3 lg:table-cell">Created</th>
                <th class="px-4 py-3 text-right">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y">
              <tr
                v-for="d in domains"
                :key="d.id"
                class="transition-colors hover:bg-gray-50 dark:hover:bg-gray-700 dark:bg-gray-800/50"
              >
                <td class="px-4 py-3">
                  <div class="flex items-center gap-2">
                    <Globe class="h-4 w-4 text-gray-400 dark:text-gray-500" />
                    <span class="font-mono text-sm text-gray-900 dark:text-white">{{ d.domain }}</span>
                  </div>
                </td>
                <td class="px-4 py-3">
                  <span v-if="d.isDefault" class="inline-flex items-center gap-1 text-xs font-medium text-yellow-600">
                    <Star class="h-3.5 w-3.5 fill-yellow-500 text-yellow-500" />
                    Default
                  </span>
                  <span v-else class="text-xs text-gray-400 dark:text-gray-500">—</span>
                </td>
                <td class="hidden whitespace-nowrap px-4 py-3 text-gray-500 dark:text-gray-400 dark:text-gray-500 lg:table-cell">
                  {{ formatDate(d.createdAt) }}
                </td>
                <td class="px-4 py-3">
                  <div class="flex items-center justify-end gap-1">
                    <a
                      :href="'https://' + d.domain"
                      target="_blank"
                      rel="noopener noreferrer"
                      class="rounded p-1.5 text-gray-400 dark:text-gray-500 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800 dark:bg-gray-800 hover:text-gray-600 dark:text-gray-400 dark:text-gray-500 dark:hover:text-gray-300 dark:text-gray-600"
                      title="Open domain"
                    >
                      <ExternalLink class="h-4 w-4" />
                    </a>
                    <button
                      class="rounded p-1.5 text-gray-400 dark:text-gray-500 transition-colors hover:bg-red-50 hover:text-red-600"
                      title="Remove domain"
                      @click="confirmDeleteDomain = d.domain"
                    >
                      <Trash2 class="h-4 w-4" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else class="mt-4 rounded-lg border bg-white dark:bg-gray-900 py-8 text-center text-sm text-gray-400 dark:text-gray-500">
          No domains configured for this tenant.
        </div>
      </div>
    </template>

    <!-- Delete confirmation modal -->
    <Teleport to="body">
      <div
        v-if="confirmDeleteDomain"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
        @click.self="confirmDeleteDomain = null"
      >
        <div
          ref="dialogEl"
          role="dialog"
          aria-modal="true"
          aria-labelledby="remove-domain-title"
          tabindex="-1"
          class="mx-4 w-full max-w-sm rounded-lg bg-white dark:bg-gray-900 p-6 shadow-xl"
        >
          <h3 id="remove-domain-title" class="text-lg font-semibold text-gray-900 dark:text-white">Remove Domain</h3>
          <p class="mt-2 text-sm text-gray-500 dark:text-gray-400 dark:text-gray-500">
            Are you sure you want to remove
            <span class="font-mono font-medium text-gray-700 dark:text-gray-300 dark:text-gray-600">{{ confirmDeleteDomain }}</span
            >? This action cannot be undone.
          </p>
          <div class="mt-4 flex justify-end gap-2">
            <button
              class="rounded-md px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 dark:text-gray-600 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800 dark:bg-gray-800"
              @click="confirmDeleteDomain = null"
            >
              Cancel
            </button>
            <button
              :disabled="deleting"
              class="rounded-md bg-red-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-red-700 disabled:opacity-50"
              @click="handleDeleteDomain(confirmDeleteDomain!)"
            >
              {{ deleting ? 'Removing...' : 'Remove' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
