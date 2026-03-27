<template>
  <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl overflow-hidden shadow-sm">
    <div class="p-4 border-b border-border-light dark:border-border-dark flex flex-wrap gap-4 items-center">
      <div class="relative flex-1 max-w-[180px]">
        <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]">http</span>
        <input
          v-model="searchInfo.method"
          type="text"
          :placeholder="t('method')"
          class="w-full pl-10 pr-4 py-2 bg-slate-50 dark:bg-zinc-800 border-none rounded-lg text-sm focus:ring-1 focus:ring-primary outline-none"
          @keyup.enter="$emit('search')"
        />
      </div>
      <div class="relative flex-1 max-w-[180px]">
        <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]">link</span>
        <input
          v-model="searchInfo.path"
          type="text"
          :placeholder="t('requestPath')"
          class="w-full pl-10 pr-4 py-2 bg-slate-50 dark:bg-zinc-800 border-none rounded-lg text-sm focus:ring-1 focus:ring-primary outline-none"
          @keyup.enter="$emit('search')"
        />
      </div>
      <div class="relative flex-1 max-w-[180px]">
        <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]">tag</span>
        <input
          v-model="searchInfo.status"
          type="text"
          :placeholder="t('statusCode')"
          class="w-full pl-10 pr-4 py-2 bg-slate-50 dark:bg-zinc-800 border-none rounded-lg text-sm focus:ring-1 focus:ring-primary outline-none"
          @keyup.enter="$emit('search')"
        />
      </div>
      <div class="flex gap-2">
        <button
          class="flex items-center gap-2 px-4 py-2 bg-primary hover:bg-primary-hover text-white rounded-lg text-sm font-medium"
          @click="$emit('search')"
        >
          <span class="material-icons text-[18px]">search</span>
          {{ t('searchQuery') }}
        </button>
        <button
          class="flex items-center gap-2 px-4 py-2 bg-white dark:bg-zinc-800 border border-border-light dark:border-border-dark hover:bg-slate-50 dark:hover:bg-zinc-700 rounded-lg text-sm font-medium"
          @click="$emit('reset')"
        >
          <span class="material-icons text-[18px]">refresh</span>
          {{ t('reset') }}
        </button>
      </div>
    </div>

    <div class="overflow-x-auto" v-loading="loading">
      <table class="w-full min-w-[1280px] text-left">
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
            <th class="px-6 py-4 text-right">{{ t('actions') }}</th>
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
              <span class="font-medium">{{ row.user?.userName || '-' }}</span>
              <span v-if="row.user?.nickName" class="text-slate-400">({{ row.user.nickName }})</span>
            </td>
            <td class="px-6 py-4 text-sm text-slate-500">{{ formatDate(row.CreatedAt) }}</td>
            <td class="px-6 py-4">
              <span :class="getStatusClass(row.status)" class="px-2.5 py-1 rounded-full text-xs font-bold">
                {{ row.status }}
              </span>
            </td>
            <td class="px-6 py-4 text-sm font-mono text-slate-500">{{ row.ip }}</td>
            <td class="px-6 py-4">
              <span :class="getMethodClass(row.method)" class="px-2 py-0.5 rounded text-xs font-bold">
                {{ row.method }}
              </span>
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
            <td class="px-6 py-4 text-right">
              <div class="flex justify-end items-center">
                <button
                  class="bg-red-500/10 hover:bg-red-500/20 text-red-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1"
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

    <div class="px-6 py-4 bg-slate-50 dark:bg-zinc-800/30 flex items-center justify-between border-t border-border-light dark:border-border-dark">
      <span class="text-xs text-slate-500">{{ t('totalRecords', { total }) }}</span>
      <el-pagination
        :current-page="page"
        :page-size="pageSize"
        :page-sizes="[10, 30, 50, 100]"
        :total="total"
        background
        layout="sizes, prev, pager, next, jumper"
        @current-change="$emit('page-change', $event)"
        @size-change="$emit('size-change', $event)"
      />
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'
import { formatDate } from '@/utils/format'

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

const getMethodClass = (method) => {
  const classes = {
    GET: 'bg-primary/10 text-primary',
    POST: 'bg-emerald-500/10 text-emerald-500',
    PUT: 'bg-amber-500/10 text-amber-500',
    DELETE: 'bg-red-500/10 text-red-500'
  }

  return classes[method] || 'bg-slate-500/10 text-slate-500'
}

const getStatusClass = (status) => {
  if (Number(status) >= 400) {
    return 'bg-red-500/10 text-red-500'
  }

  return 'bg-emerald-500/10 text-emerald-500'
}

const formatBody = (value) => {
  try {
    return JSON.stringify(JSON.parse(value), null, 2)
  } catch {
    return value
  }
}
</script>
