import request from './request'

export interface CreateTenantParams {
  name: string
  slug: string
}

export interface Tenant {
  id: string
  name: string
  slug: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface TenantMember {
  userId: string
  username: string
  role: string
  joinedAt: string
}

export function createTenant(params: CreateTenantParams) {
  return request.post<any, { data: Tenant }>('/tenant/', params)
}

export function listTenants() {
  return request.get<any, { data: Tenant[] }>('/tenant/')
}

export function getTenant(id: string) {
  return request.get<any, { data: Tenant }>(`/tenant/${id}`)
}

export function updateTenant(id: string, params: Partial<CreateTenantParams>) {
  return request.patch<any, { data: Tenant }>(`/tenant/${id}`, params)
}

export function listTenantMembers(id: string) {
  return request.get<any, { data: TenantMember[] }>(`/tenant/${id}/members`)
}

export function inviteMember(tenantId: string, username: string) {
  return request.post(`/tenant/${tenantId}/invitations`, { username })
}

export function updateMemberRole(tenantId: string, userId: string, role: string) {
  return request.patch(`/tenant/${tenantId}/members/${userId}/role`, { role })
}

export function removeMember(tenantId: string, userId: string) {
  return request.delete(`/tenant/${tenantId}/members/${userId}`)
}

export function leaveTenant(tenantId: string) {
  return request.post(`/tenant/${tenantId}/leave`)
}
