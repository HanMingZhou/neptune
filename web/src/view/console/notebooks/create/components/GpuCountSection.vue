<template>
  <div
    v-if="selectedProduct && selectedProduct.gpuCount > 0"
    class="console-create-card console-create-card--section w-full"
  >
    <h3 class="console-create-card__title mb-4">
      <span class="console-create-card__title-marker"></span>
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
