<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getShortLink, getShortLinkData } from '@/api/short-link'
import type { ShortLinkData, RequestHistory } from '@/types/api'
import {
  ArrowLeft,
  ExternalLink,
  Copy,
  Check,
  Pencil,
  Calendar,
  Globe,
  Monitor,
} from 'lucide-vue-next'

defineOptions({ name: 'ShortLinkDetailPage' })

const router = useRouter()
const route = useRoute()
const id = route.params.id as string

const link = ref<ShortLinkData | null>(null)
const histories = ref<RequestHistory[]>([])
const loading = ref(false)
const notFound = ref(false)
const copied = ref(false)

const daysAgo = ref(7)
const startDate = ref('')
const endDate = ref('')

function setDateRange(days: number) {
  daysAgo.value = days
  const end = new Date()
  const start = new Date()
  start.setDate(end.getDate() - days + 1)
  startDate.value = formatDateParam(start)
  endDate.value = formatDateParam(end)
  fetchHistory()
}

function formatDateParam(d: Date) {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
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

function formatShortDate(d: string) {
  return new Date(d).toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
}

// Aggregate histories by day for the visit trend chart
const dailyVisits = computed(() => {
  const map = new Map<string, number>()
  // Fill all days in range with 0
  const start = new Date(startDate.value)
  const end = new Date(endDate.value)
  for (let d = new Date(start); d <= end; d.setDate(d.getDate() + 1)) {
    map.set(formatDateParam(d), 0)
  }
  for (const h of histories.value) {
    const day = h.time.slice(0, 10)
    map.set(day, (map.get(day) || 0) + 1)
  }
  const entries = [...map.entries()].sort((a, b) => a[0].localeCompare(b[0]))
  return {
    labels: entries.map((e) => formatShortDate(e[0])),
    values: entries.map((e) => e[1]),
  }
})

const maxVisit = computed(() => Math.max(1, ...dailyVisits.value.values))

// Aggregate OS from user agents
const osDistribution = computed(() => {
  const map = new Map<string, number>()
  for (const h of histories.value) {
    const os = parseOS(h.ua)
    map.set(os, (map.get(os) || 0) + 1)
  }
  return [...map.entries()]
    .sort((a, b) => b[1] - a[1])
    .slice(0, 6)
    .map(([name, count]) => ({ name, count }))
})

const osMax = computed(() => Math.max(1, ...osDistribution.value.map((o) => o.count)))

// Aggregate IPs
const ipDistribution = computed(() => {
  const map = new Map<string, number>()
  for (const h of histories.value) {
    map.set(h.ip, (map.get(h.ip) || 0) + 1)
  }
  return [...map.entries()]
    .sort((a, b) => b[1] - a[1])
    .slice(0, 6)
    .map(([ip, count]) => ({ ip, count }))
})

const ipMax = computed(() => Math.max(1, ...ipDistribution.value.map((i) => i.count)))

const totalVisits = computed(() => histories.value.length)
const uniqueIps = computed(() => new Set(histories.value.map((h) => h.ip)).size)

function parseOS(ua: string): string {
  if (ua.includes('Windows')) return 'Windows'
  if (ua.includes('Mac OS X')) return 'macOS'
  if (ua.includes('Android')) return 'Android'
  if (ua.includes('iPhone') || ua.includes('iPad')) return 'iOS'
  if (ua.includes('Linux')) return 'Linux'
  if (ua.includes('CrOS')) return 'ChromeOS'
  return 'Other'
}

function copyLink() {
  if (!link.value) return
  const host = window.location.host
  navigator.clipboard.writeText(`${host}/${link.value.id}`).then(() => {
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  })
}

async function fetchData() {
  loading.value = true
  try {
    const data = await getShortLink(id)
    link.value = data.shortLink
  } catch {
    notFound.value = true
  } finally {
    loading.value = false
  }
}

async function fetchHistory() {
  if (!startDate.value || !endDate.value) return
  try {
    const data = await getShortLinkData(id, startDate.value, endDate.value)
    histories.value = data.histories || []
  } catch {
    histories.value = []
  }
}

onMounted(() => {
  setDateRange(7)
  fetchData()
})
</script>

<template>
  <div>
    <div class="flex items-center gap-3">
      <button
        class="rounded p-1.5 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600"
        @click="router.push({ name: 'short-links' })"
      >
        <ArrowLeft class="h-5 w-5" />
      </button>
      <div class="flex-1">
        <h1 class="text-xl font-semibold text-gray-900">Short Link Detail</h1>
        <p v-if="link" class="mt-0.5 text-sm text-gray-500">
          <span class="font-mono font-medium text-gray-700">{{ link.id }}</span>
        </p>
      </div>
      <button
        v-if="link"
        class="inline-flex items-center gap-1.5 rounded-md bg-blue-600 px-3 py-1.5 text-sm font-medium text-white transition-colors hover:bg-blue-700"
        @click="router.push({ name: 'short-link-edit', params: { id: link.id } })"
      >
        <Pencil class="h-3.5 w-3.5" />
        Edit
      </button>
    </div>

    <div v-if="loading" class="mt-6 text-center text-gray-400">Loading...</div>

    <div v-else-if="notFound" class="mt-6 text-center text-gray-400">
      Short link not found.
      <button
        class="ml-2 text-blue-600 hover:underline"
        @click="router.push({ name: 'short-links' })"
      >
        Back to list
      </button>
    </div>

    <template v-else-if="link">
      <!-- Basic info card -->
      <div class="mt-4 rounded-lg border bg-white p-5">
        <div class="grid gap-4 sm:grid-cols-2">
          <div>
            <label class="text-xs font-medium uppercase text-gray-400">Short Link ID</label>
            <div class="mt-1 flex items-center gap-2">
              <span class="font-mono text-sm font-medium text-gray-900">{{ link.id }}</span>
              <button
                class="rounded p-0.5 text-gray-400 transition-colors hover:text-blue-600"
                @click="copyLink"
              >
                <Copy v-if="!copied" class="h-3.5 w-3.5" />
                <Check v-else class="h-3.5 w-3.5 text-green-500" />
              </button>
            </div>
          </div>
          <div>
            <label class="text-xs font-medium uppercase text-gray-400">Target URL</label>
            <a
              :href="link.url"
              target="_blank"
              rel="noopener noreferrer"
              class="mt-1 flex items-center gap-1 text-sm text-blue-600 hover:underline"
            >
              {{ link.url }}
              <ExternalLink class="h-3 w-3 shrink-0" />
            </a>
          </div>
          <div>
            <label class="text-xs font-medium uppercase text-gray-400">Description</label>
            <p class="mt-1 text-sm text-gray-700">{{ link.description || '—' }}</p>
          </div>
          <div>
            <label class="text-xs font-medium uppercase text-gray-400">Status</label>
            <span
              class="mt-1 inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium"
              :class="link.isEnable ? 'bg-green-50 text-green-700' : 'bg-gray-100 text-gray-500'"
            >
              <span
                class="h-1.5 w-1.5 rounded-full"
                :class="link.isEnable ? 'bg-green-500' : 'bg-gray-400'"
              />
              {{ link.isEnable ? 'Enabled' : 'Disabled' }}
            </span>
          </div>
          <div>
            <label class="text-xs font-medium uppercase text-gray-400">Created By</label>
            <p class="mt-1 text-sm text-gray-700">{{ link.createdBy }}</p>
          </div>
          <div>
            <label class="text-xs font-medium uppercase text-gray-400">Created</label>
            <p class="mt-1 text-sm text-gray-700">{{ formatDate(link.createTime) }}</p>
          </div>
        </div>
      </div>

      <!-- Analytics section -->
      <div class="mt-6">
        <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <h2 class="text-lg font-semibold text-gray-900">Access Analytics</h2>
          <div class="flex items-center gap-1 rounded-lg border bg-white p-1">
            <button
              v-for="d in [7, 14, 30]"
              :key="d"
              class="rounded-md px-3 py-1 text-xs font-medium transition-colors"
              :class="daysAgo === d ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-100'"
              @click="setDateRange(d)"
            >
              {{ d }}d
            </button>
          </div>
        </div>

        <!-- Stats cards -->
        <div class="mt-4 grid grid-cols-2 gap-3 sm:grid-cols-4">
          <div class="rounded-lg border bg-white p-4">
            <div class="flex items-center gap-2 text-gray-400">
              <Calendar class="h-4 w-4" />
              <span class="text-xs font-medium uppercase">Total Visits</span>
            </div>
            <p class="mt-2 text-2xl font-bold text-gray-900">{{ totalVisits }}</p>
          </div>
          <div class="rounded-lg border bg-white p-4">
            <div class="flex items-center gap-2 text-gray-400">
              <Globe class="h-4 w-4" />
              <span class="text-xs font-medium uppercase">Unique IPs</span>
            </div>
            <p class="mt-2 text-2xl font-bold text-gray-900">{{ uniqueIps }}</p>
          </div>
          <div class="rounded-lg border bg-white p-4">
            <div class="flex items-center gap-2 text-gray-400">
              <Monitor class="h-4 w-4" />
              <span class="text-xs font-medium uppercase">Platforms</span>
            </div>
            <p class="mt-2 text-2xl font-bold text-gray-900">{{ osDistribution.length }}</p>
          </div>
          <div class="rounded-lg border bg-white p-4">
            <div class="flex items-center gap-2 text-gray-400">
              <Calendar class="h-4 w-4" />
              <span class="text-xs font-medium uppercase">Avg/Day</span>
            </div>
            <p class="mt-2 text-2xl font-bold text-gray-900">
              {{ daysAgo > 0 ? (totalVisits / daysAgo).toFixed(1) : 0 }}
            </p>
          </div>
        </div>

        <!-- Charts row -->
        <div class="mt-4 grid gap-4 lg:grid-cols-2">
          <!-- Visit trend chart -->
          <div class="rounded-lg border bg-white p-5">
            <h3 class="mb-4 text-sm font-semibold text-gray-700">Visit Trend</h3>
            <div
              v-if="dailyVisits.values.length === 0"
              class="py-8 text-center text-sm text-gray-400"
            >
              No data for this period.
            </div>
            <div v-else class="relative h-48 w-full">
              <!-- Y-axis labels -->
              <div
                class="absolute left-0 top-0 flex h-full flex-col justify-between text-right text-xs text-gray-400"
              >
                <span>{{ maxVisit }}</span>
                <span>{{ Math.round(maxVisit / 2) }}</span>
                <span>0</span>
              </div>
              <!-- Chart area -->
              <div class="ml-8 flex h-full items-end gap-1">
                <div
                  v-for="(val, i) in dailyVisits.values"
                  :key="i"
                  class="group relative flex-1 cursor-default"
                >
                  <div
                    class="w-full rounded-t bg-blue-500 transition-colors hover:bg-blue-600"
                    :style="{
                      height: `${(val / maxVisit) * 100}%`,
                      minHeight: val > 0 ? '2px' : '0',
                    }"
                  />
                  <span
                    v-if="dailyVisits.labels.length <= 14 || i % 2 === 0"
                    class="mt-1 block text-center text-[10px] text-gray-400"
                  >
                    {{ dailyVisits.labels[i] }}
                  </span>
                  <!-- Tooltip -->
                  <div
                    class="pointer-events-none absolute bottom-full left-1/2 mb-1 -translate-x-1/2 rounded bg-gray-800 px-2 py-1 text-xs text-white opacity-0 transition-opacity group-hover:opacity-100"
                  >
                    {{ val }} visit{{ val !== 1 ? 's' : '' }}
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- OS distribution -->
          <div class="rounded-lg border bg-white p-5">
            <h3 class="mb-4 text-sm font-semibold text-gray-700">OS Distribution</h3>
            <div v-if="osDistribution.length === 0" class="py-8 text-center text-sm text-gray-400">
              No data for this period.
            </div>
            <div v-else class="space-y-3">
              <div v-for="os in osDistribution" :key="os.name" class="flex items-center gap-3">
                <span class="w-20 shrink-0 truncate text-sm text-gray-700">{{ os.name }}</span>
                <div class="flex-1">
                  <div class="h-5 overflow-hidden rounded-full bg-gray-100">
                    <div
                      class="h-full rounded-full bg-indigo-500 transition-all"
                      :style="{ width: `${(os.count / osMax) * 100}%` }"
                    />
                  </div>
                </div>
                <span class="w-10 text-right text-sm font-medium text-gray-600">{{
                  os.count
                }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- IP distribution -->
        <div class="mt-4 rounded-lg border bg-white p-5">
          <h3 class="mb-4 text-sm font-semibold text-gray-700">Top IPs</h3>
          <div v-if="ipDistribution.length === 0" class="py-8 text-center text-sm text-gray-400">
            No data for this period.
          </div>
          <div v-else>
            <div class="space-y-3">
              <div v-for="item in ipDistribution" :key="item.ip" class="flex items-center gap-3">
                <span class="w-36 shrink-0 truncate font-mono text-sm text-gray-700">{{
                  item.ip
                }}</span>
                <div class="flex-1">
                  <div class="h-5 overflow-hidden rounded-full bg-gray-100">
                    <div
                      class="h-full rounded-full bg-emerald-500 transition-all"
                      :style="{ width: `${(item.count / ipMax) * 100}%` }"
                    />
                  </div>
                </div>
                <span class="w-10 text-right text-sm font-medium text-gray-600">{{
                  item.count
                }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
