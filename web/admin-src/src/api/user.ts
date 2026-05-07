import request from './request'

export interface LoginParams {
  username: string
  password: string
}

export interface LoginResult {
  token: string
}

export interface AuthInfoUser {
  id: string
  username: string
  isSuper: boolean
}

export interface AuthInfoTenant {
  tenantId: string
  tenantName: string
  role: string
}

export interface AuthInfo {
  user: AuthInfoUser
  tenants: AuthInfoTenant[]
}

export interface UserInfo {
  username: string
  role: string
  isSuper: boolean
}

export interface ChangePasswordParams {
  password: string
  newPassword: string
}

export function login(params: LoginParams) {
  return request.post<any, { data: LoginResult }>('/auth/login', params)
}

export function getAuthInfo() {
  return request.get<any, { data: AuthInfo }>('/auth/info')
}

export function getUserInfo() {
  return request.get<any, { data: UserInfo }>('/user/info')
}

export function logout() {
  return request.post('/user/logout')
}

export function changePassword(params: ChangePasswordParams) {
  return request.post('/auth/change-password', params)
}
