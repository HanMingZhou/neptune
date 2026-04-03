<template>
  <div class="console-detail-page min-h-screen">
    <DetailHeader
      :action-loading="actionLoading"
      :get-status-class="getStatusClass"
      :get-status-label="getStatusLabel"
      :is-running="isRunning"
      :is-terminal="isTerminal"
      :job="job"
      @back="goBack"
      @delete="handleDelete"
      @open-tensorboard="openTensorboard"
      @stop="handleStop"
    />

    <DetailTabs
      :active-tab="activeTab"
      :tabs="tabs"
      @update:active-tab="activeTab = $event"
    />

    <TrainingDetailPanels
      :active-tab="activeTab"
      :download-loading="downloadLoading"
      :format-time="formatTime"
      :get-framework-label="getFrameworkLabel"
      :get-pay-type-label="getPayTypeLabel"
      :get-pod-status-class="getPodStatusClass"
      :get-status-label="getStatusLabel"
      :is-running="isRunning"
      :is-terminal="isTerminal"
      :job="job"
      :logs="logs"
      :logs-connected="logsConnected"
      :logs-loading="logsLoading"
      :log-task-name="logOptions.taskName"
      :master-task-label="masterTaskLabel"
      :master-task-name="masterTaskName"
      :pod-count="podCount"
      :pod-index="logOptions.podIndex"
      :pods="pods"
      :pods-loading="podsLoading"
      :set-logs-ref="setLogsRef"
      :set-terminal-ref="setTerminalRef"
      :terminal-connected="terminalConnected"
      :terminal-task-name="terminalOptions.taskName"
      @clear-logs="clearLogs"
      @connect-logs="connectLogStream"
      @connect-terminal="connectTerminal"
      @disconnect-logs="disconnectLogStream"
      @disconnect-terminal="disconnectTerminal"
      @download-logs="handleDownloadLogs"
      @fit-terminal="fitTerminal"
      @refresh-pods="fetchPods"
      @update:log-task-name="logOptions.taskName = $event"
      @update:pod-index="logOptions.podIndex = $event"
      @update:terminal-task-name="terminalOptions.taskName = $event"
      @view-pod-logs="viewPodLogs"
    />
  </div>
</template>

<script setup>
import DetailHeader from './components/DetailHeader.vue'
import DetailTabs from '@/components/detailPage/DetailTabs.vue'
import TrainingDetailPanels from './components/TrainingDetailPanels.vue'
import { useTrainingDetail } from './composables/useTrainingDetail'

const {
  actionLoading,
  activeTab,
  clearLogs,
  connectLogStream,
  connectTerminal,
  disconnectLogStream,
  disconnectTerminal,
  downloadLoading,
  fetchPods,
  fitTerminal,
  formatTime,
  getFrameworkLabel,
  getPayTypeLabel,
  getPodStatusClass,
  getStatusClass,
  getStatusLabel,
  goBack,
  handleDelete,
  handleDownloadLogs,
  handleStop,
  isRunning,
  isTerminal,
  job,
  logOptions,
  logs,
  logsConnected,
  logsLoading,
  masterTaskLabel,
  masterTaskName,
  openTensorboard,
  podCount,
  pods,
  podsLoading,
  setLogsRef,
  setTerminalRef,
  tabs,
  terminalConnected,
  terminalOptions,
  viewPodLogs
} = useTrainingDetail()
</script>
