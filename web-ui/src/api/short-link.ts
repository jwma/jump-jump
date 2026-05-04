import { request } from './http'
import type {
  CreateShortLinkRequest,
  UpdateShortLinkRequest,
  ListShortLinksResponse,
  GetShortLinkResponse,
  CreateShortLinkResponse,
  UpdateShortLinkResponse,
  ShortLinkDataResponse,
} from '@/types/api'

export function listShortLinks(page?: number, pageSize?: number) {
  return request<ListShortLinksResponse>('GET', '/short-link/', { page, pageSize })
}

export function getShortLink(id: string) {
  return request<GetShortLinkResponse>('GET', `/short-link/${id}`)
}

export function createShortLink(data: CreateShortLinkRequest) {
  return request<CreateShortLinkResponse>('POST', '/short-link/', data)
}

export function updateShortLink(id: string, data: UpdateShortLinkRequest) {
  return request<UpdateShortLinkResponse>('PATCH', `/short-link/${id}`, data)
}

export function deleteShortLink(id: string) {
  return request<void>('DELETE', `/short-link/${id}`)
}

export function getShortLinkData(id: string, startDate: string, endDate: string) {
  return request<ShortLinkDataResponse>('GET', `/short-link/${id}/data`, { startDate, endDate })
}
