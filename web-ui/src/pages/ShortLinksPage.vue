<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { listShortLinks, deleteShortLink, updateShortLink } from '@/api/short-link'
import { useToast } from '@/composables/useToast'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import type { ShortLinkData } from '@/types/api'
import {
  Plus,
  Search,
  ExternalLink,
  Copy,
  Trash2,
  Pencil,
  ChevronLeft,
  ChevronRight,
  Check,
  Eye,
  ArrowUpDown,
  ArrowUp,
  ArrowDown,
  Link,
  Loader2,
} from 'lucide-vue-next'

defineOptions({ name: 'ShortLinksPage' })

const router = useRouter()
const toast = useToast()
const { t } = useI18n()

const links = ref<ShortLinkData[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)
const pageJump = ref('')

const searchId = ref('')
const searchUrl = ref('')
const filterEnabled = ref<'' | 'true' | 'false'>('')

const selectedIds = ref<Set<string>>(new Set())
const deleting = ref<string | null>(null)
const togglingId = ref<string | null>(null)
const copiedId = ref<string | null>(null)
const confirmDeleteId = ref<string | null>(null)
const batchUpdating = ref(false)

type SortField = 'createTime' | 'updateTime' | 'id'
type SortDirection = 'asc' | 'desc'
const sortField = ref<SortField>('createTime')
const sortDirection = ref<SortDirection>('desc')

const filteredLinks = computed(() => {
  let result = links.value
  if (searchId.value) {
    const q = searchId.value.toLowerCase()
    result = result.filter((l) => l.id.toLowerCase().includes(q))
  }
  if (searchUrl.value) {
    const q = searchUrl.value.toLowerCase()
    result = result.filter((l) => l.url.toLowerCase().includes(q))
  }
  if (filterEnabled.value) {
    const enabled = filterEnabled.value === 'true'
    result = result.filter((l) => l.isEnable === enabled)
  }
  return result
})

const sortedLinks = computed(() => {
  const arr = [...filteredLinks.value]
  const dir = sortDirection.value === 'asc' ? 1 : -1
  return arr.sort((a, b) => {
    switch (sortField.value) {
      case 'createTime':
        return dir * (new Date(a.createTime).getTime() - new Date(b.createTime).getTime())
      case 'updateTime':
        return dir * (new Date(a.updateTime).getTime() - new Date(b.updateTime).getTime())
      case 'id':
        return dir * a.id.localeCompare(b.id)
      default:
        return 0
    }
  })
})

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize.value)))
const allSelected = computed(
  () =>
    filteredLinks.value.length > 0 && filteredLinks.value.every((l) => selectedIds.value.has(l.id)),
)

const pageNumbers = computed(() => {
  const pages: (number | '...')[] = []
  const tp = totalPages.value
  const p = page.value
  if (tp <= 7) {
    for (let i = 1; i <= tp; i++) pages.push(i)
  } else {
    pages.push(1)
    if (p > 3) pages.push('...')
    const start = Math.max(2, p - 1)
    const end = Math.min(tp - 1, p + 1)
    for (let i = start; i <= end; i++) pages.push(i)
    if (p < tp - 2) pages.push('...')
    pages.push(tp)
  }
  return pages
})

const pageSizeOptions = [10, 20, 50, 100]

async function fetchLinks() {
  loading.value = true
  try {
    const data = await listShortLinks(page.value, pageSize.value)
    links.value = data.shortLinks
    total.value = data.total
    selectedIds.value.clear()
  } catch {
    // error handled by interceptor
  } finally {
    loading.value = false
  }
}

function toggleSort(field: SortField) {
  if (sortField.value === field) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortField.value = field
    sortDirection.value = 'desc'
  }
}

function SortIcon({ field }: { field: SortField }) {
  if (sortField.value !== field) return ArrowUpDown
  return sortDirection.value === 'asc' ? ArrowUp : ArrowDown
}

function toggleSelect(id: string) {
  if (selectedIds.value.has(id)) {
    selectedIds.value.delete(id)
  } else {
    selectedIds.value.add(id)
  }
}

function toggleSelectAll() {
  if (allSelected.value) {
    selectedIds.value.clear()
  } else {
    filteredLinks.value.forEach((l) => selectedIds.value.add(l.id))
  }
}

