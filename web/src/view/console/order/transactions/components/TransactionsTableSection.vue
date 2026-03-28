<template>
  <div>
    <div class="overflow-x-auto">
      <table class="w-full text-left">
        <thead>
          <tr class="bg-slate-50 dark:bg-zinc-800/50 border-b border-slate-100 dark:border-border-dark text-slate-500 text-xs font-bold uppercase tracking-wider">
            <th class="px-4 py-3">{{ t('transactionId') }}</th>
            <th class="px-4 py-3">{{ t('time') }}</th>
            <th class="px-4 py-3">{{ t('type') }}</th>
            <th class="px-4 py-3">{{ t('order.resourceInfo') }}</th>
            <th class="px-4 py-3">{{ t('order.orderNo') }}</th>
            <th class="px-4 py-3">{{ t('remark') }}</th>
            <th class="px-4 py-3 text-right">{{ t('amount') }}</th>
            <th class="px-4 py-3 text-right">{{ t('availableBalance') }}</th>
            <th class="px-4 py-3">{{ t('status') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 dark:divide-border-dark">
          <tr v-for="tx in items" :key="tx.id" class="hover:bg-slate-50 dark:hover:bg-zinc-800 transition-colors">
            <td class="px-4 py-3 text-xs font-mono text-slate-700 dark:text-slate-200 whitespace-nowrap">{{ tx.id }}</td>
            <td class="px-4 py-3 text-xs text-slate-500 dark:text-slate-400 whitespace-nowrap">{{ tx.time }}</td>
            <td class="px-4 py-3">
              <span :class="tx.typeStyle" class="px-2 py-0.5 rounded text-[10px] font-black uppercase whitespace-nowrap">
                {{ tx.typeLabel }}
              </span>
            </td>
            <td class="px-4 py-3 text-xs text-slate-700 dark:text-slate-200 whitespace-nowrap">
              {{ tx.resourceName ? `${tx.resourceName} (${tx.resourceTypeText})` : '-' }}
            </td>
            <td class="px-4 py-3 text-xs font-mono text-slate-500 dark:text-slate-400 whitespace-nowrap">{{ tx.orderNo || '-' }}</td>
            <td class="px-4 py-3 text-xs text-slate-500 dark:text-slate-400">{{ tx.remark || '-' }}</td>
            <td class="px-4 py-3 text-right whitespace-nowrap">
              <span :class="tx.rawAmount > 0 ? 'text-emerald-500' : 'text-slate-900 dark:text-white'" class="text-sm font-bold font-mono">
                {{ tx.amount }}
              </span>
            </td>
            <td class="px-4 py-3 text-right text-xs font-mono text-slate-600 dark:text-slate-300 whitespace-nowrap">
              ¥{{ tx.balanceAfter.toFixed(2) }}
            </td>
            <td class="px-4 py-3">
              <span v-if="tx.statusText" :class="tx.statusStyle" class="px-2 py-0.5 rounded text-[10px] font-black uppercase whitespace-nowrap">
                {{ tx.statusText }}
              </span>
              <span v-else class="text-emerald-500 material-icons text-[16px]">check_circle</span>
            </td>
          </tr>
          <tr v-if="items.length === 0">
            <td colspan="9" class="px-6 py-12 text-center text-slate-400 text-sm">{{ t('noData') }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="px-6 py-4 flex flex-col sm:flex-row items-center justify-between border-t border-slate-100 dark:border-border-dark bg-slate-50/50 dark:bg-zinc-800/20">
      <div class="flex-1 flex justify-end">
        <el-pagination
          v-model:current-page="pageModel"
          v-model:page-size="pageSizeModel"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="$emit('page-change', $event)"
          @size-change="$emit('size-change', $event)"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'

const props = defineProps({
  items: {
    type: Array,
    default: () => []
  },
  page: {
    type: Number,
    default: 1
  },
  pageSize: {
    type: Number,
    default: 10
  },
  total: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits(['page-change', 'size-change', 'update:page', 'update:page-size'])
const t = inject('t', (key) => key)

const pageModel = computed({
  get: () => props.page,
  set: (value) => emit('update:page', value)
})

const pageSizeModel = computed({
  get: () => props.pageSize,
  set: (value) => emit('update:page-size', value)
})
</script>
