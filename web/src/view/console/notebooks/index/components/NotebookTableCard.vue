<template>
  <TableCard :page-size="pageSize">
    <template #toolbar>
      <BaseFilterBar plain wrap-main actions-class="list-toolbar-actions--push">
        <div class="list-search-field">
          <input
            :value="searchQuery"
            :placeholder="t('searchInstancePlaceholder')"
            class="list-search-input !w-full"
            type="text"
            @input="handleSearchInput"
            @keyup.enter="$emit('refresh')"
          />
          <span
            class="material-icons absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400 text-[16px]"
            >search</span
          >
        </div>

        <el-select
          :model-value="statusFilter"
          :placeholder="`${t('status')}: ${t('all')}`"
          class="list-filter-select !w-[168px]"
          clearable
          size="small"
          @update:model-value="$emit('status-change', $event)"
        >
          <el-option :label="t('RUNNING')" value="RUNNING" />
          <el-option :label="t('STOPPED')" value="STOPPED" />
          <el-option :label="t('CREATING')" value="CREATING" />
        </el-select>

        <template #actions>
          <button
            class="list-toolbar-button list-toolbar-button--secondary list-toolbar-button--icon"
            @click="$emit('show-key-settings')"
          >
            <span class="material-icons text-[20px]">vpn_key</span>
          </button>
        </template>
      </BaseFilterBar>
    </template>

    <div
      class="console-table-scroll console-table-scroll--fill overflow-x-auto"
    >
      <table
        class="console-table console-table--resource-dense w-full min-w-[1180px]"
      >
        <thead>
          <tr>
            <th>{{ t('name') }} / {{ t('instanceId') }}</th>
            <th>{{ t('status') }}</th>
            <th>{{ t('spec') }}</th>
            <th>{{ t('gpu') }}</th>
            <th class="text-center">{{ t('sshLogin') }}</th>
            <th class="text-center">{{ t('quickTools') }}</th>
            <th class="console-actions-header">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td colspan="7" class="px-6 py-12 text-center text-slate-400">
              <div class="flex items-center justify-center gap-2">
                <div
                  class="animate-spin rounded-full h-5 w-5 border-b-2 border-primary"
                ></div>
                {{ t('loading') }}
              </div>
            </td>
          </tr>
          <tr v-else-if="notebooks.length === 0">
            <td colspan="7" class="px-6 py-12 text-center text-slate-400">
              <div class="console-list-empty">
                <span class="material-icons console-list-empty__icon"
                  >inbox</span
                >
                <p class="console-list-empty__text">{{ t('noData') }}</p>
                <button
                  class="list-row-button list-row-button--neutral"
                  @click="$emit('create')"
                >
                  {{ t('create') }}{{ t('notebook') }}
                </button>
              </div>
            </td>
          </tr>
          <tr v-for="item in notebooks" v-else :key="item.id" class="group">
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
              <ResourceStatusCell
                :label="getStatusText(item.status)"
                show-pulse
                :pulse-class="
                  item.status === 'RUNNING'
                    ? 'animate-pulse bg-emerald-500'
                    : 'bg-current'
                "
                :status-class="getStatusStyle(item.status)"
              />
            </td>
            <td>
              <div class="flex items-center gap-1">
                <span class="material-icons text-[14px] text-slate-400"
                  >memory</span
                >
                <span class="is-metric">{{ item.cpu }} {{ t('cpu') }}</span>
                <span class="text-slate-300 mx-1">|</span>
                <span class="is-metric">{{ item.memory }} GB</span>
              </div>
            </td>
            <td>
              <span
                v-if="item.gpuCount > 0 || item.gpu > 0"
                class="console-badge console-badge--neutral is-code"
              >
                {{ item.gpuCount || item.gpu }} × {{ item.gpuModel || 'GPU' }}
              </span>
              <span
                v-else-if="hasVGpuSpec(item)"
                class="console-badge console-badge--neutral is-code"
              >
                vGPU · {{ formatVGpuSpec(item) }}
              </span>
              <span v-else class="is-muted text-[11px]">{{
                t('CPU ONLY')
              }}</span>
            </td>
            <td class="text-center">
              <button
                v-if="item.status === 'RUNNING'"
                class="list-row-button list-row-button--neutral min-w-[88px]"
                @click="$emit('show-ssh', item)"
              >
                <span class="material-icons text-[14px]">terminal</span>
                {{ t('connect') }}
              </button>
              <span v-else class="text-slate-400 text-sm">-</span>
            </td>
            <td class="text-center">
              <div class="console-resource-tools">
                <a
                  v-if="item.status === 'RUNNING' && item.jupyterUrl"
                  :href="item.jupyterUrl"
                  class="list-row-button list-row-button--warning min-w-[112px] justify-start"
                  target="_blank"
                >
                  <span class="material-icons text-[14px]">science</span>
                  <span class="flex-1 text-center">Jupyter</span>
                </a>
                <a
                  v-if="
                    item.status === 'RUNNING' &&
                    item.enableTensorboard &&
                    item.tensorboardUrl
                  "
                  :href="item.tensorboardUrl"
                  class="list-row-button list-row-button--info min-w-[112px] justify-start"
                  target="_blank"
                >
                  <span class="material-icons text-[14px]">analytics</span>
                  <span class="flex-1 text-center pr-[14px]">TensorBoard</span>
                </a>
                <span
                  v-if="
                    item.status !== 'RUNNING' ||
                    (!item.jupyterUrl &&
                      !(item.enableTensorboard && item.tensorboardUrl))
                  "
                  class="is-muted text-sm"
                  >-</span
                >
              </div>
            </td>
            <td class="console-actions-cell">
              <div class="list-row-actions">
                <button
                  v-if="item.status === 'STOPPED'"
                  :disabled="btnLoading[item.id]"
                  class="list-row-button list-row-button--success disabled:opacity-50"
                  @click="$emit('start', item)"
                >
                  <span class="material-icons text-[14px]">play_arrow</span>
                  {{ t('start') }}
                </button>
                <button
                  v-if="item.status === 'RUNNING'"
                  :disabled="btnLoading[item.id]"
                  class="list-row-button list-row-button--warning disabled:opacity-50"
                  @click="$emit('stop', item)"
                >
                  <span class="material-icons text-[14px]">stop</span>
                  {{ t('stop') }}
                </button>
                <button
                  v-if="['STOPPED', 'CLOSED'].includes(normalizeStatus(item.status))"
                  :disabled="btnLoading[item.id]"
                  class="list-row-button list-row-button--info disabled:opacity-50"
                  @click="$emit('edit', item)"
                >
                  <span class="material-icons text-[14px]">edit</span>
                  {{ t('edit') }}
                </button>
                <button
                  v-if="item.status !== 'DELETING'"
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
        :page-sizes="[15, 20, 50, 100]"
        layout="sizes, prev, pager, next, jumper"
        @current-change="$emit('page-change', $event)"
        @size-change="$emit('size-change', $event)"
      />
    </template>
  </TableCard>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import BaseFilterBar from '@/components/listPage/BaseFilterBar.vue'
