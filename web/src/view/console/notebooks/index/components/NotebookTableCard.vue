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
            :model-value="statusFilter"
            :placeholder="`${t('status')}: ${t('all')}`"
            class="list-filter-select"
            clearable
            size="small"
            @update:model-value="$emit('status-change', $event)"
          >
            <el-option :label="t('RUNNING')" value="RUNNING" />
            <el-option :label="t('STOPPED')" value="STOPPED" />
            <el-option :label="t('CREATING')" value="CREATING" />
          </el-select>
        </div>

        <div class="ml-auto flex gap-2">
          <button class="p-2 text-slate-400 hover:text-primary transition-colors" @click="$emit('show-key-settings')">
            <span class="material-icons text-[20px]">vpn_key</span>
          </button>
        </div>
      </div>
    </template>

    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="bg-gray-50 dark:bg-zinc-900/50 border-b border-gray-100 dark:border-gray-800 text-slate-500 text-[11px] font-bold uppercase tracking-wider">
            <th class="px-6 py-4">{{ t('instanceId') }} / {{ t('name') }}</th>
            <th class="px-6 py-4">{{ t('status') }}</th>
            <th class="px-6 py-4">{{ t('spec') }}</th>
            <th class="px-6 py-4">{{ t('gpu') }}</th>
            <th class="px-6 py-4 text-center">{{ t('sshLogin') }}</th>
            <th class="px-6 py-4 text-center">{{ t('quickTools') }}</th>
            <th class="px-6 py-4 text-center">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-100 dark:divide-border-dark text-[13px]">
          <tr v-if="loading">
            <td colspan="7" class="px-6 py-12 text-center text-slate-400">
              <div class="flex items-center justify-center gap-2">
                <div class="animate-spin rounded-full h-5 w-5 border-b-2 border-primary"></div>
                {{ t('loading') }}
              </div>
            </td>
          </tr>
          <tr v-else-if="notebooks.length === 0">
            <td colspan="7" class="px-6 py-12 text-center text-slate-400">
              <div class="space-y-4 flex flex-col items-center">
                <span class="material-icons text-6xl text-slate-200 dark:text-zinc-700">inbox</span>
                <p>{{ t('noData') }}</p>
                <button class="text-primary hover:underline font-bold" @click="$emit('create')">{{ t('create') }}{{ t('notebook') }}</button>
              </div>
            </td>
          </tr>
          <tr
            v-for="item in notebooks"
            v-else
            :key="item.id"
            class="hover:bg-gray-50/80 dark:hover:bg-zinc-800/50 transition-colors group"
          >
            <td class="px-6 py-4 text-center">
              <div class="flex flex-col gap-1 items-center">
                <span class="font-bold text-primary hover:underline cursor-pointer text-sm" @click="$emit('detail', item)">{{ item.displayName || item.instanceName }}</span>
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
            <td class="px-6 py-4 text-center">
              <div class="flex justify-center">
                <span :class="getStatusStyle(item.status)" class="px-2.5 py-1 rounded-full text-[11px] font-bold flex items-center gap-1.5 border border-transparent">
                  <span :class="item.status === 'RUNNING' ? 'animate-pulse bg-emerald-500' : 'bg-current'" class="size-1.5 rounded-full"></span>
                  {{ getStatusText(item.status) }}
                </span>
              </div>
            </td>
            <td class="px-6 py-4 text-center">
              <div class="flex items-center justify-center gap-1">
                <span class="material-icons text-[14px] text-slate-400">memory</span>
                <span>{{ item.cpu }} {{ t('cpu') }}</span>
                <span class="text-slate-300 mx-1">|</span>
                <span>{{ item.memory }} GB</span>
              </div>
            </td>
            <td class="px-6 py-4 text-center">
              <span v-if="item.gpu" class="bg-primary/10 text-primary px-2 py-0.5 rounded text-[11px] font-mono font-bold">{{ item.gpu }}</span>
              <span v-else class="text-slate-400 text-[11px] font-bold tracking-tight uppercase">{{ t('cpuOnly') }}</span>
            </td>
            <td class="px-6 py-4 text-center">
              <button
                v-if="item.status === 'RUNNING'"
                class="whitespace-nowrap inline-flex items-center gap-1 text-[11px] font-bold text-slate-500 hover:text-primary transition-colors bg-gray-100 dark:bg-zinc-800 hover:bg-primary/10 px-2 py-1 rounded"
                @click="$emit('show-ssh', item)"
              >
                <span class="material-icons text-[14px]">terminal</span>
                {{ t('connect') }}
              </button>
              <span v-else class="text-slate-400 text-sm">-</span>
            </td>
            <td class="px-6 py-4 text-center">
              <div class="flex flex-col items-center gap-1.5">
                <a
                  v-if="item.status === 'RUNNING' && item.jupyterUrl"
                  :href="item.jupyterUrl"
                  class="whitespace-nowrap inline-flex items-center gap-1.5 text-[11px] font-bold text-orange-600 hover:text-orange-700 bg-orange-500/10 hover:bg-orange-500/20 px-2.5 py-1.5 rounded transition-colors w-[110px] justify-start"
                  target="_blank"
                >
                  <span class="material-icons text-[14px]">science</span>
                  <span class="flex-1 text-center">Jupyter</span>
                </a>
                <a
                  v-if="item.status === 'RUNNING' && item.enableTensorboard && item.tensorboardUrl"
                  :href="item.tensorboardUrl"
                  class="whitespace-nowrap inline-flex items-center gap-1.5 text-[11px] font-bold text-primary hover:text-primary-hover bg-primary/10 hover:bg-primary/20 px-2.5 py-1.5 rounded transition-colors w-[110px] justify-start"
                  target="_blank"
                >
                  <span class="material-icons text-[14px]">analytics</span>
                  <span class="flex-1 text-center pr-[14px]">TensorBoard</span>
                </a>
                <span v-if="item.status !== 'RUNNING' || (!item.jupyterUrl && !(item.enableTensorboard && item.tensorboardUrl))" class="text-slate-400 text-sm">-</span>
              </div>
            </td>
            <td class="px-6 py-4 text-center">
              <div class="flex justify-center gap-2 items-center">
                <button
                  v-if="item.status !== 'RUNNING' && item.status !== 'PENDING' && item.status !== 'CREATING' && item.status !== 'DELETING'"
                  :disabled="btnLoading[item.id]"
                  class="bg-primary/10 hover:bg-primary/20 text-primary px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1 disabled:opacity-50"
                  @click="$emit('start', item)"
                >
                  <span class="material-icons text-[14px]">play_arrow</span>
                  {{ t('start') }}
                </button>
                <button
                  v-if="item.status === 'RUNNING'"
                  :disabled="btnLoading[item.id]"
                  class="bg-amber-500/10 hover:bg-amber-500/20 text-amber-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1 disabled:opacity-50"
                  @click="$emit('stop', item)"
                >
                  <span class="material-icons text-[14px]">stop</span>
                  {{ t('stop') }}
                </button>
                <button
                  v-if="item.status !== 'DELETING'"
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
      <div class="flex items-center justify-between">
        <span class="text-xs text-slate-500">{{ t('totalRecords', { total }) }}</span>
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
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
  getStatusStyle: {
    type: Function,
    required: true
  },
  getStatusText: {
    type: Function,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  },
  notebooks: {
    type: Array,
    default: () => []
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
  statusFilter: {
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
  'page-change',
  'refresh',
  'search-change',
  'show-key-settings',
  'show-ssh',
  'size-change',
  'start',
  'status-change',
  'stop'
])

const t = inject('t', (key) => key)
</script>
