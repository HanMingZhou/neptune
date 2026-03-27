<template>
  <div class="space-y-6">
    <BillingSection
      :pay-type="form.payType"
      :pay-types="payTypes"
      @update:pay-type="form.payType = $event"
    />

    <FilterChipSection
      :filters="filters"
      :groups="filterGroups"
      @change="changeFilter($event.key, $event.value)"
    />

    <WorkerCountSection
      :available-capacity="availableCapacity"
      :framework-type="form.frameworkType"
      :max-worker-count="maxWorkerCount"
      :selected-product="selectedProduct"
      :show-worker-count="showWorkerCount"
      :worker-count="form.workerCount"
      @decrease-worker="decreaseWorker"
      @increase-worker="increaseWorker"
      @update:worker-count="form.workerCount = Math.max(2, Math.min($event || 2, maxWorkerCount))"
    />

    <ProductTable
      :format-price="formatPrice"
      :price-unit-text="priceUnitText"
      :products="products"
      :selected-product-id="selectedProduct?.id"
      @select-product="selectProduct"
    />

    <StorageHintSection />

    <SummarySection
      :price-unit-text="priceUnitText"
      :selected-product="selectedProduct"
      :total-price="totalPrice"
      :total-resources="totalResources"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import BillingSection from '@/components/createPage/BillingSection.vue'
import FilterChipSection from '@/components/createPage/FilterChipSection.vue'
import ProductTable from './ProductTable.vue'
import StorageHintSection from './StorageHintSection.vue'
import SummarySection from './SummarySection.vue'
import WorkerCountSection from './WorkerCountSection.vue'

const props = defineProps({
  areas: {
    type: Array,
    default: () => []
  },
  availableCapacity: {
    type: Number,
    default: 0
  },
  changeFilter: {
    type: Function,
    required: true
  },
  cpuModels: {
    type: Array,
    default: () => []
  },
  decreaseWorker: {
    type: Function,
    required: true
  },
  filters: {
    type: Object,
    required: true
  },
  form: {
    type: Object,
    required: true
  },
  formatPrice: {
    type: Function,
    required: true
  },
  gpuModels: {
    type: Array,
    default: () => []
  },
  increaseWorker: {
    type: Function,
    required: true
  },
  maxWorkerCount: {
    type: Number,
    required: true
  },
  payTypes: {
    type: Array,
    default: () => []
  },
  priceUnitText: {
    type: String,
    required: true
  },
  products: {
    type: Array,
    default: () => []
  },
  selectProduct: {
    type: Function,
    required: true
  },
  selectedProduct: {
    type: Object,
    default: null
  },
  showWorkerCount: {
    type: Boolean,
    default: false
  },
  totalPrice: {
    type: String,
    required: true
  },
  totalResources: {
    type: Object,
    default: null
  }
})

const filterGroups = computed(() => [
  { key: 'area', labelKey: 'selectRegion', options: props.areas },
  {
    key: 'gpuModel',
    labelKey: 'selectGpuModel',
    options: props.gpuModels,
    optionValueKey: 'model',
    optionLabelKey: 'model',
    optionCountKey: 'available'
  },
  {
    key: 'cpuModel',
    labelKey: 'selectCpuModel',
    options: props.cpuModels,
    optionValueKey: 'model',
    optionLabelKey: 'model',
    optionCountKey: 'available'
  }
])
</script>
