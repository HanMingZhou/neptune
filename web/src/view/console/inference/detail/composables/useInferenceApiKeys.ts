import { ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  createInferenceApiKey,
  deleteInferenceApiKey,
  getInferenceApiKeyList
} from '@/api/inference'
import { getErrorMessage } from '@/utils/resourceValidators'
import type { ApiResponse } from '@/utils/request'
import type {
  InferenceApiKey,
  ListData,
  Translator
} from '@/types/consoleResource'

interface CreateInferenceApiKeyResponse {
  apiKey?: string
}

interface UseInferenceApiKeysOptions {
  getServiceId: () => number
  t: Translator
}

export const useInferenceApiKeys = ({
  getServiceId,
  t
}: UseInferenceApiKeysOptions) => {
  const showApiKeyDialog = ref(false)
  const apiKeys = ref<InferenceApiKey[]>([])
  const newKeyName = ref('')
  const newlyCreatedKey = ref('')
  const creatingApiKey = ref(false)
  const apiKeysLoading = ref(false)
  const deletingApiKeyId = ref<number | null>(null)

  const fetchApiKeys = async (): Promise<void> => {
    const serviceId = getServiceId()
    if (!serviceId) {
      return
    }

    apiKeysLoading.value = true
    try {
      const res = (await getInferenceApiKeyList({ serviceId })) as ApiResponse<
        ListData<InferenceApiKey>
      >

      if (res.code === 0) {
        apiKeys.value = res.data?.list || []
      }
    } catch (error) {
      console.error('Failed to fetch API keys', error)
      ElMessage.error(getErrorMessage(error, t('error')))
    } finally {
      apiKeysLoading.value = false
    }
  }

  const handleCreateApiKey = async (): Promise<void> => {
    if (!newKeyName.value.trim() || creatingApiKey.value) {
      return
    }

    const serviceId = getServiceId()
    if (!serviceId) {
      return
    }

    creatingApiKey.value = true
    try {
      const res = (await createInferenceApiKey({
        serviceId,
        name: newKeyName.value.trim()
      })) as ApiResponse<CreateInferenceApiKeyResponse>

      if (res.code === 0) {
        newlyCreatedKey.value = res.data?.apiKey || ''
        newKeyName.value = ''
        ElMessage.success(t('success'))
        void fetchApiKeys()
      } else {
        ElMessage.error(getErrorMessage(res, t('error')))
      }
    } catch (error) {
      ElMessage.error(getErrorMessage(error, t('error')))
    } finally {
      creatingApiKey.value = false
    }
  }

  const handleDeleteApiKey = async (key: InferenceApiKey): Promise<void> => {
    try {
      await ElMessageBox.confirm(
        t('inference.confirmDeleteKey', { name: key.name || '-' }),
        t('tip'),
        { type: 'warning' }
      )

      deletingApiKeyId.value = Number(key.id)
      const res = await deleteInferenceApiKey({ id: key.id })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        void fetchApiKeys()
      } else {
        ElMessage.error(getErrorMessage(res, t('error')))
      }
    } catch (_error) {
    } finally {
      deletingApiKeyId.value = null
    }
  }

  watch(showApiKeyDialog, (value) => {
    if (value) {
      newlyCreatedKey.value = ''
      void fetchApiKeys()
    } else {
      newKeyName.value = ''
    }
  })

  return {
    apiKeys,
    apiKeysLoading,
    creatingApiKey,
    deletingApiKeyId,
    handleCreateApiKey,
    handleDeleteApiKey,
    newKeyName,
    newlyCreatedKey,
    showApiKeyDialog
  }
}
