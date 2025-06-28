<script setup lang="ts">
import Card from '../common/Card.vue'
import { CheckCircleIcon, ExclamationTriangleIcon, XCircleIcon } from '@heroicons/vue/24/outline'

interface Props {
  status?: 'healthy' | 'warning' | 'error'
  lastChecked?: string
}

const { status = 'healthy', lastChecked = 'Just now' } = defineProps<Props>()

const statusConfig = {
  healthy: {
    icon: CheckCircleIcon,
    color: 'text-green-500',
    message: 'All systems operational'
  },
  warning: {
    icon: ExclamationTriangleIcon,
    color: 'text-yellow-500',
    message: 'Some issues detected'
  },
  error: {
    icon: XCircleIcon,
    color: 'text-red-500',
    message: 'System errors detected'
  }
}
</script>

<template>
  <Card title="System Health">
    <div class="flex items-center justify-center h-32 bg-gray-50 dark:bg-gray-700 rounded-lg">
      <div class="text-center">
        <component 
          :is="statusConfig[status].icon" 
          class="mx-auto h-8 w-8" 
          :class="statusConfig[status].color" 
        />
        <p class="mt-2 text-sm text-gray-900 dark:text-white font-medium">{{ statusConfig[status].message }}</p>
        <p class="text-xs text-gray-500 dark:text-gray-400">Last checked: {{ lastChecked }}</p>
      </div>
    </div>
  </Card>
</template>