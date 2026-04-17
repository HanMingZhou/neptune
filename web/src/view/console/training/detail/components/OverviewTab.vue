<template>
  <div class="space-y-3 md:space-y-4">
    <div class="console-detail-card rounded-xl p-4 md:p-5 space-y-5">
      <section class="space-y-3">
        <h3
          class="flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100"
        >
          <span class="w-1 h-4 bg-primary rounded"></span>
          {{ t('basicInfo') }}
        </h3>
        <div class="detail-info-grid detail-info-grid--flat detail-info-grid--inline">
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('displayName') }}</div>
            <div :title="job.displayName || '-'" class="detail-info-value">
              {{ job.displayName || '-' }}
            </div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('taskName') }}</div>
            <div
              :title="job.jobName || '-'"
              class="detail-info-value detail-info-value--important detail-info-value--mono"
            >
              {{ job.jobName || '-' }}
            </div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('framework') }}</div>
            <div
              :title="getFrameworkLabel(job.frameworkType)"
              class="detail-info-value"
            >
              {{ getFrameworkLabel(job.frameworkType) }}
            </div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('status') }}</div>
            <div
              :title="getStatusLabel(job.status)"
              class="detail-info-value detail-info-value--important"
            >
              {{ getStatusLabel(job.status) }}
            </div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('createdAt') }}</div>
            <div
              :title="formatTime(job.createdAt) || '-'"
              class="detail-info-value"
            >
              {{ formatTime(job.createdAt) }}
            </div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('startedAt') }}</div>
            <div
              :title="formatTime(job.startedAt) || '-'"
              class="detail-info-value"
            >
              {{ formatTime(job.startedAt) || '-' }}
            </div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('finishedAt') }}</div>
            <div
              :title="formatTime(job.finishedAt) || '-'"
              class="detail-info-value"
            >
              {{ formatTime(job.finishedAt) || '-' }}
            </div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('duration') }}</div>
            <div :title="job.duration || '-'" class="detail-info-value">
              {{ job.duration || '-' }}
            </div>
          </div>
        </div>
      </section>

      <section
        class="space-y-3 border-t border-border-light pt-5 dark:border-border-dark"
      >
        <h3
          class="flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100"
        >
          <span class="w-1 h-4 bg-primary rounded"></span>
          {{ t('resourceConfig') }}
        </h3>
        <div
          class="detail-info-grid detail-info-grid--flat detail-info-grid--inline"
        >
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('totalGpu') }}</div>
            <div
              :title="String(job.totalGpuCount || 0)"
              class="detail-info-value detail-info-value--important"
            >
              {{ job.totalGpuCount || 0 }}
            </div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('workerCount') }}</div>
            <div
              :title="String(job.workerCount || 1)"
              class="detail-info-value detail-info-value--important"
            >
              {{ job.workerCount || 1 }}
            </div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('gpuType') }}</div>
            <div
              :title="job.gpuModel || job.gpuType || '-'"
              class="detail-info-value"
            >
              {{ job.gpuModel || job.gpuType || '-' }}
            </div>
          </div>
          <div class="detail-info-item detail-info-item--wide detail-info-item--stack">
            <div class="detail-info-label">{{ t('image') }}</div>
            <div
              class="detail-code-surface detail-code-surface--soft detail-code-surface--single"
            >
              <code
                :title="job.imageName || job.image || '-'"
                class="detail-info-value--mono"
              >
                {{ job.imageName || job.image || '-' }}
              </code>
            </div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('nodeConfig') }}</div>
            <div
              :title="`${job.cpu || 0}${t('cpu')} / ${job.memory || 0}GB`"
              class="detail-info-value"
            >
              {{ job.cpu }}{{ t('cpu') }} / {{ job.memory }}GB
            </div>
          </div>
          <div v-if="job.workerGpu > 0" class="detail-info-item">
            <div class="detail-info-label">{{ t('gpuPerWorker') }}</div>
            <div
              :title="String(job.workerGpu)"
              class="detail-info-value"
            >
              {{ job.workerGpu }}
            </div>
          </div>
        </div>
      </section>

      <section
        class="space-y-3 border-t border-border-light pt-5 dark:border-border-dark"
      >
        <h3
          class="flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100"
        >
          <span class="w-1 h-4 bg-primary rounded"></span>
          {{ t('orderAndCluster') }}
        </h3>
        <div
          class="detail-info-grid detail-info-grid--flat detail-info-grid--inline"
        >
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('clusterName') }}</div>
            <div :title="job.clusterName || '-'" class="detail-info-value">
              {{ job.clusterName || '-' }}
            </div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('area') }}</div>
            <div :title="job.area || '-'" class="detail-info-value">
              {{ job.area || '-' }}
            </div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('payType') }}</div>
            <div
              :title="getPayTypeLabel(job.payType)"
              class="detail-info-value"
            >
              {{ getPayTypeLabel(job.payType) }}
            </div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('unitPrice') }}</div>
            <div
              :title="`￥${job.price?.toFixed(4) || '0.0000'} / ${t('unitHour')}`"
              class="detail-info-value detail-info-value--important"
            >
              ￥{{ job.price?.toFixed(4) || '0.0000' }} / {{ t('unitHour') }}
            </div>
          </div>
        </div>
      </section>

      <section
        class="space-y-3 border-t border-border-light pt-5 dark:border-border-dark"
      >
        <h3
          class="flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100"
        >
          <span class="w-1 h-4 bg-primary rounded"></span>
          {{ t('startupCommand') }}
        </h3>
        <div class="detail-code-surface detail-code-surface--soft">
          <pre class="text-sm whitespace-pre-wrap break-all">{{
            job.startupCommand || '-'
          }}</pre>
        </div>
      </section>

      <section
        v-if="job.enableTensorboard"
        class="space-y-3 border-t border-border-light pt-5 dark:border-border-dark"
      >
        <h3
          class="flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100"
        >
          <span class="w-1 h-4 bg-primary rounded"></span>
          TensorBoard
        </h3>
        <div class="detail-info-grid detail-info-grid--flat detail-info-grid--inline">
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('logPath') }}</div>
            <div
              :title="job.tensorboardLogPath || '-'"
              class="detail-info-value detail-info-value--mono"
            >
              {{ job.tensorboardLogPath || '-' }}
            </div>
          </div>
          <div v-if="job.tensorboardUrl" class="detail-info-item">
            <div class="detail-info-label">{{ t('accessLink') }}</div>
            <a
              :href="job.tensorboardUrl"
              target="_blank"
              :title="job.tensorboardUrl"
              class="detail-info-value detail-info-link"
            >
              <span class="material-icons text-base">open_in_new</span>
              {{ t('openTensorboard') }}
            </a>
          </div>
        </div>
      </section>

      <section
        class="space-y-3 border-t border-border-light pt-5 dark:border-border-dark"
      >
        <h3
          class="flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100"
        >
          <span class="w-1 h-4 bg-primary rounded"></span>
          {{ t('dataMount') }}
        </h3>
        <div
          v-if="job.mounts && job.mounts.length > 0"
          class="detail-table-shell"
        >
          <table class="console-table">
            <thead>
              <tr>
                <th>{{ t('storage') }}</th>
                <th>{{ t('pvc') }}</th>
                <th>{{ t('mountPath') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(mount, index) in job.mounts" :key="index">
                <td>{{ mount.name || '-' }}</td>
                <td class="detail-info-value--mono">{{ mount.pvcName || '-' }}</td>
                <td class="detail-info-value--mono">{{ mount.mountPath }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else class="text-sm text-slate-400 text-center py-6">
          {{ t('noData') }}
        </div>
      </section>

      <section
        class="space-y-3 border-t border-border-light pt-5 dark:border-border-dark"
      >
        <h3
          class="flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100"
        >
          <span class="w-1 h-4 bg-primary rounded"></span>
          {{ t('environmentVariables') }}
        </h3>
        <div v-if="job.envs && job.envs.length > 0" class="detail-table-shell">
          <table class="console-table">
            <thead>
              <tr>
                <th>{{ t('variableName') }}</th>
                <th>{{ t('variableValue') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(env, index) in job.envs" :key="index">
                <td class="detail-info-value--mono text-primary">
                  {{ env.name }}
                </td>
                <td class="detail-info-value--mono">{{ env.value }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else class="text-sm text-slate-400 text-center py-6">
          {{ t('noData') }}
        </div>
      </section>
    </div>

    <div
      v-if="job.errorMsg"
      class="bg-surface-light dark:bg-surface-dark border border-red-200 dark:border-red-900 rounded-xl p-4 md:p-5"
    >
      <h3
        class="mb-4 flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100 text-red-500"
      >
        <span class="material-icons">error</span>
        {{ t('errorInfo') }}
      </h3>
      <div
        class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-3"
      >
        <pre
          class="text-sm text-red-600 dark:text-red-400 font-mono whitespace-pre-wrap break-all"
          >{{ job.errorMsg }}</pre
        >
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type { ConsoleTrainingDetail, Translator } from '@/types/consoleResource'

type FormatTime = (value?: string | number | null) => string | null
type StatusLabelGetter = (status?: string) => string
type PayTypeLabelGetter = (type?: number | string) => string
type TrainingFrameworkGetter = (type?: string) => string

defineProps<{
  formatTime: FormatTime
  getFrameworkLabel: TrainingFrameworkGetter
  getPayTypeLabel: PayTypeLabelGetter
  getStatusLabel: StatusLabelGetter
  job: Partial<ConsoleTrainingDetail>
}>()

const t = inject<Translator>('t', (key: string) => key)
</script>
