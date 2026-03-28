<template>
  <div class="rounded-xl border border-border-light bg-surface-light p-6 shadow-sm dark:border-border-dark dark:bg-surface-dark">
    <h3 class="mb-6 text-sm font-bold uppercase tracking-widest text-slate-500 dark:text-slate-300">{{ t('resourceTrend') }}</h3>
    <div v-if="items.length > 0" class="flex h-64 items-end justify-between px-2 pb-2">
      <div v-for="item in items" :key="item.date" class="group mx-0.5 flex flex-1 flex-col items-center gap-2">
        <div
          class="relative w-full max-w-[40px] cursor-pointer rounded-t-md bg-primary/20 transition-all group-hover:bg-primary/40"
          :style="{ height: `${Math.max(item.barHeight, 4)}px` }"
        >
          <div class="absolute -top-8 left-1/2 -translate-x-1/2 whitespace-nowrap rounded border border-border-light bg-white px-2 py-0.5 text-xs font-bold text-primary opacity-0 shadow-sm transition-opacity group-hover:opacity-100 dark:border-border-dark dark:bg-zinc-800">
            {{ item.runningTasks }} {{ t('tasks') }}
          </div>
        </div>
        <span class="text-xs text-slate-500 dark:text-slate-400">{{ item.dateLabel }}</span>
      </div>
    </div>
    <div v-else class="flex h-64 items-center justify-center text-sm text-slate-400">
      <span class="material-icons mr-2 text-4xl opacity-30">bar_chart</span>
      {{ t('noData') }}
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
