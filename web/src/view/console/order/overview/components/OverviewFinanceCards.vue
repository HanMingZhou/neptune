<template>
  <div class="grid grid-cols-1 gap-6 md:grid-cols-3">
    <div class="group relative overflow-hidden rounded-2xl bg-primary p-8 text-white shadow-xl shadow-primary/20">
      <span class="material-icons absolute -bottom-4 -right-4 text-[120px] opacity-10 transition-transform group-hover:rotate-12">account_balance_wallet</span>
      <p class="mb-2 text-xs font-black uppercase tracking-widest opacity-80">{{ t('availableBalance') }} (CNY)</p>
      <div class="mb-8 flex items-baseline gap-2">
        <span class="text-4xl font-black">¥{{ balance.toLocaleString(undefined, { minimumFractionDigits: 2 }) }}</span>
      </div>
      <div class="flex gap-3">
        <button
          class="list-row-button list-row-button--neutral"
          @click="$emit('recharge')"
        >
          {{ t('rechargeNow') }}
        </button>
      </div>
    </div>

    <div class="rounded-2xl border border-slate-200 bg-white p-8 shadow-sm dark:border-border-dark dark:bg-surface-dark">
      <p class="mb-2 text-xs font-black uppercase tracking-widest text-slate-400">{{ t('mtdEstimate') }} (MTD)</p>
      <div class="mb-1 flex items-baseline gap-2">
        <span class="text-4xl font-black text-slate-900 dark:text-white">
          ¥{{ mtdEstimate.amount.toLocaleString(undefined, { minimumFractionDigits: 2 }) }}
        </span>
        <span :class="mtdEstimate.change >= 0 ? 'text-red-500' : 'text-emerald-500'" class="text-xs font-bold">
          {{ mtdEstimate.change >= 0 ? '+' : '' }}{{ mtdEstimate.change }}%
        </span>
      </div>
      <p class="mb-8 text-[10px] text-slate-500">
        {{ t('compareLastMonth') }}: ¥{{ mtdEstimate.lastMonth.toLocaleString(undefined, { minimumFractionDigits: 2 }) }}
      </p>
      <div class="flex items-center justify-between rounded-xl border border-slate-100 bg-slate-50 p-4 dark:border-border-dark dark:bg-zinc-800/50">
        <span class="text-xs font-bold text-slate-600 dark:text-slate-400">{{ t('payAsYouGoStatus') }}</span>
        <span class="text-xs font-mono font-bold text-emerald-500">{{ t('statusHealthy') }}</span>
      </div>
    </div>

    <div class="rounded-2xl border border-slate-200 bg-white p-8 shadow-sm dark:border-border-dark dark:bg-surface-dark">
      <p class="mb-2 text-xs font-black uppercase tracking-widest text-slate-400">{{ t('lastMonthSettlement') }}</p>
      <div class="mb-1 flex items-baseline gap-2">
        <span class="text-4xl font-black text-slate-400">
          ¥{{ lastMonthSettlement.amount.toLocaleString(undefined, { minimumFractionDigits: 2 }) }}
        </span>
      </div>
      <p class="mb-8 text-[10px] text-slate-500">{{ t('orderPeriod') }}: {{ lastMonthSettlement.period }}</p>
      <button class="list-row-button list-row-button--info" @click="$emit('view-transactions')">
        {{ t('viewSettlementDetails') }}
        <span class="material-icons text-[14px]">arrow_forward</span>
      </button>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  balance: {
    type: Number,
    default: 0
  },
  lastMonthSettlement: {
    type: Object,
    required: true
  },
  mtdEstimate: {
    type: Object,
    required: true
  }
})

defineEmits(['recharge', 'view-transactions'])

const t = inject('t', (key) => key)
</script>
