<template>
  <ManagementListShell>
    <template #filters>
      <ProductToolbar
        :active-tab="activeTab"
        :areas="areas"
        :clusters="clusters"
        :filter-area="filterArea"
        :filter-cluster-id="filterClusterId"
        :filter-keyword="filterKeyword"
        :loading="loading"
        @create="emit('create')"
        @refresh="emit('refresh', $event)"
        @reset="emit('reset')"
        @search="emit('search')"
        @update:active-tab="emit('update:active-tab', $event)"
        @update:filter-area="emit('update:filter-area', $event)"
        @update:filter-cluster-id="emit('update:filter-cluster-id', $event)"
        @update:filter-keyword="emit('update:filter-keyword', $event)"
      />
    </template>

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
  </ManagementListShell>
</template>

<script setup lang="ts">
import ManagementListShell from '@/components/listPage/ManagementListShell.vue'
import ComputeProductsTable from './ComputeProductsTable.vue'
import ProductToolbar from './ProductToolbar.vue'
import StorageProductsTable from './StorageProductsTable.vue'
import type { ResourceId } from '@/types/consoleResource'
import type { CmsClusterOption, CmsProductRow } from '@/types/superAdmin'

withDefaults(
  defineProps<{
    activeTab?: number
    areas?: string[]
    clusters?: CmsClusterOption[]
    currentPage?: number
    filterArea?: ResourceId | ''
    filterClusterId?: ResourceId | ''
    filterKeyword?: string
    loading?: boolean
    pageSize?: number
    products?: CmsProductRow[]
    total?: number
  }>(),
  {
    activeTab: 1,
    areas: () => [],
    clusters: () => [],
    currentPage: 1,
    filterArea: '',
    filterClusterId: '',
    filterKeyword: '',
    loading: false,
    pageSize: 20,
    products: () => [],
    total: 0
  }
)

const emit = defineEmits<{
  'adjust-price': [row: CmsProductRow]
  create: []
  delete: [row: CmsProductRow]
  edit: [row: CmsProductRow]
  'page-change': [page: number]
  refresh: [silent?: boolean]
  reset: []
  search: []
  'size-change': [pageSize: number]
  'update:active-tab': [value: number]
  'update:filter-area': [value: ResourceId | '']
  'update:filter-cluster-id': [value: ResourceId | '']
  'update:filter-keyword': [value: string]
  'update:page': [value: number]
  'update:page-size': [value: number]
}>()
</script>
