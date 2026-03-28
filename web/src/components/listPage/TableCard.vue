<template>
  <div class="overflow-hidden rounded-[22px] border border-border-light bg-white/95 shadow-[0_28px_70px_-52px_rgba(15,23,42,0.58)] backdrop-blur-sm dark:border-border-dark dark:bg-surface-dark/95">
    <div v-if="$slots.toolbar" class="border-b border-border-light/80 p-5 dark:border-border-dark">
      <slot name="toolbar" />
    </div>

    <slot />

    <div v-if="$slots.footer" class="border-t border-border-light/80 bg-slate-50/75 px-6 py-4 dark:border-border-dark dark:bg-zinc-900/35">
      <slot name="footer" />
    </div>
  </div>
</template>
