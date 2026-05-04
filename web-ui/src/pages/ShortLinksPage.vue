<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { listShortLinks, deleteShortLink, updateShortLink } from '@/api/short-link'
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
} from 'lucide-vue-next'

defineOptions({ name: 'ShortLinksPage' })

const router = useRouter()

const links = ref<ShortLinkData[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)

const searchId = ref('')
const searchUrl = ref('')
const filterEnabled = ref<'' | 'true' | 'false'>('')

const selectedIds = ref<Set<string>>(new Set())
const deleting = ref<string | null>(null)
const togglingId = ref<string | null>(null)
const copiedId = ref<string | null>(null)
const confirmDeleteId = ref<string | null>(null)

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

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize.value)))
const allSelected = computed(
  () =>
    filteredLinks.value.length > 0 && filteredLinks.value.every((l) => selectedIds.value.has(l.id)),
)

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
    if (filteredLinks.value.length === 0 && page.value > 1) {
      page.value--
    }
    confirmDeleteId.value = null
  } catch {
    // error handled by interceptor
  } finally {
    deleting.value = null
  }
}

async function handleBatchDelete() {
  if (selectedIds.value.size === 0) return
  deleting.value = 'batch'
  try {
    await Promise.all([...selectedIds.value].map((id) => deleteShortLink(id)))
    links.value = links.value.filter((l) => !selectedIds.value.has(l.id))
    total.value -= selectedIds.value.size
    selectedIds.value.clear()
  } catch {
    await fetchLinks()
  } finally {
    deleting.value = null
  }
}

function copyLink(id: string) {
  const url = `${window.location.origin}/${id}`
  navigator.clipboard.writeText(url).then(() => {
    copiedId.value = id
    setTimeout(() => {
      copiedId.value = null
    }, 2000)
  })
}

function goPage(p: number) {
  if (p < 1 || p > totalPages.value) return
  page.value = p
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
        <h1 class="text-xl font-semibold text-gray-900">Short Links</h1>
        <p class="mt-1 text-sm text-gray-500">{{ total }} link{{ total !== 1 ? 's' : '' }} total</p>
      </div>
      <router-link
        :to="{ name: 'short-link-create' }"
        class="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700"
      >
        <Plus class="h-4 w-4" />
        Create Link
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
          placeholder="Search by ID..."
          class="w-full rounded-md border border-gray-300 py-2 pl-9 pr-3 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
        />
      </div>
      <div class="relative flex-1">
        <Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
        <input
          v-model="searchUrl"
          type="text"
          placeholder="Search by URL..."
          class="w-full rounded-md border border-gray-300 py-2 pl-9 pr-3 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
        />
      </div>
      <select
        v-model="filterEnabled"
        class="rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
      >
        <option value="">All Status</option>
        <option value="true">Enabled</option>
        <option value="false">Disabled</option>
      </select>
    </div>

    <!-- Batch actions -->
    <div v-if="selectedIds.size > 0" class="mt-3 flex items-center gap-3">
      <span class="text-sm text-gray-600">{{ selectedIds.size }} selected</span>
      <button
        :disabled="deleting === 'batch'"
        class="inline-flex items-center gap-1.5 rounded-md bg-red-50 px-3 py-1.5 text-sm font-medium text-red-600 transition-colors hover:bg-red-100 disabled:opacity-50"
        @click="handleBatchDelete"
      >
        <Trash2 class="h-3.5 w-3.5" />
        {{ deleting === 'batch' ? 'Deleting...' : 'Delete Selected' }}
      </button>
    </div>

    <!-- Table -->
    <div class="mt-4 overflow-hidden rounded-lg border bg-white">
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
              <th class="px-4 py-3">ID</th>
              <th class="px-4 py-3">URL</th>
              <th class="hidden px-4 py-3 md:table-cell">Description</th>
              <th class="px-4 py-3">Status</th>
              <th class="hidden px-4 py-3 lg:table-cell">Created</th>
              <th class="hidden px-4 py-3 lg:table-cell">Updated</th>
              <th class="px-4 py-3 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y">
            <tr v-if="loading">
              <td colspan="8" class="px-4 py-8 text-center text-gray-400">Loading...</td>
            </tr>
            <tr v-else-if="filteredLinks.length === 0">
              <td colspan="8" class="px-4 py-8 text-center text-gray-400">No short links found.</td>
            </tr>
            <tr
              v-for="link in filteredLinks"
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
                    title="Copy short link"
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
                  <span
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
                    title="View details"
                    @click="router.push({ name: 'short-link-detail', params: { id: link.id } })"
                  >
                    <Eye class="h-4 w-4" />
                  </button>
                  <button
                    class="rounded p-1.5 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600"
                    title="Edit"
                    @click="router.push({ name: 'short-link-edit', params: { id: link.id } })"
                  >
                    <Pencil class="h-4 w-4" />
                  </button>
                  <button
                    class="rounded p-1.5 text-gray-400 transition-colors hover:bg-red-50 hover:text-red-600"
                    title="Delete"
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

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="flex items-center justify-between border-t px-4 py-3">
        <span class="text-sm text-gray-500">
          Page {{ page }} of {{ totalPages }} ({{ total }} items)
        </span>
        <div class="flex items-center gap-1">
          <button
            :disabled="page <= 1"
            class="rounded p-1.5 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600 disabled:opacity-30"
            @click="goPage(page - 1)"
          >
            <ChevronLeft class="h-5 w-5" />
          </button>
          <button
            :disabled="page >= totalPages"
            class="rounded p-1.5 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600 disabled:opacity-30"
            @click="goPage(page + 1)"
          >
            <ChevronRight class="h-5 w-5" />
          </button>
        </div>
      </div>
    </div>

    <!-- Delete confirmation modal -->
    <Teleport to="body">
      <div
        v-if="confirmDeleteId"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
        @click.self="confirmDeleteId = null"
      >
        <div class="w-full max-w-sm rounded-lg bg-white p-6 shadow-xl">
          <h3 class="text-lg font-semibold text-gray-900">Delete Short Link</h3>
          <p class="mt-2 text-sm text-gray-500">
            Are you sure you want to delete
            <span class="font-mono font-medium text-gray-700">{{ confirmDeleteId }}</span
            >? This action cannot be undone.
          </p>
          <div class="mt-4 flex justify-end gap-2">
            <button
              class="rounded-md px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-100"
              @click="confirmDeleteId = null"
            >
              Cancel
            </button>
            <button
              :disabled="deleting === confirmDeleteId"
              class="rounded-md bg-red-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-red-700 disabled:opacity-50"
              @click="handleDelete(confirmDeleteId!)"
            >
              {{ deleting === confirmDeleteId ? 'Deleting...' : 'Delete' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
