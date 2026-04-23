import {
  computed,
  inject,
  nextTick,
  onMounted,
  onUnmounted,
  ref,
  watch
} from 'vue'
import { AxiosError } from 'axios'
import { useRoute, useRouter, type LocationQueryValue } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  deleteInferenceService,
  getInferenceServiceDetail,
  getInferenceServicePods,
  startInferenceService,
  stopInferenceService
} from '@/api/inference'
import { useUserStore } from '@/pinia/modules/user'
import type { ApiResponse } from '@/utils/request'
import type {
  ConsoleInferenceDetail,
  ConsolePod,
  Translator
} from '@/types/consoleResource'
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

const getQueryValue = (
  value: LocationQueryValue | LocationQueryValue[] | null | undefined
): string => (Array.isArray(value) ? value[0] || '' : value || '')

const extractErrorMessage = (error: unknown, fallback: string): string => {
  if (error instanceof AxiosError) {
    return error.response?.data?.msg || error.message || fallback
  }
  if (error instanceof Error) {
    return error.message || fallback
  }
  return fallback
}

export const useInferenceDetail = () => {
  const t = inject<Translator>('t', (key: string) => key)
  const route = useRoute()
  const router = useRouter()
  const userStore = useUserStore()

  const tabs = computed(() => createInferenceDetailTabs(t))
  const resourceId = computed(() => Number(getQueryValue(route.query.id)) || 0)

  const activeTab = ref(getQueryValue(route.query.tab) || 'overview')
  const service = ref<Partial<ConsoleInferenceDetail>>({})
  const pods = ref<ConsolePod[]>([])
  const podsLoading = ref(false)
  const actionLoading = ref(false)
  const logsRef = ref<HTMLElement | null>(null)
  const terminalRef = ref<HTMLElement | null>(null)
  const selectedPod = ref('')
  const terminalPod = ref('')
  let refreshTimer: ReturnType<typeof setInterval> | number | null = null

  const isRunning = computed(() => service.value.status === 'RUNNING')
  const isRunningOrPending = computed(() =>
    ['RUNNING', 'PENDING'].includes(service.value.status || '')
  )
  const isStopped = computed(() => service.value.status === 'STOPPED')
  const isTerminalState = computed(() =>
    ['STOPPED', 'FAILED'].includes(service.value.status || '')
  )
  const canConnectTerminal = computed(() =>
    ['RUNNING', 'PENDING', 'CREATING'].includes(service.value.status || '')
  )
  const apiEndpoint = computed(() => buildInferenceApiEndpoint(service.value))
  const curlExample = computed(() =>
    buildInferenceCurlExample(service.value, apiEndpoint.value)
  )

  const setLogsRef = (element: HTMLElement | null): void => {
    logsRef.value = element
  }

  const setTerminalRef = (element: HTMLElement | null): void => {
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
    getResourceId: () => resourceId.value,
    getToken: () => userStore.token,
    logsRef,
    selectedPod,
    t
  })

  const {
    connectTerminal,
    disconnectTerminal,
    fitTerminal,
    terminalConnected
  } = useInferenceTerminal({
    canConnectTerminal,
    getResourceId: () => resourceId.value,
    getToken: () => userStore.token,
    t,
    terminalPod,
    terminalRef
  })

  const {
    apiKeys,
    apiKeysLoading,
    creatingApiKey,
    deletingApiKeyId,
    handleCreateApiKey,
    handleDeleteApiKey,
    newKeyName,
    newlyCreatedKey,
    showApiKeyDialog
  } = useInferenceApiKeys({
    getServiceId: () => resourceId.value,
    t
  })

  const fetchDetail = async (): Promise<void> => {
    if (!resourceId.value) {
      return
    }

    try {
      const res = (await getInferenceServiceDetail({
        id: resourceId.value
      })) as ApiResponse<ConsoleInferenceDetail>
      if (res.code === 0) {
        service.value = res.data || {}
      }
    } catch (_error) {
      console.error('获取详情失败', _error)
    }
  }

  const fetchPods = async (): Promise<void> => {
    if (!resourceId.value) {
      return
    }

    podsLoading.value = true
    try {
      const res = (await getInferenceServicePods({
        id: resourceId.value
      })) as ApiResponse<ConsolePod[]>
      if (res.code === 0) {
        pods.value = res.data || []
        if (pods.value.length && !selectedPod.value) {
          selectedPod.value = pods.value[0].name
        }
        if (pods.value.length && !terminalPod.value) {
          terminalPod.value = pods.value[0].name
        }
      }
    } catch (_error) {
      console.error('获取实例列表失败', _error)
    } finally {
      podsLoading.value = false
    }
  }

  const handleStart = async (): Promise<void> => {
    if (actionLoading.value) return
    try {
      await ElMessageBox.confirm(
        t('confirmStart', {
          name: service.value.displayName || service.value.instanceName || '-'
        }),
        t('tip'),
        { type: 'info' }
      )
      actionLoading.value = true
      const res = await startInferenceService({ id: resourceId.value })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        void fetchDetail()
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (error) {
      if (error === 'cancel' || error === 'close') {
        return
      }
      ElMessage.error(extractErrorMessage(error, t('error')))
    } finally {
      actionLoading.value = false
    }
  }

  const handleStop = async (): Promise<void> => {
    if (actionLoading.value) return
    try {
      await ElMessageBox.confirm(
        t('confirmStop', {
          name: service.value.displayName || service.value.instanceName || '-'
        }),
        t('tip'),
        { type: 'warning' }
      )
      actionLoading.value = true
      const res = await stopInferenceService({ id: resourceId.value })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        void fetchDetail()
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (_error) {
    } finally {
      actionLoading.value = false
    }
  }

  const handleDelete = async (): Promise<void> => {
    if (actionLoading.value) return
    try {
      await ElMessageBox.confirm(
        t('confirmDelete', {
          name: service.value.displayName || service.value.instanceName || '-'
        }),
        t('tip'),
        { type: 'warning' }
      )
      actionLoading.value = true
      const res = await deleteInferenceService({ id: resourceId.value })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        void router.push({ name: 'inference' })
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (_error) {
    } finally {
      actionLoading.value = false
    }
  }

  const goBack = (): void => router.go(-1)

  const viewPodLogs = (pod: ConsolePod): void => {
    selectedPod.value = pod.name || ''
    activeTab.value = 'logs'
    void nextTick(() => {
      connectLogStream()
    })
  }

  const fallbackCopyText = (text: string): boolean => {
    if (typeof document === 'undefined') {
      return false
    }

    const textarea = document.createElement('textarea')
    textarea.value = text
    textarea.setAttribute('readonly', 'true')
    textarea.style.position = 'fixed'
    textarea.style.opacity = '0'
    textarea.style.pointerEvents = 'none'
    textarea.style.left = '-9999px'
    document.body.appendChild(textarea)
    textarea.focus()
    textarea.select()

    let copied = false
    try {
      copied = document.execCommand('copy')
    } catch (_error) {
      copied = false
    }

    document.body.removeChild(textarea)
    return copied
  }

  const copyText = async (text: string): Promise<void> => {
    if (!text) {
      ElMessage.error(t('error'))
      return
    }

    try {
      if (navigator.clipboard?.writeText && window.isSecureContext) {
        await navigator.clipboard.writeText(text)
      } else if (!fallbackCopyText(text)) {
        throw new Error('copy failed')
      }

      ElMessage.success(t('success'))
    } catch (_error) {
      ElMessage.error(t('error'))
    }
  }
  const formatCommand = formatInferenceCommand
  const formatTime = formatInferenceTime
  const getStatusLabel = (status?: string) => getInferenceStatusLabel(t, status)
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
      void fetchPods()
    }
  })

  watch(selectedPod, () => {
    if (logsConnected.value) {
      disconnectLogStream()
      connectLogStream()
    }
  })

  onMounted(() => {
    void fetchDetail()
    void fetchPods()
    refreshTimer = window.setInterval(() => {
      void fetchDetail()
    }, 5000)
  })

  onUnmounted(() => {
    if (refreshTimer !== null) {
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
    apiKeysLoading,
    canConnectTerminal,
    clearLogs,
    connectLogStream,
    connectTerminal,
    copyText,
    creatingApiKey,
    curlExample,
    deletingApiKeyId,
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
