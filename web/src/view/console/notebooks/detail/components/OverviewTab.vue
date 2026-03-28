<template>
  <div class="space-y-4 md:space-y-5">
    <section class="rounded-xl border border-border-light bg-surface-light p-4 md:p-5 dark:border-border-dark dark:bg-surface-dark">
      <h3 class="mb-4 flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('basicInfo') }}
      </h3>
      <div class="grid grid-cols-1 gap-x-5 gap-y-4 sm:grid-cols-2 lg:grid-cols-4">
        <div class="space-y-1">
          <div class="text-xs text-slate-500">{{ t('name') }}</div>
          <div class="text-sm font-semibold text-slate-900 dark:text-slate-100">{{ notebook.displayName || '-' }}</div>
        </div>
        <div class="space-y-1">
          <div class="text-xs text-slate-500">{{ t('instanceId') }}</div>
          <div class="text-sm font-semibold font-mono text-slate-900 break-all dark:text-slate-100">{{ notebook.instanceName || '-' }}</div>
        </div>
        <div class="space-y-1">
          <div class="text-xs text-slate-500">{{ t('status') }}</div>
          <div class="text-sm font-semibold text-slate-900 dark:text-slate-100">{{ getStatusLabel(notebook.status) }}</div>
        </div>
        <div class="space-y-1">
          <div class="text-xs text-slate-500">{{ t('createdAt') }}</div>
          <div class="text-sm font-semibold text-slate-900 dark:text-slate-100">{{ formatTime(notebook.createdAt) }}</div>
        </div>
        <div class="space-y-1 sm:col-span-2 lg:col-span-2">
          <div class="text-xs text-slate-500">{{ t('image') }}</div>
          <div class="inline-flex max-w-full rounded-md bg-slate-100 px-2.5 py-1.5 text-[11px] font-mono leading-relaxed break-all text-slate-700 dark:bg-zinc-800 dark:text-slate-200">
            {{ notebook.imageName || '-' }}
          </div>
        </div>
      </div>
    </section>

    <section class="rounded-xl border border-border-light bg-surface-light p-4 md:p-5 dark:border-border-dark dark:bg-surface-dark">
      <h3 class="mb-4 flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('resourceConfig') }}
      </h3>
      <div class="grid grid-cols-1 gap-x-5 gap-y-4 sm:grid-cols-2 lg:grid-cols-3">
        <div class="space-y-1">
          <div class="text-xs text-slate-500">{{ t('cpu') }}</div>
          <div class="text-sm font-semibold text-slate-900 dark:text-slate-100">{{ notebook.cpu }} {{ t('cpu') }}</div>
        </div>
        <div class="space-y-1">
          <div class="text-xs text-slate-500">{{ t('memory') }}</div>
          <div class="text-sm font-semibold text-slate-900 dark:text-slate-100">{{ notebook.memory }} GB</div>
        </div>
        <div class="space-y-1">
          <div class="text-xs text-slate-500">{{ t('gpu') }}</div>
          <div class="text-sm font-semibold text-slate-900 dark:text-slate-100">
            <span v-if="notebook.gpuCount" class="text-primary font-bold">{{ notebook.gpuCount }} × {{ notebook.gpuModel }}</span>
            <span v-else class="text-slate-400">CPU {{ t('ONLY') }}</span>
          </div>
        </div>
      </div>
    </section>

    <section class="rounded-xl border border-border-light bg-surface-light p-4 md:p-5 dark:border-border-dark dark:bg-surface-dark">
      <h3 class="mb-4 flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('order') }}
      </h3>
      <div class="grid grid-cols-1 gap-x-5 gap-y-4 sm:grid-cols-2">
        <div class="space-y-1">
          <div class="text-xs text-slate-500">{{ t('payType') }}</div>
          <div class="text-sm font-semibold text-slate-900 dark:text-slate-100">{{ getPayTypeLabel(notebook.payType) }}</div>
        </div>
        <div class="space-y-1">
          <div class="text-xs text-slate-500">{{ t('unitPrice') }}</div>
          <div class="text-sm font-semibold text-red-500">￥{{ notebook.price?.toFixed(4) || '0.0000' }} / {{ getUnitPriceLabel(notebook.payType) }}</div>
        </div>
      </div>
    </section>

    <section
      v-if="notebook.enableTensorboard"
      class="rounded-xl border border-border-light bg-surface-light p-4 md:p-5 dark:border-border-dark dark:bg-surface-dark"
    >
      <h3 class="mb-4 flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
        <span class="w-1 h-4 bg-primary rounded"></span>
        TensorBoard
      </h3>
      <div class="grid grid-cols-1 gap-x-5 gap-y-4 md:grid-cols-2">
        <div class="space-y-1">
          <div class="text-xs text-slate-500">{{ t('logPath') }}</div>
          <div class="text-sm font-semibold font-mono text-slate-900 dark:text-slate-100">{{ notebook.tensorboardLogPath || 'logs' }}</div>
        </div>
        <div v-if="notebook.tensorboardUrl" class="space-y-1">
          <div class="text-xs text-slate-500">{{ t('accessLink') }}</div>
          <a :href="notebook.tensorboardUrl" target="_blank" class="inline-flex items-center gap-1 text-sm text-primary hover:underline">
            <span class="material-icons text-base">open_in_new</span>
            {{ t('openTensorboard') }}
          </a>
        </div>
      </div>
    </section>

    <section class="rounded-xl border border-border-light bg-surface-light p-4 md:p-5 dark:border-border-dark dark:bg-surface-dark">
      <h3 class="mb-4 flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('dataMount') }}
      </h3>
      <div v-if="notebook.volumeMounts && notebook.volumeMounts.length > 0" class="overflow-x-auto">
        <table class="w-full text-left">
          <thead>
            <tr class="border-b border-border-light text-xs font-bold uppercase tracking-wider text-slate-500 dark:border-border-dark">
              <th class="px-4 py-2">{{ t('storage') }}</th>
              <th class="px-4 py-2">{{ t('mountPath') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border-light dark:divide-border-dark">
            <tr v-for="(mount, index) in notebook.volumeMounts" :key="index" class="text-sm">
              <td class="px-4 py-2">{{ mount.pvcName || '-' }}</td>
              <td class="px-4 py-2 font-mono">{{ mount.mountsPath }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-else class="py-5 text-center text-sm text-slate-400">{{ t('noData') }}</div>
    </section>
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
