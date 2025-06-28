<script setup lang="ts">
import { ref, computed } from 'vue'
import { XMarkIcon, KeyIcon, EyeIcon, EyeSlashIcon } from '@heroicons/vue/24/outline'

interface Props {
  open: boolean
}

interface CreateAccessKeyRequest {
  name: string
  description: string
  accessKey: string
  secretKey: string
  policy: string
  targetUser: string
  expiration: string
}

interface CreateAccessKeyResponse {
  accessKey: string
  secretKey: string
  sessionToken?: string
  expiration?: string
  name?: string
  description?: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:open': [value: boolean]
  'created': [accessKey: CreateAccessKeyResponse]
}>()

// Form state
const form = ref<CreateAccessKeyRequest>({
  name: '',
  description: '',
  accessKey: '',
  secretKey: '',
  policy: '',
  targetUser: '',
  expiration: ''
})

// UI state
const loading = ref(false)
const error = ref<string | null>(null)
const showSecretKey = ref(false)
const showCustomKeys = ref(false)

// Modal control
const isOpen = computed({
  get: () => props.open,
  set: (value) => emit('update:open', value)
})

// Reset form when modal opens/closes
const resetForm = () => {
  form.value = {
    name: '',
    description: '',
    accessKey: '',
    secretKey: '',
    policy: '',
    targetUser: '',
    expiration: ''
  }
  error.value = null
  showSecretKey.value = false
  showCustomKeys.value = false
}

// Close modal
const closeModal = () => {
  isOpen.value = false
  resetForm()
}

// Create service account
const createServiceAccount = async () => {
  loading.value = true
  error.value = null

  try {
    // Prepare request payload
    const payload: Partial<CreateAccessKeyRequest> = {
      name: form.value.name.trim(),
      description: form.value.description.trim() || undefined,
      policy: form.value.policy.trim() || undefined,
      targetUser: form.value.targetUser.trim() || undefined,
      expiration: form.value.expiration || undefined
    }

    // Only include custom keys if specified
    if (showCustomKeys.value && form.value.accessKey.trim()) {
      payload.accessKey = form.value.accessKey.trim()
    }
    if (showCustomKeys.value && form.value.secretKey.trim()) {
      payload.secretKey = form.value.secretKey.trim()
    }

    const response = await fetch('/api/access-keys', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(payload)
    })

    if (!response.ok) {
      let errorMessage = 'Failed to create service account'
      try {
        const errorText = await response.text()
        if (errorText) {
          // Try to parse JSON error response first
          try {
            const errorJson = JSON.parse(errorText)
            errorMessage = errorJson.message || errorJson.error || errorText
          } catch {
            // If not JSON, use the text directly
            errorMessage = errorText
          }
        }
      } catch {
        // If reading response fails, use status text
        errorMessage = `Request failed with status ${response.status}: ${response.statusText}`
      }
      throw new Error(errorMessage)
    }

    const result: CreateAccessKeyResponse = await response.json()
    
    // Emit success event
    emit('created', result)
    
    // Close modal
    closeModal()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'An unexpected error occurred'
  } finally {
    loading.value = false
  }
}

// Generate random access key
const generateAccessKey = () => {
  // Generate a random access key (20 characters, alphanumeric)
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
  let result = ''
  for (let i = 0; i < 20; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  form.value.accessKey = result
}

// Generate random secret key
const generateSecretKey = () => {
  // Generate a random secret key (40 characters, base64-like)
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/'
  let result = ''
  for (let i = 0; i < 40; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  form.value.secretKey = result
}
</script>

<template>
  <!-- Modal overlay -->
  <div v-if="isOpen" class="fixed inset-0 z-50 overflow-y-auto">
    <div class="flex min-h-screen items-center justify-center p-4">
      <!-- Backdrop -->
      <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" @click="closeModal"></div>
      
      <!-- Modal panel -->
      <div class="relative w-full max-w-lg transform overflow-hidden rounded-lg bg-white dark:bg-gray-800 shadow-xl transition-all">
        <!-- Header -->
        <div class="flex items-center justify-between border-b border-gray-200 dark:border-gray-700 px-6 py-4">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
            Create Service Account
          </h3>
          <button
            @click="closeModal"
            class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
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

          <!-- Form -->
          <form @submit.prevent="createServiceAccount" class="space-y-4">
            <!-- Name (required) -->
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                Service Account Name
              </label>
              <input
                v-model="form.name"
                type="text"
                class="mt-1 block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-3 py-2 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                placeholder="Enter service account name (optional)"
              />
              <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                Optional name for the service account. If not provided, MinIO will generate one.
              </p>
            </div>

            <!-- Description -->
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                Description
              </label>
              <textarea
                v-model="form.description"
                rows="2"
                class="mt-1 block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-3 py-2 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                placeholder="Optional description"
              />
            </div>

            <!-- Custom Access/Secret Keys Toggle -->
            <div>
              <label class="flex items-center">
                <input
                  v-model="showCustomKeys"
                  type="checkbox"
                  class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                />
                <span class="ml-2 text-sm text-gray-700 dark:text-gray-300">
                  Specify custom access and secret keys
                </span>
              </label>
              <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                If unchecked, MinIO will generate random keys automatically
              </p>
            </div>

            <!-- Custom Keys Section -->
            <div v-if="showCustomKeys" class="space-y-4 border-l-4 border-blue-200 dark:border-blue-800 pl-4">
              <!-- Access Key -->
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                  Access Key
                </label>
                <div class="mt-1 flex rounded-md shadow-sm">
                  <input
                    v-model="form.accessKey"
                    type="text"
                    class="block w-full rounded-l-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-3 py-2 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                    placeholder="Leave empty for auto-generation"
                  />
                  <button
                    type="button"
                    @click="generateAccessKey"
                    class="inline-flex items-center rounded-r-md border border-l-0 border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-600 px-3 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-500"
                  >
                    <KeyIcon class="h-4 w-4" />
                  </button>
                </div>
              </div>

              <!-- Secret Key -->
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                  Secret Key
                </label>
                <div class="mt-1 flex rounded-md shadow-sm">
                  <input
                    v-model="form.secretKey"
                    :type="showSecretKey ? 'text' : 'password'"
                    class="block w-full rounded-l-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-3 py-2 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                    placeholder="Leave empty for auto-generation"
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

            <!-- Policy -->
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                Policy (JSON)
              </label>
              <textarea
                v-model="form.policy"
                rows="4"
                class="mt-1 block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-3 py-2 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 font-mono text-sm"
                placeholder='{"Version": "2012-10-17", "Statement": [...]}'
              />
              <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                Optional IAM policy in JSON format. Leave empty to inherit parent user permissions.
              </p>
            </div>

            <!-- Target User -->
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                Target User
              </label>
              <input
                v-model="form.targetUser"
                type="text"
                class="mt-1 block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-3 py-2 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                placeholder="Username this service account belongs to"
              />
            </div>

            <!-- Expiration -->
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                Expiration Date
              </label>
              <input
                v-model="form.expiration"
                type="datetime-local"
                class="mt-1 block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-3 py-2 text-gray-900 dark:text-white focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
              />
              <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                Optional expiration date and time. Leave empty for no expiration.
              </p>
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
            @click="createServiceAccount"
            :disabled="loading"
            class="px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center"
          >
            <span v-if="loading" class="mr-2">
              <svg class="animate-spin h-4 w-4" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none" />
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
              </svg>
            </span>
            {{ loading ? 'Creating...' : 'Create Service Account' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>