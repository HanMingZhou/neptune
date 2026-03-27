<template>
  <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
    <h3 class="text-base font-bold mb-4 flex items-center gap-2">
      <span class="w-1 h-4 bg-primary rounded"></span>
      {{ t('storageConfig') }}
    </h3>
    <p class="text-sm text-slate-500">{{ t('trainingStorageHint') }}</p>
  </div>
</template>

<script setup>
import { inject } from 'vue'

const t = inject('t', (key) => key)
</script>
