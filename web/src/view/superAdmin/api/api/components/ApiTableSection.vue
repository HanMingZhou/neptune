<template>
  <ApiTableCard
    :all-selected="allSelected"
    :api-group-options="apiGroupOptions"
    :get-method-class="getMethodClass"
    :items="items"
    :loading="loading"
    :method-options="methodOptions"
    :page="page"
    :page-size="pageSize"
    :search-api-group="searchApiGroup"
    :search-description="searchDescription"
    :search-method="searchMethod"
    :search-path="searchPath"
    :selected-ids="selectedIds"
    :total="total"
    @delete="emit('delete', $event)"
    @edit="emit('edit', $event)"
    @page-change="emit('page-change', $event)"
    @reset="emit('reset')"
    @search="emit('search')"
    @size-change="emit('size-change', $event)"
    @toggle-select="emit('toggle-select', $event)"
    @toggle-select-all="emit('toggle-select-all', $event)"
    @update:search-api-group="emit('update:search-api-group', $event)"
    @update:search-description="emit('update:search-description', $event)"
    @update:search-method="emit('update:search-method', $event)"
    @update:search-path="emit('update:search-path', $event)"
  />
</template>

<script setup>
import ApiTableCard from './ApiTableCard.vue'

defineProps({
  allSelected: {
    type: Boolean,
    default: false
  },
  apiGroupOptions: {
    type: Array,
    default: () => []
  },
  getMethodClass: {
    type: Function,
    required: true
  },
  items: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  methodOptions: {
    type: Array,
    default: () => []
  },
  page: {
    type: Number,
    default: 1
  },
  pageSize: {
    type: Number,
    default: 10
  },
  searchApiGroup: {
    type: String,
    default: ''
  },
  searchDescription: {
    type: String,
    default: ''
  },
  searchMethod: {
    type: String,
    default: ''
  },
  searchPath: {
    type: String,
    default: ''
  },
  selectedIds: {
    type: Array,
    default: () => []
  },
  total: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits([
  'delete',
  'edit',
  'page-change',
  'reset',
  'search',
  'size-change',
  'toggle-select',
  'toggle-select-all',
  'update:search-api-group',
  'update:search-description',
  'update:search-method',
  'update:search-path'
])
</script>
