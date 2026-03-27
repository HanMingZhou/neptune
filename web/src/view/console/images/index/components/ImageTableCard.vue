<template>
  <TableCard>
    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="bg-slate-50 dark:bg-zinc-800/50 border-b border-border-light dark:border-border-dark text-slate-500 text-xs font-bold uppercase tracking-wider">
            <th class="px-6 py-4 text-left">{{ t('id') }}</th>
            <th class="px-6 py-4 text-left">{{ t('name') }}</th>
            <th class="px-6 py-4 text-center">{{ t('imageType') }}</th>
            <th class="px-6 py-4 text-center">{{ t('imageUsageType') }}</th>
            <th class="px-6 py-4 text-left">{{ t('imageAddr') }}</th>
            <th class="px-6 py-4 text-left">{{ t('imageArea') }}</th>
            <th class="px-6 py-4 text-left">{{ t('imageSize') }}</th>
            <th class="px-6 py-4 text-left">{{ t('createdAt') }}</th>
            <th class="px-6 py-4 text-center">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
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
          <tr
            v-for="row in items"
            v-else
            :key="row.id"
            class="hover:bg-slate-50 dark:hover:bg-zinc-800/40 transition-colors"
          >
            <td class="px-6 py-4 text-sm font-mono text-slate-500">{{ row.id }}</td>
            <td class="px-6 py-4">
              <span class="text-sm font-bold text-primary">{{ row.name }}</span>
            </td>
            <td class="px-6 py-4 text-center">
              <span
                v-if="row.type === 1"
                class="px-2.5 py-1 rounded-full text-xs font-bold bg-blue-500/10 text-blue-600"
              >
                {{ t('systemImage') }}
              </span>
              <span v-else class="px-2.5 py-1 rounded-full text-xs font-bold bg-amber-500/10 text-amber-600">
                {{ t('customImage') }}
              </span>
            </td>
            <td class="px-6 py-4 text-center">
              <span :class="usageTypeBadgeClass(row.usageType)" class="px-2.5 py-1 rounded-full text-xs font-bold">
                {{ usageTypeLabel(row.usageType) }}
              </span>
            </td>
            <td class="px-6 py-4">
              <el-tooltip v-if="row.image" :content="row.image" placement="top" :show-after="300">
                <code class="text-xs bg-slate-100 dark:bg-zinc-800 text-slate-600 dark:text-slate-400 px-2 py-1 rounded font-mono truncate max-w-[260px] inline-block">
                  {{ row.image }}
                </code>
              </el-tooltip>
              <span v-else class="text-xs text-slate-400">-</span>
            </td>
            <td class="px-6 py-4 text-sm text-slate-600">{{ row.area || '-' }}</td>
            <td class="px-6 py-4 text-sm text-slate-600">{{ row.size || '-' }}</td>
            <td class="px-6 py-4 text-sm text-slate-500">{{ row.createTime }}</td>
            <td class="px-6 py-4 text-center">
              <div class="flex justify-center gap-2 items-center">
                <button
                  class="bg-primary/10 hover:bg-primary/20 text-primary px-2 py-1 rounded-sm text-xs font-bold transition-colors"
                  @click="$emit('edit', row)"
                >
                  {{ t('edit') }}
                </button>
                <button
                  class="bg-red-500/10 hover:bg-red-500/20 text-red-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors"
                  @click="$emit('delete', row)"
                >
                  {{ t('delete') }}
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <template #footer>
      <div v-if="total > 0" class="flex justify-end">
        <el-pagination
          background
          layout="prev, pager, next"
          :total="total"
          :page-size="pageSize"
          v-model:current-page="pageModel"
          @current-change="$emit('page-change')"
        />
      </div>
    </template>
  </TableCard>
</template>

<script setup>
import { computed, inject } from 'vue'
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

const emit = defineEmits(['delete', 'edit', 'page-change', 'update:current-page'])
const t = inject('t', (key) => key)

const pageModel = computed({
  get: () => props.currentPage,
  set: (value) => emit('update:current-page', value)
})

const usageTypeLabel = (value) => {
  const labelMap = {
    1: t('usageNotebook'),
    2: t('usageTrain'),
    3: t('usageInfer')
  }

  return labelMap[value] || '-'
}

const usageTypeBadgeClass = (value) => {
  const classMap = {
    1: 'bg-emerald-500/10 text-emerald-600',
    2: 'bg-purple-500/10 text-purple-600',
    3: 'bg-cyan-500/10 text-cyan-600'
  }

  return classMap[value] || 'bg-slate-500/10 text-slate-500'
}
</script>
