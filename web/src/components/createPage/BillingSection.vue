<template>
  <div
    class="console-create-card bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6"
  >
    <h3 class="text-base font-bold mb-4 flex items-center gap-2">
      <span class="w-1 h-4 bg-primary rounded"></span>
      {{ t(titleKey) }}
    </h3>
    <div class="flex gap-3 flex-wrap">
      <button
        v-for="item in payTypes"
        :key="item.value"
        :class="[
          'px-5 py-2 rounded-lg text-sm font-medium border transition-all',
          payType === item.value
            ? 'bg-primary text-white border-primary shadow-lg shadow-primary/20'
            : 'bg-white dark:bg-zinc-800 border-border-light dark:border-border-dark hover:border-primary hover:text-primary'
        ]"
        @click="$emit('update:payType', item.value)"
      >
        {{ t(item.labelKey) }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type { Translator } from '@/types/consoleResource'

interface PayTypeOption {
  value: number
  labelKey: string
}

withDefaults(
  defineProps<{
    payType: number
    payTypes?: PayTypeOption[]
    titleKey?: string
  }>(),
  {
    payTypes: () => [],
    titleKey: 'orderMethod'
  }
)

defineEmits<{
  'update:payType': [value: number]
}>()

const t = inject<Translator>('t', (key: string) => key)
</script>
