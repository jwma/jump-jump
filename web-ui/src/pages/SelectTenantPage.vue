<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { createTenant } from '@/api/tenant'
import { ArrowRight, Plus, Shield, Building2, Loader2 } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const auth = useAuthStore()

const loading = ref(!auth.authUser)
const showCreateDialog = ref(false)
const createLoading = ref(false)
const createForm = reactive({ name: '', slug: '' })
const createError = ref('')
const slugManuallyEdited = ref(false)

function generateSlug(name: string): string {
  return name
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, '')
    .replace(/[\s_]+/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-+|-+$/g, '')
}

watch(
  () => createForm.name,
  (name) => {
    if (!slugManuallyEdited.value) {
      createForm.slug = generateSlug(name)
    }
  },
)

onMounted(async () => {
  if (!auth.authUser) {
    try {
      await auth.fetchAuthInfo()
    } catch {
      router.push({ name: 'login' })
    } finally {
      loading.value = false
    }
  }
})

function selectTenant(tenantId: string) {
  auth.selectTenant(tenantId)
  const redirect = (route.query.redirect as string) || '/'
  router.push(redirect)
}

function enterSuperAdmin() {
  auth.clearTenant()
  const redirect = (route.query.redirect as string) || '/'
  router.push(redirect)
}

function openCreateDialog() {
  createForm.name = ''
  createForm.slug = ''
  createError.value = ''
  slugManuallyEdited.value = false
  showCreateDialog.value = true
}

async function handleCreateTenant() {
  createError.value = ''
  createLoading.value = true
  try {
    const tenant = await createTenant({ name: createForm.name, slug: createForm.slug })
    await auth.fetchAuthInfo()
    showCreateDialog.value = false
    createForm.name = ''
    createForm.slug = ''
    auth.selectTenant(tenant.id)
    router.push({ name: 'dashboard' })
  } catch (e: unknown) {
    createError.value = (e as Error).message || t('selectTenant.failedToCreate')
  } finally {
    createLoading.value = false
  }
}

async function handleLogout() {
  await auth.logout()
  router.push({ name: 'login' })
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-50">
    <div class="w-full max-w-md">
      <template v-if="loading">
        <div class="flex flex-col items-center gap-3">
          <Loader2 class="h-8 w-8 animate-spin text-blue-600" />
          <p class="text-sm text-gray-500">{{ t('common.loading') }}</p>
        </div>
      </template>
      <template v-else>
        <h1 class="mb-2 text-center text-2xl font-bold">{{ t('selectTenant.title') }}</h1>
        <p class="mb-6 text-center text-sm text-gray-500">
          {{ t('selectTenant.welcome', { username: auth.authUser?.username }) }}
        </p>

        <div class="space-y-3">
          <button
            v-for="tenant in auth.tenants"
            :key="tenant.tenantId"
            class="flex w-full items-center justify-between rounded-lg border bg-white p-4 text-left shadow-sm transition hover:border-blue-300 hover:shadow"
            @click="selectTenant(tenant.tenantId)"
          >
            <div class="flex items-center gap-3">
              <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-50">
                <Building2 class="h-5 w-5 text-blue-600" />
              </div>
              <div>
                <div class="font-medium text-gray-900">{{ tenant.tenantName }}</div>
                <span
                  class="inline-block rounded-full px-2 py-0.5 text-xs font-medium"
                  :class="tenant.role === 'admin' ? 'bg-red-50 text-red-700' : 'bg-gray-100 text-gray-600'"
                >
                  {{ tenant.role }}
                </span>
              </div>
            </div>
            <ArrowRight class="h-4 w-4 text-gray-400" />
          </button>
        </div>

        <button
          v-if="auth.isSuper"
          class="mt-3 flex w-full items-center justify-between rounded-lg border-2 border-dashed border-amber-300 bg-amber-50 p-4 text-left transition hover:border-amber-400"
          @click="enterSuperAdmin"
        >
          <div class="flex items-center gap-3">
            <Shield class="h-5 w-5 text-amber-600" />
            <span class="font-medium text-amber-700">{{ t('selectTenant.superAdminPanel') }}</span>
          </div>
          <ArrowRight class="h-4 w-4 text-amber-400" />
        </button>

        <button
          class="mt-4 flex w-full items-center justify-center gap-2 rounded-lg border border-dashed border-gray-300 bg-white p-3 text-sm font-medium text-gray-600 transition hover:border-gray-400 hover:text-gray-800"
          @click="openCreateDialog"
        >
          <Plus class="h-4 w-4" />
          {{ t('selectTenant.createNewTenant') }}
        </button>

        <div class="mt-6 text-center">
          <button class="text-sm text-gray-400 hover:text-gray-600" @click="handleLogout">
            {{ t('selectTenant.signOut') }}
          </button>
        </div>

        <!-- Create tenant dialog -->
        <div
          v-if="showCreateDialog"
          class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
          @click.self="showCreateDialog = false"
        >
          <div class="w-full max-w-sm rounded-lg border bg-white p-6 shadow-lg">
            <h2 class="mb-4 text-lg font-semibold">{{ t('selectTenant.createTenant') }}</h2>
            <form @submit.prevent="handleCreateTenant">
              <div class="mb-3">
                <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('selectTenant.name') }}</label>
                <input
                  v-model="createForm.name"
                  type="text"
                  required
                  autofocus
                  class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
                />
              </div>
              <div class="mb-3">
                <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('selectTenant.slug') }}</label>
                <input
                  v-model="createForm.slug"
                  type="text"
                  required
                  pattern="[a-z0-9-]+"
                  class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
                  @input="slugManuallyEdited = true"
                />
                <p class="mt-1 text-xs text-gray-400">{{ t('selectTenant.slugHint') }}</p>
              </div>
              <p v-if="createError" class="mb-3 text-sm text-red-600">{{ createError }}</p>
              <div class="flex justify-end gap-2">
                <button
                  type="button"
                  class="rounded-md border px-3 py-1.5 text-sm text-gray-600 hover:bg-gray-50"
                  @click="showCreateDialog = false"
                >
                  {{ t('common.cancel') }}
                </button>
                <button
                  type="submit"
                  :disabled="createLoading"
                  class="rounded-md bg-blue-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
                >
                  {{ createLoading ? t('selectTenant.creating') : t('selectTenant.create') }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>
