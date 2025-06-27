<script setup lang="ts">
import { ref, onMounted } from 'vue'

const serverInfo = ref(null)
const loading = ref(true)
const error = ref('')
const sidebarOpen = ref(false)

const menuItems = [
  { name: 'Dashboard', icon: 'M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586l-7 7-7-7V4z', current: true },
  { name: 'Access Keys', icon: 'M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z', current: false },
  { name: 'Site Replication', icon: 'M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z', current: false },
]

onMounted(async () => {
  try {
    const response = await fetch('/api/server-info')
    if (!response.ok) {
      throw new Error('Failed to fetch server info')
    }
    serverInfo.value = await response.json()
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <!-- Mobile sidebar overlay -->
    <div v-if="sidebarOpen" class="fixed inset-0 z-40 lg:hidden">
      <div class="fixed inset-0 bg-gray-600 bg-opacity-75" @click="sidebarOpen = false"></div>
    </div>

    <!-- Sidebar -->
    <div class="fixed inset-y-0 left-0 z-50 w-64 bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 transform transition-transform duration-300 ease-in-out lg:translate-x-0" :class="sidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'">
      <div class="flex items-center h-16 px-4 border-b border-gray-200 dark:border-gray-700">
        <h1 class="text-xl font-semibold text-gray-900 dark:text-white">
          MinIO Lite Admin
        </h1>
      </div>
      
      <nav class="mt-4 px-2">
        <div class="space-y-1">
          <a
            v-for="item in menuItems"
            :key="item.name"
            href="#"
            class="group flex items-center px-2 py-2 text-sm font-medium rounded-md transition-colors"
            :class="item.current 
              ? 'bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-200' 
              : 'text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:hover:text-white'"
          >
            <svg class="mr-3 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="item.icon" />
            </svg>
            {{ item.name }}
          </a>
        </div>
      </nav>
    </div>

    <!-- Main content area -->
    <div class="lg:pl-64">
      <!-- Top navigation -->
      <div class="sticky top-0 z-10 bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 lg:hidden">
        <div class="flex items-center justify-between h-16 px-4">
          <button
            @click="sidebarOpen = !sidebarOpen"
            class="text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white"
          >
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
          <h1 class="text-lg font-semibold text-gray-900 dark:text-white">Dashboard</h1>
          <div></div>
        </div>
      </div>

      <!-- Page content -->
      <main class="py-6 px-4 sm:px-6 lg:px-8">
        <!-- Dashboard Header -->
        <div class="mb-6">
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Dashboard</h1>
          <p class="text-gray-600 dark:text-gray-400">MinIO server overview and status</p>
        </div>

        <!-- Server Status Grid -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
          <!-- Server Info Card -->
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
            <h2 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Server Information</h2>
            
            <div v-if="loading" class="flex items-center justify-center py-8">
              <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            </div>

            <div v-else-if="error" class="text-red-600 dark:text-red-400 py-4">
              Error: {{ error }}
            </div>

            <div v-else-if="serverInfo" class="space-y-4">
              <div class="flex justify-between items-center py-2 border-b border-gray-100 dark:border-gray-700">
                <span class="text-sm font-medium text-gray-500 dark:text-gray-400">Mode</span>
                <span class="text-sm font-semibold text-gray-900 dark:text-white">{{ serverInfo.mode }}</span>
              </div>
              <div class="flex justify-between items-center py-2 border-b border-gray-100 dark:border-gray-700">
                <span class="text-sm font-medium text-gray-500 dark:text-gray-400">Region</span>
                <span class="text-sm font-semibold text-gray-900 dark:text-white">{{ serverInfo.region || 'Not set' }}</span>
              </div>
              <div class="flex justify-between items-center py-2">
                <span class="text-sm font-medium text-gray-500 dark:text-gray-400">Deployment ID</span>
                <span class="text-xs font-mono text-gray-900 dark:text-white break-all max-w-40">{{ serverInfo.deploymentId }}</span>
              </div>
            </div>
          </div>

          <!-- Disk Usage Chart Placeholder -->
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
            <h2 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Disk Usage</h2>
            <div class="flex items-center justify-center h-48 bg-gray-50 dark:bg-gray-700 rounded-lg">
              <div class="text-center">
                <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                </svg>
                <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">Disk usage chart</p>
                <p class="text-xs text-gray-400 dark:text-gray-500">Coming soon</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Statistics Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <svg class="h-8 w-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">Total Buckets</dt>
                  <dd class="text-lg font-semibold text-gray-900 dark:text-white">-</dd>
                </dl>
              </div>
            </div>
          </div>

          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <svg class="h-8 w-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">Total Objects</dt>
                  <dd class="text-lg font-semibold text-gray-900 dark:text-white">-</dd>
                </dl>
              </div>
            </div>
          </div>

          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <svg class="h-8 w-8 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 4V2a1 1 0 011-1h8a1 1 0 011 1v2m-9 4v10a2 2 0 002 2h6a2 2 0 002-2V8M7 8h10M7 8V6a1 1 0 011-1h8a1 1 0 011 1v2" />
                </svg>
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">Storage Used</dt>
                  <dd class="text-lg font-semibold text-gray-900 dark:text-white">-</dd>
                </dl>
              </div>
            </div>
          </div>

          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <svg class="h-8 w-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">Server Status</dt>
                  <dd class="text-lg font-semibold text-green-600">Online</dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <!-- Recent Activity / Health Check -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
          <h2 class="text-lg font-medium text-gray-900 dark:text-white mb-4">System Health</h2>
          <div class="flex items-center justify-center h-32 bg-gray-50 dark:bg-gray-700 rounded-lg">
            <div class="text-center">
              <svg class="mx-auto h-8 w-8 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <p class="mt-2 text-sm text-gray-900 dark:text-white font-medium">All systems operational</p>
              <p class="text-xs text-gray-500 dark:text-gray-400">Last checked: Just now</p>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>
