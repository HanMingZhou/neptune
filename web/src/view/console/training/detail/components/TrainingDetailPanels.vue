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

<script setup lang="ts">
import PodsTab from '@/components/detailPage/PodsTab.vue'
import LogsTab from './LogsTab.vue'
import OverviewTab from './OverviewTab.vue'
import TerminalTab from './TerminalTab.vue'
import type { ConsolePod, ConsoleTrainingDetail } from '@/types/consoleResource'

type FormatTime = (value?: string | number | null) => string | null
type StatusLabelGetter = (status?: string) => string
type PayTypeLabelGetter = (type?: number | string) => string
type StatusClassGetter = (status?: string) => string
type TrainingFrameworkGetter = (type?: string) => string
type ElementRefSetter = (element: HTMLElement | null) => void

defineProps<{
  activeTab: string
  downloadLoading: boolean
  formatTime: FormatTime
  getFrameworkLabel: TrainingFrameworkGetter
  getPayTypeLabel: PayTypeLabelGetter
  getPodStatusClass: StatusClassGetter
  getStatusLabel: StatusLabelGetter
  isRunning: boolean
  isTerminal: boolean
  job: Partial<ConsoleTrainingDetail>
  logs: string
  logsConnected: boolean
  logsLoading: boolean
  logTaskName: string
  masterTaskLabel: string
  masterTaskName: string
  podCount: number
  podIndex: number
  pods: ConsolePod[]
  podsLoading: boolean
  setLogsRef: ElementRefSetter
  setTerminalRef: ElementRefSetter
  terminalConnected: boolean
  terminalTaskName: string
}>()

const emit = defineEmits<{
  'clear-logs': []
  'connect-logs': []
  'connect-terminal': []
  'disconnect-logs': []
  'disconnect-terminal': []
  'download-logs': []
  'fit-terminal': []
  'refresh-pods': []
  'update:log-task-name': [value: string]
  'update:pod-index': [value: number]
  'update:terminal-task-name': [value: string]
  'view-pod-logs': [pod: ConsolePod]
}>()
</script>
