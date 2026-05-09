<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { getTenant, updateTenant, listDomains, addDomain, removeDomain } from '@/api/tenant'
import { getTenantConfig, updateTenantConfig } from '@/api/config'
import { useToast } from '@/composables/useToast'
import { useUnsavedChanges } from '@/composables/useUnsavedChanges'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import type { Tenant, TenantDomain } from '@/types/api'
import {
  Settings,
  Loader2,
  Save,
  Hash,
  AlertTriangle,
  Eye,
  Globe,
  Plus,
  Star,
  ExternalLink,
  Trash2,
  Building2,
} from 'lucide-vue-next'

defineOptions({ name: 'SettingsPage' })

const auth = useAuthStore()
const toast = useToast()

const loading = ref(true)
const pageError = ref('')

// Tenant info
const tenant = ref<Tenant | null>(null)
const editName = ref('')
const editSlug = ref('')
const tenantSaving = ref(false)
const tenantError = ref('')

// Domains
const domains = ref<TenantDomain[]>([])
const domainsLoading = ref(false)
const newDomain = ref('')
const newDomainIsDefault = ref(false)
const addingDomain = ref(false)
const confirmDeleteDomain = ref<string | null>(null)
const deleting = ref(false)

// Config
const idLength = ref(6)
const idMinimumLength = ref(2)
const idMaximumLength = ref(10)
const idSaving = ref(false)
const idError = ref('')

const notFoundMode = ref<'content' | 'redirect'>('content')
const notFoundValue = ref('')
const notFoundSaving = ref(false)
const notFoundError = ref('')

// Unsaved changes tracking
const originalTenant = ref({ name: '', slug: '' })
const originalConfig = ref({ idLength: 6, idMinimumLength: 2, idMaximumLength: 10, notFoundMode: 'content' as string, notFoundValue: '' })

const isDirty = computed(() => {
  if (!isAdmin.value || loading.value) return false
  return (
    (!!tenant.value && (editName.value !== originalTenant.value.name || editSlug.value !== originalTenant.value.slug)) ||
    idLength.value !== originalConfig.value.idLength ||
    idMinimumLength.value !== originalConfig.value.idMinimumLength ||
    idMaximumLength.value !== originalConfig.value.idMaximumLength ||
    notFoundMode.value !== originalConfig.value.notFoundMode ||
    notFoundValue.value !== originalConfig.value.notFoundValue
  )
})

useUnsavedChanges(isDirty)

// Confirm dialog
const confirmOpen = ref(false)
const confirmTitle = ref('')
const confirmMessage = ref('')
const confirmAction = ref<(() => Promise<void>) | null>(null)
const confirmSaving = ref(false)

const isAdmin = computed(() => auth.isAdmin)

const idValidationError = computed(() => {
  if (!Number.isInteger(idMinimumLength.value) || !Number.isInteger(idLength.value) || !Number.isInteger(idMaximumLength.value)) return 'All values must be integers'
  if (idMinimumLength.value < 2) return 'Minimum length must be at least 2'
  if (idMaximumLength.value > 10) return 'Maximum length must be at most 10'
  if (idMinimumLength.value >= idLength.value) return 'Minimum length must be less than default length'
  if (idLength.value >= idMaximumLength.value) return 'Default length must be less than maximum length'
  return ''
})

const notFoundValidationError = computed(() => {
  if (!notFoundValue.value.trim()) return notFoundMode.value === 'content' ? 'Content text is required' : 'Redirect URL is required'
  if (notFoundMode.value === 'redirect') {
    try { new URL(notFoundValue.value); return '' } catch { return 'Please enter a valid URL' }
  }
  return ''
})

const idCanSave = computed(() => !idValidationError.value && !idSaving.value && !loading.value)
const notFoundCanSave = computed(() => !notFoundValidationError.value && !notFoundSaving.value && !loading.value)

const sampleIdSeed = ref(0)

