<script setup lang="ts">
import { ref, watch, nextTick, onMounted, onBeforeUnmount } from 'vue'
import { useGlobalSearch } from '@/composables/useGlobalSearch'
import type { SearchResult } from '@/composables/useGlobalSearch'
import { Search, Link, Users, Loader2 } from 'lucide-vue-next'

const { query, results, loading, open, hasResults, openSearch, closeSearch, selectResult } = useGlobalSearch()

const inputRef = ref<HTMLInputElement | null>(null)
const activeIndex = ref(-1)

watch(open, async (isOpen) => {
  if (isOpen) {
    await nextTick()
    inputRef.value?.focus()
  }
})

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    closeSearch()
    return
  }
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    activeIndex.value = Math.min(activeIndex.value + 1, results.value.length - 1)
    return
  }
  if (e.key === 'ArrowUp') {
    e.preventDefault()
    activeIndex.value = Math.max(activeIndex.value - 1, -1)
    return
  }
  if (e.key === 'Enter' && activeIndex.value >= 0 && results.value[activeIndex.value]) {
    e.preventDefault()
    selectResult(results.value[activeIndex.value])
    return
  }
}

function getResultIcon(type: SearchResult['type']) {
  return type === 'short-link' ? Link : Users
}

onMounted(() => {
  document.addEventListener('open-global-search', openSearch)
})

onBeforeUnmount(() => {
  document.removeEventListener('open-global-search', openSearch)
  closeSearch()
})
</script>

<template>
  <!-- Search trigger button -->
  <button
    class="flex items-center gap-2 rounded-lg border border-gray-200 bg-white px-3 py-1.5 text-sm text-gray-400 transition-colors hover:border-gray-300 hover:text-gray-500"
    @click="openSearch"
  >
    <Search class="h-4 w-4" />
    <span class="hidden sm:inline">Search...</span>
    <kbd class="hidden rounded border border-gray-200 bg-gray-50 px-1.5 py-0.5 text-[10px] font-medium text-gray-400 sm:inline">Ctrl+K</kbd>
  </button>

  <!-- Search modal -->
  <Teleport to="body">
    <div
      v-if="open"
      class="fixed inset-0 z-[60] flex items-start justify-center pt-[15vh]"
    >
      <div
        class="fixed inset-0 bg-black/50"
        @click="closeSearch"
      />
      <div
        class="relative w-full max-w-lg overflow-hidden rounded-xl border border-gray-200 bg-white shadow-2xl"
        @keydown="handleKeydown"
      >
        <!-- Search input -->
        <div class="flex items-center gap-3 border-b border-gray-100 px-4 py-3">
          <Search class="h-5 w-5 shrink-0 text-gray-400" />
          <input
            ref="inputRef"
            v-model="query"
            type="text"
            placeholder="Search short links, members..."
            class="flex-1 bg-transparent text-sm text-gray-900 placeholder-gray-400 focus:outline-none"
          />
          <Loader2 v-if="loading" class="h-4 w-4 animate-spin text-gray-400" />
          <kbd v-else class="rounded border border-gray-200 bg-gray-50 px-1.5 py-0.5 text-[10px] font-medium text-gray-400">ESC</kbd>
        </div>

        <!-- Results -->
        <div v-if="query.trim()" class="max-h-72 overflow-y-auto">
          <div v-if="!loading && !hasResults" class="px-4 py-8 text-center">
            <p class="text-sm text-gray-400">No results found for "{{ query }}"</p>
          </div>
          <ul v-else-if="hasResults">
            <li
              v-for="(result, index) in results"
              :key="result.type + result.id"
              :class="[
                'flex cursor-pointer items-center gap-3 px-4 py-2.5 transition-colors',
                index === activeIndex ? 'bg-blue-50' : 'hover:bg-gray-50',
              ]"
              @click="selectResult(result)"
              @mouseenter="activeIndex = index"
            >
              <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-gray-100">
                <component :is="getResultIcon(result.type)" class="h-4 w-4 text-gray-500" />
              </div>
              <div class="min-w-0 flex-1">
                <p class="text-sm font-medium text-gray-900">{{ result.label }}</p>
                <p class="truncate text-xs text-gray-400">{{ result.subtitle }}</p>
              </div>
              <span class="shrink-0 rounded-full bg-gray-100 px-2 py-0.5 text-[10px] font-medium text-gray-500">
                {{ result.type === 'short-link' ? 'Link' : 'Member' }}
              </span>
            </li>
          </ul>
        </div>

        <!-- Empty state -->
        <div v-else class="px-4 py-6 text-center">
          <p class="text-sm text-gray-400">Start typing to search...</p>
        </div>
      </div>
    </div>
  </Teleport>
</template>
