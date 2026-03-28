<template>
  <div class="console-page-container px-6 pt-5 pb-3">
    <div class="rounded-xl border border-border-light bg-surface-light px-5 py-4 dark:border-border-dark dark:bg-surface-dark md:px-6">
      <div class="flex flex-wrap items-center justify-between gap-4">
        <div class="flex min-w-0 items-center gap-3 md:gap-4">
          <button class="inline-flex shrink-0 items-center gap-1 text-sm text-slate-500 transition-colors hover:text-primary" @click="$emit('back')">
            <span class="material-icons text-lg">arrow_back</span>
            {{ t('back') }}
          </button>
          <h2 class="truncate text-lg font-semibold text-slate-900 dark:text-slate-100 md:text-xl">{{ title }}</h2>
          <span :class="['shrink-0 rounded-full px-3 py-1 text-xs font-bold', getStatusClass(status)]">
            {{ getStatusLabel(status) }}
          </span>
        </div>
        <div class="flex flex-wrap items-center justify-end gap-2.5">
          <slot name="actions" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  getStatusClass: {
    type: Function,
    required: true
  },
  getStatusLabel: {
    type: Function,
    required: true
  },
  status: {
    type: String,
    default: ''
  },
  title: {
    type: String,
    default: ''
  }
})

defineEmits(['back'])

const t = inject('t', (key) => key)
</script>
