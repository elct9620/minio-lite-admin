<script setup lang="ts">
import type { AccessKeyInfo } from '../../composables/useAccessKeys'

interface Props {
  accessKey: AccessKeyInfo
  getStatusColor: (status: string) => string
  getTypeDisplayName: (type: string) => string
  getTypeColor: (type: string) => string
}

defineProps<Props>()

const formatDate = (dateString: string | undefined) => {
  if (!dateString) return null
  try {
    return new Date(dateString).toLocaleDateString()
  } catch {
    return dateString
  }
}
</script>

<template>
  <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4 hover:shadow-md transition-shadow">
    <div class="flex items-start justify-between">
      <div class="flex-1 min-w-0">
        <!-- Access Key ID -->
        <div class="flex items-center gap-2 mb-2">
          <h3 class="text-sm font-medium text-gray-900 dark:text-white truncate">
            {{ accessKey.accessKey }}
          </h3>
          <span 
            :class="getTypeColor(accessKey.type)"
            class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium"
          >
            {{ getTypeDisplayName(accessKey.type) }}
          </span>
        </div>

        <!-- Parent User -->
        <div class="text-sm text-gray-500 dark:text-gray-400 mb-2">
          <span class="font-medium">Parent User:</span> {{ accessKey.parentUser }}
        </div>

        <!-- Name and Description -->
        <div v-if="accessKey.name || accessKey.description" class="space-y-1 mb-2">
          <div v-if="accessKey.name" class="text-sm">
            <span class="font-medium text-gray-700 dark:text-gray-300">Name:</span>
            <span class="text-gray-900 dark:text-white ml-1">{{ accessKey.name }}</span>
          </div>
          <div v-if="accessKey.description" class="text-sm text-gray-600 dark:text-gray-400">
            {{ accessKey.description }}
          </div>
        </div>

        <!-- Metadata -->
        <div class="flex items-center gap-4 text-xs text-gray-500 dark:text-gray-400">
          <span 
            :class="getStatusColor(accessKey.accountStatus)"
            class="font-medium"
          >
            {{ accessKey.accountStatus.charAt(0).toUpperCase() + accessKey.accountStatus.slice(1) }}
          </span>
          
          <span v-if="accessKey.impliedPolicy" class="inline-flex items-center">
            <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Implied Policy
          </span>

          <span v-if="accessKey.expiration" class="inline-flex items-center text-yellow-600 dark:text-yellow-400">
            <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Expires: {{ formatDate(accessKey.expiration) }}
          </span>
        </div>
      </div>

      <!-- Actions -->
      <div class="flex items-center space-x-2 ml-4">
        <button
          class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
          title="View Details"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
          </svg>
        </button>
        
        <button
          v-if="accessKey.type === 'serviceAccount'"
          class="text-gray-400 hover:text-red-600 dark:hover:text-red-400 transition-colors"
          title="Delete Service Account"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>