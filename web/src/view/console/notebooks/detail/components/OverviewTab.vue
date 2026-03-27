<template>
  <div class="space-y-6">
    <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
      <h3 class="text-base font-bold mb-4 flex items-center gap-2">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('basicInfo') }}
      </h3>
      <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-6">
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('name') }}</div>
          <div class="text-sm font-medium">{{ notebook.displayName || '-' }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('instanceId') }}</div>
          <div class="text-sm font-medium font-mono break-all">{{ notebook.instanceName || '-' }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('status') }}</div>
          <div class="text-sm font-medium">{{ getStatusLabel(notebook.status) }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('createdAt') }}</div>
          <div class="text-sm font-medium">{{ formatTime(notebook.createdAt) }}</div>
        </div>
        <div class="col-span-2 md:col-span-3 lg:col-span-1">
          <div class="text-xs text-slate-400 mb-1">{{ t('image') }}</div>
          <div class="text-sm font-medium bg-slate-100 dark:bg-zinc-800 px-2 py-1 rounded inline-flex font-mono text-[11px] break-all max-w-full">
            {{ notebook.imageName || '-' }}
          </div>
        </div>
      </div>
    </div>

    <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
      <h3 class="text-base font-bold mb-4 flex items-center gap-2">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('resourceConfig') }}
      </h3>
      <div class="grid grid-cols-2 md:grid-cols-3 gap-6">
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('cpu') }}</div>
          <div class="text-sm font-medium">{{ notebook.cpu }} {{ t('cpu') }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('memory') }}</div>
          <div class="text-sm font-medium">{{ notebook.memory }} GB</div>
        </div>
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('gpu') }}</div>
          <div class="text-sm font-medium">
            <span v-if="notebook.gpuCount" class="text-primary font-bold">{{ notebook.gpuCount }} × {{ notebook.gpuModel }}</span>
            <span v-else class="text-slate-400">CPU {{ t('ONLY') }}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
      <h3 class="text-base font-bold mb-4 flex items-center gap-2">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('order') }}
      </h3>
      <div class="grid grid-cols-2 md:grid-cols-3 gap-6">
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('payType') }}</div>
          <div class="text-sm font-medium">{{ getPayTypeLabel(notebook.payType) }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('unitPrice') }}</div>
          <div class="text-lg font-bold text-red-500">￥{{ notebook.price?.toFixed(4) || '0.0000' }} / {{ getUnitPriceLabel(notebook.payType) }}</div>
        </div>
      </div>
    </div>

    <div v-if="notebook.enableTensorboard" class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
      <h3 class="text-base font-bold mb-4 flex items-center gap-2">
        <span class="w-1 h-4 bg-primary rounded"></span>
        TensorBoard
      </h3>
      <div class="grid grid-cols-2 md:grid-cols-3 gap-6">
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('logPath') }}</div>
          <div class="text-sm font-medium font-mono">{{ notebook.tensorboardLogPath || 'logs' }}</div>
        </div>
        <div v-if="notebook.tensorboardUrl">
          <div class="text-xs text-slate-400 mb-1">{{ t('accessLink') }}</div>
          <a :href="notebook.tensorboardUrl" target="_blank" class="text-sm text-primary hover:underline flex items-center gap-1">
            <span class="material-icons text-base">open_in_new</span>
            {{ t('openTensorboard') }}
          </a>
        </div>
      </div>
    </div>

    <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
      <h3 class="text-base font-bold mb-4 flex items-center gap-2">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('dataMount') }}
      </h3>
      <div v-if="notebook.volumeMounts && notebook.volumeMounts.length > 0" class="overflow-x-auto">
        <table class="w-full text-left">
          <thead>
            <tr class="text-xs text-slate-400 font-bold uppercase tracking-wider border-b border-border-light dark:border-border-dark">
              <th class="px-4 py-3">{{ t('storage') }}</th>
              <th class="px-4 py-3">{{ t('mountPath') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border-light dark:divide-border-dark">
            <tr v-for="(mount, index) in notebook.volumeMounts" :key="index" class="text-sm">
              <td class="px-4 py-3">{{ mount.pvcName || '-' }}</td>
              <td class="px-4 py-3 font-mono">{{ mount.mountsPath }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-else class="text-sm text-slate-400 text-center py-8">{{ t('noData') }}</div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  formatTime: {
    type: Function,
    required: true
  },
  getPayTypeLabel: {
    type: Function,
    required: true
  },
  getStatusLabel: {
    type: Function,
    required: true
  },
  getUnitPriceLabel: {
    type: Function,
    required: true
  },
  notebook: {
    type: Object,
    default: () => ({})
  }
})

const t = inject('t', (key) => key)
</script>
