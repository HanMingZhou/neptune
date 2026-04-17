import type { FormRules } from 'element-plus'
import type { Translator } from '@/types/consoleResource'
import type {
  CmsBatchCreateComputeProductPayload,
  CmsNodeRow,
  CmsNodeSelectionState,
  CmsProductForm,
  CmsProductPriceItem,
  CmsProductPriceForm,
  CmsProductPricePayload,
  CmsProductResourceType,
  CmsProductRow,
  CmsProductType
} from '@/types/superAdmin'
import { formatVGpuSpec, getVGpuNumber, hasVGpuSpec } from '@/utils/vgpu'

export const createDefaultProductForm = (
  productType: CmsProductType = 1
): CmsProductForm => ({
  id: null,
  productType,
  name: '',
  description: '',
  clusterId: null,
  area: '',
  nodeName: '',
  nodeType: '',
  cpuModel: '',
  cpu: 0,
  memory: 0,
  gpuModel: '',
  gpuCount: 0,
  gpuMemory: 0,
  driverVersion: '',
  cudaVersion: '',
  vGpuNumber: 0,
  vGpuCount: 0,
  vGpuMemory: 0,
  vGpuCores: 0,
  priceHourly: 0,
  priceDaily: 0,
  priceWeekly: 0,
  priceMonthly: 0,
  status: 1,
  maxInstances: 0,
  storageClass: '',
  storagePriceGb: 0
})

export const createDefaultPriceForm = (): CmsProductPriceForm => ({
  id: null,
  priceHourly: 0,
  priceDaily: 0,
  priceWeekly: 0,
  priceMonthly: 0
})

export const buildComputePriceItems = (
  source: Pick<
    CmsProductForm | CmsProductPriceForm,
    'priceHourly' | 'priceDaily' | 'priceWeekly' | 'priceMonthly'
  >
): CmsProductPriceItem[] => [
  { priceType: 1, price: source.priceHourly || 0 },
  { priceType: 2, price: source.priceDaily || 0 },
  { priceType: 3, price: source.priceWeekly || 0 },
  { priceType: 4, price: source.priceMonthly || 0 }
]

export const createComputeRules = (
  t: Translator
): FormRules<CmsProductForm> => ({
  name: [{ required: true, message: t('inputProductName'), trigger: 'blur' }],
  clusterId: [
    { required: true, message: t('selectCluster'), trigger: 'change' }
  ],
  nodeName: [{ required: true, message: t('selectNode'), trigger: 'change' }],
  area: [{ required: true, message: t('inputArea'), trigger: 'blur' }],
  cpu: [{ required: true, message: t('inputCpu'), trigger: 'blur' }],
  memory: [{ required: true, message: t('inputMemory'), trigger: 'blur' }],
  priceHourly: [{ required: true, message: t('inputPrice'), trigger: 'blur' }]
})

export const createStorageRules = (
  t: Translator
): FormRules<CmsProductForm> => ({
  name: [{ required: true, message: t('inputProductName'), trigger: 'blur' }],
  clusterId: [
    { required: true, message: t('selectCluster'), trigger: 'change' }
  ],
  area: [{ required: true, message: t('inputArea'), trigger: 'blur' }],
  storageClass: [
    { required: true, message: t('inputStorageClass'), trigger: 'blur' }
  ],
  storagePriceGb: [
    { required: true, message: t('inputStoragePrice'), trigger: 'blur' }
  ]
})

export const normalizeClusterNodes = <T extends CmsNodeRow>(
  nodes: T[] = []
): T[] =>
  nodes.map((node) => ({
    ...node,
    cpu: node.cpuAllocatable || 0,
    memory: node.memoryAllocatable || 0
  }))

export const resolveProductResourceType = (
  product: Partial<CmsProductRow> = {}
): CmsProductResourceType => {
  if (hasVGpuSpec(product)) {
    return 'vgpu'
  }

  if ((product.gpuCount || 0) > 0 || Boolean(product.gpuModel)) {
    return 'gpu'
  }

  return 'cpu'
}

export const buildProductFormFromRow = (
  row: Partial<CmsProductRow> = {},
  fallbackProductType: CmsProductType = 1
): CmsProductForm => ({
  ...createDefaultProductForm(row.productType || fallbackProductType),
  id: row.id ?? null,
  productType: row.productType || fallbackProductType,
  name: row.name || '',
  description: row.description || '',
  clusterId: row.clusterId ?? null,
  area: row.area || '',
  nodeName: row.nodeName || '',
  nodeType: row.nodeType || '',
  cpuModel: row.cpuModel || '',
  cpu: row.cpu || 0,
  memory: row.memory || 0,
  gpuModel: row.gpuModel || '',
  gpuCount: row.gpuCount || 0,
  gpuMemory: row.gpuMemory || 0,
  driverVersion: row.driverVersion || '',
  cudaVersion: row.cudaVersion || '',
  vGpuNumber: getVGpuNumber(row),
  vGpuCount: getVGpuNumber(row),
  vGpuMemory: row.vGpuMemory || 0,
  vGpuCores: row.vGpuCores || 0,
  priceHourly: row.priceHourly || 0,
  priceDaily: row.priceDaily || 0,
  priceWeekly: row.priceWeekly || 0,
  priceMonthly: row.priceMonthly || 0,
  status: row.status ?? 1,
  maxInstances: row.maxInstances || 0,
  storageClass: row.storageClass || '',
  storagePriceGb: row.storagePriceGb || 0
})