async function handleToggleEnable(link: ShortLinkData) {
  togglingId.value = link.id
  try {
    await updateShortLink(link.id, {
      url: link.url,
      description: link.description,
      isEnable: !link.isEnable,
    })
    link.isEnable = !link.isEnable
  } catch {
    // error handled by interceptor
  } finally {
    togglingId.value = null
  }
}

async function handleDelete(id: string) {
  deleting.value = id
  try {
    await deleteShortLink(id)
    links.value = links.value.filter((l) => l.id !== id)
    total.value--
    confirmDeleteId.value = null
    toast.success(t('shortLinks.deleted'))
    if (links.value.length === 0 && page.value > 1) {
      page.value--
      fetchLinks()
    }
  } catch (e: unknown) {
    toast.error((e as Error).message || t('shortLinks.deleteFailed'))
    deleting.value = null
  }
}

async function handleBatchDelete() {
  if (selectedIds.value.size === 0) return
  deleting.value = 'batch'
  const count = selectedIds.value.size
  try {
    await Promise.all([...selectedIds.value].map((id) => deleteShortLink(id)))
    links.value = links.value.filter((l) => !selectedIds.value.has(l.id))
    total.value -= count
    selectedIds.value.clear()
    toast.success(t('shortLinks.deletedCount', { count, suffix: count !== 1 ? 's' : '' }))
  } catch (e: unknown) {
    toast.error((e as Error).message || t('shortLinks.deleteSomeFailed'))
    await fetchLinks()
  } finally {
    deleting.value = null
  }
}

async function handleBatchToggleEnable(enable: boolean) {
  if (selectedIds.value.size === 0) return
  batchUpdating.value = true
  const ids = [...selectedIds.value]
  try {
    let updated = 0
    await Promise.all(
      ids.map(async (id) => {
        const link = links.value.find((l) => l.id === id)
        if (link && link.isEnable !== enable) {
          await updateShortLink(id, {
            url: link.url,
            description: link.description,
            isEnable: enable,
          })
          link.isEnable = enable
          updated++
        }
      }),
    )
    toast.success(t(enable ? 'shortLinks.batchEnabled' : 'shortLinks.batchDisabled', { count: updated, suffix: updated !== 1 ? 's' : '' }))
  } catch (e: unknown) {
    toast.error((e as Error).message || t('shortLinks.enableDisableFailed', { action: enable ? t('common.enabled') : t('common.disabled') }))
  } finally {
    batchUpdating.value = false
  }
}

function copyLink(id: string) {
  const url = `${window.location.origin}/${id}`
  navigator.clipboard.writeText(url).then(
    () => {
      copiedId.value = id
      setTimeout(() => {
        copiedId.value = null
      }, 2000)
    },
    () => {
      // clipboard not available
    },
  )
}

function goPage(p: number) {
  if (p < 1 || p > totalPages.value) return
  page.value = p
  fetchLinks()
}

function handlePageJump() {
  const n = parseInt(pageJump.value, 10)
  if (!isNaN(n) && n >= 1 && n <= totalPages.value) {
    goPage(n)
  }
  pageJump.value = ''
}

