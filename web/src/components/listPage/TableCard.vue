<template>
  <div class="console-workspace-panel overflow-hidden border border-border-light bg-surface-light rounded-xl shadow-sm dark:border-border-dark dark:bg-surface-dark">
    <div v-if="$slots.toolbar" class="console-workspace-panel__toolbar border-b border-border-light p-4 dark:border-border-dark">
      <slot name="toolbar" />
    </div>

    <slot />

    <div v-if="$slots.footer" class="console-list-footer">
      <slot name="footer" />
    </div>
  </div>
</template>
