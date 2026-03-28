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
  formatTime: {
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
  getUnitPriceLabel: {
    type: Function,
    required: true
  },
  isInstanceRunning: {
    type: Boolean,
    default: false
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
  notebook: {
    type: Object,
    required: true
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
  }
})

const emit = defineEmits([
  'clear-logs',
  'connect-logs',
  'connect-terminal',
  'disconnect-logs',
  'disconnect-terminal',
  'fit-terminal',
  'refresh-pods',
  'view-pod-logs'
])
</script>
