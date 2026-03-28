<template>
  <div class="grid grid-cols-1 gap-6 md:grid-cols-2 xl:grid-cols-4">
    <div
      v-for="card in items"
      :key="card.key"
      class="min-h-[136px] rounded-xl border border-border-light bg-surface-light p-5 transition-all hover:-translate-y-1 hover:shadow-lg dark:border-border-dark dark:bg-surface-dark"
      :class="card.hoverClass"
    >
      <div class="mb-3 flex items-start justify-between">
        <span class="text-sm text-slate-600 dark:text-slate-300">{{ t(card.titleKey) }}</span>
        <div class="flex size-10 items-center justify-center rounded-xl" :class="card.iconContainerClass">
          <span class="material-icons text-[18px]" :class="card.iconClass">{{ card.icon }}</span>
        </div>
      </div>
      <div class="mb-1 text-3xl font-black leading-none text-slate-900 dark:text-white">{{ card.value }}</div>
      <div class="text-xs text-slate-500 dark:text-slate-400">{{ card.sub }}</div>
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

const t = inject('t', (key) => key)
</script>
