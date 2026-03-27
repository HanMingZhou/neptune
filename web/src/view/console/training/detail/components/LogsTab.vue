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
          style="width: 150px;"
          @update:model-value="$emit('update:log-task-name', $event)"
        >
          <el-option v-if="frameworkType !== 'STANDALONE'" :label="masterTaskLabel" :value="masterTaskName" />
          <el-option v-if="workerCount > 0 || frameworkType === 'STANDALONE'" label="Worker" value="worker" />
        </el-select>
        <el-select
          :disabled="logsConnected"
          :model-value="podIndex"
          :placeholder="t('podIndex')"
          style="width: 120px;"
          @update:model-value="$emit('update:pod-index', $event)"
        >
          <el-option v-for="n in podCount" :key="n - 1" :label="`${t('instances')} ${n - 1}`" :value="n - 1" />
        </el-select>
    </template>
    <template #controls-suffix>
        <button
          v-if="isTerminal"
          :disabled="downloadLoading"
          class="ml-auto px-4 py-2 bg-emerald-500 hover:bg-emerald-600 text-white rounded-lg text-sm font-bold flex items-center gap-2 disabled:opacity-50"
          @click="$emit('download')"
        >
          <span class="material-icons text-lg">download</span>
          {{ downloadLoading ? t('loading') : t('downloadFullLogs') }}
        </button>
    </template>
  </LogStreamPanel>
</template>

<script setup>
import { inject } from 'vue'
import LogStreamPanel from '@/components/detailPage/LogStreamPanel.vue'

defineProps({
  downloadLoading: {
    type: Boolean,
    default: false
  },
  frameworkType: {
    type: String,
    default: ''
  },
  isTerminal: {
    type: Boolean,
    default: false
  },
  logTaskName: {
    type: String,
    default: ''
  },
  logs: {
    type: String,
    default: ''
  },
  logsConnected: {
    type: Boolean,
    default: false
  },
  logsLoading: {
    type: Boolean,
    default: false
  },
  masterTaskLabel: {
    type: String,
    default: ''
  },
  masterTaskName: {
    type: String,
    default: ''
  },
  podCount: {
    type: Number,
    default: 0
  },
  podIndex: {
    type: Number,
    default: 0
  },
  setLogsRef: {
    type: Function,
    required: true
  },
  workerCount: {
    type: Number,
    default: 0
  }
})

defineEmits([
  'clear',
  'connect',
  'disconnect',
  'download',
  'update:log-task-name',
  'update:pod-index'
])

const t = inject('t', (key) => key)
</script>
