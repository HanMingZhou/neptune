import { computed, reactive, ref, watch, type Ref } from 'vue'
import { useDebounceFn } from '@vueuse/core'
import { ElMessage } from 'element-plus'
import {
  createCMSComputeProductsBatch,
  createCMSProduct,
  getCMSProductNodeCandidates,
  updateCMSProduct,
  updateCMSProductPrice
} from '@/api/cms'
import type { Translator } from '@/types/consoleResource'
import type {
  CmsBatchCreateComputeProductResult,
  CmsClusterOption,
  CmsNodeRow,
  CmsProductForm,
  CmsProductNodeCandidate,
  CmsProductPriceForm,
  CmsProductResourceType,
  CmsProductRow,
  CmsProductType
} from '@/types/superAdmin'
import type { ApiResponse } from '@/utils/request'
import {
  buildNodeContextFromProduct,
  buildProductFormFromRow,
  createComputeRules,
  createDefaultPriceForm,
  createDefaultProductForm,
  createStorageRules,
  getNodeSelectionState,
  normalizeClusterNodes,
  resolveProductResourceType,
  sanitizeBatchComputeProductPayload,
  sanitizeProductPricePayload,
  sanitizeProductPayload
} from './productPageUtils'
import { getVGpuNumber } from '@/utils/vgpu'

interface UseCmsProductDialogsOptions {
  activeTab: Ref<CmsProductType>
  clusters: Ref<CmsClusterOption[]>
  fetchAreas: () => Promise<void>
  fetchProducts: (silent?: boolean) => Promise<void>
  t: Translator
}

const DEFAULT_NODE_LIMITS = {
  cpu: 256,
  memory: 2048,
  gpuCount: 16,
  gpuMemory: 256,
  vGpuCount: 100,
  vGpuMemory: 256,
  vGpuCores: 100
} as const

