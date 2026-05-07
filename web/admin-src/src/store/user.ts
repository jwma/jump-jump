import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  login as loginApi,
  getAuthInfo as getAuthInfoApi,
  getUserInfo as getUserInfoApi,
  logout as logoutApi,
  type LoginParams,
  type AuthInfoUser,
  type AuthInfoTenant,
  type UserInfo,
} from '@/api/user'
import { listInvitations, type Invitation } from '@/api/invitation'
import { getToken, setToken, removeToken, getTenantId, setTenantId, removeTenantId } from '@/utils/auth'

export const useUserStore = defineStore('user', () => {
  const token = ref<string | null>(getToken())
  const authUser = ref<AuthInfoUser | null>(null)
  const tenants = ref<AuthInfoTenant[]>([])
  const currentTenantId = ref<string | null>(getTenantId())
  const userInfo = ref<UserInfo | null>(null)
  const invitations = ref<Invitation[]>([])

  const isSuper = computed(() => authUser.value?.isSuper ?? false)
  const currentTenant = computed(() =>
    tenants.value.find((t) => t.tenantId === currentTenantId.value) ?? null,
  )
  const pendingInvitationCount = computed(() =>
    invitations.value.filter((i) => i.status === 'pending').length,
  )

  async function login(params: LoginParams) {
    const res = await loginApi(params)
    token.value = res.data.token
    setToken(res.data.token)
  }

  async function fetchAuthInfo() {
    const res = await getAuthInfoApi()
    authUser.value = res.data.user
    tenants.value = res.data.tenants
    return res.data
  }

  async function fetchUserInfo() {
    const res = await getUserInfoApi()
    userInfo.value = res.data
    return res.data
  }

  async function fetchInvitations() {
    const res = await listInvitations()
    invitations.value = res.data ?? []
  }

  function selectTenant(tenantId: string) {
    currentTenantId.value = tenantId
    setTenantId(tenantId)
  }

  function clearTenant() {
    currentTenantId.value = null
    removeTenantId()
  }

  async function logout() {
    try {
      await logoutApi()
    } finally {
      resetState()
    }
  }

  function resetState() {
    token.value = null
    authUser.value = null
    tenants.value = []
    currentTenantId.value = null
    userInfo.value = null
    invitations.value = []
    removeToken()
    removeTenantId()
  }

  return {
    token,
    authUser,
    tenants,
    currentTenantId,
    currentTenant,
    userInfo,
    invitations,
    isSuper,
    pendingInvitationCount,
    login,
    fetchAuthInfo,
    fetchUserInfo,
    fetchInvitations,
    selectTenant,
    clearTenant,
    logout,
    resetState,
  }
})
