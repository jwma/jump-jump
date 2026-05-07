import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as apiLogin, getAuthInfo, getUserInfo, logout as apiLogout } from '@/api/user'
import { listInvitations } from '@/api/invitation'
import type { UserInfo, AuthInfoUser, AuthInfoTenant, Invitation } from '@/types/api'
import { UserRole } from '@/types/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const user = ref<UserInfo | null>(null)
  const authUser = ref<AuthInfoUser | null>(null)
  const tenants = ref<AuthInfoTenant[]>([])
  const currentTenantId = ref<string | null>(localStorage.getItem('tenant_id') || null)
  const invitations = ref<Invitation[]>([])

  const isLoggedIn = computed(() => !!token.value)
  const isSuper = computed(() => authUser.value?.isSuper ?? false)
  const isAdmin = computed(() => user.value?.role === UserRole.Admin)
  const username = computed(() => authUser.value?.username || user.value?.username || '')
  const currentTenant = computed(() =>
    tenants.value.find((t) => t.tenantId === currentTenantId.value) ?? null,
  )
  const pendingInvitationCount = computed(() =>
    invitations.value.filter((i) => i.status === 'pending').length,
  )

  async function login(username: string, password: string) {
    const data = await apiLogin({ username, password })
    token.value = data.token
    localStorage.setItem('token', data.token)
  }

  async function fetchAuthInfo() {
    const data = await getAuthInfo()
    authUser.value = data.user
    tenants.value = data.tenants
    return data
  }

  async function fetchUser() {
    user.value = await getUserInfo()
  }

  async function fetchInvitations() {
    const data = await listInvitations()
    invitations.value = data ?? []
  }

  function selectTenant(tenantId: string) {
    currentTenantId.value = tenantId
    localStorage.setItem('tenant_id', tenantId)
  }

  function clearTenant() {
    currentTenantId.value = null
    localStorage.removeItem('tenant_id')
  }

  function clearAuth() {
    token.value = ''
    user.value = null
    authUser.value = null
    tenants.value = []
    currentTenantId.value = null
    invitations.value = []
    localStorage.removeItem('token')
    localStorage.removeItem('tenant_id')
  }

  async function logout() {
    try {
      await apiLogout()
    } finally {
      clearAuth()
    }
  }

  return {
    token,
    user,
    authUser,
    tenants,
    currentTenantId,
    currentTenant,
    invitations,
    isLoggedIn,
    isSuper,
    isAdmin,
    username,
    pendingInvitationCount,
    login,
    fetchAuthInfo,
    fetchUser,
    fetchInvitations,
    selectTenant,
    clearTenant,
    clearAuth,
    logout,
  }
})
