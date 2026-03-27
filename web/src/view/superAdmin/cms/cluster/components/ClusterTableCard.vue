<template>
  <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl overflow-hidden shadow-sm">
    <div class="overflow-x-auto">
      <table class="w-full min-w-[1400px]" v-loading="loading">
        <thead>
          <tr class="bg-slate-50 dark:bg-zinc-800/50 border-b border-border-light dark:border-border-dark text-slate-500 text-xs font-bold uppercase tracking-wider">
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
            <th class="px-6 py-4 text-center">{{ t('actions') }}</th>
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
              <span class="text-sm font-bold text-primary">{{ row.name }}</span>
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
              <code
                v-if="row.apiServer"
                class="text-xs bg-slate-100 dark:bg-zinc-800 text-slate-600 dark:text-slate-400 px-2 py-1 rounded font-mono"
              >
                {{ row.apiServer }}
              </code>
              <span v-else class="text-xs text-slate-400">-</span>
            </td>
            <td class="px-6 py-4">
              <code
                v-if="row.harborAddr"
                class="text-xs bg-emerald-500/10 text-emerald-600 px-2 py-1 rounded font-mono"
              >
                {{ row.harborAddr }}
              </code>
              <span v-else class="text-xs text-slate-400">-</span>
            </td>
            <td class="px-6 py-4">
              <code
                v-if="row.storageClass"
                class="text-xs bg-blue-500/10 text-blue-600 px-2 py-1 rounded"
              >
                {{ row.storageClass }}
              </code>
              <span v-else class="text-xs text-slate-400">-</span>
            </td>
            <td class="px-6 py-4">
              <button
                v-if="row.kubeConfig"
                class="text-xs bg-slate-500/10 hover:bg-slate-500/20 text-slate-600 px-2 py-1 rounded-sm font-bold transition-colors"
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
              <span
                v-if="row.status === 1"
                class="px-2.5 py-1 rounded-full text-xs font-bold bg-emerald-500/10 text-emerald-500"
              >
                <span class="size-1.5 bg-emerald-500 rounded-full inline-block mr-1.5 mb-0.5"></span>
                {{ t('enable') }}
              </span>
              <span
                v-else
                class="px-2.5 py-1 rounded-full text-xs font-bold bg-slate-500/10 text-slate-500"
              >
                <span class="size-1.5 bg-slate-400 rounded-full inline-block mr-1.5 mb-0.5"></span>
                {{ t('disable') }}
              </span>
            </td>
            <td class="px-6 py-4 text-sm text-slate-500">{{ row.createdAt }}</td>
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
          <tr v-if="items.length === 0 && !loading">
            <td colspan="13" class="px-6 py-10 text-center text-slate-400 text-sm">
              {{ t('noData') }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

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
