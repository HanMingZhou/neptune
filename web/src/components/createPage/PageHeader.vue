<template>
  <div class="bg-surface-light dark:bg-surface-dark border-b border-border-light dark:border-border-dark px-6 py-4 flex justify-between items-center">
    <h2 class="text-xl font-bold">{{ t(titleKey) }}</h2>
    <button
      class="px-4 py-2 border border-border-light dark:border-border-dark rounded-lg text-sm hover:bg-slate-50 dark:hover:bg-zinc-800 transition-colors"
      @click="$emit('back')"
    >
      {{ t(backLabelKey) }}
    </button>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  backLabelKey: {
    type: String,
    default: 'back'
  },
  titleKey: {
    type: String,
    required: true
  }
})

defineEmits(['back'])

const t = inject('t', (key) => key)
</script>
