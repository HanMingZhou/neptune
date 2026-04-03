<template>
  <TableCard>
    <div class="flex border-b border-slate-100 dark:border-border-dark">
      <button
        v-for="tab in tabs"
        :key="tab.value"
        :class="activeTab === tab.value ? 'text-primary border-b-2 border-primary' : 'text-slate-500 hover:text-slate-700'"
        class="px-8 py-4 text-xs font-black uppercase tracking-widest transition-colors"
        @click="$emit('update:activeTab', tab.value)"
      >
        {{ tab.label }}
      </button>
    </div>

    <div v-if="activeTab === 'records'" class="overflow-x-auto">
      <table class="w-full text-left">
        <thead>
          <tr class="bg-slate-50/50 dark:bg-zinc-800/30 text-slate-500 text-[10px] font-black uppercase tracking-widest border-b border-slate-100 dark:border-border-dark">
            <th class="px-8 py-5">{{ t('order.requestIdDate') }}</th>
            <th class="px-6 py-5">{{ t('order.invoiceTitle') }}</th>
            <th class="px-6 py-5">{{ t('order.amount') }}</th>
            <th class="px-6 py-5">{{ t('order.status') }}</th>
            <th class="px-6 py-5 text-right">{{ t('order.actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 dark:divide-border-dark">
          <tr v-for="invoice in invoices" :key="invoice.id" class="text-sm hover:bg-slate-50/50 dark:hover:bg-zinc-800/20 transition-colors">
            <td class="px-8 py-5">
              <div class="flex flex-col">
                <span class="font-mono text-xs font-bold text-slate-700 dark:text-slate-200">{{ invoice.id }}</span>
                <span class="text-[10px] text-slate-400 mt-1">{{ invoice.date }}</span>
              </div>
            </td>
            <td class="px-6 py-5">
              <div class="flex flex-col">
                <span class="font-bold truncate max-w-[200px] text-slate-700 dark:text-slate-200">{{ invoice.title }}</span>
                <span class="text-[10px] text-slate-500 mt-0.5">{{ invoice.type }}</span>
              </div>
            </td>
            <td class="px-6 py-5 font-black text-slate-900 dark:text-white">{{ invoice.amount }}</td>
            <td class="px-6 py-5">
              <span
                :class="invoice.status === 'Sent' ? 'bg-emerald-500/10 text-emerald-500' : 'bg-amber-500/10 text-amber-500'"
                class="px-2 py-0.5 rounded text-[10px] font-black uppercase"
              >
                {{ invoice.status === 'Sent' ? t('order.sent') : t('order.processing') }}
              </span>
            </td>
            <td class="px-6 py-5 text-right">
              <button v-if="invoice.status === 'Sent'" class="list-row-button list-row-button--info">
                {{ t('order.downloadElectronic') }}
              </button>
              <button v-else class="list-row-button list-row-button--neutral">
                {{ t('order.viewProgress') }}
              </button>
            </td>
          </tr>
          <tr v-if="invoices.length === 0">
            <td colspan="5" class="px-8 py-12 text-center text-slate-400 text-sm">{{ t('noData') }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div v-else class="p-20 text-center text-slate-400">
      <span class="material-icons text-6xl opacity-20 mb-4">construction</span>
      <p>{{ t('order.comingSoon') }}</p>
    </div>
  </TableCard>
</template>

<script setup>
import { computed, inject } from 'vue'
import TableCard from '@/components/listPage/TableCard.vue'

defineProps({
  activeTab: {
    type: String,
    default: 'records'
  },
  invoices: {
    type: Array,
    default: () => []
  }
})

defineEmits(['update:activeTab'])

const t = inject('t', (key) => key)

const tabs = computed(() => [
  { value: 'records', label: t('order.invoiceRecords') },
  { value: 'titles', label: t('order.invoiceTitles') },
  { value: 'addresses', label: t('order.invoiceAddresses') }
])
</script>
