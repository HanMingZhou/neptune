<template>
  <TableCard>
    <template #toolbar>
      <div class="flex flex-wrap gap-4 items-center">
        <div class="flex items-center gap-3 flex-1">
          <div class="relative">
            <input
              :value="searchQuery"
              :placeholder="t('searchInstancePlaceholder')"
              class="list-search-input"
              type="text"
              @input="$emit('search-change', $event.target.value)"
              @keyup.enter="$emit('refresh')"
            />
            <span class="material-icons absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400 text-[16px]">search</span>
          </div>

          <el-select
            :model-value="filterStatus"
            :placeholder="t('statusAll')"
            class="list-filter-select"
            clearable
            size="small"
            @update:model-value="$emit('status-change', $event)"
          >
            <el-option :label="t('RUNNING')" value="RUNNING" />
            <el-option :label="t('PENDING')" value="PENDING" />
            <el-option :label="t('STOPPED')" value="STOPPED" />
            <el-option :label="t('FAILED')" value="FAILED" />
          </el-select>

          <el-select
            :model-value="filterFramework"
            :placeholder="t('frameworkAll')"
            class="list-filter-select"
            clearable
            size="small"
            @update:model-value="$emit('framework-change', $event)"
          >
            <el-option label="SGLang" value="SGLANG" />
            <el-option label="vLLM" value="VLLM" />
          </el-select>
        </div>
      </div>
    </template>

    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="bg-gray-50 dark:bg-zinc-900/50 border-b border-gray-100 dark:border-gray-800 text-slate-500 text-[11px] font-bold uppercase tracking-wider">
            <th class="px-6 py-3">{{ t('name') }}</th>
            <th class="px-6 py-3">{{ t('framework') }}</th>
            <th class="px-6 py-3">{{ t('status') }}</th>
            <th class="px-6 py-3">{{ t('gpu') }}</th>
            <th class="px-6 py-3">{{ t('deployMode') }}</th>
            <th class="px-6 py-3">{{ t('createdAt') }}</th>
            <th class="px-6 py-3 text-center">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-100 dark:divide-border-dark text-[13px]">
          <tr v-if="loading">
            <td colspan="7" class="px-6 py-10 text-center text-slate-400">
              <div class="flex items-center justify-center gap-2">
                <div class="animate-spin rounded-full h-5 w-5 border-b-2 border-primary"></div>
                {{ t('loading') }}
              </div>
            </td>
          </tr>
          <tr v-else-if="services.length === 0">
            <td colspan="7" class="px-6 py-10 text-center text-slate-400">
              <div class="space-y-4 flex flex-col items-center">
                <span class="material-icons text-6xl text-slate-200 dark:text-zinc-700">rocket_launch</span>
                <p>{{ t('noData') }}</p>
                <button class="text-primary hover:underline font-bold" @click="$emit('create')">{{ t('createInference') }}</button>
              </div>
            </td>
          </tr>
          <tr
            v-for="item in services"
            v-else
            :key="item.id"
            class="hover:bg-gray-50 dark:hover:bg-zinc-800/50 transition-colors group"
          >
            <td class="px-6 py-3 text-center">
              <div class="flex flex-col gap-1 items-center">
                <span class="font-bold text-primary hover:underline cursor-pointer text-sm" @click="$emit('detail', item)">
                  {{ item.displayName }}
                </span>
                <div class="flex items-center gap-1 group/copy">
                  <span class="text-[11px] font-mono text-slate-400">{{ item.instanceName }}</span>
                  <button
                    :title="t('copy')"
                    class="opacity-0 group-hover/copy:opacity-100 transition-opacity p-0.5 rounded hover:bg-slate-100 dark:hover:bg-zinc-700 text-slate-400 hover:text-primary"
                    @click.stop="$emit('copy', item.instanceName)"
                  >
                    <span class="material-icons text-[12px]">content_copy</span>
                  </button>
                </div>
              </div>
            </td>
            <td class="px-6 py-3 text-center">
              <span v-if="item.framework" :class="getFrameworkStyle(item.framework)" class="px-2.5 py-1 rounded-full text-[11px] font-bold border border-transparent">
                {{ item.framework }}
              </span>
              <span v-else class="text-slate-400 text-[11px]">-</span>
            </td>
            <td class="px-6 py-3 text-center">
              <div class="flex justify-center">
                <span :class="getStatusStyle(item.status)" class="px-2.5 py-1 rounded-full text-[11px] font-bold flex items-center gap-1.5 border border-transparent">
                  <span v-if="item.status === 'RUNNING'" class="size-1.5 rounded-full animate-pulse bg-emerald-500"></span>
                  {{ getStatusLabel(item.status) }}
                </span>
              </div>
            </td>
            <td class="px-6 py-3 text-center">
              <span v-if="item.gpu" class="bg-slate-100 dark:bg-zinc-800 px-2 py-0.5 rounded text-[11px] font-mono font-bold">{{ item.gpu }} GPU</span>
              <span v-else class="text-slate-400 text-[11px]">CPU ONLY</span>
            </td>
            <td class="px-6 py-3 text-center">
              <span :class="getDeployTypeStyle(item.deployType)" class="px-2 py-0.5 rounded text-[11px] font-bold">
                {{ item.deployType === 'STANDALONE' ? t('standalone') : t('distributed') }}
              </span>
            </td>
            <td class="px-6 py-3 text-center text-slate-500 text-[12px] font-mono">
              {{ formatTime(item.createdAt) }}
            </td>
            <td class="px-6 py-3 text-center">
              <div class="flex justify-center gap-2 items-center">
                <button class="bg-primary/10 hover:bg-primary/20 text-primary px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1" @click="$emit('detail', item)">
                  <span class="material-icons text-[14px]">info</span>
                  {{ t('details') }}
                </button>
                <button class="bg-blue-500/10 hover:bg-blue-500/20 text-blue-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1" @click="$emit('logs', item)">
                  <span class="material-icons text-[14px]">list_alt</span>
                  {{ t('logs') }}
                </button>
                <button
                  v-if="['RUNNING', 'PENDING'].includes(item.status)"
                  :disabled="btnLoading[item.id]"
                  class="bg-amber-500/10 hover:bg-amber-500/20 text-amber-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1 disabled:opacity-50 w-[60px] justify-center"
                  @click="$emit('stop', item)"
                >
                  <span class="material-icons text-[14px]">stop</span>
                  {{ t('stop') }}
                </button>
                <button
                  v-else
                  :disabled="!['STOPPED', 'FAILED'].includes(item.status) || btnLoading[item.id]"
                  :class="{ 'hover:bg-emerald-500/20': ['STOPPED', 'FAILED'].includes(item.status) && !btnLoading[item.id] }"
                  class="bg-emerald-500/10 text-emerald-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1 disabled:opacity-30 disabled:cursor-not-allowed w-[60px] justify-center"
                  @click="$emit('start', item)"
                >
                  <span class="material-icons text-[14px]">play_arrow</span>
                  {{ t('start') }}
                </button>
                <button
                  :disabled="!['STOPPED', 'FAILED'].includes(item.status) || btnLoading[item.id]"
                  :class="{ 'hover:bg-red-500/20': ['STOPPED', 'FAILED'].includes(item.status) && !btnLoading[item.id] }"
                  class="bg-red-500/10 text-red-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1 disabled:opacity-30 disabled:cursor-not-allowed w-[60px] justify-center"
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
      <div class="flex items-center justify-between text-xs text-slate-500">
        <span>{{ t('totalRecords', { total }) }}</span>
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 20, 50]"
          :total="total"
          class="!p-0"
          layout="sizes, prev, pager, next"
          size="small"
          @current-change="$emit('page-change', $event)"
          @size-change="$emit('size-change', $event)"
        />
      </div>
    </template>
  </TableCard>
</template>

<script setup>
import { inject } from 'vue'
import TableCard from '@/components/listPage/TableCard.vue'

defineProps({
  btnLoading: {
    type: Object,
    default: () => ({})
  },
  filterFramework: {
    type: String,
    default: ''
  },
  filterStatus: {
    type: String,
    default: ''
  },
  getDeployTypeStyle: {
    type: Function,
    required: true
  },
  getFrameworkStyle: {
    type: Function,
    required: true
  },
  getStatusLabel: {
    type: Function,
    required: true
  },
  getStatusStyle: {
    type: Function,
    required: true
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
  searchQuery: {
    type: String,
    default: ''
  },
  services: {
    type: Array,
    default: () => []
  },
  total: {
    type: Number,
    default: 0
  }
})

defineEmits([
  'copy',
  'create',
  'delete',
  'detail',
  'framework-change',
  'logs',
  'page-change',
  'refresh',
  'search-change',
  'size-change',
  'start',
  'status-change',
  'stop'
])

const t = inject('t', (key) => key)

const formatTime = (time) => {
  if (!time) return '-'
  const date = new Date(time)
  return `${date.getMonth() + 1}/${date.getDate()} ${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
}
</script>
