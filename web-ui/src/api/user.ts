import { request } from './http'
import type { LoginRequest, LoginResponseData, UserInfo, ChangePasswordRequest } from '@/types/api'

export function login(data: LoginRequest) {
  return request<LoginResponseData>('POST', '/user/login', data)
}

export function getUserInfo() {
  return request<UserInfo>('GET', '/user/info')
}

export function logout() {
  return request<void>('POST', '/user/logout')
}

export function changePassword(data: ChangePasswordRequest) {
  return request<void>('POST', '/user/change-password', data)
}
