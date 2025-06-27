<script setup lang="ts">
import Card from '../common/Card.vue'
import LoadingSpinner from '../common/LoadingSpinner.vue'
import ErrorMessage from '../common/ErrorMessage.vue'
import { useDataUsage } from '../../composables/useDataUsage'
import { formatBytes, formatPercentage } from '../../utils/formatBytes'

const { dataUsage, loading, error } = useDataUsage()
</script>

<template>
  <Card title="Disk Usage">
    <LoadingSpinner v-if="loading" text="Loading disk usage..." />
    
    <ErrorMessage v-else-if="error" :message="error" />

    <div v-else-if="dataUsage" class="space-y-6">
      <!-- Progress Bar -->
      <div class="space-y-2">
        <div class="flex justify-between text-sm">
          <span class="text-gray-500 dark:text-gray-400">Storage Used</span>
          <span class="font-medium text-gray-900 dark:text-white">
            {{ formatBytes(dataUsage.totalUsedCapacity) }} of {{ formatBytes(dataUsage.totalCapacity) }}
          </span>
        </div>
        <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-3">
          <div 
            class="bg-gradient-to-r from-blue-500 to-blue-600 h-3 rounded-full transition-all duration-300"
            :style="{ width: `${Math.min(dataUsage.usagePercentage, 100)}%` }"
          ></div>
        </div>
        <div class="flex justify-between text-xs text-gray-500 dark:text-gray-400">
          <span>{{ formatPercentage(dataUsage.usagePercentage) }} Used</span>
          <span>{{ formatBytes(dataUsage.totalFreeCapacity) }} Free</span>
        </div>
      </div>

      <!-- Disk Status Grid -->
      <div class="grid grid-cols-3 gap-4">
        <div class="text-center p-3 bg-green-50 dark:bg-green-900/20 rounded-lg">
          <div class="text-lg font-semibold text-green-600 dark:text-green-400">
            {{ dataUsage.onlineDisks }}
          </div>
          <div class="text-xs text-green-600 dark:text-green-400 font-medium">Online</div>
        </div>
        <div class="text-center p-3 bg-red-50 dark:bg-red-900/20 rounded-lg">
          <div class="text-lg font-semibold text-red-600 dark:text-red-400">
            {{ dataUsage.offlineDisks }}
          </div>
          <div class="text-xs text-red-600 dark:text-red-400 font-medium">Offline</div>
        </div>
        <div class="text-center p-3 bg-yellow-50 dark:bg-yellow-900/20 rounded-lg">
          <div class="text-lg font-semibold text-yellow-600 dark:text-yellow-400">
            {{ dataUsage.healingDisks }}
          </div>
          <div class="text-xs text-yellow-600 dark:text-yellow-400 font-medium">Healing</div>
        </div>
      </div>

      <!-- Storage Pools -->
      <div class="pt-2 border-t border-gray-200 dark:border-gray-700">
        <div class="flex justify-between items-center text-sm">
          <span class="text-gray-500 dark:text-gray-400">Storage Pools</span>
          <span class="font-medium text-gray-900 dark:text-white">{{ dataUsage.poolsCount }}</span>
        </div>
      </div>
    </div>

    <!-- Fallback for no data -->
    <div v-else class="flex items-center justify-center h-48 bg-gray-50 dark:bg-gray-700 rounded-lg">
      <div class="text-center">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
        <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">No disk usage data available</p>
      </div>
    </div>
  </Card>
</template>