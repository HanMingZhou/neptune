<template>
  <div class="p-4 border-b border-slate-100 dark:border-border-dark flex flex-wrap items-center gap-4">
    <div class="relative max-w-md">
      <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-sm">search</span>
      <input
        v-model="searchQueryModel"
        type="text"
        :placeholder="`${t('order.productService')}...`"
        class="list-search-input !w-full"
        @keyup.enter="$emit('search')"
      />
    </div>

    <div class="list-tab-group">
      <button @click="chargeFilterModel = 'all'" :class="{ 'is-active': chargeFilter === 'all' }">{{ t('all') }}</button>
      <button @click="chargeFilterModel = '按量计费'" :class="{ 'is-active': chargeFilter === '按量计费' }">{{ t('payHourly') }}</button>
      <button @click="chargeFilterModel = '预付费'" :class="{ 'is-active': chargeFilter === '预付费' }">{{ t('order.prepaid') }}</button>
    </div>

    <OrderDateRangePicker
      v-model="dateRangeModel"
      :start-placeholder="t('startDate')"
      :end-placeholder="t('endDate')"
      @change="$emit('date-change')"
    />
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import OrderDateRangePicker from '@/components/orderPage/OrderDateRangePicker.vue'

const props = defineProps({
  chargeFilter: {
    type: String,
    default: 'all'
  },
  dateRange: {
    type: Array,
    default: () => []
  },
  searchQuery: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['date-change', 'search', 'update:charge-filter', 'update:date-range', 'update:search-query'])
const t = inject('t', (key) => key)

const searchQueryModel = computed({
  get: () => props.searchQuery,
  set: (value) => emit('update:search-query', value)
})

const chargeFilterModel = computed({
  get: () => props.chargeFilter,
  set: (value) => emit('update:charge-filter', value)
})

const dateRangeModel = computed({
  get: () => props.dateRange,
  set: (value) => emit('update:date-range', value)
})
</script>
