<template>
  <div
    :class="[
      'console-workspace-panel overflow-hidden border border-border-light bg-surface-light rounded-xl shadow-sm dark:border-border-dark dark:bg-surface-dark',
      {
        'console-workspace-panel--paged': hasPagedViewport
      }
    ]"
  >
    <div
      v-if="$slots.toolbar"
      class="console-workspace-panel__toolbar border-b border-border-light p-4 dark:border-border-dark"
    >
      <slot name="toolbar" />
    </div>

    <div class="console-workspace-panel__body">
      <slot />
    </div>

    <div v-if="$slots.footer" class="console-list-footer">
      <slot name="footer" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    pageSize?: number
    scrollable?: boolean
  }>(),
  {
    pageSize: 0,
    scrollable: false
  }
)

const hasPagedViewport = computed(() => {
  if (props.scrollable) {
    return true
  }

  const pageSize = Number(props.pageSize)

  return Number.isFinite(pageSize) && pageSize > 0
})
</script>
