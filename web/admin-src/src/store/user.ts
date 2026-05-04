import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as loginApi, getUserInfo as getUserInfoApi, logout as logoutApi, type LoginParams, type UserInfo } from '@/api/user'
import { getToken, setToken, removeToken } from '@/utils/auth'

export const useUserStore = defineStore('user', () => {
  const token = ref<string | null>(getToken())
  const userInfo = ref<UserInfo | null>(null)

  async function login(params: LoginParams) {
    const res = await loginApi(params)
    token.value = res.data.token
    setToken(res.data.token)
  }

  async function fetchUserInfo() {
    const res = await getUserInfoApi()
    userInfo.value = res.data
    return res.data
  }

  async function logout() {
    try {
      await logoutApi()
    } finally {
      token.value = null
      userInfo.value = null
      removeToken()
    }
  }

  function resetState() {
    token.value = null
    userInfo.value = null
    removeToken()
  }

  return { token, userInfo, login, fetchUserInfo, logout, resetState }
})
