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

<script setup lang="ts">
import { computed } from 'vue'
import type { ConsoleProduct, FilterOption } from '@/types/consoleResource'
import BillingSection from '@/components/createPage/BillingSection.vue'
import FilterChipSection from '@/components/createPage/FilterChipSection.vue'
import ProductTable from './ProductTable.vue'
import StorageHintSection from './StorageHintSection.vue'
import SummarySection from './SummarySection.vue'
import type {
  TrainingCreateForm,
  TrainingFilters,
  TrainingPayType,
  TrainingTotalResources
} from '../composables/useTrainingCreate'

interface PayTypeOption {
  value: TrainingPayType
  labelKey: string
}

const props = withDefaults(
  defineProps<{
    areas?: string[]
    changeFilter: (key: keyof TrainingFilters, value: string) => void
    cpuModels?: FilterOption[]
    filters: TrainingFilters
    form: TrainingCreateForm
    formatPrice: (product: ConsoleProduct | null | undefined) => string
    gpuModels?: FilterOption[]
    payTypes?: PayTypeOption[]
    priceUnitText: string
    products?: ConsoleProduct[]
    selectProduct: (product: ConsoleProduct) => void
    selectedProduct?: ConsoleProduct | null
    totalPrice: string
    totalResources?: TrainingTotalResources | null
  }>(),
  {
    areas: () => [],
    cpuModels: () => [],
    gpuModels: () => [],
    payTypes: () => [],
    products: () => [],
    selectedProduct: null,
    totalResources: null
  }
)

const filterGroups = computed(() => [
  { key: 'area', labelKey: 'selectRegion', options: props.areas },
  {
    key: 'gpuModel',
    labelKey: 'selectGpuModel',
    options: props.gpuModels,
    optionMetaKey: 'meta',
    optionMetaFieldsKey: 'metaFields',
    optionValueKey: 'key',
    optionLabelKey: 'label',
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
