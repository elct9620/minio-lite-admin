<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { KeyIcon, CheckCircleIcon, BeakerIcon, ClockIcon } from '@heroicons/vue/24/outline'
import LoadingSpinner from '../components/common/LoadingSpinner.vue'
import ErrorMessage from '../components/common/ErrorMessage.vue'
import AccessKeyCard from '../components/dashboard/AccessKeyCard.vue'
import CreateAccessKeyModal from '../components/dashboard/CreateAccessKeyModal.vue'
import ConfirmDeleteModal from '../components/dashboard/ConfirmDeleteModal.vue'
import EditAccessKeyModal from '../components/dashboard/EditAccessKeyModal.vue'
import { useAccessKeys } from '../composables/useAccessKeys'

const { 
  accessKeys, 
  loading, 
  error, 
  total,
  fetchAccessKeys,
  userAccessKeys,
  serviceAccountKeys,
  stsKeys,
  enabledKeys,
  getStatusColor,
  getTypeDisplayName,
  getTypeColor
} = useAccessKeys()

const selectedFilter = ref<'all' | 'users' | 'serviceAccounts' | 'sts'>('all')
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const accessKeyToEdit = ref<typeof accessKeys.value[0] | null>(null)
const accessKeyToDelete = ref<{ accessKey: string; name?: string } | null>(null)

const filteredAccessKeys = computed(() => {
  switch (selectedFilter.value) {
    case 'users':
      return userAccessKeys.value
    case 'serviceAccounts':
      return serviceAccountKeys.value
    case 'sts':
      return stsKeys.value
    default:
      return accessKeys.value
  }
})

const handleFilterChange = async (filter: typeof selectedFilter.value) => {
  selectedFilter.value = filter
  await fetchAccessKeys({ type: filter })
}

const handleAccessKeyCreated = async () => {
  // Refresh the access keys list after creation
  await fetchAccessKeys({ type: selectedFilter.value })
}

const handleEditClick = (accessKey: typeof accessKeys.value[0]) => {
  accessKeyToEdit.value = accessKey
  showEditModal.value = true
}

const handleDeleteClick = (accessKey: typeof accessKeys.value[0]) => {
  accessKeyToDelete.value = {
    accessKey: accessKey.accessKey,
    name: accessKey.name
  }
  showDeleteModal.value = true
}

const handleEditConfirmed = async () => {
  // Refresh the access keys list after update
  await fetchAccessKeys({ type: selectedFilter.value })
  showEditModal.value = false
  accessKeyToEdit.value = null
}

const handleDeleteConfirmed = async () => {
  // Refresh the access keys list after deletion
  await fetchAccessKeys({ type: selectedFilter.value })
  showDeleteModal.value = false
  accessKeyToDelete.value = null
}

onMounted(() => {
  fetchAccessKeys()
})
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
        <button 
          @click="showCreateModal = true"
          class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors"
        >
          Create Access Key
        </button>
      </div>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-blue-100 dark:bg-blue-900 rounded-lg flex items-center justify-center">
              <KeyIcon class="w-5 h-5 text-blue-600 dark:text-blue-400" />
            </div>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400">Total Keys</p>
            <p class="text-2xl font-semibold text-gray-900 dark:text-white">{{ total }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-green-100 dark:bg-green-900 rounded-lg flex items-center justify-center">
              <CheckCircleIcon class="w-5 h-5 text-green-600 dark:text-green-400" />
            </div>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400">Enabled</p>
            <p class="text-2xl font-semibold text-gray-900 dark:text-white">{{ enabledKeys.length }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-purple-100 dark:bg-purple-900 rounded-lg flex items-center justify-center">
              <BeakerIcon class="w-5 h-5 text-purple-600 dark:text-purple-400" />
            </div>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400">Service Accounts</p>
            <p class="text-2xl font-semibold text-gray-900 dark:text-white">{{ serviceAccountKeys.length }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-orange-100 dark:bg-orange-900 rounded-lg flex items-center justify-center">
              <ClockIcon class="w-5 h-5 text-orange-600 dark:text-orange-400" />
            </div>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400">STS Tokens</p>
            <p class="text-2xl font-semibold text-gray-900 dark:text-white">{{ stsKeys.length }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Filter Tabs -->
    <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
      <div class="border-b border-gray-200 dark:border-gray-700">
        <nav class="-mb-px flex space-x-8 px-6" aria-label="Tabs">
          <button
            v-for="filter in [
              { key: 'all', label: 'All Keys', count: total },
              { key: 'users', label: 'Users', count: userAccessKeys.length },
              { key: 'serviceAccounts', label: 'Service Accounts', count: serviceAccountKeys.length },
              { key: 'sts', label: 'STS Tokens', count: stsKeys.length }
            ]"
            :key="filter.key"
            @click="handleFilterChange(filter.key as any)"
            :class="[
              'whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm transition-colors',
              selectedFilter === filter.key
                ? 'border-blue-500 text-blue-600 dark:text-blue-400'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'
            ]"
          >
            {{ filter.label }}
            <span class="ml-2 py-0.5 px-2 rounded-full text-xs bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-gray-100">
              {{ filter.count }}
            </span>
          </button>
        </nav>
      </div>

      <!-- Access Keys Content -->
      <div class="p-6">
        <LoadingSpinner v-if="loading" text="Loading access keys..." />
        
        <ErrorMessage v-else-if="error" :message="error" />

        <div v-else-if="filteredAccessKeys.length > 0" class="space-y-4">
          <AccessKeyCard
            v-for="accessKey in filteredAccessKeys"
            :key="accessKey.accessKey"
            :access-key="accessKey"
            :get-status-color="getStatusColor"
            :get-type-display-name="getTypeDisplayName"
            :get-type-color="getTypeColor"
            @edit="handleEditClick"
            @delete="handleDeleteClick"
          />
        </div>

        <div v-else class="flex items-center justify-center h-48 bg-gray-50 dark:bg-gray-700 rounded-lg">
          <div class="text-center">
            <KeyIcon class="mx-auto h-12 w-12 text-gray-400" />
            <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">
              {{ selectedFilter === 'all' ? 'No access keys found' : `No ${selectedFilter.replace(/([A-Z])/g, ' $1').toLowerCase()} found` }}
            </p>
            <p class="text-xs text-gray-400 dark:text-gray-500">
              {{ selectedFilter === 'all' ? 'Create your first access key to get started' : `Try a different filter or create a new ${selectedFilter.replace(/([A-Z])/g, ' $1').toLowerCase()}` }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Access Key Modal -->
    <CreateAccessKeyModal 
      v-model:open="showCreateModal" 
      @created="handleAccessKeyCreated"
    />

    <!-- Edit Access Key Modal -->
    <EditAccessKeyModal
      v-model:open="showEditModal"
      :access-key="accessKeyToEdit"
      @updated="handleEditConfirmed"
    />

    <!-- Delete Confirmation Modal -->
    <ConfirmDeleteModal
      v-if="accessKeyToDelete"
      v-model:open="showDeleteModal"
      :access-key="accessKeyToDelete.accessKey"
      :access-key-name="accessKeyToDelete.name"
      @confirmed="handleDeleteConfirmed"
    />
  </div>
</template>