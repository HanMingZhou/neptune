import type { Terminal } from '@xterm/xterm'

const FLUSH_INTERVAL_MS = 16
const MAX_PENDING_CHARS = 512 * 1024
const WRITE_BATCH_CHARS = 16 * 1024
const OVERFLOW_WARNING =
  '\r\n\x1b[33m[neptune] Terminal output is too large or too fast. Some content may be omitted. Try less/head/sed -n for large files.\x1b[0m\r\n'

interface TerminalOutputBufferOptions {
  onOverflow?: () => void
}

export interface TerminalOutputBuffer {
  dispose: () => void
  flushNow: () => void
  push: (data: string) => void
}

export const createTerminalOutputBuffer = (
  terminal: Terminal,
  options: TerminalOutputBufferOptions = {}
): TerminalOutputBuffer => {
  let pendingChunks: string[] = []
  let pendingChars = 0
  let flushTimer: number | null = null
  let writing = false
  let disposed = false
  let overflowNotified = false

  const clearFlushTimer = (): void => {
    if (flushTimer !== null) {
      window.clearTimeout(flushTimer)
      flushTimer = null
    }
  }

  const resetOverflowStateIfIdle = (): void => {
    if (!writing && pendingChars === 0) {
      overflowNotified = false
    }
  }

  const enqueue = (chunk: string): void => {
    if (!chunk) {
      return
    }

    pendingChunks.push(chunk)
    pendingChars += chunk.length
  }

  const notifyOverflow = (): void => {
    if (overflowNotified) {
      return
    }

    overflowNotified = true
    enqueue(OVERFLOW_WARNING)
    options.onOverflow?.()
  }

  const appendChunk = (chunk: string): void => {
    if (!chunk) {
      return
    }

    const remainingCapacity = MAX_PENDING_CHARS - pendingChars
    if (remainingCapacity <= 0) {
      notifyOverflow()
      return
    }

    if (chunk.length <= remainingCapacity) {
      enqueue(chunk)
      return
    }

    enqueue(chunk.slice(0, remainingCapacity))
    notifyOverflow()
  }

  const takeBatch = (): string => {
    let batch = ''

    while (pendingChunks.length > 0 && batch.length < WRITE_BATCH_CHARS) {
      const nextChunk = pendingChunks[0]
      const remaining = WRITE_BATCH_CHARS - batch.length

      if (nextChunk.length <= remaining) {
        batch += nextChunk
        pendingChunks.shift()
        pendingChars -= nextChunk.length
        continue
      }

      batch += nextChunk.slice(0, remaining)
      pendingChunks[0] = nextChunk.slice(remaining)
      pendingChars -= remaining
    }

    return batch
  }

  const drain = (): void => {
    if (disposed || writing || pendingChars === 0) {
      resetOverflowStateIfIdle()
      return
    }

    const batch = takeBatch()
    if (!batch) {
      resetOverflowStateIfIdle()
      return
    }

    writing = true
    terminal.write(batch, () => {
      writing = false
      if (disposed) {
        return
      }

      if (pendingChars > 0) {
        scheduleDrain()
        return
      }

      resetOverflowStateIfIdle()
    })
  }

  const scheduleDrain = (): void => {
    if (disposed || writing || flushTimer !== null) {
      return
    }

    flushTimer = window.setTimeout(() => {
      flushTimer = null
      drain()
    }, FLUSH_INTERVAL_MS)
  }

  return {
    push: (data: string): void => {
      if (disposed || !data) {
        return
      }

      appendChunk(data)
      scheduleDrain()
    },
    flushNow: (): void => {
      if (disposed) {
        return
      }

      clearFlushTimer()
      drain()
    },
    dispose: (): void => {
      disposed = true
      clearFlushTimer()
      pendingChunks = []
      pendingChars = 0
      overflowNotified = false
    }
  }
}
