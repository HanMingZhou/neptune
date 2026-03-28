<template>
  <div class="space-y-4 md:space-y-5">
    <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-4 md:p-5">
      <h3 class="mb-4 flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('basicInfo') }}
      </h3>
      <div class="grid grid-cols-1 gap-x-5 gap-y-4 sm:grid-cols-2 lg:grid-cols-4">
        <div>
          <div class="text-xs text-slate-500 mb-1">{{ t('name') }}</div>
          <div class="text-sm font-medium">{{ service.displayName || '-' }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-500 mb-1">{{ t('instanceName') }}</div>
          <div class="text-sm font-medium font-mono">{{ service.instanceName }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-500 mb-1">{{ t('inference.framework') }}</div>
          <div class="text-sm font-medium">{{ service.framework || t('inference.customCommandMode') }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-500 mb-1">{{ t('deployMode') }}</div>
          <div class="text-sm font-medium">{{ service.deployType === 'STANDALONE' ? t('standalone') : t('distributed') }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-500 mb-1">{{ t('status') }}</div>
          <div class="text-sm font-medium">{{ getStatusLabel(service.status) }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-500 mb-1">{{ t('inference.authType') }}</div>
          <div class="text-sm font-medium">{{ service.authType === 1 ? 'JWT Token' : 'API Key' }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-500 mb-1">{{ t('createdAt') }}</div>
          <div class="text-sm font-medium">{{ formatTime(service.createdAt) }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-500 mb-1">{{ t('startedAt') }}</div>
          <div class="text-sm font-medium">{{ formatTime(service.startedAt) || '-' }}</div>
        </div>
      </div>
    </div>

    <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-4 md:p-5">
      <h3 class="mb-4 flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('resourceConfig') }}
      </h3>
      <div class="grid grid-cols-1 gap-x-5 gap-y-4 sm:grid-cols-2 lg:grid-cols-4">
        <div>
          <div class="text-xs text-slate-500 mb-1">{{ t('gpu') }}</div>
          <div class="text-xl font-semibold text-primary">{{ service.gpu || 0 }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-500 mb-1">{{ t('gpuModel') }}</div>
          <div class="text-sm font-medium">{{ service.gpuModel || 'CPU Only' }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-500 mb-1">{{ t('cpu') }} / {{ t('memory') }}</div>
          <div class="text-sm font-medium">{{ service.cpu }} Core / {{ service.memory }} GB</div>
        </div>
        <div class="col-span-2 md:col-span-1">
          <div class="text-xs text-slate-500 mb-1">{{ t('image') }}</div>
          <div class="text-xs font-mono bg-slate-100 dark:bg-zinc-800 px-2 py-1 rounded inline-flex break-all max-w-full leading-relaxed">
            {{ service.imageName }}
          </div>
        </div>
      </div>
    </div>

    <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-4 md:p-5">
      <h3 class="mb-4 flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('inference.serviceConfig') }}
      </h3>
      <div class="grid grid-cols-1 gap-x-5 gap-y-4 sm:grid-cols-2 lg:grid-cols-4">
        <div>
          <div class="text-xs text-slate-500 mb-1">{{ t('inference.modelPath') }}</div>
          <div class="text-sm font-medium font-mono">{{ service.modelMountPath || '/model' }}{{ service.modelPath }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-500 mb-1">{{ t('inference.servicePort') }}</div>
          <div class="text-sm font-medium">{{ service.servicePort }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-500 mb-1">{{ t('inference.maxTokens') }}</div>
          <div class="text-sm font-medium">{{ service.maxTokens }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-500 mb-1">{{ t('inference.maxConcurrency') }}</div>
          <div class="text-sm font-medium">{{ service.maxConcurrency || t('inference.noLimit') }}</div>
        </div>
      </div>
      <div v-if="service.deployType === 'DISTRIBUTED'" class="mt-4 pt-4 border-t border-border-light dark:border-border-dark grid grid-cols-1 gap-x-5 gap-y-4 sm:grid-cols-2">
        <div>
          <div class="text-xs text-slate-500 mb-1">{{ t('inference.workerCount') }}</div>
          <div class="text-sm font-medium">{{ service.workerCount }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-500 mb-1">{{ t('inference.autoRestart') }}</div>
          <div class="text-sm font-medium">{{ service.autoRestart ? `${t('yes')} (${service.restartCount}/${service.maxRestarts})` : t('no') }}</div>
        </div>
      </div>
      <div v-if="service.command?.length || service.args?.length || service.extraArgs?.length" class="mt-4 pt-4 border-t border-border-light dark:border-border-dark">
        <div class="text-xs text-slate-400 mb-2">{{ t('inference.startupCommand') }}</div>
        <div class="bg-zinc-900 rounded-lg p-3">
          <code class="text-sm text-emerald-400 font-mono whitespace-pre-wrap break-all">{{ formatCommand(service) }}</code>
        </div>
      </div>
    </div>

    <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-4 md:p-5">
      <h3 class="mb-4 flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('inference.apiAccess') }}
      </h3>
      <div v-if="service.status !== 'RUNNING'" class="text-sm text-amber-500 bg-amber-50 dark:bg-amber-900/20 px-3 py-2 rounded-lg flex items-center gap-2">
        <span class="material-icons text-base">info</span>
        {{ t('inference.apiNotReady') }}
      </div>
      <template v-else>
        <div class="bg-zinc-900 rounded-lg p-3 mb-4">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs text-slate-400">{{ t('inference.apiEndpoint') }}</span>
            <button class="text-primary text-xs hover:underline flex items-center gap-1" @click="$emit('copy', apiEndpoint)">
              <span class="material-icons text-sm">content_copy</span>
              {{ t('copy') }}
            </button>
          </div>
          <code class="text-sm text-emerald-400 font-mono break-all">{{ apiEndpoint }}</code>
        </div>
        <div class="bg-zinc-900 rounded-lg p-3 mb-4">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs text-slate-400">{{ t('inference.curlExample') }}</span>
            <button class="text-primary text-xs hover:underline flex items-center gap-1" @click="$emit('copy', curlExample)">
              <span class="material-icons text-sm">content_copy</span>
              {{ t('copy') }}
            </button>
          </div>
          <code class="text-xs text-slate-300 font-mono whitespace-pre-wrap break-all">{{ curlExample }}</code>
        </div>
        <div class="flex gap-3">
          <button
            v-if="service.authType === 2"
            class="px-4 py-2 bg-primary hover:bg-primary-hover text-white rounded-lg text-sm font-bold flex items-center gap-2"
            @click="$emit('manage-api-key')"
          >
            <span class="material-icons text-lg">vpn_key</span>
            {{ t('inference.manageKey') }}
          </button>
        </div>
      </template>
    </div>

    <div v-if="service.mounts?.length" class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-4 md:p-5">
      <h3 class="mb-4 flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('inference.dataMount') }}
      </h3>
      <table class="w-full text-left text-sm">
        <thead>
          <tr class="text-xs text-slate-400 font-bold border-b border-border-light dark:border-border-dark">
            <th class="px-4 py-2.5">PVC</th>
            <th class="px-4 py-2.5">{{ t('inference.mountPath') }}</th>
            <th class="px-4 py-2.5">{{ t('inference.subPath') }}</th>
            <th class="px-4 py-2.5">{{ t('readOnly') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
          <tr v-for="(mount, index) in service.mounts" :key="index">
            <td class="px-4 py-2.5 font-mono">{{ mount.pvcName }}</td>
            <td class="px-4 py-2.5 font-mono">{{ mount.mountPath }}</td>
            <td class="px-4 py-2.5 font-mono">{{ mount.subPath || '-' }}</td>
            <td class="px-4 py-2.5">{{ mount.readOnly ? t('yes') : t('no') }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="service.envs?.length" class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-4 md:p-5">
      <h3 class="mb-4 flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('envVars') }}
      </h3>
      <table class="w-full text-left text-sm">
        <thead>
          <tr class="text-xs text-slate-400 font-bold border-b border-border-light dark:border-border-dark">
            <th class="px-4 py-2.5">{{ t('variableName') }}</th>
            <th class="px-4 py-2.5">{{ t('variableValue') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
          <tr v-for="(env, index) in service.envs" :key="index">
            <td class="px-4 py-2.5 font-mono">{{ env.name }}</td>
            <td class="px-4 py-2.5 font-mono">{{ env.value }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="service.errorMsg" class="bg-surface-light dark:bg-surface-dark border border-red-200 dark:border-red-900 rounded-xl p-4 md:p-5">
      <h3 class="mb-4 flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100 text-red-500">
        <span class="material-icons">error</span>
        {{ t('errorInfo') }}
      </h3>
      <pre class="text-sm text-red-600 dark:text-red-400 font-mono whitespace-pre-wrap break-all bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-3">{{ service.errorMsg }}</pre>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  apiEndpoint: {
    type: String,
    default: ''
  },
  curlExample: {
    type: String,
    default: ''
  },
  formatCommand: {
    type: Function,
    required: true
  },
  formatTime: {
    type: Function,
    required: true
  },
  getStatusLabel: {
    type: Function,
    required: true
  },
  service: {
    type: Object,
    default: () => ({})
  }
})

defineEmits(['copy', 'manage-api-key'])

const t = inject('t', (key) => key)
</script>
