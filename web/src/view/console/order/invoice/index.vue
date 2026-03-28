<template>
  <div class="console-page-container space-y-8">
    <PageIntro :title="t('order.invoice')" :description="t('order.invoiceDesc')">
      <template #actions>
        <RefreshButton :loading="loading" @refresh="fetchData" />
        <button
          disabled
          class="bg-slate-300 dark:bg-zinc-600 text-white px-6 py-3 rounded-xl font-bold text-sm flex items-center gap-2 cursor-not-allowed opacity-60"
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
