<template>
  <div class="console-page-container px-6 pb-6 pt-1 md:pt-2">
    <OverviewTab
      v-show="activeTab === 'overview'"
      :format-time="formatTime"
      :get-pay-type-label="getPayTypeLabel"
      :get-status-label="getStatusLabel"
      :get-unit-price-label="getUnitPriceLabel"
      :notebook="notebook"
    />

    <LogsTab
      v-show="activeTab === 'logs'"
      :logs="logs"
      :logs-connected="logsConnected"
      :logs-loading="logsLoading"
      :set-logs-ref="setLogsRef"
      @clear="emit('clear-logs')"
      @connect="emit('connect-logs')"
      @disconnect="emit('disconnect-logs')"
    />

    <TerminalTab
      v-show="activeTab === 'terminal'"
      :is-instance-running="isInstanceRunning"
      :notebook="notebook"
      :set-terminal-ref="setTerminalRef"
      :terminal-connected="terminalConnected"
      @connect="emit('connect-terminal')"
      @disconnect="emit('disconnect-terminal')"
      @fit="emit('fit-terminal')"
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
import type { ConsoleNotebookDetail, ConsolePod } from '@/types/consoleResource'

type FormatTime = (value?: string | number | null) => string
type StatusLabelGetter = (status?: string) => string
type PayTypeLabelGetter = (type?: number | string) => string
type StatusClassGetter = (status?: string) => string
type ElementRefSetter = (element: HTMLElement | null) => void

defineProps<{
  activeTab: string
  formatTime: FormatTime
  getPayTypeLabel: PayTypeLabelGetter
  getPodStatusClass: StatusClassGetter
  getStatusLabel: StatusLabelGetter
  getUnitPriceLabel: PayTypeLabelGetter
  isInstanceRunning: boolean
  logs: string
  logsConnected: boolean
  logsLoading: boolean
  notebook: Partial<ConsoleNotebookDetail>
  pods: ConsolePod[]
  podsLoading: boolean
  setLogsRef: ElementRefSetter
  setTerminalRef: ElementRefSetter
  terminalConnected: boolean
}>()

const emit = defineEmits<{
  'clear-logs': []
  'connect-logs': []
  'connect-terminal': []
  'disconnect-logs': []
  'disconnect-terminal': []
  'fit-terminal': []
  'refresh-pods': []
  'view-pod-logs': [pod: ConsolePod]
}>()
</script>
