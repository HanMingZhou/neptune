import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { getTrainingTerminalWsUrl } from '@/api/training'

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

export const useTrainingTerminal = ({ canConnectTerminal, route, statusText, t, terminalOptions, terminalRef, userStore }) => {
  const terminalConnected = ref(false)

  let terminal = null
  let fitAddon = null
  let resizeObserver = null
  let ws = null

  const disconnectTerminal = () => {
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

  const connectTerminal = async () => {
    if (!canConnectTerminal.value) {
      ElMessage.warning(`${t('instanceStatus')}: ${statusText.value}`)
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
    terminal.open(terminalRef.value)

    resizeObserver = new ResizeObserver(() => {
      fitAddon.fit()
    })

    if (terminalRef.value) {
      resizeObserver.observe(terminalRef.value)
    }

    setTimeout(() => {
      fitAddon.fit()
    }, 50)

    const wsUrl = getTrainingTerminalWsUrl(
      Number(route.query.id),
      userStore.token,
      terminalOptions.taskName || undefined
    )

    ws = new WebSocket(wsUrl)

    ws.onopen = () => {
      terminalConnected.value = true
      fitAddon.fit()
      terminal.focus()
    }

    ws.onmessage = (event) => {
      terminal.write(event.data)
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

    terminal.onData((data) => {
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(data)
      }
    })
  }

  const fitTerminal = () => {
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
