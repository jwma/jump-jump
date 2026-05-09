<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getShortLink, updateShortLink } from '@/api/short-link'
import { ApiError } from '@/api/http'
import { useToast } from '@/composables/useToast'
import { useUnsavedChanges } from '@/composables/useUnsavedChanges'
import { ArrowLeft } from 'lucide-vue-next'

defineOptions({ name: 'ShortLinkEditPage' })

const router = useRouter()
const route = useRoute()
const toast = useToast()
const id = route.params.id as string

const url = ref('')
const description = ref('')
const isEnable = ref(true)
const createdBy = ref('')
const createTime = ref('')
const loading = ref(false)
const saving = ref(false)
const error = ref('')
const notFound = ref(false)

const originalData = ref({ url: '', description: '', isEnable: true })

const isDirty = computed(() =>
  url.value !== originalData.value.url ||
  description.value !== originalData.value.description ||
  isEnable.value !== originalData.value.isEnable
)

useUnsavedChanges(isDirty)

function normalizeUrl(value: string): string {
  const trimmed = value.trim()
  if (!trimmed) return ''
  if (/^https?:\/\//i.test(trimmed)) return trimmed
  return `https://${trimmed}`
}

const urlError = computed(() => {
  if (!url.value) return ''
  try {
    new URL(normalizeUrl(url.value))
    return ''
  } catch {
    return 'Please enter a valid URL (e.g. https://example.com)'
  }
})

const canSubmit = computed(() => !!url.value && !urlError.value && !saving.value)

async function fetchData() {
  loading.value = true
  try {
    const data = await getShortLink(id)
    const sl = data.shortLink
    url.value = sl.url
    description.value = sl.description
    isEnable.value = sl.isEnable
    createdBy.value = sl.createdBy
    createTime.value = sl.createTime
    originalData.value = { url: sl.url, description: sl.description, isEnable: sl.isEnable }
  } catch (e) {
    if (e instanceof ApiError && e.status === 404) {
      notFound.value = true
    } else {
      error.value = t('shortLinkEdit.loadFailed')
    }
  } finally {
    loading.value = false
  }
}

async function handleSubmit() {
  if (!canSubmit.value) return
  error.value = ''
  saving.value = true

  try {
    const normalizedUrl = normalizeUrl(url.value)
    await updateShortLink(id, {
      url: normalizedUrl,
      description: description.value,
      isEnable: isEnable.value,
    })
    originalData.value = { url: normalizedUrl, description: description.value, isEnable: isEnable.value }
    toast.success(t('shortLinkEdit.updated'))
  } catch {
    error.value = t('shortLinkEdit.updateFailed')
  } finally {
    saving.value = false
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

onMounted(fetchData)
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
      <div>
        <h1 class="text-xl font-semibold text-gray-900">{{ ('shortLinkEdit.title') }}</h1>
        <p class="mt-0.5 text-sm text-gray-500">
          Editing <span class="font-mono font-medium text-gray-700">{{ id }}</span>
        </p>
      </div>
    </div>

    <!-- Skeleton loading -->
    <div v-if="loading" class="mt-6 max-w-lg rounded-lg border bg-white p-6">
      <div class="mb-5 rounded-md bg-gray-50 p-3">
        <div class="grid grid-cols-2 gap-2">
          <div class="h-4 w-32 animate-pulse rounded bg-gray-200" />
          <div class="h-4 w-28 animate-pulse rounded bg-gray-200" />
          <div class="col-span-2 h-4 w-48 animate-pulse rounded bg-gray-200" />
        </div>
      </div>
      <div class="space-y-5">
        <div>
          <div class="mb-1 h-4 w-16 animate-pulse rounded bg-gray-200" />
          <div class="h-9 w-full animate-pulse rounded-md bg-gray-200" />
        </div>
        <div>
          <div class="mb-1 h-4 w-20 animate-pulse rounded bg-gray-200" />
          <div class="h-20 w-full animate-pulse rounded-md bg-gray-200" />
        </div>
        <div class="flex items-center gap-3">
          <div class="h-4 w-14 animate-pulse rounded bg-gray-200" />
          <div class="h-6 w-11 animate-pulse rounded-full bg-gray-200" />
        </div>
      </div>
    </div>

    <div v-else-if="notFound" class="mt-6 text-center text-gray-400">
      Short link not found.
      <button
        class="ml-2 text-blue-600 hover:underline"
        @click="router.push({ name: 'short-links' })"
      >
        Back to list
      </button>
    </div>

    <form
      v-else
      class="mt-6 max-w-lg rounded-lg border bg-white p-6"
      @submit.prevent="handleSubmit"
    >
      <!-- Read-only info -->
      <div class="mb-5 rounded-md bg-gray-50 p-3">
        <div class="grid grid-cols-2 gap-2 text-sm">
          <div>
            <span class="text-gray-500">{{ ('shortLinkEdit.infoId') }}</span>
            <span class="ml-1 font-mono font-medium text-gray-900">{{ id }}</span>
          </div>
          <div>
            <span class="text-gray-500">{{ ('shortLinkEdit.infoCreatedBy') }}</span>
            <span class="ml-1 text-gray-700">{{ createdBy }}</span>
          </div>
          <div class="col-span-2">
            <span class="text-gray-500">{{ ('shortLinkEdit.infoCreated') }}</span>
            <span class="ml-1 text-gray-700">{{ formatDate(createTime) }}</span>
          </div>
        </div>
      </div>

      <div class="space-y-5">
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700">
            Target URL <span class="text-red-500">*</span>
          </label>
          <input
            v-model="url"
            type="text"
            required
            placeholder="example.com/long-url or https://example.com"
            class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
            :class="urlError ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''"
          />
          <p v-if="urlError" class="mt-1 text-xs text-red-600">{{ urlError }}</p>
          <p v-else class="mt-1 text-xs text-gray-400">
            https:// will be added automatically if omitted.
          </p>
        </div>

        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700">{{ ('shortLinkEdit.description') }}</label>
          <textarea
            v-model="description"
            rows="3"
            class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
          />
        </div>

        <div class="flex items-center gap-3">
          <label class="text-sm font-medium text-gray-700">{{ ('shortLinkEdit.enabled') }}</label>
          <button
            type="button"
            role="switch"
            :aria-checked="isEnable"
            class="relative inline-flex h-6 w-11 shrink-0 rounded-full border-2 border-transparent transition-colors"
            :class="isEnable ? 'bg-blue-600' : 'bg-gray-200'"
            @click="isEnable = !isEnable"
          >
            <span
              class="pointer-events-none inline-block h-5 w-5 rounded-full bg-white shadow-sm ring-0 transition-transform"
              :class="isEnable ? 'translate-x-5' : 'translate-x-0'"
            />
          </button>
        </div>
      </div>

      <p v-if="error" class="mt-4 text-sm text-red-600">{{ error }}</p>

      <div class="mt-6 flex items-center gap-3">
        <button
          type="submit"
          :disabled="!canSubmit"
          class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:opacity-50"
        >
          {{ saving ? 'Saving...' : 'Save Changes' }}
        </button>
        <button
          type="button"
          class="rounded-md px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-100"
          @click="router.push({ name: 'short-links' })"
        >
          Cancel
        </button>
      </div>
    </form>
  </div>
</template>
