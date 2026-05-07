import request from './request'

export interface Invitation {
  id: string
  tenantId: string
  tenantName: string
  inviterUsername: string
  status: string
  createdAt: string
}

export function listInvitations() {
  return request.get<any, { data: Invitation[] }>('/invitations/')
}

export function acceptInvitation(id: string) {
  return request.post(`/invitations/${id}/accept`)
}

export function rejectInvitation(id: string) {
  return request.post(`/invitations/${id}/reject`)
}
