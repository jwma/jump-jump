<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
import { ref, computed, onMounted } from 'vue'
import {
  listSuperUsers,
  createSuperUser,
  resetSuperUserPassword,
  setSuperUserStatus,
  getSuperUser,
} from '@/api/super'
import type { SuperUser, SuperUserDetail } from '@/types/api'
import {
  Users,
  Search,
  Plus,
  KeyRound,
  Loader2,
  ChevronLeft,
  ChevronRight,
  Eye,
  X,
  Shield,
} from 'lucide-vue-next'

defineOptions({ name: 'SuperUsersPage' })

const users = ref<SuperUser[]>([])
const loading = ref(false)
const search = ref('')
const page = ref(1)
const pageSize = 20
const total = ref(0)

const showCreateModal = ref(false)
const createLoading = ref(false)
const newUsername = ref('')
const newPassword = ref('')
const createError = ref('')

const showResetModal = ref(false)
const resetUserId = ref('')
const resetUsername = ref('')
const resetNewPassword = ref('')
const resetLoading = ref(false)
const resetError = ref('')

const showDetailModal = ref(false)
const detailLoading = ref(false)
const userDetail = ref<SuperUserDetail | null>(null)

const toggleLoading = ref<string | null>(null)

async function fetchUsers() {
  loading.value = true
  try {
    const data = await listSuperUsers(page.value, pageSize, search.value)
    users.value = data.users || []
    total.value = data.total
  } catch {
    // error handled by interceptor
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  page.value = 1
  fetchUsers()
}

function goToPage(p: number) {
  if (p < 1 || p > totalPages.value) return
  page.value = p
  fetchUsers()
}

const totalPages = computed(() => Math.ceil(total.value / pageSize) || 1)

async function handleCreate() {
  if (!newUsername.value.trim() || !newPassword.value.trim()) return
  createLoading.value = true
  createError.value = ''
  try {
    await createSuperUser({
      username: newUsername.value.trim(),
      password: newPassword.value.trim(),
    })
    showCreateModal.value = false
    newUsername.value = ''
    newPassword.value = ''
    fetchUsers()
  } catch (e) {
    createError.value = e instanceof Error ? e.message : 'Failed to create user'
  } finally {
    createLoading.value = false
  }
}

function openResetModal(user: SuperUser) {
  resetUserId.value = user.id
  resetUsername.value = user.username
  resetNewPassword.value = ''
  resetError.value = ''
  showResetModal.value = true
}

async function handleReset() {
  if (!resetNewPassword.value.trim()) return
  resetLoading.value = true
  resetError.value = ''
  try {
    await resetSuperUserPassword(resetUserId.value, { newPassword: resetNewPassword.value.trim() })
    showResetModal.value = false
  } catch (e) {
    resetError.value = e instanceof Error ? e.message : 'Failed to reset password'
  } finally {
    resetLoading.value = false
  }
}

async function handleToggleStatus(user: SuperUser) {
  toggleLoading.value = user.id
  try {
    await setSuperUserStatus(user.id, { isActive: !user.isActive })
    fetchUsers()
  } catch {
    // error handled by interceptor
  } finally {
    toggleLoading.value = null
  }
}

async function openDetailModal(user: SuperUser) {
  detailLoading.value = true
  showDetailModal.value = true
  try {
    userDetail.value = await getSuperUser(user.id)
  } catch {
    // error handled by interceptor
  } finally {
    detailLoading.value = false
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

onMounted(fetchUsers)
</script>

<template>
  <div>
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-xl font-semibold text-gray-900 dark:text-white">
          {{ t('super.users.title') }}
        </h1>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          {{ t('super.users.userCount', { count: total, suffix: total !== 1 ? 's' : '' }) }}
        </p>
      </div>
      <button
        class="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700"
        @click="showCreateModal = true"
      >
        <Plus class="h-4 w-4" />
        Create User
      </button>
    </div>

    <!-- Search -->
    <div
      class="mt-4 flex flex-col gap-3 rounded-lg border bg-white dark:bg-gray-900 p-4 sm:flex-row sm:items-center"
    >
      <div class="relative flex-1">
        <Search
          class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400 dark:text-gray-500"
        />
        <input
          v-model="search"
          type="text"
          :placeholder="t('super.users.searchPlaceholder')"
          class="w-full rounded-md border border-gray-300 dark:border-gray-600 py-2 pl-9 pr-3 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
          @keydown.enter="handleSearch"
        />
      </div>
      <button
        class="rounded-md bg-gray-100 dark:bg-gray-800 px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 transition-colors hover:bg-gray-200 dark:bg-gray-700"
        @click="handleSearch"
      >
        Search
      </button>
    </div>

    <!-- Table -->
    <div class="mt-4 overflow-hidden rounded-lg border bg-white dark:bg-gray-900">
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr
              class="border-b bg-gray-50 dark:bg-gray-800/50 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400"
            >
              <th class="px-4 py-3">{{ t('super.users.username') }}</th>
              <th class="px-4 py-3">{{ t('common.status') }}</th>
              <th class="hidden px-4 py-3 md:table-cell">{{ t('super.users.superAdmin') }}</th>
              <th class="hidden px-4 py-3 lg:table-cell">{{ t('common.created') }}</th>
              <th class="px-4 py-3 text-right">{{ t('common.actions') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y">
            <tr v-if="loading">
              <td colspan="5" class="px-4 py-8 text-center text-gray-400 dark:text-gray-500">
                <Loader2 class="inline h-5 w-5 animate-spin" />
              </td>
            </tr>
            <tr v-else-if="users.length === 0">
              <td colspan="5" class="px-4 py-12 text-center">
                <Users class="mx-auto h-10 w-10 text-gray-300 dark:text-gray-600" />
                <p class="mt-3 text-sm font-medium text-gray-500 dark:text-gray-400">
                  {{ t('super.users.noUsers') }}
                </p>
                <p class="mt-1 text-sm text-gray-400 dark:text-gray-500">
                  {{ t('super.users.createToStart') }}
                </p>
                <button
                  class="mt-3 inline-flex items-center gap-1.5 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700"
                  @click="showCreateModal = true"
                >
                  <Plus class="h-4 w-4" />
                  Create User
                </button>
              </td>
            </tr>
            <tr
              v-for="user in users"
              :key="user.id"
              class="transition-colors hover:bg-gray-50 dark:hover:bg-gray-700 dark:bg-gray-800/50"
            >
              <td class="px-4 py-3">
                <div class="flex items-center gap-2">
                  <div
                    class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-gray-100 dark:bg-gray-800 text-xs font-bold text-gray-600 dark:text-gray-400"
                  >
                    {{ user.username.charAt(0).toUpperCase() }}
                  </div>
                  <div>
                    <p class="font-medium text-gray-900 dark:text-white">{{ user.username }}</p>
                    <p class="text-xs text-gray-400 dark:text-gray-500">{{ user.id }}</p>
                  </div>
                </div>
              </td>
              <td class="px-4 py-3">
                <button
                  :disabled="toggleLoading === user.id || user.isSuper"
                  class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium transition-colors"
                  :class="[
                    user.isActive
                      ? 'bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-400 dark:hover:bg-green-800/30 hover:bg-green-100'
                      : 'bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400 hover:bg-gray-200 dark:bg-gray-700',
                    user.isSuper || toggleLoading === user.id
                      ? 'cursor-not-allowed opacity-50'
                      : 'cursor-pointer',
                  ]"
                  @click="handleToggleStatus(user)"
                >
                  <Loader2 v-if="toggleLoading === user.id" class="h-3 w-3 animate-spin" />
                  <span
                    v-else
                    class="h-1.5 w-1.5 rounded-full"
                    :class="user.isActive ? 'bg-green-500' : 'bg-gray-400'"
                  />
                  {{ user.isActive ? t('common.active') : t('common.disabled') }}
                </button>
              </td>
              <td class="hidden px-4 py-3 md:table-cell">
                <span
                  v-if="user.isSuper"
                  class="inline-flex items-center gap-1 text-xs font-medium text-indigo-600"
                >
                  <Shield class="h-3.5 w-3.5" />
                  Super
                </span>
                <span v-else class="text-xs text-gray-400 dark:text-gray-500">-</span>
              </td>
              <td
                class="hidden whitespace-nowrap px-4 py-3 text-gray-500 dark:text-gray-400 lg:table-cell"
              >
                {{ formatDate(user.createdAt) }}
              </td>
              <td class="px-4 py-3">
                <div class="flex items-center justify-end gap-1">
                  <button
                    class="rounded p-1.5 text-gray-400 dark:text-gray-500 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800 hover:text-gray-600 dark:hover:text-gray-300"
                    title="View details"
                    @click="openDetailModal(user)"
                  >
                    <Eye class="h-4 w-4" />
                  </button>
                  <button
                    class="rounded p-1.5 text-gray-400 dark:text-gray-500 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800 hover:text-gray-600 dark:hover:text-gray-300"
                    title="Reset password"
                    @click="openResetModal(user)"
                  >
                    <KeyRound class="h-4 w-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="total > pageSize" class="flex items-center justify-between border-t px-4 py-3">
        <p class="text-xs text-gray-500 dark:text-gray-400">
          {{
            t('super.users.pageOf', {
              current: page,
              total: Math.ceil(total / pageSize),
              count: total,
            })
          }}
        </p>
        <div class="flex items-center gap-1">
          <button
            :disabled="page <= 1"
            class="rounded p-1.5 text-gray-400 dark:text-gray-500 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800 hover:text-gray-600 dark:hover:text-gray-300 disabled:opacity-50"
            @click="goToPage(page - 1)"
          >
            <ChevronLeft class="h-4 w-4" />
          </button>
          <button
            :disabled="page >= Math.ceil(total / pageSize)"
            class="rounded p-1.5 text-gray-400 dark:text-gray-500 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800 hover:text-gray-600 dark:hover:text-gray-300 disabled:opacity-50"
            @click="goToPage(page + 1)"
          >
            <ChevronRight class="h-4 w-4" />
          </button>
        </div>
      </div>
    </div>

    <!-- Create User Modal -->
    <Teleport to="body">
      <div
        v-if="showCreateModal"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
        @click.self="showCreateModal = false"
      >
        <div class="mx-4 w-full max-w-md rounded-lg bg-white dark:bg-gray-900 p-6 shadow-xl">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('super.users.createUser') }}
            </h3>
            <button
              class="rounded p-1 text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300"
              @click="showCreateModal = false"
            >
              <X class="h-5 w-5" />
            </button>
          </div>
          <div class="mt-4 space-y-4">
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{
                t('super.users.username')
              }}</label>
              <input
                v-model="newUsername"
                type="text"
                class="w-full rounded-md border border-gray-300 dark:border-gray-600 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
                @keydown.enter="handleCreate"
              />
            </div>
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{
                t('auth.password')
              }}</label>
              <input
                v-model="newPassword"
                type="password"
                class="w-full rounded-md border border-gray-300 dark:border-gray-600 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
                @keydown.enter="handleCreate"
              />
            </div>
            <p v-if="createError" class="text-sm text-red-600">{{ createError }}</p>
          </div>
          <div class="mt-5 flex justify-end gap-2">
            <button
              class="rounded-md px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800"
              @click="showCreateModal = false"
            >
              Cancel
            </button>
            <button
              :disabled="createLoading || !newUsername.trim() || !newPassword.trim()"
              class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:opacity-50"
              @click="handleCreate"
            >
              <Loader2 v-if="createLoading" class="mr-1 inline h-4 w-4 animate-spin" />
              Create
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Reset Password Modal -->
    <Teleport to="body">
      <div
        v-if="showResetModal"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
        @click.self="showResetModal = false"
      >
        <div class="mx-4 w-full max-w-md rounded-lg bg-white dark:bg-gray-900 p-6 shadow-xl">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('super.users.resetPassword') }}
            </h3>
            <button
              class="rounded p-1 text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300"
              @click="showResetModal = false"
            >
              <X class="h-5 w-5" />
            </button>
          </div>
          <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">
            {{ t('super.users.resetFor', { username: resetUsername }) }}
          </p>
          <div class="mt-4">
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{
              t('super.users.newPassword')
            }}</label>
            <input
              v-model="resetNewPassword"
              type="password"
              class="w-full rounded-md border border-gray-300 dark:border-gray-600 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
              @keydown.enter="handleReset"
            />
            <p v-if="resetError" class="mt-2 text-sm text-red-600">{{ resetError }}</p>
          </div>
          <div class="mt-5 flex justify-end gap-2">
            <button
              class="rounded-md px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800"
              @click="showResetModal = false"
            >
              Cancel
            </button>
            <button
              :disabled="resetLoading || !resetNewPassword.trim()"
              class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:opacity-50"
              @click="handleReset"
            >
              <Loader2 v-if="resetLoading" class="mr-1 inline h-4 w-4 animate-spin" />
              Reset
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- User Detail Modal -->
    <Teleport to="body">
      <div
        v-if="showDetailModal"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
        @click.self="showDetailModal = false"
      >
        <div class="mx-4 w-full max-w-lg rounded-lg bg-white dark:bg-gray-900 p-6 shadow-xl">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('super.users.userDetail') }}
            </h3>
            <button
              class="rounded p-1 text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300"
              @click="showDetailModal = false"
            >
              <X class="h-5 w-5" />
            </button>
          </div>

          <div v-if="detailLoading" class="mt-4 flex justify-center py-8">
            <Loader2 class="h-5 w-5 animate-spin text-gray-400 dark:text-gray-500" />
          </div>

          <template v-else-if="userDetail">
            <div class="mt-4 space-y-3">
              <div class="flex items-center gap-3">
                <div
                  class="flex h-10 w-10 items-center justify-center rounded-full bg-gray-100 dark:bg-gray-800 text-sm font-bold text-gray-600 dark:text-gray-400"
                >
                  {{ userDetail.username.charAt(0).toUpperCase() }}
                </div>
                <div>
                  <p class="font-medium text-gray-900 dark:text-white">{{ userDetail.username }}</p>
                  <p class="text-xs text-gray-400 dark:text-gray-500 font-mono">
                    {{ userDetail.id }}
                  </p>
                </div>
                <span
                  class="ml-auto inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium"
                  :class="
                    userDetail.isActive
                      ? 'bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-400'
                      : 'bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400'
                  "
                >
                  <span
                    class="h-1.5 w-1.5 rounded-full"
                    :class="userDetail.isActive ? 'bg-green-500' : 'bg-gray-400'"
                  />
                  {{ userDetail.isActive ? t('common.active') : t('common.disabled') }}
                </span>
              </div>

              <div class="rounded-lg border bg-gray-50 dark:bg-gray-800/50 p-3 text-sm">
                <div class="grid gap-2 sm:grid-cols-2">
                  <div>
                    <span class="text-xs text-gray-400 dark:text-gray-500">{{
                      t('super.users.superAdmin')
                    }}</span>
                    <p
                      class="font-medium"
                      :class="
                        userDetail.isSuper ? 'text-indigo-600' : 'text-gray-600 dark:text-gray-400'
                      "
                    >
                      {{ userDetail.isSuper ? t('common.yes') : t('common.no') }}
                    </p>
                  </div>
                  <div>
                    <span class="text-xs text-gray-400 dark:text-gray-500">{{
                      t('common.created')
                    }}</span>
                    <p class="font-medium text-gray-600 dark:text-gray-400">
                      {{ formatDate(userDetail.createdAt) }}
                    </p>
                  </div>
                </div>
              </div>

              <div>
                <h4 class="mb-2 text-sm font-medium text-gray-700 dark:text-gray-300">
                  {{ t('super.users.tenantMemberships') }}
                </h4>
                <div
                  v-if="userDetail.tenants.length === 0"
                  class="text-sm text-gray-400 dark:text-gray-500"
                >
                  {{ t('super.users.noMemberships') }}
                </div>
                <div v-else class="space-y-1.5">
                  <div
                    v-for="t in userDetail.tenants"
                    :key="t.tenantId"
                    class="flex items-center justify-between rounded-md border dark:border-gray-700 bg-white dark:bg-gray-800 px-3 py-2"
                  >
                    <span class="font-mono text-sm text-gray-700 dark:text-gray-300">{{
                      t.tenantId
                    }}</span>
                    <span
                      class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium"
                      :class="
                        t.role === 'admin'
                          ? 'bg-blue-50 dark:bg-blue-900/20 text-blue-700 dark:text-blue-400'
                          : 'bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400'
                      "
                    >
                      {{ t.role }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </div>
      </div>
    </Teleport>
  </div>
</template>
