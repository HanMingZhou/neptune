<template>
  <div class="space-y-6">
    <BillingSection
      :pay-type="payType"
      :pay-types="payTypes"
      @update:pay-type="emit('update:pay-type', $event)"
    />

    <FilterChipSection
      :filters="filters"
      :groups="filterGroups"
      @change="changeFilter($event.key, $event.value)"
    />

    <GpuCountSection
      :model-value="gpuCount"
      :selected-product="selectedProduct"
      @update:model-value="emit('update:gpu-count', $event)"
    />

    <ProductTable
      :format-price="formatPrice"
      :price-unit-text="priceUnitText"
      :products="products"
      :selected-product-id="selectedProduct?.id"
      @select-product="selectProduct"
    />

    <VolumeMountSection
      :available-volumes="availableVolumes"
      :selected-volume-id="selectedVolumeId"
      :volume-mount-path="volumeMountPath"
      @update:selected-volume-id="emit('update:selected-volume-id', $event)"
      @update:volume-mount-path="emit('update:volume-mount-path', $event)"
      @volume-change="onVolumeChange"
    />

    <SummarySection
      :gpu-count="gpuCount"
      :price-unit-text="priceUnitText"
      :selected-product="selectedProduct"
      :selected-volume-id="selectedVolumeId"
      :selected-volume-name="selectedVolumeName"
      :total-price="totalPrice"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import BillingSection from '@/components/createPage/BillingSection.vue'
import FilterChipSection from '@/components/createPage/FilterChipSection.vue'
import GpuCountSection from './GpuCountSection.vue'
import ProductTable from './ProductTable.vue'
import SummarySection from './SummarySection.vue'
import VolumeMountSection from './VolumeMountSection.vue'

const props = defineProps({
  areas: {
    type: Array,
    default: () => []
  },
  availableVolumes: {
    type: Array,
    default: () => []
  },
  changeFilter: {
    type: Function,
    required: true
  },
  cpuModels: {
    type: Array,
    default: () => []
  },
  filters: {
    type: Object,
    required: true
  },
  formatPrice: {
    type: Function,
    required: true
  },
  gpuCount: {
    type: Number,
    required: true
  },
  gpuModels: {
    type: Array,
    default: () => []
  },
  onVolumeChange: {
    type: Function,
    required: true
  },
  payType: {
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
  selectedProduct: {
    type: Object,
    default: null
  },
  selectedVolumeId: {
    type: [Number, String],
    default: null
  },
  selectedVolumeName: {
    type: String,
    default: ''
  },
  selectProduct: {
    type: Function,
    required: true
  },
  totalPrice: {
    type: String,
    required: true
  },
  volumeMountPath: {
    type: String,
    default: ''
  }
})

const emit = defineEmits([
  'update:gpu-count',
  'update:pay-type',
  'update:selected-volume-id',
  'update:volume-mount-path'
])

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
