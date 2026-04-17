<template>
  <ManagementListShell class="min-h-0 flex-1">
    <template #filters>
      <ProductToolbar
        :active-tab="activeTab"
        :areas="areas"
        :clusters="clusters"
        :filter-area="filterArea"
        :filter-available-max="filterAvailableMax"
        :filter-available-min="filterAvailableMin"
        :filter-cluster-id="filterClusterId"
        :filter-resource-type="filterResourceType"
        :filter-gpu-model="filterGpuModel"
        :filter-keyword="filterKeyword"
        :filter-max-instances-max="filterMaxInstancesMax"
        :filter-max-instances-min="filterMaxInstancesMin"
        :filter-price-field="filterPriceField"
        :filter-price-max="filterPriceMax"
        :filter-price-min="filterPriceMin"
        :filter-used-capacity-max="filterUsedCapacityMax"
        :filter-used-capacity-min="filterUsedCapacityMin"
        :gpu-models="gpuModels"
        :loading="loading"
        @create="emit('create')"
        @refresh="emit('refresh', $event)"
        @reset="emit('reset')"
        @search="emit('search')"
        @update:active-tab="emit('update:active-tab', $event)"
        @update:filter-area="emit('update:filter-area', $event)"
        @update:filter-available-max="emit('update:filter-available-max', $event)"
        @update:filter-available-min="emit('update:filter-available-min', $event)"
        @update:filter-cluster-id="emit('update:filter-cluster-id', $event)"
        @update:filter-resource-type="emit('update:filter-resource-type', $event)"
        @update:filter-gpu-model="emit('update:filter-gpu-model', $event)"
        @update:filter-keyword="emit('update:filter-keyword', $event)"
        @update:filter-max-instances-max="
          emit('update:filter-max-instances-max', $event)
        "
        @update:filter-max-instances-min="
          emit('update:filter-max-instances-min', $event)
        "
        @update:filter-price-field="emit('update:filter-price-field', $event)"
        @update:filter-price-max="emit('update:filter-price-max', $event)"
        @update:filter-price-min="emit('update:filter-price-min', $event)"
        @update:filter-used-capacity-max="
          emit('update:filter-used-capacity-max', $event)
        "
        @update:filter-used-capacity-min="
          emit('update:filter-used-capacity-min', $event)
        "
      />
    </template>

    <ComputeProductsTable
      v-if="activeTab === 1"
      class="min-h-0 flex-1"
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
      class="min-h-0 flex-1"
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
import type {
  CmsClusterOption,
  CmsProductFilterResourceType,
  CmsProductPriceField,
  CmsProductRow
} from '@/types/superAdmin'
import type { FilterOption } from '@/types/consoleResource'

withDefaults(
  defineProps<{
    activeTab?: number
    areas?: string[]
    clusters?: CmsClusterOption[]
    currentPage?: number
    filterArea?: ResourceId | ''
    filterAvailableMax?: number
    filterAvailableMin?: number
    filterClusterId?: ResourceId | ''
    filterResourceType?: CmsProductFilterResourceType
    filterGpuModel?: string
    filterKeyword?: string
    filterMaxInstancesMax?: number
    filterMaxInstancesMin?: number
    filterPriceField?: CmsProductPriceField
    filterPriceMax?: number
    filterPriceMin?: number
    filterUsedCapacityMax?: number
    filterUsedCapacityMin?: number
    gpuModels?: FilterOption[]
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
    filterAvailableMax: undefined,
    filterAvailableMin: undefined,
    filterClusterId: '',
    filterResourceType: '',
    filterGpuModel: '',
    filterKeyword: '',
    filterMaxInstancesMax: undefined,
    filterMaxInstancesMin: undefined,
    filterPriceField: 1,
    filterPriceMax: undefined,
    filterPriceMin: undefined,
    filterUsedCapacityMax: undefined,
    filterUsedCapacityMin: undefined,
    gpuModels: () => [],
    loading: false,
    pageSize: 15,
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
  'update:filter-available-max': [value?: number]
  'update:filter-available-min': [value?: number]
  'update:filter-cluster-id': [value: ResourceId | '']
  'update:filter-resource-type': [value: CmsProductFilterResourceType]
  'update:filter-gpu-model': [value: string]
  'update:filter-keyword': [value: string]
  'update:filter-max-instances-max': [value?: number]
  'update:filter-max-instances-min': [value?: number]
  'update:filter-price-field': [value: CmsProductPriceField]
  'update:filter-price-max': [value?: number]
  'update:filter-price-min': [value?: number]
  'update:filter-used-capacity-max': [value?: number]
  'update:filter-used-capacity-min': [value?: number]
  'update:page': [value: number]
  'update:page-size': [value: number]
}>()
</script>

