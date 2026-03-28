<template>
  <div class="console-page-container w-full space-y-8">
    <PageIntro :title="t('dashboardTitle')" :description="t('dashboardDesc')">
      <template #actions>
        <DashboardHeaderActions :loading="loading" @create="goToCreate" @refresh="fetchData" />
      </template>
    </PageIntro>

    <DashboardStatsGrid :items="statCards" />

    <div class="grid grid-cols-1 gap-8 xl:grid-cols-12">
      <div class="space-y-6 xl:col-span-8">
        <ResourceTrendCard :items="usageTrends" />
        <RecentInstancesCard
          :items="recentInstances"
          :get-status-class="getStatusClass"
          :get-type-class="getTypeClass"
          @open-detail="goToDetail"
          @view-all="goToNotebooks"
        />
      </div>

      <div class="space-y-6 xl:col-span-4">
        <QuickEntryCard :items="quickEntries" @select="handleQuickEntry" />
        <QuotaCard />
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject, onMounted } from 'vue'
import PageIntro from '@/components/listPage/PageIntro.vue'
import DashboardHeaderActions from './components/DashboardHeaderActions.vue'
import DashboardStatsGrid from './components/DashboardStatsGrid.vue'
import QuotaCard from './components/QuotaCard.vue'
import QuickEntryCard from './components/QuickEntryCard.vue'
import RecentInstancesCard from './components/RecentInstancesCard.vue'
import ResourceTrendCard from './components/ResourceTrendCard.vue'
import { useDashboardPage } from './composables/useDashboardPage'

const t = inject('t', (key) => key)

const {
  fetchData,
  getStatusClass,
  getTypeClass,
  goToCreate,
  goToDetail,
  goToNotebooks,
  handleQuickEntry,
  loading,
  quickEntries,
  recentInstances,
  statCards,
  usageTrends
} = useDashboardPage({ t })

onMounted(() => {
  fetchData()
})
</script>
