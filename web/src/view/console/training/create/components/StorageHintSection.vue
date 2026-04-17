<template>
  <div
    class="console-create-card console-create-card--section w-full"
  >
    <h3 class="console-create-card__title mb-4">
      <span class="console-create-card__title-marker"></span>
      {{ t('storageConfig') }}
    </h3>
    <p class="text-sm text-slate-500">{{ t('trainingStorageHint') }}</p>
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'

import type { Translator } from '@/types/consoleResource'
const t = inject<Translator>('t', (key: string) => key)
</script>
