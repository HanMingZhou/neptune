import { inject, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getAccessLogList,
  getActiveSessionList,
  killSession
} from '@/api/account'
import type { ApiResponse } from '@/utils/request'
import type { PageListData, Translator } from '@/types/consoleResource'
import type { AccessLogRecord, AccessLogSession } from '@/types/account'

interface UseAccessLogOptions {
  t?: Translator
}

export function useAccessLog({ t }: UseAccessLogOptions = {}) {
  const translate = t || inject<Translator>('t', (key: string) => key)

  const loading = ref(false)
  const activeSessions = ref<AccessLogSession[]>([])
  const loginLogs = ref<AccessLogRecord[]>([])
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)
  const searchIp = ref('')
  const filterStatus = ref('')

  const fetchSessions = async (): Promise<void> => {
    try {
      const res = (await getActiveSessionList({
        page: 1,
        pageSize: 20
      })) as ApiResponse<PageListData<AccessLogSession>>
      if (res.code === 0) {
        activeSessions.value = res.data?.list || []
      }
    } catch (error) {
      console.error(error)
    }
  }

  const fetchLogs = async (showLoading = true): Promise<void> => {
    if (showLoading) {
      loading.value = true
    }

    try {
      const res = (await getAccessLogList({
        page: page.value,
        pageSize: pageSize.value,
        ip: searchIp.value || undefined,
        status: filterStatus.value || undefined
      })) as ApiResponse<PageListData<AccessLogRecord>>

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

  const fetchAll = async (showLoading = true): Promise<void> => {
    await Promise.all([fetchSessions(), fetchLogs(showLoading)])
  }

  watch([searchIp, filterStatus], () => {
    page.value = 1
    void fetchLogs(true)
  })

  const handleRefresh = (silent = false): Promise<void> => fetchAll(!silent)

  const handleSearch = (): void => {
    page.value = 1
    void fetchLogs(true)
  }

  const handlePageChange = (value: number): void => {
    page.value = value
    void fetchLogs(true)
  }

  const handleSizeChange = (value: number): void => {
    pageSize.value = value
    page.value = 1
    void fetchLogs(true)
  }

  const handleKill = (id: number): void => {
    void ElMessageBox.confirm(
      translate('audit.confirmKill'),
      translate('tip'),
      {
        confirmButtonText: translate('confirm'),
        cancelButtonText: translate('cancel'),
        type: 'warning'
      }
    )
      .then(async () => {
        const res = (await killSession({ logId: id })) as ApiResponse<unknown>
        if (res.code === 0) {
          ElMessage.success(translate('success'))
          await fetchAll(false)
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
