<template>
  <div class="rounded-2xl border border-slate-200 bg-white p-8 dark:border-border-dark dark:bg-surface-dark">
    <h3 class="mb-6 text-sm font-bold uppercase tracking-widest text-slate-400">{{ t('rankingMonthly') }}</h3>
    <div class="space-y-6">
      <div v-for="(item, index) in items" :key="item.name" class="space-y-2">
        <div class="flex justify-between text-xs font-bold">
          <span class="text-slate-700 dark:text-slate-300">{{ item.name }}</span>
          <span class="text-slate-900 dark:text-white">¥{{ item.amount.toLocaleString(undefined, { minimumFractionDigits: 2 }) }}</span>
        </div>
        <div class="h-1.5 w-full overflow-hidden rounded-full bg-slate-100 dark:bg-zinc-800">
          <div :class="getRankColor(index)" class="h-full transition-all duration-500" :style="{ width: `${item.percent}%` }"></div>
        </div>
      </div>
    </div>
    <button
      class="mt-8 w-full list-toolbar-button list-toolbar-button--secondary"
      @click="$emit('view-usage')"
    >
      {{ t('viewFullUsageDetails') }}
    </button>
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

defineEmits(['view-usage'])

const t = inject('t', (key) => key)

const rankColors = ['bg-primary', 'bg-amber-500', 'bg-purple-500', 'bg-slate-400']

const getRankColor = (index) => rankColors[index] || 'bg-slate-400'
</script>
