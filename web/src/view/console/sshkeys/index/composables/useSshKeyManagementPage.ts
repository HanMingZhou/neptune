import { computed, reactive, ref, watch } from 'vue'
import type { FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  createSSHKey,
  deleteSSHKey,
  getSSHKeyList,
  setDefaultSSHKey
} from '@/api/sshkey'
import type { Translator } from '@/types/consoleResource'
import type { ApiResponse } from '@/utils/request'
import type { SshKeyCreateForm, SshKeyListItem } from '@/types/sshkey'
import { getErrorMessage } from '@/utils/resourceValidators'

interface UseSshKeyManagementPageOptions {
  t?: Translator
}

const createDefaultForm = (): SshKeyCreateForm => ({
  name: '',
  publicKey: ''
})

const isDialogCancel = (error: unknown): error is 'cancel' | 'close' =>
  error === 'cancel' || error === 'close'

export function useSshKeyManagementPage({
  t
}: UseSshKeyManagementPageOptions = {}) {
  const translate: Translator = t || ((key: string) => key)

  const page = ref(1)
  const pageSize = ref(15)
  const loading = ref(false)
  const keys = ref<SshKeyListItem[]>([])
  const total = ref(0)
  const searchName = ref('')
  const showCreateDialog = ref(false)
  const creating = ref(false)
  const createForm = reactive<SshKeyCreateForm>(createDefaultForm())

  const rules = computed<FormRules<SshKeyCreateForm>>(() => ({
    name: [
      { required: true, message: translate('inputName'), trigger: 'blur' }
    ],
    publicKey: [
      {
        required: true,
        message: translate('publicKeyPlaceholder'),
        trigger: 'blur'
      }
    ]
  }))

  const resetCreateForm = (): void => {
    Object.assign(createForm, createDefaultForm())
  }

  const loadKeys = async (silent = false): Promise<void> => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = (await getSSHKeyList({
        page: page.value,
        pageSize: pageSize.value,
        name: searchName.value || undefined
      })) as ApiResponse<{ list?: SshKeyListItem[]; total?: number }>

      if (res.code === 0) {
        keys.value = res.data?.list ?? []
        total.value = res.data?.total ?? 0
      } else {
        ElMessage.error(res.msg || translate('error'))
      }
    } catch (error: unknown) {
      console.error('Failed to fetch SSH keys:', error)
      ElMessage.error(getErrorMessage(error, translate('error')))
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const initialize = async (): Promise<void> => {
    await loadKeys()
  }

  const searchKeys = async (): Promise<void> => {
    page.value = 1
    await loadKeys(true)
  }

  const handleCurrentChange = (value: number): void => {
    page.value = value
    void loadKeys()
  }

  const handleSizeChange = (value: number): void => {
    page.value = 1
    pageSize.value = value
    void loadKeys()
  }

  const openCreateDialog = (): void => {
    resetCreateForm()
    showCreateDialog.value = true
  }

  const closeCreateDialog = (): void => {
    showCreateDialog.value = false
    resetCreateForm()
  }

  const handleCreate = async (): Promise<void> => {
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
        page.value = 1
        await loadKeys()
        return
      }

      ElMessage.error(res.msg || translate('error'))
    } catch (error: unknown) {
      console.error('Failed to create SSH key:', error)
      ElMessage.error(getErrorMessage(error, translate('error')))
    } finally {
      creating.value = false
    }
  }

  const handleDelete = async (key: SshKeyListItem): Promise<void> => {
    try {
      await ElMessageBox.confirm(
        translate('confirmDeleteSshKey'),
        translate('tip'),
        {
          type: 'warning',
          confirmButtonText: translate('delete'),
          cancelButtonText: translate('cancel')
        }
      )

      const res = await deleteSSHKey({ id: key.id })
      if (res.code === 0) {
        ElMessage.success(res.msg || translate('success'))
        if (keys.value.length === 1 && page.value > 1) {
          page.value -= 1
        }
        await loadKeys()
        return
      }

      ElMessage.error(res.msg || translate('error'))
    } catch (error) {
      if (!isDialogCancel(error)) {
        ElMessage.error(getErrorMessage(error, translate('error')))
      }
    }
  }

  const setDefault = async (key: SshKeyListItem): Promise<void> => {
    try {
      const res = await setDefaultSSHKey({ id: key.id })
      if (res.code === 0) {
        ElMessage.success(res.msg || translate('success'))
        await loadKeys(true)
        return
      }

      ElMessage.error(res.msg || translate('error'))
    } catch (error: unknown) {
      console.error('Failed to set default SSH key:', error)
      ElMessage.error(getErrorMessage(error, translate('error')))
    }
  }

  watch(searchName, (value) => {
    if (!value) {
      page.value = 1
      void loadKeys()
    }
  })

  return {
    closeCreateDialog,
    createForm,
    creating,
    handleCurrentChange,
    handleCreate,
    handleDelete,
    handleSizeChange,
    initialize,
    keys,
    loading,
    loadKeys,
    openCreateDialog,
    page,
    pageSize,
    rules,
    searchKeys,
    searchName,
    setDefault,
    showCreateDialog,
    total
  }
}

