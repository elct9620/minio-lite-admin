<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

interface MenuItem {
  name: string
  icon: any
  path: string
}

interface Props {
  open: boolean
  menuItems: MenuItem[]
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:open': [value: boolean]
}>()

const route = useRoute()

const isOpen = computed({
  get: () => props.open,
  set: (value) => emit('update:open', value)
})
</script>

<template>
  <!-- Mobile sidebar overlay -->
  <div v-if="isOpen" class="fixed inset-0 z-40 lg:hidden">
    <div class="fixed inset-0 bg-gray-600 bg-opacity-75" @click="isOpen = false"></div>
  </div>

  <!-- Sidebar -->
  <div 
    class="fixed inset-y-0 left-0 z-50 w-64 bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 transform transition-transform duration-300 ease-in-out lg:translate-x-0" 
    :class="isOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'"
  >
    <div class="flex items-center h-16 px-4 border-b border-gray-200 dark:border-gray-700">
      <h1 class="text-xl font-semibold text-gray-900 dark:text-white">
        MinIO Lite Admin
      </h1>
    </div>
    
    <nav class="mt-4 px-2">
      <div class="space-y-1">
        <router-link
          v-for="item in menuItems"
          :key="item.name"
          :to="item.path"
          class="group flex items-center px-2 py-2 text-sm font-medium rounded-md transition-colors"
          :class="route.path === item.path
            ? 'bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-200' 
            : 'text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:hover:text-white'"
          @click="isOpen = false"
        >
          <component :is="item.icon" class="mr-3 h-5 w-5" />
          {{ item.name }}
        </router-link>
      </div>
    </nav>
  </div>
</template>