<template>
  <div class="console-page-container space-y-8">
    <BaseTableToolbar
      :description="t('order.invoiceDesc')"
      :loading="loading"
      :title="t('order.invoice')"
      @refresh="fetchData"
    >
      <template #actions>
        <button
          disabled
          class="list-toolbar-button list-toolbar-button--secondary"
        >
          <span class="material-icons">add_circle</span>
          {{ t('order.applyInvoice') }}
        </button>
      </template>
    </BaseTableToolbar>

    <InvoiceStatsCards />

    <InvoiceTabsCard v-model:active-tab="activeTab" :invoices="invoices" />

    <ApplyInvoiceDialog
      v-model="showApplyDialog"
      :form="applyForm"
      :rules="applyRules"
      :submitting="submitting"
      @closed="resetApplyForm"
      @submit="submitApply"
    />
  </div>
</template>

<script setup lang="ts">
import { inject, onMounted } from 'vue'
import BaseTableToolbar from '@/components/listPage/BaseTableToolbar.vue'
import ApplyInvoiceDialog from './components/ApplyInvoiceDialog.vue'
import InvoiceStatsCards from './components/InvoiceStatsCards.vue'
import InvoiceTabsCard from './components/InvoiceTabsCard.vue'
import { useInvoicePage } from './composables/useInvoicePage'
import type { Translator } from '@/types/consoleResource'

const t = inject<Translator>('t', (key: string) => key)

const {
  activeTab,
  applyForm,
  applyRules,
  fetchData,
  invoices,
  loading,
  resetApplyForm,
  showApplyDialog,
  submitApply,
  submitting
} = useInvoicePage({ t })

onMounted(() => {
  void fetchData()
})
</script>
