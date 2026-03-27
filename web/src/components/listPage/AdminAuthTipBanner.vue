<template>
  <div class="bg-amber-500/10 border border-amber-500/20 rounded-lg px-4 py-3 flex items-center gap-3">
    <span class="material-icons text-amber-500">info</span>
    <span class="text-sm text-amber-700 dark:text-amber-400">{{ t('authTip') }}</span>
  </div>
</template>

<script setup>
import { inject } from 'vue'

const t = inject('t', (key) => key)
</script>
