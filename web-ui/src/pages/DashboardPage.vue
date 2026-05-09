<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { listShortLinks, getShortLinkData } from '@/api/short-link'
import type { ShortLinkData, DailyStats, RequestHistory } from '@/types/api'
import {
  Link,
  LinkIcon,
  MousePointerClick,
  BarChart3,
  Plus,
  ExternalLink,
  Copy,
  Check,
  TrendingUp,
  ArrowRight,
  Clock,
  RefreshCw,
  TrendingDown,
  Minus,
} from 'lucide-vue-next'
import { defineAsyncComponent } from 'vue'

const VChart = defineAsyncComponent(async () => {
  const [
    { use },
    { CanvasRenderer },
    { LineChart },
    { GridComponent, TooltipComponent, LegendComponent },
  ] = await Promise.all([
    import('echarts/core'),
    import('echarts/renderers'),
    import('echarts/charts'),
    import('echarts/components'),
  ])
  use([CanvasRenderer, LineChart, GridComponent, TooltipComponent, LegendComponent])
  return (await import('vue-echarts')).default
})

defineOptions({ name: 'DashboardPage' })

const router = useRouter()
const auth = useAuthStore()

const loading = ref(true)
const refreshing = ref(false)
const totalLinks = ref(0)
const activeLinks = ref(0)
const todayVisits = ref(0)
const periodVisits = ref(0)
const recentLinks = ref<ShortLinkData[]>([])
const topLinks = ref<{ link: ShortLinkData; pv: number; uv: number }[]>([])
const trendDays = ref<7 | 30>(7)
const trendData = ref<DailyStats[]>([])
const chartLoading = ref(false)
const copiedId = ref<string | null>(null)
const prevTodayVisits = ref(0)
const prevPeriodVisits = ref(0)

const greeting = computed(() => {
  const h = new Date().getHours()
  if (h < 6) return 'Good night'
  if (h < 12) return 'Good morning'
  if (h < 18) return 'Good afternoon'
  return 'Good evening'
})

const todayStr = computed(() => {
  const d = new Date()
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
})

const yesterdayStr = computed(() => {
  const d = new Date()
  d.setDate(d.getDate() - 1)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
})

const periodLabel = computed(() => {
  return `Last ${trendDays.value} Days Visits`
})

const todayTrend = computed(() => {
  return computeTrend(todayVisits.value, prevTodayVisits.value)
})

const periodTrend = computed(() => {
  return computeTrend(periodVisits.value, prevPeriodVisits.value)
})

function computeTrend(current: number, previous: number) {
  if (previous === 0) return current > 0 ? { pct: 100, dir: 'up' as const } : null
  const pct = Math.round(((current - previous) / previous) * 100)
  if (pct > 0) return { pct, dir: 'up' as const }
  if (pct < 0) return { pct: Math.abs(pct), dir: 'down' as const }
  return { pct: 0, dir: 'flat' as const }
}

