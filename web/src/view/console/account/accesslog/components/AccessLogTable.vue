<template>
  <TableCard :page-size="pageSize">
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
            <span
              class="material-icons absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400 text-[16px]"
            >
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

    <div
      class="console-table-scroll console-table-scroll--fill overflow-x-auto"
    >
      <table class="console-table console-table--compact w-full min-w-[860px]">
        <thead>
          <tr>
            <th>{{ t('time') }}</th>
            <th>IP</th>
            <th>{{ t('location') }}</th>
            <th>{{ t('device') }}</th>
            <th>{{ t('status') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
          <tr v-if="loading">
            <td colspan="5" class="px-6 py-12 text-center text-slate-400">
              <div class="flex items-center justify-center gap-2">
                <div
                  class="animate-spin rounded-full h-5 w-5 border-b-2 border-primary"
                ></div>
                {{ t('loading') }}
              </div>
            </td>
          </tr>
          <tr v-else-if="records.length === 0">
            <td colspan="5" class="px-6 py-12 text-center text-slate-400">
              {{ t('noData') }}
            </td>
          </tr>
          <tr
            v-for="record in records"
            v-else
            :key="record.id"
            class="hover:bg-slate-50 dark:hover:bg-zinc-800/40 transition-colors"
          >
            <td class="is-code is-secondary whitespace-nowrap">
              {{ record.time }}
            </td>
            <td class="is-code">
              {{ record.ip }}
            </td>
            <td class="is-secondary">
              {{ record.location || '-' }}
            </td>
            <td class="is-secondary">
              {{ record.device || '-' }}
            </td>
            <td>
              <ListToneBadge
                :label="record.status === 'Success' ? t('success') : t('failed')"
                :tone="record.status === 'Success' ? 'success' : 'danger'"
              />
            </td>
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
        :page-sizes="[15, 20, 50, 100]"
        layout="sizes, prev, pager, next, jumper"
        @current-change="$emit('page-change', $event)"
        @size-change="$emit('size-change', $event)"
      />
    </template>
  </TableCard>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import ListPaginationBar from '@/components/listPage/ListPaginationBar.vue'
import ListToneBadge from '@/components/listPage/ListToneBadge.vue'
import TableCard from '@/components/listPage/TableCard.vue'
import type { AccessLogRecord } from '@/types/account'
import type { Translator } from '@/types/consoleResource'

const props = defineProps<{
  filterStatus: string
  loading: boolean
  page: number
  pageSize: number
  records: AccessLogRecord[]
  searchIp: string
  total: number
}>()

const emit = defineEmits<{
  'page-change': [value: number]
  search: []
  'size-change': [value: number]
  'update:filter-status': [value: string]
  'update:page': [value: number]
  'update:page-size': [value: number]
  'update:search-ip': [value: string]
}>()

const t = inject<Translator>('t', (key: string) => key)

const searchIpModel = computed({
  get: () => props.searchIp,
  set: (value: string) => emit('update:search-ip', value)
})

const filterStatusModel = computed({
  get: () => props.filterStatus,
  set: (value: string) => emit('update:filter-status', value)
})

const pageModel = computed({
  get: () => props.page,
  set: (value: number) => emit('update:page', value)
})

const pageSizeModel = computed({
  get: () => props.pageSize,
  set: (value: number) => emit('update:page-size', value)
})
</script>

