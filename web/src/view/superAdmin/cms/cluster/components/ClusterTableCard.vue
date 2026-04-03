<template>
  <TableCard>
    <div class="overflow-x-auto">
      <table class="console-table min-w-[1400px]" v-loading="loading">
        <thead>
          <tr>
            <th class="px-6 py-4 text-left">{{ t('id') }}</th>
            <th class="px-6 py-4 text-left">{{ t('clusterName') }}</th>
            <th class="px-6 py-4 text-left">{{ t('desc') }}</th>
            <th class="px-6 py-4 text-left">{{ t('area') }}</th>
            <th class="px-6 py-4 text-left">{{ t('internalIp') }}</th>
            <th class="px-6 py-4 text-left">{{ t('clusterApiServer') }}</th>
            <th class="px-6 py-4 text-left">{{ t('clusterHarbor') }}</th>
            <th class="px-6 py-4 text-left">{{ t('storageClass') }}</th>
            <th class="px-6 py-4 text-left">{{ t('clusterKubeconfig') }}</th>
            <th class="px-6 py-4 text-center">{{ t('clusterNodeCount') }}</th>
            <th class="px-6 py-4 text-center">{{ t('status') }}</th>
            <th class="px-6 py-4 text-left">{{ t('createdAt') }}</th>
            <th class="px-6 py-4 console-actions-header">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
          <tr
            v-for="row in items"
            :key="row.id"
            class="hover:bg-slate-50 dark:hover:bg-zinc-800/40 transition-colors"
          >
            <td class="px-6 py-4 text-sm font-mono text-slate-500">{{ row.id }}</td>
            <td class="px-6 py-4">
              <span class="is-primary">{{ row.name }}</span>
            </td>
            <td class="px-6 py-4 text-sm text-slate-500 max-w-[220px]">
              <el-tooltip v-if="row.description" :content="row.description" placement="top" :show-after="300">
                <span class="line-clamp-2">{{ row.description }}</span>
              </el-tooltip>
              <span v-else>-</span>
            </td>
            <td class="px-6 py-4 text-sm text-slate-600">{{ row.area || '-' }}</td>
            <td class="px-6 py-4 text-sm text-slate-600 font-mono">{{ row.internalIp || '-' }}</td>
            <td class="px-6 py-4">
              <span
                v-if="row.apiServer"
                class="console-badge console-badge--neutral is-code"
              >
                {{ row.apiServer }}
              </span>
              <span v-else class="text-xs text-slate-400">-</span>
            </td>
            <td class="px-6 py-4">
              <span
                v-if="row.harborAddr"
                class="console-badge console-badge--neutral is-code"
              >
                {{ row.harborAddr }}
              </span>
              <span v-else class="text-xs text-slate-400">-</span>
            </td>
            <td class="px-6 py-4">
              <span
                v-if="row.storageClass"
                class="console-badge console-badge--neutral is-code"
              >
                {{ row.storageClass }}
              </span>
              <span v-else class="text-xs text-slate-400">-</span>
            </td>
            <td class="px-6 py-4">
              <button
                v-if="row.kubeConfig"
                class="list-row-button list-row-button--neutral"
                @click="$emit('view-kubeconfig', row)"
              >
                {{ t('view') }}
              </button>
              <span v-else class="text-xs text-slate-400">-</span>
            </td>
            <td class="px-6 py-4 text-center">
              <span class="text-sm font-bold text-slate-700 dark:text-slate-200">{{ row.nodeCount }}</span>
            </td>
            <td class="px-6 py-4 text-center">
              <ListToneBadge
                :label="row.status === 1 ? t('enable') : t('disable')"
                :tone="row.status === 1 ? 'success' : 'neutral'"
              />
            </td>
            <td class="px-6 py-4 text-sm text-slate-500">{{ row.createdAt }}</td>
            <td class="px-6 py-4 console-actions-cell">
              <div class="list-row-actions">
                <button
                  class="list-row-button list-row-button--info"
                  @click="$emit('edit', row)"
                >
                  {{ t('edit') }}
                </button>
                <button
                  class="list-row-button list-row-button--danger"
                  @click="$emit('delete', row)"
                >
                  {{ t('delete') }}
                </button>
              </div>
            </td>
          </tr>
          <tr v-if="items.length === 0 && !loading">
            <td colspan="13" class="px-6 py-10 text-center text-slate-400 text-sm">
              {{ t('noData') }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </TableCard>
</template>

<script setup>
import { inject } from 'vue'
import ListToneBadge from '@/components/listPage/ListToneBadge.vue'
import TableCard from '@/components/listPage/TableCard.vue'

defineProps({
  items: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  }
})

defineEmits(['delete', 'edit', 'view-kubeconfig'])

const t = inject('t', (key) => key)
</script>
