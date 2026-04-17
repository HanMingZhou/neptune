<template>
  <TableCard :page-size="pageSize">
    <div
      class="console-table-scroll console-table-scroll--fill overflow-x-auto"
    >
      <table class="console-table min-w-[1180px]" v-loading="loading">
        <thead>
          <tr>
            <th class="px-6 py-4 text-left">{{ t('nodeName') }}</th>
            <th class="px-6 py-4 text-left">{{ t('internalIp') }}</th>
            <th class="px-6 py-4 text-left">{{ t('clusterName') }}</th>
            <th class="px-6 py-4 text-left">{{ t('nodeRole') }}</th>
            <th class="px-6 py-4 text-left">{{ t('area') }}</th>
            <th class="px-6 py-4 text-center">{{ t('spec') }} (可用/总量)</th>
            <th class="px-6 py-4 text-center">{{ t('status') }}</th>
            <th class="px-6 py-4 console-actions-header">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
          <tr
            v-for="row in items"
            :key="row.nodeName"
            class="hover:bg-slate-50 dark:hover:bg-zinc-800/40 transition-colors"
          >
            <td class="px-6 py-4">
              <div class="flex flex-col">
                <span
                  class="console-badge console-badge--neutral is-code w-fit"
                  >{{ row.nodeName }}</span
                >
                <span class="text-[10px] text-slate-400 mt-1">{{
                  row.area || '-'
                }}</span>
              </div>
            </td>
            <td class="px-6 py-4 text-sm text-slate-600 font-mono">
              {{ row.internalIp || 'N/A' }}
            </td>
            <td class="px-6 py-4 text-sm text-slate-600">
              {{ row.clusterName || '-' }}
            </td>
            <td class="px-6 py-4">
              <span class="console-badge console-badge--neutral">
                {{
                  row.nodeRole === 'master' ? t('roleMaster') : t('roleWorker')
                }}
              </span>
            </td>
            <td class="px-6 py-4 text-sm text-slate-600">
              {{ row.area || '-' }}
            </td>
            <td class="px-6 py-4 text-center">
              <div class="flex flex-col items-center space-y-1.5">
                <div
                  v-if="row.gpuModel"
                  class="flex flex-wrap items-center justify-center gap-2"
                >
                  <span
                    class="rounded border border-border-light bg-slate-100 px-1.5 py-0.5 text-[10px] font-semibold leading-none text-slate-500 dark:border-border-dark dark:bg-zinc-800 dark:text-slate-300"
                    >GPU</span
                  >
                  <span class="is-primary">{{ row.gpuModel }}</span>
                  <span class="is-code text-xs text-slate-500">
                    ({{ row.gpuAvailable }}/{{ row.gpuCount }})
                  </span>
                </div>

                <div class="flex flex-wrap items-center justify-center gap-4">
                  <div class="flex items-center justify-center gap-1.5">
                    <span class="material-icons text-[14px] text-slate-400"
                      >developer_board</span
                    >
                    <span
                      class="text-xs text-slate-600 dark:text-slate-400 font-mono"
                      >{{ row.cpuAvailable }}/{{ row.cpuAllocatable }}核</span
                    >
                  </div>
                  <div class="flex items-center justify-center gap-1.5">
                    <span class="material-icons text-[14px] text-slate-400"
                      >memory</span
                    >
                    <span
                      class="text-xs text-slate-600 dark:text-slate-400 font-mono"
                      >{{ row.memoryAvailable }}/{{
                        row.memoryAllocatable
                      }}GB</span
                    >
                  </div>
                </div>
              </div>
            </td>
            <td class="px-6 py-4 text-center">
              <ListToneBadge
                :label="row.schedulable ? t('schedulable') : t('unschedulable')"
                :tone="row.schedulable ? 'success' : 'warning'"
              />
            </td>
            <td class="px-6 py-4 console-actions-cell">
              <div class="list-row-actions">
                <button
                  v-if="!row.schedulable"
                  class="list-row-button list-row-button--success"
                  @click="$emit('uncordon', row)"
                >
                  <span class="material-icons text-[16px]">login</span>
                  {{ t('joinCluster') }}
                </button>
                <button
                  v-if="row.schedulable"
                  :disabled="isMasterNode(row)"
                  :title="isMasterNode(row) ? t('masterNodeEvictDisabled') : ''"
                  class="list-row-button list-row-button--danger disabled:opacity-50"
                  @click="$emit('drain', row)"
                >
                  <span class="material-icons text-[16px]">logout</span>
                  {{ t('evictFromCluster') }}
                </button>
              </div>
            </td>
          </tr>
          <tr v-if="items.length === 0 && !loading">
            <td
              colspan="8"
              class="px-6 py-10 text-center text-slate-400 text-sm"
            >
              {{ t('noData') }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <template #footer>
      <ListPaginationBar
        :current-page="page"
        :page-size="pageSize"
        :page-sizes="[15, 20, 50, 100]"
        :total="total"
        :total-text="t('totalRecords', { total })"
        layout="sizes, prev, pager, next, jumper"
        @current-change="$emit('page-change', $event)"
        @size-change="$emit('size-change', $event)"
      />
    </template>
  </TableCard>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import ListPaginationBar from '@/components/listPage/ListPaginationBar.vue'
import ListToneBadge from '@/components/listPage/ListToneBadge.vue'
import TableCard from '@/components/listPage/TableCard.vue'
import type { Translator } from '@/types/consoleResource'
import type { CmsNodeRow } from '@/types/superAdmin'

withDefaults(
  defineProps<{
    items?: CmsNodeRow[]
    loading?: boolean
    page?: number
    pageSize?: number
    total?: number
  }>(),
  {
    items: () => [],
    loading: false,
    page: 1,
    pageSize: 15,
    total: 0
  }
)

defineEmits<{
  drain: [row: CmsNodeRow]
  'page-change': [value: number]
  'size-change': [value: number]
  uncordon: [row: CmsNodeRow]
}>()

const t = inject<Translator>('t', (key: string) => key)

const isMasterNode = (row: CmsNodeRow): boolean => row.nodeRole === 'master'
</script>
