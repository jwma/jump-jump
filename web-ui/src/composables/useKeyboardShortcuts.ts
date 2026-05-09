import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRouter, useRoute } from 'vue-router'

interface ShortcutDef {
  key: string
  ctrl?: boolean
  description: string
  action: () => void
  page?: string
}

export const showHelp = ref(false)

let installed = false
let handleKeydown: ((e: KeyboardEvent) => void) | null = null

export function useKeyboardShortcuts() {
  const router = useRouter()
  const route = useRoute()

  const shortcuts: ShortcutDef[] = [
    {
      key: 'k',
      ctrl: true,
      description: 'Open global search',
      action: () => {
        document.dispatchEvent(new CustomEvent('open-global-search'))
      },
    },
    {
      key: '?',
      description: 'Show keyboard shortcuts',
      action: () => {
        showHelp.value = !showHelp.value
      },
    },
    {
      key: 'n',
      description: 'Create new short link',
      page: 'short-links',
      action: () => {
        if (route.name === 'short-links') {
          router.push({ name: 'short-link-create' })
        }
      },
    },
    {
      key: '/',
      description: 'Focus search on current page',
      page: 'short-links',
      action: () => {
        const searchInput = document.querySelector<HTMLInputElement>(
          '[data-shortcut-search]',
        )
        if (searchInput) {
          searchInput.focus()
        }
      },
    },
    {
      key: 'Escape',
      description: 'Close panels / Cancel search',
      action: () => {
        showHelp.value = false
        const searchInput = document.querySelector<HTMLInputElement>(
          '[data-shortcut-search]',
        )
        if (searchInput && document.activeElement === searchInput) {
          searchInput.blur()
        }
      },
    },
  ]

  if (!installed) {
    installed = true

    handleKeydown = (e: KeyboardEvent) => {
      const target = e.target as HTMLElement
      const isInput = target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable

      for (const shortcut of shortcuts) {
        if (shortcut.page && route.name !== shortcut.page) continue

        if (shortcut.key === 'Escape') {
          if (shortcut.key === e.key) {
            shortcut.action()
            return
          }
          continue
        }

        if (isInput) continue

        const keyMatch = shortcut.key.toLowerCase() === e.key.toLowerCase()
        const ctrlMatch = shortcut.ctrl ? (e.ctrlKey || e.metaKey) : !(e.ctrlKey || e.metaKey)

        if (keyMatch && ctrlMatch) {
          e.preventDefault()
          shortcut.action()
          return
        }
      }
    }

    onMounted(() => document.addEventListener('keydown', handleKeydown!))
    onBeforeUnmount(() => {
      if (handleKeydown) {
        document.removeEventListener('keydown', handleKeydown)
      }
      installed = false
      handleKeydown = null
    })
  }

  return { showHelp, shortcuts }
}
