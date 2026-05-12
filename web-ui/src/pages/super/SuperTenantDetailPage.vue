<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getTenant, listDomains, addDomain, removeDomain } from '@/api/tenant'
import {
  listSuperTenantMembers,
  listSuperTenantShortLinks,
  setSuperTenantStatus,
} from '@/api/super'
import { ApiError } from '@/api/http'
import type { Tenant, TenantDomain, ShortLinkData } from '@/types/api'
import type { SuperTenantMember } from '@/types/api'
import { ArrowLeft, Globe, Plus, Trash2, Star, Loader2, ExternalLink } from 'lucide-vue-next'

defineOptions({ name: 'SuperTenantDetailPage' })

const router = useRouter()
const route = useRoute()
const tenantId = route.params.id as string

const tenant = ref<Tenant | null>(null)
const loading = ref(false)
const notFound = ref(false)
const pageError = ref('')

const activeTab = ref<'info' | 'members' | 'short-links' | 'domains'>('info')

// Members
const members = ref<SuperTenantMember[]>([])
const membersLoading = ref(false)

// Short links
const shortLinks = ref<ShortLinkData[]>([])
const shortLinksTotal = ref(0)
const shortLinksPage = ref(1)
const shortLinksLoading = ref(false)

// Domains
const domains = ref<TenantDomain[]>([])
const domainsLoading = ref(false)
const newDomain = ref('')
const addingDomain = ref(false)
const confirmDeleteDomain = ref<string | null>(null)
const deletingDomain = ref(false)

// Status toggle
const toggleLoading = ref(false)

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

async function fetchMembers() {
  membersLoading.value = true
  try {
    const data = await listSuperTenantMembers(tenantId)
    members.value = data || []
  } catch {
    members.value = []
  } finally {
    membersLoading.value = false
  }
}

