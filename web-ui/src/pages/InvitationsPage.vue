<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { listInvitations, acceptInvitation, rejectInvitation } from '@/api/invitation'
import type { Invitation } from '@/types/api'
import { ArrowLeft, Check, X } from 'lucide-vue-next'

const router = useRouter()
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
      <button class="rounded p-1 text-gray-400 hover:text-gray-600" @click="router.push({ name: 'dashboard' })">
        <ArrowLeft class="h-5 w-5" />
      </button>
      <h1 class="text-xl font-semibold text-gray-900">Invitations</h1>
    </div>

    <div v-if="loading" class="py-12 text-center text-gray-400">Loading...</div>

    <div v-else-if="invitations.length === 0" class="py-12 text-center text-gray-400">
      No pending invitations.
    </div>

    <div v-else class="space-y-3">
      <div
        v-for="inv in invitations"
        :key="inv.id"
        class="flex items-center justify-between rounded-lg border bg-white p-4 shadow-sm"
      >
        <div>
          <div class="font-medium text-gray-900">{{ inv.tenantName }}</div>
          <div class="text-sm text-gray-500">
            Invited by {{ inv.inviterUsername }} &middot; {{ formatTime(inv.createdAt) }}
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
              Accept
            </button>
            <button
              :disabled="actionLoading[inv.id]"
              class="flex items-center gap-1 rounded-md border px-3 py-1.5 text-sm font-medium text-gray-600 hover:bg-gray-50 disabled:opacity-50"
              @click="handleReject(inv.id)"
            >
              <X class="h-4 w-4" />
              Reject
            </button>
          </template>
          <span
            v-else
            class="rounded-full px-2.5 py-0.5 text-xs font-medium"
            :class="inv.status === 'accepted' ? 'bg-green-50 text-green-700' : 'bg-gray-100 text-gray-500'"
          >
            {{ inv.status === 'accepted' ? 'Accepted' : 'Rejected' }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>
