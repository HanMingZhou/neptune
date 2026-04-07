import { ref, type ComputedRef, type Ref } from 'vue'
import type { Terminal as XTermTerminal } from '@xterm/xterm'
import type { FitAddon } from '@xterm/addon-fit'
import { ElMessage } from 'element-plus'
import { getInferenceTerminalWsUrl } from '@/api/inference'
import type { Translator } from '@/types/consoleResource'

const createTerminalTheme = () => ({
  background: '#1e1e1e',
  foreground: '#d4d4d4',
  cursor: '#ffffff',
  selectionBackground: '#264f78',
  black: '#000000',
  red: '#cd3131',
  green: '#0dbc79',
  yellow: '#e5e510',
  blue: '#2472c8',
  magenta: '#bc3fbc',
  cyan: '#11a8cd',
  white: '#e5e5e5',
  brightBlack: '#666666',
  brightRed: '#f14c4c',
  brightGreen: '#23d18b',
  brightYellow: '#f5f543',
  brightBlue: '#3b8eea',
  brightMagenta: '#d670d6',
  brightCyan: '#29b8db',
  brightWhite: '#e5e5e5'
})

interface DisposableLike {
  dispose: () => void
}

interface UseInferenceTerminalOptions {
  canConnectTerminal: ComputedRef<boolean>
  getResourceId: () => number
  getToken: () => string
  t: Translator
  terminalPod: Ref<string>
  terminalRef: Ref<HTMLElement | null>
}

export const useInferenceTerminal = ({
  canConnectTerminal,
  getResourceId,
  getToken,
  t,
  terminalPod,
  terminalRef
}: UseInferenceTerminalOptions) => {
  const terminalConnected = ref(false)

  let terminal: XTermTerminal | null = null
  let fitAddon: FitAddon | null = null
  let resizeObserver: ResizeObserver | null = null
  let ws: WebSocket | null = null
  let dataSubscription: DisposableLike | null = null

  const disconnectTerminal = (): void => {
    if (dataSubscription) {
      dataSubscription.dispose()
      dataSubscription = null
    }

    if (ws) {
      ws.onclose = null
      ws.onerror = null
      ws.onmessage = null
      ws.onopen = null
      ws.close()
      ws = null
    }

    if (terminal) {
      terminal.dispose()
      terminal = null
    }

    if (resizeObserver) {
      resizeObserver.disconnect()
      resizeObserver = null
    }

    terminalConnected.value = false
  }

  const connectTerminal = async (): Promise<void> => {
    if (!canConnectTerminal.value) {
      ElMessage.warning(t('inference.serviceNotReady'))
      return
    }

    const resourceId = getResourceId()
    const token = getToken()
    const terminalContainer = terminalRef.value

    if (!resourceId || !token || !terminalContainer) {
      return
    }

    const { Terminal } = await import('@xterm/xterm')
    await import('@xterm/xterm/css/xterm.css')
    const { FitAddon } = await import('@xterm/addon-fit')

    disconnectTerminal()

    terminal = new Terminal({
      cursorBlink: true,
      fontSize: 14,
      lineHeight: 1.2,
      fontFamily: 'JetBrains Mono, Menlo, Monaco, Consolas, monospace',
      theme: createTerminalTheme()
    })

    fitAddon = new FitAddon()
    terminal.loadAddon(fitAddon)
    terminal.open(terminalContainer)

    resizeObserver = new ResizeObserver(() => {
      fitAddon?.fit()
    })

    resizeObserver.observe(terminalContainer)

    window.setTimeout(() => {
      fitAddon?.fit()
    }, 50)

    const wsUrl = getInferenceTerminalWsUrl(
      resourceId,
      token,
      terminalPod.value || undefined
    )

    ws = new WebSocket(wsUrl)

    ws.onopen = () => {
      terminalConnected.value = true
      fitAddon?.fit()
      terminal?.focus()
    }

    ws.onmessage = (event: MessageEvent<string>) => {
      terminal?.write(event.data)
    }

    ws.onclose = () => {
      terminalConnected.value = false

      if (terminal) {
        terminal.write('\r\n\x1b[31mConnection disconnected.\x1b[0m\r\n')
      }
    }

    ws.onerror = () => {
      ElMessage.error(t('error'))
    }

    dataSubscription = terminal.onData((data) => {
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(data)
      }
    })
  }

  const fitTerminal = (): void => {
    if (fitAddon) {
      fitAddon.fit()
    }
  }

  return {
    connectTerminal,
    disconnectTerminal,
    fitTerminal,
    terminalConnected
  }
}
