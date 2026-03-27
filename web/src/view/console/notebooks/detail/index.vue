<template>
  <div class="min-h-screen bg-background-light dark:bg-background-dark">
    <DetailHeader
      :action-loading="actionLoading"
      :get-status-class="getStatusClass"
      :get-status-label="getStatusLabel"
      :notebook="notebook"
      @back="goBack"
      @delete="handleDelete"
      @start="handleStart"
      @stop="handleStop"
    />

    <DetailTabs
      :active-tab="activeTab"
      :tabs="tabs"
      @update:active-tab="activeTab = $event"
    />

    <NotebookDetailPanels
      :active-tab="activeTab"
      :format-time="formatTime"
      :get-pay-type-label="getPayTypeLabel"
      :get-pod-status-class="getPodStatusClass"
      :get-status-label="getStatusLabel"
      :get-unit-price-label="getUnitPriceLabel"
      :is-instance-running="isInstanceRunning"
      :logs="logs"
      :logs-connected="logsConnected"
      :logs-loading="logsLoading"
      :notebook="notebook"
      :pods="pods"
      :pods-loading="podsLoading"
      :set-logs-ref="setLogsRef"
      :set-terminal-ref="setTerminalRef"
      :terminal-connected="terminalConnected"
      @clear-logs="clearLogs"
      @connect-logs="connectLogStream"
      @connect-terminal="connectTerminal"
      @disconnect-logs="disconnectLogStream"
      @disconnect-terminal="disconnectTerminal"
      @fit-terminal="fitTerminal"
      @refresh-pods="fetchPods"
      @view-pod-logs="viewPodLogs"
    />
  </div>
</template>

<script setup>
import DetailHeader from './components/DetailHeader.vue'
import DetailTabs from '@/components/detailPage/DetailTabs.vue'
import NotebookDetailPanels from './components/NotebookDetailPanels.vue'
import { useNotebookDetail } from './composables/useNotebookDetail'

const {
  actionLoading,
  activeTab,
  clearLogs,
  connectLogStream,
  connectTerminal,
  disconnectLogStream,
  disconnectTerminal,
  fetchPods,
  fitTerminal,
  formatTime,
  getPayTypeLabel,
  getPodStatusClass,
  getStatusClass,
  getStatusLabel,
  getUnitPriceLabel,
  goBack,
  handleDelete,
  handleStart,
  handleStop,
  isInstanceRunning,
  logs,
  logsConnected,
  logsLoading,
  notebook,
  pods,
  podsLoading,
  setLogsRef,
  setTerminalRef,
  tabs,
  terminalConnected,
  viewPodLogs
} = useNotebookDetail()
</script>
