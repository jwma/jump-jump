<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { listMembers, inviteMember, updateMemberRole, removeMember, leaveTenant } from '@/api/member'
import { useToast } from '@/composables/useToast'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import type { TenantMember } from '@/types/api'
import { UserRole } from '@/types/api'
import {
  Users,
  UserPlus,
  Loader2,
  Shield,
  User,
  ArrowLeftRight,
  Trash2,
  X,
  Search,
} from 'lucide-vue-next'

defineOptions({ name: 'MembersPage' })

const router = useRouter()
const { t } = useI18n()
const auth = useAuthStore()
const toast = useToast()

const members = ref<TenantMember[]>([])
const loading = ref(false)
const searchQuery = ref('')

const filteredMembers = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return members.value
  return members.value.filter((m) => m.username.toLowerCase().includes(q))
})

// Invite modal
const showInvite = ref(false)
const inviteUsername = ref('')
const inviteLoading = ref(false)
const inviteError = ref('')

// Role change confirm
const confirmRoleChange = ref<{ userId: string; username: string; currentRole: string; newRole: string } | null>(null)
const roleLoading = ref(false)
const roleError = ref('')

// Remove member confirm
const confirmRemove = ref<{ userId: string; username: string } | null>(null)
const removeLoading = ref(false)
const removeError = ref('')

// Leave tenant confirm
const confirmLeave = ref(false)
const leaveLoading = ref(false)
const leaveError = ref('')

// Invite modal ESC handling
function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    showInvite.value = false
  }
}

watch(showInvite, (isOpen) => {
  if (isOpen) {
    document.addEventListener('keydown', handleKeydown)
  } else {
    document.removeEventListener('keydown', handleKeydown)
  }
})

onBeforeUnmount(() => {
  document.removeEventListener('keydown', handleKeydown)
})

async function fetchMembers() {
  if (!auth.currentTenantId) return
  loading.value = true
  try {
    members.value = (await listMembers(auth.currentTenantId)) ?? []
  } catch {
    members.value = []
  } finally {
    loading.value = false
  }
}

function openInviteModal() {
  inviteUsername.value = ''
  inviteError.value = ''
  showInvite.value = true
  nextTick(() => {
    const input = document.getElementById('invite-username-input') as HTMLInputElement
    input?.focus()
  })
}

async function handleInvite() {
  if (!auth.currentTenantId || !inviteUsername.value.trim()) return
  inviteLoading.value = true
  inviteError.value = ''
  try {
    await inviteMember(auth.currentTenantId, inviteUsername.value.trim())
    toast.success(t('members.invitationSent', { username: inviteUsername.value.trim() }))
    inviteUsername.value = ''
    showInvite.value = false
    await fetchMembers()
  } catch (e: unknown) {
    inviteError.value = (e as Error).message || t('members.failedToInvite')
  } finally {
    inviteLoading.value = false
  }
}

function requestRoleChange(member: TenantMember) {
  const newRole = member.role === UserRole.Admin ? UserRole.Member : UserRole.Admin
  roleError.value = ''
  confirmRoleChange.value = {
    userId: member.userId,
    username: member.username,
    currentRole: member.role,
    newRole,
  }
}

async function executeRoleChange() {
  if (!auth.currentTenantId || !confirmRoleChange.value) return
  roleLoading.value = true
  roleError.value = ''
  try {
    await updateMemberRole(auth.currentTenantId, confirmRoleChange.value.userId, confirmRoleChange.value.newRole)
    confirmRoleChange.value = null
    await fetchMembers()
  } catch (e: unknown) {
    roleError.value = (e as Error).message || t('members.failedToUpdateRole')
  } finally {
    roleLoading.value = false
  }
}

function requestRemoveMember(member: TenantMember) {
  removeError.value = ''
  confirmRemove.value = { userId: member.userId, username: member.username }
}

async function executeRemoveMember() {
  if (!auth.currentTenantId || !confirmRemove.value) return
  removeLoading.value = true
  removeError.value = ''
  try {
    await removeMember(auth.currentTenantId, confirmRemove.value.userId)
    confirmRemove.value = null
    await fetchMembers()
  } catch (e: unknown) {
    removeError.value = (e as Error).message || t('members.failedToRemoveMember')
  } finally {
    removeLoading.value = false
  }
}

function requestLeaveTenant() {
  leaveError.value = ''
  confirmLeave.value = true
}

