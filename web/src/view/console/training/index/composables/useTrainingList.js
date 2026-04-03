import { inject, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { deleteTrainingJob, getTrainingJobList, stopTrainingJob } from '@/api/training'

export const useTrainingList = () => {
  const t = inject('t', (key) => key)
  const router = useRouter()

  const loading = ref(false)
  const jobs = ref([])
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(20)
  const searchQuery = ref('')
  const filterStatus = ref('')
  const btnLoading = ref({})
  const isInitialLoad = ref(true)

  const syncJobs = (nextList = []) => {
    const currentMap = new Map(jobs.value.map((item) => [item.id, item]))
    jobs.value = nextList.map((item) => {
      const existing = currentMap.get(item.id)
      return existing ? Object.assign(existing, item) : item
    })
  }

  const fetchList = async (showLoading = false) => {
    if (showLoading || isInitialLoad.value) {
      loading.value = true
    }
    try {
      const res = await getTrainingJobList({
        page: page.value,
        pageSize: pageSize.value,
        name: searchQuery.value || undefined,
        status: filterStatus.value || undefined
      })
      if (res.code === 0) {
        syncJobs(res.data?.list || [])
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

  watch([searchQuery, filterStatus], () => {
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

  const getFrameworkLabel = (type) => {
    const map = { PYTORCH_DDP: 'PyTorchDDP', STANDALONE: 'StandAlone', MPI: 'MPI' }
    return map[type] || type
  }

  const getFrameworkStyle = (type) => {
    const map = {
      PYTORCH_DDP: 'console-badge--warning',
      STANDALONE: 'console-badge--info',
      MPI: 'console-badge--neutral'
    }
    return map[type] || 'console-badge--neutral'
  }

  const getStatusLabel = (status) => {
    return t(status) || status
  }

  const getStatusStyle = (status) => {
    const map = {
      RUNNING: 'console-badge--success',
      PENDING: 'console-badge--warning',
      SUCCEEDED: 'console-badge--info',
      FAILED: 'console-badge--danger',
      KILLED: 'console-badge--neutral'
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

  const goToCreate = () => router.push({ name: 'trainingCreate' })
  const goToDetail = (item) => router.push({ name: 'trainingDetail', query: { id: item.id } })
  const viewLogs = (item) => router.push({ name: 'trainingDetail', query: { id: item.id, tab: 'logs' } })

  const openTensorboard = (item) => {
    if (item.tensorboardUrl) {
      window.open(item.tensorboardUrl, '_blank')
    }
  }

  const handleStop = async (item) => {
    if (btnLoading.value[item.id]) return
    try {
      const name = item.displayName || item.jobName
      btnLoading.value[item.id] = true
      await ElMessageBox.confirm(t('confirmStop', { name }), t('tip'), {
        confirmButtonText: t('confirm'),
        cancelButtonText: t('cancel'),
        type: 'warning'
      })
      const res = await stopTrainingJob({ id: item.id })
      if (res.code === 0) {
        ElMessage.success(t('success'))
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

  const handleDelete = async (item) => {
    if (btnLoading.value[item.id]) return
    try {
      const name = item.displayName || item.jobName
      btnLoading.value[item.id] = true
      await ElMessageBox.confirm(t('confirmDelete', { name }), t('tip'), {
        confirmButtonText: t('confirm'),
        cancelButtonText: t('cancel'),
        type: 'warning'
      })
      const res = await deleteTrainingJob({ id: item.id })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        item.status = 'DELETING'
        fetchList()
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (error) {
    } finally {
      if (item.id) {
        btnLoading.value[item.id] = false
      }
    }
  }

  onMounted(() => {
    fetchList()
  })

  return {
    btnLoading,
    copyText,
    filterStatus,
    getFrameworkLabel,
    getFrameworkStyle,
    getStatusLabel,
    getStatusStyle,
    goToCreate,
    goToDetail,
    handleDelete,
    handlePageChange,
    handleRefresh,
    handleSizeChange,
    handleStop,
    jobs,
    loading,
    openTensorboard,
    page,
    pageSize,
    searchQuery,
    total,
    viewLogs
  }
}
