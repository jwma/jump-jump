const TOKEN_KEY = 'jj_admin_token'
const TENANT_ID_KEY = 'jj_admin_tenant_id'
const REMEMBER_KEY = 'jj_admin_remember_username'

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string) {
  localStorage.setItem(TOKEN_KEY, token)
}

export function removeToken() {
  localStorage.removeItem(TOKEN_KEY)
}

export function getTenantId(): string | null {
  return localStorage.getItem(TENANT_ID_KEY)
}

export function setTenantId(tenantId: string) {
  localStorage.setItem(TENANT_ID_KEY, tenantId)
}

export function removeTenantId() {
  localStorage.removeItem(TENANT_ID_KEY)
}

export function getRememberedUsername(): string {
  return localStorage.getItem(REMEMBER_KEY) || ''
}

export function setRememberedUsername(username: string) {
  localStorage.setItem(REMEMBER_KEY, username)
}

export function removeRememberedUsername() {
  localStorage.removeItem(REMEMBER_KEY)
}
