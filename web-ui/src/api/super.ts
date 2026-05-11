import { request } from './http'
import type {
  SuperUser,
  SuperUserDetail,
  ListSuperUsersResponse,
  CreateSuperUserRequest,
  ResetPasswordRequest,
  SetStatusRequest,
  ListSuperTenantsResponse,
  SuperTenantMember,
  ListShortLinksResponse,
  SuperTenantDomain,
} from '@/types/api'

// Users
export function listSuperUsers(page = 1, pageSize = 20, query = '') {
  const params: Record<string, string | number> = { page, pageSize }
  if (query) params.query = query
  return request<ListSuperUsersResponse>('GET', '/super/users', params)
}

export function getSuperUser(id: string) {
  return request<SuperUserDetail>('GET', `/super/users/${id}`)
}

export function createSuperUser(data: CreateSuperUserRequest) {
  return request<SuperUser>('POST', '/super/users', data)
}

export function resetSuperUserPassword(id: string, data: ResetPasswordRequest) {
  return request<void>('POST', `/super/users/${id}/reset-password`, data)
}

export function setSuperUserStatus(id: string, data: SetStatusRequest) {
  return request<void>('PATCH', `/super/users/${id}/status`, data)
}

// Tenants
export function listSuperTenants(page = 1, pageSize = 20) {
  return request<ListSuperTenantsResponse>('GET', '/super/tenants', { page, pageSize })
}

export function setSuperTenantStatus(id: string, data: SetStatusRequest) {
  return request<void>('PATCH', `/super/tenants/${id}/status`, data)
}

export function listSuperTenantMembers(tenantId: string) {
  return request<SuperTenantMember[]>('GET', `/super/tenants/${tenantId}/members`)
}

export function listSuperTenantShortLinks(tenantId: string, page = 1, pageSize = 20) {
  return request<ListShortLinksResponse>('GET', `/super/tenants/${tenantId}/short-links`, {
    page,
    pageSize,
  })
}

export function listSuperTenantDomains(tenantId: string) {
  return request<SuperTenantDomain[]>('GET', `/super/tenants/${tenantId}/domains`)
}
