<template>
  <div class="console-page-container px-6 pb-6 pt-1 md:pt-2">
    <OverviewTab
      v-show="activeTab === 'overview'"
      :api-endpoint="apiEndpoint"
      :curl-example="curlExample"
      :format-command="formatCommand"
      :format-time="formatTime"
      :get-status-label="getStatusLabel"
      :service="service"
      @copy="emit('copy', $event)"
      @manage-api-key="emit('manage-api-key')"
    />

    <LogsTab
      v-show="activeTab === 'logs'"
      :logs="logs"
      :logs-connected="logsConnected"
      :logs-loading="logsLoading"
      :pods="pods"
      :selected-pod="selectedPod"
      :set-logs-ref="setLogsRef"
      @clear="emit('clear-logs')"
      @connect="emit('connect-logs')"
      @disconnect="emit('disconnect-logs')"
      @update:selected-pod="emit('update:selected-pod', $event)"
    />

    <TerminalTab
      v-show="activeTab === 'terminal'"
      :can-connect-terminal="canConnectTerminal"
      :pods="pods"
      :service="service"
      :set-terminal-ref="setTerminalRef"
      :terminal-connected="terminalConnected"
      :terminal-pod="terminalPod"
      @connect="emit('connect-terminal')"
      @disconnect="emit('disconnect-terminal')"
      @fit="emit('fit-terminal')"
      @update:terminal-pod="emit('update:terminal-pod', $event)"
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
import type {
  ConsoleInferenceDetail,
  ConsolePod
} from '@/types/consoleResource'

type FormatTime = (value?: string | number | null) => string
type StatusLabelGetter = (status?: string) => string
type StatusClassGetter = (status?: string) => string
type FormatCommand = (service?: Partial<ConsoleInferenceDetail>) => string
type ElementRefSetter = (element: HTMLElement | null) => void

defineProps<{
  activeTab: string
  apiEndpoint: string
  canConnectTerminal: boolean
  curlExample: string
  formatCommand: FormatCommand
  formatTime: FormatTime
  getPodStatusClass: StatusClassGetter
  getStatusLabel: StatusLabelGetter
  logs: string
  logsConnected: boolean
  logsLoading: boolean
  pods: ConsolePod[]
  podsLoading: boolean
  selectedPod: string
  service: Partial<ConsoleInferenceDetail>
  setLogsRef: ElementRefSetter
  setTerminalRef: ElementRefSetter
  terminalConnected: boolean
  terminalPod: string
}>()

const emit = defineEmits<{
  'clear-logs': []
  'connect-logs': []
  'connect-terminal': []
  copy: [value: string]
  'disconnect-logs': []
  'disconnect-terminal': []
  'fit-terminal': []
  'manage-api-key': []
  'refresh-pods': []
  'update:selected-pod': [value: string]
  'update:terminal-pod': [value: string]
  'view-pod-logs': [pod: ConsolePod]
}>()
</script>
