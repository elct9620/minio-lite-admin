import { ref, onMounted, readonly } from 'vue'

export interface DiskDetail {
  endpoint: string
  rootDisk: boolean
  path: string
  state: string
  uuid: string
  major: number
  minor: number
  totalSpace: number
  usedSpace: number
  availSpace: number
  fsType: string
  mount: string
  pool: number
  set: number
}

export interface DataUsage {
  totalCapacity: number
  totalFreeCapacity: number
  totalUsedCapacity: number
  usagePercentage: number
  onlineDisks: number
  offlineDisks: number
  healingDisks: number
  poolsCount: number
  objectsCount: number
  bucketsCount: number
  diskDetails: DiskDetail[]
}

export function useDataUsage() {
  const dataUsage = ref<DataUsage | null>(null)
  const loading = ref(true)
  const error = ref('')

  const fetchDataUsage = async () => {
    try {
      loading.value = true
      error.value = ''
      
      const response = await fetch('/api/data-usage')
      if (!response.ok) {
        throw new Error('Failed to fetch data usage')
      }
      
      dataUsage.value = await response.json()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Unknown error occurred'
    } finally {
      loading.value = false
    }
  }

  onMounted(fetchDataUsage)

  return {
    dataUsage: readonly(dataUsage),
    loading: readonly(loading),
    error: readonly(error),
    refetch: fetchDataUsage
  }
}