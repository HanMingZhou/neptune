<template>
  <div
    v-if="selectedProduct"
    class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6"
  >
    <h3 class="text-base font-bold mb-4 flex items-center gap-2">
      <span class="w-1 h-4 bg-primary rounded"></span>
      {{ t('resourceSpecConfirm') }}
    </h3>
    <div
      class="grid grid-cols-2 md:grid-cols-4 gap-6 p-4 bg-slate-50 dark:bg-zinc-800/50 rounded-lg"
    >
      <div>
        <div class="text-xs text-slate-400 mb-1">{{ t('totalGPU') }}</div>
        <div class="text-2xl font-bold text-primary">
          {{ totalResources?.gpu || 0 }}
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
