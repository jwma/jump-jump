import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as apiLogin, getUserInfo, logout as apiLogout } from '@/api/user'
import type { UserInfo } from '@/types/api'
import { UserRole } from '@/types/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const user = ref<UserInfo | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === UserRole.Admin)
  const username = computed(() => user.value?.username || '')

  async function login(username: string, password: string) {
    const data = await apiLogin({ username, password })
    token.value = data.token
    localStorage.setItem('token', data.token)
    await fetchUser()
  }

  async function fetchUser() {
    user.value = await getUserInfo()
  }

  function clearAuth() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
  }

  async function logout() {
    try {
      await apiLogout()
    } finally {
      clearAuth()
    }
  }

  return { token, user, isLoggedIn, isAdmin, username, login, fetchUser, clearAuth, logout }
})