function formatDate(daysAgo: number) {
  const d = new Date()
  d.setDate(d.getDate() - daysAgo)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

function formatShortDate(dateStr: string) {
  const d = new Date(dateStr)
  return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
}

function extractDate(isoOrDateStr: string): string {
  return isoOrDateStr.slice(0, 10)
}

function aggregateFromHistories(histories: RequestHistory[]): {
  daily: Record<string, { pv: number; uv: number }>
  totalPv: number
  todayPv: number
} {
  const daily: Record<string, { pv: number; uv: number }> = {}
  const uvByDate: Record<string, Set<string>> = {}
  let totalPv = 0
  let todayPv = 0

  for (const h of histories) {
    const date = extractDate(h.time)
    if (!daily[date]) {
      daily[date] = { pv: 0, uv: 0 }
      uvByDate[date] = new Set()
    }
    daily[date].pv++
    uvByDate[date].add(h.ip)
    totalPv++
    if (date === todayStr.value) {
      todayPv++
    }
  }

  for (const [date, ips] of Object.entries(uvByDate)) {
    daily[date].uv = ips.size
  }

  return { daily, totalPv, todayPv }
}

interface LinkStatsResult {
  allDaily: Record<string, { pv: number; uv: number }>
  linkDaily: { link: ShortLinkData; daily: Record<string, { pv: number; uv: number }> }[]
  totalPv: number
  todayPv: number
}

async function fetchLinkStats(
  links: ShortLinkData[],
  startDate: string,
  endDate: string,
): Promise<LinkStatsResult> {
  const allDaily: Record<string, { pv: number; uv: number }> = {}
  const linkDaily: { link: ShortLinkData; daily: Record<string, { pv: number; uv: number }> }[] = []
  let totalPv = 0
  let todayPv = 0

  const results = await Promise.allSettled(
    links.map((link) => getShortLinkData(link.id, startDate, endDate)),
  )

  for (let i = 0; i < results.length; i++) {
    const result = results[i]
    if (result.status === 'fulfilled' && result.value) {
      const { daily, totalPv: linkTotalPv, todayPv: linkTodayPv } = aggregateFromHistories(
        result.value.histories,
      )

      for (const [date, stats] of Object.entries(daily)) {
        if (!allDaily[date]) {
          allDaily[date] = { pv: 0, uv: 0 }
        }
        allDaily[date].pv += stats.pv
        allDaily[date].uv += stats.uv
      }

      totalPv += linkTotalPv
      todayPv += linkTodayPv

      if (Object.keys(daily).length > 0) {
        linkDaily.push({ link: links[i], daily })
      }
    }
  }

  return { allDaily, linkDaily, totalPv, todayPv }
}

function computeTopLinks(
  linkDaily: { link: ShortLinkData; daily: Record<string, { pv: number; uv: number }> }[],
  startDate: string,
): { link: ShortLinkData; pv: number; uv: number }[] {
  const result: { link: ShortLinkData; pv: number; uv: number }[] = []
  for (const entry of linkDaily) {
    let pv = 0
    let uv = 0
    for (const [date, stats] of Object.entries(entry.daily)) {
      if (date >= startDate) {
        pv += stats.pv
        uv += stats.uv
      }
    }
    if (pv > 0) {
      result.push({ link: entry.link, pv, uv })
    }
  }
  return result.sort((a, b) => b.pv - a.pv).slice(0, 10)
}

function dailyToTrend(allDaily: Record<string, { pv: number; uv: number }>): DailyStats[] {
  return Object.entries(allDaily)
    .sort(([a], [b]) => a.localeCompare(b))
    .map(([date, { pv, uv }]) => ({ date, pv, uv }))
}

const chartOption = computed(() => {
  const dates = trendData.value.map((d) => formatShortDate(d.date))
  const pvData = trendData.value.map((d) => d.pv)
  const uvData = trendData.value.map((d) => d.uv)

  return {
    tooltip: {
      trigger: 'axis' as const,
      backgroundColor: '#fff',
      borderColor: '#e5e7eb',
      borderWidth: 1,
      textStyle: { color: '#374151', fontSize: 13 },
      axisPointer: { type: 'shadow' as const },
    },
    legend: {
      data: ['PV', 'UV'],
      bottom: 0,
      textStyle: { color: '#6b7280', fontSize: 12 },
      itemWidth: 16,
      itemHeight: 2,
      itemGap: 24,
    },
    grid: {
      top: 16,
      right: 12,
      bottom: 40,
      left: 36,
    },
    xAxis: {
      type: 'category' as const,
      data: dates,
      boundaryGap: false,
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
        name: 'PV',
        type: 'line' as const,
        smooth: true,
        symbol: 'circle',
        symbolSize: 6,
        lineStyle: { width: 2.5, color: '#2563eb' },
        itemStyle: { color: '#2563eb', borderWidth: 2, borderColor: '#fff' },
        areaStyle: {
          color: {
            type: 'linear' as const,
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(37,99,235,0.15)' },
              { offset: 1, color: 'rgba(37,99,235,0)' },
            ],
          },
        },
        data: pvData,
      },
      {
        name: 'UV',
        type: 'line' as const,
        smooth: true,
        symbol: 'circle',
        symbolSize: 6,
        lineStyle: { width: 2.5, color: '#8b5cf6' },
        itemStyle: { color: '#8b5cf6', borderWidth: 2, borderColor: '#fff' },
        areaStyle: {
          color: {
            type: 'linear' as const,
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(139,92,246,0.15)' },
              { offset: 1, color: 'rgba(139,92,246,0)' },
            ],
          },
        },
        data: uvData,
      },
    ],
  }
})

async function fetchAllPages() {
  let allLinks: ShortLinkData[] = []
  let page = 1
  let total = 0

  while (true) {
    const data = await listShortLinks(page, 100)
    total = data.total
    allLinks = allLinks.concat(data.shortLinks)
    if (allLinks.length >= total || data.shortLinks.length === 0) break
    page++
  }

  return { links: allLinks, total }
}

