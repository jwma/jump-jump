import { request } from './http'
import type { Invitation } from '@/types/api'

export function listInvitations() {
  return request<Invitation[]>('GET', '/invitations/')
}

export function acceptInvitation(id: string) {
  return request<void>('POST', `/invitations/${id}/accept`)
}

export function rejectInvitation(id: string) {
  return request<void>('POST', `/invitations/${id}/reject`)
}