function handlePageSizeChange() {
  page.value = 1
  fetchLinks()
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

onMounted(fetchLinks)
</script>

<template>
  <div>
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-xl font-semibold text-gray-900">{{ t('shortLinks.title') }}</h1>
        <p class="mt-1 text-sm text-gray-500">{{ t('shortLinks.linkCountTotal', { count: total, suffix: total !== 1 ? 's' : '' }) }}</p>
      </div>
      <router-link
        :to="{ name: 'short-link-create' }"
        class="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700"
      >
        <Plus class="h-4 w-4" />
        {{ t('shortLinks.createLink') }}
      </router-link>
    </div>

    <!-- Search & Filters -->
    <div
      class="mt-4 flex flex-col gap-3 rounded-lg border bg-white p-4 sm:flex-row sm:items-center"
    >
      <div class="relative flex-1">
        <Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
        <input
          v-model="searchId"
          type="text"
          :placeholder="t('shortLinks.searchById')"
          class="w-full rounded-md border border-gray-300 py-2 pl-9 pr-3 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
        />
      </div>
      <div class="relative flex-1">
        <Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
        <input
          v-model="searchUrl"
          type="text"
          :placeholder="t('shortLinks.searchByUrl')"
          class="w-full rounded-md border border-gray-300 py-2 pl-9 pr-3 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
        />
      </div>
      <select
        v-model="filterEnabled"
        class="rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none sm:w-auto"
      >
        <option value="">{{ t('shortLinks.allStatus') }}</option>
        <option value="true">{{ t('common.enabled') }}</option>
        <option value="false">{{ t('common.disabled') }}</option>
      </select>
    </div>

    <!-- Batch actions -->
    <div v-if="selectedIds.size > 0" class="mt-3 flex flex-wrap items-center gap-2">
      <span class="text-sm text-gray-600">{{ t('shortLinks.selected', { count: selectedIds.size }) }}</span>
      <button
        :disabled="batchUpdating"
        class="inline-flex items-center gap-1.5 rounded-md bg-green-50 px-3 py-1.5 text-sm font-medium text-green-600 transition-colors hover:bg-green-100 disabled:opacity-50"
        @click="handleBatchToggleEnable(true)"
      >
        {{ t('shortLinks.enableSelected') }}
      </button>
      <button
        :disabled="batchUpdating"
        class="inline-flex items-center gap-1.5 rounded-md bg-gray-100 px-3 py-1.5 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-200 disabled:opacity-50"
        @click="handleBatchToggleEnable(false)"
      >
        {{ t('shortLinks.disableSelected') }}
      </button>
      <button
        :disabled="deleting === 'batch'"
        class="inline-flex items-center gap-1.5 rounded-md bg-red-50 px-3 py-1.5 text-sm font-medium text-red-600 transition-colors hover:bg-red-100 disabled:opacity-50"
        @click="handleBatchDelete"
      >
        <Trash2 class="h-3.5 w-3.5" />
        {{ deleting === 'batch' ? t('common.deleting') : t('shortLinks.deleteSelected') }}
      </button>
    </div>

    <!-- Mobile card layout (hidden on md+) -->
    <div class="mt-4 space-y-3 md:hidden">
      <!-- Loading skeleton -->
      <template v-if="loading">
        <div v-for="i in 5" :key="`skel-${i}`" class="rounded-lg border bg-white p-4">
          <div class="flex items-center justify-between">
            <div class="h-4 w-24 animate-pulse rounded bg-gray-200" />
            <div class="h-5 w-16 animate-pulse rounded-full bg-gray-200" />
          </div>
          <div class="mt-2 h-3 w-full animate-pulse rounded bg-gray-100" />
          <div class="mt-2 flex items-center justify-between">
            <div class="h-3 w-20 animate-pulse rounded bg-gray-100" />
            <div class="h-6 w-6 animate-pulse rounded bg-gray-100" />
          </div>
        </div>
      </template>
      <!-- Empty state -->
      <div v-else-if="filteredLinks.length === 0 && total === 0" class="rounded-lg border bg-white py-12 text-center">
        <Link class="mx-auto h-10 w-10 text-gray-300" />
        <p class="mt-3 text-sm font-medium text-gray-500">{{ t('shortLinks.noLinksYet') }}</p>
        <router-link
          :to="{ name: 'short-link-create' }"
          class="mt-3 inline-flex items-center gap-1.5 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
        >
          <Plus class="h-4 w-4" />
          {{ t('shortLinks.createFirst') }}
        </router-link>
      </div>
      <div v-else-if="filteredLinks.length === 0" class="rounded-lg border bg-white py-8 text-center text-sm text-gray-400">
        {{ t('shortLinks.noMatch') }}
      </div>
      <!-- Link cards -->
      <div
        v-for="link in sortedLinks"
        :key="link.id"
        class="rounded-lg border bg-white transition-colors"
        :class="{ 'border-blue-200 bg-blue-50/30': selectedIds.has(link.id) }"
      >
        <div class="flex items-center justify-between p-4">
          <div class="flex items-center gap-2">
            <input
              type="checkbox"
              :checked="selectedIds.has(link.id)"
              class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
              @change="toggleSelect(link.id)"
            />
            <span class="font-mono text-sm font-medium text-gray-900">{{ link.id }}</span>
            <button
              class="rounded p-0.5 text-gray-400 hover:text-blue-600"
              :title="t('shortLinkCreate.copyTitle')"
              aria-label="Copy short link"
              @click="copyLink(link.id)"
            >
              <Copy v-if="copiedId !== link.id" class="h-3.5 w-3.5" />
              <Check v-else class="h-3.5 w-3.5 text-green-500" />
            </button>
          </div>
          <button
            :disabled="togglingId === link.id"
            class="inline-flex items-center gap-1.5 rounded-full px-2 py-0.5 text-xs font-medium"
            :class="link.isEnable ? 'bg-green-50 text-green-700' : 'bg-gray-100 text-gray-500'"
            @click="handleToggleEnable(link)"
          >
            <Loader2 v-if="togglingId === link.id" class="h-3 w-3 animate-spin" />
            <span v-else class="h-1.5 w-1.5 rounded-full" :class="link.isEnable ? 'bg-green-500' : 'bg-gray-400'" />
            {{ link.isEnable ? t('common.on') : t('common.off') }}
          </button>
        </div>
        <div class="border-t px-4 py-3">
          <a
            :href="link.url"
            target="_blank"
            rel="noopener noreferrer"
            class="flex items-center gap-1 text-sm text-blue-600 hover:underline"
          >
            <span class="truncate">{{ link.url }}</span>
            <ExternalLink class="h-3 w-3 shrink-0" />
          </a>
          <p v-if="link.description" class="mt-1 truncate text-xs text-gray-400">{{ link.description }}</p>
        </div>
        <div class="flex items-center justify-between border-t px-4 py-2">
          <span class="text-xs text-gray-400">{{ formatDate(link.createTime) }}</span>
          <div class="flex items-center gap-1">
            <button
              class="rounded p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600"
              :title="t('common.viewDetails')"
              aria-label="View details"
              @click="router.push({ name: 'short-link-detail', params: { id: link.id } })"
            >
              <Eye class="h-4 w-4" />
            </button>
            <button
              class="rounded p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600"
              :title="t('common.edit')"
              aria-label="Edit"
              @click="router.push({ name: 'short-link-edit', params: { id: link.id } })"
            >
              <Pencil class="h-4 w-4" />
            </button>
            <button
              class="rounded p-1.5 text-gray-400 hover:bg-red-50 hover:text-red-600"
              :title="t('common.delete')"
              aria-label="Delete"
              @click="confirmDeleteId = link.id"
            >
              <Trash2 class="h-4 w-4" />
            </button>
          </div>
        </div>
      </div>

      <!-- Mobile pagination -->
      <div
        v-if="totalPages > 1"
        class="mt-3 flex items-center justify-between"
      >
        <span class="text-xs text-gray-400">
          {{ (page - 1) * pageSize + 1 }}–{{ Math.min(page * pageSize, total) }} of {{ total }}
        </span>
        <div class="flex items-center gap-1">
          <button
            :disabled="page <= 1"
            class="rounded p-1.5 text-gray-400 transition-colors hover:bg-gray-100 disabled:opacity-30"
            aria-label="Previous page"
            @click="goPage(page - 1)"
          >
            <ChevronLeft class="h-5 w-5" />
          </button>
          <span class="min-w-[60px] text-center text-sm font-medium text-gray-700">{{ page }} / {{ totalPages }}</span>
          <button
            :disabled="page >= totalPages"
            class="rounded p-1.5 text-gray-400 transition-colors hover:bg-gray-100 disabled:opacity-30"
            aria-label="Next page"
            @click="goPage(page + 1)"
          >
            <ChevronRight class="h-5 w-5" />
          </button>
        </div>
      </div>
    </div>

    <!-- Desktop table (hidden below md) -->
    <div class="mt-4 hidden overflow-hidden rounded-lg border bg-white md:block">
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr
              class="border-b bg-gray-50 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
            >
              <th class="w-10 px-4 py-3">
                <input
                  type="checkbox"
                  :checked="allSelected"
                  class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  @change="toggleSelectAll"
                />
              </th>
              <th class="cursor-pointer select-none px-4 py-3" @click="toggleSort('id')">
                <div class="flex items-center gap-1">
                  ID
                  <component :is="SortIcon({ field: 'id' })" class="h-3.5 w-3.5" />
                </div>
              </th>
              <th class="px-4 py-3">{{ ('shortLinkDetail.targetUrl') }}</th>
              <th class="hidden px-4 py-3 md:table-cell">{{ ('common.description') }}</th>
              <th class="px-4 py-3">{{ ('common.status') }}</th>
              <th
                class="hidden cursor-pointer select-none px-4 py-3 lg:table-cell"
                @click="toggleSort('createTime')"
              >
                <div class="flex items-center gap-1">
                  Created
                  <component :is="SortIcon({ field: 'createTime' })" class="h-3.5 w-3.5" />
                </div>
              </th>
              <th
                class="hidden cursor-pointer select-none px-4 py-3 lg:table-cell"
                @click="toggleSort('updateTime')"
              >
                <div class="flex items-center gap-1">
                  Updated
                  <component :is="SortIcon({ field: 'updateTime' })" class="h-3.5 w-3.5" />
                </div>
              </th>
              <th class="px-4 py-3 text-right">{{ ('common.actions') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y">
            <!-- Skeleton loading -->
            <template v-if="loading">
              <tr v-for="i in 5" :key="`skeleton-${i}`">
                <td class="px-4 py-3">
                  <div class="h-4 w-4 animate-pulse rounded bg-gray-200" />
                </td>
                <td class="px-4 py-3">
                  <div class="h-4 w-20 animate-pulse rounded bg-gray-200" />
                </td>
                <td class="px-4 py-3">
                  <div class="h-4 w-40 animate-pulse rounded bg-gray-200" />
                </td>
                <td class="hidden px-4 py-3 md:table-cell">
                  <div class="h-4 w-32 animate-pulse rounded bg-gray-200" />
                </td>
                <td class="px-4 py-3">
                  <div class="h-5 w-16 animate-pulse rounded-full bg-gray-200" />
                </td>
                <td class="hidden px-4 py-3 lg:table-cell">
                  <div class="h-4 w-28 animate-pulse rounded bg-gray-200" />
                </td>
                <td class="hidden px-4 py-3 lg:table-cell">
                  <div class="h-4 w-28 animate-pulse rounded bg-gray-200" />
                </td>
                <td class="px-4 py-3">
                  <div class="ml-auto h-4 w-20 animate-pulse rounded bg-gray-200" />
                </td>
              </tr>
            </template>
            <!-- Empty state -->
            <tr v-else-if="filteredLinks.length === 0 && total === 0">
              <td colspan="8" class="px-4 py-12 text-center">
                <Link class="mx-auto h-10 w-10 text-gray-300" />
                <p class="mt-3 text-sm font-medium text-gray-500">{{ t('shortLinks.noLinksYet') }}</p>
                <p class="mt-1 text-sm text-gray-400">
                  Create your first short link to get started.
                </p>
                <router-link
                  :to="{ name: 'short-link-create' }"
                  class="mt-3 inline-flex items-center gap-1.5 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700"
                >
                  <Plus class="h-4 w-4" />
                  Create your first link
                </router-link>
              </td>
            </tr>
            <tr v-else-if="filteredLinks.length === 0">
              <td colspan="8" class="px-4 py-8 text-center text-gray-400">
                {{ t('shortLinks.noMatch') }}
              </td>
            </tr>
            <!-- Data rows -->
            <tr
              v-for="link in sortedLinks"
              :key="link.id"
              class="transition-colors hover:bg-gray-50"
              :class="{ 'bg-blue-50/50': selectedIds.has(link.id) }"
            >
              <td class="px-4 py-3">
                <input
                  type="checkbox"
                  :checked="selectedIds.has(link.id)"
                  class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  @change="toggleSelect(link.id)"
                />
              </td>
              <td class="px-4 py-3">
                <div class="flex items-center gap-1.5">
                  <span class="font-mono text-xs font-medium text-gray-900">{{ link.id }}</span>
                  <button
                    class="rounded p-0.5 text-gray-400 transition-colors hover:text-blue-600"
                    :title="t('shortLinkCreate.copyTitle')"
                    @click="copyLink(link.id)"
                  >
                    <Copy v-if="copiedId !== link.id" class="h-3.5 w-3.5" />
                    <Check v-else class="h-3.5 w-3.5 text-green-500" />
                  </button>
                </div>
              </td>
              <td class="max-w-[200px] truncate px-4 py-3">
                <a
                  :href="link.url"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="inline-flex items-center gap-1 text-blue-600 hover:underline"
                >
                  <span class="truncate">{{ link.url }}</span>
                  <ExternalLink class="h-3 w-3 shrink-0" />
                </a>
              </td>
              <td class="hidden max-w-[180px] truncate px-4 py-3 text-gray-500 md:table-cell">
                {{ link.description || '—' }}
              </td>
              <td class="px-4 py-3">
                <button
                  :disabled="togglingId === link.id"
                  class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium transition-colors"
                  :class="
                    link.isEnable
                      ? 'bg-green-50 text-green-700 hover:bg-green-100'
                      : 'bg-gray-100 text-gray-500 hover:bg-gray-200'
                  "
                  @click="handleToggleEnable(link)"
                >
                  <Loader2 v-if="togglingId === link.id" class="h-3 w-3 animate-spin" />
                  <span
                    v-else
                    class="h-1.5 w-1.5 rounded-full"
                    :class="link.isEnable ? 'bg-green-500' : 'bg-gray-400'"
                  />
                  {{ link.isEnable ? 'Enabled' : 'Disabled' }}
                </button>
              </td>
              <td class="hidden whitespace-nowrap px-4 py-3 text-gray-500 lg:table-cell">
                {{ formatDate(link.createTime) }}
              </td>
              <td class="hidden whitespace-nowrap px-4 py-3 text-gray-500 lg:table-cell">
                {{ formatDate(link.updateTime) }}
              </td>
              <td class="px-4 py-3">
                <div class="flex items-center justify-end gap-1">
                  <button
                    class="rounded p-1.5 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600"
                    :title="t('common.viewDetails')"
                    @click="router.push({ name: 'short-link-detail', params: { id: link.id } })"
                  >
                    <Eye class="h-4 w-4" />
                  </button>
                  <button
                    class="rounded p-1.5 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600"
                    :title="t('common.edit')"
                    @click="router.push({ name: 'short-link-edit', params: { id: link.id } })"
                  >
                    <Pencil class="h-4 w-4" />
                  </button>
                  <button
                    class="rounded p-1.5 text-gray-400 transition-colors hover:bg-red-50 hover:text-red-600"
                    :title="t('common.delete')"
                    @click="confirmDeleteId = link.id"
                  >
                    <Trash2 class="h-4 w-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Enhanced Pagination -->
      <div
        v-if="totalPages > 1 || total > 0"
        class="flex flex-col items-center justify-between gap-3 border-t px-4 py-3 sm:flex-row"
      >
        <div class="flex items-center gap-2">
          <span class="text-sm text-gray-500">Show</span>
          <select
            v-model="pageSize"
            class="rounded-md border border-gray-300 px-2 py-1 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
            @change="handlePageSizeChange"
          >
            <option v-for="size in pageSizeOptions" :key="size" :value="size">
              {{ size }}
            </option>
          </select>
          <span class="text-sm text-gray-500">per page</span>
          <span class="ml-2 text-sm text-gray-400">
            {{ (page - 1) * pageSize + 1 }}–{{ Math.min(page * pageSize, total) }} of {{ total }}
          </span>
        </div>
        <div v-if="totalPages > 1" class="flex items-center gap-1">
          <button
            :disabled="page <= 1"
            class="rounded p-1.5 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600 disabled:opacity-30"
            @click="goPage(page - 1)"
          >
            <ChevronLeft class="h-5 w-5" />
          </button>
          <template v-for="(p, idx) in pageNumbers" :key="`page-${idx}`">
            <span v-if="p === '...'" class="px-1 text-sm text-gray-400">...</span>
            <button
              v-else
              :class="[
                'min-w-[32px] rounded px-2 py-1 text-sm font-medium transition-colors',
                p === page
                  ? 'bg-blue-600 text-white'
                  : 'text-gray-600 hover:bg-gray-100',
              ]"
              @click="goPage(p as number)"
            >
              {{ p }}
            </button>
          </template>
          <button
            :disabled="page >= totalPages"
            class="rounded p-1.5 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600 disabled:opacity-30"
            @click="goPage(page + 1)"
          >
            <ChevronRight class="h-5 w-5" />
          </button>
          <div class="ml-2 flex items-center gap-1">
            <span class="text-sm text-gray-400">Go to</span>
            <input
              v-model="pageJump"
              type="number"
              :min="1"
              :max="totalPages"
              class="w-14 rounded-md border border-gray-300 px-2 py-1 text-center text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
              @keyup.enter="handlePageJump"
            />
          </div>
        </div>
      </div>
    </div>

    <ConfirmDialog
      :open="!!confirmDeleteId"
      title="Delete Short Link"
      variant="danger"
      :confirm-text="deleting === confirmDeleteId ? 'Deleting...' : 'Delete'"
      :loading="deleting === confirmDeleteId"
      @confirm="handleDelete(confirmDeleteId!)"
      @cancel="confirmDeleteId = null"
    >
      <p class="mt-2 text-sm text-gray-500">
        Are you sure you want to delete
        <span class="font-mono font-medium text-gray-700">{{ confirmDeleteId }}</span
        >? This action cannot be undone.
      </p>
    </ConfirmDialog>
  </div>
</template>
