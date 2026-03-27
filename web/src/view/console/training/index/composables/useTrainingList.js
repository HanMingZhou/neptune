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
        jobs.value = res.data?.list || []
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
    const map = { PYTORCH_DDP: 'PyTorch DDP', STANDALONE: 'Standalone', MPI: 'MPI' }
    return map[type] || type
  }

  const getFrameworkStyle = (type) => {
    const map = {
      PYTORCH_DDP: 'bg-orange-500/10 text-orange-500 border-orange-500/20',
      STANDALONE: 'bg-blue-500/10 text-blue-500 border-blue-500/20',
      MPI: 'bg-purple-500/10 text-purple-500 border-purple-500/20'
    }
    return map[type] || 'bg-slate-500/10 text-slate-500 border-slate-500/20'
  }

  const getStatusLabel = (status) => {
    return t(status) || status
  }

  const getStatusStyle = (status) => {
    const map = {
      RUNNING: 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20',
      PENDING: 'bg-amber-500/10 text-amber-500 border-amber-500/20',
      SUCCEEDED: 'bg-blue-500/10 text-blue-500 border-blue-500/20',
      FAILED: 'bg-red-500/10 text-red-500 border-red-500/20',
      KILLED: 'bg-slate-500/10 text-slate-500 border-slate-500/20'
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
