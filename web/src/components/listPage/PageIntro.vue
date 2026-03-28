<template>
  <div class="space-y-6">
    <div class="flex flex-wrap gap-2 text-[13px] font-medium text-slate-400 dark:text-slate-500">
      <template v-for="(item, index) in breadcrumbs" :key="`${item}-${index}`">
        <span :class="index === breadcrumbs.length - 1 ? 'font-semibold text-slate-900 dark:text-slate-100' : ''">{{ item }}</span>
        <span v-if="index < breadcrumbs.length - 1" class="material-icons text-[16px]">chevron_right</span>
      </template>
    </div>

    <div class="flex flex-wrap gap-3 justify-between items-end">
      <div>
        <h1 class="text-[2.35rem] leading-none font-[780] tracking-[-0.055em] text-slate-950 dark:text-slate-50">{{ title }}</h1>
        <p class="mt-2 max-w-3xl text-[0.95rem] leading-7 text-slate-500 dark:text-slate-400">{{ description }}</p>
      </div>
      <div class="flex flex-wrap gap-3">
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
