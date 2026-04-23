import { computed, reactive, ref } from 'vue'
import type { FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  createCluster,
  deleteCluster,
  getClusterList,
  updateCluster
} from '@/api/cluster'
import type { Translator } from '@/types/consoleResource'
import type {
  CmsClusterForm,
  CmsClusterListData,
  CmsClusterRow
} from '@/types/superAdmin'
import type { ApiResponse } from '@/utils/request'
import { getErrorMessage } from '@/utils/resourceValidators'

interface UseClusterManagementPageOptions {
  t?: Translator
}

const createDefaultForm = (): CmsClusterForm => ({
  id: null,
  name: '',
  area: '',
  description: '',
  kubeconfig: '',
  apiServer: '',
  status: 1,
  harborAddr: '',
  storageClass: ''
})

const copyText = async (text: string): Promise<boolean> => {
  if (typeof navigator === 'undefined' || !navigator.clipboard?.writeText) {
    return false
  }

  await navigator.clipboard.writeText(text)
  return true
}

const isDialogCancel = (error: unknown): error is 'cancel' | 'close' =>
  error === 'cancel' || error === 'close'

export function useClusterManagementPage({
  t
}: UseClusterManagementPageOptions = {}) {
  const translate: Translator = t || ((key: string) => key)

  const clusters = ref<CmsClusterRow[]>([])
  const filterKeyword = ref('')
  const filterStatus = ref<number | undefined>(undefined)
  const form = reactive<CmsClusterForm>(createDefaultForm())
  const isEdit = ref(false)
  const loading = ref(false)
  const page = ref(1)
  const pageSize = ref(15)
  const showDialog = ref(false)
  const showKubeConfigDialog = ref(false)
  const submitting = ref(false)
  const total = ref(0)
  const viewingKubeConfig = ref('')

  const dialogTitle = computed(() =>
    isEdit.value ? translate('clusterEdit') : translate('clusterAdd')
  )
  const formRules = computed<FormRules<CmsClusterForm>>(() => ({
    name: [{ required: true, message: translate('inputName'), trigger: 'blur' }]
  }))

  const resetForm = (): void => {
    Object.assign(form, createDefaultForm())
  }

  const fetchClusters = async (silent = false): Promise<void> => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = (await getClusterList({
        keyword: filterKeyword.value || undefined,
        page: page.value,
        pageSize: pageSize.value,
        status: filterStatus.value
      })) as ApiResponse<CmsClusterListData>

      if (res.code === 0) {
        clusters.value = res.data?.list ?? []
        total.value = res.data?.total ?? 0
        page.value = res.data?.page ?? page.value
        pageSize.value = res.data?.pageSize ?? pageSize.value
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error: unknown) {
      console.error('Failed to fetch clusters:', error)
      ElMessage.error(getErrorMessage(error, translate('failed')))
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const initialize = async (): Promise<void> => {
    await fetchClusters()
  }

  const handleSearch = (): void => {
    page.value = 1
    void fetchClusters()
  }

  const handleResetFilters = (): void => {
    page.value = 1
    filterKeyword.value = ''
    filterStatus.value = undefined
    void fetchClusters()
  }

  const handleCurrentChange = (value: number): void => {
    page.value = value
    void fetchClusters()
  }

  const handleSizeChange = (value: number): void => {
    page.value = 1
    pageSize.value = value
    void fetchClusters()
  }

  const openCreateDialog = (): void => {
    resetForm()
    isEdit.value = false
    showDialog.value = true
  }

  const openEditDialog = (row: CmsClusterRow): void => {
    resetForm()
    isEdit.value = true
    Object.assign(form, {
      id: row.id,
      name: row.name,
      area: row.area || '',
      description: row.description || '',
      kubeconfig: row.kubeConfig || '',
      apiServer: row.apiServer || '',
      status: row.status ?? 1,
      harborAddr: row.harborAddr || '',
      storageClass: row.storageClass || ''
    })
    showDialog.value = true
  }

  const closeDialog = (): void => {
    showDialog.value = false
    resetForm()
  }

  const handleSubmit = async (): Promise<void> => {
    submitting.value = true

    try {
      const payload: Partial<CmsClusterForm> &
        Pick<CmsClusterForm, 'name' | 'status'> = { ...form }
      if (!payload.kubeconfig) {
        delete payload.kubeconfig
      }

      const api = isEdit.value ? updateCluster : createCluster
      const res = await api(payload)

      if (res.code === 0) {
        ElMessage.success(res.msg || translate('success'))
        closeDialog()
        if (!isEdit.value) {
          page.value = 1
        }
        await fetchClusters()
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error: unknown) {
      console.error('Failed to submit cluster form:', error)
      ElMessage.error(getErrorMessage(error, translate('failed')))
    } finally {
      submitting.value = false
    }
  }

  const handleDelete = async (row: CmsClusterRow): Promise<void> => {
    try {
      await ElMessageBox.confirm(
        translate('confirmDelete', { name: row.name }),
        translate('tip'),
        {
          confirmButtonText: translate('confirm'),
          cancelButtonText: translate('cancel'),
          type: 'warning'
        }
      )

      const res = await deleteCluster({ id: row.id })
      if (res.code === 0) {
        ElMessage.success(res.msg || translate('success'))
        if (clusters.value.length === 1 && page.value > 1) {
          page.value -= 1
        }
        await fetchClusters()
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error: unknown) {
      if (!isDialogCancel(error)) {
        ElMessage.error(getErrorMessage(error, translate('failed')))
      }
    }
  }

  const viewKubeConfig = (row: CmsClusterRow): void => {
    viewingKubeConfig.value = row.kubeConfig || ''
    showKubeConfigDialog.value = true
  }

  const closeKubeConfigDialog = (): void => {
    showKubeConfigDialog.value = false
    viewingKubeConfig.value = ''
  }

  const copyKubeConfig = async (): Promise<void> => {
    try {
      const copied = await copyText(viewingKubeConfig.value)
      if (copied) {
        ElMessage.success(translate('copied'))
      } else {
        ElMessage.error(translate('copyFailed'))
      }
    } catch {
      ElMessage.error(translate('copyFailed'))
    }
  }

  return {
    closeDialog,
    closeKubeConfigDialog,
    clusters,
    copyKubeConfig,
    dialogTitle,
    fetchClusters,
    filterKeyword,
    filterStatus,
    form,
    formRules,
    handleCurrentChange,
    handleDelete,
    handleResetFilters,
    handleSearch,
    handleSizeChange,
    handleSubmit,
    initialize,
    isEdit,
    loading,
    openCreateDialog,
    openEditDialog,
    page,
    pageSize,
    showDialog,
    showKubeConfigDialog,
    submitting,
    total,
    viewKubeConfig,
    viewingKubeConfig
  }
}

