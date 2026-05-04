import http from './http'
import type { ApiResponse } from '@/types/api'
import type {
  GetConfigResponse,
  UpdateLandingHostsRequest,
  IdConfig,
  ShortLinkNotFoundConfig,
} from '@/types/api'

export function getConfig() {
  return http.get<ApiResponse<GetConfigResponse>>('/config').then((r) => r.data.data)
}

export function updateLandingHosts(data: UpdateLandingHostsRequest) {
  return http.patch<ApiResponse>('/config/landing-hosts', data).then((r) => r.data)
}

export function updateIdConfig(data: IdConfig) {
  return http.patch<ApiResponse>('/config/id-length', data).then((r) => r.data)
}

export function updateShortLinkNotFoundConfig(data: ShortLinkNotFoundConfig) {
  return http.patch<ApiResponse>('/config/short-link-404-handling', data).then((r) => r.data)
}
