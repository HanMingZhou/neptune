import { computed, inject, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  deleteInferenceService,
  getInferenceServiceDetail,
  getInferenceServicePods,
  startInferenceService,
  stopInferenceService
} from '@/api/inference'
import { useUserStore } from '@/pinia/modules/user'
import {
  buildInferenceApiEndpoint,
  buildInferenceCurlExample,
  createInferenceDetailTabs,
  formatInferenceCommand,
  formatInferenceTime,
  getInferencePodStatusClass,
  getInferenceStatusClass,
  getInferenceStatusLabel
} from './inferenceDetailUtils'
import { useInferenceApiKeys } from './useInferenceApiKeys'
import { useInferenceLogStream } from './useInferenceLogStream'
import { useInferenceTerminal } from './useInferenceTerminal'

export const useInferenceDetail = () => {
  const t = inject('t', (key) => key)
  const route = useRoute()
  const router = useRouter()
  const userStore = useUserStore()

  const tabs = computed(() => createInferenceDetailTabs(t))

  const activeTab = ref(route.query.tab || 'overview')
  const service = ref({})
  const pods = ref([])
  const podsLoading = ref(false)
  const actionLoading = ref(false)
  const logsRef = ref(null)
  const terminalRef = ref(null)
  const selectedPod = ref('')
  const terminalPod = ref('')
  let refreshTimer = null

  const isRunning = computed(() => service.value.status === 'RUNNING')
  const isRunningOrPending = computed(() => ['RUNNING', 'PENDING'].includes(service.value.status))
  const isStopped = computed(() => service.value.status === 'STOPPED')
  const isTerminalState = computed(() => ['STOPPED', 'FAILED'].includes(service.value.status))
  const canConnectTerminal = computed(() => ['RUNNING', 'PENDING', 'CREATING'].includes(service.value.status))
  const apiEndpoint = computed(() => buildInferenceApiEndpoint(service.value))
  const curlExample = computed(() => buildInferenceCurlExample(service.value, apiEndpoint.value))

  const setLogsRef = (element) => {
    logsRef.value = element
  }

  const setTerminalRef = (element) => {
    terminalRef.value = element
  }

  const {
    clearLogs,
    connectLogStream,
    disconnectLogStream,
    logs,
    logsConnected,
    logsLoading
  } = useInferenceLogStream({
    logsRef,
    route,
    selectedPod,
    t,
    userStore
  })

  const {
    connectTerminal,
    disconnectTerminal,
    fitTerminal,
    terminalConnected
  } = useInferenceTerminal({
    canConnectTerminal,
    route,
    t,
    terminalPod,
    terminalRef,
    userStore
  })

  const {
    apiKeys,
    handleCreateApiKey,
    handleDeleteApiKey,
    newKeyName,
    newlyCreatedKey,
    showApiKeyDialog
  } = useInferenceApiKeys({
    route,
    t
  })

  const fetchDetail = async () => {
    try {
      const res = await getInferenceServiceDetail({ id: route.query.id })
      if (res.code === 0) {
        service.value = res.data || {}
      }
    } catch (error) {
      console.error('获取详情失败', error)
    }
  }

  const fetchPods = async () => {
    podsLoading.value = true
    try {
      const res = await getInferenceServicePods({ id: route.query.id })
      if (res.code === 0) {
        pods.value = res.data || []
        if (pods.value.length && !selectedPod.value) {
          selectedPod.value = pods.value[0].name
        }
        if (pods.value.length && !terminalPod.value) {
          terminalPod.value = pods.value[0].name
        }
      }
    } catch (error) {
      console.error('获取实例列表失败', error)
    } finally {
      podsLoading.value = false
    }
  }

  const handleStart = async () => {
    if (actionLoading.value) return
    try {
      await ElMessageBox.confirm(t('confirmStart', { name: service.value.displayName }), t('tip'), { type: 'info' })
      actionLoading.value = true
      const res = await startInferenceService({ id: Number(route.query.id) })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        fetchDetail()
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (error) {
    } finally {
      actionLoading.value = false
    }
  }

  const handleStop = async () => {
    if (actionLoading.value) return
    try {
      await ElMessageBox.confirm(t('confirmStop', { name: service.value.displayName }), t('tip'), { type: 'warning' })
      actionLoading.value = true
      const res = await stopInferenceService({ id: Number(route.query.id) })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        fetchDetail()
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (error) {
    } finally {
      actionLoading.value = false
    }
  }

  const handleDelete = async () => {
    if (actionLoading.value) return
    try {
      await ElMessageBox.confirm(t('confirmDelete', { name: service.value.displayName }), t('tip'), { type: 'warning' })
      actionLoading.value = true
      const res = await deleteInferenceService({ id: Number(route.query.id) })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        router.push({ name: 'inference' })
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (error) {
    } finally {
      actionLoading.value = false
    }
  }

  const goBack = () => router.go(-1)

  const viewPodLogs = (pod) => {
    selectedPod.value = pod.name
    activeTab.value = 'logs'
    nextTick(() => connectLogStream())
  }

  const copyText = (text) => {
    navigator.clipboard.writeText(text).then(() => {
      ElMessage.success(t('success'))
    }).catch(() => {
      ElMessage.error(t('error'))
    })
  }
  const formatCommand = formatInferenceCommand
  const formatTime = formatInferenceTime
  const getStatusLabel = (status) => getInferenceStatusLabel(t, status)
  const getStatusClass = getInferenceStatusClass
  const getPodStatusClass = getInferencePodStatusClass

  watch(activeTab, (tab) => {
    if (tab === 'logs' && isRunning.value) {
      connectLogStream()
    } else if (tab !== 'logs') {
      disconnectLogStream()
    }

    if (tab === 'terminal') {
      nextTick(() => {
        fitTerminal()
      })
    }

    if (tab === 'pods') {
      fetchPods()
    }
  })

  watch(selectedPod, () => {
    if (logsConnected.value) {
      disconnectLogStream()
      connectLogStream()
    }
  })

  onMounted(() => {
    fetchDetail()
    fetchPods()
    refreshTimer = setInterval(fetchDetail, 5000)
  })

  onUnmounted(() => {
    if (refreshTimer) {
      clearInterval(refreshTimer)
    }
    disconnectLogStream()
    disconnectTerminal()
  })

  return {
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
  }
}
