<template>
  <div class="console-page-container space-y-8">
    <PageIntro :title="t('overview')" :description="t('order.overviewDesc')">
      <template #actions>
        <RefreshButton :loading="loading" @refresh="fetchData" />
      </template>
    </PageIntro>

    <OverviewFinanceCards
      :balance="balance"
      :last-month-settlement="lastMonthSettlement"
      :mtd-estimate="mtdEstimate"
      @recharge="handleRecharge"
      @view-transactions="goToTransactions"
    />

    <div class="grid grid-cols-1 gap-8 lg:grid-cols-3">
      <div class="lg:col-span-2">
        <OverviewTrendChart v-model:trend-period="trendPeriod" :chart-data="chartData" />
      </div>
      <OverviewRankingCard :items="rankingData" @view-usage="goToUsage" />
    </div>
  </div>
</template>

<script setup>
import { inject, onMounted } from 'vue'
import RefreshButton from '@/components/RefreshButton/index.vue'
import PageIntro from '@/components/listPage/PageIntro.vue'
import OverviewFinanceCards from './components/OverviewFinanceCards.vue'
import OverviewRankingCard from './components/OverviewRankingCard.vue'
import OverviewTrendChart from './components/OverviewTrendChart.vue'
import { useOrderOverview } from './composables/useOrderOverview'

const t = inject('t', (key) => key)

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
  fetchData()
})
</script>