async function fetchDashboardData(isRefresh = false) {
  if (isRefresh) {
    refreshing.value = true
  } else {
    loading.value = true
  }
  try {
    const { links, total } = await fetchAllPages()
    totalLinks.value = total
    activeLinks.value = links.filter((l) => l.isEnable).length

    recentLinks.value = [...links]
      .sort((a, b) => new Date(b.createTime).getTime() - new Date(a.createTime).getTime())
      .slice(0, 8)

    // Fetch 2x the period to compute trends (current period + previous period)
    const extendedStartDate = formatDate(trendDays.value * 2 - 1)
    const endDate = formatDate(0)
    const currentStartDate = formatDate(trendDays.value - 1)

    const stats = await fetchLinkStats(links, extendedStartDate, endDate)

    // Split daily data into current and previous periods
    let currentPeriodPv = 0
    let previousPeriodPv = 0
    const currentDaily: Record<string, { pv: number; uv: number }> = {}

    for (const [date, val] of Object.entries(stats.allDaily)) {
      if (date >= currentStartDate) {
        currentDaily[date] = val
        currentPeriodPv += val.pv
      } else {
        previousPeriodPv += val.pv
      }
    }

    todayVisits.value = stats.todayPv
    periodVisits.value = currentPeriodPv
    prevPeriodVisits.value = previousPeriodPv

    // Yesterday visits for today's trend
    const yesterdayData = stats.allDaily[yesterdayStr.value]
    prevTodayVisits.value = yesterdayData ? yesterdayData.pv : 0

    trendData.value = dailyToTrend(currentDaily)
    topLinks.value = computeTopLinks(stats.linkDaily, currentStartDate)
  } catch {
    // error handled by interceptor
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

async function fetchTrendData() {
  chartLoading.value = true
  try {
    const startDate = formatDate(trendDays.value - 1)
    const endDate = formatDate(0)
    const prevStartDate = formatDate(trendDays.value * 2 - 1)

    const { links } = await fetchAllPages()

    // Fetch extended range for trend comparison
    const stats = await fetchLinkStats(links, prevStartDate, endDate)

    let currentPeriodPv = 0
    let previousPeriodPv = 0
    const currentDaily: Record<string, { pv: number; uv: number }> = {}

    for (const [date, val] of Object.entries(stats.allDaily)) {
      if (date >= startDate) {
        currentDaily[date] = val
        currentPeriodPv += val.pv
      } else {
        previousPeriodPv += val.pv
      }
    }

    periodVisits.value = currentPeriodPv
    prevPeriodVisits.value = previousPeriodPv
    trendData.value = dailyToTrend(currentDaily)
    topLinks.value = computeTopLinks(stats.linkDaily, startDate)
  } catch {
    // error handled by interceptor
  } finally {
    chartLoading.value = false
  }
}

function switchTrend(days: 7 | 30) {
  if (trendDays.value === days) return
  trendDays.value = days
  fetchTrendData()
}

function handleRefresh() {
  fetchDashboardData(true)
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

onMounted(() => fetchDashboardData())
</script>

<template>
  <div class="space-y-6">
    <!-- Welcome area -->
    <div class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <h1 class="text-xl font-bold text-gray-900 dark:text-white sm:text-2xl">{{ greeting }}, {{ auth.username }}</h1>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          Here's an overview of your short links performance.
        </p>
      </div>
      <div class="flex flex-wrap items-center gap-2">
        <button
          :disabled="loading || refreshing"
          class="inline-flex items-center gap-2 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 px-3 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 shadow-sm transition-colors hover:bg-gray-50 dark:hover:bg-gray-700 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:outline-none disabled:opacity-50 sm:px-4 sm:py-2.5"
          aria-label="Refresh dashboard"
          @click="handleRefresh"
        >
          <RefreshCw class="h-4 w-4" :class="{ 'animate-spin': refreshing }" />
          <span class="hidden sm:inline">Refresh</span>
        </button>
        <router-link
          :to="{ name: 'short-link-create' }"
          class="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-3 py-2 text-sm font-medium text-white shadow-sm transition-colors hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:outline-none sm:px-4 sm:py-2.5"
        >
          <Plus class="h-4 w-4" />
          Create
        </router-link>
      </div>
    </div>

    <!-- Stats cards -->
    <div class="grid grid-cols-2 gap-3 sm:gap-4 lg:grid-cols-4">
      <!-- Total Links -->
      <div class="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 p-3 shadow-sm sm:p-5">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400 sm:text-sm">Total Links</p>
            <p class="mt-1 text-xl font-bold text-gray-900 dark:text-white sm:text-2xl">
              <span
                v-if="loading"
                class="inline-block h-7 w-16 animate-pulse rounded bg-gray-200 dark:bg-gray-700"
              />
              <template v-else>{{ totalLinks.toLocaleString() }}</template>
            </p>
          </div>
          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-50 sm:h-10 sm:w-10">
            <Link class="h-4 w-4 text-blue-600 sm:h-5 sm:w-5" />
          </div>
        </div>
      </div>

      <!-- Active Links -->
      <div class="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 p-3 shadow-sm sm:p-5">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400 sm:text-sm">Active Links</p>
            <p class="mt-1 text-xl font-bold text-gray-900 dark:text-white sm:text-2xl">
              <span
                v-if="loading"
                class="inline-block h-7 w-16 animate-pulse rounded bg-gray-200 dark:bg-gray-700"
              />
              <template v-else>{{ activeLinks.toLocaleString() }}</template>
            </p>
          </div>
          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-green-50 sm:h-10 sm:w-10">
            <LinkIcon class="h-4 w-4 text-green-600 sm:h-5 sm:w-5" />
          </div>
        </div>
      </div>

      <!-- Today's Visits -->
      <div class="rounded-xl border border-gray-200 bg-white p-3 shadow-sm sm:p-5">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-xs font-medium text-gray-500 sm:text-sm">Today's Visits</p>
            <div class="mt-1 flex items-baseline gap-2">
              <p class="text-xl font-bold text-gray-900 sm:text-2xl">
                <span
                  v-if="loading"
                  class="inline-block h-7 w-16 animate-pulse rounded bg-gray-200"
                />
                <template v-else>{{ todayVisits.toLocaleString() }}</template>
              </p>
              <span
                v-if="!loading && todayTrend"
                class="inline-flex items-center gap-0.5 text-xs font-medium"
                :class="todayTrend.dir === 'up' ? 'text-green-600' : todayTrend.dir === 'down' ? 'text-red-500' : 'text-gray-400'"
              >
                <TrendingUp v-if="todayTrend.dir === 'up'" class="h-3 w-3" />
                <TrendingDown v-else-if="todayTrend.dir === 'down'" class="h-3 w-3" />
                <Minus v-else class="h-3 w-3" />
                {{ todayTrend.pct }}%
              </span>
            </div>
          </div>
          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-purple-50 sm:h-10 sm:w-10">
            <MousePointerClick class="h-4 w-4 text-purple-600 sm:h-5 sm:w-5" />
          </div>
        </div>
        <p v-if="!loading && todayTrend" class="mt-1 text-xs text-gray-400">vs yesterday</p>
      </div>

      <!-- Period Visits -->
      <div class="rounded-xl border border-gray-200 bg-white p-3 shadow-sm sm:p-5">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-xs font-medium text-gray-500 sm:text-sm">{{ periodLabel }}</p>
            <div class="mt-1 flex items-baseline gap-2">
              <p class="text-xl font-bold text-gray-900 sm:text-2xl">
                <span
                  v-if="loading"
                  class="inline-block h-7 w-16 animate-pulse rounded bg-gray-200"
                />
                <template v-else>{{ periodVisits.toLocaleString() }}</template>
              </p>
              <span
                v-if="!loading && periodTrend"
                class="inline-flex items-center gap-0.5 text-xs font-medium"
                :class="periodTrend.dir === 'up' ? 'text-green-600' : periodTrend.dir === 'down' ? 'text-red-500' : 'text-gray-400'"
              >
                <TrendingUp v-if="periodTrend.dir === 'up'" class="h-3 w-3" />
                <TrendingDown v-else-if="periodTrend.dir === 'down'" class="h-3 w-3" />
                <Minus v-else class="h-3 w-3" />
                {{ periodTrend.pct }}%
              </span>
            </div>
          </div>
          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-orange-50 sm:h-10 sm:w-10">
            <BarChart3 class="h-4 w-4 text-orange-600 sm:h-5 sm:w-5" />
          </div>
        </div>
        <p v-if="!loading && periodTrend" class="mt-1 text-xs text-gray-400">vs previous {{ trendDays }} days</p>
      </div>
    </div>

    <!-- Chart + Recent Links -->
    <div class="grid grid-cols-1 gap-4 lg:grid-cols-5 lg:gap-6">
      <!-- Visit Trend Chart -->
      <div class="rounded-xl border border-gray-200 bg-white p-3 shadow-sm sm:p-5 lg:col-span-3">
        <div class="mb-4 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <TrendingUp class="h-4 w-4 text-gray-500" />
            <h2 class="text-base font-semibold text-gray-900">Visit Trends</h2>
          </div>
          <div class="flex rounded-lg border border-gray-200 p-0.5">
            <button
              :class="[
                'rounded-md px-3 py-1 text-xs font-medium transition-colors',
                trendDays === 7 ? 'bg-blue-600 text-white' : 'text-gray-500 hover:text-gray-700',
              ]"
              @click="switchTrend(7)"
            >
              7 Days
            </button>
            <button
              :class="[
                'rounded-md px-3 py-1 text-xs font-medium transition-colors',
                trendDays === 30 ? 'bg-blue-600 text-white' : 'text-gray-500 hover:text-gray-700',
              ]"
              @click="switchTrend(30)"
            >
              30 Days
            </button>
          </div>
        </div>
        <div v-if="chartLoading" class="h-64">
          <!-- Chart skeleton -->
          <div class="flex h-full items-end gap-3 px-2 pt-4">
            <div class="flex-1 animate-pulse rounded-t bg-gray-100" style="height: 55%" />
            <div class="flex-1 animate-pulse rounded-t bg-gray-100" style="height: 75%" />
            <div class="flex-1 animate-pulse rounded-t bg-gray-100" style="height: 45%" />
            <div class="flex-1 animate-pulse rounded-t bg-gray-100" style="height: 85%" />
            <div class="flex-1 animate-pulse rounded-t bg-gray-100" style="height: 65%" />
            <div class="flex-1 animate-pulse rounded-t bg-gray-100" style="height: 40%" />
            <div class="flex-1 animate-pulse rounded-t bg-gray-100" style="height: 70%" />
          </div>
          <div class="mt-3 flex justify-between px-2">
            <div class="h-2.5 w-10 animate-pulse rounded bg-gray-100" />
            <div class="h-2.5 w-10 animate-pulse rounded bg-gray-100" />
            <div class="h-2.5 w-10 animate-pulse rounded bg-gray-100" />
            <div class="h-2.5 w-10 animate-pulse rounded bg-gray-100" />
          </div>
        </div>
        <div v-else class="h-48 sm:h-64">
          <VChart
            :option="chartOption"
            :autoresize="true"
            class="h-full w-full"
            aria-label="Visit trend chart showing PV and UV over time"
            role="img"
          />
        </div>
      </div>

      <!-- Recent Links -->
      <div class="rounded-xl border border-gray-200 bg-white p-3 shadow-sm sm:p-5 lg:col-span-2">
        <div class="mb-4 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <Clock class="h-4 w-4 text-gray-500" />
            <h2 class="text-base font-semibold text-gray-900">Recent Links</h2>
          </div>
          <router-link
            :to="{ name: 'short-links' }"
            class="inline-flex items-center gap-1 text-xs font-medium text-blue-600 hover:text-blue-700"
          >
            View all
            <ArrowRight class="h-3 w-3" />
          </router-link>
        </div>
        <div v-if="loading" class="space-y-3">
          <div v-for="i in 5" :key="i" class="h-12 animate-pulse rounded-lg bg-gray-100" />
        </div>
        <div v-else-if="recentLinks.length === 0" class="py-8 text-center">
          <Link class="mx-auto h-8 w-8 text-gray-300" />
          <p class="mt-2 text-sm text-gray-400">No short links yet</p>
          <router-link
            :to="{ name: 'short-link-create' }"
            class="mt-2 inline-flex items-center gap-1 text-sm font-medium text-blue-600 hover:text-blue-700"
          >
            <Plus class="h-3.5 w-3.5" />
            Create your first link
          </router-link>
        </div>
        <ul v-else class="divide-y divide-gray-100">
          <li
            v-for="link in recentLinks"
            :key="link.id"
            class="group flex items-center justify-between py-2.5 first:pt-0 last:pb-0"
          >
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <span class="font-mono text-sm font-medium text-gray-900">{{ link.id }}</span>
                <span
                  class="inline-flex items-center rounded-full px-1.5 py-0.5 text-[10px] font-medium"
                  :class="
                    link.isEnable ? 'bg-green-50 text-green-700' : 'bg-gray-100 text-gray-500'
                  "
                >
                  {{ link.isEnable ? 'Active' : 'Inactive' }}
                </span>
              </div>
              <p class="mt-0.5 truncate text-xs text-gray-500">{{ link.url }}</p>
            </div>
            <div
              class="ml-3 flex shrink-0 items-center gap-1 opacity-0 transition-opacity group-hover:opacity-100"
            >
              <button
                class="rounded p-1 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600"
                title="Copy short link"
                @click="copyLink(link.id)"
              >
                <Copy v-if="copiedId !== link.id" class="h-3.5 w-3.5" />
                <Check v-else class="h-3.5 w-3.5 text-green-500" />
              </button>
              <button
                class="rounded p-1 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600"
                title="View details"
                @click="router.push({ name: 'short-link-detail', params: { id: link.id } })"
              >
                <ExternalLink class="h-3.5 w-3.5" />
              </button>
            </div>
          </li>
        </ul>
      </div>
    </div>

    <!-- Top 10 Links -->
    <div class="rounded-xl border border-gray-200 bg-white shadow-sm">
      <div class="flex items-center justify-between border-b border-gray-100 px-5 py-4">
        <div class="flex items-center gap-2">
          <BarChart3 class="h-4 w-4 text-gray-500" />
          <h2 class="text-base font-semibold text-gray-900">Top Links by Visits</h2>
        </div>
        <span class="text-xs text-gray-400">Based on {{ trendDays }}-day data</span>
      </div>
      <div v-if="loading" class="p-6">
        <div class="space-y-3">
          <div v-for="i in 5" :key="i" class="h-10 animate-pulse rounded-lg bg-gray-100" />
        </div>
      </div>
      <div v-else-if="topLinks.length === 0" class="px-5 py-12 text-center">
        <BarChart3 class="mx-auto h-8 w-8 text-gray-300" />
        <p class="mt-2 text-sm text-gray-400">No visit data available yet</p>
      </div>
      <div v-else>
        <!-- Mobile cards -->
        <div class="divide-y divide-gray-50 md:hidden">
          <div
            v-for="(item, index) in topLinks"
            :key="item.link.id"
            class="flex items-center gap-3 px-4 py-3"
          >
            <span
              :class="[
                'inline-flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-xs font-semibold',
                index < 3 ? 'bg-blue-100 text-blue-700' : 'bg-gray-100 text-gray-500',
              ]"
            >
              {{ index + 1 }}
            </span>
            <div class="min-w-0 flex-1">
              <router-link
                :to="{ name: 'short-link-detail', params: { id: item.link.id } }"
                class="font-mono text-sm font-medium text-blue-600 hover:text-blue-700"
              >
                {{ item.link.id }}
              </router-link>
              <p class="truncate text-xs text-gray-400">{{ item.link.url }}</p>
            </div>
            <div class="shrink-0 text-right text-sm">
              <p class="font-semibold text-gray-900">{{ item.pv.toLocaleString() }}</p>
              <p class="text-xs text-gray-400">{{ item.uv.toLocaleString() }} uv</p>
            </div>
          </div>
        </div>
        <!-- Desktop table -->
        <div class="hidden overflow-x-auto md:block">
          <table class="w-full text-sm">
            <thead>
              <tr
                class="border-b border-gray-100 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
              >
                <th class="w-12 px-5 py-3">#</th>
                <th class="px-5 py-3">Short Link</th>
                <th class="px-5 py-3">Destination URL</th>
                <th class="px-5 py-3 text-right">PV</th>
                <th class="px-5 py-3 text-right">UV</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-50">
              <tr
                v-for="(item, index) in topLinks"
                :key="item.link.id"
                class="transition-colors hover:bg-gray-50"
              >
                <td class="px-5 py-3">
                  <span
                    :class="[
                      'inline-flex h-6 w-6 items-center justify-center rounded-full text-xs font-semibold',
                      index < 3 ? 'bg-blue-100 text-blue-700' : 'bg-gray-100 text-gray-500',
                    ]"
                  >
                    {{ index + 1 }}
                  </span>
                </td>
                <td class="px-5 py-3">
                  <router-link
                    :to="{ name: 'short-link-detail', params: { id: item.link.id } }"
                    class="inline-flex items-center gap-1 font-mono text-sm font-medium text-blue-600 hover:text-blue-700"
                  >
                    {{ item.link.id }}
                  </router-link>
                </td>
                <td class="max-w-[280px] truncate px-5 py-3 text-gray-500">
                  {{ item.link.url }}
                </td>
                <td class="whitespace-nowrap px-5 py-3 text-right font-semibold text-gray-900">
                  {{ item.pv.toLocaleString() }}
                </td>
                <td class="whitespace-nowrap px-5 py-3 text-right text-gray-600">
                  {{ item.uv.toLocaleString() }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>
