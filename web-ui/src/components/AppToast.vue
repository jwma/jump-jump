<script setup lang="ts">
import { useToast } from '@/composables/useToast'
import { CheckCircle, XCircle, AlertTriangle, Info, X } from 'lucide-vue-next'

const { toasts, remove } = useToast()

const iconMap = {
  success: CheckCircle,
  error: XCircle,
  warning: AlertTriangle,
  info: Info,
} as const

const styleMap = {
  success:
    'bg-green-50 border-green-200 text-green-800 dark:bg-green-900/90 dark:border-green-700 dark:text-green-200',
  error:
    'bg-red-50 border-red-200 text-red-800 dark:bg-red-900/90 dark:border-red-700 dark:text-red-200',
  warning:
    'bg-amber-50 border-amber-200 text-amber-800 dark:bg-amber-900/90 dark:border-amber-700 dark:text-amber-200',
  info: 'bg-blue-50 border-blue-200 text-blue-800 dark:bg-blue-900/90 dark:border-blue-700 dark:text-blue-200',
} as const

const iconStyleMap = {
  success: 'text-green-500 dark:text-green-400',
  error: 'text-red-500 dark:text-red-400',
  warning: 'text-amber-500 dark:text-amber-400',
  info: 'text-blue-500 dark:text-blue-400',
} as const
</script>

<template>
  <Teleport to="body">
    <div class="fixed right-4 top-4 z-[100] flex flex-col gap-2" aria-live="polite">
      <TransitionGroup name="toast">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          class="flex w-80 items-start gap-2.5 rounded-lg border p-3 shadow-lg"
          :class="styleMap[toast.type]"
        >
          <component
            :is="iconMap[toast.type]"
            class="mt-0.5 h-4 w-4 shrink-0"
            :class="iconStyleMap[toast.type]"
          />
          <p class="flex-1 text-sm">{{ toast.message }}</p>
          <button
            type="button"
            aria-label="Dismiss notification"
            class="shrink-0 rounded p-0.5 opacity-60 transition-opacity hover:opacity-100"
            @click="remove(toast.id)"
          >
            <X class="h-3.5 w-3.5" />
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped>
.toast-enter-active {
  transition: all 0.3s ease-out;
}
.toast-leave-active {
  transition: all 0.2s ease-in;
}
.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}
.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
.toast-move {
  transition: transform 0.2s ease;
}
</style>
