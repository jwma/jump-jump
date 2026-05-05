import { request } from './http'
import type { Tenant, TenantDomain, CreateTenantRequest, AddDomainRequest } from '@/types/api'

export function listTenants() {
  return request<Tenant[]>('GET', '/tenant/')
}

export function getTenant(id: string) {
  return request<Tenant>('GET', `/tenant/${id}`)
}

export function createTenant(data: CreateTenantRequest) {
  return request<Tenant>('POST', '/tenant/', data)
}

export function listDomains(tenantId: string) {
  return request<TenantDomain[]>('GET', `/tenant/${tenantId}/domains`)
}

export function addDomain(tenantId: string, data: AddDomainRequest) {
  return request<TenantDomain[]>('POST', `/tenant/${tenantId}/domains`, data)
}

export function removeDomain(tenantId: string, domain: string) {
  return request<TenantDomain[]>('DELETE', `/tenant/${tenantId}/domains?domain=${encodeURIComponent(domain)}`)
}
