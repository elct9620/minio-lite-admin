<script setup lang="ts">
import Card from '../common/Card.vue'
import LoadingSpinner from '../common/LoadingSpinner.vue'
import ErrorMessage from '../common/ErrorMessage.vue'

interface ServerInfo {
  mode: string
  region?: string
  deploymentId: string
}

interface Props {
  loading?: boolean
  error?: string
  serverInfo?: ServerInfo | null
}

withDefaults(defineProps<Props>(), {
  loading: false
})
</script>

<template>
  <Card title="Server Information">
    <LoadingSpinner v-if="loading" text="Loading server information..." />
    
    <ErrorMessage v-else-if="error" :message="error" />

    <div v-else-if="serverInfo" class="space-y-4">
      <div class="flex justify-between items-center py-2 border-b border-gray-100 dark:border-gray-700">
        <span class="text-sm font-medium text-gray-500 dark:text-gray-400">Mode</span>
        <span class="text-sm font-semibold text-gray-900 dark:text-white">{{ serverInfo.mode }}</span>
      </div>
      <div class="flex justify-between items-center py-2 border-b border-gray-100 dark:border-gray-700">
        <span class="text-sm font-medium text-gray-500 dark:text-gray-400">Region</span>
        <span class="text-sm font-semibold text-gray-900 dark:text-white">{{ serverInfo.region || 'Not set' }}</span>
      </div>
      <div class="flex justify-between items-center py-2">
        <span class="text-sm font-medium text-gray-500 dark:text-gray-400">Deployment ID</span>
        <span class="text-xs font-mono text-gray-900 dark:text-white break-all max-w-40">{{ serverInfo.deploymentId }}</span>
      </div>
    </div>
  </Card>
</template>