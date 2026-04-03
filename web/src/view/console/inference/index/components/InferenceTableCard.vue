<template>
  <TableCard>
    <template #toolbar>
      <div class="list-filter-bar">
        <div class="list-toolbar-main">
          <div class="list-search-field">
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
            class="list-filter-select !w-[168px]"
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
            class="list-filter-select !w-[168px]"
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
      <table class="console-table console-table--resource-dense w-full min-w-[1180px]">
        <thead>
          <tr>
            <th>{{ t('name') }}</th>
            <th>{{ t('framework') }}</th>
            <th>{{ t('status') }}</th>
            <th>{{ t('gpu') }}</th>
            <th>{{ t('deployMode') }}</th>
            <th>{{ t('createdAt') }}</th>
            <th class="console-actions-header">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody>
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
              <div class="console-list-empty">
                <span class="material-icons console-list-empty__icon">rocket_launch</span>
                <p class="console-list-empty__text">{{ t('noData') }}</p>
                <button class="list-row-button list-row-button--neutral" @click="$emit('create')">{{ t('createInference') }}</button>
              </div>
            </td>
          </tr>
          <tr
            v-for="item in services"
            v-else
            :key="item.id"
            class="group"
          >
            <td>
              <ResourceIdentityCell
                :copy-title="t('copy')"
                :copy-value="item.instanceName"
                :primary="item.displayName || item.instanceName"
                :secondary="item.instanceName"
                @copy="$emit('copy', $event)"
                @primary-click="$emit('detail', item)"
              />
            </td>
            <td>
              <span v-if="item.framework" :class="getFrameworkStyle(item.framework)" class="console-badge text-[11px] border border-transparent">
                {{ item.framework }}
              </span>
              <span v-else class="is-muted text-[11px]">-</span>
            </td>
            <td>
              <ResourceStatusCell
                :label="getStatusLabel(item.status)"
                :show-pulse="item.status === 'RUNNING'"
                pulse-class="animate-pulse bg-emerald-500"
                :status-class="getStatusStyle(item.status)"
              />
            </td>
            <td>
              <span v-if="item.gpu" class="console-badge console-badge--neutral is-code">{{ item.gpu }} GPU</span>
              <span v-else class="is-muted text-[11px]">CPU ONLY</span>
            </td>
            <td>
              <span :class="getDeployTypeStyle(item.deployType)" class="console-badge text-[11px]">
                {{ item.deployType === 'STANDALONE' ? t('standalone') : t('distributed') }}
              </span>
            </td>
            <td class="text-slate-500 text-[12px] font-mono is-secondary">
              {{ formatTime(item.createdAt) }}
            </td>
            <td class="console-actions-cell">
              <div class="list-row-actions">
                <button class="list-row-button list-row-button--neutral" @click="$emit('detail', item)">
                  <span class="material-icons text-[14px]">info</span>
                  {{ t('details') }}
                </button>
                <button class="list-row-button list-row-button--info" @click="$emit('logs', item)">
                  <span class="material-icons text-[14px]">list_alt</span>
                  {{ t('logs') }}
                </button>
                <button
                  v-if="['RUNNING', 'PENDING'].includes(item.status)"
                  :disabled="btnLoading[item.id]"
                  class="list-row-button list-row-button--warning disabled:opacity-50 min-w-[64px]"
                  @click="$emit('stop', item)"
                >
                  <span class="material-icons text-[14px]">stop</span>
                  {{ t('stop') }}
                </button>
                <button
                  v-else
                  :disabled="!['STOPPED', 'FAILED'].includes(item.status) || btnLoading[item.id]"
                  class="list-row-button list-row-button--success disabled:opacity-30 disabled:cursor-not-allowed min-w-[64px]"
                  @click="$emit('start', item)"
                >
                  <span class="material-icons text-[14px]">play_arrow</span>
                  {{ t('start') }}
                </button>
                <button
                  :disabled="!['STOPPED', 'FAILED'].includes(item.status) || btnLoading[item.id]"
                  class="list-row-button list-row-button--danger disabled:opacity-30 disabled:cursor-not-allowed min-w-[64px]"
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
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        :total-text="t('totalRecords', { total })"
        :page-sizes="[10, 20, 50]"
        layout="sizes, prev, pager, next, jumper"
        @current-change="$emit('page-change', $event)"
        @size-change="$emit('size-change', $event)"
      />
    </template>
  </TableCard>
</template>

<script setup>
import { inject } from 'vue'
import ListPaginationBar from '@/components/listPage/ListPaginationBar.vue'
import ResourceIdentityCell from '@/components/listPage/ResourceIdentityCell.vue'
import ResourceStatusCell from '@/components/listPage/ResourceStatusCell.vue'
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
