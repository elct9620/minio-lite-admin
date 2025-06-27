<script setup lang="ts">
import AppLayout from './components/layout/AppLayout.vue'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const pageTitle = computed(() => route.meta?.title as string || 'MinIO Lite Admin')
</script>

<template>
  <AppLayout :page-title="pageTitle">
    <router-view />
  </AppLayout>
</template>