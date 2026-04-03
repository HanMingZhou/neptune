<template>
  <TableCard>
    <div class="overflow-x-auto">
      <table class="console-table console-table--compact w-full min-w-[1180px]">
        <thead>
          <tr>
            <th>{{ t('id') }}</th>
            <th>{{ t('name') }}</th>
            <th class="text-center">{{ t('imageType') }}</th>
            <th class="text-center">{{ t('imageUsageType') }}</th>
            <th>{{ t('imageAddr') }}</th>
            <th>{{ t('imageArea') }}</th>
            <th>{{ t('imageSize') }}</th>
            <th>{{ t('createdAt') }}</th>
            <th class="console-actions-header">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td colspan="9" class="px-6 py-12 text-center text-slate-400">
              <div class="flex items-center justify-center gap-2">
                <div class="animate-spin rounded-full h-5 w-5 border-b-2 border-primary"></div>
                {{ t('loading') }}
              </div>
            </td>
          </tr>
          <tr v-else-if="items.length === 0">
            <td colspan="9" class="px-6 py-10 text-center text-slate-400 text-sm">{{ t('noData') }}</td>
          </tr>
          <tr v-for="row in items" v-else :key="row.id">
            <td class="is-code text-slate-500">{{ row.id }}</td>
            <td>
              <span class="is-primary text-sm">{{ row.name }}</span>
            </td>
            <td class="text-center">
              <ListToneBadge
                :label="row.type === 1 ? t('systemImage') : t('customImage')"
                :tone="row.type === 1 ? 'info' : 'warning'"
              />
            </td>
            <td class="text-center">
              <ListToneBadge :label="usageTypeLabel(row.usageType)" :tone="usageTypeBadgeTone(row.usageType)" />
            </td>
            <td>
              <el-tooltip v-if="row.image" :content="row.image" placement="top" :show-after="300">
                <code class="is-code is-secondary inline-block max-w-[260px] truncate">
                  {{ row.image }}
                </code>
              </el-tooltip>
              <span v-else class="text-xs text-slate-400">-</span>
            </td>
            <td>{{ row.area || '-' }}</td>
            <td>{{ row.size || '-' }}</td>
            <td class="is-secondary">{{ row.createTime }}</td>
            <td class="console-actions-cell">
              <div class="list-row-actions">
                <button
                  class="list-row-button list-row-button--info"
                  @click="$emit('edit', row)"
                >
                  <span class="material-icons text-[14px]">edit</span>
                  {{ t('edit') }}
                </button>
                <button
                  class="list-row-button list-row-button--danger"
                  @click="$emit('delete', row)"
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
        v-model:page-size="pageSizeModel"
        :total="total"
        :total-text="t('totalRecords', { total })"
        :page-sizes="[10, 20, 50, 100]"
        layout="sizes, prev, pager, next, jumper"
        :hide-when-empty="true"
        @current-change="$emit('page-change')"
        @size-change="$emit('size-change', $event)"
      />
    </template>
  </TableCard>
</template>

<script setup>
import { computed, inject } from 'vue'
import ListPaginationBar from '@/components/listPage/ListPaginationBar.vue'
import ListToneBadge from '@/components/listPage/ListToneBadge.vue'
import TableCard from '@/components/listPage/TableCard.vue'

const props = defineProps({
  currentPage: {
    type: Number,
    default: 1
  },
  items: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  pageSize: {
    type: Number,
    default: 10
  },
  total: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits(['delete', 'edit', 'page-change', 'size-change', 'update:current-page', 'update:page-size'])
const t = inject('t', (key) => key)

const pageModel = computed({
  get: () => props.currentPage,
  set: (value) => emit('update:current-page', value)
})

const pageSizeModel = computed({
  get: () => props.pageSize,
  set: (value) => emit('update:page-size', value)
})

const usageTypeLabel = (value) => {
  const labelMap = {
    1: t('usageNotebook'),
    2: t('usageTrain'),
    3: t('usageInfer')
  }

  return labelMap[value] || '-'
}

const usageTypeBadgeTone = (value) => {
  const toneMap = {
    1: 'success',
    2: 'warning',
    3: 'info'
  }

  return toneMap[value] || 'neutral'
}
</script>
