<script setup lang="ts">
import { ref, watch, computed, nextTick, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import { listDomains } from '@/api/tenant'
import { useAuthStore } from '@/stores/auth'
import { generateQRCodeDataURL, downloadQRCode } from '@/utils/qrcode'
import { Copy, Check, Download, Loader2, X, QrCode } from 'lucide-vue-next'
import type { TenantDomain } from '@/types/api'

const props = defineProps<{
  shortLinkId: string
  open: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const { t } = useI18n()
const auth = useAuthStore()

const domains = ref<TenantDomain[]>([])
const selectedDomain = ref('')
const qrDataUrl = ref('')
const loadingDomains = ref(false)
const loadingQR = ref(false)
const loadError = ref(false)
const copied = ref(false)
const closeBtnRef = ref<HTMLButtonElement | null>(null)
let previousFocusEl: HTMLElement | null = null

const qrUrl = computed(() => {
  if (!selectedDomain.value) return ''
  const protocol = selectedDomain.value.startsWith('localhost') ? 'http://' : `${window.location.protocol}//`
  return `${protocol}${selectedDomain.value}/${props.shortLinkId}`
})

const hasMultipleDomains = computed(() => domains.value.length > 1)

async function fetchDomains() {
  if (!auth.currentTenantId) return
  loadingDomains.value = true
  loadError.value = false
  try {
    const data = await listDomains(auth.currentTenantId)
    domains.value = data
    const defaultDomain = data.find((d) => d.isDefault)
    if (defaultDomain) {
      selectedDomain.value = defaultDomain.domain
    } else if (data.length > 0) {
      selectedDomain.value = data[0].domain
    }
  } catch {
    loadError.value = true
  } finally {
    loadingDomains.value = false
  }
}

async function generateQR() {
  if (!qrUrl.value) {
    qrDataUrl.value = ''
    return
  }
  loadingQR.value = true
  try {
    qrDataUrl.value = await generateQRCodeDataURL(qrUrl.value)
  } catch {
    qrDataUrl.value = ''
  } finally {
    loadingQR.value = false
  }
}

function copyLink() {
  navigator.clipboard.writeText(qrUrl.value).then(
    () => {
      copied.value = true
      setTimeout(() => {
        copied.value = false
      }, 2000)
    },
    () => {},
  )
}

function handleDownload() {
  if (!qrDataUrl.value) return
  downloadQRCode(qrDataUrl.value, `qrcode-${props.shortLinkId}.png`)
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    emit('close')
  }
}

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen) {
      previousFocusEl = document.activeElement as HTMLElement
      document.addEventListener('keydown', handleKeydown)
      fetchDomains()
      nextTick(() => closeBtnRef.value?.focus())
    } else {
      document.removeEventListener('keydown', handleKeydown)
      if (previousFocusEl) {
        previousFocusEl.focus()
        previousFocusEl = null
      }
    }
  },
)

watch(selectedDomain, () => {
  generateQR()
})

watch(() => props.shortLinkId, () => {
  generateQR()
})

onBeforeUnmount(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<template>
  <Teleport to="body">
    <div
      v-if="open"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
      @click.self="emit('close')"
    >
      <div
        role="dialog"
        aria-modal="true"
        tabindex="-1"
        class="mx-4 w-full max-w-sm rounded-lg bg-white p-6 shadow-xl dark:bg-gray-900"
      >
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ t('qrCode.title') }}
          </h3>
          <button
            ref="closeBtnRef"
            class="rounded p-1 text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 hover:text-gray-600 dark:hover:text-gray-300"
            :aria-label="t('common.cancel')"
            @click="emit('close')"
          >
            <X class="h-5 w-5" />
          </button>
        </div>

        <!-- Loading state -->
        <div v-if="loadingDomains" class="mt-6 flex items-center justify-center py-8">
          <Loader2 class="h-6 w-6 animate-spin text-gray-400" />
        </div>

        <!-- Error state -->
        <div v-else-if="loadError" class="mt-4 text-center">
          <p class="text-sm text-red-500">{{ t('qrCode.loadError') }}</p>
          <button
            class="mt-2 text-sm text-blue-600 hover:underline"
            @click="fetchDomains"
          >
            {{ t('common.refresh') }}
          </button>
        </div>

        <!-- No domains -->
        <div v-else-if="domains.length === 0" class="mt-4 text-center">
          <QrCode class="mx-auto h-10 w-10 text-gray-300 dark:text-gray-600" />
          <p class="mt-3 text-sm text-gray-500 dark:text-gray-400">{{ t('qrCode.noDomains') }}</p>
        </div>

        <!-- QR Code content -->
        <template v-else>
          <div class="mt-4 flex justify-center">
            <div
              v-if="loadingQR"
              class="flex h-[256px] w-[256px] items-center justify-center"
            >
              <Loader2 class="h-6 w-6 animate-spin text-gray-400" />
            </div>
            <img
              v-else-if="qrDataUrl"
              :src="qrDataUrl"
              :alt="t('qrCode.imageAlt')"
              class="h-[256px] w-[256px] rounded"
            />
          </div>

          <!-- Short link URL with copy -->
          <div class="mt-4">
            <label class="text-xs font-medium uppercase text-gray-400 dark:text-gray-500">
              {{ t('shortLinkDetail.shortLink') }}
            </label>
            <div class="mt-1 flex items-center gap-2">
              <span class="flex-1 truncate rounded-md border border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-800 px-3 py-2 font-mono text-sm text-gray-700 dark:text-gray-300">
                {{ qrUrl }}
              </span>
              <button
                class="shrink-0 rounded-md p-2 text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 hover:text-blue-600"
                :title="t('shortLinkCreate.copyTitle')"
                @click="copyLink"
              >
                <Copy v-if="!copied" class="h-4 w-4" />
                <Check v-else class="h-4 w-4 text-green-500" />
              </button>
            </div>
          </div>

          <!-- Domain selector -->
          <div v-if="hasMultipleDomains" class="mt-3">
            <label class="text-xs font-medium uppercase text-gray-400 dark:text-gray-500">
              {{ t('common.domain') }}
            </label>
            <select
              v-model="selectedDomain"
              class="mt-1 w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
            >
              <option
                v-for="d in domains"
                :key="d.id"
                :value="d.domain"
              >
                {{ d.domain }}{{ d.isDefault ? ` (${t('common.default')})` : '' }}
              </option>
            </select>
          </div>

          <!-- Download button -->
          <button
            :disabled="!qrDataUrl"
            class="mt-4 flex w-full items-center justify-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:opacity-50"
            @click="handleDownload"
          >
            <Download class="h-4 w-4" />
            {{ t('qrCode.downloadPng') }}
          </button>
        </template>
      </div>
    </div>
  </Teleport>
</template>
