<template>
  <div class="max-w-[1400px] mx-auto space-y-8">
    <PageIntro :title="t('transactions')" :description="t('transactionsDesc')">
      <template #actions>
        <TransactionsHeaderActions :balance="balance" :loading="loading" @refresh="fetchData" />
      </template>
    </PageIntro>

    <div class="bg-white dark:bg-surface-dark border border-slate-200 dark:border-border-dark rounded-xl overflow-hidden shadow-sm">
      <TransactionsFiltersBar
        v-model:search-query="searchQuery"
        v-model:filter-type="filterType"
        v-model:date-range="dateRange"
        @date-change="handleDateChange"
        @search="handleSearch"
      />

      <TransactionsTableSection
        v-model:page="page"
        v-model:page-size="pageSize"
        :items="transactions"
        :total="total"
        @page-change="handlePageChange"
        @size-change="handleSizeChange"
      />
    </div>
  </div>
</template>

<script setup>
import { inject, onMounted } from 'vue'
import PageIntro from '@/components/listPage/PageIntro.vue'
import TransactionsFiltersBar from './components/TransactionsFiltersBar.vue'
import TransactionsHeaderActions from './components/TransactionsHeaderActions.vue'
import TransactionsTableSection from './components/TransactionsTableSection.vue'
import { useOrderTransactions } from './composables/useOrderTransactions'

const t = inject('t', (key) => key)

const {
  balance,
  dateRange,
  fetchData,
  filterType,
  handleDateChange,
  handlePageChange,
  handleSearch,
  handleSizeChange,
  loading,
  page,
  pageSize,
  searchQuery,
  total,
  transactions
} = useOrderTransactions({ t })

onMounted(() => {
  fetchData()
})
</script>
