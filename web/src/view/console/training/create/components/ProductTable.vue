<template>
  <div
    class="console-create-card console-create-card--table w-full"
  >
    <div class="console-create-card__header">
      <h3 class="console-create-card__title">
        <span class="console-create-card__title-marker"></span>
        {{ t('selectConfig') }}
      </h3>
      <p class="mt-1 text-xs text-slate-400 dark:text-slate-500">
        {{ t('gpuMemoryPerCardHint') }}
      </p>
    </div>
    <div class="console-create-card__table-scroll">
      <table class="w-full min-w-[1360px] table-fixed text-left">
        <thead>
          <tr
            class="border-b border-border-light bg-slate-50 text-xs font-semibold tracking-wide text-slate-500 dark:border-border-dark dark:bg-zinc-800/50"
          >
            <th class="w-12 px-6 py-3.5"></th>
            <th class="w-56 px-6 py-3.5">{{ t('cluster') }}</th>
            <th class="w-[320px] px-6 py-3.5">{{ t('gpuConfig') }}</th>
            <th class="w-36 px-6 py-3.5">{{ t('cpuMemory') }}</th>
            <th class="w-28 px-6 py-3.5">{{ t('area') }}</th>
            <th class="w-28 px-6 py-3.5">{{ t('remainingInventory') }}</th>
            <th class="w-44 px-6 py-3.5">
              {{ t('driverVersion') }}/{{ t('cudaVersion') }}
            </th>
            <th class="w-36 px-6 py-3.5">{{ t('price') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
          <tr
            v-for="product in products"
            :key="product.id"
            :class="[
              'transition-colors',
              isSelectableProduct(product)
                ? 'cursor-pointer'
                : 'cursor-not-allowed opacity-50',
              isSelectableProduct(product)
                ? selectedProductId === product.id
                  ? 'bg-primary/5'
                  : 'hover:bg-slate-50 dark:hover:bg-zinc-800/40'
                : ''
            ]"
            @click="isSelectableProduct(product) && $emit('select-product', product)"
          >
            <td class="px-6 py-4 align-top">
              <div
                :class="[
                  'w-4 h-4 rounded-full border-2 flex items-center justify-center transition-all',
                  selectedProductId === product.id && isSelectableProduct(product)
                    ? 'border-primary'
                    : 'border-slate-300 dark:border-zinc-600'
                ]"
              >
                <div
                  v-if="selectedProductId === product.id && isSelectableProduct(product)"
                  class="w-2 h-2 rounded-full bg-primary"
                ></div>
              </div>
            </td>
            <td class="px-6 py-4 align-top">
              <div class="text-sm font-semibold text-slate-700 dark:text-slate-200">
                {{ product.clusterName || '-' }}
              </div>
              <div class="mt-1 truncate text-xs text-slate-400">
                {{ product.name || product.nodeType }}
              </div>
            </td>
            <td class="px-6 py-4 align-top">
              <div v-if="hasGpuSpec(product)" class="space-y-2">
                <div class="text-sm font-semibold text-primary">
                  {{ product.gpuModel }}
                </div>
                <div class="flex flex-wrap gap-2">
                  <span
                    v-for="entry in getGpuFieldEntries(product)"
                    :key="entry.key"
                    class="inline-flex items-center gap-1 rounded-md border border-slate-200 bg-slate-50 px-2 py-1 text-[11px] leading-none text-slate-600"
                  >
                    <span class="text-slate-400">{{ entry.label }}</span>
                    <span class="font-semibold text-slate-700">{{
                      entry.value
                    }}</span>
                  </span>
                </div>
              </div>
              <div v-else-if="hasVGpuSpec(product)" class="space-y-2">
                <div class="flex items-center gap-2">
                  <div class="text-sm font-semibold text-primary">
                    {{ product.gpuModel || t('gpu') }}
                  </div>
                  <span
                    class="inline-flex items-center rounded-full bg-cyan-50 px-2 py-0.5 text-[11px] font-semibold text-cyan-700"
                  >
                    vGPU
                  </span>
                </div>
                <div class="flex flex-wrap gap-2">
                  <span
                    v-for="entry in getVGpuFieldEntries(product)"
                    :key="entry.key"
                    class="inline-flex items-center gap-1 rounded-md border border-slate-200 bg-slate-50 px-2 py-1 text-[11px] leading-none text-slate-600"
                  >
                    <span class="text-slate-400">{{ entry.label }}</span>
                    <span class="font-semibold text-slate-700">{{
                      entry.value
                    }}</span>
                  </span>
                </div>
              </div>
              <div v-else class="text-sm text-slate-400">CPU ONLY</div>
            </td>
            <td class="px-6 py-4 align-top text-sm">
              {{ product.cpu }} / {{ product.memory }}GB
            </td>
            <td class="px-6 py-4 align-top text-sm">{{ product.area }}</td>
            <td class="px-6 py-4 align-top">
              <span
                :class="[
                  'text-sm font-bold',
                  (product.available || 0) <= 0
                    ? 'text-red-500'
                    : 'text-emerald-500'
                ]"
              >
                {{ product.available || 0 }}
              </span>
            </td>
            <td class="px-6 py-4 align-top text-sm text-slate-600 dark:text-slate-400">
              <div>Driver {{ product.driverVersion || '-' }}</div>
              <div class="mt-1">CUDA {{ product.cudaVersion || '-' }}</div>
            </td>
            <td class="px-6 py-4 align-top">
              <div class="text-lg font-bold text-red-500">
                ¥{{ formatPrice(product) }}
              </div>
              <div class="text-xs text-slate-400">/{{ priceUnitText }}</div>
            </td>
          </tr>
          <tr v-if="products.length === 0">
            <td colspan="8" class="px-6 py-12 text-center text-slate-400">
              {{ t('noData') }}
            </td>
          </tr>
        </tbody>
      </table>
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

withDefaults(
  defineProps<{
    formatPrice: (product: ConsoleProduct | null | undefined) => string
    priceUnitText: string
    products?: ConsoleProduct[]
    selectedProductId?: ResourceId | null
  }>(),
  {
    products: () => [],
    selectedProductId: null
  }
)

defineEmits<{
  'select-product': [product: ConsoleProduct]
}>()

const t = inject<Translator>('t', (key: string) => key)

const getGpuFieldEntries = (product: ConsoleProduct) =>
  buildGpuFieldEntries(product, t)

const getVGpuFieldEntries = (product: ConsoleProduct) =>
  buildVGpuFieldEntries(product, t)

const isSelectableProduct = (product: ConsoleProduct): boolean =>
  (product.available || 0) > 0
</script>
