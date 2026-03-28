<template>
  <div class="rounded-xl border border-border-light bg-surface-light p-6 shadow-sm dark:border-border-dark dark:bg-surface-dark">
    <h3 class="mb-4 text-sm font-bold uppercase tracking-widest text-slate-500 dark:text-slate-300">{{ t('quickEntry') }}</h3>
    <div class="grid grid-cols-2 gap-3">
      <button
        v-for="entry in items"
        :key="entry.key"
        class="flex flex-col items-center justify-center rounded-xl border border-border-light p-4 transition-all hover:border-primary dark:border-border-dark"
        @click="$emit('select', entry.key)"
      >
        <span class="material-icons mb-2 text-[22px] text-primary">{{ entry.icon }}</span>
        <span class="text-xs font-bold">{{ t(entry.labelKey) }}</span>
      </button>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  items: {
    type: Array,
    default: () => []
  }
})

defineEmits(['select'])

const t = inject('t', (key) => key)
</script>
