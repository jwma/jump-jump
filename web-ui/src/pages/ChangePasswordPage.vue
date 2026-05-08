<script setup lang="ts">
import { ref } from 'vue'
import { changePassword } from '@/api/user'
import { useToast } from '@/composables/useToast'

defineOptions({ name: 'ChangePasswordPage' })
const toast = useToast()
const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const error = ref('')

async function handleSubmit() {
  error.value = ''

  if (newPassword.value !== confirmPassword.value) {
    error.value = 'New passwords do not match.'
    return
  }

  if (newPassword.value.length < 6) {
    error.value = 'New password must be at least 6 characters.'
    return
  }

  loading.value = true
  try {
    await changePassword({ password: currentPassword.value, newPassword: newPassword.value })
    toast.success('Password changed successfully')
    currentPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
  } catch {
    error.value = 'Failed to change password. Please check your current password.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div>
    <h1 class="text-xl font-semibold text-gray-900">Change Password</h1>
    <p class="mt-1 text-sm text-gray-500">Update your account password.</p>

    <form class="mt-6 max-w-md rounded-lg border bg-white p-6" @submit.prevent="handleSubmit">
      <div class="space-y-4">
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700">Current Password</label>
          <input
            v-model="currentPassword"
            type="password"
            required
            class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
          />
        </div>
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700">New Password</label>
          <input
            v-model="newPassword"
            type="password"
            required
            minlength="6"
            class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
          />
        </div>
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700">Confirm New Password</label>
          <input
            v-model="confirmPassword"
            type="password"
            required
            minlength="6"
            class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
          />
        </div>
      </div>

      <p v-if="error" class="mt-4 text-sm text-red-600">{{ error }}</p>

      <button
        type="submit"
        :disabled="loading"
        class="mt-4 rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:opacity-50"
      >
        {{ loading ? 'Changing...' : 'Change Password' }}
      </button>
    </form>
  </div>
</template>
