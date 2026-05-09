<script setup lang="ts">
import { useI18n } from 'vue-i18n'

defineOptions({ name: 'PreferencesPage' })

const { t, locale } = useI18n()

function setLocale(lang: string) {
  locale.value = lang
  localStorage.setItem('locale', lang)
}
</script>

<template>
  <div>
    <h1 class="text-xl font-semibold text-gray-900">{{ t('preferences.title') }}</h1>
    <p class="mt-1 text-sm text-gray-500">{{ t('preferences.subtitle') }}</p>
    <div class="mt-6 rounded-lg border bg-white p-6">
      <p class="text-sm text-gray-400">{{ t('preferences.comingSoon') }}</p>
    </div>
    <div class="mt-6 rounded-lg border bg-white p-6">
      <h2 class="text-sm font-semibold text-gray-900">{{ t('preferences.language') }}</h2>
      <p class="mt-1 text-xs text-gray-400">{{ t('preferences.languageDesc') }}</p>
      <div class="mt-3 flex gap-2">
        <button
          v-for="loc in ['en', 'zh']"
          :key="loc"
          class="rounded-md px-3 py-1.5 text-sm font-medium transition-colors"
          :class="locale === loc ? 'bg-blue-600 text-white' : 'border text-gray-600 hover:bg-gray-50'"
          @click="setLocale(loc)"
        >
          {{ loc === 'en' ? t('preferences.en') : t('preferences.zh') }}
        </button>
      </div>
    </div>
  </div>
</template>
