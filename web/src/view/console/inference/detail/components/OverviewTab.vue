<template>
  <div class="space-y-3 md:space-y-4">
    <div class="console-detail-card rounded-xl p-4 md:p-5 space-y-5">
      <section class="space-y-3">
        <h3 class="flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
          <span class="w-1 h-4 bg-primary rounded"></span>
          {{ t('basicInfo') }}
        </h3>
        <div class="detail-info-grid detail-info-grid--flat">
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('name') }}</div>
            <div class="detail-info-value">{{ service.displayName || '-' }}</div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('instanceName') }}</div>
            <div class="detail-info-value detail-info-value--mono">{{ service.instanceName }}</div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('inference.framework') }}</div>
            <div class="detail-info-value">{{ service.framework || t('inference.customCommandMode') }}</div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('deployMode') }}</div>
            <div class="detail-info-value">{{ service.deployType === 'STANDALONE' ? t('standalone') : t('distributed') }}</div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('status') }}</div>
            <div class="detail-info-value">{{ getStatusLabel(service.status) }}</div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('inference.authType') }}</div>
            <div class="detail-info-value">{{ service.authType === 1 ? 'JWT Token' : 'API Key' }}</div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('createdAt') }}</div>
            <div class="detail-info-value">{{ formatTime(service.createdAt) }}</div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('startedAt') }}</div>
            <div class="detail-info-value">{{ formatTime(service.startedAt) || '-' }}</div>
          </div>
        </div>
      </section>

      <section class="space-y-3 border-t border-border-light pt-5 dark:border-border-dark">
        <h3 class="flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
          <span class="w-1 h-4 bg-primary rounded"></span>
          {{ t('resourceConfig') }}
        </h3>
        <div class="detail-info-grid detail-info-grid--flat">
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('gpu') }}</div>
            <div class="detail-info-value text-primary">{{ service.gpu || 0 }}</div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('gpuModel') }}</div>
            <div class="detail-info-value">{{ service.gpuModel || 'CPU Only' }}</div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('cpu') }} / {{ t('memory') }}</div>
            <div class="detail-info-value">{{ service.cpu }} Core / {{ service.memory }} GB</div>
          </div>
          <div class="detail-info-item detail-info-item--wide">
            <div class="detail-info-label">{{ t('image') }}</div>
            <div class="detail-inline-chip detail-info-value--mono break-all">{{ service.imageName }}</div>
          </div>
        </div>
      </section>

      <section class="space-y-3 border-t border-border-light pt-5 dark:border-border-dark">
        <h3 class="flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
          <span class="w-1 h-4 bg-primary rounded"></span>
          {{ t('inference.serviceConfig') }}
        </h3>
        <div class="detail-info-grid detail-info-grid--flat">
          <div class="detail-info-item detail-info-item--wide">
            <div class="detail-info-label">{{ t('inference.modelPath') }}</div>
            <div class="detail-info-value detail-info-value--mono">{{ service.modelMountPath || '/model' }}{{ service.modelPath }}</div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('inference.servicePort') }}</div>
            <div class="detail-info-value">{{ service.servicePort }}</div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('inference.maxTokens') }}</div>
            <div class="detail-info-value">{{ service.maxTokens }}</div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('inference.maxConcurrency') }}</div>
            <div class="detail-info-value">{{ service.maxConcurrency || t('inference.noLimit') }}</div>
          </div>
        </div>
        <div v-if="service.deployType === 'DISTRIBUTED'" class="detail-info-grid detail-info-grid--flat mt-4 border-t border-border-light pt-4 dark:border-border-dark">
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('inference.workerCount') }}</div>
            <div class="detail-info-value">{{ service.workerCount }}</div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('inference.autoRestart') }}</div>
            <div class="detail-info-value">{{ service.autoRestart ? `${t('yes')} (${service.restartCount}/${service.maxRestarts})` : t('no') }}</div>
          </div>
        </div>
        <div v-if="service.command?.length || service.args?.length || service.extraArgs?.length" class="mt-4 pt-4 border-t border-border-light dark:border-border-dark">
          <div class="text-xs text-slate-400 mb-2">{{ t('inference.startupCommand') }}</div>
          <div class="detail-code-surface">
            <code class="text-sm whitespace-pre-wrap break-all">{{ formatCommand(service) }}</code>
          </div>
        </div>
      </section>
      <section class="space-y-3 border-t border-border-light pt-5 dark:border-border-dark">
        <h3 class="flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
          <span class="w-1 h-4 bg-primary rounded"></span>
          {{ t('inference.apiAccess') }}
        </h3>
        <div v-if="service.status !== 'RUNNING'" class="text-sm text-amber-500 bg-amber-50 dark:bg-amber-900/20 px-3 py-2 rounded-lg flex items-center gap-2">
          <span class="material-icons text-base">info</span>
          {{ t('inference.apiNotReady') }}
        </div>
        <template v-else>
          <div class="detail-code-surface mb-4">
            <div class="flex items-center justify-between mb-2">
              <span class="text-xs text-slate-400">{{ t('inference.apiEndpoint') }}</span>
              <button class="text-primary text-xs hover:underline flex items-center gap-1" @click="$emit('copy', apiEndpoint)">
                <span class="material-icons text-sm">content_copy</span>
                {{ t('copy') }}
              </button>
            </div>
            <code class="text-sm break-all">{{ apiEndpoint }}</code>
          </div>
          <div class="detail-code-surface mb-4">
            <div class="flex items-center justify-between mb-2">
              <span class="text-xs text-slate-400">{{ t('inference.curlExample') }}</span>
              <button class="text-primary text-xs hover:underline flex items-center gap-1" @click="$emit('copy', curlExample)">
                <span class="material-icons text-sm">content_copy</span>
                {{ t('copy') }}
              </button>
            </div>
            <code class="text-xs whitespace-pre-wrap break-all">{{ curlExample }}</code>
          </div>
          <div class="flex gap-3">
            <button
              v-if="service.authType === 2"
              class="list-toolbar-button list-toolbar-button--primary"
              @click="$emit('manage-api-key')"
            >
              <span class="material-icons text-[18px]">vpn_key</span>
              {{ t('inference.manageKey') }}
            </button>
          </div>
        </template>
      </section>

      <section
        v-if="service.mounts?.length || service.envs?.length"
        class="space-y-5 border-t border-border-light pt-5 dark:border-border-dark"
      >
        <section v-if="service.mounts?.length" class="space-y-3">
          <h3 class="flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
            <span class="w-1 h-4 bg-primary rounded"></span>
            {{ t('inference.dataMount') }}
          </h3>
          <div class="detail-table-shell">
            <table class="console-table">
              <thead>
                <tr>
                  <th>PVC</th>
                  <th>{{ t('inference.mountPath') }}</th>
                  <th>{{ t('inference.subPath') }}</th>
                  <th>{{ t('readOnly') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(mount, index) in service.mounts" :key="index">
                  <td class="detail-info-value--mono">{{ mount.pvcName }}</td>
                  <td class="detail-info-value--mono">{{ mount.mountPath }}</td>
                  <td class="detail-info-value--mono">{{ mount.subPath || '-' }}</td>
                  <td>{{ mount.readOnly ? t('yes') : t('no') }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>

        <section
          v-if="service.envs?.length"
          class="space-y-3"
          :class="service.mounts?.length ? 'border-t border-border-light pt-5 dark:border-border-dark' : ''"
        >
          <h3 class="flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
            <span class="w-1 h-4 bg-primary rounded"></span>
            {{ t('envVars') }}
          </h3>
          <div class="detail-table-shell">
            <table class="console-table">
              <thead>
                <tr>
                  <th>{{ t('variableName') }}</th>
                  <th>{{ t('variableValue') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(env, index) in service.envs" :key="index">
                  <td class="detail-info-value--mono">{{ env.name }}</td>
                  <td class="detail-info-value--mono">{{ env.value }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>
      </section>
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

