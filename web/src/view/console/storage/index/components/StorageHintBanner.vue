<template>
  <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg px-4 py-3 text-sm text-blue-700 dark:text-blue-300 flex items-center gap-2">
    <span class="material-icons text-[20px]">lightbulb</span>
    {{ t('storageMountDesc') }}:
    <code class="bg-blue-100 dark:bg-blue-800 px-1 rounded">/home/notebook/neptune</code>
  </div>
</template>

<script setup>
import { inject } from 'vue'

const t = inject('t', (key) => key)
</script>
