import { ref, onMounted, readonly } from 'vue'

interface ServerInfo {
  mode: string
  region?: string
  deploymentId: string
}

export function useServerInfo() {
  const serverInfo = ref<ServerInfo | null>(null)
  const loading = ref(true)
  const error = ref('')

  const fetchServerInfo = async () => {
    try {
      loading.value = true
      error.value = ''
      
      const response = await fetch('/api/server-info')
      if (!response.ok) {
        throw new Error('Failed to fetch server info')
      }
      
      serverInfo.value = await response.json()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Unknown error occurred'
    } finally {
      loading.value = false
    }
  }

  onMounted(fetchServerInfo)

  return {
    serverInfo: readonly(serverInfo),
    loading: readonly(loading),
    error: readonly(error),
    refetch: fetchServerInfo
  }
}