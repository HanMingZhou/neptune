<template>
  <div
    v-if="selectedProduct"
    class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6"
  >
    <h3 class="text-base font-bold mb-4 flex items-center gap-2">
      <span class="w-1 h-4 bg-primary rounded"></span>
      {{ t('resourceSpecConfirm') }}
    </h3>
    <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4 p-4 bg-slate-50 dark:bg-zinc-800/50 rounded-lg">
      <div>
        <div class="text-xs text-slate-400 mb-1">{{ t('hostId') }}</div>
        <div class="text-sm font-bold">{{ selectedProduct.name || selectedProduct.nodeType }}</div>
      </div>
      <div>
        <div class="text-xs text-slate-400 mb-1">GPU</div>
        <div class="text-sm font-bold">
          <template v-if="selectedProduct.gpuCount > 0">{{ selectedProduct.gpuModel }} × {{ gpuCount }}</template>
          <template v-else-if="selectedProduct.vGpuCount > 0">vGPU ({{ selectedProduct.vGpuMemory }}GB)</template>
          <template v-else>CPU ONLY</template>
        </div>
      </div>
      <div>
        <div class="text-xs text-slate-400 mb-1">{{ t('cpuMemory') }}</div>
        <div class="text-sm font-bold">{{ selectedProduct.cpu }} / {{ selectedProduct.memory }}GB</div>
      </div>
      <div>
        <div class="text-xs text-slate-400 mb-1">{{ t('systemDisk') }}</div>
        <div class="text-sm font-bold">10GB ({{ t('free') }})</div>
      </div>
      <div v-if="selectedVolumeId">
        <div class="text-xs text-slate-400 mb-1">{{ t('dataMount') }}</div>
        <div class="text-sm font-bold">{{ selectedVolumeName }}</div>
      </div>
      <div>
        <div class="text-xs text-slate-400 mb-1">{{ t('unitPrice') }}</div>
        <div class="text-lg font-bold text-red-500">¥{{ totalPrice }}/{{ priceUnitText }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  gpuCount: {
    type: Number,
    required: true
  },
  priceUnitText: {
    type: String,
    required: true
  },
  selectedProduct: {
    type: Object,
    default: null
  },
  selectedVolumeId: {
    type: [Number, String],
    default: null
  },
  selectedVolumeName: {
    type: String,
    default: ''
  },
  totalPrice: {
    type: String,
    required: true
  }
})

const t = inject('t', (key) => key)
</script>
