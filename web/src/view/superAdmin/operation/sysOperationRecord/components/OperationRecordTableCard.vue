<template>
  <TableCard>
    <template #toolbar>
      <div class="list-filter-bar">
      <div class="list-filter-field list-filter-field--compact max-w-[180px]">
        <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]">http</span>
        <input
          v-model="searchInfo.method"
          type="text"
          :placeholder="t('method')"
          class="list-search-input !w-full"
          @keyup.enter="$emit('search')"
        />
      </div>
      <div class="list-filter-field list-filter-field--compact max-w-[180px]">
        <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]">link</span>
        <input
          v-model="searchInfo.path"
          type="text"
          :placeholder="t('requestPath')"
          class="list-search-input !w-full"
          @keyup.enter="$emit('search')"
        />
      </div>
      <div class="list-filter-field list-filter-field--compact max-w-[180px]">
        <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]">tag</span>
        <input
          v-model="searchInfo.status"
          type="text"
          :placeholder="t('statusCode')"
          class="list-search-input !w-full"
          @keyup.enter="$emit('search')"
        />
      </div>
      <div class="list-toolbar-actions">
        <button
          class="list-toolbar-button list-toolbar-button--primary"
          @click="$emit('search')"
        >
          <span class="material-icons text-[18px]">search</span>
          {{ t('searchQuery') }}
        </button>
        <button
          class="list-toolbar-button list-toolbar-button--secondary"
          @click="$emit('reset')"
        >
          <span class="material-icons text-[18px]">refresh</span>
          {{ t('reset') }}
        </button>
      </div>
      </div>
    </template>

    <div class="overflow-x-auto" v-loading="loading">
      <table class="console-table w-full min-w-[1280px] text-left">
        <thead>
          <tr class="bg-slate-50 dark:bg-zinc-800/50 border-b border-border-light dark:border-border-dark text-slate-500 text-xs font-bold uppercase tracking-wider">
            <th class="px-6 py-4 w-12">
              <input type="checkbox" :checked="allSelected" @change="$emit('toggle-select-all', $event.target.checked)" class="rounded" />
            </th>
            <th class="px-6 py-4">{{ t('operator') }}</th>
            <th class="px-6 py-4">{{ t('date') }}</th>
            <th class="px-6 py-4">{{ t('statusCode') }}</th>
            <th class="px-6 py-4">{{ t('requestIp') }}</th>
            <th class="px-6 py-4">{{ t('method') }}</th>
            <th class="px-6 py-4">{{ t('requestPath') }}</th>
            <th class="px-6 py-4">{{ t('requestBody') }}</th>
            <th class="px-6 py-4">{{ t('responseBody') }}</th>
            <th class="px-6 py-4 console-actions-header">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
          <tr v-if="items.length === 0">
            <td colspan="10" class="px-6 py-12 text-center text-slate-400">
              <div class="space-y-2">
                <span class="material-icons text-4xl">inbox</span>
                <p>{{ t('noOperationRecords') }}</p>
              </div>
            </td>
          </tr>
          <tr
            v-for="row in items"
            :key="row.ID"
            class="hover:bg-slate-50 dark:hover:bg-zinc-800/40 transition-colors"
          >
            <td class="px-6 py-4">
              <input
                type="checkbox"
                :checked="selectedIds.includes(row.ID)"
                @change="$emit('toggle-select', row)"
                class="rounded"
              />
            </td>
            <td class="px-6 py-4 text-sm">
              <span class="is-key">{{ row.user?.userName || '-' }}</span>
              <span v-if="row.user?.nickName" class="is-muted">({{ row.user.nickName }})</span>
            </td>
            <td class="px-6 py-4 text-sm text-slate-500">{{ formatDate(row.CreatedAt) }}</td>
            <td class="px-6 py-4">
              <ListToneBadge :label="String(row.status)" :tone="getStatusTone(row.status)" />
            </td>
            <td class="px-6 py-4 text-sm font-mono text-slate-500">{{ row.ip }}</td>
            <td class="px-6 py-4">
              <ListToneBadge :label="row.method" :tone="getMethodTone(row.method)" />
            </td>
            <td class="px-6 py-4 text-sm text-slate-600 dark:text-slate-400 max-w-[220px] truncate" :title="row.path">
              {{ row.path }}
            </td>
            <td class="px-6 py-4">
              <el-popover v-if="row.body" placement="left-start" :width="444">
                <div class="bg-slate-900 rounded-lg p-4 max-h-[400px] overflow-auto">
                  <pre class="text-green-400 text-xs font-mono whitespace-pre-wrap">{{ formatBody(row.body) }}</pre>
                </div>
                <template #reference>
                  <span class="material-icons text-primary cursor-pointer text-[18px]">visibility</span>
                </template>
              </el-popover>
              <span v-else class="text-slate-400 text-sm">-</span>
            </td>
            <td class="px-6 py-4">
              <el-popover v-if="row.resp" placement="left-start" :width="444">
                <div class="bg-slate-900 rounded-lg p-4 max-h-[400px] overflow-auto">
                  <pre class="text-green-400 text-xs font-mono whitespace-pre-wrap">{{ formatBody(row.resp) }}</pre>
                </div>
                <template #reference>
                  <span class="material-icons text-primary cursor-pointer text-[18px]">visibility</span>
                </template>
              </el-popover>
              <span v-else class="text-slate-400 text-sm">-</span>
            </td>
            <td class="px-6 py-4 console-actions-cell">
              <div class="list-row-actions">
                <button
                  class="list-row-button list-row-button--danger"
                  @click="$emit('delete', row)"
                >
                  <span class="material-icons text-[16px]">delete</span>
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
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        :total-text="t('totalRecords', { total })"
        :page-sizes="[10, 30, 50, 100]"
        layout="sizes, prev, pager, next, jumper"
        @current-change="$emit('page-change', $event)"
        @size-change="$emit('size-change', $event)"
      />
    </template>
  </TableCard>
</template>

<script setup>
import { inject } from 'vue'
import ListToneBadge from '@/components/listPage/ListToneBadge.vue'
import { formatDate } from '@/utils/format'
import ListPaginationBar from '@/components/listPage/ListPaginationBar.vue'
import TableCard from '@/components/listPage/TableCard.vue'

defineProps({
  allSelected: {
    type: Boolean,
    default: false
  },
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
    default: 10
  },
  searchInfo: {
    type: Object,
    required: true
  },
  selectedIds: {
    type: Array,
    default: () => []
  },
  total: {
    type: Number,
    default: 0
  }
})

defineEmits([
  'delete',
  'page-change',
  'reset',
  'search',
  'size-change',
  'toggle-select',
  'toggle-select-all'
])

const t = inject('t', (key) => key)

const getMethodTone = (method) => {
  const tones = {
    GET: 'info',
    POST: 'success',
    PUT: 'warning',
    DELETE: 'danger'
  }

  return tones[method] || 'neutral'
}

const getStatusTone = (status) => {
  if (Number(status) >= 400) {
    return 'danger'
  }

  return 'success'
}

const formatBody = (value) => {
  try {
    return JSON.stringify(JSON.parse(value), null, 2)
  } catch {
    return value
  }
}
</script>
