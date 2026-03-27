<template>
  <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl overflow-hidden shadow-sm">
    <div class="overflow-x-auto">
      <table class="w-full min-w-[1180px]" v-loading="loading">
        <thead>
          <tr class="bg-slate-50 dark:bg-zinc-800/50 border-b border-border-light dark:border-border-dark text-slate-500 text-xs font-bold uppercase tracking-wider">
            <th class="px-6 py-4 text-left">{{ t('nodeName') }}</th>
            <th class="px-6 py-4 text-left">{{ t('internalIp') }}</th>
            <th class="px-6 py-4 text-left">{{ t('clusterName') }}</th>
            <th class="px-6 py-4 text-left">{{ t('nodeRole') }}</th>
            <th class="px-6 py-4 text-left">{{ t('area') }}</th>
            <th class="px-6 py-4 text-left">{{ t('spec') }} (可用/总量)</th>
            <th class="px-6 py-4 text-center">{{ t('status') }}</th>
            <th class="px-6 py-4 text-center">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
          <tr
            v-for="row in items"
            :key="row.nodeName"
            class="hover:bg-slate-50 dark:hover:bg-zinc-800/40 transition-colors"
          >
            <td class="px-6 py-4">
              <div class="flex flex-col">
                <code class="text-sm font-bold text-primary bg-primary/5 px-2 py-1 rounded w-fit">{{ row.nodeName }}</code>
                <span class="text-[10px] text-slate-400 mt-1">{{ row.area || '-' }}</span>
              </div>
            </td>
            <td class="px-6 py-4 text-sm text-slate-600 font-mono">{{ row.internalIp || 'N/A' }}</td>
            <td class="px-6 py-4 text-sm text-slate-600">{{ row.clusterName || '-' }}</td>
            <td class="px-6 py-4">
              <span
                v-if="row.nodeRole === 'master'"
                class="px-2 py-0.5 rounded-full text-xs font-bold bg-blue-500/10 text-blue-600"
              >
                {{ t('roleMaster') }}
              </span>
              <span
                v-else
                class="px-2 py-0.5 rounded-full text-xs font-bold bg-slate-500/10 text-slate-500"
              >
                {{ t('roleWorker') }}
              </span>
            </td>
            <td class="px-6 py-4 text-sm text-slate-600">{{ row.area || '-' }}</td>
            <td class="px-6 py-4">
              <div class="space-y-1.5">
                <div v-if="row.gpuModel" class="flex items-center gap-2">
                  <span class="text-[10px] bg-amber-500 text-white px-1 rounded font-bold leading-none">GPU</span>
                  <span class="text-xs font-bold text-slate-700 dark:text-slate-200">{{ row.gpuModel }}</span>
                  <span
                    :class="row.gpuAvailable > 0 ? 'text-emerald-500' : 'text-slate-400'"
                    class="text-xs font-mono"
                  >
                    ({{ row.gpuAvailable }}/{{ row.gpuCount }})
                  </span>
                </div>

                <div class="flex items-center gap-4">
                  <div class="flex items-center gap-1.5">
                    <span class="material-icons text-[14px] text-slate-400">developer_board</span>
                    <span class="text-xs text-slate-600 dark:text-slate-400 font-mono">{{ row.cpuAvailable }}/{{ row.cpuAllocatable }}核</span>
                  </div>
                  <div class="flex items-center gap-1.5">
                    <span class="material-icons text-[14px] text-slate-400">memory</span>
                    <span class="text-xs text-slate-600 dark:text-slate-400 font-mono">{{ row.memoryAvailable }}/{{ row.memoryAllocatable }}GB</span>
                  </div>
                </div>
              </div>
            </td>
            <td class="px-6 py-4 text-center">
              <span
                v-if="row.schedulable"
                class="px-2.5 py-1 rounded-full text-xs font-bold bg-emerald-500/10 text-emerald-500"
              >
                <span class="size-1.5 bg-emerald-500 rounded-full inline-block mr-1.5 mb-0.5"></span>
                {{ t('schedulable') }}
              </span>
              <span
                v-else
                class="px-2.5 py-1 rounded-full text-xs font-bold bg-amber-500/10 text-amber-500"
              >
                <span class="size-1.5 bg-amber-500 rounded-full inline-block mr-1.5 mb-0.5"></span>
                {{ t('unschedulable') }}
              </span>
            </td>
            <td class="px-6 py-4 text-center">
              <div class="flex justify-center gap-2 items-center">
                <button
                  v-if="!row.schedulable"
                  class="bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-600 px-3 py-1.5 rounded-lg text-xs font-bold transition-all flex items-center gap-1"
                  @click="$emit('uncordon', row)"
                >
                  <span class="material-icons text-[16px]">login</span>
                  {{ t('joinCluster') }}
                </button>
                <button
                  v-if="row.schedulable"
                  class="bg-red-500/10 hover:bg-red-500/20 text-red-600 px-3 py-1.5 rounded-lg text-xs font-bold transition-all flex items-center gap-1"
                  @click="$emit('drain', row)"
                >
                  <span class="material-icons text-[16px]">logout</span>
                  {{ t('evictFromCluster') }}
                </button>
              </div>
            </td>
          </tr>
          <tr v-if="items.length === 0 && !loading">
            <td colspan="8" class="px-6 py-10 text-center text-slate-400 text-sm">
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

defineEmits(['drain', 'uncordon'])

const t = inject('t', (key) => key)
</script>
