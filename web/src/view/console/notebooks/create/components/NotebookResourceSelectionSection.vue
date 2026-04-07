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

<script setup lang="ts">
import { computed } from 'vue'
import type {
  ConsoleProduct,
  ConsoleVolume,
  FilterOption,
  ResourceId
} from '@/types/consoleResource'
import BillingSection from '@/components/createPage/BillingSection.vue'
import FilterChipSection from '@/components/createPage/FilterChipSection.vue'
import GpuCountSection from './GpuCountSection.vue'
import ProductTable from './ProductTable.vue'
import SummarySection from './SummarySection.vue'
import VolumeMountSection from './VolumeMountSection.vue'
import type {
  NotebookFilters,
  NotebookPayType
} from '../composables/useNotebookCreate'

interface PayTypeOption {
  value: NotebookPayType
  labelKey: string
}

const props = withDefaults(
  defineProps<{
    areas?: string[]
    availableVolumes?: ConsoleVolume[]
    changeFilter: (key: keyof NotebookFilters, value: string) => void
    cpuModels?: FilterOption[]
    filters: NotebookFilters
    formatPrice: (product: ConsoleProduct | null | undefined) => string
    gpuCount: number
    gpuModels?: FilterOption[]
    onVolumeChange: (value: ResourceId | null) => void
    payType: NotebookPayType
    payTypes?: PayTypeOption[]
    priceUnitText: string
    products?: ConsoleProduct[]
    selectedProduct?: ConsoleProduct | null
    selectedVolumeId?: ResourceId | null
    selectedVolumeName?: string
    selectProduct: (product: ConsoleProduct) => void
    totalPrice: string
    volumeMountPath?: string
  }>(),
  {
    areas: () => [],
    availableVolumes: () => [],
    cpuModels: () => [],
    gpuModels: () => [],
    payTypes: () => [],
    products: () => [],
    selectedProduct: null,
    selectedVolumeId: null,
    selectedVolumeName: '',
    volumeMountPath: ''
  }
)

const emit = defineEmits<{
  'update:gpu-count': [value: number]
  'update:pay-type': [value: NotebookPayType]
  'update:selected-volume-id': [value: ResourceId | null]
  'update:volume-mount-path': [value: string]
}>()

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
