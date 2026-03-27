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

<script setup>
import { computed } from 'vue'
import BillingSection from '@/components/createPage/BillingSection.vue'
import FilterChipSection from '@/components/createPage/FilterChipSection.vue'
import ProductTable from './ProductTable.vue'

const props = defineProps({
  areas: {
    type: Array,
    default: () => []
  },
  changeFilter: {
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
  gpuModelsList: {
    type: Array,
    default: () => []
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
  }
})

const emit = defineEmits(['update:pay-type', 'update:product-id'])

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
