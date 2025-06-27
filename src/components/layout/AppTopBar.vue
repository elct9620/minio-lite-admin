<script setup lang="ts">
interface Props {
  title: string
}

defineProps<Props>()
const emit = defineEmits<{
  'toggle-sidebar': []
}>()
</script>

<template>
  <div class="sticky top-0 z-10 bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 lg:hidden">
    <div class="flex items-center justify-between h-16 px-4">
      <button
        @click="emit('toggle-sidebar')"
        class="text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white"
      >
        <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </button>
      <h1 class="text-lg font-semibold text-gray-900 dark:text-white">{{ title }}</h1>
      <div></div>
    </div>
  </div>
</template>