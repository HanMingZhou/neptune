<template>
  <div
    v-if="selectedProduct && selectedProduct.gpuCount > 0"
    class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6"
  >
    <h3 class="text-base font-bold mb-4 flex items-center gap-2">
      <span class="w-1 h-4 bg-primary rounded"></span>
      {{ t('gpuCount') }}
    </h3>
    <div class="flex gap-3">
      <button
        v-for="count in Math.min(selectedProduct.gpuCount, 8)"
        :key="count"
        :class="[
          'w-11 h-11 rounded-lg text-sm font-bold border transition-all',
          modelValue === count
            ? 'bg-primary text-white border-primary'
            : 'bg-white dark:bg-zinc-800 border-border-light dark:border-border-dark hover:border-primary hover:text-primary'
        ]"
        @click="$emit('update:modelValue', count)"
      >
        {{ count }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type { ConsoleProduct, Translator } from '@/types/consoleResource'

withDefaults(
  defineProps<{
    modelValue: number
    selectedProduct?: ConsoleProduct | null
  }>(),
  {
    selectedProduct: null
  }
)

defineEmits<{
  'update:modelValue': [value: number]
}>()

const t = inject<Translator>('t', (key: string) => key)
</script>
