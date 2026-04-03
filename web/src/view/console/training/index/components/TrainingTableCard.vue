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
            :placeholder="`${t('status')}: ${t('all')}`"
            class="list-filter-select !w-[168px]"
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
      <table class="console-table console-table--resource-dense w-full min-w-[1180px]">
        <thead>
          <tr>
            <th>{{ t('name') }}</th>
            <th>{{ t('spec') }}</th>
            <th>{{ t('status') }}</th>
            <th>{{ t('gpu') }}</th>
            <th>{{ t('workerCount') }}</th>
            <th>{{ t('createdAt') }}</th>
            <th>{{ t('duration') }}</th>
            <th class="console-actions-header">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td colspan="8" class="px-6 py-10 text-center text-slate-400">
              <div class="flex items-center justify-center gap-2">
                <div class="animate-spin rounded-full h-5 w-5 border-b-2 border-primary"></div>
                {{ t('loading') }}
              </div>
            </td>
          </tr>
          <tr v-else-if="jobs.length === 0">
            <td colspan="8" class="px-6 py-10 text-center text-slate-400">
              <div class="console-list-empty">
                <span class="material-icons console-list-empty__icon">model_training</span>
                <p class="console-list-empty__text">{{ t('noData') }}</p>
                <button class="list-row-button list-row-button--neutral" @click="$emit('create')">{{ t('createTraining') }}</button>
              </div>
            </td>
          </tr>
          <tr
            v-for="item in jobs"
            v-else
            :key="item.id"
            class="group"
          >
            <td>
              <ResourceIdentityCell
                :copy-title="t('copy')"
                :copy-value="item.displayName && item.displayName !== item.jobName ? item.jobName : ''"
                :primary="item.displayName || item.jobName"
                :secondary="item.jobName"
                @copy="$emit('copy', $event)"
                @primary-click="$emit('detail', item)"
              />
            </td>
            <td>
              <span :class="getFrameworkStyle(item.frameworkType)" class="console-badge text-[11px] border border-transparent">
                {{ getFrameworkLabel(item.frameworkType) }}
              </span>
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
              <span v-if="item.totalGpuCount" class="console-badge console-badge--neutral is-code">{{ item.totalGpuCount }} GPU</span>
              <span v-else class="is-muted text-[11px]">{{ t('CPU ONLY') }}</span>
            </td>
            <td class="font-mono text-xs is-metric">
              {{ item.workerCount || 1 }}x
            </td>
            <td class="text-slate-500 text-[12px] font-mono is-secondary">
              {{ formatTime(item.createdAt) }}
            </td>
            <td class="font-mono text-xs is-metric">
              {{ item.duration || '-' }}
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
                  v-if="item.enableTensorboard"
                  :disabled="!item.tensorboardUrl"
                  class="list-row-button list-row-button--warning disabled:opacity-40 disabled:cursor-not-allowed"
                  @click="$emit('open-tensorboard', item)"
                >
                  <span class="material-icons text-[14px]">assessment</span>
                  TensorBoard
                </button>
                <button
                  v-if="['RUNNING', 'PENDING', 'CREATING'].includes(item.status)"
                  :disabled="btnLoading[item.id]"
                  class="list-row-button list-row-button--warning disabled:opacity-50"
                  @click="$emit('stop', item)"
                >
                  <span class="material-icons text-[14px]">stop</span>
                  {{ t('stop') }}
                </button>
                <button
                  v-if="['SUCCEEDED', 'FAILED', 'KILLED'].includes(item.status)"
                  :disabled="btnLoading[item.id]"
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
