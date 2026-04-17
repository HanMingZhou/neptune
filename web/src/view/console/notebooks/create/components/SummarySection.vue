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
      class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4 p-4 bg-slate-50 dark:bg-zinc-800/50 rounded-lg"
    >
      <div>
        <div class="text-xs text-slate-400 mb-1">{{ t('hostId') }}</div>
        <div class="text-sm font-bold">
          {{ selectedProduct.name || selectedProduct.nodeType }}
        </div>
      </div>
      <div class="flex flex-col items-center text-center">
        <div class="text-xs text-slate-400 mb-1">GPU</div>
        <div
          v-if="hasGpuSpec(selectedProduct)"
          class="flex flex-col items-center space-y-2 text-center"
        >
          <div class="text-sm font-bold">
            {{ selectedProduct.gpuModel }}
          </div>
          <div class="flex flex-wrap justify-center gap-2">
            <span
              v-for="entry in getGpuFieldEntries()"
              :key="entry.key"
              class="inline-flex items-center gap-1 rounded-md border border-slate-200 bg-white px-2 py-1 text-[11px] leading-none text-slate-600"
            >
              <span class="text-slate-400">{{ entry.label }}</span>
              <span class="font-semibold text-slate-700">{{ entry.value }}</span>
            </span>
          </div>
        </div>
        <div
          v-else-if="hasVGpuSpec(selectedProduct)"
          class="flex flex-col items-center space-y-2 text-center"
        >
          <div class="flex items-center justify-center gap-2 text-sm font-bold">
            <span>{{ selectedProduct.gpuModel || t('gpu') }}</span>
            <span
              class="inline-flex items-center rounded-full bg-cyan-50 px-2 py-0.5 text-[11px] font-bold text-cyan-700"
            >
              vGPU
            </span>
          </div>
          <div class="flex flex-wrap justify-center gap-2">
            <span
              v-for="entry in getVGpuFieldEntries(selectedProduct)"
              :key="entry.key"
              class="inline-flex items-center gap-1 rounded-md border border-slate-200 bg-white px-2 py-1 text-[11px] leading-none text-slate-600"
            >
              <span class="text-slate-400">{{ entry.label }}</span>
              <span class="font-semibold text-slate-700">{{ entry.value }}</span>
            </span>
          </div>
        </div>
        <div v-else class="text-sm font-bold">CPU ONLY</div>
      </div>
      <div>
        <div class="text-xs text-slate-400 mb-1">{{ t('cpuMemory') }}</div>
        <div class="text-sm font-bold">
          {{ selectedProduct.cpu }} / {{ selectedProduct.memory }}GB
        </div>
      </div>
      <div>
        <div class="text-xs text-slate-400 mb-1">{{ t('systemDisk') }}</div>
        <div class="text-sm font-bold">50GB ({{ t('free') }})</div>
      </div>
      <div v-if="selectedVolumeId">
        <div class="text-xs text-slate-400 mb-1">{{ t('dataMount') }}</div>
        <div class="text-sm font-bold">{{ selectedVolumeName }}</div>
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
import type {
  ConsoleProduct,
  ResourceId,
  Translator
} from '@/types/consoleResource'
import {
  buildGpuFieldEntries,
  buildVGpuFieldEntries,
  hasGpuSpec,
  hasVGpuSpec
} from '@/utils/vgpu'

const t = inject<Translator>('t', (key: string) => key)

const props = defineProps<{
  gpuCount: number
  priceUnitText: string
  selectedProduct: ConsoleProduct | null
  selectedVolumeId: ResourceId | null
  selectedVolumeName: string
  totalPrice: string
}>()

const getGpuFieldEntries = () =>
  buildGpuFieldEntries(
    {
      gpuCount: props.gpuCount,
      gpuMemory: props.selectedProduct?.gpuMemory
    },
    t
  )

const getVGpuFieldEntries = (product: ConsoleProduct) =>
  buildVGpuFieldEntries(product, t)
</script>
