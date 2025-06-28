<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { XMarkIcon, KeyIcon, EyeIcon, EyeSlashIcon, ExclamationTriangleIcon } from '@heroicons/vue/24/outline'
import type { AccessKeyInfo } from '../../composables/useAccessKeys'

interface Props {
  open: boolean
  accessKey: AccessKeyInfo | null
}

interface UpdateAccessKeyRequest {
  newPolicy?: string
  newSecretKey?: string
  newStatus?: string
  newName?: string
  newDescription?: string
  newExpiration?: number
}

interface UpdateAccessKeyResponse {
  accessKey: string
  message: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:open': [value: boolean]
  'updated': [accessKey: AccessKeyInfo]
}>()

// Form state
const form = ref<UpdateAccessKeyRequest>({
  newPolicy: '',
  newSecretKey: '',
  newStatus: '',
  newName: '',
  newDescription: '',
  newExpiration: ''
})

// UI state
const loading = ref(false)
const error = ref<string | null>(null)
const showSecretKey = ref(false)
const showRotateSecret = ref(false)
const showAdvancedOptions = ref(false)

// Helper functions for datetime conversion
const convertFromTimestamp = (timestamp: number): string => {
  if (!timestamp) return ''
  try {
    const date = new Date(timestamp * 1000) // Convert from seconds to milliseconds
    // Convert to local datetime string in format YYYY-MM-DDTHH:MM
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    return `${year}-${month}-${day}T${hours}:${minutes}`
  } catch {
    return ''
  }
}

const convertToTimestamp = (datetimeLocal: string): number | undefined => {
  if (!datetimeLocal) return undefined
  // Convert datetime-local to timestamp (seconds)
  const date = new Date(datetimeLocal)
  return Math.floor(date.getTime() / 1000)
}

// Modal control
const isOpen = computed({
  get: () => props.open,
  set: (value) => emit('update:open', value)
})

// Watch for changes in accessKey prop to populate form
watch(() => props.accessKey, (newAccessKey) => {
  if (newAccessKey) {
    form.value = {
      newPolicy: '',
      newSecretKey: '',
      newStatus: newAccessKey.accountStatus,
      newName: newAccessKey.name || '',
      newDescription: newAccessKey.description || '',
      newExpiration: convertFromTimestamp(newAccessKey.expiration as any || 0)
    }
    // Reset UI state
    showSecretKey.value = false
    showRotateSecret.value = false
    showAdvancedOptions.value = false
  }
}, { immediate: true })

// Reset form when modal opens/closes
const resetForm = () => {
  form.value = {
    newPolicy: '',
    newSecretKey: '',
    newStatus: '',
    newName: '',
    newDescription: '',
    newExpiration: ''
  }
  error.value = null
  showSecretKey.value = false
  showRotateSecret.value = false
  showAdvancedOptions.value = false
}

// Close modal
const closeModal = () => {
  isOpen.value = false
  resetForm()
}

