<template>
  <div class="bg-white dark:bg-surface-dark border border-gray-200 dark:border-border-dark rounded-xl overflow-hidden shadow-sm">
    <div v-if="$slots.toolbar" class="p-4 border-b border-gray-100 dark:border-border-dark">
      <slot name="toolbar" />
    </div>

    <slot />

    <div v-if="$slots.footer" class="px-6 py-3 border-t border-gray-200 dark:border-border-dark bg-gray-50 dark:bg-zinc-900/30">
      <slot name="footer" />
    </div>
  </div>
</template>