const normalizedIdLengths = computed(() => {
  const min = Math.max(2, Math.round(idMinimumLength.value))
  const def = Math.max(min + 1, Math.round(idLength.value))
  const max = Math.min(10, Math.max(def + 1, Math.round(idMaximumLength.value)))
  return { min, default: def, max }
})

function generateId(len: number): string {
  const chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
  let result = ''
  for (let i = 0; i < len; i++) result += chars.charAt(Math.floor(Math.random() * chars.length))
  return result
}

const sampleIds = computed(() => {
  void sampleIdSeed.value
  if (idValidationError.value) return null
  const n = normalizedIdLengths.value
  return {
    min: generateId(n.min),
    default: generateId(n.default),
    max: generateId(n.max),
  }
})

const previewText = computed(() => {
  if (notFoundMode.value === 'content') return notFoundValue.value || '(No content configured)'
  return `Redirect to: ${notFoundValue.value || '(No URL configured)'}`
})

async function fetchAll() {
  if (!auth.currentTenantId) return
  loading.value = true
  pageError.value = ''
  try {
    const [t, configData] = await Promise.all([
      getTenant(auth.currentTenantId),
      getTenantConfig(auth.currentTenantId),
    ])
    tenant.value = t
    editName.value = t.name
    editSlug.value = t.slug
    originalTenant.value = { name: t.name, slug: t.slug }

    idLength.value = configData.config.idConfig.idLength
    idMinimumLength.value = configData.config.idConfig.idMinimumLength
    idMaximumLength.value = configData.config.idConfig.idMaximumLength
    notFoundMode.value = configData.config.shortLinkNotFoundConfig.mode
    notFoundValue.value = configData.config.shortLinkNotFoundConfig.value
    originalConfig.value = {
      idLength: configData.config.idConfig.idLength,
      idMinimumLength: configData.config.idConfig.idMinimumLength,
      idMaximumLength: configData.config.idConfig.idMaximumLength,
      notFoundMode: configData.config.shortLinkNotFoundConfig.mode,
      notFoundValue: configData.config.shortLinkNotFoundConfig.value,
    }
  } catch {
    pageError.value = 'Failed to load settings.'
  } finally {
    loading.value = false
  }

  fetchDomains()
}

async function fetchDomains() {
  if (!auth.currentTenantId) return
  domainsLoading.value = true
  try {
    domains.value = (await listDomains(auth.currentTenantId)) ?? []
  } catch {
    domains.value = []
  } finally {
    domainsLoading.value = false
  }
}

// Tenant info save
async function handleSaveTenant() {
  if (!auth.currentTenantId) return
  tenantSaving.value = true
  tenantError.value = ''
  try {
    const updated = await updateTenant(auth.currentTenantId, {
      name: editName.value,
      slug: editSlug.value,
    })
    tenant.value = updated
    editName.value = updated.name
    editSlug.value = updated.slug
    originalTenant.value = { name: updated.name, slug: updated.slug }
    toast.success('Tenant information saved successfully')
    await auth.fetchAuthInfo()
  } catch (e: unknown) {
    tenantError.value = (e as Error).message || 'Failed to update tenant.'
  } finally {
    tenantSaving.value = false
  }
}

// Domain management
async function handleAddDomain() {
  if (!auth.currentTenantId || !newDomain.value.trim()) return
  addingDomain.value = true
  try {
    const data = await addDomain(auth.currentTenantId, { domain: newDomain.value.trim(), isDefault: newDomainIsDefault.value })
    domains.value = data || []
    newDomain.value = ''
    newDomainIsDefault.value = false
  } catch {
    // handled by interceptor
  } finally {
    addingDomain.value = false
  }
}

async function handleDeleteDomain(domain: string) {
  if (!auth.currentTenantId) return
  deleting.value = true
  try {
    const data = await removeDomain(auth.currentTenantId, domain)
    domains.value = data || []
  } catch {
    // handled by interceptor
  } finally {
    deleting.value = false
    confirmDeleteDomain.value = null
  }
}

