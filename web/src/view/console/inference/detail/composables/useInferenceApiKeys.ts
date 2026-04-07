import { ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  createInferenceApiKey,
  deleteInferenceApiKey,
  getInferenceApiKeyList
} from '@/api/inference'
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

  const fetchApiKeys = async (): Promise<void> => {
    const serviceId = getServiceId()
    if (!serviceId) {
      return
    }

    try {
      const res = (await getInferenceApiKeyList({ serviceId })) as ApiResponse<
        ListData<InferenceApiKey>
      >

      if (res.code === 0) {
        apiKeys.value = res.data?.list || []
      }
    } catch (_error) {
      console.error('Failed to fetch API keys', _error)
    }
  }

  const handleCreateApiKey = async (): Promise<void> => {
    if (!newKeyName.value.trim()) {
      return
    }

    const serviceId = getServiceId()
    if (!serviceId) {
      return
    }

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
        ElMessage.error(res.msg || t('error'))
      }
    } catch (_error) {
      ElMessage.error(t('error'))
    }
  }

  const handleDeleteApiKey = async (key: InferenceApiKey): Promise<void> => {
    try {
      await ElMessageBox.confirm(
        t('inference.confirmDeleteKey', { name: key.name || '-' }),
        t('tip'),
        { type: 'warning' }
      )

      const res = await deleteInferenceApiKey({ id: key.id })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        void fetchApiKeys()
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (_error) {}
  }

  watch(showApiKeyDialog, (value) => {
    if (value) {
      newlyCreatedKey.value = ''
      void fetchApiKeys()
    }
  })

  return {
    apiKeys,
    handleCreateApiKey,
    handleDeleteApiKey,
    newKeyName,
    newlyCreatedKey,
    showApiKeyDialog
  }
}
