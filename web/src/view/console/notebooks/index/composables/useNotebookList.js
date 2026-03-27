import { inject, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { deleteNotebook, getNotebookList, startNotebook, stopNotebook } from '@/api/notebook'

export const useNotebookList = () => {
  const t = inject('t', (key) => key)
  const router = useRouter()

  const loading = ref(false)
  const notebooks = ref([])
  const showSSHDialog = ref(false)
  const currentNotebook = ref({})
  const showPassword = ref(false)
  const searchQuery = ref('')
  const statusFilter = ref('')
  const page = ref(1)
  const pageSize = ref(20)
  const total = ref(0)
  const btnLoading = ref({})

  let timer = null

  const fetchNotebooks = async (silent = false) => {
    if (!silent) loading.value = true
    try {
      const res = await getNotebookList({
        page: page.value,
        pageSize: pageSize.value,
        displayName: searchQuery.value || undefined,
        status: statusFilter.value || undefined
      })
      if (res.code === 0) {
        notebooks.value = res.data?.list || []
        total.value = res.data?.total || 0
      } else if (!silent) {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (error) {
      console.error(error)
      if (!silent) {
        ElMessage.error(t('error'))
      }
    } finally {
      if (!silent) loading.value = false
    }
  }

  const handlePageChange = (value) => {
    page.value = value
    fetchNotebooks()
  }

  const handleSizeChange = (value) => {
    pageSize.value = value
    page.value = 1
    fetchNotebooks()
  }

  watch([searchQuery, statusFilter], () => {
    page.value = 1
    fetchNotebooks()
  })

  const startPolling = () => {
    stopPolling()
    timer = setInterval(() => {
      fetchNotebooks(true)
    }, 5000)
  }

  const stopPolling = () => {
    if (timer) {
      clearInterval(timer)
      timer = null
    }
  }

  const handleDelete = async (item) => {
    if (btnLoading.value[item.id]) return
    try {
      btnLoading.value[item.id] = true
      await ElMessageBox.confirm(t('confirmDelete', { name: item.displayName }), t('tip'), {
        confirmButtonText: t('confirm'),
        cancelButtonText: t('cancel'),
        type: 'warning'
      })
      const res = await deleteNotebook({ id: item.id })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        item.status = 'DELETING'
      } else {
        ElMessage.error(res.msg || t('operationFailed'))
      }
    } catch (error) {
      if (error !== 'cancel') {
        console.error(error)
        ElMessage.error(t('operationFailed'))
      }
    } finally {
      btnLoading.value[item.id] = false
    }
  }

  const handleStart = async (item) => {
    if (btnLoading.value[item.id]) return
    if (['RUNNING', 'PENDING', 'CREATING', 'DELETING'].includes(item.status)) return
    try {
      btnLoading.value[item.id] = true
      const res = await startNotebook({ id: item.id })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        item.status = 'PENDING'
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (error) {
      console.error(error)
      ElMessage.error(t('operationFailed'))
    } finally {
      btnLoading.value[item.id] = false
    }
  }

  const handleStop = async (item) => {
    if (btnLoading.value[item.id]) return
    if (item.status !== 'RUNNING') return
    try {
      btnLoading.value[item.id] = true
      const res = await stopNotebook({ id: item.id })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        item.status = 'STOPPED'
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (error) {
      console.error(error)
      ElMessage.error(t('operationFailed'))
    } finally {
      btnLoading.value[item.id] = false
    }
  }

  const showSSHInfo = (item) => {
    currentNotebook.value = item
    showPassword.value = false
    showSSHDialog.value = true
  }

  const copyText = (text) => {
    navigator.clipboard.writeText(text).then(() => {
      ElMessage.success(t('copied'))
    }).catch(() => {
      ElMessage.error(t('copyFailed'))
    })
  }

  const showKeySettings = () => {
    router.push({ name: 'sshkeys' }).catch(() => {
      router.push('/layout/sshkeys')
    })
  }

  const goToCreate = () => {
    router.push({ name: 'notebookCreate' }).catch(() => {
      router.push('/layout/notebooks/create')
    })
  }

  const goToDetail = (item) => {
    router.push({ name: 'notebookDetail', query: { id: item.id } }).catch(() => {
      router.push({ path: '/layout/notebooks/detail', query: { id: item.id } })
    })
  }

  const getStatusText = (status) => {
    return t(status) || status
  }

  const getStatusStyle = (status) => {
    const map = {
      RUNNING: 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20',
      STOPPED: 'bg-slate-500/10 text-slate-500 border-slate-500/20',
      CREATING: 'bg-blue-500/10 text-blue-500 border-blue-500/20',
      PENDING: 'bg-amber-500/10 text-amber-500 border-amber-500/20',
      FAILED: 'bg-red-500/10 text-red-500 border-red-500/20',
      SUCCEEDED: 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20',
      DELETING: 'bg-red-500/10 text-red-500 border-red-500/20',
      UPDATING: 'bg-blue-500/10 text-blue-500 border-blue-500/20',
      UPDATE_FAILED: 'bg-red-500/10 text-red-500 border-red-500/20',
      DELETE_FAILED: 'bg-red-500/10 text-red-500 border-red-500/20',
      UPDATED: 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20'
    }
    return map[status] || 'bg-slate-500/10 text-slate-500'
  }

  onMounted(() => {
    fetchNotebooks()
    startPolling()
  })

  onUnmounted(() => {
    stopPolling()
  })

  return {
    btnLoading,
    copyText,
    currentNotebook,
    fetchNotebooks,
    getStatusStyle,
    getStatusText,
    goToCreate,
    goToDetail,
    handleDelete,
    handlePageChange,
    handleSizeChange,
    handleStart,
    handleStop,
    loading,
    notebooks,
    page,
    pageSize,
    searchQuery,
    showKeySettings,
    showPassword,
    showSSHDialog,
    showSSHInfo,
    statusFilter,
    total
  }
}
