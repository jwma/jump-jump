<script setup lang="ts">
import { useKeyboardShortcuts } from '@/composables/useKeyboardShortcuts'
import { X, Keyboard } from 'lucide-vue-next'

const { showHelp } = useKeyboardShortcuts()

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
    <div
      v-if="showHelp"
      class="fixed inset-0 z-[70] flex items-center justify-center"
    >
      <div class="fixed inset-0 bg-black/50" @click="close" />
      <div
        class="relative w-full max-w-md overflow-hidden rounded-xl border border-gray-200 bg-white shadow-2xl"
        @keydown="handleKeydown"
      >
        <div class="flex items-center justify-between border-b px-5 py-4">
          <div class="flex items-center gap-2">
            <Keyboard class="h-5 w-5 text-gray-500" />
            <h3 class="text-base font-semibold text-gray-900">Keyboard Shortcuts</h3>
          </div>
          <button
            class="rounded p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-600"
            @click="close"
          >
            <X class="h-4 w-4" />
          </button>
        </div>

        <div class="space-y-1 px-5 py-4">
          <div class="flex items-center justify-between py-2">
            <span class="text-sm text-gray-600">Open global search</span>
            <div class="flex items-center gap-1">
              <kbd class="rounded border border-gray-200 bg-gray-50 px-2 py-0.5 text-xs font-medium text-gray-500">Ctrl</kbd>
              <kbd class="rounded border border-gray-200 bg-gray-50 px-2 py-0.5 text-xs font-medium text-gray-500">K</kbd>
            </div>
          </div>

          <div class="border-t border-gray-100" />

          <div class="flex items-center justify-between py-2">
            <span class="text-sm text-gray-600">Show this help panel</span>
            <kbd class="rounded border border-gray-200 bg-gray-50 px-2 py-0.5 text-xs font-medium text-gray-500">?</kbd>
          </div>

          <div class="border-t border-gray-100" />

          <p class="pt-2 text-xs font-medium uppercase text-gray-400">Short Links Page</p>

          <div class="flex items-center justify-between py-2">
            <span class="text-sm text-gray-600">Create new short link</span>
            <kbd class="rounded border border-gray-200 bg-gray-50 px-2 py-0.5 text-xs font-medium text-gray-500">N</kbd>
          </div>

          <div class="flex items-center justify-between py-2">
            <span class="text-sm text-gray-600">Focus search</span>
            <kbd class="rounded border border-gray-200 bg-gray-50 px-2 py-0.5 text-xs font-medium text-gray-500">/</kbd>
          </div>

          <div class="border-t border-gray-100" />

          <div class="flex items-center justify-between py-2">
            <span class="text-sm text-gray-600">Close panels / Cancel</span>
            <kbd class="rounded border border-gray-200 bg-gray-50 px-2 py-0.5 text-xs font-medium text-gray-500">Esc</kbd>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>
