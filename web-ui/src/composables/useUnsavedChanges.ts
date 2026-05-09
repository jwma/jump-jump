import { onMounted, onBeforeUnmount, type Ref } from 'vue'
import { onBeforeRouteLeave } from 'vue-router'

export function useUnsavedChanges(isDirty: Ref<boolean> | (() => boolean)) {
  const getDirty = typeof isDirty === 'function' ? isDirty : () => isDirty.value

  function handleBeforeUnload(e: BeforeUnloadEvent) {
    if (getDirty()) {
      e.preventDefault()
      e.returnValue = ''
    }
  }

  onMounted(() => {
    window.addEventListener('beforeunload', handleBeforeUnload)
  })

  onBeforeUnmount(() => {
    window.removeEventListener('beforeunload', handleBeforeUnload)
  })

  onBeforeRouteLeave(() => {
    if (getDirty()) {
      return window.confirm('You have unsaved changes. Are you sure you want to leave?')
    }
  })
}
