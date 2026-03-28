<template>
  <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl overflow-hidden shadow-sm">
    <div class="list-filter-bar border-b border-border-light p-4 dark:border-border-dark">
      <div class="list-filter-field list-filter-field--compact max-w-[160px]">
        <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]">link</span>
        <input
          v-model="searchPathModel"
          type="text"
          :placeholder="t('path')"
          class="list-search-input !w-full"
          @keyup.enter="$emit('search')"
        />
      </div>

      <div class="list-filter-field list-filter-field--compact max-w-[160px]">
        <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]">description</span>
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

      <div class="list-toolbar-actions">
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
      </div>
    </div>

    <div class="overflow-x-auto" v-loading="loading">
      <table class="w-full text-left">
        <thead>
          <tr class="bg-slate-50 dark:bg-zinc-800/50 border-b border-border-light dark:border-border-dark text-slate-500 text-xs font-bold uppercase tracking-wider">
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
            <th class="px-6 py-4 text-right">{{ t('actions') }}</th>
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
            <td class="px-6 py-4 text-sm text-slate-500 font-mono">{{ row.ID }}</td>
            <td class="px-6 py-4 text-sm font-mono text-primary">{{ row.path }}</td>
            <td class="px-6 py-4">
              <span class="px-2 py-0.5 rounded text-xs font-bold bg-slate-100 dark:bg-zinc-700 text-slate-600 dark:text-slate-300">
                {{ row.apiGroup }}
              </span>
            </td>
            <td class="px-6 py-4 text-sm text-slate-600 dark:text-slate-400">{{ row.description }}</td>
            <td class="px-6 py-4">
              <span :class="getMethodClass(row.method)" class="px-2 py-0.5 rounded text-xs font-bold">
                {{ row.method }}
              </span>
            </td>
            <td class="px-6 py-4 text-right">
              <div class="flex justify-end gap-3 items-center">
                <button
                  class="bg-primary/10 hover:bg-primary/20 text-primary px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1"
                  @click="$emit('edit', row)"
                >
                  <span class="material-icons text-[16px]">edit</span>
                  {{ t('edit') }}
                </button>
                <button
                  class="bg-red-500/10 hover:bg-red-500/20 text-red-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1"
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

    <div class="px-6 py-4 bg-slate-50 dark:bg-zinc-800/30 flex items-center justify-between border-t border-border-light dark:border-border-dark">
      <span class="text-xs text-slate-500">{{ t('totalRecords', { total }) }}</span>
      <el-pagination
        :current-page="page"
        :page-size="pageSize"
        :page-sizes="[10, 30, 50, 100]"
        :total="total"
        background
        layout="sizes, prev, pager, next, jumper"
        @current-change="$emit('page-change', $event)"
        @size-change="$emit('size-change', $event)"
      />
    </div>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'

const props = defineProps({
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

const t = inject('t', (key) => key)

const searchApiGroupModel = computed({
  get: () => props.searchApiGroup,
  set: (value) => emit('update:search-api-group', value || '')
})

const searchDescriptionModel = computed({
  get: () => props.searchDescription,
  set: (value) => emit('update:search-description', value)
})

const searchMethodModel = computed({
  get: () => props.searchMethod,
  set: (value) => emit('update:search-method', value || '')
})

const searchPathModel = computed({
  get: () => props.searchPath,
  set: (value) => emit('update:search-path', value)
})
</script>
