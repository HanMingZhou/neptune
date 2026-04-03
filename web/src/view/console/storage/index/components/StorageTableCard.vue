<template>
  <TableCard>
    <div class="overflow-x-auto">
      <table class="console-table">
        <thead>
          <tr>
            <th>{{ t('name') }}</th>
            <th>{{ t('storageProduct') }}</th>
            <th>{{ t('capacity') }}</th>
            <th>{{ t('status') }}</th>
            <th>{{ t('area') }}</th>
            <th>{{ t('createdAt') }}</th>
            <th class="console-actions-header">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody>
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
                <button class="list-row-button list-row-button--neutral" @click="$emit('create')">
                  {{ t('create') }}{{ t('storage') }}
                </button>
              </div>
            </td>
          </tr>
          <tr v-for="item in items" v-else :key="item.id">
            <td><span class="is-primary text-sm">{{ item.name }}</span></td>
            <td>{{ item.productName || '-' }}</td>
            <td class="is-code">{{ item.size }}</td>
            <td>
              <span
                :class="item.status === 'Bound' ? 'console-badge console-badge--success' : 'console-badge console-badge--neutral'"
              >
                {{ t(item.status) || item.status }}
              </span>
            </td>
            <td>{{ item.area }}</td>
            <td class="text-slate-500">{{ item.createdAt }}</td>
            <td class="console-actions-cell">
              <div class="list-row-actions">
                <button
                  class="list-row-button list-row-button--info"
                  @click="$emit('expand', item)"
                >
                  <span class="material-icons text-[14px]">add_circle_outline</span>
                  {{ t('expand') }}
                </button>
                <button
                  :disabled="Boolean(btnLoading[item.id])"
                  class="list-row-button list-row-button--danger disabled:opacity-50"
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
      <ListPaginationBar
        v-model:current-page="pageModel"
        :total="total"
        :total-text="t('totalRecords', { total })"
        :page-size="10"
        :page-sizes="[10]"
        layout="prev, pager, next"
        @current-change="$emit('refresh')"
      />
    </template>
  </TableCard>
</template>

<script setup>
import { computed, inject } from 'vue'
import ListPaginationBar from '@/components/listPage/ListPaginationBar.vue'
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
  'update:page'
])

const t = inject('t', (key) => key)

const pageModel = computed({
  get: () => props.page,
  set: (value) => emit('update:page', value)
})
</script>
