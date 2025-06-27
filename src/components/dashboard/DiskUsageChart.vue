<script setup lang="ts">
import Card from '../common/Card.vue'

// Props for future implementation when real data is available
interface Props {
  usedSpace?: number
  totalSpace?: number
}

withDefaults(defineProps<Props>(), {
  usedSpace: 0,
  totalSpace: 0
})
</script>

<template>
  <Card title="Disk Usage">
    <div class="flex items-center justify-center h-48 bg-gray-50 dark:bg-gray-700 rounded-lg">
      <div class="text-center">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
        <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">Disk usage chart</p>
        <p class="text-xs text-gray-400 dark:text-gray-500">Coming soon</p>
      </div>
    </div>
  </Card>
</template>