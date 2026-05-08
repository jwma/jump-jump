/// <reference types="vite/client" />

import 'vue-router'

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

declare module 'vue-router' {
  interface RouteMeta {
    title?: string
    requiresAuth?: boolean
    requiresTenant?: boolean
    requiresSuper?: boolean
    requiredRole?: string
  }
}