// Config save
function openConfirm(title: string, message: string, action: () => Promise<void>) {
  confirmTitle.value = title
  confirmMessage.value = message
  confirmAction.value = action
  confirmSaving.value = false
  confirmOpen.value = true
}

function closeConfirm() {
  confirmOpen.value = false
  confirmAction.value = null
  confirmSaving.value = false
}

async function executeConfirm() {
  if (!confirmAction.value || confirmSaving.value) return
  confirmSaving.value = true
  try {
    await confirmAction.value()
  } finally {
    closeConfirm()
  }
}

function handleSaveIdConfig() {
  openConfirm(
    'Save ID Length Configuration',
    `Set default length to ${idLength.value}, min ${idMinimumLength.value}, max ${idMaximumLength.value}?`,
    doSaveIdConfig,
  )
}

async function doSaveIdConfig() {
  if (!auth.currentTenantId) return
  idSaving.value = true
  idError.value = ''
  try {
    await updateTenantConfig(auth.currentTenantId!, { idLength: idLength.value, idMinimumLength: idMinimumLength.value, idMaximumLength: idMaximumLength.value })
    originalConfig.value = { ...originalConfig.value, idLength: idLength.value, idMinimumLength: idMinimumLength.value, idMaximumLength: idMaximumLength.value }
    toast.success('ID length configuration saved successfully')
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
  if (!auth.currentTenantId) return
  notFoundSaving.value = true
  notFoundError.value = ''
  try {
    await updateTenantConfig(auth.currentTenantId!, { notFoundMode: notFoundMode.value, notFoundValue: notFoundValue.value })
    originalConfig.value = { ...originalConfig.value, notFoundMode: notFoundMode.value, notFoundValue: notFoundValue.value }
    toast.success('404 handling configuration saved successfully')
  } catch {
    notFoundError.value = 'Failed to save 404 handling configuration.'
  } finally {
    notFoundSaving.value = false
  }
}

function formatDate(d: string) {
  return new Date(d).toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function refreshSampleIds() {
  sampleIdSeed.value++
}

onMounted(fetchAll)
</script>

<template>
  <div class="space-y-4 sm:space-y-6">
    <!-- Header -->
    <div>
      <h1 class="flex items-center gap-2 text-xl font-semibold text-gray-900 dark:text-white">
        <Settings class="h-5 w-5 text-gray-500 dark:text-gray-400 dark:text-gray-500" />
        Settings
      </h1>
      <p class="mt-1 text-sm text-gray-500 dark:text-gray-400 dark:text-gray-500">Manage tenant settings, domains, and system configuration.</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-20">
      <Loader2 class="h-6 w-6 animate-spin text-gray-400 dark:text-gray-500" />
    </div>

    <div v-else-if="pageError" class="rounded-lg border border-red-200 bg-red-50 p-4">
      <p class="text-sm text-red-600">{{ pageError }}</p>
    </div>

    <template v-else>
      <!-- Tenant Basic Info (admin only) -->
      <div v-if="isAdmin && tenant" class="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 shadow-sm">
        <div class="flex items-center gap-2 border-b border-gray-100 dark:border-gray-800 px-4 py-3 sm:px-5 sm:py-4">
          <Building2 class="h-4 w-4 text-blue-600" />
          <h2 class="text-base font-semibold text-gray-900 dark:text-white">Tenant Information</h2>
        </div>
        <div class="p-4 sm:p-5">
          <div class="space-y-4">
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300 dark:text-gray-600">Name</label>
              <input
                v-model="editName"
                type="text"
                class="w-full rounded-md border border-gray-300 dark:border-gray-600 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
              />
            </div>
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300 dark:text-gray-600">Slug</label>
              <input
                v-model="editSlug"
                type="text"
                pattern="[a-z0-9-]+"
                class="w-full rounded-md border border-gray-300 dark:border-gray-600 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
              />
              <p class="mt-1 text-xs text-gray-400 dark:text-gray-500">Lowercase letters, numbers, and hyphens</p>
            </div>
            <p v-if="tenantError" class="text-xs text-red-600">{{ tenantError }}</p>
          </div>
          <div class="mt-5">
            <button
              :disabled="tenantSaving"
              class="inline-flex items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:opacity-50"
              @click="handleSaveTenant"
            >
              <Loader2 v-if="tenantSaving" class="h-4 w-4 animate-spin" />
              <Save v-else class="h-4 w-4" />
              {{ tenantSaving ? 'Saving...' : 'Save Tenant Info' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Domain Management (admin only) -->
      <div v-if="isAdmin" class="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 shadow-sm">
        <div class="flex items-center gap-2 border-b border-gray-100 dark:border-gray-800 px-4 py-3 sm:px-5 sm:py-4">
          <Globe class="h-4 w-4 text-green-600" />
          <h2 class="text-base font-semibold text-gray-900 dark:text-white">Domains</h2>
          <span class="rounded-full bg-gray-100 dark:bg-gray-800 px-2 py-0.5 text-xs font-medium text-gray-600 dark:text-gray-400 dark:text-gray-500">
            {{ domains.length }}
          </span>
        </div>
        <div class="p-4 sm:p-5">
          <!-- Add domain form -->
          <div class="rounded-lg border bg-gray-50 dark:bg-gray-800/50 p-4">
            <div class="flex flex-col gap-3 sm:flex-row sm:items-end">
              <div class="flex-1">
                <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300 dark:text-gray-600">Domain</label>
                <input
                  v-model="newDomain"
                  type="text"
                  placeholder="example.com"
                  class="w-full rounded-md border border-gray-300 dark:border-gray-600 px-3 py-2 font-mono text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
                  @keydown.enter.prevent="handleAddDomain"
                />
              </div>
              <div class="flex items-center gap-3">
                <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 dark:text-gray-600">
                  <input v-model="newDomainIsDefault" type="checkbox" class="h-4 w-4 rounded border-gray-300 dark:border-gray-600 text-blue-600 focus:ring-blue-500" />
                  Default
                </label>
                <button
                  :disabled="!newDomain.trim() || addingDomain"
                  class="inline-flex items-center gap-1.5 rounded-md bg-blue-600 px-3 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:opacity-50"
                  @click="handleAddDomain"
                >
                  <Loader2 v-if="addingDomain" class="h-3.5 w-3.5 animate-spin" />
                  <Plus v-else class="h-3.5 w-3.5" />
                  Add
                </button>
              </div>
            </div>
          </div>

          <!-- Domain list -->
          <div v-if="domainsLoading" class="mt-4 flex items-center justify-center py-8">
            <Loader2 class="h-5 w-5 animate-spin text-gray-400 dark:text-gray-500" />
          </div>
          <div v-else-if="domains.length > 0" class="mt-4 overflow-hidden rounded-lg border">
            <!-- Mobile cards -->
            <div class="divide-y lg:hidden">
              <div
                v-for="d in domains"
                :key="d.id"
                class="flex items-center justify-between px-4 py-3"
              >
                <div class="min-w-0 flex-1">
                  <div class="flex items-center gap-2">
                    <Globe class="h-4 w-4 text-gray-400 dark:text-gray-500" />
                    <span class="font-mono text-sm text-gray-900 dark:text-white">{{ d.domain }}</span>
                  </div>
                  <div class="mt-1 flex items-center gap-2 pl-6">
                    <span v-if="d.isDefault" class="inline-flex items-center gap-1 text-xs font-medium text-yellow-600">
                      <Star class="h-3 w-3 fill-yellow-500 text-yellow-500" />
                      Default
                    </span>
                    <span class="text-xs text-gray-400 dark:text-gray-500">{{ formatDate(d.createdAt) }}</span>
                  </div>
                </div>
                <div class="flex items-center gap-1">
                  <a
                    :href="'https://' + d.domain"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="rounded p-1.5 text-gray-400 dark:text-gray-500 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800 dark:bg-gray-800 hover:text-gray-600 dark:text-gray-400 dark:text-gray-500 dark:hover:text-gray-300 dark:text-gray-600"
                    title="Open domain"
                    aria-label="Open domain"
                  >
                    <ExternalLink class="h-4 w-4" />
                  </a>
                  <button
                    class="rounded p-1.5 text-gray-400 dark:text-gray-500 transition-colors hover:bg-red-50 hover:text-red-600"
                    title="Delete domain"
                    aria-label="Delete domain"
                    @click="confirmDeleteDomain = d.domain"
                  >
                    <Trash2 class="h-4 w-4" />
                  </button>
                </div>
              </div>
            </div>
            <!-- Desktop table -->
            <table class="hidden w-full text-sm lg:table">
              <thead>
                <tr class="border-b bg-gray-50 dark:bg-gray-800/50 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400 dark:text-gray-500">
                  <th class="px-4 py-3">Domain</th>
                  <th class="px-4 py-3">Default</th>
                  <th class="hidden px-4 py-3 lg:table-cell">Created</th>
                  <th class="px-4 py-3 text-right">Actions</th>
                </tr>
              </thead>
              <tbody class="divide-y">
                <tr v-for="d in domains" :key="d.id" class="transition-colors hover:bg-gray-50 dark:hover:bg-gray-700 dark:bg-gray-800/50">
                  <td class="px-4 py-3">
                    <div class="flex items-center gap-2">
                      <Globe class="h-4 w-4 text-gray-400 dark:text-gray-500" />
                      <span class="font-mono text-sm text-gray-900 dark:text-white">{{ d.domain }}</span>
                    </div>
                  </td>
                  <td class="px-4 py-3">
                    <span v-if="d.isDefault" class="inline-flex items-center gap-1 text-xs font-medium text-yellow-600">
                      <Star class="h-3.5 w-3.5 fill-yellow-500 text-yellow-500" />
                      Default
                    </span>
                    <span v-else class="text-xs text-gray-400 dark:text-gray-500">&mdash;</span>
                  </td>
                  <td class="hidden whitespace-nowrap px-4 py-3 text-gray-500 dark:text-gray-400 dark:text-gray-500 lg:table-cell">{{ formatDate(d.createdAt) }}</td>
                  <td class="px-4 py-3">
                    <div class="flex items-center justify-end gap-1">
                      <a
                        :href="'https://' + d.domain"
                        target="_blank"
                        rel="noopener noreferrer"
                        class="rounded p-1.5 text-gray-400 dark:text-gray-500 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800 dark:bg-gray-800 hover:text-gray-600 dark:text-gray-400 dark:text-gray-500 dark:hover:text-gray-300 dark:text-gray-600"
                      >
                        <ExternalLink class="h-4 w-4" />
                      </a>
                      <button
                        class="rounded p-1.5 text-gray-400 dark:text-gray-500 transition-colors hover:bg-red-50 hover:text-red-600"
                        @click="confirmDeleteDomain = d.domain"
                      >
                        <Trash2 class="h-4 w-4" />
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <div v-else class="mt-4 py-6 text-center text-sm text-gray-400 dark:text-gray-500">
            No domains configured for this tenant.
          </div>
        </div>
      </div>

      <!-- ID Length Config (admin only) -->
      <div v-if="isAdmin" class="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 shadow-sm">
        <div class="flex items-center gap-2 border-b border-gray-100 dark:border-gray-800 px-4 py-3 sm:px-5 sm:py-4">
          <Hash class="h-4 w-4 text-blue-600" />
          <h2 class="text-base font-semibold text-gray-900 dark:text-white">ID Length Configuration</h2>
        </div>
        <div class="p-4 sm:p-5">
          <div class="max-w-md space-y-4">
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300 dark:text-gray-600">Default ID Length</label>
              <input
                v-model.number="idLength"
                type="number"
                step="1"
                :min="idMinimumLength + 1"
                :max="idMaximumLength - 1"
                class="w-full rounded-md border border-gray-300 dark:border-gray-600 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
              />
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300 dark:text-gray-600">Minimum Length</label>
                <input
                  v-model.number="idMinimumLength"
                  type="number"
                  step="1"
                  :min="2"
                  :max="idLength - 1"
                  class="w-full rounded-md border border-gray-300 dark:border-gray-600 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
                />
              </div>
              <div>
                <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300 dark:text-gray-600">Maximum Length</label>
                <input
                  v-model.number="idMaximumLength"
                  type="number"
                  step="1"
                  :min="idLength + 1"
                  :max="10"
                  class="w-full rounded-md border border-gray-300 dark:border-gray-600 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
                />
              </div>
            </div>
            <p v-if="idValidationError" class="text-xs text-red-600">{{ idValidationError }}</p>
            <p v-if="idError" class="text-xs text-red-600">{{ idError }}</p>
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

      <!-- 404 Config (admin only) -->
      <div v-if="isAdmin" class="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 shadow-sm">
        <div class="flex items-center gap-2 border-b border-gray-100 dark:border-gray-800 px-4 py-3 sm:px-5 sm:py-4">
          <AlertTriangle class="h-4 w-4 text-orange-500" />
          <h2 class="text-base font-semibold text-gray-900 dark:text-white">404 Handling Configuration</h2>
        </div>
        <div class="p-4 sm:p-5">
          <div class="max-w-md space-y-4">
            <div>
              <label id="mode-label" class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300 dark:text-gray-600">Handling Mode</label>
              <div role="radiogroup" aria-labelledby="mode-label" class="flex rounded-lg border border-gray-200 dark:border-gray-700 p-0.5">
                <button
                  role="radio"
                  :aria-checked="notFoundMode === 'content'"
                  :class="['flex-1 rounded-md px-4 py-2 text-sm font-medium transition-colors', notFoundMode === 'content' ? 'bg-blue-600 text-white' : 'text-gray-500 dark:text-gray-400 dark:text-gray-500 hover:text-gray-700 dark:hover:text-gray-200 dark:text-gray-300 dark:text-gray-600']"
                  @click="notFoundMode = 'content'"
                >
                  Display Content
                </button>
                <button
                  role="radio"
                  :aria-checked="notFoundMode === 'redirect'"
                  :class="['flex-1 rounded-md px-4 py-2 text-sm font-medium transition-colors', notFoundMode === 'redirect' ? 'bg-blue-600 text-white' : 'text-gray-500 dark:text-gray-400 dark:text-gray-500 hover:text-gray-700 dark:hover:text-gray-200 dark:text-gray-300 dark:text-gray-600']"
                  @click="notFoundMode = 'redirect'"
                >
                  Redirect
                </button>
              </div>
            </div>
            <div v-if="notFoundMode === 'content'">
              <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300 dark:text-gray-600">Content Text</label>
              <textarea
                v-model="notFoundValue"
                rows="4"
                placeholder="The page content to display when a short link is not found..."
                class="w-full rounded-md border border-gray-300 dark:border-gray-600 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
              />
            </div>
            <div v-else>
              <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300 dark:text-gray-600">Redirect URL</label>
              <input
                v-model="notFoundValue"
                type="url"
                placeholder="https://example.com/404"
                class="w-full rounded-md border border-gray-300 dark:border-gray-600 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
              />
            </div>
            <p v-if="notFoundValidationError" class="text-xs text-red-600">{{ notFoundValidationError }}</p>
            <p v-if="notFoundError" class="text-xs text-red-600">{{ notFoundError }}</p>
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

      <!-- Config Preview (admin only) -->
      <div v-if="isAdmin" class="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 shadow-sm">
        <div class="flex items-center gap-2 border-b border-gray-100 dark:border-gray-800 px-4 py-3 sm:px-5 sm:py-4">
          <Eye class="h-4 w-4 text-purple-600" />
          <h2 class="text-base font-semibold text-gray-900 dark:text-white">Configuration Preview</h2>
        </div>
        <div class="p-4 sm:p-5">
          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div class="rounded-lg border border-gray-100 dark:border-gray-800 bg-gray-50 dark:bg-gray-800/50 p-4">
              <div class="mb-3 flex items-center justify-between">
                <p class="text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400 dark:text-gray-500">ID Length Examples</p>
                <button v-if="sampleIds" class="rounded p-1 text-gray-400 dark:text-gray-500 transition-colors hover:bg-gray-200 dark:bg-gray-700 hover:text-gray-600 dark:text-gray-400 dark:text-gray-500 dark:hover:text-gray-300 dark:text-gray-600" @click="refreshSampleIds">
                  <Loader2 class="h-3.5 w-3.5" />
                </button>
              </div>
              <template v-if="sampleIds">
                <div class="space-y-2">
                  <div class="flex items-center gap-3">
                    <span class="w-16 text-xs text-gray-400 dark:text-gray-500">Min ({{ normalizedIdLengths.min }})</span>
                    <span class="font-mono text-sm font-semibold text-gray-600 dark:text-gray-400 dark:text-gray-500">{{ sampleIds.min }}</span>
                  </div>
                  <div class="flex items-center gap-3">
                    <span class="w-16 text-xs font-medium text-blue-600">Default ({{ normalizedIdLengths.default }})</span>
                    <span class="font-mono text-sm font-bold text-blue-600">{{ sampleIds.default }}</span>
                  </div>
                  <div class="flex items-center gap-3">
                    <span class="w-16 text-xs text-gray-400 dark:text-gray-500">Max ({{ normalizedIdLengths.max }})</span>
                    <span class="font-mono text-sm font-semibold text-gray-600 dark:text-gray-400 dark:text-gray-500">{{ sampleIds.max }}</span>
                  </div>
                </div>
              </template>
              <p v-else class="text-xs text-gray-400 dark:text-gray-500">Fix validation errors to see examples.</p>
            </div>
            <div class="rounded-lg border border-gray-100 dark:border-gray-800 bg-gray-50 dark:bg-gray-800/50 p-4">
              <p class="mb-2 text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400 dark:text-gray-500">404 Handling</p>
              <span
                class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium"
                :class="notFoundMode === 'content' ? 'bg-blue-50 text-blue-700' : 'bg-orange-50 text-orange-700'"
              >
                {{ notFoundMode === 'content' ? 'Display Content' : 'Redirect' }}
              </span>
              <p class="mt-2 line-clamp-3 text-xs text-gray-600 dark:text-gray-400 dark:text-gray-500">{{ previewText }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Non-admin message -->
      <div v-if="!isAdmin" class="rounded-lg border bg-white dark:bg-gray-900 p-8 text-center text-sm text-gray-400 dark:text-gray-500">
        Only admins can manage tenant settings.
      </div>
    </template>

    <!-- Delete domain confirm -->
    <ConfirmDialog
      :open="!!confirmDeleteDomain"
      title="Remove Domain"
      variant="danger"
      :confirm-text="deleting ? 'Removing...' : 'Remove'"
      :loading="deleting"
      @confirm="handleDeleteDomain(confirmDeleteDomain!)"
      @cancel="confirmDeleteDomain = null"
    >
      <p class="mt-2 text-sm text-gray-500 dark:text-gray-400 dark:text-gray-500">
        Are you sure you want to remove
        <span class="font-mono font-medium text-gray-700 dark:text-gray-300 dark:text-gray-600">{{ confirmDeleteDomain }}</span>?
      </p>
    </ConfirmDialog>

    <!-- Generic confirm dialog -->
    <ConfirmDialog
      :open="confirmOpen"
      :title="confirmTitle"
      :message="confirmMessage"
      :confirm-text="confirmSaving ? 'Saving...' : 'Confirm'"
      :loading="confirmSaving"
      @confirm="executeConfirm"
      @cancel="closeConfirm"
    />
  </div>
</template>
