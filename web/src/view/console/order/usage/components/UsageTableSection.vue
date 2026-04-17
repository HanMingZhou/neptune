<template>
  <div class="flex flex-col overflow-hidden">
    <div
      class="console-table-scroll console-table-scroll--fill flex-1 overflow-x-auto"
    >
      <table class="console-table w-full">
        <thead>
          <tr
            class="bg-slate-50/50 dark:bg-zinc-800/30 text-slate-500 text-[10px] font-black uppercase tracking-widest border-b border-slate-100 dark:border-border-dark"
          >
            <th class="px-8 py-5">{{ t('order.resourceType') }}</th>
            <th class="px-6 py-5">{{ t('order.instanceName') }}</th>
            <th class="px-6 py-5">{{ t('order.orderNo') }}</th>
            <th class="px-6 py-5">{{ t('order.chargeType') }}</th>
            <th class="px-6 py-5">{{ t('order.usageVolume') }}</th>
            <th class="px-6 py-5">{{ t('order.unitPrice') }}</th>
            <th class="px-6 py-5">{{ t('status') }}</th>
            <th class="px-6 py-5">{{ t('order.createdAt') }}</th>
            <th class="px-6 py-5">{{ t('order.updatedAt') }}</th>
            <th class="px-6 py-5">{{ t('order.subtotal') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 dark:divide-border-dark">
          <tr
            v-for="item in items"
            :key="`${item.id}${item.source}`"
            class="text-sm hover:bg-slate-50/50 dark:hover:bg-zinc-800/20 transition-colors"
          >
            <td class="px-8 py-5">
              <span class="font-bold text-slate-700 dark:text-slate-200">{{
                getResourceTypeText(item.resourceType)
              }}</span>
            </td>
            <td class="px-6 py-5">
              <span class="text-slate-600 dark:text-slate-300 text-xs">{{
                item.resourceName || '-'
              }}</span>
            </td>
            <td class="px-6 py-5">
              <span
                class="text-slate-500 dark:text-slate-400 text-xs font-mono"
                >{{ item.orderNo || '-' }}</span
              >
            </td>
            <td class="px-6 py-5">
              <span
                :class="getChargeTypeStyle(item.source)"
                class="order-tone-chip"
              >
                {{ item.chargeTypeName }}
              </span>
            </td>
            <td
              class="px-6 py-5 font-mono text-slate-600 dark:text-slate-300 text-xs"
            >
              {{ formatUsage(item) }}
            </td>
            <td class="px-6 py-5 text-xs text-slate-600 dark:text-slate-400">
              ¥{{ item.unitPrice.toFixed(6) }} / {{ getUnitText(item) }}
            </td>
            <td class="px-6 py-5">
              <span :class="getStatusStyle(item)" class="order-tone-chip">
                {{ item.statusText }}
              </span>
            </td>
            <td class="px-6 py-5 text-xs text-slate-500 dark:text-slate-400">
              {{ item.createdAt || '-' }}
            </td>
            <td class="px-6 py-5 text-xs text-slate-500 dark:text-slate-400">
              {{ item.updatedAt || '-' }}
            </td>
            <td class="px-6 py-5">
              <span
                class="console-amount-chip console-amount-chip--neutral font-mono"
              >
                ¥{{ item.amount.toFixed(6) }}
              </span>
            </td>
          </tr>
          <tr v-if="items.length === 0">
            <td
              colspan="10"
              class="px-8 py-16 text-center text-slate-400 text-sm"
            >
              <span class="material-icons text-4xl mb-2 block opacity-30"
                >receipt_long</span
              >
              {{ t('noData') }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="console-list-footer">
      <ListPaginationBar
        v-model:current-page="pageModel"
        v-model:page-size="pageSizeModel"
        :total="chargeFilter === 'all' ? total : items.length"
        :show-total="false"
        justify="end"
        :page-sizes="[15, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="$emit('page-change', $event)"
        @size-change="$emit('size-change', $event)"
      />
    </div>

    <div class="p-8 bg-slate-50/30 dark:bg-zinc-900/10 text-center">
      <p class="text-xs text-slate-400 leading-relaxed max-w-lg mx-auto">
        {{ t('order.usageHint') }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import ListPaginationBar from '@/components/listPage/ListPaginationBar.vue'
import type { Translator } from '@/types/consoleResource'
import type { UsageListItem } from '@/types/order'

type UsageStyleGetter = (source?: string) => string
type ResourceTypeGetter = (type?: number) => string
type UsageStatusGetter = (item: UsageListItem) => string
type UsageTextGetter = (item: UsageListItem) => string

const props = defineProps<{
  chargeFilter: string
  getChargeTypeStyle: UsageStyleGetter
  getResourceTypeText: ResourceTypeGetter
  getStatusStyle: UsageStatusGetter
  getUnitText: UsageTextGetter
  items: UsageListItem[]
  formatUsage: UsageTextGetter
  page: number
  pageSize: number
  total: number
}>()

const emit = defineEmits<{
  'page-change': [value: number]
  'size-change': [value: number]
  'update:page': [value: number]
  'update:page-size': [value: number]
}>()
const t = inject<Translator>('t', (key: string) => key)

const pageModel = computed({
  get: () => props.page,
  set: (value: number) => emit('update:page', value)
})

const pageSizeModel = computed({
  get: () => props.pageSize,
  set: (value: number) => emit('update:page-size', value)
})

</script>

