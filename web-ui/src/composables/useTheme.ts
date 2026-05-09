import { ref, watch } from 'vue'

type Theme = 'light' | 'dark' | 'system'

const theme = ref<Theme>(loadTheme())
const darkMode = ref(false)

function loadTheme(): Theme {
  if (typeof localStorage !== 'undefined') {
    return (localStorage.getItem('theme') as Theme) || 'system'
  }
  return 'system'
}

function applyDark(dark: boolean) {
  darkMode.value = dark
  document.documentElement.classList.toggle('dark', dark)
}

function resolveTheme(t: Theme) {
  if (t === 'system') {
    return window.matchMedia('(prefers-color-scheme: dark)').matches
  }
  return t === 'dark'
}

function initTheme() {
  applyDark(resolveTheme(theme.value))

  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
    if (theme.value === 'system') {
      applyDark(e.matches)
    }
  })
}

watch(theme, (val) => {
  localStorage.setItem('theme', val)
  applyDark(resolveTheme(val))
})

export function useTheme() {
  return { theme, darkMode, initTheme }
}
