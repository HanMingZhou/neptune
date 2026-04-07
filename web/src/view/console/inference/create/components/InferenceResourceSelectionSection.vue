<template>
  <div class="space-y-6">
    <BillingSection
      :pay-type="form.payType"
      :pay-types="payTypes"
      @update:pay-type="emit('update:pay-type', $event)"
    />

    <FilterChipSection
      :filters="filters"
      :groups="filterGroups"
      @change="changeFilter($event.key, $event.value)"
    />

    <ProductTable
      :format-price="formatPrice"
      :price-unit-text="priceUnitText"
      :product-id="form.productId"
      :products="products"
      @update:product-id="emit('update:product-id', $event)"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type {
  ConsoleProduct,
  FilterOption,
  ResourceId
} from '@/types/consoleResource'
import BillingSection from '@/components/createPage/BillingSection.vue'
import FilterChipSection from '@/components/createPage/FilterChipSection.vue'
import ProductTable from './ProductTable.vue'
import type {
  InferenceCreateForm,
  InferenceFilters,
  InferencePayType
} from '../composables/useInferenceCreate'

interface PayTypeOption {
  value: InferencePayType
  labelKey: string
}

const props = withDefaults(
  defineProps<{
    areas?: string[]
    changeFilter: (key: keyof InferenceFilters, value: string) => void
    filters: InferenceFilters
    form: InferenceCreateForm
    formatPrice: (product: ConsoleProduct | null | undefined) => string
    gpuModelsList?: FilterOption[]
    payTypes?: PayTypeOption[]
    priceUnitText: string
    products?: ConsoleProduct[]
  }>(),
  {
    areas: () => [],
    gpuModelsList: () => [],
    payTypes: () => [],
    products: () => []
  }
)

const emit = defineEmits<{
  'update:pay-type': [value: InferencePayType]
  'update:product-id': [value: ResourceId]
}>()

const filterGroups = computed(() => [
  { key: 'area', labelKey: 'selectRegion', options: props.areas },
  {
    key: 'gpuModel',
    labelKey: 'selectGpuModel',
    options: props.gpuModelsList,
    optionValueKey: 'model',
    optionLabelKey: 'model',
    optionCountKey: 'available'
  }
])
</script>
