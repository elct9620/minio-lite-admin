<script setup lang="ts">
import { ref } from 'vue'
import { XMarkIcon, ExclamationTriangleIcon } from '@heroicons/vue/24/outline'

interface Props {
  open: boolean
  accessKey: string
  accessKeyName?: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:open': [value: boolean]
  'confirmed': [accessKey: string]
}>()

// UI state
const loading = ref(false)
const error = ref<string | null>(null)

// Close modal
const closeModal = () => {
  if (!loading.value) {
    error.value = null
    emit('update:open', false)
  }
}

// Confirm deletion
const confirmDelete = async () => {
  loading.value = true
  error.value = null

  try {
    const response = await fetch(`/api/access-keys/${encodeURIComponent(props.accessKey)}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(errorText || 'Failed to delete service account')
    }

    // Emit success event
    emit('confirmed', props.accessKey)
    
    // Close modal
    closeModal()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'An unexpected error occurred'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <!-- Modal overlay -->
  <div v-if="open" class="fixed inset-0 z-50 overflow-y-auto">
    <div class="flex min-h-screen items-center justify-center p-4">
      <!-- Backdrop -->
      <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" @click="closeModal"></div>
      
      <!-- Modal panel -->
      <div class="relative w-full max-w-md transform overflow-hidden rounded-lg bg-white dark:bg-gray-800 shadow-xl transition-all">
        <!-- Header -->
        <div class="flex items-center justify-between border-b border-gray-200 dark:border-gray-700 px-6 py-4">
          <div class="flex items-center">
            <ExclamationTriangleIcon class="h-6 w-6 text-red-600 dark:text-red-400 mr-3" />
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              Delete Service Account
            </h3>
          </div>
          <button
            @click="closeModal"
            :disabled="loading"
            class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 disabled:opacity-50"
          >
            <XMarkIcon class="h-5 w-5" />
          </button>
        </div>

        <!-- Content -->
        <div class="px-6 py-4">
          <!-- Error message -->
          <div v-if="error" class="mb-4 rounded-md bg-red-50 dark:bg-red-900/50 p-4">
            <p class="text-sm text-red-800 dark:text-red-200">{{ error }}</p>
          </div>

          <!-- Warning message -->
          <div class="mb-6">
            <div class="rounded-md bg-yellow-50 dark:bg-yellow-900/50 p-4">
              <div class="flex">
                <ExclamationTriangleIcon class="h-5 w-5 text-yellow-400" />
                <div class="ml-3">
                  <h3 class="text-sm font-medium text-yellow-800 dark:text-yellow-200">
                    Warning: This action cannot be undone
                  </h3>
                  <p class="mt-1 text-sm text-yellow-700 dark:text-yellow-300">
                    Are you sure you want to delete the service account?
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- Service account details -->
          <div class="mb-6 rounded-md bg-gray-50 dark:bg-gray-700 p-4">
            <dl class="space-y-2">
              <div>
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Access Key</dt>
                <dd class="text-sm font-mono text-gray-900 dark:text-white break-all">{{ accessKey }}</dd>
              </div>
              <div v-if="accessKeyName">
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Name</dt>
                <dd class="text-sm text-gray-900 dark:text-white">{{ accessKeyName }}</dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Footer -->
        <div class="flex justify-end space-x-3 border-t border-gray-200 dark:border-gray-700 px-6 py-4">
          <button
            @click="closeModal"
            :disabled="loading"
            type="button"
            class="inline-flex items-center rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 shadow-sm hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50"
          >
            Cancel
          </button>
          <button
            @click="confirmDelete"
            :disabled="loading"
            type="button"
            class="inline-flex items-center rounded-md border border-transparent bg-red-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 disabled:opacity-50"
          >
            <span v-if="loading">Deleting...</span>
            <span v-else>Delete Service Account</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>