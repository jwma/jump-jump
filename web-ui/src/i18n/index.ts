import { createI18n } from 'vue-i18n'
import en from './locales/en'
import zh from './locales/zh'

const saved = typeof localStorage !== 'undefined' ? localStorage.getItem('locale') : null

const i18n = createI18n({
  legacy: false,
  locale: saved || 'en',
  fallbackLocale: 'en',
  messages: { en, zh },
})

export default i18n
