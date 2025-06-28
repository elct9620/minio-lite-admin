<script setup lang="ts">
import { ref, computed } from 'vue'
import { XMarkIcon, CheckIcon, ClipboardIcon, ExclamationTriangleIcon } from '@heroicons/vue/24/outline'

interface Props {
  open: boolean
  accessKey: string
  secretKey: string
  title?: string
  description?: string
  isRotation?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  title: 'Access Key Created',
  description: 'Your new access key has been created successfully. Please save these credentials securely.',
  isRotation: false
})

const emit = defineEmits<{
  'update:open': [value: boolean]
  'closed': []
}>()

// UI state
const copiedAccessKey = ref(false)
const copiedSecretKey = ref(false)

// Modal control
const isOpen = computed({
  get: () => props.open,
  set: (value) => emit('update:open', value)
})

// Copy to clipboard functionality
const copyToClipboard = async (text: string, type: 'accessKey' | 'secretKey') => {
  try {
    await navigator.clipboard.writeText(text)
    if (type === 'accessKey') {
      copiedAccessKey.value = true
      setTimeout(() => { copiedAccessKey.value = false }, 2000)
    } else {
      copiedSecretKey.value = true
      setTimeout(() => { copiedSecretKey.value = false }, 2000)
    }
  } catch (err) {
    console.error('Failed to copy to clipboard:', err)
  }
}

// Close modal
const closeModal = () => {
  isOpen.value = false
  emit('closed')
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
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <div class="w-8 h-8 bg-green-100 dark:bg-green-900 rounded-lg flex items-center justify-center">
                <CheckIcon class="w-5 h-5 text-green-600 dark:text-green-400" />
              </div>
            </div>
            <div class="ml-3">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
                {{ title }}
              </h3>
            </div>
          </div>
          <button
            @click="closeModal"
            class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
          >
            <XMarkIcon class="h-5 w-5" />
          </button>
        </div>

        <!-- Content -->
        <div class="px-6 py-4">
          <!-- Warning message -->
          <div class="mb-6 rounded-md bg-yellow-50 dark:bg-yellow-900/50 p-4">
            <div class="flex">
              <ExclamationTriangleIcon class="h-5 w-5 text-yellow-400 flex-shrink-0" />
              <div class="ml-3">
                <h3 class="text-sm font-medium text-yellow-800 dark:text-yellow-200">
                  Important Security Notice
                </h3>
                <p class="mt-1 text-sm text-yellow-700 dark:text-yellow-300">
                  {{ description }} This is the only time the secret key will be displayed. 
                  {{ isRotation ? 'The previous secret key is now invalid.' : '' }}
                </p>
              </div>
            </div>
          </div>

          <!-- Credentials -->
          <div class="space-y-4">
            <!-- Access Key -->
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Access Key
              </label>
              <div class="flex rounded-md shadow-sm">
                <input
                  :value="accessKey"
                  type="text"
                  readonly
                  class="block w-full rounded-l-md border border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 px-3 py-2 text-gray-900 dark:text-white font-mono text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                />
                <button
                  type="button"
                  @click="copyToClipboard(accessKey, 'accessKey')"
                  class="inline-flex items-center rounded-r-md border border-l-0 border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-600 px-3 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-500 transition-colors"
                >
                  <CheckIcon v-if="copiedAccessKey" class="h-4 w-4 text-green-600" />
                  <ClipboardIcon v-else class="h-4 w-4" />
                </button>
              </div>
            </div>

            <!-- Secret Key -->
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Secret Key
              </label>
              <div class="flex rounded-md shadow-sm">
                <input
                  :value="secretKey"
                  type="text"
                  readonly
                  class="block w-full rounded-l-md border border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 px-3 py-2 text-gray-900 dark:text-white font-mono text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                />
                <button
                  type="button"
                  @click="copyToClipboard(secretKey, 'secretKey')"
                  class="inline-flex items-center rounded-r-md border border-l-0 border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-600 px-3 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-500 transition-colors"
                >
                  <CheckIcon v-if="copiedSecretKey" class="h-4 w-4 text-green-600" />
                  <ClipboardIcon v-else class="h-4 w-4" />
                </button>
              </div>
            </div>
          </div>

          <!-- Instructions -->
          <div class="mt-6 p-4 bg-blue-50 dark:bg-blue-900/50 rounded-md">
            <h4 class="text-sm font-medium text-blue-800 dark:text-blue-200 mb-2">
              Next Steps
            </h4>
            <ul class="text-sm text-blue-700 dark:text-blue-300 space-y-1">
              <li>• Save these credentials in a secure location</li>
              <li>• Configure your MinIO client with these credentials</li>
              <li>• Test the connection to ensure it works properly</li>
              <li v-if="isRotation">• Update any applications using the old secret key</li>
            </ul>
          </div>
        </div>

        <!-- Footer -->
        <div class="flex justify-end border-t border-gray-200 dark:border-gray-700 px-6 py-4">
          <button
            @click="closeModal"
            class="px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700 transition-colors"
          >
            I've Saved the Credentials
          </button>
        </div>
      </div>
    </div>
  </div>
</template>