import { request } from './http'
import type { GetTenantConfigResponse } from '@/types/api'

export interface UpdateTenantConfigRequest {
  idLength?: number
  idMinimumLength?: number
  idMaximumLength?: number
  notFoundMode?: string
  notFoundValue?: string
}

export function getTenantConfig(tenantId: string) {
  return request<GetTenantConfigResponse>('GET', `/tenant/${tenantId}/config`)
}

export function updateTenantConfig(tenantId: string, data: UpdateTenantConfigRequest) {
  return request<GetTenantConfigResponse>('PATCH', `/tenant/${tenantId}/config`, data)
}
