<script setup lang="ts">
import { XCircleIcon, ExclamationTriangleIcon, InformationCircleIcon } from '@heroicons/vue/24/outline'

interface Props {
  message: string
  type?: 'error' | 'warning' | 'info'
}

withDefaults(defineProps<Props>(), {
  type: 'error'
})

const typeClasses = {
  error: 'text-red-600 dark:text-red-400',
  warning: 'text-yellow-600 dark:text-yellow-400',
  info: 'text-blue-600 dark:text-blue-400'
}

const icons = {
  error: XCircleIcon,
  warning: ExclamationTriangleIcon,
  info: InformationCircleIcon
}
</script>

<template>
  <div class="flex items-center gap-3">
    <component 
      :is="icons[type]" 
      class="h-5 w-5 flex-shrink-0" 
      :class="typeClasses[type]" 
    />
    <p :class="typeClasses[type]">{{ message }}</p>
  </div>
</template>