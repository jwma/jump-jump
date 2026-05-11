<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { showHelp } from '@/composables/useKeyboardShortcuts'
import { X, Keyboard } from 'lucide-vue-next'

const { t } = useI18n()

function close() {
  showHelp.value = false
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    close()
  }
}
</script>

<template>
  <Teleport to="body">
    <div v-if="showHelp" class="fixed inset-0 z-[70] flex items-center justify-center">
      <div class="fixed inset-0 bg-black/50" @click="close" />
      <div
        class="relative w-full max-w-md overflow-hidden rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 shadow-2xl"
        @keydown="handleKeydown"
      >
        <div class="flex items-center justify-between border-b dark:border-gray-700 px-5 py-4">
          <div class="flex items-center gap-2">
            <Keyboard class="h-5 w-5 text-gray-500 dark:text-gray-400" />
            <h3 class="text-base font-semibold text-gray-900 dark:text-white">
              {{ t('shortcuts.title') }}
            </h3>
          </div>
          <button
            class="rounded p-1 text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 hover:text-gray-600 dark:hover:text-gray-300"
            @click="close"
          >
            <X class="h-4 w-4" />
          </button>
        </div>

        <div class="space-y-1 px-5 py-4">
          <div class="flex items-center justify-between py-2">
            <span class="text-sm text-gray-600 dark:text-gray-400">{{ t('shortcuts.openGlobalSearch') }}</span>
            <div class="flex items-center gap-1">
              <kbd
                class="rounded border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 px-2 py-0.5 text-xs font-medium text-gray-500 dark:text-gray-400"
                >Ctrl</kbd
              >
              <kbd
                class="rounded border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 px-2 py-0.5 text-xs font-medium text-gray-500 dark:text-gray-400"
                >K</kbd
              >
            </div>
          </div>

          <div class="border-t border-gray-100 dark:border-gray-800" />

          <div class="flex items-center justify-between py-2">
            <span class="text-sm text-gray-600 dark:text-gray-400">{{ t('shortcuts.showHelp') }}</span>
            <kbd
              class="rounded border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 px-2 py-0.5 text-xs font-medium text-gray-500 dark:text-gray-400"
              >?</kbd
            >
          </div>

          <div class="border-t border-gray-100 dark:border-gray-800" />

          <p class="pt-2 text-xs font-medium uppercase text-gray-400 dark:text-gray-500">
            {{ t('shortcuts.shortLinksPage') }}
          </p>

          <div class="flex items-center justify-between py-2">
            <span class="text-sm text-gray-600 dark:text-gray-400">{{ t('shortcuts.createNewLink') }}</span>
            <kbd
              class="rounded border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 px-2 py-0.5 text-xs font-medium text-gray-500 dark:text-gray-400"
              >N</kbd
            >
          </div>

          <div class="flex items-center justify-between py-2">
            <span class="text-sm text-gray-600 dark:text-gray-400">{{ t('shortcuts.focusSearch') }}</span>
            <kbd
              class="rounded border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 px-2 py-0.5 text-xs font-medium text-gray-500 dark:text-gray-400"
              >/</kbd
              >
          </div>

          <div class="border-t border-gray-100 dark:border-gray-800" />

          <div class="flex items-center justify-between py-2">
            <span class="text-sm text-gray-600 dark:text-gray-400">{{ t('shortcuts.closePanels') }}</span>
            <kbd
              class="rounded border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 px-2 py-0.5 text-xs font-medium text-gray-500 dark:text-gray-400"
              >Esc</kbd
            >
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>
