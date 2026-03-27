import { inject, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  deleteInferenceService,
  getInferenceServiceList,
  startInferenceService,
  stopInferenceService
} from '@/api/inference'

export const useInferenceList = () => {
  const t = inject('t', (key) => key)
  const router = useRouter()

  const loading = ref(false)
  const services = ref([])
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(20)
  const searchQuery = ref('')
  const filterStatus = ref('')
  const filterFramework = ref('')
  const btnLoading = ref({})
  const isInitialLoad = ref(true)

  let refreshTimer = null

  const fetchList = async (showLoading = false) => {
    if (showLoading || isInitialLoad.value) {
      loading.value = true
    }
    try {
      const res = await getInferenceServiceList({
        page: page.value,
        pageSize: pageSize.value,
        name: searchQuery.value || undefined,
        status: filterStatus.value || undefined,
        framework: filterFramework.value || undefined
      })
      if (res.code === 0) {
        services.value = res.data?.list || []
        total.value = res.data?.total || 0
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (error) {
      console.error(error)
      ElMessage.error(t('error'))
    } finally {
      loading.value = false
      isInitialLoad.value = false
    }
  }

  watch([searchQuery, filterStatus, filterFramework], () => {
    page.value = 1
    fetchList(true)
  })

  const handleRefresh = (silent = false) => fetchList(!silent)

  const handlePageChange = (value) => {
    page.value = value
    fetchList(true)
  }

  const handleSizeChange = (value) => {
    pageSize.value = value
    page.value = 1
    fetchList(true)
  }

  const getFrameworkStyle = (type) => {
    const map = {
      SGLANG: 'bg-purple-500/10 text-purple-500 border-purple-500/20',
      VLLM: 'bg-orange-500/10 text-orange-500 border-orange-500/20'
    }
    return map[type] || 'bg-slate-500/10 text-slate-500 border-slate-500/20'
  }

  const getDeployTypeStyle = (type) => {
    return type === 'STANDALONE'
      ? 'bg-blue-500/10 text-blue-500'
      : 'bg-emerald-500/10 text-emerald-500'
  }

  const getStatusLabel = (status) => {
    return t(status) || status
  }

  const getStatusStyle = (status) => {
    const map = {
      RUNNING: 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20',
      PENDING: 'bg-amber-500/10 text-amber-500 border-amber-500/20',
      CREATING: 'bg-amber-500/10 text-amber-500 border-amber-500/20',
      STOPPED: 'bg-slate-500/10 text-slate-500 border-slate-500/20',
      FAILED: 'bg-red-500/10 text-red-500 border-red-500/20',
      DELETING: 'bg-blue-500/10 text-blue-500 border-blue-500/20',
      RESTARTING: 'bg-blue-500/10 text-blue-500 border-blue-500/20'
    }
    return map[status] || 'bg-slate-500/10 text-slate-500 border-slate-500/20'
  }

  const copyText = (text) => {
    navigator.clipboard.writeText(text).then(() => {
      ElMessage.success(t('copied') || 'Copied')
    }).catch(() => {
      ElMessage.error(t('copyFailed') || 'Copy failed')
    })
  }

  const goToCreate = () => router.push({ name: 'inferenceCreate' })
  const goToDetail = (item) => router.push({ name: 'inferenceDetail', query: { id: item.id } })
  const viewLogs = (item) => router.push({ name: 'inferenceDetail', query: { id: item.id, tab: 'logs' } })

  const handleStart = async (item) => {
    if (btnLoading.value[item.id]) return
    try {
      btnLoading.value[item.id] = true
      await ElMessageBox.confirm(t('confirmStart', { name: item.displayName }), t('tip'), { type: 'warning' })
      const res = await startInferenceService({ id: item.id })
      if (res.code === 0) {
        item.status = 'PENDING'
        fetchList()
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (error) {
    } finally {
      btnLoading.value[item.id] = false
    }
  }

  const handleStop = async (item) => {
    if (btnLoading.value[item.id]) return
    try {
      btnLoading.value[item.id] = true
      await ElMessageBox.confirm(t('confirmStop', { name: item.displayName }), t('tip'), {
        confirmButtonText: t('confirm'),
        cancelButtonText: t('cancel'),
        type: 'warning'
      })
      const res = await stopInferenceService({ id: item.id })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        fetchList(true)
      }
    } catch (error) {
      if (error !== 'cancel') console.error(error)
    } finally {
      btnLoading.value[item.id] = false
    }
  }

  const handleDelete = async (item) => {
    if (btnLoading.value[item.id]) return
    try {
      btnLoading.value[item.id] = true
      await ElMessageBox.confirm(t('confirmDeleteInference', { name: item.displayName }), t('tip'), {
        confirmButtonText: t('confirm'),
        cancelButtonText: t('cancel'),
        type: 'warning'
      })
      const res = await deleteInferenceService({ id: item.id })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        fetchList(true)
      }
    } catch (error) {
      if (error !== 'cancel') console.error(error)
    } finally {
      btnLoading.value[item.id] = false
    }
  }

  const startPolling = () => {
    stopPolling()
    refreshTimer = setInterval(() => {
      fetchList(false)
    }, 10000)
  }

  const stopPolling = () => {
    if (refreshTimer) {
      clearInterval(refreshTimer)
      refreshTimer = null
    }
  }

  onMounted(() => {
    fetchList()
    startPolling()
  })

  onUnmounted(() => {
    stopPolling()
  })

  return {
    btnLoading,
    copyText,
    filterFramework,
    filterStatus,
    getDeployTypeStyle,
    getFrameworkStyle,
    getStatusLabel,
    getStatusStyle,
    goToCreate,
    goToDetail,
    handleDelete,
    handlePageChange,
    handleRefresh,
    handleSizeChange,
    handleStart,
    handleStop,
    loading,
    page,
    pageSize,
    searchQuery,
    services,
    total,
    viewLogs
  }
}
