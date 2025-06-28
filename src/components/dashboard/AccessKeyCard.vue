<script setup lang="ts">
import { CheckCircleIcon, ClockIcon, TrashIcon, PencilIcon } from '@heroicons/vue/24/outline'
import type { AccessKeyInfo } from '../../composables/useAccessKeys'

interface Props {
  accessKey: AccessKeyInfo
  getStatusColor: (status: string) => string
  getTypeDisplayName: (type: string) => string
  getTypeColor: (type: string) => string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'delete': [accessKey: AccessKeyInfo]
  'edit': [accessKey: AccessKeyInfo]
}>()

const formatDate = (dateString: string | undefined) => {
  if (!dateString) return null
  try {
    const date = new Date(dateString)
    // Check if the date is the Unix epoch (1970-01-01) or invalid
    if (date.getTime() <= 0 || date.getFullYear() === 1970) {
      return null
    }
    return date.toLocaleDateString()
  } catch {
    return null
  }
}

const formatExpiration = (dateString: string | undefined) => {
  const formattedDate = formatDate(dateString)
  return formattedDate || 'Never expires'
}

const handleDeleteClick = () => {
  emit('delete', props.accessKey)
}

const handleEditClick = () => {
  emit('edit', props.accessKey)
}
</script>

<template>
  <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4 hover:shadow-md transition-shadow">
    <div class="flex items-start justify-between">
      <div class="flex-1 min-w-0">
        <!-- Primary Display: Name or Access Key -->
        <div class="flex items-center gap-2 mb-2">
          <h3 class="text-sm font-medium text-gray-900 dark:text-white truncate">
            {{ accessKey.name || accessKey.accessKey }}
          </h3>
          <span 
            :class="getTypeColor(accessKey.type)"
            class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium"
          >
            {{ getTypeDisplayName(accessKey.type) }}
          </span>
        </div>

        <!-- Access Key ID (if name exists) -->
        <div v-if="accessKey.name" class="text-sm text-gray-500 dark:text-gray-400 mb-2">
          <span class="font-medium">Access Key:</span> 
          <span class="font-mono text-xs">{{ accessKey.accessKey }}</span>
        </div>

        <!-- Parent User (only for service accounts and STS tokens) -->
        <div v-if="accessKey.type !== 'user'" class="text-sm text-gray-500 dark:text-gray-400 mb-2">
          <span class="font-medium">Parent User:</span> {{ accessKey.parentUser }}
        </div>

        <!-- Description -->
        <div v-if="accessKey.description" class="text-sm text-gray-600 dark:text-gray-400 mb-2">
          {{ accessKey.description }}
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
            <CheckCircleIcon class="w-3 h-3 mr-1" />
            Implied Policy
          </span>

          <span class="inline-flex items-center" :class="formatDate(accessKey.expiration) ? 'text-yellow-600 dark:text-yellow-400' : 'text-gray-500 dark:text-gray-400'">
            <ClockIcon class="w-3 h-3 mr-1" />
            {{ formatExpiration(accessKey.expiration) }}
          </span>
        </div>
      </div>

      <!-- Actions -->
      <div class="flex items-center space-x-2 ml-4">
        <button
          v-if="accessKey.type === 'serviceAccount'"
          @click="handleEditClick"
          class="text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
          title="Edit Service Account"
        >
          <PencilIcon class="w-5 h-5" />
        </button>
        
        <button
          v-if="accessKey.type === 'serviceAccount'"
          @click="handleDeleteClick"
          class="text-gray-400 hover:text-red-600 dark:hover:text-red-400 transition-colors"
          title="Delete Service Account"
        >
          <TrashIcon class="w-5 h-5" />
        </button>
      </div>
    </div>
  </div>
</template>