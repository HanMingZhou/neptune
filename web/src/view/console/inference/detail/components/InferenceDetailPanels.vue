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
  apiEndpoint: {
    type: String,
    default: ''
  },
  canConnectTerminal: {
    type: Boolean,
    default: false
  },
  curlExample: {
    type: String,
    default: ''
  },
  formatCommand: {
    type: Function,
    required: true
  },
  formatTime: {
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
  pods: {
    type: Array,
    default: () => []
  },
  podsLoading: {
    type: Boolean,
    default: false
  },
  selectedPod: {
    type: String,
    default: ''
  },
  service: {
    type: Object,
    required: true
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
  terminalPod: {
    type: String,
    default: ''
  }
})

const emit = defineEmits([
  'clear-logs',
  'connect-logs',
  'connect-terminal',
  'copy',
  'disconnect-logs',
  'disconnect-terminal',
  'fit-terminal',
  'manage-api-key',
  'refresh-pods',
  'update:selected-pod',
  'update:terminal-pod',
  'view-pod-logs'
])
</script>
