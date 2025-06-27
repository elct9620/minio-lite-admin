<script setup lang="ts">
import Card from '../common/Card.vue'

interface Props {
  status?: 'healthy' | 'warning' | 'error'
  lastChecked?: string
}

withDefaults(defineProps<Props>(), {
  status: 'healthy',
  lastChecked: 'Just now'
})

const statusConfig = {
  healthy: {
    icon: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
    color: 'text-green-500',
    message: 'All systems operational'
  },
  warning: {
    icon: 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z',
    color: 'text-yellow-500',
    message: 'Some issues detected'
  },
  error: {
    icon: 'M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z',
    color: 'text-red-500',
    message: 'System errors detected'
  }
}
</script>

<template>
  <Card title="System Health">
    <div class="flex items-center justify-center h-32 bg-gray-50 dark:bg-gray-700 rounded-lg">
      <div class="text-center">
        <svg class="mx-auto h-8 w-8" :class="statusConfig[status].color" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="statusConfig[status].icon" />
        </svg>
        <p class="mt-2 text-sm text-gray-900 dark:text-white font-medium">{{ statusConfig[status].message }}</p>
        <p class="text-xs text-gray-500 dark:text-gray-400">Last checked: {{ lastChecked }}</p>
      </div>
    </div>
  </Card>
</template>