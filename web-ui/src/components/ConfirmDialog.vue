<script setup lang="ts">
import { ref, watch, nextTick, onBeforeUnmount } from 'vue'

const props = withDefaults(defineProps<{
  open: boolean
  title: string
  message?: string
  variant?: 'danger' | 'primary'
  confirmText?: string
  cancelText?: string
  loading?: boolean
}>(), {
  variant: 'primary',
  confirmText: 'Confirm',
  cancelText: 'Cancel',
  loading: false,
})

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()

const confirmBtnRef = ref<HTMLButtonElement | null>(null)
let previousFocusEl: HTMLElement | null = null

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    emit('cancel')
  }
}

watch(() => props.open, (isOpen) => {
  if (isOpen) {
    previousFocusEl = document.activeElement as HTMLElement
    document.addEventListener('keydown', handleKeydown)
    nextTick(() => confirmBtnRef.value?.focus())
  } else {
    document.removeEventListener('keydown', handleKeydown)
    if (previousFocusEl) {
      previousFocusEl.focus()
      previousFocusEl = null
    }
  }
})

onBeforeUnmount(() => {
  document.removeEventListener('keydown', handleKeydown)
  if (previousFocusEl) {
    previousFocusEl.focus()
  }
})
</script>

<template>
  <Teleport to="body">
    <div
      v-if="open"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
      @click.self="emit('cancel')"
    >
      <div
        role="dialog"
        aria-modal="true"
        tabindex="-1"
        class="mx-4 w-full max-w-sm rounded-lg bg-white p-6 shadow-xl"
      >
        <h3 class="text-lg font-semibold text-gray-900">{{ title }}</h3>
        <slot>
          <p v-if="message" class="mt-2 text-sm text-gray-600">{{ message }}</p>
        </slot>
        <div class="mt-4 flex justify-end gap-2">
          <button
            class="rounded-md px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100"
            @click="emit('cancel')"
          >
            {{ cancelText }}
          </button>
          <button
            ref="confirmBtnRef"
            :disabled="loading"
            class="rounded-md px-4 py-2 text-sm font-medium text-white disabled:opacity-50"
            :class="variant === 'danger' ? 'bg-red-600 hover:bg-red-700' : 'bg-blue-600 hover:bg-blue-700'"
            @click="emit('confirm')"
          >
            {{ confirmText }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>
