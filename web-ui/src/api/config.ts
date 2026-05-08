import { request } from './http'
import type { GetConfigResponse } from '@/types/api'

export interface UpdateTenantConfigRequest {
  idLength?: number
  idMinimumLength?: number
  idMaximumLength?: number
  notFoundMode?: string
  notFoundValue?: string
}

export function getTenantConfig(tenantId: string) {
  return request<GetConfigResponse>('GET', `/tenant/${tenantId}/config`)
}

export function updateTenantConfig(tenantId: string, data: UpdateTenantConfigRequest) {
  return request<GetConfigResponse>('PATCH', `/tenant/${tenantId}/config`, data)
}
