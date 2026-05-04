import request from './request'

export interface LoginParams {
  username: string
  password: string
}

export interface LoginResult {
  token: string
}

export interface UserInfo {
  username: string
  role: number
}

export interface ChangePasswordParams {
  password: string
  newPassword: string
}

export function login(params: LoginParams) {
  return request.post<any, { data: LoginResult }>('/user/login', params)
}

export function getUserInfo() {
  return request.get<any, { data: UserInfo }>('/user/info')
}

export function logout() {
  return request.post('/user/logout')
}

export function changePassword(params: ChangePasswordParams) {
  return request.post('/user/change-password', params)
}
