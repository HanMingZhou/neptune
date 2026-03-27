<template>
  <ProductToolbar
    :active-tab="activeTab"
    :areas="areas"
    :clusters="clusters"
    :filter-area="filterArea"
    :filter-cluster-id="filterClusterId"
    :filter-keyword="filterKeyword"
    :loading="loading"
    @create="emit('create')"
    @refresh="emit('refresh')"
    @reset="emit('reset')"
    @search="emit('search')"
    @update:active-tab="emit('update:active-tab', $event)"
    @update:filter-area="emit('update:filter-area', $event)"
    @update:filter-cluster-id="emit('update:filter-cluster-id', $event)"
    @update:filter-keyword="emit('update:filter-keyword', $event)"
  />

  <ComputeProductsTable
    v-if="activeTab === 1"
    :items="products"
    :loading="loading"
    :page="currentPage"
    :page-size="pageSize"
    :total="total"
    @adjust-price="emit('adjust-price', $event)"
    @delete="emit('delete', $event)"
    @edit="emit('edit', $event)"
    @page-change="emit('page-change', $event)"
    @size-change="emit('size-change', $event)"
    @update:page="emit('update:page', $event)"
    @update:page-size="emit('update:page-size', $event)"
  />

  <StorageProductsTable
    v-else
    :items="products"
    :loading="loading"
    :page="currentPage"
    :page-size="pageSize"
    :total="total"
    @delete="emit('delete', $event)"
    @edit="emit('edit', $event)"
    @page-change="emit('page-change', $event)"
    @size-change="emit('size-change', $event)"
    @update:page="emit('update:page', $event)"
    @update:page-size="emit('update:page-size', $event)"
  />
</template>

<script setup>
import ComputeProductsTable from './ComputeProductsTable.vue'
import ProductToolbar from './ProductToolbar.vue'
import StorageProductsTable from './StorageProductsTable.vue'

defineProps({
  activeTab: {
    type: Number,
    default: 1
  },
  areas: {
    type: Array,
    default: () => []
  },
  clusters: {
    type: Array,
    default: () => []
  },
  currentPage: {
    type: Number,
    default: 1
  },
  filterArea: {
    type: [String, Number],
    default: ''
  },
  filterClusterId: {
    type: [String, Number],
    default: ''
  },
  filterKeyword: {
    type: String,
    default: ''
  },
  loading: {
    type: Boolean,
    default: false
  },
  pageSize: {
    type: Number,
    default: 20
  },
  products: {
    type: Array,
    default: () => []
  },
  total: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits([
  'adjust-price',
  'create',
  'delete',
  'edit',
  'page-change',
  'refresh',
  'reset',
  'search',
  'size-change',
  'update:active-tab',
  'update:filter-area',
  'update:filter-cluster-id',
  'update:filter-keyword',
  'update:page',
  'update:page-size'
])
</script>
