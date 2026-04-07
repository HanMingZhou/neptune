<template>
  <div
    class="bg-amber-500/10 border border-amber-500/20 rounded-lg px-4 py-3 flex items-center gap-3"
  >
    <span class="material-icons text-amber-500">info</span>
    <span class="text-sm text-amber-700 dark:text-amber-400">{{
      t('authTip')
    }}</span>
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type { Translator } from '@/types/consoleResource'

const t = inject<Translator>('t', (key: string) => key)
</script>
