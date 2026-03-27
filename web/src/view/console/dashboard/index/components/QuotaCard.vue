<template>
  <div class="rounded-xl border border-dashed border-primary/20 bg-gradient-to-br from-primary/20 to-blue-500/5 p-6 shadow-xl">
    <h3 class="mb-2 text-lg font-black text-primary">{{ t('needQuota') }}</h3>
    <p class="mb-4 text-xs leading-relaxed text-slate-500">{{ t('quotaDesc') }}</p>
    <button class="w-full rounded-lg bg-primary py-2.5 text-xs font-black text-white shadow-lg shadow-primary/20 transition-all hover:bg-primary-hover">
      {{ t('upgradeNow') }}
    </button>
  </div>
</template>

<script setup>
import { inject } from 'vue'

const t = inject('t', (key) => key)
</script>
