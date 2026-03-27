<template>
  <div class="space-y-6">
    <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
      <h3 class="text-base font-bold mb-4 flex items-center gap-2">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('basicInfo') }}
      </h3>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('displayName') }}</div>
          <div class="text-sm font-medium">{{ job.displayName || '-' }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('taskName') }}</div>
          <div class="text-sm font-medium font-mono">{{ job.jobName }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('framework') }}</div>
          <div class="text-sm font-medium">{{ getFrameworkLabel(job.frameworkType) }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('status') }}</div>
          <div class="text-sm font-medium">{{ getStatusLabel(job.status) }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('createdAt') }}</div>
          <div class="text-sm font-medium">{{ formatTime(job.createdAt) }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('startedAt') }}</div>
          <div class="text-sm font-medium">{{ formatTime(job.startedAt) || '-' }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('finishedAt') }}</div>
          <div class="text-sm font-medium">{{ formatTime(job.finishedAt) || '-' }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('duration') }}</div>
          <div class="text-sm font-medium">{{ job.duration || '-' }}</div>
        </div>
      </div>
    </div>

    <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
      <h3 class="text-base font-bold mb-4 flex items-center gap-2">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('resourceConfig') }}
      </h3>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('totalGpu') }}</div>
          <div class="text-2xl font-bold text-primary">{{ job.totalGpuCount || 0 }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('workerCount') }}</div>
          <div class="text-2xl font-bold text-primary">{{ job.workerCount || 1 }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('gpuType') }}</div>
          <div class="text-sm font-medium">{{ job.gpuModel || job.gpuType || '-' }}</div>
        </div>
        <div class="col-span-2">
          <div class="text-xs text-slate-400 mb-1">{{ t('image') }}</div>
          <div class="text-xs font-mono bg-slate-100 dark:bg-zinc-800 px-2 py-1 rounded inline-flex break-all max-w-full leading-relaxed">
            {{ job.imageName || job.image }}
          </div>
        </div>
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('nodeConfig') }}</div>
          <div class="text-sm font-medium">{{ job.cpu }}{{ t('cpu') }} / {{ job.memory }}GB</div>
        </div>
        <div v-if="job.workerGpu > 0">
          <div class="text-xs text-slate-400 mb-1">{{ t('gpuPerWorker') }}</div>
          <div class="text-sm font-medium">{{ job.workerGpu }}</div>
        </div>
      </div>
    </div>

    <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
      <h3 class="text-base font-bold mb-4 flex items-center gap-2">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('orderAndCluster') }}
      </h3>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('clusterName') }}</div>
          <div class="text-sm font-medium">{{ job.clusterName }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('area') }}</div>
          <div class="text-sm font-medium">{{ job.area || '-' }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('payType') }}</div>
          <div class="text-sm font-medium">{{ getPayTypeLabel(job.payType) }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('unitPrice') }}</div>
          <div class="text-lg font-bold text-red-500">￥{{ job.price?.toFixed(4) || '0.0000' }} / {{ t('unitHour') }}</div>
        </div>
      </div>
    </div>

    <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
      <h3 class="text-base font-bold mb-4 flex items-center gap-2">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('startupCommand') }}
      </h3>
      <div class="bg-zinc-900 rounded-lg p-4">
        <pre class="text-sm text-slate-300 font-mono whitespace-pre-wrap break-all">{{ job.startupCommand || '-' }}</pre>
      </div>
    </div>

    <div v-if="job.enableTensorboard" class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
      <h3 class="text-base font-bold mb-4 flex items-center gap-2">
        <span class="w-1 h-4 bg-primary rounded"></span>
        TensorBoard
      </h3>
      <div class="grid grid-cols-2 md:grid-cols-3 gap-6">
        <div>
          <div class="text-xs text-slate-400 mb-1">{{ t('logPath') }}</div>
          <div class="text-sm font-medium font-mono">{{ job.tensorboardLogPath || '-' }}</div>
        </div>
        <div v-if="job.tensorboardUrl">
          <div class="text-xs text-slate-400 mb-1">{{ t('accessLink') }}</div>
          <a :href="job.tensorboardUrl" target="_blank" class="text-sm text-primary hover:underline flex items-center gap-1">
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
      <div v-if="job.mounts && job.mounts.length > 0" class="overflow-x-auto">
        <table class="w-full text-left">
          <thead>
            <tr class="text-xs text-slate-400 font-bold uppercase tracking-wider border-b border-border-light dark:border-border-dark">
              <th class="px-4 py-3">{{ t('storage') }}</th>
              <th class="px-4 py-3">{{ t('mountPath') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border-light dark:divide-border-dark">
            <tr v-for="(mount, index) in job.mounts" :key="index" class="text-sm">
              <td class="px-4 py-3">{{ mount.pvcName || '-' }}</td>
              <td class="px-4 py-3 font-mono">{{ mount.mountPath }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-else class="text-sm text-slate-400 text-center py-8">{{ t('noData') }}</div>
    </div>

    <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
      <h3 class="text-base font-bold mb-4 flex items-center gap-2">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('environmentVariables') }}
      </h3>
      <div v-if="job.envs && job.envs.length > 0" class="overflow-x-auto">
        <table class="w-full text-left">
          <thead>
            <tr class="text-xs text-slate-400 font-bold uppercase tracking-wider border-b border-border-light dark:border-border-dark">
              <th class="px-4 py-3 w-1/3">{{ t('variableName') }}</th>
              <th class="px-4 py-3">{{ t('variableValue') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border-light dark:divide-border-dark">
            <tr v-for="(env, index) in job.envs" :key="index" class="text-sm">
              <td class="px-4 py-3 font-mono text-primary font-bold">{{ env.name }}</td>
              <td class="px-4 py-3 font-mono">{{ env.value }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-else class="text-sm text-slate-400 text-center py-8">{{ t('noData') }}</div>
    </div>

    <div v-if="job.errorMsg" class="bg-surface-light dark:bg-surface-dark border border-red-200 dark:border-red-900 rounded-xl p-6">
      <h3 class="text-base font-bold mb-4 flex items-center gap-2 text-red-500">
        <span class="material-icons">error</span>
        {{ t('errorInfo') }}
      </h3>
      <div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
        <pre class="text-sm text-red-600 dark:text-red-400 font-mono whitespace-pre-wrap break-all">{{ job.errorMsg }}</pre>
      </div>
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
  getFrameworkLabel: {
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
  job: {
    type: Object,
    default: () => ({})
  }
})

const t = inject('t', (key) => key)
</script>
