import { computed, reactive, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { createSSHKey, deleteSSHKey, getSSHKeyList, setDefaultSSHKey } from '@/api/sshkey'

const createDefaultForm = () => ({
  name: '',
  publicKey: ''
})

const isDialogCancel = (error) => error === 'cancel' || error === 'close'

export function useSshKeyManagementPage({ t }) {
  const translate = t || ((key) => key)

  const loading = ref(false)
  const keys = ref([])
  const searchName = ref('')
  const showCreateDialog = ref(false)
  const creating = ref(false)
  const createForm = reactive(createDefaultForm())

  const rules = computed(() => ({
    name: [{ required: true, message: translate('inputName'), trigger: 'blur' }],
    publicKey: [{ required: true, message: translate('publicKeyPlaceholder'), trigger: 'blur' }]
  }))

  const resetCreateForm = () => {
    Object.assign(createForm, createDefaultForm())
  }

  const loadKeys = async (silent = false) => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = await getSSHKeyList({
        page: 1,
        pageSize: 100,
        name: searchName.value || undefined
      })

      if (res.code === 0) {
        keys.value = res.data?.list || []
      } else {
        ElMessage.error(res.msg || translate('error'))
      }
    } catch (error) {
      console.error('Failed to fetch SSH keys:', error)
      ElMessage.error(translate('error'))
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const initialize = async () => {
    await loadKeys()
  }

  const searchKeys = async () => {
    await loadKeys(true)
  }

  const openCreateDialog = () => {
    resetCreateForm()
    showCreateDialog.value = true
  }

  const closeCreateDialog = () => {
    showCreateDialog.value = false
    resetCreateForm()
  }

  const handleCreate = async () => {
    const trimmedPublicKey = createForm.publicKey.trim()
    if (!trimmedPublicKey) {
      ElMessage.warning(translate('publicKeyPlaceholder'))
      return
    }

    creating.value = true

    try {
      const res = await createSSHKey({
        name: createForm.name,
        publicKey: trimmedPublicKey
      })

      if (res.code === 0) {
        ElMessage.success(res.msg || translate('success'))
        closeCreateDialog()
        await loadKeys()
        return
      }

      ElMessage.error(res.msg || translate('error'))
    } catch (error) {
      console.error('Failed to create SSH key:', error)
      ElMessage.error(error?.response?.data?.msg || translate('error'))
    } finally {
      creating.value = false
    }
  }

  const handleDelete = async (key) => {
    try {
      await ElMessageBox.confirm(translate('confirmDeleteSshKey'), translate('tip'), {
        type: 'warning',
        confirmButtonText: translate('delete'),
        cancelButtonText: translate('cancel')
      })

      const res = await deleteSSHKey({ id: key.id })
      if (res.code === 0) {
        ElMessage.success(res.msg || translate('success'))
        await loadKeys()
        return
      }

      ElMessage.error(res.msg || translate('error'))
    } catch (error) {
      if (!isDialogCancel(error)) {
        ElMessage.error(translate('error'))
      }
    }
  }

  const setDefault = async (key) => {
    try {
      const res = await setDefaultSSHKey({ id: key.id })
      if (res.code === 0) {
        ElMessage.success(res.msg || translate('success'))
        await loadKeys(true)
        return
      }

      ElMessage.error(res.msg || translate('error'))
    } catch (error) {
      console.error('Failed to set default SSH key:', error)
      ElMessage.error(translate('error'))
    }
  }

  watch(searchName, (value) => {
    if (!value) {
      loadKeys()
    }
  })

  return {
    closeCreateDialog,
    createForm,
    creating,
    handleCreate,
    handleDelete,
    initialize,
    keys,
    loading,
    loadKeys,
    openCreateDialog,
    rules,
    searchKeys,
    searchName,
    setDefault,
    showCreateDialog
  }
}
