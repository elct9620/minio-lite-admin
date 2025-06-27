<script setup lang="ts">
interface Props {
  title?: string
  padding?: boolean
}

withDefaults(defineProps<Props>(), {
  padding: true
})
</script>

<template>
  <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700">
    <div v-if="title" class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
      <h2 class="text-lg font-medium text-gray-900 dark:text-white">{{ title }}</h2>
    </div>
    <div :class="{ 'p-6': padding && !title, 'p-6 pt-4': padding && title }">
      <slot />
    </div>
  </div>
</template>