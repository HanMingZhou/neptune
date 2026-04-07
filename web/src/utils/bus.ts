import mitt from 'mitt'

export interface KeepAliveHistoryItem {
  name: string
}

interface BusEvents {
  [key: string]: unknown
  [key: symbol]: unknown
  setKeepAlive: KeepAliveHistoryItem[]
  closeThisPage: undefined
  showLoading: undefined
  closeLoading: undefined
  'show-error': {
    code: string
    message: string
    fn: () => void
  }
}

export const emitter = mitt<BusEvents>()
