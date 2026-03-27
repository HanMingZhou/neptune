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
            :placeholder="`${t('status')}: ${t('all')}`"
            class="list-filter-select"
            clearable
            size="small"
            @update:model-value="$emit('status-change', $event)"
          >
            <el-option :label="t('Running')" value="RUNNING" />
            <el-option :label="t('Pending')" value="PENDING" />
            <el-option :label="t('Succeeded')" value="SUCCEEDED" />
            <el-option :label="t('Failed')" value="FAILED" />
            <el-option :label="t('Stopped')" value="KILLED" />
          </el-select>
        </div>
      </div>
    </template>

    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="bg-gray-50 dark:bg-zinc-900/50 border-b border-gray-100 dark:border-gray-800 text-slate-500 text-[11px] font-bold uppercase tracking-wider">
            <th class="px-6 py-4">{{ t('name') }}</th>
            <th class="px-6 py-4">{{ t('spec') }}</th>
            <th class="px-6 py-4">{{ t('status') }}</th>
            <th class="px-6 py-4">{{ t('gpu') }}</th>
            <th class="px-6 py-4">{{ t('workerCount') }}</th>
            <th class="px-6 py-4">{{ t('createdAt') }}</th>
            <th class="px-6 py-4">{{ t('duration') }}</th>
            <th class="px-6 py-4 text-center">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-100 dark:divide-border-dark text-[13px]">
          <tr v-if="loading">
            <td colspan="8" class="px-6 py-12 text-center text-slate-400">
              <div class="flex items-center justify-center gap-2">
                <div class="animate-spin rounded-full h-5 w-5 border-b-2 border-primary"></div>
                {{ t('loading') }}
              </div>
            </td>
          </tr>
          <tr v-else-if="jobs.length === 0">
            <td colspan="8" class="px-6 py-12 text-center text-slate-400">
              <div class="space-y-4 flex flex-col items-center">
                <span class="material-icons text-6xl text-slate-200 dark:text-zinc-700">model_training</span>
                <p>{{ t('noData') }}</p>
                <button class="text-primary hover:underline font-bold" @click="$emit('create')">{{ t('createTraining') }}</button>
              </div>
            </td>
          </tr>
          <tr
            v-for="item in jobs"
            v-else
            :key="item.id"
            class="hover:bg-gray-50 dark:hover:bg-zinc-800/50 transition-colors group"
          >
            <td class="px-6 py-4 text-center">
              <div class="flex flex-col gap-1 items-center">
                <span class="font-bold text-primary hover:underline cursor-pointer text-sm" @click="$emit('detail', item)">
                  {{ item.displayName || item.jobName }}
                </span>
                <div v-if="item.displayName && item.displayName !== item.jobName" class="flex items-center gap-1 group/copy">
                  <span class="text-[11px] font-mono text-slate-400">{{ item.jobName }}</span>
                  <button
                    :title="t('copy')"
                    class="opacity-0 group-hover/copy:opacity-100 transition-opacity p-0.5 rounded hover:bg-slate-100 dark:hover:bg-zinc-700 text-slate-400 hover:text-primary"
                    @click.stop="$emit('copy', item.jobName)"
                  >
                    <span class="material-icons text-[12px]">content_copy</span>
                  </button>
                </div>
              </div>
            </td>
            <td class="px-6 py-4 text-center">
              <span :class="getFrameworkStyle(item.frameworkType)" class="px-2.5 py-1 rounded-full text-[11px] font-bold border border-transparent">
                {{ getFrameworkLabel(item.frameworkType) }}
              </span>
            </td>
            <td class="px-6 py-4 text-center">
              <div class="flex justify-center">
                <span :class="getStatusStyle(item.status)" class="px-2.5 py-1 rounded-full text-[11px] font-bold flex items-center gap-1.5 border border-transparent">
                  <span v-if="item.status === 'RUNNING'" class="size-1.5 rounded-full animate-pulse bg-emerald-500"></span>
                  {{ getStatusLabel(item.status) }}
                </span>
              </div>
            </td>
            <td class="px-6 py-4 text-center">
              <span v-if="item.totalGpuCount" class="bg-slate-100 dark:bg-zinc-800 px-2 py-0.5 rounded text-[11px] font-mono font-bold">{{ item.totalGpuCount }} GPU</span>
              <span v-else class="text-slate-400 text-[11px] font-bold tracking-tight uppercase">{{ t('cpuOnly') }}</span>
            </td>
            <td class="px-6 py-4 text-center font-mono text-xs">
              {{ item.workerCount || 1 }}x
            </td>
            <td class="px-6 py-4 text-center text-slate-500 text-[12px] font-mono">
              {{ formatTime(item.createdAt) }}
            </td>
            <td class="px-6 py-4 text-center font-mono text-xs">
              {{ item.duration || '-' }}
            </td>
            <td class="px-6 py-4 text-center">
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
                  v-if="item.enableTensorboard"
                  :disabled="!item.tensorboardUrl"
                  class="bg-orange-500/10 hover:bg-orange-500/20 text-orange-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1 disabled:opacity-40 disabled:cursor-not-allowed"
                  @click="$emit('open-tensorboard', item)"
                >
                  <span class="material-icons text-[14px]">assessment</span>
                  TensorBoard
                </button>
                <button
                  v-if="['RUNNING', 'PENDING', 'CREATING'].includes(item.status)"
                  :disabled="btnLoading[item.id]"
                  class="bg-amber-500/10 hover:bg-amber-500/20 text-amber-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1 disabled:opacity-50"
                  @click="$emit('stop', item)"
                >
                  <span class="material-icons text-[14px]">stop</span>
                  {{ t('stop') }}
                </button>
                <button
                  v-if="['SUCCEEDED', 'FAILED', 'KILLED'].includes(item.status)"
                  :disabled="btnLoading[item.id]"
                  class="bg-red-500/10 hover:bg-red-500/20 text-red-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1 disabled:opacity-50"
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
  filterStatus: {
    type: String,
    default: ''
  },
  getFrameworkLabel: {
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
  jobs: {
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
  searchQuery: {
    type: String,
    default: ''
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
  'logs',
  'open-tensorboard',
  'page-change',
  'refresh',
  'search-change',
  'size-change',
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
