<template>
  <div class="space-y-6">
    <div class="flex flex-wrap gap-2 text-sm text-slate-500">
      <template v-for="(item, index) in breadcrumbs" :key="`${item}-${index}`">
        <span :class="index === breadcrumbs.length - 1 ? 'text-slate-900 dark:text-white font-medium' : ''">{{ item }}</span>
        <span v-if="index < breadcrumbs.length - 1" class="material-icons text-[16px]">chevron_right</span>
      </template>
    </div>

    <div class="flex justify-between items-end">
      <div>
        <h1 class="text-3xl font-black tracking-tight text-slate-900 dark:text-white">{{ title }}</h1>
        <p class="text-slate-500 text-sm mt-1">{{ description }}</p>
      </div>
      <div class="flex gap-3">
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
