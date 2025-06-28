<script setup lang="ts">
import Card from '../common/Card.vue'
import LoadingSpinner from '../common/LoadingSpinner.vue'
import ErrorMessage from '../common/ErrorMessage.vue'
import { useDataUsage, type DiskDetail } from '../../composables/useDataUsage'
import { formatBytes } from '../../utils/formatBytes'
import { InboxIcon } from '@heroicons/vue/24/outline'

const { dataUsage, loading, error } = useDataUsage()

function getDiskStatusColor(state: string): string {
  switch (state.toLowerCase()) {
    case 'ok':
    case 'online':
      return 'text-green-600 dark:text-green-400'
    case 'offline':
      return 'text-red-600 dark:text-red-400'
    case 'healing':
      return 'text-yellow-600 dark:text-yellow-400'
    default:
      return 'text-gray-600 dark:text-gray-400'
  }
}

function getDiskStatusBg(state: string): string {
  switch (state.toLowerCase()) {
    case 'ok':
    case 'online':
      return 'bg-green-50 dark:bg-green-900/20'
    case 'offline':
      return 'bg-red-50 dark:bg-red-900/20'
    case 'healing':
      return 'bg-yellow-50 dark:bg-yellow-900/20'
    default:
      return 'bg-gray-50 dark:bg-gray-900/20'
  }
}

function getUsagePercentage(disk: DiskDetail): number {
  if (disk.totalSpace === 0) return 0
  return (disk.usedSpace / disk.totalSpace) * 100
}

function getUsageColor(percentage: number): string {
  if (percentage >= 90) return 'bg-red-500'
  if (percentage >= 75) return 'bg-yellow-500'
  return 'bg-green-500'
}
</script>

<template>
  <Card title="Disk Details">
    <LoadingSpinner v-if="loading" text="Loading disk details..." />
    
    <ErrorMessage v-else-if="error" :message="error" />

    <div v-else-if="dataUsage?.diskDetails?.length" class="space-y-4">
      <div class="grid gap-4">
        <div 
          v-for="(disk, index) in dataUsage.diskDetails" 
          :key="`${disk.endpoint}-${disk.path}-${index}`"
          class="border border-gray-200 dark:border-gray-700 rounded-lg p-4"
        >
          <!-- Disk Header -->
          <div class="flex justify-between items-start mb-3">
            <div class="flex-1 min-w-0">
              <h4 class="text-sm font-medium text-gray-900 dark:text-white truncate">
                {{ disk.endpoint }}
              </h4>
              <p class="text-xs text-gray-500 dark:text-gray-400 truncate">
                {{ disk.path }}
              </p>
            </div>
            <div class="flex items-center space-x-2">
              <span :class="getDiskStatusBg(disk.state)" class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium">
                <span :class="getDiskStatusColor(disk.state)">{{ disk.state }}</span>
              </span>
              <span v-if="disk.rootDisk" class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400">
                Root
              </span>
            </div>
          </div>

          <!-- Usage Bar -->
          <div class="space-y-2">
            <div class="flex justify-between text-xs">
              <span class="text-gray-500 dark:text-gray-400">
                {{ formatBytes(disk.usedSpace) }} / {{ formatBytes(disk.totalSpace) }}
              </span>
              <span class="text-gray-900 dark:text-white font-medium">
                {{ getUsagePercentage(disk).toFixed(1) }}%
              </span>
            </div>
            <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
              <div 
                :class="getUsageColor(getUsagePercentage(disk))"
                class="h-2 rounded-full transition-all duration-300"
                :style="{ width: `${Math.min(getUsagePercentage(disk), 100)}%` }"
              ></div>
            </div>
          </div>

          <!-- Disk Info Grid -->
          <div class="grid grid-cols-2 md:grid-cols-4 gap-3 mt-3 pt-3 border-t border-gray-100 dark:border-gray-700">
            <div class="text-center">
              <div class="text-xs text-gray-500 dark:text-gray-400">Pool</div>
              <div class="text-sm font-medium text-gray-900 dark:text-white">{{ disk.pool }}</div>
            </div>
            <div class="text-center">
              <div class="text-xs text-gray-500 dark:text-gray-400">Set</div>
              <div class="text-sm font-medium text-gray-900 dark:text-white">{{ disk.set }}</div>
            </div>
            <div class="text-center">
              <div class="text-xs text-gray-500 dark:text-gray-400">FS Type</div>
              <div class="text-sm font-medium text-gray-900 dark:text-white">{{ disk.fsType || 'N/A' }}</div>
            </div>
            <div class="text-center">
              <div class="text-xs text-gray-500 dark:text-gray-400">Available</div>
              <div class="text-sm font-medium text-gray-900 dark:text-white">{{ formatBytes(disk.availSpace) }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="flex items-center justify-center h-32 bg-gray-50 dark:bg-gray-700 rounded-lg">
      <div class="text-center">
        <InboxIcon class="mx-auto h-8 w-8 text-gray-400" />
        <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">No disk details available</p>
      </div>
    </div>
  </Card>
</template>