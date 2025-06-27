<script setup lang="ts">
import { ref } from 'vue'
import AppSidebar from './AppSidebar.vue'
import AppTopBar from './AppTopBar.vue'

interface Props {
  pageTitle?: string
}

withDefaults(defineProps<Props>(), {
  pageTitle: 'Dashboard'
})

const sidebarOpen = ref(false)

const menuItems = [
  { 
    name: 'Dashboard', 
    icon: 'M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586l-7 7-7-7V4z', 
    path: '/dashboard' 
  },
  { 
    name: 'Access Keys', 
    icon: 'M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1721 9z', 
    path: '/access-keys' 
  },
  { 
    name: 'Site Replication', 
    icon: 'M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z', 
    path: '/site-replication' 
  }
]
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <AppSidebar 
      v-model:open="sidebarOpen" 
      :menu-items="menuItems" 
    />

    <!-- Main content area -->
    <div class="lg:pl-64">
      <AppTopBar 
        :title="pageTitle" 
        @toggle-sidebar="sidebarOpen = !sidebarOpen" 
      />

      <!-- Page content -->
      <main class="py-6 px-4 sm:px-6 lg:px-8">
        <slot />
      </main>
    </div>
  </div>
</template>