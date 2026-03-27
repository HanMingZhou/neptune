<template>
  <section class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl overflow-hidden shadow-sm">
    <div class="px-6 py-4 border-b border-border-light dark:border-border-dark bg-slate-50/50 dark:bg-zinc-800/50">
      <h2 class="text-xs font-black uppercase tracking-widest text-slate-400">{{ title }}</h2>
    </div>
    <slot />
  </section>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true
  }
})
</script>