async function fetchShortLinks() {
  shortLinksLoading.value = true
  try {
    const data = await listSuperTenantShortLinks(tenantId, shortLinksPage.value, 20)
    shortLinks.value = data.shortLinks || []
    shortLinksTotal.value = data.total
  } catch {
    shortLinks.value = []
  } finally {
    shortLinksLoading.value = false
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

async function handleToggleStatus() {
  if (!tenant.value) return
  toggleLoading.value = true
  try {
    await setSuperTenantStatus(tenantId, { isActive: !tenant.value.isActive })
    tenant.value = { ...tenant.value, isActive: !tenant.value.isActive }
  } catch {
    // error handled by interceptor
  } finally {
    toggleLoading.value = false
  }
}

async function handleAddDomain() {
  if (!newDomain.value.trim()) return
  addingDomain.value = true
  try {
    const data = await addDomain(tenantId, { domain: newDomain.value.trim(), isDefault: false })
    domains.value = data || []
    newDomain.value = ''
  } catch {
    // error handled by interceptor
  } finally {
    addingDomain.value = false
  }
}

async function handleDeleteDomain(domain: string) {
  deletingDomain.value = true
  try {
    const data = await removeDomain(tenantId, domain)
    domains.value = data || []
  } catch {
    // error handled by interceptor
  } finally {
    deletingDomain.value = false
    confirmDeleteDomain.value = null
  }
}

function onTabChange(tab: 'info' | 'members' | 'short-links' | 'domains') {
  activeTab.value = tab
  if (tab === 'members' && members.value.length === 0) fetchMembers()
  if (tab === 'short-links' && shortLinks.value.length === 0) fetchShortLinks()
  if (tab === 'domains' && domains.value.length === 0) fetchDomains()
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

onMounted(fetchTenant)
</script>

<template>
  <div>
    <div class="flex items-center gap-3">
      <button
        class="rounded p-1.5 text-gray-400 dark:text-gray-500 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800 hover:text-gray-600 dark:hover:text-gray-300"
        @click="router.push({ name: 'super-tenants' })"
      >
        <ArrowLeft class="h-5 w-5" />
      </button>
      <div class="flex-1">
        <h1 class="text-xl font-semibold text-gray-900 dark:text-white">
          {{ t('super.tenantDetail.title') }}
        </h1>
        <p v-if="tenant" class="mt-0.5 text-sm text-gray-500 dark:text-gray-400">
          <span class="font-mono font-medium text-gray-700 dark:text-gray-300">{{
            tenant.id
          }}</span>
        </p>
      </div>
    </div>

    <div v-if="loading" class="mt-6 space-y-4">
      <div class="rounded-lg border bg-white dark:bg-gray-900 p-5">
        <div class="flex items-center gap-3">
          <div class="h-10 w-10 animate-pulse rounded-lg bg-gray-200 dark:bg-gray-700" />
          <div class="flex-1">
            <div class="h-5 w-40 animate-pulse rounded bg-gray-200 dark:bg-gray-700" />
            <div class="mt-2 h-4 w-24 animate-pulse rounded bg-gray-200 dark:bg-gray-700" />
          </div>
        </div>
        <div class="mt-4 grid gap-4 sm:grid-cols-3">
          <div class="h-10 animate-pulse rounded bg-gray-100 dark:bg-gray-800" />
          <div class="h-10 animate-pulse rounded bg-gray-100 dark:bg-gray-800" />
          <div class="h-10 animate-pulse rounded bg-gray-100 dark:bg-gray-800" />
        </div>
      </div>
    </div>

    <div v-else-if="notFound" class="mt-6 text-center text-gray-400 dark:text-gray-500">
      Tenant not found.
      <button
        class="ml-2 text-blue-600 hover:underline"
        @click="router.push({ name: 'super-tenants' })"
      >
        Back to list
      </button>
    </div>

    <div
      v-else-if="pageError"
      class="mt-6 rounded-lg border border-red-200 dark:border-red-800/50 bg-red-50 dark:bg-red-900/20 p-4"
    >
      <p class="text-sm text-red-600">{{ pageError }}</p>
    </div>

    <template v-else-if="tenant">
      <!-- Tenant header card -->
      <div class="mt-4 rounded-lg border bg-white dark:bg-gray-900 p-5">
        <div class="flex items-center gap-3">
          <div
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-blue-50 dark:bg-blue-900/20 text-sm font-bold text-blue-600 dark:text-blue-400"
          >
            {{ tenant.name.charAt(0).toUpperCase() }}
          </div>
          <div>
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ tenant.name }}</h2>
            <p class="text-sm text-gray-500 dark:text-gray-400">{{ tenant.slug }}</p>
          </div>
          <span
            class="ml-auto inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium"
            :class="
              tenant.isActive
                ? 'bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-400'
                : 'bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400'
            "
          >
            <span
              class="h-1.5 w-1.5 rounded-full"
              :class="tenant.isActive ? 'bg-green-500' : 'bg-gray-400'"
            />
            {{ tenant.isActive ? t('common.active') : t('common.inactive') }}
          </span>
          <button
            :disabled="toggleLoading"
            class="rounded-md px-3 py-1.5 text-xs font-medium transition-colors"
            :class="
              tenant.isActive
                ? 'bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 hover:bg-red-100 dark:hover:bg-red-900/30'
                : 'bg-green-50 dark:bg-green-900/20 text-green-600 dark:text-green-400 hover:bg-green-100 dark:hover:bg-green-900/30'
            "
            @click="handleToggleStatus"
          >
            <Loader2 v-if="toggleLoading" class="mr-1 inline h-3 w-3 animate-spin" />
            {{ tenant.isActive ? t('common.disabled') : t('common.enabled') }}
          </button>
        </div>
        <div class="mt-4 grid gap-4 sm:grid-cols-3">
          <div>
            <label class="text-xs font-medium uppercase text-gray-400 dark:text-gray-500">{{
              t('super.tenantDetail.tenantId')
            }}</label>
            <p class="mt-1 font-mono text-sm text-gray-900 dark:text-white">{{ tenant.id }}</p>
          </div>
          <div>
            <label class="text-xs font-medium uppercase text-gray-400 dark:text-gray-500">{{
              t('common.created')
            }}</label>
            <p class="mt-1 text-sm text-gray-700 dark:text-gray-300">
              {{ formatDate(tenant.createdAt) }}
            </p>
          </div>
          <div>
            <label class="text-xs font-medium uppercase text-gray-400 dark:text-gray-500">{{
              t('common.updated')
            }}</label>
            <p class="mt-1 text-sm text-gray-700 dark:text-gray-300">
              {{ formatDate(tenant.updatedAt) }}
            </p>
          </div>
        </div>
      </div>

      <!-- Tabs -->
      <div class="mt-4 border-b">
        <nav class="flex gap-6">
          <button
            v-for="tab in ['info', 'members', 'short-links', 'domains'] as const"
            :key="tab"
            class="whitespace-nowrap border-b-2 pb-3 pt-2 text-sm font-medium transition-colors"
            :class="
              activeTab === tab
                ? 'border-blue-600 text-blue-600'
                : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'
            "
            @click="onTabChange(tab)"
          >
            {{ tab === 'short-links' ? 'Short Links' : tab.charAt(0).toUpperCase() + tab.slice(1) }}
          </button>
        </nav>
      </div>

      <!-- Tab content -->
      <div class="mt-4">
        <!-- Info tab -->
        <div v-if="activeTab === 'info'" class="rounded-lg border bg-white dark:bg-gray-900 p-5">
          <div class="grid gap-4 sm:grid-cols-2">
            <div>
              <label class="text-xs font-medium uppercase text-gray-400 dark:text-gray-500">{{
                t('common.name')
              }}</label>
              <p class="mt-1 text-sm text-gray-900 dark:text-white">{{ tenant.name }}</p>
            </div>
            <div>
              <label class="text-xs font-medium uppercase text-gray-400 dark:text-gray-500">{{
                t('common.slug')
              }}</label>
              <p class="mt-1 font-mono text-sm text-gray-900 dark:text-white">{{ tenant.slug }}</p>
            </div>
          </div>
        </div>

        <!-- Members tab -->
        <div v-if="activeTab === 'members'">
          <div
            v-if="membersLoading"
            class="flex items-center justify-center rounded-lg border bg-white dark:bg-gray-900 py-8"
          >
            <Loader2 class="h-5 w-5 animate-spin text-gray-400 dark:text-gray-500" />
          </div>
          <div v-else class="overflow-hidden rounded-lg border bg-white dark:bg-gray-900">
            <table class="w-full text-sm">
              <thead>
                <tr
                  class="border-b bg-gray-50 dark:bg-gray-800/50 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400"
                >
                  <th class="px-4 py-3">{{ t('super.users.username') }}</th>
                  <th class="px-4 py-3">{{ t('members.role') }}</th>
                  <th class="hidden px-4 py-3 lg:table-cell">{{ t('members.joined') }}</th>
                </tr>
              </thead>
              <tbody class="divide-y">
                <tr v-if="members.length === 0">
                  <td colspan="3" class="px-4 py-8 text-center text-gray-400 dark:text-gray-500">
                    No members in this tenant.
                  </td>
                </tr>
                <tr
                  v-for="m in members"
                  :key="m.userId"
                  class="transition-colors hover:bg-gray-50 dark:hover:bg-gray-700 dark:bg-gray-800/50"
                >
                  <td class="px-4 py-3">
                    <div class="flex items-center gap-2">
                      <div
                        class="flex h-7 w-7 items-center justify-center rounded-full bg-gray-100 dark:bg-gray-800 text-xs font-bold text-gray-600 dark:text-gray-400"
                      >
                        {{ m.username.charAt(0).toUpperCase() }}
                      </div>
                      <span class="font-medium text-gray-900 dark:text-white">{{
                        m.username
                      }}</span>
                    </div>
                  </td>
                  <td class="px-4 py-3">
                    <span
                      class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium"
                      :class="
                        m.role === 'admin'
                          ? 'bg-blue-50 dark:bg-blue-900/20 text-blue-700 dark:text-blue-400'
                          : 'bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400'
                      "
                    >
                      {{ m.role }}
                    </span>
                  </td>
                  <td
                    class="hidden whitespace-nowrap px-4 py-3 text-gray-500 dark:text-gray-400 lg:table-cell"
                  >
                    {{ formatDate(m.joinedAt) }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Short Links tab -->
        <div v-if="activeTab === 'short-links'">
          <div
            v-if="shortLinksLoading"
            class="flex items-center justify-center rounded-lg border bg-white dark:bg-gray-900 py-8"
          >
            <Loader2 class="h-5 w-5 animate-spin text-gray-400 dark:text-gray-500" />
          </div>
          <div v-else class="overflow-hidden rounded-lg border bg-white dark:bg-gray-900">
            <table class="w-full text-sm">
              <thead>
                <tr
                  class="border-b bg-gray-50 dark:bg-gray-800/50 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400"
                >
                  <th class="px-4 py-3">{{ t('common.name') }}</th>
                  <th class="hidden px-4 py-3 md:table-cell">
                    {{ t('shortLinkDetail.targetUrl') }}
                  </th>
                  <th class="px-4 py-3">{{ t('common.status') }}</th>
                  <th class="hidden px-4 py-3 lg:table-cell">
                    {{ t('shortLinkDetail.createdBy') }}
                  </th>
                  <th class="hidden px-4 py-3 lg:table-cell">{{ t('common.created') }}</th>
                </tr>
              </thead>
              <tbody class="divide-y">
                <tr v-if="shortLinks.length === 0">
                  <td colspan="5" class="px-4 py-8 text-center text-gray-400 dark:text-gray-500">
                    No short links in this tenant.
                  </td>
                </tr>
                <tr
                  v-for="sl in shortLinks"
                  :key="sl.id"
                  class="transition-colors hover:bg-gray-50 dark:hover:bg-gray-700 dark:bg-gray-800/50"
                >
                  <td class="px-4 py-3 font-mono text-sm font-medium text-gray-900 dark:text-white">
                    {{ sl.id }}
                  </td>
                  <td
                    class="hidden max-w-[240px] truncate px-4 py-3 text-gray-500 dark:text-gray-400 md:table-cell"
                  >
                    {{ sl.url }}
                  </td>
                  <td class="px-4 py-3">
                    <span
                      class="inline-flex items-center gap-1.5 rounded-full px-2 py-0.5 text-xs font-medium"
                      :class="
                        sl.isEnable
                          ? 'bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-400'
                          : 'bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400'
                      "
                    >
                      <span
                        class="h-1.5 w-1.5 rounded-full"
                        :class="sl.isEnable ? 'bg-green-500' : 'bg-gray-400'"
                      />
                      {{ sl.isEnable ? t('common.active') : t('common.inactive') }}
                    </span>
                  </td>
                  <td class="hidden px-4 py-3 text-gray-500 dark:text-gray-400 lg:table-cell">
                    {{ sl.createdBy }}
                  </td>
                  <td
                    class="hidden whitespace-nowrap px-4 py-3 text-gray-500 dark:text-gray-400 lg:table-cell"
                  >
                    {{ formatDate(sl.createTime) }}
                  </td>
                </tr>
              </tbody>
            </table>
            <div
              v-if="shortLinksTotal > 20"
              class="border-t px-4 py-3 text-center text-xs text-gray-400 dark:text-gray-500"
            >
              {{ t('super.tenantDetail.showingOf', { shown: 20, total: shortLinksTotal }) }}
            </div>
          </div>
        </div>

        <!-- Domains tab -->
        <div v-if="activeTab === 'domains'">
          <!-- Add domain -->
          <div class="mb-4 rounded-lg border bg-white dark:bg-gray-900 p-4">
            <div class="flex flex-col gap-3 sm:flex-row sm:items-end">
              <div class="flex-1">
                <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{
                  t('common.domain')
                }}</label>
                <input
                  v-model="newDomain"
                  type="text"
                  :placeholder="t('settings.domainPlaceholder')"
                  class="w-full rounded-md border border-gray-300 dark:border-gray-600 px-3 py-2 font-mono text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
                  @keydown.enter.prevent="handleAddDomain"
                />
              </div>
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

          <div
            v-if="domainsLoading"
            class="flex items-center justify-center rounded-lg border bg-white dark:bg-gray-900 py-8"
          >
            <Loader2 class="h-5 w-5 animate-spin text-gray-400 dark:text-gray-500" />
          </div>
          <div
            v-else-if="domains.length > 0"
            class="overflow-hidden rounded-lg border bg-white dark:bg-gray-900"
          >
            <table class="w-full text-sm">
              <thead>
                <tr
                  class="border-b bg-gray-50 dark:bg-gray-800/50 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400"
                >
                  <th class="px-4 py-3">{{ t('common.domain') }}</th>
                  <th class="px-4 py-3">{{ t('common.default') }}</th>
                  <th class="hidden px-4 py-3 lg:table-cell">{{ t('common.created') }}</th>
                  <th class="px-4 py-3 text-right">{{ t('common.actions') }}</th>
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
                      <span class="font-mono text-sm text-gray-900 dark:text-white">{{
                        d.domain
                      }}</span>
                    </div>
                  </td>
                  <td class="px-4 py-3">
                    <span
                      v-if="d.isDefault"
                      class="inline-flex items-center gap-1 text-xs font-medium text-yellow-600"
                    >
                      <Star class="h-3.5 w-3.5 fill-yellow-500 text-yellow-500" />
                      Default
                    </span>
                    <span v-else class="text-xs text-gray-400 dark:text-gray-500">-</span>
                  </td>
                  <td
                    class="hidden whitespace-nowrap px-4 py-3 text-gray-500 dark:text-gray-400 lg:table-cell"
                  >
                    {{ formatDate(d.createdAt) }}
                  </td>
                  <td class="px-4 py-3">
                    <div class="flex items-center justify-end gap-1">
                      <a
                        :href="'https://' + d.domain"
                        target="_blank"
                        rel="noopener noreferrer"
                        class="rounded p-1.5 text-gray-400 dark:text-gray-500 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800 hover:text-gray-600 dark:hover:text-gray-300"
                      >
                        <ExternalLink class="h-4 w-4" />
                      </a>
                      <button
                        class="rounded p-1.5 text-gray-400 dark:text-gray-500 transition-colors hover:bg-red-50 dark:hover:bg-red-900/20 hover:text-red-600"
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
          <div
            v-else
            class="rounded-lg border bg-white dark:bg-gray-900 py-8 text-center text-sm text-gray-400 dark:text-gray-500"
          >
            No domains configured for this tenant.
          </div>
        </div>
      </div>
    </template>

    <!-- Delete domain confirmation -->
    <Teleport to="body">
      <div
        v-if="confirmDeleteDomain"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
        @click.self="confirmDeleteDomain = null"
      >
        <div class="mx-4 w-full max-w-sm rounded-lg bg-white dark:bg-gray-900 p-6 shadow-xl">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ t('settings.removeDomain') }}
          </h3>
          <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">
            Are you sure you want to remove
            <span class="font-mono font-medium text-gray-700 dark:text-gray-300">{{
              confirmDeleteDomain
            }}</span
            >?
          </p>
          <div class="mt-4 flex justify-end gap-2">
            <button
              class="rounded-md px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800"
              @click="confirmDeleteDomain = null"
            >
              Cancel
            </button>
            <button
              :disabled="deletingDomain"
              class="rounded-md bg-red-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-red-700 disabled:opacity-50"
              @click="handleDeleteDomain(confirmDeleteDomain!)"
            >
              {{ deletingDomain ? 'Removing...' : 'Remove' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
