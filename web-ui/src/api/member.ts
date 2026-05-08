import { request } from './http'
import type { TenantMember } from '@/types/api'

export function listMembers(tenantId: string) {
  return request<TenantMember[]>('GET', `/tenant/${tenantId}/members`)
}

export function inviteMember(tenantId: string, username: string) {
  return request<void>('POST', `/tenant/${tenantId}/invitations`, { username })
}

export function updateMemberRole(tenantId: string, userId: string, role: string) {
  return request<void>('PATCH', `/tenant/${tenantId}/members/${userId}/role`, { role })
}

export function removeMember(tenantId: string, userId: string) {
  return request<void>('DELETE', `/tenant/${tenantId}/members/${userId}`)
}

export function leaveTenant(tenantId: string) {
  return request<void>('POST', `/tenant/${tenantId}/leave`)
}
