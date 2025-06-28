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
    return new Date(dateString).toLocaleDateString()
  } catch {
    return dateString
  }
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
            <CheckCircleIcon class="w-3 h-3 mr-1" />
            Implied Policy
          </span>

          <span v-if="accessKey.expiration" class="inline-flex items-center text-yellow-600 dark:text-yellow-400">
            <ClockIcon class="w-3 h-3 mr-1" />
            Expires: {{ formatDate(accessKey.expiration) }}
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