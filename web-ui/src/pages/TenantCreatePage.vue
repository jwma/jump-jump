<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { createTenant } from '@/api/tenant'
import { ArrowLeft } from 'lucide-vue-next'

defineOptions({ name: 'TenantCreatePage' })

const router = useRouter()

const name = ref('')
const slug = ref('')
const loading = ref(false)
const error = ref('')

const slugError = computed(() => {
  if (!slug.value) return ''
  if (!/^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(slug.value)) {
    return 'Slug must be lowercase, URL-friendly (letters, numbers, hyphens only)'
  }
  return ''
})

const canSubmit = computed(() => !!name.value && !!slug.value && !slugError.value && !loading.value)

watch(name, (val) => {
  slug.value = val
    .toLowerCase()
    .replace(/[^a-z0-9\s-]/g, '')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '')
})

async function handleSubmit() {
  if (!canSubmit.value) return
  error.value = ''
  loading.value = true

  try {
    const tenant = await createTenant({ name: name.value, slug: slug.value })
    router.push({ name: 'tenant-detail', params: { id: tenant.id } })
  } catch {
    error.value = 'Failed to create tenant. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div>
    <div class="flex items-center gap-3">
      <button
        class="rounded p-1.5 text-gray-400 dark:text-gray-500 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800 dark:bg-gray-800 hover:text-gray-600 dark:text-gray-400 dark:text-gray-500 dark:hover:text-gray-300 dark:text-gray-600"
        @click="router.push({ name: 'tenants' })"
      >
        <ArrowLeft class="h-5 w-5" />
      </button>
      <div>
        <h1 class="text-xl font-semibold text-gray-900 dark:text-white">Create Tenant</h1>
        <p class="mt-0.5 text-sm text-gray-500 dark:text-gray-400 dark:text-gray-500">Add a new tenant to the system.</p>
      </div>
    </div>

    <form class="mt-6 max-w-lg rounded-lg border bg-white dark:bg-gray-900 p-6" @submit.prevent="handleSubmit">
      <div class="space-y-5">
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300 dark:text-gray-600">
            Name <span class="text-red-500">*</span>
          </label>
          <input
            v-model="name"
            type="text"
            required
            placeholder="My Organization"
            class="w-full rounded-md border border-gray-300 dark:border-gray-600 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
          />
        </div>

        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300 dark:text-gray-600">
            Slug <span class="text-red-500">*</span>
          </label>
          <input
            v-model="slug"
            type="text"
            required
            placeholder="my-organization"
            class="w-full rounded-md border border-gray-300 dark:border-gray-600 px-3 py-2 font-mono text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
            :class="slugError ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''"
          />
          <p v-if="slugError" class="mt-1 text-xs text-red-600">{{ slugError }}</p>
          <p v-else class="mt-1 text-xs text-gray-400 dark:text-gray-500">URL-friendly identifier (auto-generated from name).</p>
        </div>
      </div>

      <p v-if="error" class="mt-4 text-sm text-red-600">{{ error }}</p>

      <div class="mt-6 flex items-center gap-3">
        <button
          type="submit"
          :disabled="!canSubmit"
          class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:opacity-50"
        >
          {{ loading ? 'Creating...' : 'Create Tenant' }}
        </button>
        <button
          type="button"
          class="rounded-md px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 dark:text-gray-600 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800 dark:bg-gray-800"
          @click="router.push({ name: 'tenants' })"
        >
          Cancel
        </button>
      </div>
    </form>
  </div>
</template>
