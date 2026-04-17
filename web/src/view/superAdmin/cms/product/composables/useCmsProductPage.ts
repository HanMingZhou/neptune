import { ElMessage, ElMessageBox } from 'element-plus'
import { deleteCMSProduct } from '@/api/cms'
import type { Translator } from '@/types/consoleResource'
import type { CmsProductRow } from '@/types/superAdmin'
import { useCmsProductCatalog } from './useCmsProductCatalog'
import { useCmsProductDialogs } from './useCmsProductDialogs'

interface UseCmsProductPageOptions {
  t?: Translator
}

const isDialogCancel = (error: unknown): error is 'cancel' | 'close' =>
  error === 'cancel' || error === 'close'

export function useCmsProductPage({ t }: UseCmsProductPageOptions = {}) {
  const translate: Translator = t || ((key: string) => key)

  const {
    activeTab,
    areas,
    clusters,
    currentPage,
    fetchAreas,
    fetchFilterOptions,
    fetchProducts,
    filterArea,
    filterAvailableMax,
    filterAvailableMin,
    filterClusterId,
    filterResourceType,
    filterGpuModel,
    filterKeyword,
    filterMaxInstancesMax,
    filterMaxInstancesMin,
    filterPriceField,
    filterPriceMax,
    filterPriceMin,
    filterUsedCapacityMax,
    filterUsedCapacityMin,
    handlePageChange,
    handleResetFilters,
    handleSearch,
    handleSizeChange,
    initialize,
    loading,
    pageSize,
    products,
    setActiveTab,
    gpuModels,
    total
  } = useCmsProductCatalog({ t: translate })

  const {
    activePreviewNodeName,
    canSubmit,
    clearSelectedNodes,
    clusterNodes,
    computeRules,
    dialogTitle,
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
    showPriceDialog,
    showProductDialog,
    storageRules,
    submitButtonText,
    submitting,
    toggleNodeSelection
  } = useCmsProductDialogs({
    activeTab,
    clusters,
    fetchAreas,
    fetchProducts,
    t: translate
  })

  const handleDelete = async (row: CmsProductRow): Promise<void> => {
    if (row.id === null) {
      ElMessage.error(translate('failed'))
      return
    }

    try {
      await ElMessageBox.confirm(
        translate('deleteProductConfirm'),
        translate('tip'),
        {
          confirmButtonText: translate('confirm'),
          cancelButtonText: translate('cancel'),
          type: 'warning'
        }
      )

      const res = await deleteCMSProduct({ id: row.id })
      if (res.code === 0) {
        ElMessage.success(translate('success'))
        await fetchProducts()
      }
    } catch (error: unknown) {
      if (!isDialogCancel(error)) {
        ElMessage.error(translate('failed'))
      }
    }
  }

  return {
    activePreviewNodeName,
    activeTab,
    areas,
    canSubmit,
    clearSelectedNodes,
    clusterNodes,
    clusters,
    computeRules,
    currentPage,
    dialogTitle,
    fetchFilterOptions,
    fetchProducts,
    filterArea,
    filterAvailableMax,
    filterAvailableMin,
    filterClusterId,
    filterResourceType,
    filterGpuModel,
    filterKeyword,
    filterMaxInstancesMax,
    filterMaxInstancesMin,
    filterPriceField,
    filterPriceMax,
    filterPriceMin,
    filterUsedCapacityMax,
    filterUsedCapacityMin,
    handleClusterChange,
    handleDelete,
    handlePageChange,
    handleResetFilters,
    handleResourceTypeChange,
    handleSearch,
    handleSizeChange,
    handleStorageClusterChange,
    handleSubmitProduct,
    handleUpdatePrice,
    initialize,
    isEdit,
    loading,
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
    pageSize,
    priceForm,
    products,
    productForm,
    resetPriceForm,
    resourceType,
    selectedNodeCount,
    selectedNodeLabel,
    selectedNodeNames,
    setActiveTab,
    showPriceDialog,
    showProductDialog,
    storageRules,
    submitButtonText,
    submitting,
    toggleNodeSelection,
    gpuModels,
    total
  }
}
