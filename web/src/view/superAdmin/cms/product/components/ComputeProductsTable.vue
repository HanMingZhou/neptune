<template>
  <TableCard>
    <div class="overflow-x-auto">
      <table class="console-table console-table--compact w-full min-w-[1320px]" v-loading="loading">
        <thead>
          <tr>
            <th>{{ t('id') }}</th>
            <th>{{ t('productName') }}</th>
            <th>{{ t('nodeName') }}</th>
            <th>{{ t('area') }}</th>
            <th>{{ t('spec') }}</th>
            <th class="text-center">{{ t('inventory') }}</th>
            <th>{{ t('prices') }}({{ t('priceHourly') }}/{{ t('priceDaily') }}/{{ t('priceWeekly') }}/{{ t('priceMonthly') }})</th>
            <th class="text-center">{{ t('status') }}</th>
            <th class="console-actions-header">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
          <tr v-for="row in items" :key="row.id" class="transition-colors hover:bg-slate-50 dark:hover:bg-zinc-800/40">
            <td class="is-code is-secondary">{{ row.id }}</td>
            <td><span class="is-primary">{{ row.name }}</span></td>
            <td><span class="is-secondary is-code">{{ row.nodeName }}</span></td>
            <td class="is-secondary">{{ row.area }}</td>
            <td>
              <div class="text-sm">
                <span v-if="row.gpuModel" class="is-primary">{{ row.gpuModel }} × {{ row.gpuCount }}</span>
                <span v-else-if="row.vGpuNumber > 0" class="is-primary">vGPU: {{ row.vGpuMemory }}GB / {{ row.vGpuCores }}%</span>
                <span v-else class="text-slate-400">CPU ONLY</span>
                <div class="mt-1 text-xs text-slate-500">{{ row.cpu }}核 / {{ row.memory }}GB</div>
              </div>
            </td>
            <td class="text-center">
              <div class="text-sm">
                <span class="is-primary">{{ row.available ?? 0 }}</span>
                <span class="mx-1 text-slate-300">/</span>
                <span class="text-slate-500">{{ row.maxInstances || 0 }}</span>
                <div class="text-xs text-slate-400">
                  <template v-if="row.gpuCount > 0">{{ t('gpu') }}</template>
                  <template v-else-if="row.vGpuNumber > 0">vGPU</template>
                  <template v-else>{{ t('instances') }}</template>
                </div>
              </div>
            </td>
            <td>
              <div class="text-xs text-slate-600 dark:text-slate-300">
                <span class="is-primary">¥{{ row.priceHourly?.toFixed(2) || '0.00' }}</span> /
                <span class="is-primary">¥{{ row.priceDaily?.toFixed(2) || '0.00' }}</span> /
                <span class="is-primary">¥{{ row.priceWeekly?.toFixed(2) || '0.00' }}</span> /
                <span class="is-primary">¥{{ row.priceMonthly?.toFixed(2) || '0.00' }}</span>
              </div>
            </td>
            <td class="text-center">
              <ListToneBadge :label="row.status === 1 ? t('onShelf') : t('offShelf')" :tone="row.status === 1 ? 'success' : 'neutral'" />
            </td>
            <td class="console-actions-cell">
              <div class="list-row-actions">
                <button class="list-row-button list-row-button--info" @click="$emit('edit', row)">{{ t('edit') }}</button>
                <button class="list-row-button list-row-button--warning" @click="$emit('adjust-price', row)">{{ t('prices') }}</button>
                <button class="list-row-button list-row-button--danger" @click="$emit('delete', row)">{{ t('delete') }}</button>
              </div>
            </td>
          </tr>
          <tr v-if="items.length === 0">
            <td colspan="9" class="px-6 py-12 text-center text-sm text-slate-400">{{ t('noData') }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <template #footer>
      <ListPaginationBar
        v-model:current-page="pageModel"
        v-model:page-size="pageSizeModel"
        :total="total"
        :total-text="t('totalRecords', { total })"
        :page-sizes="[10, 20, 50, 100]"
        layout="sizes, prev, pager, next, jumper"
        @current-change="$emit('page-change', $event)"
        @size-change="$emit('size-change', $event)"
      />
    </template>
  </TableCard>
</template>

<script setup>
import { computed, inject } from 'vue'
import ListPaginationBar from '@/components/listPage/ListPaginationBar.vue'
import ListToneBadge from '@/components/listPage/ListToneBadge.vue'
import TableCard from '@/components/listPage/TableCard.vue'

const props = defineProps({
  items: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  page: {
    type: Number,
    default: 1
  },
  pageSize: {
    type: Number,
    default: 20
  },
  total: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits(['adjust-price', 'delete', 'edit', 'page-change', 'size-change', 'update:page', 'update:page-size'])
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
