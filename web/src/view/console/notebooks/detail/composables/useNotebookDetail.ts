import {
  computed,
  inject,
  nextTick,
  onMounted,
  onUnmounted,
  ref,
  watch
} from 'vue'
import { useRoute, useRouter, type LocationQueryValue } from 'vue-router'
import type { Terminal as XTermTerminal } from '@xterm/xterm'
import type { FitAddon } from '@xterm/addon-fit'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  deleteNotebook,
  getNotebookDetail,
  getNotebookPods,
  startNotebook,
  stopNotebook
} from '@/api/notebook'
import { useUserStore } from '@/pinia/modules/user'
import type { ApiResponse } from '@/utils/request'
import type {
  ConsoleNotebookDetail,
  ConsoleNotebookVolumeMount,
  ConsolePod,
  DetailTab,
  Translator
} from '@/types/consoleResource'

interface DisposableLike {
  dispose: () => void
}

type NotebookLike = Partial<ConsoleNotebookDetail>

const getQueryValue = (
  value: LocationQueryValue | LocationQueryValue[] | null | undefined
): string => (Array.isArray(value) ? value[0] || '' : value || '')

const normalizeStatus = (status: unknown): string =>
  `${status || ''}`.trim().toUpperCase()

const normalizeVolumeMounts = (
  volumeMounts: ConsoleNotebookDetail['volumeMounts']
): ConsoleNotebookVolumeMount[] => {
  if (Array.isArray(volumeMounts)) return volumeMounts

  if (typeof volumeMounts === 'string' && volumeMounts.trim()) {
    try {
      const parsed = JSON.parse(volumeMounts) as unknown
      return Array.isArray(parsed)
        ? (parsed as ConsoleNotebookVolumeMount[])
        : []
    } catch (_error) {
      console.warn('解析挂载卷数据失败', _error)
    }
  }

  return []
}

const normalizeNotebookDetail = (
  data: ConsoleNotebookDetail = {}
): NotebookLike => ({
  ...data,
  price: Number.isFinite(Number(data.price)) ? Number(data.price) : 0,
  volumeMounts: normalizeVolumeMounts(data.volumeMounts)
})

const processLogData = (data: string): string => {
  let processed = data.replace(/\r(?!\n)/g, '\n')
  processed = processed.replace(/\n{3,}/g, '\n\n')
  return processed
}

