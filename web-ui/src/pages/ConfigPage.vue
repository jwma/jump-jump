<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getConfig, updateIdConfig, updateShortLinkNotFoundConfig } from '@/api/config'
import type { IdConfig, ShortLinkNotFoundConfig } from '@/types/api'
import { Settings, Hash, AlertTriangle, Loader2, Save, Check, Eye } from 'lucide-vue-next'

defineOptions({ name: 'ConfigPage' })

const loading = ref(true)
const pageError = ref('')

// ID config form
const idLength = ref(6)
const idMinimumLength = ref(2)
const idMaximumLength = ref(10)
const idSaving = ref(false)
const idSaveSuccess = ref(false)
const idError = ref('')

// 404 config form
const notFoundMode = ref<'content' | 'redirect'>('content')
const notFoundValue = ref('')
const notFoundSaving = ref(false)
const notFoundSaveSuccess = ref(false)
const notFoundError = ref('')

// Confirm dialog
const confirmOpen = ref(false)
const confirmTitle = ref('')
const confirmMessage = ref('')
const confirmAction = ref<(() => Promise<void>) | null>(null)

const idValidationError = computed(() => {
  if (idMinimumLength.value >= idLength.value) {
    return 'Minimum length must be less than default length'
  }
  if (idLength.value >= idMaximumLength.value) {
    return 'Default length must be less than maximum length'
  }
  if (idMinimumLength.value < 2) {
    return 'Minimum length must be at least 2'
  }
  if (idMaximumLength.value > 10) {
    return 'Maximum length must be at most 10'
  }
  return ''
})

const notFoundValidationError = computed(() => {
  if (!notFoundValue.value.trim()) {
    return notFoundMode.value === 'content' ? 'Content text is required' : 'Redirect URL is required'
  }
  if (notFoundMode.value === 'redirect') {
    try {
      new URL(notFoundValue.value)
      return ''
    } catch {
      return 'Please enter a valid URL (e.g. https://example.com/404)'
    }
  }
  return ''
})

const idCanSave = computed(
  () => !idValidationError.value && !idSaving.value && !loading.value,
)

const notFoundCanSave = computed(
  () => !notFoundValidationError.value && !notFoundSaving.value && !loading.value,
)

