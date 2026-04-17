import type {
  FilterOption,
  ProductFilterData,
  Translator,
  VGpuFilterOption
} from '@/types/consoleResource'
import { buildVGpuFieldEntries, getVGpuNumber } from '@/utils/vgpu'

const getTranslatedLabel = (
  t: Translator | undefined,
  key: string,
  fallback: string
): string => {
  const translated = typeof t === 'function' ? t(key) : key
  return translated && translated !== key ? translated : fallback
}

export const buildVGpuFilterKey = (
  option: Pick<VGpuFilterOption, 'model' | 'vGpuNumber' | 'vGpuMemory' | 'vGpuCores'>
): string =>
  `vgpu:${option.model?.trim() || 'GPU'}:${option.vGpuNumber || 0}:${option.vGpuMemory || 0}:${option.vGpuCores || 0}`

const buildGpuFilterKey = (option: Pick<FilterOption, 'model'>): string =>
  `gpu:${option.model}`

const sortGpuFilterOptions = (left: FilterOption, right: FilterOption): number => {
  const leftModel = left.gpuModel || left.model || ''
  const rightModel = right.gpuModel || right.model || ''
  const modelCompare = leftModel.localeCompare(rightModel)
  if (modelCompare !== 0) return modelCompare

  const leftTypeWeight = left.resourceType === 'vgpu' ? 1 : 0
  const rightTypeWeight = right.resourceType === 'vgpu' ? 1 : 0
  if (leftTypeWeight !== rightTypeWeight) return leftTypeWeight - rightTypeWeight

  if ((left.vGpuMemory || 0) !== (right.vGpuMemory || 0)) {
    return (left.vGpuMemory || 0) - (right.vGpuMemory || 0)
  }
  if ((left.vGpuNumber || 0) !== (right.vGpuNumber || 0)) {
    return (left.vGpuNumber || 0) - (right.vGpuNumber || 0)
  }
  return (left.vGpuCores || 0) - (right.vGpuCores || 0)
}

export const buildGpuResourceFilterOptions = (
  data?: ProductFilterData | null,
  t?: Translator
): FilterOption[] => {
  const gpuLabel = getTranslatedLabel(t, 'gpu', 'GPU')

  const gpuOptions = (data?.gpuModels || []).map<FilterOption>((option) => ({
    ...option,
    key: buildGpuFilterKey(option),
    label: option.model ? `${option.model} / ${gpuLabel}` : gpuLabel,
    resourceType: 'gpu',
    gpuModel: option.model
  }))

  const vgpuOptions = (data?.vgpuModels || []).map<FilterOption>((option) => {
    const gpuModel = option.model?.trim() || 'GPU'
    const metaFields = buildVGpuFieldEntries(option, t)

    return {
      model: gpuModel,
      available: option.available,
      total: option.total,
      key: buildVGpuFilterKey(option),
      label: `${gpuModel} vGPU`,
      meta: metaFields.map((entry) => `${entry.label} ${entry.value}`).join(' / '),
      metaFields,
      resourceType: 'vgpu',
      gpuModel,
      vGpuNumber: getVGpuNumber(option),
      vGpuMemory: option.vGpuMemory || 0,
      vGpuCores: option.vGpuCores || 0
    }
  })

  return [...gpuOptions, ...vgpuOptions].sort(sortGpuFilterOptions)
}

export const findGpuResourceFilterOption = (
  options: FilterOption[],
  value: string
): FilterOption | undefined =>
  options.find((option) => (option.key || option.model) === value)