export const useNotebookDetail = () => {
  const t = inject<Translator>('t', (key: string) => key)
  const route = useRoute()
  const router = useRouter()
  const userStore = useUserStore()

  const resourceId = computed(() => Number(getQueryValue(route.query.id)) || 0)
  const tabs = computed<DetailTab[]>(() => [
    { key: 'overview', label: t('overview'), icon: 'description' },
    { key: 'logs', label: t('logs'), icon: 'article' },
    { key: 'terminal', label: t('terminal'), icon: 'terminal' },
    { key: 'pods', label: t('instanceList'), icon: 'dns' }
  ])

  const activeTab = ref(getQueryValue(route.query.tab) || 'overview')
  const notebook = ref<NotebookLike>({})
  const logs = ref('')
  const pods = ref<ConsolePod[]>([])
  const logsLoading = ref(false)
  const logsConnected = ref(false)
  const podsLoading = ref(false)
  const terminalConnected = ref(false)
  const actionLoading = ref(false)
  const logsRef = ref<HTMLElement | null>(null)
  const terminalRef = ref<HTMLElement | null>(null)

  let refreshTimer: number | null = null
  let logsWs: WebSocket | null = null
  let terminal: XTermTerminal | null = null
  let fitAddon: FitAddon | null = null
  let resizeObserver: ResizeObserver | null = null
  let ws: WebSocket | null = null
  let dataSubscription: DisposableLike | null = null

  const setLogsRef = (element: HTMLElement | null): void => {
    logsRef.value = element
  }

  const setTerminalRef = (element: HTMLElement | null): void => {
    terminalRef.value = element
  }

  const isInstanceRunning = computed(() => {
    const status = normalizeStatus(notebook.value.status)
    return status === 'RUNNING' || status === 'READY'
  })

  const fetchDetail = async (): Promise<void> => {
    if (!resourceId.value) {
      return
    }

    try {
      const res = (await getNotebookDetail({
        id: resourceId.value
      })) as ApiResponse<ConsoleNotebookDetail>
      if (res.code === 0) {
        notebook.value = normalizeNotebookDetail(res.data || {})
      }
    } catch (_error) {
      console.error('获取详情失败', _error)
    }
  }

  const connectLogStream = (): void => {
    if (logsConnected.value || logsLoading.value || logsWs) return

    const token = userStore.token
    if (!resourceId.value || !token) {
      return
    }

    logsLoading.value = true
    logs.value = ''

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const baseApi = import.meta.env.VITE_BASE_API || ''
    const wsUrl = `${protocol}//${window.location.host}${baseApi}/api/v1/notebook/log/stream?id=${resourceId.value}&token=${token}`

    logsWs = new WebSocket(wsUrl)

    logsWs.onopen = () => {
      logsConnected.value = true
      logsLoading.value = false
    }

    logsWs.onmessage = async (event: MessageEvent<string>) => {
      logs.value += processLogData(event.data)
      await nextTick()

      if (logsRef.value) {
        logsRef.value.scrollTop = logsRef.value.scrollHeight
      }
    }

    logsWs.onclose = () => {
      logsWs = null
      logsConnected.value = false
      logsLoading.value = false
    }

    logsWs.onerror = (error) => {
      console.error('WebSocket 错误', error)
      logsWs = null
      logsConnected.value = false
      logsLoading.value = false
      ElMessage.error(t('error'))
    }
  }

  const disconnectLogStream = (): void => {
    if (logsWs) {
      logsWs.close()
      logsWs = null
    }
    logsConnected.value = false
    logsLoading.value = false
  }

  const clearLogs = (): void => {
    logs.value = ''
  }

  const fetchPods = async (): Promise<void> => {
    if (!resourceId.value) {
      return
    }

    podsLoading.value = true
    try {
      const res = (await getNotebookPods({
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

  const fitTerminal = (): void => {
    fitAddon?.fit()
  }

  const disconnectTerminal = (): void => {
    if (dataSubscription) {
      dataSubscription.dispose()
      dataSubscription = null
    }

    if (ws) {
      ws.onclose = null
      ws.onerror = null
      ws.onmessage = null
      ws.onopen = null
      ws.close()
      ws = null
    }

    if (terminal) {
      terminal.dispose()
      terminal = null
    }

    if (resizeObserver) {
      resizeObserver.disconnect()
      resizeObserver = null
    }

    terminalConnected.value = false
  }

  const connectTerminal = async (): Promise<void> => {
    if (!isInstanceRunning.value) {
      ElMessage.warning(
        `${t('instanceStatus')}: ${notebook.value.status || '-'}`
      )
      return
    }

    const token = userStore.token
    const terminalContainer = terminalRef.value
    if (!resourceId.value || !token || !terminalContainer) {
      return
    }

    const { Terminal } = await import('@xterm/xterm')
    await import('@xterm/xterm/css/xterm.css')
    const { FitAddon } = await import('@xterm/addon-fit')

    disconnectTerminal()

    terminal = new Terminal({
      cursorBlink: true,
      fontSize: 14,
      lineHeight: 1.2,
      fontFamily: 'JetBrains Mono, Menlo, Monaco, Consolas, monospace',
      theme: {
        background: '#1e1e1e',
        foreground: '#d4d4d4',
        cursor: '#ffffff',
        selectionBackground: '#264f78',
        black: '#000000',
        red: '#cd3131',
        green: '#0dbc79',
        yellow: '#e5e510',
        blue: '#2472c8',
        magenta: '#bc3fbc',
        cyan: '#11a8cd',
        white: '#e5e5e5',
        brightBlack: '#666666',
        brightRed: '#f14c4c',
        brightGreen: '#23d18b',
        brightYellow: '#f5f543',
        brightBlue: '#3b8eea',
        brightMagenta: '#d670d6',
        brightCyan: '#29b8db',
        brightWhite: '#e5e5e5'
      }
    })

    fitAddon = new FitAddon()
    terminal.loadAddon(fitAddon)
    terminal.open(terminalContainer)

    resizeObserver = new ResizeObserver(() => {
      fitAddon?.fit()
    })
    resizeObserver.observe(terminalContainer)

    window.setTimeout(() => {
      fitAddon?.fit()
    }, 50)

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const baseApi = import.meta.env.VITE_BASE_API || ''
    const wsUrl = `${protocol}//${window.location.host}${baseApi}/api/v1/notebook/terminal?id=${resourceId.value}&token=${token}`

    ws = new WebSocket(wsUrl)

    ws.onopen = () => {
      terminalConnected.value = true
      fitAddon?.fit()
      terminal?.focus()
    }

    ws.onmessage = (event: MessageEvent<string>) => {
      terminal?.write(event.data)
    }

    ws.onclose = () => {
      terminalConnected.value = false
      terminal?.write('\r\n\x1b[31mConnection disconnected.\x1b[0m\r\n')
    }

    ws.onerror = () => {
      ElMessage.error(t('error'))
    }

    dataSubscription = terminal.onData((data) => {
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(data)
      }
    })
  }

  const viewPodLogs = (): void => {
    activeTab.value = 'logs'
    connectLogStream()
  }

  const handleStart = async (): Promise<void> => {
    if (actionLoading.value) return

    try {
      await ElMessageBox.confirm(t('confirm'), t('tip'), { type: 'info' })
      actionLoading.value = true
      const res = await startNotebook({ id: resourceId.value })
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

  const handleStop = async (): Promise<void> => {
    if (actionLoading.value) return

    try {
      await ElMessageBox.confirm(t('confirm'), t('tip'), { type: 'warning' })
      actionLoading.value = true
      const res = await stopNotebook({ id: resourceId.value })
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
        t('confirmDelete', { name: notebook.value.displayName || '-' }),
        t('tip'),
        { type: 'warning' }
      )
      actionLoading.value = true
      const res = await deleteNotebook({ id: resourceId.value })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        void router.push({ name: 'notebooks' })
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (_error) {
    } finally {
      actionLoading.value = false
    }
  }

  const goBack = (): void => router.go(-1)

  const getPayTypeLabel = (type?: number | string): string => {
    const map: Record<number, string> = {
      1: 'payHourly',
      2: 'payDaily',
      3: 'payWeekly',
      4: 'payMonthly'
    }
    return t(map[Number(type)] || 'payHourly')
  }

  const getUnitPriceLabel = (type?: number | string): string => {
    const map: Record<number, string> = {
      1: 'unitHour',
      2: 'unitDay',
      3: 'unitWeek',
      4: 'unitMonth'
    }
    return t(map[Number(type)] || 'unitHour')
  }

  const getStatusLabel = (status?: string): string =>
    t(normalizeStatus(status)) || status || '-'

  const getStatusClass = (status?: string): string => {
    const normalized = normalizeStatus(status)
    if (normalized === 'RUNNING' || normalized === 'READY')
      return 'bg-emerald-500/10 text-emerald-500'
    if (normalized === 'PENDING' || normalized === 'CREATING')
      return 'bg-amber-500/10 text-amber-500'
    if (normalized === 'FAILED' || normalized === 'ERROR')
      return 'bg-red-500/10 text-red-500'
    if (normalized === 'STOPPED') return 'bg-slate-500/10 text-slate-500'
    return 'bg-slate-500/10 text-slate-500'
  }

  const getPodStatusClass = (status?: string): string => {
    const normalized = `${status || ''}`.trim().toLowerCase()
    const map: Record<string, string> = {
      running: 'bg-emerald-500/10 text-emerald-500',
      pending: 'bg-amber-500/10 text-amber-500',
      failed: 'bg-red-500/10 text-red-500'
    }
    return map[normalized] || 'bg-slate-500/10 text-slate-500'
  }

  const formatTime = (time?: string | number | null): string => {
    if (!time) return '-'
    return new Date(time).toLocaleString()
  }

  watch(activeTab, (tab) => {
    if (tab === 'logs' && isInstanceRunning.value) {
      connectLogStream()
    } else if (tab !== 'logs') {
      disconnectLogStream()
    }

    if (tab === 'terminal') {
      void nextTick(() => {
        fitTerminal()
      })
    }

    if (tab === 'pods') {
      void fetchPods()
    }
  })

  onMounted(() => {
    void fetchDetail()
    refreshTimer = window.setInterval(fetchDetail, 5000)
  })

  onUnmounted(() => {
    if (refreshTimer !== null) clearInterval(refreshTimer)
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
  }
}
