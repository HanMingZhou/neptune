<template>
  <div class="p-4 border-b border-gray-100 dark:border-border-dark flex flex-wrap gap-4 items-center">
    <div class="flex items-center gap-3 flex-1">
      <div class="relative">
        <input
          v-model="searchQueryModel"
          type="text"
          :placeholder="`${t('transactionId')}...`"
          class="list-search-input"
          @keyup.enter="$emit('search')"
        />
        <span class="material-icons absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400 text-[16px]">search</span>
      </div>

      <div class="list-tab-group">
        <button @click="filterTypeModel = 'all'" :class="{ 'is-active': filterType === 'all' }">{{ t('allTypes') }}</button>
        <button @click="filterTypeModel = 'Recharge'" :class="{ 'is-active': filterType === 'Recharge' }">{{ t('recharge') }}</button>
        <button @click="filterTypeModel = 'Consumption'" :class="{ 'is-active': filterType === 'Consumption' }">{{ t('consumption') }}</button>
      </div>

      <OrderDateRangePicker
        v-model="dateRangeModel"
        :range-separator="t('to')"
        :start-placeholder="t('startDate')"
        :end-placeholder="t('endDate')"
        @change="$emit('date-change')"
      />
    </div>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import OrderDateRangePicker from '@/components/orderPage/OrderDateRangePicker.vue'

const props = defineProps({
  dateRange: {
    type: Array,
    default: () => []
  },
  filterType: {
    type: String,
    default: 'all'
  },
  searchQuery: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['date-change', 'search', 'update:date-range', 'update:filter-type', 'update:search-query'])
const t = inject('t', (key) => key)

const searchQueryModel = computed({
  get: () => props.searchQuery,
  set: (value) => emit('update:search-query', value)
})

const filterTypeModel = computed({
  get: () => props.filterType,
  set: (value) => emit('update:filter-type', value)
})

const dateRangeModel = computed({
  get: () => props.dateRange,
  set: (value) => emit('update:date-range', value)
})
</script>
