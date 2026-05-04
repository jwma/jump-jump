import http from './http'
import type { ApiResponse } from '@/types/api'
import type {
  CreateShortLinkRequest,
  UpdateShortLinkRequest,
  ListShortLinksResponse,
  GetShortLinkResponse,
  CreateShortLinkResponse,
  UpdateShortLinkResponse,
  ShortLinkHistoryResponse,
  ShortLinkDataResponse,
} from '@/types/api'

export function listShortLinks(page?: number, pageSize?: number) {
  return http
    .get<ApiResponse<ListShortLinksResponse>>('/short-link/', { params: { page, pageSize } })
    .then((r) => r.data.data)
}

export function getShortLink(id: string) {
  return http.get<ApiResponse<GetShortLinkResponse>>(`/short-link/${id}`).then((r) => r.data.data)
}

export function createShortLink(data: CreateShortLinkRequest) {
  return http
    .post<ApiResponse<CreateShortLinkResponse>>('/short-link/', data)
    .then((r) => r.data.data)
}

export function updateShortLink(id: string, data: UpdateShortLinkRequest) {
  return http
    .patch<ApiResponse<UpdateShortLinkResponse>>(`/short-link/${id}`, data)
    .then((r) => r.data.data)
}

export function deleteShortLink(id: string) {
  return http.delete<ApiResponse>(`/short-link/${id}`).then((r) => r.data)
}

export function getShortLinkHistory(id: string) {
  return http
    .get<ApiResponse<ShortLinkHistoryResponse>>(`/short-link/${id}/history`)
    .then((r) => r.data.data)
}

export function getShortLinkData(id: string, startDate: string, endDate: string) {
  return http
    .get<ApiResponse<ShortLinkDataResponse>>(`/short-link/${id}/data`, {
      params: { startDate, endDate },
    })
    .then((r) => r.data.data)
}
