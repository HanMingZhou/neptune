<template>
  <TableCard>
    <template #toolbar>
      <BaseFilterBar plain>
        <div class="list-filter-field list-filter-field--compact max-w-[160px]">
          <span
            class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]"
            >link</span
          >
          <input
            v-model="searchPathModel"
            type="text"
            :placeholder="t('path')"
            class="list-search-input !w-full"
            @keyup.enter="$emit('search')"
          />
        </div>

        <div class="list-filter-field list-filter-field--compact max-w-[160px]">
          <span
            class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]"
            >description</span
          >
          <input
            v-model="searchDescriptionModel"
            type="text"
            :placeholder="t('apiDesc')"
            class="list-search-input !w-full"
            @keyup.enter="$emit('search')"
          />
        </div>

        <el-select
          v-model="searchApiGroupModel"
          clearable
          :placeholder="t('apiGroup')"
          class="list-filter-select !w-[160px]"
        >
          <el-option
            v-for="item in apiGroupOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>

        <el-select
          v-model="searchMethodModel"
          clearable
          :placeholder="t('method')"
          class="list-filter-select !w-[140px]"
        >
          <el-option
            v-for="item in methodOptions"
            :key="item.value"
            :label="`${item.label}(${item.value})`"
            :value="item.value"
          />
        </el-select>

        <template #actions>
          <button
            class="list-toolbar-button list-toolbar-button--primary"
            @click="$emit('search')"
          >
            <span class="material-icons text-[18px]">search</span>
            {{ t('searchQuery') }}
          </button>
          <button
            class="list-toolbar-button list-toolbar-button--secondary"
            @click="$emit('reset')"
          >
            <span class="material-icons text-[18px]">refresh</span>
            {{ t('reset') }}
          </button>
        </template>
      </BaseFilterBar>
    </template>

    <div class="overflow-x-auto" v-loading="loading">
      <table class="console-table w-full text-left">
        <thead>
          <tr>
            <th class="px-6 py-4 w-12">
              <input
                type="checkbox"
                class="rounded"
                :checked="allSelected"
                @change="$emit('toggle-select-all', $event.target.checked)"
              />
            </th>
            <th class="px-6 py-4">{{ t('id') }}</th>
            <th class="px-6 py-4">{{ t('apiPath') }}</th>
            <th class="px-6 py-4">{{ t('apiGroup') }}</th>
            <th class="px-6 py-4">{{ t('apiDesc') }}</th>
            <th class="px-6 py-4">{{ t('method') }}</th>
            <th class="px-6 py-4 console-actions-header">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
          <tr v-if="items.length === 0">
            <td colspan="7" class="px-6 py-12 text-center text-slate-400">
              <div class="space-y-2">
                <span class="material-icons text-4xl">inbox</span>
                <p>{{ t('noData') }}</p>
              </div>
            </td>
          </tr>
          <tr
            v-for="row in items"
            :key="row.ID"
            class="hover:bg-slate-50 dark:hover:bg-zinc-800/40 transition-colors"
          >
            <td class="px-6 py-4">
              <input
                type="checkbox"
                class="rounded"
                :checked="selectedIds.includes(row.ID)"
                @change="$emit('toggle-select', row)"
              />
            </td>
            <td class="px-6 py-4 text-sm is-secondary is-code">{{ row.ID }}</td>
            <td class="px-6 py-4 text-sm is-primary is-code">{{ row.path }}</td>
            <td class="px-6 py-4">
              <span
                class="px-2 py-0.5 rounded text-xs font-bold bg-slate-100 dark:bg-zinc-700 text-slate-600 dark:text-slate-300"
              >
                {{ row.apiGroup }}
              </span>
            </td>
            <td class="px-6 py-4 text-sm is-secondary">
              {{ row.description }}
            </td>
            <td class="px-6 py-4">
              <span
                :class="getMethodClass(row.method)"
                class="px-2 py-0.5 rounded text-xs font-bold"
              >
                {{ row.method }}
              </span>
            </td>
            <td class="px-6 py-4 console-actions-cell">
              <div class="list-row-actions">
                <button
                  class="list-row-button list-row-button--info"
                  @click="$emit('edit', row)"
                >
                  <span class="material-icons text-[16px]">edit</span>
                  {{ t('edit') }}
                </button>
                <button
                  class="list-row-button list-row-button--danger"
                  @click="$emit('delete', row)"
                >
                  <span class="material-icons text-[16px]">delete</span>
                  {{ t('delete') }}
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <template #footer>
      <ListPaginationBar
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        :total-text="t('totalRecords', { total })"
        :page-sizes="[10, 30, 50, 100]"
        layout="sizes, prev, pager, next, jumper"
        @current-change="$emit('page-change', $event)"
        @size-change="$emit('size-change', $event)"
      />
    </template>
  </TableCard>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import BaseFilterBar from '@/components/listPage/BaseFilterBar.vue'
import ListPaginationBar from '@/components/listPage/ListPaginationBar.vue'
import TableCard from '@/components/listPage/TableCard.vue'
import type { ResourceId, Translator } from '@/types/consoleResource'
import type {
  ApiListItem,
  ApiMethodOption,
  LabelValueOption
} from '@/types/superAdmin'

const props = withDefaults(
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
    pageSize: 10,
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

const t = inject<Translator>('t', (key: string) => key)

const searchApiGroupModel = computed({
  get: () => props.searchApiGroup,
  set: (value) => emit('update:search-api-group', value || '')
})

const searchDescriptionModel = computed({
  get: () => props.searchDescription,
  set: (value: string) => emit('update:search-description', value)
})

const searchMethodModel = computed({
  get: () => props.searchMethod,
  set: (value?: string) => emit('update:search-method', value || '')
})

const searchPathModel = computed({
  get: () => props.searchPath,
  set: (value: string) => emit('update:search-path', value)
})
</script>
