<template>
  <div class="max-w-[1400px] mx-auto space-y-8">
    <PageIntro :title="t('order.usage')" :description="t('order.usageDesc')">
      <template #actions>
        <RefreshButton :loading="loading" @refresh="fetchData" />
        <button disabled class="px-4 py-2 bg-slate-300 dark:bg-zinc-600 text-white rounded-lg text-xs font-bold flex items-center gap-2 cursor-not-allowed opacity-60">
          <span class="material-icons text-sm">file_download</span>
          {{ t('order.exportOrder') }}
        </button>
      </template>
    </PageIntro>

    <div class="bg-white dark:bg-surface-dark border border-slate-200 dark:border-border-dark rounded-2xl overflow-hidden shadow-sm">
      <UsageFiltersBar
        v-model:search-query="searchQuery"
        v-model:charge-filter="chargeFilter"
        v-model:date-range="dateRange"
        @date-change="handleDateChange"
        @search="handleSearch"
      />

      <UsageSummaryBar :main-item="mainItem" :total-expenditure="totalExpenditure" />

      <UsageTableSection
        v-model:page="page"
        v-model:page-size="pageSize"
        :charge-filter="chargeFilter"
        :get-charge-type-style="getChargeTypeStyle"
        :get-resource-type-text="getResourceTypeText"
        :get-status-style="getStatusStyle"
        :get-unit-text="getUnitText"
        :items="filteredUsageItems"
        :format-usage="formatUsage"
        :total="total"
        @page-change="handlePageChange"
        @size-change="handleSizeChange"
      />
    </div>
  </div>
</template>

<script setup>
import { inject, onMounted } from 'vue'
import RefreshButton from '@/components/RefreshButton/index.vue'
import PageIntro from '@/components/listPage/PageIntro.vue'
import UsageFiltersBar from './components/UsageFiltersBar.vue'
import UsageSummaryBar from './components/UsageSummaryBar.vue'
import UsageTableSection from './components/UsageTableSection.vue'
import { useOrderUsage } from './composables/useOrderUsage'

const t = inject('t', (key) => key)

const {
  chargeFilter,
  dateRange,
  fetchData,
  filteredUsageItems,
  getChargeTypeStyle,
  getResourceTypeText,
  getStatusStyle,
  getUnitText,
  formatUsage,
  handleDateChange,
  handlePageChange,
  handleSearch,
  handleSizeChange,
  loading,
  mainItem,
  page,
  pageSize,
  searchQuery,
  total,
  totalExpenditure
} = useOrderUsage({ t })

onMounted(() => {
  fetchData()
})
</script>
