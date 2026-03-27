export const createDefaultProductForm = (productType = 1) => ({
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

export const createDefaultPriceForm = () => ({
  id: null,
  priceHourly: 0,
  priceDaily: 0,
  priceWeekly: 0,
  priceMonthly: 0
})

export const createComputeRules = (t) => ({
  name: [{ required: true, message: t('inputProductName'), trigger: 'blur' }],
  clusterId: [{ required: true, message: t('selectCluster'), trigger: 'change' }],
  nodeName: [{ required: true, message: t('selectNode'), trigger: 'change' }],
  area: [{ required: true, message: t('inputArea'), trigger: 'blur' }],
  cpu: [{ required: true, message: t('inputCpu'), trigger: 'blur' }],
  memory: [{ required: true, message: t('inputMemory'), trigger: 'blur' }],
  priceHourly: [{ required: true, message: t('inputPrice'), trigger: 'blur' }]
})

export const createStorageRules = (t) => ({
  name: [{ required: true, message: t('inputProductName'), trigger: 'blur' }],
  clusterId: [{ required: true, message: t('selectCluster'), trigger: 'change' }],
  area: [{ required: true, message: t('inputArea'), trigger: 'blur' }],
  storageClass: [{ required: true, message: t('inputStorageClass'), trigger: 'blur' }],
  storagePriceGb: [{ required: true, message: t('inputStoragePrice'), trigger: 'blur' }]
})

export const normalizeClusterNodes = (nodes = []) => (
  nodes.map((node) => ({
    ...node,
    cpu: node.cpuAllocatable || 0,
    memory: node.memoryAllocatable || 0
  }))
)

export const resolveProductResourceType = (product = {}) => {
  if (product.vGpuCount > 0 || product.vGpuMemory > 0 || product.vGpuCores > 0) {
    return 'vgpu'
  }

  if (product.gpuCount > 0 || product.gpuModel) {
    return 'gpu'
  }

  return 'cpu'
}

export const buildProductFormFromRow = (row = {}, fallbackProductType = 1) => ({
  ...createDefaultProductForm(row.productType || fallbackProductType),
  id: row.id,
  productType: row.productType || fallbackProductType,
  name: row.name,
  description: row.description,
  clusterId: row.clusterId,
  area: row.area,
  nodeName: row.nodeName || '',
  nodeType: row.nodeType || '',
  cpuModel: row.cpuModel || '',
  cpu: row.cpu || 0,
  memory: row.memory || 0,
  gpuModel: row.gpuModel || '',
  gpuCount: row.gpuCount || 0,
  gpuMemory: row.gpuMemory || 0,
  vGpuCount: row.vGpuCount || 0,
  vGpuMemory: row.vGpuMemory || 0,
  vGpuCores: row.vGpuCores || 0,
  priceHourly: row.priceHourly || 0,
  priceDaily: row.priceDaily || 0,
  priceWeekly: row.priceWeekly || 0,
  priceMonthly: row.priceMonthly || 0,
  status: row.status,
  maxInstances: row.maxInstances || 0,
  storageClass: row.storageClass || '',
  storagePriceGb: row.storagePriceGb || 0
})

export const getNodeSelectionState = (node = {}) => {
  const baseFields = {
    nodeName: node.nodeName,
    area: node.area || '',
    cpu: node.cpu || 0,
    memory: node.memory || 0,
    cpuModel: node.cpuModel || ''
  }

  if (node.vGpuNumber > 0) {
    return {
      resourceType: 'vgpu',
      fields: {
        ...baseFields,
        gpuModel: '',
        gpuCount: 0,
        gpuMemory: 0,
        vGpuCount: node.vGpuNumber || 0,
        vGpuMemory: node.vGpuMemory || 0,
        vGpuCores: node.vGpuCores || 0
      },
      suggestedName: `vGPU ${node.vGpuNumber}`
    }
  }

  if (node.gpuCount > 0) {
    return {
      resourceType: 'gpu',
      fields: {
        ...baseFields,
        gpuModel: node.gpuModel || '',
        gpuCount: node.gpuCount || 0,
        gpuMemory: node.gpuMemory || 0,
        vGpuCount: 0,
        vGpuMemory: 0,
        vGpuCores: 0
      },
      suggestedName: `${node.gpuModel || 'GPU'} × ${node.gpuCount || 0}`
    }
  }

  return {
    resourceType: 'cpu',
    fields: {
      ...baseFields,
      gpuModel: '',
      gpuCount: 0,
      gpuMemory: 0,
      vGpuCount: 0,
      vGpuMemory: 0,
      vGpuCores: 0
    },
    suggestedName: `CPU ${node.cpu || 0}核`
  }
}

export const sanitizeProductPayload = (productForm, resourceType, isEdit) => {
  const data = { ...productForm }

  if (resourceType !== 'gpu') {
    data.gpuModel = ''
    data.gpuCount = 0
    data.gpuMemory = 0
  }

  if (resourceType !== 'vgpu') {
    data.vGpuCount = 0
    data.vGpuMemory = 0
    data.vGpuCores = 0
  }

  if (!isEdit) {
    data.maxInstances = 0
  }

  return data
}