export const useCmsProductDialogs = ({
  activeTab,
  clusters,
  fetchAreas,
  fetchProducts,
  t
}: UseCmsProductDialogsOptions) => {
  const loadingNodes = ref(false)
  const submitting = ref(false)
  const clusterNodes = ref<CmsProductNodeCandidate[]>([])
  const showProductDialog = ref(false)
  const showPriceDialog = ref(false)
  const isEdit = ref(false)
  const resourceType = ref<CmsProductResourceType>('cpu')
  const selectedNode = ref<CmsNodeRow | null>(null)
  const selectedNodeNames = ref<string[]>([])
  const activePreviewNodeName = ref('')
  const productForm = reactive<CmsProductForm>(createDefaultProductForm())
  const priceForm = reactive<CmsProductPriceForm>(createDefaultPriceForm())

  const computeRules = computed(() => createComputeRules(t))
  const storageRules = computed(() => createStorageRules(t))

  const createPreviewNode = computed<CmsProductNodeCandidate | null>(() => {
    if (isEdit.value) {
      return null
    }

    const byPreview =
      clusterNodes.value.find(
        (node) => node.nodeName === activePreviewNodeName.value
      ) || null
    if (byPreview) {
      return byPreview
    }

    const firstSelected = selectedNodeNames.value[0]
    if (firstSelected) {
      return (
        clusterNodes.value.find((node) => node.nodeName === firstSelected) ||
        null
      )
    }

    return clusterNodes.value[0] || null
  })

  const previewNode = computed<CmsNodeRow | CmsProductNodeCandidate | null>(
    () => (isEdit.value ? selectedNode.value : createPreviewNode.value)
  )

  const selectedNodeCandidates = computed(() =>
    clusterNodes.value.filter((node) =>
      selectedNodeNames.value.includes(node.nodeName)
    )
  )

  const selectedNodeCount = computed(() => selectedNodeNames.value.length)

  const editLimitNode = computed<CmsNodeRow | null>(() =>
    isEdit.value ? selectedNode.value : null
  )

  const nodeMaxCpu = computed(() =>
    editLimitNode.value
      ? (editLimitNode.value.cpu ?? 0)
      : DEFAULT_NODE_LIMITS.cpu
  )
  const nodeMaxMemory = computed(() =>
    editLimitNode.value
      ? (editLimitNode.value.memory ?? 0)
      : DEFAULT_NODE_LIMITS.memory
  )
  const nodeMaxGpuCount = computed(() =>
    editLimitNode.value
      ? (editLimitNode.value.gpuCount ?? 0)
      : DEFAULT_NODE_LIMITS.gpuCount
  )
  const nodeMaxGpuMemory = computed(() =>
    editLimitNode.value
      ? (editLimitNode.value.gpuMemory ?? 0)
      : DEFAULT_NODE_LIMITS.gpuMemory
  )
  const nodeMaxVGpuCount = computed(() =>
    editLimitNode.value
      ? getVGpuNumber(editLimitNode.value)
      : DEFAULT_NODE_LIMITS.vGpuCount
  )
  const nodeMaxVGpuMemory = computed(() =>
    editLimitNode.value
      ? (editLimitNode.value.vGpuMemory ?? 0)
      : DEFAULT_NODE_LIMITS.vGpuMemory
  )
  const nodeMaxVGpuCores = computed(() =>
    editLimitNode.value
      ? (editLimitNode.value.vGpuCores ?? 0)
      : DEFAULT_NODE_LIMITS.vGpuCores
  )

  const dialogTitle = computed(() => {
    if (isEdit.value) {
      return t('edit')
    }

    return productForm.productType === 1
      ? t('newComputeProduct')
      : t('newStorageProduct')
  })

  const submitButtonText = computed(() => {
    if (isEdit.value) {
      return t('save')
    }

    if (productForm.productType === 1 && selectedNodeCount.value > 1) {
      return t('createForSelectedNodes')
    }

    return t('create')
  })

  const selectedNodeLabel = computed(() => {
    if (selectedNodeCount.value === 0) {
      return t('selectAtLeastOneNode')
    }

    return selectedNodeNames.value.join(' / ')
  })

  const canSubmit = computed(() => {
    if (isEdit.value) {
      return true
    }

    if (productForm.productType === 1) {
      return selectedNodeNames.value.length > 0
    }

    return Boolean(
      productForm.name && productForm.area && productForm.storageClass
    )
  })

  const hasExplicitResourceSpec = (): boolean =>
    Boolean(
      productForm.cpu ||
      productForm.memory ||
      productForm.gpuCount ||
      productForm.gpuMemory ||
      productForm.vGpuNumber ||
      productForm.vGpuMemory ||
      productForm.vGpuCores
    )

  const syncMissingResourceFieldsFromNode = (
    type: CmsProductResourceType,
    node: CmsNodeRow | CmsProductNodeCandidate | null
  ): void => {
    if (!node) {
      return
    }

    if (!productForm.area) {
      productForm.area = node.area || ''
    }
    if (!productForm.cpuModel) {
      productForm.cpuModel = node.cpuModel || ''
    }
    if (!productForm.driverVersion) {
      productForm.driverVersion = node.driverVersion || ''
    }
    if (!productForm.cudaVersion) {
      productForm.cudaVersion = node.cudaVersion || ''
    }

    if (type === 'gpu') {
      if (!productForm.gpuModel) {
        productForm.gpuModel = node.gpuModel || ''
      }
      if ((productForm.gpuCount || 0) === 0) {
        productForm.gpuCount = node.gpuCount || 0
      }
      if ((productForm.gpuMemory || 0) === 0) {
        productForm.gpuMemory = node.gpuMemory || 0
      }
      return
    }

    if (type === 'vgpu') {
      const vGpuNumber = getVGpuNumber(node)
      if (!productForm.gpuModel) {
        productForm.gpuModel = node.gpuModel || ''
      }
      if ((productForm.vGpuNumber || 0) === 0) {
        productForm.vGpuNumber = vGpuNumber
      }
      if ((productForm.vGpuCount || 0) === 0) {
        productForm.vGpuCount = vGpuNumber
      }
      if ((productForm.vGpuMemory || 0) === 0) {
        productForm.vGpuMemory = node.vGpuMemory || 0
      }
      if ((productForm.vGpuCores || 0) === 0) {
        productForm.vGpuCores = node.vGpuCores || 0
      }
    }
  }

  const sortSelectionByCandidates = (nodeNames: string[]): string[] => {
    const selected = new Set(nodeNames)
    return clusterNodes.value
      .filter((node) => selected.has(node.nodeName))
      .map((node) => node.nodeName)
  }

  const syncSelectedNodeField = (): void => {
    productForm.nodeName = selectedNodeNames.value[0] || ''
  }

  const resetProductForm = (): void => {
    Object.assign(productForm, createDefaultProductForm(activeTab.value))
    resourceType.value = 'cpu'
    selectedNode.value = null
    selectedNodeNames.value = []
    activePreviewNodeName.value = ''
    clusterNodes.value = []
  }

  const resetPriceForm = (): void => {
    Object.assign(priceForm, createDefaultPriceForm())
  }

  const refreshPreviewPointer = (): void => {
    if (
      activePreviewNodeName.value &&
      clusterNodes.value.some(
        (node) => node.nodeName === activePreviewNodeName.value
      )
    ) {
      return
    }

    activePreviewNodeName.value =
      selectedNodeNames.value[0] || clusterNodes.value[0]?.nodeName || ''
  }

  const fetchClusterNodes = async (
    clusterId: CmsProductForm['clusterId']
  ): Promise<void> => {
    if (!clusterId) {
      clusterNodes.value = []
      selectedNodeNames.value = []
      activePreviewNodeName.value = ''
      syncSelectedNodeField()
      return
    }

    loadingNodes.value = true

    try {
      const res = (await getCMSProductNodeCandidates({
        clusterId,
        resourceType: resourceType.value,
        cpu: productForm.cpu || 0,
        memory: productForm.memory || 0,
        gpuCount: productForm.gpuCount || 0,
        gpuMemory: productForm.gpuMemory || 0,
        vGpuNumber: productForm.vGpuCount || productForm.vGpuNumber || 0,
        vGpuMemory: productForm.vGpuMemory || 0,
        vGpuCores: productForm.vGpuCores || 0,
        excludeProductId:
          isEdit.value && productForm.id ? productForm.id : undefined
      })) as ApiResponse<CmsProductNodeCandidate[]>

      if (res.code !== 0) {
        ElMessage.error(res.msg || t('failed'))
        return
      }

      clusterNodes.value = normalizeClusterNodes(res.data ?? [])

      if (!isEdit.value) {
        const previousSelected = [...selectedNodeNames.value]
        const selectable = new Set(
          clusterNodes.value
            .filter((node) => node.canCreateComputeProduct)
            .map((node) => node.nodeName)
        )
        selectedNodeNames.value = sortSelectionByCandidates(
          selectedNodeNames.value.filter((nodeName) => selectable.has(nodeName))
        )
        syncSelectedNodeField()

        const removedNodeNames = previousSelected.filter(
          (nodeName) => !selectedNodeNames.value.includes(nodeName)
        )
        if (removedNodeNames.length > 0) {
          const previewText = removedNodeNames.slice(0, 3).join(' / ')
          const suffix =
            removedNodeNames.length > 3
              ? t('selectionAdjustedMore', {
                  count: removedNodeNames.length - 3
                })
              : ''
          ElMessage.warning(
            t('selectionAdjustedBySpec', {
              nodes: `${previewText}${suffix}`
            })
          )
        }
      } else if (productForm.nodeName) {
        selectedNode.value =
          clusterNodes.value.find(
            (node) => node.nodeName === productForm.nodeName
          ) || selectedNode.value
      }

      refreshPreviewPointer()
    } catch (error: unknown) {
      console.error(error)
      ElMessage.error(t('failed'))
    } finally {
      loadingNodes.value = false
    }
  }

  const applyNodeDefaults = (node: CmsProductNodeCandidate): void => {
    const selection = getNodeSelectionState(node)

    if (!hasExplicitResourceSpec()) {
      resourceType.value = selection.resourceType
      Object.assign(productForm, selection.fields)
    } else {
      syncMissingResourceFieldsFromNode(resourceType.value, node)
    }

    syncMissingResourceFieldsFromNode(resourceType.value, node)

    if (!productForm.name) {
      productForm.name = selection.suggestedName
    }
  }

  const hydrateEditNodeContext = async (row: CmsProductRow): Promise<void> => {
    if (row.productType !== 1 || !row.clusterId || !row.nodeName) {
      return
    }

    await fetchClusterNodes(row.clusterId)

    const matchedNode =
      clusterNodes.value.find((node) => node.nodeName === row.nodeName) || null

    if (
      !matchedNode ||
      !showProductDialog.value ||
      !isEdit.value ||
      productForm.id !== row.id
    ) {
      return
    }

    selectedNode.value = matchedNode
    activePreviewNodeName.value = matchedNode.nodeName
    syncMissingResourceFieldsFromNode(resourceType.value, matchedNode)
  }

  const openCreateDialog = (): void => {
    resetProductForm()
    isEdit.value = false
    showProductDialog.value = true
  }

  const openEditDialog = (row: CmsProductRow): void => {
    Object.assign(productForm, buildProductFormFromRow(row, activeTab.value))
    selectedNode.value = buildNodeContextFromProduct(row)
    selectedNodeNames.value = row.nodeName ? [row.nodeName] : []
    activePreviewNodeName.value = row.nodeName || ''
    clusterNodes.value = []
    isEdit.value = true
    resourceType.value = resolveProductResourceType(row)
    showProductDialog.value = true

    void hydrateEditNodeContext(row)
  }

  const openPriceDialog = (row: CmsProductRow): void => {
    Object.assign(priceForm, {
      id: row.id,
      priceHourly: row.priceHourly || 0,
      priceDaily: row.priceDaily || 0,
      priceWeekly: row.priceWeekly || 0,
      priceMonthly: row.priceMonthly || 0
    })
    showPriceDialog.value = true
  }

  const handleResourceTypeChange = (type: CmsProductResourceType): void => {
    const node = previewNode.value

    if (type === 'gpu') {
      productForm.gpuModel = node?.gpuModel || ''
      productForm.gpuCount = node?.gpuCount || 0
      productForm.gpuMemory = node?.gpuMemory || 0
      productForm.vGpuNumber = 0
      productForm.vGpuCount = 0
      productForm.vGpuMemory = 0
      productForm.vGpuCores = 0
    } else if (type === 'vgpu') {
      const vGpuNumber = getVGpuNumber(node)
      productForm.gpuModel = node?.gpuModel || ''
      productForm.gpuCount = 0
      productForm.gpuMemory = 0
      productForm.vGpuNumber = vGpuNumber
      productForm.vGpuCount = vGpuNumber
      productForm.vGpuMemory = node?.vGpuMemory || 0
      productForm.vGpuCores = node?.vGpuCores || 0
    } else {
      productForm.gpuModel = ''
      productForm.gpuCount = 0
      productForm.gpuMemory = 0
      productForm.vGpuNumber = 0
      productForm.vGpuCount = 0
      productForm.vGpuMemory = 0
      productForm.vGpuCores = 0
    }
  }

  const handleClusterChange = (
    clusterId: CmsProductForm['clusterId']
  ): void => {
    productForm.nodeName = ''
    selectedNode.value = null
    selectedNodeNames.value = []
    activePreviewNodeName.value = ''

    if (!clusterId) {
      clusterNodes.value = []
      return
    }

    const cluster = clusters.value.find((item) => item.id === clusterId)
    if (cluster?.area) {
      productForm.area = cluster.area
    }

    void fetchClusterNodes(clusterId)
  }

  const handleStorageClusterChange = (
    clusterId: CmsProductForm['clusterId']
  ): void => {
    if (!clusterId) {
      return
    }

    const cluster = clusters.value.find((item) => item.id === clusterId)
    if (cluster?.area) {
      productForm.area = cluster.area
    }
  }

  const previewClusterNode = (node: CmsProductNodeCandidate): void => {
    activePreviewNodeName.value = node.nodeName
  }

  const clearSelectedNodes = (): void => {
    selectedNodeNames.value = []
    syncSelectedNodeField()
  }

  const toggleNodeSelection = (node: CmsProductNodeCandidate): void => {
    activePreviewNodeName.value = node.nodeName

    if (!node.canCreateComputeProduct) {
      return
    }

    const nextSelection = new Set(selectedNodeNames.value)
    if (nextSelection.has(node.nodeName)) {
      nextSelection.delete(node.nodeName)
    } else {
      nextSelection.add(node.nodeName)
    }

    selectedNodeNames.value = sortSelectionByCandidates([...nextSelection])
    syncSelectedNodeField()

    if (nextSelection.has(node.nodeName)) {
      applyNodeDefaults(node)
    }
  }

  const debouncedFetchNodes = useDebounceFn(
    (clusterId: CmsProductForm['clusterId']) => {
      void fetchClusterNodes(clusterId)
    },
    260
  )

  watch(
    () =>
      [
        productForm.cpu,
        productForm.memory,
        productForm.gpuCount,
        productForm.gpuMemory,
        productForm.vGpuCount,
        productForm.vGpuNumber,
        productForm.vGpuMemory,
        productForm.vGpuCores,
        resourceType.value
      ] as const,
    () => {
      if (
        productForm.clusterId &&
        !isEdit.value &&
        productForm.productType === 1
      ) {
        debouncedFetchNodes(productForm.clusterId)
      }
    }
  )

  const handleSubmitProduct = async (): Promise<void> => {
    submitting.value = true

    try {
      if (!isEdit.value && productForm.productType === 1) {
        if (selectedNodeNames.value.length === 0) {
          ElMessage.warning(t('selectAtLeastOneNode'))
          return
        }

        const res = (await createCMSComputeProductsBatch(
          sanitizeBatchComputeProductPayload(
            productForm,
            resourceType.value,
            selectedNodeNames.value
          )
        )) as ApiResponse<CmsBatchCreateComputeProductResult>

        if (res.code === 0) {
          ElMessage.success(
            t('batchCreateComputeSuccess', {
              count: res.data?.createdCount || selectedNodeNames.value.length
            })
          )
          showProductDialog.value = false
          await Promise.all([fetchProducts(), fetchAreas()])
        } else {
          ElMessage.error(res.msg || t('failed'))
        }
        return
      }

      const api = isEdit.value ? updateCMSProduct : createCMSProduct
      const res = await api(
        sanitizeProductPayload(productForm, resourceType.value, isEdit.value)
      )

      if (res.code === 0) {
        ElMessage.success(t('success'))
        showProductDialog.value = false
        await Promise.all([fetchProducts(), fetchAreas()])
      } else {
        ElMessage.error(res.msg || t('failed'))
      }
    } catch (error: unknown) {
      console.error(error)
      ElMessage.error(t('failed'))
    } finally {
      submitting.value = false
    }
  }

  const handleUpdatePrice = async (): Promise<void> => {
    submitting.value = true

    try {
      const res = await updateCMSProductPrice(sanitizeProductPricePayload(priceForm))
      if (res.code === 0) {
        ElMessage.success(t('success'))
        showPriceDialog.value = false
        await fetchProducts()
      } else {
        ElMessage.error(res.msg || t('failed'))
      }
    } catch (error: unknown) {
      console.error(error)
      ElMessage.error(t('failed'))
    } finally {
      submitting.value = false
    }
  }

  return {
    activePreviewNodeName,
    canSubmit,
    clusterNodes,
    clearSelectedNodes,
    computeRules,
    dialogTitle,
    fetchClusterNodes,
    handleClusterChange,
    handleResourceTypeChange,
    handleStorageClusterChange,
    handleSubmitProduct,
    handleUpdatePrice,
    isEdit,
    loadingNodes,
    nodeMaxCpu,
    nodeMaxGpuCount,
    nodeMaxGpuMemory,
    nodeMaxMemory,
    nodeMaxVGpuCores,
    nodeMaxVGpuCount,
    nodeMaxVGpuMemory,
    openCreateDialog,
    openEditDialog,
    openPriceDialog,
    previewClusterNode,
    previewNode,
    priceForm,
    productForm,
    resetPriceForm,
    resourceType,
    selectedNodeCount,
    selectedNodeLabel,
    selectedNodeNames,
    selectedNodeCandidates,
    showPriceDialog,
    showProductDialog,
    storageRules,
    submitButtonText,
    submitting,
    toggleNodeSelection
  }
}
