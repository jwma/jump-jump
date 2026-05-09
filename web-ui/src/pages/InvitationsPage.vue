<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { listInvitations, acceptInvitation, rejectInvitation } from '@/api/invitation'
import type { Invitation } from '@/types/api'
import { ArrowLeft, Check, X } from 'lucide-vue-next'

const router = useRouter()
const { t } = useI18n()
const auth = useAuthStore()
const loading = ref(false)
const invitations = ref<Invitation[]>([])
const actionLoading = ref<Record<string, boolean>>({})

onMounted(async () => {
  loading.value = true
  try {
    invitations.value = (await listInvitations()) ?? []
  } catch {
    // handled by interceptor
  } finally {
    loading.value = false
  }
})

function formatTime(dateStr: string) {
  return new Date(dateStr).toLocaleString()
}

async function handleAccept(id: string) {
  actionLoading.value[id] = true
  try {
    await acceptInvitation(id)
    invitations.value = invitations.value.filter((i) => i.id !== id)
    await auth.fetchAuthInfo()
    await auth.fetchInvitations()
  } catch {
    // handled by interceptor
  } finally {
    actionLoading.value[id] = false
  }
}

async function handleReject(id: string) {
  actionLoading.value[id] = true
  try {
    await rejectInvitation(id)
    invitations.value = invitations.value.filter((i) => i.id !== id)
    await auth.fetchInvitations()
  } catch {
    // handled by interceptor
  } finally {
    actionLoading.value[id] = false
  }
}
</script>

<template>
  <div>
    <div class="mb-6 flex items-center gap-3">
      <button class="rounded p-1 text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300" @click="router.push({ name: 'dashboard' })">
        <ArrowLeft class="h-5 w-5" />
      </button>
      <h1 class="text-xl font-semibold text-gray-900 dark:text-white">{{ t('invitations.title') }}</h1>
    </div>

    <div v-if="loading" class="py-12 text-center text-gray-400 dark:text-gray-500">{{ t('common.loading') }}</div>

    <div v-else-if="invitations.length === 0" class="py-12 text-center text-gray-400 dark:text-gray-500">
      {{ t('invitations.empty') }}
    </div>

    <div v-else class="space-y-3">
      <div
        v-for="inv in invitations"
        :key="inv.id"
        class="flex items-center justify-between rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 p-4 shadow-sm"
      >
        <div>
          <div class="font-medium text-gray-900 dark:text-white">{{ inv.tenantName }}</div>
          <div class="text-sm text-gray-500 dark:text-gray-400">
            {{ t('invitations.invitedBy', { username: inv.inviterUsername }) }} &middot; {{ formatTime(inv.createdAt) }}
          </div>
        </div>
        <div class="flex items-center gap-2">
          <template v-if="inv.status === 'pending'">
            <button
              :disabled="actionLoading[inv.id]"
              class="flex items-center gap-1 rounded-md bg-blue-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
              @click="handleAccept(inv.id)"
            >
              <Check class="h-4 w-4" />
              {{ t('invitations.accept') }}
            </button>
            <button
              :disabled="actionLoading[inv.id]"
              class="flex items-center gap-1 rounded-md border border-gray-200 dark:border-gray-700 px-3 py-1.5 text-sm font-medium text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50"
              @click="handleReject(inv.id)"
            >
              <X class="h-4 w-4" />
              {{ t('invitations.reject') }}
            </button>
          </template>
          <span
            v-else
            class="rounded-full px-2.5 py-0.5 text-xs font-medium"
            :class="inv.status === 'accepted' ? 'bg-green-50 text-green-700' : 'bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400'"
          >
            {{ inv.status === 'accepted' ? t('invitations.accepted') : t('invitations.rejected') }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>
