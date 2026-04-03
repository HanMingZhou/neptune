import { inject, onMounted, ref, watch } from 'vue'
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

  const syncServices = (nextList = []) => {
    const currentMap = new Map(services.value.map((item) => [item.id, item]))
    services.value = nextList.map((item) => {
      const existing = currentMap.get(item.id)
      return existing ? Object.assign(existing, item) : item
    })
  }

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
        syncServices(res.data?.list || [])
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
      SGLANG: 'console-badge--neutral',
      VLLM: 'console-badge--warning'
    }
    return map[type] || 'console-badge--neutral'
  }

  const getDeployTypeStyle = (type) => {
    return type === 'STANDALONE'
      ? 'console-badge--info'
      : 'console-badge--success'
  }

  const getStatusLabel = (status) => {
    return t(status) || status
  }

  const getStatusStyle = (status) => {
    const map = {
      RUNNING: 'console-badge--success',
      PENDING: 'console-badge--warning',
      CREATING: 'console-badge--warning',
      STOPPED: 'console-badge--neutral',
      FAILED: 'console-badge--danger',
      DELETING: 'console-badge--info',
      RESTARTING: 'console-badge--info'
    }
    return map[status] || 'console-badge--neutral'
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

  onMounted(() => {
    fetchList()
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
