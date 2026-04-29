import { inject, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { AxiosError } from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  deleteInferenceService,
  getInferenceServiceList,
  startInferenceService,
  stopInferenceService
} from '@/api/inference'
import type { ApiResponse } from '@/utils/request'
import { getErrorMessage } from '@/utils/resourceValidators'
import type {
  ConsoleInferenceService,
  PageListData,
  ResourceId,
  Translator
} from '@/types/consoleResource'

type ButtonLoadingState = Partial<Record<ResourceId, boolean>>

const normalizeStatus = (status: unknown): string =>
  `${status || ''}`.trim().toUpperCase()

const extractErrorMessage = (error: unknown, fallback: string): string => {
  if (error instanceof AxiosError) {
    return error.response?.data?.msg || error.message || fallback
  }
  if (error instanceof Error) {
    return error.message || fallback
  }
  return fallback
}

export const useInferenceList = () => {
  const t = inject<Translator>('t', (key) => key)
  const router = useRouter()

  const loading = ref(false)
  const services = ref<ConsoleInferenceService[]>([])
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(15)
  const searchQuery = ref('')
  const filterStatus = ref('')
  const filterFramework = ref('')
  const btnLoading = ref<ButtonLoadingState>({})
  const isInitialLoad = ref(true)

  const syncServices = (nextList: ConsoleInferenceService[] = []): void => {
    const currentMap = new Map<ResourceId, ConsoleInferenceService>(
      services.value.map((item) => [item.id, item] as const)
    )

    services.value = nextList.map((item) => {
      const existing = currentMap.get(item.id)
      return existing ? Object.assign(existing, item) : item
    })
  }

  const fetchList = async (showLoading = false): Promise<void> => {
    if (showLoading || isInitialLoad.value) {
      loading.value = true
    }

    try {
      const res = (await getInferenceServiceList({
        page: page.value,
        pageSize: pageSize.value,
        name: searchQuery.value || undefined,
        status: filterStatus.value || undefined,
        framework: filterFramework.value || undefined
      })) as ApiResponse<PageListData<ConsoleInferenceService>>

      if (res.code === 0) {
        syncServices(res.data?.list || [])
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

  watch([searchQuery, filterStatus, filterFramework], () => {
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

  const getFrameworkStyle = (type?: string): string => {
    const frameworkStyleMap: Record<string, string> = {
      SGLANG: 'console-badge--neutral',
      VLLM: 'console-badge--warning'
    }

    return frameworkStyleMap[type || ''] || 'console-badge--neutral'
  }

  const getDeployTypeStyle = (type?: string): string => {
    return type === 'STANDALONE'
      ? 'console-badge--info'
      : 'console-badge--success'
  }

  const getStatusLabel = (status?: string): string => {
    return t(normalizeStatus(status)) || status || '-'
  }

  const getStatusStyle = (status?: string): string => {
    const statusStyleMap: Record<string, string> = {
      RUNNING: 'console-badge--success',
      PENDING: 'console-badge--warning',
      CREATING: 'console-badge--warning',
      STOPPED: 'console-badge--neutral',
      FAILED: 'console-badge--danger',
      DELETING: 'console-badge--info',
      RESTARTING: 'console-badge--info'
    }

    return statusStyleMap[normalizeStatus(status)] || 'console-badge--neutral'
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
    void router.push({ name: 'inferenceCreate' })
  }

  const goToDetail = (item: ConsoleInferenceService): void => {
    void router.push({ name: 'inferenceDetail', query: { id: item.id } })
  }

  const goToEdit = (item: ConsoleInferenceService): void => {
    if (normalizeStatus(item.status) !== 'STOPPED') return
    void router.push({ name: 'inferenceCreate', query: { id: item.id } })
  }

  const viewLogs = (item: ConsoleInferenceService): void => {
    void router.push({
      name: 'inferenceDetail',
      query: { id: item.id, tab: 'logs' }
    })
  }

  const handleStart = async (item: ConsoleInferenceService): Promise<void> => {
    if (btnLoading.value[item.id]) return
    if (!['STOPPED', 'FAILED'].includes(normalizeStatus(item.status))) return

    try {
      btnLoading.value[item.id] = true
      await ElMessageBox.confirm(
        t('confirmStart', {
          name: item.displayName || item.instanceName || '-'
        }),
        t('tip'),
        {
          type: 'warning'
        }
      )

      const res = await startInferenceService({ id: item.id })
      if (res.code === 0) {
        item.status = 'PENDING'
        await fetchList()
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (error) {
      if (error === 'cancel' || error === 'close') {
        return
      }
      ElMessage.error(extractErrorMessage(error, t('error')))
    } finally {
      btnLoading.value[item.id] = false
    }
  }

  const handleStop = async (item: ConsoleInferenceService): Promise<void> => {
    if (btnLoading.value[item.id]) return
    if (!['RUNNING', 'PENDING'].includes(normalizeStatus(item.status))) return

    try {
      btnLoading.value[item.id] = true
      await ElMessageBox.confirm(
        t('confirmStop', {
          name: item.displayName || item.instanceName || '-'
        }),
        t('tip'),
        {
          confirmButtonText: t('confirm'),
          cancelButtonText: t('cancel'),
          type: 'warning'
        }
      )
      const res = await stopInferenceService({ id: item.id })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        await fetchList(true)
      }
    } catch (_error) {
      if (_error !== 'cancel') console.error(_error)
    } finally {
      btnLoading.value[item.id] = false
    }
  }

  const handleDelete = async (item: ConsoleInferenceService): Promise<void> => {
    if (btnLoading.value[item.id]) return
    if (!['STOPPED', 'FAILED'].includes(normalizeStatus(item.status))) return

    try {
      btnLoading.value[item.id] = true
      await ElMessageBox.confirm(
        t('confirmDeleteInference', {
          name: item.displayName || item.instanceName || '-'
        }),
        t('tip'),
        {
          confirmButtonText: t('confirm'),
          cancelButtonText: t('cancel'),
          type: 'warning'
        }
      )
      const res = await deleteInferenceService({ id: item.id })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        await fetchList(true)
      }
    } catch (_error) {
      if (_error !== 'cancel') console.error(_error)
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
    filterFramework,
    filterStatus,
    getDeployTypeStyle,
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

