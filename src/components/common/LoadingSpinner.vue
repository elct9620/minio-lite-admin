<script setup lang="ts">
interface Props {
  size?: 'sm' | 'md' | 'lg'
  text?: string
}

withDefaults(defineProps<Props>(), {
  size: 'md'
})

const sizeClasses = {
  sm: 'h-4 w-4',
  md: 'h-8 w-8',
  lg: 'h-12 w-12'
}
</script>

<template>
  <div class="flex flex-col items-center justify-center">
    <div
      class="animate-spin rounded-full border-b-2 border-blue-600"
      :class="sizeClasses[size]"
    ></div>
    <p v-if="text" class="mt-2 text-sm text-gray-500 dark:text-gray-400">{{ text }}</p>
  </div>
</template>