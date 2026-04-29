import { inject, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  deleteTrainingJob,
  getTrainingJobList,
  stopTrainingJob
} from '@/api/training'
import type { ApiResponse } from '@/utils/request'
import { getErrorMessage } from '@/utils/resourceValidators'
import type {
  ConsoleTrainingJob,
  PageListData,
  ResourceId,
  Translator
} from '@/types/consoleResource'

type ButtonLoadingState = Partial<Record<ResourceId, boolean>>

const isEditableStatus = (status?: string): boolean =>
  ['KILLED', 'SUCCEEDED'].includes(`${status || ''}`.toUpperCase())

export const useTrainingList = () => {
  const t = inject<Translator>('t', (key) => key)
  const router = useRouter()

  const loading = ref(false)
  const jobs = ref<ConsoleTrainingJob[]>([])
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(15)
  const searchQuery = ref('')
  const filterStatus = ref('')
  const btnLoading = ref<ButtonLoadingState>({})
  const isInitialLoad = ref(true)

  const syncJobs = (nextList: ConsoleTrainingJob[] = []): void => {
    const currentMap = new Map<ResourceId, ConsoleTrainingJob>(
      jobs.value.map((item) => [item.id, item] as const)
    )

    jobs.value = nextList.map((item) => {
      const existing = currentMap.get(item.id)
      return existing ? Object.assign(existing, item) : item
    })
  }

  const fetchList = async (showLoading = false): Promise<void> => {
    if (showLoading || isInitialLoad.value) {
      loading.value = true
    }

    try {
      const res = (await getTrainingJobList({
        page: page.value,
        pageSize: pageSize.value,
        name: searchQuery.value || undefined,
        status: filterStatus.value || undefined
      })) as ApiResponse<PageListData<ConsoleTrainingJob>>

      if (res.code === 0) {
        syncJobs(res.data?.list || [])
        total.value = res.data?.total || 0
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (error) {
      console.error(error)
      ElMessage.error(getErrorMessage(error, t('error')))
    } finally {
      loading.value = false
      isInitialLoad.value = false
    }
  }

  watch([searchQuery, filterStatus], () => {
    page.value = 1
    void fetchList(true)
  })

  const handleRefresh = (silent = false): Promise<void> => fetchList(!silent)

  const handlePageChange = (value: number): void => {
    page.value = value
    void fetchList(true)
  }

  const handleSizeChange = (value: number): void => {
    pageSize.value = value
    page.value = 1
    void fetchList(true)
  }

  const getFrameworkLabel = (type?: string): string => {
    const frameworkLabelMap: Record<string, string> = {
      PYTORCH_DDP: 'PyTorchDDP',
      STANDALONE: 'StandAlone',
      MPI: 'MPI'
    }

    return frameworkLabelMap[type || ''] || type || '-'
  }

  const getFrameworkStyle = (type?: string): string => {
    const frameworkStyleMap: Record<string, string> = {
      PYTORCH_DDP: 'console-badge--warning',
      STANDALONE: 'console-badge--info',
      MPI: 'console-badge--neutral'
    }

    return frameworkStyleMap[type || ''] || 'console-badge--neutral'
  }

  const getStatusLabel = (status?: string): string => {
    return t(status || '') || status || '-'
  }

  const getStatusStyle = (status?: string): string => {
    const statusStyleMap: Record<string, string> = {
      RUNNING: 'console-badge--success',
      PENDING: 'console-badge--warning',
      SUCCEEDED: 'console-badge--info',
      FAILED: 'console-badge--danger',
      KILLED: 'console-badge--neutral'
    }

    return statusStyleMap[status || ''] || 'console-badge--neutral'
  }

  const copyText = (text?: string): void => {
    if (!text) return

    void navigator.clipboard
      .writeText(text)
      .then(() => {
        ElMessage.success(t('copied') || 'Copied')
      })
      .catch(() => {
        ElMessage.error(t('copyFailed') || 'Copy failed')
      })
  }

  const goToCreate = (): void => {
    void router.push({ name: 'trainingCreate' })
  }

  const goToDetail = (item: ConsoleTrainingJob): void => {
    void router.push({ name: 'trainingDetail', query: { id: item.id } })
  }

  const goToEdit = (item: ConsoleTrainingJob): void => {
    if (!isEditableStatus(item.status)) return
    void router.push({ name: 'trainingCreate', query: { id: item.id } })
  }

  const viewLogs = (item: ConsoleTrainingJob): void => {
    void router.push({
      name: 'trainingDetail',
      query: { id: item.id, tab: 'logs' }
    })
  }

  const openTensorboard = (item: ConsoleTrainingJob): void => {
    if (item.tensorboardUrl) {
      window.open(item.tensorboardUrl, '_blank')
    }
  }

  const handleStop = async (item: ConsoleTrainingJob): Promise<void> => {
    if (btnLoading.value[item.id]) return

    try {
      const name = item.displayName || item.jobName || '-'
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
        await fetchList()
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (_error) {
    } finally {
      btnLoading.value[item.id] = false
    }
  }

  const handleDelete = async (item: ConsoleTrainingJob): Promise<void> => {
    if (btnLoading.value[item.id]) return

    try {
      const name = item.displayName || item.jobName || '-'
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
        await fetchList()
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (_error) {
    } finally {
      btnLoading.value[item.id] = false
    }
  }

  onMounted(() => {
    void fetchList()
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
    goToEdit,
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

