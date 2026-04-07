import { nextTick, ref, type Ref } from 'vue'
import { ElMessage } from 'element-plus'
import { getInferenceLogStreamWsUrl } from '@/api/inference'
import type { Translator } from '@/types/consoleResource'
import { processInferenceLogData } from './inferenceDetailUtils'

interface UseInferenceLogStreamOptions {
  getResourceId: () => number
  getToken: () => string
  logsRef: Ref<HTMLElement | null>
  selectedPod: Ref<string>
  t: Translator
}

export const useInferenceLogStream = ({
  getResourceId,
  getToken,
  logsRef,
  selectedPod,
  t
}: UseInferenceLogStreamOptions) => {
  const logs = ref('')
  const logsLoading = ref(false)
  const logsConnected = ref(false)

  let logsWs: WebSocket | null = null

  const connectLogStream = (): void => {
    if (logsConnected.value || logsLoading.value || logsWs) {
      return
    }

    const resourceId = getResourceId()
    const token = getToken()
    if (!resourceId || !token) {
      return
    }

    logsLoading.value = true
    logs.value = ''

    const wsUrl = getInferenceLogStreamWsUrl(
      resourceId,
      token,
      selectedPod.value || undefined
    )

    logsWs = new WebSocket(wsUrl)

    logsWs.onopen = () => {
      logsConnected.value = true
      logsLoading.value = false
    }

    logsWs.onmessage = async (event: MessageEvent<string>) => {
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

  const disconnectLogStream = (): void => {
    if (logsWs) {
      logsWs.close()
      logsWs = null
    }

    logsConnected.value = false
    logsLoading.value = false
  }

  const clearLogs = (): void => {
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
