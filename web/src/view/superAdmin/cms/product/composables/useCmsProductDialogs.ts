import { computed, reactive, ref, watch, type Ref } from 'vue'
import { useDebounceFn } from '@vueuse/core'
import { ElMessage } from 'element-plus'
import {
  createCMSProduct,
  getCMSClusterNodes,
  updateCMSProduct,
  updateCMSProductPrice
} from '@/api/cms'
import type { Translator } from '@/types/consoleResource'
import type {
  CmsClusterOption,
  CmsNodeRow,
  CmsProductForm,
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
  sanitizeProductPayload
} from './productPageUtils'

interface UseCmsProductDialogsOptions {
  activeTab: Ref<CmsProductType>
  clusters: Ref<CmsClusterOption[]>
  fetchAreas: () => Promise<void>
  fetchProducts: (silent?: boolean) => Promise<void>
  t: Translator
}

export const useCmsProductDialogs = ({
  activeTab,
  clusters,
  fetchAreas,
  fetchProducts,
  t
}: UseCmsProductDialogsOptions) => {
  const loadingNodes = ref(false)
  const submitting = ref(false)
  const clusterNodes = ref<CmsNodeRow[]>([])
  const showProductDialog = ref(false)
  const showPriceDialog = ref(false)
  const isEdit = ref(false)
  const resourceType = ref<CmsProductResourceType>('cpu')
  const selectedNode = ref<CmsNodeRow | null>(null)
  const productForm = reactive<CmsProductForm>(createDefaultProductForm())
  const priceForm = reactive<CmsProductPriceForm>(createDefaultPriceForm())

  const computeRules = computed(() => createComputeRules(t))
  const storageRules = computed(() => createStorageRules(t))
  const nodeMaxCpu = computed(() =>
    selectedNode.value ? (selectedNode.value.cpu ?? 0) : 256
  )
  const nodeMaxMemory = computed(() =>
    selectedNode.value ? (selectedNode.value.memory ?? 0) : 2048
  )
  const nodeMaxGpuCount = computed(() =>
    selectedNode.value ? (selectedNode.value.gpuCount ?? 0) : 16
  )
  const nodeMaxGpuMemory = computed(() =>
    selectedNode.value ? (selectedNode.value.gpuMemory ?? 0) : 256
  )
  const nodeMaxVGpuCount = computed(() =>
    selectedNode.value ? (selectedNode.value.vGpuNumber ?? 0) : 100
  )
  const nodeMaxVGpuMemory = computed(() =>
    selectedNode.value ? (selectedNode.value.vGpuMemory ?? 0) : 256
  )
  const nodeMaxVGpuCores = computed(() =>
    selectedNode.value ? (selectedNode.value.vGpuCores ?? 0) : 100
  )

  const dialogTitle = computed(() => {
    if (isEdit.value) {
      return t('edit')
    }

    return productForm.productType === 1
      ? t('newComputeProduct')
      : t('newStorageProduct')
  })

  const canSubmit = computed(() => {
    if (isEdit.value) {
      return true
    }

    if (productForm.productType === 1) {
      return Boolean(productForm.nodeName)
    }

    return Boolean(
      productForm.name && productForm.area && productForm.storageClass
    )
  })

  const fetchClusterNodes = async (
    clusterId: CmsProductForm['clusterId']
  ): Promise<void> => {
    if (!clusterId) {
      clusterNodes.value = []
      return
    }

    loadingNodes.value = true

    try {
      const res = (await getCMSClusterNodes({
        clusterId,
        cpu: productForm.cpu || 0,
        memory: productForm.memory || 0,
        gpuCount: productForm.gpuCount || 0
      })) as ApiResponse<CmsNodeRow[]>

      if (res.code === 0) {
        clusterNodes.value = normalizeClusterNodes(res.data ?? [])
      } else {
        ElMessage.error(res.msg || t('failed'))
      }
    } catch (error: unknown) {
      console.error(error)
      ElMessage.error(t('failed'))
    } finally {
      loadingNodes.value = false
    }
  }

  const resetProductForm = (): void => {
    Object.assign(productForm, createDefaultProductForm(activeTab.value))
    resourceType.value = 'cpu'
    selectedNode.value = null
    clusterNodes.value = []
  }

  const resetPriceForm = (): void => {
    Object.assign(priceForm, createDefaultPriceForm())
  }

  const syncMissingResourceFieldsFromNode = (
    type: CmsProductResourceType,
    node: CmsNodeRow | null
  ): void => {
    if (!node) {
      return
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
      if ((productForm.vGpuCount || 0) === 0) {
        productForm.vGpuCount = node.vGpuNumber || node.vGpuCount || 0
      }

      if ((productForm.vGpuMemory || 0) === 0) {
        productForm.vGpuMemory = node.vGpuMemory || 0
      }

      if ((productForm.vGpuCores || 0) === 0) {
        productForm.vGpuCores = node.vGpuCores || 0
      }
    }
  }

  const hydrateEditNodeContext = async (
    row: CmsProductRow
  ): Promise<void> => {
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
    const node = selectedNode.value

    if (type === 'gpu') {
      productForm.gpuModel = node?.gpuModel || ''
      productForm.gpuCount = node?.gpuCount || 0
      productForm.gpuMemory = node?.gpuMemory || 0
      productForm.vGpuCount = 0
      productForm.vGpuMemory = 0
      productForm.vGpuCores = 0
    } else if (type === 'vgpu') {
      productForm.gpuModel = ''
      productForm.gpuCount = 0
      productForm.gpuMemory = 0
      productForm.vGpuCount = node?.vGpuNumber || 0
      productForm.vGpuMemory = node?.vGpuMemory || 0
      productForm.vGpuCores = node?.vGpuCores || 0
    } else {
      productForm.gpuModel = ''
      productForm.gpuCount = 0
      productForm.gpuMemory = 0
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

    if (clusterId) {
      void fetchClusterNodes(clusterId)
    } else {
      clusterNodes.value = []
    }
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

  const selectNode = (node: CmsNodeRow): void => {
    const selection = getNodeSelectionState(node)

    selectedNode.value = node
    resourceType.value = selection.resourceType
    Object.assign(productForm, selection.fields)

    if (!productForm.name) {
      productForm.name = selection.suggestedName
    }
  }

  const debouncedFetchNodes = useDebounceFn(
    (clusterId: CmsProductForm['clusterId']) => {
      void fetchClusterNodes(clusterId)
    },
    300
  )

  watch(
    () => [productForm.cpu, productForm.memory, productForm.gpuCount] as const,
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
      const api = isEdit.value ? updateCMSProduct : createCMSProduct
      const res = await api(
        sanitizeProductPayload(productForm, resourceType.value, isEdit.value)
      )

      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: t('success')
        })
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
      const res = await updateCMSProductPrice(priceForm)
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: t('success')
        })
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
    canSubmit,
    clusterNodes,
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
    priceForm,
    productForm,
    resetPriceForm,
    resourceType,
    selectNode,
    showPriceDialog,
    showProductDialog,
    storageRules,
    submitting
  }
}
