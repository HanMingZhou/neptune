<template>
  <TableCard :page-size="pageSize">
    <div
      class="console-table-scroll console-table-scroll--fill overflow-x-auto"
    >
      <table class="console-table">
        <thead>
          <tr>
            <th>{{ t('name') }}</th>
            <th>{{ t('pvc') }}</th>
            <th>{{ t('storageProduct') }}</th>
            <th>{{ t('capacity') }}</th>
            <th>{{ t('status') }}</th>
            <th>{{ t('area') }}</th>
            <th>{{ t('createdAt') }}</th>
            <th class="console-actions-header">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td colspan="8" class="px-6 py-12 text-center text-slate-400">
              <div class="flex items-center justify-center gap-2">
                <div
                  class="animate-spin rounded-full h-5 w-5 border-b-2 border-primary"
                ></div>
                {{ t('loading') }}
              </div>
            </td>
          </tr>
          <tr v-else-if="items.length === 0">
            <td colspan="8" class="px-6 py-12 text-center text-slate-400">
              <div class="space-y-2">
                <span class="material-icons text-4xl">folder_open</span>
                <p>{{ t('noData') }}</p>
                <button
                  class="list-row-button list-row-button--neutral"
                  @click="$emit('create')"
                >
                  {{ t('create') }}{{ t('storage') }}
                </button>
              </div>
            </td>
          </tr>
          <tr v-for="item in items" v-else :key="item.id">
            <td>
              <span class="is-primary text-sm">{{ item.name }}</span>
            </td>
            <td class="detail-info-value--mono">{{ item.pvcName || '-' }}</td>
            <td>{{ item.productName || '-' }}</td>
            <td class="is-code">{{ item.size }}</td>
            <td>
              <span :class="resolveStatusBadgeClass(item.status)">
                {{ t(item.status || '') || item.status || '-' }}
              </span>
            </td>
            <td>{{ item.area || '-' }}</td>
            <td class="text-slate-500">{{ item.createdAt || '-' }}</td>
            <td class="console-actions-cell">
              <div class="list-row-actions">
                <button
                  class="list-row-button list-row-button--neutral"
                  @click="$emit('edit', item)"
                >
                  <span class="material-icons text-[14px]">edit</span>
                  {{ t('edit') }}
                </button>
                <button
                  class="list-row-button list-row-button--info"
                  @click="$emit('expand', item)"
                >
                  <span class="material-icons text-[14px]"
                    >add_circle_outline</span
                  >
                  {{ t('expand') }}
                </button>
                <button
                  :disabled="Boolean(btnLoading[String(item.id)])"
                  class="list-row-button list-row-button--danger disabled:opacity-50"
                  @click="$emit('delete', item)"
                >
                  <span class="material-icons text-[14px]">delete</span>
                  {{ t('delete') }}
                </button>
              </div>
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
        :page-size="pageSize"
        :page-sizes="[15, 20, 50, 100]"
        layout="sizes, prev, pager, next, jumper"
        @current-change="$emit('refresh')"
        @size-change="$emit('refresh')"
      />
    </template>
  </TableCard>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import type { Translator } from '@/types/consoleResource'
import type { StorageListItem } from '@/types/storage'
import ListPaginationBar from '@/components/listPage/ListPaginationBar.vue'
import TableCard from '@/components/listPage/TableCard.vue'

const props = withDefaults(
  defineProps<{
    btnLoading?: Record<string, boolean>
    items?: StorageListItem[]
    loading?: boolean
    page?: number
    pageSize?: number
    total?: number
  }>(),
  {
    btnLoading: () => ({}),
    items: () => [],
    loading: false,
    page: 1,
    pageSize: 15,
    total: 0
  }
)

const emit = defineEmits<{
  create: []
  delete: [item: StorageListItem]
  edit: [item: StorageListItem]
  expand: [item: StorageListItem]
  refresh: []
  'update:page': [value: number]
  'update:page-size': [value: number]
}>()

const t = inject<Translator>('t', (key: string) => key)

const resolveStatusBadgeClass = (status?: string): string => {
  const raw = `${status || ''}`.trim()
  const normalized = raw.toUpperCase()

  if (
    raw === '使用中' ||
    normalized === 'IN_USE' ||
    normalized === 'INUSE' ||
    normalized === 'BOUND' ||
    normalized === 'READY' ||
    normalized === 'RUNNING' ||
    normalized === 'SUCCEEDED'
  ) {
    return 'console-badge console-badge--success'
  }

  if (
    normalized === 'PENDING' ||
    normalized === 'CREATING' ||
    normalized === 'EXPANDING' ||
    normalized === 'RESIZE_PENDING' ||
    normalized === 'INQUEUE'
  ) {
    return 'console-badge console-badge--warning'
  }

  if (normalized === 'FAILED' || normalized === 'ERROR') {
    return 'console-badge console-badge--danger'
  }

  if (raw === '未使用' || normalized === 'UNUSED' || normalized === 'STOPPED') {
    return 'console-badge console-badge--info'
  }

  return 'console-badge console-badge--neutral'
}

const pageModel = computed({
  get: () => props.page,
  set: (value: number) => emit('update:page', value)
})

const pageSizeModel = computed({
  get: () => props.pageSize,
  set: (value: number) => emit('update:page-size', value)
})
</script>

