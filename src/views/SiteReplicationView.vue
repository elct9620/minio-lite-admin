<script setup lang="ts">
import Card from '../components/common/Card.vue'
</script>

<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Site Replication</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
            Configure and manage MinIO site replication settings
          </p>
        </div>
        <button class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors">
          Configure Replication
        </button>
      </div>
    </div>

    <!-- Replication Status -->
    <Card title="Replication Status">
      <div class="flex items-center justify-center h-48 bg-gray-50 dark:bg-gray-700 rounded-lg">
        <div class="text-center">
          <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
          </svg>
          <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">Site replication status</p>
          <p class="text-xs text-gray-400 dark:text-gray-500">Coming soon</p>
        </div>
      </div>
    </Card>

    <!-- Replication Sites -->
    <Card title="Replication Sites">
      <div class="flex items-center justify-center h-48 bg-gray-50 dark:bg-gray-700 rounded-lg">
        <div class="text-center">
          <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
          </svg>
          <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">Replication site management</p>
          <p class="text-xs text-gray-400 dark:text-gray-500">Coming soon</p>
        </div>
      </div>
    </Card>
  </div>
</template>