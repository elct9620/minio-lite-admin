<script setup lang="ts">
import ServerInfoCard from '../components/dashboard/ServerInfoCard.vue'
import DiskUsageChart from '../components/dashboard/DiskUsageChart.vue'
import StatisticsCard from '../components/dashboard/StatisticsCard.vue'
import SystemHealthCard from '../components/dashboard/SystemHealthCard.vue'
import DiskDetailsCard from '../components/dashboard/DiskDetailsCard.vue'
import { useServerInfo } from '../composables/useServerInfo'
import { useDataUsage } from '../composables/useDataUsage'
import { formatBytes } from '../utils/formatBytes'
import { computed } from 'vue'
import { ArchiveBoxIcon, DocumentIcon, ServerIcon, BoltIcon } from '@heroicons/vue/24/outline'

const { serverInfo, loading, error } = useServerInfo()
const { dataUsage, loading: dataLoading } = useDataUsage()

const statisticsData = computed(() => [
  {
    title: 'Total Buckets',
    value: dataUsage.value?.bucketsCount?.toLocaleString() || '-',
    icon: ArchiveBoxIcon,
    iconColor: 'text-blue-600'
  },
  {
    title: 'Total Objects',
    value: dataUsage.value?.objectsCount?.toLocaleString() || '-',
    icon: DocumentIcon,
    iconColor: 'text-green-600'
  },
  {
    title: 'Storage Used',
    value: dataUsage.value ? formatBytes(dataUsage.value.totalUsedCapacity) : '-',
    icon: ServerIcon,
    iconColor: 'text-yellow-600'
  },
  {
    title: 'Server Status',
    value: serverInfo.value ? 'Online' : (loading.value || dataLoading.value ? 'Checking...' : 'Offline'),
    icon: BoltIcon,
    iconColor: serverInfo.value ? 'text-green-600' : (loading.value || dataLoading.value ? 'text-yellow-600' : 'text-red-600')
  }
])
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
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
      <SystemHealthCard />
      <DiskDetailsCard />
    </div>
  </div>
</template>