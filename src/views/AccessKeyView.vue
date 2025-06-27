<script setup lang="ts">
import Card from '../components/common/Card.vue'
</script>

<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Access Keys</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
            Manage MinIO access keys and service accounts
          </p>
        </div>
        <button class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors">
          Create Access Key
        </button>
      </div>
    </div>

    <!-- Access Keys List -->
    <Card title="Access Keys">
      <div class="flex items-center justify-center h-48 bg-gray-50 dark:bg-gray-700 rounded-lg">
        <div class="text-center">
          <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
          </svg>
          <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">Access key management</p>
          <p class="text-xs text-gray-400 dark:text-gray-500">Coming soon</p>
        </div>
      </div>
    </Card>

    <!-- Service Accounts -->
    <Card title="Service Accounts">
      <div class="flex items-center justify-center h-48 bg-gray-50 dark:bg-gray-700 rounded-lg">
        <div class="text-center">
          <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
          </svg>
          <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">Service account management</p>
          <p class="text-xs text-gray-400 dark:text-gray-500">Coming soon</p>
        </div>
      </div>
    </Card>
  </div>
</template>