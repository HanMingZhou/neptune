import { inject, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  deleteNotebook,
  getNotebookList,
  startNotebook,
  stopNotebook
} from '@/api/notebook'
import type { ApiResponse } from '@/utils/request'
import type {
  ConsoleNotebook,
  PageListData,
  ResourceId,
  Translator
} from '@/types/consoleResource'

type ButtonLoadingState = Partial<Record<ResourceId, boolean>>

const normalizeStatus = (status: unknown): string =>
  `${status || ''}`.trim().toUpperCase()

export const useNotebookList = () => {
  const t = inject<Translator>('t', (key) => key)
  const router = useRouter()

  const loading = ref(false)
  const notebooks = ref<ConsoleNotebook[]>([])
  const showSSHDialog = ref(false)
  const currentNotebook = ref<ConsoleNotebook>({})
  const showPassword = ref(false)
  const searchQuery = ref('')
  const statusFilter = ref('')
  const page = ref(1)
  const pageSize = ref(20)
  const total = ref(0)
  const btnLoading = ref<ButtonLoadingState>({})

  const syncNotebooks = (nextList: ConsoleNotebook[] = []): void => {
    const currentMap = new Map<ResourceId, ConsoleNotebook>(
      notebooks.value
        .filter(
          (item): item is ConsoleNotebook & { id: ResourceId } =>
            item.id !== undefined
        )
        .map((item) => [item.id, item] as const)
    )

    notebooks.value = nextList.map((item) => {
      if (item.id === undefined) return item
      const existing = currentMap.get(item.id)
      return existing ? Object.assign(existing, item) : item
    })
  }

  const fetchNotebooks = async (silent = false): Promise<void> => {
    if (!silent) loading.value = true

    try {
      const res = (await getNotebookList({
        page: page.value,
        pageSize: pageSize.value,
        displayName: searchQuery.value || undefined,
        status: statusFilter.value || undefined
      })) as ApiResponse<PageListData<ConsoleNotebook>>

      if (res.code === 0) {
        syncNotebooks(res.data?.list || [])
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
      loading.value = false
    }
  }

  const handlePageChange = (value: number): void => {
    page.value = value
    void fetchNotebooks()
  }

  const handleSizeChange = (value: number): void => {
    pageSize.value = value
    page.value = 1
    void fetchNotebooks()
  }

  watch([searchQuery, statusFilter], () => {
    page.value = 1
    void fetchNotebooks()
  })

  const handleDelete = async (item: ConsoleNotebook): Promise<void> => {
    if (item.id === undefined || btnLoading.value[item.id]) return

    try {
      btnLoading.value[item.id] = true
      await ElMessageBox.confirm(
        t('confirmDelete', {
          name: item.displayName || item.instanceName || '-'
        }),
        t('tip'),
        {
          confirmButtonText: t('confirm'),
          cancelButtonText: t('cancel'),
          type: 'warning'
        }
      )

      const res = await deleteNotebook({ id: item.id })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        item.status = 'DELETING'
        await fetchNotebooks(true)
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

  const handleStart = async (item: ConsoleNotebook): Promise<void> => {
    if (item.id === undefined || btnLoading.value[item.id]) return

    const normalizedStatus = normalizeStatus(item.status)
    if (normalizedStatus !== 'STOPPED') return

    try {
      btnLoading.value[item.id] = true
      const res = await startNotebook({ id: item.id })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        item.status = 'PENDING'
        await fetchNotebooks(true)
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

  const handleStop = async (item: ConsoleNotebook): Promise<void> => {
    if (item.id === undefined || btnLoading.value[item.id]) return
    if (normalizeStatus(item.status) !== 'RUNNING') return

    try {
      btnLoading.value[item.id] = true
      const res = await stopNotebook({ id: item.id })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        item.status = 'STOPPED'
        await fetchNotebooks(true)
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

  const showSSHInfo = (item: ConsoleNotebook): void => {
    currentNotebook.value = item
    showPassword.value = false
    showSSHDialog.value = true
  }

  const copyText = (text?: string): void => {
    if (!text) return

    void navigator.clipboard
      .writeText(text)
      .then(() => {
        ElMessage.success(t('copied'))
      })
      .catch(() => {
        ElMessage.error(t('copyFailed'))
      })
  }

  const showKeySettings = (): void => {
    void router.push({ name: 'sshkeys' }).catch(() => {
      return router.push('/layout/sshkeys')
    })
  }

  const goToCreate = (): void => {
    void router.push({ name: 'notebookCreate' }).catch(() => {
      return router.push('/layout/notebooks/create')
    })
  }

  const goToDetail = (item: ConsoleNotebook): void => {
    if (item.id === undefined) return

    void router
      .push({ name: 'notebookDetail', query: { id: item.id } })
      .catch(() => {
        return router.push({
          path: '/layout/notebooks/detail',
          query: { id: item.id }
        })
      })
  }

  const getStatusText = (status?: string): string => {
    return t(normalizeStatus(status)) || status || '-'
  }

  const getStatusStyle = (status?: string): string => {
    const statusStyleMap: Record<string, string> = {
      RUNNING: 'console-badge--success',
      STOPPED: 'console-badge--neutral',
      CREATING: 'console-badge--info',
      PENDING: 'console-badge--warning',
      FAILED: 'console-badge--danger',
      SUCCEEDED: 'console-badge--success',
      DELETING: 'console-badge--danger',
      UPDATING: 'console-badge--info',
      UPDATE_FAILED: 'console-badge--danger',
      DELETE_FAILED: 'console-badge--danger',
      UPDATED: 'console-badge--success'
    }

    return statusStyleMap[normalizeStatus(status)] || 'console-badge--neutral'
  }

  onMounted(() => {
    void fetchNotebooks()
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
