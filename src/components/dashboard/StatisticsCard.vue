<script setup lang="ts">
interface Props {
  title: string
  value: string | number
  icon: string
  iconColor?: string
}

withDefaults(defineProps<Props>(), {
  iconColor: 'text-blue-600'
})
</script>

<template>
  <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
    <div class="flex items-center">
      <div class="flex-shrink-0">
        <svg class="h-8 w-8" :class="iconColor" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="icon" />
        </svg>
      </div>
      <div class="ml-5 w-0 flex-1">
        <dl>
          <dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">{{ title }}</dt>
          <dd class="text-lg font-semibold text-gray-900 dark:text-white">{{ value }}</dd>
        </dl>
      </div>
    </div>
  </div>
</template>