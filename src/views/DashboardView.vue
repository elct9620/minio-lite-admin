<script setup lang="ts">
import ServerInfoCard from '../components/dashboard/ServerInfoCard.vue'
import DiskUsageChart from '../components/dashboard/DiskUsageChart.vue'
import StatisticsCard from '../components/dashboard/StatisticsCard.vue'
import SystemHealthCard from '../components/dashboard/SystemHealthCard.vue'
import { useServerInfo } from '../composables/useServerInfo'

const { serverInfo, loading, error } = useServerInfo()

const statisticsData = [
  {
    title: 'Total Buckets',
    value: '-',
    icon: 'M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10',
    iconColor: 'text-blue-600'
  },
  {
    title: 'Total Objects',
    value: '-',
    icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z',
    iconColor: 'text-green-600'
  },
  {
    title: 'Storage Used',
    value: '-',
    icon: 'M7 4V2a1 1 0 011-1h8a1 1 0 011 1v2m-9 4v10a2 2 0 002 2h6a2 2 0 002-2V8M7 8h10M7 8V6a1 1 0 011-1h8a1 1 0 011 1v2',
    iconColor: 'text-yellow-600'
  },
  {
    title: 'Server Status',
    value: 'Online',
    icon: 'M13 10V3L4 14h7v7l9-11h-7z',
    iconColor: 'text-green-600'
  }
]
</script>

<template>
  <div>
    <!-- Dashboard Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Dashboard</h1>
      <p class="text-gray-600 dark:text-gray-400">MinIO server overview and status</p>
    </div>

    <!-- Server Status Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
      <ServerInfoCard 
        :loading="loading" 
        :error="error" 
        :server-info="serverInfo" 
      />
      <DiskUsageChart />
    </div>

    <!-- Statistics Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
      <StatisticsCard
        v-for="stat in statisticsData"
        :key="stat.title"
        :title="stat.title"
        :value="stat.value"
        :icon="stat.icon"
        :icon-color="stat.iconColor"
      />
    </div>

    <!-- System Health -->
    <SystemHealthCard />
  </div>
</template>