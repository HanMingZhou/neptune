import type { Translator } from '@/types/consoleResource'

type NullableNumber = number | null | undefined

export interface VGpuSpecLike {
  vGpuNumber?: NullableNumber
  vGpuCount?: NullableNumber
  vGpuMemory?: NullableNumber
  vGpuCores?: NullableNumber
}

export interface GpuSpecLike {
  gpuCount?: NullableNumber
  gpuMemory?: NullableNumber
}

interface FormatVGpuSpecOptions {
  concise?: boolean
  detailed?: boolean
  separator?: string
  t?: Translator
}

export interface VGpuFieldEntry {
  key: 'count' | 'memory' | 'cores'
  label: string
  value: string
}

export const getVGpuNumber = (resource?: VGpuSpecLike | null): number =>
  resource?.vGpuNumber || resource?.vGpuCount || 0

export const hasVGpuSpec = (resource?: VGpuSpecLike | null): boolean =>
  getVGpuNumber(resource) > 0 ||
  (resource?.vGpuMemory || 0) > 0 ||
  (resource?.vGpuCores || 0) > 0

export const hasGpuSpec = (resource?: GpuSpecLike | null): boolean =>
  (resource?.gpuCount || 0) > 0 || (resource?.gpuMemory || 0) > 0

const trimTrailingZeros = (value: string): string =>
  value.replace(/\.0+$|(\.\d*[1-9])0+$/, '$1')

const getTranslatedLabel = (
  t: Translator | undefined,
  key: string,
  fallback: string
): string => {
  const translated = typeof t === 'function' ? t(key) : key
  return translated && translated !== key ? translated : fallback
}

const getTranslatedLabelWithFallback = (
  t: Translator | undefined,
  primaryKey: string,
  fallbackKey: string,
  fallback: string
): string => {
  const primaryLabel = getTranslatedLabel(t, primaryKey, '')
  if (primaryLabel) {
    return primaryLabel
  }

  return getTranslatedLabel(t, fallbackKey, fallback)
}

const getCountShortLabel = (t?: Translator): string =>
  getTranslatedLabelWithFallback(t, 'gpuCountShort', 'vGpuCountShort', 'Count')

const getMemoryShortLabel = (t?: Translator): string =>
  getTranslatedLabelWithFallback(t, 'gpuMemoryShort', 'vGpuMemoryShort', 'Memory')

export const formatVGpuMemory = (
  memory?: NullableNumber,
  options: Pick<FormatVGpuSpecOptions, 'concise' | 'detailed'> = {}
): string => {
  const value = memory || 0
  if (value <= 0) return ''

  const { concise = false, detailed = false } = options
  if (value < 1000) {
    return detailed ? `${value}MB 显存` : `${value}MB`
  }

  const gbValue = trimTrailingZeros((value / 1000).toFixed(value % 1000 === 0 ? 0 : 1))
  if (detailed) {
    return `${gbValue}GB 显存`
  }
  if (concise) {
    return `${gbValue}GB`
  }
  return `${gbValue}GB`
}

export const formatGpuMemory = (memory?: NullableNumber): string => {
  const value = memory || 0
  if (value <= 0) return ''

  return `${trimTrailingZeros(String(value))}GB`
}

export const formatVGpuSpec = (
  resource?: VGpuSpecLike | null,
  options: FormatVGpuSpecOptions = {}
): string => {
  const parts: string[] = []
  const number = getVGpuNumber(resource)
  const { concise = false, detailed = false, separator = ' / ' } = options

  if (number > 0) {
    if (detailed) {
      parts.push(`${number} 个 vGPU`)
    } else {
      parts.push(`${number} vGPU`)
    }
  }

  if ((resource?.vGpuMemory || 0) > 0) {
    parts.push(formatVGpuMemory(resource?.vGpuMemory, { concise, detailed }))
  }

  if ((resource?.vGpuCores || 0) > 0) {
    if (detailed) {
      parts.push(`${resource?.vGpuCores}% 算力`)
    } else {
      parts.push(`${resource?.vGpuCores}%`)
    }
  }

  return parts.join(separator)
}

export const buildVGpuFieldEntries = (
  resource?: VGpuSpecLike | null,
  t?: Translator
): VGpuFieldEntry[] => {
  const entries: VGpuFieldEntry[] = []
  const number = getVGpuNumber(resource)

  if (number > 0) {
    entries.push({
      key: 'count',
      label: getCountShortLabel(t),
      value: String(number)
    })
  }

  const memory = formatVGpuMemory(resource?.vGpuMemory, { concise: true })
  if (memory) {
    entries.push({
      key: 'memory',
      label: getMemoryShortLabel(t),
      value: memory
    })
  }

  if ((resource?.vGpuCores || 0) > 0) {
    entries.push({
      key: 'cores',
      label: getTranslatedLabelWithFallback(
        t,
        'vGpuCoresShort',
        'vGpuCores',
        'Compute'
      ),
      value: `${resource?.vGpuCores}%`
    })
  }

  return entries
}

export const buildGpuFieldEntries = (
  resource?: GpuSpecLike | null,
  t?: Translator
): VGpuFieldEntry[] => {
  const entries: VGpuFieldEntry[] = []

  if ((resource?.gpuCount || 0) > 0) {
    entries.push({
      key: 'count',
      label: getCountShortLabel(t),
      value: String(resource?.gpuCount || 0)
    })
  }

  const memory = formatGpuMemory(resource?.gpuMemory)
  if (memory) {
    entries.push({
      key: 'memory',
      label: getMemoryShortLabel(t),
      value: memory
    })
  }

  return entries
}
