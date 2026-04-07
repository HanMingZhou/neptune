<template>
  <ProductDialogs
    :model-value="showProductDialog"
    :price-visible="showPriceDialog"
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
    :storage-rules="storageRules"
    :submitting="submitting"
    @cluster-change="emit('cluster-change', $event)"
    @node-select="emit('node-select', $event)"
    @resource-type-change="emit('resource-type-change', $event)"
    @storage-cluster-change="emit('storage-cluster-change', $event)"
    @submit-price="emit('submit-price')"
    @submit-product="emit('submit-product')"
    @update:model-value="emit('update:product-dialog-visible', $event)"
    @update:price-visible="emit('update:price-dialog-visible', $event)"
    @update:resource-type="emit('update:resource-type', $event)"
  />
</template>

<script setup lang="ts">
import ProductDialogs from './ProductDialogs.vue'
import type { FormRules } from 'element-plus'
import type {
  CmsClusterOption,
  CmsNodeRow,
  CmsProductForm,
  CmsProductPriceForm,
  CmsProductResourceType
} from '@/types/superAdmin'

withDefaults(
  defineProps<{
    canSubmit?: boolean
    clusterNodes?: CmsNodeRow[]
    clusters?: CmsClusterOption[]
    computeRules?: FormRules<CmsProductForm>
    dialogTitle?: string
    isEdit?: boolean
    loadingNodes?: boolean
    nodeMaxCpu?: number
    nodeMaxGpuCount?: number
    nodeMaxGpuMemory?: number
    nodeMaxMemory?: number
    nodeMaxVGpuCores?: number
    nodeMaxVGpuCount?: number
    nodeMaxVGpuMemory?: number
    priceForm: CmsProductPriceForm
    productForm: CmsProductForm
    resourceType?: CmsProductResourceType
    showPriceDialog?: boolean
    showProductDialog?: boolean
    storageRules?: FormRules<CmsProductForm>
    submitting?: boolean
  }>(),
  {
    canSubmit: false,
    clusterNodes: () => [],
    clusters: () => [],
    computeRules: () => ({}),
    dialogTitle: '',
    isEdit: false,
    loadingNodes: false,
    nodeMaxCpu: 0,
    nodeMaxGpuCount: 0,
    nodeMaxGpuMemory: 0,
    nodeMaxMemory: 0,
    nodeMaxVGpuCores: 0,
    nodeMaxVGpuCount: 0,
    nodeMaxVGpuMemory: 0,
    resourceType: 'cpu',
    showPriceDialog: false,
    showProductDialog: false,
    storageRules: () => ({}),
    submitting: false
  }
)

const emit = defineEmits<{
  'cluster-change': [clusterId: CmsProductForm['clusterId']]
  'node-select': [node: CmsNodeRow]
  'resource-type-change': [type: CmsProductResourceType]
  'storage-cluster-change': [clusterId: CmsProductForm['clusterId']]
  'submit-price': []
  'submit-product': []
  'update:price-dialog-visible': [value: boolean]
  'update:product-dialog-visible': [value: boolean]
  'update:resource-type': [value: CmsProductResourceType]
}>()
</script>
