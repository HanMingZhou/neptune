import type { Terminal } from '@xterm/xterm'

interface TerminalClipboardOptions {
  notifyCopyError?: () => void
  notifyPasteError?: () => void
  sendData: (data: string) => void
}

export interface TerminalClipboardBridge {
  dispose: () => void
}

const getViewportElement = (container: HTMLElement): HTMLElement | null =>
  container.querySelector<HTMLElement>('.xterm-viewport')

const copyWithExecCommand = (text: string): boolean => {
  if (typeof document === 'undefined') {
    return false
  }

  const textarea = document.createElement('textarea')
  textarea.value = text
  textarea.setAttribute('readonly', 'true')
  textarea.style.position = 'fixed'
  textarea.style.opacity = '0'
  textarea.style.pointerEvents = 'none'
  textarea.style.left = '-9999px'
  textarea.style.top = '0'
  document.body.appendChild(textarea)
  textarea.focus()
  textarea.select()

  let copied = false
  try {
    copied = document.execCommand('copy')
  } catch (_error) {
    copied = false
  }

  document.body.removeChild(textarea)
  return copied
}

const writeClipboardText = async (text: string): Promise<boolean> => {
  if (!text) {
    return true
  }

  if (typeof navigator !== 'undefined' && navigator.clipboard?.writeText) {
    try {
      await navigator.clipboard.writeText(text)
      return true
    } catch (_error) {}
  }

  return copyWithExecCommand(text)
}

const readClipboardText = async (): Promise<string | null> => {
  if (typeof navigator === 'undefined' || !navigator.clipboard?.readText) {
    return null
  }

  try {
    return await navigator.clipboard.readText()
  } catch (_error) {
    return null
  }
}

const isCopyShortcut = (
  event: KeyboardEvent,
  hasSelection: boolean
): boolean => {
  if (event.altKey) {
    return false
  }

  const key = event.key.toLowerCase()
  if (key !== 'c') {
    return false
  }

  if (event.metaKey && !event.ctrlKey) {
    return hasSelection
  }

  if (!event.ctrlKey) {
    return false
  }

  if (event.shiftKey) {
    return hasSelection
  }

  return hasSelection
}

const isPasteShortcut = (event: KeyboardEvent): boolean => {
  if (event.altKey) {
    return false
  }

  const key = event.key.toLowerCase()
  if (key === 'insert' && event.shiftKey) {
    return true
  }

  if (key !== 'v') {
    return false
  }

  if (event.metaKey && !event.ctrlKey) {
    return true
  }

  return event.ctrlKey
}

export const createTerminalClipboardBridge = (
  terminal: Terminal,
  container: HTMLElement,
  options: TerminalClipboardOptions
): TerminalClipboardBridge => {
  const copySelection = async (): Promise<void> => {
    const selection = terminal.getSelection()
    if (!selection) {
      return
    }

    const copied = await writeClipboardText(selection)
    if (!copied) {
      options.notifyCopyError?.()
      return
    }

    terminal.focus()
  }

  const pasteText = (text: string): void => {
    if (!text) {
      return
    }

    options.sendData(text)
    terminal.focus()
  }

  const handleCopyEvent = (event: ClipboardEvent): void => {
    const selection = terminal.getSelection()
    if (!selection) {
      return
    }

    event.preventDefault()
    if (event.clipboardData) {
      event.clipboardData.setData('text/plain', selection)
      terminal.focus()
      return
    }

    void copySelection()
  }

  const handlePasteEvent = (event: ClipboardEvent): void => {
    const text = event.clipboardData?.getData('text/plain')
    if (typeof text !== 'string') {
      return
    }

    event.preventDefault()
    pasteText(text)
  }

  const handleWheelEvent = (event: WheelEvent): void => {
    const viewport = getViewportElement(container)
    if (!viewport) {
      return
    }

    event.stopPropagation()

    if (!event.cancelable || event.deltaY === 0) {
      return
    }

    const maxScrollTop = viewport.scrollHeight - viewport.clientHeight
    if (maxScrollTop <= 0) {
      event.preventDefault()
      return
    }

    const scrollTop = viewport.scrollTop
    const isScrollingUpPastTop = event.deltaY < 0 && scrollTop <= 0
    const isScrollingDownPastBottom =
      event.deltaY > 0 && scrollTop >= maxScrollTop - 1

    if (isScrollingUpPastTop || isScrollingDownPastBottom) {
      event.preventDefault()
    }
  }

  container.addEventListener('copy', handleCopyEvent, true)
  container.addEventListener('paste', handlePasteEvent, true)
  container.addEventListener('wheel', handleWheelEvent, {
    capture: true,
    passive: false
  })

  terminal.attachCustomKeyEventHandler((event: KeyboardEvent) => {
    if (event.type !== 'keydown') {
      return true
    }

    const hasSelection = terminal.hasSelection()

    if (isCopyShortcut(event, hasSelection)) {
      event.preventDefault()
      void copySelection()
      return false
    }

    if (isPasteShortcut(event)) {
      if (typeof navigator === 'undefined' || !navigator.clipboard?.readText) {
        return true
      }

      event.preventDefault()
      void readClipboardText().then((text) => {
        if (typeof text !== 'string') {
          options.notifyPasteError?.()
          return
        }

        pasteText(text)
      })
      return false
    }

    return true
  })

  return {
    dispose: (): void => {
      container.removeEventListener('copy', handleCopyEvent, true)
      container.removeEventListener('paste', handlePasteEvent, true)
      container.removeEventListener('wheel', handleWheelEvent, true)
    }
  }
}
