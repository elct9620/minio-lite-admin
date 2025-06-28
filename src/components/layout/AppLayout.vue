<script setup lang="ts">
import { ref } from 'vue'
import { HomeIcon, KeyIcon, ArrowsRightLeftIcon } from '@heroicons/vue/24/outline'
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
    icon: HomeIcon, 
    path: '/dashboard' 
  },
  { 
    name: 'Access Keys', 
    icon: KeyIcon, 
    path: '/access-keys' 
  },
  { 
    name: 'Site Replication', 
    icon: ArrowsRightLeftIcon, 
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