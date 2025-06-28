<script setup lang="ts">
import { Bars3Icon } from '@heroicons/vue/24/outline'

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
        <Bars3Icon class="h-6 w-6" />
      </button>
      <h1 class="text-lg font-semibold text-gray-900 dark:text-white">{{ title }}</h1>
      <div></div>
    </div>
  </div>
</template>