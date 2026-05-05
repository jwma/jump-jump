import http from './http'
import type { ApiResponse } from '@/types/api'
import type {
  GetConfigResponse,
  UpdateLandingHostsRequest,
  IdConfig,
  ShortLinkNotFoundConfig,
} from '@/types/api'

function checkCode<T>(res: { data: ApiResponse<T> }): T {
  if (res.data.code !== 0) {
    return Promise.reject(new Error(res.data.msg || 'Request failed')) as never
  }
  return res.data.data
}

export function getConfig() {
  return http.get<ApiResponse<GetConfigResponse>>('/config').then(checkCode)
}

export function updateLandingHosts(data: UpdateLandingHostsRequest) {
  return http.patch<ApiResponse>('/config/landing-hosts', data).then(checkCode)
}

export function updateIdConfig(data: IdConfig) {
  return http.patch<ApiResponse>('/config/id-length', data).then(checkCode)
}

export function updateShortLinkNotFoundConfig(data: ShortLinkNotFoundConfig) {
  return http.patch<ApiResponse>('/config/short-link-404-handling', data).then(checkCode)
}
