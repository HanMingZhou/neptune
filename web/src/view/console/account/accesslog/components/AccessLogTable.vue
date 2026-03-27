<template>
  <TableCard>
    <template #toolbar>
      <div class="flex flex-wrap gap-4 items-center">
        <div class="flex items-center gap-3 flex-1">
          <div class="relative">
            <input
              v-model="searchIpModel"
              type="text"
              :placeholder="t('audit.searchIp')"
              class="list-search-input"
              @keyup.enter="$emit('search')"
            />
            <span class="material-icons absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400 text-[16px]">
              search
            </span>
          </div>
          <el-select
            v-model="filterStatusModel"
            :placeholder="t('audit.allStatus')"
            class="list-filter-select"
            size="small"
            clearable
          >
            <el-option :label="t('success')" value="Success" />
            <el-option :label="t('failed')" value="Failed" />
          </el-select>
        </div>
      </div>
    </template>

    <div class="overflow-x-auto">
      <table class="w-full text-left">
        <thead>
          <tr class="bg-gray-50 dark:bg-zinc-900/50 border-b border-gray-100 dark:border-gray-800 text-slate-500 text-[11px] font-bold uppercase tracking-wider">
            <th class="px-6 py-4">{{ t('time') }}</th>
            <th class="px-6 py-4">IP</th>
            <th class="px-6 py-4">{{ t('location') }}</th>
            <th class="px-6 py-4">{{ t('device') }}</th>
            <th class="px-6 py-4">{{ t('status') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-100 dark:divide-border-dark text-[13px]">
          <tr v-if="loading">
            <td colspan="5" class="px-6 py-12 text-center text-slate-400">
              <div class="flex items-center justify-center gap-2">
                <div class="animate-spin rounded-full h-5 w-5 border-b-2 border-primary"></div>
                {{ t('loading') }}
              </div>
            </td>
          </tr>
          <tr v-else-if="records.length === 0">
            <td colspan="5" class="px-6 py-12 text-center text-slate-400">{{ t('noData') }}</td>
          </tr>
          <tr
            v-for="record in records"
            v-else
            :key="record.id"
            class="hover:bg-gray-50 dark:hover:bg-zinc-800/50 transition-colors"
          >
            <td class="px-6 py-4 text-slate-500 whitespace-nowrap text-xs font-mono">{{ record.time }}</td>
            <td class="px-6 py-4 text-xs font-mono text-slate-700 dark:text-slate-200">{{ record.ip }}</td>
            <td class="px-6 py-4 text-xs text-slate-500">{{ record.location || '-' }}</td>
            <td class="px-6 py-4 text-xs text-slate-500">{{ record.device || '-' }}</td>
            <td class="px-6 py-4">
              <span
                :class="record.status === 'Success' ? 'bg-emerald-500/10 text-emerald-500' : 'bg-red-500/10 text-red-500'"
                class="px-2 py-0.5 rounded text-[10px] font-black uppercase"
              >
                {{ record.status === 'Success' ? t('success') : t('failed') }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <template #footer>
      <div class="flex items-center justify-between text-xs text-slate-500">
        <span>{{ t('totalRecords', { total }) }}</span>
        <el-pagination
          v-model:current-page="pageModel"
          v-model:page-size="pageSizeModel"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="sizes, prev, pager, next"
          size="small"
          class="!p-0"
          @current-change="$emit('page-change', $event)"
          @size-change="$emit('size-change', $event)"
        />
      </div>
    </template>
  </TableCard>
</template>

<script setup>
import { computed, inject } from 'vue'
import TableCard from '@/components/listPage/TableCard.vue'

const props = defineProps({
  filterStatus: {
    type: String,
    default: ''
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
    default: 10
  },
  records: {
    type: Array,
    default: () => []
  },
  searchIp: {
    type: String,
    default: ''
  },
  total: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits([
  'page-change',
  'search',
  'size-change',
  'update:filter-status',
  'update:page',
  'update:page-size',
  'update:search-ip'
])

const t = inject('t', (key) => key)

const searchIpModel = computed({
  get: () => props.searchIp,
  set: (value) => emit('update:search-ip', value)
})

const filterStatusModel = computed({
  get: () => props.filterStatus,
  set: (value) => emit('update:filter-status', value)
})

const pageModel = computed({
  get: () => props.page,
  set: (value) => emit('update:page', value)
})

const pageSizeModel = computed({
  get: () => props.pageSize,
  set: (value) => emit('update:page-size', value)
})
</script>
