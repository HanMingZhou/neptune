import { ElMessage, ElMessageBox } from 'element-plus'
import { deleteCMSProduct } from '@/api/cms'
import { useCmsProductCatalog } from './useCmsProductCatalog'
import { useCmsProductDialogs } from './useCmsProductDialogs'

export function useCmsProductPage({ t }) {
  const translate = t || ((key) => key)

  const {
    activeTab,
    areas,
    clusters,
    currentPage,
    fetchAreas,
    fetchProducts,
    filterArea,
    filterClusterId,
    filterKeyword,
    handlePageChange,
    handleResetFilters,
    handleSearch,
    handleSizeChange,
    initialize,
    loading,
    pageSize,
    products,
    setActiveTab,
    total
  } = useCmsProductCatalog({ t: translate })

  const {
    canSubmit,
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
    priceForm,
    productForm,
    resetPriceForm,
    resourceType,
    selectNode,
    showPriceDialog,
    showProductDialog,
    storageRules,
    submitting
  } = useCmsProductDialogs({
    activeTab,
    clusters,
    fetchAreas,
    fetchProducts,
    t: translate
  })

  const handleDelete = async (row) => {
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
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error(translate('failed'))
      }
    }
  }

  return {
    activeTab,
    areas,
    canSubmit,
    clusterNodes,
    clusters,
    computeRules,
    currentPage,
    dialogTitle,
    fetchProducts,
    filterArea,
    filterClusterId,
    filterKeyword,
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
    pageSize,
    priceForm,
    products,
    productForm,
    resetPriceForm,
    resourceType,
    selectNode,
    setActiveTab,
    showPriceDialog,
    showProductDialog,
    storageRules,
    submitting,
    total
  }
}
