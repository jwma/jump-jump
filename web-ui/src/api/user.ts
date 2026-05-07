import { request } from './http'
import type { LoginRequest, LoginResponseData, UserInfo, AuthInfo, ChangePasswordRequest } from '@/types/api'

export function login(data: LoginRequest) {
  return request<LoginResponseData>('POST', '/auth/login', data)
}

export function getAuthInfo() {
  return request<AuthInfo>('GET', '/auth/info')
}

export function getUserInfo() {
  return request<UserInfo>('GET', '/user/info')
}

export function logout() {
  return request<void>('POST', '/user/logout')
}

export function changePassword(data: ChangePasswordRequest) {
  return request<void>('POST', '/auth/change-password', data)
}
