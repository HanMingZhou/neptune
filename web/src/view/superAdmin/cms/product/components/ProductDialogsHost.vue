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

<script setup>
import ProductDialogs from './ProductDialogs.vue'

defineProps({
  canSubmit: {
    type: Boolean,
    default: false
  },
  clusterNodes: {
    type: Array,
    default: () => []
  },
  clusters: {
    type: Array,
    default: () => []
  },
  computeRules: {
    type: Object,
    default: () => ({})
  },
  dialogTitle: {
    type: String,
    default: ''
  },
  isEdit: {
    type: Boolean,
    default: false
  },
  loadingNodes: {
    type: Boolean,
    default: false
  },
  nodeMaxCpu: {
    type: Number,
    default: 0
  },
  nodeMaxGpuCount: {
    type: Number,
    default: 0
  },
  nodeMaxGpuMemory: {
    type: Number,
    default: 0
  },
  nodeMaxMemory: {
    type: Number,
    default: 0
  },
  nodeMaxVGpuCores: {
    type: Number,
    default: 0
  },
  nodeMaxVGpuCount: {
    type: Number,
    default: 0
  },
  nodeMaxVGpuMemory: {
    type: Number,
    default: 0
  },
  priceForm: {
    type: Object,
    required: true
  },
  productForm: {
    type: Object,
    required: true
  },
  resourceType: {
    type: String,
    default: 'cpu'
  },
  showPriceDialog: {
    type: Boolean,
    default: false
  },
  showProductDialog: {
    type: Boolean,
    default: false
  },
  storageRules: {
    type: Object,
    default: () => ({})
  },
  submitting: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits([
  'cluster-change',
  'node-select',
  'resource-type-change',
  'storage-cluster-change',
  'submit-price',
  'submit-product',
  'update:price-dialog-visible',
  'update:product-dialog-visible',
  'update:resource-type'
])
</script>
