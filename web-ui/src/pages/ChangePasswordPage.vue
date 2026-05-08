<script setup lang="ts">
import { ref, computed } from 'vue'
import { changePassword } from '@/api/user'
import { useToast } from '@/composables/useToast'
import { Eye, EyeOff, CheckCircle2 } from 'lucide-vue-next'

defineOptions({ name: 'ChangePasswordPage' })
const toast = useToast()
const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const error = ref('')
const success = ref(false)
const showCurrentPassword = ref(false)
const showNewPassword = ref(false)
const showConfirmPassword = ref(false)

const passwordStrength = computed(() => {
  const pw = newPassword.value
  if (!pw) return { score: 0, label: '', color: '' }

  let score = 0
  if (pw.length >= 6) score++
  if (pw.length >= 10) score++
  if (/[a-z]/.test(pw) && /[A-Z]/.test(pw)) score++
  if (/\d/.test(pw)) score++
  if (/[^a-zA-Z0-9]/.test(pw)) score++

  const levels = [
    { label: 'Very weak', color: 'bg-red-500' },
    { label: 'Weak', color: 'bg-orange-500' },
    { label: 'Fair', color: 'bg-yellow-500' },
    { label: 'Good', color: 'bg-blue-500' },
    { label: 'Strong', color: 'bg-green-500' },
  ]

  const level = levels[Math.min(score, levels.length) - 1] || levels[0]
  return { score, ...level }
})

async function handleSubmit() {
  error.value = ''
  success.value = false

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
    toast.success('Password changed successfully!')
    success.value = true
    currentPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
    setTimeout(() => (success.value = false), 4000)
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
          <div class="relative">
            <input
              v-model="currentPassword"
              :type="showCurrentPassword ? 'text' : 'password'"
              required
              class="w-full rounded-md border border-gray-300 px-3 py-2 pr-10 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
            />
            <button
              type="button"
              class="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400 hover:text-gray-600"
              @click="showCurrentPassword = !showCurrentPassword"
            >
              <Eye v-if="!showCurrentPassword" class="h-4 w-4" />
              <EyeOff v-else class="h-4 w-4" />
            </button>
          </div>
        </div>
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700">New Password</label>
          <div class="relative">
            <input
              v-model="newPassword"
              :type="showNewPassword ? 'text' : 'password'"
              required
              minlength="6"
              class="w-full rounded-md border border-gray-300 px-3 py-2 pr-10 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
            />
            <button
              type="button"
              class="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400 hover:text-gray-600"
              @click="showNewPassword = !showNewPassword"
            >
              <Eye v-if="!showNewPassword" class="h-4 w-4" />
              <EyeOff v-else class="h-4 w-4" />
            </button>
          </div>
          <div v-if="newPassword" class="mt-2">
            <div class="flex gap-1">
              <div
                v-for="i in 5"
                :key="i"
                class="h-1.5 flex-1 rounded-full"
                :class="i <= passwordStrength.score ? passwordStrength.color : 'bg-gray-200'"
              />
            </div>
            <p class="mt-1 text-xs text-gray-500">
              Strength: <span class="font-medium">{{ passwordStrength.label }}</span>
            </p>
          </div>
        </div>
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700">Confirm New Password</label>
          <div class="relative">
            <input
              v-model="confirmPassword"
              :type="showConfirmPassword ? 'text' : 'password'"
              required
              minlength="6"
              class="w-full rounded-md border border-gray-300 px-3 py-2 pr-10 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
              :class="confirmPassword && newPassword !== confirmPassword ? 'border-red-300' : ''"
            />
            <button
              type="button"
              class="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400 hover:text-gray-600"
              @click="showConfirmPassword = !showConfirmPassword"
            >
              <Eye v-if="!showConfirmPassword" class="h-4 w-4" />
              <EyeOff v-else class="h-4 w-4" />
            </button>
          </div>
          <p
            v-if="confirmPassword && newPassword !== confirmPassword"
            class="mt-1 text-xs text-red-500"
          >
            Passwords do not match
          </p>
        </div>
      </div>

      <p v-if="error" class="mt-4 text-sm text-red-600">{{ error }}</p>

      <div
        v-if="success"
        class="mt-4 flex items-center gap-2 rounded-md bg-green-50 px-3 py-2 text-sm text-green-700"
      >
        <CheckCircle2 class="h-4 w-4" />
        Password changed successfully!
      </div>

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
