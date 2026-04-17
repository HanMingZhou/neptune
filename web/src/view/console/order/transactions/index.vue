<template>
  <div class="console-page-container flex min-h-full flex-col gap-8">
    <BaseTableToolbar
      :description="t('transactionsDesc')"
      :show-refresh="false"
      :title="t('transactions')"
    >
      <template #actions>
        <TransactionsHeaderActions
          :balance="balance"
          :loading="loading"
          :can-recharge="canSystemRecharge"
          :disabled-reason="rechargeDisabledReason"
          @refresh="fetchData"
          @recharge="openRechargeDialog"
        />
      </template>
    </BaseTableToolbar>

    <div
      class="flex min-h-0 flex-1 flex-col bg-white dark:bg-surface-dark border border-slate-200 dark:border-border-dark rounded-xl overflow-hidden shadow-sm"
    >
      <TransactionsFiltersBar
        v-model:search-query="searchQuery"
        v-model:filter-type="filterType"
        v-model:date-range="dateRange"
        @date-change="handleDateChange"
        @search="handleSearch"
      />

      <TransactionsTableSection
        class="min-h-0 flex-1"
        v-model:page="page"
        v-model:page-size="pageSize"
        :items="transactions"
        :total="total"
        @page-change="handlePageChange"
        @size-change="handleSizeChange"
      />
    </div>

    <RechargeBalanceDialog
      v-model="rechargeDialogVisible"
      :form="rechargeForm"
      :rules="rechargeRules"
      :submitting="rechargeSubmitting"
      @closed="closeRechargeDialog"
      @submit="submitRecharge"
    />
  </div>
</template>

<script setup lang="ts">
import { inject, onMounted } from 'vue'
import BaseTableToolbar from '@/components/listPage/BaseTableToolbar.vue'
import RechargeBalanceDialog from './components/RechargeBalanceDialog.vue'
import TransactionsFiltersBar from './components/TransactionsFiltersBar.vue'
import TransactionsHeaderActions from './components/TransactionsHeaderActions.vue'
import TransactionsTableSection from './components/TransactionsTableSection.vue'
import { useOrderTransactions } from './composables/useOrderTransactions'
import type { Translator } from '@/types/consoleResource'

const t = inject<Translator>('t', (key: string) => key)

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
  canSystemRecharge,
  closeRechargeDialog,
  openRechargeDialog,
  rechargeDialogVisible,
  rechargeDisabledReason,
  rechargeForm,
  rechargeRules,
  rechargeSubmitting,
  searchQuery,
  submitRecharge,
  total,
  transactions
} = useOrderTransactions({ t })

onMounted(() => {
  void fetchData()
})
</script>
