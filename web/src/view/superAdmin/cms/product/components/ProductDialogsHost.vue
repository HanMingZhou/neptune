<template>
  <ProductDialogs
    :active-preview-node-name="activePreviewNodeName"
    :model-value="showProductDialog"
    :price-visible="showPriceDialog"
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
    :storage-rules="storageRules"
    :submit-button-text="submitButtonText"
    :submitting="submitting"
    @cluster-change="emit('cluster-change', $event)"
    @node-preview="emit('node-preview', $event)"
    @node-toggle="emit('node-toggle', $event)"
    @selection-clear="emit('selection-clear')"
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
  CmsProductNodeCandidate,
  CmsProductForm,
  CmsProductPriceForm,
  CmsProductResourceType
} from '@/types/superAdmin'

withDefaults(
  defineProps<{
    activePreviewNodeName?: string
    canSubmit?: boolean
    clusterNodes?: CmsProductNodeCandidate[]
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
    previewNode?: CmsNodeRow | CmsProductNodeCandidate | null
    priceForm: CmsProductPriceForm
    productForm: CmsProductForm
    resourceType?: CmsProductResourceType
    selectedNodeCount?: number
    selectedNodeLabel?: string
    selectedNodeNames?: string[]
    showPriceDialog?: boolean
    showProductDialog?: boolean
    storageRules?: FormRules<CmsProductForm>
    submitButtonText?: string
    submitting?: boolean
  }>(),
  {
    activePreviewNodeName: '',
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
    previewNode: null,
    resourceType: 'cpu',
    selectedNodeCount: 0,
    selectedNodeLabel: '',
    selectedNodeNames: () => [],
    showPriceDialog: false,
    showProductDialog: false,
    storageRules: () => ({}),
    submitButtonText: '',
    submitting: false
  }
)

const emit = defineEmits<{
  'cluster-change': [clusterId: CmsProductForm['clusterId']]
  'node-preview': [node: CmsProductNodeCandidate]
  'node-toggle': [node: CmsProductNodeCandidate]
  'selection-clear': []
  'resource-type-change': [type: CmsProductResourceType]
  'storage-cluster-change': [clusterId: CmsProductForm['clusterId']]
  'submit-price': []
  'submit-product': []
  'update:price-dialog-visible': [value: boolean]
  'update:product-dialog-visible': [value: boolean]
  'update:resource-type': [value: CmsProductResourceType]
}>()
</script>
