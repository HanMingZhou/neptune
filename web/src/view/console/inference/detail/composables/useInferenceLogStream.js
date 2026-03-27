import { nextTick, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { getInferenceLogStreamWsUrl } from '@/api/inference'
import { processInferenceLogData } from './inferenceDetailUtils'

export const useInferenceLogStream = ({ logsRef, route, selectedPod, t, userStore }) => {
  const logs = ref('')
  const logsLoading = ref(false)
  const logsConnected = ref(false)

  let logsWs = null

  const connectLogStream = () => {
    if (logsConnected.value || logsLoading.value || logsWs) {
      return
    }

    logsLoading.value = true
    logs.value = ''

    const wsUrl = getInferenceLogStreamWsUrl(
      Number(route.query.id),
      userStore.token,
      selectedPod.value || undefined
    )

    logsWs = new WebSocket(wsUrl)

    logsWs.onopen = () => {
      logsConnected.value = true
      logsLoading.value = false
    }

    logsWs.onmessage = async (event) => {
      logs.value += processInferenceLogData(event.data)
      await nextTick()

      if (logsRef.value) {
        logsRef.value.scrollTop = logsRef.value.scrollHeight
      }
    }

    logsWs.onclose = () => {
      logsWs = null
      logsConnected.value = false
      logsLoading.value = false
    }

    logsWs.onerror = (error) => {
      console.error('WebSocket error', error)
      logsWs = null
      logsConnected.value = false
      logsLoading.value = false
      ElMessage.error(t('error'))
    }
  }

  const disconnectLogStream = () => {
    if (logsWs) {
      logsWs.close()
      logsWs = null
    }

    logsConnected.value = false
    logsLoading.value = false
  }

  const clearLogs = () => {
    logs.value = ''
  }

  return {
    clearLogs,
    connectLogStream,
    disconnectLogStream,
    logs,
    logsConnected,
    logsLoading
  }
}
