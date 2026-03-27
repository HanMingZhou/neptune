import { computed, reactive, ref, watch } from 'vue'
import { useDebounceFn } from '@vueuse/core'
import { ElMessage } from 'element-plus'
import {
  createCMSProduct,
  getCMSClusterNodes,
  updateCMSProduct,
  updateCMSProductPrice
} from '@/api/cms'
import {
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

export const useCmsProductDialogs = ({ activeTab, clusters, fetchAreas, fetchProducts, t }) => {
  const loadingNodes = ref(false)
  const submitting = ref(false)
  const clusterNodes = ref([])
  const showProductDialog = ref(false)
  const showPriceDialog = ref(false)
  const isEdit = ref(false)
  const resourceType = ref('cpu')
  const selectedNode = ref(null)
  const productForm = reactive(createDefaultProductForm())
  const priceForm = reactive(createDefaultPriceForm())

  const computeRules = computed(() => createComputeRules(t))
  const storageRules = computed(() => createStorageRules(t))
  const nodeMaxCpu = computed(() => selectedNode.value ? (selectedNode.value.cpu ?? 0) : 256)
  const nodeMaxMemory = computed(() => selectedNode.value ? (selectedNode.value.memory ?? 0) : 2048)
  const nodeMaxGpuCount = computed(() => selectedNode.value ? (selectedNode.value.gpuCount ?? 0) : 16)
  const nodeMaxGpuMemory = computed(() => selectedNode.value ? (selectedNode.value.gpuMemory ?? 0) : 256)
  const nodeMaxVGpuCount = computed(() => selectedNode.value ? (selectedNode.value.vGpuNumber ?? 0) : 100)
  const nodeMaxVGpuMemory = computed(() => selectedNode.value ? (selectedNode.value.vGpuMemory ?? 0) : 256)
  const nodeMaxVGpuCores = computed(() => selectedNode.value ? (selectedNode.value.vGpuCores ?? 0) : 100)

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

    return Boolean(productForm.name && productForm.area && productForm.storageClass)
  })

  const fetchClusterNodes = async (clusterId) => {
    if (!clusterId) {
      clusterNodes.value = []
      return
    }

    loadingNodes.value = true

    try {
      const res = await getCMSClusterNodes({
        clusterId,
        cpu: productForm.cpu || 0,
        memory: productForm.memory || 0,
        gpuCount: productForm.gpuCount || 0
      })

      if (res.code === 0) {
        clusterNodes.value = normalizeClusterNodes(res.data || [])
      } else {
        ElMessage.error(res.msg || t('failed'))
      }
    } catch (error) {
      console.error(error)
      ElMessage.error(t('failed'))
    } finally {
      loadingNodes.value = false
    }
  }

  const resetProductForm = () => {
    Object.assign(productForm, createDefaultProductForm(activeTab.value))
    resourceType.value = 'cpu'
    selectedNode.value = null
    clusterNodes.value = []
  }

  const resetPriceForm = () => {
    Object.assign(priceForm, createDefaultPriceForm())
  }

  const openCreateDialog = () => {
    resetProductForm()
    isEdit.value = false
    showProductDialog.value = true
  }

  const openEditDialog = (row) => {
    Object.assign(productForm, buildProductFormFromRow(row, activeTab.value))
    selectedNode.value = null
    clusterNodes.value = []
    isEdit.value = true
    resourceType.value = resolveProductResourceType(row)
    showProductDialog.value = true
  }

  const openPriceDialog = (row) => {
    Object.assign(priceForm, {
      id: row.id,
      priceHourly: row.priceHourly || 0,
      priceDaily: row.priceDaily || 0,
      priceWeekly: row.priceWeekly || 0,
      priceMonthly: row.priceMonthly || 0
    })
    showPriceDialog.value = true
  }

  const handleResourceTypeChange = (type) => {
    if (type !== 'gpu') {
      productForm.gpuModel = ''
      productForm.gpuCount = 0
      productForm.gpuMemory = 0
    }

    if (type !== 'vgpu') {
      productForm.vGpuCount = 0
      productForm.vGpuMemory = 0
      productForm.vGpuCores = 0
    }
  }

  const handleClusterChange = (clusterId) => {
    productForm.nodeName = ''
    selectedNode.value = null

    if (clusterId) {
      fetchClusterNodes(clusterId)
    } else {
      clusterNodes.value = []
    }
  }

  const handleStorageClusterChange = (clusterId) => {
    if (!clusterId) {
      return
    }

    const cluster = clusters.value.find((item) => item.id === clusterId)
    if (cluster?.area) {
      productForm.area = cluster.area
    }
  }

  const selectNode = (node) => {
    const selection = getNodeSelectionState(node)

    selectedNode.value = node
    resourceType.value = selection.resourceType
    Object.assign(productForm, selection.fields)

    if (!productForm.name) {
      productForm.name = selection.suggestedName
    }
  }

  const debouncedFetchNodes = useDebounceFn((clusterId) => {
    fetchClusterNodes(clusterId)
  }, 300)

  watch(
    () => [productForm.cpu, productForm.memory, productForm.gpuCount],
    () => {
      if (productForm.clusterId && !isEdit.value && productForm.productType === 1) {
        debouncedFetchNodes(productForm.clusterId)
      }
    }
  )

  const handleSubmitProduct = async () => {
    submitting.value = true

    try {
      const api = isEdit.value ? updateCMSProduct : createCMSProduct
      const res = await api(sanitizeProductPayload(productForm, resourceType.value, isEdit.value))

      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: t('success')
        })
        showProductDialog.value = false
        await Promise.all([
          fetchProducts(),
          fetchAreas()
        ])
      } else {
        ElMessage.error(res.msg || t('failed'))
      }
    } catch (error) {
      console.error(error)
      ElMessage.error(t('failed'))
    } finally {
      submitting.value = false
    }
  }

  const handleUpdatePrice = async () => {
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
    } catch (error) {
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
