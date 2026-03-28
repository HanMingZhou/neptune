<template>
  <div class="console-page-container space-y-6">
    <ProductManagementHeader />

    <ProductTableSection
      :active-tab="activeTab"
      :areas="areas"
      :clusters="clusters"
      :current-page="currentPage"
      :filter-area="filterArea"
      :filter-cluster-id="filterClusterId"
      :filter-keyword="filterKeyword"
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
      @update:filter-cluster-id="filterClusterId = $event"
      @update:filter-keyword="filterKeyword = $event"
      @update:page="currentPage = $event"
      @update:page-size="pageSize = $event"
    />

    <ProductDialogsHost
      :can-submit="canSubmit"
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
      :price-form="priceForm"
      :product-form="productForm"
      :resource-type="resourceType"
      :show-price-dialog="showPriceDialog"
      :show-product-dialog="showProductDialog"
      :storage-rules="storageRules"
      :submitting="submitting"
      @cluster-change="handleClusterChange"
      @node-select="selectNode"
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

<script setup>
import { inject, onMounted } from 'vue'
import ProductDialogsHost from './components/ProductDialogsHost.vue'
import ProductManagementHeader from './components/ProductManagementHeader.vue'
import ProductTableSection from './components/ProductTableSection.vue'
import { useCmsProductPage } from './composables/useCmsProductPage'

const t = inject('t', (key) => key)

const {
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
  resourceType,
  selectNode,
  setActiveTab,
  showPriceDialog,
  showProductDialog,
  storageRules,
  submitting,
  total
} = useCmsProductPage({ t })

onMounted(() => {
  initialize()
})
</script>
