import axios from 'axios'
import type { ApiResponse } from '@/types/api'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

export class ApiError extends Error {
  status: number
  code: number

  constructor(message: string, status: number = 0, code: number = 0) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
  }
}

const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/v1',
  timeout: 10000,
})

http.interceptors.request.use((config) => {
  const auth = useAuthStore()
  if (auth.token) {
    config.headers.Authorization = `Bearer ${auth.token}`
  }
  if (auth.currentTenantId && !config.url?.startsWith('/super/')) {
    config.headers['X-Tenant-ID'] = auth.currentTenantId
  }
  return config
})

http.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      const auth = useAuthStore()
      auth.clearAuth()
      router.push({ name: 'login' })
    }
    return Promise.reject(new ApiError(
      error.response?.data?.msg || error.message || 'Request failed',
      error.response?.status || 0,
      error.response?.data?.code || 0,
    ))
  },
)

export async function request<T>(method: string, url: string, data?: unknown): Promise<T> {
  const res = await http.request<ApiResponse<T>>({
    method,
    url,
    data,
    params: method === 'GET' ? data : undefined,
  })
  if (res.data.code !== 0) {
    return Promise.reject(new ApiError(res.data.msg || 'Request failed', res.status, res.data.code))
  }
  return res.data.data
}

export default http
