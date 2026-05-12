<script setup lang="ts">
import { ref, computed, onMounted, defineAsyncComponent, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { getShortLink, getShortLinkData } from '@/api/short-link'
import { listDomains } from '@/api/tenant'
import { useAuthStore } from '@/stores/auth'
import { generateQRCodeDataURL, downloadQRCode } from '@/utils/qrcode'
import type { ShortLinkData, RequestHistory, TenantDomain } from '@/types/api'
import {
  ArrowLeft,
  ExternalLink,
  Copy,
  Check,
  Pencil,
  Calendar,
  Globe,
  Monitor,
  Loader2,
  LayoutGrid,
  FileText,
  QrCode,
  Download,
} from 'lucide-vue-next'

const { t } = useI18n()

const VChart = defineAsyncComponent(async () => {
  const [
    { use },
    { CanvasRenderer },
    { BarChart, PieChart },
    { GridComponent, TooltipComponent, LegendComponent },
  ] = await Promise.all([
    import('echarts/core'),
    import('echarts/renderers'),
    import('echarts/charts'),
    import('echarts/components'),
  ])
  use([CanvasRenderer, BarChart, PieChart, GridComponent, TooltipComponent, LegendComponent])
  return (await import('vue-echarts')).default
})

defineOptions({ name: 'ShortLinkDetailPage' })

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const id = route.params.id as string

const link = ref<ShortLinkData | null>(null)
const histories = ref<RequestHistory[]>([])
const landingHost = ref('')
const loading = ref(false)
const chartLoading = ref(false)
const notFound = ref(false)
const copied = ref(false)
const activeTab = ref<'trend' | 'records'>('trend')

const qrDomains = ref<TenantDomain[]>([])
const qrSelectedDomain = ref('')
const qrDataUrl = ref('')
const qrLoading = ref(false)
const qrCopied = ref(false)
const qrLoadError = ref(false)

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

const shortLinkUrl = computed(() => {
  if (!link.value) return ''
  if (landingHost.value) {
    const protocol = landingHost.value.startsWith('localhost')
      ? 'http://'
      : `${window.location.protocol}//`
    return `${protocol}${landingHost.value}/${link.value.id}`
  }
  return `${window.location.origin}/${link.value.id}`
})

const qrUrl = computed(() => {
  if (!link.value || !qrSelectedDomain.value) return ''
  const protocol = qrSelectedDomain.value.startsWith('localhost')
    ? 'http://'
    : `${window.location.protocol}//`
  return `${protocol}${qrSelectedDomain.value}/${link.value.id}`
})

const qrHasMultipleDomains = computed(() => qrDomains.value.length > 1)

const dailyVisits = computed(() => {
  const map = new Map<string, number>()
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

const browserDistribution = computed(() => {
  const map = new Map<string, number>()
  for (const h of histories.value) {
    const browser = h.browser || parseBrowser(h.ua)
    if (browser) {
      map.set(browser, (map.get(browser) || 0) + 1)
    }
  }
  return [...map.entries()]
    .sort((a, b) => b[1] - a[1])
    .slice(0, 6)
    .map(([name, count]) => ({ name, count }))
})

const refererDistribution = computed(() => {
  const map = new Map<string, number>()
  for (const h of histories.value) {
    const ref = h.referer || 'Direct'
    const source = extractRefererSource(ref)
    map.set(source, (map.get(source) || 0) + 1)
  }
  return [...map.entries()]
    .sort((a, b) => b[1] - a[1])
    .slice(0, 8)
    .map(([name, count]) => ({ name, count }))
})

const ipDistribution = computed(() => {
  const map = new Map<string, number>()
  for (const h of histories.value) {
    map.set(h.ip, (map.get(h.ip) || 0) + 1)
  }
  const total = histories.value.length || 1
  return [...map.entries()]
    .sort((a, b) => b[1] - a[1])
    .slice(0, 10)
    .map(([ip, count]) => ({ ip, count, percent: ((count / total) * 100).toFixed(1) }))
})

const totalVisits = computed(() => histories.value.length)
const uniqueIps = computed(() => new Set(histories.value.map((h) => h.ip)).size)

function parseOS(ua: string): string {
  if (ua.includes('Windows')) return 'Windows'
  if (ua.includes('iPhone') || ua.includes('iPad')) return 'iOS'
  if (ua.includes('Mac OS X')) return 'macOS'
  if (ua.includes('Android')) return 'Android'
  if (ua.includes('CrOS')) return 'ChromeOS'
  if (ua.includes('Linux')) return 'Linux'
  return 'Other'
}

function parseBrowser(ua: string): string {
  if (ua.includes('Edg/')) return 'Edge'
  if (ua.includes('OPR/') || ua.includes('Opera')) return 'Opera'
  if (ua.includes('Firefox/')) return 'Firefox'
  if (ua.includes('Chrome/')) return 'Chrome'
  if (ua.includes('Safari/') && !ua.includes('Chrome')) return 'Safari'
  return 'Other'
}

function extractRefererSource(referer: string): string {
  if (!referer || referer === 'Direct') return 'Direct'
  try {
    const url = new URL(referer)
    return url.hostname.replace(/^www\./, '')
  } catch {
    return referer.length > 40 ? referer.slice(0, 40) + '...' : referer
  }
}

const tooltipBase = {
  backgroundColor: '#fff',
  borderColor: '#e5e7eb',
  borderWidth: 1,
  textStyle: { color: '#374151', fontSize: 13 },
}

const visitChartOption = computed(() => ({
  tooltip: {
    ...tooltipBase,
    trigger: 'axis' as const,
    formatter: (params: { name: string; value: number }[]) => {
      const p = params[0]
      return `${p.name}<br/><b>${p.value}</b> ${t('shortLinkDetail.visitsWithCount', { count: p.value })}`
    },
  },
  grid: { top: 16, right: 16, bottom: 32, left: 48 },
  xAxis: {
    type: 'category' as const,
    data: dailyVisits.value.labels,
    axisLine: { lineStyle: { color: '#e5e7eb' } },
    axisTick: { show: false },
    axisLabel: { color: '#9ca3af', fontSize: 11 },
  },
  yAxis: {
    type: 'value' as const,
    splitLine: { lineStyle: { color: '#f3f4f6', type: 'dashed' as const } },
    axisLine: { show: false },
    axisTick: { show: false },
    axisLabel: { color: '#9ca3af', fontSize: 11 },
  },
  series: [
    {
      type: 'bar' as const,
      data: dailyVisits.value.values,
      barWidth: '60%',
      itemStyle: {
        color: '#2563eb',
        borderRadius: [4, 4, 0, 0],
      },
    },
  ],
}))

function makePieOption(
  data: { name: string; count: number }[],
  formatter: (name: string, value: number, total: number) => string,
) {
  const total = data.reduce((s, d) => s + d.count, 0) || 1
  return {
    tooltip: {
      ...tooltipBase,
      trigger: 'item' as const,
      formatter: (params: { name: string; value: number }) =>
        formatter(params.name, params.value, total),
    },
    legend: {
      bottom: 0,
      textStyle: { color: '#6b7280', fontSize: 12 },
      itemWidth: 12,
      itemHeight: 12,
    },
    series: [
      {
        type: 'pie' as const,
        radius: ['40%', '70%'],
        center: ['50%', '45%'],
        avoidLabelOverlap: false,
        label: { show: false },
        emphasis: {
          label: { show: true, fontSize: 14, fontWeight: 'bold' as const },
        },
        data: data.map((d) => ({ name: d.name, value: d.count })),
        color: [
          '#2563eb',
          '#8b5cf6',
          '#06b6d4',
          '#f59e0b',
          '#ef4444',
          '#6b7280',
          '#10b981',
          '#ec4899',
        ],
      },
    ],
  }
}

const osChartOption = computed(() =>
  makePieOption(
    osDistribution.value,
    (name, value, total) =>
      `${name}<br/><b>${value}</b> ${t('shortLinkDetail.visitsWithPercent', { count: value, percent: ((value / total) * 100).toFixed(1) })}`,
  ),
)

const browserChartOption = computed(() =>
  makePieOption(
    browserDistribution.value,
    (name, value, total) =>
      `${name}<br/><b>${value}</b> ${t('shortLinkDetail.visitsWithPercent', { count: value, percent: ((value / total) * 100).toFixed(1) })}`,
  ),
)

const refererChartOption = computed(() =>
  makePieOption(
    refererDistribution.value,
    (name, value, total) =>
      `${name}<br/><b>${value}</b> ${t('shortLinkDetail.visitsWithPercent', { count: value, percent: ((value / total) * 100).toFixed(1) })}`,
  ),
)

function copyLink() {
  if (!link.value) return
  navigator.clipboard.writeText(shortLinkUrl.value).then(
    () => {
      copied.value = true
      setTimeout(() => {
        copied.value = false
      }, 2000)
    },
    () => {},
  )
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

async function fetchLandingHost() {
  if (!auth.currentTenantId) return
  try {
    const domains = await listDomains(auth.currentTenantId)
    qrDomains.value = domains
    qrLoadError.value = false
    const defaultDomain = (domains as TenantDomain[]).find((d) => d.isDefault)
    if (defaultDomain) {
      landingHost.value = defaultDomain.domain
      qrSelectedDomain.value = defaultDomain.domain
    } else if (domains.length > 0) {
      landingHost.value = (domains as TenantDomain[])[0].domain
      qrSelectedDomain.value = (domains as TenantDomain[])[0].domain
    }
  } catch {
    qrLoadError.value = true
  }
}

async function generateQR() {
  if (!qrUrl.value) {
    qrDataUrl.value = ''
    return
  }
  qrLoading.value = true
  try {
    qrDataUrl.value = await generateQRCodeDataURL(qrUrl.value)
  } catch {
    qrDataUrl.value = ''
  } finally {
    qrLoading.value = false
  }
}

function handleDownloadQR() {
  if (!qrDataUrl.value || !link.value) return
  downloadQRCode(qrDataUrl.value, `qrcode-${link.value.id}.png`)
}

function copyQRLink() {
  navigator.clipboard.writeText(qrUrl.value).then(
    () => {
      qrCopied.value = true
      setTimeout(() => {
        qrCopied.value = false
      }, 2000)
    },
    () => {},
  )
}

async function fetchHistory() {
  if (!startDate.value || !endDate.value) return
  chartLoading.value = true
  try {
    const data = await getShortLinkData(id, startDate.value, endDate.value)
    histories.value = data.histories || []
  } catch {
    histories.value = []
  } finally {
    chartLoading.value = false
  }
}

onMounted(() => {
  setDateRange(7)
  fetchData()
  fetchLandingHost()
})

watch([() => link.value, qrSelectedDomain], () => {
  generateQR()
})
</script>

<template>
  <div>
    <div class="flex items-center gap-3">
      <button
        class="rounded p-1.5 text-gray-400 dark:text-gray-500 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800 hover:text-gray-600 dark:hover:text-gray-300"
        @click="router.push({ name: 'short-links' })"
      >
        <ArrowLeft class="h-5 w-5" />
      </button>
      <div class="flex-1">
        <h1 class="text-xl font-semibold text-gray-900 dark:text-white">
          {{ t('shortLinkDetail.title') }}
        </h1>
        <p v-if="link" class="mt-0.5 text-sm text-gray-500 dark:text-gray-400">
          <span class="font-mono font-medium text-blue-600">{{ shortLinkUrl }}</span>
        </p>
      </div>
      <button
        v-if="link"
        class="inline-flex items-center gap-1.5 rounded-md bg-blue-600 px-3 py-1.5 text-sm font-medium text-white transition-colors hover:bg-blue-700"
        @click="router.push({ name: 'short-link-edit', params: { id: link.id } })"
      >
        <Pencil class="h-3.5 w-3.5" />
        {{ t('shortLinkDetail.edit') }}
      </button>
    </div>

    <div v-if="loading" class="mt-6 text-center text-gray-400 dark:text-gray-500">
      {{ t('shortLinkDetail.loading') }}
    </div>

    <div v-else-if="notFound" class="mt-6 text-center text-gray-400 dark:text-gray-500">
      {{ t('shortLinkDetail.notFound') }}
      <button
        class="ml-2 text-blue-600 hover:underline"
        @click="router.push({ name: 'short-links' })"
      >
        {{ t('common.backToList') }}
      </button>
    </div>

    <template v-else-if="link">
      <!-- Basic info card -->
      <div class="mt-4 rounded-lg border dark:border-gray-700 bg-white dark:bg-gray-900 p-5">
        <div class="grid gap-4 sm:grid-cols-2">
          <div>
            <label class="text-xs font-medium uppercase text-gray-400 dark:text-gray-500">{{
              t('shortLinkDetail.shortLink')
            }}</label>
            <div class="mt-1 flex items-center gap-2">
              <a
                :href="shortLinkUrl"
                target="_blank"
                rel="noopener noreferrer"
                class="font-mono text-sm font-medium text-blue-600 hover:underline"
              >
                {{ shortLinkUrl }}
              </a>
              <button
                class="rounded p-0.5 text-gray-400 dark:text-gray-500 transition-colors hover:text-blue-600"
                @click="copyLink"
              >
                <Copy v-if="!copied" class="h-3.5 w-3.5" />
                <Check v-else class="h-3.5 w-3.5 text-green-500" />
              </button>
            </div>
          </div>
          <div>
            <label class="text-xs font-medium uppercase text-gray-400 dark:text-gray-500">{{
              t('shortLinkDetail.targetUrl')
            }}</label>
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
            <label class="text-xs font-medium uppercase text-gray-400 dark:text-gray-500">{{
              t('shortLinkDetail.description')
            }}</label>
            <p class="mt-1 text-sm text-gray-700 dark:text-gray-300">
              {{ link.description || '—' }}
            </p>
          </div>
          <div>
            <label class="text-xs font-medium uppercase text-gray-400 dark:text-gray-500">{{
              t('shortLinkDetail.status')
            }}</label>
            <span
              class="mt-1 inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium"
              :class="
                link.isEnable
                  ? 'bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-400'
                  : 'bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400'
              "
            >
              <span
                class="h-1.5 w-1.5 rounded-full"
                :class="link.isEnable ? 'bg-green-500' : 'bg-gray-400'"
              />
              {{ link.isEnable ? t('common.enabled') : t('common.disabled') }}
            </span>
          </div>
          <div>
            <label class="text-xs font-medium uppercase text-gray-400 dark:text-gray-500">{{
              t('shortLinkDetail.createdBy')
            }}</label>
            <p class="mt-1 text-sm text-gray-700 dark:text-gray-300">{{ link.createdBy }}</p>
          </div>
          <div>
            <label class="text-xs font-medium uppercase text-gray-400 dark:text-gray-500">{{
              t('shortLinkDetail.created')
            }}</label>
            <p class="mt-1 text-sm text-gray-700 dark:text-gray-300">
              {{ formatDate(link.createTime) }}
            </p>
          </div>
        </div>
      </div>

      <!-- QR Code card -->
      <div class="mt-4 rounded-lg border dark:border-gray-700 bg-white dark:bg-gray-900 p-5">
        <h3
          class="mb-4 flex items-center gap-2 text-sm font-semibold text-gray-700 dark:text-gray-300"
        >
          <QrCode class="h-4 w-4" />
          {{ t('qrCode.title') }}
        </h3>
        <template v-if="qrLoadError">
          <div class="py-6 text-center">
            <p class="text-sm text-red-500">{{ t('qrCode.loadError') }}</p>
            <button class="mt-2 text-sm text-blue-600 hover:underline" @click="fetchLandingHost()">
              {{ t('common.refresh') }}
            </button>
          </div>
        </template>
        <template v-else-if="qrDomains.length === 0">
          <div class="py-6 text-center">
            <QrCode class="mx-auto h-8 w-8 text-gray-300 dark:text-gray-600" />
            <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">{{ t('qrCode.noDomains') }}</p>
          </div>
        </template>
        <template v-else>
          <div class="flex flex-col items-center gap-4 sm:flex-row sm:items-start">
            <div class="flex shrink-0 items-center justify-center">
              <div v-if="qrLoading" class="flex h-[192px] w-[192px] items-center justify-center">
                <Loader2 class="h-6 w-6 animate-spin text-gray-400" />
              </div>
              <img
                v-else-if="qrDataUrl"
                :src="qrDataUrl"
                :alt="t('qrCode.imageAlt')"
                class="h-[192px] w-[192px] rounded"
              />
            </div>
            <div class="flex w-full flex-col gap-3">
              <div v-if="qrHasMultipleDomains">
                <label class="text-xs font-medium uppercase text-gray-400 dark:text-gray-500">
                  {{ t('common.domain') }}
                </label>
                <select
                  v-model="qrSelectedDomain"
                  class="mt-1 w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
                >
                  <option v-for="d in qrDomains" :key="d.id" :value="d.domain">
                    {{ d.domain }}{{ d.isDefault ? ` (${t('common.default')})` : '' }}
                  </option>
                </select>
              </div>
              <div class="flex flex-wrap gap-2">
                <button
                  :disabled="!qrDataUrl"
                  class="inline-flex items-center gap-1.5 rounded-md bg-blue-600 px-3 py-1.5 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:opacity-50"
                  @click="handleDownloadQR"
                >
                  <Download class="h-3.5 w-3.5" />
                  {{ t('qrCode.downloadPng') }}
                </button>
                <button
                  class="inline-flex items-center gap-1.5 rounded-md border border-gray-300 dark:border-gray-600 px-3 py-1.5 text-sm font-medium text-gray-700 dark:text-gray-300 transition-colors hover:bg-gray-50 dark:hover:bg-gray-800"
                  @click="copyQRLink"
                >
                  <Copy v-if="!qrCopied" class="h-3.5 w-3.5" />
                  <Check v-else class="h-3.5 w-3.5 text-green-500" />
                  {{ qrCopied ? t('qrCode.copied') : t('qrCode.copyLink') }}
                </button>
              </div>
            </div>
          </div>
        </template>
      </div>

      <!-- Analytics section -->
      <div class="mt-6">
        <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ t('shortLinkDetail.accessAnalytics') }}
          </h2>
          <div class="flex items-center gap-2">
            <div
              class="flex items-center gap-1 rounded-lg border dark:border-gray-700 bg-white dark:bg-gray-900 p-1"
            >
              <button
                v-for="d in [7, 14, 30]"
                :key="d"
                class="rounded-md px-3 py-1 text-xs font-medium transition-colors"
                :class="
                  daysAgo === d
                    ? 'bg-blue-600 text-white'
                    : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800'
                "
                @click="setDateRange(d)"
              >
                {{ d }}d
              </button>
            </div>
          </div>
        </div>

        <!-- Stats cards -->
        <div class="mt-4 grid grid-cols-2 gap-3 sm:grid-cols-4">
          <div class="rounded-lg border dark:border-gray-700 bg-white dark:bg-gray-900 p-4">
            <div class="flex items-center gap-2 text-gray-400 dark:text-gray-500">
              <Calendar class="h-4 w-4" />
              <span class="text-xs font-medium uppercase">{{
                t('shortLinkDetail.totalVisits')
              }}</span>
            </div>
            <p class="mt-2 text-2xl font-bold text-gray-900 dark:text-white">{{ totalVisits }}</p>
          </div>
          <div class="rounded-lg border dark:border-gray-700 bg-white dark:bg-gray-900 p-4">
            <div class="flex items-center gap-2 text-gray-400 dark:text-gray-500">
              <Globe class="h-4 w-4" />
              <span class="text-xs font-medium uppercase">{{
                t('shortLinkDetail.uniqueIps')
              }}</span>
            </div>
            <p class="mt-2 text-2xl font-bold text-gray-900 dark:text-white">{{ uniqueIps }}</p>
          </div>
          <div class="rounded-lg border dark:border-gray-700 bg-white dark:bg-gray-900 p-4">
            <div class="flex items-center gap-2 text-gray-400 dark:text-gray-500">
              <Monitor class="h-4 w-4" />
              <span class="text-xs font-medium uppercase">{{
                t('shortLinkDetail.platforms')
              }}</span>
            </div>
            <p class="mt-2 text-2xl font-bold text-gray-900 dark:text-white">
              {{ osDistribution.length }}
            </p>
          </div>
          <div class="rounded-lg border dark:border-gray-700 bg-white dark:bg-gray-900 p-4">
            <div class="flex items-center gap-2 text-gray-400 dark:text-gray-500">
              <Calendar class="h-4 w-4" />
              <span class="text-xs font-medium uppercase">{{ t('shortLinkDetail.avgDay') }}</span>
            </div>
            <p class="mt-2 text-2xl font-bold text-gray-900 dark:text-white">
              {{ daysAgo > 0 ? (totalVisits / daysAgo).toFixed(1) : 0 }}
            </p>
          </div>
        </div>

        <!-- Tab switcher -->
        <div
          class="mt-6 flex gap-1 rounded-lg border dark:border-gray-700 bg-white dark:bg-gray-900 p-1 w-fit"
        >
          <button
            class="inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-xs font-medium transition-colors"
            :class="
              activeTab === 'trend'
                ? 'bg-blue-600 text-white'
                : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800'
            "
            @click="activeTab = 'trend'"
          >
            <LayoutGrid class="h-3.5 w-3.5" />
            {{ t('shortLinkDetail.charts') }}
          </button>
          <button
            class="inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-xs font-medium transition-colors"
            :class="
              activeTab === 'records'
                ? 'bg-blue-600 text-white'
                : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800'
            "
            @click="activeTab = 'records'"
          >
            <FileText class="h-3.5 w-3.5" />
            {{ t('shortLinkDetail.accessRecords') }}
          </button>
        </div>

        <!-- Charts tab -->
        <template v-if="activeTab === 'trend'">
          <!-- Visit trend chart -->
          <div class="mt-4 rounded-lg border dark:border-gray-700 bg-white dark:bg-gray-900 p-5">
            <h3 class="mb-4 text-sm font-semibold text-gray-700 dark:text-gray-300">
              {{ t('shortLinkDetail.visitTrend') }}
            </h3>
            <div v-if="chartLoading" class="flex h-56 items-center justify-center">
              <Loader2 class="h-5 w-5 animate-spin text-gray-400 dark:text-gray-500" />
            </div>
            <div
              v-else-if="totalVisits === 0"
              class="flex h-56 items-center justify-center text-sm text-gray-400 dark:text-gray-500"
            >
              {{ t('shortLinkDetail.noData') }}
            </div>
            <div v-else class="h-56">
              <VChart
                :option="visitChartOption"
                :autoresize="true"
                class="h-full w-full"
                role="img"
                aria-label="Visit trend bar chart"
              />
            </div>
          </div>

          <!-- OS & Browser distribution charts -->
          <div class="mt-4 grid gap-4 lg:grid-cols-2">
            <div class="rounded-lg border dark:border-gray-700 bg-white dark:bg-gray-900 p-5">
              <h3 class="mb-4 text-sm font-semibold text-gray-700 dark:text-gray-300">
                {{ t('shortLinkDetail.osDistribution') }}
              </h3>
              <div v-if="chartLoading" class="flex h-56 items-center justify-center">
                <Loader2 class="h-5 w-5 animate-spin text-gray-400 dark:text-gray-500" />
              </div>
              <div
                v-else-if="osDistribution.length === 0"
                class="flex h-56 items-center justify-center text-sm text-gray-400 dark:text-gray-500"
              >
                {{ t('shortLinkDetail.noData') }}
              </div>
              <div v-else class="h-56">
                <VChart
                  :option="osChartOption"
                  :autoresize="true"
                  class="h-full w-full"
                  role="img"
                  aria-label="OS distribution pie chart"
                />
              </div>
            </div>

            <div class="rounded-lg border dark:border-gray-700 bg-white dark:bg-gray-900 p-5">
              <h3 class="mb-4 text-sm font-semibold text-gray-700 dark:text-gray-300">
                {{ t('shortLinkDetail.browserDistribution') }}
              </h3>
              <div v-if="chartLoading" class="flex h-56 items-center justify-center">
                <Loader2 class="h-5 w-5 animate-spin text-gray-400 dark:text-gray-500" />
              </div>
              <div
                v-else-if="browserDistribution.length === 0"
                class="flex h-56 items-center justify-center text-sm text-gray-400 dark:text-gray-500"
              >
                {{ t('shortLinkDetail.noData') }}
              </div>
              <div v-else class="h-56">
                <VChart
                  :option="browserChartOption"
                  :autoresize="true"
                  class="h-full w-full"
                  role="img"
                  aria-label="Browser distribution pie chart"
                />
              </div>
            </div>
          </div>

          <!-- Referer distribution -->
          <div class="mt-4 rounded-lg border dark:border-gray-700 bg-white dark:bg-gray-900 p-5">
            <h3 class="mb-4 text-sm font-semibold text-gray-700 dark:text-gray-300">
              {{ t('shortLinkDetail.refererSources') }}
            </h3>
            <div v-if="chartLoading" class="flex h-56 items-center justify-center">
              <Loader2 class="h-5 w-5 animate-spin text-gray-400 dark:text-gray-500" />
            </div>
            <div
              v-else-if="refererDistribution.length === 0"
              class="flex h-56 items-center justify-center text-sm text-gray-400 dark:text-gray-500"
            >
              {{ t('shortLinkDetail.noData') }}
            </div>
            <div v-else class="h-64">
              <VChart
                :option="refererChartOption"
                :autoresize="true"
                class="h-full w-full"
                role="img"
                aria-label="Referer source distribution pie chart"
              />
            </div>
          </div>

          <!-- IP distribution table -->
          <div class="mt-4 rounded-lg border dark:border-gray-700 bg-white dark:bg-gray-900 p-5">
            <h3 class="mb-4 text-sm font-semibold text-gray-700 dark:text-gray-300">
              {{ t('shortLinkDetail.topIps') }}
            </h3>
            <div
              v-if="ipDistribution.length === 0"
              class="py-8 text-center text-sm text-gray-400 dark:text-gray-500"
            >
              {{ t('shortLinkDetail.noData') }}
            </div>
            <table v-else class="w-full text-left text-sm">
              <thead>
                <tr class="border-b text-xs font-medium uppercase text-gray-400 dark:text-gray-500">
                  <th class="pb-2 pr-4">{{ t('shortLinkDetail.ipAddress') }}</th>
                  <th class="pb-2 pr-4 text-right">{{ t('shortLinkDetail.visits') }}</th>
                  <th class="pb-2 text-right">{{ t('shortLinkDetail.percentage') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="item in ipDistribution"
                  :key="item.ip"
                  class="border-b border-gray-50 dark:border-gray-800 last:border-0"
                >
                  <td class="py-2 pr-4 font-mono text-gray-700 dark:text-gray-300">
                    {{ item.ip }}
                  </td>
                  <td class="py-2 pr-4 text-right font-medium text-gray-600 dark:text-gray-400">
                    {{ item.count }}
                  </td>
                  <td class="py-2 text-right">
                    <span class="inline-flex items-center gap-2">
                      <span class="text-gray-500 dark:text-gray-400">{{ item.percent }}%</span>
                      <span
                        class="inline-block h-1.5 w-16 overflow-hidden rounded-full bg-gray-100 dark:bg-gray-800"
                      >
                        <span
                          class="inline-block h-full rounded-full bg-emerald-500"
                          :style="{ width: `${item.percent}%` }"
                        />
                      </span>
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </template>

        <!-- Access records tab -->
        <template v-if="activeTab === 'records'">
          <div class="mt-4 rounded-lg border dark:border-gray-700 bg-white dark:bg-gray-900">
            <div v-if="chartLoading" class="flex h-40 items-center justify-center">
              <Loader2 class="h-5 w-5 animate-spin text-gray-400 dark:text-gray-500" />
            </div>
            <div
              v-else-if="histories.length === 0"
              class="flex h-40 items-center justify-center text-sm text-gray-400 dark:text-gray-500"
            >
              {{ t('shortLinkDetail.noRecords') }}
            </div>
            <template v-else>
              <div class="overflow-x-auto">
                <table class="w-full text-left text-sm">
                  <thead>
                    <tr
                      class="border-b bg-gray-50 dark:bg-gray-800/50 text-xs font-medium uppercase text-gray-400 dark:text-gray-500"
                    >
                      <th class="px-4 py-3">{{ t('shortLinkDetail.time') }}</th>
                      <th class="px-4 py-3">{{ t('shortLinkDetail.ip') }}</th>
                      <th class="px-4 py-3">{{ t('shortLinkDetail.browser') }}</th>
                      <th class="px-4 py-3">{{ t('shortLinkDetail.os') }}</th>
                      <th class="px-4 py-3">{{ t('shortLinkDetail.referer') }}</th>
                      <th class="px-4 py-3">{{ t('shortLinkDetail.userAgent') }}</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr
                      v-for="h in histories"
                      :key="h.id"
                      class="border-b border-gray-50 dark:border-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 dark:bg-gray-800/50/50"
                    >
                      <td class="whitespace-nowrap px-4 py-2.5 text-gray-600 dark:text-gray-400">
                        {{ formatDate(h.time) }}
                      </td>
                      <td
                        class="whitespace-nowrap px-4 py-2.5 font-mono text-gray-700 dark:text-gray-300"
                      >
                        {{ h.ip }}
                      </td>
                      <td class="whitespace-nowrap px-4 py-2.5 text-gray-600 dark:text-gray-400">
                        {{ h.browser || parseBrowser(h.ua) || '—' }}
                      </td>
                      <td class="whitespace-nowrap px-4 py-2.5 text-gray-600 dark:text-gray-400">
                        {{ parseOS(h.ua) }}
                      </td>
                      <td
                        class="max-w-[200px] truncate px-4 py-2.5 text-gray-600 dark:text-gray-400"
                      >
                        {{
                          h.referer ? extractRefererSource(h.referer) : t('shortLinkDetail.direct')
                        }}
                      </td>
                      <td
                        class="max-w-[300px] truncate px-4 py-2.5 text-xs text-gray-400 dark:text-gray-500"
                        :title="h.ua"
                      >
                        {{ h.ua }}
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
              <div class="border-t px-4 py-3 text-xs text-gray-400 dark:text-gray-500">
                {{
                  t('shortLinkDetail.recordCount', {
                    count: histories.length,
                    suffix: histories.length !== 1 ? 's' : '',
                  })
                }}
              </div>
            </template>
          </div>
        </template>
      </div>
    </template>
  </div>
</template>
