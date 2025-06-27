<script setup lang="ts">
import { ref, onMounted } from 'vue'

const serverInfo = ref(null)
const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    const response = await fetch('/api/server-info')
    if (!response.ok) {
      throw new Error('Failed to fetch server info')
    }
    serverInfo.value = await response.json()
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <!-- Header -->
    <header class="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <div class="flex items-center">
            <h1 class="text-xl font-semibold text-gray-900 dark:text-white">
              MinIO Lite Admin
            </h1>
          </div>
          <div class="flex items-center space-x-4">
            <span class="text-sm text-gray-500 dark:text-gray-400">Dashboard</span>
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
      <!-- Server Status Card -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6 mb-6">
        <h2 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Server Status</h2>
        
        <div v-if="loading" class="flex items-center justify-center py-8">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>

        <div v-else-if="error" class="text-red-600 dark:text-red-400 py-4">
          Error: {{ error }}
        </div>

        <div v-else-if="serverInfo" class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
            <div class="text-sm font-medium text-gray-500 dark:text-gray-400">Mode</div>
            <div class="mt-1 text-lg font-semibold text-gray-900 dark:text-white">
              {{ serverInfo.mode }}
            </div>
          </div>
          <div class="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
            <div class="text-sm font-medium text-gray-500 dark:text-gray-400">Region</div>
            <div class="mt-1 text-lg font-semibold text-gray-900 dark:text-white">
              {{ serverInfo.region || 'Not set' }}
            </div>
          </div>
          <div class="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
            <div class="text-sm font-medium text-gray-500 dark:text-gray-400">Deployment ID</div>
            <div class="mt-1 text-sm font-mono text-gray-900 dark:text-white break-all">
              {{ serverInfo.deploymentId }}
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
        <h2 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Quick Actions</h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <button class="p-4 text-left border border-gray-200 dark:border-gray-600 rounded-lg hover:border-blue-500 dark:hover:border-blue-400 transition-colors">
            <div class="text-sm font-medium text-gray-900 dark:text-white">Disk Status</div>
            <div class="text-xs text-gray-500 dark:text-gray-400 mt-1">Monitor storage</div>
          </button>
          <button class="p-4 text-left border border-gray-200 dark:border-gray-600 rounded-lg hover:border-blue-500 dark:hover:border-blue-400 transition-colors">
            <div class="text-sm font-medium text-gray-900 dark:text-white">Access Keys</div>
            <div class="text-xs text-gray-500 dark:text-gray-400 mt-1">Manage credentials</div>
          </button>
          <button class="p-4 text-left border border-gray-200 dark:border-gray-600 rounded-lg hover:border-blue-500 dark:hover:border-blue-400 transition-colors">
            <div class="text-sm font-medium text-gray-900 dark:text-white">Replication</div>
            <div class="text-xs text-gray-500 dark:text-gray-400 mt-1">Site configuration</div>
          </button>
          <button class="p-4 text-left border border-gray-200 dark:border-gray-600 rounded-lg hover:border-blue-500 dark:hover:border-blue-400 transition-colors">
            <div class="text-sm font-medium text-gray-900 dark:text-white">Settings</div>
            <div class="text-xs text-gray-500 dark:text-gray-400 mt-1">System preferences</div>
          </button>
        </div>
      </div>
    </main>
  </div>
</template>
