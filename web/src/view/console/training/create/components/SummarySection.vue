<template>
  <div
    v-if="selectedProduct"
    class="console-create-card console-create-card--section w-full"
  >
    <h3 class="console-create-card__title mb-4">
      <span class="console-create-card__title-marker"></span>
      {{ t('resourceSpecConfirm') }}
    </h3>
    <div
      class="grid grid-cols-2 md:grid-cols-4 gap-6 p-4 bg-slate-50 dark:bg-zinc-800/50 rounded-lg"
    >
      <div class="flex flex-col items-center text-center">
        <div class="text-xs text-slate-400 mb-1">
          {{ hasVGpuSpec(selectedProduct) ? 'vGPU' : t('totalGPU') }}
        </div>
        <div class="text-2xl font-bold text-primary">
          {{
            hasVGpuSpec(selectedProduct) && getVGpuNumber(selectedProduct) <= 0
              ? '-'
              : (totalResources?.gpu || 0)
          }}
        </div>
      </div>
      <div>
        <div class="text-xs text-slate-400 mb-1">{{ t('totalNodes') }}</div>
        <div class="text-2xl font-bold text-primary">
          {{ totalResources?.nodes || 0 }}
        </div>
      </div>
      <div>
        <div class="text-xs text-slate-400 mb-1">{{ t('nodeConfig') }}</div>
        <div class="text-sm font-bold">
          {{ selectedProduct.cpu }} / {{ selectedProduct.memory }}GB
        </div>
      </div>
      <div>
        <div class="text-xs text-slate-400 mb-1">{{ t('totalPrice') }}</div>
        <div class="text-lg font-bold text-red-500">
          ¥{{ totalPrice }}/{{ priceUnitText }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type { ConsoleProduct, Translator } from '@/types/consoleResource'
import { getVGpuNumber, hasVGpuSpec } from '@/utils/vgpu'

interface TrainingTotalResources {
  gpu: number
  nodes: number
}

defineProps<{
  priceUnitText: string
  selectedProduct: ConsoleProduct | null
  totalPrice: string
  totalResources: TrainingTotalResources | null
}>()

const t = inject<Translator>('t', (key: string) => key)
</script>
