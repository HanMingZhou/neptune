import { ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  createInferenceApiKey,
  deleteInferenceApiKey,
  getInferenceApiKeyList
} from '@/api/inference'

export const useInferenceApiKeys = ({ route, t }) => {
  const showApiKeyDialog = ref(false)
  const apiKeys = ref([])
  const newKeyName = ref('')
  const newlyCreatedKey = ref('')

  const fetchApiKeys = async () => {
    try {
      const res = await getInferenceApiKeyList({ serviceId: Number(route.query.id) })

      if (res.code === 0) {
        apiKeys.value = res.data?.list || []
      }
    } catch (error) {
      console.error('Failed to fetch API keys', error)
    }
  }

  const handleCreateApiKey = async () => {
    if (!newKeyName.value.trim()) {
      return
    }

    try {
      const res = await createInferenceApiKey({
        serviceId: Number(route.query.id),
        name: newKeyName.value.trim()
      })

      if (res.code === 0) {
        newlyCreatedKey.value = res.data?.apiKey || ''
        newKeyName.value = ''
        ElMessage.success(t('success'))
        fetchApiKeys()
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (error) {
      ElMessage.error(t('error'))
    }
  }

  const handleDeleteApiKey = async (key) => {
    try {
      await ElMessageBox.confirm(
        t('inference.confirmDeleteKey', { name: key.name }),
        t('tip'),
        { type: 'warning' }
      )

      const res = await deleteInferenceApiKey({ id: key.id })
      if (res.code === 0) {
        ElMessage.success(t('success'))
        fetchApiKeys()
      } else {
        ElMessage.error(res.msg || t('error'))
      }
    } catch (error) {
    }
  }

  watch(showApiKeyDialog, (value) => {
    if (value) {
      newlyCreatedKey.value = ''
      fetchApiKeys()
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
