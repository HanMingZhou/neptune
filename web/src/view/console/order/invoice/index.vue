<template>
  <div class="console-page-container space-y-8">
    <PageIntro :title="t('order.invoice')" :description="t('order.invoiceDesc')">
      <template #actions>
        <RefreshButton :loading="loading" @refresh="fetchData" />
        <button
          disabled
          class="list-toolbar-button list-toolbar-button--secondary"
        >
          <span class="material-icons">add_circle</span>
          {{ t('order.applyInvoice') }}
        </button>
      </template>
    </PageIntro>

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

<script setup>
import { inject, onMounted } from 'vue'
import RefreshButton from '@/components/RefreshButton/index.vue'
import PageIntro from '@/components/listPage/PageIntro.vue'
import ApplyInvoiceDialog from './components/ApplyInvoiceDialog.vue'
import InvoiceStatsCards from './components/InvoiceStatsCards.vue'
import InvoiceTabsCard from './components/InvoiceTabsCard.vue'
import { useInvoicePage } from './composables/useInvoicePage'

const t = inject('t', (key) => key)

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
  fetchData()
})
</script>
