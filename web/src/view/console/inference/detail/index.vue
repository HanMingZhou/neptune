<template>
  <div class="console-detail-page min-h-screen">
    <DetailHeader
      :action-loading="actionLoading"
      :get-status-class="getStatusClass"
      :get-status-label="getStatusLabel"
      :is-running-or-pending="isRunningOrPending"
      :is-stopped="isStopped"
      :is-terminal-state="isTerminalState"
      :service="service"
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

    <InferenceDetailPanels
      :active-tab="activeTab"
      :api-endpoint="apiEndpoint"
      :can-connect-terminal="canConnectTerminal"
      :curl-example="curlExample"
      :format-command="formatCommand"
      :format-time="formatTime"
      :get-pod-status-class="getPodStatusClass"
      :get-status-label="getStatusLabel"
      :logs="logs"
      :logs-connected="logsConnected"
      :logs-loading="logsLoading"
      :pods="pods"
      :pods-loading="podsLoading"
      :selected-pod="selectedPod"
      :service="service"
      :set-logs-ref="setLogsRef"
      :set-terminal-ref="setTerminalRef"
      :terminal-connected="terminalConnected"
      :terminal-pod="terminalPod"
      @clear-logs="clearLogs"
      @connect-logs="connectLogStream"
      @connect-terminal="connectTerminal"
      @copy="copyText"
      @disconnect-logs="disconnectLogStream"
      @disconnect-terminal="disconnectTerminal"
      @fit-terminal="fitTerminal"
      @manage-api-key="showApiKeyDialog = true"
      @refresh-pods="fetchPods"
      @update:selected-pod="selectedPod = $event"
      @update:terminal-pod="terminalPod = $event"
      @view-pod-logs="viewPodLogs"
    />

    <InferenceDetailDialogsHost
      :api-keys="apiKeys"
      :format-time="formatTime"
      :new-key-name="newKeyName"
      :newly-created-key="newlyCreatedKey"
      :show-api-key-dialog="showApiKeyDialog"
      @close-api-key-dialog="showApiKeyDialog = false"
      @copy="copyText"
      @create-api-key="handleCreateApiKey"
      @delete-api-key="handleDeleteApiKey"
      @update:new-key-name="newKeyName = $event"
    />
  </div>
</template>

<script setup>
import DetailHeader from './components/DetailHeader.vue'
import InferenceDetailDialogsHost from './components/InferenceDetailDialogsHost.vue'
import InferenceDetailPanels from './components/InferenceDetailPanels.vue'
import DetailTabs from '@/components/detailPage/DetailTabs.vue'
import { useInferenceDetail } from './composables/useInferenceDetail'

const {
  actionLoading,
  activeTab,
  apiEndpoint,
  apiKeys,
  canConnectTerminal,
  clearLogs,
  connectLogStream,
  connectTerminal,
  copyText,
  curlExample,
  disconnectLogStream,
  disconnectTerminal,
  fetchPods,
  fitTerminal,
  formatCommand,
  formatTime,
  getPodStatusClass,
  getStatusClass,
  getStatusLabel,
  goBack,
  handleCreateApiKey,
  handleDelete,
  handleDeleteApiKey,
  handleStart,
  handleStop,
  isRunningOrPending,
  isStopped,
  isTerminalState,
  logs,
  logsConnected,
  logsLoading,
  newKeyName,
  newlyCreatedKey,
  pods,
  podsLoading,
  selectedPod,
  service,
  setLogsRef,
  setTerminalRef,
  showApiKeyDialog,
  tabs,
  terminalConnected,
  terminalPod,
  viewPodLogs
} = useInferenceDetail()
</script>
