import { nextTick, ref, type Ref } from 'vue'
import { ElMessage } from 'element-plus'
import { getTrainingLogStreamWsUrl } from '@/api/training'
import type { Translator } from '@/types/consoleResource'
import { processTrainingLogData } from './trainingDetailUtils'

interface TrainingLogOptions {
  taskName: string
  podIndex: number
}

interface UseTrainingLogStreamOptions {
  getResourceId: () => number
  getToken: () => string
  logOptions: TrainingLogOptions
  logsRef: Ref<HTMLElement | null>
  t: Translator
}

export const useTrainingLogStream = ({
  getResourceId,
  getToken,
  logOptions,
  logsRef,
  t
}: UseTrainingLogStreamOptions) => {
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

    const wsUrl = getTrainingLogStreamWsUrl(
      resourceId,
      token,
      logOptions.taskName || undefined,
      logOptions.podIndex
    )

    logsWs = new WebSocket(wsUrl)

    logsWs.onopen = () => {
      logsConnected.value = true
      logsLoading.value = false
    }

    logsWs.onmessage = async (event: MessageEvent<string>) => {
      logs.value += processTrainingLogData(event.data)
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
      console.error('Training log websocket error', error)
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
