import {
  computed,
  inject,
  nextTick,
  onMounted,
  onUnmounted,
  reactive,
  ref,
  watch
} from 'vue'
import { useRoute, useRouter, type LocationQueryValue } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  deleteTrainingJob,
  downloadTrainingJobLogs,
  getTrainingJobDetail,
  getTrainingJobPods,
  stopTrainingJob
} from '@/api/training'
import { useUserStore } from '@/pinia/modules/user'
import type { ApiResponse } from '@/utils/request'
import type {
  ConsolePod,
  ConsoleTrainingDetail,
  Translator
} from '@/types/consoleResource'
import {
  createTrainingDetailTabs,
  formatTrainingTime,
  getDefaultTrainingTaskName,
  getTrainingFrameworkLabel,
  getTrainingMasterTaskLabel,
  getTrainingMasterTaskName,
  getTrainingPayTypeLabel,
  getTrainingPodCount,
  getTrainingPodStatusClass,
  getTrainingStatusClass,
  getTrainingStatusLabel,
  resolveTrainingLogTaskName
} from './trainingDetailUtils'
import { useTrainingLogStream } from './useTrainingLogStream'
import { useTrainingTerminal } from './useTrainingTerminal'

interface TrainingLogOptions {
  taskName: string
  podIndex: number
}

interface TrainingTerminalOptions {
  taskName: string
}

const getQueryValue = (
  value: LocationQueryValue | LocationQueryValue[] | null | undefined
): string => (Array.isArray(value) ? value[0] || '' : value || '')

export const useTrainingDetail = () => {
  const t = inject<Translator>('t', (key: string) => key)
  const route = useRoute()
  const router = useRouter()
  const userStore = useUserStore()

  const tabs = computed(() => createTrainingDetailTabs(t))
  const resourceId = computed(() => Number(getQueryValue(route.query.id)) || 0)

  const activeTab = ref(getQueryValue(route.query.tab) || 'overview')
  const job = ref<Partial<ConsoleTrainingDetail>>({})
  const pods = ref<ConsolePod[]>([])
  const podsLoading = ref(false)
  const downloadLoading = ref(false)
  const actionLoading = ref(false)
  const logsRef = ref<HTMLElement | null>(null)
  const terminalRef = ref<HTMLElement | null>(null)
  const logOptions = reactive<TrainingLogOptions>({
    taskName: '',
    podIndex: 0
  })
  const terminalOptions = reactive<TrainingTerminalOptions>({
    taskName: ''
  })
  let refreshTimer: ReturnType<typeof setInterval> | number | null = null

  const masterTaskName = computed(() => getTrainingMasterTaskName(job.value))
  const masterTaskLabel = computed(() => getTrainingMasterTaskLabel(job.value))
  const isRunning = computed(() =>
    ['RUNNING', 'PENDING', 'CREATING'].includes(job.value.status || '')
  )
  const isTerminal = computed(() =>
    ['SUCCEEDED', 'FAILED', 'KILLED'].includes(job.value.status || '')
  )
  const podCount = computed(() =>
    getTrainingPodCount(job.value, logOptions.taskName)
  )
  const statusText = computed(() => job.value.status || '')

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
  } = useTrainingLogStream({
    getResourceId: () => resourceId.value,
    getToken: () => userStore.token,
    logsRef,
    logOptions,
    t
  })

  const {
    connectTerminal,
    disconnectTerminal,
    fitTerminal,
    terminalConnected
  } = useTrainingTerminal({
    canConnectTerminal: isRunning,
    getResourceId: () => resourceId.value,
    getToken: () => userStore.token,
    statusText,
    t,
    terminalOptions,
    terminalRef
  })

  const fetchDetail = async (): Promise<void> => {
    if (!resourceId.value) {
      return
    }

    try {
      const res = (await getTrainingJobDetail({
        id: resourceId.value
      })) as ApiResponse<ConsoleTrainingDetail>
      if (res.code === 0) {
        job.value = res.data || {}

        const defaultTaskName = getDefaultTrainingTaskName(job.value)

        if (!logOptions.taskName) {
          logOptions.taskName = defaultTaskName
        }

        if (!terminalOptions.taskName) {
          terminalOptions.taskName = defaultTaskName
        }
      }
    } catch (_error) {
      console.error('获取详情失败', _error)
    }
  }

  const handleDownloadLogs = async (): Promise<void> => {
    if (!resourceId.value) {
      return
    }

    downloadLoading.value = true

    try {
      const res = (await downloadTrainingJobLogs({
        id: resourceId.value,
        taskName: logOptions.taskName,
        podIndex: logOptions.podIndex
      })) as unknown
      const blob =
        res instanceof Blob
          ? res
          : new Blob([res as BlobPart], { type: 'text/plain' })
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = `${job.value.displayName || job.value.jobName || 'training'}-${logOptions.taskName}-${logOptions.podIndex}.log`
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
      ElMessage.success(t('success'))
    } catch (_error) {
      console.error('下载日志失败', _error)
      ElMessage.error(t('error'))
    } finally {
      downloadLoading.value = false
    }
  }

  const fetchPods = async (): Promise<void> => {
    if (!resourceId.value) {
      return
    }

    podsLoading.value = true

    try {
      const res = (await getTrainingJobPods({
        id: resourceId.value
      })) as ApiResponse<ConsolePod[]>
      if (res.code === 0) {
        pods.value = res.data || []
      }
    } catch (_error) {
      console.error('获取实例列表失败', _error)
    } finally {
      podsLoading.value = false
    }
  }

  const viewPodLogs = (pod: ConsolePod): void => {
    logOptions.taskName = resolveTrainingLogTaskName(job.value, pod)
    activeTab.value = 'logs'
    connectLogStream()
  }

  const handleStop = async (): Promise<void> => {
    if (actionLoading.value) {
      return
    }

    try {
      await ElMessageBox.confirm(t('confirm'), t('tip'), { type: 'warning' })
      actionLoading.value = true
      const res = await stopTrainingJob({ id: resourceId.value })

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
    if (actionLoading.value) {
      return
    }

    try {
      await ElMessageBox.confirm(
        t('confirmDelete', {
          name: job.value.displayName || job.value.jobName || '-'
        }),
        t('tip'),
        { type: 'warning' }
      )
      actionLoading.value = true
      const res = await deleteTrainingJob({ id: resourceId.value })

      if (res.code === 0) {
        ElMessage.success(t('success'))
        void router.push({ name: 'training' })
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (_error) {
    } finally {
      actionLoading.value = false
    }
  }

  const goBack = (): void => router.go(-1)

  const openTensorboard = (): void => {
    if (job.value.tensorboardUrl) {
      window.open(job.value.tensorboardUrl, '_blank')
    }
  }

  const formatTime = formatTrainingTime
  const getFrameworkLabel = getTrainingFrameworkLabel
  const getPayTypeLabel = (type?: number | string) =>
    getTrainingPayTypeLabel(t, type)
  const getStatusLabel = (status?: string) => getTrainingStatusLabel(t, status)
  const getStatusClass = getTrainingStatusClass
  const getPodStatusClass = getTrainingPodStatusClass

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

  watch(
    () => [logOptions.taskName, logOptions.podIndex] as const,
    () => {
      if (logsConnected.value) {
        disconnectLogStream()
        connectLogStream()
      }
    }
  )

  watch(
    () => logOptions.taskName,
    () => {
      logOptions.podIndex = 0
    }
  )

  onMounted(() => {
    void fetchDetail()
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
  }
}
