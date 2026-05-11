import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { listShortLinks } from '@/api/short-link'
import { listMembers } from '@/api/member'
import type { ShortLinkData, TenantMember } from '@/types/api'

export interface SearchResult {
  type: 'short-link' | 'member'
  id: string
  label: string
  subtitle: string
  route: { name: string; params?: Record<string, string> }
}

const query = ref('')
const results = ref<SearchResult[]>([])
const loading = ref(false)
const open = ref(false)
let debounceTimer: ReturnType<typeof setTimeout> | null = null

export function useGlobalSearch() {
  const router = useRouter()
  const auth = useAuthStore()

  const hasResults = computed(() => results.value.length > 0)

  async function search(q: string) {
    if (!q.trim()) {
      results.value = []
      loading.value = false
      return
    }

    if (!auth.currentTenantId) {
      results.value = []
      loading.value = false
      return
    }

    loading.value = true
    const found: SearchResult[] = []
    const lower = q.toLowerCase()

    try {
      const linkData = await listShortLinks(1, 20)
      const matchingLinks = (linkData.shortLinks || []).filter(
        (l: ShortLinkData) =>
          l.id.toLowerCase().includes(lower) || l.url.toLowerCase().includes(lower),
      )
      for (const link of matchingLinks.slice(0, 5)) {
        found.push({
          type: 'short-link',
          id: link.id,
          label: link.id,
          subtitle: link.url,
          route: { name: 'short-link-detail', params: { id: link.id } },
        })
      }
    } catch {
      // ignore
    }

    try {
      const members = await listMembers(auth.currentTenantId!)
      const matchingMembers = (members || []).filter((m: TenantMember) =>
        m.username.toLowerCase().includes(lower),
      )
      for (const member of matchingMembers.slice(0, 3)) {
        found.push({
          type: 'member',
          id: member.userId,
          label: member.username,
          subtitle: member.role,
          route: { name: 'members' },
        })
      }
    } catch {
      // ignore
    }

    results.value = found
    loading.value = false
  }

  watch(query, (val) => {
    if (debounceTimer) clearTimeout(debounceTimer)
    if (!val.trim()) {
      results.value = []
      loading.value = false
      return
    }
    loading.value = true
    debounceTimer = setTimeout(() => search(val), 300)
  })

  function openSearch() {
    open.value = true
  }

  function closeSearch() {
    open.value = false
    query.value = ''
    results.value = []
  }

  function selectResult(result: SearchResult) {
    closeSearch()
    router.push(result.route)
  }

  return {
    query,
    results,
    loading,
    open,
    hasResults,
    openSearch,
    closeSearch,
    selectResult,
  }
}
