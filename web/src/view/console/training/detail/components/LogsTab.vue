<template>
  <LogStreamPanel
    :clear-button-class="'text-slate-400 hover:text-slate-600 dark:hover:text-slate-300'"
    :logs="logs"
    :logs-connected="logsConnected"
    :logs-loading="logsLoading"
    :set-logs-ref="setLogsRef"
    :toolbar-class="'flex items-center gap-4 flex-wrap'"
    @clear="$emit('clear')"
    @connect="$emit('connect')"
    @disconnect="$emit('disconnect')"
  >
    <template #controls-prefix>
      <el-select
        :disabled="logsConnected"
        :model-value="logTaskName"
        :placeholder="t('pleaseSelect')"
        style="width: 150px"
        @update:model-value="$emit('update:log-task-name', $event)"
      >
        <el-option
          v-if="frameworkType !== 'STANDALONE'"
          :label="masterTaskLabel"
          :value="masterTaskName"
        />
        <el-option
          v-if="workerCount > 0 || frameworkType === 'STANDALONE'"
          label="Worker"
          value="worker"
        />
      </el-select>
      <el-select
        :disabled="logsConnected"
        :model-value="podIndex"
        :placeholder="t('podIndex')"
        style="width: 120px"
        @update:model-value="$emit('update:pod-index', $event)"
      >
        <el-option
          v-for="n in podCount"
          :key="n - 1"
          :label="`${t('instances')} ${n - 1}`"
          :value="n - 1"
        />
      </el-select>
    </template>
    <template #controls-suffix>
      <button
        v-if="isTerminal"
        :disabled="downloadLoading"
        class="detail-header-action detail-header-action--primary ml-auto disabled:opacity-50"
        @click="$emit('download')"
      >
        <span class="material-icons text-lg">download</span>
        {{ downloadLoading ? t('loading') : t('downloadFullLogs') }}
      </button>
    </template>
  </LogStreamPanel>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type { Translator } from '@/types/consoleResource'
import LogStreamPanel from '@/components/detailPage/LogStreamPanel.vue'

withDefaults(
  defineProps<{
    downloadLoading?: boolean
    frameworkType?: string
    isTerminal?: boolean
    logTaskName?: string
    logs?: string
    logsConnected?: boolean
    logsLoading?: boolean
    masterTaskLabel?: string
    masterTaskName?: string
    podCount?: number
    podIndex?: number
    setLogsRef: (element: HTMLElement | null) => void
    workerCount?: number
  }>(),
  {
    downloadLoading: false,
    frameworkType: '',
    isTerminal: false,
    logTaskName: '',
    logs: '',
    logsConnected: false,
    logsLoading: false,
    masterTaskLabel: '',
    masterTaskName: '',
    podCount: 0,
    podIndex: 0,
    workerCount: 0
  }
)

defineEmits<{
  clear: []
  connect: []
  disconnect: []
  download: []
  'update:log-task-name': [value: string]
  'update:pod-index': [value: number]
}>()

const t = inject<Translator>('t', (key: string) => key)
</script>