async function executeLeaveTenant() {
  if (!auth.currentTenantId) return
  leaveLoading.value = true
  leaveError.value = ''
  try {
    await leaveTenant(auth.currentTenantId)
    auth.clearTenant()
    await auth.fetchAuthInfo()
    router.push({ name: 'select-tenant' })
  } catch (e: unknown) {
    leaveError.value = (e as Error).message || t('members.failedToLeave')
  } finally {
    leaveLoading.value = false
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

onMounted(fetchMembers)
</script>

<template>
  <div>
    <!-- Header -->
    <div class="mb-6 flex items-center justify-between">
      <div class="flex items-center gap-2">
        <Users class="h-5 w-5 text-gray-500" />
        <h1 class="text-xl font-semibold text-gray-900">{{ t('members.title') }}</h1>
        <span v-if="!loading" class="rounded-full bg-gray-100 px-2 py-0.5 text-xs font-medium text-gray-600">
          {{ members.length }}
        </span>
      </div>
      <div class="flex items-center gap-2">
        <!-- Search -->
        <div class="relative">
          <Search class="absolute left-2.5 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="t('members.searchPlaceholder')"
            class="w-48 rounded-md border border-gray-300 py-2 pl-8 pr-3 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
          />
        </div>
        <button
          v-if="auth.isAdmin"
          class="inline-flex items-center gap-1.5 rounded-md bg-blue-600 px-3 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700"
          @click="openInviteModal"
        >
          <UserPlus class="h-4 w-4" />
          {{ t('members.inviteMember') }}
        </button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-20">
      <Loader2 class="h-6 w-6 animate-spin text-gray-400" />
    </div>

    <!-- Members table -->
    <div v-else class="overflow-hidden rounded-lg border bg-white">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b bg-gray-50 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
            <th class="px-4 py-3">{{ t('members.user') }}</th>
            <th class="px-4 py-3">{{ t('members.role') }}</th>
            <th class="hidden px-4 py-3 sm:table-cell">{{ t('members.joined') }}</th>
            <th v-if="auth.isAdmin" class="px-4 py-3 text-right">{{ t('members.actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y">
          <tr v-for="m in filteredMembers" :key="m.userId" class="transition-colors hover:bg-gray-50">
            <td class="px-4 py-3">
              <div class="flex items-center gap-3">
                <div
                  class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-xs font-bold"
                  :class="m.role === UserRole.Admin ? 'bg-blue-100 text-blue-700' : 'bg-gray-100 text-gray-600'"
                >
                  {{ m.username.charAt(0).toUpperCase() }}
                </div>
                <div>
                  <span class="font-medium text-gray-900">{{ m.username }}</span>
                  <span v-if="m.userId === auth.authUser?.id" class="ml-1 text-xs text-gray-400">{{ t('members.you') }}</span>
                </div>
              </div>
            </td>
            <td class="px-4 py-3">
              <span
                class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium"
                :class="m.role === UserRole.Admin ? 'bg-blue-50 text-blue-700' : 'bg-gray-100 text-gray-600'"
              >
                <Shield v-if="m.role === UserRole.Admin" class="h-3 w-3" />
                <User v-else class="h-3 w-3" />
                {{ m.role }}
              </span>
            </td>
            <td class="hidden whitespace-nowrap px-4 py-3 text-gray-500 sm:table-cell">
              {{ formatDate(m.joinedAt) }}
            </td>
            <td v-if="auth.isAdmin" class="px-4 py-3">
              <div v-if="m.userId !== auth.authUser?.id" class="flex items-center justify-end gap-1">
                <button
                  class="rounded p-1.5 text-gray-400 transition-colors hover:bg-blue-50 hover:text-blue-600"
                  :title="t('members.changeRole')"
                  @click="requestRoleChange(m)"
                >
                  <ArrowLeftRight class="h-4 w-4" />
                </button>
                <button
                  class="rounded p-1.5 text-gray-400 transition-colors hover:bg-red-50 hover:text-red-600"
                  :title="t('members.removeMember')"
                  @click="requestRemoveMember(m)"
                >
                  <Trash2 class="h-4 w-4" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-if="filteredMembers.length === 0 && members.length > 0" class="py-12 text-center text-sm text-gray-400">
        {{ t('members.noMatch', { query: searchQuery }) }}
      </div>
      <div v-else-if="members.length === 0" class="py-12 text-center">
        <Users class="mx-auto h-10 w-10 text-gray-300" />
        <p class="mt-3 text-sm font-medium text-gray-500">{{ t('members.noMembersYet') }}</p>
        <p class="mt-1 text-sm text-gray-400">{{ t('members.inviteToCollaborate') }}</p>
        <button
          v-if="auth.isAdmin"
          class="mt-3 inline-flex items-center gap-1.5 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700"
          @click="openInviteModal"
        >
          <UserPlus class="h-4 w-4" />
          {{ t('members.inviteMember') }}
        </button>
      </div>
    </div>

    <!-- Leave tenant (subtle placement) -->
    <div v-if="!loading" class="mt-6 text-right">
      <button
        class="text-sm text-gray-400 transition-colors hover:text-gray-600"
        @click="requestLeaveTenant"
      >
        {{ t('members.leaveTenant') }}
      </button>
    </div>

    <!-- Invite Modal -->
    <Teleport to="body">
      <div
        v-if="showInvite"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
        @click.self="showInvite = false"
      >
        <div class="mx-4 w-full max-w-sm rounded-lg bg-white p-6 shadow-xl">
          <div class="mb-4 flex items-center justify-between">
            <h3 class="text-lg font-semibold text-gray-900">{{ t('members.inviteMember') }}</h3>
            <button class="rounded p-1 text-gray-400 hover:text-gray-600" @click="showInvite = false">
              <X class="h-4 w-4" />
            </button>
          </div>

          <form @submit.prevent="handleInvite">
            <div class="mb-3">
              <label for="invite-username-input" class="mb-1 block text-sm font-medium text-gray-700">{{ t('members.username') }}</label>
              <input
                id="invite-username-input"
                v-model="inviteUsername"
                type="text"
                required
                :placeholder="t('members.invitePlaceholder')"
                class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
              />
              <p class="mt-1 text-xs text-gray-400">{{ t('members.inviteHint') }}</p>
            </div>

            <p v-if="inviteError" class="mb-3 text-sm text-red-600">{{ inviteError }}</p>

            <div class="flex justify-end gap-2">
              <button
                type="button"
                class="rounded-md border px-3 py-1.5 text-sm text-gray-600 hover:bg-gray-50"
                @click="showInvite = false"
              >
                {{ t('common.cancel') }}
              </button>
              <button
                type="submit"
                :disabled="inviteLoading || !inviteUsername.trim()"
                class="rounded-md bg-blue-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
              >
                {{ inviteLoading ? t('members.sending') : t('members.sendInvite') }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Role Change Confirm -->
    <ConfirmDialog
      :open="!!confirmRoleChange"
      :title="t('members.changeRole')"
      :confirm-text="roleLoading ? t('members.changing') : t('members.confirm')"
      :loading="roleLoading"
      @confirm="executeRoleChange"
      @cancel="confirmRoleChange = null"
    >
      <p class="mt-2 text-sm text-gray-600">
        {{ t('members.changeRoleQuestion', {
          username: confirmRoleChange?.username,
          currentRole: confirmRoleChange?.currentRole,
          newRole: confirmRoleChange?.newRole
        }) }}
      </p>
      <div class="mt-3 rounded-lg border bg-gray-50 p-3 text-xs text-gray-600">
        <p class="font-medium text-gray-700">
          {{ confirmRoleChange?.newRole === UserRole.Admin ? t('members.admin') : t('members.member') }} {{ t('members.permissions') }}:
        </p>
        <ul v-if="confirmRoleChange?.newRole === UserRole.Admin" class="mt-1 list-inside list-disc space-y-0.5 text-gray-500">
          <li>{{ t('members.adminPermSettings') }}</li>
          <li>{{ t('members.adminPermMembers') }}</li>
          <li>{{ t('members.adminPermRoles') }}</li>
          <li>{{ t('members.adminPermLinks') }}</li>
        </ul>
        <ul v-else class="mt-1 list-inside list-disc space-y-0.5 text-gray-500">
          <li>{{ t('members.memberPermLinks') }}</li>
          <li>{{ t('members.memberPermView') }}</li>
        </ul>
      </div>
      <p v-if="roleError" class="mt-2 text-sm text-red-600">{{ roleError }}</p>
    </ConfirmDialog>

    <!-- Remove Member Confirm -->
    <ConfirmDialog
      :open="!!confirmRemove"
      :title="t('members.removeMember')"
      variant="danger"
      :confirm-text="removeLoading ? t('members.removing') : t('members.remove')"
      :loading="removeLoading"
      @confirm="executeRemoveMember"
      @cancel="confirmRemove = null"
    >
      <p class="mt-2 text-sm text-gray-600">
        {{ t('members.removeConfirm', { username: confirmRemove?.username }) }}
      </p>
      <p v-if="removeError" class="mt-2 text-sm text-red-600">{{ removeError }}</p>
    </ConfirmDialog>

    <!-- Leave Tenant Confirm -->
    <ConfirmDialog
      :open="confirmLeave"
      :title="t('members.leaveTenant')"
      variant="danger"
      :confirm-text="leaveLoading ? t('members.leaving') : t('members.leave')"
      :loading="leaveLoading"
      @confirm="executeLeaveTenant"
      @cancel="confirmLeave = false"
    >
      <p class="mt-2 text-sm text-gray-600">
        {{ t('members.leaveConfirm', { tenantName: auth.currentTenant?.tenantName }) }}
      </p>
      <p v-if="leaveError" class="mt-2 text-sm text-red-600">{{ leaveError }}</p>
    </ConfirmDialog>
  </div>
</template>
