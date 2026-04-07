<template>
  <div class="console-page-container px-6 pt-6">
    <div
      class="console-create-card w-full px-6 py-4 flex justify-between items-center"
    >
      <h2 class="text-xl font-bold">{{ t(titleKey) }}</h2>
      <button
        class="px-4 py-2 border border-border-light bg-surface-light text-sm font-medium transition-colors hover:bg-slate-50 dark:border-border-dark dark:bg-surface-dark dark:hover:bg-zinc-800"
        @click="$emit('back')"
      >
        {{ t(backLabelKey) }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type { Translator } from '@/types/consoleResource'

withDefaults(
  defineProps<{
    backLabelKey?: string
    titleKey: string
  }>(),
  {
    backLabelKey: 'back'
  }
)

defineEmits<{
  back: []
}>()

const t = inject<Translator>('t', (key: string) => key)
</script>
