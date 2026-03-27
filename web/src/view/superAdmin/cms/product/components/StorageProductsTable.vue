<template>
  <div class="overflow-hidden rounded-xl border border-border-light bg-surface-light shadow-sm dark:border-border-dark dark:bg-surface-dark">
    <div class="overflow-x-auto">
      <table class="w-full" v-loading="loading">
        <thead>
          <tr class="border-b border-border-light bg-slate-50 text-xs font-bold uppercase tracking-wider text-slate-500 dark:border-border-dark dark:bg-zinc-800/50">
            <th class="px-6 py-4">{{ t('id') }}</th>
            <th class="px-6 py-4">{{ t('productName') }}</th>
            <th class="px-6 py-4">{{ t('storageClass') }}</th>
            <th class="px-6 py-4">{{ t('area') }}</th>
            <th class="px-6 py-4">{{ t('prices') }}</th>
            <th class="px-6 py-4 text-center">{{ t('status') }}</th>
            <th class="px-6 py-4">{{ t('paramDesc') }}</th>
            <th class="px-6 py-4 text-center">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
          <tr v-for="row in items" :key="row.id" class="transition-colors hover:bg-slate-50 dark:hover:bg-zinc-800/40">
            <td class="px-6 py-4 text-sm font-mono text-slate-500">{{ row.id }}</td>
            <td class="px-6 py-4 text-sm font-bold text-primary">{{ row.name }}</td>
            <td class="px-6 py-4">
              <code class="rounded bg-purple-500/10 px-2 py-1 text-xs text-purple-600">{{ row.storageClass }}</code>
            </td>
            <td class="px-6 py-4 text-sm text-slate-600 dark:text-slate-300">{{ row.area }}</td>
            <td class="px-6 py-4">
              <span class="text-sm font-bold text-amber-600">¥{{ row.storagePriceGb?.toFixed(4) || '0' }}/GB/日</span>
            </td>
            <td class="px-6 py-4 text-center">
              <span
                class="rounded-full px-2.5 py-1 text-xs font-bold"
                :class="row.status === 1 ? 'bg-emerald-500/10 text-emerald-500' : 'bg-slate-500/10 text-slate-500'"
              >
                {{ row.status === 1 ? t('onShelf') : t('offShelf') }}
              </span>
            </td>
            <td class="max-w-[200px] truncate px-6 py-4 text-sm text-slate-500 dark:text-slate-400">{{ row.description }}</td>
            <td class="px-6 py-4 text-center">
              <div class="flex items-center justify-center gap-2">
                <button class="rounded-sm bg-primary/10 px-2 py-1 text-xs font-bold text-primary transition-colors hover:bg-primary/20" @click="$emit('edit', row)">{{ t('edit') }}</button>
                <button class="rounded-sm bg-red-500/10 px-2 py-1 text-xs font-bold text-red-600 transition-colors hover:bg-red-500/20" @click="$emit('delete', row)">{{ t('delete') }}</button>
              </div>
            </td>
          </tr>
          <tr v-if="items.length === 0">
            <td colspan="8" class="px-6 py-12 text-center text-sm text-slate-400">{{ t('noData') }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="flex items-center justify-between border-t border-border-light bg-slate-50 px-6 py-4 dark:border-border-dark dark:bg-zinc-800/30">
      <span class="text-xs text-slate-500">{{ t('totalRecords', { total }) }}</span>
      <el-pagination
        v-model:current-page="pageModel"
        v-model:page-size="pageSizeModel"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="sizes, prev, pager, next"
        @current-change="$emit('page-change', $event)"
        @size-change="$emit('size-change', $event)"
      />
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

const emit = defineEmits(['delete', 'edit', 'page-change', 'size-change', 'update:page', 'update:page-size'])
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