import ListPaginationBar from '@/components/listPage/ListPaginationBar.vue'
import ResourceIdentityCell from '@/components/listPage/ResourceIdentityCell.vue'
import ResourceStatusCell from '@/components/listPage/ResourceStatusCell.vue'
import TableCard from '@/components/listPage/TableCard.vue'
import type { ConsoleNotebook, Translator } from '@/types/consoleResource'
import { formatVGpuSpec, hasVGpuSpec } from '@/utils/vgpu'

const normalizeStatus = (status?: string): string =>
  `${status || ''}`.trim().toUpperCase()

withDefaults(
  defineProps<{
    btnLoading?: Record<string | number, boolean>
    getStatusStyle: (status?: string) => string
    getStatusText: (status?: string) => string
    loading?: boolean
    notebooks?: ConsoleNotebook[]
    page?: number
    pageSize?: number
    searchQuery?: string
    statusFilter?: string
    total?: number
  }>(),
  {
    btnLoading: () => ({}),
    loading: false,
    notebooks: () => [],
    page: 1,
    pageSize: 15,
    searchQuery: '',
    statusFilter: '',
    total: 0
  }
)

const emit = defineEmits<{
  copy: [value: string]
  create: []
  delete: [item: ConsoleNotebook]
  detail: [item: ConsoleNotebook]
  edit: [item: ConsoleNotebook]
  'page-change': [page: number]
  refresh: []
  'search-change': [value: string]
  'show-key-settings': []
  'show-ssh': [item: ConsoleNotebook]
  'size-change': [pageSize: number]
  start: [item: ConsoleNotebook]
  'status-change': [value?: string]
  stop: [item: ConsoleNotebook]
}>()

const t = inject<Translator>('t', (key: string) => key)

const handleSearchInput = (event: Event) => {
  emit('search-change', (event.target as HTMLInputElement).value)
}
</script>

