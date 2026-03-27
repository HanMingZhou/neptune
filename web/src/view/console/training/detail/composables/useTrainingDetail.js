import { computed, inject, nextTick, onMounted, onUnmounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  deleteTrainingJob,
  downloadTrainingJobLogs,
  getTrainingJobDetail,
  getTrainingJobPods,
  stopTrainingJob
} from '@/api/training'
import { useUserStore } from '@/pinia/modules/user'
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

export const useTrainingDetail = () => {
  const t = inject('t', (key) => key)
  const route = useRoute()
  const router = useRouter()
  const userStore = useUserStore()

  const tabs = computed(() => createTrainingDetailTabs(t))

  const activeTab = ref(route.query.tab || 'overview')
  const job = ref({})
  const pods = ref([])
  const podsLoading = ref(false)
  const downloadLoading = ref(false)
  const actionLoading = ref(false)
  const logsRef = ref(null)
  const terminalRef = ref(null)
  const logOptions = reactive({
    taskName: '',
    podIndex: 0
  })
  const terminalOptions = reactive({
    taskName: ''
  })
  let refreshTimer = null

  const masterTaskName = computed(() => getTrainingMasterTaskName(job.value))
  const masterTaskLabel = computed(() => getTrainingMasterTaskLabel(job.value))
  const isRunning = computed(() => ['RUNNING', 'PENDING', 'CREATING'].includes(job.value.status))
  const isTerminal = computed(() => ['SUCCEEDED', 'FAILED', 'KILLED'].includes(job.value.status))
  const podCount = computed(() => getTrainingPodCount(job.value, logOptions.taskName))
  const statusText = computed(() => job.value.status || '')

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
  } = useTrainingLogStream({
    logsRef,
    logOptions,
    route,
    t,
    userStore
  })

  const {
    connectTerminal,
    disconnectTerminal,
    fitTerminal,
    terminalConnected
  } = useTrainingTerminal({
    canConnectTerminal: isRunning,
    route,
    statusText,
    t,
    terminalOptions,
    terminalRef,
    userStore
  })

  const fetchDetail = async () => {
    try {
      const res = await getTrainingJobDetail({ id: route.query.id })
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
    } catch (error) {
      console.error('获取详情失败', error)
    }
  }

  const handleDownloadLogs = async () => {
    downloadLoading.value = true

    try {
      const res = await downloadTrainingJobLogs({
        id: route.query.id,
        taskName: logOptions.taskName,
        podIndex: logOptions.podIndex
      })
      const blob = new Blob([res], { type: 'text/plain' })
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = `${job.value.displayName || job.value.jobName || 'training'}-${logOptions.taskName}-${logOptions.podIndex}.log`
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
      ElMessage.success(t('success'))
    } catch (error) {
      console.error('下载日志失败', error)
      ElMessage.error(t('error'))
    } finally {
      downloadLoading.value = false
    }
  }

  const fetchPods = async () => {
    podsLoading.value = true

    try {
      const res = await getTrainingJobPods({ id: route.query.id })
      if (res.code === 0) {
        pods.value = res.data || []
      }
    } catch (error) {
      console.error('获取实例列表失败', error)
    } finally {
      podsLoading.value = false
    }
  }

  const viewPodLogs = (pod) => {
    logOptions.taskName = resolveTrainingLogTaskName(job.value, pod)
    activeTab.value = 'logs'
    connectLogStream()
  }

  const handleStop = async () => {
    if (actionLoading.value) {
      return
    }

    try {
      await ElMessageBox.confirm(t('confirm'), t('tip'), { type: 'warning' })
      actionLoading.value = true
      const res = await stopTrainingJob({ id: Number(route.query.id) })

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
    if (actionLoading.value) {
      return
    }

    try {
      await ElMessageBox.confirm(
        t('confirmDelete', { name: job.value.displayName || job.value.jobName }),
        t('tip'),
        { type: 'warning' }
      )
      actionLoading.value = true
      const res = await deleteTrainingJob({ id: Number(route.query.id) })

      if (res.code === 0) {
        ElMessage.success(t('success'))
        router.push({ name: 'training' })
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (error) {
    } finally {
      actionLoading.value = false
    }
  }

  const goBack = () => router.go(-1)

  const openTensorboard = () => {
    if (job.value.tensorboardUrl) {
      window.open(job.value.tensorboardUrl, '_blank')
    }
  }

  const formatTime = formatTrainingTime
  const getFrameworkLabel = getTrainingFrameworkLabel
  const getPayTypeLabel = (type) => getTrainingPayTypeLabel(t, type)
  const getStatusLabel = (status) => getTrainingStatusLabel(t, status)
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
      fetchPods()
    }
  })

  watch(() => [logOptions.taskName, logOptions.podIndex], () => {
    if (logsConnected.value) {
      disconnectLogStream()
      connectLogStream()
    }
  })

  watch(() => logOptions.taskName, () => {
    logOptions.podIndex = 0
  })

  onMounted(() => {
    fetchDetail()
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
