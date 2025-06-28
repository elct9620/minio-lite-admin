<script setup lang="ts">
import type { Component } from 'vue'

interface Props {
  title: string
  value: string | number
  icon: Component
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
        <component :is="icon" class="h-8 w-8" :class="iconColor" />
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