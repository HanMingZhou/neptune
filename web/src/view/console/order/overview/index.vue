<template>
  <div class="console-page-container space-y-8">
    <BaseTableToolbar
      :description="t('order.overviewDesc')"
      :loading="loading"
      :title="t('overview')"
      @refresh="fetchData"
    />

    <OverviewFinanceCards
      :balance="balance"
      :last-month-settlement="lastMonthSettlement"
      :mtd-estimate="mtdEstimate"
      @recharge="handleRecharge"
      @view-transactions="goToTransactions"
    />

    <div class="grid grid-cols-1 gap-8 lg:grid-cols-3">
      <div class="lg:col-span-2">
        <OverviewTrendChart
          v-model:trend-period="trendPeriod"
          :chart-data="chartData"
        />
      </div>
      <OverviewRankingCard :items="rankingData" @view-usage="goToUsage" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineAsyncComponent, inject, onMounted } from 'vue'
import BaseTableToolbar from '@/components/listPage/BaseTableToolbar.vue'
import OverviewFinanceCards from './components/OverviewFinanceCards.vue'
import OverviewRankingCard from './components/OverviewRankingCard.vue'
import { useOrderOverview } from './composables/useOrderOverview'
import type { Translator } from '@/types/consoleResource'

const OverviewTrendChart = defineAsyncComponent(
  () => import('./components/OverviewTrendChart.vue')
)

const t = inject<Translator>('t', (key: string) => key)

const {
  balance,
  chartData,
  fetchData,
  goToTransactions,
  goToUsage,
  handleRecharge,
  lastMonthSettlement,
  loading,
  mtdEstimate,
  rankingData,
  trendPeriod
} = useOrderOverview({ t })

onMounted(() => {
  void fetchData()
})
</script>
