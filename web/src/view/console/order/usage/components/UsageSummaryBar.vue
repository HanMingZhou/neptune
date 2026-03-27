<template>
  <div class="p-6 border-b border-slate-100 dark:border-border-dark flex items-center gap-6">
    <div class="flex-1">
      <h3 class="text-xs font-black uppercase text-slate-400 mb-1">{{ t('order.totalExpenditure') }}</h3>
      <p class="text-2xl font-black text-slate-900 dark:text-white">
        ¥{{ totalExpenditure.toLocaleString(undefined, { minimumFractionDigits: 2 }) }}
      </p>
    </div>
    <div class="h-10 w-px bg-slate-100 dark:bg-border-dark"></div>
    <div class="flex-1">
      <h3 class="text-xs font-black uppercase text-slate-400 mb-1">{{ t('order.mainExpenditureItem') }}</h3>
      <p class="text-sm font-bold text-slate-700 dark:text-slate-300">{{ mainItem.name }} ({{ mainItem.percent }}%)</p>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  mainItem: {
    type: Object,
    required: true
  },
  totalExpenditure: {
    type: Number,
    default: 0
  }
})

const t = inject('t', (key) => key)
</script>
