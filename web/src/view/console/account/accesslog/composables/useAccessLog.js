import { inject, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getAccessLogList, getActiveSessionList, killSession } from '@/api/account'

export function useAccessLog({ t }) {
  const translate = t || inject('t', (key) => key)

  const loading = ref(false)
  const activeSessions = ref([])
  const loginLogs = ref([])
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)
  const searchIp = ref('')
  const filterStatus = ref('')

  const fetchSessions = async () => {
    try {
      const res = await getActiveSessionList({ page: 1, pageSize: 20 })
      if (res.code === 0) {
        activeSessions.value = res.data?.list || []
      }
    } catch (error) {
      console.error(error)
    }
  }

  const fetchLogs = async (showLoading = true) => {
    if (showLoading) {
      loading.value = true
    }

    try {
      const res = await getAccessLogList({
        page: page.value,
        pageSize: pageSize.value,
        ip: searchIp.value || undefined,
        status: filterStatus.value || undefined
      })

      if (res.code === 0) {
        loginLogs.value = res.data?.list || []
        total.value = res.data?.total || 0
      }
    } catch (error) {
      console.error(error)
    } finally {
      loading.value = false
    }
  }

  const fetchAll = async (showLoading = true) => {
    await Promise.all([fetchSessions(), fetchLogs(showLoading)])
  }

  watch([searchIp, filterStatus], () => {
    page.value = 1
    fetchLogs(true)
  })

  const handleRefresh = (silent = false) => fetchAll(!silent)
  const handleSearch = () => {
    page.value = 1
    fetchLogs(true)
  }
  const handlePageChange = (value) => {
    page.value = value
    fetchLogs(true)
  }
  const handleSizeChange = (value) => {
    pageSize.value = value
    page.value = 1
    fetchLogs(true)
  }

  const handleKill = (id) => {
    ElMessageBox.confirm(translate('audit.confirmKill'), translate('tip'), {
      confirmButtonText: translate('confirm'),
      cancelButtonText: translate('cancel'),
      type: 'warning'
    })
      .then(async () => {
        const res = await killSession({ logId: id })
        if (res.code === 0) {
          ElMessage.success(translate('success'))
          fetchAll(false)
        }
      })
      .catch(() => {})
  }

  return {
    activeSessions,
    fetchAll,
    filterStatus,
    handleKill,
    handlePageChange,
    handleRefresh,
    handleSearch,
    handleSizeChange,
    loading,
    loginLogs,
    page,
    pageSize,
    searchIp,
    total
  }
}
