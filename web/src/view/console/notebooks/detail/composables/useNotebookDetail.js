import { computed, inject, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  deleteNotebook,
  getNotebookDetail,
  getNotebookPods,
  startNotebook,
  stopNotebook
} from '@/api/notebook'
import { useUserStore } from '@/pinia/modules/user'

export const useNotebookDetail = () => {
  const t = inject('t', (key) => key)
  const route = useRoute()
  const router = useRouter()

  const tabs = computed(() => [
    { key: 'overview', label: t('overview'), icon: 'description' },
    { key: 'logs', label: t('logs'), icon: 'article' },
    { key: 'terminal', label: t('terminal'), icon: 'terminal' },
    { key: 'pods', label: t('instanceList'), icon: 'dns' }
  ])

  const activeTab = ref(route.query.tab || 'overview')
  const notebook = ref({})
  const logs = ref('')
  const pods = ref([])
  const logsLoading = ref(false)
  const logsConnected = ref(false)
  const podsLoading = ref(false)
  const terminalConnected = ref(false)
  const actionLoading = ref(false)
  const logsRef = ref(null)
  const terminalRef = ref(null)

  let refreshTimer = null
  let logsWs = null
  let terminal = null
  let fitAddon = null
  let resizeObserver = null
  let ws = null

  const setLogsRef = (element) => {
    logsRef.value = element
  }

  const setTerminalRef = (element) => {
    terminalRef.value = element
  }

  const isInstanceRunning = computed(() => {
    const status = notebook.value.status?.toLowerCase()
    return status === 'running' || status === 'ready'
  })

  const fetchDetail = async () => {
    try {
      const res = await getNotebookDetail({ id: route.query.id })
      if (res.code === 0) {
        notebook.value = res.data || {}
      }
    } catch (error) {
      console.error('获取详情失败', error)
    }
  }

  const processLogData = (data) => {
    let processed = data.replace(/\r(?!\n)/g, '\n')
    processed = processed.replace(/\n{3,}/g, '\n\n')
    return processed
  }

  const connectLogStream = () => {
    if (logsConnected.value) return

    logsLoading.value = true
    logs.value = ''

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const baseApi = import.meta.env.VITE_BASE_API || ''
    const userStore = useUserStore()
    const wsUrl = `${protocol}//${window.location.host}${baseApi}/api/v1/notebook/log/stream?id=${route.query.id}&token=${userStore.token}`

    logsWs = new WebSocket(wsUrl)

    logsWs.onopen = () => {
      logsConnected.value = true
      logsLoading.value = false
    }

    logsWs.onmessage = async (event) => {
      const processedData = processLogData(event.data)
      logs.value += processedData
      await nextTick()
      if (logsRef.value) {
        logsRef.value.scrollTop = logsRef.value.scrollHeight
      }
    }

    logsWs.onclose = () => {
      logsConnected.value = false
      logsLoading.value = false
    }

    logsWs.onerror = (error) => {
      console.error('WebSocket 错误', error)
      logsConnected.value = false
      logsLoading.value = false
      ElMessage.error(t('error'))
    }
  }

  const disconnectLogStream = () => {
    if (logsWs) {
      logsWs.close()
      logsWs = null
    }
    logsConnected.value = false
  }

  const clearLogs = () => {
    logs.value = ''
  }

  const fetchPods = async () => {
    podsLoading.value = true
    try {
      const res = await getNotebookPods({ id: route.query.id })
      if (res.code === 0) {
        pods.value = res.data || []
      }
    } catch (error) {
      console.error('获取实例列表失败', error)
    } finally {
      podsLoading.value = false
    }
  }

  const fitTerminal = () => {
    if (fitAddon) {
      fitAddon.fit()
    }
  }

  const disconnectTerminal = () => {
    if (ws) {
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

  const connectTerminal = async () => {
    if (!isInstanceRunning.value) {
      ElMessage.warning(`${t('instanceStatus')}: ${notebook.value.status}`)
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
    terminal.open(terminalRef.value)

    resizeObserver = new ResizeObserver(() => {
      fitAddon.fit()
    })
    if (terminalRef.value) {
      resizeObserver.observe(terminalRef.value)
    }

    setTimeout(() => {
      fitAddon.fit()
    }, 50)

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const baseApi = import.meta.env.VITE_BASE_API || ''
    const userStore = useUserStore()
    const wsUrl = `${protocol}//${window.location.host}${baseApi}/api/v1/notebook/terminal?id=${route.query.id}&token=${userStore.token}`

    ws = new WebSocket(wsUrl)

    ws.onopen = () => {
      terminalConnected.value = true
      fitAddon.fit()
      terminal.focus()
    }

    ws.onmessage = (event) => {
      terminal.write(event.data)
    }

    ws.onclose = () => {
      terminalConnected.value = false
      if (terminal) {
        terminal.write('\r\n\x1b[31mConnection disconnected.\x1b[0m\r\n')
      }
    }

    ws.onerror = () => {
      ElMessage.error(t('error'))
    }

    terminal.onData((data) => {
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(data)
      }
    })
  }

  const viewPodLogs = () => {
    activeTab.value = 'logs'
    connectLogStream()
  }

  const handleStart = async () => {
    if (actionLoading.value) return
    try {
      await ElMessageBox.confirm(t('confirm'), t('tip'), { type: 'info' })
      actionLoading.value = true
      const res = await startNotebook({ id: Number(route.query.id) })
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
      await ElMessageBox.confirm(t('confirm'), t('tip'), { type: 'warning' })
      actionLoading.value = true
      const res = await stopNotebook({ id: Number(route.query.id) })
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
      await ElMessageBox.confirm(t('confirmDelete', { name: notebook.value.displayName }), t('tip'), { type: 'warning' })
      actionLoading.value = true
      const res = await deleteNotebook({ id: Number(route.query.id) })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        router.push({ name: 'notebooks' })
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (error) {
    } finally {
      actionLoading.value = false
    }
  }

  const goBack = () => router.go(-1)

  const getPayTypeLabel = (type) => {
    const map = { 1: 'payHourly', 2: 'payDaily', 3: 'payWeekly', 4: 'payMonthly' }
    return t(map[type] || 'payHourly')
  }

  const getUnitPriceLabel = (type) => {
    const map = { 1: 'unitHour', 2: 'unitDay', 3: 'unitWeek', 4: 'unitMonth' }
    return t(map[type] || 'unitHour')
  }

  const getStatusLabel = (status) => t(status) || status

  const getStatusClass = (status) => {
    const normalized = status?.toLowerCase()
    if (normalized === 'running' || normalized === 'ready') return 'bg-emerald-500/10 text-emerald-500'
    if (normalized === 'pending' || normalized === 'creating') return 'bg-amber-500/10 text-amber-500'
    if (normalized === 'failed' || normalized === 'error') return 'bg-red-500/10 text-red-500'
    if (normalized === 'stopped') return 'bg-slate-500/10 text-slate-500'
    return 'bg-slate-500/10 text-slate-500'
  }

  const getPodStatusClass = (status) => {
    const normalized = status?.toLowerCase()
    const map = {
      running: 'bg-emerald-500/10 text-emerald-500',
      pending: 'bg-amber-500/10 text-amber-500',
      failed: 'bg-red-500/10 text-red-500'
    }
    return map[normalized] || 'bg-slate-500/10 text-slate-500'
  }

  const formatTime = (time) => {
    if (!time) return '-'
    return new Date(time).toLocaleString()
  }

  watch(activeTab, (tab) => {
    if (tab === 'logs' && notebook.value.status === 'Running') {
      connectLogStream()
    } else if (tab !== 'logs') {
      disconnectLogStream()
    }

    if (tab === 'terminal') {
      nextTick(() => {
        if (terminal && fitAddon) fitAddon.fit()
      })
    }

    if (tab === 'pods') fetchPods()
  })

  onMounted(() => {
    fetchDetail()
    refreshTimer = setInterval(fetchDetail, 5000)
  })

  onUnmounted(() => {
    if (refreshTimer) clearInterval(refreshTimer)
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
