<template>
  <div class="console-page-container px-6 pb-6 pt-1 md:pt-2">
    <OverviewTab
      v-show="activeTab === 'overview'"
      :format-time="formatTime"
      :get-framework-label="getFrameworkLabel"
      :get-pay-type-label="getPayTypeLabel"
      :get-status-label="getStatusLabel"
      :job="job"
    />

    <LogsTab
      v-show="activeTab === 'logs'"
      :download-loading="downloadLoading"
      :framework-type="job.frameworkType"
      :is-terminal="isTerminal"
      :log-task-name="logTaskName"
      :logs="logs"
      :logs-connected="logsConnected"
      :logs-loading="logsLoading"
      :master-task-label="masterTaskLabel"
      :master-task-name="masterTaskName"
      :pod-count="podCount"
      :pod-index="podIndex"
      :set-logs-ref="setLogsRef"
      :worker-count="job.workerCount"
      @clear="emit('clear-logs')"
      @connect="emit('connect-logs')"
      @disconnect="emit('disconnect-logs')"
      @download="emit('download-logs')"
      @update:log-task-name="emit('update:log-task-name', $event)"
      @update:pod-index="emit('update:pod-index', $event)"
    />

    <TerminalTab
      v-show="activeTab === 'terminal'"
      :framework-type="job.frameworkType"
      :is-running="isRunning"
      :job-name="job.jobName"
      :master-task-label="masterTaskLabel"
      :master-task-name="masterTaskName"
      :set-terminal-ref="setTerminalRef"
      :status-label="getStatusLabel(job.status)"
      :terminal-connected="terminalConnected"
      :terminal-task-name="terminalTaskName"
      :worker-count="job.workerCount"
      @connect="emit('connect-terminal')"
      @disconnect="emit('disconnect-terminal')"
      @fit="emit('fit-terminal')"
      @update:terminal-task-name="emit('update:terminal-task-name', $event)"
    />

    <div v-show="activeTab === 'pods'">
      <PodsTab
        :get-pod-status-class="getPodStatusClass"
        :pods="pods"
        :pods-loading="podsLoading"
        @refresh="emit('refresh-pods')"
        @view-logs="emit('view-pod-logs', $event)"
      />
    </div>
  </div>
</template>

<script setup>
import PodsTab from '@/components/detailPage/PodsTab.vue'
import LogsTab from './LogsTab.vue'
import OverviewTab from './OverviewTab.vue'
import TerminalTab from './TerminalTab.vue'

defineProps({
  activeTab: {
    type: String,
    default: 'overview'
  },
  downloadLoading: {
    type: Boolean,
    default: false
  },
  formatTime: {
    type: Function,
    required: true
  },
  getFrameworkLabel: {
    type: Function,
    required: true
  },
  getPayTypeLabel: {
    type: Function,
    required: true
  },
  getPodStatusClass: {
    type: Function,
    required: true
  },
  getStatusLabel: {
    type: Function,
    required: true
  },
  isRunning: {
    type: Boolean,
    default: false
  },
  isTerminal: {
    type: Boolean,
    default: false
  },
  job: {
    type: Object,
    required: true
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
  logTaskName: {
    type: String,
    default: ''
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
  pods: {
    type: Array,
    default: () => []
  },
  podsLoading: {
    type: Boolean,
    default: false
  },
  setLogsRef: {
    type: Function,
    required: true
  },
  setTerminalRef: {
    type: Function,
    required: true
  },
  terminalConnected: {
    type: Boolean,
    default: false
  },
  terminalTaskName: {
    type: String,
    default: ''
  }
})

const emit = defineEmits([
  'clear-logs',
  'connect-logs',
  'connect-terminal',
  'disconnect-logs',
  'disconnect-terminal',
  'download-logs',
  'fit-terminal',
  'refresh-pods',
  'update:log-task-name',
  'update:pod-index',
  'update:terminal-task-name',
  'view-pod-logs'
])
</script>
