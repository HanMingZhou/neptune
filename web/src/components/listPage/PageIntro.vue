<template>
  <div class="space-y-3">
    <div class="flex flex-wrap gap-2 text-xs font-medium text-slate-400 dark:text-slate-500">
      <template v-for="(item, index) in breadcrumbs" :key="`${item}-${index}`">
        <span :class="index === breadcrumbs.length - 1 ? 'font-semibold text-slate-900 dark:text-slate-100' : ''">{{ item }}</span>
        <span v-if="index < breadcrumbs.length - 1" class="material-icons text-[16px]">chevron_right</span>
      </template>
    </div>

    <div class="flex flex-wrap items-end justify-between gap-3">
      <div>
        <h1 class="text-[1.75rem] font-[720] leading-tight tracking-[-0.025em] text-slate-950 dark:text-slate-50">{{ title }}</h1>
        <p class="mt-1.5 max-w-3xl text-sm leading-6 text-slate-500 dark:text-slate-400">{{ description }}</p>
      </div>
      <div class="page-header-actions">
        <slot name="actions" />
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  breadcrumbs: {
    type: Array,
    default: () => []
  },
  description: {
    type: String,
    default: ''
  },
  title: {
    type: String,
    default: ''
  }
})
</script>
