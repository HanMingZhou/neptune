<template>
  <TableCard>
    <template #toolbar>
      <div class="flex flex-wrap gap-4 items-center">
        <div class="relative flex-1 max-w-xs">
          <span class="material-icons absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400 text-[16px]">search</span>
          <input
            v-model="searchNameModel"
            type="text"
            :placeholder="t('searchStoragePlaceholder')"
            class="list-search-input"
            @keyup.enter="$emit('refresh')"
          />
        </div>
        <el-select
          v-model="searchStatusModel"
          :placeholder="`${t('status')}: ${t('all')}`"
          clearable
          class="list-filter-select"
          size="small"
        >
          <el-option :label="t('Creating')" value="Creating" />
          <el-option :label="t('Bound')" value="Bound" />
          <el-option :label="t('Error')" value="Error" />
        </el-select>
      </div>
    </template>

    <div class="overflow-x-auto">
      <table class="w-full text-left">
        <thead>
          <tr class="bg-slate-50 dark:bg-zinc-800/50 border-b border-border-light dark:border-border-dark text-slate-500 text-xs font-bold uppercase tracking-wider">
            <th class="px-6 py-4">{{ t('name') }}</th>
            <th class="px-6 py-4">{{ t('storageProduct') }}</th>
            <th class="px-6 py-4">{{ t('capacity') }}</th>
            <th class="px-6 py-4">{{ t('status') }}</th>
            <th class="px-6 py-4">{{ t('area') }}</th>
            <th class="px-6 py-4">{{ t('createdAt') }}</th>
            <th class="px-6 py-4 text-center">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
          <tr v-if="loading">
            <td colspan="7" class="px-6 py-12 text-center text-slate-400">
              <div class="flex items-center justify-center gap-2">
                <div class="animate-spin rounded-full h-5 w-5 border-b-2 border-primary"></div>
                {{ t('loading') }}
              </div>
            </td>
          </tr>
          <tr v-else-if="items.length === 0">
            <td colspan="7" class="px-6 py-12 text-center text-slate-400">
              <div class="space-y-2">
                <span class="material-icons text-4xl">folder_open</span>
                <p>{{ t('noData') }}</p>
                <button class="text-primary hover:underline" @click="$emit('create')">
                  {{ t('create') }}{{ t('storage') }}
                </button>
              </div>
            </td>
          </tr>
          <tr
            v-for="item in items"
            v-else
            :key="item.id"
            class="hover:bg-slate-50 dark:hover:bg-zinc-800/40 transition-colors"
          >
            <td class="px-6 py-4 font-bold text-sm">{{ item.name }}</td>
            <td class="px-6 py-4 text-sm text-slate-600 dark:text-slate-400">{{ item.productName || '-' }}</td>
            <td class="px-6 py-4 text-sm text-slate-600 dark:text-slate-400 font-mono">{{ item.size }}</td>
            <td class="px-6 py-4">
              <span
                :class="item.status === 'Bound' ? 'bg-emerald-500/10 text-emerald-500' : 'bg-slate-500/10 text-slate-500'"
                class="px-2.5 py-1 rounded-full text-xs font-bold"
              >
                {{ t(item.status) || item.status }}
              </span>
            </td>
            <td class="px-6 py-4 text-sm text-slate-600 dark:text-slate-400">{{ item.area }}</td>
            <td class="px-6 py-4 text-sm text-slate-600 dark:text-slate-400">{{ item.createdAt }}</td>
            <td class="px-6 py-4 text-center">
              <div class="flex justify-center gap-2 items-center">
                <button
                  class="bg-primary/10 hover:bg-primary/20 text-primary px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1"
                  @click="$emit('expand', item)"
                >
                  <span class="material-icons text-[14px]">add_circle_outline</span>
                  {{ t('expand') }}
                </button>
                <button
                  :disabled="Boolean(btnLoading[item.id])"
                  class="bg-red-500/10 hover:bg-red-500/20 text-red-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1 disabled:opacity-50"
                  @click="$emit('delete', item)"
                >
                  <span class="material-icons text-[14px]">delete</span>
                  {{ t('delete') }}
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <template #footer>
      <div class="flex items-center justify-between">
        <span>{{ t('totalRecords', { total }) }}</span>
        <el-pagination
          v-model:current-page="pageModel"
          :total="total"
          layout="prev, pager, next"
          @current-change="$emit('refresh')"
        />
      </div>
    </template>
  </TableCard>
</template>

<script setup>
import { computed, inject } from 'vue'
import TableCard from '@/components/listPage/TableCard.vue'

const props = defineProps({
  btnLoading: {
    type: Object,
    default: () => ({})
  },
  items: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  page: {
    type: Number,
    default: 1
  },
  searchName: {
    type: String,
    default: ''
  },
  searchStatus: {
    type: String,
    default: ''
  },
  total: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits([
  'create',
  'delete',
  'expand',
  'refresh',
  'update:page',
  'update:search-name',
  'update:search-status'
])

const t = inject('t', (key) => key)

const searchNameModel = computed({
  get: () => props.searchName,
  set: (value) => emit('update:search-name', value)
})

const searchStatusModel = computed({
  get: () => props.searchStatus,
  set: (value) => emit('update:search-status', value)
})

const pageModel = computed({
  get: () => props.page,
  set: (value) => emit('update:page', value)
})
</script>
