import { ref, computed } from 'vue'

export interface AccessKeyInfo {
  accessKey: string
  parentUser: string
  accountStatus: string
  type: 'user' | 'serviceAccount' | 'sts'
  name?: string
  description?: string
  expiration?: string
  createdAt?: string
  impliedPolicy: boolean
}

export interface AccessKeysResponse {
  accessKeys: AccessKeyInfo[]
  total: number
}

export interface AccessKeysOptions {
  type?: 'all' | 'users' | 'serviceAccounts' | 'sts'
  user?: string
}

export function useAccessKeys() {
  const accessKeys = ref<AccessKeyInfo[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const total = ref(0)

  const fetchAccessKeys = async (options: AccessKeysOptions = {}) => {
    loading.value = true
    error.value = null
    
    try {
      const params = new URLSearchParams()
      if (options.type && options.type !== 'all') {
        params.append('type', options.type)
      }
      if (options.user) {
        params.append('user', options.user)
      }

      const url = `/api/access-keys${params.toString() ? `?${params.toString()}` : ''}`
      const response = await fetch(url)
      
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`)
      }
      
      const data: AccessKeysResponse = await response.json()
      accessKeys.value = data.accessKeys
      total.value = data.total
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Unknown error occurred'
      accessKeys.value = []
      total.value = 0
    } finally {
      loading.value = false
    }
  }

  // Computed properties for filtered access keys
  const userAccessKeys = computed(() => 
    accessKeys.value.filter(key => key.type === 'user')
  )

  const serviceAccountKeys = computed(() => 
    accessKeys.value.filter(key => key.type === 'serviceAccount')
  )

  const stsKeys = computed(() => 
    accessKeys.value.filter(key => key.type === 'sts')
  )

  const enabledKeys = computed(() => 
    accessKeys.value.filter(key => key.accountStatus === 'enabled')
  )

  const disabledKeys = computed(() => 
    accessKeys.value.filter(key => key.accountStatus === 'disabled')
  )

  // Helper function to get status color
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'enabled':
        return 'text-green-600 dark:text-green-400'
      case 'disabled':
        return 'text-red-600 dark:text-red-400'
      default:
        return 'text-gray-600 dark:text-gray-400'
    }
  }

  // Helper function to get type display name
  const getTypeDisplayName = (type: string) => {
    switch (type) {
      case 'user':
        return 'User'
      case 'serviceAccount':
        return 'Service Account'
      case 'sts':
        return 'STS Token'
      default:
        return type
    }
  }

  // Helper function to get type color
  const getTypeColor = (type: string) => {
    switch (type) {
      case 'user':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200'
      case 'serviceAccount':
        return 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200'
      case 'sts':
        return 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200'
      default:
        return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200'
    }
  }

  return {
    accessKeys,
    loading,
    error,
    total,
    fetchAccessKeys,
    userAccessKeys,
    serviceAccountKeys,
    stsKeys,
    enabledKeys,
    disabledKeys,
    getStatusColor,
    getTypeDisplayName,
    getTypeColor
  }
}