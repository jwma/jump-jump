export interface ApiResponse<T = unknown> {
  msg: string
  code: number
  data: T
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponseData {
  token: string
}

export interface UserInfo {
  username: string
  role: number
}

export interface ChangePasswordRequest {
  password: string
  newPassword: string
}

export interface ShortLinkData {
  id: string
  url: string
  description: string
  isEnable: boolean
  createdBy: string
  createTime: string
  updateTime: string
}

export interface CreateShortLinkRequest {
  id?: string
  url: string
  description: string
  isEnable: boolean
  idLength?: number
}

export interface UpdateShortLinkRequest {
  url: string
  description: string
  isEnable: boolean
}

export interface GetShortLinkResponse {
  shortLink: ShortLinkData
}

export interface CreateShortLinkResponse {
  shortLink: ShortLinkData
}

export interface UpdateShortLinkResponse {
  shortLink: ShortLinkData
}

export interface ListShortLinksResponse {
  shortLinks: ShortLinkData[]
  total: number
}

export interface RequestHistory {
  id: string
  url: string
  ip: string
  ua: string
  time: string
}

export interface ShortLinkHistoryResponse {
  histories: RequestHistory[]
}

export interface IdConfig {
  idLength: number
  idMinimumLength: number
  idMaximumLength: number
}

export interface ShortLinkNotFoundConfig {
  mode: 'content' | 'redirect'
  value: string
}

export interface SystemConfig {
  landingHosts: string[]
  idConfig: IdConfig
  shortLinkNotFoundConfig: ShortLinkNotFoundConfig
}

export interface GetConfigResponse {
  config: SystemConfig
}

export interface UpdateLandingHostsRequest {
  hosts: string[]
}

export interface DailyReport {
  pv: number
  uv: number
  os: Record<string, number>
}

export interface DailyReportItem {
  date: string
  report: DailyReport
}

export interface DailyStats {
  date: string
  pv: number
  uv: number
}

export interface ShortLinkDataResponse {
  histories: RequestHistory[]
}

export interface Tenant {
  id: string
  name: string
  slug: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface TenantDomain {
  id: string
  tenantId: string
  domain: string
  isDefault: boolean
  createdAt: string
}

export interface CreateTenantRequest {
  name: string
  slug: string
}

export interface AddDomainRequest {
  domain: string
  isDefault: boolean
}

export const UserRole = {
  User: 1,
  Admin: 2,
} as const
