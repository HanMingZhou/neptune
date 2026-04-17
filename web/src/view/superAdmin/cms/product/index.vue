<template>
  <div class="console-page-container flex min-h-full flex-col gap-6">
    <ProductManagementHeader />

    <ProductTableSection
      class="min-h-0 flex-1"
      :active-tab="activeTab"
      :areas="areas"
      :clusters="clusters"
      :current-page="currentPage"
      :filter-area="filterArea"
      :filter-available-max="filterAvailableMax"
      :filter-available-min="filterAvailableMin"
      :filter-cluster-id="filterClusterId"
      :filter-resource-type="filterResourceType"
      :filter-gpu-model="filterGpuModel"
      :filter-keyword="filterKeyword"
      :filter-max-instances-max="filterMaxInstancesMax"
      :filter-max-instances-min="filterMaxInstancesMin"
      :filter-price-field="filterPriceField"
      :filter-price-max="filterPriceMax"
      :filter-price-min="filterPriceMin"
      :filter-used-capacity-max="filterUsedCapacityMax"
      :filter-used-capacity-min="filterUsedCapacityMin"
      :gpu-models="gpuModels"
      :loading="loading"
      :page-size="pageSize"
      :products="products"
      :total="total"
      @adjust-price="openPriceDialog"
      @create="openCreateDialog"
      @delete="handleDelete"
      @edit="openEditDialog"
      @page-change="handlePageChange"
      @refresh="fetchProducts"
      @reset="handleResetFilters"
      @search="handleSearch"
      @size-change="handleSizeChange"
      @update:active-tab="setActiveTab"
      @update:filter-area="filterArea = $event"
      @update:filter-available-max="filterAvailableMax = $event"
      @update:filter-available-min="filterAvailableMin = $event"
      @update:filter-cluster-id="filterClusterId = $event"
      @update:filter-resource-type="filterResourceType = $event"
      @update:filter-gpu-model="filterGpuModel = $event"
      @update:filter-keyword="filterKeyword = $event"
      @update:filter-max-instances-max="filterMaxInstancesMax = $event"
      @update:filter-max-instances-min="filterMaxInstancesMin = $event"
      @update:filter-price-field="filterPriceField = $event"
      @update:filter-price-max="filterPriceMax = $event"
      @update:filter-price-min="filterPriceMin = $event"
      @update:filter-used-capacity-max="filterUsedCapacityMax = $event"
      @update:filter-used-capacity-min="filterUsedCapacityMin = $event"
      @update:page="currentPage = $event"
      @update:page-size="pageSize = $event"
    />

    <ProductDialogsHost
      :active-preview-node-name="activePreviewNodeName"
      :can-submit="canSubmit"
      :selected-node-label="selectedNodeLabel"
      :cluster-nodes="clusterNodes"
      :clusters="clusters"
      :compute-rules="computeRules"
      :dialog-title="dialogTitle"
      :is-edit="isEdit"
      :loading-nodes="loadingNodes"
      :node-max-cpu="nodeMaxCpu"
      :node-max-gpu-count="nodeMaxGpuCount"
      :node-max-gpu-memory="nodeMaxGpuMemory"
      :node-max-memory="nodeMaxMemory"
      :node-max-v-gpu-cores="nodeMaxVGpuCores"
      :node-max-v-gpu-count="nodeMaxVGpuCount"
      :node-max-v-gpu-memory="nodeMaxVGpuMemory"
      :preview-node="previewNode"
      :price-form="priceForm"
      :product-form="productForm"
      :resource-type="resourceType"
      :selected-node-count="selectedNodeCount"
      :selected-node-names="selectedNodeNames"
      :show-price-dialog="showPriceDialog"
      :show-product-dialog="showProductDialog"
      :storage-rules="storageRules"
      :submit-button-text="submitButtonText"
      :submitting="submitting"
      @cluster-change="handleClusterChange"
      @node-preview="previewClusterNode"
      @node-toggle="toggleNodeSelection"
      @selection-clear="clearSelectedNodes"
      @resource-type-change="handleResourceTypeChange"
      @storage-cluster-change="handleStorageClusterChange"
      @submit-price="handleUpdatePrice"
      @submit-product="handleSubmitProduct"
      @update:price-dialog-visible="showPriceDialog = $event"
      @update:product-dialog-visible="showProductDialog = $event"
      @update:resource-type="resourceType = $event"
    />
  </div>
</template>

<script setup lang="ts">
import { inject, onMounted } from 'vue'
import ProductDialogsHost from './components/ProductDialogsHost.vue'
import ProductManagementHeader from './components/ProductManagementHeader.vue'
import ProductTableSection from './components/ProductTableSection.vue'
import { useCmsProductPage } from './composables/useCmsProductPage'

const t = inject('t', (key) => key)

const {
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
  gpuModels,
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
  total
} = useCmsProductPage({ t })

onMounted(() => {
  void initialize()
})
</script>