const sampleId = computed(() => {
  const chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
  let result = ''
  for (let i = 0; i < idLength.value; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  return result
})

const previewText = computed(() => {
  if (notFoundMode.value === 'content') {
    return notFoundValue.value || '(No content configured)'
  }
  return `Redirect to: ${notFoundValue.value || '(No URL configured)'}`
})

async function fetchConfig() {
  loading.value = true
  pageError.value = ''
  try {
    const data = await getConfig()
    idLength.value = data.config.idConfig.idLength
    idMinimumLength.value = data.config.idConfig.idMinimumLength
    idMaximumLength.value = data.config.idConfig.idMaximumLength
    notFoundMode.value = data.config.shortLinkNotFoundConfig.mode
    notFoundValue.value = data.config.shortLinkNotFoundConfig.value
  } catch {
    pageError.value = 'Failed to load configuration. Please refresh the page.'
  } finally {
    loading.value = false
  }
}

function openConfirm(title: string, message: string, action: () => Promise<void>) {
  confirmTitle.value = title
  confirmMessage.value = message
  confirmAction.value = action
  confirmOpen.value = true
}

function closeConfirm() {
  confirmOpen.value = false
  confirmAction.value = null
}

async function executeConfirm() {
  if (confirmAction.value) {
    await confirmAction.value()
  }
  closeConfirm()
}

function handleSaveIdConfig() {
  openConfirm(
    'Save ID Length Configuration',
    `Set default length to ${idLength.value}, min ${idMinimumLength.value}, max ${idMaximumLength.value}?`,
    doSaveIdConfig,
  )
}

async function doSaveIdConfig() {
  idSaving.value = true
  idError.value = ''
  idSaveSuccess.value = false
  try {
    await updateIdConfig({
      idLength: idLength.value,
      idMinimumLength: idMinimumLength.value,
      idMaximumLength: idMaximumLength.value,
    } as IdConfig)
    idSaveSuccess.value = true
    setTimeout(() => {
      idSaveSuccess.value = false
    }, 3000)
  } catch {
    idError.value = 'Failed to save ID length configuration.'
  } finally {
    idSaving.value = false
  }
}

function handleSaveNotFoundConfig() {
  const modeLabel = notFoundMode.value === 'content' ? 'display content' : 'redirect'
  openConfirm(
    'Save 404 Handling Configuration',
    `Set 404 handling to ${modeLabel}?`,
    doSaveNotFoundConfig,
  )
}

async function doSaveNotFoundConfig() {
  notFoundSaving.value = true
  notFoundError.value = ''
  notFoundSaveSuccess.value = false
  try {
    await updateShortLinkNotFoundConfig({
      mode: notFoundMode.value,
      value: notFoundValue.value,
    } as ShortLinkNotFoundConfig)
    notFoundSaveSuccess.value = true
    setTimeout(() => {
      notFoundSaveSuccess.value = false
    }, 3000)
  } catch {
    notFoundError.value = 'Failed to save 404 handling configuration.'
  } finally {
    notFoundSaving.value = false
  }
}

function refreshSampleId() {
  // Trigger re-computation by toggling a dependency
  const current = idLength.value
  idLength.value = current + 1
  idLength.value = current
}

onMounted(fetchConfig)
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div>
      <h1 class="flex items-center gap-2 text-xl font-semibold text-gray-900">
        <Settings class="h-5 w-5 text-gray-500" />
        System Configuration
      </h1>
      <p class="mt-1 text-sm text-gray-500">Manage system-wide settings for short link ID generation and 404 handling.</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-20">
      <Loader2 class="h-6 w-6 animate-spin text-gray-400" />
    </div>

    <!-- Error -->
    <div v-else-if="pageError" class="rounded-lg border border-red-200 bg-red-50 p-4">
      <p class="text-sm text-red-600">{{ pageError }}</p>
    </div>

    <template v-else>
      <!-- ID Length Configuration -->
      <div class="rounded-xl border border-gray-200 bg-white shadow-sm">
        <div class="flex items-center gap-2 border-b border-gray-100 px-5 py-4">
          <Hash class="h-4 w-4 text-blue-600" />
          <h2 class="text-base font-semibold text-gray-900">ID Length Configuration</h2>
        </div>
        <div class="p-5">
          <div class="max-w-md space-y-4">
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700">
                Default ID Length
              </label>
              <input
                v-model.number="idLength"
                type="number"
                :min="idMinimumLength + 1"
                :max="idMaximumLength - 1"
                class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
              />
              <p class="mt-1 text-xs text-gray-400">The length used when auto-generating short link IDs.</p>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="mb-1 block text-sm font-medium text-gray-700">
                  Minimum Length
                </label>
                <input
                  v-model.number="idMinimumLength"
                  type="number"
                  :min="2"
                  :max="idLength - 1"
                  class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
                />
              </div>
              <div>
                <label class="mb-1 block text-sm font-medium text-gray-700">
                  Maximum Length
                </label>
                <input
                  v-model.number="idMaximumLength"
                  type="number"
                  :min="idLength + 1"
                  :max="10"
                  class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
                />
              </div>
            </div>

            <p v-if="idValidationError" class="text-xs text-red-600">{{ idValidationError }}</p>
            <p v-if="idError" class="text-xs text-red-600">{{ idError }}</p>
            <p v-if="idSaveSuccess" class="flex items-center gap-1 text-xs text-green-600">
              <Check class="h-3.5 w-3.5" />
              Configuration saved successfully.
            </p>
          </div>

          <div class="mt-5">
            <button
              :disabled="!idCanSave"
              class="inline-flex items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:opacity-50"
              @click="handleSaveIdConfig"
            >
              <Loader2 v-if="idSaving" class="h-4 w-4 animate-spin" />
              <Save v-else class="h-4 w-4" />
              {{ idSaving ? 'Saving...' : 'Save ID Configuration' }}
            </button>
          </div>
        </div>
      </div>

      <!-- 404 Handling Configuration -->
      <div class="rounded-xl border border-gray-200 bg-white shadow-sm">
        <div class="flex items-center gap-2 border-b border-gray-100 px-5 py-4">
          <AlertTriangle class="h-4 w-4 text-orange-500" />
          <h2 class="text-base font-semibold text-gray-900">404 Handling Configuration</h2>
        </div>
        <div class="p-5">
          <div class="max-w-md space-y-4">
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700">Handling Mode</label>
              <div class="flex rounded-lg border border-gray-200 p-0.5">
                <button
                  :class="[
                    'flex-1 rounded-md px-4 py-2 text-sm font-medium transition-colors',
                    notFoundMode === 'content'
                      ? 'bg-blue-600 text-white'
                      : 'text-gray-500 hover:text-gray-700',
                  ]"
                  @click="notFoundMode = 'content'"
                >
                  Display Content
                </button>
                <button
                  :class="[
                    'flex-1 rounded-md px-4 py-2 text-sm font-medium transition-colors',
                    notFoundMode === 'redirect'
                      ? 'bg-blue-600 text-white'
                      : 'text-gray-500 hover:text-gray-700',
                  ]"
                  @click="notFoundMode = 'redirect'"
                >
                  Redirect
                </button>
              </div>
            </div>

            <div v-if="notFoundMode === 'content'">
              <label class="mb-1 block text-sm font-medium text-gray-700">
                Content Text
              </label>
              <textarea
                v-model="notFoundValue"
                rows="4"
                placeholder="The page content to display when a short link is not found..."
                class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
              />
              <p class="mt-1 text-xs text-gray-400">This content will be shown to visitors when they access a non-existent short link.</p>
            </div>

            <div v-else>
              <label class="mb-1 block text-sm font-medium text-gray-700">
                Redirect URL
              </label>
              <input
                v-model="notFoundValue"
                type="url"
                placeholder="https://example.com/404"
                class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
              />
              <p class="mt-1 text-xs text-gray-400">Visitors will be redirected to this URL when they access a non-existent short link.</p>
            </div>

            <p v-if="notFoundValidationError" class="text-xs text-red-600">{{ notFoundValidationError }}</p>
            <p v-if="notFoundError" class="text-xs text-red-600">{{ notFoundError }}</p>
            <p v-if="notFoundSaveSuccess" class="flex items-center gap-1 text-xs text-green-600">
              <Check class="h-3.5 w-3.5" />
              Configuration saved successfully.
            </p>
          </div>

          <div class="mt-5">
            <button
              :disabled="!notFoundCanSave"
              class="inline-flex items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:opacity-50"
              @click="handleSaveNotFoundConfig"
            >
              <Loader2 v-if="notFoundSaving" class="h-4 w-4 animate-spin" />
              <Save v-else class="h-4 w-4" />
              {{ notFoundSaving ? 'Saving...' : 'Save 404 Configuration' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Configuration Preview -->
      <div class="rounded-xl border border-gray-200 bg-white shadow-sm">
        <div class="flex items-center gap-2 border-b border-gray-100 px-5 py-4">
          <Eye class="h-4 w-4 text-purple-600" />
          <h2 class="text-base font-semibold text-gray-900">Configuration Preview</h2>
        </div>
        <div class="p-5">
          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <!-- ID Length Preview -->
            <div class="rounded-lg border border-gray-100 bg-gray-50 p-4">
              <p class="mb-2 text-xs font-medium uppercase tracking-wider text-gray-500">ID Length</p>
              <div class="flex items-center gap-3">
                <span class="font-mono text-lg font-bold text-blue-600">{{ sampleId }}</span>
                <button
                  class="rounded p-1 text-gray-400 transition-colors hover:bg-gray-200 hover:text-gray-600"
                  title="Regenerate sample"
                  @click="refreshSampleId"
                >
                  <Loader2 class="h-3.5 w-3.5" />
                </button>
              </div>
              <p class="mt-2 text-xs text-gray-500">
                Length: {{ idLength }} characters (min {{ idMinimumLength }}, max {{ idMaximumLength }})
              </p>
            </div>

            <!-- 404 Handling Preview -->
            <div class="rounded-lg border border-gray-100 bg-gray-50 p-4">
              <p class="mb-2 text-xs font-medium uppercase tracking-wider text-gray-500">404 Handling</p>
              <div class="flex items-center gap-2">
                <span
                  class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium"
                  :class="
                    notFoundMode === 'content'
                      ? 'bg-blue-50 text-blue-700'
                      : 'bg-orange-50 text-orange-700'
                  "
                >
                  {{ notFoundMode === 'content' ? 'Display Content' : 'Redirect' }}
                </span>
              </div>
              <p class="mt-2 line-clamp-3 text-xs text-gray-600">{{ previewText }}</p>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Confirm Dialog -->
    <Teleport to="body">
      <div
        v-if="confirmOpen"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
        @click.self="closeConfirm"
      >
        <div class="mx-4 w-full max-w-sm rounded-lg bg-white p-6 shadow-xl">
          <h3 class="text-lg font-semibold text-gray-900">{{ confirmTitle }}</h3>
          <p class="mt-2 text-sm text-gray-600">{{ confirmMessage }}</p>
          <div class="mt-5 flex justify-end gap-3">
            <button
              class="rounded-md px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-100"
              @click="closeConfirm"
            >
              Cancel
            </button>
            <button
              class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700"
              @click="executeConfirm"
            >
              Confirm
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
