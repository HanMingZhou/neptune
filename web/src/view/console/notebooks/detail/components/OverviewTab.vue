<template>
  <div class="space-y-3 md:space-y-4">
    <section class="console-detail-card rounded-xl p-4 md:p-5 space-y-5">
      <section class="space-y-3">
        <h3
          class="flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100"
        >
          <span class="w-1 h-4 bg-primary rounded"></span>
          {{ t('basicInfo') }}
        </h3>
        <div class="detail-info-grid detail-info-grid--flat detail-info-grid--inline">
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('name') }}</div>
            <div
              :title="notebook.displayName || '-'"
              class="detail-info-value"
            >
              {{ notebook.displayName || '-' }}
            </div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('instanceId') }}</div>
            <div
              :title="notebook.instanceName || '-'"
              class="detail-info-value detail-info-value--important detail-info-value--mono"
            >
              {{ notebook.instanceName || '-' }}
            </div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('status') }}</div>
            <div
              :title="getStatusLabel(notebook.status)"
              class="detail-info-value detail-info-value--important"
            >
              {{ getStatusLabel(notebook.status) }}
            </div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('createdAt') }}</div>
            <div
              :title="formatTime(notebook.createdAt)"
              class="detail-info-value"
            >
              {{ formatTime(notebook.createdAt) }}
            </div>
          </div>
          <div class="detail-info-item detail-info-item--wide detail-info-item--stack">
            <div class="detail-info-label">{{ t('image') }}</div>
            <div
              class="detail-code-surface detail-code-surface--soft detail-code-surface--single"
            >
              <code
                :title="notebook.imageName || '-'"
                class="detail-info-value--mono"
              >
                {{ notebook.imageName || '-' }}
              </code>
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
            <div class="detail-info-label">{{ t('cpu') }}</div>
            <div
              :title="`${notebook.cpu || 0} ${t('cpu')}`"
              class="detail-info-value"
            >
              {{ notebook.cpu }} {{ t('cpu') }}
            </div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('memory') }}</div>
            <div
              :title="`${notebook.memory || 0} GB`"
              class="detail-info-value"
            >
              {{ notebook.memory }} GB
            </div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('gpu') }}</div>
            <div
              :title="
                notebook.gpuCount
                  ? `${notebook.gpuCount} × ${notebook.gpuModel || 'GPU'}`
                  : hasVGpuSpec(notebook)
                    ? `vGPU / ${formatVGpuSpec(notebook, { detailed: true, t })}`
                    : t('CPU ONLY')
              "
              class="detail-info-value detail-info-value--important"
            >
              <span v-if="notebook.gpuCount">
                {{ notebook.gpuCount }} × {{ notebook.gpuModel }}
              </span>
              <span v-else-if="hasVGpuSpec(notebook)">
                {{ notebook.gpuModel || t('gpu') }} vGPU ·
                {{ formatVGpuSpec(notebook, { detailed: true, t }) }}
              </span>
              <span v-else>{{ t('CPU ONLY') }}</span>
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
          {{ t('order') }}
        </h3>
        <div
          class="detail-info-grid detail-info-grid--flat detail-info-grid--inline"
        >
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('payType') }}</div>
            <div
              :title="getPayTypeLabel(notebook.payType)"
              class="detail-info-value"
            >
              {{ getPayTypeLabel(notebook.payType) }}
            </div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('unitPrice') }}</div>
            <div
              :title="`￥${formatPrice(notebook.price)} / ${getUnitPriceLabel(notebook.payType)}`"
              class="detail-info-value detail-info-value--important"
            >
              ￥{{ formatPrice(notebook.price) }} /
              {{ getUnitPriceLabel(notebook.payType) }}
            </div>
          </div>
        </div>
      </section>

      <section
        v-if="notebook.enableTensorboard"
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
              :title="notebook.tensorboardLogPath || 'logs'"
              class="detail-info-value detail-info-value--mono"
            >
              {{ notebook.tensorboardLogPath || 'logs' }}
            </div>
          </div>
          <div v-if="notebook.tensorboardUrl" class="detail-info-item">
            <div class="detail-info-label">{{ t('accessLink') }}</div>
            <a
              :href="notebook.tensorboardUrl"
              target="_blank"
              :title="notebook.tensorboardUrl"
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
        <div v-if="volumeMounts.length > 0" class="detail-table-shell">
          <table class="console-table">
            <thead>
              <tr>
                <th>{{ t('storage') }}</th>
                <th>{{ t('pvc') }}</th>
                <th>{{ t('mountPath') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(mount, index) in volumeMounts" :key="index">
                <td>{{ mount.name || '-' }}</td>
                <td class="detail-info-value--mono">{{ mount.pvcName || '-' }}</td>
                <td class="detail-info-value--mono">
                  {{ mount.mountsPath || mount.mountPath || '-' }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else class="py-5 text-center text-sm text-slate-400">
          {{ t('noData') }}
        </div>
      </section>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import type {
  ConsoleNotebookDetail,
  ConsoleNotebookVolumeMount,
  Translator
} from '@/types/consoleResource'
import { formatVGpuSpec, hasVGpuSpec } from '@/utils/vgpu'

type FormatTime = (value?: string | number | null) => string
type StatusLabelGetter = (status?: string) => string
type PayTypeLabelGetter = (type?: number | string) => string

const t = inject<Translator>('t', (key: string) => key)

const props = defineProps<{
  formatTime: FormatTime
  getPayTypeLabel: PayTypeLabelGetter
  getStatusLabel: StatusLabelGetter
  getUnitPriceLabel: PayTypeLabelGetter
  notebook: Partial<ConsoleNotebookDetail>
}>()

const volumeMounts = computed<ConsoleNotebookVolumeMount[]>(() =>
  Array.isArray(props.notebook.volumeMounts) ? props.notebook.volumeMounts : []
)

const formatPrice = (value?: number | string): string => {
  const price = Number(value)
  return Number.isFinite(price) ? price.toFixed(4) : '0.0000'
}
</script>
