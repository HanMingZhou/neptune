import { computed, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { createCluster, deleteCluster, getClusterList, updateCluster } from '@/api/cluster'

const createDefaultForm = () => ({
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

const copyText = async (text) => {
  if (typeof navigator === 'undefined' || !navigator.clipboard?.writeText) {
    return false
  }

  await navigator.clipboard.writeText(text)
  return true
}

const isDialogCancel = (error) => error === 'cancel' || error === 'close'

export function useClusterManagementPage({ t }) {
  const translate = t || ((key) => key)

  const clusters = ref([])
  const filterKeyword = ref('')
  const filterStatus = ref(undefined)
  const form = reactive(createDefaultForm())
  const isEdit = ref(false)
  const loading = ref(false)
  const showDialog = ref(false)
  const showKubeConfigDialog = ref(false)
  const submitting = ref(false)
  const viewingKubeConfig = ref('')

  const dialogTitle = computed(() => (isEdit.value ? translate('clusterEdit') : translate('clusterAdd')))
  const formRules = computed(() => ({
    name: [{ required: true, message: translate('inputName'), trigger: 'blur' }]
  }))

  const resetForm = () => {
    Object.assign(form, createDefaultForm())
  }

  const fetchClusters = async (silent = false) => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = await getClusterList({
        keyword: filterKeyword.value || undefined,
        status: filterStatus.value
      })

      if (res.code === 0) {
        clusters.value = res.data?.list || []
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error) {
      console.error('Failed to fetch clusters:', error)
      ElMessage.error(translate('failed'))
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const initialize = async () => {
    await fetchClusters()
  }

  const handleResetFilters = () => {
    filterKeyword.value = ''
    filterStatus.value = undefined
    fetchClusters()
  }

  const openCreateDialog = () => {
    resetForm()
    isEdit.value = false
    showDialog.value = true
  }

  const openEditDialog = (row) => {
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

  const closeDialog = () => {
    showDialog.value = false
    resetForm()
  }

  const handleSubmit = async () => {
    submitting.value = true

    try {
      const payload = { ...form }
      if (!payload.kubeconfig) {
        delete payload.kubeconfig
      }

      const api = isEdit.value ? updateCluster : createCluster
      const res = await api(payload)

      if (res.code === 0) {
        ElMessage.success(res.msg || translate('success'))
        closeDialog()
        await fetchClusters()
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error) {
      console.error('Failed to submit cluster form:', error)
      ElMessage.error(translate('failed'))
    } finally {
      submitting.value = false
    }
  }

  const handleDelete = async (row) => {
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
        await fetchClusters()
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error) {
      if (!isDialogCancel(error)) {
        ElMessage.error(translate('failed'))
      }
    }
  }

  const viewKubeConfig = (row) => {
    viewingKubeConfig.value = row.kubeConfig || ''
    showKubeConfigDialog.value = true
  }

  const closeKubeConfigDialog = () => {
    showKubeConfigDialog.value = false
    viewingKubeConfig.value = ''
  }

  const copyKubeConfig = async () => {
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
    handleDelete,
    handleResetFilters,
    handleSubmit,
    initialize,
    isEdit,
    loading,
    openCreateDialog,
    openEditDialog,
    showDialog,
    showKubeConfigDialog,
    submitting,
    viewKubeConfig,
    viewingKubeConfig
  }
}