// Update service account
const updateServiceAccount = async () => {
  if (!props.accessKey) {
    error.value = 'No access key selected'
    return
  }

  loading.value = true
  error.value = null

  try {
    // Prepare request payload - only include fields that have been modified
    const payload: UpdateAccessKeyRequest = {}

    // Only include fields that have changed from original values
    if (form.value.newName?.trim() !== (props.accessKey.name || '')) {
      payload.newName = form.value.newName?.trim()
    }
    
    if (form.value.newDescription?.trim() !== (props.accessKey.description || '')) {
      payload.newDescription = form.value.newDescription?.trim()
    }
    
    if (form.value.newStatus !== props.accessKey.accountStatus) {
      payload.newStatus = form.value.newStatus
    }
    
    if (form.value.newPolicy?.trim()) {
      payload.newPolicy = form.value.newPolicy.trim()
    }
    
    if (showRotateSecret.value && form.value.newSecretKey?.trim()) {
      payload.newSecretKey = form.value.newSecretKey.trim()
    }
    
    if (form.value.newExpiration !== convertFromTimestamp(props.accessKey.expiration as any || 0)) {
      payload.newExpiration = convertToTimestamp(form.value.newExpiration)
    }

    // Only proceed if there are actual changes
    if (Object.keys(payload).length === 0) {
      error.value = 'No changes to save'
      return
    }

    const response = await fetch(`/api/access-keys/${encodeURIComponent(props.accessKey.accessKey)}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(payload)
    })

    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(errorText || 'Failed to update service account')
    }

    const result: UpdateAccessKeyResponse = await response.json()
    
    // Emit success event
    emit('updated', props.accessKey)
    
    // Close modal
    closeModal()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'An unexpected error occurred'
  } finally {
    loading.value = false
  }
}

// Generate random secret key
const generateSecretKey = () => {
  // Generate a random secret key (40 characters, base64-like)
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/'
  let result = ''
  for (let i = 0; i < 40; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  form.value.newSecretKey = result
}

// Validate JSON policy
const validatePolicy = (policy: string): boolean => {
  if (!policy.trim()) return true // Empty policy is valid
  try {
    JSON.parse(policy)
    return true
  } catch {
    return false
  }
}

// Policy validation
const isPolicyValid = computed(() => validatePolicy(form.value.newPolicy || ''))
</script>

<template>
  <!-- Modal overlay -->
  <div v-if="isOpen && accessKey" class="fixed inset-0 z-50 overflow-y-auto">
    <div class="flex min-h-screen items-center justify-center p-4">
      <!-- Backdrop -->
      <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" @click="closeModal"></div>
      
      <!-- Modal panel -->
      <div class="relative w-full max-w-2xl transform overflow-hidden rounded-lg bg-white dark:bg-gray-800 shadow-xl transition-all">
        <!-- Header -->
        <div class="flex items-center justify-between border-b border-gray-200 dark:border-gray-700 px-6 py-4">
          <div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              Edit Service Account
            </h3>
            <p class="text-sm text-gray-500 dark:text-gray-400 font-mono">
              {{ accessKey.accessKey }}
            </p>
          </div>
          <button
            @click="closeModal"
            class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
          >
            <XMarkIcon class="h-5 w-5" />
          </button>
        </div>

        <!-- Content -->
        <div class="px-6 py-4 max-h-96 overflow-y-auto">
          <!-- Error message -->
          <div v-if="error" class="mb-4 rounded-md bg-red-50 dark:bg-red-900/50 p-4">
            <p class="text-sm text-red-800 dark:text-red-200">{{ error }}</p>
          </div>

          <!-- Form -->
          <form @submit.prevent="updateServiceAccount" class="space-y-4">
            <!-- Basic Information -->
            <div class="space-y-4">
              <h4 class="text-sm font-medium text-gray-900 dark:text-white">Basic Information</h4>
              
              <!-- Name -->
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                  Name
                </label>
                <input
                  v-model="form.newName"
                  type="text"
                  class="mt-1 block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-3 py-2 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                  :placeholder="accessKey.name || 'Enter service account name'"
                />
              </div>

              <!-- Description -->
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                  Description
                </label>
                <textarea
                  v-model="form.newDescription"
                  rows="2"
                  class="mt-1 block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-3 py-2 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                  :placeholder="accessKey.description || 'Enter description'"
                />
              </div>

              <!-- Status -->
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                  Status
                </label>
                <select
                  v-model="form.newStatus"
                  class="mt-1 block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-3 py-2 text-gray-900 dark:text-white focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                >
                  <option value="enabled">Enabled</option>
                  <option value="disabled">Disabled</option>
                </select>
              </div>
            </div>

            <!-- Secret Key Rotation -->
            <div class="border-t border-gray-200 dark:border-gray-700 pt-4">
              <div class="flex items-center justify-between">
                <h4 class="text-sm font-medium text-gray-900 dark:text-white">Secret Key Rotation</h4>
                <button
                  type="button"
                  @click="showRotateSecret = !showRotateSecret"
                  class="text-sm text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300"
                >
                  {{ showRotateSecret ? 'Cancel Rotation' : 'Rotate Secret Key' }}
                </button>
              </div>
              
              <div v-if="showRotateSecret" class="mt-4 space-y-4 border-l-4 border-yellow-200 dark:border-yellow-800 pl-4">
                <div class="rounded-md bg-yellow-50 dark:bg-yellow-900/50 p-4">
                  <div class="flex">
                    <ExclamationTriangleIcon class="h-5 w-5 text-yellow-400" />
                    <div class="ml-3">
                      <h3 class="text-sm font-medium text-yellow-800 dark:text-yellow-200">
                        Warning: Secret Key Rotation
                      </h3>
                      <p class="mt-1 text-sm text-yellow-700 dark:text-yellow-300">
                        Rotating the secret key will invalidate the current credentials. Update your applications accordingly.
                      </p>
                    </div>
                  </div>
                </div>
                
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    New Secret Key
                  </label>
                  <div class="mt-1 flex rounded-md shadow-sm">
                    <input
                      v-model="form.newSecretKey"
                      :type="showSecretKey ? 'text' : 'password'"
                      class="block w-full rounded-l-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-3 py-2 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                      placeholder="Leave empty to auto-generate"
                    />
                    <button
                      type="button"
                      @click="generateSecretKey"
                      class="inline-flex items-center border border-l-0 border-r-0 border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-600 px-3 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-500"
                    >
                      <KeyIcon class="h-4 w-4" />
                    </button>
                    <button
                      type="button"
                      @click="showSecretKey = !showSecretKey"
                      class="inline-flex items-center rounded-r-md border border-l-0 border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-600 px-3 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-500"
                    >
                      <EyeIcon v-if="!showSecretKey" class="h-4 w-4" />
                      <EyeSlashIcon v-else class="h-4 w-4" />
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <!-- Advanced Options -->
            <div class="border-t border-gray-200 dark:border-gray-700 pt-4">
              <div class="flex items-center justify-between">
                <h4 class="text-sm font-medium text-gray-900 dark:text-white">Advanced Options</h4>
                <button
                  type="button"
                  @click="showAdvancedOptions = !showAdvancedOptions"
                  class="text-sm text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300"
                >
                  {{ showAdvancedOptions ? 'Hide Advanced' : 'Show Advanced' }}
                </button>
              </div>
              
              <div v-if="showAdvancedOptions" class="mt-4 space-y-4">
                <!-- Policy -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    IAM Policy (JSON)
                  </label>
                  <textarea
                    v-model="form.newPolicy"
                    rows="6"
                    :class="[
                      'mt-1 block w-full rounded-md border px-3 py-2 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:ring-1 font-mono text-sm',
                      isPolicyValid 
                        ? 'border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 focus:border-blue-500 focus:ring-blue-500'
                        : 'border-red-300 dark:border-red-600 bg-red-50 dark:bg-red-900/20 focus:border-red-500 focus:ring-red-500'
                    ]"
                    placeholder='{"Version": "2012-10-17", "Statement": [...]}'
                  />
                  <p v-if="!isPolicyValid" class="mt-1 text-sm text-red-600 dark:text-red-400">
                    Invalid JSON format
                  </p>
                  <p v-else class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                    Optional IAM policy in JSON format. Leave empty to inherit parent user permissions.
                  </p>
                </div>

                <!-- Expiration -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    Expiration Date
                  </label>
                  <input
                    v-model="form.newExpiration"
                    type="datetime-local"
                    class="mt-1 block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-3 py-2 text-gray-900 dark:text-white focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                  />
                  <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                    Optional expiration date and time. Leave empty for no expiration.
                  </p>
                </div>
              </div>
            </div>
          </form>
        </div>

        <!-- Footer -->
        <div class="flex justify-end space-x-3 border-t border-gray-200 dark:border-gray-700 px-6 py-4">
          <button
            type="button"
            @click="closeModal"
            class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white"
          >
            Cancel
          </button>
          <button
            @click="updateServiceAccount"
            :disabled="loading || !isPolicyValid"
            class="px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center"
          >
            <span v-if="loading" class="mr-2">
              <svg class="animate-spin h-4 w-4" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none" />
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
              </svg>
            </span>
            {{ loading ? 'Updating...' : 'Update Service Account' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>