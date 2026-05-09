<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { createShortLink } from '@/api/short-link'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import { useUnsavedChanges } from '@/composables/useUnsavedChanges'
import { UserRole } from '@/types/api'
import type { ShortLinkData } from '@/types/api'
import { ArrowLeft, Copy, Check, ExternalLink } from 'lucide-vue-next'

defineOptions({ name: 'ShortLinkCreatePage' })

const router = useRouter()
const auth = useAuthStore()
const toast = useToast()

const url = ref('')
const customId = ref('')
const description = ref('')
const isEnable = ref(true)
const loading = ref(false)
const error = ref('')
const createdLink = ref<ShortLinkData | null>(null)
const copiedShortLink = ref(false)

const isDirty = computed(() =>
  !createdLink.value && (!!url.value || !!customId.value || !!description.value || !isEnable.value)
)

useUnsavedChanges(isDirty)

const isAdmin = computed(() => auth.user?.role === UserRole.Admin)

const CUSTOM_ID_REGEX = /^[a-zA-Z0-9_-]+$/
const CUSTOM_ID_MAX_LENGTH = 64

const customIdError = computed(() => {
  if (!customId.value) return ''
  if (customId.value.length < 2) return 'Custom ID must be at least 2 characters'
  if (customId.value.length > CUSTOM_ID_MAX_LENGTH)
    return `Custom ID must be at most ${CUSTOM_ID_MAX_LENGTH} characters`
  if (!CUSTOM_ID_REGEX.test(customId.value))
    return 'Custom ID can only contain letters, numbers, hyphens, and underscores'
  return ''
})

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

const canSubmit = computed(
  () =>
    !!url.value &&
    !urlError.value &&
    !customIdError.value &&
    !loading.value,
)

function getShortLinkUrl(id: string) {
  return `${window.location.origin}/${id}`
}

function copyCreatedLink() {
  if (!createdLink.value) return
  navigator.clipboard.writeText(getShortLinkUrl(createdLink.value.id)).then(
    () => {
      copiedShortLink.value = true
      setTimeout(() => {
        copiedShortLink.value = false
      }, 2000)
    },
    () => {
      toast.error('Failed to copy link to clipboard')
    },
  )
}

async function handleSubmit() {
  if (!canSubmit.value) return
  error.value = ''
  loading.value = true

  try {
    const normalizedUrl = normalizeUrl(url.value)
    const result = await createShortLink({
      url: normalizedUrl,
      description: description.value,
      isEnable: isEnable.value,
      ...(isAdmin.value && customId.value ? { id: customId.value } : {}),
    })
    createdLink.value = result.shortLink
    toast.success('Short link created successfully')
  } catch {
    error.value = 'Failed to create short link. Please try again.'
  } finally {
    loading.value = false
  }
}

function handleCreateAnother() {
  url.value = ''
  customId.value = ''
  description.value = ''
  isEnable.value = true
  error.value = ''
  createdLink.value = null
}

function goToLinks() {
  router.push({ name: 'short-links' })
}
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
        <h1 class="text-xl font-semibold text-gray-900">Create Short Link</h1>
        <p class="mt-0.5 text-sm text-gray-500">Generate a new short link for your URL.</p>
      </div>
    </div>

    <!-- Success state -->
    <div
      v-if="createdLink"
      class="mt-6 max-w-lg rounded-lg border border-green-200 bg-green-50 p-6"
    >
      <h2 class="text-lg font-semibold text-green-800">Short Link Created!</h2>
      <div class="mt-3 flex items-center gap-2 rounded-md bg-white p-3 shadow-sm">
        <span class="font-mono text-sm font-medium text-gray-900">
          {{ getShortLinkUrl(createdLink.id) }}
        </span>
        <button
          class="rounded p-1 text-gray-400 transition-colors hover:text-green-600"
          title="Copy short link"
          @click="copyCreatedLink"
        >
          <Copy v-if="!copiedShortLink" class="h-4 w-4" />
          <Check v-else class="h-4 w-4 text-green-500" />
        </button>
        <a
          :href="getShortLinkUrl(createdLink.id)"
          target="_blank"
          rel="noopener noreferrer"
          class="rounded p-1 text-gray-400 transition-colors hover:text-blue-600"
          title="Open short link"
        >
          <ExternalLink class="h-4 w-4" />
        </a>
      </div>
      <p class="mt-2 text-sm text-gray-600">
        Destination:
        <a
          :href="createdLink.url"
          target="_blank"
          rel="noopener noreferrer"
          class="text-blue-600 hover:underline"
        >
          {{ createdLink.url }}
        </a>
      </p>
      <div class="mt-4 flex items-center gap-3">
        <button
          class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700"
          @click="handleCreateAnother"
        >
          Create Another
        </button>
        <button
          class="rounded-md px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-100"
          @click="goToLinks"
        >
          View All Links
        </button>
      </div>
    </div>

    <!-- Create form -->
    <form v-else class="mt-6 max-w-lg rounded-lg border bg-white p-6" @submit.prevent="handleSubmit">
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

        <div v-if="isAdmin">
          <label class="mb-1 block text-sm font-medium text-gray-700">
            Custom ID
            <span class="text-xs font-normal text-gray-400">(optional)</span>
          </label>
          <input
            v-model="customId"
            type="text"
            placeholder="my-custom-id"
            maxlength="64"
            class="w-full rounded-md border border-gray-300 px-3 py-2 font-mono text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
            :class="customIdError ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''"
          />
          <p v-if="customIdError" class="mt-1 text-xs text-red-600">{{ customIdError }}</p>
          <p v-else class="mt-1 text-xs text-gray-400">
            Letters, numbers, hyphens, underscores. 2–64 characters. Leave empty to auto-generate.
          </p>
        </div>

        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700">
            Description
            <span class="text-xs font-normal text-gray-400">(optional)</span>
          </label>
          <textarea
            v-model="description"
            rows="3"
            placeholder="Describe where this link points to..."
            class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
          />
        </div>

        <div class="flex items-center gap-3">
          <label class="text-sm font-medium text-gray-700">Enabled</label>
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
          {{ loading ? 'Creating...' : 'Create Short Link' }}
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
