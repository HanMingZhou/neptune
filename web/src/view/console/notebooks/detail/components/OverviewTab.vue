<template>
  <div class="space-y-3 md:space-y-4">
    <section class="console-detail-card rounded-xl p-4 md:p-5 space-y-5">
      <section class="space-y-3">
        <h3 class="flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
          <span class="w-1 h-4 bg-primary rounded"></span>
          {{ t('basicInfo') }}
        </h3>
        <div class="detail-info-grid detail-info-grid--flat">
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('name') }}</div>
            <div class="detail-info-value">{{ notebook.displayName || '-' }}</div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('instanceId') }}</div>
            <div class="detail-info-value detail-info-value--mono">{{ notebook.instanceName || '-' }}</div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('status') }}</div>
            <div class="detail-info-value">{{ getStatusLabel(notebook.status) }}</div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('createdAt') }}</div>
            <div class="detail-info-value">{{ formatTime(notebook.createdAt) }}</div>
          </div>
          <div class="detail-info-item detail-info-item--wide">
            <div class="detail-info-label">{{ t('image') }}</div>
            <div class="detail-inline-chip detail-info-value--mono break-all">{{ notebook.imageName || '-' }}</div>
          </div>
        </div>
      </section>

      <section class="space-y-3 border-t border-border-light pt-5 dark:border-border-dark">
        <h3 class="flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
          <span class="w-1 h-4 bg-primary rounded"></span>
          {{ t('resourceConfig') }}
        </h3>
        <div class="detail-info-grid detail-info-grid--flat detail-info-grid--triple-align">
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('cpu') }}</div>
            <div class="detail-info-value">{{ notebook.cpu }} {{ t('cpu') }}</div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('memory') }}</div>
            <div class="detail-info-value">{{ notebook.memory }} GB</div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('gpu') }}</div>
            <div class="detail-info-value">
              <span v-if="notebook.gpuCount" class="text-primary font-bold">{{ notebook.gpuCount }} × {{ notebook.gpuModel }}</span>
              <span v-else class="text-slate-400">{{ t('CPU ONLY') }}</span>
            </div>
          </div>
        </div>
      </section>

      <section class="space-y-3 border-t border-border-light pt-5 dark:border-border-dark">
        <h3 class="flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
          <span class="w-1 h-4 bg-primary rounded"></span>
          {{ t('order') }}
        </h3>
        <div class="detail-info-grid detail-info-grid--flat detail-info-grid--triple-align">
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('payType') }}</div>
            <div class="detail-info-value">{{ getPayTypeLabel(notebook.payType) }}</div>
          </div>
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('unitPrice') }}</div>
            <div class="detail-info-value text-red-500">￥{{ formatPrice(notebook.price) }} / {{ getUnitPriceLabel(notebook.payType) }}</div>
          </div>
        </div>
      </section>

      <section
        v-if="notebook.enableTensorboard"
        class="space-y-3 border-t border-border-light pt-5 dark:border-border-dark"
      >
        <h3 class="flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
          <span class="w-1 h-4 bg-primary rounded"></span>
          TensorBoard
        </h3>
        <div class="detail-info-grid detail-info-grid--flat">
          <div class="detail-info-item">
            <div class="detail-info-label">{{ t('logPath') }}</div>
            <div class="detail-info-value detail-info-value--mono">{{ notebook.tensorboardLogPath || 'logs' }}</div>
          </div>
          <div v-if="notebook.tensorboardUrl" class="detail-info-item">
            <div class="detail-info-label">{{ t('accessLink') }}</div>
            <a :href="notebook.tensorboardUrl" target="_blank" class="inline-flex items-center gap-1 text-sm text-primary hover:underline">
              <span class="material-icons text-base">open_in_new</span>
              {{ t('openTensorboard') }}
            </a>
          </div>
        </div>
      </section>
      <section class="space-y-3 border-t border-border-light pt-5 dark:border-border-dark">
        <h3 class="flex items-center gap-2 text-[15px] font-semibold text-slate-900 dark:text-slate-100">
          <span class="w-1 h-4 bg-primary rounded"></span>
          {{ t('dataMount') }}
        </h3>
        <div v-if="volumeMounts.length > 0" class="detail-table-shell">
          <table class="console-table">
            <thead>
              <tr>
                <th>{{ t('storage') }}</th>
                <th>{{ t('mountPath') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(mount, index) in volumeMounts" :key="index">
                <td>{{ mount.pvcName || '-' }}</td>
                <td class="detail-info-value--mono">{{ mount.mountsPath || mount.mountPath || '-' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else class="py-5 text-center text-sm text-slate-400">{{ t('noData') }}</div>
      </section>
    </section>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'

const t = inject('t', (key) => key)

const props = defineProps({
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

const volumeMounts = computed(() => Array.isArray(props.notebook.volumeMounts) ? props.notebook.volumeMounts : [])

const formatPrice = (value) => {
  const price = Number(value)
  return Number.isFinite(price) ? price.toFixed(4) : '0.0000'
}
</script>
