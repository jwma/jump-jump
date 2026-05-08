<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { createShortLink } from '@/api/short-link'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import { UserRole } from '@/types/api'
import { ArrowLeft } from 'lucide-vue-next'

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

const isAdmin = computed(() => auth.user?.role === UserRole.Admin)

const urlError = computed(() => {
  if (!url.value) return ''
  try {
    new URL(url.value)
    return ''
  } catch {
    return 'Please enter a valid URL (e.g. https://example.com)'
  }
})

const canSubmit = computed(() => !!url.value && !urlError.value && !loading.value)

async function handleSubmit() {
  if (!canSubmit.value) return
  error.value = ''
  loading.value = true

  try {
    await createShortLink({
      url: url.value,
      description: description.value,
      isEnable: isEnable.value,
      ...(isAdmin.value && customId.value ? { id: customId.value } : {}),
    })
    toast.success('Short link created successfully')
    router.push({ name: 'short-links' })
  } catch {
    error.value = 'Failed to create short link. Please try again.'
  } finally {
    loading.value = false
  }
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

    <form class="mt-6 max-w-lg rounded-lg border bg-white p-6" @submit.prevent="handleSubmit">
      <div class="space-y-5">
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700">
            Target URL <span class="text-red-500">*</span>
          </label>
          <input
            v-model="url"
            type="url"
            required
            placeholder="https://example.com/long-url"
            class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
            :class="urlError ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''"
          />
          <p v-if="urlError" class="mt-1 text-xs text-red-600">{{ urlError }}</p>
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
            class="w-full rounded-md border border-gray-300 px-3 py-2 font-mono text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
          />
          <p class="mt-1 text-xs text-gray-400">Leave empty to auto-generate a short ID.</p>
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
