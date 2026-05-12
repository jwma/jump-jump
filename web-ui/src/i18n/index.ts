import { createI18n } from 'vue-i18n'
import en from './locales/en'
import zh from './locales/zh'

function detectLocale(): string {
  const supported = ['en', 'zh']
  const saved = typeof localStorage !== 'undefined' ? localStorage.getItem('locale') : null
  if (saved && supported.includes(saved)) return saved
  if (typeof navigator !== 'undefined' && navigator.language) {
    const lang = navigator.language.toLowerCase()
    if (lang.startsWith('zh')) return 'zh'
  }
  return 'en'
}

const i18n = createI18n({
  legacy: false,
  locale: detectLocale(),
  fallbackLocale: 'en',
  messages: { en, zh },
})

export default i18n