export const buildNodeContextFromProduct = (
  product: Partial<CmsProductRow> = {}
): CmsNodeRow => ({
  nodeName: product.nodeName || '',
  area: product.area || '',
  cpu: product.cpu || 0,
  memory: product.memory || 0,
  cpuModel: product.cpuModel || '',
  gpuModel: product.gpuModel || '',
  gpuCount: product.gpuCount || 0,
  gpuMemory: product.gpuMemory || 0,
  driverVersion: product.driverVersion || '',
  cudaVersion: product.cudaVersion || '',
  vGpuNumber: getVGpuNumber(product),
  vGpuCount: getVGpuNumber(product),
  vGpuMemory: product.vGpuMemory || 0,
  vGpuCores: product.vGpuCores || 0
})

export const getNodeSelectionState = (
  node: Partial<CmsNodeRow> = {}
): CmsNodeSelectionState => {
  const baseFields: CmsNodeSelectionState['fields'] = {
    nodeName: node.nodeName || '',
    area: node.area || '',
    cpu: node.cpu || 0,
    memory: node.memory || 0,
    cpuModel: node.cpuModel || '',
    gpuModel: node.gpuModel || '',
    gpuCount: 0,
    gpuMemory: 0,
    driverVersion: node.driverVersion || '',
    cudaVersion: node.cudaVersion || '',
    vGpuNumber: 0,
    vGpuCount: 0,
    vGpuMemory: 0,
    vGpuCores: 0
  }

  if (hasVGpuSpec(node)) {
    return {
      resourceType: 'vgpu',
      fields: {
        ...baseFields,
        vGpuNumber: getVGpuNumber(node),
        vGpuCount: getVGpuNumber(node),
        vGpuMemory: node.vGpuMemory || 0,
        vGpuCores: node.vGpuCores || 0
      },
      suggestedName: formatVGpuSpec(node)
        ? `vGPU ${formatVGpuSpec(node)}`
        : 'vGPU'
    }
  }

  if ((node.gpuCount || 0) > 0) {
    return {
      resourceType: 'gpu',
      fields: {
        ...baseFields,
        gpuModel: node.gpuModel || '',
        gpuCount: node.gpuCount || 0,
        gpuMemory: node.gpuMemory || 0
      },
      suggestedName: `${node.gpuModel || 'GPU'} × ${node.gpuCount || 0}`
    }
  }

  return {
    resourceType: 'cpu',
    fields: baseFields,
    suggestedName: `CPU ${node.cpu || 0}核`
  }
}

export const sanitizeProductPayload = (
  productForm: CmsProductForm,
  resourceType: CmsProductResourceType,
  isEdit: boolean
): CmsProductForm => {
  const data: CmsProductForm = { ...productForm }
  const normalizedVGpuNumber =
    productForm.vGpuCount || productForm.vGpuNumber || 0

  if (resourceType === 'cpu') {
    data.gpuModel = ''
  }

  if (resourceType !== 'gpu') {
    data.gpuCount = 0
    data.gpuMemory = 0
  }

  if (resourceType !== 'vgpu') {
    data.vGpuNumber = 0
    data.vGpuCount = 0
    data.vGpuMemory = 0
    data.vGpuCores = 0
  } else {
    data.vGpuNumber = normalizedVGpuNumber
    data.vGpuCount = normalizedVGpuNumber
  }

  if (!isEdit) {
    data.maxInstances = 0
  }

  if (data.productType === 1) {
    data.prices = buildComputePriceItems(data)
  } else {
    data.prices = undefined
  }

  return data
}

export const sanitizeBatchComputeProductPayload = (
  productForm: CmsProductForm,
  resourceType: CmsProductResourceType,
  nodeNames: string[]
): CmsBatchCreateComputeProductPayload => {
  const data = sanitizeProductPayload(productForm, resourceType, false)

  return {
    productType: data.productType,
    nodeNames,
    name: data.name,
    description: data.description,
    clusterId: data.clusterId || 0,
    area: data.area,
    nodeType: data.nodeType,
    cpuModel: data.cpuModel,
    cpu: data.cpu,
    memory: data.memory,
    gpuModel: data.gpuModel,
    gpuCount: data.gpuCount,
    gpuMemory: data.gpuMemory,
    vGpuNumber: data.vGpuNumber || 0,
    vGpuMemory: data.vGpuMemory,
    vGpuCores: data.vGpuCores,
    prices: data.prices || buildComputePriceItems(data),
    driverVersion: data.driverVersion || '',
    cudaVersion: data.cudaVersion || '',
    systemDisk: 0,
    dataDisk: 0,
    status: data.status,
    maxInstances: data.maxInstances || 0
  }
}

export const sanitizeProductPricePayload = (
  priceForm: CmsProductPriceForm
): CmsProductPricePayload => ({
  id: priceForm.id,
  prices: buildComputePriceItems(priceForm)
})
