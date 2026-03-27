<template>
  <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
    <div class="bg-white dark:bg-surface-dark border border-slate-200 dark:border-border-dark p-6 rounded-xl shadow-sm">
      <h3 class="text-[10px] font-black uppercase text-slate-400 mb-2">{{ t('order.pendingInvoiceAmount') }}</h3>
      <p class="text-sm font-medium text-slate-400 italic">{{ t('order.comingSoon') }}</p>
    </div>
    <div class="bg-white dark:bg-surface-dark border border-slate-200 dark:border-border-dark p-6 rounded-xl shadow-sm">
      <h3 class="text-[10px] font-black uppercase text-slate-400 mb-2">{{ t('order.invoiceTitleStatus') }}</h3>
      <p class="text-sm font-medium text-slate-400 italic">{{ t('order.comingSoon') }}</p>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

const t = inject('t', (key) => key)
</script>
