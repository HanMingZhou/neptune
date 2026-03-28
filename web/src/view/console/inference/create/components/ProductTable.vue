<template>
  <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl overflow-hidden">
    <div class="p-6 border-b border-border-light dark:border-border-dark flex justify-between items-center">
      <h3 class="text-base font-bold flex items-center gap-2">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('selectConfig') }}
      </h3>
    </div>
    <div class="overflow-x-auto">
      <table class="w-full text-left">
        <thead>
          <tr class="bg-slate-50 dark:bg-zinc-800/50 border-b border-border-light dark:border-border-dark text-slate-500 text-xs font-bold uppercase tracking-wider">
            <th class="px-6 py-4 w-12"></th>
            <th class="px-6 py-4">{{ t('cluster') }}</th>
            <th class="px-6 py-4">{{ t('gpuConfig') }}</th>
            <th class="px-6 py-4">{{ t('cpuMemory') }}</th>
            <th class="px-6 py-4">{{ t('remainingInventory') }}</th>
            <th class="px-6 py-4">{{ t('disk') }}</th>
            <th class="px-6 py-4">{{ t('driverVersion') }}</th>
            <th class="px-6 py-4">{{ t('price') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
          <tr
            v-for="product in products"
            :key="product.id"
            :class="[
              'cursor-pointer transition-colors',
              productId === product.id ? 'bg-primary/5' : 'hover:bg-slate-50 dark:hover:bg-zinc-800/40'
            ]"
            @click="$emit('update:productId', product.id)"
          >
            <td class="px-6 py-4">
              <div
                :class="[
                  'w-4 h-4 rounded-full border-2 flex items-center justify-center transition-all',
                  productId === product.id ? 'border-primary' : 'border-slate-300 dark:border-zinc-600'
                ]"
              >
                <div v-if="productId === product.id" class="w-2 h-2 rounded-full bg-primary"></div>
              </div>
            </td>
            <td class="px-6 py-4">
              <div class="font-bold text-sm">{{ product.clusterName || '-' }}</div>
              <div class="text-xs text-slate-400 mt-1">{{ product.name || product.nodeType }}</div>
            </td>
            <td class="px-6 py-4">
              <div v-if="product.gpuCount > 0">
                <div class="font-bold text-primary text-sm">{{ product.gpuModel }}</div>
                <div class="text-xs text-slate-400 mt-1">{{ product.gpuMemory }}GB × {{ product.gpuCount }}</div>
              </div>
              <div v-else class="text-slate-400 text-sm">CPU ONLY</div>
            </td>
            <td class="px-6 py-4 text-sm">{{ product.cpu }} / {{ product.memory }}GB</td>
            <td class="px-6 py-4">
              <span :class="['font-bold text-sm', (product.available || 0) <= 0 ? 'text-red-500' : 'text-emerald-500']">
                {{ product.available || 0 }}
              </span>
            </td>
            <td class="px-6 py-4 text-sm text-slate-600 dark:text-slate-400">
              <div>{{ t('systemDisk') }}: {{ product.systemDisk || 10 }}GB</div>
            </td>
            <td class="px-6 py-4 text-sm text-slate-600 dark:text-slate-400">
              <div>CUDA {{ product.cudaVersion || '-' }}</div>
            </td>
            <td class="px-6 py-4">
              <div class="text-lg font-bold text-red-500">¥{{ formatPrice(product) }}</div>
              <div class="text-xs text-slate-400">/{{ priceUnitText }}</div>
            </td>
          </tr>
          <tr v-if="products.length === 0">
            <td colspan="8" class="px-6 py-12 text-center text-slate-400">{{ t('noData') }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  formatPrice: {
    type: Function,
    required: true
  },
  priceUnitText: {
    type: String,
    required: true
  },
  productId: {
    type: [Number, String],
    default: ''
  },
  products: {
    type: Array,
    default: () => []
  }
})

defineEmits(['update:productId'])

const t = inject('t', (key) => key)
</script>
