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

<script setup lang="ts">
import ApiTableCard from './ApiTableCard.vue'
import type { ResourceId } from '@/types/consoleResource'
import type {
  ApiListItem,
  ApiMethodOption,
  LabelValueOption
} from '@/types/superAdmin'

withDefaults(
  defineProps<{
    allSelected?: boolean
    apiGroupOptions?: LabelValueOption[]
    getMethodClass: (method?: string) => string
    items?: ApiListItem[]
    loading?: boolean
    methodOptions?: ApiMethodOption[]
    page?: number
    pageSize?: number
    searchApiGroup?: string
    searchDescription?: string
    searchMethod?: string
    searchPath?: string
    selectedIds?: ResourceId[]
    total?: number
  }>(),
  {
    allSelected: false,
    apiGroupOptions: () => [],
    items: () => [],
    loading: false,
    methodOptions: () => [],
    page: 1,
    pageSize: 15,
    searchApiGroup: '',
    searchDescription: '',
    searchMethod: '',
    searchPath: '',
    selectedIds: () => [],
    total: 0
  }
)

const emit = defineEmits<{
  delete: [row: ApiListItem]
  edit: [row: ApiListItem]
  'page-change': [page: number]
  reset: []
  search: []
  'size-change': [pageSize: number]
  'toggle-select': [row: ApiListItem]
  'toggle-select-all': [checked: boolean]
  'update:search-api-group': [value: string]
  'update:search-description': [value: string]
  'update:search-method': [value: string]
  'update:search-path': [value: string]
}>()
</script>

